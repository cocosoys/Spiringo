package middleware

import (
	"context"
	"testing"
	"time"
)

// 中文：TestTokenBucketLimiterIsPerKey 验证相关行为符合预期。
// English: TestTokenBucketLimiterIsPerKey verifies the related behavior.
func TestTokenBucketLimiterIsPerKey(t *testing.T) {
	limiter := NewTokenBucketLimiter(1, 1)
	ctx := context.Background()

	if !limiter.Allow(ctx, "client-a") {
		t.Fatal("first client-a request should pass")
	}
	if limiter.Allow(ctx, "client-a") {
		t.Fatal("second immediate client-a request should be limited")
	}
	if !limiter.Allow(ctx, "client-b") {
		t.Fatal("client-b should have an independent bucket")
	}
}

// 中文：TestSlidingWindowLimiterDefaults 验证相关行为符合预期。
// English: TestSlidingWindowLimiterDefaults verifies the related behavior.
func TestSlidingWindowLimiterDefaults(t *testing.T) {
	limiter := NewSlidingWindowLimiter(0, 0)
	if !limiter.Allow(context.Background(), "client") {
		t.Fatal("default sliding window should allow first request")
	}
	if limiter.Allow(context.Background(), "client") {
		t.Fatal("default sliding window should limit second request in same window")
	}
}

// 中文：TestLeakyBucketLimiterRefills 验证相关行为符合预期。
// English: TestLeakyBucketLimiterRefills verifies the related behavior.
func TestLeakyBucketLimiterRefills(t *testing.T) {
	limiter := NewLeakyBucketLimiter(100, 1)
	ctx := context.Background()
	if !limiter.Allow(ctx, "client") {
		t.Fatal("first request should pass")
	}
	if limiter.Allow(ctx, "client") {
		t.Fatal("second immediate request should be limited")
	}
	time.Sleep(20 * time.Millisecond)
	if !limiter.Allow(ctx, "client") {
		t.Fatal("bucket should leak enough to allow another request")
	}
}
