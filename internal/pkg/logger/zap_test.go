package logger

import (
	"testing"

	"github.com/spiringo/spiringo/internal/core/config"
)

// 中文：TestNewZapSlogCreatesUsableLogger 验证相关行为符合预期。
// English: TestNewZapSlogCreatesUsableLogger verifies the related behavior.
func TestNewZapSlogCreatesUsableLogger(t *testing.T) {
	l, syncFn, err := NewZapSlog(config.LogConfig{
		Driver: "zap",
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	})
	if err != nil {
		t.Fatal(err)
	}
	if l == nil || syncFn == nil {
		t.Fatal("expected logger and sync function")
	}

	l.Info("zap logger ready", "component", "test")
}
