package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/refund"
	"github.com/stripe/stripe-go/v82/webhook"
)

// 中文：StripeChannel 定义当前包使用的数据结构或接口。
// English: StripeChannel defines a data structure or interface used by this package.
// StripeChannel Stripe支付通道
type StripeChannel struct {
	// 中文：SecretKey 保存当前结构中的配置或数据值。
	// English: SecretKey stores a configuration or data value for this struct.
	SecretKey string
	// 中文：WebhookSecret 保存当前结构中的配置或数据值。
	// English: WebhookSecret stores a configuration or data value for this struct.
	WebhookSecret string
	// 中文：NotifyURL 保存当前结构中的配置或数据值。
	// English: NotifyURL stores a configuration or data value for this struct.
	NotifyURL string
}

// 中文：NewStripeChannel 创建并返回对应组件实例。
// English: NewStripeChannel creates and returns the corresponding component instance.
func NewStripeChannel(secretKey, notifyURL string) *StripeChannel {
	return &StripeChannel{SecretKey: secretKey, NotifyURL: notifyURL}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (c *StripeChannel) Name() string { return "stripe" }

// 中文：initStripe 执行当前包中的对应流程。
// English: initStripe executes the corresponding workflow in this package.
// initStripe 设置全局API Key
func (c *StripeChannel) initStripe() {
	stripe.Key = c.SecretKey
}

// 中文：CreatePayment 执行当前包中的对应流程。
// English: CreatePayment executes the corresponding workflow in this package.
func (c *StripeChannel) CreatePayment(ctx context.Context, outTradeNo, subject string, amount int64, scene, notifyURL, returnURL, openID string) (*PayResult, error) {
	c.initStripe()

	successURL := returnURL
	if successURL == "" {
		successURL = c.NotifyURL + "/success"
	}
	cancelURL := c.NotifyURL + "/cancel"

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(subject),
					},
					UnitAmount: stripe.Int64(amount),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Metadata: map[string]string{
			"out_trade_no": outTradeNo,
		},
	}

	s, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe create session: %w", err)
	}

	return &PayResult{
		PayURL:  s.URL,
		TradeNo: string(s.ID),
	}, nil
}

// 中文：VerifyCallback 执行当前包中的对应流程。
// English: VerifyCallback executes the corresponding workflow in this package.
func (c *StripeChannel) VerifyCallback(ctx context.Context, rawData []byte) (*CallbackResult, error) {
	if c.WebhookSecret == "" {
		return nil, fmt.Errorf("stripe: webhook secret not configured")
	}
	return c.VerifyCallbackWithHeader(rawData, "")
}

// 中文：VerifyCallbackWithRequest 执行当前包中的对应流程。
// English: VerifyCallbackWithRequest executes the corresponding workflow in this package.
func (c *StripeChannel) VerifyCallbackWithRequest(_ context.Context, req *http.Request, rawData []byte) (*CallbackResult, error) {
	if req == nil {
		return nil, fmt.Errorf("stripe: original http request is required")
	}
	return c.VerifyCallbackWithHeader(rawData, req.Header.Get("Stripe-Signature"))
}

// 中文：VerifyCallbackWithHeader 执行当前包中的对应流程。
// English: VerifyCallbackWithHeader executes the corresponding workflow in this package.
// VerifyCallbackWithHeader 使用HTTP header验签（推荐方式）
func (c *StripeChannel) VerifyCallbackWithHeader(rawData []byte, sigHeader string) (*CallbackResult, error) {
	if c.WebhookSecret == "" {
		return nil, fmt.Errorf("stripe: webhook secret not configured")
	}
	if sigHeader == "" {
		return nil, fmt.Errorf("stripe: missing Stripe-Signature header")
	}

	event, err := webhook.ConstructEvent(rawData, sigHeader, c.WebhookSecret)
	if err != nil {
		return nil, fmt.Errorf("stripe: verify webhook: %w", err)
	}

	var outTradeNo, tradeNo string
	var amount int64
	status := "failed"

	switch event.Type {
	case "checkout.session.completed":
		var cs stripe.CheckoutSession
		if err := unmarshalEventObject(event.Data.Object, &cs); err == nil {
			outTradeNo = cs.Metadata["out_trade_no"]
			tradeNo = string(cs.ID)
			if cs.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
				status = "paid"
			}
			amount = cs.AmountTotal
		}
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		if err := unmarshalEventObject(event.Data.Object, &pi); err == nil {
			outTradeNo = pi.Metadata["out_trade_no"]
			tradeNo = string(pi.ID)
			status = "paid"
			amount = pi.Amount
		}
	}

	return &CallbackResult{
		OutTradeNo: outTradeNo,
		TradeNo:    tradeNo,
		Status:     status,
		Amount:     amount,
		RawData:    rawData,
	}, nil
}

