package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/modules/auth/dto"
	"github.com/spiringo/spiringo/internal/modules/auth/service"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：AuthHandler 定义当前包使用的数据结构或接口。
// English: AuthHandler defines a data structure or interface used by this package.
// AuthHandler 认证HTTP处理器
type AuthHandler struct {
	// 中文：svc 保存当前结构中的配置或数据值。
	// English: svc stores a configuration or data value for this struct.
	svc *service.AuthService
}

// 中文：NewAuthHandler 创建并返回对应组件实例。
// English: NewAuthHandler creates and returns the corresponding component instance.
// NewAuthHandler 创建认证处理器
func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// 中文：Login 执行当前包中的对应流程。
// English: Login executes the corresponding workflow in this package.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	token, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		types.Fail(c, err)
		return
	}

	types.OK(c, token)
}

// 中文：Register 执行当前包中的对应流程。
// English: Register executes the corresponding workflow in this package.
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	token, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		types.Fail(c, err)
		return
	}

	types.OK(c, token)
}

// 中文：Logout 执行当前包中的对应流程。
// English: Logout executes the corresponding workflow in this package.
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从Header中提取token
	tokenString := extractToken(c.GetHeader("Authorization"))
	if tokenString == "" {
		types.Fail(c, types.ErrTokenInvalid)
		return
	}

	if err := h.svc.Logout(c.Request.Context(), tokenString); err != nil {
		types.Fail(c, err)
		return
	}

	types.OK(c, nil)
}

// 中文：Refresh 执行当前包中的对应流程。
// English: Refresh executes the corresponding workflow in this package.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		types.Fail(c, types.ErrBadRequest.WithMessage(err.Error()))
		return
	}

	token, err := h.svc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		types.Fail(c, err)
		return
	}

	types.OK(c, token)
}

// 中文：OAuthLogin 执行当前包中的对应流程。
// English: OAuthLogin executes the corresponding workflow in this package.
func (h *AuthHandler) OAuthLogin(c *gin.Context) {
	provider := c.Param("provider")
	redirectURL := c.Query("redirect_url")

	url, err := h.svc.GetOAuthAuthorizeURL(c.Request.Context(), provider, redirectURL)
	if err != nil {
		types.Fail(c, err)
		return
	}

	c.Redirect(302, url)
}

// 中文：OAuthCallback 执行当前包中的对应流程。
// English: OAuthCallback executes the corresponding workflow in this package.
func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")
	redirectURL := c.Query("redirect_url")

	token, err := h.svc.HandleOAuthCallback(c.Request.Context(), provider, code, state, redirectURL)
	if err != nil {
		types.Fail(c, err)
		return
	}

	types.OK(c, token)
}

// 中文：extractToken 执行当前包中的对应流程。
// English: extractToken executes the corresponding workflow in this package.
// extractToken 从Authorization头提取token
func extractToken(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return authHeader
}
