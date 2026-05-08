package oauth

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// 中文：TestTelegramProviderVerifiesCallback 验证相关行为符合预期。
// English: TestTelegramProviderVerifiesCallback verifies the related behavior.
func TestTelegramProviderVerifiesCallback(t *testing.T) {
	token := "123456:secret"
	payload := map[string]any{
		"id":         int64(42),
		"first_name": "Ada",
		"last_name":  "Lovelace",
		"username":   "ada",
		"auth_date":  int64(1700000000),
	}
	payload["hash"] = telegramTestHash(token, []string{
		"auth_date=1700000000",
		"first_name=Ada",
		"id=42",
		"last_name=Lovelace",
		"username=ada",
	})
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	user, err := NewTelegramProvider(token).GetUserInfo(context.Background(), string(body))
	if err != nil {
		t.Fatalf("GetUserInfo returned error: %v", err)
	}
	if user.Provider != "telegram" || user.ProviderUID != "42" || strings.TrimSpace(user.Nickname) != "Ada Lovelace" {
		t.Fatalf("unexpected user info: %+v", user)
	}
}

// 中文：TestReadOAuthResponseRejectsHTTPError 验证相关行为符合预期。
// English: TestReadOAuthResponseRejectsHTTPError verifies the related behavior.
func TestReadOAuthResponseRejectsHTTPError(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusBadGateway,
		Status:     "502 Bad Gateway",
		Body:       ioNopCloser{Reader: bytes.NewReader([]byte("upstream failed"))},
	}
	_, err := readOAuthResponse(resp)
	if err == nil || !strings.Contains(err.Error(), "502") {
		t.Fatalf("expected status error, got %v", err)
	}
}

// 中文：TestWechatProviderRefreshTokenCallsRefreshEndpoint 验证相关行为符合预期。
// English: TestWechatProviderRefreshTokenCallsRefreshEndpoint verifies the related behavior.
func TestWechatProviderRefreshTokenCallsRefreshEndpoint(t *testing.T) {
	withOAuthHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", req.Method)
		}
		if req.URL.Host != "api.weixin.qq.com" || req.URL.Path != "/sns/oauth2/refresh_token" {
			t.Fatalf("url = %s", req.URL.String())
		}
		if req.URL.Query().Get("appid") != "wx-app" || req.URL.Query().Get("refresh_token") != "refresh-1" {
			t.Fatalf("unexpected query: %s", req.URL.RawQuery)
		}
		return oauthTestResponse(http.StatusOK, `{"access_token":"access-1","openid":"openid-1"}`), nil
	})

	if err := NewWechatProvider("wx-app", "secret").RefreshToken(context.Background(), "refresh-1"); err != nil {
		t.Fatalf("RefreshToken returned error: %v", err)
	}
}

// 中文：TestQQProviderRefreshTokenAcceptsQueryResponse 验证相关行为符合预期。
// English: TestQQProviderRefreshTokenAcceptsQueryResponse verifies the related behavior.
func TestQQProviderRefreshTokenAcceptsQueryResponse(t *testing.T) {
	withOAuthHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", req.Method)
		}
		if req.URL.Host != "graph.qq.com" || req.URL.Query().Get("grant_type") != "refresh_token" {
			t.Fatalf("url = %s", req.URL.String())
		}
		return oauthTestResponse(http.StatusOK, `access_token=qq-access&expires_in=777`), nil
	})

	if err := NewQQProvider("qq-app", "secret").RefreshToken(context.Background(), "refresh-1"); err != nil {
		t.Fatalf("RefreshToken returned error: %v", err)
	}
}

// 中文：TestGoogleProviderRefreshTokenPostsRefreshGrant 验证相关行为符合预期。
// English: TestGoogleProviderRefreshTokenPostsRefreshGrant verifies the related behavior.
func TestGoogleProviderRefreshTokenPostsRefreshGrant(t *testing.T) {
	withOAuthHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", req.Method)
		}
		if req.URL.String() != "https://oauth2.googleapis.com/token" {
			t.Fatalf("url = %s", req.URL.String())
		}
		if err := req.ParseForm(); err != nil {
			t.Fatalf("ParseForm returned error: %v", err)
		}
		if req.PostForm.Get("client_id") != "google-client" ||
			req.PostForm.Get("client_secret") != "secret" ||
			req.PostForm.Get("refresh_token") != "refresh-1" ||
			req.PostForm.Get("grant_type") != "refresh_token" {
			t.Fatalf("unexpected form: %v", req.PostForm)
		}
		return oauthTestResponse(http.StatusOK, `{"access_token":"google-access","expires_in":3600}`), nil
	})

	if err := NewGoogleProvider("google-client", "secret").RefreshToken(context.Background(), "refresh-1"); err != nil {
		t.Fatalf("RefreshToken returned error: %v", err)
	}
}

