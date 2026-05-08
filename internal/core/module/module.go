package module

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/core/config"
	"github.com/spiringo/spiringo/internal/core/di"
	"github.com/spiringo/spiringo/internal/core/event"
)

// 中文：ModuleState 定义当前包使用的数据结构或接口。
// English: ModuleState defines a data structure or interface used by this package.
// ModuleState 模块状态
type ModuleState int

// 中文：ModuleStateInactive、ModuleStateInitializing、ModuleStateActive、... 声明当前包使用的常量。
// English: ModuleStateInactive、ModuleStateInitializing、ModuleStateActive、... declares constants used by this package.
const (
	ModuleStateInactive ModuleState = iota
	ModuleStateInitializing
	ModuleStateActive
	ModuleStateStopping
	ModuleStateStopped
)

// 中文：String 执行当前包中的对应流程。
// English: String executes the corresponding workflow in this package.
func (s ModuleState) String() string {
	switch s {
	case ModuleStateInactive:
		return "inactive"
	case ModuleStateInitializing:
		return "initializing"
	case ModuleStateActive:
		return "active"
	case ModuleStateStopping:
		return "stopping"
	case ModuleStateStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

// 中文：Module 定义当前包使用的数据结构或接口。
// English: Module defines a data structure or interface used by this package.
// Module 模块接口 —— 所有业务模块必须实现
type Module interface {
	// 中文：Name 声明该接口需要实现的行为。
	// English: Name declares behavior required by this interface.
	// Name 模块唯一标识
	Name() string
	// 中文：Dependencies 声明该接口需要实现的行为。
	// English: Dependencies declares behavior required by this interface.
	// Dependencies 声明依赖的模块名列表
	Dependencies() []string
	// 中文：Config 声明该接口需要实现的行为。
	// English: Config declares behavior required by this interface.
	// Config 返回模块的配置结构体指针
	Config() any
	// 中文：Init 声明该接口需要实现的行为。
	// English: Init declares behavior required by this interface.
	// Init 初始化模块
	Init(app *App) error
	// 中文：Start 声明该接口需要实现的行为。
	// English: Start declares behavior required by this interface.
	// Start 启动模块
	Start(ctx context.Context) error
	// 中文：Stop 声明该接口需要实现的行为。
	// English: Stop declares behavior required by this interface.
	// Stop 优雅停止模块
	Stop(ctx context.Context) error
	// 中文：State 声明该接口需要实现的行为。
	// English: State declares behavior required by this interface.
	// State 返回当前模块状态
	State() ModuleState
}

// 中文：BaseModule 定义当前包使用的数据结构或接口。
// English: BaseModule defines a data structure or interface used by this package.
// BaseModule 模块基类（减少样板代码）
type BaseModule struct {
	// 中文：name 保存当前结构中的配置或数据值。
	// English: name stores a configuration or data value for this struct.
	name string
	// 中文：deps 保存当前结构中的配置或数据值。
	// English: deps stores a configuration or data value for this struct.
	deps []string
	// 中文：state 保存当前结构中的配置或数据值。
	// English: state stores a configuration or data value for this struct.
	state ModuleState
	// 中文：stateMu 保存当前结构中的配置或数据值。
	// English: stateMu stores a configuration or data value for this struct.
	stateMu sync.RWMutex
}

// 中文：stateSetter 定义当前包使用的数据结构或接口。
// English: stateSetter defines a data structure or interface used by this package.
type stateSetter interface {
	// 中文：SetState 声明该接口需要实现的行为。
	// English: SetState declares behavior required by this interface.
	SetState(ModuleState)
}

// 中文：ModuleSnapshot 定义当前包使用的数据结构或接口。
// English: ModuleSnapshot defines a data structure or interface used by this package.
// NewBaseModule 创建模块基类
// ModuleSnapshot is a stable runtime view of one registered module.
type ModuleSnapshot struct {
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `json:"name"`
	// 中文：State 保存当前结构中的配置或数据值。
	// English: State stores a configuration or data value for this struct.
	State string `json:"state"`
	// 中文：Dependencies 保存当前结构中的配置或数据值。
	// English: Dependencies stores a configuration or data value for this struct.
	Dependencies []string `json:"dependencies,omitempty"`
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `json:"enabled"`
	// 中文：Skipped 保存当前结构中的配置或数据值。
	// English: Skipped stores a configuration or data value for this struct.
	Skipped bool `json:"skipped"`
}

// 中文：NewBaseModule 创建并返回对应组件实例。
// English: NewBaseModule creates and returns the corresponding component instance.
// NewBaseModule creates a base module.
func NewBaseModule(name string, deps ...string) *BaseModule {
	return &BaseModule{
		name:  name,
		deps:  deps,
		state: ModuleStateInactive,
	}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (m *BaseModule) Name() string { return m.name }

// 中文：Dependencies 执行当前包中的对应流程。
// English: Dependencies executes the corresponding workflow in this package.
func (m *BaseModule) Dependencies() []string { return m.deps }

// 中文：Config 执行当前包中的对应流程。
// English: Config executes the corresponding workflow in this package.
func (m *BaseModule) Config() any { return nil }

// 中文：Init 执行当前包中的对应流程。
// English: Init executes the corresponding workflow in this package.
func (m *BaseModule) Init(_ *App) error { return nil }

// 中文：Start 执行当前包中的对应流程。
// English: Start executes the corresponding workflow in this package.
func (m *BaseModule) Start(ctx context.Context) error {
	if ctx != nil {
		return ctx.Err()
	}
	return nil
}

// 中文：Stop 执行当前包中的对应流程。
// English: Stop executes the corresponding workflow in this package.
func (m *BaseModule) Stop(ctx context.Context) error {
	if ctx != nil {
		return ctx.Err()
	}
	return nil
}

// 中文：State 执行当前包中的对应流程。
// English: State executes the corresponding workflow in this package.
func (m *BaseModule) State() ModuleState {
	m.stateMu.RLock()
	defer m.stateMu.RUnlock()
	return m.state
}

// 中文：SetState 执行当前包中的对应流程。
// English: SetState executes the corresponding workflow in this package.
func (m *BaseModule) SetState(s ModuleState) {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	m.state = s
}

// 中文：Routable 定义当前包使用的数据结构或接口。
// English: Routable defines a data structure or interface used by this package.
// Routable 可注册路由的模块接口
type Routable interface {
	// 中文：Routes 声明该接口需要实现的行为。
	// English: Routes declares behavior required by this interface.
	Routes(r *gin.RouterGroup)
}

// 中文：Migratable 定义当前包使用的数据结构或接口。
// English: Migratable defines a data structure or interface used by this package.
// Migratable 可执行数据库迁移的模块接口
type Migratable interface {
	// 中文：Migrations 声明该接口需要实现的行为。
	// English: Migrations declares behavior required by this interface.
	Migrations() []Migration
}

// 中文：MigrationStore 定义当前包使用的数据结构或接口。
// English: MigrationStore defines a data structure or interface used by this package.
// MigrationStore records applied migrations and skips already-applied IDs.
type MigrationStore interface {
	// 中文：RunMigrations 声明该接口需要实现的行为。
	// English: RunMigrations declares behavior required by this interface.
	RunMigrations(ctx context.Context, moduleName string, migrations []Migration) error
}

// 中文：EventSubscriber 定义当前包使用的数据结构或接口。
// English: EventSubscriber defines a data structure or interface used by this package.
// EventSubscriber 可订阅事件的模块接口
type EventSubscriber interface {
	// 中文：Subscriptions 声明该接口需要实现的行为。
	// English: Subscriptions declares behavior required by this interface.
	Subscriptions() []EventSubscription
}

// 中文：EventSubscription 定义当前包使用的数据结构或接口。
// English: EventSubscription defines a data structure or interface used by this package.
// EventSubscription 事件订阅
type EventSubscription struct {
	// 中文：Topic 保存当前结构中的配置或数据值。
	// English: Topic stores a configuration or data value for this struct.
	Topic string
	// 中文：Handler 保存当前结构中的配置或数据值。
	// English: Handler stores a configuration or data value for this struct.
	Handler event.Handler
}

// 中文：Migration 定义当前包使用的数据结构或接口。
// English: Migration defines a data structure or interface used by this package.
// Migration 数据库迁移
type Migration struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string
	// 中文：Up 保存当前结构中的配置或数据值。
	// English: Up stores a configuration or data value for this struct.
	Up func(ctx context.Context) error
	// 中文：Down 保存当前结构中的配置或数据值。
	// English: Down stores a configuration or data value for this struct.
	Down func(ctx context.Context) error
}

// 中文：App 定义当前包使用的数据结构或接口。
// English: App defines a data structure or interface used by this package.
// App 应用上下文（传递给模块Init方法）
type App struct {
	// 中文：Config 保存当前结构中的配置或数据值。
	// English: Config stores a configuration or data value for this struct.
	Config *config.Manager
	// 中文：DI 保存当前结构中的配置或数据值。
	// English: DI stores a configuration or data value for this struct.
	DI *di.Container
	// 中文：EventBus 保存当前结构中的配置或数据值。
	// English: EventBus stores a configuration or data value for this struct.
	EventBus *event.Bus
	// 中文：Router 保存当前结构中的配置或数据值。
	// English: Router stores a configuration or data value for this struct.
	Router *gin.Engine
	// 中文：Modules 保存当前结构中的配置或数据值。
	// English: Modules stores a configuration or data value for this struct.
	Modules *Registry
	// 中文：Migrate 保存当前结构中的配置或数据值。
	// English: Migrate stores a configuration or data value for this struct.
	Migrate MigrationStore
}

// 中文：Registry 定义当前包使用的数据结构或接口。
// English: Registry defines a data structure or interface used by this package.
// Registry 模块注册中心
type Registry struct {
	// 中文：modules 保存当前结构中的配置或数据值。
	// English: modules stores a configuration or data value for this struct.
	modules map[string]Module
	// 中文：initOrder 保存当前结构中的配置或数据值。
	// English: initOrder stores a configuration or data value for this struct.
	initOrder []string
	// 中文：skippedModules 保存当前结构中的配置或数据值。
	// English: skippedModules stores a configuration or data value for this struct.
	skippedModules map[string]bool // 在InitAll中被跳过的模块（未启用）
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.RWMutex
}

// 中文：NewRegistry 创建并返回对应组件实例。
// English: NewRegistry creates and returns the corresponding component instance.
// NewRegistry 创建模块注册中心
func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

// 中文：Register 执行当前包中的对应流程。
// English: Register executes the corresponding workflow in this package.
// Register 注册模块
func (r *Registry) Register(m Module) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := m.Name()
	if _, exists := r.modules[name]; exists {
		return fmt.Errorf("module %q already registered", name)
	}
	r.modules[name] = m
	return nil
}

// 中文：MustRegister 执行当前包中的对应流程。
// English: MustRegister executes the corresponding workflow in this package.
// MustRegister 注册模块，重复则panic
func (r *Registry) MustRegister(m Module) {
	if err := r.Register(m); err != nil {
		panic(err)
	}
}

// 中文：Get 执行当前包中的对应流程。
// English: Get executes the corresponding workflow in this package.
// Get 按名获取模块
func (r *Registry) Get(name string) (Module, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, ok := r.modules[name]
	if !ok {
		return nil, fmt.Errorf("module %q not found", name)
	}
	return m, nil
}

// 中文：ResolveOrder 执行当前包中的对应流程。
// English: ResolveOrder executes the corresponding workflow in this package.
// ResolveOrder 拓扑排序解决依赖
func (r *Registry) ResolveOrder() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	order := make([]string, 0)

	var visit func(name string) error
	visit = func(name string) error {
		if visited[name] {
			return nil
		}
		if visiting[name] {
			return fmt.Errorf("circular dependency detected at module %q", name)
		}

		m, ok := r.modules[name]
		if !ok {
			return fmt.Errorf("dependency %q not found (required by another module)", name)
		}

		visiting[name] = true
		for _, dep := range m.Dependencies() {
			if err := visit(dep); err != nil {
				return err
			}
		}
		visiting[name] = false
		visited[name] = true
		order = append(order, name)
		return nil
	}

	for name := range r.modules {
		if err := visit(name); err != nil {
			return err
		}
	}

	r.initOrder = order
	return nil
}

