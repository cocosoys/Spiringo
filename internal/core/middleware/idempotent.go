package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：IdempotentConfig 定义当前包使用的数据结构或接口。
// English: IdempotentConfig defines a data structure or interface used by this package.
// IdempotentConfig 幂等配置
type IdempotentConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled"`
	// 中文：Header 保存当前结构中的配置或数据值。
	// English: Header stores a configuration or data value for this struct.
	Header string `yaml:"header"`
}

// 中文：idempotentStore 定义当前包使用的数据结构或接口。
// English: idempotentStore defines a data structure or interface used by this package.
// idempotentStore 幂等键存储
type idempotentStore struct {
	// 中文：store 保存当前结构中的配置或数据值。
	// English: store stores a configuration or data value for this struct.
	store map[string]time.Time
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.RWMutex
	// 中文：ttl 保存当前结构中的配置或数据值。
	// English: ttl stores a configuration or data value for this struct.
	ttl time.Duration
}

// 中文：newIdempotentStore 执行当前包中的对应流程。
// English: newIdempotentStore executes the corresponding workflow in this package.
func newIdempotentStore(ttl time.Duration) *idempotentStore {
	s := &idempotentStore{
		store: make(map[string]time.Time),
		ttl:   ttl,
	}
	go s.cleanup()
	return s
}

// 中文：Set 执行当前包中的对应流程。
// English: Set executes the corresponding workflow in this package.
func (s *idempotentStore) Set(key string) bool {
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.purgeExpiredLocked(now)

	if _, exists := s.store[key]; exists {
		return false
	}
	s.store[key] = now
	return true
}

// 中文：purgeExpiredLocked 执行当前包中的对应流程。
// English: purgeExpiredLocked executes the corresponding workflow in this package.
func (s *idempotentStore) purgeExpiredLocked(now time.Time) {
	for k, t := range s.store {
		if now.Sub(t) > s.ttl {
			delete(s.store, k)
		}
	}
}

// 中文：cleanup 执行当前包中的对应流程。
// English: cleanup executes the corresponding workflow in this package.
func (s *idempotentStore) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		s.purgeExpiredLocked(time.Now())
		s.mu.Unlock()
	}
}

// 中文：Idempotent 执行当前包中的对应流程。
// English: Idempotent executes the corresponding workflow in this package.
// Idempotent 幂等性中间件
func Idempotent(header string) gin.HandlerFunc {
	store := newIdempotentStore(5 * time.Minute)
	if header == "" {
		header = "X-Idempotent-Key"
	}

	return func(c *gin.Context) {
		if isIdempotencySafeMethod(c.Request.Method) {
			c.Next()
			return
		}

		key := strings.TrimSpace(c.GetHeader(header))

		// 没有幂等键，跳过检查（非强制场景）
		if key == "" {
			c.Next()
			return
		}

		if !store.Set(idempotencyScope(c, header, key)) {
			c.JSON(http.StatusConflict, types.Response{
				Code:    10008,
				Message: "duplicate request",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 中文：isIdempotencySafeMethod 执行当前包中的对应流程。
// English: isIdempotencySafeMethod executes the corresponding workflow in this package.
func isIdempotencySafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

// 中文：idempotencyScope 执行当前包中的对应流程。
// English: idempotencyScope executes the corresponding workflow in this package.
func idempotencyScope(c *gin.Context, header, key string) string {
	routePath := c.FullPath()
	if routePath == "" {
		routePath = c.Request.URL.Path
	}

	ctx := c.Request.Context()
	parts := []string{
		c.Request.Method,
		routePath,
		types.GetTenantID(ctx),
		types.GetUserID(ctx),
		header,
		key,
	}
	return strings.Join(parts, "\x00")
}
