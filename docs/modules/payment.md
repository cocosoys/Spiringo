# 支付模块

支付模块以 `Channel` 接口统一微信、支付宝、银联、云闪付、数字人民币、Stripe 和 PayPal。

## 配置

```yaml
modules:
  payment:
    enabled: true
    default_notify_url: ""
    wechat:
      enabled: false
      app_id: ""
      mch_id: ""
      api_v3_key: ""
    stripe:
      enabled: false
      secret_key: ""
      webhook_secret: ""
```

## 路由

- `POST /api/v1/payment/create`
- `POST /api/v1/payment/callback/:channel`
- `POST /api/v1/payment/refund`
- `GET /api/v1/payment/query/:id`
- `POST /api/v1/payment/close/:id`

## 事件

- `payment.created`
- `payment.success`
- `payment.failed`
- `payment.refunded`
- `payment.closed`
- `payment.fulfillment_requested`

回调响应由通道的 `CallbackSuccess/CallbackFail` 决定，避免把所有网关强行转换成同一种 JSON。
