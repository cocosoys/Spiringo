package model

import "github.com/spiringo/spiringo/internal/pkg/orm"

// 中文：Role 定义当前包使用的数据结构或接口。
// English: Role defines a data structure or interface used by this package.
// Role 角色
type Role struct {
	// 中文：orm.TenantBaseModel 嵌入复用该类型提供的能力。
	// English: orm.TenantBaseModel embeds reusable behavior from that type.
	orm.TenantBaseModel
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `gorm:"size:64;not null" json:"name"`
	// 中文：Code 保存当前结构中的配置或数据值。
	// English: Code stores a configuration or data value for this struct.
	Code string `gorm:"size:64;not null" json:"code"`
	// 中文：Description 保存当前结构中的配置或数据值。
	// English: Description stores a configuration or data value for this struct.
	Description string `gorm:"size:256" json:"description,omitempty"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `gorm:"size:32;not null;default:active" json:"status"`
}

// 中文：TableName 执行当前包中的对应流程。
// English: TableName executes the corresponding workflow in this package.
func (Role) TableName() string { return "roles" }

// 中文：Permission 定义当前包使用的数据结构或接口。
// English: Permission defines a data structure or interface used by this package.
// Permission 权限
type Permission struct {
	// 中文：orm.TenantBaseModel 嵌入复用该类型提供的能力。
	// English: orm.TenantBaseModel embeds reusable behavior from that type.
	orm.TenantBaseModel
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `gorm:"size:64;not null" json:"name"`
	// 中文：Code 保存当前结构中的配置或数据值。
	// English: Code stores a configuration or data value for this struct.
	Code string `gorm:"size:128;not null" json:"code"`
	// 中文：Resource 保存当前结构中的配置或数据值。
	// English: Resource stores a configuration or data value for this struct.
	Resource string `gorm:"size:64;not null" json:"resource"`
	// 中文：Action 保存当前结构中的配置或数据值。
	// English: Action stores a configuration or data value for this struct.
	Action string `gorm:"size:32;not null" json:"action"`
	// 中文：ParentID 保存当前结构中的配置或数据值。
	// English: ParentID stores a configuration or data value for this struct.
	ParentID string `gorm:"size:36;index" json:"parent_id,omitempty"`
	// 中文：SortOrder 保存当前结构中的配置或数据值。
	// English: SortOrder stores a configuration or data value for this struct.
	SortOrder int `gorm:"default:0" json:"sort_order"`
}

// 中文：TableName 执行当前包中的对应流程。
// English: TableName executes the corresponding workflow in this package.
func (Permission) TableName() string { return "permissions" }

// 中文：RolePermission 定义当前包使用的数据结构或接口。
// English: RolePermission defines a data structure or interface used by this package.
// RolePermission 角色-权限关联
type RolePermission struct {
	// 中文：orm.BaseModel 嵌入复用该类型提供的能力。
	// English: orm.BaseModel embeds reusable behavior from that type.
	orm.BaseModel
	// 中文：RoleID 保存当前结构中的配置或数据值。
	// English: RoleID stores a configuration or data value for this struct.
	RoleID string `gorm:"size:36;index;not null" json:"role_id"`
	// 中文：PermissionID 保存当前结构中的配置或数据值。
	// English: PermissionID stores a configuration or data value for this struct.
	PermissionID string `gorm:"size:36;index;not null" json:"permission_id"`
}

// 中文：TableName 执行当前包中的对应流程。
// English: TableName executes the corresponding workflow in this package.
func (RolePermission) TableName() string { return "role_permissions" }

// 中文：UserRole 定义当前包使用的数据结构或接口。
// English: UserRole defines a data structure or interface used by this package.
// UserRole 用户-角色关联
type UserRole struct {
	// 中文：orm.BaseModel 嵌入复用该类型提供的能力。
	// English: orm.BaseModel embeds reusable behavior from that type.
	orm.BaseModel
	// 中文：UserID 保存当前结构中的配置或数据值。
	// English: UserID stores a configuration or data value for this struct.
	UserID string `gorm:"size:36;index;not null" json:"user_id"`
	// 中文：RoleID 保存当前结构中的配置或数据值。
	// English: RoleID stores a configuration or data value for this struct.
	RoleID string `gorm:"size:36;index;not null" json:"role_id"`
}

// 中文：TableName 执行当前包中的对应流程。
// English: TableName executes the corresponding workflow in this package.
func (UserRole) TableName() string { return "user_roles" }
