package orm

import (
	"context"
	"errors"
	"testing"

	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：TestTenantRouterDefaultsToSharedDB 验证相关行为符合预期。
// English: TestTenantRouterDefaultsToSharedDB verifies the related behavior.
func TestTenantRouterDefaultsToSharedDB(t *testing.T) {
	shared := &DB{}
	router := NewTenantRouter(shared)

	tdb, route, err := router.Route(context.Background())
	if err != nil {
		t.Fatalf("Route returned error: %v", err)
	}
	if tdb.DB() != shared {
		t.Fatal("route without tenant should use shared db")
	}
	if route.Strategy != TenantStrategySharedDB {
		t.Fatalf("strategy = %s", route.Strategy)
	}
}

// 中文：TestTenantRouterRoutesDatabaseStrategy 验证相关行为符合预期。
// English: TestTenantRouterRoutesDatabaseStrategy verifies the related behavior.
func TestTenantRouterRoutesDatabaseStrategy(t *testing.T) {
	shared := &DB{}
	tenantDB := &DB{}
	router := NewTenantRouter(shared)
	if err := router.RegisterDatabase("tenant-db-1", tenantDB); err != nil {
		t.Fatalf("RegisterDatabase returned error: %v", err)
	}
	if err := router.RegisterRoute(TenantRoute{
		TenantID: "tenant-1",
		Strategy: TenantStrategyDatabase,
		Database: "tenant-db-1",
	}); err != nil {
		t.Fatalf("RegisterRoute returned error: %v", err)
	}

	ctx := types.WithTenantID(context.Background(), "tenant-1")
	tdb, route, err := router.Route(ctx)
	if err != nil {
		t.Fatalf("Route returned error: %v", err)
	}
	if tdb.DB() != tenantDB {
		t.Fatal("database strategy should use registered tenant database")
	}
	if route.Database != "tenant-db-1" {
		t.Fatalf("route = %+v", route)
	}
}

// 中文：TestTenantRouterRoutesSchemaStrategy 验证相关行为符合预期。
// English: TestTenantRouterRoutesSchemaStrategy verifies the related behavior.
func TestTenantRouterRoutesSchemaStrategy(t *testing.T) {
	shared := &DB{}
	schemaDB := &DB{}
	router := NewTenantRouter(shared)
	if err := router.RegisterSchema("tenant_1", schemaDB); err != nil {
		t.Fatalf("RegisterSchema returned error: %v", err)
	}
	if err := router.RegisterRoute(TenantRoute{
		TenantID: "tenant-1",
		Strategy: TenantStrategySchema,
		Schema:   "tenant_1",
	}); err != nil {
		t.Fatalf("RegisterRoute returned error: %v", err)
	}

	ctx := types.WithTenantID(context.Background(), "tenant-1")
	tdb, route, err := router.Route(ctx)
	if err != nil {
		t.Fatalf("Route returned error: %v", err)
	}
	if tdb.DB() != schemaDB {
		t.Fatal("schema strategy should use registered schema db")
	}
	if route.Schema != "tenant_1" {
		t.Fatalf("route = %+v", route)
	}
}

// 中文：TestTenantRouterReportsMissingDatabase 验证相关行为符合预期。
// English: TestTenantRouterReportsMissingDatabase verifies the related behavior.
func TestTenantRouterReportsMissingDatabase(t *testing.T) {
	router := NewTenantRouter(&DB{})
	if err := router.RegisterRoute(TenantRoute{
		TenantID: "tenant-1",
		Strategy: TenantStrategyDatabase,
		Database: "missing",
	}); err != nil {
		t.Fatalf("RegisterRoute returned error: %v", err)
	}

	ctx := types.WithTenantID(context.Background(), "tenant-1")
	if _, _, err := router.Route(ctx); !errors.Is(err, ErrTenantDBNotFound) {
		t.Fatalf("err = %v, want ErrTenantDBNotFound", err)
	}
}

// 中文：TestMultiTenantDBFacadeRoutesTenantDB 验证相关行为符合预期。
// English: TestMultiTenantDBFacadeRoutesTenantDB verifies the related behavior.
func TestMultiTenantDBFacadeRoutesTenantDB(t *testing.T) {
	shared := &DB{}
	tenantDB := &DB{}
	multi := NewMultiTenantDB(shared)
	if err := multi.RegisterDatabase("tenant-db-1", tenantDB); err != nil {
		t.Fatalf("RegisterDatabase returned error: %v", err)
	}
	if err := multi.RegisterRoute(TenantRoute{
		TenantID: "tenant-1",
		Strategy: TenantStrategyDatabase,
		Database: "tenant-db-1",
	}); err != nil {
		t.Fatalf("RegisterRoute returned error: %v", err)
	}

	tdb, err := multi.GetDB(types.WithTenantID(context.Background(), "tenant-1"))
	if err != nil {
		t.Fatalf("GetDB returned error: %v", err)
	}
	if tdb.DB() != tenantDB {
		t.Fatal("GetDB should route to tenant database")
	}
}
