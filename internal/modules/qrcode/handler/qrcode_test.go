package handler

import (
	"encoding/base64"
	"strings"
	"testing"
)

// 中文：TestDecodeParseImageBase64AcceptsDataURL 验证相关行为符合预期。
// English: TestDecodeParseImageBase64AcceptsDataURL verifies the related behavior.
func TestDecodeParseImageBase64AcceptsDataURL(t *testing.T) {
	want := []byte("png-bytes")
	got, err := decodeParseImageBase64("data:image/png;base64," + base64.StdEncoding.EncodeToString(want))
	if err != nil {
		t.Fatalf("decodeParseImageBase64 returned error: %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("decoded = %q, want %q", got, want)
	}
}

// 中文：TestDecodeParseImageBase64RejectsEmptyAndInvalid 验证相关行为符合预期。
// English: TestDecodeParseImageBase64RejectsEmptyAndInvalid verifies the related behavior.
func TestDecodeParseImageBase64RejectsEmptyAndInvalid(t *testing.T) {
	if _, err := decodeParseImageBase64(""); err == nil {
		t.Fatal("expected empty value error")
	}
	if _, err := decodeParseImageBase64("not-base64"); err == nil {
		t.Fatal("expected invalid base64 error")
	}
}

// 中文：TestReadLimitedImageRejectsOversizedPayload 验证相关行为符合预期。
// English: TestReadLimitedImageRejectsOversizedPayload verifies the related behavior.
func TestReadLimitedImageRejectsOversizedPayload(t *testing.T) {
	payload := strings.NewReader(strings.Repeat("x", maxParseImageBytes+1))
	if _, err := readLimitedImage(payload); err == nil {
		t.Fatal("expected oversized image error")
	}
}
