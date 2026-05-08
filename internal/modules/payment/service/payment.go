package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/spiringo/spiringo/internal/core/event"
	"github.com/spiringo/spiringo/internal/modules/payment/channel"
	"github.com/spiringo/spiringo/internal/modules/payment/dto"
	"github.com/spiringo/spiringo/internal/modules/payment/model"
	"github.com/spiringo/spiringo/internal/modules/payment/repository"
	"github.com/spiringo/spiringo/pkg/types"
	"gorm.io/gorm"
)

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
// Config 支付服务配置
type Config struct {
	// 中文：Wechat 保存当前结构中的配置或数据值。
	// English: Wechat stores a configuration or data value for this struct.
	Wechat struct {
		Enabled  bool
		AppID    string
		MchID    string
		APIKey   string
		CertPath string
	}
	// 中文：Alipay 保存当前结构中的配置或数据值。
	// English: Alipay stores a configuration or data value for this struct.
	Alipay struct {
		Enabled    bool
		AppID      string
		PrivateKey string
		PublicKey  string
	}
	// 中文：Stripe 保存当前结构中的配置或数据值。
	// English: Stripe stores a configuration or data value for this struct.
	Stripe struct {
		Enabled   bool
		SecretKey string
	}
	// 中文：PayPal 保存当前结构中的配置或数据值。
	// English: PayPal stores a configuration or data value for this struct.
	PayPal struct {
		Enabled      bool
		ClientID     string
		ClientSecret string
		Sandbox      bool
	}
	// 中文：UnionPay 保存当前结构中的配置或数据值。
	// English: UnionPay stores a configuration or data value for this struct.
	UnionPay struct {
		Enabled bool
		MchID   string
		APIKey  string
	}
	// 中文：CloudPay 保存当前结构中的配置或数据值。
	// English: CloudPay stores a configuration or data value for this struct.
	CloudPay struct {
		Enabled bool
		MchID   string
		APIKey  string
	}
	// 中文：DigitalRMB 保存当前结构中的配置或数据值。
	// English: DigitalRMB stores a configuration or data value for this struct.
	DigitalRMB struct {
		Enabled    bool
		AppID      string
		MerchantID string
		APIKey     string
	}
	// 中文：DefaultNotifyURL 保存当前结构中的配置或数据值。
	// English: DefaultNotifyURL stores a configuration or data value for this struct.
	DefaultNotifyURL string
}

// 中文：PaymentService 定义当前包使用的数据结构或接口。
// English: PaymentService defines a data structure or interface used by this package.
// PaymentService 支付业务逻辑
type PaymentService struct {
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config Config
	// 中文：eventBus 保存当前结构中的配置或数据值。
	// English: eventBus stores a configuration or data value for this struct.
	eventBus *event.Bus
	// 中文：repo 保存当前结构中的配置或数据值。
	// English: repo stores a configuration or data value for this struct.
	repo *repository.PaymentRepository
	// 中文：chReg 保存当前结构中的配置或数据值。
	// English: chReg stores a configuration or data value for this struct.
	chReg *channel.Registry
}

// 中文：NewPaymentService 创建并返回对应组件实例。
// English: NewPaymentService creates and returns the corresponding component instance.
// NewPaymentService 创建支付服务
func NewPaymentService(config Config, eventBus *event.Bus, repo *repository.PaymentRepository, chReg *channel.Registry) *PaymentService {
	return &PaymentService{config: config, eventBus: eventBus, repo: repo, chReg: chReg}
}