// 中文：InitAll 执行当前包中的对应流程。
// English: InitAll executes the corresponding workflow in this package.
// InitAll 按依赖顺序初始化所有已启用的模块
func (r *Registry) InitAll(app *App) error {
	if err := r.ResolveOrder(); err != nil {
		return err
	}

	r.skippedModules = make(map[string]bool)

	for _, name := range r.initOrder {
		m := r.modules[name]

		// 检查模块是否启用
		enabledKey := "modules." + name + ".enabled"
		if app.Config.IsSet(enabledKey) && !app.Config.GetBool(enabledKey) {
			r.skippedModules[name] = true
			continue
		}
		if skippedDep := r.firstSkippedDependency(m); skippedDep != "" {
			r.skippedModules[name] = true
			slog.Warn("module skipped because dependency is disabled", "module", name, "dependency", skippedDep)
			continue
		}

		// 绑定模块专属配置
		if cfg := m.Config(); cfg != nil {
			if err := app.Config.Unmarshal("modules."+name, cfg); err != nil {
				return fmt.Errorf("bind config for module %s: %w", name, err)
			}
		}

		// 更新状态
		if setter, ok := m.(stateSetter); ok {
			setter.SetState(ModuleStateInitializing)
		}

		// 初始化
		if err := m.Init(app); err != nil {
			return fmt.Errorf("init module %s: %w", name, err)
		}

		// 注册路由
		if routable, ok := m.(Routable); ok {
			group := app.Router.Group("/api/v1/" + name)
			routable.Routes(group)
		}

		// 注册事件订阅
		if subscriber, ok := m.(EventSubscriber); ok {
			for _, sub := range subscriber.Subscriptions() {
				if err := app.EventBus.Subscribe(sub.Topic, sub.Handler); err != nil {
					return fmt.Errorf("subscribe module %s event %s: %w", name, sub.Topic, err)
				}
			}
		}

		// 执行数据库迁移
		if migratable, ok := m.(Migratable); ok {
			if err := runModuleMigrations(context.Background(), app, name, migratable.Migrations()); err != nil {
				return fmt.Errorf("run migrations for module %s: %w", name, err)
			}
		}

		if setter, ok := m.(stateSetter); ok {
			setter.SetState(ModuleStateActive)
		}
	}

	return nil
}

