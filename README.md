# Spiringo

[中文](README.md) | [English](README_EN.md)

Spiringo 是一个基于 Go 与 Gin 的模块化单体后端基座，目标是为企业内部系统、SaaS 多租户平台、支付/认证/权限类业务和二次开发项目提供可直接扩展的工程骨架。项目采用“模块化单体优先，必要时可拆分微服务”的设计路线，将应用生命周期、配置、依赖注入、中间件、数据库、缓存、消息队列、可观测性、业务模块和部署模板统一放在一个清晰的代码结构中。

## 项目定位

Spiringo 不是一个只包含示例代码的空模板，而是一个已经具备实际业务模块和基础设施适配层的后端平台基座。它适合以下场景：

- 快速启动 Go/Gin 后端项目，并保留后续模块拆分空间。
- 构建带用户、租户、RBAC/ABAC、认证、支付、二维码和通知能力的业务系统。
- 为内部框架沉淀统一的配置、日志、链路追踪、指标、告警、缓存、锁、MQ、存储和搜索抽象。
- 通过 CLI 生成模块、CRUD、支付通道、OAuth Provider 和迁移文件，降低重复样板代码成本。

## 核心能力

- **模块生命周期**：模块注册、依赖排序、启用/禁用、初始化、启动、停止、迁移和运行快照。
- **配置管理**：本地 YAML/TOML、环境变量、Nacos、Consul，支持配置监听和热更新回调。
- **依赖注入**：轻量级 DI 容器，支持按类型、名称和接口注入。
- **HTTP 服务**：Gin 封装、健康检查、就绪检查、系统报告接口。
- **中间件链**：Recovery、Request ID、CORS、I18n、限流、租户上下文、JWT 认证、幂等、熔断、指标、Trace。
- **数据与基础设施**：GORM ORM、MySQL/PostgreSQL/SQLite、读写分离、多租户路由、分片、MongoDB、Redis 缓存、多级缓存、分布式锁、对象存储、搜索引擎。
- **消息与任务**：内存任务队列、Redis Streams、RabbitMQ、Kafka，以及事件总线到 MQ 的桥接。
- **可观测性**：结构化日志、Zap 适配、Prometheus 指标、Markdown 指标报告、OpenTelemetry/OTLP、Sentry/Webhook 告警。
- **业务模块**：认证、用户、租户、RBAC/ABAC、支付、二维码、通知。
- **部署资产**：Docker、Docker Compose、Kubernetes、Helm、Prometheus/Grafana、Serverless 示例。

## 业务模块

| 模块 | 说明 |
| --- | --- |
| `auth` | 用户登录、注册、JWT、刷新令牌、OAuth 授权与绑定。 |
| `user` | 用户资料、密码哈希、默认租户管理员、认证模块用户接口。 |
| `tenant` | 租户创建、查询、更新、删除和租户事件。 |
| `rbac` | 角色、权限、用户角色关系、默认权限、RBAC 与 ABAC 访问控制。 |
| `payment` | 微信支付、支付宝、银联、云闪付、数字人民币、Stripe、PayPal；支持创建支付、回调、退款、关闭和事件发布。 |
| `qrcode` | 二维码生成、解析、短链跳转、扫码日志和统计。 |
| `notification` | Webhook、Email、站内信、事件订阅、收件箱列表与已读标记。 |

## 目录结构

```text
cmd/
  spiringo/              主服务与迁移命令入口
  spiringo-cli/          代码生成 CLI
  spiringo-serverless/   Serverless 入口示例
configs/                多环境配置
docs/                   架构、API、模块、CLI、部署文档
deployments/            Docker、Compose、Kubernetes、Helm、观测和 Serverless 模板
internal/core/          框架核心：App、Config、DI、Event、Middleware、Module、Server、Migration
internal/modules/       业务模块：auth、user、tenant、rbac、payment、qrcode、notification
internal/pkg/           通用基础设施包：ORM、Cache、Lock、MQ、Storage、Search、Trace、Metrics 等
pkg/                    对外通用类型、响应、错误码和常量
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 运行测试

```bash
go test ./...
```

在 Windows 工作区中，如果使用 SQLite/CGO 相关包，建议显式指定：

```powershell
$env:GOARCH='amd64'
$env:CGO_ENABLED='1'
$env:GOCACHE=(Join-Path (Get-Location) '.cache\go-build')
go test ./...
```

### 3. 构建入口程序

```bash
go build ./cmd/spiringo ./cmd/spiringo-cli ./cmd/spiringo-serverless
```

### 4. 启动开发环境服务

```bash
go run ./cmd/spiringo -env development -config configs
```

### 5. 仅执行迁移

```bash
go run ./cmd/spiringo migrate up -env development -config configs
```

## CLI 用法

Spiringo 提供 `spiringo-cli` 用于生成常见工程代码：

```bash
go run ./cmd/spiringo-cli module order
go run ./cmd/spiringo-cli crud order product
go run ./cmd/spiringo-cli payment-channel bankpay
go run ./cmd/spiringo-cli oauth-provider github
go run ./cmd/spiringo-cli migrate create add_product_table
```

详细说明见 [CLI 文档](docs/cli.zh-CN.md)。

## 配置与环境

配置文件位于 `configs/`，项目提供开发、测试、预发、生产等环境配置示例：

- `configs/config.yaml`
- `configs/config.development.yaml`
- `configs/config.test.yaml`
- `configs/config.staging.yaml`
- `configs/config.production.yaml`

运行时通过 `-env` 和 `-config` 指定环境与配置目录。配置管理器会合并本地文件、环境变量以及可选的 Nacos/Consul 配置源。

## 系统接口

应用启动后会注册基础系统接口：

- `GET /health`：存活检查。
- `GET /ready`：就绪检查，会检查数据库连接。
- `GET /system/report`：Markdown 运行报告。
- `GET /system/report.json`：结构化运行报告。
- `GET /metrics`：Prometheus 指标。
- `GET /metrics/report`：Markdown 指标报告。

## 文档入口

- [架构蓝图](docs/architecture.md)
- [模块化教程目录](docs/module-tutorials/目录.md)
- [快速上手](docs/getting-started.md)
- [OpenAPI](docs/api/openapi.yaml)
- [CLI 文档](docs/cli.zh-CN.md)
- [部署文档](docs/deployment.zh-CN.md)
- [认证模块](docs/modules/auth.md)
- [用户模块](docs/modules/user.md)
- [租户模块](docs/modules/tenant.md)
- [RBAC 模块](docs/modules/rbac.md)
- [支付模块](docs/modules/payment.md)
- [二维码模块](docs/modules/qrcode.md)
- [通知模块](docs/modules/notification.md)

## 部署

项目内置多种部署资产：

- Dockerfile：`deployments/docker/`
- Docker Compose：`deployments/docker-compose/`
- Kubernetes：`deployments/kubernetes/`
- Helm Chart：`deployments/helm/spiringo/`
- Prometheus/Grafana：`deployments/observability/`
- Serverless 示例：`deployments/serverless/`

生产环境需要根据实际数据库、Redis、MQ、对象存储、搜索引擎、支付通道、OAuth Provider 和告警服务配置真实凭据。

## 当前状态

项目主体代码已经覆盖架构蓝图中的核心框架层、主要业务模块、CLI 生成器、部署模板和文档体系。第三方支付、OAuth、邮件、Webhook、对象存储、搜索、MQ 等外部能力已经有适配代码，但上线前仍需要使用真实凭据和真实服务环境进行联调。