// 中文：CreateOrder 执行当前包中的对应流程。
// English: CreateOrder executes the corresponding workflow in this package.
// CreateOrder 创建支付订单
func (s *PaymentService) CreateOrder(ctx context.Context, req dto.CreatePayReq) (*dto.PayOrderResp, error) {
	// 检查商户订单号唯一性
	existing, _ := s.repo.GetOrderByOutTradeNo(ctx, req.OutTradeNo)
	if existing != nil {
		return nil, types.ErrPaymentDuplicate
	}

	// 获取支付通道
	ch, ok := s.chReg.Get(req.Channel)
	if !ok {
		return nil, types.ErrPaymentChannel.WithMessagef("channel %s is not enabled", req.Channel)
	}

	currency := req.Currency
	if currency == "" {
		currency = "CNY"
	}

	order := &model.PaymentOrder{
		OutTradeNo: req.OutTradeNo,
		Channel:    req.Channel,
		Scene:      req.Scene,
		Amount:     req.Amount,
		Currency:   currency,
		Subject:    req.Subject,
		Status:     string(model.PayStatusPending),
		NotifyURL:  req.NotifyURL,
		ReturnURL:  req.ReturnURL,
		ExpiredAt:  timePtr(time.Now().Add(30 * time.Minute)),
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("create payment order: %w", err)
	}

	// 调用支付通道SDK创建预支付
	payResult, err := ch.CreatePayment(ctx, req.OutTradeNo, req.Subject, req.Amount, req.Scene, req.NotifyURL, req.ReturnURL, req.OpenID)
	if err != nil {
		return nil, types.ErrPaymentFailed.WithMessage(err.Error())
	}

	// 保存通道交易号
	if payResult.TradeNo != "" {
		order.TradeNo = payResult.TradeNo
		if err := s.repo.UpdateOrder(ctx, order); err != nil {
			return nil, fmt.Errorf("update order with trade_no: %w", err)
		}
	}

	resp := &dto.PayOrderResp{
		ID:         order.ID,
		OutTradeNo: order.OutTradeNo,
		TradeNo:    order.TradeNo,
		Channel:    order.Channel,
		Scene:      order.Scene,
		Amount:     order.Amount,
		Currency:   order.Currency,
		Subject:    order.Subject,
		Status:     order.Status,
		PayURL:     payResult.PayURL,
		QrCode:     payResult.QrCode,
		PayParams:  payResult.Params,
	}

	s.publishAsync(ctx, event.NewEvent(event.EventPaymentCreated, &event.PaymentEventPayload{
		OrderID:    order.ID,
		OutTradeNo: order.OutTradeNo,
		TradeNo:    order.TradeNo,
		Channel:    order.Channel,
		Amount:     order.Amount,
		Currency:   order.Currency,
		Subject:    order.Subject,
		TenantID:   order.TenantID,
	}))

	return resp, nil
}

// 中文：HandleCallback 执行当前包中的对应流程。
// English: HandleCallback executes the corresponding workflow in this package.
// HandleCallback 处理支付回调
func (s *PaymentService) HandleCallback(ctx context.Context, channelName string, rawData []byte) error {
	return s.HandleCallbackRequest(ctx, channelName, nil, rawData)
}

