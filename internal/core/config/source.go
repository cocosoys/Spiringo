package config

import "context"

// 中文：Source 定义当前包使用的数据结构或接口。
// English: Source defines a data structure or interface used by this package.
// Source 配置源接口
type Source interface {
	// 中文：Name 声明该接口需要实现的行为。
	// English: Name declares behavior required by this interface.
	// Name 配置源名称
	Name() string
	// 中文：Priority 声明该接口需要实现的行为。
	// English: Priority declares behavior required by this interface.
	// Priority 优先级（越大越优先）
	Priority() int
	// 中文：Read 声明该接口需要实现的行为。
	// English: Read declares behavior required by this interface.
	// Read 读取配置
	Read() (map[string]any, error)
	// 中文：Watch 声明该接口需要实现的行为。
	// English: Watch declares behavior required by this interface.
	// Watch 监听配置变更（如不支持返回nil）
	Watch(ctx context.Context, onChange func(key string, value any)) error
	// 中文：Close 声明该接口需要实现的行为。
	// English: Close declares behavior required by this interface.
	// Close 关闭
	Close() error
}
