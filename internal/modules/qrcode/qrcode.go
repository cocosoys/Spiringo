package qrcode

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/core/di"
	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/modules/qrcode/handler"
	"github.com/spiringo/spiringo/internal/modules/qrcode/repository"
	"github.com/spiringo/spiringo/internal/modules/qrcode/service"
	"github.com/spiringo/spiringo/internal/pkg/orm"
	"github.com/spiringo/spiringo/internal/pkg/storage"
)

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
// Config 二维码模块配置
type Config struct {
	// 中文：DefaultSize 保存当前结构中的配置或数据值。
	// English: DefaultSize stores a configuration or data value for this struct.
	DefaultSize int `yaml:"default_size" mapstructure:"default_size"`
	// 中文：DefaultLevel 保存当前结构中的配置或数据值。
	// English: DefaultLevel stores a configuration or data value for this struct.
	DefaultLevel string `yaml:"default_level" mapstructure:"default_level"`
	// 中文：OSSPrefix 保存当前结构中的配置或数据值。
	// English: OSSPrefix stores a configuration or data value for this struct.
	OSSPrefix string `yaml:"oss_prefix" mapstructure:"oss_prefix"`
	// 中文：BucketName 保存当前结构中的配置或数据值。
	// English: BucketName stores a configuration or data value for this struct.
	BucketName string `yaml:"bucket_name" mapstructure:"bucket_name"`
}

// 中文：QRCodeModule 定义当前包使用的数据结构或接口。
// English: QRCodeModule defines a data structure or interface used by this package.
// QRCodeModule 二维码模块
type QRCodeModule struct {
	// 中文：*module.BaseModule 嵌入复用该类型提供的能力。
	// English: *module.BaseModule embeds reusable behavior from that type.
	*module.BaseModule
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config Config
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *service.QRCodeService
	// 中文：migrateDB 保存当前结构中的配置或数据值。
	// English: migrateDB stores a configuration or data value for this struct.
	migrateDB *orm.DB
}

// 中文：NewQRCodeModule 创建并返回对应组件实例。
// English: NewQRCodeModule creates and returns the corresponding component instance.
// NewQRCodeModule 创建二维码模块
func NewQRCodeModule() *QRCodeModule {
	return &QRCodeModule{
		BaseModule: module.NewBaseModule("qrcode", "tenant"),
	}
}

// 中文：Config 执行当前包中的对应流程。
// English: Config executes the corresponding workflow in this package.
func (m *QRCodeModule) Config() any { return &m.config }

// 中文：Init 执行当前包中的对应流程。
// English: Init executes the corresponding workflow in this package.
func (m *QRCodeModule) Init(app *module.App) error {
	db, err := di.Resolve[*orm.DB](app.DI)
	if err != nil {
		return fmt.Errorf("qrcode module init: %w", err)
	}
	m.migrateDB = db
	tdb := orm.NewTenantDB(db)
	repo := repository.NewQRCodeRepository(tdb, db)

	svcCfg := service.Config{
		DefaultSize:  m.config.DefaultSize,
		DefaultLevel: m.config.DefaultLevel,
		OSSPrefix:    m.config.OSSPrefix,
		BucketName:   m.config.BucketName,
	}
	if svcCfg.DefaultSize == 0 {
		svcCfg.DefaultSize = 256
	}
	if svcCfg.DefaultLevel == "" {
		svcCfg.DefaultLevel = "medium"
	}

	// 可选：从DI获取Storage实例
	var s storage.Storage
	if di.Has[storage.Storage](app.DI) {
		s, _ = di.Resolve[storage.Storage](app.DI)
	}

	m.svc = service.NewQRCodeService(svcCfg, repo, s)
	return nil
}

// 中文：Routes 执行当前包中的对应流程。
// English: Routes executes the corresponding workflow in this package.
func (m *QRCodeModule) Routes(r *gin.RouterGroup) {
	h := handler.NewQRCodeHandler(m.svc)
	r.POST("/generate", h.Generate)
	r.POST("/parse", h.Parse)
	r.GET("/s/:code", h.Redirect)
	r.GET("/stats/:code", h.Stats)
}

// 中文：Start 执行当前包中的对应流程。
// English: Start executes the corresponding workflow in this package.
func (m *QRCodeModule) Start(_ context.Context) error {
	if m.svc == nil {
		return fmt.Errorf("qrcode service is not initialized")
	}
	if m.migrateDB == nil {
		return fmt.Errorf("qrcode migration database is not initialized")
	}
	return nil
}

// 中文：Stop 执行当前包中的对应流程。
// English: Stop executes the corresponding workflow in this package.
func (m *QRCodeModule) Stop(ctx context.Context) error {
	if ctx != nil {
		return ctx.Err()
	}
	return nil
}
