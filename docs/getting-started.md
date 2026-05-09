# Spiringo 快速开始

本文档对应 `docs/architecture.md` 中的开发流程，目标是在不依赖外部服务的情况下先完成代码级验证，再按需接入数据库、Redis、MQ、支付通道和部署环境。

## 环境要求

- Go 1.24 或更高版本
- Windows 本地验证建议使用 `GOARCH=amd64` 和 `CGO_ENABLED=1`
- 可选组件：MySQL/PostgreSQL、Redis、RabbitMQ/Kafka、MinIO/Ceph、Elasticsearch、MongoDB、Nacos/Consul

## 代码级验证

```powershell
$env:GOARCH='amd64'
$env:CGO_ENABLED='1'
$env:GOCACHE=(Join-Path (Get-Location) '.gocache')
go test ./...
go build ./cmd/spiringo ./cmd/spiringo-cli ./cmd/spiringo-serverless
```

## 启动主服务

```bash
go run ./cmd/spiringo -env local -config configs
```

常用环境：

| 环境 | 标记位置 | 用途 |
| --- | --- | --- |
| `local` | `app.env` + `environments.local` | 本地单机运行，建议使用内存组件或本机服务 |
| `dev` | `app.env` + `environments.dev` | 共享开发或容器开发，可由环境变量覆盖为 MySQL/Redis/MQ |
| `prod` | `app.env` + `environments.prod` | 生产发布模式，真实凭据必须通过环境变量、Secret 或配置中心注入 |

兼容别名：`development` 会映射到 `dev`，`production` 会映射到 `prod`。

配置加载顺序为：

1. `configs/config.yaml`
2. `environments.<app.env>` 环境段
3. 配置中心 source
4. `SP_` 前缀环境变量

`-env` 和 `APP_ENV` 只覆盖当前运行环境标记，不会加载额外的 `config.<env>.yaml` 文件。

服务地址集中在 `server` 段：

```yaml
server:
  host: "127.0.0.1"
  port: 8080
  addr: "" # 可选：填写后覆盖 host + port
  public_url: "http://127.0.0.1:8080"
  api_base_url: "http://127.0.0.1:8080/api/v1"
```

本地或新项目优先修改 `host`、`port`、`public_url`、`api_base_url`。容器或云平台也可以用 `SP_SERVER_HOST`、`SP_SERVER_PORT`、`SP_SERVER_PUBLIC_URL`、`SP_SERVER_API_BASE_URL` 覆盖。

## 数据库迁移

```bash
go run ./cmd/spiringo migrate up -env local -config configs
```

也可以使用脚本：

```bash
ENV=local scripts/migrate.sh
```

## 生成代码

```bash
go run ./cmd/spiringo-cli module order
go run ./cmd/spiringo-cli crud order Product
go run ./cmd/spiringo-cli payment-channel bank
go run ./cmd/spiringo-cli oauth-provider github
go run ./cmd/spiringo-cli migrate create add_product_table
```

脚本入口：

```bash
TYPE=module MODULE=order scripts/generate.sh
TYPE=crud MODULE=order MODEL=Product scripts/generate.sh
TYPE=payment-channel MODULE=bank scripts/generate.sh
TYPE=oauth-provider MODULE=github scripts/generate.sh
```

## 主要 API 入口

- 认证：`/api/v1/auth`
- 用户：`/api/v1/user`
- 租户：`/api/v1/tenant`
- RBAC/ABAC：`/api/v1/rbac`
- 支付：`/api/v1/payment`
- 二维码：`/api/v1/qrcode`
- 通知：`/api/v1/notification`
- 系统报告：`/system/report`、`/system/report.json`
- 指标：`/metrics`、`/metrics/report`

完整接口概要见 [api/openapi.yaml](api/openapi.yaml)。
