package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spiringo/spiringo/internal/core/config"
	"github.com/spiringo/spiringo/internal/core/di"
	"github.com/spiringo/spiringo/internal/core/event"
	"github.com/spiringo/spiringo/internal/core/middleware"
	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/core/server"
	alertpkg "github.com/spiringo/spiringo/internal/pkg/alert"
	cachepkg "github.com/spiringo/spiringo/internal/pkg/cache"
	lockpkg "github.com/spiringo/spiringo/internal/pkg/lock"
	loggerpkg "github.com/spiringo/spiringo/internal/pkg/logger"
	metricspkg "github.com/spiringo/spiringo/internal/pkg/metrics"
	mqpkg "github.com/spiringo/spiringo/internal/pkg/mq"
	"github.com/spiringo/spiringo/internal/pkg/orm"
	queuepkg "github.com/spiringo/spiringo/internal/pkg/queue"
	searchpkg "github.com/spiringo/spiringo/internal/pkg/search"
	storagepkg "github.com/spiringo/spiringo/internal/pkg/storage"
	tracepkg "github.com/spiringo/spiringo/internal/pkg/trace"
)

// 中文：App 定义当前包使用的数据结构或接口。
// English: App defines a data structure or interface used by this package.
// App owns the framework lifecycle and shared infrastructure.
type App struct {
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config *config.Manager
	// 中文：di 保存当前结构中的配置或数据值。
	// English: di stores a configuration or data value for this struct.
	di *di.Container
	// 中文：eventBus 保存当前结构中的配置或数据值。
	// English: eventBus stores a configuration or data value for this struct.
	eventBus *event.Bus
	// 中文：server 保存当前结构中的配置或数据值。
	// English: server stores a configuration or data value for this struct.
	server *server.Server
	// 中文：registry 保存当前结构中的配置或数据值。
	// English: registry stores a configuration or data value for this struct.
	registry *module.Registry
	// 中文：logger 保存当前结构中的配置或数据值。
	// English: logger stores a configuration or data value for this struct.
	logger *slog.Logger
	// 中文：loggerSync 保存当前结构中的配置或数据值。
	// English: loggerSync stores a configuration or data value for this struct.
	loggerSync func() error
	// 中文：configDir 保存当前结构中的配置或数据值。
	// English: configDir stores a configuration or data value for this struct.
	configDir string
	// 中文：env 保存当前结构中的配置或数据值。
	// English: env stores a configuration or data value for this struct.
	env string

	// 中文：db 保存当前结构中的配置或数据值。
	// English: db stores a configuration or data value for this struct.
	db *orm.DB
	// 中文：document 保存当前结构中的配置或数据值。
	// English: document stores a configuration or data value for this struct.
	document orm.DocumentStore
	// 中文：cache 保存当前结构中的配置或数据值。
	// English: cache stores a configuration or data value for this struct.
	cache cachepkg.Cache
	// 中文：lock 保存当前结构中的配置或数据值。
	// English: lock stores a configuration or data value for this struct.
	lock lockpkg.Lock
	// 中文：lockRedis 保存当前结构中的配置或数据值。
	// English: lockRedis stores a configuration or data value for this struct.
	lockRedis *redis.Client
	// 中文：metrics 保存当前结构中的配置或数据值。
	// English: metrics stores a configuration or data value for this struct.
	metrics *metricspkg.Registry
	// 中文：alerts 保存当前结构中的配置或数据值。
	// English: alerts stores a configuration or data value for this struct.
	alerts alertpkg.Notifier
	// 中文：mq 保存当前结构中的配置或数据值。
	// English: mq stores a configuration or data value for this struct.
	mq mqpkg.MQ
	// 中文：mqRedis 保存当前结构中的配置或数据值。
	// English: mqRedis stores a configuration or data value for this struct.
	mqRedis *redis.Client
	// 中文：queue 保存当前结构中的配置或数据值。
	// English: queue stores a configuration or data value for this struct.
	queue queuepkg.Queue
	// 中文：search 保存当前结构中的配置或数据值。
	// English: search stores a configuration or data value for this struct.
	search searchpkg.Engine
	// 中文：storage 保存当前结构中的配置或数据值。
	// English: storage stores a configuration or data value for this struct.
	storage storagepkg.Storage
	// 中文：tracer 保存当前结构中的配置或数据值。
	// English: tracer stores a configuration or data value for this struct.
	tracer *tracepkg.Tracer
}

// 中文：RuntimeReport 定义当前包使用的数据结构或接口。
// English: RuntimeReport defines a data structure or interface used by this package.
// RuntimeReport describes the live framework surface exposed by /system/report.
type RuntimeReport struct {
	// 中文：AppName 保存当前结构中的配置或数据值。
	// English: AppName stores a configuration or data value for this struct.
	AppName string `json:"app_name"`
	// 中文：Env 保存当前结构中的配置或数据值。
	// English: Env stores a configuration or data value for this struct.
	Env string `json:"env"`
	// 中文：Server 保存当前结构中的配置或数据值。
	// English: Server stores a configuration or data value for this struct.
	Server ServerSnapshot `json:"server"`
	// 中文：Modules 保存当前结构中的配置或数据值。
	// English: Modules stores a configuration or data value for this struct.
	Modules []module.ModuleSnapshot `json:"modules"`
	// 中文：Infrastructure 保存当前结构中的配置或数据值。
	// English: Infrastructure stores a configuration or data value for this struct.
	Infrastructure InfrastructureSnapshot `json:"infrastructure"`
	// 中文：Metrics 保存当前结构中的配置或数据值。
	// English: Metrics stores a configuration or data value for this struct.
	Metrics *metricspkg.Snapshot `json:"metrics,omitempty"`
}

// 中文：ServerSnapshot 定义当前包使用的数据结构或接口。
// English: ServerSnapshot defines a data structure or interface used by this package.
// ServerSnapshot records resolved runtime network addresses for deployment checks.
type ServerSnapshot struct {
	// 中文：ListenAddr 保存当前结构中的配置或数据值。
	// English: ListenAddr stores a configuration or data value for this struct.
	ListenAddr string `json:"listen_addr"`
	// 中文：PublicURL 保存当前结构中的配置或数据值。
	// English: PublicURL stores a configuration or data value for this struct.
	PublicURL string `json:"public_url,omitempty"`
	// 中文：APIBaseURL 保存当前结构中的配置或数据值。
	// English: APIBaseURL stores a configuration or data value for this struct.
	APIBaseURL string `json:"api_base_url,omitempty"`
}