// 中文：HandleCallbackRequest 执行当前包中的对应流程。
// English: HandleCallbackRequest executes the corresponding workflow in this package.
// HandleCallbackRequest handles payment callbacks that require original request
// headers for signature verification.
func (s *PaymentService) HandleCallbackRequest(ctx context.Context, channelName string, req *http.Request, rawData []byte) error {
	// 记录回调日志
	callbackLog := &model.CallbackLog{
		Channel: channelName,
		RawData: string(rawData),
	}
	if err := s.repo.CreateCallbackLog(ctx, callbackLog); err != nil {
		return fmt.Errorf("save callback log: %w", err)
	}

	// 获取通道并验证签名
	ch, ok := s.chReg.Get(channelName)
	if !ok {
		return types.ErrPaymentChannel.WithMessagef("channel %s not found", channelName)
	}

	var result *channel.CallbackResult
	var err error
	if req != nil {
		if verifier, ok := ch.(channel.HTTPCallbackVerifier); ok {
			result, err = verifier.VerifyCallbackWithRequest(ctx, req, rawData)
		} else {
			result, err = ch.VerifyCallback(ctx, rawData)
		}
	} else {
		result, err = ch.VerifyCallback(ctx, rawData)
	}
	if err != nil {
		return types.ErrCallbackVerify.WithMessage(err.Error())
	}
	if result == nil {
		return types.ErrCallbackVerify.WithMessage("empty callback result")
	}
	if result.OutTradeNo == "" {
		return types.ErrCallbackVerify.WithMessage("callback missing merchant order id")
	}
	callbackLog.TradeNo = result.TradeNo

	// 查找并更新订单
	order, err := s.repo.GetOrderByOutTradeNo(ctx, result.OutTradeNo)
	if err != nil {
		return fmt.Errorf("order not found: %s", result.OutTradeNo)
	}
	if result.Amount > 0 && result.Amount != order.Amount {
		return types.ErrCallbackVerify.WithMessagef("callback amount mismatch: got %d, want %d", result.Amount, order.Amount)
	}
	if isTerminalPaymentStatus(order.Status) {
		callbackLog.Processed = true
		if err := s.repo.UpdateCallbackLog(ctx, callbackLog); err != nil {
			return fmt.Errorf("update duplicate callback log: %w", err)
		}
		return nil
	}

	if result.TradeNo != "" {
		order.TradeNo = result.TradeNo
	}

	var postCommitEvents []*event.Event
	now := time.Now()
	switch result.Status {
	case "paid":
		order.Status = string(model.PayStatusPaid)
		order.PaidAt = &now
		payload := paymentEventPayload(order)
		postCommitEvents = append(postCommitEvents,
			event.NewEvent(event.EventPaymentSuccess, payload),
			event.NewEvent(event.EventPaymentFulfillmentRequested, payload),
		)
	case "failed":
		order.Status = string(model.PayStatusFailed)
		postCommitEvents = append(postCommitEvents, event.NewEvent(event.EventPaymentFailed, paymentEventPayload(order)))
	}

	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		return fmt.Errorf("update order after callback: %w", err)
	}
	callbackLog.Processed = true
	if err := s.repo.UpdateCallbackLog(ctx, callbackLog); err != nil {
		return fmt.Errorf("update callback log: %w", err)
	}
	for _, postCommitEvent := range postCommitEvents {
		s.publishAsync(ctx, postCommitEvent)
	}

	return nil
}

// 中文：isTerminalPaymentStatus 执行当前包中的对应流程。
// English: isTerminalPaymentStatus executes the corresponding workflow in this package.
func isTerminalPaymentStatus(status string) bool {
	switch status {
	case string(model.PayStatusPaid), string(model.PayStatusRefunded), string(model.PayStatusClosed):
		return true
	default:
		return false
	}
}

// 中文：Refund 执行当前包中的对应流程。
// English: Refund executes the corresponding workflow in this package.
// Refund 发起退款
func (s *PaymentService) Refund(ctx context.Context, req dto.RefundReq) (*model.RefundOrder, error) {
	if req.RefundAmount > req.TotalAmount {
		return nil, types.ErrRefundAmount
	}
	existingRefund, err := s.repo.GetRefundByOutRefundNo(ctx, req.OutRefundNo)
	if err == nil {
		if sameRefundRequest(existingRefund, req) {
			return existingRefund, nil
		}
		return nil, types.ErrPaymentDuplicate.WithMessage("refund idempotency key conflicts with an existing refund")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("query refund by out_refund_no: %w", err)
	}

	// 查找原订单
	order, err := s.repo.GetOrderByOutTradeNo(ctx, req.OutTradeNo)
	if err != nil {
		return nil, types.ErrPaymentNotFound
	}
	if !canRefundOrder(order.Status) {
		return nil, types.ErrPaymentFailed.WithMessagef("cannot refund payment order in %s status", order.Status)
	}
	if req.TotalAmount != order.Amount {
		return nil, types.ErrRefundAmount.WithMessagef("refund total amount mismatch: got %d, want %d", req.TotalAmount, order.Amount)
	}

	refundableAmount, err := s.refundableAmount(ctx, order)
	if err != nil {
		return nil, err
	}
	if req.RefundAmount > refundableAmount {
		return nil, types.ErrRefundAmount.WithMessagef("refund amount exceeds refundable amount: requested %d, remaining %d", req.RefundAmount, refundableAmount)
	}

	refund := &model.RefundOrder{
		OutRefundNo:  req.OutRefundNo,
		OutTradeNo:   req.OutTradeNo,
		TotalAmount:  req.TotalAmount,
		RefundAmount: req.RefundAmount,
		Reason:       req.Reason,
		Status:       string(model.PayStatusPending),
		Channel:      order.Channel,
	}

	if err := s.repo.CreateRefund(ctx, refund); err != nil {
		return nil, fmt.Errorf("create refund order: %w", err)
	}

	// 调用通道SDK发起退款
	ch, ok := s.chReg.Get(order.Channel)
	if !ok {
		return nil, types.ErrPaymentChannel.WithMessagef("channel %s not found", order.Channel)
	}

	refundResult, err := ch.Refund(ctx, req.OutTradeNo, req.OutRefundNo, req.TotalAmount, req.RefundAmount, req.Reason)
	if err != nil {
		refund.Status = string(model.PayStatusFailed)
		s.repo.UpdateRefund(ctx, refund)
		return nil, types.ErrPaymentFailed.WithMessage(err.Error())
	}

	refund.RefundNo = refundResult.RefundNo
	refund.Status = refundResult.Status
	if refundResult.Status == "success" {
		now := time.Now()
		refund.RefundedAt = &now
	}

	if err := s.repo.UpdateRefund(ctx, refund); err != nil {
		return nil, fmt.Errorf("update refund order: %w", err)
	}
	if err := s.updateOrderAfterRefund(ctx, order, refund); err != nil {
		return nil, err
	}

	return refund, nil
}