// 中文：StartAll 执行当前包中的对应流程。
// English: StartAll executes the corresponding workflow in this package.
// StartAll 启动所有模块
func (r *Registry) StartAll(ctx context.Context) error {
	for _, name := range r.initOrder {
		// 跳过未启用模块
		if r.skippedModules != nil && r.skippedModules[name] {
			continue
		}

		m := r.modules[name]

		if err := m.Start(ctx); err != nil {
			return fmt.Errorf("start module %s: %w", name, err)
		}
	}
	return nil
}

// 中文：StopAll 执行当前包中的对应流程。
// English: StopAll executes the corresponding workflow in this package.
// StopAll 逆序停止所有模块
func (r *Registry) StopAll(ctx context.Context) error {
	for i := len(r.initOrder) - 1; i >= 0; i-- {
		name := r.initOrder[i]
		if r.skippedModules != nil && r.skippedModules[name] {
			continue
		}
		m := r.modules[name]

		if setter, ok := m.(stateSetter); ok {
			setter.SetState(ModuleStateStopping)
		}

		if err := m.Stop(ctx); err != nil {
			// 停止失败只记录，不中断其他模块
			slog.Warn("stop module failed", "module", name, "error", err)
		}

		if setter, ok := m.(stateSetter); ok {
			setter.SetState(ModuleStateStopped)
		}
	}
	return nil
}

