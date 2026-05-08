# Spiringo 架构蓝图

> 基于 Go 的全栈微服务平台基座，模块化单体架构，面向二次开发，兼顾内部复用与开源。

---

## 目录

- [一、项目定位与技术选型](#一项目定位与技术选型)
- [二、完整目录结构](#二完整目录结构)
- [三、核心框架层设计](#三核心框架层设计)
  - [3.1 模块注册与生命周期管理](#31-模块注册与生命周期管理)
  - [3.2 依赖注入方案](#32-依赖注入方案)
  - [3.3 中间件链设计](#33-中间件链设计)
  - [3.4 配置管理架构](#34-配置管理架构)
- [四、业务模块详细内部结构](#四业务模块详细内部结构)
  - [4.1 支付模块](#41-支付模块)
  - [4.2 认证模块](#42-认证模块)
- [五、中间件抽象层设计](#五中间件抽象层设计)
- [六、关键接口定义](#六关键接口定义)
  - [6.1 支付统一接口](#61-支付统一接口)
  - [6.2 OAuth 统一接口](#62-oauth-统一接口)
  - [6.3 二维码接口](#63-二维码接口)
- [七、模块间通信机制](#七模块间通信机制)
- [八、多租户实现方案](#八多租户实现方案)
- [九、构建与开发流程](#九构建与开发流程)
  - [9.1 模块启用/禁用](#91-模块启用禁用)
  - [9.2 多环境管理](#92-多环境管理)
  - [9.3 CLI 代码生成器](#93-cli-代码生成器)
  - [9.4 开发流程总览](#94-开发流程总览)
  - [9.5 Makefile 常用命令](#95-makefile-常用命令)

---

## 一、项目定位与技术选型

| 维度 | 决策 |
|------|------|
| **语言** | Go |
| **HTTP框架** | Gin |
| **架构模式** | 模块化单体（Monolithic Modular），可按需拆分微服务 |
| **API风格** | RESTful JSON |
| **配置管理** | 本地YAML/TOML兜底 + 可插拔配置中心（Nacos/Consul），支持热更新 |

### 存储与中间件

| 类型 | 支持 |
|------|------|
| **关系型数据库** | MySQL / PostgreSQL（读写分离、连接池、分库分表） |
| **缓存** | Redis（多级缓存 + 穿透/击穿/雪崩防护） |
| **消息队列** | RabbitMQ / Kafka（异步削峰、批量处理） |
| **搜索引擎** | Elasticsearch |
| **对象存储** | MinIO / Ceph |
| **文档数据库** | MongoDB |

### 核心业务模块

- **支付通道**：微信支付（全场景）、支付宝、银联、云闪付、数字人民币 + Stripe、PayPal
- **用户与权限**：基础认证 + 第三方OAuth + RBAC/ABAC + 多租户SaaS
- **二维码**：全功能（生成+解析+短链+统计）
- **高并发**：限流防刷 + 熔断降级 + 分布式锁 + 幂等 + 异步任务队列 + 多级缓存三防 + MQ削峰 + 数据层优化

### 可观测性

结构化日志 + Prometheus/Grafana + OpenTelemetry + 异常告警 + Spark性能报告 + 简易综合报告

### 开发者体验

CLI代码生成器 + 完善文档 + 示例

### 部署

二进制 / Docker Compose / K8s(Helm) / CI-CD / Serverless + 多环境隔离（开发/测试/预发/生产）

---

## 二、完整目录结构

```
spiringo/
├── cmd/
│   └── spiringo/                    # 主服务入口
│       └── main.go                  # 启动、初始化、优雅关停
│
├── internal/
│   ├── core/                        # ===== 核心框架层 =====
│   │   ├── app/                     # 应用生命周期
│   │   │   ├── app.go               # App结构体：启动/停止/信号处理
│   │   │   └── options.go           # App配置项
│   │   ├── module/                  # 模块注册与生命周期
│   │   │   ├── module.go            # Module接口定义
│   │   │   ├── registry.go          # 模块注册中心
│   │   │   └── lifecycle.go         # 模块生命周期管理
│   │   ├── server/                  # HTTP服务引擎
│   │   │   ├── server.go            # Gin引擎封装
│   │   │   ├── router.go            # 路由注册管理
│   │   │   └── options.go           # 服务配置项
│   │   ├── middleware/              # 全局中间件
│   │   │   ├── recovery.go          # 崩溃恢复
│   │   │   ├── cors.go              # 跨域
│   │   │   ├── ratelimit/           # 限流
│   │   │   │   ├── ratelimit.go     # 限流中间件
│   │   │   │   ├── token_bucket.go  # 令牌桶
│   │   │   │   ├── sliding_window.go# 滑动窗口
│   │   │   │   └── leaky_bucket.go  # 漏桶
│   │   │   ├── circuitbreaker/      # 熔断
│   │   │   │   └── circuitbreaker.go
│   │   │   ├── idempotent.go        # 幂等性
│   │   │   ├── request_id.go        # 请求ID
│   │   │   ├── i18n.go              # 国际化
│   │   │   └── tenant.go            # 多租户上下文注入
│   │   ├── config/                  # 配置管理
│   │   │   ├── config.go            # 配置管理器
│   │   │   ├── source.go            # 配置源接口
│   │   │   ├── file.go              # 文件配置源（YAML/TOML）
│   │   │   ├── nacos.go             # Nacos配置源
│   │   │   ├── consul.go            # Consul配置源
│   │   │   ├── env.go               # 环境变量配置源
│   │   │   ├── watcher.go           # 配置热更新监听
│   │   │   └── options.go           # 配置选项
│   │   ├── event/                   # 事件总线
│   │   │   ├── bus.go               # 事件总线实现
│   │   │   ├── event.go             # 事件定义
│   │   │   └── handler.go           # 事件处理器接口
│   │   ├── di/                      # 依赖注入
│   │   │   ├── container.go         # IoC容器
│   │   │   ├── provider.go          # 服务提供者接口
│   │   │   └── inject.go            # 注入辅助
│   │   └── migrate/                 # 数据库迁移
│   │       ├── migrate.go           # 迁移管理器
│   │       └── migrator.go          # 迁移接口
│   │
│   ├── pkg/                         # ===== 公共工具包 =====
│   │   ├── orm/                     # 数据库抽象
│   │   │   ├── orm.go               # ORM接口与工厂
│   │   │   ├── mysql.go             # MySQL适配
│   │   │   ├── postgres.go          # PostgreSQL适配
│   │   │   ├── mongo.go             # MongoDB适配
│   │   │   ├── readwrite.go         # 读写分离
│   │   │   ├── sharding.go          # 分库分表
│   │   │   └── options.go
│   │   ├── cache/                   # 缓存抽象
│   │   │   ├── cache.go             # 缓存接口
│   │   │   ├── memory.go            # 内存缓存
│   │   │   ├── redis.go             # Redis适配
│   │   │   ├── multilevel.go        # 多级缓存
│   │   │   ├── protector/           # 缓存防护
│   │   │   │   ├── penetration.go   # 防穿透（布隆过滤器/空值缓存）
│   │   │   │   ├── breakdown.go     # 防击穿（互斥锁/逻辑过期）
│   │   │   │   └── avalanche.go     # 防雪崩（随机TTL/预热）
│   │   │   └── options.go
│   │   ├── lock/                    # 分布式锁
│   │   │   ├── lock.go              # 锁接口
│   │   │   ├── redis.go             # Redis锁实现
│   │   │   ├── zookeeper.go         # ZK锁实现
│   │   │   └── options.go
│   │   ├── mq/                      # 消息队列抽象
│   │   │   ├── mq.go                # MQ接口
│   │   │   ├── rabbitmq.go          # RabbitMQ适配
│   │   │   ├── kafka.go             # Kafka适配
│   │   │   └── options.go
│   │   ├── search/                  # 搜索引擎抽象
│   │   │   ├── search.go            # 搜索接口
│   │   │   ├── elasticsearch.go     # ES适配
│   │   │   └── options.go
│   │   ├── oss/                     # 对象存储抽象
│   │   │   ├── oss.go               # OSS接口
│   │   │   ├── minio.go             # MinIO适配
│   │   │   ├── ceph.go              # Ceph适配
│   │   │   └── options.go
│   │   ├── logger/                  # 日志系统
│   │   │   ├── logger.go            # 日志接口
│   │   │   ├── zap.go               # Zap实现
│   │   │   └── options.go
│   │   ├── trace/                   # 链路追踪
│   │   │   ├── trace.go             # 追踪接口
│   │   │   ├── otel.go              # OpenTelemetry实现
│   │   │   └── options.go
│   │   ├── metrics/                 # 指标采集
│   │   │   ├── metrics.go           # 指标接口
│   │   │   ├── prometheus.go        # Prometheus实现
│   │   │   └── options.go
│   │   ├── alert/                   # 告警
│   │   │   ├── alert.go             # 告警接口
│   │   │   ├── sentry.go            # Sentry实现
│   │   │   ├── webhook.go           # Webhook实现
│   │   │   └── options.go
│   │   ├── crypto/                  # 加解密
│   │   │   ├── aes.go
│   │   │   ├── rsa.go
│   │   │   ├── hash.go
│   │   │   └── jwt.go
│   │   ├── queue/                   # 任务队列
│   │   │   ├── queue.go             # 队列接口
│   │   │   ├── async.go             # 异步任务
│   │   │   ├── delayed.go           # 延迟任务
│   │   │   └── options.go
│   │   ├── httpx/                   # HTTP客户端
│   │   │   ├── client.go            # 封装http.Client
│   │   │   ├── retry.go             # 重试策略
│   │   │   └── options.go
│   │   ├── validator/               # 参数校验
│   │   │   └── validator.go
│   │   ├── pagination/              # 分页
│   │   │   └── pagination.go
│   │   ├── convert/                 # 类型转换
│   │   │   └── convert.go
│   │   ├── snowflake/               # ID生成器
│   │   │   └── snowflake.go
│   │   └── utils/                   # 通用工具
│   │       ├── str.go
│   │       ├── time.go
│   │       ├── file.go
│   │       └── net.go
│   │
│   └── modules/                     # ===== 业务模块 =====
│       ├── auth/                    # 认证模块
│       │   ├── auth.go              # 模块注册入口
│       │   ├── handler/             # HTTP处理器
│       │   │   ├── login.go
│       │   │   ├── register.go
│       │   │   ├── logout.go
│       │   │   └── oauth.go
│       │   ├── service/             # 业务逻辑
│       │   │   ├── auth.go
│       │   │   ├── jwt.go
│       │   │   └── oauth/
│       │   │       ├── oauth.go     # OAuth统一接口
│       │   │       ├── wechat.go
│       │   │       ├── qq.go
│       │   │       ├── discord.go
│       │   │       ├── telegram.go
│       │   │       ├── x.go
│       │   │       ├── google.go
│       │   │       ├── dingtalk.go
│       │   │       ├── douyin.go
│       │   │       └── work_wechat.go
│       │   ├── repository/          # 数据访问
│       │   │   └── auth.go
│       │   ├── model/               # 数据模型
│       │   │   └── auth.go
│       │   ├── dto/                 # 请求/响应DTO
│       │   │   └── auth.go
│       │   └── config.go            # 模块配置
│       │
│       ├── user/                    # 用户模块
│       │   ├── user.go
│       │   ├── handler/
│       │   ├── service/
│       │   ├── repository/
│       │   ├── model/
│       │   ├── dto/
│       │   └── config.go
│       │
│       ├── tenant/                  # 多租户模块
│       │   ├── tenant.go
│       │   ├── handler/
│       │   ├── service/
│       │   ├── repository/
│       │   ├── model/
│       │   ├── dto/
│       │   └── config.go
│       │
│       ├── rbac/                    # 权限管理模块
│       │   ├── rbac.go
│       │   ├── handler/
│       │   ├── service/
│       │   ├── repository/
│       │   ├── model/
│       │   ├── dto/
│       │   └── config.go
│       │
│       ├── payment/                 # 支付模块
│       │   ├── payment.go           # 模块注册入口
│       │   ├── handler/
│       │   │   ├── pay.go           # 统一下单
│       │   │   ├── callback.go      # 支付回调
│       │   │   ├── refund.go        # 退款
│       │   │   └── query.go         # 查询
│       │   ├── service/
│       │   │   ├── payment.go       # 支付业务逻辑
│       │   │   └── channel/         # 支付通道实现
│       │   │       ├── channel.go   # 通道接口
│       │   │       ├── registry.go  # 通道注册表
│       │   │       ├── wechat/      # 微信支付
│       │   │       │   ├── wechat.go
│       │   │       │   ├── jsapi.go
│       │   │       │   ├── app.go
│       │   │       │   ├── h5.go
│       │   │       │   ├── native.go
│       │   │       │   └── mini.go
│       │   │       ├── alipay/
│       │   │       │   └── alipay.go
│       │   │       ├── unionpay/
│       │   │       │   └── unionpay.go
│       │   │       ├── cloudpay/
│       │   │       │   └── cloudpay.go
│       │   │       ├── digital_rmb/
│       │   │       │   └── digital_rmb.go
│       │   │       ├── stripe/
│       │   │       │   └── stripe.go
│       │   │       └── paypal/
│       │   │           └── paypal.go
│       │   ├── repository/
│       │   ├── model/
│       │   │   ├── order.go         # 支付订单
│       │   │   ├── refund.go        # 退款记录
│       │   │   └── callback.go      # 回调记录
│       │   ├── dto/
│       │   └── config.go
│       │
│       ├── qrcode/                  # 二维码模块
│       │   ├── qrcode.go
│       │   ├── handler/
│       │   ├── service/
│       │   │   ├── generator.go     # 二维码生成
│       │   │   ├── parser.go        # 二维码解析
│       │   │   └── shortlink.go     # 短链跳转
│       │   ├── repository/
│       │   ├── model/
│       │   ├── dto/
│       │   └── config.go
│       │
│       └── notification/            # 通知模块（短信/邮件/站内信）
│           ├── notification.go
│           ├── handler/
│           ├── service/
│           │   ├── sms.go
│           │   ├── email.go
│           │   └── inbox.go
│           ├── repository/
│           ├── model/
│           ├── dto/
│           └── config.go
│
├── pkg/                             # ===== 对外暴露公共包 =====
│   ├── types/                       # 公共类型
│   │   ├── response.go              # 统一响应结构
│   │   ├── errors.go                # 统一错误码
│   │   ├── context.go               # 上下文Key
│   │   └── pagination.go            # 分页类型
│   └── constants/                   # 公共常量
│       └── constants.go
│
├── configs/                         # ===== 配置文件 =====
│   ├── config.yaml                  # 默认配置
│   ├── config.test.yaml             # 测试环境
│   └── config.example.yaml          # 示例配置
│
├── deployments/                     # ===== 部署相关 =====
│   ├── docker/
│   │   └── Dockerfile
│   ├── docker-compose/
│   │   ├── docker-compose.yaml
│   │   └── docker-compose.dev.yaml
│   ├── k8s/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── configmap.yaml
│   │   └── ingress.yaml
│   └── helm/
│       └── spiringo/
│           ├── Chart.yaml
│           ├── values.yaml
│           └── templates/
│
├── scripts/                         # ===== 脚本 =====
│   ├── build.sh
│   ├── migrate.sh
│   └── generate.sh
│
├── docs/                            # ===== 文档 =====
│   ├── architecture.md              # 架构文档（本文件）
│   ├── getting-started.md
│   ├── modules/
│   │   ├── payment.md
│   │   ├── auth.md
│   │   └── ...
│   ├── deployment.md
│   └── api/
│       └── openapi.yaml
│
├── .github/                         # ===== CI/CD =====
│   └── workflows/
│       ├── ci.yaml
│       ├── cd.yaml
│       └── release.yaml
│
├── cli/                             # ===== CLI代码生成器 =====
│   └── spiringo/
│       ├── main.go
│       ├── cmd/
│       │   ├── new.go               # 新建项目
│       │   ├── module.go            # 新建模块
│       │   ├── crud.go              # 生成CRUD
│       │   └── migrate.go           # 生成迁移
│       └── templates/               # 代码模板
│           ├── module/
│           ├── handler/
│           ├── service/
│           ├── repository/
│           └── model/
│
├── go.mod
├── go.sum
├── Makefile
├── .gitignore
├── .golangci.yaml
└── README.md
```

---

## 三、核心框架层设计

### 3.1 模块注册与生命周期管理

**设计思路**：采用「接口驱动 + 自动注册」模式。每个模块实现 `Module` 接口，通过 `init()` 函数或显式注册将自身注册到全局 Registry。App 启动时按依赖顺序初始化模块，关停时逆序销毁。

```go
// internal/core/module/module.go

// ModuleState 模块状态
type ModuleState int

const (
    ModuleStateInactive ModuleState = iota
    ModuleStateInitializing
    ModuleStateActive
    ModuleStateStopping
    ModuleStateStopped
)

// Module 模块接口 —— 所有业务模块必须实现
type Module interface {
    // Name 模块唯一标识，如 "payment"、"auth"
    Name() string

    // Dependencies 声明依赖的模块名列表（用于确定初始化顺序）
    Dependencies() []string

    // Config 返回模块的配置结构体指针（框架会自动绑定配置）
    Config() any

    // Init 初始化模块（注册路由、绑定事件、初始化连接等）
    Init(app *App) error

    // Start 启动模块（启动goroutine、消费者等）
    Start(ctx context.Context) error

    // Stop 优雅停止模块
    Stop(ctx context.Context) error

    // State 返回当前模块状态
    State() ModuleState
}

// OptionalModule 可选接口 —— 模块可选择实现
type Routable interface {
    // Routes 注册HTTP路由
    Routes(r *gin.RouterGroup)
}

type Migratable interface {
    // Migrations 返回数据库迁移列表
    Migrations() []Migration
}

type EventSubscriber interface {
    // Subscriptions 返回事件订阅列表
    Subscriptions() []EventSubscription
}
```

```go
// internal/core/module/registry.go

type Registry struct {
    modules   map[string]Module
    initOrder []string // 拓扑排序后的初始化顺序
    mu        sync.RWMutex
}

func (r *Registry) Register(m Module) error          // 注册模块
func (r *Registry) MustRegister(m Module)             // 注册模块，重复则panic
func (r *Registry) Get(name string) (Module, error)   // 按名获取模块
func (r *Registry) ResolveOrder() error               // 拓扑排序解决依赖
func (r *Registry) InitAll(app *App) error            // 按依赖顺序初始化所有模块
func (r *Registry) StartAll(ctx context.Context) error // 启动所有模块
func (r *Registry) StopAll(ctx context.Context) error  // 逆序停止所有模块
```

**模块注册示例**：

```go
// internal/modules/payment/payment.go

func init() {
    registry.MustRegister(&PaymentModule{})
}

type PaymentModule struct {
    state   module.ModuleState
    config  PaymentConfig
    service *PaymentService
}

func (m *PaymentModule) Name() string            { return "payment" }
func (m *PaymentModule) Dependencies() []string   { return []string{"auth", "user", "tenant"} }
func (m *PaymentModule) Config() any             { return &m.config }

func (m *PaymentModule) Init(app *module.App) error {
    m.service = NewPaymentService(m.config, app.DI())
    return nil
}

func (m *PaymentModule) Routes(r *gin.RouterGroup) {
    pay := r.Group("/payment")
    pay.POST("/create", m.service.CreateOrder)
    pay.POST("/callback/:channel", m.service.HandleCallback)
    pay.POST("/refund", m.service.Refund)
}

func (m *PaymentModule) Start(ctx context.Context) error { ... }
func (m *PaymentModule) Stop(ctx context.Context) error  { ... }
```

### 3.2 依赖注入方案

**设计思路**：不引入重量级DI框架（如wire/dig），采用轻量级 Service Container + 构造函数注入。容器提供按类型/名称获取实例的能力，模块通过 `Init(app)` 获取容器引用。

```go
// internal/core/di/container.go

type Container struct {
    instances map[reflect.Type]any        // 按类型存单例
    named     map[string]any              // 按名称存（适用于多实现场景）
    factories map[reflect.Type]func() any // 延迟工厂
    mu        sync.RWMutex
}

func (c *Container) Provide(instance any)                                              // 注册单例
func (c *Container) ProvideNamed(name string, instance any)                            // 按名称注册
func (c *Container) ProvideFactory(instanceType reflect.Type, factory func() any)      // 注册延迟工厂
func Resolve[T any](c *Container) (T, error)                                           // 按类型解析
func ResolveNamed[T any](c *Container, name string) (T, error)                         // 按名称解析
func MustResolve[T any](c *Container) T                                                // 解析或panic
```

**使用方式**：

```go
// 在核心层注册基础设施
container.Provide(dbInstance)                          // 按类型自动注册 *gorm.DB
container.ProvideNamed("redis_cache", redisCache)
container.ProvideNamed("memory_cache", memoryCache)

// 在模块中解析
func (m *PaymentModule) Init(app *module.App) error {
    db := di.MustResolve[*gorm.DB](app.DI())
    cache := di.MustResolveNamed[cache.Cache](app.DI(), "redis_cache")
    m.service = NewPaymentService(db, cache)
    return nil
}
```

**理由**：Go 生态中，构造函数注入比注解注入更地道。轻量容器避免了代码生成依赖，同时提供了多实现场景下的按名称解析能力。

### 3.3 中间件链设计

**设计思路**：分层中间件——全局中间件（框架层）+ 模块级中间件（模块层），统一通过 Gin 的中间件机制注册。

```go
// internal/core/middleware/chain.go

// MiddlewareConfig 控制中间件的启用/禁用和参数
type MiddlewareConfig struct {
    Recovery     bool              `yaml:"recovery"`
    CORS         CORSConfig        `yaml:"cors"`
    RateLimit    RateLimitConfig   `yaml:"rate_limit"`
    CircuitBreak CBConfig          `yaml:"circuit_break"`
    RequestID    bool              `yaml:"request_id"`
    Tenant       bool              `yaml:"tenant"`
    Idempotent   IdempotentConfig  `yaml:"idempotent"`
    I18n         bool              `yaml:"i18n"`
    Auth         AuthConfig        `yaml:"auth"`
}

// SetupMiddlewares 按配置注册全局中间件
func SetupMiddlewares(r *gin.Engine, cfg MiddlewareConfig, container *di.Container) {
    // 注册顺序（从外到内）：
    // 1. Recovery（最外层，兜底panic）
    // 2. RequestID（为后续日志/追踪提供ID）
    // 3. CORS
    // 4. I18n
    // 5. RateLimit（在业务逻辑前限流）
    // 6. Tenant（注入租户上下文，必须在Auth前）
    // 7. Auth（认证鉴权）
    // 8. Idempotent（幂等性校验）
    // 9. CircuitBreaker（熔断保护）
}
```

**模块级中间件**：模块可以在 `Routes()` 方法中为特定路由组添加中间件：

```go
func (m *PaymentModule) Routes(r *gin.RouterGroup) {
    pay := r.Group("/payment")
    pay.Use(middleware.Idempotent())  // 支付接口强制幂等
    pay.POST("/create", m.service.CreateOrder)
}
```

框架层幂等实现将业务幂等键按 HTTP 方法、路由模板、租户、用户与 Header 名组合成作用域，避免不同接口、不同租户或不同用户复用同一业务键时互相误判为重复请求。

### 3.4 配置管理架构

**设计思路**：配置源抽象 + 优先级合并 + 热更新推送。

```go
// internal/core/config/config.go

type Manager struct {
    sources  []Source       // 按优先级排列的配置源（高→低）
    values   *viper.Viper   // 合并后的配置值
    watchers []Watcher      // 热更新监听器
    mu       sync.RWMutex
    onChange func(key string, value any) // 配置变更回调
}

func (m *Manager) Load() error                                  // 加载所有配置源并合并
func (m *Manager) Unmarshal(key string, target any) error        // 将配置绑定到结构体
func (m *Manager) Watch(key string, callback func(value any))    // 监听配置变更
func (m *Manager) Get(key string) any                            // 获取配置值
```

```go
// internal/core/config/source.go

// Source 配置源接口
type Source interface {
    Name() string                                                      // 配置源名称
    Priority() int                                                     // 优先级（越大越优先）
    Read() (map[string]any, error)                                     // 读取配置
    Watch(ctx context.Context, onChange func(key string, value any)) error // 监听配置变更
    Close() error                                                      // 关闭
}
```

**配置文件结构** (`configs/config.yaml`)：

```yaml
app:
  name: spiringo
  env: development   # development | test | staging | production
  debug: true

server:
  addr: ":8080"
  mode: "debug"      # debug | release | test

config_center:        # 配置中心（可选）
  enabled: false
  type: "nacos"       # nacos | consul
  nacos:
    server_addr: "127.0.0.1:8848"
    namespace: "public"
    group: "DEFAULT_GROUP"
    data_id: "spiringo"

database:
  default:
    driver: "mysql"
    dsn: "root:password@tcp(127.0.0.1:3306)/spiringo?charset=utf8mb4&parseTime=True"
    max_idle: 10
    max_open: 100
    read_replicas:
      - dsn: "root:password@tcp(127.0.0.1:3307)/spiringo?charset=utf8mb4&parseTime=True"

redis:
  default:
    addr: "127.0.0.1:6379"
    password: ""
    db: 0

modules:
  auth:
    enabled: true
    jwt:
      secret: "your-secret-key"
      expire: "24h"
  payment:
    enabled: true
    wechat:
      app_id: ""
      mch_id: ""
      api_key: ""
  tenant:
    enabled: true
    strategy: "shared_db"  # shared_db | schema | database

middleware:
  rate_limit:
    enabled: true
    strategy: "sliding_window"
    rate: 100
    burst: 200
  circuit_break:
    enabled: true
    threshold: 0.5
```

---

## 四、业务模块详细内部结构

### 4.1 支付模块

```
payment/
├── payment.go              # 模块注册入口，实现 Module 接口
├── config.go               # PaymentConfig 结构体
├── handler/                # HTTP层：只做参数绑定/校验/响应
│   ├── pay.go              #   POST /payment/create     → 创建支付单
│   ├── callback.go         #   POST /payment/callback/:channel → 支付回调
│   ├── refund.go           #   POST /payment/refund     → 申请退款
│   └── query.go            #   GET  /payment/query/:id  → 查询订单
├── service/                # 业务层：编排业务逻辑
│   ├── payment.go          #   PaymentService：创建订单/退款/查询
│   └── channel/            #   通道层：各支付通道实现
│       ├── channel.go      #     Channel 接口定义
│       ├── registry.go     #     通道注册表
│       ├── wechat/         #     微信支付
│       │   ├── wechat.go   #     WechatChannel 实现 Channel 接口
│       │   ├── jsapi.go    #     JSAPI支付场景
│       │   ├── app.go      #     APP支付场景
│       │   ├── h5.go       #     H5支付场景
│       │   ├── native.go   #     Native扫码支付
│       │   └── mini.go     #     小程序支付
│       ├── alipay/
│       │   └── alipay.go
│       ├── unionpay/
│       │   └── unionpay.go
│       ├── cloudpay/
│       │   └── cloudpay.go
│       ├── digital_rmb/
│       │   └── digital_rmb.go
│       ├── stripe/
│       │   └── stripe.go
│       └── paypal/
│           └── paypal.go
├── repository/             # 数据层：只做CRUD
│   ├── order.go            #   支付订单 CRUD
│   ├── refund.go           #   退款记录 CRUD
│   └── callback_log.go     #   回调日志 CRUD
├── model/                  # 数据模型：与数据库表对应
│   ├── order.go            #   PaymentOrder
│   ├── refund.go           #   RefundRecord
│   └── callback_log.go     #   CallbackLog
└── dto/                    # 请求/响应传输对象
    ├── request.go          #   CreatePayReq / RefundReq / QueryReq
    └── response.go         #   CreatePayResp / PayOrderResp / RefundResp
```

**支付核心流程**：

```
用户请求 → Handler(参数校验) → Service(业务编排)
  → 1. 创建本地订单（repository）
  → 2. 幂等性检查（分布式锁）
  → 3. 选择通道（channel registry）
  → 4. 调用通道创建预支付（channel.CreatePayment）
  → 5. 发布"订单创建"事件（event bus）
  → 6. 返回支付参数

回调流程：
第三方回调 → Handler(验签) → Service(回调处理)
  → 1. 验证签名（channel.VerifyCallback）
  → 2. 幂等处理（检查订单状态）
  → 3. 更新订单状态
  → 4. 发布"支付成功"事件
  → 5. 通知业务模块（如订单模块）

退款流程：
退款请求 → Handler(参数校验) → Service(退款编排)
  → 1. 按 out_refund_no 做业务幂等，重复请求返回已有退款单
  → 2. 校验原订单状态、总金额一致性与累计可退金额
  → 3. 创建本地退款单并调用通道退款
  → 4. 更新退款单与原订单状态
  → 5. 退款成功时发布"退款完成"事件
```

### 4.2 认证模块

```
auth/
├── auth.go                 # 模块注册入口
├── config.go               # AuthConfig（JWT配置、OAuth配置等）
├── handler/
│   ├── login.go            #   POST /auth/login          → 账密登录
│   ├── register.go         #   POST /auth/register       → 注册
│   ├── logout.go           #   POST /auth/logout         → 登出
│   ├── refresh.go          #   POST /auth/refresh        → 刷新Token
│   ├── oauth.go            #   GET  /auth/oauth/:provider → 第三方登录跳转
│   └── oauth_callback.go   #   GET  /auth/oauth/:provider/callback → 回调
├── service/
│   ├── auth.go             #   AuthService：登录/注册/登出逻辑
│   ├── jwt.go              #   JWTService：生成/验证/刷新Token
│   └── oauth/
│       ├── oauth.go        #     OAuthProvider 接口
│       ├── registry.go     #     Provider注册表
│       ├── wechat.go       #     微信OAuth
│       ├── qq.go           #     QQ OAuth
│       ├── discord.go      #     Discord OAuth
│       ├── telegram.go     #     Telegram OAuth
│       ├── x.go            #     X(Twitter) OAuth
│       ├── google.go       #     Google OAuth
│       ├── dingtalk.go     #     钉钉 OAuth
│       ├── douyin.go       #     抖音 OAuth
│       └── work_wechat.go  #     企业微信 OAuth
├── repository/
│   └── auth.go             #   用户认证信息CRUD
├── model/
│   └── auth.go             #   AuthIdentity / OAuthBinding
└── dto/
    ├── request.go          #   LoginReq / RegisterReq / OAuthReq
    └── response.go         #   LoginResp（含AccessToken/RefreshToken）
```

**OAuth 统一流程**：

```
1. GET /auth/oauth/:provider  → 302重定向到第三方授权页
2. 用户授权后回调 /auth/oauth/:provider/callback?code=xxx
3. Service 用 code 换取 access_token
4. Service 用 access_token 获取用户信息
5. 查找或创建本地用户（关联 OAuthBinding）
6. 生成 JWT Token 返回
```

---

## 五、中间件抽象层设计

**设计原则**：每个中间件类型定义统一接口 → 提供默认实现 → 通过配置切换实现 → 通过 DI 容器注入到需要的地方。

### 统一抽象模式

所有抽象遵循同一模式：

```go
// 1. 定义接口（接口只关注"做什么"，不关注"怎么做"）
type XXX interface {
    Method1(...)
    Method2(...)
}

// 2. 提供默认实现
type redisXXX struct { ... }
func NewRedisXXX(cfg XXXConfig) XXX { ... }

type memoryXXX struct { ... }
func NewMemoryXXX(cfg XXXConfig) XXX { ... }

// 3. 通过工厂函数根据配置选择实现
func NewXXX(cfg XXXConfig) XXX {
    switch cfg.Driver {
    case "redis":    return NewRedisXXX(cfg)
    case "memory":   return NewMemoryXXX(cfg)
    case "kafka":    return NewKafkaXXX(cfg)
    default:         return NewDefaultXXX(cfg)
    }
}

// 4. 注册到DI容器，模块按需解析
container.Provide(NewXXX(cfg))
```

### 数据库抽象（ORM层）

```go
// internal/pkg/orm/orm.go

type DB interface {
    // 基础CRUD
    Create(ctx context.Context, model any) error
    Update(ctx context.Context, model any) error
    Delete(ctx context.Context, model any, conds ...any) error
    First(ctx context.Context, model any, conds ...any) error
    Find(ctx context.Context, models any, conds ...any) error
    Count(ctx context.Context, model any, conds ...any) (int64, error)

    // 条件构建
    Where(query any, args ...any) DB
    Order(value string) DB
    Limit(limit int) DB
    Offset(offset int) DB

    // 事务
    Transaction(ctx context.Context, fn func(tx DB) error) error

    // 原始查询（复杂场景兜底）
    Raw(sql string, values ...any) DB
    Exec(sql string, values ...any) error

    // 读写分离：强制走主库
    Master() DB
}
```

**理由**：不直接暴露 `*gorm.DB`，而是包装一层。好处：1. 读写分离对调用者透明；2. 未来可替换 ORM 实现；3. 自动注入租户ID过滤条件。

### 缓存抽象

```go
// internal/pkg/cache/cache.go

type Cache interface {
    Get(ctx context.Context, key string, dest any) error
    Set(ctx context.Context, key string, value any, ttl time.Duration) error
    Delete(ctx context.Context, keys ...string) error
    Exists(ctx context.Context, key string) (bool, error)

    // 批量操作
    MGet(ctx context.Context, keys []string, dest any) error
    MSet(ctx context.Context, kv map[string]any, ttl time.Duration) error

    // 过期/自增
    Expire(ctx context.Context, key string, ttl time.Duration) error
    Incr(ctx context.Context, key string) (int64, error)
    IncrBy(ctx context.Context, key string, delta int64) (int64, error)

    // 关闭
    Close() error
}

// MultiLevelCache 多级缓存（L1内存 + L2 Redis）
type MultiLevelCache struct {
    l1 Cache  // 进程内缓存（BigCache/FreeCache）
    l2 Cache  // 分布式缓存（Redis）
    protector.Protector  // 穿透/击穿/雪崩防护
}
```

### 消息队列抽象

```go
// internal/pkg/mq/mq.go

type Message struct {
    Topic   string
    Key     string
    Value   []byte
    Headers map[string]string
}

type MQ interface {
    Publish(ctx context.Context, msg *Message) error                       // 发布
    PublishBatch(ctx context.Context, msgs []*Message) error               // 批量发布
    Subscribe(ctx context.Context, topic string, handler func(msg *Message) error) error // 订阅
    Close() error                                                          // 关闭
}
```

### 搜索引擎抽象

```go
// internal/pkg/search/search.go

type SearchQuery struct {
    Index   string
    Keyword string
    Fields  []string
    Filter  map[string]any
    Sort    []SortField
    From    int
    Size    int
}

type SearchResult struct {
    Total int64
    Hits  []map[string]any
}

type Search interface {
    Index(ctx context.Context, index string, id string, doc any) error
    BulkIndex(ctx context.Context, index string, docs []any) error
    Delete(ctx context.Context, index string, id string) error
    Search(ctx context.Context, query *SearchQuery) (*SearchResult, error)
    Aggregation(ctx context.Context, query *SearchQuery, aggFields []string) (map[string]any, error)
}
```

### 对象存储抽象

```go
// internal/pkg/oss/oss.go

type OSS interface {
    Put(ctx context.Context, key string, reader io.Reader, size int64) error
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
    PresignedURL(ctx context.Context, key string, expire time.Duration) (string, error)
    Exists(ctx context.Context, key string) (bool, error)
    List(ctx context.Context, prefix string, limit int) ([]string, error)
}
```

### 分布式锁抽象

```go
// internal/pkg/lock/lock.go

type Lock interface {
    Lock(ctx context.Context, key string, ttl time.Duration) (LockHolder, error)
    TryLock(ctx context.Context, key string, ttl time.Duration) (LockHolder, error)
}

type LockHolder interface {
    Unlock(ctx context.Context) error
    Renew(ctx context.Context, ttl time.Duration) error  // 续期
}
```

---

## 六、关键接口定义

### 6.1 支付统一接口

```go
// internal/modules/payment/service/channel/channel.go

type PayScene string

const (
    PaySceneJSAPI  PayScene = "jsapi"   // 微信公众号/小程序
    PaySceneAPP    PayScene = "app"     // APP支付
    PaySceneH5     PayScene = "h5"      // H5网页支付
    PaySceneNative PayScene = "native"  // 扫码支付（商家码）
    PaySceneQrCode PayScene = "qrcode"  // 扫码支付（用户扫）
)

type PayStatus string

const (
    PayStatusPending   PayStatus = "pending"
    PayStatusPaid      PayStatus = "paid"
    PayStatusFailed    PayStatus = "failed"
    PayStatusClosed    PayStatus = "closed"
    PayStatusRefunding PayStatus = "refunding"
    PayStatusRefunded  PayStatus = "refunded"
)

// CreatePaymentRequest 统一创建支付请求
type CreatePaymentRequest struct {
    OutTradeNo  string            // 商户订单号（幂等键）
    Amount      int64             // 金额（单位：分）
    Currency    string            // 币种：CNY/USD/EUR
    Subject     string            // 订单标题
    Description string            // 订单描述
    Scene       PayScene          // 支付场景
    ChannelCode string            // 通道编码：wechat/alipay/unionpay/stripe/paypal
    NotifyURL   string            // 回调地址
    ReturnURL   string            // 前端跳转地址
    Metadata    map[string]string // 扩展参数（如 openid 等）
    TenantID    string            // 租户ID
}

// CreatePaymentResponse 统一创建支付响应
type CreatePaymentResponse struct {
    TradeNo     string            // 平台流水号
    OutTradeNo  string            // 商户订单号
    PayURL      string            // 支付跳转URL（H5/微信跳转）
    QrCode      string            // 二维码内容（Native扫码）
    PrepayID    string            // 预支付ID
    PayParams   map[string]string // 客户端调起支付所需参数（如微信JSAPI参数）
    ExpireAt    time.Time         // 过期时间
}

// RefundRequest 统一退款请求
type RefundRequest struct {
    OutTradeNo   string // 原商户订单号
    OutRefundNo  string // 退款单号
    TotalAmount  int64  // 原订单金额（分）
    RefundAmount int64  // 退款金额（分）
    Reason       string // 退款原因
    NotifyURL    string // 退款回调地址
}

// RefundResponse 统一退款响应
type RefundResponse struct {
    RefundNo    string    // 退款流水号
    OutRefundNo string    // 商户退款单号
    Status      PayStatus // 退款状态
    RefundAt    time.Time // 退款时间
}

// CallbackResult 统一回调解析结果
type CallbackResult struct {
    OutTradeNo string            // 商户订单号
    TradeNo    string            // 平台流水号
    Status     PayStatus         // 支付状态
    Amount     int64             // 实付金额
    PaidAt     time.Time         // 支付时间
    RawData    map[string]any    // 原始回调数据
    Metadata   map[string]string // 扩展信息
}

// Channel 支付通道接口 —— 每个支付通道必须实现
type Channel interface {
    Code() string   // 通道编码（如 "wechat"、"alipay"）
    Name() string   // 通道名称（如 "微信支付"）

    CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error)
    QueryPayment(ctx context.Context, outTradeNo string) (*CallbackResult, error)
    ClosePayment(ctx context.Context, outTradeNo string) error
    Refund(ctx context.Context, req *RefundRequest) (*RefundResponse, error)
    VerifyCallback(ctx context.Context, r *http.Request) (*CallbackResult, error)
    CallbackSuccess() any   // 回调成功响应（各通道格式不同）
    CallbackFail() any      // 回调失败响应
}

// ChannelRegistry 通道注册表
type ChannelRegistry struct {
    channels map[string]Channel
}

func (r *ChannelRegistry) Register(ch Channel)                // 注册通道
func (r *ChannelRegistry) Get(code string) (Channel, error)   // 获取通道
func (r *ChannelRegistry) List() []Channel                    // 列出所有通道
```

### 6.2 OAuth 统一接口

```go
// internal/modules/auth/service/oauth/oauth.go

type OAuthUser struct {
    Provider   string         // 提供商标识
    ProviderID string         // 第三方用户ID
    Username   string         // 用户名
    Nickname   string         // 昵称
    Avatar     string         // 头像URL
    Email      string         // 邮箱
    Phone      string         // 手机号
    RawData    map[string]any // 原始数据
}

type OAuthProvider interface {
    Name() string  // 提供者名称（如 "wechat"、"google"）
    AuthURL(ctx context.Context, state string, redirectURL string) (string, error)  // 获取授权跳转URL
    GetUser(ctx context.Context, code string, redirectURL string) (*OAuthUser, error) // 用授权码换取用户信息
    RefreshToken(ctx context.Context, refreshToken string) error  // 刷新access_token（可选实现）
}
```

### 6.3 二维码接口

```go
// internal/modules/qrcode/service/qrcode.go

type QRCodeStyle struct {
    Size            int    // 尺寸（像素）
    ForegroundColor string // 前景色
    BackgroundColor string // 背景色
    LogoURL         string // Logo图片URL
    LogoSize        int    // Logo尺寸
    Level           string // 容错级别：L/M/Q/H
}

type QRCodeResult struct {
    Content    []byte    // 二维码内容
    ImageURL   string    // 图片URL（如上传到OSS）
    ImageBytes []byte    // 图片字节（可选）
    ExpireAt   time.Time // 过期时间
}

type ShortLinkResult struct {
    ShortCode string    // 短码
    ShortURL  string    // 短链URL
    TargetURL string    // 目标URL
    ScanCount int64     // 扫码次数
    ExpireAt  time.Time
}

type Generator interface {
    Generate(ctx context.Context, content string, style *QRCodeStyle) (*QRCodeResult, error)
    Parse(ctx context.Context, imageData []byte) (string, error)
}

type ShortLinkService interface {
    Create(ctx context.Context, targetURL string, expire time.Duration) (*ShortLinkResult, error)
    Resolve(ctx context.Context, shortCode string) (string, error) // 短码→长链
    Stats(ctx context.Context, shortCode string) (*ShortLinkResult, error)
}
```

---

## 七、模块间通信机制

**设计决策：事件总线 + 接口依赖，双轨并行**

**理由**：
- **同步调用**用接口依赖（DI注入）—— 清晰、类型安全、可追踪
- **异步解耦**用事件总线 —— 松耦合、可扩展、适合跨模块通知

### 事件总线

```go
// internal/core/event/bus.go

type Event struct {
    Topic     string            // 事件主题
    Payload   any               // 事件数据
    Timestamp time.Time         // 事件时间
    Source    string            // 产生事件的模块名
    Metadata  map[string]string // 扩展元数据
}

type Handler func(ctx context.Context, event *Event) error

type Subscription struct {
    Topic   string
    Handler Handler
    Module  string  // 订阅者模块名（用于生命周期管理）
}

type EventBus interface {
    Publish(ctx context.Context, event *Event) error       // 发布事件（同步）
    PublishAsync(ctx context.Context, event *Event) error  // 异步发布（通过任务队列）
    Subscribe(topic string, handler Handler) error         // 订阅事件
    Unsubscribe(topic string, handler Handler) error       // 取消订阅
}
```

**内置事件主题规范**（采用 `module.action` 命名）：

```go
const (
    // 用户模块
    EventUserCreated    = "user.created"
    EventUserUpdated    = "user.updated"
    EventUserDeleted    = "user.deleted"

    // 认证模块
    EventAuthLogin      = "auth.login"
    EventAuthLogout     = "auth.logout"
    EventAuthOAuthBound = "auth.oauth_bound"

    // 支付模块
    EventPaymentCreated  = "payment.created"
    EventPaymentSuccess  = "payment.success"
    EventPaymentFailed   = "payment.failed"
    EventPaymentRefunded = "payment.refunded"
    EventPaymentClosed   = "payment.closed"

    // 租户模块
    EventTenantCreated   = "tenant.created"
    EventTenantSuspended = "tenant.suspended"
)
```

### 同步调用：接口依赖

模块间需要同步调用时，通过 DI 容器注入依赖模块的 Service 接口：

```go
type UserService interface {
    GetUserByID(ctx context.Context, id string) (*User, error)
}

type PaymentService struct {
    userSvc UserService  // 依赖用户模块的接口，而非具体实现
}

func (m *PaymentModule) Init(app *module.App) error {
    userSvc := di.MustResolve[UserService](app.DI())
    m.service = NewPaymentService(userSvc, ...)
    return nil
}
```

**关键原则**：
- **依赖接口，不依赖实现**：模块只引用其他模块的 Service 接口定义
- **单向依赖**：避免循环依赖，通过事件总线解决反向通知需求
- **模块必须在 `Dependencies()` 中声明依赖**：确保初始化顺序正确

### 通信方式选择指南

| 场景 | 方式 | 示例 |
|------|------|------|
| 模块A需要模块B的数据 | 接口依赖 | 支付模块查询用户信息 |
| 模块A完成操作需通知模块B | 事件总线 | 支付成功→通知订单模块 |
| 需要异步处理/削峰 | 异步事件 | 支付回调→异步更新订单 |
| 跨模块数据聚合 | 接口依赖 | 列表页同时展示用户+订单信息 |

---

## 八、多租户实现方案

**推荐方案：共享数据库 + 租户ID（默认）+ 可升级到独立Schema**

**理由**：
- **共享数据库 + tenant_id**：运维成本最低，适合大多数SaaS场景，代码侵入最小
- 提供平滑升级路径：租户规模增长后可迁移到独立Schema
- 通过ORM层自动注入 `tenant_id`，业务代码几乎无感知

### 租户隔离策略

```go
// internal/modules/tenant/config.go

type TenantStrategy string

const (
    StrategySharedDB TenantStrategy = "shared_db"  // 共享数据库+tenant_id
    StrategySchema   TenantStrategy = "schema"      // 独立Schema（PostgreSQL适用）
    StrategyDatabase TenantStrategy = "database"    // 独立数据库
)
```

### 租户上下文

```go
// internal/core/middleware/tenant.go

type TenantContext struct {
    TenantID   string
    TenantName string
    Strategy   TenantStrategy
    DBConn     string  // 独立数据库时的连接标识
}

func WithTenant(ctx context.Context, tc *TenantContext) context.Context {
    return context.WithValue(ctx, tenantKey{}, tc)
}

func FromContext(ctx context.Context) *TenantContext {
    tc, _ := ctx.Value(tenantKey{}).(*TenantContext)
    return tc
}

// TenantMiddleware 租户中间件
// 从请求头/域名/subdomain解析租户标识
func TenantMiddleware(tenantSvc TenantService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 策略1：从Header解析 X-Tenant-ID
        tenantID := c.GetHeader("X-Tenant-ID")

        // 策略2：从子域名解析 {tenant}.example.com
        if tenantID == "" {
            tenantID = extractFromSubdomain(c.Request.Host)
        }

        // 策略3：从JWT Token解析
        if tenantID == "" {
            tenantID = extractFromToken(c)
        }

        if tenantID == "" {
            c.Next()
            return
        }

        tenant, err := tenantSvc.GetByID(c.Request.Context(), tenantID)
        if err != nil {
            response.Fail(c, errors.ErrTenantNotFound)
            c.Abort()
            return
        }

        ctx := WithTenant(c.Request.Context(), &TenantContext{
            TenantID:   tenant.ID,
            TenantName: tenant.Name,
            Strategy:   tenant.Strategy,
        })
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
```

### ORM层自动注入 tenant_id

```go
// internal/pkg/orm/tenant.go

type TenantDB struct {
    DB DB
}

func (tdb *TenantDB) Find(ctx context.Context, models any, conds ...any) error {
    tc := FromContext(ctx)
    if tc != nil {
        tdb.DB = tdb.DB.Where("tenant_id = ?", tc.TenantID)
    }
    return tdb.DB.Find(ctx, models, conds...)
}

func (tdb *TenantDB) Create(ctx context.Context, model any) error {
    tc := FromContext(ctx)
    if tc != nil {
        setTenantID(model, tc.TenantID)
    }
    return tdb.DB.Create(ctx, model)
}
```

### 独立Schema/数据库支持

```go
// internal/pkg/orm/multitenant.go

type MultiTenantDB struct {
    defaultDB DB            // 默认数据库（共享模式用）
    tenantDBs map[string]DB // 租户→独立数据库连接（schema/database模式）
    strategy  TenantStrategy
    mu        sync.RWMutex
}

func (m *MultiTenantDB) GetDB(ctx context.Context) DB {
    tc := FromContext(ctx)
    if tc == nil || m.strategy == StrategySharedDB {
        return m.defaultDB
    }

    m.mu.RLock()
    db, ok := m.tenantDBs[tc.TenantID]
    m.mu.RUnlock()

    if ok {
        return db
    }

    return m.createTenantDB(tc) // 懒加载
}
```

### 数据模型约定

所有支持多租户的模型必须嵌入 `TenantModel`：

```go
// internal/pkg/orm/model.go

type TenantModel struct {
    TenantID string `gorm:"index;not null" json:"tenant_id"`
}

type BaseModel struct {
    ID        string    `gorm:"primaryKey;size:36" json:"id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type TenantBaseModel struct {
    BaseModel
    TenantModel
}
```

---

## 九、构建与开发流程

### 9.1 模块启用/禁用

**通过配置文件控制**，无需修改代码：

```yaml
# configs/config.yaml
modules:
  auth:
    enabled: true       # 设为 false 即禁用整个认证模块
    jwt:
      secret: "xxx"
  payment:
    enabled: true
    wechat:
      enabled: true     # 细粒度：可单独禁用某支付通道
    alipay:
      enabled: true
    stripe:
      enabled: false    # 暂不需要国际支付
  qrcode:
    enabled: true
  tenant:
    enabled: true
  notification:
    enabled: false      # 暂不需要通知模块
```

**框架层实现**：

```go
// cmd/spiringo/main.go

func main() {
    app := core.NewApp(
        core.WithConfigFile("configs/config.yaml"),
    )

    // 注册所有模块（无论是否启用都注册，框架根据配置决定是否初始化）
    module.RegisterModules(
        &auth.AuthModule{},
        &user.UserModule{},
        &tenant.TenantModule{},
        &rbac.RBACModule{},
        &payment.PaymentModule{},
        &qrcode.QRCodeModule{},
        &notification.NotificationModule{},
    )

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

```go
// internal/core/module/registry.go

// InitAll 只初始化 enabled=true 的模块
func (r *Registry) InitAll(app *App) error {
    for _, name := range r.initOrder {
        m := r.modules[name]

        // 检查模块是否启用
        if !app.Config().GetBool("modules." + name + ".enabled") {
            app.Logger().Info("module disabled, skipping", "module", name)
            continue
        }

        // 绑定模块专属配置
        if cfg := m.Config(); cfg != nil {
            if err := app.Config().Unmarshal("modules."+name, cfg); err != nil {
                return fmt.Errorf("bind config for module %s: %w", name, err)
            }
        }

        // 初始化模块
        if err := m.Init(app); err != nil {
            return fmt.Errorf("init module %s: %w", name, err)
        }

        // 注册路由（如果模块实现了 Routable）
        if routable, ok := m.(Routable); ok {
            group := app.Router().Group("/api/v1/" + name)
            routable.Routes(group)
        }

        // 注册事件订阅
        if subscriber, ok := m.(EventSubscriber); ok {
            for _, sub := range subscriber.Subscriptions() {
                app.EventBus().Subscribe(sub.Topic, sub.Handler)
            }
        }
    }
    return nil
}
```

### 9.2 多环境管理

配置文件按环境分离：

```
configs/config.yaml           — 默认/开发环境
configs/config.test.yaml      — 测试环境覆盖
configs/config.staging.yaml   — 预发环境覆盖
configs/config.prod.yaml      — 生产环境覆盖
```

**环境切换方式**：

```bash
# 方式1：环境变量
export APP_ENV=production
./spiringo

# 方式2：命令行参数
./spiringo -env production

# 方式3：配置中心（自动按环境拉取，Nacos namespace = 环境标识）
```

**配置合并优先级**（从高到低）：

```
环境变量 > 配置中心 > 环境配置文件(config.{env}.yaml) > 默认配置文件(config.yaml) > 代码默认值
```

### 9.3 CLI 代码生成器

```bash
# 安装CLI
go install github.com/spiringo/cli/spiringo@latest

# 新建项目（基于Spiringo模板）
spiringo new my-project

# 新建业务模块
spiringo module order
# → 生成 internal/modules/order/ 完整目录结构
# → 自动在 main.go 中注册模块

# 生成CRUD代码
spiringo crud order Product
# → 生成 handler/service/repository/model/dto/migration

# 生成支付通道
spiringo payment-channel bank
# → 生成 Channel 接口实现骨架

# 生成OAuth Provider
spiringo oauth-provider github
# → 生成 OAuthProvider 接口实现骨架

# 生成数据库迁移
spiringo migrate create add_product_table
```

### 9.4 开发流程总览

```
1. spiringo new my-project        → 初始化项目
2. 修改 config.yaml               → 配置数据库/Redis/支付通道等
3. spiringo module order          → 创建业务模块
4. spiringo crud order Product    → 生成CRUD代码
5. 编辑 service/product.go        → 编写业务逻辑
6. go run cmd/spiringo/main.go    → 本地启动调试
7. spiringo migrate create xxx    → 创建迁移
8. go test ./...                  → 运行测试
9. docker-compose up -d           → 本地全栈启动
10. helm install spiringo ...     → K8s部署
```

### 9.5 Makefile 常用命令

```makefile
.PHONY: all build run test lint migrate generate docker

build:
	go build -o bin/spiringo cmd/spiringo/main.go

run:
	go run cmd/spiringo/main.go -env development

test:
	go test -v -cover ./...

lint:
	golangci-lint run ./...

migrate:
	go run cmd/spiringo/main.go migrate up

generate:
	spiringo module $(MODULE) || spiringo crud $(MODULE) $(MODEL)

docker:
	docker build -t spiringo:latest -f deployments/docker/Dockerfile .

compose-up:
	docker-compose -f deployments/docker-compose/docker-compose.yaml up -d

helm-deploy:
	helm upgrade --install spiringo deployments/helm/spiringo \
		--namespace spiringo --create-namespace \
		-f deployments/helm/spiringo/values-$(ENV).yaml
```

---

## 架构核心决策一览

| 设计维度 | 决策 | 理由 |
|---------|------|------|
| **整体架构** | 模块化单体 | 单一代码库易开发，模块化设计可按需拆分 |
| **模块注册** | 接口驱动 + init()自动注册 | 开发者只需实现Module接口并注册，零配置 |
| **依赖注入** | 轻量DI容器 + 泛型Resolve | 无代码生成依赖，Go地道风格 |
| **模块通信** | 同步用接口依赖 + 异步用事件总线 | 解耦与类型安全兼顾 |
| **配置管理** | 本地文件兜底 + 可插拔配置中心 | 开发时简单，生产时可接入Nacos/Consul |
| **中间件抽象** | 统一接口 + 工厂模式 + DI注入 | 实现可替换，业务代码零感知 |
| **多租户** | 共享DB+tenant_id（默认）+ 可升级 | 运维成本低，有升级路径 |
| **支付通道** | Channel接口 + Registry | 新增通道只需实现接口并注册 |
| **模块启禁** | YAML配置控制 | 零代码修改，配置即生效 |
