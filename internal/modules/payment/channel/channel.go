package channel

import (
	"context"
	"net/http"
	"sort"
)

// 中文：PayResult 定义当前包使用的数据结构或接口。
// English: PayResult defines a data structure or interface used by this package.
// PayResult is the normalized result returned after creating a payment.
type PayResult struct {
	// 中文：TradeNo 保存当前结构中的配置或数据值。
	// English: TradeNo stores a configuration or data value for this struct.
	TradeNo string `json:"trade_no,omitempty"`
	// 中文：PayURL 保存当前结构中的配置或数据值。
	// English: PayURL stores a configuration or data value for this struct.
	PayURL string `json:"pay_url,omitempty"`
	// 中文：QrCode 保存当前结构中的配置或数据值。
	// English: QrCode stores a configuration or data value for this struct.
	QrCode string `json:"qr_code,omitempty"`
	// 中文：PrepayID 保存当前结构中的配置或数据值。
	// English: PrepayID stores a configuration or data value for this struct.
	PrepayID string `json:"prepay_id,omitempty"`
	// 中文：Params 保存当前结构中的配置或数据值。
	// English: Params stores a configuration or data value for this struct.
	Params any `json:"params,omitempty"`
}

// 中文：RefundResult 定义当前包使用的数据结构或接口。
// English: RefundResult defines a data structure or interface used by this package.
// RefundResult is the normalized result returned after creating a refund.
type RefundResult struct {
	// 中文：RefundNo 保存当前结构中的配置或数据值。
	// English: RefundNo stores a configuration or data value for this struct.
	RefundNo string `json:"refund_no,omitempty"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `json:"status"`
}

// 中文：CallbackResult 定义当前包使用的数据结构或接口。
// English: CallbackResult defines a data structure or interface used by this package.
// CallbackResult is the normalized payment state parsed from a channel callback.
type CallbackResult struct {
	// 中文：OutTradeNo 保存当前结构中的配置或数据值。
	// English: OutTradeNo stores a configuration or data value for this struct.
	OutTradeNo string `json:"out_trade_no"`
	// 中文：TradeNo 保存当前结构中的配置或数据值。
	// English: TradeNo stores a configuration or data value for this struct.
	TradeNo string `json:"trade_no"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `json:"status"`
	// 中文：Amount 保存当前结构中的配置或数据值。
	// English: Amount stores a configuration or data value for this struct.
	Amount int64 `json:"amount"`
	// 中文：RawData 保存当前结构中的配置或数据值。
	// English: RawData stores a configuration or data value for this struct.
	RawData []byte `json:"raw_data"`
}

// 中文：Channel 定义当前包使用的数据结构或接口。
// English: Channel defines a data structure or interface used by this package.
// Channel defines the payment channel contract used by the payment service.
type Channel interface {
	// 中文：Name 声明该接口需要实现的行为。
	// English: Name declares behavior required by this interface.
	Name() string
	// 中文：CreatePayment 声明该接口需要实现的行为。
	// English: CreatePayment declares behavior required by this interface.
	CreatePayment(ctx context.Context, outTradeNo, subject string, amount int64, scene, notifyURL, returnURL, openID string) (*PayResult, error)
	// 中文：VerifyCallback 声明该接口需要实现的行为。
	// English: VerifyCallback declares behavior required by this interface.
	VerifyCallback(ctx context.Context, rawData []byte) (*CallbackResult, error)
	// 中文：Refund 声明该接口需要实现的行为。
	// English: Refund declares behavior required by this interface.
	Refund(ctx context.Context, outTradeNo, outRefundNo string, totalAmount, refundAmount int64, reason string) (*RefundResult, error)
	// 中文：QueryPayment 声明该接口需要实现的行为。
	// English: QueryPayment declares behavior required by this interface.
	QueryPayment(ctx context.Context, outTradeNo string) (*CallbackResult, error)
	// 中文：ClosePayment 声明该接口需要实现的行为。
	// English: ClosePayment declares behavior required by this interface.
	ClosePayment(ctx context.Context, outTradeNo string) error
	// 中文：CallbackSuccess 声明该接口需要实现的行为。
	// English: CallbackSuccess declares behavior required by this interface.
	CallbackSuccess() any
	// 中文：CallbackFail 声明该接口需要实现的行为。
	// English: CallbackFail declares behavior required by this interface.
	CallbackFail() any
}

// 中文：HTTPCallbackVerifier 定义当前包使用的数据结构或接口。
// English: HTTPCallbackVerifier defines a data structure or interface used by this package.
// HTTPCallbackVerifier is implemented by channels that need original webhook
// headers or request metadata for signature verification.
type HTTPCallbackVerifier interface {
	// 中文：VerifyCallbackWithRequest 声明该接口需要实现的行为。
	// English: VerifyCallbackWithRequest declares behavior required by this interface.
	VerifyCallbackWithRequest(ctx context.Context, req *http.Request, rawData []byte) (*CallbackResult, error)
}

// 中文：Registry 定义当前包使用的数据结构或接口。
// English: Registry defines a data structure or interface used by this package.
// Registry stores enabled payment channels by channel name.
type Registry struct {
	// 中文：channels 保存当前结构中的配置或数据值。
	// English: channels stores a configuration or data value for this struct.
	channels map[string]Channel
}

// 中文：NewRegistry 创建并返回对应组件实例。
// English: NewRegistry creates and returns the corresponding component instance.
// NewRegistry creates a channel registry.
func NewRegistry() *Registry {
	return &Registry{channels: make(map[string]Channel)}
}

// 中文：Register 执行当前包中的对应流程。
// English: Register executes the corresponding workflow in this package.
// Register stores a channel by its Name. Nil channels are ignored.
func (r *Registry) Register(c Channel) {
	if r == nil || c == nil {
		return
	}
	r.channels[c.Name()] = c
}

// 中文：Get 执行当前包中的对应流程。
// English: Get executes the corresponding workflow in this package.
// Get returns a channel by name.
func (r *Registry) Get(name string) (Channel, bool) {
	if r == nil {
		return nil, false
	}
	c, ok := r.channels[name]
	return c, ok
}

// 中文：List 执行当前包中的对应流程。
// English: List executes the corresponding workflow in this package.
// List returns all registered channels sorted by channel name.
func (r *Registry) List() []Channel {
	if r == nil || len(r.channels) == 0 {
		return nil
	}
	names := make([]string, 0, len(r.channels))
	for name := range r.channels {
		names = append(names, name)
	}
	sort.Strings(names)

	channels := make([]Channel, 0, len(names))
	for _, name := range names {
		channels = append(channels, r.channels[name])
	}
	return channels
}
