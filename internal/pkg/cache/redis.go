package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// 中文：RedisCache 定义当前包使用的数据结构或接口。
// English: RedisCache defines a data structure or interface used by this package.
// RedisCache Redis缓存实现
type RedisCache struct {
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *redis.Client
}

// 中文：NewRedisCache 创建并返回对应组件实例。
// English: NewRedisCache creates and returns the corresponding component instance.
// NewRedisCache 创建Redis缓存
func NewRedisCache(addr, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisCache{client: client}
}

// 中文：NewRedisCacheFromClient 创建并返回对应组件实例。
// English: NewRedisCacheFromClient creates and returns the corresponding component instance.
// NewRedisCacheFromClient 从已有client创建
func NewRedisCacheFromClient(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// 中文：Client 执行当前包中的对应流程。
// English: Client executes the corresponding workflow in this package.
// Client 获取原始Redis客户端
func (r *RedisCache) Client() *redis.Client {
	return r.client
}

// 中文：Get 执行当前包中的对应流程。
// English: Get executes the corresponding workflow in this package.
func (r *RedisCache) Get(ctx context.Context, key string, dest any) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrKeyNotFound
		}
		return fmt.Errorf("redis get: %w", err)
	}
	if dp, ok := dest.(*any); ok {
		*dp = val
	}
	return nil
}

// 中文：Set 执行当前包中的对应流程。
// English: Set executes the corresponding workflow in this package.
func (r *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// 中文：MGet 执行当前包中的对应流程。
// English: MGet executes the corresponding workflow in this package.
func (r *RedisCache) MGet(ctx context.Context, keys []string, dest any) error {
	if len(keys) == 0 {
		return assignCacheMap(dest, map[string]any{})
	}
	raw, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return fmt.Errorf("redis mget: %w", err)
	}
	values := make(map[string]any, len(keys))
	for i, key := range keys {
		if i < len(raw) && raw[i] != nil {
			values[key] = raw[i]
		}
	}
	return assignCacheMap(dest, values)
}

// 中文：MSet 执行当前包中的对应流程。
// English: MSet executes the corresponding workflow in this package.
func (r *RedisCache) MSet(ctx context.Context, kv map[string]any, ttl time.Duration) error {
	if len(kv) == 0 {
		return nil
	}
	pipe := r.client.Pipeline()
	for key, value := range kv {
		if key == "" {
			continue
		}
		pipe.Set(ctx, key, value, ttl)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("redis mset: %w", err)
	}
	return nil
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// 中文：Exists 执行当前包中的对应流程。
// English: Exists executes the corresponding workflow in this package.
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

// 中文：Expire 执行当前包中的对应流程。
// English: Expire executes the corresponding workflow in this package.
func (r *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}

// 中文：Incr 执行当前包中的对应流程。
// English: Incr executes the corresponding workflow in this package.
func (r *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// 中文：IncrBy 执行当前包中的对应流程。
// English: IncrBy executes the corresponding workflow in this package.
func (r *RedisCache) IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	return r.client.IncrBy(ctx, key, delta).Result()
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// 中文：Ping 执行当前包中的对应流程。
// English: Ping executes the corresponding workflow in this package.
// Ping 检查连接
func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
