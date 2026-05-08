package orm

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：TenantStrategySharedDB、TenantStrategySchema、TenantStrategyDatabase 声明当前包使用的常量。
// English: TenantStrategySharedDB、TenantStrategySchema、TenantStrategyDatabase declares constants used by this package.
const (
	TenantStrategySharedDB = "shared_db"
	TenantStrategySchema   = "schema"
	TenantStrategyDatabase = "database"
)

// 中文：ErrTenantRouteNotFound、ErrTenantDBNotFound、ErrTenantSchemaNotFound 声明当前包使用的变量。
// English: ErrTenantRouteNotFound、ErrTenantDBNotFound、ErrTenantSchemaNotFound declares variables used by this package.
var (
	ErrTenantRouteNotFound  = errors.New("tenant route not found")
	ErrTenantDBNotFound     = errors.New("tenant database not found")
	ErrTenantSchemaNotFound = errors.New("tenant schema not found")
)

// 中文：TenantRoute 定义当前包使用的数据结构或接口。
// English: TenantRoute defines a data structure or interface used by this package.
type TenantRoute struct {
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string
	// 中文：Strategy 保存当前结构中的配置或数据值。
	// English: Strategy stores a configuration or data value for this struct.
	Strategy string
	// 中文：Schema 保存当前结构中的配置或数据值。
	// English: Schema stores a configuration or data value for this struct.
	Schema string
	// 中文：Database 保存当前结构中的配置或数据值。
	// English: Database stores a configuration or data value for this struct.
	Database string
}

// 中文：TenantRouter 定义当前包使用的数据结构或接口。
// English: TenantRouter defines a data structure or interface used by this package.
type TenantRouter struct {
	// 中文：sharedDB 保存当前结构中的配置或数据值。
	// English: sharedDB stores a configuration or data value for this struct.
	sharedDB *DB
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.RWMutex
	// 中文：routes 保存当前结构中的配置或数据值。
	// English: routes stores a configuration or data value for this struct.
	routes map[string]TenantRoute
	// 中文：databases 保存当前结构中的配置或数据值。
	// English: databases stores a configuration or data value for this struct.
	databases map[string]*DB
	// 中文：schemas 保存当前结构中的配置或数据值。
	// English: schemas stores a configuration or data value for this struct.
	schemas map[string]*DB
}

// 中文：NewTenantRouter 创建并返回对应组件实例。
// English: NewTenantRouter creates and returns the corresponding component instance.
func NewTenantRouter(sharedDB *DB) *TenantRouter {
	return &TenantRouter{
		sharedDB:  sharedDB,
		routes:    make(map[string]TenantRoute),
		databases: make(map[string]*DB),
		schemas:   make(map[string]*DB),
	}
}

// 中文：RegisterRoute 执行当前包中的对应流程。
// English: RegisterRoute executes the corresponding workflow in this package.
func (r *TenantRouter) RegisterRoute(route TenantRoute) error {
	if route.TenantID == "" {
		return fmt.Errorf("tenant id is required")
	}
	if route.Strategy == "" {
		route.Strategy = TenantStrategySharedDB
	}
	switch route.Strategy {
	case TenantStrategySharedDB, TenantStrategySchema, TenantStrategyDatabase:
	default:
		return fmt.Errorf("unsupported tenant strategy: %s", route.Strategy)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes[route.TenantID] = route
	return nil
}

// 中文：RegisterDatabase 执行当前包中的对应流程。
// English: RegisterDatabase executes the corresponding workflow in this package.
func (r *TenantRouter) RegisterDatabase(name string, db *DB) error {
	if name == "" {
		return fmt.Errorf("database name is required")
	}
	if db == nil {
		return fmt.Errorf("database db is required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.databases[name] = db
	return nil
}

// 中文：RegisterSchema 执行当前包中的对应流程。
// English: RegisterSchema executes the corresponding workflow in this package.
func (r *TenantRouter) RegisterSchema(name string, db *DB) error {
	if name == "" {
		return fmt.Errorf("schema name is required")
	}
	if db == nil {
		return fmt.Errorf("schema db is required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.schemas[name] = db
	return nil
}

// 中文：Route 执行当前包中的对应流程。
// English: Route executes the corresponding workflow in this package.
func (r *TenantRouter) Route(ctx context.Context) (*TenantDB, TenantRoute, error) {
	if r == nil {
		return nil, TenantRoute{}, ErrTenantRouteNotFound
	}
	tenantID := types.GetTenantID(ctx)
	if tenantID == "" {
		return NewTenantDB(r.sharedDB), TenantRoute{Strategy: TenantStrategySharedDB}, nil
	}

	r.mu.RLock()
	route, ok := r.routes[tenantID]
	if !ok {
		r.mu.RUnlock()
		return NewTenantDB(r.sharedDB), TenantRoute{TenantID: tenantID, Strategy: TenantStrategySharedDB}, nil
	}
	db, err := r.resolve(route)
	r.mu.RUnlock()
	if err != nil {
		return nil, route, err
	}
	return NewTenantDB(db), route, nil
}

// 中文：resolve 执行当前包中的对应流程。
// English: resolve executes the corresponding workflow in this package.
func (r *TenantRouter) resolve(route TenantRoute) (*DB, error) {
	switch route.Strategy {
	case "", TenantStrategySharedDB:
		return r.sharedDB, nil
	case TenantStrategySchema:
		name := firstNonEmpty(route.Schema, route.TenantID)
		db := r.schemas[name]
		if db == nil {
			return nil, fmt.Errorf("%w: %s", ErrTenantSchemaNotFound, name)
		}
		return db, nil
	case TenantStrategyDatabase:
		name := firstNonEmpty(route.Database, route.TenantID)
		db := r.databases[name]
		if db == nil {
			return nil, fmt.Errorf("%w: %s", ErrTenantDBNotFound, name)
		}
		return db, nil
	default:
		return nil, fmt.Errorf("unsupported tenant strategy: %s", route.Strategy)
	}
}

// 中文：firstNonEmpty 执行当前包中的对应流程。
// English: firstNonEmpty executes the corresponding workflow in this package.
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
