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
go run ./cmd/spiringo -env development -config configs
```

环境文件加载顺序为：

1. `configs/config.yaml`
2. `configs/config.{env}.yaml`
3. 配置中心 source
4. `SP_` 前缀环境变量

## 数据库迁移

```bash
go run ./cmd/spiringo migrate up -env development -config configs
```

也可以使用脚本：

```bash
ENV=development scripts/migrate.sh
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
