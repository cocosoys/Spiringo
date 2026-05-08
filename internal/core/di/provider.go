package di

import (
	"fmt"
	"reflect"
)

// 中文：Provider 定义当前包使用的数据结构或接口。
// English: Provider defines a data structure or interface used by this package.
// Provider registers one or more dependencies into a container.
type Provider interface {
	// 中文：Provide 声明该接口需要实现的行为。
	// English: Provide declares behavior required by this interface.
	Provide(c *Container) error
}

// 中文：ProviderFunc 定义当前包使用的数据结构或接口。
// English: ProviderFunc defines a data structure or interface used by this package.
// ProviderFunc adapts a function into a Provider.
type ProviderFunc func(c *Container) error

// 中文：Provide 执行当前包中的对应流程。
// English: Provide executes the corresponding workflow in this package.
func (f ProviderFunc) Provide(c *Container) error {
	if f == nil {
		return fmt.Errorf("di provider func is nil")
	}
	return f(c)
}

// 中文：Register 执行当前包中的对应流程。
// English: Register executes the corresponding workflow in this package.
// Register applies providers in order.
func Register(c *Container, providers ...Provider) error {
	if c == nil {
		return fmt.Errorf("di container is nil")
	}
	for _, provider := range providers {
		if provider == nil {
			continue
		}
		if err := provider.Provide(c); err != nil {
			return err
		}
	}
	return nil
}

// 中文：Singleton 执行当前包中的对应流程。
// English: Singleton executes the corresponding workflow in this package.
// Singleton registers an already-created dependency under its requested type.
func Singleton[T any](instance T) Provider {
	return ProviderFunc(func(c *Container) error {
		ProvideAs[T](c, instance)
		return nil
	})
}

// 中文：Concrete 执行当前包中的对应流程。
// English: Concrete executes the corresponding workflow in this package.
// Concrete registers an already-created dependency under its concrete type.
func Concrete(instance any) Provider {
	return ProviderFunc(func(c *Container) error {
		if instance == nil {
			return fmt.Errorf("di concrete instance is nil")
		}
		c.Provide(instance)
		return nil
	})
}

// 中文：Named 执行当前包中的对应流程。
// English: Named executes the corresponding workflow in this package.
// Named registers an instance under a string key.
func Named(name string, instance any) Provider {
	return ProviderFunc(func(c *Container) error {
		c.ProvideNamed(name, instance)
		return nil
	})
}

// 中文：Factory 执行当前包中的对应流程。
// English: Factory executes the corresponding workflow in this package.
// Factory registers a lazy factory under the requested type.
func Factory[T any](factory func() T) Provider {
	return ProviderFunc(func(c *Container) error {
		if factory == nil {
			return fmt.Errorf("di factory is nil")
		}
		c.ProvideFactory(typeOf[T](), func() any {
			return factory()
		})
		return nil
	})
}

// 中文：FactoryOf 执行当前包中的对应流程。
// English: FactoryOf executes the corresponding workflow in this package.
// FactoryOf registers a lazy factory under an explicit reflection type.
func FactoryOf(instanceType reflect.Type, factory func() any) Provider {
	return ProviderFunc(func(c *Container) error {
		if instanceType == nil {
			return fmt.Errorf("di factory type is nil")
		}
		if factory == nil {
			return fmt.Errorf("di factory is nil")
		}
		c.ProvideFactory(instanceType, factory)
		return nil
	})
}
