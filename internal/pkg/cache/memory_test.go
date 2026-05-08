package cache

import (
	"context"
	"testing"
	"time"
)

// 中文：TestMemoryCache_SetAndGet 验证相关行为符合预期。
// English: TestMemoryCache_SetAndGet verifies the related behavior.
func TestMemoryCache_SetAndGet(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	err := c.Set(ctx, "key1", "value1", 0)
	if err != nil {
		t.Fatal(err)
	}

	var val any
	err = c.Get(ctx, "key1", &val)
	if err != nil {
		t.Fatal(err)
	}
	if val != "value1" {
		t.Errorf("expected 'value1', got '%v'", val)
	}
}

// 中文：TestMemoryCache_MSetAndMGet 验证相关行为符合预期。
// English: TestMemoryCache_MSetAndMGet verifies the related behavior.
func TestMemoryCache_MSetAndMGet(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	if err := c.MSet(ctx, map[string]any{
		"key1": "value1",
		"key2": "value2",
	}, time.Minute); err != nil {
		t.Fatal(err)
	}

	var got map[string]string
	if err := c.MGet(ctx, []string{"key1", "key2", "missing"}, &got); err != nil {
		t.Fatal(err)
	}
	if got["key1"] != "value1" || got["key2"] != "value2" {
		t.Fatalf("unexpected MGet result: %+v", got)
	}
	if _, ok := got["missing"]; ok {
		t.Fatalf("missing key should not be returned: %+v", got)
	}
}

// 中文：TestMemoryCache_GetNotFound 验证相关行为符合预期。
// English: TestMemoryCache_GetNotFound verifies the related behavior.
func TestMemoryCache_GetNotFound(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	var val any
	err := c.Get(ctx, "nonexistent", &val)
	if err == nil {
		t.Error("expected error for nonexistent key")
	}
}

// 中文：TestMemoryCache_Delete 验证相关行为符合预期。
// English: TestMemoryCache_Delete verifies the related behavior.
func TestMemoryCache_Delete(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	c.Set(ctx, "key1", "value1", 0)
	c.Delete(ctx, "key1")

	var val any
	err := c.Get(ctx, "key1", &val)
	if err == nil {
		t.Error("expected error after delete")
	}
}

// 中文：TestMemoryCache_Exists 验证相关行为符合预期。
// English: TestMemoryCache_Exists verifies the related behavior.
func TestMemoryCache_Exists(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	c.Set(ctx, "key1", "value1", 0)

	exists, _ := c.Exists(ctx, "key1")
	if !exists {
		t.Error("expected key1 to exist")
	}

	exists, _ = c.Exists(ctx, "key2")
	if exists {
		t.Error("expected key2 to not exist")
	}
}

// 中文：TestMemoryCache_TTLExpiration 验证相关行为符合预期。
// English: TestMemoryCache_TTLExpiration verifies the related behavior.
func TestMemoryCache_TTLExpiration(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	c.Set(ctx, "expiring", "value", 100*time.Millisecond)

	var val any
	err := c.Get(ctx, "expiring", &val)
	if err != nil || val != "value" {
		t.Error("key should exist before TTL expires")
	}

	time.Sleep(150 * time.Millisecond)

	err = c.Get(ctx, "expiring", &val)
	if err == nil {
		t.Error("expected key to expire after TTL")
	}
}

// 中文：TestMemoryCache_Incr 验证相关行为符合预期。
// English: TestMemoryCache_Incr verifies the related behavior.
func TestMemoryCache_Incr(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	c.Set(ctx, "counter", int64(0), 0)

	val, _ := c.Incr(ctx, "counter")
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}

	val, _ = c.Incr(ctx, "counter")
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
}

// 中文：TestMemoryCache_IncrBy 验证相关行为符合预期。
// English: TestMemoryCache_IncrBy verifies the related behavior.
func TestMemoryCache_IncrBy(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	c.Set(ctx, "counter", int64(10), 0)

	val, _ := c.IncrBy(ctx, "counter", 5)
	if val != 15 {
		t.Errorf("expected 15, got %d", val)
	}
}

// 中文：TestMemoryCache_Expire 验证相关行为符合预期。
// English: TestMemoryCache_Expire verifies the related behavior.
func TestMemoryCache_Expire(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	c.Set(ctx, "key1", "value1", 0)
	c.Expire(ctx, "key1", 50*time.Millisecond)

	var val any
	err := c.Get(ctx, "key1", &val)
	if err != nil {
		t.Error("key should exist before expiry")
	}

	time.Sleep(100 * time.Millisecond)
	err = c.Get(ctx, "key1", &val)
	if err == nil {
		t.Error("expected key to be expired")
	}
}

// 中文：TestMemoryCache_Close 验证相关行为符合预期。
// English: TestMemoryCache_Close verifies the related behavior.
func TestMemoryCache_Close(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	c.Set(ctx, "key1", "value1", 0)
	c.Close()

	var val any
	err := c.Get(ctx, "key1", &val)
	if err == nil {
		t.Error("expected error after close")
	}
}

// 中文：TestMultiLevelCache_MGetBackfillsL1 验证相关行为符合预期。
// English: TestMultiLevelCache_MGetBackfillsL1 verifies the related behavior.
func TestMultiLevelCache_MGetBackfillsL1(t *testing.T) {
	ctx := context.Background()
	l1 := NewMemoryCache()
	l2 := NewMemoryCache()
	c := NewMultiLevelCache(l1, l2)

	if err := l2.MSet(ctx, map[string]any{"key1": "value1"}, time.Minute); err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := c.MGet(ctx, []string{"key1", "missing"}, &got); err != nil {
		t.Fatal(err)
	}
	if got["key1"] != "value1" {
		t.Fatalf("unexpected MGet result: %+v", got)
	}

	var cached any
	if err := l1.Get(ctx, "key1", &cached); err != nil {
		t.Fatalf("expected L1 backfill: %v", err)
	}
	if cached != "value1" {
		t.Fatalf("L1 cached value = %v, want value1", cached)
	}
}

// 中文：_、_、_ 声明当前包使用的变量。
// English: _、_、_ declares variables used by this package.
var (
	_ Cache = (*MemoryCache)(nil)
	_ Cache = (*RedisCache)(nil)
	_ Cache = (*MultiLevelCache)(nil)
)
