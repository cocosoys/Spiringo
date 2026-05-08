package dto

import "github.com/spiringo/spiringo/pkg/types"

// 中文：SendReq 定义当前包使用的数据结构或接口。
// English: SendReq defines a data structure or interface used by this package.
type SendReq struct {
	// 中文：Event 保存当前结构中的配置或数据值。
	// English: Event stores a configuration or data value for this struct.
	Event string `json:"event" binding:"required"`
	// 中文：Severity 保存当前结构中的配置或数据值。
	// English: Severity stores a configuration or data value for this struct.
	Severity string `json:"severity" binding:"omitempty,oneof=info warning error critical"`
	// 中文：Subject 保存当前结构中的配置或数据值。
	// English: Subject stores a configuration or data value for this struct.
	Subject string `json:"subject" binding:"required,max=256"`
	// 中文：Content 保存当前结构中的配置或数据值。
	// English: Content stores a configuration or data value for this struct.
	Content string `json:"content" binding:"omitempty,max=4096"`
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string `json:"tenant_id" binding:"omitempty,max=64"`
	// 中文：RecipientID 保存当前结构中的配置或数据值。
	// English: RecipientID stores a configuration or data value for this struct.
	RecipientID string `json:"recipient_id" binding:"omitempty,max=64"`
	// 中文：Payload 保存当前结构中的配置或数据值。
	// English: Payload stores a configuration or data value for this struct.
	Payload map[string]any `json:"payload"`
}

// 中文：InboxListReq 定义当前包使用的数据结构或接口。
// English: InboxListReq defines a data structure or interface used by this package.
type InboxListReq struct {
	// 中文：types.PaginationRequest 嵌入复用该类型提供的能力。
	// English: types.PaginationRequest embeds reusable behavior from that type.
	types.PaginationRequest
	// 中文：Event 保存当前结构中的配置或数据值。
	// English: Event stores a configuration or data value for this struct.
	Event string `form:"event" binding:"omitempty,max=128"`
	// 中文：RecipientID 保存当前结构中的配置或数据值。
	// English: RecipientID stores a configuration or data value for this struct.
	RecipientID string `form:"recipient_id" binding:"omitempty,max=64"`
	// 中文：UnreadOnly 保存当前结构中的配置或数据值。
	// English: UnreadOnly stores a configuration or data value for this struct.
	UnreadOnly bool `form:"unread_only"`
}
