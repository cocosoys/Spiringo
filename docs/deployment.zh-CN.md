# Spiringo 部署指南

本文档对应蓝图中的二进制、Docker Compose、Kubernetes/Helm、CI/CD、Serverless 与多环境部署路径。

## 本地二进制

```bash
go build -o bin/spiringo ./cmd/spiringo
APP_ENV=local ./bin/spiringo -config configs
```

Windows PowerShell：

```powershell
go build -o .\bin\spiringo.exe .\cmd\spiringo
$env:APP_ENV = "local"
.\bin\spiringo.exe -config configs
```

## Docker Compose

生产风格全栈：

```bash
docker compose -f deployments/docker-compose/docker-compose.yaml up -d
```

开发全栈：

```bash
docker compose -f deployments/docker-compose/docker-compose.dev.yaml up -d
```

默认包含 MySQL、Redis、MinIO、Elasticsearch、Prometheus 和 Grafana。Grafana 默认账号为 `admin` / `spiringo`。

## Helm 多环境

```bash
helm lint deployments/helm/spiringo -f deployments/helm/spiringo/values-dev.yaml
helm upgrade --install spiringo deployments/helm/spiringo -f deployments/helm/spiringo/values-prod.yaml
```

当前提供：

- `values-dev.yaml`
- `values-test.yaml`（复用 `dev` 应用配置段，测试差异通过 Helm 环境变量覆盖）
- `values-staging.yaml`（复用 `prod` 应用配置段，预发差异通过 Helm 环境变量覆盖）
- `values-prod.yaml`

环境名约定为 `local`、`dev`、`prod`；历史名称 `development` 和 `production` 会分别映射到 `dev` 和 `prod`。生产环境必须替换 Secret、镜像仓库和数据库/缓存连接。

## Serverless

容器型 Serverless 使用独立入口：

```bash
docker build -t spiringo-serverless:latest -f deployments/docker/Dockerfile.serverless .
kubectl apply -f deployments/serverless/knative-service.yaml
```

该入口会读取平台注入的 `PORT` 并设置 `SP_SERVER_PORT`，仍然允许用 `SP_SERVER_ADDR` 直接覆盖完整监听地址。对外访问地址请通过 `SP_SERVER_PUBLIC_URL` 和 `SP_SERVER_API_BASE_URL` 集中配置。AWS Lambda 事件模型需要额外 adapter，当前仓库未引入对应依赖。

## Trace / OTLP

默认 trace exporter 为 `logger`。如果要导出到 OTLP/HTTP collector，可使用配置文件：

```yaml
trace:
  enabled: true
  exporter: "otlp" # logger, otlp, both
  service:
    name: "spiringo"
  otlp:
    endpoint: "http://otel-collector:4318/v1/traces"
    timeout: "5s"
    headers: {}
```

容器或 Helm 环境可设置：

```bash
SP_TRACE_ENABLED=true
SP_TRACE_EXPORTER=otlp
SP_TRACE_OTLP_ENDPOINT=http://otel-collector:4318/v1/traces
```

## 常用 Make 命令

```bash
make test
make build
make build-cli
make build-serverless
make docker-build
make docker-dev
make helm-template ENV=dev
make helm-deploy ENV=prod
```

## 运行时报告

应用启动后会暴露两个轻量级排障接口：

- `GET /system/report`：Markdown 综合报告，包含应用环境、基础设施启用状态、模块生命周期状态和指标摘要。
- `GET /system/report.json`：同一份报告的 JSON 结构，适合接入巡检脚本或控制台。

如果启用了 `metrics`，仍然保留 `GET /metrics` 和 `GET /metrics/report`；其中 `/system/report` 面向整体运行态，`/metrics/report` 只面向指标明细。
