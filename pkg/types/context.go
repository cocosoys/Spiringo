package types

import "context"

// 中文：contextKey 定义当前包使用的数据结构或接口。
// English: contextKey defines a data structure or interface used by this package.
// contextKey 上下文键类型（防止外部包冲突）
type contextKey string

// 中文：CtxKeyRequestID、CtxKeyTenantID、CtxKeyTenantName、... 声明当前包使用的常量。
// English: CtxKeyRequestID、CtxKeyTenantID、CtxKeyTenantName、... declares constants used by this package.
const (
	// CtxKeyRequestID 请求ID
	CtxKeyRequestID contextKey = "request_id"
	// CtxKeyTenantID 租户ID
	CtxKeyTenantID contextKey = "tenant_id"
	// CtxKeyTenantName 租户名称
	CtxKeyTenantName contextKey = "tenant_name"
	// CtxKeyUserID 用户ID
	CtxKeyUserID contextKey = "user_id"
	// CtxKeyUsername 用户名
	CtxKeyUsername contextKey = "username"
	// CtxKeyLanguage 请求语言
	CtxKeyLanguage contextKey = "language"
)

// 中文：WithRequestID 执行当前包中的对应流程。
// English: WithRequestID executes the corresponding workflow in this package.
// WithRequestID 向context注入请求ID
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CtxKeyRequestID, id)
}

// 中文：GetRequestID 执行当前包中的对应流程。
// English: GetRequestID executes the corresponding workflow in this package.
// GetRequestID 从context获取请求ID
func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(CtxKeyRequestID).(string); ok {
		return v
	}
	return ""
}

// 中文：WithTenantID 执行当前包中的对应流程。
// English: WithTenantID executes the corresponding workflow in this package.
// WithTenantID 向context注入租户ID
func WithTenantID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CtxKeyTenantID, id)
}

// 中文：GetTenantID 执行当前包中的对应流程。
// English: GetTenantID executes the corresponding workflow in this package.
// GetTenantID 从context获取租户ID
func GetTenantID(ctx context.Context) string {
	if v, ok := ctx.Value(CtxKeyTenantID).(string); ok {
		return v
	}
	return ""
}

// 中文：WithTenantName 执行当前包中的对应流程。
// English: WithTenantName executes the corresponding workflow in this package.
// WithTenantName 向context注入租户名称
func WithTenantName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, CtxKeyTenantName, name)
}

// 中文：GetTenantName 执行当前包中的对应流程。
// English: GetTenantName executes the corresponding workflow in this package.
// GetTenantName 从context获取租户名称
func GetTenantName(ctx context.Context) string {
	if v, ok := ctx.Value(CtxKeyTenantName).(string); ok {
		return v
	}
	return ""
}

// 中文：WithUserID 执行当前包中的对应流程。
// English: WithUserID executes the corresponding workflow in this package.
// WithUserID 向context注入用户ID
func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CtxKeyUserID, id)
}

// 中文：GetUserID 执行当前包中的对应流程。
// English: GetUserID executes the corresponding workflow in this package.
// GetUserID 从context获取用户ID
func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(CtxKeyUserID).(string); ok {
		return v
	}
	return ""
}

// 中文：WithUsername 执行当前包中的对应流程。
// English: WithUsername executes the corresponding workflow in this package.
// WithUsername 向context注入用户名
func WithUsername(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, CtxKeyUsername, name)
}

// 中文：GetUsername 执行当前包中的对应流程。
// English: GetUsername executes the corresponding workflow in this package.
// GetUsername 从context获取用户名
func GetUsername(ctx context.Context) string {
	if v, ok := ctx.Value(CtxKeyUsername).(string); ok {
		return v
	}
	return ""
}

// 中文：WithLanguage 执行当前包中的对应流程。
// English: WithLanguage executes the corresponding workflow in this package.
// WithLanguage 向context注入请求语言
func WithLanguage(ctx context.Context, language string) context.Context {
	return context.WithValue(ctx, CtxKeyLanguage, language)
}

// 中文：GetLanguage 执行当前包中的对应流程。
// English: GetLanguage executes the corresponding workflow in this package.
// GetLanguage 从context获取请求语言
func GetLanguage(ctx context.Context) string {
	if v, ok := ctx.Value(CtxKeyLanguage).(string); ok {
		return v
	}
	return ""
}
