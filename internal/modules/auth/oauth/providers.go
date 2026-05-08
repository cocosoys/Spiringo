package oauth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// 中文：WechatProvider 定义当前包使用的数据结构或接口。
// English: WechatProvider defines a data structure or interface used by this package.
// WechatProvider 微信OAuth
type WechatProvider struct {
	// 中文：AppID 保存当前结构中的配置或数据值。
	// English: AppID stores a configuration or data value for this struct.
	AppID string
	// 中文：AppSecret 保存当前结构中的配置或数据值。
	// English: AppSecret stores a configuration or data value for this struct.
	AppSecret string
}

// 中文：NewWechatProvider 创建并返回对应组件实例。
// English: NewWechatProvider creates and returns the corresponding component instance.
func NewWechatProvider(appID, appSecret string) *WechatProvider {
	return &WechatProvider{AppID: appID, AppSecret: appSecret}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (p *WechatProvider) Name() string { return "wechat" }

// 中文：AuthorizeURL 执行当前包中的对应流程。
// English: AuthorizeURL executes the corresponding workflow in this package.
func (p *WechatProvider) AuthorizeURL(state, redirectURL string) string {
	return fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect",
		url.QueryEscape(p.AppID),
		url.QueryEscape(redirectURL),
		url.QueryEscape(state),
	)
}

// 中文：GetUserInfo 执行当前包中的对应流程。
// English: GetUserInfo executes the corresponding workflow in this package.
func (p *WechatProvider) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	// Step 1: 用code换access_token
	tokenURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		p.AppID, p.AppSecret, code,
	)
	tokenBody, err := httpGet(ctx, tokenURL)
	if err != nil {
		return nil, fmt.Errorf("wechat get access_token: %w", err)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		OpenID      string `json:"openid"`
		UnionID     string `json:"unionid"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("wechat parse token response: %w", err)
	}
	if tokenResp.ErrCode != 0 {
		return nil, fmt.Errorf("wechat api error: %d %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	// Step 2: 用access_token获取用户信息
	userURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		tokenResp.AccessToken, tokenResp.OpenID,
	)
	userBody, err := httpGet(ctx, userURL)
	if err != nil {
		return nil, fmt.Errorf("wechat get userinfo: %w", err)
	}

	var userResp struct {
		OpenID     string `json:"openid"`
		UnionID    string `json:"unionid"`
		Nickname   string `json:"nickname"`
		HeadImgURL string `json:"headimgurl"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
	if err := json.Unmarshal(userBody, &userResp); err != nil {
		return nil, fmt.Errorf("wechat parse userinfo response: %w", err)
	}
	if userResp.ErrCode != 0 {
		return nil, fmt.Errorf("wechat userinfo error: %d %s", userResp.ErrCode, userResp.ErrMsg)
	}

	return &UserInfo{
		Provider:    "wechat",
		ProviderUID: userResp.OpenID,
		OpenID:      userResp.OpenID,
		UnionID:     userResp.UnionID,
		Nickname:    userResp.Nickname,
		Avatar:      userResp.HeadImgURL,
	}, nil
}

// 中文：QQProvider 定义当前包使用的数据结构或接口。
// English: QQProvider defines a data structure or interface used by this package.
// QQProvider QQ OAuth
type QQProvider struct {
	// 中文：AppID 保存当前结构中的配置或数据值。
	// English: AppID stores a configuration or data value for this struct.
	AppID string
	// 中文：AppSecret 保存当前结构中的配置或数据值。
	// English: AppSecret stores a configuration or data value for this struct.
	AppSecret string
}

// 中文：NewQQProvider 创建并返回对应组件实例。
// English: NewQQProvider creates and returns the corresponding component instance.
func NewQQProvider(appID, appSecret string) *QQProvider {
	return &QQProvider{AppID: appID, AppSecret: appSecret}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (p *QQProvider) Name() string { return "qq" }

// 中文：AuthorizeURL 执行当前包中的对应流程。
// English: AuthorizeURL executes the corresponding workflow in this package.
func (p *QQProvider) AuthorizeURL(state, redirectURL string) string {
	return fmt.Sprintf(
		"https://graph.qq.com/oauth2.0/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=get_user_info&state=%s",
		url.QueryEscape(p.AppID),
		url.QueryEscape(redirectURL),
		url.QueryEscape(state),
	)
}

// 中文：GetUserInfo 执行当前包中的对应流程。
// English: GetUserInfo executes the corresponding workflow in this package.
func (p *QQProvider) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	return p.GetUserInfoWithRedirect(ctx, code, "")
}

