//go:build !(windows && 386)

package app

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/spiringo/spiringo/internal/core/di"
	"github.com/spiringo/spiringo/internal/core/module"
	cachepkg "github.com/spiringo/spiringo/internal/pkg/cache"
	metricspkg "github.com/spiringo/spiringo/internal/pkg/metrics"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：TestAppInitInfrastructureRegistersDatabaseAndCache 验证相关行为符合预期。
// English: TestAppInitInfrastructureRegistersDatabaseAndCache verifies the related behavior.
func TestAppInitInfrastructureRegistersDatabaseAndCache(t *testing.T) {
	if runtime.GOOS == "windows" && runtime.GOARCH == "386" {
		t.Skip("sqlite cgo cannot link with the local 64-bit Windows gcc while GOARCH=386")
	}

	dir := t.TempDir()
	configData := []byte(`
log:
  level: "error"
  format: "text"
database:
  default:
    driver: "sqlite"
    dsn: "file:app_test?mode=memory&cache=shared"
    max_idle: 1
    max_open: 1
    conn_max_lifetime: "0s"
cache:
  driver: "memory"
storage:
  enabled: false
`)
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), configData, 0o600); err != nil {
		t.Fatal(err)
	}

	a := New(WithConfigDir(dir), WithEnv("test"))
	a.config.SetConfigDir(dir)
	a.config.SetEnv("test")
	if err := a.config.Load(); err != nil {
		t.Fatal(err)
	}
	if err := a.initLogger(); err != nil {
		t.Fatal(err)
	}
	if err := a.initInfrastructure(); err != nil {
		t.Fatal(err)
	}
	defer a.closeInfrastructure()

	if _, err := di.Resolve[*orm.DB](a.di); err != nil {
		t.Fatalf("expected database to be registered: %v", err)
	}
	if _, err := di.Resolve[cachepkg.Cache](a.di); err != nil {
		t.Fatalf("expected cache interface to be registered: %v", err)
	}
}

// 中文：TestAppMigrateRunsRegisteredModuleMigrations 验证相关行为符合预期。
// English: TestAppMigrateRunsRegisteredModuleMigrations verifies the related behavior.
func TestAppMigrateRunsRegisteredModuleMigrations(t *testing.T) {
	if runtime.GOOS == "windows" && runtime.GOARCH == "386" {
		t.Skip("sqlite cgo cannot link with the local 64-bit Windows gcc while GOARCH=386")
	}

	dir := t.TempDir()
	configData := []byte(`
log:
  level: "error"
  format: "text"
database:
  default:
    driver: "sqlite"
    dsn: "` + filepath.ToSlash(filepath.Join(dir, "migrate.db")) + `"
    max_idle: 1
    max_open: 1
cache:
  driver: "memory"
modules:
  demo:
    enabled: true
`)
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), configData, 0o600); err != nil {
		t.Fatal(err)
	}

	var migrated atomic.Int32
	a := New(WithConfigDir(dir), WithEnv("test"))
	a.RegisterModules(&migrateTestModule{BaseModule: module.NewBaseModule("demo"), migrated: &migrated})
	if err := a.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}

	if got := migrated.Load(); got != 1 {
		t.Fatalf("migration count = %d", got)
	}
}

// 中文：TestAppRuntimeReportIncludesModulesAndMetrics 验证相关行为符合预期。
// English: TestAppRuntimeReportIncludesModulesAndMetrics verifies the related behavior.
func TestAppRuntimeReportIncludesModulesAndMetrics(t *testing.T) {
	a := New(WithEnv("test"))
	a.config.Set("app.name", "demo")

	m := module.NewBaseModule("tenant")
	m.SetState(module.ModuleStateActive)
	a.registry.MustRegister(m)
	if err := a.registry.ResolveOrder(); err != nil {
		t.Fatal(err)
	}

	a.metrics = metricspkg.NewRegistry("demo")
	a.metrics.IncCounter("requests_total", nil)

	report := a.RuntimeReport()
	if report.AppName != "demo" || report.Env != "test" {
		t.Fatalf("unexpected app identity: %+v", report)
	}
	if len(report.Modules) != 1 || report.Modules[0].Name != "tenant" || report.Modules[0].State != module.ModuleStateActive.String() {
		t.Fatalf("unexpected modules: %+v", report.Modules)
	}
	if !report.Infrastructure.Metrics || report.Metrics == nil || len(report.Metrics.Counters) != 1 {
		t.Fatalf("expected metrics snapshot, got infrastructure=%+v metrics=%+v", report.Infrastructure, report.Metrics)
	}

	var b strings.Builder
	if err := a.WriteRuntimeReport(&b); err != nil {
		t.Fatal(err)
	}
	output := b.String()
	for _, want := range []string{"# Spiringo Runtime Report", "| `tenant` | `active` | `enabled` | - |", "Counters: `1`"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected report to contain %q, got:\n%s", want, output)
		}
	}
}

