package cache

import (
	"context"
	"time"
)

// 中文：Cache 定义当前包使用的数据结构或接口。
// English: Cache defines a data structure or interface used by this package.
// Cache 缓存接口
type Cache interface {
	// 中文：Get 声明该接口需要实现的行为。
	// English: Get declares behavior required by this interface.
	Get(ctx context.Context, key string, dest any) error
	// 中文：Set 声明该接口需要实现的行为。
	// English: Set declares behavior required by this interface.
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	// 中文：Delete 声明该接口需要实现的行为。
	// English: Delete declares behavior required by this interface.
	Delete(ctx context.Context, keys ...string) error
	// 中文：Exists 声明该接口需要实现的行为。
	// English: Exists declares behavior required by this interface.
	Exists(ctx context.Context, key string) (bool, error)
	// 中文：MGet 声明该接口需要实现的行为。
	// English: MGet declares behavior required by this interface.
	MGet(ctx context.Context, keys []string, dest any) error
	// 中文：MSet 声明该接口需要实现的行为。
	// English: MSet declares behavior required by this interface.
	MSet(ctx context.Context, kv map[string]any, ttl time.Duration) error
	// 中文：Expire 声明该接口需要实现的行为。
	// English: Expire declares behavior required by this interface.
	Expire(ctx context.Context, key string, ttl time.Duration) error
	// 中文：Incr 声明该接口需要实现的行为。
	// English: Incr declares behavior required by this interface.
	Incr(ctx context.Context, key string) (int64, error)
	// 中文：IncrBy 声明该接口需要实现的行为。
	// English: IncrBy declares behavior required by this interface.
	IncrBy(ctx context.Context, key string, delta int64) (int64, error)
	// 中文：Close 声明该接口需要实现的行为。
	// English: Close declares behavior required by this interface.
	Close() error
}
