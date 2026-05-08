package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/pkg/metrics"
)

// 中文：TestMetricsMiddlewareRecordsRequest 验证相关行为符合预期。
// English: TestMetricsMiddlewareRecordsRequest verifies the related behavior.
func TestMetricsMiddlewareRecordsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	registry := metrics.NewRegistry("app")
	router := gin.New()
	router.Use(Metrics(registry))
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var out strings.Builder
	if err := registry.WritePrometheus(&out); err != nil {
		t.Fatalf("WritePrometheus returned error: %v", err)
	}
	text := out.String()
	if !strings.Contains(text, `app_http_requests_total{method="GET",path="/health",status="204"} 1`) {
		t.Fatalf("metrics output missing request counter:\n%s", text)
	}
	if !strings.Contains(text, `app_http_request_duration_seconds_count{method="GET",path="/health",status="204"} 1`) {
		t.Fatalf("metrics output missing duration count:\n%s", text)
	}
}
