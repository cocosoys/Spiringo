package channel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 中文：TestCloudPayCreatePaymentPostsSignedPayload 验证相关行为符合预期。
// English: TestCloudPayCreatePaymentPostsSignedPayload verifies the related behavior.
func TestCloudPayCreatePaymentPostsSignedPayload(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["mch_id"] != "mch-1" {
			t.Fatalf("mch_id = %v, want mch-1", payload["mch_id"])
		}
		if payload["out_trade_no"] != "order-1" {
			t.Fatalf("out_trade_no = %v, want order-1", payload["out_trade_no"])
		}
		if payload["sign"] == "" {
			t.Fatal("missing sign")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"trade_no":"cp-1","pay_url":"https://pay.example/order-1","qr_code":"qr-data"}}`))
	}))
	defer server.Close()

	ch := NewCloudPayChannel("mch-1", "secret", server.URL, "https://notify.example/callback")
	result, err := ch.CreatePayment(context.Background(), "order-1", "subject", 100, "qrcode", "", "", "")
	if err != nil {
		t.Fatalf("CreatePayment returned error: %v", err)
	}
	if gotPath != "/payments" {
		t.Fatalf("path = %q, want /payments", gotPath)
	}
	if result.TradeNo != "cp-1" || result.PayURL == "" || result.QrCode != "qr-data" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

// 中文：TestCloudPayVerifyCallbackChecksSignature 验证相关行为符合预期。
// English: TestCloudPayVerifyCallbackChecksSignature verifies the related behavior.
func TestCloudPayVerifyCallbackChecksSignature(t *testing.T) {
	ch := NewCloudPayChannel("mch-1", "secret", "https://gateway.example", "")
	fields := map[string]string{
		"out_trade_no": "order-1",
		"trade_no":     "cp-1",
		"status":       "success",
		"amount":       "100",
	}
	signature := gatewaySign("secret", fields)
	body := `{"out_trade_no":"order-1","trade_no":"cp-1","status":"success","amount":100,"sign":"` + signature + `"}`

	result, err := ch.VerifyCallback(context.Background(), []byte(body))
	if err != nil {
		t.Fatalf("VerifyCallback returned error: %v", err)
	}
	if result.OutTradeNo != "order-1" || result.TradeNo != "cp-1" || result.Status != "paid" || result.Amount != 100 {
		t.Fatalf("unexpected callback result: %+v", result)
	}

	if _, err := ch.VerifyCallback(context.Background(), []byte(strings.Replace(body, signature, "BAD", 1))); err == nil {
		t.Fatal("VerifyCallback returned nil for bad signature")
	}
}

// 中文：TestCloudPayClosePaymentPostsSignedPayload 验证相关行为符合预期。
// English: TestCloudPayClosePaymentPostsSignedPayload verifies the related behavior.
func TestCloudPayClosePaymentPostsSignedPayload(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["mch_id"] != "mch-1" {
			t.Fatalf("mch_id = %v, want mch-1", payload["mch_id"])
		}
		if payload["out_trade_no"] != "order-close-1" {
			t.Fatalf("out_trade_no = %v, want order-close-1", payload["out_trade_no"])
		}
		if payload["sign"] == "" {
			t.Fatal("missing sign")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	ch := NewCloudPayChannel("mch-1", "secret", server.URL, "")
	if err := ch.ClosePayment(context.Background(), "order-close-1"); err != nil {
		t.Fatalf("ClosePayment returned error: %v", err)
	}
	if gotPath != "/payments/close" {
		t.Fatalf("path = %q, want /payments/close", gotPath)
	}
}

// 中文：TestDigitalRMBRefundAndQuery 验证相关行为符合预期。
// English: TestDigitalRMBRefundAndQuery verifies the related behavior.
func TestDigitalRMBRefundAndQuery(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/refunds":
			_, _ = w.Write([]byte(`{"success":true,"data":{"refund_no":"dr-refund-1","status":"success"}}`))
		case "/payments/query":
			_, _ = w.Write([]byte(`{"success":true,"data":{"out_trade_no":"order-2","trade_no":"dr-1","status":"paid","amount":299}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ch := NewDigitalRMBChannel("app-1", "merchant-1", "secret", server.URL, "wallet-1", "")
	refund, err := ch.Refund(context.Background(), "order-2", "refund-2", 299, 100, "partial")
	if err != nil {
		t.Fatalf("Refund returned error: %v", err)
	}
	if refund.RefundNo != "dr-refund-1" || refund.Status != "success" {
		t.Fatalf("unexpected refund result: %+v", refund)
	}

	query, err := ch.QueryPayment(context.Background(), "order-2")
	if err != nil {
		t.Fatalf("QueryPayment returned error: %v", err)
	}
	if query.OutTradeNo != "order-2" || query.TradeNo != "dr-1" || query.Status != "paid" || query.Amount != 299 {
		t.Fatalf("unexpected query result: %+v", query)
	}
	if len(paths) != 2 || paths[0] != "/refunds" || paths[1] != "/payments/query" {
		t.Fatalf("paths = %+v", paths)
	}
}

// 中文：TestDigitalRMBClosePaymentPostsSignedPayload 验证相关行为符合预期。
// English: TestDigitalRMBClosePaymentPostsSignedPayload verifies the related behavior.
func TestDigitalRMBClosePaymentPostsSignedPayload(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload["app_id"] != "app-1" || payload["merchant_id"] != "merchant-1" || payload["wallet_id"] != "wallet-1" {
			t.Fatalf("unexpected identity payload: %+v", payload)
		}
		if payload["out_trade_no"] != "order-close-2" {
			t.Fatalf("out_trade_no = %v, want order-close-2", payload["out_trade_no"])
		}
		if payload["sign"] == "" {
			t.Fatal("missing sign")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	ch := NewDigitalRMBChannel("app-1", "merchant-1", "secret", server.URL, "wallet-1", "")
	if err := ch.ClosePayment(context.Background(), "order-close-2"); err != nil {
		t.Fatalf("ClosePayment returned error: %v", err)
	}
	if gotPath != "/payments/close" {
		t.Fatalf("path = %q, want /payments/close", gotPath)
	}
}
