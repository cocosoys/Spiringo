package model

import (
	"time"

	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：PayStatus 定义当前包使用的数据结构或接口。
// English: PayStatus defines a data structure or interface used by this package.
// PayStatus 支付状态
type PayStatus string

// 中文：PayStatusPending、PayStatusPaid、PayStatusFailed、... 声明当前包使用的常量。
// English: PayStatusPending、PayStatusPaid、PayStatusFailed、... declares constants used by this package.
const (
	PayStatusPending   PayStatus = "pending"
	PayStatusPaid      PayStatus = "paid"
	PayStatusFailed    PayStatus = "failed"
	PayStatusClosed    PayStatus = "closed"
	PayStatusRefunding PayStatus = "refunding"
	PayStatusRefunded  PayStatus = "refunded"
)

// 中文：PaymentOrder 定义当前包使用的数据结构或接口。
// English: PaymentOrder defines a data structure or interface used by this package.
// PaymentOrder 支付订单
type PaymentOrder struct {
	// 中文：orm.TenantBaseModel 嵌入复用该类型提供的能力。
	// English: orm.TenantBaseModel embeds reusable behavior from that type.
	orm.TenantBaseModel
	// 中文：OutTradeNo 保存当前结构中的配置或数据值。
	// English: OutTradeNo stores a configuration or data value for this struct.
	OutTradeNo string `gorm:"size:64;uniqueIndex;not null" json:"out_trade_no"`
	// 中文：TradeNo 保存当前结构中的配置或数据值。
	// English: TradeNo stores a configuration or data value for this struct.
	TradeNo string `gorm:"size:64;index" json:"trade_no,omitempty"`
	// 中文：Channel 保存当前结构中的配置或数据值。
	// English: Channel stores a configuration or data value for this struct.
	Channel string `gorm:"size:32;not null" json:"channel"`
	// 中文：Scene 保存当前结构中的配置或数据值。
	// English: Scene stores a configuration or data value for this struct.
	Scene string `gorm:"size:32;not null" json:"scene"`
	// 中文：Amount 保存当前结构中的配置或数据值。
	// English: Amount stores a configuration or data value for this struct.
	Amount int64 `gorm:"not null" json:"amount"`
	// 中文：Currency 保存当前结构中的配置或数据值。
	// English: Currency stores a configuration or data value for this struct.
	Currency string `gorm:"size:8;not null;default:CNY" json:"currency"`
	// 中文：Subject 保存当前结构中的配置或数据值。
	// English: Subject stores a configuration or data value for this struct.
	Subject string `gorm:"size:256" json:"subject"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `gorm:"size:32;not null;default:pending" json:"status"`
	// 中文：UserID 保存当前结构中的配置或数据值。
	// English: UserID stores a configuration or data value for this struct.
	UserID string `gorm:"size:36;index" json:"user_id"`
	// 中文：NotifyURL 保存当前结构中的配置或数据值。
	// English: NotifyURL stores a configuration or data value for this struct.
	NotifyURL string `gorm:"size:512" json:"notify_url"`
	// 中文：ReturnURL 保存当前结构中的配置或数据值。
	// English: ReturnURL stores a configuration or data value for this struct.
	ReturnURL string `gorm:"size:512" json:"return_url"`
	// 中文：PaidAt 保存当前结构中的配置或数据值。
	// English: PaidAt stores a configuration or data value for this struct.
	PaidAt *time.Time `json:"paid_at,omitempty"`
	// 中文：ExpiredAt 保存当前结构中的配置或数据值。
	// English: ExpiredAt stores a configuration or data value for this struct.
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
	// 中文：PayParams 保存当前结构中的配置或数据值。
	// English: PayParams stores a configuration or data value for this struct.
	PayParams string `gorm:"type:text" json:"pay_params,omitempty"` // JSON
}

// 中文：TableName 执行当前包中的对应流程。
// English: TableName executes the corresponding workflow in this package.
func (PaymentOrder) TableName() string { return "payment_orders" }

// 中文：RefundOrder 定义当前包使用的数据结构或接口。
// English: RefundOrder defines a data structure or interface used by this package.
// RefundOrder 退款订单
type RefundOrder struct {
	// 中文：orm.TenantBaseModel 嵌入复用该类型提供的能力。
	// English: orm.TenantBaseModel embeds reusable behavior from that type.
	orm.TenantBaseModel
	// 中文：OutRefundNo 保存当前结构中的配置或数据值。
	// English: OutRefundNo stores a configuration or data value for this struct.
	OutRefundNo string `gorm:"size:64;uniqueIndex;not null" json:"out_refund_no"`
	// 中文：OutTradeNo 保存当前结构中的配置或数据值。
	// English: OutTradeNo stores a configuration or data value for this struct.
	OutTradeNo string `gorm:"size:64;index;not null" json:"out_trade_no"`
	// 中文：RefundNo 保存当前结构中的配置或数据值。
	// English: RefundNo stores a configuration or data value for this struct.
	RefundNo string `gorm:"size:64;index" json:"refund_no,omitempty"`
	// 中文：TotalAmount 保存当前结构中的配置或数据值。
	// English: TotalAmount stores a configuration or data value for this struct.
	TotalAmount int64 `gorm:"not null" json:"total_amount"`
	// 中文：RefundAmount 保存当前结构中的配置或数据值。
	// English: RefundAmount stores a configuration or data value for this struct.
	RefundAmount int64 `gorm:"not null" json:"refund_amount"`
	// 中文：Reason 保存当前结构中的配置或数据值。
	// English: Reason stores a configuration or data value for this struct.
	Reason string `gorm:"size:256" json:"reason"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `gorm:"size:32;not null;default:pending" json:"status"`
	// 中文：Channel 保存当前结构中的配置或数据值。
	// English: Channel stores a configuration or data value for this struct.
	Channel string `gorm:"size:32;not null" json:"channel"`
	// 中文：RefundedAt 保存当前结构中的配置或数据值。
	// English: RefundedAt stores a configuration or data value for this struct.
	RefundedAt *time.Time `json:"refunded_at,omitempty"`
}

// 中文：TableName 执行当前包中的对应流程。
// English: TableName executes the corresponding workflow in this package.
func (RefundOrder) TableName() string { return "refund_orders" }

// 中文：CallbackLog 定义当前包使用的数据结构或接口。
// English: CallbackLog defines a data structure or interface used by this package.
// CallbackLog 回调日志
type CallbackLog struct {
	// 中文：orm.BaseModel 嵌入复用该类型提供的能力。
	// English: orm.BaseModel embeds reusable behavior from that type.
	orm.BaseModel
	// 中文：Channel 保存当前结构中的配置或数据值。
	// English: Channel stores a configuration or data value for this struct.
	Channel string `gorm:"size:32;not null" json:"channel"`
	// 中文：TradeNo 保存当前结构中的配置或数据值。
	// English: TradeNo stores a configuration or data value for this struct.
	TradeNo string `gorm:"size:64;index" json:"trade_no,omitempty"`
	// 中文：RawData 保存当前结构中的配置或数据值。
	// English: RawData stores a configuration or data value for this struct.
	RawData string `gorm:"type:text" json:"raw_data"`
	// 中文：Processed 保存当前结构中的配置或数据值。
	// English: Processed stores a configuration or data value for this struct.
	Processed bool `gorm:"default:false" json:"processed"`
}

// 中文：TableName 执行当前包中的对应流程。
// English: TableName executes the corresponding workflow in this package.
func (CallbackLog) TableName() string { return "callback_logs" }
