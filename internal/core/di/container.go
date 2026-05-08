package di

import (
	"fmt"
	"reflect"
	"sync"
)

// 中文：Container 定义当前包使用的数据结构或接口。
// English: Container defines a data structure or interface used by this package.
// Container is a small concurrency-safe dependency container.
type Container struct {
	// 中文：instances 保存当前结构中的配置或数据值。
	// English: instances stores a configuration or data value for this struct.
	instances map[reflect.Type]any
	// 中文：named 保存当前结构中的配置或数据值。
	// English: named stores a configuration or data value for this struct.
	named map[string]any
	// 中文：factories 保存当前结构中的配置或数据值。
	// English: factories stores a configuration or data value for this struct.
	factories map[reflect.Type]func() any
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.RWMutex
}

// 中文：NewContainer 创建并返回对应组件实例。
// English: NewContainer creates and returns the corresponding component instance.
// NewContainer creates an empty dependency container.
func NewContainer() *Container {
	return &Container{
		instances: make(map[reflect.Type]any),
		named:     make(map[string]any),
		factories: make(map[reflect.Type]func() any),
	}
}

// 中文：Provide 执行当前包中的对应流程。
// English: Provide executes the corresponding workflow in this package.
// Provide registers a singleton instance under its concrete type.
func (c *Container) Provide(instance any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.instances[reflect.TypeOf(instance)] = instance
}

// 中文：ProvideAs 执行当前包中的对应流程。
// English: ProvideAs executes the corresponding workflow in this package.
// ProvideAs registers an instance under an explicit abstraction type.
func ProvideAs[T any](c *Container, instance T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.instances[typeOf[T]()] = instance
}

// 中文：ProvideNamed 执行当前包中的对应流程。
// English: ProvideNamed executes the corresponding workflow in this package.
// ProvideNamed registers an instance by name.
func (c *Container) ProvideNamed(name string, instance any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.named[name] = instance
}

// 中文：ProvideFactory 执行当前包中的对应流程。
// English: ProvideFactory executes the corresponding workflow in this package.
// ProvideFactory registers a lazy factory for a type.
func (c *Container) ProvideFactory(instanceType reflect.Type, factory func() any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.factories[instanceType] = factory
}

// 中文：Resolve 执行当前包中的对应流程。
// English: Resolve executes the corresponding workflow in this package.
// Resolve resolves an instance by type.
func Resolve[T any](c *Container) (T, error) {
	var zero T
	targetType := typeOf[T]()

	c.mu.RLock()
	if inst, ok := c.instances[targetType]; ok {
		c.mu.RUnlock()
		return inst.(T), nil
	}
	if factory, ok := c.factories[targetType]; ok {
		c.mu.RUnlock()
		result := factory()
		if t, ok := result.(T); ok {
			return t, nil
		}
		return zero, fmt.Errorf("factory for type %v returned wrong type", targetType)
	}
	c.mu.RUnlock()

	return zero, fmt.Errorf("no instance registered for type %v", targetType)
}

// 中文：MustResolve 执行当前包中的对应流程。
// English: MustResolve executes the corresponding workflow in this package.
// MustResolve resolves an instance and panics on failure.
func MustResolve[T any](c *Container) T {
	inst, err := Resolve[T](c)
	if err != nil {
		panic(err)
	}
	return inst
}

// 中文：ResolveNamed 执行当前包中的对应流程。
// English: ResolveNamed executes the corresponding workflow in this package.
// ResolveNamed resolves an instance by name.
func ResolveNamed[T any](c *Container, name string) (T, error) {
	var zero T

	c.mu.RLock()
	defer c.mu.RUnlock()

	if inst, ok := c.named[name]; ok {
		if t, ok := inst.(T); ok {
			return t, nil
		}
		return zero, fmt.Errorf("instance named %q is not of type %T", name, zero)
	}
	return zero, fmt.Errorf("no instance registered with name %q", name)
}

// 中文：MustResolveNamed 执行当前包中的对应流程。
// English: MustResolveNamed executes the corresponding workflow in this package.
// MustResolveNamed resolves a named instance and panics on failure.
func MustResolveNamed[T any](c *Container, name string) T {
	inst, err := ResolveNamed[T](c, name)
	if err != nil {
		panic(err)
	}
	return inst
}

// 中文：Has 执行当前包中的对应流程。
// English: Has executes the corresponding workflow in this package.
// Has checks whether a type has been registered.
func Has[T any](c *Container) bool {
	targetType := typeOf[T]()

	c.mu.RLock()
	defer c.mu.RUnlock()

	if _, ok := c.instances[targetType]; ok {
		return true
	}
	_, ok := c.factories[targetType]
	return ok
}

// 中文：HasNamed 执行当前包中的对应流程。
// English: HasNamed executes the corresponding workflow in this package.
// HasNamed checks whether a named instance has been registered.
func (c *Container) HasNamed(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.named[name]
	return ok
}

// 中文：typeOf 执行当前包中的对应流程。
// English: typeOf executes the corresponding workflow in this package.
func typeOf[T any]() reflect.Type {
	var zero T
	if t := reflect.TypeOf(zero); t != nil {
		return t
	}
	return reflect.TypeOf((*T)(nil)).Elem()
}