// 中文：TestMiddlewareConfigBackfillsAuthSecretAndPublicPaths 验证相关行为符合预期。
// English: TestMiddlewareConfigBackfillsAuthSecretAndPublicPaths verifies the related behavior.
func TestMiddlewareConfigBackfillsAuthSecretAndPublicPaths(t *testing.T) {
	a := New(WithEnv("test"))
	a.config.Set("middleware.auth.enabled", true)
	a.config.Set("modules.auth.jwt.secret", "secret-from-auth-module")
	a.config.Set("metrics.path", "/internal/metrics")

	cfg, err := a.middlewareConfig()
	if err != nil {
		t.Fatal(err)
	}

	if !cfg.Auth.Enabled {
		t.Fatal("expected global auth to be enabled")
	}
	if cfg.Auth.JWTSecret != "secret-from-auth-module" {
		t.Fatalf("jwt secret = %q, want fallback secret", cfg.Auth.JWTSecret)
	}
	for _, want := range []string{"/api/v1/auth/login", "/api/v1/payment/callback/*", "/internal/metrics"} {
		if !containsString(cfg.Auth.PublicPaths, want) {
			t.Fatalf("public paths missing %q: %#v", want, cfg.Auth.PublicPaths)
		}
	}
}

// 中文：containsString 执行当前包中的对应流程。
// English: containsString executes the corresponding workflow in this package.
func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

// 中文：migrateTestModule 定义当前包使用的数据结构或接口。
// English: migrateTestModule defines a data structure or interface used by this package.
type migrateTestModule struct {
	// 中文：*module.BaseModule 嵌入复用该类型提供的能力。
	// English: *module.BaseModule embeds reusable behavior from that type.
	*module.BaseModule
	// 中文：migrated 保存当前结构中的配置或数据值。
	// English: migrated stores a configuration or data value for this struct.
	migrated *atomic.Int32
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (m *migrateTestModule) Name() string {
	if m.BaseModule == nil {
		m.BaseModule = module.NewBaseModule("demo")
	}
	return m.BaseModule.Name()
}

// 中文：Dependencies 执行当前包中的对应流程。
// English: Dependencies executes the corresponding workflow in this package.
func (m *migrateTestModule) Dependencies() []string { return nil }

// 中文：Config 执行当前包中的对应流程。
// English: Config executes the corresponding workflow in this package.
func (m *migrateTestModule) Config() any { return nil }

// 中文：Init 执行当前包中的对应流程。
// English: Init executes the corresponding workflow in this package.
func (m *migrateTestModule) Init(_ *module.App) error {
	if m.BaseModule == nil {
		m.BaseModule = module.NewBaseModule("demo")
	}
	return nil
}

// 中文：Start 执行当前包中的对应流程。
// English: Start executes the corresponding workflow in this package.
func (m *migrateTestModule) Start(ctx context.Context) error { return ctx.Err() }

// 中文：Stop 执行当前包中的对应流程。
// English: Stop executes the corresponding workflow in this package.
func (m *migrateTestModule) Stop(ctx context.Context) error { return ctx.Err() }

// 中文：State 执行当前包中的对应流程。
// English: State executes the corresponding workflow in this package.
func (m *migrateTestModule) State() module.ModuleState {
	if m.BaseModule == nil {
		m.BaseModule = module.NewBaseModule("demo")
	}
	return m.BaseModule.State()
}

// 中文：Migrations 执行当前包中的对应流程。
// English: Migrations executes the corresponding workflow in this package.
func (m *migrateTestModule) Migrations() []module.Migration {
	return []module.Migration{
		{
			ID: "demo_001",
			Up: func(ctx context.Context) error {
				m.migrated.Add(1)
				return ctx.Err()
			},
		},
	}
}
