package alert

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 中文：TestWebhookNotifierSendsAlert 验证相关行为符合预期。
// English: TestWebhookNotifierSendsAlert verifies the related behavior.
func TestWebhookNotifierSendsAlert(t *testing.T) {
	var got Alert
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Token") != "secret" {
			t.Fatalf("X-Token = %q, want secret", r.Header.Get("X-Token"))
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode alert: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	notifier := NewWebhookNotifier(server.URL, 0, map[string]string{"X-Token": "secret"})
	err := notifier.Notify(context.Background(), Alert{
		Title:    "event failed",
		Message:  "boom",
		Severity: SeverityWarning,
	})
	if err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
	if got.Title != "event failed" || got.Severity != SeverityWarning {
		t.Fatalf("unexpected alert: %+v", got)
	}
}
