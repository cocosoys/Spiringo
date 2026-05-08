package channel

import (
	"context"
	"strings"
	"testing"
)

// 中文：TestUnionPayClosePaymentPostsVoidTransaction 验证相关行为符合预期。
// English: TestUnionPayClosePaymentPostsVoidTransaction verifies the related behavior.
func TestUnionPayClosePaymentPostsVoidTransaction(t *testing.T) {
	ch := &UnionPayChannel{
		MchID:     "merchant-1",
		APIKey:    "secret",
		NotifyURL: "https://example.test/notify",
		Sandbox:   true,
	}

	var endpoint string
	var params map[string]string
	oldPost := unionPayHTTPPost
	unionPayHTTPPost = func(_ context.Context, gotEndpoint string, gotParams map[string]string) (map[string]string, error) {
		endpoint = gotEndpoint
		params = gotParams
		return map[string]string{"respCode": "00", "queryId": "void-1"}, nil
	}
	t.Cleanup(func() {
		unionPayHTTPPost = oldPost
	})

	if err := ch.ClosePayment(context.Background(), "QUERY-123"); err != nil {
		t.Fatalf("ClosePayment returned error: %v", err)
	}
	if endpoint != "https://gateway.test.95516.com/gateway/api/backTransReq.do" {
		t.Fatalf("endpoint = %q", endpoint)
	}
	if params["txnType"] != "31" || params["txnSubType"] != "00" || params["origQryId"] != "QUERY-123" {
		t.Fatalf("unexpected params: %v", params)
	}
	if params["signature"] == "" {
		t.Fatalf("signature is empty")
	}
	if !strings.HasPrefix(params["orderId"], "V") || len(params["orderId"]) > 40 {
		t.Fatalf("orderId = %q", params["orderId"])
	}
}

// 中文：TestUnionPayClosePaymentReturnsGatewayFailure 验证相关行为符合预期。
// English: TestUnionPayClosePaymentReturnsGatewayFailure verifies the related behavior.
func TestUnionPayClosePaymentReturnsGatewayFailure(t *testing.T) {
	ch := &UnionPayChannel{MchID: "merchant-1", APIKey: "secret", Sandbox: true}

	oldPost := unionPayHTTPPost
	unionPayHTTPPost = func(context.Context, string, map[string]string) (map[string]string, error) {
		return map[string]string{"respCode": "34", "respMsg": "denied"}, nil
	}
	t.Cleanup(func() {
		unionPayHTTPPost = oldPost
	})

	err := ch.ClosePayment(context.Background(), "QUERY-123")
	if err == nil || !strings.Contains(err.Error(), "34 denied") {
		t.Fatalf("ClosePayment error = %v, want gateway failure", err)
	}
}