// 中文：sameRefundRequest 执行当前包中的对应流程。
// English: sameRefundRequest executes the corresponding workflow in this package.
func sameRefundRequest(refund *model.RefundOrder, req dto.RefundReq) bool {
	if refund == nil {
		return false
	}
	return refund.OutTradeNo == req.OutTradeNo &&
		refund.TotalAmount == req.TotalAmount &&
		refund.RefundAmount == req.RefundAmount
}

// 中文：canRefundOrder 执行当前包中的对应流程。
// English: canRefundOrder executes the corresponding workflow in this package.
func canRefundOrder(status string) bool {
	switch status {
	case string(model.PayStatusPaid), string(model.PayStatusRefunding):
		return true
	default:
		return false
	}
}

// 中文：refundableAmount 执行当前包中的对应流程。
// English: refundableAmount executes the corresponding workflow in this package.
func (s *PaymentService) refundableAmount(ctx context.Context, order *model.PaymentOrder) (int64, error) {
	refunds, err := s.repo.ListRefundsByOutTradeNo(ctx, order.OutTradeNo)
	if err != nil {
		return 0, fmt.Errorf("list existing refunds: %w", err)
	}

	used := int64(0)
	for _, refund := range refunds {
		if refundBlocksRefundableAmount(refund.Status) {
			used += refund.RefundAmount
		}
	}
	remaining := order.Amount - used
	if remaining < 0 {
		return 0, nil
	}
	return remaining, nil
}

// 中文：refundBlocksRefundableAmount 执行当前包中的对应流程。
// English: refundBlocksRefundableAmount executes the corresponding workflow in this package.
func refundBlocksRefundableAmount(status string) bool {
	switch status {
	case "success", "refunding", string(model.PayStatusPending):
		return true
	default:
		return false
	}
}

// 中文：QueryOrder 执行当前包中的对应流程。
// English: QueryOrder executes the corresponding workflow in this package.
// QueryOrder 查询订单
func (s *PaymentService) QueryOrder(ctx context.Context, id string) (*model.PaymentOrder, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, types.ErrPaymentNotFound
	}
	return order, nil
}

// 中文：CloseOrder 执行当前包中的对应流程。
// English: CloseOrder executes the corresponding workflow in this package.
// CloseOrder closes an unpaid order at the upstream channel and marks it closed locally.
func (s *PaymentService) CloseOrder(ctx context.Context, id string) (*model.PaymentOrder, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, types.ErrPaymentNotFound
	}

	switch order.Status {
	case string(model.PayStatusClosed):
		return order, nil
	case string(model.PayStatusPaid), string(model.PayStatusRefunding), string(model.PayStatusRefunded):
		return nil, types.ErrPaymentFailed.WithMessagef("cannot close payment order in %s status", order.Status)
	case string(model.PayStatusFailed):
		return order, nil
	}

	ch, ok := s.chReg.Get(order.Channel)
	if !ok {
		return nil, types.ErrPaymentChannel.WithMessagef("channel %s not found", order.Channel)
	}

	closeID := paymentCloseID(order)
	if err := ch.ClosePayment(ctx, closeID); err != nil {
		return nil, types.ErrPaymentFailed.WithMessage(err.Error())
	}

	order.Status = string(model.PayStatusClosed)
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("update closed payment order: %w", err)
	}
	s.publishAsync(ctx, event.NewEvent(event.EventPaymentClosed, paymentEventPayload(order)))
	return order, nil
}

