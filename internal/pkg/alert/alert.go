package alert

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

// 中文：Severity 定义当前包使用的数据结构或接口。
// English: Severity defines a data structure or interface used by this package.
type Severity string

// 中文：SeverityInfo、SeverityWarning、SeverityCritical 声明当前包使用的常量。
// English: SeverityInfo、SeverityWarning、SeverityCritical declares constants used by this package.
const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// 中文：Alert 定义当前包使用的数据结构或接口。
// English: Alert defines a data structure or interface used by this package.
type Alert struct {
	// 中文：Title 保存当前结构中的配置或数据值。
	// English: Title stores a configuration or data value for this struct.
	Title string `json:"title"`
	// 中文：Message 保存当前结构中的配置或数据值。
	// English: Message stores a configuration or data value for this struct.
	Message string `json:"message"`
	// 中文：Severity 保存当前结构中的配置或数据值。
	// English: Severity stores a configuration or data value for this struct.
	Severity Severity `json:"severity"`
	// 中文：Source 保存当前结构中的配置或数据值。
	// English: Source stores a configuration or data value for this struct.
	Source string `json:"source,omitempty"`
	// 中文：Labels 保存当前结构中的配置或数据值。
	// English: Labels stores a configuration or data value for this struct.
	Labels map[string]string `json:"labels,omitempty"`
	// 中文：Timestamp 保存当前结构中的配置或数据值。
	// English: Timestamp stores a configuration or data value for this struct.
	Timestamp time.Time `json:"timestamp"`
}

// 中文：Notifier 定义当前包使用的数据结构或接口。
// English: Notifier defines a data structure or interface used by this package.
type Notifier interface {
	// 中文：Notify 声明该接口需要实现的行为。
	// English: Notify declares behavior required by this interface.
	Notify(ctx context.Context, alert Alert) error
}

// 中文：Manager 定义当前包使用的数据结构或接口。
// English: Manager defines a data structure or interface used by this package.
type Manager struct {
	// 中文：sinks 保存当前结构中的配置或数据值。
	// English: sinks stores a configuration or data value for this struct.
	sinks []Notifier
}

// 中文：NewManager 创建并返回对应组件实例。
// English: NewManager creates and returns the corresponding component instance.
func NewManager(sinks ...Notifier) *Manager {
	return &Manager{sinks: sinks}
}

// 中文：Notify 执行当前包中的对应流程。
// English: Notify executes the corresponding workflow in this package.
func (m *Manager) Notify(ctx context.Context, a Alert) error {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	var errs []error
	for _, sink := range m.sinks {
		if sink == nil {
			continue
		}
		if err := sink.Notify(ctx, a); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (m *Manager) Close() error {
	var errs []error
	for _, sink := range m.sinks {
		if closer, ok := sink.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

// 中文：LoggerNotifier 定义当前包使用的数据结构或接口。
// English: LoggerNotifier defines a data structure or interface used by this package.
type LoggerNotifier struct {
	// 中文：logger 保存当前结构中的配置或数据值。
	// English: logger stores a configuration or data value for this struct.
	logger *slog.Logger
}

// 中文：NewLoggerNotifier 创建并返回对应组件实例。
// English: NewLoggerNotifier creates and returns the corresponding component instance.
func NewLoggerNotifier(logger *slog.Logger) *LoggerNotifier {
	if logger == nil {
		logger = slog.Default()
	}
	return &LoggerNotifier{logger: logger}
}

// 中文：Notify 执行当前包中的对应流程。
// English: Notify executes the corresponding workflow in this package.
func (n *LoggerNotifier) Notify(ctx context.Context, a Alert) error {
	args := []any{
		"severity", a.Severity,
		"source", a.Source,
		"labels", a.Labels,
	}
	switch a.Severity {
	case SeverityCritical:
		n.logger.ErrorContext(ctx, a.Title, append(args, "message", a.Message)...)
	case SeverityWarning:
		n.logger.WarnContext(ctx, a.Title, append(args, "message", a.Message)...)
	default:
		n.logger.InfoContext(ctx, a.Title, append(args, "message", a.Message)...)
	}
	return nil
}
