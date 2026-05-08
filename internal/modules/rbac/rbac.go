package rbac

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/core/di"
	"github.com/spiringo/spiringo/internal/core/middleware"
	"github.com/spiringo/spiringo/internal/core/module"
	authsvc "github.com/spiringo/spiringo/internal/modules/auth/service"
	"github.com/spiringo/spiringo/internal/modules/rbac/handler"
	"github.com/spiringo/spiringo/internal/modules/rbac/repository"
	"github.com/spiringo/spiringo/internal/modules/rbac/service"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
type Config struct {
	// 中文：AuthRequired 保存当前结构中的配置或数据值。
	// English: AuthRequired stores a configuration or data value for this struct.
	AuthRequired bool `yaml:"auth_required" mapstructure:"auth_required"`
}

// 中文：RBACModule 定义当前包使用的数据结构或接口。
// English: RBACModule defines a data structure or interface used by this package.
type RBACModule struct {
	// 中文：*module.BaseModule 嵌入复用该类型提供的能力。
	// English: *module.BaseModule embeds reusable behavior from that type.
	*module.BaseModule
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config Config
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *service.RBACService
	// 中文：migrateDB 保存当前结构中的配置或数据值。
	// English: migrateDB stores a configuration or data value for this struct.
	migrateDB *orm.DB
	// 中文：tokenValidator 保存当前结构中的配置或数据值。
	// English: tokenValidator stores a configuration or data value for this struct.
	tokenValidator middleware.TokenValidator
}

// 中文：NewRBACModule 创建并返回对应组件实例。
// English: NewRBACModule creates and returns the corresponding component instance.
func NewRBACModule() *RBACModule {
	return &RBACModule{
		BaseModule: module.NewBaseModule("rbac", "auth", "user", "tenant"),
		config: Config{
			AuthRequired: true,
		},
	}
}

// 中文：Config 执行当前包中的对应流程。
// English: Config executes the corresponding workflow in this package.
func (m *RBACModule) Config() any { return &m.config }

// 中文：Init 执行当前包中的对应流程。
// English: Init executes the corresponding workflow in this package.
func (m *RBACModule) Init(app *module.App) error {
	db, err := di.Resolve[*orm.DB](app.DI)
	if err != nil {
		return fmt.Errorf("rbac module init: %w", err)
	}
	m.migrateDB = db
	tdb := orm.NewTenantDB(db)
	repo := repository.NewRBACRepository(tdb, db)
	m.svc = service.NewRBACService(repo)

	if m.config.AuthRequired {
		validator, err := di.Resolve[*authsvc.AuthService](app.DI)
		if err != nil {
			return fmt.Errorf("rbac auth middleware init: %w", err)
		}
		m.tokenValidator = validator
	}

	return nil
}

// 中文：Routes 执行当前包中的对应流程。
// English: Routes executes the corresponding workflow in this package.
func (m *RBACModule) Routes(r *gin.RouterGroup) {
	h := handler.NewRBACHandler(m.svc)
	if !m.config.AuthRequired {
		r.GET("/roles", h.ListRoles)
		r.POST("/roles", h.CreateRole)
		r.PUT("/roles/:id", h.UpdateRole)
		r.DELETE("/roles/:id", h.DeleteRole)
		r.GET("/permissions", h.ListPermissions)
		r.POST("/roles/:id/permissions", h.AssignPermissions)
		return
	}

	protected := r.Group("", middleware.Auth(m.tokenValidator))
	protected.GET("/roles", middleware.RequirePermission(m.svc, "rbac.roles", "read"), h.ListRoles)
	protected.POST("/roles", middleware.RequirePermission(m.svc, "rbac.roles", "create"), h.CreateRole)
	protected.PUT("/roles/:id", middleware.RequirePermission(m.svc, "rbac.roles", "update"), h.UpdateRole)
	protected.DELETE("/roles/:id", middleware.RequirePermission(m.svc, "rbac.roles", "delete"), h.DeleteRole)
	protected.GET("/permissions", middleware.RequirePermission(m.svc, "rbac.permissions", "read"), h.ListPermissions)
	protected.POST("/roles/:id/permissions", middleware.RequirePermission(m.svc, "rbac.permissions", "assign"), h.AssignPermissions)
}

// 中文：Start 执行当前包中的对应流程。
// English: Start executes the corresponding workflow in this package.
func (m *RBACModule) Start(_ context.Context) error {
	if m.svc == nil {
		return fmt.Errorf("rbac service is not initialized")
	}
	if m.migrateDB == nil {
		return fmt.Errorf("rbac migration database is not initialized")
	}
	if m.config.AuthRequired && m.tokenValidator == nil {
		return fmt.Errorf("rbac token validator is not initialized")
	}
	return nil
}

// 中文：Stop 执行当前包中的对应流程。
// English: Stop executes the corresponding workflow in this package.
func (m *RBACModule) Stop(ctx context.Context) error {
	if ctx != nil {
		return ctx.Err()
	}
	return nil
}
