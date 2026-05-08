package model

import "github.com/spiringo/spiringo/internal/pkg/orm"

// 中文：QRCodeRecord 定义当前包使用的数据结构或接口。
// English: QRCodeRecord defines a data structure or interface used by this package.
// QRCodeRecord 二维码记录
type QRCodeRecord struct {
	// 中文：orm.TenantBaseModel 嵌入复用该类型提供的能力。
	// English: orm.TenantBaseModel embeds reusable behavior from that type.
	orm.TenantBaseModel
	// 中文：Content 保存当前结构中的配置或数据值。
	// English: Content stores a configuration or data value for this struct.
	Content string `gorm:"type:text;not null" json:"content"`
	// 中文：ShortCode 保存当前结构中的配置或数据值。
	// English: ShortCode stores a configuration or data value for this struct.
	ShortCode string `gorm:"size:16;uniqueIndex" json:"short_code,omitempty"`
	// 中文：ImageURL 保存当前结构中的配置或数据值。
	// English: ImageURL stores a configuration or data value for this struct.
	ImageURL string `gorm:"size:512" json:"image_url,omitempty"`
	// 中文：Size 保存当前结构中的配置或数据值。
	// English: Size stores a configuration or data value for this struct.
	Size int `gorm:"default:256" json:"size"`
	// 中文：Level 保存当前结构中的配置或数据值。
	// English: Level stores a configuration or data value for this struct.
	Level string `gorm:"size:4;default:M" json:"level"`
	// 中文：ForegroundColor 保存当前结构中的配置或数据值。
	// English: ForegroundColor stores a configuration or data value for this struct.
	ForegroundColor string `gorm:"size:16" json:"foreground_color,omitempty"`
	// 中文：BackgroundColor 保存当前结构中的配置或数据值。
	// English: BackgroundColor stores a configuration or data value for this struct.
	BackgroundColor string `gorm:"size:16" json:"background_color,omitempty"`
	// 中文：LogoURL 保存当前结构中的配置或数据值。
	// English: LogoURL stores a configuration or data value for this struct.
	LogoURL string `gorm:"size:512" json:"logo_url,omitempty"`
	// 中文：LogoSize 保存当前结构中的配置或数据值。
	// English: LogoSize stores a configuration or data value for this struct.
	LogoSize int `json:"logo_size,omitempty"`
	// 中文：ScanCount 保存当前结构中的配置或数据值。
	// English: ScanCount stores a configuration or data value for this struct.
	ScanCount int64 `gorm:"default:0" json:"scan_count"`
	// 中文：ExpiredAt 保存当前结构中的配置或数据值。
	// English: ExpiredAt stores a configuration or data value for this struct.
	ExpiredAt *string `json:"expired_at,omitempty"`
}

// 中文：TableName 执行当前包中的对应流程。
// English: TableName executes the corresponding workflow in this package.
func (QRCodeRecord) TableName() string { return "qrcode_records" }

// 中文：ScanLog 定义当前包使用的数据结构或接口。
// English: ScanLog defines a data structure or interface used by this package.
// ScanLog 扫码记录
type ScanLog struct {
	// 中文：orm.BaseModel 嵌入复用该类型提供的能力。
	// English: orm.BaseModel embeds reusable behavior from that type.
	orm.BaseModel
	// 中文：ShortCode 保存当前结构中的配置或数据值。
	// English: ShortCode stores a configuration or data value for this struct.
	ShortCode string `gorm:"size:16;index;not null" json:"short_code"`
	// 中文：IP 保存当前结构中的配置或数据值。
	// English: IP stores a configuration or data value for this struct.
	IP string `gorm:"size:45" json:"ip"`
	// 中文：UserAgent 保存当前结构中的配置或数据值。
	// English: UserAgent stores a configuration or data value for this struct.
	UserAgent string `gorm:"size:512" json:"user_agent"`
	// 中文：UserID 保存当前结构中的配置或数据值。
	// English: UserID stores a configuration or data value for this struct.
	UserID string `gorm:"size:36;index" json:"user_id,omitempty"`
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string `gorm:"size:36;index" json:"tenant_id,omitempty"`
}

// 中文：TableName 执行当前包中的对应流程。
// English: TableName executes the corresponding workflow in this package.
func (ScanLog) TableName() string { return "qrcode_scan_logs" }
