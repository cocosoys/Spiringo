package model

import "github.com/spiringo/spiringo/internal/pkg/orm"

// 中文：User 定义当前包使用的数据结构或接口。
// English: User defines a data structure or interface used by this package.
// User 用户模型
type User struct {
	// 中文：orm.TenantBaseModel 嵌入复用该类型提供的能力。
	// English: orm.TenantBaseModel embeds reusable behavior from that type.
	orm.TenantBaseModel
	// 中文：Username 保存当前结构中的配置或数据值。
	// English: Username stores a configuration or data value for this struct.
	Username string `gorm:"size:64;uniqueIndex:idx_tenant_username;not null" json:"username"`
	// 中文：Email 保存当前结构中的配置或数据值。
	// English: Email stores a configuration or data value for this struct.
	Email string `gorm:"size:128;index" json:"email,omitempty"`
	// 中文：Phone 保存当前结构中的配置或数据值。
	// English: Phone stores a configuration or data value for this struct.
	Phone string `gorm:"size:20;index" json:"phone,omitempty"`
	// 中文：Password 保存当前结构中的配置或数据值。
	// English: Password stores a configuration or data value for this struct.
	Password string `gorm:"size:256;not null" json:"-"`
	// 中文：Nickname 保存当前结构中的配置或数据值。
	// English: Nickname stores a configuration or data value for this struct.
	Nickname string `gorm:"size:64" json:"nickname,omitempty"`
	// 中文：Avatar 保存当前结构中的配置或数据值。
	// English: Avatar stores a configuration or data value for this struct.
	Avatar string `gorm:"size:512" json:"avatar,omitempty"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `gorm:"size:32;not null;default:active" json:"status"`
}