// 中文：firstSkippedDependency 执行当前包中的对应流程。
// English: firstSkippedDependency executes the corresponding workflow in this package.
func (r *Registry) firstSkippedDependency(m Module) string {
	if r == nil || m == nil || r.skippedModules == nil {
		return ""
	}
	for _, dep := range m.Dependencies() {
		if r.skippedModules[dep] {
			return dep
		}
	}
	return ""
}

// 中文：ListModules 执行当前包中的对应流程。
// English: ListModules executes the corresponding workflow in this package.
// ListModules 列出所有已注册模块
func (r *Registry) ListModules() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// 中文：Snapshots 执行当前包中的对应流程。
// English: Snapshots executes the corresponding workflow in this package.
// RegisterModules 批量注册模块
// Snapshots returns deterministic module lifecycle snapshots.
func (r *Registry) Snapshots() []ModuleSnapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := append([]string(nil), r.initOrder...)
	seen := make(map[string]struct{}, len(names))
	for _, name := range names {
		seen[name] = struct{}{}
	}

	remaining := make([]string, 0, len(r.modules))
	for name := range r.modules {
		if _, ok := seen[name]; !ok {
			remaining = append(remaining, name)
		}
	}
	sort.Strings(remaining)
	names = append(names, remaining...)

	snapshots := make([]ModuleSnapshot, 0, len(names))
	for _, name := range names {
		m, ok := r.modules[name]
		if !ok {
			continue
		}
		skipped := r.skippedModules != nil && r.skippedModules[name]
		snapshots = append(snapshots, ModuleSnapshot{
			Name:         name,
			State:        m.State().String(),
			Dependencies: append([]string(nil), m.Dependencies()...),
			Enabled:      !skipped,
			Skipped:      skipped,
		})
	}
	return snapshots
}

// 中文：RegisterModules 执行当前包中的对应流程。
// English: RegisterModules executes the corresponding workflow in this package.
func (r *Registry) RegisterModules(modules ...Module) {
	for _, m := range modules {
		r.MustRegister(m)
	}
}

// 中文：runModuleMigrations 执行当前包中的对应流程。
// English: runModuleMigrations executes the corresponding workflow in this package.
func runModuleMigrations(ctx context.Context, app *App, moduleName string, migrations []Migration) error {
	if app != nil && app.Migrate != nil {
		return app.Migrate.RunMigrations(ctx, moduleName, migrations)
	}
	for _, mg := range migrations {
		if mg.Up != nil {
			if err := mg.Up(ctx); err != nil {
				return fmt.Errorf("run migration %s: %w", mg.ID, err)
			}
		}
	}
	return nil
}