// 中文：Refund 执行当前包中的对应流程。
// English: Refund executes the corresponding workflow in this package.
func (c *StripeChannel) Refund(ctx context.Context, outTradeNo, outRefundNo string, totalAmount, refundAmount int64, reason string) (*RefundResult, error) {
	c.initStripe()

	params := &stripe.RefundParams{
		Metadata: map[string]string{
			"out_trade_no":  outTradeNo,
			"out_refund_no": outRefundNo,
		},
	}

	if refundAmount > 0 && refundAmount < totalAmount {
		params.Amount = stripe.Int64(refundAmount)
	}

	if outTradeNo != "" {
		params.PaymentIntent = stripe.String(outTradeNo)
	}

	r, err := refund.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe refund: %w", err)
	}

	status := "refunding"
	if r.Status == stripe.RefundStatusSucceeded {
		status = "success"
	}

	return &RefundResult{
		RefundNo: string(r.ID),
		Status:   status,
	}, nil
}

// 中文：QueryPayment 执行当前包中的对应流程。
// English: QueryPayment executes the corresponding workflow in this package.
func (c *StripeChannel) QueryPayment(ctx context.Context, outTradeNo string) (*CallbackResult, error) {
	c.initStripe()

	// 通过Checkout Session ID查询
	s, err := session.Get(outTradeNo, nil)
	if err != nil {
		// 尝试作为PaymentIntent查询
		pi, err2 := paymentintent.Get(outTradeNo, nil)
		if err2 != nil {
			return nil, fmt.Errorf("stripe query: session err=%v, pi err=%v", err, err2)
		}
		status := "pending"
		if pi.Status == stripe.PaymentIntentStatusSucceeded {
			status = "paid"
		} else if pi.Status == stripe.PaymentIntentStatusCanceled {
			status = "failed"
		}
		return &CallbackResult{
			OutTradeNo: pi.Metadata["out_trade_no"],
			TradeNo:    string(pi.ID),
			Status:     status,
			Amount:     pi.Amount,
		}, nil
	}

	status := "pending"
	if s.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
		status = "paid"
	}

	return &CallbackResult{
		OutTradeNo: s.Metadata["out_trade_no"],
		TradeNo:    string(s.ID),
		Status:     status,
		Amount:     s.AmountTotal,
	}, nil
}

// 中文：ClosePayment 执行当前包中的对应流程。
// English: ClosePayment executes the corresponding workflow in this package.
// SetWebhookSecret 设置Webhook签名密钥
func (c *StripeChannel) ClosePayment(_ context.Context, outTradeNo string) error {
	c.initStripe()
	if outTradeNo == "" {
		return fmt.Errorf("stripe: checkout session or payment intent id is required")
	}

	if strings.HasPrefix(outTradeNo, "cs_") {
		if _, err := session.Expire(outTradeNo, nil); err != nil {
			return fmt.Errorf("stripe expire checkout session: %w", err)
		}
		return nil
	}
	if strings.HasPrefix(outTradeNo, "pi_") {
		if _, err := paymentintent.Cancel(outTradeNo, nil); err != nil {
			return fmt.Errorf("stripe cancel payment intent: %w", err)
		}
		return nil
	}

	_, sessionErr := session.Expire(outTradeNo, nil)
	if sessionErr == nil {
		return nil
	}
	_, intentErr := paymentintent.Cancel(outTradeNo, nil)
	if intentErr == nil {
		return nil
	}
	return fmt.Errorf("stripe close payment: session err=%v, payment_intent err=%v", sessionErr, intentErr)
}

// 中文：CallbackSuccess 执行当前包中的对应流程。
// English: CallbackSuccess executes the corresponding workflow in this package.
func (c *StripeChannel) CallbackSuccess() any { return "success" }

// 中文：CallbackFail 执行当前包中的对应流程。
// English: CallbackFail executes the corresponding workflow in this package.
func (c *StripeChannel) CallbackFail() any { return "fail" }

// 中文：SetWebhookSecret 执行当前包中的对应流程。
// English: SetWebhookSecret executes the corresponding workflow in this package.
func (c *StripeChannel) SetWebhookSecret(secret string) {
	c.WebhookSecret = secret
}

// 中文：unmarshalEventObject 执行当前包中的对应流程。
// English: unmarshalEventObject executes the corresponding workflow in this package.
// unmarshalEventObject 将 webhook event data object (map[string]interface{}) 解析为目标结构体
func unmarshalEventObject(obj map[string]interface{}, target any) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
