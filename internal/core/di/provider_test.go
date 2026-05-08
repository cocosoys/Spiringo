package di

import "testing"

// 中文：providerTestService 定义当前包使用的数据结构或接口。
// English: providerTestService defines a data structure or interface used by this package.
type providerTestService interface {
	// 中文：Name 声明该接口需要实现的行为。
	// English: Name declares behavior required by this interface.
	Name() string
}

// 中文：providerTestImpl 定义当前包使用的数据结构或接口。
// English: providerTestImpl defines a data structure or interface used by this package.
type providerTestImpl struct {
	// 中文：name 保存当前结构中的配置或数据值。
	// English: name stores a configuration or data value for this struct.
	name string
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (s providerTestImpl) Name() string { return s.name }

// 中文：TestProviderRegisterAndInject 验证相关行为符合预期。
// English: TestProviderRegisterAndInject verifies the related behavior.
func TestProviderRegisterAndInject(t *testing.T) {
	c := NewContainer()
	err := Register(c,
		Singleton[providerTestService](providerTestImpl{name: "primary"}),
		Named("secondary", providerTestImpl{name: "secondary"}),
		Factory(func() *providerTestImpl {
			return &providerTestImpl{name: "factory"}
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	var svc providerTestService
	if err := Inject(c, &svc); err != nil {
		t.Fatal(err)
	}
	if svc.Name() != "primary" {
		t.Fatalf("service name = %q", svc.Name())
	}

	var named providerTestImpl
	if err := InjectNamed(c, "secondary", &named); err != nil {
		t.Fatal(err)
	}
	if named.Name() != "secondary" {
		t.Fatalf("named service = %q", named.Name())
	}

	var factory *providerTestImpl
	if err := Inject(c, &factory); err != nil {
		t.Fatal(err)
	}
	if factory.Name() != "factory" {
		t.Fatalf("factory service = %q", factory.Name())
	}
}

// 中文：TestInjectRejectsNilTarget 验证相关行为符合预期。
// English: TestInjectRejectsNilTarget verifies the related behavior.
func TestInjectRejectsNilTarget(t *testing.T) {
	if err := Inject[string](NewContainer(), nil); err == nil {
		t.Fatal("expected nil target to fail")
	}
}

// 中文：TestProviderRejectsInvalidInputs 验证相关行为符合预期。
// English: TestProviderRejectsInvalidInputs verifies the related behavior.
func TestProviderRejectsInvalidInputs(t *testing.T) {
	if err := Register(nil); err == nil {
		t.Fatal("expected nil container to fail")
	}

	if err := Register(NewContainer(), ProviderFunc(nil)); err == nil {
		t.Fatal("expected nil provider func to fail")
	}

	if err := Register(NewContainer(), Concrete(nil)); err == nil {
		t.Fatal("expected nil concrete instance to fail")
	}

	if err := Register(NewContainer(), Factory[*providerTestImpl](nil)); err == nil {
		t.Fatal("expected nil factory to fail")
	}
}

// 中文：TestInjectRejectsNilContainer 验证相关行为符合预期。
// English: TestInjectRejectsNilContainer verifies the related behavior.
func TestInjectRejectsNilContainer(t *testing.T) {
	var target string
	if err := Inject(nil, &target); err == nil {
		t.Fatal("expected nil container to fail")
	}

	if err := InjectNamed(nil, "name", &target); err == nil {
		t.Fatal("expected nil container for named injection to fail")
	}
}
