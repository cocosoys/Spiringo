# RBAC/ABAC 模块

RBAC 模块提供角色、权限和角色权限绑定；ABAC 引擎支持基于请求属性、租户和用户上下文的策略判定。

## 配置

```yaml
modules:
  rbac:
    enabled: true
    auth_required: true
```

## 路由

- `GET /api/v1/rbac/roles`
- `POST /api/v1/rbac/roles`
- `PUT /api/v1/rbac/roles/:id`
- `DELETE /api/v1/rbac/roles/:id`
- `GET /api/v1/rbac/permissions`
- `POST /api/v1/rbac/roles/:id/permissions`

当 `auth_required=true` 时，路由会经过 JWT 校验和权限中间件。
