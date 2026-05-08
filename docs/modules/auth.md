# 认证模块

认证模块提供账号注册、登录、登出、刷新令牌和第三方 OAuth 登录。

## 配置

```yaml
modules:
  auth:
    enabled: true
    jwt:
      secret: "change-me"
      access_expire: "2h"
      refresh_expire: "168h"
      issuer: "spiringo"
    oauth:
      google:
        enabled: false
        client_id: ""
        client_secret: ""
```

## 路由

- `POST /api/v1/auth/login`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/logout`
- `POST /api/v1/auth/refresh`
- `GET /api/v1/auth/oauth/:provider`
- `GET /api/v1/auth/oauth/:provider/callback`

OAuth provider 统一实现 `Name/AuthURL/GetUser/RefreshToken`，内置微信、QQ、Google、Discord、Telegram、X、钉钉、抖音和企业微信。
需要严格校验回调地址的平台可在登录和回调阶段都传入同一个 `redirect_url`，服务会把它贯通到授权码换 token 的 provider 调用。
