# 二维码模块

二维码模块提供生成、解析、短链跳转和扫码统计。

## 配置

```yaml
modules:
  qrcode:
    enabled: true
    default_size: 256
    default_level: "M"
    oss_prefix: "qrcode/"
    bucket_name: ""
```

`default_level` 和请求中的 `level` 支持二维码标准缩写 `L/M/Q/H`，也兼容 `low/medium/high/highest`。
`oss_prefix` 在接入对象存储时作为对象 key 前缀；未接入对象存储时可作为图片 URL 前缀拼接返回。

## 路由

- `POST /api/v1/qrcode/generate`
- `POST /api/v1/qrcode/parse`
- `GET /api/v1/qrcode/s/:code`
- `GET /api/v1/qrcode/stats/:code`

服务同时暴露蓝图接口 `Generator` 和 `ShortLinkService` 适配器，便于其他模块直接依赖接口调用。
