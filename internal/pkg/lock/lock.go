package lock

import (
	"context"
	"time"
)

// 中文：LockHolder 定义当前包使用的数据结构或接口。
// English: LockHolder defines a data structure or interface used by this package.
// LockHolder 锁持有者
type LockHolder interface {
	// 中文：Unlock 声明该接口需要实现的行为。
	// English: Unlock declares behavior required by this interface.
	Unlock(ctx context.Context) error
	// 中文：Renew 声明该接口需要实现的行为。
	// English: Renew declares behavior required by this interface.
	Renew(ctx context.Context, ttl time.Duration) error
}

// 中文：Lock 定义当前包使用的数据结构或接口。
// English: Lock defines a data structure or interface used by this package.
// Lock 分布式锁接口
type Lock interface {
	// 中文：Lock 声明该接口需要实现的行为。
	// English: Lock declares behavior required by this interface.
	Lock(ctx context.Context, key string, ttl time.Duration) (LockHolder, error)
	// 中文：TryLock 声明该接口需要实现的行为。
	// English: TryLock declares behavior required by this interface.
	TryLock(ctx context.Context, key string, ttl time.Duration) (LockHolder, error)
}
