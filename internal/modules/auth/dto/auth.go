package dto

// 中文：LoginReq 定义当前包使用的数据结构或接口。
// English: LoginReq defines a data structure or interface used by this package.
// LoginReq 登录请求
type LoginReq struct {
	// 中文：Username 保存当前结构中的配置或数据值。
	// English: Username stores a configuration or data value for this struct.
	Username string `json:"username" binding:"required,min=2,max=64"`
	// 中文：Password 保存当前结构中的配置或数据值。
	// English: Password stores a configuration or data value for this struct.
	Password string `json:"password" binding:"required,min=6,max=64"`
}

// 中文：RegisterReq 定义当前包使用的数据结构或接口。
// English: RegisterReq defines a data structure or interface used by this package.
// RegisterReq 注册请求
type RegisterReq struct {
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

// 中文：RefreshReq 定义当前包使用的数据结构或接口。
// English: RefreshReq defines a data structure or interface used by this package.
// RefreshReq 刷新Token请求
type RefreshReq struct {
	// 中文：RefreshToken 保存当前结构中的配置或数据值。
	// English: RefreshToken stores a configuration or data value for this struct.
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// 中文：TokenResp 定义当前包使用的数据结构或接口。
// English: TokenResp defines a data structure or interface used by this package.
// TokenResp Token响应
type TokenResp struct {
	// 中文：AccessToken 保存当前结构中的配置或数据值。
	// English: AccessToken stores a configuration or data value for this struct.
	AccessToken string `json:"access_token"`
	// 中文：RefreshToken 保存当前结构中的配置或数据值。
	// English: RefreshToken stores a configuration or data value for this struct.
	RefreshToken string `json:"refresh_token"`
	// 中文：ExpiresIn 保存当前结构中的配置或数据值。
	// English: ExpiresIn stores a configuration or data value for this struct.
	ExpiresIn int64 `json:"expires_in"` // 秒
	// 中文：TokenType 保存当前结构中的配置或数据值。
	// English: TokenType stores a configuration or data value for this struct.
	TokenType string `json:"token_type"`
}

// 中文：OAuthCallbackReq 定义当前包使用的数据结构或接口。
// English: OAuthCallbackReq defines a data structure or interface used by this package.
// OAuthCallbackReq OAuth回调请求
type OAuthCallbackReq struct {
	// 中文：Code 保存当前结构中的配置或数据值。
	// English: Code stores a configuration or data value for this struct.
	Code string `form:"code" binding:"required"`
	// 中文：State 保存当前结构中的配置或数据值。
	// English: State stores a configuration or data value for this struct.
	State string `form:"state"`
}
