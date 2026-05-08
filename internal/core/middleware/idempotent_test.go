package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：TestIdempotentRejectsDuplicateScopedRequest 验证相关行为符合预期。
// English: TestIdempotentRejectsDuplicateScopedRequest verifies the related behavior.
func TestIdempotentRejectsDuplicateScopedRequest(t *testing.T) {
	router := newIdempotentTestRouter()

	if code := performIdempotentRequest(router, http.MethodPost, "/orders/1", "order-key", "tenant-a", "user-a"); code != http.StatusNoContent {
		t.Fatalf("first request status = %d, want %d", code, http.StatusNoContent)
	}
	if code := performIdempotentRequest(router, http.MethodPost, "/orders/1", "order-key", "tenant-a", "user-a"); code != http.StatusConflict {
		t.Fatalf("duplicate request status = %d, want %d", code, http.StatusConflict)
	}
}

// 中文：TestIdempotentScopesKeyByRouteTenantAndUser 验证相关行为符合预期。
// English: TestIdempotentScopesKeyByRouteTenantAndUser verifies the related behavior.
func TestIdempotentScopesKeyByRouteTenantAndUser(t *testing.T) {
	router := newIdempotentTestRouter()

	cases := []struct {
		name   string
		method string
		path   string
		key    string
		tenant string
		user   string
		want   int
	}{
		{name: "first scoped request", method: http.MethodPost, path: "/orders/1", key: "shared-key", tenant: "tenant-a", user: "user-a", want: http.StatusNoContent},
		{name: "same key different route", method: http.MethodPost, path: "/refunds", key: "shared-key", tenant: "tenant-a", user: "user-a", want: http.StatusNoContent},
		{name: "same key different tenant", method: http.MethodPost, path: "/orders/1", key: "shared-key", tenant: "tenant-b", user: "user-a", want: http.StatusNoContent},
		{name: "same key different user", method: http.MethodPost, path: "/orders/1", key: "shared-key", tenant: "tenant-a", user: "user-b", want: http.StatusNoContent},
		{name: "duplicate original scope", method: http.MethodPost, path: "/orders/1", key: "shared-key", tenant: "tenant-a", user: "user-a", want: http.StatusConflict},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if code := performIdempotentRequest(router, tc.method, tc.path, tc.key, tc.tenant, tc.user); code != tc.want {
				t.Fatalf("status = %d, want %d", code, tc.want)
			}
		})
	}
}

// 中文：TestIdempotentIgnoresSafeMethods 验证相关行为符合预期。
// English: TestIdempotentIgnoresSafeMethods verifies the related behavior.
func TestIdempotentIgnoresSafeMethods(t *testing.T) {
	router := newIdempotentTestRouter()

	for i := 0; i < 2; i++ {
		if code := performIdempotentRequest(router, http.MethodGet, "/status", "safe-key", "tenant-a", "user-a"); code != http.StatusNoContent {
			t.Fatalf("GET request %d status = %d, want %d", i+1, code, http.StatusNoContent)
		}
	}
}

// 中文：newIdempotentTestRouter 执行当前包中的对应流程。
// English: newIdempotentTestRouter executes the corresponding workflow in this package.
func newIdempotentTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		ctx := c.Request.Context()
		if tenantID := c.GetHeader("X-Tenant-ID"); tenantID != "" {
			ctx = types.WithTenantID(ctx, tenantID)
		}
		if userID := c.GetHeader("X-User-ID"); userID != "" {
			ctx = types.WithUserID(ctx, userID)
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})
	router.Use(Idempotent(""))

	router.POST("/orders/:id", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	router.POST("/refunds", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	router.GET("/status", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	return router
}

// 中文：performIdempotentRequest 执行当前包中的对应流程。
// English: performIdempotentRequest executes the corresponding workflow in this package.
func performIdempotentRequest(router *gin.Engine, method, path, key, tenantID, userID string) int {
	req := httptest.NewRequest(method, path, nil)
	if key != "" {
		req.Header.Set("X-Idempotent-Key", key)
	}
	if tenantID != "" {
		req.Header.Set("X-Tenant-ID", tenantID)
	}
	if userID != "" {
		req.Header.Set("X-User-ID", userID)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec.Code
}
