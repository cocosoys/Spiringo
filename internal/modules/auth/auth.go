package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/core/di"
	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/modules/auth/handler"
	"github.com/spiringo/spiringo/internal/modules/auth/oauth"
	"github.com/spiringo/spiringo/internal/modules/auth/repository"
	"github.com/spiringo/spiringo/internal/modules/auth/service"
	"github.com/spiringo/spiringo/internal/pkg/cache"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
// Config 认证模块配置
type Config struct {
	// 中文：JWT 保存当前结构中的配置或数据值。
	// English: JWT stores a configuration or data value for this struct.
	JWT struct {
		Secret        string `yaml:"secret" mapstructure:"secret"`
		AccessExpire  string `yaml:"access_expire" mapstructure:"access_expire"`
		RefreshExpire string `yaml:"refresh_expire" mapstructure:"refresh_expire"`
		Issuer        string `yaml:"issuer" mapstructure:"issuer"`
	} `yaml:"jwt" mapstructure:"jwt"`
	// 中文：OAuth 保存当前结构中的配置或数据值。
	// English: OAuth stores a configuration or data value for this struct.
	OAuth map[string]OAuthProviderConfig `yaml:"oauth" mapstructure:"oauth"`
}

// 中文：OAuthProviderConfig 定义当前包使用的数据结构或接口。
// English: OAuthProviderConfig defines a data structure or interface used by this package.
// OAuthProviderConfig OAuth提供者配置
type OAuthProviderConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：AppID 保存当前结构中的配置或数据值。
	// English: AppID stores a configuration or data value for this struct.
	AppID string `yaml:"app_id" mapstructure:"app_id"`
	// 中文：AppSecret 保存当前结构中的配置或数据值。
	// English: AppSecret stores a configuration or data value for this struct.
	AppSecret string `yaml:"app_secret" mapstructure:"app_secret"`
	// 中文：ClientID 保存当前结构中的配置或数据值。
	// English: ClientID stores a configuration or data value for this struct.
	ClientID string `yaml:"client_id" mapstructure:"client_id"`
	// 中文：ClientSecret 保存当前结构中的配置或数据值。
	// English: ClientSecret stores a configuration or data value for this struct.
	ClientSecret string `yaml:"client_secret" mapstructure:"client_secret"`
	// 中文：RedirectURL 保存当前结构中的配置或数据值。
	// English: RedirectURL stores a configuration or data value for this struct.
	RedirectURL string `yaml:"redirect_url" mapstructure:"redirect_url"`
	// 中文：CorpID 保存当前结构中的配置或数据值。
	// English: CorpID stores a configuration or data value for this struct.
	CorpID string `yaml:"corp_id" mapstructure:"corp_id"`
	// 中文：CorpSecret 保存当前结构中的配置或数据值。
	// English: CorpSecret stores a configuration or data value for this struct.
	CorpSecret string `yaml:"corp_secret" mapstructure:"corp_secret"`
	// 中文：AgentID 保存当前结构中的配置或数据值。
	// English: AgentID stores a configuration or data value for this struct.
	AgentID string `yaml:"agent_id" mapstructure:"agent_id"`
	// 中文：BotToken 保存当前结构中的配置或数据值。
	// English: BotToken stores a configuration or data value for this struct.
	BotToken string `yaml:"bot_token" mapstructure:"bot_token"`
}

// 中文：AuthModule 定义当前包使用的数据结构或接口。
// English: AuthModule defines a data structure or interface used by this package.
// AuthModule 认证模块
type AuthModule struct {
	// 中文：*module.BaseModule 嵌入复用该类型提供的能力。
	// English: *module.BaseModule embeds reusable behavior from that type.
	*module.BaseModule
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config Config
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *service.AuthService
	// 中文：migrateDB 保存当前结构中的配置或数据值。
	// English: migrateDB stores a configuration or data value for this struct.
	migrateDB *orm.DB
}

// 中文：NewAuthModule 创建并返回对应组件实例。
// English: NewAuthModule creates and returns the corresponding component instance.
// NewAuthModule 创建认证模块
func NewAuthModule() *AuthModule {
	return &AuthModule{
		BaseModule: module.NewBaseModule("auth", "user", "tenant"),
	}
}

// 中文：Config 执行当前包中的对应流程。
// English: Config executes the corresponding workflow in this package.
func (m *AuthModule) Config() any { return &m.config }

