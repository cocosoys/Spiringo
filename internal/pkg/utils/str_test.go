package utils

import (
	"strings"
	"testing"
)

// 中文：TestGenerateRandomString 验证相关行为符合预期。
// English: TestGenerateRandomString verifies the related behavior.
func TestGenerateRandomString(t *testing.T) {
	s1 := GenerateRandomString(16)
	s2 := GenerateRandomString(16)

	if len(s1) != 16 {
		t.Errorf("expected length 16, got %d", len(s1))
	}
	if s1 == s2 {
		t.Error("two random strings should not be equal")
	}
}

// 中文：TestContainsString 验证相关行为符合预期。
// English: TestContainsString verifies the related behavior.
func TestContainsString(t *testing.T) {
	slice := []string{"a", "b", "c"}
	if !ContainsString(slice, "b") {
		t.Error("expected to find 'b' in slice")
	}
	if ContainsString(slice, "d") {
		t.Error("expected not to find 'd' in slice")
	}
}

// 中文：TestDefaultString 验证相关行为符合预期。
// English: TestDefaultString verifies the related behavior.
func TestDefaultString(t *testing.T) {
	if DefaultString("", "default") != "default" {
		t.Error("expected default value for empty string")
	}
	if DefaultString("value", "default") != "value" {
		t.Error("expected original value for non-empty string")
	}
}

// 中文：TestMaskString 验证相关行为符合预期。
// English: TestMaskString verifies the related behavior.
func TestMaskString(t *testing.T) {
	result := MaskString("1234567890", 3, 4)
	if !strings.HasPrefix(result, "123") {
		t.Errorf("expected prefix '123', got '%s'", result)
	}
	if !strings.HasSuffix(result, "7890") {
		t.Errorf("expected suffix '7890', got '%s'", result)
	}
}
