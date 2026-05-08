package cache

import (
	"context"
	"sync"
	"time"
)

// 中文：memoryItem 定义当前包使用的数据结构或接口。
// English: memoryItem defines a data structure or interface used by this package.
type memoryItem struct {
	// 中文：value 保存当前结构中的配置或数据值。
	// English: value stores a configuration or data value for this struct.
	value any
	// 中文：expiredAt 保存当前结构中的配置或数据值。
	// English: expiredAt stores a configuration or data value for this struct.
	expiredAt time.Time
}

// 中文：isExpired 执行当前包中的对应流程。
// English: isExpired executes the corresponding workflow in this package.
func (i *memoryItem) isExpired() bool {
	return !i.expiredAt.IsZero() && time.Now().After(i.expiredAt)
}

// 中文：MemoryCache 定义当前包使用的数据结构或接口。
// English: MemoryCache defines a data structure or interface used by this package.
// MemoryCache 内存缓存实现
type MemoryCache struct {
	// 中文：items 保存当前结构中的配置或数据值。
	// English: items stores a configuration or data value for this struct.
	items map[string]*memoryItem
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.RWMutex
}

// 中文：NewMemoryCache 创建并返回对应组件实例。
// English: NewMemoryCache creates and returns the corresponding component instance.
// NewMemoryCache 创建内存缓存
func NewMemoryCache() *MemoryCache {
	mc := &MemoryCache{
		items: make(map[string]*memoryItem),
	}
	go mc.cleanup()
	return mc
}

// 中文：Get 执行当前包中的对应流程。
// English: Get executes the corresponding workflow in this package.
func (m *MemoryCache) Get(_ context.Context, key string, dest any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.items[key]
	if !ok || item.isExpired() {
		return ErrKeyNotFound
	}

	// 简单赋值：dest需要是*any类型
	if dp, ok := dest.(*any); ok {
		*dp = item.value
	}
	return nil
}

// 中文：Set 执行当前包中的对应流程。
// English: Set executes the corresponding workflow in this package.
func (m *MemoryCache) Set(_ context.Context, key string, value any, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var expiredAt time.Time
	if ttl > 0 {
		expiredAt = time.Now().Add(ttl)
	}

	m.items[key] = &memoryItem{
		value:     value,
		expiredAt: expiredAt,
	}
	return nil
}

// 中文：MGet 执行当前包中的对应流程。
// English: MGet executes the corresponding workflow in this package.
func (m *MemoryCache) MGet(_ context.Context, keys []string, dest any) error {
	values := make(map[string]any, len(keys))

	m.mu.RLock()
	for _, key := range keys {
		item, ok := m.items[key]
		if !ok || item.isExpired() {
			continue
		}
		values[key] = item.value
	}
	m.mu.RUnlock()

	return assignCacheMap(dest, values)
}

// 中文：MSet 执行当前包中的对应流程。
// English: MSet executes the corresponding workflow in this package.
func (m *MemoryCache) MSet(_ context.Context, kv map[string]any, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var expiredAt time.Time
	if ttl > 0 {
		expiredAt = time.Now().Add(ttl)
	}
	for key, value := range kv {
		if key == "" {
			continue
		}
		m.items[key] = &memoryItem{
			value:     value,
			expiredAt: expiredAt,
		}
	}
	return nil
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (m *MemoryCache) Delete(_ context.Context, keys ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, key := range keys {
		delete(m.items, key)
	}
	return nil
}

// 中文：Exists 执行当前包中的对应流程。
// English: Exists executes the corresponding workflow in this package.
func (m *MemoryCache) Exists(_ context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.items[key]
	if !ok || item.isExpired() {
		return false, nil
	}
	return true, nil
}

// 中文：Expire 执行当前包中的对应流程。
// English: Expire executes the corresponding workflow in this package.
func (m *MemoryCache) Expire(_ context.Context, key string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.items[key]
	if !ok {
		return ErrKeyNotFound
	}
	item.expiredAt = time.Now().Add(ttl)
	return nil
}

// 中文：Incr 执行当前包中的对应流程。
// English: Incr executes the corresponding workflow in this package.
func (m *MemoryCache) Incr(_ context.Context, key string) (int64, error) {
	return m.IncrBy(context.Background(), key, 1)
}

// 中文：IncrBy 执行当前包中的对应流程。
// English: IncrBy executes the corresponding workflow in this package.
func (m *MemoryCache) IncrBy(_ context.Context, key string, delta int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.items[key]
	if !ok || item.isExpired() {
		m.items[key] = &memoryItem{value: delta}
		return delta, nil
	}

	var current int64
	if v, ok := item.value.(int64); ok {
		current = v
	}
	current += delta
	item.value = current
	return current, nil
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (m *MemoryCache) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = make(map[string]*memoryItem)
	return nil
}

// 中文：cleanup 执行当前包中的对应流程。
// English: cleanup executes the corresponding workflow in this package.
func (m *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for k, v := range m.items {
			if !v.expiredAt.IsZero() && now.After(v.expiredAt) {
				delete(m.items, k)
			}
		}
		m.mu.Unlock()
	}
}

// 中文：ErrKeyNotFound 声明当前包使用的变量。
// English: ErrKeyNotFound declares variables used by this package.
// ErrKeyNotFound 键不存在错误
var ErrKeyNotFound = &cacheError{msg: "key not found"}

// 中文：cacheError 定义当前包使用的数据结构或接口。
// English: cacheError defines a data structure or interface used by this package.
type cacheError struct {
	// 中文：msg 保存当前结构中的配置或数据值。
	// English: msg stores a configuration or data value for this struct.
	msg string
}

// 中文：Error 执行当前包中的对应流程。
// English: Error executes the corresponding workflow in this package.
func (e *cacheError) Error() string { return e.msg }
