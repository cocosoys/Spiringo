package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// 中文：TestCircuitBreakOpensAfterFailureThreshold 验证相关行为符合预期。
// English: TestCircuitBreakOpensAfterFailureThreshold verifies the related behavior.
func TestCircuitBreakOpensAfterFailureThreshold(t *testing.T) {
	gin.SetMode(gin.TestMode)

	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 0.5,
		MinimumRequests:  2,
		OpenTimeout:      time.Hour,
	})
	router := gin.New()
	router.Use(CircuitBreak(breaker))
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})

	assertStatus(t, router, http.StatusInternalServerError)
	assertStatus(t, router, http.StatusInternalServerError)
	assertStatus(t, router, http.StatusServiceUnavailable)
}

// 中文：TestCircuitBreakHalfOpenSuccessClosesCircuit 验证相关行为符合预期。
// English: TestCircuitBreakHalfOpenSuccessClosesCircuit verifies the related behavior.
func TestCircuitBreakHalfOpenSuccessClosesCircuit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 1,
		MinimumRequests:  1,
		OpenTimeout:      time.Millisecond,
	})
	fail := true
	router := gin.New()
	router.Use(CircuitBreak(breaker))
	router.GET("/", func(c *gin.Context) {
		if fail {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})

	assertStatus(t, router, http.StatusInternalServerError)
	assertStatus(t, router, http.StatusServiceUnavailable)
	time.Sleep(2 * time.Millisecond)
	fail = false
	assertStatus(t, router, http.StatusOK)
	assertStatus(t, router, http.StatusOK)
}

// 中文：assertStatus 执行当前包中的对应流程。
// English: assertStatus executes the corresponding workflow in this package.
func assertStatus(t *testing.T, router http.Handler, want int) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != want {
		t.Fatalf("status = %d, want %d", resp.Code, want)
	}
}