// 中文：InfrastructureSnapshot 定义当前包使用的数据结构或接口。
// English: InfrastructureSnapshot defines a data structure or interface used by this package.
// InfrastructureSnapshot records which optional infrastructure components are active.
type InfrastructureSnapshot struct {
	// 中文：Database 保存当前结构中的配置或数据值。
	// English: Database stores a configuration or data value for this struct.
	Database bool `json:"database"`
	// 中文：Document 保存当前结构中的配置或数据值。
	// English: Document stores a configuration or data value for this struct.
	Document bool `json:"document"`
	// 中文：Cache 保存当前结构中的配置或数据值。
	// English: Cache stores a configuration or data value for this struct.
	Cache bool `json:"cache"`
	// 中文：Lock 保存当前结构中的配置或数据值。
	// English: Lock stores a configuration or data value for this struct.
	Lock bool `json:"lock"`
	// 中文：Metrics 保存当前结构中的配置或数据值。
	// English: Metrics stores a configuration or data value for this struct.
	Metrics bool `json:"metrics"`
	// 中文：Alerts 保存当前结构中的配置或数据值。
	// English: Alerts stores a configuration or data value for this struct.
	Alerts bool `json:"alerts"`
	// 中文：MQ 保存当前结构中的配置或数据值。
	// English: MQ stores a configuration or data value for this struct.
	MQ bool `json:"mq"`
	// 中文：Queue 保存当前结构中的配置或数据值。
	// English: Queue stores a configuration or data value for this struct.
	Queue bool `json:"queue"`
	// 中文：Search 保存当前结构中的配置或数据值。
	// English: Search stores a configuration or data value for this struct.
	Search bool `json:"search"`
	// 中文：Storage 保存当前结构中的配置或数据值。
	// English: Storage stores a configuration or data value for this struct.
	Storage bool `json:"storage"`
	// 中文：Trace 保存当前结构中的配置或数据值。
	// English: Trace stores a configuration or data value for this struct.
	Trace bool `json:"trace"`
}

// 中文：Option 定义当前包使用的数据结构或接口。
// English: Option defines a data structure or interface used by this package.
// Option configures an App.
type Option func(*App)

// 中文：WithConfigDir 执行当前包中的对应流程。
// English: WithConfigDir executes the corresponding workflow in this package.
// WithConfigDir sets the configuration directory.
func WithConfigDir(dir string) Option {
	return func(a *App) {
		a.configDir = dir
	}
}

// 中文：WithEnv 执行当前包中的对应流程。
// English: WithEnv executes the corresponding workflow in this package.
// WithEnv sets the runtime environment name.
func WithEnv(env string) Option {
	return func(a *App) {
		a.env = config.NormalizeEnv(env)
	}
}

