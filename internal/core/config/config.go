package config

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// 中文：Manager 定义当前包使用的数据结构或接口。
// English: Manager defines a data structure or interface used by this package.
// Manager 配置管理器
type Manager struct {
	// 中文：viper 保存当前结构中的配置或数据值。
	// English: viper stores a configuration or data value for this struct.
	viper *viper.Viper
	// 中文：sources 保存当前结构中的配置或数据值。
	// English: sources stores a configuration or data value for this struct.
	sources []Source
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.RWMutex
	// 中文：env 保存当前结构中的配置或数据值。
	// English: env stores a configuration or data value for this struct.
	env string
	// 中文：configDir 保存当前结构中的配置或数据值。
	// English: configDir stores a configuration or data value for this struct.
	configDir string
	// 中文：onChange 保存当前结构中的配置或数据值。
	// English: onChange stores a configuration or data value for this struct.
	onChange map[string][]func(any)
}

// 中文：NewManager 创建并返回对应组件实例。
// English: NewManager creates and returns the corresponding component instance.
// NewManager 创建配置管理器
func NewManager() *Manager {
	v := viper.New()
	v.SetEnvPrefix("SP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetTypeByDefaultValue(true)

	return &Manager{
		viper:    v,
		sources:  make([]Source, 0),
		onChange: make(map[string][]func(any)),
	}
}

// 中文：SetConfigDir 执行当前包中的对应流程。
// English: SetConfigDir executes the corresponding workflow in this package.
// SetConfigDir 设置配置文件目录
func (m *Manager) SetConfigDir(dir string) {
	m.configDir = dir
}

// 中文：SetEnv 执行当前包中的对应流程。
// English: SetEnv executes the corresponding workflow in this package.
// SetEnv 设置环境
func (m *Manager) SetEnv(env string) {
	m.env = NormalizeEnv(env)
}

// 中文：AddSource 执行当前包中的对应流程。
// English: AddSource executes the corresponding workflow in this package.
// AddSource 添加配置源
func (m *Manager) AddSource(source Source) {
	m.sources = append(m.sources, source)
}

// 中文：Load 执行当前包中的对应流程。
// English: Load executes the corresponding workflow in this package.
// Load 加载所有配置
// 优先级：环境变量 > 配置中心 > 当前环境段 > 默认配置文件 > 代码默认值
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 1. 加载默认配置文件
	m.viper.SetConfigName("config")
	m.viper.SetConfigType("yaml")
	if m.configDir != "" {
		m.viper.AddConfigPath(m.configDir)
	}
	if err := m.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("read config: %w", err)
		}
	}

	if m.env == "" {
		m.env = NormalizeEnv(m.viper.GetString("app.env"))
	}
	if m.env == "" {
		m.env = "local"
	}
	m.applyEnvironmentProfileLocked(m.env)
	m.viper.Set("app.env", m.env)

	// 2. Load pluggable sources from low to high priority so higher sources win.
	sources := append([]Source(nil), m.sources...)
	sort.SliceStable(sources, func(i, j int) bool {
		return sources[i].Priority() < sources[j].Priority()
	})
	for _, src := range sources {
		values, err := src.Read()
		if err != nil {
			return fmt.Errorf("read config source %s: %w", src.Name(), err)
		}
		for k, v := range values {
			m.viper.Set(k, v)
		}
	}

	// 3. Environment variables are always the final override layer.
	envValues, err := NewEnvSource(100).Read()
	if err != nil {
		return fmt.Errorf("read config source env: %w", err)
	}
	for k, v := range envValues {
		m.viper.Set(k, v)
	}

	return nil
}

// 中文：applyEnvironmentProfileLocked 将单一配置文件中的 environments.<env> 覆盖到根配置。
// English: applyEnvironmentProfileLocked overlays environments.<env> from the single config file onto root settings.
func (m *Manager) applyEnvironmentProfileLocked(env string) {
	profile := m.viper.GetStringMap("environments." + env)
	for key, value := range flattenSettings(profile, "") {
		if strings.HasPrefix(key, "environments.") {
			continue
		}
		m.viper.Set(key, value)
	}
}

// 中文：Get 执行当前包中的对应流程。
// English: Get executes the corresponding workflow in this package.
// Get 获取配置值
func (m *Manager) Get(key string) any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.Get(key)
}