// 中文：GetUserInfoWithRedirect 执行当前包中的对应流程。
// English: GetUserInfoWithRedirect executes the corresponding workflow in this package.
func (p *QQProvider) GetUserInfoWithRedirect(ctx context.Context, code string, redirectURL string) (*UserInfo, error) {
	// Step 1: 用code换access_token
	q := url.Values{}
	q.Set("client_id", p.AppID)
	q.Set("client_secret", p.AppSecret)
	q.Set("code", code)
	q.Set("grant_type", "authorization_code")
	if redirectURL != "" {
		q.Set("redirect_uri", redirectURL)
	}
	tokenURL := "https://graph.qq.com/oauth2.0/token?" + q.Encode()
	tokenBody, err := httpGet(ctx, tokenURL)
	if err != nil {
		return nil, fmt.Errorf("qq get access_token: %w", err)
	}

	// QQ返回的是callback格式
	tokenStr := string(tokenBody)
	accessToken := extractParam(tokenStr, "access_token")
	if accessToken == "" {
		return nil, fmt.Errorf("qq: access_token not found in response")
	}

	// Step 2: 获取openid
	openIDURL := fmt.Sprintf("https://graph.qq.com/oauth2.0/me?access_token=%s", accessToken)
	openIDBody, err := httpGet(ctx, openIDURL)
	if err != nil {
		return nil, fmt.Errorf("qq get openid: %w", err)
	}

	openIDStr := string(openIDBody)
	openID := extractJSONValue(openIDStr, "openid")
	if openID == "" {
		return nil, fmt.Errorf("qq: openid not found in response")
	}

	// Step 3: 获取用户信息
	userURL := fmt.Sprintf(
		"https://graph.qq.com/user/get_user_info?access_token=%s&oauth_consumer_key=%s&openid=%s",
		accessToken, p.AppID, openID,
	)
	userBody, err := httpGet(ctx, userURL)
	if err != nil {
		return nil, fmt.Errorf("qq get userinfo: %w", err)
	}

	var userResp struct {
		Nickname  string `json:"nickname"`
		FigureURL string `json:"figureurl_qq_2"`
		Ret       int    `json:"ret"`
		Msg       string `json:"msg"`
	}
	if err := json.Unmarshal(userBody, &userResp); err != nil {
		return nil, fmt.Errorf("qq parse userinfo: %w", err)
	}
	if userResp.Ret != 0 {
		return nil, fmt.Errorf("qq userinfo error: %d %s", userResp.Ret, userResp.Msg)
	}

	return &UserInfo{
		Provider:    "qq",
		ProviderUID: openID,
		OpenID:      openID,
		Nickname:    userResp.Nickname,
		Avatar:      userResp.FigureURL,
	}, nil
}

// 中文：GoogleProvider 定义当前包使用的数据结构或接口。
// English: GoogleProvider defines a data structure or interface used by this package.
// GoogleProvider Google OAuth2
type GoogleProvider struct {
	// 中文：ClientID 保存当前结构中的配置或数据值。
	// English: ClientID stores a configuration or data value for this struct.
	ClientID string
	// 中文：ClientSecret 保存当前结构中的配置或数据值。
	// English: ClientSecret stores a configuration or data value for this struct.
	ClientSecret string
}

// 中文：NewGoogleProvider 创建并返回对应组件实例。
// English: NewGoogleProvider creates and returns the corresponding component instance.
func NewGoogleProvider(clientID, clientSecret string) *GoogleProvider {
	return &GoogleProvider{ClientID: clientID, ClientSecret: clientSecret}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (p *GoogleProvider) Name() string { return "google" }

// 中文：AuthorizeURL 执行当前包中的对应流程。
// English: AuthorizeURL executes the corresponding workflow in this package.
func (p *GoogleProvider) AuthorizeURL(state, redirectURL string) string {
	return fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid+email+profile&state=%s",
		url.QueryEscape(p.ClientID),
		url.QueryEscape(redirectURL),
		url.QueryEscape(state),
	)
}

