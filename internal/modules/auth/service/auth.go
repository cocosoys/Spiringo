package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/spiringo/spiringo/internal/core/event"
	"github.com/spiringo/spiringo/internal/modules/auth/dto"
	"github.com/spiringo/spiringo/internal/modules/auth/model"
	"github.com/spiringo/spiringo/internal/modules/auth/oauth"
	"github.com/spiringo/spiringo/internal/modules/auth/repository"
	"github.com/spiringo/spiringo/internal/pkg/cache"
	"github.com/spiringo/spiringo/internal/pkg/crypto"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：UserServiceInterface 定义当前包使用的数据结构或接口。
// English: UserServiceInterface defines a data structure or interface used by this package.
// UserServiceInterface 用户服务接口（避免循环依赖，auth模块不直接import user模块）
type UserServiceInterface interface {
	// 中文：VerifyPasswordForAuth 声明该接口需要实现的行为。
	// English: VerifyPasswordForAuth declares behavior required by this interface.
	VerifyPasswordForAuth(ctx context.Context, username, password string) (userID, tenantID string, err error)
	// 中文：CreateUserForAuth 声明该接口需要实现的行为。
	// English: CreateUserForAuth declares behavior required by this interface.
	CreateUserForAuth(ctx context.Context, username, email, phone, password, nickname string) (userID, tenantID string, err error)
}

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
// Config 认证服务配置
type Config struct {
	// 中文：JWT 保存当前结构中的配置或数据值。
	// English: JWT stores a configuration or data value for this struct.
	JWT struct {
		Secret        string
		AccessExpire  time.Duration
		RefreshExpire time.Duration
		Issuer        string
	}
}

// 中文：AuthService 定义当前包使用的数据结构或接口。
// English: AuthService defines a data structure or interface used by this package.
// AuthService 认证业务逻辑
type AuthService struct {
	// 中文：config 保存当前结构中的配置或数据值。
	// English: config stores a configuration or data value for this struct.
	config Config
	// 中文：eventBus 保存当前结构中的配置或数据值。
	// English: eventBus stores a configuration or data value for this struct.
	eventBus *event.Bus
	// 中文：authRepo 保存当前结构中的配置或数据值。
	// English: authRepo stores a configuration or data value for this struct.
	authRepo *repository.AuthRepository
	// 中文：userSvc 保存当前结构中的配置或数据值。
	// English: userSvc stores a configuration or data value for this struct.
	userSvc UserServiceInterface
	// 中文：cache 保存当前结构中的配置或数据值。
	// English: cache stores a configuration or data value for this struct.
	cache cache.Cache
	// 中文：oauthReg 保存当前结构中的配置或数据值。
	// English: oauthReg stores a configuration or data value for this struct.
	oauthReg *oauth.Registry
}

// 中文：NewAuthService 创建并返回对应组件实例。
// English: NewAuthService creates and returns the corresponding component instance.
// NewAuthService 创建认证服务
func NewAuthService(config Config, eventBus *event.Bus, authRepo *repository.AuthRepository, userSvc UserServiceInterface, c cache.Cache, oauthReg *oauth.Registry) *AuthService {
	return &AuthService{
		config:   config,
		eventBus: eventBus,
		authRepo: authRepo,
		userSvc:  userSvc,
		cache:    c,
		oauthReg: oauthReg,
	}
}

