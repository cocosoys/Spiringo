package cache

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"reflect"
	"sync"
	"time"
)

// 中文：emptyMarker 声明当前包使用的常量。
// English: emptyMarker declares constants used by this package.
const emptyMarker = "__spiringo_cache_empty__"

// 中文：Loader 定义当前包使用的数据结构或接口。
// English: Loader defines a data structure or interface used by this package.
type Loader func(ctx context.Context) (any, error)

// 中文：ProtectOptions 定义当前包使用的数据结构或接口。
// English: ProtectOptions defines a data structure or interface used by this package.
type ProtectOptions struct {
	// 中文：TTL 保存当前结构中的配置或数据值。
	// English: TTL stores a configuration or data value for this struct.
	TTL time.Duration
	// 中文：EmptyTTL 保存当前结构中的配置或数据值。
	// English: EmptyTTL stores a configuration or data value for this struct.
	EmptyTTL time.Duration
	// 中文：TTLJitter 保存当前结构中的配置或数据值。
	// English: TTLJitter stores a configuration or data value for this struct.
	TTLJitter time.Duration
	// 中文：CacheEmpty 保存当前结构中的配置或数据值。
	// English: CacheEmpty stores a configuration or data value for this struct.
	CacheEmpty bool
}

// 中文：Protector 定义当前包使用的数据结构或接口。
// English: Protector defines a data structure or interface used by this package.
type Protector struct {
	// 中文：cache 保存当前结构中的配置或数据值。
	// English: cache stores a configuration or data value for this struct.
	cache Cache
	// 中文：group 保存当前结构中的配置或数据值。
	// English: group stores a configuration or data value for this struct.
	group flightGroup
}

// 中文：NewProtector 创建并返回对应组件实例。
// English: NewProtector creates and returns the corresponding component instance.
func NewProtector(cache Cache) *Protector {
	return &Protector{cache: cache}
}

// 中文：GetOrLoad 执行当前包中的对应流程。
// English: GetOrLoad executes the corresponding workflow in this package.
func (p *Protector) GetOrLoad(ctx context.Context, key string, dest any, loader Loader, opts ProtectOptions) error {
	if p == nil || p.cache == nil {
		return fmt.Errorf("cache protector requires cache")
	}
	if key == "" {
		return fmt.Errorf("cache key is required")
	}
	if loader == nil {
		return fmt.Errorf("cache loader is required")
	}

	if ok, err := p.getCached(ctx, key, dest); err != nil || ok {
		return err
	}

	value, err := p.group.Do(key, func() (any, error) {
		var cached any
		if ok, err := p.getCached(ctx, key, &cached); err != nil || ok {
			return cached, err
		}

		loaded, err := loader(ctx)
		if err != nil {
			if opts.CacheEmpty && errors.Is(err, ErrKeyNotFound) {
				return emptyMarker, p.cache.Set(ctx, key, emptyMarker, opts.emptyTTL())
			}
			return nil, err
		}
		if loaded == nil {
			if opts.CacheEmpty {
				return emptyMarker, p.cache.Set(ctx, key, emptyMarker, opts.emptyTTL())
			}
			return nil, ErrKeyNotFound
		}
		if err := p.cache.Set(ctx, key, loaded, opts.ttl(key)); err != nil {
			return nil, fmt.Errorf("cache loaded value: %w", err)
		}
		return loaded, nil
	})
	if err != nil {
		return err
	}
	if isEmptyMarker(value) {
		return ErrKeyNotFound
	}
	return assign(dest, value)
}

// 中文：getCached 执行当前包中的对应流程。
// English: getCached executes the corresponding workflow in this package.
func (p *Protector) getCached(ctx context.Context, key string, dest any) (bool, error) {
	var cached any
	if err := p.cache.Get(ctx, key, &cached); err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return false, nil
		}
		return false, err
	}
	if isEmptyMarker(cached) {
		return true, ErrKeyNotFound
	}
	return true, assign(dest, cached)
}

