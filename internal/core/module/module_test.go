package module

import (
	"context"
	"testing"

	"github.com/spiringo/spiringo/internal/core/config"
)

// 中文：mockModule 定义当前包使用的数据结构或接口。
// English: mockModule defines a data structure or interface used by this package.
type mockModule struct {
	// 中文：*BaseModule 嵌入复用该类型提供的能力。
	// English: *BaseModule embeds reusable behavior from that type.
	*BaseModule
}

// 中文：newMockModule 执行当前包中的对应流程。
// English: newMockModule executes the corresponding workflow in this package.
func newMockModule(name string, deps ...string) *mockModule {
	return &mockModule{
		BaseModule: NewBaseModule(name, deps...),
	}
}

// 中文：Config 执行当前包中的对应流程。
// English: Config executes the corresponding workflow in this package.
func (m *mockModule) Config() any { return nil }

// 中文：TestBaseModule_Name 验证相关行为符合预期。
// English: TestBaseModule_Name verifies the related behavior.
func TestBaseModule_Name(t *testing.T) {
	m := NewBaseModule("test")
	if m.Name() != "test" {
		t.Errorf("expected 'test', got '%s'", m.Name())
	}
}

// 中文：TestBaseModule_Dependencies 验证相关行为符合预期。
// English: TestBaseModule_Dependencies verifies the related behavior.
func TestBaseModule_Dependencies(t *testing.T) {
	m := NewBaseModule("child", "parent1", "parent2")
	deps := m.Dependencies()
	if len(deps) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(deps))
	}
	if deps[0] != "parent1" || deps[1] != "parent2" {
		t.Errorf("unexpected dependencies: %v", deps)
	}
}

// 中文：TestBaseModule_State 验证相关行为符合预期。
// English: TestBaseModule_State verifies the related behavior.
func TestBaseModule_State(t *testing.T) {
	m := NewBaseModule("test")
	if m.State() != ModuleStateInactive {
		t.Errorf("expected inactive state, got %s", m.State())
	}
	m.SetState(ModuleStateActive)
	if m.State() != ModuleStateActive {
		t.Errorf("expected active state, got %s", m.State())
	}
}

// 中文：TestRegistry_RegisterAndResolve 验证相关行为符合预期。
// English: TestRegistry_RegisterAndResolve verifies the related behavior.
func TestRegistry_RegisterAndResolve(t *testing.T) {
	registry := NewRegistry()

	m1 := newMockModule("base")
	m2 := newMockModule("child", "base")

	if err := registry.Register(m1); err != nil {
		t.Fatal(err)
	}
	if err := registry.Register(m2); err != nil {
		t.Fatal(err)
	}

	if err := registry.ResolveOrder(); err != nil {
		t.Fatal(err)
	}

	// Verify init order: base should come before child
	order := registry.initOrder
	baseIdx := -1
	childIdx := -1
	for i, name := range order {
		if name == "base" {
			baseIdx = i
		}
		if name == "child" {
			childIdx = i
		}
	}

	if baseIdx == -1 || childIdx == -1 {
		t.Fatal("modules not found in resolved order")
	}
	if baseIdx >= childIdx {
		t.Errorf("expected 'base' before 'child', got base at %d, child at %d", baseIdx, childIdx)
	}
}

// 中文：TestRegistry_CircularDependency 验证相关行为符合预期。
// English: TestRegistry_CircularDependency verifies the related behavior.
func TestRegistry_CircularDependency(t *testing.T) {
	registry := NewRegistry()

	m1 := newMockModule("a", "b")
	m2 := newMockModule("b", "a")

	registry.Register(m1)
	registry.Register(m2)

	err := registry.ResolveOrder()
	if err == nil {
		t.Error("expected error for circular dependency")
	}
}

// 中文：TestRegistry_MissingDependency 验证相关行为符合预期。
// English: TestRegistry_MissingDependency verifies the related behavior.
func TestRegistry_MissingDependency(t *testing.T) {
	registry := NewRegistry()

	m := newMockModule("child", "nonexistent")
	registry.Register(m)

	err := registry.ResolveOrder()
	if err == nil {
		t.Error("expected error for missing dependency")
	}
}

// 中文：TestRegistry_DuplicateRegistration 验证相关行为符合预期。
// English: TestRegistry_DuplicateRegistration verifies the related behavior.
func TestRegistry_DuplicateRegistration(t *testing.T) {
	registry := NewRegistry()

	m1 := newMockModule("test")
	registry.Register(m1)

	m2 := newMockModule("test")
	err := registry.Register(m2)
	if err == nil {
		t.Error("expected error for duplicate module registration")
	}
}

// 中文：TestRegistry_Get 验证相关行为符合预期。
// English: TestRegistry_Get verifies the related behavior.
func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()
	m := newMockModule("test")
	registry.Register(m)

	got, err := registry.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name() != "test" {
		t.Errorf("expected 'test', got '%s'", got.Name())
	}
}

