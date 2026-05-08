package oss

import "testing"

// 中文：TestBlueprintOSSPackageAliasesStorageImplementations 验证相关行为符合预期。
// English: TestBlueprintOSSPackageAliasesStorageImplementations verifies the related behavior.
func TestBlueprintOSSPackageAliasesStorageImplementations(t *testing.T) {
	if _, err := NewCephStorage(CephConfig{Endpoint: "127.0.0.1:7480"}); err != nil {
		t.Fatalf("NewCephStorage returned error: %v", err)
	}
}
