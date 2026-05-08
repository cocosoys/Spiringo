package di

import (
	"reflect"
	"testing"
)

// 中文：TestService 定义当前包使用的数据结构或接口。
// English: TestService defines a data structure or interface used by this package.
type TestService struct {
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string
}

// 中文：TestInterface 定义当前包使用的数据结构或接口。
// English: TestInterface defines a data structure or interface used by this package.
type TestInterface interface {
	// 中文：GetName 声明该接口需要实现的行为。
	// English: GetName declares behavior required by this interface.
	GetName() string
}

// 中文：GetName 执行当前包中的对应流程。
// English: GetName executes the corresponding workflow in this package.
func (s *TestService) GetName() string {
	return s.Name
}

// 中文：TestContainer_ProvideAndResolve 验证相关行为符合预期。
// English: TestContainer_ProvideAndResolve verifies the related behavior.
func TestContainer_ProvideAndResolve(t *testing.T) {
	c := NewContainer()
	svc := &TestService{Name: "hello"}
	c.Provide(svc)

	resolved, err := Resolve[*TestService](c)
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Name != "hello" {
		t.Errorf("expected 'hello', got '%s'", resolved.Name)
	}
}

// 中文：TestContainer_ResolveNotFound 验证相关行为符合预期。
// English: TestContainer_ResolveNotFound verifies the related behavior.
func TestContainer_ResolveNotFound(t *testing.T) {
	c := NewContainer()
	_, err := Resolve[*TestService](c)
	if err == nil {
		t.Error("expected error for unresolved type")
	}
}

// 中文：TestContainer_Has 验证相关行为符合预期。
// English: TestContainer_Has verifies the related behavior.
func TestContainer_Has(t *testing.T) {
	c := NewContainer()
	if Has[*TestService](c) {
		t.Error("expected Has to return false for empty container")
	}

	c.Provide(&TestService{Name: "test"})
	if !Has[*TestService](c) {
		t.Error("expected Has to return true after Provide")
	}
}

// 中文：TestContainer_ProvideNamed 验证相关行为符合预期。
// English: TestContainer_ProvideNamed verifies the related behavior.
func TestContainer_ProvideNamed(t *testing.T) {
	c := NewContainer()
	c.ProvideNamed("service_a", &TestService{Name: "A"})
	c.ProvideNamed("service_b", &TestService{Name: "B"})

	a, err := ResolveNamed[*TestService](c, "service_a")
	if err != nil {
		t.Fatal(err)
	}
	if a.Name != "A" {
		t.Errorf("expected 'A', got '%s'", a.Name)
	}
}

// 中文：TestContainer_ProvideFactory 验证相关行为符合预期。
// English: TestContainer_ProvideFactory verifies the related behavior.
func TestContainer_ProvideFactory(t *testing.T) {
	c := NewContainer()
	c.ProvideFactory(reflect.TypeOf(&TestService{}), func() any {
		return &TestService{Name: "factory"}
	})

	resolved, err := Resolve[*TestService](c)
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Name != "factory" {
		t.Errorf("expected 'factory', got '%s'", resolved.Name)
	}
}

// 中文：TestContainer_ProvideAsInterface 验证相关行为符合预期。
// English: TestContainer_ProvideAsInterface verifies the related behavior.
func TestContainer_ProvideAsInterface(t *testing.T) {
	c := NewContainer()
	ProvideAs[TestInterface](c, &TestService{Name: "interface"})

	resolved, err := Resolve[TestInterface](c)
	if err != nil {
		t.Fatal(err)
	}
	if resolved.GetName() != "interface" {
		t.Errorf("expected 'interface', got '%s'", resolved.GetName())
	}
	if !Has[TestInterface](c) {
		t.Error("expected Has to return true for explicitly provided interface")
	}
}

// 中文：TestContainer_HasNamed 验证相关行为符合预期。
// English: TestContainer_HasNamed verifies the related behavior.
func TestContainer_HasNamed(t *testing.T) {
	c := NewContainer()
	c.ProvideNamed("svc", &TestService{Name: "test"})

	if !c.HasNamed("svc") {
		t.Error("expected HasNamed to return true")
	}
	if c.HasNamed("nonexistent") {
		t.Error("expected HasNamed to return false for nonexistent name")
	}
}

// 中文：TestContainer_MustResolve_Panic 验证相关行为符合预期。
// English: TestContainer_MustResolve_Panic verifies the related behavior.
func TestContainer_MustResolve_Panic(t *testing.T) {
	c := NewContainer()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for MustResolve with unregistered type")
		}
	}()
	MustResolve[*TestService](c)
}
