package dto

// 中文：CreateUserReq 定义当前包使用的数据结构或接口。
// English: CreateUserReq defines a data structure or interface used by this package.
// CreateUserReq 创建用户请求
type CreateUserReq struct {
	// 中文：Username 保存当前结构中的配置或数据值。
	// English: Username stores a configuration or data value for this struct.
	Username string `json:"username" binding:"required,min=3,max=64"`
	// 中文：Email 保存当前结构中的配置或数据值。
	// English: Email stores a configuration or data value for this struct.
	Email string `json:"email" binding:"omitempty,email,max=128"`
	// 中文：Phone 保存当前结构中的配置或数据值。
	// English: Phone stores a configuration or data value for this struct.
	Phone string `json:"phone" binding:"omitempty,max=20"`
	// 中文：Password 保存当前结构中的配置或数据值。
	// English: Password stores a configuration or data value for this struct.
	Password string `json:"password" binding:"required,min=6,max=64"`
	// 中文：Nickname 保存当前结构中的配置或数据值。
	// English: Nickname stores a configuration or data value for this struct.
	Nickname string `json:"nickname" binding:"omitempty,max=64"`
}

// 中文：UpdateUserReq 定义当前包使用的数据结构或接口。
// English: UpdateUserReq defines a data structure or interface used by this package.
// UpdateUserReq 更新用户请求
type UpdateUserReq struct {
	// 中文：Email 保存当前结构中的配置或数据值。
	// English: Email stores a configuration or data value for this struct.
	Email string `json:"email" binding:"omitempty,email,max=128"`
	// 中文：Phone 保存当前结构中的配置或数据值。
	// English: Phone stores a configuration or data value for this struct.
	Phone string `json:"phone" binding:"omitempty,max=20"`
	// 中文：Nickname 保存当前结构中的配置或数据值。
	// English: Nickname stores a configuration or data value for this struct.
	Nickname string `json:"nickname" binding:"omitempty,max=64"`
	// 中文：Avatar 保存当前结构中的配置或数据值。
	// English: Avatar stores a configuration or data value for this struct.
	Avatar string `json:"avatar" binding:"omitempty,max=512"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `json:"status" binding:"omitempty,oneof=active disabled"`
}

// 中文：UserResp 定义当前包使用的数据结构或接口。
// English: UserResp defines a data structure or interface used by this package.
// UserResp 用户响应
type UserResp struct {
	// 中文：ID 保存当前结构中的配置或数据值。
	// English: ID stores a configuration or data value for this struct.
	ID string `json:"id"`
	// 中文：Username 保存当前结构中的配置或数据值。
	// English: Username stores a configuration or data value for this struct.
	Username string `json:"username"`
	// 中文：Email 保存当前结构中的配置或数据值。
	// English: Email stores a configuration or data value for this struct.
	Email string `json:"email,omitempty"`
	// 中文：Phone 保存当前结构中的配置或数据值。
	// English: Phone stores a configuration or data value for this struct.
	Phone string `json:"phone,omitempty"`
	// 中文：Nickname 保存当前结构中的配置或数据值。
	// English: Nickname stores a configuration or data value for this struct.
	Nickname string `json:"nickname,omitempty"`
	// 中文：Avatar 保存当前结构中的配置或数据值。
	// English: Avatar stores a configuration or data value for this struct.
	Avatar string `json:"avatar,omitempty"`
	// 中文：Status 保存当前结构中的配置或数据值。
	// English: Status stores a configuration or data value for this struct.
	Status string `json:"status"`
	// 中文：TenantID 保存当前结构中的配置或数据值。
	// English: TenantID stores a configuration or data value for this struct.
	TenantID string `json:"tenant_id,omitempty"`
}
