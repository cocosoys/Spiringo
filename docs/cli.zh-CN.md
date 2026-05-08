# Spiringo CLI 使用说明

`spiringo-cli` 用于生成项目、业务模块、CRUD、迁移、支付通道和 OAuth Provider。

## 新建项目

```bash
spiringo-cli new billing-app -module example.com/acme/billing
```

会复制当前模板并重写 `go.mod` 的 module 路径。

## 生成模块

```bash
spiringo-cli module inventory
```

会生成：

- `internal/modules/inventory/inventory.go`
- `handler/`
- `service/`
- `repository/`
- `model/`
- `dto/`
- `migrations.go`

如果目标项目存在 `internal/modules/builtin/builtin.go`，CLI 会自动加入 import 和 `inventory.NewInventoryModule()` 注册。
生成的列表接口会在 DTO 层统一处理分页默认值：`page <= 0` 时使用第 1 页，`page_size <= 0` 时使用 20，超过 100 时自动收敛到 100。

## 生成 CRUD

```bash
spiringo-cli crud inventory Product
```

会生成模型、DTO、Repository、Service、Handler 与迁移函数。生成代码会读取目标项目 `go.mod` 的 module 路径，避免硬编码模板仓库路径。
生成的 Handler 和 Service 使用同一套分页默认值，响应中的分页元数据会与实际查询窗口保持一致。

## 生成迁移

```bash
spiringo-cli migrate create add_product_table
```

会生成时间戳前缀的迁移文件。迁移模板只提供结构，具体 schema 变更由开发者填写。

## 生成支付通道

```bash
spiringo-cli payment-channel bank-pay
```

生成的通道是通用 HTTP 网关骨架，默认调用：

- `/payment/create`
- `/payment/refund`
- `/payment/query`

回调验签使用 HMAC-SHA256，与内置 gateway helper 保持一致。真实银行或三方网关字段不一致时，只需要调整生成文件中的 payload 和 response 映射。

## 生成 OAuth Provider

```bash
spiringo-cli oauth-provider github
```

生成的 Provider 是标准 OAuth2 授权码模式骨架，包含授权 URL、带 `redirect_uri` 的 token exchange、refresh token 和 userinfo 读取逻辑。不同平台的字段差异可在生成文件中调整 `firstOAuthString` 的字段列表。
