package event

// 预定义事件主题

// 中文：EventUserCreated、EventUserUpdated、EventUserDeleted、... 声明当前包使用的常量。
// English: EventUserCreated、EventUserUpdated、EventUserDeleted、... declares constants used by this package.
const (
	// 用户模块
	EventUserCreated = "user.created"
	EventUserUpdated = "user.updated"
	EventUserDeleted = "user.deleted"

	// 认证模块
	EventAuthLogin      = "auth.login"
	EventAuthLogout     = "auth.logout"
	EventAuthOAuthBound = "auth.oauth_bound"

	// 支付模块
	EventPaymentCreated              = "payment.created"
	EventPaymentSuccess              = "payment.success"
	EventPaymentFailed               = "payment.failed"
	EventPaymentRefunded             = "payment.refunded"
	EventPaymentClosed               = "payment.closed"
	EventPaymentFulfillmentRequested = "payment.fulfillment_requested"

	// 租户模块
	EventTenantCreated   = "tenant.created"
	EventTenantSuspended = "tenant.suspended"
	EventTenantActivated = "tenant.activated"
)

// 中文：PaymentEventPayload 定义当前包使用的数据结构或接口。
// English: PaymentEventPayload defines a data structure or interface used by this package.
// PaymentEventPayload 支付事件载荷
type PaymentEventPayload struct {
	// 中文：OrderID 保存当前结构中的配置或数据值。
	// English: OrderID stores a configuration or data value for this struct.
	OrderID string `json:"order_id"`
	// 中文：OutTradeNo 保存当前结构中的配置或数据值。
	// English: OutTradeNo stores a configuration or data value for this struct.
	OutTradeNo string `json:"out_trade_no"`
	// 中文：TradeNo 保存当前结构中的配置或数据值。
	// English: TradeNo stores a configuration or data value for this struct.
	TradeNo string `json:"trade_no"`
	// 中文：Channel 保存当前结构中的配置或数据值。
	// English: Channel stores a configuration or data value for this struct.
	Channel string `json:"channel"`
	// 中文：Amount 保存当前结构中的配置或数据值。
	// English: Amount stores a configuration or data value for this struct.
	Amount int64 `json:"amount"`
	// 中文：Currency 保存当前结构中的配置或数据值。
	// English: Currency stores a configuration or data value for this struct.
	Currency string `json:"currency"`
	// 中文：Subject 保存当前结构中的配置或数据值。
	// English: Subject stores a configuration or data value for this struct.
	Subject string `json:"subject"`
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string `json:"tenant_id"`
}

// 中文：UserEventPayload 定义当前包使用的数据结构或接口。
// English: UserEventPayload defines a data structure or interface used by this package.
// UserEventPayload 用户事件载荷
type UserEventPayload struct {
	// 中文：UserID 保存当前结构中的配置或数据值。
	// English: UserID stores a configuration or data value for this struct.
	UserID string `json:"user_id"`
	// 中文：Username 保存当前结构中的配置或数据值。
	// English: Username stores a configuration or data value for this struct.
	Username string `json:"username"`
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string `json:"tenant_id"`
}

// 中文：TenantEventPayload 定义当前包使用的数据结构或接口。
// English: TenantEventPayload defines a data structure or interface used by this package.
// TenantEventPayload 租户事件载荷
type TenantEventPayload struct {
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string `json:"tenant_id"`
	// 中文：TenantName 保存当前结构中的配置或数据值。
	// English: TenantName stores a configuration or data value for this struct.
	TenantName string `json:"tenant_name"`
	// 中文：Strategy 保存当前结构中的配置或数据值。
	// English: Strategy stores a configuration or data value for this struct.
	Strategy string `json:"strategy"`
}

// 中文：AuthEventPayload 定义当前包使用的数据结构或接口。
// English: AuthEventPayload defines a data structure or interface used by this package.
// AuthEventPayload 认证事件载荷
type AuthEventPayload struct {
	// 中文：UserID 保存当前结构中的配置或数据值。
	// English: UserID stores a configuration or data value for this struct.
	UserID string `json:"user_id"`
	// 中文：IP 保存当前结构中的配置或数据值。
	// English: IP stores a configuration or data value for this struct.
	IP string `json:"ip"`
	// 中文：UserAgent 保存当前结构中的配置或数据值。
	// English: UserAgent stores a configuration or data value for this struct.
	UserAgent string `json:"user_agent"`
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string `json:"tenant_id"`
}
