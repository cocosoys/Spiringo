package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/core/config"
	"github.com/spiringo/spiringo/internal/core/middleware"
)

// 中文：Server 定义当前包使用的数据结构或接口。
// English: Server defines a data structure or interface used by this package.
// Server HTTP服务
type Server struct {
	// 中文：engine 保存当前结构中的配置或数据值。
	// English: engine stores a configuration or data value for this struct.
	engine *gin.Engine
	// 中文：server 保存当前结构中的配置或数据值。
	// English: server stores a configuration or data value for this struct.
	server *http.Server
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config config.ServerConfig
	// 中文：errCh 保存当前结构中的配置或数据值。
	// English: errCh stores a configuration or data value for this struct.
	errCh chan error
}

// 中文：New 创建并返回对应组件实例。
// English: New creates and returns the corresponding component instance.
// New 创建HTTP服务
func New(cfg config.ServerConfig, mwCfg config.MiddlewareConfig) *Server {
	// 设置Gin模式
	gin.SetMode(cfg.Mode)

	engine := gin.New()

	// 注册全局中间件
	setupMiddlewares(engine, mwCfg)

	return &Server{
		engine: engine,
		config: cfg,
	}
}

// 中文：Engine 执行当前包中的对应流程。
// English: Engine executes the corresponding workflow in this package.
// Engine 获取Gin引擎
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

// 中文：Start 执行当前包中的对应流程。
// English: Start executes the corresponding workflow in this package.
// Start 启动HTTP服务
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    s.config.Addr,
		Handler: s.engine,
	}
	s.errCh = make(chan error, 1)

	listener, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.config.Addr, err)
	}

	go func() {
		if err := s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			s.errCh <- fmt.Errorf("http server error: %w", err)
		}
		close(s.errCh)
	}()

	return nil
}

// 中文：Errors 执行当前包中的对应流程。
// English: Errors executes the corresponding workflow in this package.
func (s *Server) Errors() <-chan error {
	return s.errCh
}

// 中文：Stop 执行当前包中的对应流程。
// English: Stop executes the corresponding workflow in this package.
// Stop 优雅停止HTTP服务
func (s *Server) Stop(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

// 中文：Addr 执行当前包中的对应流程。
// English: Addr executes the corresponding workflow in this package.
// Addr 获取监听地址
func (s *Server) Addr() string {
	return s.config.Addr
}

// 中文：setupMiddlewares 执行当前包中的对应流程。
// English: setupMiddlewares executes the corresponding workflow in this package.
// setupMiddlewares 注册全局中间件
func setupMiddlewares(r *gin.Engine, cfg config.MiddlewareConfig) {
	// 1. Recovery（最外层，兜底panic）
	if cfg.Recovery {
		r.Use(middleware.Recovery())
	}

	// 2. RequestID
	if cfg.RequestID {
		r.Use(middleware.RequestID())
	}

	// 3. CORS
	if cfg.CORS {
		r.Use(middleware.CORS(middleware.DefaultCORSConfig()))
	}

	// 4. I18n
	if cfg.I18n {
		r.Use(middleware.I18n(middleware.DefaultI18nConfig()))
	}

	// 5. RateLimit
	if cfg.RateLimit.Enabled {
		limiter := middleware.NewRateLimiter(middleware.RateLimitConfig{
			Strategy: cfg.RateLimit.Strategy,
			Rate:     cfg.RateLimit.Rate,
			Burst:    cfg.RateLimit.Burst,
		})
		r.Use(middleware.RateLimit(limiter, middleware.IPKeyFunc))
	}

	// 6. Tenant
	if cfg.Tenant {
		r.Use(middleware.TenantMiddleware())
	}

	// 7. Auth
	if cfg.Auth.Enabled {
		r.Use(middleware.JWTAuthWithOptions(cfg.Auth.JWTSecret, middleware.AuthOptions{
			PublicPaths: cfg.Auth.PublicPaths,
		}))
	}

	// 8. Idempotent
	if cfg.Idempotent.Enabled {
		r.Use(middleware.Idempotent(cfg.Idempotent.Header))
	}

	// 9. CircuitBreak
	if cfg.CircuitBreak.Enabled {
		timeout, _ := time.ParseDuration(cfg.CircuitBreak.Timeout)
		r.Use(middleware.CircuitBreak(middleware.NewCircuitBreaker(middleware.CircuitBreakerConfig{
			FailureThreshold: cfg.CircuitBreak.Threshold,
			OpenTimeout:      timeout,
		})))
	}
}
