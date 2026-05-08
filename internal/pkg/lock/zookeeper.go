package lock

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

// 中文：ZooKeeperConfig 定义当前包使用的数据结构或接口。
// English: ZooKeeperConfig defines a data structure or interface used by this package.
type ZooKeeperConfig struct {
	// 中文：Servers 保存当前结构中的配置或数据值。
	// English: Servers stores a configuration or data value for this struct.
	Servers []string
	// 中文：Root 保存当前结构中的配置或数据值。
	// English: Root stores a configuration or data value for this struct.
	Root string
	// 中文：SessionTimeout 保存当前结构中的配置或数据值。
	// English: SessionTimeout stores a configuration or data value for this struct.
	SessionTimeout time.Duration
}

// 中文：ZooKeeperLock 定义当前包使用的数据结构或接口。
// English: ZooKeeperLock defines a data structure or interface used by this package.
type ZooKeeperLock struct {
	// 中文：conn 保存当前结构中的配置或数据值。
	// English: conn stores a configuration or data value for this struct.
	conn *zk.Conn
	// 中文：root 保存当前结构中的配置或数据值。
	// English: root stores a configuration or data value for this struct.
	root string
	// 中文：acl 保存当前结构中的配置或数据值。
	// English: acl stores a configuration or data value for this struct.
	acl []zk.ACL
}

// 中文：zooKeeperLockHolder 定义当前包使用的数据结构或接口。
// English: zooKeeperLockHolder defines a data structure or interface used by this package.
type zooKeeperLockHolder struct {
	// 中文：conn 保存当前结构中的配置或数据值。
	// English: conn stores a configuration or data value for this struct.
	conn *zk.Conn
	// 中文：path 保存当前结构中的配置或数据值。
	// English: path stores a configuration or data value for this struct.
	path string
}

// 中文：NewZooKeeperLock 创建并返回对应组件实例。
// English: NewZooKeeperLock creates and returns the corresponding component instance.
func NewZooKeeperLock(cfg ZooKeeperConfig) (*ZooKeeperLock, error) {
	if len(cfg.Servers) == 0 {
		return nil, fmt.Errorf("zookeeper servers are required")
	}
	if cfg.Root == "" {
		cfg.Root = "/spiringo/locks"
	}
	if cfg.SessionTimeout <= 0 {
		cfg.SessionTimeout = 5 * time.Second
	}

	conn, events, err := zk.Connect(cfg.Servers, cfg.SessionTimeout, zk.WithLogInfo(false))
	if err != nil {
		return nil, fmt.Errorf("zookeeper connect: %w", err)
	}
	if err := waitZooKeeperSession(events, cfg.SessionTimeout); err != nil {
		conn.Close()
		return nil, err
	}

	manager := &ZooKeeperLock{conn: conn, root: path.Clean(cfg.Root), acl: zk.WorldACL(zk.PermAll)}
	if err := manager.ensurePath(manager.root); err != nil {
		conn.Close()
		return nil, err
	}
	return manager, nil
}

// 中文：Lock 执行当前包中的对应流程。
// English: Lock executes the corresponding workflow in this package.
func (z *ZooKeeperLock) Lock(ctx context.Context, key string, ttl time.Duration) (LockHolder, error) {
	for {
		holder, err := z.TryLock(ctx, key, ttl)
		if err == nil {
			return holder, nil
		}
		if !errors.Is(err, ErrLockFailed) {
			return nil, err
		}

		lockPath := z.lockPath(key)
		_, _, watch, err := z.conn.ExistsW(lockPath)
		if err != nil && err != zk.ErrNoNode {
			return nil, fmt.Errorf("zookeeper watch lock: %w", err)
		}
		if err == zk.ErrNoNode {
			continue
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-watch:
		}
	}
}

// 中文：TryLock 执行当前包中的对应流程。
// English: TryLock executes the corresponding workflow in this package.
func (z *ZooKeeperLock) TryLock(ctx context.Context, key string, _ time.Duration) (LockHolder, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	lockPath := z.lockPath(key)
	created, err := z.conn.Create(lockPath, []byte(time.Now().Format(time.RFC3339Nano)), zk.FlagEphemeral, z.acl)
	if err == zk.ErrNodeExists {
		return nil, ErrLockFailed
	}
	if err != nil {
		return nil, fmt.Errorf("zookeeper create lock: %w", err)
	}
	return &zooKeeperLockHolder{conn: z.conn, path: created}, nil
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (z *ZooKeeperLock) Close() error {
	if z.conn != nil {
		z.conn.Close()
	}
	return nil
}

// 中文：Unlock 执行当前包中的对应流程。
// English: Unlock executes the corresponding workflow in this package.
func (h *zooKeeperLockHolder) Unlock(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := h.conn.Delete(h.path, -1)
	if err == zk.ErrNoNode {
		return nil
	}
	return err
}

// 中文：Renew 执行当前包中的对应流程。
// English: Renew executes the corresponding workflow in this package.
func (h *zooKeeperLockHolder) Renew(ctx context.Context, _ time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	ok, _, err := h.conn.Exists(h.path)
	if err != nil {
		return err
	}
	if !ok {
		return ErrLockFailed
	}
	return nil
}

// 中文：lockPath 执行当前包中的对应流程。
// English: lockPath executes the corresponding workflow in this package.
func (z *ZooKeeperLock) lockPath(key string) string {
	name := strings.NewReplacer("/", "_", "\\", "_", ":", "_", " ", "_").Replace(strings.TrimSpace(key))
	if name == "" {
		name = "default"
	}
	return path.Join(z.root, name)
}

// 中文：ensurePath 执行当前包中的对应流程。
// English: ensurePath executes the corresponding workflow in this package.
func (z *ZooKeeperLock) ensurePath(p string) error {
	p = path.Clean(p)
	if p == "/" || p == "." {
		return nil
	}
	current := ""
	for _, part := range strings.Split(strings.TrimPrefix(p, "/"), "/") {
		current += "/" + part
		ok, _, err := z.conn.Exists(current)
		if err != nil {
			return fmt.Errorf("zookeeper check path %s: %w", current, err)
		}
		if ok {
			continue
		}
		if _, err := z.conn.Create(current, nil, zk.FlagPersistent, z.acl); err != nil && err != zk.ErrNodeExists {
			return fmt.Errorf("zookeeper create path %s: %w", current, err)
		}
	}
	return nil
}

// 中文：waitZooKeeperSession 执行当前包中的对应流程。
// English: waitZooKeeperSession executes the corresponding workflow in this package.
func waitZooKeeperSession(events <-chan zk.Event, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case event, ok := <-events:
			if !ok {
				return fmt.Errorf("zookeeper event stream closed")
			}
			if event.State == zk.StateConnected || event.State == zk.StateHasSession {
				return nil
			}
		case <-timer.C:
			return fmt.Errorf("zookeeper session timeout")
		}
	}
}
