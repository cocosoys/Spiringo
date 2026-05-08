package constants

// 中文：EnvDevelopment、EnvTest、EnvStaging、... 声明当前包使用的常量。
// English: EnvDevelopment、EnvTest、EnvStaging、... declares constants used by this package.
const (
	EnvDevelopment = "development"
	EnvTest        = "test"
	EnvStaging     = "staging"
	EnvProduction  = "production"
	EnvProd        = "prod"
)

// 中文：HeaderRequestID、HeaderTenantID、HeaderTenantName、... 声明当前包使用的常量。
// English: HeaderRequestID、HeaderTenantID、HeaderTenantName、... declares constants used by this package.
const (
	HeaderRequestID      = "X-Request-ID"
	HeaderTenantID       = "X-Tenant-ID"
	HeaderTenantName     = "X-Tenant-Name"
	HeaderTenantStrategy = "X-Tenant-Strategy"
	HeaderIdempotentKey  = "X-Idempotent-Key"
	HeaderLanguage       = "X-Language"
)

// 中文：ContentTypeJSON、ContentTypeForm、ContentTypeText 声明当前包使用的常量。
// English: ContentTypeJSON、ContentTypeForm、ContentTypeText declares constants used by this package.
const (
	ContentTypeJSON = "application/json"
	ContentTypeForm = "application/x-www-form-urlencoded"
	ContentTypeText = "text/plain; charset=utf-8"
)

// 中文：ModuleAuth、ModuleUser、ModuleTenant、... 声明当前包使用的常量。
// English: ModuleAuth、ModuleUser、ModuleTenant、... declares constants used by this package.
const (
	ModuleAuth         = "auth"
	ModuleUser         = "user"
	ModuleTenant       = "tenant"
	ModuleRBAC         = "rbac"
	ModulePayment      = "payment"
	ModuleQRCode       = "qrcode"
	ModuleNotification = "notification"
)
