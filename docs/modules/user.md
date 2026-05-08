# 用户模块

用户模块提供基础用户 CRUD，并为认证模块注册 `auth_user_service` 命名依赖。

## 配置

```yaml
modules:
  user:
    enabled: true
    default_admin:
      enabled: true
      username: "admin_%s"
      password: "changeme"
```

## 路由

- `GET /api/v1/user`
- `POST /api/v1/user`
- `GET /api/v1/user/:id`
- `PUT /api/v1/user/:id`
- `DELETE /api/v1/user/:id`

创建用户会发布 `user.created`，更新和删除分别发布 `user.updated`、`user.deleted`。
