package service

import "testing"

// 中文：TestNormalizeDefaultAdminConfigKeepsExplicitDisable 验证相关行为符合预期。
// English: TestNormalizeDefaultAdminConfigKeepsExplicitDisable verifies the related behavior.
func TestNormalizeDefaultAdminConfigKeepsExplicitDisable(t *testing.T) {
	cfg := normalizeDefaultAdminConfig(DefaultAdminConfig{
		Enabled: false,
	})

	if cfg.Enabled {
		t.Fatal("expected explicit disable to be preserved")
	}
	if cfg.Username != "admin_%s" || cfg.Password != "changeme" || cfg.Nickname != "System Admin" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

// 中文：TestDefaultAdminConfigSeedFormatsTenantValues 验证相关行为符合预期。
// English: TestDefaultAdminConfigSeedFormatsTenantValues verifies the related behavior.
func TestDefaultAdminConfigSeedFormatsTenantValues(t *testing.T) {
	cfg := normalizeDefaultAdminConfig(DefaultAdminConfig{
		Enabled:       true,
		Username:      "owner_%s",
		EmailTemplate: "owner+%s@example.com",
		Nickname:      "Owner",
		Password:      "secret",
	})

	seed := cfg.seed("tenant-1")
	if seed.username != "owner_tenant-1" {
		t.Fatalf("username = %q, want owner_tenant-1", seed.username)
	}
	if seed.email != "owner+tenant-1@example.com" {
		t.Fatalf("email = %q, want owner+tenant-1@example.com", seed.email)
	}
	if seed.nickname != "Owner" || seed.password != "secret" {
		t.Fatalf("unexpected seed: %+v", seed)
	}
}
