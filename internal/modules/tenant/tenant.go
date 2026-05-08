package tenant

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/core/di"
	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/modules/tenant/handler"
	"github.com/spiringo/spiringo/internal/modules/tenant/repository"
	"github.com/spiringo/spiringo/internal/modules/tenant/service"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
// Config 租户模块配置
type Config struct {
	// 中文：Strategy 保存当前结构中的配置或数据值。
	// English: Strategy stores a configuration or data value for this struct.
	Strategy string `yaml:"strategy" mapstructure:"strategy"` // shared_db | schema | database
}

// 中文：TenantModule 定义当前包使用的数据结构或接口。
// English: TenantModule defines a data structure or interface used by this package.
// TenantModule 多租户模块
type TenantModule struct {
	// 中文：*module.BaseModule 嵌入复用该类型提供的能力。
	// English: *module.BaseModule embeds reusable behavior from that type.
	*module.BaseModule
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config Config
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *service.TenantService
	// 中文：migrateDB 保存当前结构中的配置或数据值。
	// English: migrateDB stores a configuration or data value for this struct.
	migrateDB *orm.DB
}

// 中文：NewTenantModule 创建并返回对应组件实例。
// English: NewTenantModule creates and returns the corresponding component instance.
// NewTenantModule 创建租户模块
func NewTenantModule() *TenantModule {
	return &TenantModule{
		BaseModule: module.NewBaseModule("tenant"),
	}
}

// 中文：Config 执行当前包中的对应流程。
// English: Config executes the corresponding workflow in this package.
func (m *TenantModule) Config() any { return &m.config }

// 中文：Init 执行当前包中的对应流程。
// English: Init executes the corresponding workflow in this package.
func (m *TenantModule) Init(app *module.App) error {
	db, err := di.Resolve[*orm.DB](app.DI)
	if err != nil {
		return fmt.Errorf("tenant module init: %w", err)
	}
	m.migrateDB = db
	repo := repository.NewTenantRepository(db)
	m.svc = service.NewTenantService(repo, app.EventBus)
	return nil
}

// 中文：Routes 执行当前包中的对应流程。
// English: Routes executes the corresponding workflow in this package.
func (m *TenantModule) Routes(r *gin.RouterGroup) {
	h := handler.NewTenantHandler(m.svc)
	r.GET("", h.List)
	r.GET("/:id", h.Get)
	r.POST("", h.Create)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}

// 中文：Start 执行当前包中的对应流程。
// English: Start executes the corresponding workflow in this package.
func (m *TenantModule) Start(_ context.Context) error {
	if m.svc == nil {
		return fmt.Errorf("tenant service is not initialized")
	}
	if m.migrateDB == nil {
		return fmt.Errorf("tenant migration database is not initialized")
	}
	return nil
}

// 中文：Stop 执行当前包中的对应流程。
// English: Stop executes the corresponding workflow in this package.
func (m *TenantModule) Stop(ctx context.Context) error {
	if ctx != nil {
		return ctx.Err()
	}
	return nil
}
