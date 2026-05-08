package oauth

import (
	"context"
	"fmt"
	"net/url"
)

// 中文：AuthURL 执行当前包中的对应流程。
// English: AuthURL executes the corresponding workflow in this package.
func (p *WechatProvider) AuthURL(_ context.Context, state string, redirectURL string) (string, error) {
	return p.AuthorizeURL(state, redirectURL), nil
}

// 中文：GetUser 执行当前包中的对应流程。
// English: GetUser executes the corresponding workflow in this package.
func (p *WechatProvider) GetUser(ctx context.Context, code string, _ string) (*OAuthUser, error) {
	info, err := p.GetUserInfo(ctx, code)
	return toOAuthUser(info), err
}

// 中文：RefreshToken 执行当前包中的对应流程。
// English: RefreshToken executes the corresponding workflow in this package.
func (p *WechatProvider) RefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return fmt.Errorf("wechat: refresh token is required")
	}
	endpoint := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s",
		url.QueryEscape(p.AppID),
		url.QueryEscape(refreshToken),
	)
	body, err := httpGet(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("wechat refresh token: %w", err)
	}
	return requireOAuthJSONAccessToken("wechat", body, "access_token")
}

// 中文：AuthURL 执行当前包中的对应流程。
// English: AuthURL executes the corresponding workflow in this package.
func (p *QQProvider) AuthURL(_ context.Context, state string, redirectURL string) (string, error) {
	return p.AuthorizeURL(state, redirectURL), nil
}

// 中文：GetUser 执行当前包中的对应流程。
// English: GetUser executes the corresponding workflow in this package.
func (p *QQProvider) GetUser(ctx context.Context, code string, redirectURL string) (*OAuthUser, error) {
	info, err := p.GetUserInfoWithRedirect(ctx, code, redirectURL)
	return toOAuthUser(info), err
}

// 中文：RefreshToken 执行当前包中的对应流程。
// English: RefreshToken executes the corresponding workflow in this package.
func (p *QQProvider) RefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return fmt.Errorf("qq: refresh token is required")
	}
	endpoint := fmt.Sprintf(
		"https://graph.qq.com/oauth2.0/token?grant_type=refresh_token&client_id=%s&client_secret=%s&refresh_token=%s",
		url.QueryEscape(p.AppID),
		url.QueryEscape(p.AppSecret),
		url.QueryEscape(refreshToken),
	)
	body, err := httpGet(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("qq refresh token: %w", err)
	}
	return requireOAuthQueryAccessToken("qq", body)
}

// 中文：AuthURL 执行当前包中的对应流程。
// English: AuthURL executes the corresponding workflow in this package.
func (p *GoogleProvider) AuthURL(_ context.Context, state string, redirectURL string) (string, error) {
	return p.AuthorizeURL(state, redirectURL), nil
}

// 中文：GetUser 执行当前包中的对应流程。
// English: GetUser executes the corresponding workflow in this package.
func (p *GoogleProvider) GetUser(ctx context.Context, code string, redirectURL string) (*OAuthUser, error) {
	info, err := p.GetUserInfoWithRedirect(ctx, code, redirectURL)
	return toOAuthUser(info), err
}

// 中文：RefreshToken 执行当前包中的对应流程。
// English: RefreshToken executes the corresponding workflow in this package.
func (p *GoogleProvider) RefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return fmt.Errorf("google: refresh token is required")
	}
	body, err := httpPostForm(ctx, "https://oauth2.googleapis.com/token", map[string]string{
		"client_id":     p.ClientID,
		"client_secret": p.ClientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	})
	if err != nil {
		return fmt.Errorf("google refresh token: %w", err)
	}
	return requireOAuthJSONAccessToken("google", body, "access_token")
}

// 中文：AuthURL 执行当前包中的对应流程。
// English: AuthURL executes the corresponding workflow in this package.
func (p *DiscordProvider) AuthURL(_ context.Context, state string, redirectURL string) (string, error) {
	return p.AuthorizeURL(state, redirectURL), nil
}

// 中文：GetUser 执行当前包中的对应流程。
// English: GetUser executes the corresponding workflow in this package.
func (p *DiscordProvider) GetUser(ctx context.Context, code string, redirectURL string) (*OAuthUser, error) {
	info, err := p.GetUserInfoWithRedirect(ctx, code, redirectURL)
	return toOAuthUser(info), err
}

// 中文：RefreshToken 执行当前包中的对应流程。
// English: RefreshToken executes the corresponding workflow in this package.
func (p *DiscordProvider) RefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return fmt.Errorf("discord: refresh token is required")
	}
	body, err := httpPostForm(ctx, "https://discord.com/api/oauth2/token", map[string]string{
		"client_id":     p.ClientID,
		"client_secret": p.ClientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	})
	if err != nil {
		return fmt.Errorf("discord refresh token: %w", err)
	}
	return requireOAuthJSONAccessToken("discord", body, "access_token")
}

// 中文：AuthURL 执行当前包中的对应流程。
// English: AuthURL executes the corresponding workflow in this package.
func (p *TelegramProvider) AuthURL(_ context.Context, state string, redirectURL string) (string, error) {
	return p.AuthorizeURL(state, redirectURL), nil
}