// 中文：emptyTTL 执行当前包中的对应流程。
// English: emptyTTL executes the corresponding workflow in this package.
func (o ProtectOptions) emptyTTL() time.Duration {
	if o.EmptyTTL > 0 {
		return o.EmptyTTL
	}
	if o.TTL > 0 && o.TTL < time.Minute {
		return o.TTL
	}
	return time.Minute
}

// 中文：ttl 执行当前包中的对应流程。
// English: ttl executes the corresponding workflow in this package.
func (o ProtectOptions) ttl(key string) time.Duration {
	if o.TTL <= 0 || o.TTLJitter <= 0 {
		return o.TTL
	}
	return o.TTL + stableJitter(key, o.TTLJitter)
}

// 中文：stableJitter 执行当前包中的对应流程。
// English: stableJitter executes the corresponding workflow in this package.
func stableJitter(key string, max time.Duration) time.Duration {
	if max <= 0 {
		return 0
	}
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	return time.Duration(h.Sum64() % uint64(max))
}

// 中文：isEmptyMarker 执行当前包中的对应流程。
// English: isEmptyMarker executes the corresponding workflow in this package.
func isEmptyMarker(value any) bool {
	s, ok := value.(string)
	return ok && s == emptyMarker
}

// 中文：assign 执行当前包中的对应流程。
// English: assign executes the corresponding workflow in this package.
func assign(dest any, value any) error {
	if dest == nil {
		return nil
	}
	target := reflect.ValueOf(dest)
	if target.Kind() != reflect.Ptr || target.IsNil() {
		return fmt.Errorf("cache destination must be a non-nil pointer")
	}
	elem := target.Elem()
	if !elem.CanSet() {
		return fmt.Errorf("cache destination cannot be set")
	}
	if value == nil {
		elem.Set(reflect.Zero(elem.Type()))
		return nil
	}
	source := reflect.ValueOf(value)
	if source.Type().AssignableTo(elem.Type()) {
		elem.Set(source)
		return nil
	}
	if source.Type().ConvertibleTo(elem.Type()) {
		elem.Set(source.Convert(elem.Type()))
		return nil
	}
	if elem.Kind() == reflect.Interface && source.Type().Implements(elem.Type()) {
		elem.Set(source)
		return nil
	}
	return fmt.Errorf("cache value type %T cannot be assigned to %s", value, elem.Type())
}

// 中文：flightGroup 定义当前包使用的数据结构或接口。
// English: flightGroup defines a data structure or interface used by this package.
type flightGroup struct {
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.Mutex
	// 中文：calls 保存当前结构中的配置或数据值。
	// English: calls stores a configuration or data value for this struct.
	calls map[string]*flightCall
}

// 中文：flightCall 定义当前包使用的数据结构或接口。
// English: flightCall defines a data structure or interface used by this package.
type flightCall struct {
	// 中文：wg 保存当前结构中的配置或数据值。
	// English: wg stores a configuration or data value for this struct.
	wg sync.WaitGroup
	// 中文：value 保存当前结构中的配置或数据值。
	// English: value stores a configuration or data value for this struct.
	value any
	// 中文：err 保存当前结构中的配置或数据值。
	// English: err stores a configuration or data value for this struct.
	err error
}

// 中文：Do 执行当前包中的对应流程。
// English: Do executes the corresponding workflow in this package.
func (g *flightGroup) Do(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()
	if g.calls == nil {
		g.calls = make(map[string]*flightCall)
	}
	if c := g.calls[key]; c != nil {
		g.mu.Unlock()
		c.wg.Wait()
		return c.value, c.err
	}
	c := &flightCall{}
	c.wg.Add(1)
	g.calls[key] = c
	g.mu.Unlock()

	c.value, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.calls, key)
	g.mu.Unlock()
	return c.value, c.err
}