// 中文：New 创建并返回对应组件实例。
// English: New creates and returns the corresponding component instance.
// New creates an application instance.
func New(opts ...Option) *App {
	a := &App{
		config:   config.NewManager(),
		di:       di.NewContainer(),
		eventBus: event.NewBus(4),
		registry: module.NewRegistry(),
		logger:   slog.Default(),
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// 中文：Config 执行当前包中的对应流程。
// English: Config executes the corresponding workflow in this package.
// Config returns the configuration manager.
func (a *App) Config() *config.Manager {
	return a.config
}

// 中文：DI 执行当前包中的对应流程。
// English: DI executes the corresponding workflow in this package.
// DI returns the dependency container.
func (a *App) DI() *di.Container {
	return a.di
}

// 中文：EventBus 执行当前包中的对应流程。
// English: EventBus executes the corresponding workflow in this package.
// EventBus returns the application event bus.
func (a *App) EventBus() *event.Bus {
	return a.eventBus
}

// 中文：Registry 执行当前包中的对应流程。
// English: Registry executes the corresponding workflow in this package.
// Registry returns the module registry.
func (a *App) Registry() *module.Registry {
	return a.registry
}

// 中文：Logger 执行当前包中的对应流程。
// English: Logger executes the corresponding workflow in this package.
// Logger returns the application logger.
func (a *App) Logger() *slog.Logger {
	return a.logger
}

// 中文：SetLogger 执行当前包中的对应流程。
// English: SetLogger executes the corresponding workflow in this package.
// SetLogger replaces the application logger.
func (a *App) SetLogger(l *slog.Logger) {
	a.logger = l
	slog.SetDefault(l)
	a.di.Provide(l)
}

// 中文：RegisterModules 执行当前包中的对应流程。
// English: RegisterModules executes the corresponding workflow in this package.
// RegisterModules registers business modules.
func (a *App) RegisterModules(modules ...module.Module) {
	a.registry.RegisterModules(modules...)
}

// 中文：Migrate 执行当前包中的对应流程。
// English: Migrate executes the corresponding workflow in this package.
// Migrate initializes infrastructure and runs enabled module migrations without starting HTTP service.
func (a *App) Migrate(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := a.loadRuntimeConfig(); err != nil {
		return err
	}
	if err := a.initLogger(); err != nil {
		return err
	}
	if err := a.initTrace(); err != nil {
		return err
	}
	if err := a.initInfrastructure(); err != nil {
		return err
	}
	defer a.closeInfrastructure()

	moduleApp := &module.App{
		Config:   a.config,
		DI:       a.di,
		EventBus: a.eventBus,
		Router:   gin.New(),
		Modules:  a.registry,
		Migrate:  newMigrationStore(a.db),
	}
	if err := a.registry.InitAll(moduleApp); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	return ctx.Err()
}

// 中文：Run 执行当前包中的对应流程。
// English: Run executes the corresponding workflow in this package.
// Run starts the application and blocks until SIGINT or SIGTERM.
func (a *App) Run() error {
	if err := a.loadRuntimeConfig(); err != nil {
		return err
	}
	if err := a.initLogger(); err != nil {
		return err
	}
	if err := a.initTrace(); err != nil {
		return err
	}
	if err := a.initInfrastructure(); err != nil {
		return err
	}
	defer a.closeInfrastructure()

	serverCfg, err := a.serverConfig()
	if err != nil {
		return err
	}
	mwCfg, err := a.middlewareConfig()
	if err != nil {
		return err
	}

	a.server = server.New(serverCfg, mwCfg)
	if a.tracer != nil {
		a.server.Engine().Use(middleware.Trace(a.tracer))
	}
	if err := a.installMetrics(); err != nil {
		return err
	}
	if err := a.initAlert(); err != nil {
		return err
	}
	moduleApp := &module.App{
		Config:   a.config,
		DI:       a.di,
		EventBus: a.eventBus,
		Router:   a.server.Engine(),
		Modules:  a.registry,
		Migrate:  newMigrationStore(a.db),
	}

	if err := a.registry.InitAll(moduleApp); err != nil {
		return fmt.Errorf("init modules: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := a.config.Watch(ctx); err != nil {
		a.logger.Warn("config watch not started", "error", err)
	}
	a.eventBus.Start(ctx)
	if a.queue != nil {
		a.queue.Start(ctx)
	}

	if err := a.registry.StartAll(ctx); err != nil {
		return fmt.Errorf("start modules: %w", err)
	}

	a.registerHealthRoutes()
	a.registerSystemRoutes()

	a.logger.Info("starting server", "addr", serverCfg.Addr, "env", a.env)
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	if err := a.waitStop(); err != nil {
		return err
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	a.logger.Info("shutting down")
	_ = a.server.Stop(shutdownCtx)
	_ = a.registry.StopAll(shutdownCtx)
	a.eventBus.Stop()

	a.logger.Info("server stopped")
	return nil
}

// 中文：loadRuntimeConfig 执行当前包中的对应流程。
// English: loadRuntimeConfig executes the corresponding workflow in this package.
func (a *App) loadRuntimeConfig() error {
	if a.configDir == "" {
		a.configDir = "configs"
	}
	a.config.SetConfigDir(a.configDir)
	if env := firstNonEmpty(a.env, os.Getenv("APP_ENV")); env != "" {
		a.config.SetEnv(env)
	}

	if err := a.config.Load(); err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	a.env = firstNonEmpty(a.config.Env(), "local")
	a.config.SetEnv(a.env)
	a.config.Set("app.env", a.env)
	if err := a.initConfigCenter(); err != nil {
		return err
	}
	return nil
}

// 中文：initConfigCenter 执行当前包中的对应流程。
// English: initConfigCenter executes the corresponding workflow in this package.
func (a *App) initConfigCenter() error {
	var cfg config.ConfigCenterConfig
	if err := a.config.Unmarshal("config_center", &cfg); err != nil {
		return fmt.Errorf("parse config_center config: %w", err)
	}
	if !cfg.Enabled {
		return nil
	}

	switch strings.ToLower(strings.TrimSpace(cfg.Type)) {
	case "nacos":
		a.config.AddSource(config.NewNacosSource(cfg.Nacos, 50))
	case "consul":
		a.config.AddSource(config.NewConsulSource(cfg.Consul, 50))
	default:
		return fmt.Errorf("unsupported config_center type: %s", cfg.Type)
	}

	if err := a.config.Load(); err != nil {
		return fmt.Errorf("reload config with config_center: %w", err)
	}
	return nil
}

// 中文：initLogger 执行当前包中的对应流程。
// English: initLogger executes the corresponding workflow in this package.
func (a *App) initLogger() error {
	var cfg config.LogConfig
	if err := a.config.Unmarshal("log", &cfg); err != nil {
		return fmt.Errorf("parse log config: %w", err)
	}

	switch strings.ToLower(strings.TrimSpace(cfg.Driver)) {
	case "", "slog":
	case "zap":
		l, syncFn, err := loggerpkg.NewZapSlog(cfg)
		if err != nil {
			return fmt.Errorf("init zap logger: %w", err)
		}
		a.loggerSync = syncFn
		a.SetLogger(l)
		return nil
	default:
		return fmt.Errorf("unsupported log driver: %s", cfg.Driver)
	}

	var out io.Writer = os.Stdout
	if strings.EqualFold(cfg.Output, "stderr") {
		out = os.Stderr
	}

	opts := &slog.HandlerOptions{Level: parseSlogLevel(cfg.Level)}
	var handler slog.Handler
	if strings.EqualFold(cfg.Format, "json") {
		handler = slog.NewJSONHandler(out, opts)
	} else {
		handler = slog.NewTextHandler(out, opts)
	}

	a.SetLogger(slog.New(handler))
	return nil
}

// 中文：initTrace 执行当前包中的对应流程。
// English: initTrace executes the corresponding workflow in this package.
func (a *App) initTrace() error {
	var cfg config.TraceConfig
	if err := a.config.Unmarshal("trace", &cfg); err != nil {
		return fmt.Errorf("parse trace config: %w", err)
	}
	if !cfg.Enabled {
		return nil
	}
	exporter, err := a.traceExporter(cfg)
	if err != nil {
		return err
	}
	tracer := tracepkg.NewTracer(exporter)
	a.tracer = tracer
	a.di.Provide(tracer)
	a.logger.Info("trace initialized", "exporter", firstNonEmpty(cfg.Exporter, "logger"))
	return nil
}

// 中文：traceExporter 执行当前包中的对应流程。
// English: traceExporter executes the corresponding workflow in this package.
func (a *App) traceExporter(cfg config.TraceConfig) (tracepkg.Exporter, error) {
	exporterType := strings.ToLower(strings.TrimSpace(cfg.Exporter))
	if exporterType == "" {
		exporterType = "logger"
	}
	loggerExporter := tracepkg.NewLoggerExporter(a.logger)
	switch exporterType {
	case "logger":
		return loggerExporter, nil
	case "otlp":
		return a.otlpExporter(cfg)
	case "both", "logger+otlp", "otlp+logger":
		otlpExporter, err := a.otlpExporter(cfg)
		if err != nil {
			return nil, err
		}
		return tracepkg.NewMultiExporter(loggerExporter, otlpExporter), nil
	default:
		return nil, fmt.Errorf("unsupported trace exporter: %s", cfg.Exporter)
	}
}

// 中文：otlpExporter 执行当前包中的对应流程。
// English: otlpExporter executes the corresponding workflow in this package.
func (a *App) otlpExporter(cfg config.TraceConfig) (tracepkg.Exporter, error) {
	timeout, err := parseOptionalDuration(cfg.OTLP.Timeout)
	if err != nil {
		return nil, fmt.Errorf("parse trace.otlp.timeout: %w", err)
	}
	serviceName := firstNonEmpty(cfg.Service.Name, a.config.GetString("app.name"), "spiringo")
	exporter, err := tracepkg.NewOTLPHTTPExporter(tracepkg.OTLPHTTPConfig{
		Endpoint:    cfg.OTLP.Endpoint,
		ServiceName: serviceName,
		Timeout:     timeout,
		Headers:     cfg.OTLP.Headers,
	})
	if err != nil {
		return nil, fmt.Errorf("init otlp trace exporter: %w", err)
	}
	return exporter, nil
}

// 中文：initInfrastructure 执行当前包中的对应流程。
// English: initInfrastructure executes the corresponding workflow in this package.
func (a *App) initInfrastructure() error {
	if err := a.initDatabase(); err != nil {
		return err
	}
	if err := a.initDocumentStore(); err != nil {
		return err
	}
	if err := a.initCache(); err != nil {
		return err
	}
	if err := a.initLock(); err != nil {
		return err
	}
	if err := a.initStorage(); err != nil {
		return err
	}
	if err := a.initSearch(); err != nil {
		return err
	}
	if err := a.initMQ(); err != nil {
		return err
	}
	if err := a.initTaskQueue(); err != nil {
		return err
	}
	return nil
}

// 中文：initDatabase 执行当前包中的对应流程。
// English: initDatabase executes the corresponding workflow in this package.
func (a *App) initDatabase() error {
	if !a.config.IsSet("database.default.driver") && !a.config.IsSet("database.default.dsn") {
		return nil
	}

	var cfg config.DatabaseConfig
	if err := a.config.Unmarshal("database.default", &cfg); err != nil {
		return fmt.Errorf("parse database config: %w", err)
	}
	if cfg.Driver == "" {
		return fmt.Errorf("database.default.driver is required")
	}

	lifetime, err := parseOptionalDuration(cfg.ConnMaxLifetime)
	if err != nil {
		return fmt.Errorf("parse database conn_max_lifetime: %w", err)
	}
	if cfg.MaxIdle <= 0 {
		cfg.MaxIdle = 5
	}
	if cfg.MaxOpen <= 0 {
		cfg.MaxOpen = 20
	}

	db, err := orm.New(orm.Config{
		Driver:          cfg.Driver,
		DSN:             cfg.DSN,
		MaxIdle:         cfg.MaxIdle,
		MaxOpen:         cfg.MaxOpen,
		ConnMaxLifetime: lifetime,
		LogLevel:        a.config.GetString("log.level"),
		ReadReplicas:    ormReadReplicas(cfg.ReadReplicas),
	})
	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return fmt.Errorf("ping database: %w", err)
	}

	a.db = db
	a.di.Provide(db)
	a.logger.Info("database initialized", "driver", cfg.Driver, "read_replicas", len(cfg.ReadReplicas))
	return nil
}

// 中文：ormReadReplicas 执行当前包中的对应流程。
// English: ormReadReplicas executes the corresponding workflow in this package.
func ormReadReplicas(replicas []config.DBConfig) []orm.EndpointConfig {
	if len(replicas) == 0 {
		return nil
	}
	result := make([]orm.EndpointConfig, 0, len(replicas))
	for _, replica := range replicas {
		if replica.DSN == "" {
			continue
		}
		result = append(result, orm.EndpointConfig{
			Driver: replica.Driver,
			DSN:    replica.DSN,
		})
	}
	return result
}

// 中文：initDocumentStore 执行当前包中的对应流程。
// English: initDocumentStore executes the corresponding workflow in this package.
func (a *App) initDocumentStore() error {
	var cfg config.DocumentConfig
	if err := a.config.Unmarshal("document", &cfg); err != nil {
		return fmt.Errorf("parse document config: %w", err)
	}
	if !cfg.Enabled {
		return nil
	}

	driver := strings.ToLower(strings.TrimSpace(cfg.Driver))
	if driver == "" {
		driver = "mongodb"
	}
	if driver != "mongodb" {
		return fmt.Errorf("unsupported document driver: %s", cfg.Driver)
	}

	timeout, err := parseOptionalDuration(cfg.MongoDB.Timeout)
	if err != nil {
		return fmt.Errorf("parse document.mongodb.timeout: %w", err)
	}
	store, err := orm.NewMongoStore(context.Background(), orm.MongoConfig{
		URI:      cfg.MongoDB.URI,
		Database: cfg.MongoDB.Database,
		Timeout:  timeout,
	})
	if err != nil {
		return fmt.Errorf("init mongodb: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(context.Background(), firstPositiveDuration(timeout, 3*time.Second))
	defer cancel()
	if err := store.Ping(pingCtx); err != nil {
		_ = store.Close(context.Background())
		return fmt.Errorf("ping mongodb: %w", err)
	}

	a.document = store
	a.di.Provide(store)
	di.ProvideAs[orm.DocumentStore](a.di, store)
	a.logger.Info("document database initialized", "driver", driver)
	return nil
}

// 中文：initCache 执行当前包中的对应流程。
// English: initCache executes the corresponding workflow in this package.
func (a *App) initCache() error {
	var cfg config.CacheConfig
	if err := a.config.Unmarshal("cache", &cfg); err != nil {
		return fmt.Errorf("parse cache config: %w", err)
	}

	driver := strings.ToLower(strings.TrimSpace(cfg.Driver))
	if driver == "" {
		driver = "memory"
	}

	var c cachepkg.Cache
	switch driver {
	case "memory":
		c = cachepkg.NewMemoryCache()
	case "redis":
		if cfg.Redis.Addr == "" {
			if err := a.config.Unmarshal("redis.default", &cfg.Redis); err != nil {
				return fmt.Errorf("parse redis default config: %w", err)
			}
		}
		if cfg.Redis.Addr == "" {
			return fmt.Errorf("cache.redis.addr is required when cache.driver=redis")
		}
		rc := cachepkg.NewRedisCache(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := rc.Ping(ctx); err != nil {
			_ = rc.Close()
			return fmt.Errorf("ping redis cache: %w", err)
		}
		c = rc
	default:
		return fmt.Errorf("unsupported cache driver: %s", cfg.Driver)
	}

	a.cache = c
	a.di.Provide(c)
	di.ProvideAs[cachepkg.Cache](a.di, c)
	a.logger.Info("cache initialized", "driver", driver)
	return nil
}

// 中文：initLock 执行当前包中的对应流程。
// English: initLock executes the corresponding workflow in this package.
func (a *App) initLock() error {
	var cfg config.LockConfig
	if err := a.config.Unmarshal("lock", &cfg); err != nil {
		return fmt.Errorf("parse lock config: %w", err)
	}
	if !cfg.Enabled {
		return nil
	}

	driver := strings.ToLower(strings.TrimSpace(cfg.Driver))
	if driver == "" {
		driver = "redis"
	}

	var manager lockpkg.Lock
	switch driver {
	case "redis":
		if cfg.Redis.Addr == "" {
			if err := a.config.Unmarshal("redis.default", &cfg.Redis); err != nil {
				return fmt.Errorf("parse redis default config: %w", err)
			}
		}
		if cfg.Redis.Addr == "" {
			return fmt.Errorf("lock.redis.addr is required when lock.driver=redis")
		}
		client := redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := client.Ping(ctx).Err(); err != nil {
			_ = client.Close()
			return fmt.Errorf("ping lock redis: %w", err)
		}
		a.lockRedis = client
		manager = lockpkg.NewRedisLock(client)
	case "zookeeper", "zk":
		if len(cfg.ZooKeeper.Servers) == 0 {
			cfg.ZooKeeper.Servers = splitCSV(a.config.GetString("lock.zookeeper.servers"))
		}
		timeout, err := parseOptionalDuration(cfg.ZooKeeper.SessionTimeout)
		if err != nil {
			return fmt.Errorf("parse lock.zookeeper.session_timeout: %w", err)
		}
		zkLock, err := lockpkg.NewZooKeeperLock(lockpkg.ZooKeeperConfig{
			Servers:        cfg.ZooKeeper.Servers,
			Root:           cfg.ZooKeeper.Root,
			SessionTimeout: timeout,
		})
		if err != nil {
			return fmt.Errorf("init zookeeper lock: %w", err)
		}
		manager = zkLock
	default:
		return fmt.Errorf("unsupported lock driver: %s", cfg.Driver)
	}

	a.lock = manager
	a.di.Provide(manager)
	di.ProvideAs[lockpkg.Lock](a.di, manager)
	a.logger.Info("lock initialized", "driver", driver)
	return nil
}

// 中文：initStorage 执行当前包中的对应流程。
// English: initStorage executes the corresponding workflow in this package.
func (a *App) initStorage() error {
	var cfg config.StorageConfig
	if err := a.config.Unmarshal("storage", &cfg); err != nil {
		return fmt.Errorf("parse storage config: %w", err)
	}
	if !cfg.Enabled {
		return nil
	}

	storageType := strings.ToLower(strings.TrimSpace(cfg.Type))
	if storageType == "" {
		storageType = "minio"
	}
	var (
		s   storagepkg.Storage
		err error
	)
	switch storageType {
	case "minio":
		if cfg.MinIO.Endpoint == "" {
			return fmt.Errorf("storage.minio.endpoint is required when storage is enabled")
		}
		s, err = storagepkg.NewMinIOStorage(storagepkg.MinIOConfig{
			Endpoint:  cfg.MinIO.Endpoint,
			AccessKey: cfg.MinIO.AccessKey,
			SecretKey: cfg.MinIO.SecretKey,
			UseSSL:    cfg.MinIO.UseSSL,
		})
		if err != nil {
			return fmt.Errorf("init minio storage: %w", err)
		}
	case "ceph":
		if cfg.Ceph.Endpoint == "" {
			return fmt.Errorf("storage.ceph.endpoint is required when storage type=ceph")
		}
		s, err = storagepkg.NewCephStorage(storagepkg.CephConfig{
			Endpoint:  cfg.Ceph.Endpoint,
			AccessKey: cfg.Ceph.AccessKey,
			SecretKey: cfg.Ceph.SecretKey,
			UseSSL:    cfg.Ceph.UseSSL,
			PublicURL: cfg.Ceph.PublicURL,
		})
		if err != nil {
			return fmt.Errorf("init ceph storage: %w", err)
		}
	default:
		return fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}

	a.storage = s
	a.di.Provide(s)
	di.ProvideAs[storagepkg.Storage](a.di, s)
	a.logger.Info("storage initialized", "type", storageType)
	return nil
}

// 中文：initSearch 执行当前包中的对应流程。
// English: initSearch executes the corresponding workflow in this package.
func (a *App) initSearch() error {
	var cfg config.SearchConfig
	if err := a.config.Unmarshal("search", &cfg); err != nil {
		return fmt.Errorf("parse search config: %w", err)
	}
	if !cfg.Enabled {
		return nil
	}

	driver := strings.ToLower(strings.TrimSpace(cfg.Driver))
	if driver == "" {
		driver = "elasticsearch"
	}
	if driver != "elasticsearch" {
		return fmt.Errorf("unsupported search driver: %s", cfg.Driver)
	}

	timeout, err := parseOptionalDuration(cfg.Elasticsearch.Timeout)
	if err != nil {
		return fmt.Errorf("parse search.elasticsearch.timeout: %w", err)
	}
	client, err := searchpkg.NewElasticsearch(searchpkg.ElasticsearchConfig{
		Endpoint: cfg.Elasticsearch.Endpoint,
		Username: cfg.Elasticsearch.Username,
		Password: cfg.Elasticsearch.Password,
		Index:    cfg.Elasticsearch.Index,
		Timeout:  timeout,
	})
	if err != nil {
		return fmt.Errorf("init elasticsearch: %w", err)
	}

	a.search = client
	a.di.Provide(client)
	di.ProvideAs[searchpkg.Engine](a.di, client)
	a.logger.Info("search initialized", "driver", driver)
	return nil
}

// 中文：initMQ 执行当前包中的对应流程。
// English: initMQ executes the corresponding workflow in this package.
func (a *App) initMQ() error {
	var cfg config.MQConfig
	if err := a.config.Unmarshal("mq", &cfg); err != nil {
		return fmt.Errorf("parse mq config: %w", err)
	}
	if !cfg.Enabled {
		return nil
	}

	driver := strings.ToLower(strings.TrimSpace(cfg.Driver))
	if driver == "" {
		driver = "redis_stream"
	}

	var queue mqpkg.MQ
	switch driver {
	case "redis_stream", "redis":
		if cfg.Redis.Addr == "" {
			var redisCfg config.RedisConfig
			if err := a.config.Unmarshal("redis.default", &redisCfg); err != nil {
				return fmt.Errorf("parse redis default config: %w", err)
			}
			cfg.Redis.Addr = redisCfg.Addr
			cfg.Redis.Password = redisCfg.Password
			cfg.Redis.DB = redisCfg.DB
		}
		if cfg.Redis.Addr == "" {
			return fmt.Errorf("mq.redis.addr is required when mq.driver=redis_stream")
		}

		client := redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := client.Ping(ctx).Err(); err != nil {
			_ = client.Close()
			return fmt.Errorf("ping mq redis: %w", err)
		}
		a.mqRedis = client
		queue = mqpkg.NewRedisStreamMQ(client, cfg.Redis.Prefix)
	case "rabbitmq", "rabbit", "amqp":
		if cfg.RabbitMQ.URL == "" {
			return fmt.Errorf("mq.rabbitmq.url is required when mq.driver=rabbitmq")
		}
		rabbit, err := mqpkg.NewRabbitMQ(mqpkg.RabbitMQConfig{
			URL:         cfg.RabbitMQ.URL,
			Exchange:    cfg.RabbitMQ.Exchange,
			QueuePrefix: cfg.RabbitMQ.QueuePrefix,
		})
		if err != nil {
			return fmt.Errorf("init rabbitmq: %w", err)
		}
		queue = rabbit
	case "kafka":
		if len(cfg.Kafka.Brokers) == 0 {
			cfg.Kafka.Brokers = splitCSV(a.config.GetString("mq.kafka.brokers"))
		}
		kafkaMQ, err := mqpkg.NewKafkaMQ(mqpkg.KafkaConfig{
			Brokers:     cfg.Kafka.Brokers,
			ClientID:    firstNonEmpty(cfg.Kafka.ClientID, "spiringo"),
			GroupID:     cfg.Kafka.GroupID,
			TopicPrefix: cfg.Kafka.TopicPrefix,
		})
		if err != nil {
			return fmt.Errorf("init kafka mq: %w", err)
		}
		queue = kafkaMQ
	default:
		return fmt.Errorf("unsupported mq driver: %s", cfg.Driver)
	}

	bridge := event.NewMQBridge(a.eventBus, queue, event.MQBridgeConfig{Source: "spiringo"})
	bridge.Register()

	a.mq = queue
	a.di.Provide(queue)
	di.ProvideAs[mqpkg.MQ](a.di, queue)
	a.logger.Info("mq initialized", "driver", driver)
	return nil
}

// 中文：initTaskQueue 执行当前包中的对应流程。
// English: initTaskQueue executes the corresponding workflow in this package.
func (a *App) initTaskQueue() error {
	var cfg config.QueueConfig
	if err := a.config.Unmarshal("queue", &cfg); err != nil {
		return fmt.Errorf("parse queue config: %w", err)
	}
	if !cfg.Enabled {
		return nil
	}

	retryDelay, err := parseOptionalDuration(cfg.RetryDelay)
	if err != nil {
		return fmt.Errorf("parse queue.retry_delay: %w", err)
	}
	q := queuepkg.NewMemoryQueue(queuepkg.Config{
		Workers:    cfg.Workers,
		Buffer:     cfg.Buffer,
		MaxRetries: cfg.MaxRetries,
		RetryDelay: retryDelay,
	})
	q.SetErrorHandler(func(ctx context.Context, task *queuepkg.Task, err error) {
		taskName := ""
		if task != nil {
			taskName = task.Name
		}
		if a.metrics != nil {
			a.metrics.IncCounter("queue_task_errors_total", metricspkg.Labels{"task": taskName})
		}
		if a.alerts != nil {
			if notifyErr := a.alerts.Notify(ctx, alertpkg.Alert{
				Title:    "queue task failed",
				Message:  err.Error(),
				Severity: alertpkg.SeverityWarning,
				Source:   "task_queue",
				Labels:   map[string]string{"task": taskName},
			}); notifyErr != nil {
				a.logger.WarnContext(ctx, "send queue alert failed", "error", notifyErr)
			}
		}
	})

	a.queue = q
	a.di.Provide(q)
	di.ProvideAs[queuepkg.Queue](a.di, q)
	a.logger.Info("task queue initialized", "workers", cfg.Workers, "buffer", cfg.Buffer)
	return nil
}

// 中文：initAlert 执行当前包中的对应流程。
// English: initAlert executes the corresponding workflow in this package.
func (a *App) initAlert() error {
	var cfg config.AlertConfig
	if err := a.config.Unmarshal("alert", &cfg); err != nil {
		return fmt.Errorf("parse alert config: %w", err)
	}
	if !cfg.Enabled {
		return nil
	}

	sinks := make([]alertpkg.Notifier, 0, 2)
	if cfg.Logger {
		sinks = append(sinks, alertpkg.NewLoggerNotifier(a.logger))
	}
	if cfg.Webhook.Enabled {
		timeout, err := parseOptionalDuration(cfg.Webhook.Timeout)
		if err != nil {
			return fmt.Errorf("parse alert.webhook.timeout: %w", err)
		}
		sinks = append(sinks, alertpkg.NewWebhookNotifier(cfg.Webhook.URL, timeout, cfg.Webhook.Headers))
	}
	if cfg.Sentry.Enabled {
		timeout, err := parseOptionalDuration(cfg.Sentry.FlushTimeout)
		if err != nil {
			return fmt.Errorf("parse alert.sentry.flush_timeout: %w", err)
		}
		sentryNotifier, err := alertpkg.NewSentryNotifier(alertpkg.SentryConfig{
			DSN:              cfg.Sentry.DSN,
			Environment:      firstNonEmpty(cfg.Sentry.Environment, a.env, a.config.GetString("app.env")),
			Release:          cfg.Sentry.Release,
			TracesSampleRate: cfg.Sentry.TracesSampleRate,
			Debug:            cfg.Sentry.Debug,
			FlushTimeout:     timeout,
		})
		if err != nil {
			return fmt.Errorf("init sentry alert: %w", err)
		}
		sinks = append(sinks, sentryNotifier)
	}
	if len(sinks) == 0 {
		sinks = append(sinks, alertpkg.NewLoggerNotifier(a.logger))
	}

	manager := alertpkg.NewManager(sinks...)
	a.alerts = manager
	a.di.Provide(manager)
	di.ProvideAs[alertpkg.Notifier](a.di, manager)
	a.eventBus.SetErrorHandler(func(ctx context.Context, e *event.Event, err error) {
		if a.metrics != nil && e != nil {
			a.metrics.IncCounter("event_handler_errors_total", metricspkg.Labels{"topic": e.Topic})
		}
		title := "event handler failed"
		labels := map[string]string{}
		if e != nil {
			labels["topic"] = e.Topic
			labels["source"] = e.Source
		}
		if notifyErr := manager.Notify(ctx, alertpkg.Alert{
			Title:    title,
			Message:  err.Error(),
			Severity: alertpkg.SeverityWarning,
			Source:   "event_bus",
			Labels:   labels,
		}); notifyErr != nil {
			a.logger.WarnContext(ctx, "send alert failed", "error", notifyErr)
		}
	})
	return nil
}

// 中文：serverConfig 执行当前包中的对应流程。
// English: serverConfig executes the corresponding workflow in this package.
func (a *App) serverConfig() (config.ServerConfig, error) {
	var cfg config.ServerConfig
	if err := a.config.Unmarshal("server", &cfg); err != nil {
		return cfg, fmt.Errorf("parse server config: %w", err)
	}
	addr, err := resolveServerListenAddr(cfg)
	if err != nil {
		return cfg, err
	}
	cfg.Addr = addr
	if cfg.Mode == "" {
		cfg.Mode = gin.DebugMode
	}
	if cfg.PublicURL == "" {
		cfg.PublicURL = deriveServerPublicURL(cfg)
	}
	if cfg.APIBaseURL == "" && cfg.PublicURL != "" {
		cfg.APIBaseURL = strings.TrimRight(cfg.PublicURL, "/") + "/api/v1"
	}
	a.config.Set("server.addr", cfg.Addr)
	a.config.Set("server.public_url", cfg.PublicURL)
	a.config.Set("server.api_base_url", cfg.APIBaseURL)
	return cfg, nil
}

// 中文：resolveServerListenAddr 根据配置解析最终监听地址。
// English: resolveServerListenAddr resolves the final listen address from configuration.
func resolveServerListenAddr(cfg config.ServerConfig) (string, error) {
	if addr := strings.TrimSpace(cfg.Addr); addr != "" {
		return addr, nil
	}
	if cfg.Port <= 0 {
		return "", fmt.Errorf("server.port is required when server.addr is empty")
	}
	return net.JoinHostPort(strings.TrimSpace(cfg.Host), strconv.Itoa(cfg.Port)), nil
}

// 中文：deriveServerPublicURL 从监听配置推导仅用于开发兜底的公开地址。
// English: deriveServerPublicURL derives a development fallback public URL from listen settings.
func deriveServerPublicURL(cfg config.ServerConfig) string {
	host := strings.TrimSpace(cfg.Host)
	port := cfg.Port
	if port <= 0 && strings.TrimSpace(cfg.Addr) != "" {
		if splitHost, splitPort, err := net.SplitHostPort(cfg.Addr); err == nil {
			host = splitHost
			if parsed, parseErr := strconv.Atoi(splitPort); parseErr == nil {
				port = parsed
			}
		}
	}
	if port <= 0 {
		return ""
	}
	switch host {
	case "", "0.0.0.0", "::", "[::]":
		host = "127.0.0.1"
	}
	host = strings.Trim(host, "[]")
	if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	return "http://" + host + ":" + strconv.Itoa(port)
}

// 中文：middlewareConfig 执行当前包中的对应流程。
// English: middlewareConfig executes the corresponding workflow in this package.
func (a *App) middlewareConfig() (config.MiddlewareConfig, error) {
	var cfg config.MiddlewareConfig
	if err := a.config.Unmarshal("middleware", &cfg); err != nil {
		return cfg, fmt.Errorf("parse middleware config: %w", err)
	}
	if cfg.Auth.Enabled {
		if strings.TrimSpace(cfg.Auth.JWTSecret) == "" {
			cfg.Auth.JWTSecret = a.config.GetString("modules.auth.jwt.secret")
		}
		if strings.TrimSpace(cfg.Auth.JWTSecret) == "" {
			return cfg, fmt.Errorf("middleware.auth.jwt_secret or modules.auth.jwt.secret is required when middleware.auth is enabled")
		}
		if len(cfg.Auth.PublicPaths) == 0 {
			cfg.Auth.PublicPaths = defaultGlobalAuthPublicPaths(
				a.config.GetString("metrics.path"),
				a.config.GetString("metrics.report_path"),
			)
		}
	}
	return cfg, nil
}

// 中文：defaultGlobalAuthPublicPaths 执行当前包中的对应流程。
// English: defaultGlobalAuthPublicPaths executes the corresponding workflow in this package.
func defaultGlobalAuthPublicPaths(metricsPath, metricsReportPath string) []string {
	paths := []string{
		"/health",
		"/ready",
		"/system/report",
		"/system/report.json",
		"/metrics",
		"/metrics/report",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/refresh",
		"/api/v1/auth/oauth/*",
		"/api/v1/payment/callback/*",
		"/api/v1/qrcode/s/*",
	}
	paths = appendUniqueString(paths, strings.TrimSpace(metricsPath))
	paths = appendUniqueString(paths, strings.TrimSpace(metricsReportPath))
	return paths
}

// 中文：appendUniqueString 执行当前包中的对应流程。
// English: appendUniqueString executes the corresponding workflow in this package.
func appendUniqueString(values []string, value string) []string {
	if value == "" {
		return values
	}
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

// 中文：installMetrics 执行当前包中的对应流程。
// English: installMetrics executes the corresponding workflow in this package.
func (a *App) installMetrics() error {
	var cfg config.MetricsConfig
	if err := a.config.Unmarshal("metrics", &cfg); err != nil {
		return fmt.Errorf("parse metrics config: %w", err)
	}
	if !cfg.Enabled {
		return nil
	}
	if cfg.Path == "" {
		cfg.Path = "/metrics"
	}
	if cfg.ReportPath == "" {
		cfg.ReportPath = "/metrics/report"
	}

	registry := metricspkg.NewRegistry(cfg.Namespace)
	a.metrics = registry
	a.di.Provide(registry)
	a.server.Engine().Use(middleware.Metrics(registry))
	a.server.Engine().GET(cfg.Path, func(c *gin.Context) {
		c.Header("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		if err := registry.WritePrometheus(c.Writer); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
	})
	a.server.Engine().GET(cfg.ReportPath, func(c *gin.Context) {
		c.Header("Content-Type", "text/markdown; charset=utf-8")
		if err := registry.WriteReport(c.Writer); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
	})
	return nil
}

// 中文：registerSystemRoutes 执行当前包中的对应流程。
// English: registerSystemRoutes executes the corresponding workflow in this package.
func (a *App) registerSystemRoutes() {
	a.server.Engine().GET("/system/report", func(c *gin.Context) {
		c.Header("Content-Type", "text/markdown; charset=utf-8")
		if err := a.WriteRuntimeReport(c.Writer); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
	})
	a.server.Engine().GET("/system/report.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, a.RuntimeReport())
	})
}

// 中文：RuntimeReport 执行当前包中的对应流程。
// English: RuntimeReport executes the corresponding workflow in this package.
// RuntimeReport returns a structured snapshot of the app, modules, and infrastructure.
func (a *App) RuntimeReport() RuntimeReport {
	appName := "spiringo"
	env := a.env
	var serverSnapshot ServerSnapshot
	if a.config != nil {
		appName = firstNonEmpty(a.config.GetString("app.name"), appName)
		env = firstNonEmpty(env, a.config.GetString("app.env"))
		serverSnapshot = ServerSnapshot{
			ListenAddr: a.config.GetString("server.addr"),
			PublicURL:  a.config.GetString("server.public_url"),
			APIBaseURL: a.config.GetString("server.api_base_url"),
		}
	}
	if env == "" {
		env = "unknown"
	}

	var modules []module.ModuleSnapshot
	if a.registry != nil {
		modules = a.registry.Snapshots()
	}

	report := RuntimeReport{
		AppName: appName,
		Env:     env,
		Server:  serverSnapshot,
		Modules: modules,
		Infrastructure: InfrastructureSnapshot{
			Database: a.db != nil,
			Document: a.document != nil,
			Cache:    a.cache != nil,
			Lock:     a.lock != nil,
			Metrics:  a.metrics != nil,
			Alerts:   a.alerts != nil,
			MQ:       a.mq != nil,
			Queue:    a.queue != nil,
			Search:   a.search != nil,
			Storage:  a.storage != nil,
			Trace:    a.tracer != nil,
		},
	}
	if a.metrics != nil {
		snapshot := a.metrics.Snapshot()
		report.Metrics = &snapshot
	}
	return report
}

// 中文：WriteRuntimeReport 执行当前包中的对应流程。
// English: WriteRuntimeReport executes the corresponding workflow in this package.
// WriteRuntimeReport writes a Markdown version of RuntimeReport.
func (a *App) WriteRuntimeReport(w io.Writer) error {
	report := a.RuntimeReport()
	if _, err := fmt.Fprintf(w, "# Spiringo Runtime Report\n\nApp: `%s`\n\nEnv: `%s`\n\n", report.AppName, report.Env); err != nil {
		return err
	}
	if report.Server.ListenAddr != "" || report.Server.PublicURL != "" || report.Server.APIBaseURL != "" {
		if _, err := fmt.Fprintf(w, "Listen: `%s`\n\nPublic URL: `%s`\n\nAPI Base URL: `%s`\n\n", report.Server.ListenAddr, report.Server.PublicURL, report.Server.APIBaseURL); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(w, "## Infrastructure"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| Component | Status |"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| --- | --- |"); err != nil {
		return err
	}
	for _, row := range []struct {
		name   string
		active bool
	}{
		{"database", report.Infrastructure.Database},
		{"document", report.Infrastructure.Document},
		{"cache", report.Infrastructure.Cache},
		{"lock", report.Infrastructure.Lock},
		{"metrics", report.Infrastructure.Metrics},
		{"alerts", report.Infrastructure.Alerts},
		{"mq", report.Infrastructure.MQ},
		{"queue", report.Infrastructure.Queue},
		{"search", report.Infrastructure.Search},
		{"storage", report.Infrastructure.Storage},
		{"trace", report.Infrastructure.Trace},
	} {
		if _, err := fmt.Fprintf(w, "| `%s` | %s |\n", row.name, activeText(row.active)); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(w, "\n## Modules"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	if len(report.Modules) == 0 {
		if _, err := fmt.Fprintln(w, "No modules registered."); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintln(w, "| Name | State | Enabled | Dependencies |"); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "| --- | --- | --- | --- |"); err != nil {
			return err
		}
		for _, m := range report.Modules {
			deps := "-"
			if len(m.Dependencies) > 0 {
				deps = strings.Join(m.Dependencies, ", ")
			}
			if _, err := fmt.Fprintf(w, "| `%s` | `%s` | %s | %s |\n", m.Name, m.State, moduleStatusText(m), deps); err != nil {
				return err
			}
		}
	}

	if _, err := fmt.Fprintln(w, "\n## Metrics"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	if report.Metrics == nil {
		_, err := fmt.Fprintln(w, "Metrics disabled.")
		return err
	}
	_, err := fmt.Fprintf(w, "Namespace: `%s`\n\nCounters: `%d`\n\nSummaries: `%d`\n", report.Metrics.Namespace, len(report.Metrics.Counters), len(report.Metrics.Summaries))
	return err
}

// 中文：activeText 执行当前包中的对应流程。
// English: activeText executes the corresponding workflow in this package.
func activeText(active bool) string {
	if active {
		return "`active`"
	}
	return "`inactive`"
}

// 中文：moduleStatusText 执行当前包中的对应流程。
// English: moduleStatusText executes the corresponding workflow in this package.
func moduleStatusText(snapshot module.ModuleSnapshot) string {
	if snapshot.Skipped {
		return "`skipped`"
	}
	if snapshot.Enabled {
		return "`enabled`"
	}
	return "`disabled`"
}

// 中文：registerHealthRoutes 执行当前包中的对应流程。
// English: registerHealthRoutes executes the corresponding workflow in this package.
func (a *App) registerHealthRoutes() {
	a.server.Engine().GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	a.server.Engine().GET("/ready", func(c *gin.Context) {
		if a.db != nil {
			if err := a.db.Ping(); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "database": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})
}

// 中文：closeInfrastructure 执行当前包中的对应流程。
// English: closeInfrastructure executes the corresponding workflow in this package.
func (a *App) closeInfrastructure() {
	if a.document != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		if err := a.document.Close(ctx); err != nil {
			a.logger.Warn("close document database failed", "error", err)
		}
		cancel()
	}
	if a.queue != nil {
		if err := a.queue.Close(); err != nil {
			a.logger.Warn("close task queue failed", "error", err)
		}
	}
	if a.lock != nil {
		if closer, ok := a.lock.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				a.logger.Warn("close lock failed", "error", err)
			}
		}
	}
	if a.lockRedis != nil {
		if err := a.lockRedis.Close(); err != nil {
			a.logger.Warn("close lock redis failed", "error", err)
		}
	}
	if a.mq != nil {
		if err := a.mq.Close(); err != nil {
			a.logger.Warn("close mq failed", "error", err)
		}
	}
	if a.mqRedis != nil {
		if err := a.mqRedis.Close(); err != nil {
			a.logger.Warn("close mq redis failed", "error", err)
		}
	}
	if a.cache != nil {
		if err := a.cache.Close(); err != nil {
			a.logger.Warn("close cache failed", "error", err)
		}
	}
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			a.logger.Warn("close database failed", "error", err)
		}
	}
	if a.alerts != nil {
		if closer, ok := a.alerts.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				a.logger.Warn("close alerts failed", "error", err)
			}
		}
	}
	if a.loggerSync != nil {
		if err := a.loggerSync(); err != nil {
			a.logger.Debug("sync logger failed", "error", err)
		}
	}
}

// 中文：waitStop 执行当前包中的对应流程。
// English: waitStop executes the corresponding workflow in this package.
func (a *App) waitStop() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case <-quit:
		return nil
	case err, ok := <-a.server.Errors():
		if ok && err != nil {
			return err
		}
		return nil
	}
}

// 中文：parseOptionalDuration 执行当前包中的对应流程。
// English: parseOptionalDuration executes the corresponding workflow in this package.
func parseOptionalDuration(value string) (time.Duration, error) {
	if strings.TrimSpace(value) == "" {
		return 0, nil
	}
	return time.ParseDuration(value)
}

// 中文：parseSlogLevel 执行当前包中的对应流程。
// English: parseSlogLevel executes the corresponding workflow in this package.
func parseSlogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
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

// 中文：firstPositiveDuration 执行当前包中的对应流程。
// English: firstPositiveDuration executes the corresponding workflow in this package.
func firstPositiveDuration(values ...time.Duration) time.Duration {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}

// 中文：splitCSV 执行当前包中的对应流程。
// English: splitCSV executes the corresponding workflow in this package.
func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if item := strings.TrimSpace(part); item != "" {
			result = append(result, item)
		}
	}
	return result
}
