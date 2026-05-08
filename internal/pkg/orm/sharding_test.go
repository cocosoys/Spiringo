package orm

import (
	"errors"
	"testing"
)

// 中文：TestHashShardStrategyRoutesNumericKeys 验证相关行为符合预期。
// English: TestHashShardStrategyRoutesNumericKeys verifies the related behavior.
func TestHashShardStrategyRoutesNumericKeys(t *testing.T) {
	strategy, err := NewHashShardStrategy([]string{"db0", "db1"}, "orders", 16)
	if err != nil {
		t.Fatalf("new strategy: %v", err)
	}

	target, err := strategy.Target(uint64(17))
	if err != nil {
		t.Fatalf("target: %v", err)
	}
	if target.Database != "db1" {
		t.Fatalf("database = %s", target.Database)
	}
	if target.Table != "orders_01" {
		t.Fatalf("table = %s", target.Table)
	}
}

// 中文：TestHashShardStrategyRejectsEmptyKeys 验证相关行为符合预期。
// English: TestHashShardStrategyRejectsEmptyKeys verifies the related behavior.
func TestHashShardStrategyRejectsEmptyKeys(t *testing.T) {
	strategy, err := NewHashShardStrategy([]string{"db0"}, "orders", 4)
	if err != nil {
		t.Fatalf("new strategy: %v", err)
	}
	if _, err := strategy.Target(""); !errors.Is(err, ErrShardKeyRequired) {
		t.Fatalf("expected ErrShardKeyRequired, got %v", err)
	}
}

// 中文：TestShardedDBReportsMissingShard 验证相关行为符合预期。
// English: TestShardedDBReportsMissingShard verifies the related behavior.
func TestShardedDBReportsMissingShard(t *testing.T) {
	strategy, err := NewHashShardStrategy([]string{"db0"}, "orders", 4)
	if err != nil {
		t.Fatalf("new strategy: %v", err)
	}
	sharded := NewShardedDB(nil, nil, strategy)

	_, target, err := sharded.ForKey(1)
	if !errors.Is(err, ErrShardNotFound) {
		t.Fatalf("expected ErrShardNotFound, target=%+v err=%v", target, err)
	}
}

// 中文：TestShardTable 验证相关行为符合预期。
// English: TestShardTable verifies the related behavior.
func TestShardTable(t *testing.T) {
	if got := ShardTable("orders", 7); got != "orders_7" {
		t.Fatalf("table = %s", got)
	}
	if got := ShardTable("", 7); got != "" {
		t.Fatalf("empty prefix table = %s", got)
	}
}
