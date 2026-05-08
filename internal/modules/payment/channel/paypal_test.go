package channel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/plutov/paypal/v4"
)

// 中文：TestPayPalVerifyCallbackParsesOrderCompletedEvent 验证相关行为符合预期。
// English: TestPayPalVerifyCallbackParsesOrderCompletedEvent verifies the related behavior.
func TestPayPalVerifyCallbackParsesOrderCompletedEvent(t *testing.T) {
	ch := NewPayPalChannel("", "", true, "")
	raw := []byte(`{
		"event_type":"CHECKOUT.ORDER.COMPLETED",
		"resource":{
			"id":"ORDER-123",
			"status":"COMPLETED",
			"purchase_units":[{
				"custom_id":"merchant-001",
				"payments":{
					"captures":[{
						"id":"CAPTURE-123",
						"status":"COMPLETED",
						"amount":{"value":"12.34","currency_code":"USD"}
					}]
				}
			}]
		}
	}`)

	result, err := ch.VerifyCallback(context.Background(), raw)
	if err != nil {
		t.Fatalf("VerifyCallback returned error: %v", err)
	}
	if result.OutTradeNo != "merchant-001" {
		t.Fatalf("OutTradeNo = %q, want merchant-001", result.OutTradeNo)
	}
	if result.TradeNo != "ORDER-123" {
		t.Fatalf("TradeNo = %q, want ORDER-123", result.TradeNo)
	}
	if result.Status != "paid" {
		t.Fatalf("Status = %q, want paid", result.Status)
	}
	if result.Amount != 1234 {
		t.Fatalf("Amount = %d, want 1234", result.Amount)
	}
}

// 中文：TestPayPalVerifyCallbackRequiresMerchantOrderID 验证相关行为符合预期。
// English: TestPayPalVerifyCallbackRequiresMerchantOrderID verifies the related behavior.
func TestPayPalVerifyCallbackRequiresMerchantOrderID(t *testing.T) {
	ch := NewPayPalChannel("", "", true, "")
	raw := []byte(`{"event_type":"CHECKOUT.ORDER.COMPLETED","resource":{"id":"ORDER-123"}}`)

	if _, err := ch.VerifyCallback(context.Background(), raw); err == nil {
		t.Fatal("VerifyCallback returned nil error, want missing merchant order id error")
	}
}

// 中文：TestPayPalClosePaymentVoidsAuthorizations 验证相关行为符合预期。
// English: TestPayPalClosePaymentVoidsAuthorizations verifies the related behavior.
func TestPayPalClosePaymentVoidsAuthorizations(t *testing.T) {
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		calls = append(calls, req.Method+" "+req.URL.Path)
		switch req.Method + " " + req.URL.Path {
		case "GET /v2/checkout/orders/ORDER-1":
			_, _ = w.Write([]byte(`{
				"id":"ORDER-1",
				"status":"APPROVED",
				"purchase_units":[{"payments":{"authorizations":[{"id":"AUTH-1","status":"CREATED"}]}}]
			}`))
		case "POST /v2/payments/authorizations/AUTH-1/void":
			_, _ = w.Write([]byte(`{"id":"AUTH-1","status":"VOIDED"}`))
		default:
			http.NotFound(w, req)
		}
	}))
	defer server.Close()

	client, err := paypal.NewClient("client-id", "secret", server.URL)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	client.SetAccessToken("token")
	ch := &PayPalChannel{client: client}

	if err := ch.ClosePayment(context.Background(), "ORDER-1"); err != nil {
		t.Fatalf("ClosePayment returned error: %v", err)
	}
	if len(calls) != 2 ||
		calls[0] != "GET /v2/checkout/orders/ORDER-1" ||
		calls[1] != "POST /v2/payments/authorizations/AUTH-1/void" {
		t.Fatalf("unexpected calls: %v", calls)
	}
}

// 中文：TestPayPalClosePaymentRejectsCompletedOrder 验证相关行为符合预期。
// English: TestPayPalClosePaymentRejectsCompletedOrder verifies the related behavior.
func TestPayPalClosePaymentRejectsCompletedOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet || req.URL.Path != "/v2/checkout/orders/ORDER-1" {
			http.NotFound(w, req)
			return
		}
		_, _ = w.Write([]byte(`{"id":"ORDER-1","status":"COMPLETED"}`))
	}))
	defer server.Close()

	client, err := paypal.NewClient("client-id", "secret", server.URL)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	client.SetAccessToken("token")
	ch := &PayPalChannel{client: client}

	err = ch.ClosePayment(context.Background(), "ORDER-1")
	if err == nil || !strings.Contains(err.Error(), "completed order") {
		t.Fatalf("ClosePayment error = %v, want completed order error", err)
	}
}
