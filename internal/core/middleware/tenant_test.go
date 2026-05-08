package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：TestTenantContextRoundTrip 验证相关行为符合预期。
// English: TestTenantContextRoundTrip verifies the related behavior.
func TestTenantContextRoundTrip(t *testing.T) {
	ctx := WithTenant(context.Background(), &TenantContext{
		TenantID:   "tenant-1",
		TenantName: "Acme",
		Strategy:   StrategyDatabase,
		DBConn:     "tenant-db-1",
	})

	tc := FromContext(ctx)
	if tc == nil {
		t.Fatal("tenant context missing")
	}
	if tc.TenantID != "tenant-1" || tc.TenantName != "Acme" || tc.Strategy != StrategyDatabase || tc.DBConn != "tenant-db-1" {
		t.Fatalf("tenant context = %+v", tc)
	}
	if got := types.GetTenantID(ctx); got != "tenant-1" {
		t.Fatalf("legacy tenant id = %q", got)
	}
}

// 中文：TestTenantMiddlewareInjectsStructuredContext 验证相关行为符合预期。
// English: TestTenantMiddlewareInjectsStructuredContext verifies the related behavior.
func TestTenantMiddlewareInjectsStructuredContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TenantMiddleware())
	router.GET("/tenant", func(c *gin.Context) {
		tc := FromContext(c.Request.Context())
		if tc == nil {
			t.Fatal("tenant context missing")
		}
		if tc.TenantID != "tenant-1" || tc.Strategy != StrategySchema || tc.DBConn != "schema-1" {
			t.Fatalf("tenant context = %+v", tc)
		}
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/tenant", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")
	req.Header.Set("X-Tenant-Strategy", "schema")
	req.Header.Set("X-Tenant-DB-Conn", "schema-1")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

// 中文：TestTenantMiddlewareExtractsTenantFromSubdomain 验证相关行为符合预期。
// English: TestTenantMiddlewareExtractsTenantFromSubdomain verifies the related behavior.
func TestTenantMiddlewareExtractsTenantFromSubdomain(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TenantMiddleware())
	router.GET("/tenant", func(c *gin.Context) {
		tc := FromContext(c.Request.Context())
		if tc == nil {
			t.Fatal("tenant context missing")
		}
		if tc.TenantID != "acme" {
			t.Fatalf("tenant id = %q, want acme", tc.TenantID)
		}
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "http://acme.example.com/tenant", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

// 中文：TestTenantIDFromHostIgnoresReservedAndNonDomainHosts 验证相关行为符合预期。
// English: TestTenantIDFromHostIgnoresReservedAndNonDomainHosts verifies the related behavior.
func TestTenantIDFromHostIgnoresReservedAndNonDomainHosts(t *testing.T) {
	for _, host := range []string{
		"localhost:8080",
		"127.0.0.1:8080",
		"www.example.com",
		"api.example.com",
		"example.com",
	} {
		if got := tenantIDFromHost(host); got != "" {
			t.Fatalf("tenantIDFromHost(%q) = %q, want empty", host, got)
		}
	}
}
