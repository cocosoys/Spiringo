package alert

import (
	"testing"

	"github.com/getsentry/sentry-go"
)

// 中文：TestNewSentryNotifierRequiresDSN 验证相关行为符合预期。
// English: TestNewSentryNotifierRequiresDSN verifies the related behavior.
func TestNewSentryNotifierRequiresDSN(t *testing.T) {
	if _, err := NewSentryNotifier(SentryConfig{}); err == nil {
		t.Fatal("expected missing dsn to fail")
	}
}

// 中文：TestSentryLevelMapping 验证相关行为符合预期。
// English: TestSentryLevelMapping verifies the related behavior.
func TestSentryLevelMapping(t *testing.T) {
	tests := map[Severity]sentry.Level{
		SeverityInfo:     sentry.LevelInfo,
		SeverityWarning:  sentry.LevelWarning,
		SeverityCritical: sentry.LevelFatal,
	}
	for severity, want := range tests {
		if got := sentryLevel(severity); got != want {
			t.Fatalf("sentryLevel(%q) = %q, want %q", severity, got, want)
		}
	}
}
