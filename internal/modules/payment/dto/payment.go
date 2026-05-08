package dto

// 中文：CreatePayReq 定义当前包使用的数据结构或接口。
// English: CreatePayReq defines a data structure or interface used by this package.
// CreatePayReq 创建支付请求
type CreatePayReq struct {
	// 中文：OutTradeNo 保存当前结构中的配置或数据值。
	// English: OutTradeNo stores a configuration or data value for this struct.
	OutTradeNo string `json:"out_trade_no" binding:"required,min=1,max=64"`
	// 中文：Amount 保存当前结构中的配置或数据值。
	// English: Amount stores a configuration or data value for this struct.
	Amount int64 `json:"amount" binding:"required,min=1"`
	// 中文：Currency 保存当前结构中的配置或数据值。
	// English: Currency stores a configuration or data value for this struct.
	Currency string `json:"currency" binding:"omitempty,max=8"`
	// 中文：Subject 保存当前结构中的配置或数据值。
	// English: Subject stores a configuration or data value for this struct.
	Subject string `json:"subject" binding:"required,max=256"`
	// 中文：Channel 保存当前结构中的配置或数据值。
	// English: Channel stores a configuration or data value for this struct.
	Channel string `json:"channel" binding:"required,oneof=wechat alipay unionpay cloudpay digital_rmb stripe paypal"`
	// 中文：Scene 保存当前结构中的配置或数据值。
	// English: Scene stores a configuration or data value for this struct.
	Scene string `json:"scene" binding:"required,oneof=jsapi app h5 native qrcode"`
	// 中文：NotifyURL 保存当前结构中的配置或数据值。
	// English: NotifyURL stores a configuration or data value for this struct.
	NotifyURL string `json:"notify_url" binding:"omitempty,max=512"`
	// 中文：ReturnURL 保存当前结构中的配置或数据值。
	// English: ReturnURL stores a configuration or data value for this struct.
	ReturnURL string `json:"return_url" binding:"omitempty,max=512"`
	// 中文：OpenID 保存当前结构中的配置或数据值。
	// English: OpenID stores a configuration or data value for this struct.
	OpenID string `json:"open_id" binding:"omitempty,max=128"`
}

// 中文：RefundReq 定义当前包使用的数据结构或接口。
// English: RefundReq defines a data structure or interface used by this package.
// RefundReq 退款请求
type RefundReq struct {
	// 中文：OutTradeNo 保存当前结构中的配置或数据值。
	// English: OutTradeNo stores a configuration or data value for this struct.
	OutTradeNo string `json:"out_trade_no" binding:"required"`
	// 中文：OutRefundNo 保存当前结构中的配置或数据值。
	// English: OutRefundNo stores a configuration or data value for this struct.
	OutRefundNo string `json:"out_refund_no" binding:"required"`
	// 中文：TotalAmount 保存当前结构中的配置或数据值。
	// English: TotalAmount stores a configuration or data value for this struct.
	TotalAmount int64 `json:"total_amount" binding:"required,min=1"`
	// 中文：RefundAmount 保存当前结构中的配置或数据值。
	// English: RefundAmount stores a configuration or data value for this struct.
	RefundAmount int64 `json:"refund_amount" binding:"required,min=1"`
	// 中文：Reason 保存当前结构中的配置或数据值。
	// English: Reason stores a configuration or data value for this struct.
	Reason string `json:"reason" binding:"omitempty,max=256"`
}

// 中文：PayOrderResp 定义当前包使用的数据结构或接口。
// English: PayOrderResp defines a data structure or interface used by this package.
// PayOrderResp 支付订单响应
type PayOrderResp struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string `json:"id"`
	// 中文：OutTradeNo 保存当前结构中的配置或数据值。
	// English: OutTradeNo stores a configuration or data value for this struct.
	OutTradeNo string `json:"out_trade_no"`
	// 中文：TradeNo 保存当前结构中的配置或数据值。
	// English: TradeNo stores a configuration or data value for this struct.
	TradeNo string `json:"trade_no,omitempty"`
	// 中文：Channel 保存当前结构中的配置或数据值。
	// English: Channel stores a configuration or data value for this struct.
	Channel string `json:"channel"`
	// 中文：Scene 保存当前结构中的配置或数据值。
	// English: Scene stores a configuration or data value for this struct.
	Scene string `json:"scene"`
	// 中文：Amount 保存当前结构中的配置或数据值。
	// English: Amount stores a configuration or data value for this struct.
	Amount int64 `json:"amount"`
	// 中文：Currency 保存当前结构中的配置或数据值。
	// English: Currency stores a configuration or data value for this struct.
	Currency string `json:"currency"`
	// 中文：Subject 保存当前结构中的配置或数据值。
	// English: Subject stores a configuration or data value for this struct.
	Subject string `json:"subject"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `json:"status"`
	// 中文：PayURL 保存当前结构中的配置或数据值。
	// English: PayURL stores a configuration or data value for this struct.
	PayURL string `json:"pay_url,omitempty"`
	// 中文：QrCode 保存当前结构中的配置或数据值。
	// English: QrCode stores a configuration or data value for this struct.
	QrCode string `json:"qr_code,omitempty"`
	// 中文：PayParams 保存当前结构中的配置或数据值。
	// English: PayParams stores a configuration or data value for this struct.
	PayParams any `json:"pay_params,omitempty"`
}
