package channel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/plutov/paypal/v4"
)

// 中文：PayPalChannel 定义当前包使用的数据结构或接口。
// English: PayPalChannel defines a data structure or interface used by this package.
// PayPalChannel PayPal支付通道
type PayPalChannel struct {
	// 中文：ClientID 保存当前结构中的配置或数据值。
	// English: ClientID stores a configuration or data value for this struct.
	ClientID string
	// 中文：ClientSecret 保存当前结构中的配置或数据值。
	// English: ClientSecret stores a configuration or data value for this struct.
	ClientSecret string
	// 中文：Sandbox 保存当前结构中的配置或数据值。
	// English: Sandbox stores a configuration or data value for this struct.
	Sandbox bool
	// 中文：WebhookID 保存当前结构中的配置或数据值。
	// English: WebhookID stores a configuration or data value for this struct.
	WebhookID string
	// 中文：NotifyURL 保存当前结构中的配置或数据值。
	// English: NotifyURL stores a configuration or data value for this struct.
	NotifyURL string
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *paypal.Client
}

// 中文：NewPayPalChannel 创建并返回对应组件实例。
// English: NewPayPalChannel creates and returns the corresponding component instance.
func NewPayPalChannel(clientID, clientSecret string, sandbox bool, notifyURL string) *PayPalChannel {
	return &PayPalChannel{ClientID: clientID, ClientSecret: clientSecret, Sandbox: sandbox, NotifyURL: notifyURL}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (c *PayPalChannel) Name() string { return "paypal" }

// 中文：baseURL 执行当前包中的对应流程。
// English: baseURL executes the corresponding workflow in this package.
func (c *PayPalChannel) baseURL() string {
	if c.Sandbox {
		return paypal.APIBaseSandBox
	}
	return paypal.APIBaseLive
}

// 中文：initClient 执行当前包中的对应流程。
// English: initClient executes the corresponding workflow in this package.
// initClient 懒初始化PayPal客户端
func (c *PayPalChannel) initClient(ctx context.Context) error {
	if c.client != nil {
		return nil
	}

	client, err := paypal.NewClient(c.ClientID, c.ClientSecret, c.baseURL())
	if err != nil {
		return fmt.Errorf("paypal: init client: %w", err)
	}

	// 获取Access Token
	if _, err := client.GetAccessToken(ctx); err != nil {
		return fmt.Errorf("paypal: get access token: %w", err)
	}

	c.client = client
	return nil
}

// 中文：CreatePayment 执行当前包中的对应流程。
// English: CreatePayment executes the corresponding workflow in this package.
func (c *PayPalChannel) CreatePayment(ctx context.Context, outTradeNo, subject string, amount int64, scene, notifyURL, returnURL, openID string) (*PayResult, error) {
	if err := c.initClient(ctx); err != nil {
		return nil, err
	}

	// 金额：分转元
	amountUSD := fmt.Sprintf("%.2f", float64(amount)/100.0)

	successURL := returnURL
	if successURL == "" {
		successURL = c.NotifyURL + "/success"
	}
	cancelURL := c.NotifyURL + "/cancel"

	purchaseUnits := []paypal.PurchaseUnitRequest{
		{
			ReferenceID: outTradeNo,
			Amount: &paypal.PurchaseUnitAmount{
				Value:    amountUSD,
				Currency: "USD",
			},
			Description: subject,
			CustomID:    outTradeNo,
		},
	}

	appContext := &paypal.ApplicationContext{
		ReturnURL: successURL,
		CancelURL: cancelURL,
	}

	order, err := c.client.CreateOrder(ctx, paypal.OrderIntentCapture, purchaseUnits, nil, appContext)
	if err != nil {
		return nil, fmt.Errorf("paypal create order: %w", err)
	}

	// 提取approve链接
	payURL := ""
	for _, link := range order.Links {
		if link.Rel == "approve" {
			payURL = link.Href
			break
		}
	}

	return &PayResult{
		PayURL:  payURL,
		TradeNo: order.ID,
	}, nil
}

// 中文：VerifyCallback 执行当前包中的对应流程。
// English: VerifyCallback executes the corresponding workflow in this package.
func (c *PayPalChannel) VerifyCallback(ctx context.Context, rawData []byte) (*CallbackResult, error) {
	return parsePayPalCallback(rawData)
}

// 中文：VerifyCallbackWithRequest 执行当前包中的对应流程。
// English: VerifyCallbackWithRequest executes the corresponding workflow in this package.
// VerifyCallbackWithRequest 使用原始HTTP请求验证PayPal Webhook（推荐）
func (c *PayPalChannel) VerifyCallbackWithRequest(ctx context.Context, req *http.Request, rawData []byte) (*CallbackResult, error) {
	if err := c.initClient(ctx); err != nil {
		return nil, err
	}

	if c.WebhookID == "" {
		return nil, fmt.Errorf("paypal: webhook ID not configured")
	}
	if req == nil {
		return nil, fmt.Errorf("paypal: original http request is required")
	}

	req.Body = io.NopCloser(bytes.NewReader(rawData))
	verifyResp, err := c.client.VerifyWebhookSignature(ctx, req, c.WebhookID)
	if err != nil {
		return nil, fmt.Errorf("paypal: verify webhook signature: %w", err)
	}
	if !strings.EqualFold(verifyResp.VerificationStatus, "SUCCESS") {
		return nil, fmt.Errorf("paypal: webhook signature verification failed: %s", verifyResp.VerificationStatus)
	}

	return parsePayPalCallback(rawData)
}

// 中文：parsePayPalCallback 执行当前包中的对应流程。
// English: parsePayPalCallback executes the corresponding workflow in this package.
func parsePayPalCallback(rawData []byte) (*CallbackResult, error) {
	var event struct {
		EventType string `json:"event_type"`
		Resource  struct {
			ID        string `json:"id"`
			CustomID  string `json:"custom_id"`
			InvoiceID string `json:"invoice_id"`
			Status    string `json:"status"`
			Amount    struct {
				Value string `json:"value"`
			} `json:"amount"`
			PurchaseUnits []struct {
				CustomID string `json:"custom_id"`
				Payments struct {
					Captures []struct {
						ID     string `json:"id"`
						Status string `json:"status"`
						Amount struct {
							Value string `json:"value"`
						} `json:"amount"`
					} `json:"captures"`
				} `json:"payments"`
			} `json:"purchase_units"`
		} `json:"resource"`
	}
	if err := json.Unmarshal(rawData, &event); err != nil {
		return nil, fmt.Errorf("paypal: parse webhook event: %w", err)
	}

	outTradeNo := firstNonEmpty(event.Resource.CustomID, event.Resource.InvoiceID)
	tradeNo := event.Resource.ID
	amountValue := event.Resource.Amount.Value
	resourceStatus := event.Resource.Status

	for _, unit := range event.Resource.PurchaseUnits {
		if outTradeNo == "" {
			outTradeNo = unit.CustomID
		}
		for _, capture := range unit.Payments.Captures {
			if tradeNo == "" {
				tradeNo = capture.ID
			}
			if amountValue == "" {
				amountValue = capture.Amount.Value
			}
			if resourceStatus == "" {
				resourceStatus = capture.Status
			}
		}
	}
	if outTradeNo == "" {
		return nil, fmt.Errorf("paypal: callback missing merchant order id")
	}

	status := "failed"
	if event.EventType == "PAYMENT.CAPTURE.COMPLETED" ||
		event.EventType == "CHECKOUT.ORDER.COMPLETED" ||
		strings.EqualFold(resourceStatus, "COMPLETED") {
		status = "paid"
	}

	var amountCents int64
	if amountValue != "" {
		if f, err := strconv.ParseFloat(amountValue, 64); err == nil {
			amountCents = int64(f*100 + 0.5)
		}
	}

	return &CallbackResult{
		OutTradeNo: outTradeNo,
		TradeNo:    tradeNo,
		Status:     status,
		Amount:     amountCents,
		RawData:    rawData,
	}, nil
}

// 中文：firstNonEmpty 执行当前包中的对应流程。
// English: firstNonEmpty executes the corresponding workflow in this package.
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

// 中文：Refund 执行当前包中的对应流程。
// English: Refund executes the corresponding workflow in this package.
func (c *PayPalChannel) Refund(ctx context.Context, outTradeNo, outRefundNo string, totalAmount, refundAmount int64, reason string) (*RefundResult, error) {
	if err := c.initClient(ctx); err != nil {
		return nil, err
	}

	// PayPal退款需要capture ID，这里用outTradeNo作为capture ID
	amountUSD := fmt.Sprintf("%.2f", float64(refundAmount)/100.0)

	refundReq := paypal.RefundCaptureRequest{
		Amount: &paypal.Money{
			Value:    amountUSD,
			Currency: "USD",
		},
		NoteToPayer: reason,
	}

	refundResp, err := c.client.RefundCapture(ctx, outTradeNo, refundReq)
	if err != nil {
		return nil, fmt.Errorf("paypal refund: %w", err)
	}

	status := "refunding"
	if refundResp.Status == "COMPLETED" {
		status = "success"
	}

	return &RefundResult{
		RefundNo: refundResp.ID,
		Status:   status,
	}, nil
}

// 中文：QueryPayment 执行当前包中的对应流程。
// English: QueryPayment executes the corresponding workflow in this package.
func (c *PayPalChannel) QueryPayment(ctx context.Context, outTradeNo string) (*CallbackResult, error) {
	if err := c.initClient(ctx); err != nil {
		return nil, err
	}

	order, err := c.client.GetOrder(ctx, outTradeNo)
	if err != nil {
		return nil, fmt.Errorf("paypal query order: %w", err)
	}

	status := "pending"
	if order.Status == "COMPLETED" {
		status = "paid"
	} else if order.Status == "CANCELLED" || order.Status == "DENIED" {
		status = "failed"
	}

	var amountCents int64
	if len(order.PurchaseUnits) > 0 && order.PurchaseUnits[0].Amount != nil {
		if f, err := strconv.ParseFloat(order.PurchaseUnits[0].Amount.Value, 64); err == nil {
			amountCents = int64(f * 100)
		}
	}

	var customID string
	if len(order.PurchaseUnits) > 0 {
		customID = order.PurchaseUnits[0].CustomID
	}

	return &CallbackResult{
		OutTradeNo: customID,
		TradeNo:    order.ID,
		Status:     status,
		Amount:     amountCents,
	}, nil
}

// 中文：ClosePayment 执行当前包中的对应流程。
// English: ClosePayment executes the corresponding workflow in this package.
// SetWebhookID 设置Webhook ID
func (c *PayPalChannel) ClosePayment(ctx context.Context, outTradeNo string) error {
	if outTradeNo == "" {
		return fmt.Errorf("paypal: order id is required")
	}
	if err := c.initClient(ctx); err != nil {
		return err
	}

	order, err := c.client.GetOrder(ctx, outTradeNo)
	if err != nil {
		return fmt.Errorf("paypal query order before close: %w", err)
	}

	switch strings.ToUpper(order.Status) {
	case "COMPLETED":
		return fmt.Errorf("paypal: completed order cannot be closed")
	case "CANCELLED", "DENIED", "EXPIRED", "VOIDED":
		return nil
	}

	for _, authID := range paypalOrderAuthorizationIDs(order) {
		if _, err := c.client.VoidAuthorization(ctx, authID); err != nil {
			return fmt.Errorf("paypal void authorization %s: %w", authID, err)
		}
	}
	return nil
}

// 中文：CallbackSuccess 执行当前包中的对应流程。
// English: CallbackSuccess executes the corresponding workflow in this package.
func (c *PayPalChannel) CallbackSuccess() any { return "success" }

// 中文：CallbackFail 执行当前包中的对应流程。
// English: CallbackFail executes the corresponding workflow in this package.
func (c *PayPalChannel) CallbackFail() any { return "fail" }

// 中文：SetWebhookID 执行当前包中的对应流程。
// English: SetWebhookID executes the corresponding workflow in this package.
func (c *PayPalChannel) SetWebhookID(id string) {
	c.WebhookID = id
}

// 中文：paypalOrderAuthorizationIDs 执行当前包中的对应流程。
// English: paypalOrderAuthorizationIDs executes the corresponding workflow in this package.
func paypalOrderAuthorizationIDs(order *paypal.Order) []string {
	if order == nil {
		return nil
	}
	seen := make(map[string]struct{})
	ids := make([]string, 0)
	for _, unit := range order.PurchaseUnits {
		if unit.Payments == nil {
			continue
		}
		for _, authorization := range unit.Payments.Authorizations {
			if authorization.ID == "" {
				continue
			}
			if _, ok := seen[authorization.ID]; ok {
				continue
			}
			seen[authorization.ID] = struct{}{}
			ids = append(ids, authorization.ID)
		}
	}
	return ids
}