// 中文：GetString 执行当前包中的对应流程。
// English: GetString executes the corresponding workflow in this package.
// GetString 获取字符串配置
func (m *Manager) GetString(key string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.GetString(key)
}

// 中文：GetInt 执行当前包中的对应流程。
// English: GetInt executes the corresponding workflow in this package.
// GetInt 获取整数配置
func (m *Manager) GetInt(key string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.GetInt(key)
}

// 中文：GetBool 执行当前包中的对应流程。
// English: GetBool executes the corresponding workflow in this package.
// GetBool 获取布尔配置
func (m *Manager) GetBool(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.GetBool(key)
}

// 中文：GetDuration 执行当前包中的对应流程。
// English: GetDuration executes the corresponding workflow in this package.
// GetDuration 获取时间段配置
func (m *Manager) GetDuration(key string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.GetDuration(key)
}

// 中文：Unmarshal 执行当前包中的对应流程。
// English: Unmarshal executes the corresponding workflow in this package.
// Unmarshal 将配置绑定到结构体
func (m *Manager) Unmarshal(key string, target any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.UnmarshalKey(key, target)
}

// 中文：IsSet 执行当前包中的对应流程。
// English: IsSet executes the corresponding workflow in this package.
// IsSet 检查配置是否存在
func (m *Manager) IsSet(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.IsSet(key)
}

// 中文：Set 执行当前包中的对应流程。
// English: Set executes the corresponding workflow in this package.
// Set 设置配置值（运行时覆盖）
func (m *Manager) Set(key string, value any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.viper.Set(key, value)
}

// 中文：OnChange 执行当前包中的对应流程。
// English: OnChange executes the corresponding workflow in this package.
// OnChange 注册配置变更回调
func (m *Manager) OnChange(key string, callback func(value any)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onChange[key] = append(m.onChange[key], callback)
}

// 中文：Watch 执行当前包中的对应流程。
// English: Watch executes the corresponding workflow in this package.
// Watch is a compatibility entry point for both blueprint-style key watching
// and source watcher startup. Pass a context.Context to start all source
// watchers, or pass a string key plus one callback to observe that key.
func (m *Manager) Watch(target any, callbacks ...func(value any)) error {
	switch value := target.(type) {
	case context.Context:
		if len(callbacks) != 0 {
			return fmt.Errorf("watch source startup does not accept callbacks")
		}
		return m.WatchSources(value)
	case string:
		if len(callbacks) != 1 {
			return fmt.Errorf("watch key %q requires exactly one callback", value)
		}
		m.OnChange(value, callbacks[0])
		return nil
	default:
		return fmt.Errorf("unsupported watch target %T", target)
	}
}

// 中文：WatchKey 执行当前包中的对应流程。
// English: WatchKey executes the corresponding workflow in this package.
// WatchKey registers a callback for one configuration key.
func (m *Manager) WatchKey(key string, callback func(value any)) {
	m.OnChange(key, callback)
}

// 中文：WatchSources 执行当前包中的对应流程。
// English: WatchSources executes the corresponding workflow in this package.
// WatchSources starts all source watchers and applies incoming changes in-process.
func (m *Manager) WatchSources(ctx context.Context) error {
	m.mu.RLock()
	sources := append([]Source(nil), m.sources...)
	m.mu.RUnlock()

	for _, src := range sources {
		source := src
		if err := source.Watch(ctx, func(key string, value any) {
			m.applyChange(key, value)
		}); err != nil {
			return fmt.Errorf("watch config source %s: %w", source.Name(), err)
		}
	}
	return nil
}

// 中文：applyChange 执行当前包中的对应流程。
// English: applyChange executes the corresponding workflow in this package.
func (m *Manager) applyChange(key string, value any) {
	m.mu.Lock()
	m.viper.Set(key, value)
	callbacks := append([]func(any){}, m.onChange[key]...)
	m.mu.Unlock()

	for _, callback := range callbacks {
		callback(value)
	}
}

// 中文：Env 执行当前包中的对应流程。
// English: Env executes the corresponding workflow in this package.
// Env 获取当前环境
func (m *Manager) Env() string {
	return m.env
}

// 中文：AllSettings 执行当前包中的对应流程。
// English: AllSettings executes the corresponding workflow in this package.
// AllSettings 获取所有配置
func (m *Manager) AllSettings() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.AllSettings()
}
