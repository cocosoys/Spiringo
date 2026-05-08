package model

import (
	"time"

	"github.com/spiringo/spiringo/internal/pkg/orm"
)

// 中文：Tenant 定义当前包使用的数据结构或接口。
// English: Tenant defines a data structure or interface used by this package.
// Tenant 租户模型
type Tenant struct {
	// 中文：orm.BaseModel 嵌入复用该类型提供的能力。
	// English: orm.BaseModel embeds reusable behavior from that type.
	orm.BaseModel
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `gorm:"size:128;not null" json:"name"`
	// 中文：Code 保存当前结构中的配置或数据值。
	// English: Code stores a configuration or data value for this struct.
	Code string `gorm:"size:64;uniqueIndex;not null" json:"code"`
	// 中文：Strategy 保存当前结构中的配置或数据值。
	// English: Strategy stores a configuration or data value for this struct.
	Strategy string `gorm:"size:32;not null;default:shared_db" json:"strategy"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `gorm:"size:32;not null;default:active" json:"status"`
	// 中文：Domain 保存当前结构中的配置或数据值。
	// English: Domain stores a configuration or data value for this struct.
	Domain string `gorm:"size:256" json:"domain,omitempty"`
	// 中文：LogoURL 保存当前结构中的配置或数据值。
	// English: LogoURL stores a configuration or data value for this struct.
	LogoURL string `gorm:"size:512" json:"logo_url,omitempty"`
	// 中文：ExpireAt 保存当前结构中的配置或数据值。
	// English: ExpireAt stores a configuration or data value for this struct.
	ExpireAt *time.Time `json:"expire_at,omitempty"`
}