// 中文：TestRegistry_GetNotFound 验证相关行为符合预期。
// English: TestRegistry_GetNotFound verifies the related behavior.
func TestRegistry_GetNotFound(t *testing.T) {
	registry := NewRegistry()
	_, err := registry.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent module")
	}
}

// 中文：TestRegistry_ListModules 验证相关行为符合预期。
// English: TestRegistry_ListModules verifies the related behavior.
func TestRegistry_ListModules(t *testing.T) {
	registry := NewRegistry()
	registry.Register(newMockModule("a"))
	registry.Register(newMockModule("b"))

	names := registry.ListModules()
	if len(names) != 2 {
		t.Errorf("expected 2 modules, got %d", len(names))
	}
}

// 中文：TestRegistry_Snapshots 验证相关行为符合预期。
// English: TestRegistry_Snapshots verifies the related behavior.
func TestRegistry_Snapshots(t *testing.T) {
	registry := NewRegistry()
	registry.Register(newMockModule("base"))
	registry.Register(newMockModule("child", "base"))

	cfg := config.NewManager()
	cfg.Set("modules.child.enabled", false)
	if err := registry.InitAll(&App{Config: cfg}); err != nil {
		t.Fatal(err)
	}

	snapshots := registry.Snapshots()
	if len(snapshots) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(snapshots))
	}
	if snapshots[0].Name != "base" || snapshots[0].State != ModuleStateActive.String() || !snapshots[0].Enabled || snapshots[0].Skipped {
		t.Fatalf("unexpected base snapshot: %+v", snapshots[0])
	}
	if snapshots[1].Name != "child" || snapshots[1].State != ModuleStateInactive.String() || snapshots[1].Enabled || !snapshots[1].Skipped {
		t.Fatalf("unexpected child snapshot: %+v", snapshots[1])
	}
	if len(snapshots[1].Dependencies) != 1 || snapshots[1].Dependencies[0] != "base" {
		t.Fatalf("unexpected child dependencies: %+v", snapshots[1].Dependencies)
	}
}

// 中文：TestRegistry_SkipsDependentsOfDisabledModules 验证相关行为符合预期。
// English: TestRegistry_SkipsDependentsOfDisabledModules verifies the related behavior.
func TestRegistry_SkipsDependentsOfDisabledModules(t *testing.T) {
	registry := NewRegistry()
	registry.Register(newMockModule("base"))
	registry.Register(newMockModule("child", "base"))
	registry.Register(newMockModule("grandchild", "child"))

	cfg := config.NewManager()
	cfg.Set("modules.child.enabled", false)
	if err := registry.InitAll(&App{Config: cfg}); err != nil {
		t.Fatal(err)
	}

	snapshots := registry.Snapshots()
	if len(snapshots) != 3 {
		t.Fatalf("expected 3 snapshots, got %d", len(snapshots))
	}
	if snapshots[0].Name != "base" || snapshots[0].State != ModuleStateActive.String() || snapshots[0].Skipped {
		t.Fatalf("unexpected base snapshot: %+v", snapshots[0])
	}
	if snapshots[1].Name != "child" || snapshots[1].State != ModuleStateInactive.String() || !snapshots[1].Skipped {
		t.Fatalf("unexpected child snapshot: %+v", snapshots[1])
	}
	if snapshots[2].Name != "grandchild" || snapshots[2].State != ModuleStateInactive.String() || !snapshots[2].Skipped {
		t.Fatalf("unexpected grandchild snapshot: %+v", snapshots[2])
	}
}

// 中文：TestRegistry_StopAllDoesNotStopSkippedModules 验证相关行为符合预期。
// English: TestRegistry_StopAllDoesNotStopSkippedModules verifies the related behavior.
func TestRegistry_StopAllDoesNotStopSkippedModules(t *testing.T) {
	registry := NewRegistry()
	registry.Register(newMockModule("base"))
	registry.Register(newMockModule("child", "base"))

	cfg := config.NewManager()
	cfg.Set("modules.child.enabled", false)
	if err := registry.InitAll(&App{Config: cfg}); err != nil {
		t.Fatal(err)
	}
	if err := registry.StopAll(context.Background()); err != nil {
		t.Fatal(err)
	}

	snapshots := registry.Snapshots()
	if snapshots[0].Name != "base" || snapshots[0].State != ModuleStateStopped.String() {
		t.Fatalf("unexpected base snapshot after stop: %+v", snapshots[0])
	}
	if snapshots[1].Name != "child" || snapshots[1].State != ModuleStateInactive.String() || !snapshots[1].Skipped {
		t.Fatalf("disabled child should remain inactive after stop: %+v", snapshots[1])
	}
}

// 中文：TestRegistry_RegisterModules 验证相关行为符合预期。
// English: TestRegistry_RegisterModules verifies the related behavior.
func TestRegistry_RegisterModules(t *testing.T) {
	registry := NewRegistry()
	registry.RegisterModules(newMockModule("a"), newMockModule("b"))

	names := registry.ListModules()
	if len(names) != 2 {
		t.Errorf("expected 2 modules, got %d", len(names))
	}
}
