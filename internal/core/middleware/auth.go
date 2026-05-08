package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/pkg/crypto"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：GinUserIDKey、GinUsernameKey、GinTenantIDKey 声明当前包使用的常量。
// English: GinUserIDKey、GinUsernameKey、GinTenantIDKey declares constants used by this package.
const (
	GinUserIDKey   = "user_id"
	GinUsernameKey = "username"
	GinTenantIDKey = "tenant_id"
)

// 中文：TokenValidator 定义当前包使用的数据结构或接口。
// English: TokenValidator defines a data structure or interface used by this package.
type TokenValidator interface {
	// 中文：ValidateToken 声明该接口需要实现的行为。
	// English: ValidateToken declares behavior required by this interface.
	ValidateToken(ctx context.Context, tokenString string) (*crypto.JWTClaims, error)
}

// 中文：PermissionChecker 定义当前包使用的数据结构或接口。
// English: PermissionChecker defines a data structure or interface used by this package.
type PermissionChecker interface {
	// 中文：CheckPermission 声明该接口需要实现的行为。
	// English: CheckPermission declares behavior required by this interface.
	CheckPermission(ctx context.Context, userID, resource, action string) (bool, error)
}

// 中文：AuthOptions 定义当前包使用的数据结构或接口。
// English: AuthOptions defines a data structure or interface used by this package.
type AuthOptions struct {
	// 中文：PublicPaths 保存当前结构中的配置或数据值。
	// English: PublicPaths stores a configuration or data value for this struct.
	PublicPaths []string
}

// 中文：JWTAuth 执行当前包中的对应流程。
// English: JWTAuth executes the corresponding workflow in this package.
func JWTAuth(secret string) gin.HandlerFunc {
	return JWTAuthWithOptions(secret, AuthOptions{})
}

// 中文：JWTAuthWithOptions 执行当前包中的对应流程。
// English: JWTAuthWithOptions executes the corresponding workflow in this package.
func JWTAuthWithOptions(secret string, opts AuthOptions) gin.HandlerFunc {
	return AuthWithOptions(tokenValidatorFunc(func(_ context.Context, tokenString string) (*crypto.JWTClaims, error) {
		return crypto.ParseToken(secret, tokenString)
	}), opts)
}

// 中文：Auth 执行当前包中的对应流程。
// English: Auth executes the corresponding workflow in this package.
func Auth(validator TokenValidator) gin.HandlerFunc {
	return AuthWithOptions(validator, AuthOptions{})
}

// 中文：AuthWithOptions 执行当前包中的对应流程。
// English: AuthWithOptions executes the corresponding workflow in this package.
func AuthWithOptions(validator TokenValidator, opts AuthOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isPublicPath(c.Request.URL.Path, opts.PublicPaths) {
			c.Next()
			return
		}
		if validator == nil {
			abortWithError(c, types.ErrUnauthorized.WithMessage("token validator is not configured"))
			return
		}

		tokenString := extractBearerToken(c.GetHeader("Authorization"))
		if tokenString == "" {
			abortWithError(c, types.ErrUnauthorized.WithMessage("missing bearer token"))
			return
		}

		claims, err := validator.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			abortWithError(c, types.ErrTokenInvalid.WithMessage(err.Error()))
			return
		}
		if claims.UserID == "" {
			abortWithError(c, types.ErrTokenInvalid.WithMessage("token missing user id"))
			return
		}

		c.Set(GinUserIDKey, claims.UserID)
		c.Set(GinUsernameKey, claims.Username)
		if claims.TenantID != "" {
			c.Set(GinTenantIDKey, claims.TenantID)
		}

		ctx := types.WithUserID(c.Request.Context(), claims.UserID)
		if claims.Username != "" {
			ctx = types.WithUsername(ctx, claims.Username)
		}
		if claims.TenantID != "" {
			ctx = types.WithTenantID(ctx, claims.TenantID)
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// 中文：RequirePermission 执行当前包中的对应流程。
// English: RequirePermission executes the corresponding workflow in this package.
func RequirePermission(checker PermissionChecker, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if checker == nil {
			abortWithError(c, types.ErrForbidden.WithMessage("permission checker is not configured"))
			return
		}

		userID := types.GetUserID(c.Request.Context())
		if userID == "" {
			abortWithError(c, types.ErrUnauthorized)
			return
		}

		ok, err := checker.CheckPermission(c.Request.Context(), userID, resource, action)
		if err != nil {
			abortWithError(c, types.ErrPermDenied.WithMessage(err.Error()))
			return
		}
		if !ok {
			abortWithError(c, types.ErrPermDenied.WithMessagef("missing permission %s:%s", resource, action))
			return
		}
		c.Next()
	}
}

// 中文：extractBearerToken 执行当前包中的对应流程。
// English: extractBearerToken executes the corresponding workflow in this package.
func extractBearerToken(authHeader string) string {
	parts := strings.Fields(strings.TrimSpace(authHeader))
	if len(parts) == 1 {
		return parts[0]
	}
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

// 中文：isPublicPath 执行当前包中的对应流程。
// English: isPublicPath executes the corresponding workflow in this package.
func isPublicPath(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if pathMatches(pattern, path) {
			return true
		}
	}
	return false
}

// 中文：pathMatches 执行当前包中的对应流程。
// English: pathMatches executes the corresponding workflow in this package.
func pathMatches(pattern, path string) bool {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return false
	}
	if pattern == "*" || pattern == path {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(path, prefix)
	}
	return false
}

// 中文：abortWithError 执行当前包中的对应流程。
// English: abortWithError executes the corresponding workflow in this package.
func abortWithError(c *gin.Context, err error) {
	types.Fail(c, err)
	c.Abort()
}

// 中文：tokenValidatorFunc 定义当前包使用的数据结构或接口。
// English: tokenValidatorFunc defines a data structure or interface used by this package.
type tokenValidatorFunc func(context.Context, string) (*crypto.JWTClaims, error)

// 中文：ValidateToken 执行当前包中的对应流程。
// English: ValidateToken executes the corresponding workflow in this package.
func (f tokenValidatorFunc) ValidateToken(ctx context.Context, tokenString string) (*crypto.JWTClaims, error) {
	return f(ctx, tokenString)
}
