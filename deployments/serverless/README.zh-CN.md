# Spiringo Serverless 部署

本目录提供基于容器的 Serverless 部署样例，目标是 Knative 与兼容 Knative Service API 的平台，例如 Cloud Run。

## 构建镜像

```bash
docker build -t spiringo-serverless:latest -f deployments/docker/Dockerfile.serverless .
```

`cmd/spiringo-serverless` 会在未显式设置 `SP_SERVER_ADDR` / `SP_SERVER_PORT` 时读取平台注入的 `PORT`，并写入统一的 `server.port` 配置。对外访问地址建议通过 `SP_SERVER_PUBLIC_URL` 和 `SP_SERVER_API_BASE_URL` 显式设置。

## Knative

```bash
kubectl apply -f deployments/serverless/knative-service.yaml
```

部署前需要创建 `spiringo-secrets`，至少包含：

- `database-dsn`
- `redis-addr`
- `jwt-secret`

## Cloud Run

`cloud-run.yaml` 使用 Cloud Run 的 Secret Manager 引用格式。发布前需要替换：

- `gcr.io/PROJECT_ID/spiringo-serverless:latest`
- `spiringo-database-dsn`
- `spiringo-jwt-secret`

## 运行边界

当前 Serverless 入口复用完整 HTTP 应用，不新增平台 SDK。适合容器型 Serverless；如果要接 AWS Lambda API Gateway 事件，需要在允许新增依赖后引入 Lambda adapter。
