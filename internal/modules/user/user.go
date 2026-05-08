package user

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/core/di"
	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/modules/auth/service"
	"github.com/spiringo/spiringo/internal/modules/user/handler"
	"github.com/spiringo/spiringo/internal/modules/user/repository"
	usersvc "github.com/spiringo/spiringo/internal/modules/user/service"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：UserModule 定义当前包使用的数据结构或接口。
// English: UserModule defines a data structure or interface used by this package.
// UserModule 用户模块
type UserModule struct {
	// 中文：*module.BaseModule 嵌入复用该类型提供的能力。
	// English: *module.BaseModule embeds reusable behavior from that type.
	*module.BaseModule
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *usersvc.UserService
	// 中文：migrateDB 保存当前结构中的配置或数据值。
	// English: migrateDB stores a configuration or data value for this struct.
	migrateDB *orm.DB
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config Config
}

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
type Config struct {
	// 中文：DefaultAdmin 保存当前结构中的配置或数据值。
	// English: DefaultAdmin stores a configuration or data value for this struct.
	DefaultAdmin usersvc.DefaultAdminConfig `yaml:"default_admin" mapstructure:"default_admin"`
}

// 中文：NewUserModule 创建并返回对应组件实例。
// English: NewUserModule creates and returns the corresponding component instance.
// NewUserModule 创建用户模块
func NewUserModule() *UserModule {
	return &UserModule{
		BaseModule: module.NewBaseModule("user", "tenant"),
	}
}

// 中文：Config 执行当前包中的对应流程。
// English: Config executes the corresponding workflow in this package.
func (m *UserModule) Config() any { return &m.config }

// 中文：Init 执行当前包中的对应流程。
// English: Init executes the corresponding workflow in this package.
func (m *UserModule) Init(app *module.App) error {
	db, err := di.Resolve[*orm.DB](app.DI)
	if err != nil {
		return fmt.Errorf("user module init: %w", err)
	}
	m.migrateDB = db
	tdb := orm.NewTenantDB(db)
	repo := repository.NewUserRepository(tdb)
	m.svc = usersvc.NewUserService(repo, app.EventBus)
	adminCfg := m.config.DefaultAdmin
	if !app.Config.IsSet("modules.user.default_admin.enabled") {
		adminCfg.Enabled = true
	}
	m.svc.SetDefaultAdminConfig(adminCfg)

	// 将UserService注册为auth模块的UserServiceInterface
	app.DI.ProvideNamed("auth_user_service", m.svc)

	return nil
}

// 中文：Routes 执行当前包中的对应流程。
// English: Routes executes the corresponding workflow in this package.
func (m *UserModule) Routes(r *gin.RouterGroup) {
	h := handler.NewUserHandler(m.svc)
	r.GET("", h.List)
	r.GET("/:id", h.Get)
	r.POST("", h.Create)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}

// 中文：Start 执行当前包中的对应流程。
// English: Start executes the corresponding workflow in this package.
func (m *UserModule) Start(_ context.Context) error {
	if m.svc == nil {
		return fmt.Errorf("user service is not initialized")
	}
	if m.migrateDB == nil {
		return fmt.Errorf("user migration database is not initialized")
	}
	return nil
}

// 中文：Stop 执行当前包中的对应流程。
// English: Stop executes the corresponding workflow in this package.
func (m *UserModule) Stop(ctx context.Context) error {
	if ctx != nil {
		return ctx.Err()
	}
	return nil
}

// 中文：_ 声明当前包使用的变量。
// English: _ declares variables used by this package.
// 确保 UserService 实现 auth 模块的 UserServiceInterface
var _ service.UserServiceInterface = (*usersvc.UserService)(nil)
