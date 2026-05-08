package di

import "fmt"

// 中文：Inject 执行当前包中的对应流程。
// English: Inject executes the corresponding workflow in this package.
// Inject resolves T into a non-nil pointer target.
func Inject[T any](c *Container, target *T) error {
	if c == nil {
		return fmt.Errorf("di container is nil")
	}
	if target == nil {
		return fmt.Errorf("inject target is nil")
	}
	value, err := Resolve[T](c)
	if err != nil {
		return err
	}
	*target = value
	return nil
}

// 中文：MustInject 执行当前包中的对应流程。
// English: MustInject executes the corresponding workflow in this package.
// MustInject resolves T into target and panics on failure.
func MustInject[T any](c *Container, target *T) {
	if err := Inject(c, target); err != nil {
		panic(err)
	}
}

// 中文：InjectNamed 执行当前包中的对应流程。
// English: InjectNamed executes the corresponding workflow in this package.
// InjectNamed resolves a named dependency into a non-nil pointer target.
func InjectNamed[T any](c *Container, name string, target *T) error {
	if c == nil {
		return fmt.Errorf("di container is nil")
	}
	if target == nil {
		return fmt.Errorf("inject target is nil")
	}
	value, err := ResolveNamed[T](c, name)
	if err != nil {
		return err
	}
	*target = value
	return nil
}

// 中文：MustInjectNamed 执行当前包中的对应流程。
// English: MustInjectNamed executes the corresponding workflow in this package.
// MustInjectNamed resolves a named dependency into target and panics on failure.
func MustInjectNamed[T any](c *Container, name string, target *T) {
	if err := InjectNamed(c, name, target); err != nil {
		panic(err)
	}
}
