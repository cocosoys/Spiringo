package service

import (
	"testing"

	"github.com/spiringo/spiringo/internal/modules/payment/model"
)

// 中文：TestPaymentCloseIDUsesTradeNoForUnionPay 验证相关行为符合预期。
// English: TestPaymentCloseIDUsesTradeNoForUnionPay verifies the related behavior.
func TestPaymentCloseIDUsesTradeNoForUnionPay(t *testing.T) {
	got := paymentCloseID(&model.PaymentOrder{
		OutTradeNo: "order-close-001",
		TradeNo:    "query-001",
		Channel:    "unionpay",
	})
	if got != "query-001" {
		t.Fatalf("paymentCloseID = %q, want query-001", got)
	}
}
