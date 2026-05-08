package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：RateLimitConfig 定义当前包使用的数据结构或接口。
// English: RateLimitConfig defines a data structure or interface used by this package.
// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// 中文：Strategy 保存当前结构中的配置或数据值。
	// English: Strategy stores a configuration or data value for this struct.
	Strategy string `yaml:"strategy"` // token_bucket | sliding_window | leaky_bucket
	// 中文：Rate 保存当前结构中的配置或数据值。
	// English: Rate stores a configuration or data value for this struct.
	Rate int `yaml:"rate"` // 每秒请求数
	// 中文：Burst 保存当前结构中的配置或数据值。
	// English: Burst stores a configuration or data value for this struct.
	Burst int `yaml:"burst"` // 突发容量
}

// 中文：RateLimiter 定义当前包使用的数据结构或接口。
// English: RateLimiter defines a data structure or interface used by this package.
// RateLimiter 限流器接口
type RateLimiter interface {
	// 中文：Allow 声明该接口需要实现的行为。
	// English: Allow declares behavior required by this interface.
	Allow(ctx context.Context, key string) bool
}

// 中文：RateLimit 执行当前包中的对应流程。
// English: RateLimit executes the corresponding workflow in this package.
// RateLimit 限流中间件
// keyFunc 用于从请求中提取限流key（如IP、用户ID等）
func RateLimit(limiter RateLimiter, keyFunc func(c *gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c)
		if !limiter.Allow(c.Request.Context(), key) {
			c.JSON(http.StatusTooManyRequests, types.Response{
				Code:    10006,
				Message: "too many requests",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 中文：IPKeyFunc 执行当前包中的对应流程。
// English: IPKeyFunc executes the corresponding workflow in this package.
// IPKeyFunc 使用客户端IP作为限流key
func IPKeyFunc(c *gin.Context) string {
	return c.ClientIP()
}

// ---- 令牌桶实现 ----

// 中文：TokenBucketLimiter 定义当前包使用的数据结构或接口。
// English: TokenBucketLimiter defines a data structure or interface used by this package.
// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	// 中文：rate 保存当前结构中的配置或数据值。
	// English: rate stores a configuration or data value for this struct.
	rate float64
	// 中文：burst 保存当前结构中的配置或数据值。
	// English: burst stores a configuration or data value for this struct.
	burst int
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.Mutex
	// 中文：records 保存当前结构中的配置或数据值。
	// English: records stores a configuration or data value for this struct.
	records map[string]*tokenBucketRecord
}

// 中文：tokenBucketRecord 定义当前包使用的数据结构或接口。
// English: tokenBucketRecord defines a data structure or interface used by this package.
type tokenBucketRecord struct {
	// 中文：tokens 保存当前结构中的配置或数据值。
	// English: tokens stores a configuration or data value for this struct.
	tokens float64
	// 中文：lastTime 保存当前结构中的配置或数据值。
	// English: lastTime stores a configuration or data value for this struct.
	lastTime time.Time
}

// 中文：NewTokenBucketLimiter 创建并返回对应组件实例。
// English: NewTokenBucketLimiter creates and returns the corresponding component instance.
// NewTokenBucketLimiter 创建令牌桶限流器
func NewTokenBucketLimiter(rate int, burst int) *TokenBucketLimiter {
	if rate <= 0 {
		rate = 1
	}
	if burst <= 0 {
		burst = rate
	}
	return &TokenBucketLimiter{
		rate:    float64(rate),
		burst:   burst,
		records: make(map[string]*tokenBucketRecord),
	}
}

// 中文：Allow 执行当前包中的对应流程。
// English: Allow executes the corresponding workflow in this package.
func (l *TokenBucketLimiter) Allow(_ context.Context, key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if key == "" {
		key = "_global"
	}
	record, ok := l.records[key]
	if !ok {
		record = &tokenBucketRecord{tokens: float64(l.burst), lastTime: now}
		l.records[key] = record
	}
	elapsed := now.Sub(record.lastTime).Seconds()
	record.tokens += elapsed * l.rate
	if record.tokens > float64(l.burst) {
		record.tokens = float64(l.burst)
	}
	record.lastTime = now

	if record.tokens >= 1 {
		record.tokens--
		return true
	}
	return false
}

// ---- 滑动窗口实现 ----

// 中文：SlidingWindowLimiter 定义当前包使用的数据结构或接口。
// English: SlidingWindowLimiter defines a data structure or interface used by this package.
// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	// 中文：rate 保存当前结构中的配置或数据值。
	// English: rate stores a configuration or data value for this struct.
	rate int
	// 中文：window 保存当前结构中的配置或数据值。
	// English: window stores a configuration or data value for this struct.
	window time.Duration
	// 中文：records 保存当前结构中的配置或数据值。
	// English: records stores a configuration or data value for this struct.
	records map[string]*windowRecord
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.Mutex
}

// 中文：windowRecord 定义当前包使用的数据结构或接口。
// English: windowRecord defines a data structure or interface used by this package.
type windowRecord struct {
	// 中文：count 保存当前结构中的配置或数据值。
	// English: count stores a configuration or data value for this struct.
	count int
	// 中文：startAt 保存当前结构中的配置或数据值。
	// English: startAt stores a configuration or data value for this struct.
	startAt time.Time
}

// 中文：NewSlidingWindowLimiter 创建并返回对应组件实例。
// English: NewSlidingWindowLimiter creates and returns the corresponding component instance.
// NewSlidingWindowLimiter 创建滑动窗口限流器
func NewSlidingWindowLimiter(rate int, window time.Duration) *SlidingWindowLimiter {
	if rate <= 0 {
		rate = 1
	}
	if window <= 0 {
		window = time.Second
	}
	return &SlidingWindowLimiter{
		rate:    rate,
		window:  window,
		records: make(map[string]*windowRecord),
	}
}

// 中文：Allow 执行当前包中的对应流程。
// English: Allow executes the corresponding workflow in this package.
func (l *SlidingWindowLimiter) Allow(_ context.Context, key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	record, ok := l.records[key]
	if !ok || now.Sub(record.startAt) > l.window {
		l.records[key] = &windowRecord{count: 1, startAt: now}
		return true
	}

	if record.count >= l.rate {
		return false
	}

	record.count++
	return true
}

// 中文：LeakyBucketLimiter 定义当前包使用的数据结构或接口。
// English: LeakyBucketLimiter defines a data structure or interface used by this package.
// LeakyBucketLimiter 漏桶限流器（匀速出水）
type LeakyBucketLimiter struct {
	// 中文：rate 保存当前结构中的配置或数据值。
	// English: rate stores a configuration or data value for this struct.
	rate float64 // 每秒漏出水滴数
	// 中文：lastLeak 保存当前结构中的配置或数据值。
	// English: lastLeak stores a configuration or data value for this struct.
	lastLeak time.Time // 上次漏水时间
	// 中文：water 保存当前结构中的配置或数据值。
	// English: water stores a configuration or data value for this struct.
	water float64 // 当前桶中水量
	// 中文：capacity 保存当前结构中的配置或数据值。
	// English: capacity stores a configuration or data value for this struct.
	capacity float64 // 桶容量
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.Mutex
	// 中文：records 保存当前结构中的配置或数据值。
	// English: records stores a configuration or data value for this struct.
	records map[string]*leakyRecord
}

// 中文：leakyRecord 定义当前包使用的数据结构或接口。
// English: leakyRecord defines a data structure or interface used by this package.
type leakyRecord struct {
	// 中文：water 保存当前结构中的配置或数据值。
	// English: water stores a configuration or data value for this struct.
	water float64
	// 中文：lastLeak 保存当前结构中的配置或数据值。
	// English: lastLeak stores a configuration or data value for this struct.
	lastLeak time.Time
}

// 中文：NewLeakyBucketLimiter 创建并返回对应组件实例。
// English: NewLeakyBucketLimiter creates and returns the corresponding component instance.
// NewLeakyBucketLimiter 创建漏桶限流器
func NewLeakyBucketLimiter(rate float64, capacity float64) *LeakyBucketLimiter {
	if rate <= 0 {
		rate = 1
	}
	if capacity <= 0 {
		capacity = rate
	}
	return &LeakyBucketLimiter{
		rate:     rate,
		capacity: capacity,
		records:  make(map[string]*leakyRecord),
	}
}

// 中文：Allow 执行当前包中的对应流程。
// English: Allow executes the corresponding workflow in this package.
func (l *LeakyBucketLimiter) Allow(_ context.Context, key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	record, ok := l.records[key]
	if !ok {
		l.records[key] = &leakyRecord{water: 1, lastLeak: now}
		return 1 <= l.capacity
	}

	// 漏出水量
	elapsed := now.Sub(record.lastLeak).Seconds()
	record.water -= elapsed * l.rate
	if record.water < 0 {
		record.water = 0
	}
	record.lastLeak = now

	// 加一滴水
	record.water++
	return record.water <= l.capacity
}

// 中文：NewRateLimiter 创建并返回对应组件实例。
// English: NewRateLimiter creates and returns the corresponding component instance.
// NewRateLimiter 根据配置创建限流器
func NewRateLimiter(cfg RateLimitConfig) RateLimiter {
	switch cfg.Strategy {
	case "token_bucket":
		return NewTokenBucketLimiter(cfg.Rate, cfg.Burst)
	case "sliding_window":
		return NewSlidingWindowLimiter(cfg.Rate, time.Second)
	case "leaky_bucket":
		return NewLeakyBucketLimiter(float64(cfg.Rate), float64(cfg.Burst))
	default:
		return NewSlidingWindowLimiter(cfg.Rate, time.Second)
	}
}