// 中文：GetUserInfo 执行当前包中的对应流程。
// English: GetUserInfo executes the corresponding workflow in this package.
func (p *GoogleProvider) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	return p.GetUserInfoWithRedirect(ctx, code, "")
}

// 中文：GetUserInfoWithRedirect 执行当前包中的对应流程。
// English: GetUserInfoWithRedirect executes the corresponding workflow in this package.
func (p *GoogleProvider) GetUserInfoWithRedirect(ctx context.Context, code string, redirectURL string) (*UserInfo, error) {
	// Step 1: 用code换token
	form := map[string]string{
		"client_id":     p.ClientID,
		"client_secret": p.ClientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
	}
	if redirectURL != "" {
		form["redirect_uri"] = redirectURL
	}
	tokenBody, err := httpPostForm(ctx, "https://oauth2.googleapis.com/token", form)
	if err != nil {
		return nil, fmt.Errorf("google get token: %w", err)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}
	if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("google parse token: %w", err)
	}

	// Step 2: 获取用户信息
	userBody, err := httpGetWithAuth(ctx, "https://www.googleapis.com/oauth2/v2/userinfo", tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("google get userinfo: %w", err)
	}

	var userResp struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.Unmarshal(userBody, &userResp); err != nil {
		return nil, fmt.Errorf("google parse userinfo: %w", err)
	}

	return &UserInfo{
		Provider:    "google",
		ProviderUID: userResp.ID,
		Nickname:    userResp.Name,
		Avatar:      userResp.Picture,
		Email:       userResp.Email,
	}, nil
}

// 中文：DiscordProvider 定义当前包使用的数据结构或接口。
// English: DiscordProvider defines a data structure or interface used by this package.
// DiscordProvider Discord OAuth2
type DiscordProvider struct {
	// 中文：ClientID 保存当前结构中的配置或数据值。
	// English: ClientID stores a configuration or data value for this struct.
	ClientID string
	// 中文：ClientSecret 保存当前结构中的配置或数据值。
	// English: ClientSecret stores a configuration or data value for this struct.
	ClientSecret string
}

// 中文：NewDiscordProvider 创建并返回对应组件实例。
// English: NewDiscordProvider creates and returns the corresponding component instance.
func NewDiscordProvider(clientID, clientSecret string) *DiscordProvider {
	return &DiscordProvider{ClientID: clientID, ClientSecret: clientSecret}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (p *DiscordProvider) Name() string { return "discord" }

// 中文：AuthorizeURL 执行当前包中的对应流程。
// English: AuthorizeURL executes the corresponding workflow in this package.
func (p *DiscordProvider) AuthorizeURL(state, redirectURL string) string {
	return fmt.Sprintf(
		"https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify+email&state=%s",
		url.QueryEscape(p.ClientID),
		url.QueryEscape(redirectURL),
		url.QueryEscape(state),
	)
}

// 中文：GetUserInfo 执行当前包中的对应流程。
// English: GetUserInfo executes the corresponding workflow in this package.
func (p *DiscordProvider) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	return p.GetUserInfoWithRedirect(ctx, code, "")
}

// 中文：GetUserInfoWithRedirect 执行当前包中的对应流程。
// English: GetUserInfoWithRedirect executes the corresponding workflow in this package.
func (p *DiscordProvider) GetUserInfoWithRedirect(ctx context.Context, code string, redirectURL string) (*UserInfo, error) {
	form := map[string]string{
		"client_id":     p.ClientID,
		"client_secret": p.ClientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
	}
	if redirectURL != "" {
		form["redirect_uri"] = redirectURL
	}
	tokenBody, err := httpPostForm(ctx, "https://discord.com/api/oauth2/token", form)
	if err != nil {
		return nil, fmt.Errorf("discord get token: %w", err)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("discord parse token: %w", err)
	}

	userBody, err := httpGetWithAuth(ctx, "https://discord.com/api/v10/users/@me", tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("discord get userinfo: %w", err)
	}

	var userResp struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
		Email    string `json:"email"`
	}
	if err := json.Unmarshal(userBody, &userResp); err != nil {
		return nil, fmt.Errorf("discord parse userinfo: %w", err)
	}

	avatarURL := ""
	if userResp.Avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userResp.ID, userResp.Avatar)
	}

	return &UserInfo{
		Provider:    "discord",
		ProviderUID: userResp.ID,
		Nickname:    userResp.Username,
		Avatar:      avatarURL,
		Email:       userResp.Email,
	}, nil
}

