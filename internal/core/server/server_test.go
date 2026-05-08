package server

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/core/config"
)

// 中文：TestStartReturnsListenError 验证相关行为符合预期。
// English: TestStartReturnsListenError verifies the related behavior.
func TestStartReturnsListenError(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	srv := New(config.ServerConfig{Addr: ln.Addr().String(), Mode: "test"}, config.MiddlewareConfig{})
	err = srv.Start()
	if err == nil {
		t.Fatal("Start returned nil error, want listen error")
	}
	if !strings.Contains(err.Error(), "listen") {
		t.Fatalf("err = %v, want listen error", err)
	}
}

// 中文：TestStartAndStop 验证相关行为符合预期。
// English: TestStartAndStop verifies the related behavior.
func TestStartAndStop(t *testing.T) {
	srv := New(config.ServerConfig{Addr: "127.0.0.1:0", Mode: "test"}, config.MiddlewareConfig{})
	if err := srv.Start(); err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		t.Fatalf("Stop returned error: %v", err)
	}
}

// 中文：TestNewInstallsI18nMiddleware 验证相关行为符合预期。
// English: TestNewInstallsI18nMiddleware verifies the related behavior.
func TestNewInstallsI18nMiddleware(t *testing.T) {
	srv := New(
		config.ServerConfig{Addr: "127.0.0.1:0", Mode: "test"},
		config.MiddlewareConfig{I18n: true},
	)
	srv.Engine().GET("/", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Language", "en-US")
	rec := httptest.NewRecorder()
	srv.Engine().ServeHTTP(rec, req)

	if got := rec.Header().Get("Content-Language"); got != "en-US" {
		t.Fatalf("Content-Language = %q, want en-US", got)
	}
}

// 中文：TestNewInstallsGlobalAuthMiddleware 验证相关行为符合预期。
// English: TestNewInstallsGlobalAuthMiddleware verifies the related behavior.
func TestNewInstallsGlobalAuthMiddleware(t *testing.T) {
	srv := New(
		config.ServerConfig{Addr: "127.0.0.1:0", Mode: "test"},
		config.MiddlewareConfig{
			Auth: config.GlobalAuthConfig{
				Enabled:     true,
				JWTSecret:   "secret",
				PublicPaths: []string{"/public"},
			},
		},
	)
	srv.Engine().GET("/public", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	srv.Engine().GET("/private", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	publicReq := httptest.NewRequest(http.MethodGet, "/public", nil)
	publicRec := httptest.NewRecorder()
	srv.Engine().ServeHTTP(publicRec, publicReq)
	if publicRec.Code != http.StatusNoContent {
		t.Fatalf("public status = %d, want %d", publicRec.Code, http.StatusNoContent)
	}

	privateReq := httptest.NewRequest(http.MethodGet, "/private", nil)
	privateRec := httptest.NewRecorder()
	srv.Engine().ServeHTTP(privateRec, privateReq)
	if privateRec.Code != http.StatusUnauthorized {
		t.Fatalf("private status = %d, want %d", privateRec.Code, http.StatusUnauthorized)
	}
}
