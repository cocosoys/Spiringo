package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// 中文：CloudPayChannel 定义当前包使用的数据结构或接口。
// English: CloudPayChannel defines a data structure or interface used by this package.
type CloudPayChannel struct {
	// 中文：MchID 保存当前结构中的配置或数据值。
	// English: MchID stores a configuration or data value for this struct.
	MchID string
	// 中文：APIKey 保存当前结构中的配置或数据值。
	// English: APIKey stores a configuration or data value for this struct.
	APIKey string
	// 中文：GatewayURL 保存当前结构中的配置或数据值。
	// English: GatewayURL stores a configuration or data value for this struct.
	GatewayURL string
	// 中文：NotifyURL 保存当前结构中的配置或数据值。
	// English: NotifyURL stores a configuration or data value for this struct.
	NotifyURL string
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client gatewayHTTPClient
}

// 中文：NewCloudPayChannel 创建并返回对应组件实例。
// English: NewCloudPayChannel creates and returns the corresponding component instance.
func NewCloudPayChannel(mchID, apiKey, gatewayURL, notifyURL string) *CloudPayChannel {
	return &CloudPayChannel{
		MchID:      mchID,
		APIKey:     apiKey,
		GatewayURL: gatewayURL,
		NotifyURL:  notifyURL,
		client:     http.DefaultClient,
	}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (c *CloudPayChannel) Name() string { return "cloudpay" }

// 中文：CreatePayment 执行当前包中的对应流程。
// English: CreatePayment executes the corresponding workflow in this package.
func (c *CloudPayChannel) CreatePayment(ctx context.Context, outTradeNo, subject string, amount int64, scene, notifyURL, returnURL, openID string) (*PayResult, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	endpoint, err := gatewayEndpoint(c.GatewayURL, "/payments")
	if err != nil {
		return nil, fmt.Errorf("cloudpay: %w", err)
	}

	notify := notifyURL
	if notify == "" {
		notify = c.NotifyURL
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signValues := map[string]string{
		"mch_id":       c.MchID,
		"out_trade_no": outTradeNo,
		"subject":      subject,
		"amount":       strconv.FormatInt(amount, 10),
		"scene":        scene,
		"notify_url":   notify,
		"return_url":   returnURL,
		"openid":       openID,
		"timestamp":    timestamp,
	}

	payload := map[string]any{
		"mch_id":       c.MchID,
		"out_trade_no": outTradeNo,
		"subject":      subject,
		"amount":       amount,
		"currency":     "CNY",
		"scene":        scene,
		"notify_url":   notify,
		"return_url":   returnURL,
		"openid":       openID,
		"timestamp":    timestamp,
		"sign":         gatewaySign(c.APIKey, signValues),
	}

	var resp cloudPayResponse
	if err := gatewayPostJSON(ctx, c.client, endpoint, payload, &resp); err != nil {
		return nil, fmt.Errorf("cloudpay create payment: %w", err)
	}
	if !gatewayResponseOK(resp.Success, resp.Code) {
		return nil, fmt.Errorf("cloudpay create payment: code=%s, message=%s", resp.Code, resp.Message)
	}

	data := resp.paymentData()
	return &PayResult{
		TradeNo:  data.TradeNo,
		PayURL:   data.PayURL,
		QrCode:   data.QrCode,
		PrepayID: data.PrepayID,
		Params:   data.Params,
	}, nil
}

// 中文：VerifyCallback 执行当前包中的对应流程。
// English: VerifyCallback executes the corresponding workflow in this package.
func (c *CloudPayChannel) VerifyCallback(_ context.Context, rawData []byte) (*CallbackResult, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	fields, signature, err := gatewaySignedFields(rawData)
	if err != nil {
		return nil, fmt.Errorf("cloudpay: parse callback: %w", err)
	}
	if !gatewayVerifySignature(c.APIKey, fields, signature) {
		return nil, fmt.Errorf("cloudpay: signature mismatch")
	}

	return &CallbackResult{
		OutTradeNo: gatewayFirst(fields, "out_trade_no", "merchant_order_no", "order_id"),
		TradeNo:    gatewayFirst(fields, "trade_no", "transaction_id", "channel_trade_no"),
		Status:     gatewayPaymentStatus(gatewayFirst(fields, "status", "pay_status", "trade_status")),
		Amount:     gatewayParseAmount(gatewayFirst(fields, "amount", "total_amount", "paid_amount")),
		RawData:    rawData,
	}, nil
}

// 中文：Refund 执行当前包中的对应流程。
// English: Refund executes the corresponding workflow in this package.
func (c *CloudPayChannel) Refund(ctx context.Context, outTradeNo, outRefundNo string, totalAmount, refundAmount int64, reason string) (*RefundResult, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	endpoint, err := gatewayEndpoint(c.GatewayURL, "/refunds")
	if err != nil {
		return nil, fmt.Errorf("cloudpay: %w", err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signValues := map[string]string{
		"mch_id":        c.MchID,
		"out_trade_no":  outTradeNo,
		"out_refund_no": outRefundNo,
		"total_amount":  strconv.FormatInt(totalAmount, 10),
		"refund_amount": strconv.FormatInt(refundAmount, 10),
		"reason":        reason,
		"timestamp":     timestamp,
	}

	payload := map[string]any{
		"mch_id":        c.MchID,
		"out_trade_no":  outTradeNo,
		"out_refund_no": outRefundNo,
		"total_amount":  totalAmount,
		"refund_amount": refundAmount,
		"reason":        reason,
		"timestamp":     timestamp,
		"sign":          gatewaySign(c.APIKey, signValues),
	}

	var resp cloudPayResponse
	if err := gatewayPostJSON(ctx, c.client, endpoint, payload, &resp); err != nil {
		return nil, fmt.Errorf("cloudpay refund: %w", err)
	}
	if !gatewayResponseOK(resp.Success, resp.Code) {
		return nil, fmt.Errorf("cloudpay refund: code=%s, message=%s", resp.Code, resp.Message)
	}

	data := resp.refundData()
	return &RefundResult{
		RefundNo: data.RefundNo,
		Status:   gatewayRefundStatus(data.Status),
	}, nil
}

// 中文：QueryPayment 执行当前包中的对应流程。
// English: QueryPayment executes the corresponding workflow in this package.
func (c *CloudPayChannel) QueryPayment(ctx context.Context, outTradeNo string) (*CallbackResult, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	endpoint, err := gatewayEndpoint(c.GatewayURL, "/payments/query")
	if err != nil {
		return nil, fmt.Errorf("cloudpay: %w", err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signValues := map[string]string{
		"mch_id":       c.MchID,
		"out_trade_no": outTradeNo,
		"timestamp":    timestamp,
	}
	payload := map[string]any{
		"mch_id":       c.MchID,
		"out_trade_no": outTradeNo,
		"timestamp":    timestamp,
		"sign":         gatewaySign(c.APIKey, signValues),
	}

	var resp cloudPayResponse
	if err := gatewayPostJSON(ctx, c.client, endpoint, payload, &resp); err != nil {
		return nil, fmt.Errorf("cloudpay query: %w", err)
	}
	if !gatewayResponseOK(resp.Success, resp.Code) {
		return nil, fmt.Errorf("cloudpay query: code=%s, message=%s", resp.Code, resp.Message)
	}

	data := resp.paymentData()
	return &CallbackResult{
		OutTradeNo: firstNonEmptyString(data.OutTradeNo, outTradeNo),
		TradeNo:    data.TradeNo,
		Status:     gatewayPaymentStatus(data.Status),
		Amount:     data.Amount,
	}, nil
}

// 中文：ClosePayment 执行当前包中的对应流程。
// English: ClosePayment executes the corresponding workflow in this package.
func (c *CloudPayChannel) ClosePayment(ctx context.Context, outTradeNo string) error {
	if err := c.validate(); err != nil {
		return err
	}
	if outTradeNo == "" {
		return fmt.Errorf("cloudpay: out_trade_no is required")
	}
	endpoint, err := gatewayEndpoint(c.GatewayURL, "/payments/close")
	if err != nil {
		return fmt.Errorf("cloudpay: %w", err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signValues := map[string]string{
		"mch_id":       c.MchID,
		"out_trade_no": outTradeNo,
		"timestamp":    timestamp,
	}
	payload := map[string]any{
		"mch_id":       c.MchID,
		"out_trade_no": outTradeNo,
		"timestamp":    timestamp,
		"sign":         gatewaySign(c.APIKey, signValues),
	}

	var resp cloudPayResponse
	if err := gatewayPostJSON(ctx, c.client, endpoint, payload, &resp); err != nil {
		return fmt.Errorf("cloudpay close payment: %w", err)
	}
	if !gatewayResponseOK(resp.Success, resp.Code) {
		return fmt.Errorf("cloudpay close payment: code=%s, message=%s", resp.Code, resp.Message)
	}
	return nil
}

// 中文：CallbackSuccess 执行当前包中的对应流程。
// English: CallbackSuccess executes the corresponding workflow in this package.
func (c *CloudPayChannel) CallbackSuccess() any { return "success" }

// 中文：CallbackFail 执行当前包中的对应流程。
// English: CallbackFail executes the corresponding workflow in this package.
func (c *CloudPayChannel) CallbackFail() any { return "fail" }

// 中文：validate 执行当前包中的对应流程。
// English: validate executes the corresponding workflow in this package.
func (c *CloudPayChannel) validate() error {
	if c.MchID == "" {
		return fmt.Errorf("cloudpay: mch_id is required")
	}
	if c.APIKey == "" {
		return fmt.Errorf("cloudpay: api_key is required")
	}
	if c.GatewayURL == "" {
		return fmt.Errorf("cloudpay: gateway_url is required")
	}
	return nil
}

// 中文：cloudPayResponse 定义当前包使用的数据结构或接口。
// English: cloudPayResponse defines a data structure or interface used by this package.
type cloudPayResponse struct {
	// 中文：Code 保存当前结构中的配置或数据值。
	// English: Code stores a configuration or data value for this struct.
	Code string `json:"code"`
	// 中文：Message 保存当前结构中的配置或数据值。
	// English: Message stores a configuration or data value for this struct.
	Message string `json:"message"`
	// 中文：Success 保存当前结构中的配置或数据值。
	// English: Success stores a configuration or data value for this struct.
	Success bool `json:"success"`
	// 中文：Data 保存当前结构中的配置或数据值。
	// English: Data stores a configuration or data value for this struct.
	Data json.RawMessage `json:"data"`
	// 中文：TradeNo 保存当前结构中的配置或数据值。
	// English: TradeNo stores a configuration or data value for this struct.
	TradeNo string `json:"trade_no"`
	// 中文：PayURL 保存当前结构中的配置或数据值。
	// English: PayURL stores a configuration or data value for this struct.
	PayURL string `json:"pay_url"`
	// 中文：QrCode 保存当前结构中的配置或数据值。
	// English: QrCode stores a configuration or data value for this struct.
	QrCode string `json:"qr_code"`
	// 中文：PrepayID 保存当前结构中的配置或数据值。
	// English: PrepayID stores a configuration or data value for this struct.
	PrepayID string `json:"prepay_id"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `json:"status"`
	// 中文：Amount 保存当前结构中的配置或数据值。
	// English: Amount stores a configuration or data value for this struct.
	Amount int64 `json:"amount"`
	// 中文：RefundNo 保存当前结构中的配置或数据值。
	// English: RefundNo stores a configuration or data value for this struct.
	RefundNo string `json:"refund_no"`
	// 中文：OutTradeNo 保存当前结构中的配置或数据值。
	// English: OutTradeNo stores a configuration or data value for this struct.
	OutTradeNo string `json:"out_trade_no"`
	// 中文：Params 保存当前结构中的配置或数据值。
	// English: Params stores a configuration or data value for this struct.
	Params map[string]any `json:"params"`
}

// 中文：cloudPayPaymentData 定义当前包使用的数据结构或接口。
// English: cloudPayPaymentData defines a data structure or interface used by this package.
type cloudPayPaymentData struct {
	// 中文：OutTradeNo 保存当前结构中的配置或数据值。
	// English: OutTradeNo stores a configuration or data value for this struct.
	OutTradeNo string `json:"out_trade_no"`
	// 中文：TradeNo 保存当前结构中的配置或数据值。
	// English: TradeNo stores a configuration or data value for this struct.
	TradeNo string `json:"trade_no"`
	// 中文：PayURL 保存当前结构中的配置或数据值。
	// English: PayURL stores a configuration or data value for this struct.
	PayURL string `json:"pay_url"`
	// 中文：QrCode 保存当前结构中的配置或数据值。
	// English: QrCode stores a configuration or data value for this struct.
	QrCode string `json:"qr_code"`
	// 中文：PrepayID 保存当前结构中的配置或数据值。
	// English: PrepayID stores a configuration or data value for this struct.
	PrepayID string `json:"prepay_id"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `json:"status"`
	// 中文：Amount 保存当前结构中的配置或数据值。
	// English: Amount stores a configuration or data value for this struct.
	Amount int64 `json:"amount"`
	// 中文：Params 保存当前结构中的配置或数据值。
	// English: Params stores a configuration or data value for this struct.
	Params map[string]any `json:"params"`
}

// 中文：cloudPayRefundData 定义当前包使用的数据结构或接口。
// English: cloudPayRefundData defines a data structure or interface used by this package.
type cloudPayRefundData struct {
	// 中文：RefundNo 保存当前结构中的配置或数据值。
	// English: RefundNo stores a configuration or data value for this struct.
	RefundNo string `json:"refund_no"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `json:"status"`
}

// 中文：paymentData 执行当前包中的对应流程。
// English: paymentData executes the corresponding workflow in this package.
func (r cloudPayResponse) paymentData() cloudPayPaymentData {
	data := cloudPayPaymentData{
		OutTradeNo: r.OutTradeNo,
		TradeNo:    r.TradeNo,
		PayURL:     r.PayURL,
		QrCode:     r.QrCode,
		PrepayID:   r.PrepayID,
		Status:     r.Status,
		Amount:     r.Amount,
		Params:     r.Params,
	}
	if len(r.Data) > 0 {
		_ = json.Unmarshal(r.Data, &data)
	}
	return data
}

// 中文：refundData 执行当前包中的对应流程。
// English: refundData executes the corresponding workflow in this package.
func (r cloudPayResponse) refundData() cloudPayRefundData {
	data := cloudPayRefundData{RefundNo: r.RefundNo, Status: r.Status}
	if len(r.Data) > 0 {
		_ = json.Unmarshal(r.Data, &data)
	}
	return data
}

// 中文：firstNonEmptyString 执行当前包中的对应流程。
// English: firstNonEmptyString executes the corresponding workflow in this package.
func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
