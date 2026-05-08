package cache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 中文：TestProtectorLoadsAndCachesValue 验证相关行为符合预期。
// English: TestProtectorLoadsAndCachesValue verifies the related behavior.
func TestProtectorLoadsAndCachesValue(t *testing.T) {
	ctx := context.Background()
	p := NewProtector(NewMemoryCache())
	var loads atomic.Int64

	var got string
	err := p.GetOrLoad(ctx, "user:1", &got, func(context.Context) (any, error) {
		loads.Add(1)
		return "alice", nil
	}, ProtectOptions{TTL: time.Minute})
	if err != nil {
		t.Fatalf("GetOrLoad returned error: %v", err)
	}
	if got != "alice" || loads.Load() != 1 {
		t.Fatalf("got=%q loads=%d", got, loads.Load())
	}

	got = ""
	if err := p.GetOrLoad(ctx, "user:1", &got, func(context.Context) (any, error) {
		loads.Add(1)
		return "bob", nil
	}, ProtectOptions{TTL: time.Minute}); err != nil {
		t.Fatalf("GetOrLoad cached returned error: %v", err)
	}
	if got != "alice" || loads.Load() != 1 {
		t.Fatalf("cached got=%q loads=%d", got, loads.Load())
	}
}

// 中文：TestProtectorSingleFlight 验证相关行为符合预期。
// English: TestProtectorSingleFlight verifies the related behavior.
func TestProtectorSingleFlight(t *testing.T) {
	ctx := context.Background()
	p := NewProtector(NewMemoryCache())
	var loads atomic.Int64
	start := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			var got string
			if err := p.GetOrLoad(ctx, "hot-key", &got, func(context.Context) (any, error) {
				loads.Add(1)
				time.Sleep(20 * time.Millisecond)
				return "value", nil
			}, ProtectOptions{TTL: time.Minute}); err != nil {
				t.Errorf("GetOrLoad returned error: %v", err)
			}
			if got != "value" {
				t.Errorf("got = %q", got)
			}
		}()
	}
	close(start)
	wg.Wait()

	if loads.Load() != 1 {
		t.Fatalf("loads = %d, want 1", loads.Load())
	}
}

// 中文：TestProtectorCachesEmptyMiss 验证相关行为符合预期。
// English: TestProtectorCachesEmptyMiss verifies the related behavior.
func TestProtectorCachesEmptyMiss(t *testing.T) {
	ctx := context.Background()
	p := NewProtector(NewMemoryCache())
	var loads atomic.Int64

	load := func(context.Context) (any, error) {
		loads.Add(1)
		return nil, ErrKeyNotFound
	}
	var got string
	if err := p.GetOrLoad(ctx, "missing", &got, load, ProtectOptions{CacheEmpty: true, EmptyTTL: time.Minute}); !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("first err = %v, want ErrKeyNotFound", err)
	}
	if err := p.GetOrLoad(ctx, "missing", &got, load, ProtectOptions{CacheEmpty: true, EmptyTTL: time.Minute}); !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("second err = %v, want ErrKeyNotFound", err)
	}
	if loads.Load() != 1 {
		t.Fatalf("loads = %d, want 1", loads.Load())
	}
}

// 中文：TestProtectOptionsAddsStableJitter 验证相关行为符合预期。
// English: TestProtectOptionsAddsStableJitter verifies the related behavior.
func TestProtectOptionsAddsStableJitter(t *testing.T) {
	opts := ProtectOptions{TTL: time.Minute, TTLJitter: 10 * time.Second}
	first := opts.ttl("same-key")
	second := opts.ttl("same-key")
	if first != second {
		t.Fatalf("jitter should be stable: %v != %v", first, second)
	}
	if first < time.Minute || first >= time.Minute+10*time.Second {
		t.Fatalf("ttl = %v outside jitter range", first)
	}
}
