package cache

import (
	"context"
	"time"
)

// 中文：MultiLevelCache 定义当前包使用的数据结构或接口。
// English: MultiLevelCache defines a data structure or interface used by this package.
// MultiLevelCache 多级缓存（L1内存 + L2 Redis）
type MultiLevelCache struct {
	// 中文：l1 保存当前结构中的配置或数据值。
	// English: l1 stores a configuration or data value for this struct.
	l1 Cache // 进程内缓存
	// 中文：l2 保存当前结构中的配置或数据值。
	// English: l2 stores a configuration or data value for this struct.
	l2 Cache // 分布式缓存
}

// 中文：NewMultiLevelCache 创建并返回对应组件实例。
// English: NewMultiLevelCache creates and returns the corresponding component instance.
// NewMultiLevelCache 创建多级缓存
func NewMultiLevelCache(l1, l2 Cache) *MultiLevelCache {
	return &MultiLevelCache{l1: l1, l2: l2}
}

// 中文：Get 执行当前包中的对应流程。
// English: Get executes the corresponding workflow in this package.
func (m *MultiLevelCache) Get(ctx context.Context, key string, dest any) error {
	// 先查L1
	err := m.l1.Get(ctx, key, dest)
	if err == nil {
		return nil
	}

	// L1未命中，查L2
	err = m.l2.Get(ctx, key, dest)
	if err != nil {
		return err
	}

	// 回写L1（短TTL）
	_ = m.l1.Set(ctx, key, dest, 5*time.Minute)
	return nil
}

// 中文：Set 执行当前包中的对应流程。
// English: Set executes the corresponding workflow in this package.
func (m *MultiLevelCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	// 同时写入L1和L2
	_ = m.l1.Set(ctx, key, value, ttl)
	return m.l2.Set(ctx, key, value, ttl)
}

// 中文：MGet 执行当前包中的对应流程。
// English: MGet executes the corresponding workflow in this package.
func (m *MultiLevelCache) MGet(ctx context.Context, keys []string, dest any) error {
	values := make(map[string]any, len(keys))
	var l1Values map[string]any
	if err := m.l1.MGet(ctx, keys, &l1Values); err != nil {
		return err
	}
	for key, value := range l1Values {
		values[key] = value
	}

	missing := make([]string, 0, len(keys)-len(values))
	for _, key := range keys {
		if _, ok := values[key]; !ok {
			missing = append(missing, key)
		}
	}
	if len(missing) == 0 {
		return assignCacheMap(dest, values)
	}

	var l2Values map[string]any
	if err := m.l2.MGet(ctx, missing, &l2Values); err != nil {
		return err
	}
	for key, value := range l2Values {
		values[key] = value
		_ = m.l1.Set(ctx, key, value, 5*time.Minute)
	}
	return assignCacheMap(dest, values)
}

// 中文：MSet 执行当前包中的对应流程。
// English: MSet executes the corresponding workflow in this package.
func (m *MultiLevelCache) MSet(ctx context.Context, kv map[string]any, ttl time.Duration) error {
	_ = m.l1.MSet(ctx, kv, ttl)
	return m.l2.MSet(ctx, kv, ttl)
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (m *MultiLevelCache) Delete(ctx context.Context, keys ...string) error {
	_ = m.l1.Delete(ctx, keys...)
	return m.l2.Delete(ctx, keys...)
}

// 中文：Exists 执行当前包中的对应流程。
// English: Exists executes the corresponding workflow in this package.
func (m *MultiLevelCache) Exists(ctx context.Context, key string) (bool, error) {
	if ok, _ := m.l1.Exists(ctx, key); ok {
		return true, nil
	}
	return m.l2.Exists(ctx, key)
}

// 中文：Expire 执行当前包中的对应流程。
// English: Expire executes the corresponding workflow in this package.
func (m *MultiLevelCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	_ = m.l1.Expire(ctx, key, ttl)
	return m.l2.Expire(ctx, key, ttl)
}

// 中文：Incr 执行当前包中的对应流程。
// English: Incr executes the corresponding workflow in this package.
func (m *MultiLevelCache) Incr(ctx context.Context, key string) (int64, error) {
	return m.l2.Incr(ctx, key)
}

// 中文：IncrBy 执行当前包中的对应流程。
// English: IncrBy executes the corresponding workflow in this package.
func (m *MultiLevelCache) IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	return m.l2.IncrBy(ctx, key, delta)
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (m *MultiLevelCache) Close() error {
	_ = m.l1.Close()
	return m.l2.Close()
}