// 中文：TelegramProvider 定义当前包使用的数据结构或接口。
// English: TelegramProvider defines a data structure or interface used by this package.
// TelegramProvider Telegram OAuth
type TelegramProvider struct {
	// 中文：BotToken 保存当前结构中的配置或数据值。
	// English: BotToken stores a configuration or data value for this struct.
	BotToken string
}

// 中文：NewTelegramProvider 创建并返回对应组件实例。
// English: NewTelegramProvider creates and returns the corresponding component instance.
func NewTelegramProvider(botToken string) *TelegramProvider {
	return &TelegramProvider{BotToken: botToken}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (p *TelegramProvider) Name() string { return "telegram" }

// 中文：AuthorizeURL 执行当前包中的对应流程。
// English: AuthorizeURL executes the corresponding workflow in this package.
func (p *TelegramProvider) AuthorizeURL(state, redirectURL string) string {
	// Telegram Login Widget 使用 bot username，这里用通用方式
	return fmt.Sprintf(
		"https://oauth.telegram.org/auth?bot_id=%s&origin=%s&request_access=write&return_to=%s",
		url.QueryEscape(extractBotID(p.BotToken)),
		url.QueryEscape(redirectURL),
		url.QueryEscape(state),
	)
}

// 中文：GetUserInfo 执行当前包中的对应流程。
// English: GetUserInfo executes the corresponding workflow in this package.
func (p *TelegramProvider) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	// Telegram Login Widget: code 实际是前端传递的完整回调数据（JSON格式）
	// 包含 id, first_name, last_name, username, photo_url, hash 等字段
	// 需要验证 hash 后解析用户数据
	var data struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		PhotoURL  string `json:"photo_url"`
		Hash      string `json:"hash"`
		AuthDate  int64  `json:"auth_date"`
	}
	if err := json.Unmarshal([]byte(code), &data); err != nil {
		return nil, fmt.Errorf("telegram: parse callback data: %w", err)
	}

	// HMAC-SHA256 验证
	// 1. 用 SHA256(bot_token) 作为 secret
	secret := sha256.Sum256([]byte(p.BotToken))
	// 2. 构造 check_string: 按 key 字典序排列 "key=value\n"
	checkPairs := []string{
		fmt.Sprintf("auth_date=%d", data.AuthDate),
		fmt.Sprintf("first_name=%s", data.FirstName),
		fmt.Sprintf("id=%d", data.ID),
	}
	if data.LastName != "" {
		checkPairs = append(checkPairs, fmt.Sprintf("last_name=%s", data.LastName))
	}
	if data.PhotoURL != "" {
		checkPairs = append(checkPairs, fmt.Sprintf("photo_url=%s", data.PhotoURL))
	}
	if data.Username != "" {
		checkPairs = append(checkPairs, fmt.Sprintf("username=%s", data.Username))
	}
	sort.Strings(checkPairs)
	checkString := strings.Join(checkPairs, "\n")

	// 3. 计算 HMAC-SHA256(check_string, secret)
	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(checkString))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	if !strings.EqualFold(data.Hash, expectedHash) {
		return nil, fmt.Errorf("telegram: hash verification failed")
	}

	return &UserInfo{
		Provider:    "telegram",
		ProviderUID: fmt.Sprintf("%d", data.ID),
		OpenID:      fmt.Sprintf("%d", data.ID),
		Nickname:    data.FirstName + " " + data.LastName,
		Avatar:      data.PhotoURL,
	}, nil
}

// 中文：XProvider 定义当前包使用的数据结构或接口。
// English: XProvider defines a data structure or interface used by this package.
// XProvider X (Twitter) OAuth2
type XProvider struct {
	// 中文：ClientID 保存当前结构中的配置或数据值。
	// English: ClientID stores a configuration or data value for this struct.
	ClientID string
	// 中文：ClientSecret 保存当前结构中的配置或数据值。
	// English: ClientSecret stores a configuration or data value for this struct.
	ClientSecret string
}

