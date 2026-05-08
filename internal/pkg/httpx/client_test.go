package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// 中文：TestClientRetriesAndDecodesJSON 验证相关行为符合预期。
// English: TestClientRetriesAndDecodesJSON verifies the related behavior.
func TestClientRetriesAndDecodesJSON(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			http.Error(w, "try again", http.StatusServiceUnavailable)
			return
		}
		if r.Header.Get("X-App") != "spiringo" {
			http.Error(w, "missing header", http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"ok": "yes"})
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL:   server.URL,
		Headers:   map[string]string{"X-App": "spiringo"},
		Retries:   1,
		RetryWait: time.Millisecond,
	})

	var out map[string]string
	if err := client.Get(context.Background(), "/status", &out); err != nil {
		t.Fatal(err)
	}
	if attempts != 2 || out["ok"] != "yes" {
		t.Fatalf("attempts=%d out=%v", attempts, out)
	}
}

// 中文：TestClientReturnsHTTPError 验证相关行为符合预期。
// English: TestClientReturnsHTTPError verifies the related behavior.
func TestClientReturnsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad input", http.StatusBadRequest)
	}))
	defer server.Close()

	err := NewClient(Config{BaseURL: server.URL}).Get(context.Background(), "/", nil)
	var httpErr *Error
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http error, got %v", err)
	}
}

// 中文：TestURLWithQuery 验证相关行为符合预期。
// English: TestURLWithQuery verifies the related behavior.
func TestURLWithQuery(t *testing.T) {
	got := URLWithQuery("/items?keep=1", map[string]string{"q": "go"})
	if got != "/items?keep=1&q=go" {
		t.Fatalf("URLWithQuery = %q", got)
	}
}
