# 通知模块

通知模块支持 webhook、email 和站内信 inbox。站内信路径会持久化消息，并提供列表和已读标记接口。

## 配置

```yaml
modules:
  notification:
    enabled: true
    events:
      - "payment.failed"
      - "tenant.suspended"
    inbox:
      enabled: true
```

## 路由

- `POST /api/v1/notification/send`
- `GET /api/v1/notification/inbox`
- `PUT /api/v1/notification/inbox/:id/read`

## 订阅

模块会根据 `modules.notification.events` 订阅事件总线，并把事件转换为通知消息。