// 中文：NewXProvider 创建并返回对应组件实例。
// English: NewXProvider creates and returns the corresponding component instance.
func NewXProvider(clientID, clientSecret string) *XProvider {
	return &XProvider{ClientID: clientID, ClientSecret: clientSecret}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (p *XProvider) Name() string { return "x" }

// 中文：AuthorizeURL 执行当前包中的对应流程。
// English: AuthorizeURL executes the corresponding workflow in this package.
func (p *XProvider) AuthorizeURL(state, redirectURL string) string {
	return fmt.Sprintf(
		"https://twitter.com/i/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=tweet.read+users.read&state=%s&code_challenge=challenge&code_challenge_method=plain",
		url.QueryEscape(p.ClientID),
		url.QueryEscape(redirectURL),
		url.QueryEscape(state),
	)
}

// 中文：GetUserInfo 执行当前包中的对应流程。
// English: GetUserInfo executes the corresponding workflow in this package.
func (p *XProvider) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	return p.GetUserInfoWithRedirect(ctx, code, "")
}

// 中文：GetUserInfoWithRedirect 执行当前包中的对应流程。
// English: GetUserInfoWithRedirect executes the corresponding workflow in this package.
func (p *XProvider) GetUserInfoWithRedirect(ctx context.Context, code string, redirectURL string) (*UserInfo, error) {
	form := map[string]string{
		"client_id":     p.ClientID,
		"client_secret": p.ClientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
		"code_verifier": "challenge",
	}
	if redirectURL != "" {
		form["redirect_uri"] = redirectURL
	}
	tokenBody, err := httpPostForm(ctx, "https://api.twitter.com/2/oauth2/token", form)
	if err != nil {
		return nil, fmt.Errorf("x get token: %w", err)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("x parse token: %w", err)
	}

	userBody, err := httpGetWithAuth(ctx, "https://api.twitter.com/2/users/me", tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("x get userinfo: %w", err)
	}

	var userResp struct {
		Data struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"data"`
	}
	if err := json.Unmarshal(userBody, &userResp); err != nil {
		return nil, fmt.Errorf("x parse userinfo: %w", err)
	}

	return &UserInfo{
		Provider:    "x",
		ProviderUID: userResp.Data.ID,
		Nickname:    userResp.Data.Username,
	}, nil
}

// 中文：DingTalkProvider 定义当前包使用的数据结构或接口。
// English: DingTalkProvider defines a data structure or interface used by this package.
// DingTalkProvider 钉钉OAuth
type DingTalkProvider struct {
	// 中文：AppID 保存当前结构中的配置或数据值。
	// English: AppID stores a configuration or data value for this struct.
	AppID string
	// 中文：AppSecret 保存当前结构中的配置或数据值。
	// English: AppSecret stores a configuration or data value for this struct.
	AppSecret string
}

// 中文：NewDingTalkProvider 创建并返回对应组件实例。
// English: NewDingTalkProvider creates and returns the corresponding component instance.
func NewDingTalkProvider(appID, appSecret string) *DingTalkProvider {
	return &DingTalkProvider{AppID: appID, AppSecret: appSecret}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (p *DingTalkProvider) Name() string { return "dingtalk" }

// 中文：AuthorizeURL 执行当前包中的对应流程。
// English: AuthorizeURL executes the corresponding workflow in this package.
func (p *DingTalkProvider) AuthorizeURL(state, redirectURL string) string {
	return fmt.Sprintf(
		"https://login.dingtalk.com/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid&state=%s&prompt=consent",
		url.QueryEscape(p.AppID),
		url.QueryEscape(redirectURL),
		url.QueryEscape(state),
	)
}

// 中文：GetUserInfo 执行当前包中的对应流程。
// English: GetUserInfo executes the corresponding workflow in this package.
func (p *DingTalkProvider) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	// Step 1: 获取access_token
	tokenBody, err := httpPostForm(ctx, "https://api.dingtalk.com/v1.0/oauth2/accessToken", map[string]string{
		"appKey":    p.AppID,
		"appSecret": p.AppSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("dingtalk get token: %w", err)
	}

	var tokenResp struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("dingtalk parse token: %w", err)
	}

	// Step 2: 获取用户信息
	userBody, err := httpPostJSON(ctx, "https://api.dingtalk.com/v1.0/contact/users/me", map[string]string{
		"tmpCode": code,
	}, tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("dingtalk get userinfo: %w", err)
	}

	var userResp struct {
		OpenID  string `json:"openId"`
		UnionID string `json:"unionId"`
		Mobile  string `json:"mobile"`
		Nick    string `json:"nick"`
		Avatar  string `json:"avatarUrl"`
	}
	if err := json.Unmarshal(userBody, &userResp); err != nil {
		return nil, fmt.Errorf("dingtalk parse userinfo: %w", err)
	}

	return &UserInfo{
		Provider:    "dingtalk",
		ProviderUID: userResp.OpenID,
		OpenID:      userResp.OpenID,
		UnionID:     userResp.UnionID,
		Nickname:    userResp.Nick,
		Avatar:      userResp.Avatar,
	}, nil
}

// 中文：DouyinProvider 定义当前包使用的数据结构或接口。
// English: DouyinProvider defines a data structure or interface used by this package.
// DouyinProvider 抖音OAuth
type DouyinProvider struct {
	// 中文：ClientKey 保存当前结构中的配置或数据值。
	// English: ClientKey stores a configuration or data value for this struct.
	ClientKey string
	// 中文：ClientSecret 保存当前结构中的配置或数据值。
	// English: ClientSecret stores a configuration or data value for this struct.
	ClientSecret string
}

// 中文：NewDouyinProvider 创建并返回对应组件实例。
// English: NewDouyinProvider creates and returns the corresponding component instance.
func NewDouyinProvider(clientKey, clientSecret string) *DouyinProvider {
	return &DouyinProvider{ClientKey: clientKey, ClientSecret: clientSecret}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (p *DouyinProvider) Name() string { return "douyin" }

// 中文：AuthorizeURL 执行当前包中的对应流程。
// English: AuthorizeURL executes the corresponding workflow in this package.
func (p *DouyinProvider) AuthorizeURL(state, redirectURL string) string {
	return fmt.Sprintf(
		"https://open.douyin.com/platform/oauth/connect/?client_key=%s&response_type=code&scope=user_info&redirect_uri=%s&state=%s",
		url.QueryEscape(p.ClientKey),
		url.QueryEscape(redirectURL),
		url.QueryEscape(state),
	)
}

// 中文：GetUserInfo 执行当前包中的对应流程。
// English: GetUserInfo executes the corresponding workflow in this package.
func (p *DouyinProvider) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	tokenBody, err := httpPostForm(ctx, "https://open.douyin.com/oauth/access_token/", map[string]string{
		"client_key":    p.ClientKey,
		"client_secret": p.ClientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
	})
	if err != nil {
		return nil, fmt.Errorf("douyin get token: %w", err)
	}

	var tokenResp struct {
		Data struct {
			AccessToken string `json:"access_token"`
			OpenID      string `json:"open_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("douyin parse token: %w", err)
	}

	userBody, err := httpGetWithAuth(ctx, "https://open.douyin.com/oauth/userinfo/?open_id="+tokenResp.Data.OpenID, tokenResp.Data.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("douyin get userinfo: %w", err)
	}

	var userResp struct {
		Data struct {
			OpenID   string `json:"open_id"`
			UnionID  string `json:"union_id"`
			Nickname string `json:"nickname"`
			Avatar   string `json:"avatar_larger"`
		} `json:"data"`
	}
	if err := json.Unmarshal(userBody, &userResp); err != nil {
		return nil, fmt.Errorf("douyin parse userinfo: %w", err)
	}

	return &UserInfo{
		Provider:    "douyin",
		ProviderUID: userResp.Data.OpenID,
		OpenID:      userResp.Data.OpenID,
		UnionID:     userResp.Data.UnionID,
		Nickname:    userResp.Data.Nickname,
		Avatar:      userResp.Data.Avatar,
	}, nil
}

