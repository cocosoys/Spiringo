package dto

// 中文：CreateRoleReq 定义当前包使用的数据结构或接口。
// English: CreateRoleReq defines a data structure or interface used by this package.
// CreateRoleReq 创建角色请求
type CreateRoleReq struct {
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `json:"name" binding:"required,min=2,max=64"`
	// 中文：Code 保存当前结构中的配置或数据值。
	// English: Code stores a configuration or data value for this struct.
	Code string `json:"code" binding:"required,min=2,max=64"`
	// 中文：Description 保存当前结构中的配置或数据值。
	// English: Description stores a configuration or data value for this struct.
	Description string `json:"description" binding:"omitempty,max=256"`
}

// 中文：UpdateRoleReq 定义当前包使用的数据结构或接口。
// English: UpdateRoleReq defines a data structure or interface used by this package.
// UpdateRoleReq 更新角色请求
type UpdateRoleReq struct {
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `json:"name" binding:"omitempty,min=2,max=64"`
	// 中文：Description 保存当前结构中的配置或数据值。
	// English: Description stores a configuration or data value for this struct.
	Description string `json:"description" binding:"omitempty,max=256"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `json:"status" binding:"omitempty,oneof=active disabled"`
}

// 中文：AssignPermissionsReq 定义当前包使用的数据结构或接口。
// English: AssignPermissionsReq defines a data structure or interface used by this package.
// AssignPermissionsReq 分配权限请求
type AssignPermissionsReq struct {
	// 中文：PermissionIDs 保存当前结构中的配置或数据值。
	// English: PermissionIDs stores a configuration or data value for this struct.
	PermissionIDs []string `json:"permission_ids" binding:"required,min=1"`
}

// 中文：RoleResp 定义当前包使用的数据结构或接口。
// English: RoleResp defines a data structure or interface used by this package.
// RoleResp 角色响应
type RoleResp struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string `json:"id"`
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `json:"name"`
	// 中文：Code 保存当前结构中的配置或数据值。
	// English: Code stores a configuration or data value for this struct.
	Code string `json:"code"`
	// 中文：Description 保存当前结构中的配置或数据值。
	// English: Description stores a configuration or data value for this struct.
	Description string `json:"description,omitempty"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `json:"status"`
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string `json:"tenant_id,omitempty"`
}

// 中文：PermissionResp 定义当前包使用的数据结构或接口。
// English: PermissionResp defines a data structure or interface used by this package.
// PermissionResp 权限响应
type PermissionResp struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string `json:"id"`
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `json:"name"`
	// 中文：Code 保存当前结构中的配置或数据值。
	// English: Code stores a configuration or data value for this struct.
	Code string `json:"code"`
	// 中文：Resource 保存当前结构中的配置或数据值。
	// English: Resource stores a configuration or data value for this struct.
	Resource string `json:"resource"`
	// 中文：Action 保存当前结构中的配置或数据值。
	// English: Action stores a configuration or data value for this struct.
	Action string `json:"action"`
	// 中文：ParentID 保存当前结构中的配置或数据值。
	// English: ParentID stores a configuration or data value for this struct.
	ParentID string `json:"parent_id,omitempty"`
	// 中文：SortOrder 保存当前结构中的配置或数据值。
	// English: SortOrder stores a configuration or data value for this struct.
	SortOrder int `json:"sort_order"`
}
