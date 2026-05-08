package orm

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：readWriteProbe 定义当前包使用的数据结构或接口。
// English: readWriteProbe defines a data structure or interface used by this package.
type readWriteProbe struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID uint `gorm:"primaryKey"`
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string
}

// 中文：tenantScopedProbe 定义当前包使用的数据结构或接口。
// English: tenantScopedProbe defines a data structure or interface used by this package.
type tenantScopedProbe struct {
	// 中文：BaseModel 嵌入复用该类型提供的能力。
	// English: BaseModel embeds reusable behavior from that type.
	BaseModel
	// 中文：TenantModel 嵌入复用该类型提供的能力。
	// English: TenantModel embeds reusable behavior from that type.
	TenantModel
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string
}

// 中文：TestReadReplicasHandleReadQueries 验证相关行为符合预期。
// English: TestReadReplicasHandleReadQueries verifies the related behavior.
func TestReadReplicasHandleReadQueries(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	primaryDSN := filepath.Join(dir, "primary.db")
	replicaDSN := filepath.Join(dir, "replica.db")

	primarySeed, err := New(Config{Driver: "sqlite", DSN: primaryDSN})
	if err != nil {
		t.Fatalf("open primary seed: %v", err)
	}
	if err := primarySeed.AutoMigrate(&readWriteProbe{}); err != nil {
		t.Fatalf("migrate primary: %v", err)
	}
	if err := primarySeed.Create(ctx, &readWriteProbe{Name: "primary"}); err != nil {
		t.Fatalf("seed primary: %v", err)
	}
	_ = primarySeed.Close()

	replicaSeed, err := New(Config{Driver: "sqlite", DSN: replicaDSN})
	if err != nil {
		t.Fatalf("open replica seed: %v", err)
	}
	if err := replicaSeed.AutoMigrate(&readWriteProbe{}); err != nil {
		t.Fatalf("migrate replica: %v", err)
	}
	if err := replicaSeed.Create(ctx, &readWriteProbe{Name: "replica"}); err != nil {
		t.Fatalf("seed replica: %v", err)
	}
	_ = replicaSeed.Close()

	db, err := New(Config{
		Driver: "sqlite",
		DSN:    primaryDSN,
		ReadReplicas: []EndpointConfig{
			{Driver: "sqlite", DSN: replicaDSN},
		},
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	var got readWriteProbe
	if err := db.First(ctx, &got); err != nil {
		t.Fatalf("read from replica: %v", err)
	}
	if got.Name != "replica" {
		t.Fatalf("read name = %q, want replica", got.Name)
	}

	if err := db.Create(ctx, &readWriteProbe{Name: "written"}); err != nil {
		t.Fatalf("write primary: %v", err)
	}

	var primaryRows []readWriteProbe
	if err := db.DB().Order("id ASC").Find(&primaryRows).Error; err != nil {
		t.Fatalf("query primary: %v", err)
	}
	if len(primaryRows) != 2 || primaryRows[1].Name != "written" {
		t.Fatalf("primary rows = %+v, want primary plus written", primaryRows)
	}
}

// 中文：TestMasterForcesReadFromPrimary 验证相关行为符合预期。
// English: TestMasterForcesReadFromPrimary verifies the related behavior.
func TestMasterForcesReadFromPrimary(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	primaryDSN := filepath.Join(dir, "primary.db")
	replicaDSN := filepath.Join(dir, "replica.db")

	for _, seed := range []struct {
		dsn  string
		name string
	}{
		{dsn: primaryDSN, name: "primary"},
		{dsn: replicaDSN, name: "replica"},
	} {
		db, err := New(Config{Driver: "sqlite", DSN: seed.dsn})
		if err != nil {
			t.Fatalf("open seed %s: %v", seed.dsn, err)
		}
		if err := db.AutoMigrate(&readWriteProbe{}); err != nil {
			t.Fatalf("migrate seed %s: %v", seed.dsn, err)
		}
		if err := db.Create(ctx, &readWriteProbe{Name: seed.name}); err != nil {
			t.Fatalf("seed %s: %v", seed.dsn, err)
		}
		_ = db.Close()
	}

	db, err := New(Config{
		Driver:       "sqlite",
		DSN:          primaryDSN,
		ReadReplicas: []EndpointConfig{{Driver: "sqlite", DSN: replicaDSN}},
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	var got readWriteProbe
	if err := db.Master().First(ctx, &got); err != nil {
		t.Fatalf("master read: %v", err)
	}
	if got.Name != "primary" {
		t.Fatalf("master read name = %q, want primary", got.Name)
	}
}

// 中文：TestReadReplicaWhereKeepsWriteConditionsOnPrimary 验证相关行为符合预期。
// English: TestReadReplicaWhereKeepsWriteConditionsOnPrimary verifies the related behavior.
func TestReadReplicaWhereKeepsWriteConditionsOnPrimary(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	primaryDSN := filepath.Join(dir, "primary.db")
	replicaDSN := filepath.Join(dir, "replica.db")

	for _, dsn := range []string{primaryDSN, replicaDSN} {
		db, err := New(Config{Driver: "sqlite", DSN: dsn})
		if err != nil {
			t.Fatalf("open seed %s: %v", dsn, err)
		}
		if err := db.AutoMigrate(&readWriteProbe{}); err != nil {
			t.Fatalf("migrate seed %s: %v", dsn, err)
		}
		_ = db.Close()
	}

	db, err := New(Config{
		Driver:       "sqlite",
		DSN:          primaryDSN,
		ReadReplicas: []EndpointConfig{{Driver: "sqlite", DSN: replicaDSN}},
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.Create(ctx, &readWriteProbe{Name: "delete-me"}); err != nil {
		t.Fatalf("create primary: %v", err)
	}
	if err := db.Where("name = ?", "delete-me").Delete(ctx, &readWriteProbe{}); err != nil {
		t.Fatalf("delete through where: %v", err)
	}

	var count int64
	if err := db.DB().Model(&readWriteProbe{}).Where("name = ?", "delete-me").Count(&count).Error; err != nil {
		t.Fatalf("count primary: %v", err)
	}
	if count != 0 {
		t.Fatalf("primary count = %d, want 0", count)
	}
}

// 中文：TestTenantDBUpdateAppliesTenantFilter 验证相关行为符合预期。
// English: TestTenantDBUpdateAppliesTenantFilter verifies the related behavior.
func TestTenantDBUpdateAppliesTenantFilter(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()

	db, err := New(Config{Driver: "sqlite", DSN: filepath.Join(dir, "tenant.db")})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()
	if err := db.AutoMigrate(&tenantScopedProbe{}); err != nil {
		t.Fatalf("migrate tenant probe: %v", err)
	}

	tdb := NewTenantDB(db)
	ctxA := types.WithTenantID(ctx, "tenant-a")
	ctxB := types.WithTenantID(ctx, "tenant-b")
	a := &tenantScopedProbe{Name: "alpha"}
	b := &tenantScopedProbe{Name: "beta"}
	if err := tdb.Create(ctxA, a); err != nil {
		t.Fatalf("create tenant a row: %v", err)
	}
	if err := tdb.Create(ctxB, b); err != nil {
		t.Fatalf("create tenant b row: %v", err)
	}

	var own tenantScopedProbe
	if err := tdb.First(ctxA, &own, "id = ?", a.ID); err != nil {
		t.Fatalf("read tenant a row: %v", err)
	}
	own.Name = "alpha-updated"
	if err := tdb.Update(ctxA, &own); err != nil {
		t.Fatalf("update own tenant row: %v", err)
	}

	var gotA tenantScopedProbe
	if err := db.First(ctx, &gotA, "id = ?", a.ID); err != nil {
		t.Fatalf("read updated tenant a row: %v", err)
	}
	if gotA.Name != "alpha-updated" || gotA.TenantID != "tenant-a" {
		t.Fatalf("tenant a row = %+v, want updated within tenant-a", gotA)
	}

	var other tenantScopedProbe
	if err := db.First(ctx, &other, "id = ?", b.ID); err != nil {
		t.Fatalf("read tenant b row: %v", err)
	}
	other.Name = "hijacked"
	if err := tdb.Update(ctxA, &other); err == nil {
		t.Fatalf("expected cross-tenant update to fail")
	}

	var gotB tenantScopedProbe
	if err := db.First(ctx, &gotB, "id = ?", b.ID); err != nil {
		t.Fatalf("read tenant b row after denied update: %v", err)
	}
	if gotB.Name != "beta" || gotB.TenantID != "tenant-b" {
		t.Fatalf("tenant b row = %+v, want unchanged", gotB)
	}

	var total int64
	if err := db.DB().Model(&tenantScopedProbe{}).Count(&total).Error; err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if total != 2 {
		t.Fatalf("row count = %d, want 2", total)
	}
}

// 中文：TestPaginateNormalizesDefaultsAndCapsPageSize 验证相关行为符合预期。
// English: TestPaginateNormalizesDefaultsAndCapsPageSize verifies the related behavior.
func TestPaginateNormalizesDefaultsAndCapsPageSize(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()

	db, err := New(Config{Driver: "sqlite", DSN: filepath.Join(dir, "paginate.db")})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()
	if err := db.AutoMigrate(&readWriteProbe{}); err != nil {
		t.Fatalf("migrate probe: %v", err)
	}
	for i := 0; i < 5; i++ {
		if err := db.Create(ctx, &readWriteProbe{Name: "item"}); err != nil {
			t.Fatalf("seed row %d: %v", i, err)
		}
	}

	var items []readWriteProbe
	total, err := db.Paginate(ctx, &items, 0, 0)
	if err != nil {
		t.Fatalf("Paginate returned error: %v", err)
	}
	if total != 5 || len(items) != 5 {
		t.Fatalf("default pagination total=%d len=%d, want 5/5", total, len(items))
	}

	items = nil
	total, err = db.Paginate(ctx, &items, 1, 1000)
	if err != nil {
		t.Fatalf("Paginate with oversized page size returned error: %v", err)
	}
	if total != 5 || len(items) != 5 {
		t.Fatalf("oversized pagination total=%d len=%d, want 5/5", total, len(items))
	}
}