// 中文：GetUser 执行当前包中的对应流程。
// English: GetUser executes the corresponding workflow in this package.
func (p *TelegramProvider) GetUser(ctx context.Context, code string, _ string) (*OAuthUser, error) {
	info, err := p.GetUserInfo(ctx, code)
	return toOAuthUser(info), err
}

// 中文：RefreshToken 执行当前包中的对应流程。
// English: RefreshToken executes the corresponding workflow in this package.
func (p *TelegramProvider) RefreshToken(context.Context, string) error {
	return ErrRefreshTokenUnsupported
}

// 中文：AuthURL 执行当前包中的对应流程。
// English: AuthURL executes the corresponding workflow in this package.
func (p *XProvider) AuthURL(_ context.Context, state string, redirectURL string) (string, error) {
	return p.AuthorizeURL(state, redirectURL), nil
}

// 中文：GetUser 执行当前包中的对应流程。
// English: GetUser executes the corresponding workflow in this package.
func (p *XProvider) GetUser(ctx context.Context, code string, redirectURL string) (*OAuthUser, error) {
	info, err := p.GetUserInfoWithRedirect(ctx, code, redirectURL)
	return toOAuthUser(info), err
}

// 中文：RefreshToken 执行当前包中的对应流程。
// English: RefreshToken executes the corresponding workflow in this package.
func (p *XProvider) RefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return fmt.Errorf("x: refresh token is required")
	}
	body, err := httpPostForm(ctx, "https://api.twitter.com/2/oauth2/token", map[string]string{
		"client_id":     p.ClientID,
		"client_secret": p.ClientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	})
	if err != nil {
		return fmt.Errorf("x refresh token: %w", err)
	}
	return requireOAuthJSONAccessToken("x", body, "access_token")
}

// 中文：AuthURL 执行当前包中的对应流程。
// English: AuthURL executes the corresponding workflow in this package.
func (p *DingTalkProvider) AuthURL(_ context.Context, state string, redirectURL string) (string, error) {
	return p.AuthorizeURL(state, redirectURL), nil
}

// 中文：GetUser 执行当前包中的对应流程。
// English: GetUser executes the corresponding workflow in this package.
func (p *DingTalkProvider) GetUser(ctx context.Context, code string, _ string) (*OAuthUser, error) {
	info, err := p.GetUserInfo(ctx, code)
	return toOAuthUser(info), err
}

// 中文：RefreshToken 执行当前包中的对应流程。
// English: RefreshToken executes the corresponding workflow in this package.
func (p *DingTalkProvider) RefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return fmt.Errorf("dingtalk: refresh token is required")
	}
	body, err := httpPostJSON(ctx, "https://api.dingtalk.com/v1.0/oauth2/userAccessToken", map[string]string{
		"clientId":     p.AppID,
		"clientSecret": p.AppSecret,
		"refreshToken": refreshToken,
		"grantType":    "refresh_token",
	}, "")
	if err != nil {
		return fmt.Errorf("dingtalk refresh token: %w", err)
	}
	return requireOAuthJSONAccessToken("dingtalk", body, "accessToken")
}

// 中文：AuthURL 执行当前包中的对应流程。
// English: AuthURL executes the corresponding workflow in this package.
func (p *DouyinProvider) AuthURL(_ context.Context, state string, redirectURL string) (string, error) {
	return p.AuthorizeURL(state, redirectURL), nil
}

// 中文：GetUser 执行当前包中的对应流程。
// English: GetUser executes the corresponding workflow in this package.
func (p *DouyinProvider) GetUser(ctx context.Context, code string, _ string) (*OAuthUser, error) {
	info, err := p.GetUserInfo(ctx, code)
	return toOAuthUser(info), err
}

// 中文：RefreshToken 执行当前包中的对应流程。
// English: RefreshToken executes the corresponding workflow in this package.
func (p *DouyinProvider) RefreshToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return fmt.Errorf("douyin: refresh token is required")
	}
	body, err := httpPostForm(ctx, "https://open.douyin.com/oauth/refresh_token/", map[string]string{
		"client_key":    p.ClientKey,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	})
	if err != nil {
		return fmt.Errorf("douyin refresh token: %w", err)
	}
	return requireOAuthJSONAccessToken("douyin", body, "data", "access_token")
}

// 中文：AuthURL 执行当前包中的对应流程。
// English: AuthURL executes the corresponding workflow in this package.
func (p *WorkWechatProvider) AuthURL(_ context.Context, state string, redirectURL string) (string, error) {
	return p.AuthorizeURL(state, redirectURL), nil
}

// 中文：GetUser 执行当前包中的对应流程。
// English: GetUser executes the corresponding workflow in this package.
func (p *WorkWechatProvider) GetUser(ctx context.Context, code string, _ string) (*OAuthUser, error) {
	info, err := p.GetUserInfo(ctx, code)
	return toOAuthUser(info), err
}

// 中文：RefreshToken 执行当前包中的对应流程。
// English: RefreshToken executes the corresponding workflow in this package.
func (p *WorkWechatProvider) RefreshToken(context.Context, string) error {
	return ErrRefreshTokenUnsupported
}
