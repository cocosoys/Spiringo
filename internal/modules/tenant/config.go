package tenant

// 中文：Strategy 定义当前包使用的数据结构或接口。
// English: Strategy defines a data structure or interface used by this package.
// 租户隔离策略
type Strategy string

// 中文：StrategySharedDB、StrategySchema、StrategyDatabase 声明当前包使用的常量。
// English: StrategySharedDB、StrategySchema、StrategyDatabase declares constants used by this package.
const (
	StrategySharedDB Strategy = "shared_db" // 共享数据库+tenant_id
	StrategySchema   Strategy = "schema"    // 独立Schema
	StrategyDatabase Strategy = "database"  // 独立数据库
)

// 中文：TenantStatus 定义当前包使用的数据结构或接口。
// English: TenantStatus defines a data structure or interface used by this package.
// TenantStatus 租户状态
type TenantStatus string

// 中文：StatusActive、StatusSuspended、StatusExpired 声明当前包使用的常量。
// English: StatusActive、StatusSuspended、StatusExpired declares constants used by this package.
const (
	StatusActive    TenantStatus = "active"
	StatusSuspended TenantStatus = "suspended"
	StatusExpired   TenantStatus = "expired"
)