// 中文：WorkWechatProvider 定义当前包使用的数据结构或接口。
// English: WorkWechatProvider defines a data structure or interface used by this package.
// WorkWechatProvider 企业微信OAuth
type WorkWechatProvider struct {
	// 中文：CorpID 保存当前结构中的配置或数据值。
	// English: CorpID stores a configuration or data value for this struct.
	CorpID string
	// 中文：AgentID 保存当前结构中的配置或数据值。
	// English: AgentID stores a configuration or data value for this struct.
	AgentID string
	// 中文：CorpSecret 保存当前结构中的配置或数据值。
	// English: CorpSecret stores a configuration or data value for this struct.
	CorpSecret string
}

// 中文：NewWorkWechatProvider 创建并返回对应组件实例。
// English: NewWorkWechatProvider creates and returns the corresponding component instance.
func NewWorkWechatProvider(corpID, agentID, corpSecret string) *WorkWechatProvider {
	return &WorkWechatProvider{CorpID: corpID, AgentID: agentID, CorpSecret: corpSecret}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (p *WorkWechatProvider) Name() string { return "work_wechat" }

// 中文：AuthorizeURL 执行当前包中的对应流程。
// English: AuthorizeURL executes the corresponding workflow in this package.
func (p *WorkWechatProvider) AuthorizeURL(state, redirectURL string) string {
	return fmt.Sprintf(
		"https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=%s&agentid=%s&redirect_uri=%s&state=%s",
		url.QueryEscape(p.CorpID),
		url.QueryEscape(p.AgentID),
		url.QueryEscape(redirectURL),
		url.QueryEscape(state),
	)
}

// 中文：GetUserInfo 执行当前包中的对应流程。
// English: GetUserInfo executes the corresponding workflow in this package.
func (p *WorkWechatProvider) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	// Step 1: 获取access_token
	tokenURL := fmt.Sprintf(
		"https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s",
		p.CorpID, p.CorpSecret,
	)
	tokenBody, err := httpGet(ctx, tokenURL)
	if err != nil {
		return nil, fmt.Errorf("work_wechat get token: %w", err)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("work_wechat parse token: %w", err)
	}
	if tokenResp.ErrCode != 0 {
		return nil, fmt.Errorf("work_wechat token error: %d %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	// Step 2: 获取用户身份
	userURL := fmt.Sprintf(
		"https://qyapi.weixin.qq.com/cgi-bin/auth/getuserinfo?access_token=%s&code=%s",
		tokenResp.AccessToken, code,
	)
	userBody, err := httpGet(ctx, userURL)
	if err != nil {
		return nil, fmt.Errorf("work_wechat get userid: %w", err)
	}

	var userResp struct {
		UserID  string `json:"userid"`
		OpenID  string `json:"openid"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(userBody, &userResp); err != nil {
		return nil, fmt.Errorf("work_wechat parse userid: %w", err)
	}
	if userResp.ErrCode != 0 {
		return nil, fmt.Errorf("work_wechat userid error: %d %s", userResp.ErrCode, userResp.ErrMsg)
	}

	return &UserInfo{
		Provider:    "work_wechat",
		ProviderUID: userResp.UserID,
		OpenID:      userResp.OpenID,
	}, nil
}

// ---- HTTP 工具函数 ----

// 中文：httpGet 执行当前包中的对应流程。
// English: httpGet executes the corresponding workflow in this package.
func httpGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return readOAuthResponse(resp)
}

// 中文：httpGetWithAuth 执行当前包中的对应流程。
// English: httpGetWithAuth executes the corresponding workflow in this package.
func httpGetWithAuth(ctx context.Context, url, token string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return readOAuthResponse(resp)
}

// 中文：httpPostForm 执行当前包中的对应流程。
// English: httpPostForm executes the corresponding workflow in this package.
func httpPostForm(ctx context.Context, endpoint string, data map[string]string) ([]byte, error) {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return readOAuthResponse(resp)
}

// 中文：httpPostJSON 执行当前包中的对应流程。
// English: httpPostJSON executes the corresponding workflow in this package.
func httpPostJSON(ctx context.Context, url string, data any, token string) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("x-acs-dingtalk-access-token", token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return readOAuthResponse(resp)
}

// 中文：readOAuthResponse 执行当前包中的对应流程。
// English: readOAuthResponse executes the corresponding workflow in this package.
func readOAuthResponse(resp *http.Response) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("oauth http %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return data, nil
}

// 中文：requireOAuthJSONAccessToken 执行当前包中的对应流程。
// English: requireOAuthJSONAccessToken executes the corresponding workflow in this package.
// extractParam 从URL编码字符串中提取参数
func requireOAuthJSONAccessToken(provider string, body []byte, path ...string) error {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("%s parse refresh response: %w", provider, err)
	}
	if err := oauthPayloadError(provider, payload); err != nil {
		return err
	}
	value, ok := nestedString(payload, path...)
	if !ok || value == "" {
		return fmt.Errorf("%s: access_token not found in refresh response", provider)
	}
	return nil
}

// 中文：requireOAuthQueryAccessToken 执行当前包中的对应流程。
// English: requireOAuthQueryAccessToken executes the corresponding workflow in this package.
func requireOAuthQueryAccessToken(provider string, body []byte) error {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return fmt.Errorf("%s parse refresh response: %w", provider, err)
	}
	if values.Get("access_token") != "" {
		return nil
	}
	if code := values.Get("error"); code != "" {
		return fmt.Errorf("%s api error: %s %s", provider, code, values.Get("error_description"))
	}
	return fmt.Errorf("%s: access_token not found in refresh response", provider)
}

// 中文：oauthPayloadError 执行当前包中的对应流程。
// English: oauthPayloadError executes the corresponding workflow in this package.
func oauthPayloadError(provider string, payload map[string]any) error {
	if code := stringValue(payload["error"]); code != "" {
		return fmt.Errorf("%s api error: %s %s", provider, code, firstJSONMessage(payload, "error_description", "message", "errmsg", "description"))
	}
	if code, ok := nonZeroJSONCode(payload, "errcode", "error_code", "code"); ok {
		return fmt.Errorf("%s api error: %s %s", provider, code, firstJSONMessage(payload, "errmsg", "message", "description"))
	}
	if data, ok := payload["data"].(map[string]any); ok {
		if code, ok := nonZeroJSONCode(data, "errcode", "error_code", "code"); ok {
			return fmt.Errorf("%s api error: %s %s", provider, code, firstJSONMessage(data, "errmsg", "message", "description"))
		}
	}
	return nil
}

// 中文：nestedString 执行当前包中的对应流程。
// English: nestedString executes the corresponding workflow in this package.
func nestedString(payload map[string]any, path ...string) (string, bool) {
	var current any = payload
	for _, key := range path {
		m, ok := current.(map[string]any)
		if !ok {
			return "", false
		}
		current, ok = m[key]
		if !ok {
			return "", false
		}
	}
	value := stringValue(current)
	return value, value != ""
}

// 中文：nonZeroJSONCode 执行当前包中的对应流程。
// English: nonZeroJSONCode executes the corresponding workflow in this package.
func nonZeroJSONCode(payload map[string]any, keys ...string) (string, bool) {
	for _, key := range keys {
		code := stringValue(payload[key])
		if code != "" && code != "0" {
			return code, true
		}
	}
	return "", false
}

// 中文：firstJSONMessage 执行当前包中的对应流程。
// English: firstJSONMessage executes the corresponding workflow in this package.
func firstJSONMessage(payload map[string]any, keys ...string) string {
	for _, key := range keys {
		if value := stringValue(payload[key]); value != "" {
			return value
		}
	}
	return ""
}

// 中文：stringValue 执行当前包中的对应流程。
// English: stringValue executes the corresponding workflow in this package.
func stringValue(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	case int:
		return fmt.Sprintf("%d", v)
	case json.Number:
		return v.String()
	default:
		return ""
	}
}

// 中文：extractParam 执行当前包中的对应流程。
// English: extractParam executes the corresponding workflow in this package.
func extractParam(s, key string) string {
	for _, pair := range strings.Split(s, "&") {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 && kv[0] == key {
			v, _ := url.QueryUnescape(kv[1])
			return v
		}
	}
	return ""
}

// 中文：extractJSONValue 执行当前包中的对应流程。
// English: extractJSONValue executes the corresponding workflow in this package.
// extractJSONValue 从callback包裹的JSON中提取值
func extractJSONValue(s, key string) string {
	// 去除callback(...)包裹
	s = strings.TrimSpace(s)
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start < 0 || end < 0 {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s[start:end+1]), &m); err != nil {
		return ""
	}
	if v, ok := m[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// 中文：extractBotID 执行当前包中的对应流程。
// English: extractBotID executes the corresponding workflow in this package.
// extractBotID 从Telegram Bot Token中提取bot ID
func extractBotID(token string) string {
	parts := strings.SplitN(token, ":", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return token
}
