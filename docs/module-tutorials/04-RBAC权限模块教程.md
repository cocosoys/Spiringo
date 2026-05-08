# RBAC 权限模块教程

RBAC 模块提供角色、权限、角色权限绑定，并在 `auth_required=true` 时用认证模块签发的 JWT 保护权限管理接口。

## 1. 复制代码

复制模块代码：

```text
internal/modules/rbac/
  rbac.go
  migrations.go
  dto/
  model/
  repository/
  service/
  handler/
```

同时复制依赖模块：

```text
internal/modules/auth/
internal/modules/user/
internal/modules/tenant/
internal/core/middleware/
internal/pkg/orm/
pkg/types/
```

如果只是离线管理权限，可以把 `auth_required` 设为 `false`，但生产环境不建议这样做。

## 2. 复制依赖

RBAC 模块需要：

```bash
go get github.com/gin-gonic/gin
go get github.com/google/uuid
go get github.com/golang-jwt/jwt/v5
go get gorm.io/gorm
```

`github.com/golang-jwt/jwt/v5` 来自认证模块和认证中间件。

## 3. 启用配置

```yaml
modules:
  rbac:
    enabled: true
    auth_required: true
```

配置含义：

| 配置 | 含义 |
| --- | --- |
| `enabled` | 是否初始化 RBAC 模块 |
| `auth_required` | 是否为 RBAC 管理接口启用 JWT 和权限中间件 |

## 4. 注册模块

RBAC 依赖认证、用户和租户：

```go
application.RegisterModules(
    tenant.NewTenantModule(),
    user.NewUserModule(),
    auth.NewAuthModule(),
    rbac.NewRBACModule(),
)
```

依赖声明：

```go
module.NewBaseModule("rbac", "auth", "user", "tenant")
```

当 `auth_required=true` 时，RBAC 初始化阶段会从 DI 中读取：

```go
di.Resolve[*authsvc.AuthService](app.DI)
```

因此认证模块必须成功初始化。

## 5. 数据迁移

RBAC 迁移 ID：

```text
rbac_001_create_roles_table
```

迁移表：

```go
model.Role{}
model.Permission{}
model.RolePermission{}
model.UserRole{}
```

执行：

```powershell
$env:GOARCH='amd64'
$env:CGO_ENABLED='1'
go run ./cmd/spiringo migrate up -env development -config configs
```

## 6. HTTP 接口

路由前缀：

```text
/api/v1/rbac
```

接口：

| 方法 | 路径 | 权限点 | 用途 |
| --- | --- | --- | --- |
| `GET` | `/api/v1/rbac/roles` | `rbac.roles:read` | 分页查询角色 |
| `POST` | `/api/v1/rbac/roles` | `rbac.roles:create` | 创建角色 |
| `PUT` | `/api/v1/rbac/roles/:id` | `rbac.roles:update` | 更新角色 |
| `DELETE` | `/api/v1/rbac/roles/:id` | `rbac.roles:delete` | 删除角色 |
| `GET` | `/api/v1/rbac/permissions` | `rbac.permissions:read` | 分页查询权限 |
| `POST` | `/api/v1/rbac/roles/:id/permissions` | `rbac.permissions:assign` | 给角色分配权限 |

创建角色请求：

```json
{
  "name": "管理员",
  "code": "admin",
  "description": "系统管理员"
}
```

更新角色请求：

```json
{
  "name": "平台管理员",
  "description": "可管理平台基础资料",
  "status": "active"
}
```

分配权限请求：

```json
{
  "permission_ids": ["perm-id-1", "perm-id-2"]
}
```

## 7. 在业务路由中应用权限

如果你的业务模块也要使用 RBAC，可在路由中组合认证和权限中间件：

```go
protected := r.Group("", middleware.Auth(authService))
protected.POST(
    "/orders",
    middleware.RequirePermission(rbacService, "order.orders", "create"),
    h.CreateOrder,
)
```

权限模型推荐：

```text
resource: order.orders
action: create | read | update | delete | approve | export
```

## 8. ABAC 扩展

RBAC 模块还包含 ABAC 服务，可基于属性判断访问，例如：

- 当前用户是否属于当前租户。
- 当前资源是否归属当前用户。
- 当前操作是否在允许时间段内。

适合在“角色权限不足以表达规则”时补充使用。

## 9. 常见问题

- RBAC 路由返回未认证：确认请求头包含 `Authorization: Bearer <token>`。
- 有 token 但无权限：确认角色和权限绑定已写入，并且权限点的 resource/action 与中间件一致。
- 本地调试想跳过鉴权：可以临时设置 `modules.rbac.auth_required: false`，但不要用于生产。
- 复制到新项目后初始化失败：确认 `auth.NewAuthModule()` 已注册，且认证模块已把 `AuthService` 注入 DI。

