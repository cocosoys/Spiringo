package oauth

import (
	"context"
	"fmt"
)

// 中文：UserInfo 定义当前包使用的数据结构或接口。
// English: UserInfo defines a data structure or interface used by this package.
// UserInfo is the runtime user profile returned by OAuth providers.
type UserInfo struct {
	// 中文：Provider 保存当前结构中的配置或数据值。
	// English: Provider stores a configuration or data value for this struct.
	Provider string `json:"provider"`
	// 中文：ProviderUID 保存当前结构中的配置或数据值。
	// English: ProviderUID stores a configuration or data value for this struct.
	ProviderUID string `json:"provider_uid"`
	// 中文：OpenID 保存当前结构中的配置或数据值。
	// English: OpenID stores a configuration or data value for this struct.
	OpenID string `json:"open_id,omitempty"`
	// 中文：UnionID 保存当前结构中的配置或数据值。
	// English: UnionID stores a configuration or data value for this struct.
	UnionID string `json:"union_id,omitempty"`
	// 中文：Username 保存当前结构中的配置或数据值。
	// English: Username stores a configuration or data value for this struct.
	Username string `json:"username,omitempty"`
	// 中文：Nickname 保存当前结构中的配置或数据值。
	// English: Nickname stores a configuration or data value for this struct.
	Nickname string `json:"nickname,omitempty"`
	// 中文：Avatar 保存当前结构中的配置或数据值。
	// English: Avatar stores a configuration or data value for this struct.
	Avatar string `json:"avatar,omitempty"`
	// 中文：Email 保存当前结构中的配置或数据值。
	// English: Email stores a configuration or data value for this struct.
	Email string `json:"email,omitempty"`
	// 中文：Phone 保存当前结构中的配置或数据值。
	// English: Phone stores a configuration or data value for this struct.
	Phone string `json:"phone,omitempty"`
	// 中文：RawData 保存当前结构中的配置或数据值。
	// English: RawData stores a configuration or data value for this struct.
	RawData map[string]any `json:"raw_data,omitempty"`
}

// 中文：Provider 定义当前包使用的数据结构或接口。
// English: Provider defines a data structure or interface used by this package.
// Provider is the legacy runtime OAuth provider contract used by auth service.
type Provider interface {
	// 中文：Name 声明该接口需要实现的行为。
	// English: Name declares behavior required by this interface.
	Name() string
	// 中文：AuthorizeURL 声明该接口需要实现的行为。
	// English: AuthorizeURL declares behavior required by this interface.
	AuthorizeURL(state, redirectURL string) string
	// 中文：GetUserInfo 声明该接口需要实现的行为。
	// English: GetUserInfo declares behavior required by this interface.
	GetUserInfo(ctx context.Context, code string) (*UserInfo, error)
}

// 中文：RedirectProvider 定义当前包使用的数据结构或接口。
// English: RedirectProvider defines a data structure or interface used by this package.
// RedirectProvider is implemented by providers whose authorization-code token
// exchange must include the same redirect_uri used by the authorize URL.
type RedirectProvider interface {
	// 中文：Provider 声明该接口需要实现的行为。
	// English: Provider declares behavior required by this interface.
	Provider
	// 中文：GetUserInfoWithRedirect 声明该接口需要实现的行为。
	// English: GetUserInfoWithRedirect declares behavior required by this interface.
	GetUserInfoWithRedirect(ctx context.Context, code string, redirectURL string) (*UserInfo, error)
}

// 中文：OAuthUser 定义当前包使用的数据结构或接口。
// English: OAuthUser defines a data structure or interface used by this package.
// OAuthUser matches the blueprint-level OAuth user contract.
type OAuthUser struct {
	// 中文：Provider 保存当前结构中的配置或数据值。
	// English: Provider stores a configuration or data value for this struct.
	Provider string `json:"provider"`
	// 中文：ProviderID 保存当前结构中的配置或数据值。
	// English: ProviderID stores a configuration or data value for this struct.
	ProviderID string `json:"provider_id"`
	// 中文：Username 保存当前结构中的配置或数据值。
	// English: Username stores a configuration or data value for this struct.
	Username string `json:"username,omitempty"`
	// 中文：Nickname 保存当前结构中的配置或数据值。
	// English: Nickname stores a configuration or data value for this struct.
	Nickname string `json:"nickname,omitempty"`
	// 中文：Avatar 保存当前结构中的配置或数据值。
	// English: Avatar stores a configuration or data value for this struct.
	Avatar string `json:"avatar,omitempty"`
	// 中文：Email 保存当前结构中的配置或数据值。
	// English: Email stores a configuration or data value for this struct.
	Email string `json:"email,omitempty"`
	// 中文：Phone 保存当前结构中的配置或数据值。
	// English: Phone stores a configuration or data value for this struct.
	Phone string `json:"phone,omitempty"`
	// 中文：RawData 保存当前结构中的配置或数据值。
	// English: RawData stores a configuration or data value for this struct.
	RawData map[string]any `json:"raw_data,omitempty"`
}

