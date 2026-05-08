package payment

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/core/di"
	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/modules/payment/channel"
	"github.com/spiringo/spiringo/internal/modules/payment/handler"
	"github.com/spiringo/spiringo/internal/modules/payment/repository"
	"github.com/spiringo/spiringo/internal/modules/payment/service"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
// Config 支付模块配置
type Config struct {
	// 中文：Wechat 保存当前结构中的配置或数据值。
	// English: Wechat stores a configuration or data value for this struct.
	Wechat struct {
		Enabled    bool   `yaml:"enabled" mapstructure:"enabled"`
		AppID      string `yaml:"app_id" mapstructure:"app_id"`
		MchID      string `yaml:"mch_id" mapstructure:"mch_id"`
		APIV3Key   string `yaml:"api_v3_key" mapstructure:"api_v3_key"`
		APIKey     string `yaml:"api_key" mapstructure:"api_key"`
		SerialNo   string `yaml:"serial_no" mapstructure:"serial_no"`
		PrivateKey string `yaml:"private_key" mapstructure:"private_key"`
		CertPath   string `yaml:"cert_path" mapstructure:"cert_path"`
	} `yaml:"wechat" mapstructure:"wechat"`
	// 中文：Alipay 保存当前结构中的配置或数据值。
	// English: Alipay stores a configuration or data value for this struct.
	Alipay struct {
		Enabled        bool   `yaml:"enabled" mapstructure:"enabled"`
		AppID          string `yaml:"app_id" mapstructure:"app_id"`
		PrivateKey     string `yaml:"private_key" mapstructure:"private_key"`
		IsProd         bool   `yaml:"is_prod" mapstructure:"is_prod"`
		AppCertPath    string `yaml:"app_cert_path" mapstructure:"app_cert_path"`
		AlipayCertPath string `yaml:"alipay_cert_path" mapstructure:"alipay_cert_path"`
		RootCertPath   string `yaml:"root_cert_path" mapstructure:"root_cert_path"`
	} `yaml:"alipay" mapstructure:"alipay"`
	// 中文：Stripe 保存当前结构中的配置或数据值。
	// English: Stripe stores a configuration or data value for this struct.
	Stripe struct {
		Enabled       bool   `yaml:"enabled" mapstructure:"enabled"`
		SecretKey     string `yaml:"secret_key" mapstructure:"secret_key"`
		WebhookSecret string `yaml:"webhook_secret" mapstructure:"webhook_secret"`
	} `yaml:"stripe" mapstructure:"stripe"`
	// 中文：PayPal 保存当前结构中的配置或数据值。
	// English: PayPal stores a configuration or data value for this struct.
	PayPal struct {
		Enabled      bool   `yaml:"enabled" mapstructure:"enabled"`
		ClientID     string `yaml:"client_id" mapstructure:"client_id"`
		ClientSecret string `yaml:"client_secret" mapstructure:"client_secret"`
		Sandbox      bool   `yaml:"sandbox" mapstructure:"sandbox"`
		WebhookID    string `yaml:"webhook_id" mapstructure:"webhook_id"`
	} `yaml:"paypal" mapstructure:"paypal"`
	// 中文：UnionPay 保存当前结构中的配置或数据值。
	// English: UnionPay stores a configuration or data value for this struct.
	UnionPay struct {
		Enabled bool   `yaml:"enabled" mapstructure:"enabled"`
		MchID   string `yaml:"mch_id" mapstructure:"mch_id"`
		APIKey  string `yaml:"api_key" mapstructure:"api_key"`
		Sandbox bool   `yaml:"sandbox" mapstructure:"sandbox"`
	} `yaml:"unionpay" mapstructure:"unionpay"`
	// 中文：CloudPay 保存当前结构中的配置或数据值。
	// English: CloudPay stores a configuration or data value for this struct.
	CloudPay struct {
		Enabled    bool   `yaml:"enabled" mapstructure:"enabled"`
		MchID      string `yaml:"mch_id" mapstructure:"mch_id"`
		APIKey     string `yaml:"api_key" mapstructure:"api_key"`
		GatewayURL string `yaml:"gateway_url" mapstructure:"gateway_url"`
	} `yaml:"cloudpay" mapstructure:"cloudpay"`
	// 中文：DigitalRMB 保存当前结构中的配置或数据值。
	// English: DigitalRMB stores a configuration or data value for this struct.
	DigitalRMB struct {
		Enabled    bool   `yaml:"enabled" mapstructure:"enabled"`
		AppID      string `yaml:"app_id" mapstructure:"app_id"`
		MerchantID string `yaml:"merchant_id" mapstructure:"merchant_id"`
		APIKey     string `yaml:"api_key" mapstructure:"api_key"`
		GatewayURL string `yaml:"gateway_url" mapstructure:"gateway_url"`
		WalletID   string `yaml:"wallet_id" mapstructure:"wallet_id"`
	} `yaml:"digital_rmb" mapstructure:"digital_rmb"`
	// 中文：DefaultNotifyURL 保存当前结构中的配置或数据值。
	// English: DefaultNotifyURL stores a configuration or data value for this struct.
	DefaultNotifyURL string `yaml:"default_notify_url" mapstructure:"default_notify_url"`
}

