package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/spiringo/spiringo/internal/modules/notification/model"
	"github.com/spiringo/spiringo/internal/modules/notification/repository"
	"github.com/spiringo/spiringo/internal/pkg/orm"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：TestNotifySendsWebhookForEnabledEvent 验证相关行为符合预期。
// English: TestNotifySendsWebhookForEnabledEvent verifies the related behavior.
func TestNotifySendsWebhookForEnabledEvent(t *testing.T) {
	var received Message
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Notify-Token"); got != "secret" {
			t.Fatalf("X-Notify-Token = %q, want secret", got)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode webhook body: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := New(Config{
		Events: []string{"payment.failed"},
		Webhook: WebhookConfig{
			Enabled: true,
			URLs:    []string{server.URL},
			Headers: map[string]string{"X-Notify-Token": "secret"},
		},
	})
	err := svc.Notify(context.Background(), Message{
		Event:    "payment.failed",
		Severity: "error",
		Subject:  "Payment failed",
		Content:  "order failed",
		Payload:  map[string]any{"out_trade_no": "order-001"},
	})
	if err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
	if received.Event != "payment.failed" {
		t.Fatalf("webhook event = %q, want payment.failed", received.Event)
	}
	if received.Payload["out_trade_no"] != "order-001" {
		t.Fatalf("webhook payload out_trade_no = %v, want order-001", received.Payload["out_trade_no"])
	}
}

// 中文：TestNotifySkipsDisabledEvent 验证相关行为符合预期。
// English: TestNotifySkipsDisabledEvent verifies the related behavior.
func TestNotifySkipsDisabledEvent(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := New(Config{
		Events: []string{"payment.failed"},
		Webhook: WebhookConfig{
			Enabled: true,
			URLs:    []string{server.URL},
		},
	})
	if err := svc.Notify(context.Background(), Message{Event: "payment.success"}); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
	if calls != 0 {
		t.Fatalf("webhook calls = %d, want 0", calls)
	}
}

// 中文：TestNotifyReturnsWebhookStatusError 验证相关行为符合预期。
// English: TestNotifyReturnsWebhookStatusError verifies the related behavior.
func TestNotifyReturnsWebhookStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	svc := New(Config{
		Webhook: WebhookConfig{
			Enabled: true,
			URLs:    []string{server.URL},
		},
	})
	if err := svc.Notify(context.Background(), Message{Event: "payment.failed"}); err == nil {
		t.Fatal("Notify returned nil error, want webhook status error")
	}
}

// 中文：TestNotifyStoresInboxAndMarksRead 验证相关行为符合预期。
// English: TestNotifyStoresInboxAndMarksRead verifies the related behavior.
func TestNotifyStoresInboxAndMarksRead(t *testing.T) {
	db, err := orm.New(orm.Config{
		Driver: "sqlite",
		DSN:    filepath.Join(t.TempDir(), "notification_test.db"),
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()
	if err := db.AutoMigrate(&model.Notification{}); err != nil {
		t.Fatalf("migrate notification: %v", err)
	}

	repo := repository.NewNotificationRepository(orm.NewTenantDB(db), db)
	svc := New(Config{
		Events: []string{"payment.failed"},
		Inbox:  InboxConfig{Enabled: true},
	}, repo)
	ctx := types.WithTenantID(context.Background(), "tenant-1")
	if err := svc.Notify(ctx, Message{
		Event:       "payment.failed",
		Severity:    "error",
		Subject:     "Payment failed",
		Content:     "order failed",
		RecipientID: "user-1",
		Payload:     map[string]any{"out_trade_no": "order-001"},
	}); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}

	items, total, err := svc.ListInbox(ctx, InboxFilter{
		Page:        1,
		PageSize:    20,
		RecipientID: "user-1",
		UnreadOnly:  true,
	})
	if err != nil {
		t.Fatalf("ListInbox returned error: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected one inbox item, total=%d len=%d", total, len(items))
	}
	if items[0].Payload["out_trade_no"] != "order-001" {
		t.Fatalf("unexpected inbox payload: %+v", items[0].Payload)
	}

	if err := svc.MarkRead(ctx, items[0].ID, "user-1"); err != nil {
		t.Fatalf("MarkRead returned error: %v", err)
	}
	unread, total, err := svc.ListInbox(ctx, InboxFilter{Page: 1, PageSize: 20, RecipientID: "user-1", UnreadOnly: true})
	if err != nil {
		t.Fatalf("ListInbox unread returned error: %v", err)
	}
	if total != 0 || len(unread) != 0 {
		t.Fatalf("expected no unread inbox items, total=%d len=%d", total, len(unread))
	}
	all, total, err := svc.ListInbox(ctx, InboxFilter{Page: 1, PageSize: 20, RecipientID: "user-1"})
	if err != nil {
		t.Fatalf("ListInbox all returned error: %v", err)
	}
	if total != 1 || all[0].ReadAt == nil {
		t.Fatalf("expected read inbox item, total=%d item=%+v", total, all[0])
	}
}

// 中文：TestNotifyStoresInboxTenantFromMessageWhenContextIsEmpty 验证相关行为符合预期。
// English: TestNotifyStoresInboxTenantFromMessageWhenContextIsEmpty verifies the related behavior.
func TestNotifyStoresInboxTenantFromMessageWhenContextIsEmpty(t *testing.T) {
	db, err := orm.New(orm.Config{
		Driver: "sqlite",
		DSN:    filepath.Join(t.TempDir(), "notification_message_tenant.db"),
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()
	if err := db.AutoMigrate(&model.Notification{}); err != nil {
		t.Fatalf("migrate notification: %v", err)
	}

	repo := repository.NewNotificationRepository(orm.NewTenantDB(db), db)
	svc := New(Config{
		Events: []string{"payment.success"},
		Inbox:  InboxConfig{Enabled: true},
	}, repo)
	if err := svc.Notify(context.Background(), Message{
		Event:       "payment.success",
		TenantID:    "tenant-from-event",
		RecipientID: "user-2",
		Subject:     "Payment succeeded",
	}); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}

	items, total, err := svc.ListInbox(types.WithTenantID(context.Background(), "tenant-from-event"), InboxFilter{
		Page:        1,
		PageSize:    20,
		RecipientID: "user-2",
	})
	if err != nil {
		t.Fatalf("ListInbox returned error: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected tenant-scoped inbox item, total=%d len=%d", total, len(items))
	}
}
