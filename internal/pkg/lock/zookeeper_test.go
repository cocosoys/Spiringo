package lock

import "testing"

// 中文：TestNewZooKeeperLockRequiresServers 验证相关行为符合预期。
// English: TestNewZooKeeperLockRequiresServers verifies the related behavior.
func TestNewZooKeeperLockRequiresServers(t *testing.T) {
	if _, err := NewZooKeeperLock(ZooKeeperConfig{}); err == nil {
		t.Fatal("expected missing servers to fail")
	}
}

// 中文：TestZooKeeperLockPathSanitizesKeys 验证相关行为符合预期。
// English: TestZooKeeperLockPathSanitizesKeys verifies the related behavior.
func TestZooKeeperLockPathSanitizesKeys(t *testing.T) {
	locker := &ZooKeeperLock{root: "/spiringo/locks"}
	got := locker.lockPath("tenant/a:b c")
	want := "/spiringo/locks/tenant_a_b_c"
	if got != want {
		t.Fatalf("lockPath = %q, want %q", got, want)
	}
}