// 中文：TestGoogleProviderTokenExchangeUsesRedirectURL 验证相关行为符合预期。
// English: TestGoogleProviderTokenExchangeUsesRedirectURL verifies the related behavior.
func TestGoogleProviderTokenExchangeUsesRedirectURL(t *testing.T) {
	withOAuthHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.String() {
		case "https://oauth2.googleapis.com/token":
			if req.Method != http.MethodPost {
				t.Fatalf("method = %s, want POST", req.Method)
			}
			if err := req.ParseForm(); err != nil {
				t.Fatalf("ParseForm returned error: %v", err)
			}
			if req.PostForm.Get("code") != "auth-code" ||
				req.PostForm.Get("grant_type") != "authorization_code" ||
				req.PostForm.Get("redirect_uri") != "https://app.example.com/oauth/callback" {
				t.Fatalf("unexpected form: %v", req.PostForm)
			}
			return oauthTestResponse(http.StatusOK, `{"access_token":"google-access","expires_in":3600}`), nil
		case "https://www.googleapis.com/oauth2/v2/userinfo":
			if req.Header.Get("Authorization") != "Bearer google-access" {
				t.Fatalf("authorization = %q", req.Header.Get("Authorization"))
			}
			return oauthTestResponse(http.StatusOK, `{"id":"google-user","email":"u@example.com","name":"User"}`), nil
		default:
			t.Fatalf("unexpected url = %s", req.URL.String())
			return nil, nil
		}
	})

	user, err := NewGoogleProvider("google-client", "secret").GetUserInfoWithRedirect(
		context.Background(),
		"auth-code",
		"https://app.example.com/oauth/callback",
	)
	if err != nil {
		t.Fatalf("GetUserInfoWithRedirect returned error: %v", err)
	}
	if user.ProviderUID != "google-user" || user.Email != "u@example.com" {
		t.Fatalf("unexpected user: %+v", user)
	}
}

// 中文：TestDouyinProviderRefreshTokenReadsNestedAccessToken 验证相关行为符合预期。
// English: TestDouyinProviderRefreshTokenReadsNestedAccessToken verifies the related behavior.
func TestDouyinProviderRefreshTokenReadsNestedAccessToken(t *testing.T) {
	withOAuthHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", req.Method)
		}
		if req.URL.String() != "https://open.douyin.com/oauth/refresh_token/" {
			t.Fatalf("url = %s", req.URL.String())
		}
		if err := req.ParseForm(); err != nil {
			t.Fatalf("ParseForm returned error: %v", err)
		}
		if req.PostForm.Get("client_key") != "douyin-key" ||
			req.PostForm.Get("refresh_token") != "refresh-1" ||
			req.PostForm.Get("grant_type") != "refresh_token" {
			t.Fatalf("unexpected form: %v", req.PostForm)
		}
		return oauthTestResponse(http.StatusOK, `{"data":{"access_token":"douyin-access"},"error_code":0}`), nil
	})

	if err := NewDouyinProvider("douyin-key", "secret").RefreshToken(context.Background(), "refresh-1"); err != nil {
		t.Fatalf("RefreshToken returned error: %v", err)
	}
}

// 中文：TestDingTalkProviderRefreshTokenPostsUserAccessTokenRequest 验证相关行为符合预期。
// English: TestDingTalkProviderRefreshTokenPostsUserAccessTokenRequest verifies the related behavior.
func TestDingTalkProviderRefreshTokenPostsUserAccessTokenRequest(t *testing.T) {
	withOAuthHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", req.Method)
		}
		if req.URL.String() != "https://api.dingtalk.com/v1.0/oauth2/userAccessToken" {
			t.Fatalf("url = %s", req.URL.String())
		}
		var payload map[string]string
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if payload["clientId"] != "ding-app" ||
			payload["clientSecret"] != "secret" ||
			payload["refreshToken"] != "refresh-1" ||
			payload["grantType"] != "refresh_token" {
			t.Fatalf("unexpected payload: %v", payload)
		}
		return oauthTestResponse(http.StatusOK, `{"accessToken":"ding-access","expireIn":7200}`), nil
	})

	if err := NewDingTalkProvider("ding-app", "secret").RefreshToken(context.Background(), "refresh-1"); err != nil {
		t.Fatalf("RefreshToken returned error: %v", err)
	}
}

// 中文：ioNopCloser 定义当前包使用的数据结构或接口。
// English: ioNopCloser defines a data structure or interface used by this package.
type ioNopCloser struct {
	// 中文：*bytes.Reader 嵌入复用该类型提供的能力。
	// English: *bytes.Reader embeds reusable behavior from that type.
	*bytes.Reader
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (c ioNopCloser) Close() error { return nil }

// 中文：oauthRoundTripper 定义当前包使用的数据结构或接口。
// English: oauthRoundTripper defines a data structure or interface used by this package.
type oauthRoundTripper func(*http.Request) (*http.Response, error)

// 中文：RoundTrip 执行当前包中的对应流程。
// English: RoundTrip executes the corresponding workflow in this package.
func (f oauthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// 中文：withOAuthHTTPClient 执行当前包中的对应流程。
// English: withOAuthHTTPClient executes the corresponding workflow in this package.
func withOAuthHTTPClient(t *testing.T, fn oauthRoundTripper) {
	t.Helper()
	oldClient := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: fn}
	t.Cleanup(func() {
		http.DefaultClient = oldClient
	})
}

// 中文：oauthTestResponse 执行当前包中的对应流程。
// English: oauthTestResponse executes the corresponding workflow in this package.
func oauthTestResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Body:       ioNopCloser{Reader: bytes.NewReader([]byte(body))},
		Header:     make(http.Header),
	}
}

// 中文：telegramTestHash 执行当前包中的对应流程。
// English: telegramTestHash executes the corresponding workflow in this package.
func telegramTestHash(token string, pairs []string) string {
	secret := sha256.Sum256([]byte(token))
	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(strings.Join(pairs, "\n")))
	return hex.EncodeToString(mac.Sum(nil))
}