// 中文：PaymentModule 定义当前包使用的数据结构或接口。
// English: PaymentModule defines a data structure or interface used by this package.
// PaymentModule 支付模块
type PaymentModule struct {
	// 中文：*module.BaseModule 嵌入复用该类型提供的能力。
	// English: *module.BaseModule embeds reusable behavior from that type.
	*module.BaseModule
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config Config
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *service.PaymentService
	// 中文：migrateDB 保存当前结构中的配置或数据值。
	// English: migrateDB stores a configuration or data value for this struct.
	migrateDB *orm.DB
}

// 中文：NewPaymentModule 创建并返回对应组件实例。
// English: NewPaymentModule creates and returns the corresponding component instance.
// NewPaymentModule 创建支付模块
func NewPaymentModule() *PaymentModule {
	return &PaymentModule{
		BaseModule: module.NewBaseModule("payment", "tenant"),
	}
}

// 中文：Config 执行当前包中的对应流程。
// English: Config executes the corresponding workflow in this package.
func (m *PaymentModule) Config() any { return &m.config }

// 中文：Init 执行当前包中的对应流程。
// English: Init executes the corresponding workflow in this package.
func (m *PaymentModule) Init(app *module.App) error {
	db, err := di.Resolve[*orm.DB](app.DI)
	if err != nil {
		return fmt.Errorf("payment module init: %w", err)
	}
	m.migrateDB = db
	tdb := orm.NewTenantDB(db)
	repo := repository.NewPaymentRepository(tdb, db)

	// 初始化支付通道
	chReg := channel.NewRegistry()
	svcCfg := service.Config{}
	svcCfg.DefaultNotifyURL = m.config.DefaultNotifyURL

	if m.config.Wechat.Enabled {
		wechatAPIKey := firstNonEmpty(m.config.Wechat.APIV3Key, m.config.Wechat.APIKey)
		wc := channel.NewWechatChannel(
			m.config.Wechat.AppID,
			m.config.Wechat.MchID,
			wechatAPIKey,
			m.config.DefaultNotifyURL,
		)
		wc.SetSerialNo(m.config.Wechat.SerialNo)
		wc.SetPrivateKey(m.config.Wechat.PrivateKey)
		chReg.Register(wc)
		svcCfg.Wechat.Enabled = true
		svcCfg.Wechat.AppID = m.config.Wechat.AppID
		svcCfg.Wechat.MchID = m.config.Wechat.MchID
		svcCfg.Wechat.APIKey = wechatAPIKey
	}

	if m.config.Alipay.Enabled {
		ac := channel.NewAlipayChannel(
			m.config.Alipay.AppID,
			m.config.Alipay.PrivateKey,
			m.config.DefaultNotifyURL,
		)
		ac.SetProd(m.config.Alipay.IsProd)
		if hasAnyValue(m.config.Alipay.AppCertPath, m.config.Alipay.AlipayCertPath, m.config.Alipay.RootCertPath) {
			appCert, alipayCert, rootCert, err := loadAlipayCerts(
				m.config.Alipay.AppCertPath,
				m.config.Alipay.AlipayCertPath,
				m.config.Alipay.RootCertPath,
			)
			if err != nil {
				return err
			}
			ac.SetCerts(appCert, alipayCert, rootCert)
		}
		chReg.Register(ac)
		svcCfg.Alipay.Enabled = true
		svcCfg.Alipay.AppID = m.config.Alipay.AppID
		svcCfg.Alipay.PrivateKey = m.config.Alipay.PrivateKey
	}

	if m.config.Stripe.Enabled {
		sc := channel.NewStripeChannel(
			m.config.Stripe.SecretKey,
			m.config.DefaultNotifyURL,
		)
		sc.SetWebhookSecret(m.config.Stripe.WebhookSecret)
		chReg.Register(sc)
		svcCfg.Stripe.Enabled = true
		svcCfg.Stripe.SecretKey = m.config.Stripe.SecretKey
	}

	if m.config.PayPal.Enabled {
		pc := channel.NewPayPalChannel(
			m.config.PayPal.ClientID,
			m.config.PayPal.ClientSecret,
			m.config.PayPal.Sandbox,
			m.config.DefaultNotifyURL,
		)
		pc.SetWebhookID(m.config.PayPal.WebhookID)
		chReg.Register(pc)
		svcCfg.PayPal.Enabled = true
	}

	if m.config.UnionPay.Enabled {
		uc := channel.NewUnionPayChannel(
			m.config.UnionPay.MchID,
			m.config.UnionPay.APIKey,
			m.config.DefaultNotifyURL,
		)
		chReg.Register(uc)
		svcCfg.UnionPay.Enabled = true
	}

	if m.config.CloudPay.Enabled {
		cc := channel.NewCloudPayChannel(
			m.config.CloudPay.MchID,
			m.config.CloudPay.APIKey,
			m.config.CloudPay.GatewayURL,
			m.config.DefaultNotifyURL,
		)
		chReg.Register(cc)
		svcCfg.CloudPay.Enabled = true
		svcCfg.CloudPay.MchID = m.config.CloudPay.MchID
		svcCfg.CloudPay.APIKey = m.config.CloudPay.APIKey
	}

	if m.config.DigitalRMB.Enabled {
		dc := channel.NewDigitalRMBChannel(
			m.config.DigitalRMB.AppID,
			m.config.DigitalRMB.MerchantID,
			m.config.DigitalRMB.APIKey,
			m.config.DigitalRMB.GatewayURL,
			m.config.DigitalRMB.WalletID,
			m.config.DefaultNotifyURL,
		)
		chReg.Register(dc)
		svcCfg.DigitalRMB.Enabled = true
		svcCfg.DigitalRMB.AppID = m.config.DigitalRMB.AppID
		svcCfg.DigitalRMB.MerchantID = m.config.DigitalRMB.MerchantID
		svcCfg.DigitalRMB.APIKey = m.config.DigitalRMB.APIKey
	}

	m.svc = service.NewPaymentService(svcCfg, app.EventBus, repo, chReg)
	return nil
}

