package lock

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// 中文：RedisLock 定义当前包使用的数据结构或接口。
// English: RedisLock defines a data structure or interface used by this package.
// RedisLock Redis分布式锁实现
type RedisLock struct {
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *redis.Client
}

// 中文：NewRedisLock 创建并返回对应组件实例。
// English: NewRedisLock creates and returns the corresponding component instance.
// NewRedisLock 创建Redis分布式锁
func NewRedisLock(client *redis.Client) *RedisLock {
	return &RedisLock{client: client}
}

// 中文：redisLockHolder 定义当前包使用的数据结构或接口。
// English: redisLockHolder defines a data structure or interface used by this package.
type redisLockHolder struct {
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *redis.Client
	// 中文：key 保存当前结构中的配置或数据值。
	// English: key stores a configuration or data value for this struct.
	key string
	// 中文：value 保存当前结构中的配置或数据值。
	// English: value stores a configuration or data value for this struct.
	value string
}

// 中文：Lock 执行当前包中的对应流程。
// English: Lock executes the corresponding workflow in this package.
func (r *RedisLock) Lock(ctx context.Context, key string, ttl time.Duration) (LockHolder, error) {
	value := fmt.Sprintf("%d", time.Now().UnixNano())
	for {
		ok, err := r.client.SetNX(ctx, key, value, ttl).Result()
		if err != nil {
			return nil, fmt.Errorf("redis lock: %w", err)
		}
		if ok {
			return &redisLockHolder{client: r.client, key: key, value: value}, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
}

// 中文：TryLock 执行当前包中的对应流程。
// English: TryLock executes the corresponding workflow in this package.
func (r *RedisLock) TryLock(ctx context.Context, key string, ttl time.Duration) (LockHolder, error) {
	value := fmt.Sprintf("%d", time.Now().UnixNano())
	ok, err := r.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return nil, fmt.Errorf("redis trylock: %w", err)
	}
	if !ok {
		return nil, ErrLockFailed
	}
	return &redisLockHolder{client: r.client, key: key, value: value}, nil
}

// 中文：Unlock 执行当前包中的对应流程。
// English: Unlock executes the corresponding workflow in this package.
func (h *redisLockHolder) Unlock(ctx context.Context) error {
	// 使用Lua脚本确保只释放自己持有的锁
	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`)
	_, err := script.Run(ctx, h.client, []string{h.key}, h.value).Result()
	return err
}

// 中文：Renew 执行当前包中的对应流程。
// English: Renew executes the corresponding workflow in this package.
func (h *redisLockHolder) Renew(ctx context.Context, ttl time.Duration) error {
	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("expire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`)
	_, err := script.Run(ctx, h.client, []string{h.key}, h.value, int(ttl.Seconds())).Result()
	return err
}

// 中文：ErrLockFailed 声明当前包使用的变量。
// English: ErrLockFailed declares variables used by this package.
// ErrLockFailed 获取锁失败
var ErrLockFailed = &lockError{msg: "lock failed"}

// 中文：lockError 定义当前包使用的数据结构或接口。
// English: lockError defines a data structure or interface used by this package.
type lockError struct {
	// 中文：msg 保存当前结构中的配置或数据值。
	// English: msg stores a configuration or data value for this struct.
	msg string
}

// 中文：Error 执行当前包中的对应流程。
// English: Error executes the corresponding workflow in this package.
func (e *lockError) Error() string { return e.msg }
