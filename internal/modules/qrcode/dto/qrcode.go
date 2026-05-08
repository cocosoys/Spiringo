package dto

// 中文：GenerateReq 定义当前包使用的数据结构或接口。
// English: GenerateReq defines a data structure or interface used by this package.
// GenerateReq 生成二维码请求
type GenerateReq struct {
	// 中文：Content 保存当前结构中的配置或数据值。
	// English: Content stores a configuration or data value for this struct.
	Content string `json:"content" binding:"required,min=1"`
	// 中文：Size 保存当前结构中的配置或数据值。
	// English: Size stores a configuration or data value for this struct.
	Size int `json:"size" binding:"omitempty,min=64,max=1024"`
	// 中文：Level 保存当前结构中的配置或数据值。
	// English: Level stores a configuration or data value for this struct.
	Level string `json:"level" binding:"omitempty,oneof=low medium high highest L M Q H l m q h"`
	// 中文：ForegroundColor 保存当前结构中的配置或数据值。
	// English: ForegroundColor stores a configuration or data value for this struct.
	ForegroundColor string `json:"foreground_color" binding:"omitempty,max=16"`
	// 中文：BackgroundColor 保存当前结构中的配置或数据值。
	// English: BackgroundColor stores a configuration or data value for this struct.
	BackgroundColor string `json:"background_color" binding:"omitempty,max=16"`
	// 中文：LogoURL 保存当前结构中的配置或数据值。
	// English: LogoURL stores a configuration or data value for this struct.
	LogoURL string `json:"logo_url" binding:"omitempty,max=512"`
	// 中文：LogoSize 保存当前结构中的配置或数据值。
	// English: LogoSize stores a configuration or data value for this struct.
	LogoSize int `json:"logo_size" binding:"omitempty,min=16,max=512"`
}

// 中文：GenerateResp 定义当前包使用的数据结构或接口。
// English: GenerateResp defines a data structure or interface used by this package.
// GenerateResp 生成二维码响应
type GenerateResp struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string `json:"id"`
	// 中文：Content 保存当前结构中的配置或数据值。
	// English: Content stores a configuration or data value for this struct.
	Content string `json:"content"`
	// 中文：ShortCode 保存当前结构中的配置或数据值。
	// English: ShortCode stores a configuration or data value for this struct.
	ShortCode string `json:"short_code,omitempty"`
	// 中文：ImageURL 保存当前结构中的配置或数据值。
	// English: ImageURL stores a configuration or data value for this struct.
	ImageURL string `json:"image_url,omitempty"`
	// 中文：ShortURL 保存当前结构中的配置或数据值。
	// English: ShortURL stores a configuration or data value for this struct.
	ShortURL string `json:"short_url,omitempty"`
	// 中文：ImageData 保存当前结构中的配置或数据值。
	// English: ImageData stores a configuration or data value for this struct.
	ImageData []byte `json:"-"` // PNG二进制数据，不序列化到JSON
}

// 中文：StatsResp 定义当前包使用的数据结构或接口。
// English: StatsResp defines a data structure or interface used by this package.
// StatsResp 二维码统计响应
type StatsResp struct {
	// 中文：ShortCode 保存当前结构中的配置或数据值。
	// English: ShortCode stores a configuration or data value for this struct.
	ShortCode string `json:"short_code"`
	// 中文：ScanCount 保存当前结构中的配置或数据值。
	// English: ScanCount stores a configuration or data value for this struct.
	ScanCount int64 `json:"scan_count"`
	// 中文：Content 保存当前结构中的配置或数据值。
	// English: Content stores a configuration or data value for this struct.
	Content string `json:"content"`
}
