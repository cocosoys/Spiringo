package alert

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
)

// 中文：SentryConfig 定义当前包使用的数据结构或接口。
// English: SentryConfig defines a data structure or interface used by this package.
type SentryConfig struct {
	// 中文：DSN 保存当前结构中的配置或数据值。
	// English: DSN stores a configuration or data value for this struct.
	DSN string
	// 中文：Environment 保存当前结构中的配置或数据值。
	// English: Environment stores a configuration or data value for this struct.
	Environment string
	// 中文：Release 保存当前结构中的配置或数据值。
	// English: Release stores a configuration or data value for this struct.
	Release string
	// 中文：TracesSampleRate 保存当前结构中的配置或数据值。
	// English: TracesSampleRate stores a configuration or data value for this struct.
	TracesSampleRate float64
	// 中文：Debug 保存当前结构中的配置或数据值。
	// English: Debug stores a configuration or data value for this struct.
	Debug bool
	// 中文：FlushTimeout 保存当前结构中的配置或数据值。
	// English: FlushTimeout stores a configuration or data value for this struct.
	FlushTimeout time.Duration
}

// 中文：SentryNotifier 定义当前包使用的数据结构或接口。
// English: SentryNotifier defines a data structure or interface used by this package.
type SentryNotifier struct {
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *sentry.Client
	// 中文：flushTimeout 保存当前结构中的配置或数据值。
	// English: flushTimeout stores a configuration or data value for this struct.
	flushTimeout time.Duration
}

// 中文：NewSentryNotifier 创建并返回对应组件实例。
// English: NewSentryNotifier creates and returns the corresponding component instance.
func NewSentryNotifier(cfg SentryConfig) (*SentryNotifier, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("sentry dsn is required")
	}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              cfg.DSN,
		Environment:      cfg.Environment,
		Release:          cfg.Release,
		TracesSampleRate: cfg.TracesSampleRate,
		Debug:            cfg.Debug,
	})
	if err != nil {
		return nil, fmt.Errorf("init sentry: %w", err)
	}
	if cfg.FlushTimeout <= 0 {
		cfg.FlushTimeout = 2 * time.Second
	}
	return &SentryNotifier{client: client, flushTimeout: cfg.FlushTimeout}, nil
}

// 中文：Notify 执行当前包中的对应流程。
// English: Notify executes the corresponding workflow in this package.
func (n *SentryNotifier) Notify(ctx context.Context, a Alert) error {
	if n.client == nil {
		return fmt.Errorf("sentry client is not configured")
	}
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}
	scope := sentry.NewScope()
	scope.SetLevel(sentryLevel(a.Severity))
	scope.SetTag("source", a.Source)
	scope.SetTag("severity", string(a.Severity))
	for key, value := range a.Labels {
		scope.SetTag(key, value)
	}
	if eventID := n.client.CaptureMessage(strings.TrimSpace(a.Title+": "+a.Message), nil, scope); eventID == nil {
		return fmt.Errorf("sentry event was not accepted")
	}
	if deadline, ok := ctx.Deadline(); ok {
		flushCtx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()
		if !n.client.FlushWithContext(flushCtx) {
			return fmt.Errorf("sentry flush timed out")
		}
		return nil
	}
	if !n.client.Flush(n.flushTimeout) {
		return fmt.Errorf("sentry flush timed out")
	}
	return nil
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (n *SentryNotifier) Close() error {
	if n.client != nil {
		n.client.Close()
	}
	return nil
}

// 中文：sentryLevel 执行当前包中的对应流程。
// English: sentryLevel executes the corresponding workflow in this package.
func sentryLevel(severity Severity) sentry.Level {
	switch severity {
	case SeverityCritical:
		return sentry.LevelFatal
	case SeverityWarning:
		return sentry.LevelWarning
	default:
		return sentry.LevelInfo
	}
}
