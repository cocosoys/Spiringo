package types

import "fmt"

// 中文：Error 定义当前包使用的数据结构或接口。
// English: Error defines a data structure or interface used by this package.
// Error 统一错误结构
type Error struct {
	// 中文：Code 保存当前结构中的配置或数据值。
	// English: Code stores a configuration or data value for this struct.
	Code int `json:"code"`
	// 中文：Message 保存当前结构中的配置或数据值。
	// English: Message stores a configuration or data value for this struct.
	Message string `json:"message"`
	// 中文：HTTPStatus 保存当前结构中的配置或数据值。
	// English: HTTPStatus stores a configuration or data value for this struct.
	HTTPStatus int `json:"-"`
}

// 中文：Error 执行当前包中的对应流程。
// English: Error executes the corresponding workflow in this package.
func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// 中文：WithMessage 执行当前包中的对应流程。
// English: WithMessage executes the corresponding workflow in this package.
// WithMessage 返回相同错误码但不同消息的错误
func (e *Error) WithMessage(msg string) *Error {
	return &Error{
		Code:       e.Code,
		Message:    msg,
		HTTPStatus: e.HTTPStatus,
	}
}

// 中文：WithMessagef 执行当前包中的对应流程。
// English: WithMessagef executes the corresponding workflow in this package.
// WithMessagef 格式化消息
func (e *Error) WithMessagef(format string, args ...any) *Error {
	return &Error{
		Code:       e.Code,
		Message:    fmt.Sprintf(format, args...),
		HTTPStatus: e.HTTPStatus,
	}
}

// 中文：NewError 创建并返回对应组件实例。
// English: NewError creates and returns the corresponding component instance.
// NewError 创建新错误
func NewError(httpStatus, code int, message string) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// ---- 通用错误 ----

// 中文：ErrSuccess、ErrInternal、ErrBadRequest、... 声明当前包使用的变量。
// English: ErrSuccess、ErrInternal、ErrBadRequest、... declares variables used by this package.
var (
	ErrSuccess      = NewError(200, 0, "success")
	ErrInternal     = NewError(500, 10000, "internal server error")
	ErrBadRequest   = NewError(400, 10001, "bad request")
	ErrUnauthorized = NewError(401, 10002, "unauthorized")
	ErrForbidden    = NewError(403, 10003, "forbidden")
	ErrNotFound     = NewError(404, 10004, "not found")
	ErrConflict     = NewError(409, 10005, "conflict")
	ErrTooManyReq   = NewError(429, 10006, "too many requests")
	ErrValidation   = NewError(400, 10007, "validation failed")
	ErrIdempotent   = NewError(409, 10008, "duplicate request")
)

// ---- 认证相关错误 (2xxxx) ----

// 中文：ErrInvalidCredentials、ErrTokenExpired、ErrTokenInvalid、... 声明当前包使用的变量。
// English: ErrInvalidCredentials、ErrTokenExpired、ErrTokenInvalid、... declares variables used by this package.
var (
	ErrInvalidCredentials = NewError(401, 20001, "invalid credentials")
	ErrTokenExpired       = NewError(401, 20002, "token expired")
	ErrTokenInvalid       = NewError(401, 20003, "token invalid")
	ErrTokenRefreshFailed = NewError(401, 20004, "token refresh failed")
	ErrOAuthFailed        = NewError(401, 20005, "oauth authentication failed")
	ErrOAuthProvider      = NewError(400, 20006, "unsupported oauth provider")
)

// ---- 用户相关错误 (3xxxx) ----

// 中文：ErrUserNotFound、ErrUserExists、ErrUserDisabled、... 声明当前包使用的变量。
// English: ErrUserNotFound、ErrUserExists、ErrUserDisabled、... declares variables used by this package.
var (
	ErrUserNotFound  = NewError(404, 30001, "user not found")
	ErrUserExists    = NewError(409, 30002, "user already exists")
	ErrUserDisabled  = NewError(403, 30003, "user is disabled")
	ErrPasswordWrong = NewError(401, 30004, "wrong password")
	ErrPasswordWeak  = NewError(400, 30005, "password is too weak")
)

// ---- 租户相关错误 (4xxxx) ----

// 中文：ErrTenantNotFound、ErrTenantDisabled、ErrTenantExpired、... 声明当前包使用的变量。
// English: ErrTenantNotFound、ErrTenantDisabled、ErrTenantExpired、... declares variables used by this package.
var (
	ErrTenantNotFound    = NewError(404, 40001, "tenant not found")
	ErrTenantDisabled    = NewError(403, 40002, "tenant is disabled")
	ErrTenantExpired     = NewError(403, 40003, "tenant has expired")
	ErrTenantQuotaExceed = NewError(403, 40004, "tenant quota exceeded")
)

// ---- 支付相关错误 (5xxxx) ----

// 中文：ErrPaymentNotFound、ErrPaymentExpired、ErrPaymentDuplicate、... 声明当前包使用的变量。
// English: ErrPaymentNotFound、ErrPaymentExpired、ErrPaymentDuplicate、... declares variables used by this package.
var (
	ErrPaymentNotFound  = NewError(404, 50001, "payment order not found")
	ErrPaymentExpired   = NewError(400, 50002, "payment order expired")
	ErrPaymentDuplicate = NewError(409, 50003, "duplicate payment order")
	ErrPaymentChannel   = NewError(400, 50004, "unsupported payment channel")
	ErrPaymentFailed    = NewError(500, 50005, "payment failed")
	ErrRefundNotFound   = NewError(404, 50006, "refund order not found")
	ErrRefundAmount     = NewError(400, 50007, "refund amount exceeds paid amount")
	ErrCallbackVerify   = NewError(400, 50008, "callback signature verification failed")
)

// ---- 权限相关错误 (6xxxx) ----

// 中文：ErrRoleNotFound、ErrRoleExists、ErrPermDenied、... 声明当前包使用的变量。
// English: ErrRoleNotFound、ErrRoleExists、ErrPermDenied、... declares variables used by this package.
var (
	ErrRoleNotFound = NewError(404, 60001, "role not found")
	ErrRoleExists   = NewError(409, 60002, "role already exists")
	ErrPermDenied   = NewError(403, 60003, "permission denied")
	ErrPermNotFound = NewError(404, 60004, "permission not found")
)
