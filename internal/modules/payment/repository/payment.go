package repository

import (
	"context"

	"github.com/spiringo/spiringo/internal/modules/payment/model"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：PaymentRepository 定义当前包使用的数据结构或接口。
// English: PaymentRepository defines a data structure or interface used by this package.
// PaymentRepository 支付数据访问
type PaymentRepository struct {
	// 中文：tdb 保存当前结构中的配置或数据值。
	// English: tdb stores a configuration or data value for this struct.
	tdb *orm.TenantDB
	// 中文：db 保存当前结构中的配置或数据值。
	// English: db stores a configuration or data value for this struct.
	db *orm.DB
}

// 中文：NewPaymentRepository 创建并返回对应组件实例。
// English: NewPaymentRepository creates and returns the corresponding component instance.
// NewPaymentRepository 创建支付仓库
func NewPaymentRepository(tdb *orm.TenantDB, db *orm.DB) *PaymentRepository {
	return &PaymentRepository{tdb: tdb, db: db}
}

// 中文：CreateOrder 执行当前包中的对应流程。
// English: CreateOrder executes the corresponding workflow in this package.
func (r *PaymentRepository) CreateOrder(ctx context.Context, order *model.PaymentOrder) error {
	return r.tdb.Create(ctx, order)
}

// 中文：GetOrderByID 执行当前包中的对应流程。
// English: GetOrderByID executes the corresponding workflow in this package.
func (r *PaymentRepository) GetOrderByID(ctx context.Context, id string) (*model.PaymentOrder, error) {
	var order model.PaymentOrder
	if err := r.tdb.First(ctx, &order, "id = ?", id); err != nil {
		return nil, err
	}
	return &order, nil
}

// 中文：GetOrderByOutTradeNo 执行当前包中的对应流程。
// English: GetOrderByOutTradeNo executes the corresponding workflow in this package.
func (r *PaymentRepository) GetOrderByOutTradeNo(ctx context.Context, outTradeNo string) (*model.PaymentOrder, error) {
	var order model.PaymentOrder
	if err := r.tdb.First(ctx, &order, "out_trade_no = ?", outTradeNo); err != nil {
		return nil, err
	}
	return &order, nil
}

// 中文：UpdateOrder 执行当前包中的对应流程。
// English: UpdateOrder executes the corresponding workflow in this package.
func (r *PaymentRepository) UpdateOrder(ctx context.Context, order *model.PaymentOrder) error {
	return r.tdb.Update(ctx, order)
}

// 中文：CreateRefund 执行当前包中的对应流程。
// English: CreateRefund executes the corresponding workflow in this package.
func (r *PaymentRepository) CreateRefund(ctx context.Context, refund *model.RefundOrder) error {
	return r.tdb.Create(ctx, refund)
}

// 中文：GetRefundByOutRefundNo 执行当前包中的对应流程。
// English: GetRefundByOutRefundNo executes the corresponding workflow in this package.
func (r *PaymentRepository) GetRefundByOutRefundNo(ctx context.Context, outRefundNo string) (*model.RefundOrder, error) {
	var refund model.RefundOrder
	if err := r.tdb.First(ctx, &refund, "out_refund_no = ?", outRefundNo); err != nil {
		return nil, err
	}
	return &refund, nil
}

// 中文：ListRefundsByOutTradeNo 执行当前包中的对应流程。
// English: ListRefundsByOutTradeNo executes the corresponding workflow in this package.
func (r *PaymentRepository) ListRefundsByOutTradeNo(ctx context.Context, outTradeNo string) ([]model.RefundOrder, error) {
	var refunds []model.RefundOrder
	if err := r.tdb.Find(ctx, &refunds, "out_trade_no = ?", outTradeNo); err != nil {
		return nil, err
	}
	return refunds, nil
}

// 中文：UpdateRefund 执行当前包中的对应流程。
// English: UpdateRefund executes the corresponding workflow in this package.
func (r *PaymentRepository) UpdateRefund(ctx context.Context, refund *model.RefundOrder) error {
	return r.tdb.Update(ctx, refund)
}

// 中文：CreateCallbackLog 执行当前包中的对应流程。
// English: CreateCallbackLog executes the corresponding workflow in this package.
func (r *PaymentRepository) CreateCallbackLog(ctx context.Context, log *model.CallbackLog) error {
	return r.db.Create(ctx, log)
}

// 中文：UpdateCallbackLog 执行当前包中的对应流程。
// English: UpdateCallbackLog executes the corresponding workflow in this package.
func (r *PaymentRepository) UpdateCallbackLog(ctx context.Context, log *model.CallbackLog) error {
	return r.db.Update(ctx, log)
}
