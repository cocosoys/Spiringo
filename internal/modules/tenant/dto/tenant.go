package dto

import "time"

// 中文：CreateTenantReq 定义当前包使用的数据结构或接口。
// English: CreateTenantReq defines a data structure or interface used by this package.
// CreateTenantReq 创建租户请求
type CreateTenantReq struct {
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `json:"name" binding:"required,min=2,max=128"`
	// 中文：Code 保存当前结构中的配置或数据值。
	// English: Code stores a configuration or data value for this struct.
	Code string `json:"code" binding:"required,min=2,max=64"`
	// 中文：Strategy 保存当前结构中的配置或数据值。
	// English: Strategy stores a configuration or data value for this struct.
	Strategy string `json:"strategy" binding:"omitempty,oneof=shared_db schema database"`
	// 中文：Domain 保存当前结构中的配置或数据值。
	// English: Domain stores a configuration or data value for this struct.
	Domain string `json:"domain" binding:"omitempty,max=256"`
}

// 中文：UpdateTenantReq 定义当前包使用的数据结构或接口。
// English: UpdateTenantReq defines a data structure or interface used by this package.
// UpdateTenantReq 更新租户请求
type UpdateTenantReq struct {
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `json:"name" binding:"omitempty,min=2,max=128"`
	// 中文：Domain 保存当前结构中的配置或数据值。
	// English: Domain stores a configuration or data value for this struct.
	Domain string `json:"domain" binding:"omitempty,max=256"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `json:"status" binding:"omitempty,oneof=active suspended expired"`
}

// 中文：TenantResp 定义当前包使用的数据结构或接口。
// English: TenantResp defines a data structure or interface used by this package.
// TenantResp 租户响应
type TenantResp struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string `json:"id"`
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `json:"name"`
	// 中文：Code 保存当前结构中的配置或数据值。
	// English: Code stores a configuration or data value for this struct.
	Code string `json:"code"`
	// 中文：Strategy 保存当前结构中的配置或数据值。
	// English: Strategy stores a configuration or data value for this struct.
	Strategy string `json:"strategy"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `json:"status"`
	// 中文：Domain 保存当前结构中的配置或数据值。
	// English: Domain stores a configuration or data value for this struct.
	Domain string `json:"domain,omitempty"`
	// 中文：LogoURL 保存当前结构中的配置或数据值。
	// English: LogoURL stores a configuration or data value for this struct.
	LogoURL string `json:"logo_url,omitempty"`
	// 中文：ExpireAt 保存当前结构中的配置或数据值。
	// English: ExpireAt stores a configuration or data value for this struct.
	ExpireAt *time.Time `json:"expire_at,omitempty"`
	// 中文：CreatedAt 保存当前结构中的配置或数据值。
	// English: CreatedAt stores a configuration or data value for this struct.
	CreatedAt time.Time `json:"created_at"`
	// 中文：UpdatedAt 保存当前结构中的配置或数据值。
	// English: UpdatedAt stores a configuration or data value for this struct.
	UpdatedAt time.Time `json:"updated_at"`
}