// 中文：CallbackResponse 执行当前包中的对应流程。
// English: CallbackResponse executes the corresponding workflow in this package.
// CallbackResponse returns the channel-specific payload expected by payment gateways.
func (s *PaymentService) CallbackResponse(channelName string, success bool) any {
	if s == nil || s.chReg == nil {
		if success {
			return "success"
		}
		return "fail"
	}
	ch, ok := s.chReg.Get(channelName)
	if !ok {
		if success {
			return "success"
		}
		return "fail"
	}
	if success {
		return ch.CallbackSuccess()
	}
	return ch.CallbackFail()
}

// 中文：ListChannels 执行当前包中的对应流程。
// English: ListChannels executes the corresponding workflow in this package.
// ListChannels returns registered channel names in stable order.
func (s *PaymentService) ListChannels() []string {
	if s == nil || s.chReg == nil {
		return nil
	}
	channels := s.chReg.List()
	names := make([]string, 0, len(channels))
	for _, ch := range channels {
		names = append(names, ch.Name())
	}
	return names
}

// 中文：timePtr 执行当前包中的对应流程。
// English: timePtr executes the corresponding workflow in this package.
func timePtr(t time.Time) *time.Time { return &t }

// 中文：paymentCloseID 执行当前包中的对应流程。
// English: paymentCloseID executes the corresponding workflow in this package.
func paymentCloseID(order *model.PaymentOrder) string {
	if order == nil {
		return ""
	}
	switch order.Channel {
	case "stripe", "paypal", "unionpay":
		if order.TradeNo != "" {
			return order.TradeNo
		}
	}
	return order.OutTradeNo
}

// 中文：updateOrderAfterRefund 执行当前包中的对应流程。
// English: updateOrderAfterRefund executes the corresponding workflow in this package.
func (s *PaymentService) updateOrderAfterRefund(ctx context.Context, order *model.PaymentOrder, refund *model.RefundOrder) error {
	if order == nil || refund == nil {
		return nil
	}
	switch refund.Status {
	case "success":
		if refund.RefundAmount >= order.Amount || refund.RefundAmount >= refund.TotalAmount {
			order.Status = string(model.PayStatusRefunded)
		} else {
			order.Status = string(model.PayStatusRefunding)
		}
	case "refunding":
		order.Status = string(model.PayStatusRefunding)
	default:
		return nil
	}
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		return fmt.Errorf("update payment order after refund: %w", err)
	}
	if refund.Status == "success" {
		s.publishAsync(ctx, event.NewEvent(event.EventPaymentRefunded, paymentEventPayload(order)))
	}
	return nil
}

// 中文：paymentEventPayload 执行当前包中的对应流程。
// English: paymentEventPayload executes the corresponding workflow in this package.
func paymentEventPayload(order *model.PaymentOrder) *event.PaymentEventPayload {
	if order == nil {
		return &event.PaymentEventPayload{}
	}
	return &event.PaymentEventPayload{
		OrderID:    order.ID,
		OutTradeNo: order.OutTradeNo,
		TradeNo:    order.TradeNo,
		Channel:    order.Channel,
		Amount:     order.Amount,
		Currency:   order.Currency,
		Subject:    order.Subject,
		TenantID:   order.TenantID,
	}
}

// 中文：publishAsync 执行当前包中的对应流程。
// English: publishAsync executes the corresponding workflow in this package.
func (s *PaymentService) publishAsync(ctx context.Context, e *event.Event) {
	if s.eventBus != nil {
		_ = s.eventBus.PublishAsync(ctx, e)
	}
}
