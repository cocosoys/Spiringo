package model

import "github.com/spiringo/spiringo/internal/pkg/orm"

// 中文：OAuthBinding 定义当前包使用的数据结构或接口。
// English: OAuthBinding defines a data structure or interface used by this package.
// OAuthBinding OAuth绑定关系
type OAuthBinding struct {
	// 中文：orm.BaseModel 嵌入复用该类型提供的能力。
	// English: orm.BaseModel embeds reusable behavior from that type.
	orm.BaseModel
	// 中文：orm.TenantModel 嵌入复用该类型提供的能力。
	// English: orm.TenantModel embeds reusable behavior from that type.
	orm.TenantModel
	// 中文：UserID 保存当前结构中的配置或数据值。
	// English: UserID stores a configuration or data value for this struct.
	UserID string `gorm:"size:36;index;not null" json:"user_id"`
	// 中文：Provider 保存当前结构中的配置或数据值。
	// English: Provider stores a configuration or data value for this struct.
	Provider string `gorm:"size:32;not null" json:"provider"`
	// 中文：ProviderUID 保存当前结构中的配置或数据值。
	// English: ProviderUID stores a configuration or data value for this struct.
	ProviderUID string `gorm:"size:128;not null" json:"provider_uid"`
	// 中文：OpenID 保存当前结构中的配置或数据值。
	// English: OpenID stores a configuration or data value for this struct.
	OpenID string `gorm:"size:128" json:"open_id,omitempty"`
	// 中文：UnionID 保存当前结构中的配置或数据值。
	// English: UnionID stores a configuration or data value for this struct.
	UnionID string `gorm:"size:128" json:"union_id,omitempty"`
	// 中文：Nickname 保存当前结构中的配置或数据值。
	// English: Nickname stores a configuration or data value for this struct.
	Nickname string `gorm:"size:64" json:"nickname,omitempty"`
	// 中文：Avatar 保存当前结构中的配置或数据值。
	// English: Avatar stores a configuration or data value for this struct.
	Avatar string `gorm:"size:512" json:"avatar,omitempty"`
	// 中文：Email 保存当前结构中的配置或数据值。
	// English: Email stores a configuration or data value for this struct.
	Email string `gorm:"size:128" json:"email,omitempty"`
}

// 中文：TableName 执行当前包中的对应流程。
// English: TableName executes the corresponding workflow in this package.
// TableName 表名
func (OAuthBinding) TableName() string { return "oauth_bindings" }
