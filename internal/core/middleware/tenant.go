package middleware

import (
	"context"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：TenantStrategy 定义当前包使用的数据结构或接口。
// English: TenantStrategy defines a data structure or interface used by this package.
type TenantStrategy string

// 中文：StrategySharedDB、StrategySchema、StrategyDatabase 声明当前包使用的常量。
// English: StrategySharedDB、StrategySchema、StrategyDatabase declares constants used by this package.
const (
	StrategySharedDB TenantStrategy = "shared_db"
	StrategySchema   TenantStrategy = "schema"
	StrategyDatabase TenantStrategy = "database"
)

// 中文：TenantContext 定义当前包使用的数据结构或接口。
// English: TenantContext defines a data structure or interface used by this package.
// TenantContext is the structured tenant state attached to request contexts.
type TenantContext struct {
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string
	// 中文：TenantName 保存当前结构中的配置或数据值。
	// English: TenantName stores a configuration or data value for this struct.
	TenantName string
	// 中文：Strategy 保存当前结构中的配置或数据值。
	// English: Strategy stores a configuration or data value for this struct.
	Strategy TenantStrategy
	// 中文：DBConn 保存当前结构中的配置或数据值。
	// English: DBConn stores a configuration or data value for this struct.
	DBConn string
}

// 中文：tenantContextKey 定义当前包使用的数据结构或接口。
// English: tenantContextKey defines a data structure or interface used by this package.
type tenantContextKey struct{}

// 中文：WithTenant 执行当前包中的对应流程。
// English: WithTenant executes the corresponding workflow in this package.
func WithTenant(ctx context.Context, tc *TenantContext) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if tc == nil {
		return ctx
	}
	clone := *tc
	if clone.Strategy == "" {
		clone.Strategy = StrategySharedDB
	}
	if clone.TenantID != "" {
		ctx = types.WithTenantID(ctx, clone.TenantID)
	}
	if clone.TenantName != "" {
		ctx = types.WithTenantName(ctx, clone.TenantName)
	}
	return context.WithValue(ctx, tenantContextKey{}, &clone)
}

// 中文：FromContext 执行当前包中的对应流程。
// English: FromContext executes the corresponding workflow in this package.
func FromContext(ctx context.Context) *TenantContext {
	if ctx == nil {
		return nil
	}
	if tc, ok := ctx.Value(tenantContextKey{}).(*TenantContext); ok && tc != nil {
		clone := *tc
		return &clone
	}
	tenantID := types.GetTenantID(ctx)
	tenantName := types.GetTenantName(ctx)
	if tenantID == "" && tenantName == "" {
		return nil
	}
	return &TenantContext{
		TenantID:   tenantID,
		TenantName: tenantName,
		Strategy:   StrategySharedDB,
	}
}

// 中文：TenantMiddleware 执行当前包中的对应流程。
// English: TenantMiddleware executes the corresponding workflow in this package.
// TenantMiddleware injects tenant context from request headers.
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := strings.TrimSpace(c.GetHeader("X-Tenant-ID"))
		if tenantID == "" {
			tenantID = tenantIDFromHost(c.Request.Host)
		}
		tc := &TenantContext{
			TenantID:   tenantID,
			TenantName: strings.TrimSpace(c.GetHeader("X-Tenant-Name")),
			Strategy:   parseTenantStrategy(c.GetHeader("X-Tenant-Strategy")),
			DBConn:     strings.TrimSpace(c.GetHeader("X-Tenant-DB-Conn")),
		}
		if tc.TenantID != "" || tc.TenantName != "" || tc.DBConn != "" || c.GetHeader("X-Tenant-Strategy") != "" {
			c.Request = c.Request.WithContext(WithTenant(c.Request.Context(), tc))
		}
		c.Next()
	}
}

// 中文：tenantIDFromHost 执行当前包中的对应流程。
// English: tenantIDFromHost executes the corresponding workflow in this package.
func tenantIDFromHost(host string) string {
	host = strings.TrimSpace(host)
	if host == "" {
		return ""
	}
	if parsed, _, err := net.SplitHostPort(host); err == nil {
		host = parsed
	}
	host = strings.Trim(strings.TrimSuffix(host, "."), "[]")
	if host == "" || strings.EqualFold(host, "localhost") || net.ParseIP(host) != nil {
		return ""
	}
	labels := strings.Split(host, ".")
	if len(labels) < 3 {
		return ""
	}
	candidate := strings.ToLower(strings.TrimSpace(labels[0]))
	switch candidate {
	case "", "www", "api", "app", "admin":
		return ""
	default:
		return candidate
	}
}

// 中文：parseTenantStrategy 执行当前包中的对应流程。
// English: parseTenantStrategy executes the corresponding workflow in this package.
func parseTenantStrategy(value string) TenantStrategy {
	switch TenantStrategy(strings.TrimSpace(value)) {
	case StrategySchema:
		return StrategySchema
	case StrategyDatabase:
		return StrategyDatabase
	default:
		return StrategySharedDB
	}
}
