//go:build !(windows && 386)

package migrate

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/spiringo/spiringo/internal/core/module"
	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：TestStoreRecordsAppliedMigrations 验证相关行为符合预期。
// English: TestStoreRecordsAppliedMigrations verifies the related behavior.
func TestStoreRecordsAppliedMigrations(t *testing.T) {
	db := newTestDB(t, "file:migrate_store_test?mode=memory&cache=shared")
	defer db.Close()

	store := NewStore(db)
	var up atomic.Int32
	migrations := []module.Migration{
		{
			ID: "demo_001",
			Up: func(ctx context.Context) error {
				up.Add(1)
				return nil
			},
		},
	}

	if err := store.RunMigrations(context.Background(), "demo", migrations); err != nil {
		t.Fatalf("first migration run: %v", err)
	}
	if err := store.RunMigrations(context.Background(), "demo", migrations); err != nil {
		t.Fatalf("second migration run: %v", err)
	}
	if got := up.Load(); got != 1 {
		t.Fatalf("migration up count = %d", got)
	}

	applied, err := store.Applied(context.Background(), "demo_001")
	if err != nil {
		t.Fatal(err)
	}
	if !applied {
		t.Fatal("expected migration to be applied")
	}
}

// 中文：TestStoreRollbackRunsDownAndDeletesRecord 验证相关行为符合预期。
// English: TestStoreRollbackRunsDownAndDeletesRecord verifies the related behavior.
func TestStoreRollbackRunsDownAndDeletesRecord(t *testing.T) {
	db := newTestDB(t, "file:migrate_rollback_test?mode=memory&cache=shared")
	defer db.Close()

	store := NewStore(db)
	var down atomic.Int32
	migrations := []module.Migration{
		{
			ID: "demo_001",
			Up: func(ctx context.Context) error {
				return nil
			},
			Down: func(ctx context.Context) error {
				down.Add(1)
				return nil
			},
		},
	}

	if err := store.RunMigrations(context.Background(), "demo", migrations); err != nil {
		t.Fatal(err)
	}
	if err := store.Rollback(context.Background(), "demo", migrations, 1); err != nil {
		t.Fatal(err)
	}
	if got := down.Load(); got != 1 {
		t.Fatalf("migration down count = %d", got)
	}
	applied, err := store.Applied(context.Background(), "demo_001")
	if err != nil {
		t.Fatal(err)
	}
	if applied {
		t.Fatal("expected migration record to be removed")
	}
}

// 中文：newTestDB 执行当前包中的对应流程。
// English: newTestDB executes the corresponding workflow in this package.
func newTestDB(t *testing.T, dsn string) *orm.DB {
	t.Helper()
	db, err := orm.New(orm.Config{
		Driver:  "sqlite",
		DSN:     dsn,
		MaxIdle: 1,
		MaxOpen: 1,
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	return db
}