// 中文：Init 执行当前包中的对应流程。
// English: Init executes the corresponding workflow in this package.
func (m *AuthModule) Init(app *module.App) error {
	svcCfg := service.Config{}
	svcCfg.JWT.Secret = m.config.JWT.Secret
	svcCfg.JWT.Issuer = m.config.JWT.Issuer
	if d, err := time.ParseDuration(m.config.JWT.AccessExpire); err == nil {
		svcCfg.JWT.AccessExpire = d
	} else {
		svcCfg.JWT.AccessExpire = 2 * time.Hour
	}
	if d, err := time.ParseDuration(m.config.JWT.RefreshExpire); err == nil {
		svcCfg.JWT.RefreshExpire = d
	} else {
		svcCfg.JWT.RefreshExpire = 168 * time.Hour
	}

	db, err := di.Resolve[*orm.DB](app.DI)
	if err != nil {
		return fmt.Errorf("auth module init: %w", err)
	}
	m.migrateDB = db
	tdb := orm.NewTenantDB(db)
	authRepo := repository.NewAuthRepository(tdb)

	var c cache.Cache
	if di.Has[cache.Cache](app.DI) {
		c, _ = di.Resolve[cache.Cache](app.DI)
	}

	var userSvc service.UserServiceInterface
	// user模块通过ProvideNamed("auth_user_service", ...)注册，需要用ResolveNamed解析
	if usvc, err := di.ResolveNamed[service.UserServiceInterface](app.DI, "auth_user_service"); err == nil {
		userSvc = usvc
	}

	oauthReg := initOAuthRegistry(m.config.OAuth)
	m.svc = service.NewAuthService(svcCfg, app.EventBus, authRepo, userSvc, c, oauthReg)
	app.DI.Provide(m.svc)
	return nil
}

// 中文：initOAuthRegistry 执行当前包中的对应流程。
// English: initOAuthRegistry executes the corresponding workflow in this package.
func initOAuthRegistry(configs map[string]OAuthProviderConfig) *oauth.Registry {
	reg := oauth.NewRegistry()
	for name, cfg := range configs {
		if !cfg.Enabled {
			continue
		}
		switch name {
		case "wechat":
			reg.Register(oauth.NewWechatProvider(cfg.AppID, cfg.AppSecret))
		case "qq":
			reg.Register(oauth.NewQQProvider(cfg.AppID, cfg.AppSecret))
		case "google":
			reg.Register(oauth.NewGoogleProvider(cfg.ClientID, cfg.ClientSecret))
		case "discord":
			reg.Register(oauth.NewDiscordProvider(cfg.ClientID, cfg.ClientSecret))
		case "telegram":
			reg.Register(oauth.NewTelegramProvider(cfg.BotToken))
		case "x":
			reg.Register(oauth.NewXProvider(cfg.ClientID, cfg.ClientSecret))
		case "dingtalk":
			reg.Register(oauth.NewDingTalkProvider(cfg.AppID, cfg.AppSecret))
		case "douyin":
			reg.Register(oauth.NewDouyinProvider(cfg.ClientID, cfg.ClientSecret))
		case "work_wechat":
			reg.Register(oauth.NewWorkWechatProvider(cfg.CorpID, cfg.AgentID, cfg.CorpSecret))
		}
	}
	return reg
}

// 中文：Routes 执行当前包中的对应流程。
// English: Routes executes the corresponding workflow in this package.
func (m *AuthModule) Routes(r *gin.RouterGroup) {
	h := handler.NewAuthHandler(m.svc)
	r.POST("/login", h.Login)
	r.POST("/register", h.Register)
	r.POST("/logout", h.Logout)
	r.POST("/refresh", h.Refresh)
	r.GET("/oauth/:provider", h.OAuthLogin)
	r.GET("/oauth/:provider/callback", h.OAuthCallback)
}

// 中文：Start 执行当前包中的对应流程。
// English: Start executes the corresponding workflow in this package.
func (m *AuthModule) Start(_ context.Context) error {
	if m.svc == nil {
		return fmt.Errorf("auth service is not initialized")
	}
	if m.migrateDB == nil {
		return fmt.Errorf("auth migration database is not initialized")
	}
	return nil
}

// 中文：Stop 执行当前包中的对应流程。
// English: Stop executes the corresponding workflow in this package.
func (m *AuthModule) Stop(ctx context.Context) error {
	if ctx != nil {
		return ctx.Err()
	}
	return nil
}
