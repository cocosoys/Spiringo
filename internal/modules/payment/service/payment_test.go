//go:build !(windows && 386)

package service

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/spiringo/spiringo/internal/core/event"
	"github.com/spiringo/spiringo/internal/modules/payment/channel"
	"github.com/spiringo/spiringo/internal/modules/payment/dto"
	"github.com/spiringo/spiringo/internal/modules/payment/model"
	"github.com/spiringo/spiringo/internal/modules/payment/repository"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：fakeChannel 定义当前包使用的数据结构或接口。
// English: fakeChannel defines a data structure or interface used by this package.
type fakeChannel struct {
	// 中文：result 保存当前结构中的配置或数据值。
	// English: result stores a configuration or data value for this struct.
	result *channel.CallbackResult
	// 中文：err 保存当前结构中的配置或数据值。
	// English: err stores a configuration or data value for this struct.
	err error
	// 中文：closeErr 保存当前结构中的配置或数据值。
	// English: closeErr stores a configuration or data value for this struct.
	closeErr error
	// 中文：closeID 保存当前结构中的配置或数据值。
	// English: closeID stores a configuration or data value for this struct.
	closeID *string
	// 中文：refundCalls 保存当前结构中的配置或数据值。
	// English: refundCalls stores a configuration or data value for this struct.
	refundCalls *int
	// 中文：refundNo 保存当前结构中的配置或数据值。
	// English: refundNo stores a configuration or data value for this struct.
	refundNo string
	// 中文：refundStatus 保存当前结构中的配置或数据值。
	// English: refundStatus stores a configuration or data value for this struct.
	refundStatus string
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (f fakeChannel) Name() string { return "fake" }

// 中文：CreatePayment 执行当前包中的对应流程。
// English: CreatePayment executes the corresponding workflow in this package.
func (f fakeChannel) CreatePayment(context.Context, string, string, int64, string, string, string, string) (*channel.PayResult, error) {
	return &channel.PayResult{}, nil
}

// 中文：VerifyCallback 执行当前包中的对应流程。
// English: VerifyCallback executes the corresponding workflow in this package.
func (f fakeChannel) VerifyCallback(context.Context, []byte) (*channel.CallbackResult, error) {
	return f.result, f.err
}

// 中文：Refund 执行当前包中的对应流程。
// English: Refund executes the corresponding workflow in this package.
func (f fakeChannel) Refund(context.Context, string, string, int64, int64, string) (*channel.RefundResult, error) {
	if f.refundCalls != nil {
		*f.refundCalls = *f.refundCalls + 1
	}
	status := f.refundStatus
	if status == "" {
		status = "success"
	}
	return &channel.RefundResult{RefundNo: f.refundNo, Status: status}, nil
}

// 中文：QueryPayment 执行当前包中的对应流程。
// English: QueryPayment executes the corresponding workflow in this package.
func (f fakeChannel) QueryPayment(context.Context, string) (*channel.CallbackResult, error) {
	return f.result, f.err
}

// 中文：ClosePayment 执行当前包中的对应流程。
// English: ClosePayment executes the corresponding workflow in this package.
func (f fakeChannel) ClosePayment(_ context.Context, outTradeNo string) error {
	if f.closeID != nil {
		*f.closeID = outTradeNo
	}
	return f.closeErr
}

// 中文：CallbackSuccess 执行当前包中的对应流程。
// English: CallbackSuccess executes the corresponding workflow in this package.
func (f fakeChannel) CallbackSuccess() any { return "success" }

// 中文：CallbackFail 执行当前包中的对应流程。
// English: CallbackFail executes the corresponding workflow in this package.
func (f fakeChannel) CallbackFail() any { return "fail" }

// 中文：TestHandleCallbackIsIdempotentForPaidOrder 验证相关行为符合预期。
// English: TestHandleCallbackIsIdempotentForPaidOrder verifies the related behavior.
func TestHandleCallbackIsIdempotentForPaidOrder(t *testing.T) {
	ctx := context.Background()
	svc, repo, db := newPaymentServiceForTest(t, &channel.CallbackResult{
		OutTradeNo: "order-001",
		TradeNo:    "trade-001",
		Status:     "paid",
		Amount:     100,
	})

	order := &model.PaymentOrder{
		OutTradeNo: "order-001",
		Channel:    "fake",
		Scene:      "native",
		Amount:     100,
		Currency:   "CNY",
		Subject:    "test",
		Status:     string(model.PayStatusPending),
	}
	if err := repo.CreateOrder(ctx, order); err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if err := svc.HandleCallback(ctx, "fake", []byte(`paid`)); err != nil {
		t.Fatalf("first HandleCallback returned error: %v", err)
	}
	if err := svc.HandleCallback(ctx, "fake", []byte(`paid-again`)); err != nil {
		t.Fatalf("duplicate HandleCallback returned error: %v", err)
	}

	updated, err := repo.GetOrderByOutTradeNo(ctx, "order-001")
	if err != nil {
		t.Fatalf("GetOrderByOutTradeNo returned error: %v", err)
	}
	if updated.Status != string(model.PayStatusPaid) {
		t.Fatalf("order status = %q, want paid", updated.Status)
	}

	var callbackCount int64
	if err := db.DB().Model(&model.CallbackLog{}).Count(&callbackCount).Error; err != nil {
		t.Fatalf("count callback logs: %v", err)
	}
	if callbackCount != 2 {
		t.Fatalf("callback log count = %d, want 2", callbackCount)
	}
}

// 中文：TestHandleCallbackRejectsAmountMismatch 验证相关行为符合预期。
// English: TestHandleCallbackRejectsAmountMismatch verifies the related behavior.
func TestHandleCallbackRejectsAmountMismatch(t *testing.T) {
	ctx := context.Background()
	svc, repo, _ := newPaymentServiceForTest(t, &channel.CallbackResult{
		OutTradeNo: "order-002",
		TradeNo:    "trade-002",
		Status:     "paid",
		Amount:     200,
	})

	order := &model.PaymentOrder{
		OutTradeNo: "order-002",
		Channel:    "fake",
		Scene:      "native",
		Amount:     100,
		Currency:   "CNY",
		Subject:    "test",
		Status:     string(model.PayStatusPending),
	}
	if err := repo.CreateOrder(ctx, order); err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if err := svc.HandleCallback(ctx, "fake", []byte(`paid`)); err == nil {
		t.Fatal("HandleCallback returned nil error, want amount mismatch")
	}
}

// 中文：TestHandleCallbackPublishesFulfillmentRequested 验证相关行为符合预期。
// English: TestHandleCallbackPublishesFulfillmentRequested verifies the related behavior.
func TestHandleCallbackPublishesFulfillmentRequested(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc, repo, _, bus := newPaymentServiceForTestWithBus(t, &channel.CallbackResult{
		OutTradeNo: "order-003",
		TradeNo:    "trade-003",
		Status:     "paid",
		Amount:     300,
	})

	seen := make(chan *event.PaymentEventPayload, 1)
	bus.Subscribe(event.EventPaymentFulfillmentRequested, func(ctx context.Context, e *event.Event) error {
		payload, ok := e.Payload.(*event.PaymentEventPayload)
		if ok {
			seen <- payload
		}
		return nil
	})
	bus.Start(ctx)
	defer bus.Stop()

	order := &model.PaymentOrder{
		OutTradeNo: "order-003",
		Channel:    "fake",
		Scene:      "native",
		Amount:     300,
		Currency:   "CNY",
		Subject:    "fulfillment",
		Status:     string(model.PayStatusPending),
	}
	if err := repo.CreateOrder(ctx, order); err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if err := svc.HandleCallback(ctx, "fake", []byte(`paid`)); err != nil {
		t.Fatalf("HandleCallback returned error: %v", err)
	}

	select {
	case payload := <-seen:
		if payload.OutTradeNo != "order-003" || payload.TradeNo != "trade-003" || payload.Subject != "fulfillment" {
			t.Fatalf("unexpected payload: %+v", payload)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for fulfillment event")
	}
}

// 中文：TestCloseOrderUpdatesStatusAndPublishesEvent 验证相关行为符合预期。
// English: TestCloseOrderUpdatesStatusAndPublishesEvent verifies the related behavior.
func TestCloseOrderUpdatesStatusAndPublishesEvent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := orm.New(orm.Config{
		Driver: "sqlite",
		DSN:    filepath.Join(t.TempDir(), "payment_close_test.db"),
	})
	if err != nil {
		t.Fatalf("orm.New returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	if err := db.AutoMigrate(&model.PaymentOrder{}, &model.RefundOrder{}, &model.CallbackLog{}); err != nil {
		t.Fatalf("AutoMigrate returned error: %v", err)
	}

	repo := repository.NewPaymentRepository(orm.NewTenantDB(db), db)
	var closeID string
	reg := channel.NewRegistry()
	reg.Register(fakeChannel{closeID: &closeID})
	bus := event.NewBus(1)
	svc := NewPaymentService(Config{}, bus, repo, reg)

	seen := make(chan *event.PaymentEventPayload, 1)
	bus.Subscribe(event.EventPaymentClosed, func(ctx context.Context, e *event.Event) error {
		payload, ok := e.Payload.(*event.PaymentEventPayload)
		if ok {
			seen <- payload
		}
		return nil
	})
	bus.Start(ctx)
	defer bus.Stop()

	order := &model.PaymentOrder{
		OutTradeNo: "order-close-001",
		Channel:    "fake",
		Scene:      "native",
		Amount:     100,
		Currency:   "CNY",
		Subject:    "close",
		Status:     string(model.PayStatusPending),
	}
	if err := repo.CreateOrder(ctx, order); err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	closed, err := svc.CloseOrder(ctx, order.ID)
	if err != nil {
		t.Fatalf("CloseOrder returned error: %v", err)
	}
	if closeID != "order-close-001" {
		t.Fatalf("ClosePayment id = %q, want order-close-001", closeID)
	}
	if closed.Status != string(model.PayStatusClosed) {
		t.Fatalf("closed status = %q, want closed", closed.Status)
	}

	select {
	case payload := <-seen:
		if payload.OutTradeNo != "order-close-001" || payload.OrderID != order.ID {
			t.Fatalf("unexpected close payload: %+v", payload)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for close event")
	}
}

// 中文：TestRefundUpdatesOriginalOrderAndPublishesEvent 验证相关行为符合预期。
// English: TestRefundUpdatesOriginalOrderAndPublishesEvent verifies the related behavior.
func TestRefundUpdatesOriginalOrderAndPublishesEvent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc, repo, _, bus := newPaymentServiceForTestWithBus(t, nil)
	seen := make(chan *event.PaymentEventPayload, 1)
	bus.Subscribe(event.EventPaymentRefunded, func(ctx context.Context, e *event.Event) error {
		payload, ok := e.Payload.(*event.PaymentEventPayload)
		if ok {
			seen <- payload
		}
		return nil
	})
	bus.Start(ctx)
	defer bus.Stop()

	order := &model.PaymentOrder{
		OutTradeNo: "order-refund-001",
		TradeNo:    "trade-refund-001",
		Channel:    "fake",
		Scene:      "native",
		Amount:     500,
		Currency:   "CNY",
		Subject:    "refund",
		Status:     string(model.PayStatusPaid),
	}
	if err := repo.CreateOrder(ctx, order); err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	refund, err := svc.Refund(ctx, dto.RefundReq{
		OutTradeNo:   "order-refund-001",
		OutRefundNo:  "refund-001",
		TotalAmount:  500,
		RefundAmount: 500,
		Reason:       "test",
	})
	if err != nil {
		t.Fatalf("Refund returned error: %v", err)
	}
	if refund.Status != "success" {
		t.Fatalf("refund status = %q, want success", refund.Status)
	}

	updated, err := repo.GetOrderByOutTradeNo(ctx, "order-refund-001")
	if err != nil {
		t.Fatalf("GetOrderByOutTradeNo returned error: %v", err)
	}
	if updated.Status != string(model.PayStatusRefunded) {
		t.Fatalf("order status = %q, want refunded", updated.Status)
	}

	select {
	case payload := <-seen:
		if payload.OutTradeNo != "order-refund-001" || payload.TradeNo != "trade-refund-001" {
			t.Fatalf("unexpected refund payload: %+v", payload)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for refund event")
	}
}

// 中文：TestRefundReturnsExistingRefundForSameIdempotencyKey 验证相关行为符合预期。
// English: TestRefundReturnsExistingRefundForSameIdempotencyKey verifies the related behavior.
func TestRefundReturnsExistingRefundForSameIdempotencyKey(t *testing.T) {
	ctx := context.Background()
	calls := 0
	svc, repo, db := newPaymentServiceForTestWithChannel(t, fakeChannel{refundCalls: &calls, refundNo: "refund-upstream-001"})

	order := &model.PaymentOrder{
		OutTradeNo: "order-refund-dup-001",
		TradeNo:    "trade-refund-dup-001",
		Channel:    "fake",
		Scene:      "native",
		Amount:     500,
		Currency:   "CNY",
		Subject:    "refund duplicate",
		Status:     string(model.PayStatusPaid),
	}
	if err := repo.CreateOrder(ctx, order); err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	req := dto.RefundReq{
		OutTradeNo:   "order-refund-dup-001",
		OutRefundNo:  "refund-dup-001",
		TotalAmount:  500,
		RefundAmount: 100,
		Reason:       "same request",
	}
	first, err := svc.Refund(ctx, req)
	if err != nil {
		t.Fatalf("first Refund returned error: %v", err)
	}
	second, err := svc.Refund(ctx, req)
	if err != nil {
		t.Fatalf("duplicate Refund returned error: %v", err)
	}
	if first.ID != second.ID {
		t.Fatalf("duplicate refund id = %q, want %q", second.ID, first.ID)
	}
	if calls != 1 {
		t.Fatalf("refund channel calls = %d, want 1", calls)
	}

	var refundCount int64
	if err := db.DB().Model(&model.RefundOrder{}).Count(&refundCount).Error; err != nil {
		t.Fatalf("count refunds: %v", err)
	}
	if refundCount != 1 {
		t.Fatalf("refund count = %d, want 1", refundCount)
	}
}

// 中文：TestRefundRejectsCumulativeOverRefund 验证相关行为符合预期。
// English: TestRefundRejectsCumulativeOverRefund verifies the related behavior.
func TestRefundRejectsCumulativeOverRefund(t *testing.T) {
	ctx := context.Background()
	calls := 0
	svc, repo, _ := newPaymentServiceForTestWithChannel(t, fakeChannel{refundCalls: &calls})

	order := &model.PaymentOrder{
		OutTradeNo: "order-refund-over-001",
		TradeNo:    "trade-refund-over-001",
		Channel:    "fake",
		Scene:      "native",
		Amount:     500,
		Currency:   "CNY",
		Subject:    "refund over",
		Status:     string(model.PayStatusPaid),
	}
	if err := repo.CreateOrder(ctx, order); err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	if _, err := svc.Refund(ctx, dto.RefundReq{
		OutTradeNo:   "order-refund-over-001",
		OutRefundNo:  "refund-over-001",
		TotalAmount:  500,
		RefundAmount: 300,
	}); err != nil {
		t.Fatalf("first Refund returned error: %v", err)
	}
	if _, err := svc.Refund(ctx, dto.RefundReq{
		OutTradeNo:   "order-refund-over-001",
		OutRefundNo:  "refund-over-002",
		TotalAmount:  500,
		RefundAmount: 250,
	}); err == nil {
		t.Fatal("second Refund returned nil error, want over-refund error")
	}
	if calls != 1 {
		t.Fatalf("refund channel calls = %d, want 1", calls)
	}
}

// 中文：newPaymentServiceForTest 执行当前包中的对应流程。
// English: newPaymentServiceForTest executes the corresponding workflow in this package.
func newPaymentServiceForTest(t *testing.T, result *channel.CallbackResult) (*PaymentService, *repository.PaymentRepository, *orm.DB) {
	svc, repo, db, _ := newPaymentServiceForTestWithBus(t, result)
	return svc, repo, db
}

// 中文：newPaymentServiceForTestWithBus 执行当前包中的对应流程。
// English: newPaymentServiceForTestWithBus executes the corresponding workflow in this package.
func newPaymentServiceForTestWithBus(t *testing.T, result *channel.CallbackResult) (*PaymentService, *repository.PaymentRepository, *orm.DB, *event.Bus) {
	t.Helper()

	db, err := orm.New(orm.Config{
		Driver: "sqlite",
		DSN:    filepath.Join(t.TempDir(), "payment_test.db"),
	})
	if err != nil {
		t.Fatalf("orm.New returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	if err := db.AutoMigrate(&model.PaymentOrder{}, &model.RefundOrder{}, &model.CallbackLog{}); err != nil {
		t.Fatalf("AutoMigrate returned error: %v", err)
	}

	repo := repository.NewPaymentRepository(orm.NewTenantDB(db), db)
	reg := channel.NewRegistry()
	reg.Register(fakeChannel{result: result})
	bus := event.NewBus(1)
	return NewPaymentService(Config{}, bus, repo, reg), repo, db, bus
}

// 中文：newPaymentServiceForTestWithChannel 执行当前包中的对应流程。
// English: newPaymentServiceForTestWithChannel executes the corresponding workflow in this package.
func newPaymentServiceForTestWithChannel(t *testing.T, ch channel.Channel) (*PaymentService, *repository.PaymentRepository, *orm.DB) {
	t.Helper()

	db, err := orm.New(orm.Config{
		Driver: "sqlite",
		DSN:    filepath.Join(t.TempDir(), "payment_refund_test.db"),
	})
	if err != nil {
		t.Fatalf("orm.New returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	if err := db.AutoMigrate(&model.PaymentOrder{}, &model.RefundOrder{}, &model.CallbackLog{}); err != nil {
		t.Fatalf("AutoMigrate returned error: %v", err)
	}

	repo := repository.NewPaymentRepository(orm.NewTenantDB(db), db)
	reg := channel.NewRegistry()
	reg.Register(ch)
	bus := event.NewBus(1)
	return NewPaymentService(Config{}, bus, repo, reg), repo, db
}
