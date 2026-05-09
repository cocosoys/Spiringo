# Spiringo

[中文](README.md) | [English](README_EN.md)

Spiringo is a Go/Gin modular-monolith backend foundation for enterprise systems, SaaS-style multi-tenant platforms, payment/auth/RBAC-heavy applications, and secondary development. It follows a modular-monolith-first architecture while keeping clear boundaries for future service extraction. Application lifecycle, configuration, dependency injection, middleware, persistence, cache, messaging, observability, business modules, code generation, and deployment assets are organized in one coherent repository.

## Project Positioning

Spiringo is not an empty demo template. It is a practical backend foundation with working business modules and infrastructure adapters. It is useful when you need to:

- Bootstrap a Go/Gin backend while keeping the codebase modular.
- Build systems with users, tenants, RBAC/ABAC, authentication, payments, QR codes, and notifications.
- Standardize configuration, logging, tracing, metrics, alerts, cache, locks, MQ, storage, and search in an internal framework.
- Generate modules, CRUD layers, payment channels, OAuth providers, and migrations through a CLI.

## Core Features

- **Module lifecycle**: registration, dependency ordering, enable/disable controls, initialization, start/stop hooks, migrations, and runtime snapshots.
- **Configuration**: local YAML/TOML, environment variables, Nacos, Consul, and change callbacks.
- **Dependency injection**: lightweight DI container with type, name, and interface registration.
- **HTTP service**: Gin wrapper, health checks, readiness checks, and system report endpoints.
- **Middleware chain**: Recovery, Request ID, CORS, I18n, rate limiting, tenant context, JWT auth, idempotency, circuit breaker, metrics, and tracing.
- **Data and infrastructure**: GORM ORM, MySQL/PostgreSQL/SQLite, read replicas, tenant routing, sharding, MongoDB, Redis cache, multi-level cache, distributed locks, object storage, and search engine adapters.
- **Messaging and tasks**: in-memory task queue, Redis Streams, RabbitMQ, Kafka, and event-bus-to-MQ bridge.
- **Observability**: structured logging, Zap adapter, Prometheus metrics, Markdown metrics report, OpenTelemetry/OTLP, Sentry and webhook alerts.
- **Business modules**: auth, user, tenant, RBAC/ABAC, payment, QR code, and notification.
- **Deployment assets**: Docker, Docker Compose, Kubernetes, Helm, Prometheus/Grafana, and Serverless examples.

## Business Modules

| Module | Description |
| --- | --- |
| `auth` | Login, registration, JWT, refresh tokens, OAuth authorization and binding. |
| `user` | User profile, password hashing, default tenant admin, and user interface for auth. |
| `tenant` | Tenant create, read, update, delete, and tenant events. |
| `rbac` | Roles, permissions, user-role mappings, default permissions, RBAC and ABAC access checks. |
| `payment` | WeChat Pay, Alipay, UnionPay, CloudPay, Digital RMB, Stripe, and PayPal; supports payment creation, callbacks, refunds, closing, and events. |
| `qrcode` | QR code generation, parsing, short links, scan logs, and statistics. |
| `notification` | Webhook, email, inbox persistence, event subscriptions, inbox listing, and read marking. |

## Repository Layout

```text
cmd/
  spiringo/              HTTP server and migration command entrypoint
  spiringo-cli/          code generation CLI
  spiringo-serverless/   Serverless entry example
configs/                multi-environment configuration
docs/                   architecture, API, module, CLI, and deployment docs
deployments/            Docker, Compose, Kubernetes, Helm, observability, and Serverless assets
internal/core/          framework core: App, Config, DI, Event, Middleware, Module, Server, Migration
internal/modules/       business modules: auth, user, tenant, rbac, payment, qrcode, notification
internal/pkg/           infrastructure packages: ORM, Cache, Lock, MQ, Storage, Search, Trace, Metrics, etc.
pkg/                    public shared types, responses, error codes, and constants
```

## Quick Start

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Run Tests

```bash
go test ./...
```

On Windows, SQLite/CGO-related packages are more reliable with explicit environment variables:

```powershell
$env:GOARCH='amd64'
$env:CGO_ENABLED='1'
$env:GOCACHE=(Join-Path (Get-Location) '.cache\go-build')
go test ./...
```

### 3. Build Entrypoints

```bash
go build ./cmd/spiringo ./cmd/spiringo-cli ./cmd/spiringo-serverless
```

### 4. Start the Development Server

```bash
go run ./cmd/spiringo -env local -config configs
```

### 5. Run Migrations Only

```bash
go run ./cmd/spiringo migrate up -env local -config configs
```

## CLI Usage

Spiringo provides `spiringo-cli` for common code generation tasks:

```bash
go run ./cmd/spiringo-cli module order
go run ./cmd/spiringo-cli crud order product
go run ./cmd/spiringo-cli payment-channel bankpay
go run ./cmd/spiringo-cli oauth-provider github
go run ./cmd/spiringo-cli migrate create add_product_table
```

See the [CLI documentation](docs/cli.zh-CN.md) for details.

## Configuration and Environments

Configuration is centralized in `configs/config.yaml`. The active runtime environment is marked by `app.env`, which supports `local`, `dev`, and `prod`; profile-specific overrides live in the same file under `environments.<env>`:

- `configs/config.yaml`

Use `-env` or `APP_ENV` to temporarily override `app.env` and select the matching profile section. `development` maps to `dev`, and `production` maps to `prod` for compatibility. The config manager merges the local file, optional config-center sources, and `SP_` environment variables.

Server listen and external access settings are centralized under `server`: `server.host` / `server.port` build the listen address, `server.addr` remains a full-address override, and `server.public_url` plus `server.api_base_url` are intended for payment callbacks, OAuth callbacks, frontend integration, and deployment docs. When reusing the project, update this block and the matching environment profile first.

## System Endpoints

After startup, the application registers these system endpoints:

- `GET /health`: liveness check.
- `GET /ready`: readiness check, including database ping when configured.
- `GET /system/report`: Markdown runtime report.
- `GET /system/report.json`: structured runtime report.
- `GET /metrics`: Prometheus metrics.
- `GET /metrics/report`: Markdown metrics report.

## Documentation

- [Architecture Blueprint](docs/architecture.md)
- [Getting Started](docs/getting-started.md)
- [OpenAPI](docs/api/openapi.yaml)
- [CLI](docs/cli.zh-CN.md)
- [Deployment](docs/deployment.zh-CN.md)
- [Auth Module](docs/modules/auth.md)
- [User Module](docs/modules/user.md)
- [Tenant Module](docs/modules/tenant.md)
- [RBAC Module](docs/modules/rbac.md)
- [Payment Module](docs/modules/payment.md)
- [QR Code Module](docs/modules/qrcode.md)
- [Notification Module](docs/modules/notification.md)

## Deployment

The repository includes multiple deployment surfaces:

- Dockerfile: `deployments/docker/`
- Docker Compose: `deployments/docker-compose/`
- Kubernetes: `deployments/kubernetes/`
- Helm Chart: `deployments/helm/spiringo/`
- Prometheus/Grafana: `deployments/observability/`
- Serverless examples: `deployments/serverless/`

Production deployments must provide real credentials and service endpoints for databases, Redis, MQ, object storage, search, payment channels, OAuth providers, and alert sinks.

## Current Status

The codebase currently covers the main framework layer, business modules, CLI generator, deployment templates, and documentation described by the architecture blueprint. Third-party integrations such as payment gateways, OAuth providers, email, webhooks, object storage, search, and MQ have adapter implementations, but they still need real credentials and live-environment integration testing before production use.
