# 租户模块

租户模块实现共享库、独立 schema 和独立数据库三种策略的配置入口。默认策略为 `shared_db`。

## 配置

```yaml
modules:
  tenant:
    enabled: true
    strategy: "shared_db"
```

## 路由

- `GET /api/v1/tenant`
- `POST /api/v1/tenant`
- `GET /api/v1/tenant/:id`
- `PUT /api/v1/tenant/:id`
- `DELETE /api/v1/tenant/:id`

请求可通过 `X-Tenant-ID`、`X-Tenant-Name`、`X-Tenant-Strategy` 和 `X-Tenant-DB-Conn` 注入租户上下文。
