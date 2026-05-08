package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/pkg/crypto"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：fakePermissionChecker 定义当前包使用的数据结构或接口。
// English: fakePermissionChecker defines a data structure or interface used by this package.
type fakePermissionChecker struct {
	// 中文：allowed 保存当前结构中的配置或数据值。
	// English: allowed stores a configuration or data value for this struct.
	allowed bool
}

// 中文：CheckPermission 执行当前包中的对应流程。
// English: CheckPermission executes the corresponding workflow in this package.
func (f fakePermissionChecker) CheckPermission(context.Context, string, string, string) (bool, error) {
	return f.allowed, nil
}

// 中文：TestAuthInjectsClaimsIntoRequestContext 验证相关行为符合预期。
// English: TestAuthInjectsClaimsIntoRequestContext verifies the related behavior.
func TestAuthInjectsClaimsIntoRequestContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	token, _, err := crypto.GenerateToken(crypto.JWTConfig{
		Secret:        "secret",
		AccessExpire:  time.Hour,
		RefreshExpire: time.Hour,
		Issuer:        "test",
	}, "user-1", "alice", "tenant-1")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	router := gin.New()
	router.Use(JWTAuth("secret"))
	router.GET("/protected", func(c *gin.Context) {
		if got := types.GetUserID(c.Request.Context()); got != "user-1" {
			t.Fatalf("context user id = %q, want user-1", got)
		}
		if got := types.GetTenantID(c.Request.Context()); got != "tenant-1" {
			t.Fatalf("context tenant id = %q, want tenant-1", got)
		}
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

// 中文：TestAuthWithOptionsSkipsPublicPath 验证相关行为符合预期。
// English: TestAuthWithOptionsSkipsPublicPath verifies the related behavior.
func TestAuthWithOptionsSkipsPublicPath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(AuthWithOptions(nil, AuthOptions{PublicPaths: []string{"/public", "/oauth/*"}}))
	router.GET("/public", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	router.GET("/oauth/callback", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	for _, path := range []string{"/public", "/oauth/callback"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("%s status = %d, want %d", path, rec.Code, http.StatusNoContent)
		}
	}
}

// 中文：TestAuthWithOptionsProtectsPrivatePath 验证相关行为符合预期。
// English: TestAuthWithOptionsProtectsPrivatePath verifies the related behavior.
func TestAuthWithOptionsProtectsPrivatePath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(AuthWithOptions(nil, AuthOptions{PublicPaths: []string{"/public"}}))
	router.GET("/private", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

// 中文：TestRequirePermissionRejectsMissingPermission 验证相关行为符合预期。
// English: TestRequirePermissionRejectsMissingPermission verifies the related behavior.
func TestRequirePermissionRejectsMissingPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		ctx := types.WithUserID(c.Request.Context(), "user-1")
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})
	router.GET("/protected", RequirePermission(fakePermissionChecker{allowed: false}, "rbac.roles", "read"), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}
