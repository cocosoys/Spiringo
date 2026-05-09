package main

import (
	"strings"
	"testing"
)

// 中文：TestRunMigrateRequiresUp 验证相关行为符合预期。
// English: TestRunMigrateRequiresUp verifies the related behavior.
func TestRunMigrateRequiresUp(t *testing.T) {
	err := run([]string{"migrate"})
	if err == nil || !strings.Contains(err.Error(), "migrate requires: up") {
		t.Fatalf("expected migrate usage error, got %v", err)
	}
}

// 中文：TestNewApplicationParsesConfigAndEnv 验证相关行为符合预期。
// English: TestNewApplicationParsesConfigAndEnv verifies the related behavior.
func TestNewApplicationParsesConfigAndEnv(t *testing.T) {
	app, err := newApplication([]string{"-config", "configs", "-env", "dev"})
	if err != nil {
		t.Fatal(err)
	}
	if app == nil {
		t.Fatal("expected application")
	}
}
