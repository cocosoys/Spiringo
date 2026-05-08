package logger

import (
	"log/slog"
	"os"

	"github.com/spiringo/spiringo/internal/core/config"
)

// 中文：Logger 定义当前包使用的数据结构或接口。
// English: Logger defines a data structure or interface used by this package.
// Logger 日志接口
type Logger interface {
	// 中文：Debug 声明该接口需要实现的行为。
	// English: Debug declares behavior required by this interface.
	Debug(msg string, args ...any)
	// 中文：Info 声明该接口需要实现的行为。
	// English: Info declares behavior required by this interface.
	Info(msg string, args ...any)
	// 中文：Warn 声明该接口需要实现的行为。
	// English: Warn declares behavior required by this interface.
	Warn(msg string, args ...any)
	// 中文：Error 声明该接口需要实现的行为。
	// English: Error declares behavior required by this interface.
	Error(msg string, args ...any)
	// 中文：With 声明该接口需要实现的行为。
	// English: With declares behavior required by this interface.
	With(args ...any) Logger
}

// 中文：slogLogger 定义当前包使用的数据结构或接口。
// English: slogLogger defines a data structure or interface used by this package.
// slogLogger 基于slog的日志实现
type slogLogger struct {
	// 中文：inner 保存当前结构中的配置或数据值。
	// English: inner stores a configuration or data value for this struct.
	inner *slog.Logger
}

// 中文：New 创建并返回对应组件实例。
// English: New creates and returns the corresponding component instance.
// New 创建日志器
func New(cfg config.LogConfig) Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: parseLevel(cfg.Level),
	}

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &slogLogger{inner: slog.New(handler)}
}

// 中文：NewNop 创建并返回对应组件实例。
// English: NewNop creates and returns the corresponding component instance.
// NewNop 创建空日志器
func NewNop() Logger {
	return &slogLogger{inner: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.Level(-100)}))}
}

// 中文：NewFromSlog 创建并返回对应组件实例。
// English: NewFromSlog creates and returns the corresponding component instance.
// NewFromSlog 从slog.Logger创建
func NewFromSlog(l *slog.Logger) Logger {
	return &slogLogger{inner: l}
}

// 中文：Debug 执行当前包中的对应流程。
// English: Debug executes the corresponding workflow in this package.
func (l *slogLogger) Debug(msg string, args ...any) { l.inner.Debug(msg, args...) }

// 中文：Info 执行当前包中的对应流程。
// English: Info executes the corresponding workflow in this package.
func (l *slogLogger) Info(msg string, args ...any) { l.inner.Info(msg, args...) }

// 中文：Warn 执行当前包中的对应流程。
// English: Warn executes the corresponding workflow in this package.
func (l *slogLogger) Warn(msg string, args ...any) { l.inner.Warn(msg, args...) }

// 中文：Error 执行当前包中的对应流程。
// English: Error executes the corresponding workflow in this package.
func (l *slogLogger) Error(msg string, args ...any) { l.inner.Error(msg, args...) }

// 中文：With 执行当前包中的对应流程。
// English: With executes the corresponding workflow in this package.
func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{inner: l.inner.With(args...)}
}

// 中文：parseLevel 执行当前包中的对应流程。
// English: parseLevel executes the corresponding workflow in this package.
func parseLevel(level string) slog.Level {
	switch level {
	case "debug", "DEBUG":
		return slog.LevelDebug
	case "info", "INFO":
		return slog.LevelInfo
	case "warn", "WARN", "warning":
		return slog.LevelWarn
	case "error", "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
