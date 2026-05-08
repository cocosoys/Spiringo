package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：CircuitState 定义当前包使用的数据结构或接口。
// English: CircuitState defines a data structure or interface used by this package.
type CircuitState string

// 中文：CircuitClosed、CircuitOpen、CircuitHalfOpen 声明当前包使用的常量。
// English: CircuitClosed、CircuitOpen、CircuitHalfOpen declares constants used by this package.
const (
	CircuitClosed   CircuitState = "closed"
	CircuitOpen     CircuitState = "open"
	CircuitHalfOpen CircuitState = "half_open"
)

// 中文：CircuitBreakerConfig 定义当前包使用的数据结构或接口。
// English: CircuitBreakerConfig defines a data structure or interface used by this package.
type CircuitBreakerConfig struct {
	// 中文：FailureThreshold 保存当前结构中的配置或数据值。
	// English: FailureThreshold stores a configuration or data value for this struct.
	FailureThreshold float64
	// 中文：MinimumRequests 保存当前结构中的配置或数据值。
	// English: MinimumRequests stores a configuration or data value for this struct.
	MinimumRequests int
	// 中文：OpenTimeout 保存当前结构中的配置或数据值。
	// English: OpenTimeout stores a configuration or data value for this struct.
	OpenTimeout time.Duration
	// 中文：HalfOpenMaxRequest 保存当前结构中的配置或数据值。
	// English: HalfOpenMaxRequest stores a configuration or data value for this struct.
	HalfOpenMaxRequest int
}

// 中文：CircuitBreaker 定义当前包使用的数据结构或接口。
// English: CircuitBreaker defines a data structure or interface used by this package.
type CircuitBreaker struct {
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.Mutex
	// 中文：state 保存当前结构中的配置或数据值。
	// English: state stores a configuration or data value for this struct.
	state CircuitState
	// 中文：failureThreshold 保存当前结构中的配置或数据值。
	// English: failureThreshold stores a configuration or data value for this struct.
	failureThreshold float64
	// 中文：minimumRequests 保存当前结构中的配置或数据值。
	// English: minimumRequests stores a configuration or data value for this struct.
	minimumRequests int
	// 中文：openTimeout 保存当前结构中的配置或数据值。
	// English: openTimeout stores a configuration or data value for this struct.
	openTimeout time.Duration
	// 中文：halfOpenMaxRequest 保存当前结构中的配置或数据值。
	// English: halfOpenMaxRequest stores a configuration or data value for this struct.
	halfOpenMaxRequest int
	// 中文：openedAt 保存当前结构中的配置或数据值。
	// English: openedAt stores a configuration or data value for this struct.
	openedAt time.Time
	// 中文：total 保存当前结构中的配置或数据值。
	// English: total stores a configuration or data value for this struct.
	total int
	// 中文：failures 保存当前结构中的配置或数据值。
	// English: failures stores a configuration or data value for this struct.
	failures int
	// 中文：halfOpenInFlight 保存当前结构中的配置或数据值。
	// English: halfOpenInFlight stores a configuration or data value for this struct.
	halfOpenInFlight int
}

// 中文：NewCircuitBreaker 创建并返回对应组件实例。
// English: NewCircuitBreaker creates and returns the corresponding component instance.
func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	if cfg.FailureThreshold <= 0 || cfg.FailureThreshold > 1 {
		cfg.FailureThreshold = 0.5
	}
	if cfg.MinimumRequests <= 0 {
		cfg.MinimumRequests = 20
	}
	if cfg.OpenTimeout <= 0 {
		cfg.OpenTimeout = 30 * time.Second
	}
	if cfg.HalfOpenMaxRequest <= 0 {
		cfg.HalfOpenMaxRequest = 1
	}
	return &CircuitBreaker{
		state:              CircuitClosed,
		failureThreshold:   cfg.FailureThreshold,
		minimumRequests:    cfg.MinimumRequests,
		openTimeout:        cfg.OpenTimeout,
		halfOpenMaxRequest: cfg.HalfOpenMaxRequest,
	}
}

// 中文：CircuitBreak 执行当前包中的对应流程。
// English: CircuitBreak executes the corresponding workflow in this package.
func CircuitBreak(breaker *CircuitBreaker) gin.HandlerFunc {
	if breaker == nil {
		breaker = NewCircuitBreaker(CircuitBreakerConfig{})
	}
	return func(c *gin.Context) {
		if !breaker.allow() {
			c.JSON(http.StatusServiceUnavailable, types.Response{
				Code:    10010,
				Message: "service unavailable",
			})
			c.Abort()
			return
		}

		c.Next()
		breaker.record(c.Writer.Status() >= http.StatusInternalServerError)
	}
}

// 中文：allow 执行当前包中的对应流程。
// English: allow executes the corresponding workflow in this package.
func (b *CircuitBreaker) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	switch b.state {
	case CircuitOpen:
		if now.Sub(b.openedAt) < b.openTimeout {
			return false
		}
		b.state = CircuitHalfOpen
		b.halfOpenInFlight = 0
	case CircuitHalfOpen:
	default:
		b.state = CircuitClosed
	}

	if b.state == CircuitHalfOpen {
		if b.halfOpenInFlight >= b.halfOpenMaxRequest {
			return false
		}
		b.halfOpenInFlight++
	}
	return true
}

// 中文：record 执行当前包中的对应流程。
// English: record executes the corresponding workflow in this package.
func (b *CircuitBreaker) record(failed bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == CircuitHalfOpen {
		if b.halfOpenInFlight > 0 {
			b.halfOpenInFlight--
		}
		if failed {
			b.open()
			return
		}
		b.close()
		return
	}
	if b.state == CircuitOpen {
		return
	}

	b.total++
	if failed {
		b.failures++
	}
	if b.total >= b.minimumRequests && float64(b.failures)/float64(b.total) >= b.failureThreshold {
		b.open()
	}
}

// 中文：open 执行当前包中的对应流程。
// English: open executes the corresponding workflow in this package.
func (b *CircuitBreaker) open() {
	b.state = CircuitOpen
	b.openedAt = time.Now()
}

// 中文：close 执行当前包中的对应流程。
// English: close executes the corresponding workflow in this package.
func (b *CircuitBreaker) close() {
	b.state = CircuitClosed
	b.total = 0
	b.failures = 0
	b.halfOpenInFlight = 0
}