// 中文：Routes 执行当前包中的对应流程。
// English: Routes executes the corresponding workflow in this package.
func (m *PaymentModule) Routes(r *gin.RouterGroup) {
	h := handler.NewPaymentHandler(m.svc)
	r.POST("/create", h.Create)
	r.POST("/callback/:channel", h.Callback)
	r.POST("/refund", h.Refund)
	r.GET("/query/:id", h.Query)
	r.POST("/close/:id", h.Close)
}

// 中文：Start 执行当前包中的对应流程。
// English: Start executes the corresponding workflow in this package.
func (m *PaymentModule) Start(_ context.Context) error {
	if m.svc == nil {
		return fmt.Errorf("payment service is not initialized")
	}
	if m.migrateDB == nil {
		return fmt.Errorf("payment migration database is not initialized")
	}
	return nil
}

// 中文：Stop 执行当前包中的对应流程。
// English: Stop executes the corresponding workflow in this package.
func (m *PaymentModule) Stop(ctx context.Context) error {
	if ctx != nil {
		return ctx.Err()
	}
	return nil
}

// 中文：loadAlipayCerts 执行当前包中的对应流程。
// English: loadAlipayCerts executes the corresponding workflow in this package.
func loadAlipayCerts(appCertPath, alipayCertPath, rootCertPath string) ([]byte, []byte, []byte, error) {
	if appCertPath == "" || alipayCertPath == "" || rootCertPath == "" {
		return nil, nil, nil, fmt.Errorf("payment alipay cert paths must include app_cert_path, alipay_cert_path, and root_cert_path")
	}
	appCert, err := os.ReadFile(appCertPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read alipay app cert: %w", err)
	}
	alipayCert, err := os.ReadFile(alipayCertPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read alipay public cert: %w", err)
	}
	rootCert, err := os.ReadFile(rootCertPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read alipay root cert: %w", err)
	}
	return appCert, alipayCert, rootCert, nil
}

// 中文：hasAnyValue 执行当前包中的对应流程。
// English: hasAnyValue executes the corresponding workflow in this package.
func hasAnyValue(values ...string) bool {
	for _, value := range values {
		if value != "" {
			return true
		}
	}
	return false
}

// 中文：firstNonEmpty 执行当前包中的对应流程。
// English: firstNonEmpty executes the corresponding workflow in this package.
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
