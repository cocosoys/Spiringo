package orm

import "context"

// 中文：MultiTenantDB 定义当前包使用的数据结构或接口。
// English: MultiTenantDB defines a data structure or interface used by this package.
// MultiTenantDB is the blueprint-facing facade for shared, schema, and
// database-per-tenant routing.
type MultiTenantDB struct {
	// 中文：router 保存当前结构中的配置或数据值。
	// English: router stores a configuration or data value for this struct.
	router *TenantRouter
}

// 中文：NewMultiTenantDB 创建并返回对应组件实例。
// English: NewMultiTenantDB creates and returns the corresponding component instance.
func NewMultiTenantDB(defaultDB *DB) *MultiTenantDB {
	return &MultiTenantDB{router: NewTenantRouter(defaultDB)}
}

// 中文：RegisterRoute 执行当前包中的对应流程。
// English: RegisterRoute executes the corresponding workflow in this package.
func (m *MultiTenantDB) RegisterRoute(route TenantRoute) error {
	return m.router.RegisterRoute(route)
}

// 中文：RegisterDatabase 执行当前包中的对应流程。
// English: RegisterDatabase executes the corresponding workflow in this package.
func (m *MultiTenantDB) RegisterDatabase(name string, db *DB) error {
	return m.router.RegisterDatabase(name, db)
}

// 中文：RegisterSchema 执行当前包中的对应流程。
// English: RegisterSchema executes the corresponding workflow in this package.
func (m *MultiTenantDB) RegisterSchema(name string, db *DB) error {
	return m.router.RegisterSchema(name, db)
}

// 中文：GetDB 执行当前包中的对应流程。
// English: GetDB executes the corresponding workflow in this package.
func (m *MultiTenantDB) GetDB(ctx context.Context) (*TenantDB, error) {
	tdb, _, err := m.Route(ctx)
	return tdb, err
}

// 中文：Route 执行当前包中的对应流程。
// English: Route executes the corresponding workflow in this package.
func (m *MultiTenantDB) Route(ctx context.Context) (*TenantDB, TenantRoute, error) {
	if m == nil || m.router == nil {
		return nil, TenantRoute{}, ErrTenantRouteNotFound
	}
	return m.router.Route(ctx)
}