// 中文：OAuthProvider 定义当前包使用的数据结构或接口。
// English: OAuthProvider defines a data structure or interface used by this package.
// OAuthProvider is the blueprint-compatible OAuth provider interface.
type OAuthProvider interface {
	// 中文：Name 声明该接口需要实现的行为。
	// English: Name declares behavior required by this interface.
	Name() string
	// 中文：AuthURL 声明该接口需要实现的行为。
	// English: AuthURL declares behavior required by this interface.
	AuthURL(ctx context.Context, state string, redirectURL string) (string, error)
	// 中文：GetUser 声明该接口需要实现的行为。
	// English: GetUser declares behavior required by this interface.
	GetUser(ctx context.Context, code string, redirectURL string) (*OAuthUser, error)
	// 中文：RefreshToken 声明该接口需要实现的行为。
	// English: RefreshToken declares behavior required by this interface.
	RefreshToken(ctx context.Context, refreshToken string) error
}

// 中文：ErrRefreshTokenUnsupported 声明当前包使用的变量。
// English: ErrRefreshTokenUnsupported declares variables used by this package.
var ErrRefreshTokenUnsupported = fmt.Errorf("oauth refresh token is not supported by this provider")

// 中文：toOAuthUser 执行当前包中的对应流程。
// English: toOAuthUser executes the corresponding workflow in this package.
func toOAuthUser(info *UserInfo) *OAuthUser {
	if info == nil {
		return nil
	}
	return &OAuthUser{
		Provider:   info.Provider,
		ProviderID: firstNonEmpty(info.ProviderUID, info.OpenID, info.UnionID),
		Username:   info.Username,
		Nickname:   info.Nickname,
		Avatar:     info.Avatar,
		Email:      info.Email,
		Phone:      info.Phone,
		RawData:    info.RawData,
	}
}

// 中文：firstNonEmpty 执行当前包中的对应流程。
// English: firstNonEmpty executes the corresponding workflow in this package.
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

// 中文：Registry 定义当前包使用的数据结构或接口。
// English: Registry defines a data structure or interface used by this package.
// Registry stores OAuth providers by name.
type Registry struct {
	// 中文：providers 保存当前结构中的配置或数据值。
	// English: providers stores a configuration or data value for this struct.
	providers map[string]Provider
}

// 中文：NewRegistry 创建并返回对应组件实例。
// English: NewRegistry creates and returns the corresponding component instance.
// NewRegistry creates an OAuth provider registry.
func NewRegistry() *Registry {
	return &Registry{providers: make(map[string]Provider)}
}

// 中文：Register 执行当前包中的对应流程。
// English: Register executes the corresponding workflow in this package.
// Register stores a provider by its Name. Nil providers are ignored.
func (r *Registry) Register(p Provider) {
	if r == nil || p == nil {
		return
	}
	r.providers[p.Name()] = p
}

// 中文：Get 执行当前包中的对应流程。
// English: Get executes the corresponding workflow in this package.
// Get returns a provider by name.
func (r *Registry) Get(name string) (Provider, bool) {
	if r == nil {
		return nil, false
	}
	p, ok := r.providers[name]
	return p, ok
}

// 中文：Names 执行当前包中的对应流程。
// English: Names executes the corresponding workflow in this package.
// Names returns all registered provider names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.providers))
	for k := range r.providers {
		names = append(names, k)
	}
	return names
}

// 中文：defaultRegistry 声明当前包使用的变量。
// English: defaultRegistry declares variables used by this package.
var defaultRegistry = NewRegistry()

// 中文：Register 执行当前包中的对应流程。
// English: Register executes the corresponding workflow in this package.
func Register(p Provider) {
	defaultRegistry.Register(p)
}

// 中文：Get 执行当前包中的对应流程。
// English: Get executes the corresponding workflow in this package.
func Get(name string) (Provider, bool) {
	return defaultRegistry.Get(name)
}
