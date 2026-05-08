package logger

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/spiringo/spiringo/internal/core/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 中文：zapLogger 定义当前包使用的数据结构或接口。
// English: zapLogger defines a data structure or interface used by this package.
type zapLogger struct {
	// 中文：inner 保存当前结构中的配置或数据值。
	// English: inner stores a configuration or data value for this struct.
	inner *zap.SugaredLogger
}

// 中文：NewZap 创建并返回对应组件实例。
// English: NewZap creates and returns the corresponding component instance.
func NewZap(cfg config.LogConfig) (Logger, func() error, error) {
	zl, err := buildZap(cfg)
	if err != nil {
		return nil, nil, err
	}
	return &zapLogger{inner: zl.Sugar()}, zl.Sync, nil
}

// 中文：NewZapSlog 创建并返回对应组件实例。
// English: NewZapSlog creates and returns the corresponding component instance.
func NewZapSlog(cfg config.LogConfig) (*slog.Logger, func() error, error) {
	zl, err := buildZap(cfg)
	if err != nil {
		return nil, nil, err
	}
	return slog.New(&zapSlogHandler{logger: zl}), zl.Sync, nil
}

// 中文：Debug 执行当前包中的对应流程。
// English: Debug executes the corresponding workflow in this package.
func (l *zapLogger) Debug(msg string, args ...any) { l.inner.Debugw(msg, args...) }

// 中文：Info 执行当前包中的对应流程。
// English: Info executes the corresponding workflow in this package.
func (l *zapLogger) Info(msg string, args ...any) { l.inner.Infow(msg, args...) }

// 中文：Warn 执行当前包中的对应流程。
// English: Warn executes the corresponding workflow in this package.
func (l *zapLogger) Warn(msg string, args ...any) { l.inner.Warnw(msg, args...) }

// 中文：Error 执行当前包中的对应流程。
// English: Error executes the corresponding workflow in this package.
func (l *zapLogger) Error(msg string, args ...any) { l.inner.Errorw(msg, args...) }

// 中文：With 执行当前包中的对应流程。
// English: With executes the corresponding workflow in this package.
func (l *zapLogger) With(args ...any) Logger {
	return &zapLogger{inner: l.inner.With(args...)}
}

// 中文：buildZap 执行当前包中的对应流程。
// English: buildZap executes the corresponding workflow in this package.
func buildZap(cfg config.LogConfig) (*zap.Logger, error) {
	zcfg := zap.NewProductionConfig()
	if strings.EqualFold(cfg.Format, "console") || strings.EqualFold(cfg.Format, "text") {
		zcfg.Encoding = "console"
	}
	zcfg.Level = zap.NewAtomicLevelAt(toZapcoreLevel(cfg.Level))
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	switch strings.ToLower(strings.TrimSpace(cfg.Output)) {
	case "", "stdout":
		zcfg.OutputPaths = []string{"stdout"}
	case "stderr":
		zcfg.OutputPaths = []string{"stderr"}
	default:
		zcfg.OutputPaths = []string{cfg.Output}
	}
	zcfg.ErrorOutputPaths = []string{"stderr"}
	return zcfg.Build()
}

// 中文：zapSlogHandler 定义当前包使用的数据结构或接口。
// English: zapSlogHandler defines a data structure or interface used by this package.
type zapSlogHandler struct {
	// 中文：logger 保存当前结构中的配置或数据值。
	// English: logger stores a configuration or data value for this struct.
	logger *zap.Logger
	// 中文：attrs 保存当前结构中的配置或数据值。
	// English: attrs stores a configuration or data value for this struct.
	attrs []slog.Attr
	// 中文：groups 保存当前结构中的配置或数据值。
	// English: groups stores a configuration or data value for this struct.
	groups []string
}

// 中文：Enabled 执行当前包中的对应流程。
// English: Enabled executes the corresponding workflow in this package.
func (h *zapSlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.logger.Core().Enabled(slogToZapLevel(level))
}

// 中文：Handle 执行当前包中的对应流程。
// English: Handle executes the corresponding workflow in this package.
func (h *zapSlogHandler) Handle(_ context.Context, record slog.Record) error {
	fields := make([]zap.Field, 0, len(h.attrs)+record.NumAttrs())
	for _, attr := range h.attrs {
		fields = appendAttr(fields, h.groups, attr)
	}
	record.Attrs(func(attr slog.Attr) bool {
		fields = appendAttr(fields, h.groups, attr)
		return true
	})
	h.logger.Log(slogToZapLevel(record.Level), record.Message, fields...)
	return nil
}

// 中文：WithAttrs 执行当前包中的对应流程。
// English: WithAttrs executes the corresponding workflow in this package.
func (h *zapSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := *h
	next.attrs = append(append([]slog.Attr(nil), h.attrs...), attrs...)
	return &next
}

// 中文：WithGroup 执行当前包中的对应流程。
// English: WithGroup executes the corresponding workflow in this package.
func (h *zapSlogHandler) WithGroup(name string) slog.Handler {
	if strings.TrimSpace(name) == "" {
		return h
	}
	next := *h
	next.groups = append(append([]string(nil), h.groups...), name)
	return &next
}

// 中文：appendAttr 执行当前包中的对应流程。
// English: appendAttr executes the corresponding workflow in this package.
func appendAttr(fields []zap.Field, groups []string, attr slog.Attr) []zap.Field {
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return fields
	}
	if attr.Value.Kind() == slog.KindGroup {
		groupPath := append(append([]string(nil), groups...), attr.Key)
		for _, child := range attr.Value.Group() {
			fields = appendAttr(fields, groupPath, child)
		}
		return fields
	}
	key := groupedKey(groups, attr.Key)
	switch attr.Value.Kind() {
	case slog.KindString:
		return append(fields, zap.String(key, attr.Value.String()))
	case slog.KindBool:
		return append(fields, zap.Bool(key, attr.Value.Bool()))
	case slog.KindInt64:
		return append(fields, zap.Int64(key, attr.Value.Int64()))
	case slog.KindUint64:
		return append(fields, zap.Uint64(key, attr.Value.Uint64()))
	case slog.KindFloat64:
		return append(fields, zap.Float64(key, attr.Value.Float64()))
	case slog.KindDuration:
		return append(fields, zap.Duration(key, time.Duration(attr.Value.Int64())))
	case slog.KindTime:
		return append(fields, zap.Time(key, attr.Value.Time()))
	default:
		return append(fields, zap.Any(key, attr.Value.Any()))
	}
}

// 中文：groupedKey 执行当前包中的对应流程。
// English: groupedKey executes the corresponding workflow in this package.
func groupedKey(groups []string, key string) string {
	if len(groups) == 0 {
		return key
	}
	parts := append(append([]string(nil), groups...), key)
	return strings.Join(parts, ".")
}

// 中文：slogToZapLevel 执行当前包中的对应流程。
// English: slogToZapLevel executes the corresponding workflow in this package.
func slogToZapLevel(level slog.Level) zapcore.Level {
	switch {
	case level <= slog.LevelDebug:
		return zapcore.DebugLevel
	case level >= slog.LevelError:
		return zapcore.ErrorLevel
	case level >= slog.LevelWarn:
		return zapcore.WarnLevel
	default:
		return zapcore.InfoLevel
	}
}

// 中文：toZapcoreLevel 执行当前包中的对应流程。
// English: toZapcoreLevel executes the corresponding workflow in this package.
func toZapcoreLevel(level string) zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return zapcore.DebugLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