// 中文：Login 执行当前包中的对应流程。
// English: Login executes the corresponding workflow in this package.
// Login 登录
func (s *AuthService) Login(ctx context.Context, req dto.LoginReq) (*dto.TokenResp, error) {
	userID, tenantID, err := s.userSvc.VerifyPasswordForAuth(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := crypto.GenerateToken(crypto.JWTConfig{
		Secret:        s.config.JWT.Secret,
		AccessExpire:  s.config.JWT.AccessExpire,
		RefreshExpire: s.config.JWT.RefreshExpire,
		Issuer:        s.config.JWT.Issuer,
	}, userID, req.Username, tenantID)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	s.publishAsync(ctx, event.NewEvent(event.EventAuthLogin, &event.AuthEventPayload{
		UserID:   userID,
		TenantID: tenantID,
	}))

	return &dto.TokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWT.AccessExpire.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// 中文：Register 执行当前包中的对应流程。
// English: Register executes the corresponding workflow in this package.
// Register 注册
func (s *AuthService) Register(ctx context.Context, req dto.RegisterReq) (*dto.TokenResp, error) {
	userID, tenantID, err := s.userSvc.CreateUserForAuth(ctx, req.Username, req.Email, req.Phone, req.Password, req.Nickname)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := crypto.GenerateToken(crypto.JWTConfig{
		Secret:        s.config.JWT.Secret,
		AccessExpire:  s.config.JWT.AccessExpire,
		RefreshExpire: s.config.JWT.RefreshExpire,
		Issuer:        s.config.JWT.Issuer,
	}, userID, req.Username, tenantID)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &dto.TokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWT.AccessExpire.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// 中文：Logout 执行当前包中的对应流程。
// English: Logout executes the corresponding workflow in this package.
// Logout 登出（将token加入黑名单）
func (s *AuthService) Logout(ctx context.Context, tokenString string) error {
	if s.cache == nil {
		return nil
	}

	claims, err := crypto.ParseToken(s.config.JWT.Secret, tokenString)
	if err != nil {
		return nil // 已失效的token无需处理
	}

	// 将token加入黑名单，TTL等于token剩余有效时间
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl > 0 {
		blacklistKey := fmt.Sprintf("token_blacklist:%s", tokenString)
		if err := s.cache.Set(ctx, blacklistKey, "1", ttl); err != nil {
			return fmt.Errorf("failed to blacklist token: %w", err)
		}
	}

	s.publishAsync(ctx, event.NewEvent(event.EventAuthLogout, &event.AuthEventPayload{
		UserID: claims.UserID,
	}))

	return nil
}

// 中文：RefreshToken 执行当前包中的对应流程。
// English: RefreshToken executes the corresponding workflow in this package.
// RefreshToken 刷新Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResp, error) {
	claims, err := crypto.ParseToken(s.config.JWT.Secret, refreshToken)
	if err != nil {
		return nil, types.ErrTokenRefreshFailed.WithMessage("invalid refresh token")
	}

	// 检查是否在黑名单中
	if s.cache != nil {
		blacklistKey := fmt.Sprintf("token_blacklist:%s", refreshToken)
		if exists, _ := s.cache.Exists(ctx, blacklistKey); exists {
			return nil, types.ErrTokenRefreshFailed.WithMessage("token has been revoked")
		}
	}

	accessToken, newRefreshToken, err := crypto.GenerateToken(crypto.JWTConfig{
		Secret:        s.config.JWT.Secret,
		AccessExpire:  s.config.JWT.AccessExpire,
		RefreshExpire: s.config.JWT.RefreshExpire,
		Issuer:        s.config.JWT.Issuer,
	}, claims.UserID, claims.Username, claims.TenantID)
	if err != nil {
		return nil, types.ErrTokenRefreshFailed.WithMessage("generate token failed")
	}

	return &dto.TokenResp{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.config.JWT.AccessExpire.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// 中文：ValidateToken 执行当前包中的对应流程。
// English: ValidateToken executes the corresponding workflow in this package.
// ValidateToken 验证AccessToken
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*crypto.JWTClaims, error) {
	claims, err := crypto.ParseToken(s.config.JWT.Secret, tokenString)
	if err != nil {
		return nil, types.ErrTokenInvalid
	}

	// 检查黑名单
	if s.cache != nil {
		blacklistKey := fmt.Sprintf("token_blacklist:%s", tokenString)
		if exists, _ := s.cache.Exists(ctx, blacklistKey); exists {
			return nil, types.ErrTokenExpired
		}
	}

	return claims, nil
}

// 中文：GetOAuthAuthorizeURL 执行当前包中的对应流程。
// English: GetOAuthAuthorizeURL executes the corresponding workflow in this package.
// GetOAuthAuthorizeURL 获取OAuth授权URL
func (s *AuthService) GetOAuthAuthorizeURL(ctx context.Context, provider, redirectURL string) (string, error) {
	if s.oauthReg == nil {
		return "", types.ErrOAuthProvider.WithMessage("oauth not configured")
	}

	p, ok := s.oauthReg.Get(provider)
	if !ok {
		return "", types.ErrOAuthProvider.WithMessagef("oauth provider %s not registered", provider)
	}

	state := generateOAuthState()
	// 将state存入缓存用于回调校验
	if s.cache != nil {
		stateKey := fmt.Sprintf("oauth_state:%s", state)
		_ = s.cache.Set(ctx, stateKey, provider, 10*time.Minute)
	}

	return p.AuthorizeURL(state, redirectURL), nil
}

// 中文：HandleOAuthCallback 执行当前包中的对应流程。
// English: HandleOAuthCallback executes the corresponding workflow in this package.
// HandleOAuthCallback 处理OAuth回调
func (s *AuthService) HandleOAuthCallback(ctx context.Context, provider, code, state, redirectURL string) (*dto.TokenResp, error) {
	if s.oauthReg == nil {
		return nil, types.ErrOAuthFailed.WithMessage("oauth not configured")
	}

	// 验证state
	if s.cache != nil && state != "" {
		stateKey := fmt.Sprintf("oauth_state:%s", state)
		savedProvider, err := func() (string, error) {
			var val any
			if err := s.cache.Get(ctx, stateKey, &val); err != nil {
				return "", err
			}
			return fmt.Sprintf("%v", val), nil
		}()
		if err != nil || savedProvider != provider {
			return nil, types.ErrOAuthFailed.WithMessage("invalid oauth state")
		}
		_ = s.cache.Delete(ctx, stateKey)
	}

	// 获取Provider
	p, ok := s.oauthReg.Get(provider)
	if !ok {
		return nil, types.ErrOAuthProvider.WithMessagef("oauth provider %s not registered", provider)
	}

	// 用code换取用户信息
	userInfo, err := getOAuthUserInfo(ctx, p, code, redirectURL)
	if err != nil {
		return nil, types.ErrOAuthFailed.WithMessage(err.Error())
	}

	// 查找已有OAuth绑定
	binding, err := s.authRepo.GetOAuthBinding(ctx, provider, userInfo.ProviderUID)
	if err == nil && binding != nil {
		// 已绑定用户，直接签发token
		// 需要通过UserService获取用户信息
		accessToken, refreshToken, err := crypto.GenerateToken(crypto.JWTConfig{
			Secret:        s.config.JWT.Secret,
			AccessExpire:  s.config.JWT.AccessExpire,
			RefreshExpire: s.config.JWT.RefreshExpire,
			Issuer:        s.config.JWT.Issuer,
		}, binding.UserID, provider, binding.TenantID)
		if err != nil {
			return nil, types.ErrOAuthFailed.WithMessage("generate token failed")
		}

		return &dto.TokenResp{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int64(s.config.JWT.AccessExpire.Seconds()),
			TokenType:    "Bearer",
		}, nil
	}

	// 未绑定用户，创建新用户并绑定
	username := fmt.Sprintf("%s_%s", provider, userInfo.ProviderUID)
	nickname := userInfo.Nickname
	if nickname == "" {
		nickname = username
	}

	userID, tenantID, err := s.userSvc.CreateUserForAuth(ctx, username, userInfo.Email, "", generateRandomPassword(16), nickname)
	if err != nil {
		return nil, types.ErrOAuthFailed.WithMessagef("create user from oauth: %v", err)
	}

	// 创建OAuth绑定
	if err := s.BindOAuth(ctx, userID, provider, userInfo.ProviderUID, userInfo.OpenID, userInfo.UnionID, userInfo.Nickname, userInfo.Avatar); err != nil {
		return nil, types.ErrOAuthFailed.WithMessagef("bind oauth: %v", err)
	}

	// 签发token
	accessToken, refreshToken, err := crypto.GenerateToken(crypto.JWTConfig{
		Secret:        s.config.JWT.Secret,
		AccessExpire:  s.config.JWT.AccessExpire,
		RefreshExpire: s.config.JWT.RefreshExpire,
		Issuer:        s.config.JWT.Issuer,
	}, userID, username, tenantID)
	if err != nil {
		return nil, types.ErrOAuthFailed.WithMessage("generate token failed")
	}

	s.publishAsync(ctx, event.NewEvent(event.EventAuthOAuthBound, map[string]string{
		"user_id":  userID,
		"provider": provider,
	}))

	return &dto.TokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWT.AccessExpire.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// 中文：getOAuthUserInfo 执行当前包中的对应流程。
// English: getOAuthUserInfo executes the corresponding workflow in this package.
func getOAuthUserInfo(ctx context.Context, p oauth.Provider, code, redirectURL string) (*oauth.UserInfo, error) {
	if rp, ok := p.(oauth.RedirectProvider); ok {
		return rp.GetUserInfoWithRedirect(ctx, code, redirectURL)
	}
	return p.GetUserInfo(ctx, code)
}

// 中文：BindOAuth 执行当前包中的对应流程。
// English: BindOAuth executes the corresponding workflow in this package.
// BindOAuth 绑定OAuth账号
func (s *AuthService) BindOAuth(ctx context.Context, userID, provider, providerUID, openID, unionID, nickname, avatar string) error {
	binding := &model.OAuthBinding{
		UserID:      userID,
		Provider:    provider,
		ProviderUID: providerUID,
		OpenID:      openID,
		UnionID:     unionID,
		Nickname:    nickname,
		Avatar:      avatar,
	}

	if err := s.authRepo.CreateOAuthBinding(ctx, binding); err != nil {
		return fmt.Errorf("bind oauth: %w", err)
	}

	s.publishAsync(ctx, event.NewEvent(event.EventAuthOAuthBound, map[string]string{
		"user_id":  userID,
		"provider": provider,
	}))

	return nil
}

// 中文：generateOAuthState 执行当前包中的对应流程。
// English: generateOAuthState executes the corresponding workflow in this package.
// generateOAuthState 生成OAuth state参数
func generateOAuthState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// 中文：generateRandomPassword 执行当前包中的对应流程。
// English: generateRandomPassword executes the corresponding workflow in this package.
// generateRandomPassword 生成随机密码
func generateRandomPassword(length int) string {
	b := make([]byte, (length+1)/2)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)[:length]
}

// 中文：CleanUserOAuthBindings 执行当前包中的对应流程。
// English: CleanUserOAuthBindings executes the corresponding workflow in this package.
// CleanUserOAuthBindings 清理用户的所有OAuth绑定（用户删除时调用）
func (s *AuthService) CleanUserOAuthBindings(ctx context.Context, userID string) error {
	return s.authRepo.DeleteByUserID(ctx, userID)
}

// 中文：publishAsync 执行当前包中的对应流程。
// English: publishAsync executes the corresponding workflow in this package.
func (s *AuthService) publishAsync(ctx context.Context, e *event.Event) {
	if s.eventBus != nil {
		_ = s.eventBus.PublishAsync(ctx, e)
	}
}
