package oauth

import "testing"

// 中文：TestProvidersImplementBlueprintOAuthProvider 验证相关行为符合预期。
// English: TestProvidersImplementBlueprintOAuthProvider verifies the related behavior.
func TestProvidersImplementBlueprintOAuthProvider(t *testing.T) {
	providers := []OAuthProvider{
		NewWechatProvider("app", "secret"),
		NewQQProvider("app", "secret"),
		NewGoogleProvider("client", "secret"),
		NewDiscordProvider("client", "secret"),
		NewTelegramProvider("123:token"),
		NewXProvider("client", "secret"),
		NewDingTalkProvider("app", "secret"),
		NewDouyinProvider("client", "secret"),
		NewWorkWechatProvider("corp", "agent", "secret"),
	}
	for _, provider := range providers {
		if provider.Name() == "" {
			t.Fatalf("empty provider name for %T", provider)
		}
	}
}

// 中文：TestToOAuthUserUsesProviderUIDAsProviderID 验证相关行为符合预期。
// English: TestToOAuthUserUsesProviderUIDAsProviderID verifies the related behavior.
func TestToOAuthUserUsesProviderUIDAsProviderID(t *testing.T) {
	user := toOAuthUser(&UserInfo{
		Provider:    "test",
		ProviderUID: "uid-1",
		OpenID:      "open-1",
		Username:    "user",
		Nickname:    "nick",
		Email:       "user@example.com",
		Phone:       "10086",
	})
	if user.ProviderID != "uid-1" || user.Username != "user" || user.Phone != "10086" {
		t.Fatalf("unexpected OAuthUser: %+v", user)
	}
}
