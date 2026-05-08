package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：TestI18nUsesExplicitLanguageHeader 验证相关行为符合预期。
// English: TestI18nUsesExplicitLanguageHeader verifies the related behavior.
func TestI18nUsesExplicitLanguageHeader(t *testing.T) {
	router := gin.New()
	router.Use(I18n(DefaultI18nConfig()))
	router.GET("/", func(c *gin.Context) {
		if got := Language(c); got != "en-US" {
			t.Fatalf("Language() = %q, want en-US", got)
		}
		if got := types.GetLanguage(c.Request.Context()); got != "en-US" {
			t.Fatalf("context language = %q, want en-US", got)
		}
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Language", "en-US")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Content-Language"); got != "en-US" {
		t.Fatalf("Content-Language = %q, want en-US", got)
	}
}

// 中文：TestI18nNegotiatesAcceptLanguage 验证相关行为符合预期。
// English: TestI18nNegotiatesAcceptLanguage verifies the related behavior.
func TestI18nNegotiatesAcceptLanguage(t *testing.T) {
	router := gin.New()
	router.Use(I18n(DefaultI18nConfig()))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, Language(c))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Language", "fr-FR, en;q=0.8, zh-CN;q=0.2")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if got := rec.Body.String(); got != "en-US" {
		t.Fatalf("body = %q, want en-US", got)
	}
}

// 中文：TestI18nFallsBackToDefaultLanguage 验证相关行为符合预期。
// English: TestI18nFallsBackToDefaultLanguage verifies the related behavior.
func TestI18nFallsBackToDefaultLanguage(t *testing.T) {
	router := gin.New()
	router.Use(I18n(DefaultI18nConfig()))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, Language(c))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Language", "fr-FR")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if got := rec.Body.String(); got != "zh-CN" {
		t.Fatalf("body = %q, want zh-CN", got)
	}
}
