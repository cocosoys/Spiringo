package notification

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/core/di"
	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/modules/notification/handler"
	"github.com/spiringo/spiringo/internal/modules/notification/repository"
	"github.com/spiringo/spiringo/internal/modules/notification/service"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
type Config struct {
	// 中文：Events 保存当前结构中的配置或数据值。
	// English: Events stores a configuration or data value for this struct.
	Events []string `yaml:"events" mapstructure:"events"`
	// 中文：Webhook 保存当前结构中的配置或数据值。
	// English: Webhook stores a configuration or data value for this struct.
	Webhook service.WebhookConfig `yaml:"webhook" mapstructure:"webhook"`
	// 中文：Email 保存当前结构中的配置或数据值。
	// English: Email stores a configuration or data value for this struct.
	Email service.EmailConfig `yaml:"email" mapstructure:"email"`
	// 中文：Inbox 保存当前结构中的配置或数据值。
	// English: Inbox stores a configuration or data value for this struct.
	Inbox service.InboxConfig `yaml:"inbox" mapstructure:"inbox"`
}

// 中文：Module 定义当前包使用的数据结构或接口。
// English: Module defines a data structure or interface used by this package.
type Module struct {
	// 中文：*module.BaseModule 嵌入复用该类型提供的能力。
	// English: *module.BaseModule embeds reusable behavior from that type.
	*module.BaseModule
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config Config
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *service.Service
	// 中文：migrateDB 保存当前结构中的配置或数据值。
	// English: migrateDB stores a configuration or data value for this struct.
	migrateDB *orm.DB
}

// 中文：NewNotificationModule 创建并返回对应组件实例。
// English: NewNotificationModule creates and returns the corresponding component instance.
func NewNotificationModule() *Module {
	return &Module{
		BaseModule: module.NewBaseModule("notification"),
	}
}

// 中文：Config 执行当前包中的对应流程。
// English: Config executes the corresponding workflow in this package.
func (m *Module) Config() any { return &m.config }

// 中文：Init 执行当前包中的对应流程。
// English: Init executes the corresponding workflow in this package.
func (m *Module) Init(app *module.App) error {
	db, err := di.Resolve[*orm.DB](app.DI)
	if err != nil {
		return fmt.Errorf("notification module init: %w", err)
	}
	m.migrateDB = db
	repo := repository.NewNotificationRepository(orm.NewTenantDB(db), db)
	m.svc = service.New(service.Config{
		Events:  m.config.Events,
		Webhook: m.config.Webhook,
		Email:   m.config.Email,
		Inbox:   m.config.Inbox,
	}, repo)
	return nil
}

// 中文：Routes 执行当前包中的对应流程。
// English: Routes executes the corresponding workflow in this package.
func (m *Module) Routes(r *gin.RouterGroup) {
	h := handler.New(m.svc)
	r.POST("/send", h.Send)
	r.GET("/inbox", h.Inbox)
	r.PUT("/inbox/:id/read", h.MarkRead)
}

// 中文：Start 执行当前包中的对应流程。
// English: Start executes the corresponding workflow in this package.
func (m *Module) Start(_ context.Context) error {
	if m.svc == nil {
		return fmt.Errorf("notification service is not initialized")
	}
	if m.migrateDB == nil {
		return fmt.Errorf("notification migration database is not initialized")
	}
	return nil
}

// 中文：Stop 执行当前包中的对应流程。
// English: Stop executes the corresponding workflow in this package.
func (m *Module) Stop(ctx context.Context) error {
	if ctx != nil {
		return ctx.Err()
	}
	return nil
}
