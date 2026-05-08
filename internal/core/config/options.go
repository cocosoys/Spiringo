package config

// 中文：AppConfig 定义当前包使用的数据结构或接口。
// English: AppConfig defines a data structure or interface used by this package.
// AppConfig contains application-level settings.
type AppConfig struct {
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `yaml:"name" mapstructure:"name"`
	// 中文：Env 保存当前结构中的配置或数据值。
	// English: Env stores a configuration or data value for this struct.
	Env string `yaml:"env" mapstructure:"env"`
	// 中文：Debug 保存当前结构中的配置或数据值。
	// English: Debug stores a configuration or data value for this struct.
	Debug bool `yaml:"debug" mapstructure:"debug"`
}

// 中文：ServerConfig 定义当前包使用的数据结构或接口。
// English: ServerConfig defines a data structure or interface used by this package.
// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	// 中文：Addr 保存当前结构中的配置或数据值。
	// English: Addr stores a configuration or data value for this struct.
	Addr string `yaml:"addr" mapstructure:"addr"`
	// 中文：Mode 保存当前结构中的配置或数据值。
	// English: Mode stores a configuration or data value for this struct.
	Mode string `yaml:"mode" mapstructure:"mode"`
}

// 中文：LogConfig 定义当前包使用的数据结构或接口。
// English: LogConfig defines a data structure or interface used by this package.
// LogConfig contains structured logging settings.
type LogConfig struct {
	// 中文：Driver 保存当前结构中的配置或数据值。
	// English: Driver stores a configuration or data value for this struct.
	Driver string `yaml:"driver" mapstructure:"driver"` // slog or zap
	// 中文：Level 保存当前结构中的配置或数据值。
	// English: Level stores a configuration or data value for this struct.
	Level string `yaml:"level" mapstructure:"level"`
	// 中文：Format 保存当前结构中的配置或数据值。
	// English: Format stores a configuration or data value for this struct.
	Format string `yaml:"format" mapstructure:"format"`
	// 中文：Output 保存当前结构中的配置或数据值。
	// English: Output stores a configuration or data value for this struct.
	Output string `yaml:"output" mapstructure:"output"`
}

// 中文：DatabaseConfig 定义当前包使用的数据结构或接口。
// English: DatabaseConfig defines a data structure or interface used by this package.
// DatabaseConfig contains the default database connection settings.
type DatabaseConfig struct {
	// 中文：Driver 保存当前结构中的配置或数据值。
	// English: Driver stores a configuration or data value for this struct.
	Driver string `yaml:"driver" mapstructure:"driver"`
	// 中文：DSN 保存当前结构中的配置或数据值。
	// English: DSN stores a configuration or data value for this struct.
	DSN string `yaml:"dsn" mapstructure:"dsn"`
	// 中文：MaxIdle 保存当前结构中的配置或数据值。
	// English: MaxIdle stores a configuration or data value for this struct.
	MaxIdle int `yaml:"max_idle" mapstructure:"max_idle"`
	// 中文：MaxOpen 保存当前结构中的配置或数据值。
	// English: MaxOpen stores a configuration or data value for this struct.
	MaxOpen int `yaml:"max_open" mapstructure:"max_open"`
	// 中文：ConnMaxLifetime 保存当前结构中的配置或数据值。
	// English: ConnMaxLifetime stores a configuration or data value for this struct.
	ConnMaxLifetime string `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	// 中文：ReadReplicas 保存当前结构中的配置或数据值。
	// English: ReadReplicas stores a configuration or data value for this struct.
	ReadReplicas []DBConfig `yaml:"read_replicas" mapstructure:"read_replicas"`
}

// 中文：DBConfig 定义当前包使用的数据结构或接口。
// English: DBConfig defines a data structure or interface used by this package.
// DBConfig contains a simple database endpoint.
type DBConfig struct {
	// 中文：Driver 保存当前结构中的配置或数据值。
	// English: Driver stores a configuration or data value for this struct.
	Driver string `yaml:"driver" mapstructure:"driver"`
	// 中文：DSN 保存当前结构中的配置或数据值。
	// English: DSN stores a configuration or data value for this struct.
	DSN string `yaml:"dsn" mapstructure:"dsn"`
}

// 中文：RedisConfig 定义当前包使用的数据结构或接口。
// English: RedisConfig defines a data structure or interface used by this package.
// RedisConfig contains Redis connection settings.
type RedisConfig struct {
	// 中文：Addr 保存当前结构中的配置或数据值。
	// English: Addr stores a configuration or data value for this struct.
	Addr string `yaml:"addr" mapstructure:"addr"`
	// 中文：Password 保存当前结构中的配置或数据值。
	// English: Password stores a configuration or data value for this struct.
	Password string `yaml:"password" mapstructure:"password"`
	// 中文：DB 保存当前结构中的配置或数据值。
	// English: DB stores a configuration or data value for this struct.
	DB int `yaml:"db" mapstructure:"db"`
}

// 中文：DocumentConfig 定义当前包使用的数据结构或接口。
// English: DocumentConfig defines a data structure or interface used by this package.
// DocumentConfig controls document database registration.
type DocumentConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Driver 保存当前结构中的配置或数据值。
	// English: Driver stores a configuration or data value for this struct.
	Driver string `yaml:"driver" mapstructure:"driver"`
	// 中文：MongoDB 保存当前结构中的配置或数据值。
	// English: MongoDB stores a configuration or data value for this struct.
	MongoDB MongoConfig `yaml:"mongodb" mapstructure:"mongodb"`
}

// 中文：MongoConfig 定义当前包使用的数据结构或接口。
// English: MongoConfig defines a data structure or interface used by this package.
// MongoConfig contains MongoDB connection settings.
type MongoConfig struct {
	// 中文：URI 保存当前结构中的配置或数据值。
	// English: URI stores a configuration or data value for this struct.
	URI string `yaml:"uri" mapstructure:"uri"`
	// 中文：Database 保存当前结构中的配置或数据值。
	// English: Database stores a configuration or data value for this struct.
	Database string `yaml:"database" mapstructure:"database"`
	// 中文：Timeout 保存当前结构中的配置或数据值。
	// English: Timeout stores a configuration or data value for this struct.
	Timeout string `yaml:"timeout" mapstructure:"timeout"`
}

// 中文：CacheConfig 定义当前包使用的数据结构或接口。
// English: CacheConfig defines a data structure or interface used by this package.
// CacheConfig controls the Cache implementation registered into DI.
type CacheConfig struct {
	// 中文：Driver 保存当前结构中的配置或数据值。
	// English: Driver stores a configuration or data value for this struct.
	Driver string `yaml:"driver" mapstructure:"driver"`
	// 中文：Redis 保存当前结构中的配置或数据值。
	// English: Redis stores a configuration or data value for this struct.
	Redis RedisConfig `yaml:"redis" mapstructure:"redis"`
}

// 中文：MQConfig 定义当前包使用的数据结构或接口。
// English: MQConfig defines a data structure or interface used by this package.
// MQConfig controls external message-queue integration for domain events.
type MQConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Driver 保存当前结构中的配置或数据值。
	// English: Driver stores a configuration or data value for this struct.
	Driver string `yaml:"driver" mapstructure:"driver"`
	// 中文：Redis 保存当前结构中的配置或数据值。
	// English: Redis stores a configuration or data value for this struct.
	Redis RedisMQConfig `yaml:"redis" mapstructure:"redis"`
	// 中文：RabbitMQ 保存当前结构中的配置或数据值。
	// English: RabbitMQ stores a configuration or data value for this struct.
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq" mapstructure:"rabbitmq"`
	// 中文：Kafka 保存当前结构中的配置或数据值。
	// English: Kafka stores a configuration or data value for this struct.
	Kafka KafkaConfig `yaml:"kafka" mapstructure:"kafka"`
}

// 中文：RedisMQConfig 定义当前包使用的数据结构或接口。
// English: RedisMQConfig defines a data structure or interface used by this package.
// RedisMQConfig contains Redis Streams settings.
type RedisMQConfig struct {
	// 中文：Addr 保存当前结构中的配置或数据值。
	// English: Addr stores a configuration or data value for this struct.
	Addr string `yaml:"addr" mapstructure:"addr"`
	// 中文：Password 保存当前结构中的配置或数据值。
	// English: Password stores a configuration or data value for this struct.
	Password string `yaml:"password" mapstructure:"password"`
	// 中文：DB 保存当前结构中的配置或数据值。
	// English: DB stores a configuration or data value for this struct.
	DB int `yaml:"db" mapstructure:"db"`
	// 中文：Prefix 保存当前结构中的配置或数据值。
	// English: Prefix stores a configuration or data value for this struct.
	Prefix string `yaml:"prefix" mapstructure:"prefix"`
}

// 中文：RabbitMQConfig 定义当前包使用的数据结构或接口。
// English: RabbitMQConfig defines a data structure or interface used by this package.
// RabbitMQConfig contains AMQP 0-9-1 settings.
type RabbitMQConfig struct {
	// 中文：URL 保存当前结构中的配置或数据值。
	// English: URL stores a configuration or data value for this struct.
	URL string `yaml:"url" mapstructure:"url"`
	// 中文：Exchange 保存当前结构中的配置或数据值。
	// English: Exchange stores a configuration or data value for this struct.
	Exchange string `yaml:"exchange" mapstructure:"exchange"`
	// 中文：QueuePrefix 保存当前结构中的配置或数据值。
	// English: QueuePrefix stores a configuration or data value for this struct.
	QueuePrefix string `yaml:"queue_prefix" mapstructure:"queue_prefix"`
}

// 中文：KafkaConfig 定义当前包使用的数据结构或接口。
// English: KafkaConfig defines a data structure or interface used by this package.
// KafkaConfig contains Kafka producer and consumer settings.
type KafkaConfig struct {
	// 中文：Brokers 保存当前结构中的配置或数据值。
	// English: Brokers stores a configuration or data value for this struct.
	Brokers []string `yaml:"brokers" mapstructure:"brokers"`
	// 中文：ClientID 保存当前结构中的配置或数据值。
	// English: ClientID stores a configuration or data value for this struct.
	ClientID string `yaml:"client_id" mapstructure:"client_id"`
	// 中文：GroupID 保存当前结构中的配置或数据值。
	// English: GroupID stores a configuration or data value for this struct.
	GroupID string `yaml:"group_id" mapstructure:"group_id"`
	// 中文：TopicPrefix 保存当前结构中的配置或数据值。
	// English: TopicPrefix stores a configuration or data value for this struct.
	TopicPrefix string `yaml:"topic_prefix" mapstructure:"topic_prefix"`
}

// 中文：LockConfig 定义当前包使用的数据结构或接口。
// English: LockConfig defines a data structure or interface used by this package.
// LockConfig controls distributed lock registration.
type LockConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Driver 保存当前结构中的配置或数据值。
	// English: Driver stores a configuration or data value for this struct.
	Driver string `yaml:"driver" mapstructure:"driver"`
	// 中文：Redis 保存当前结构中的配置或数据值。
	// English: Redis stores a configuration or data value for this struct.
	Redis RedisConfig `yaml:"redis" mapstructure:"redis"`
	// 中文：ZooKeeper 保存当前结构中的配置或数据值。
	// English: ZooKeeper stores a configuration or data value for this struct.
	ZooKeeper ZooKeeperConfig `yaml:"zookeeper" mapstructure:"zookeeper"`
}

// 中文：ZooKeeperConfig 定义当前包使用的数据结构或接口。
// English: ZooKeeperConfig defines a data structure or interface used by this package.
// ZooKeeperConfig contains ZooKeeper lock settings.
type ZooKeeperConfig struct {
	// 中文：Servers 保存当前结构中的配置或数据值。
	// English: Servers stores a configuration or data value for this struct.
	Servers []string `yaml:"servers" mapstructure:"servers"`
	// 中文：Root 保存当前结构中的配置或数据值。
	// English: Root stores a configuration or data value for this struct.
	Root string `yaml:"root" mapstructure:"root"`
	// 中文：SessionTimeout 保存当前结构中的配置或数据值。
	// English: SessionTimeout stores a configuration or data value for this struct.
	SessionTimeout string `yaml:"session_timeout" mapstructure:"session_timeout"`
}

// 中文：QueueConfig 定义当前包使用的数据结构或接口。
// English: QueueConfig defines a data structure or interface used by this package.
// QueueConfig controls the in-process async task queue.
type QueueConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Workers 保存当前结构中的配置或数据值。
	// English: Workers stores a configuration or data value for this struct.
	Workers int `yaml:"workers" mapstructure:"workers"`
	// 中文：Buffer 保存当前结构中的配置或数据值。
	// English: Buffer stores a configuration or data value for this struct.
	Buffer int `yaml:"buffer" mapstructure:"buffer"`
	// 中文：MaxRetries 保存当前结构中的配置或数据值。
	// English: MaxRetries stores a configuration or data value for this struct.
	MaxRetries int `yaml:"max_retries" mapstructure:"max_retries"`
	// 中文：RetryDelay 保存当前结构中的配置或数据值。
	// English: RetryDelay stores a configuration or data value for this struct.
	RetryDelay string `yaml:"retry_delay" mapstructure:"retry_delay"`
}

// 中文：StorageConfig 定义当前包使用的数据结构或接口。
// English: StorageConfig defines a data structure or interface used by this package.
// StorageConfig controls object storage registration.
type StorageConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Type 保存当前结构中的配置或数据值。
	// English: Type stores a configuration or data value for this struct.
	Type string `yaml:"type" mapstructure:"type"`
	// 中文：MinIO 保存当前结构中的配置或数据值。
	// English: MinIO stores a configuration or data value for this struct.
	MinIO MinIOConfig `yaml:"minio" mapstructure:"minio"`
	// 中文：Ceph 保存当前结构中的配置或数据值。
	// English: Ceph stores a configuration or data value for this struct.
	Ceph CephConfig `yaml:"ceph" mapstructure:"ceph"`
}

// 中文：SearchConfig 定义当前包使用的数据结构或接口。
// English: SearchConfig defines a data structure or interface used by this package.
// SearchConfig controls the search engine registered into DI.
type SearchConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Driver 保存当前结构中的配置或数据值。
	// English: Driver stores a configuration or data value for this struct.
	Driver string `yaml:"driver" mapstructure:"driver"`
	// 中文：Elasticsearch 保存当前结构中的配置或数据值。
	// English: Elasticsearch stores a configuration or data value for this struct.
	Elasticsearch ElasticsearchConfig `yaml:"elasticsearch" mapstructure:"elasticsearch"`
}

// 中文：MetricsConfig 定义当前包使用的数据结构或接口。
// English: MetricsConfig defines a data structure or interface used by this package.
// MetricsConfig controls internal Prometheus-compatible metrics.
type MetricsConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Path 保存当前结构中的配置或数据值。
	// English: Path stores a configuration or data value for this struct.
	Path string `yaml:"path" mapstructure:"path"`
	// 中文：ReportPath 保存当前结构中的配置或数据值。
	// English: ReportPath stores a configuration or data value for this struct.
	ReportPath string `yaml:"report_path" mapstructure:"report_path"`
	// 中文：Namespace 保存当前结构中的配置或数据值。
	// English: Namespace stores a configuration or data value for this struct.
	Namespace string `yaml:"namespace" mapstructure:"namespace"`
}

// 中文：AlertConfig 定义当前包使用的数据结构或接口。
// English: AlertConfig defines a data structure or interface used by this package.
// AlertConfig controls operational alert sinks.
type AlertConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Logger 保存当前结构中的配置或数据值。
	// English: Logger stores a configuration or data value for this struct.
	Logger bool `yaml:"logger" mapstructure:"logger"`
	// 中文：Webhook 保存当前结构中的配置或数据值。
	// English: Webhook stores a configuration or data value for this struct.
	Webhook AlertWebhookConfig `yaml:"webhook" mapstructure:"webhook"`
	// 中文：Sentry 保存当前结构中的配置或数据值。
	// English: Sentry stores a configuration or data value for this struct.
	Sentry SentryConfig `yaml:"sentry" mapstructure:"sentry"`
}

// 中文：AlertWebhookConfig 定义当前包使用的数据结构或接口。
// English: AlertWebhookConfig defines a data structure or interface used by this package.
// AlertWebhookConfig contains webhook alert settings.
type AlertWebhookConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：URL 保存当前结构中的配置或数据值。
	// English: URL stores a configuration or data value for this struct.
	URL string `yaml:"url" mapstructure:"url"`
	// 中文：Timeout 保存当前结构中的配置或数据值。
	// English: Timeout stores a configuration or data value for this struct.
	Timeout string `yaml:"timeout" mapstructure:"timeout"`
	// 中文：Headers 保存当前结构中的配置或数据值。
	// English: Headers stores a configuration or data value for this struct.
	Headers map[string]string `yaml:"headers" mapstructure:"headers"`
}

// 中文：SentryConfig 定义当前包使用的数据结构或接口。
// English: SentryConfig defines a data structure or interface used by this package.
// SentryConfig contains Sentry alert settings.
type SentryConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：DSN 保存当前结构中的配置或数据值。
	// English: DSN stores a configuration or data value for this struct.
	DSN string `yaml:"dsn" mapstructure:"dsn"`
	// 中文：Environment 保存当前结构中的配置或数据值。
	// English: Environment stores a configuration or data value for this struct.
	Environment string `yaml:"environment" mapstructure:"environment"`
	// 中文：Release 保存当前结构中的配置或数据值。
	// English: Release stores a configuration or data value for this struct.
	Release string `yaml:"release" mapstructure:"release"`
	// 中文：TracesSampleRate 保存当前结构中的配置或数据值。
	// English: TracesSampleRate stores a configuration or data value for this struct.
	TracesSampleRate float64 `yaml:"traces_sample_rate" mapstructure:"traces_sample_rate"`
	// 中文：Debug 保存当前结构中的配置或数据值。
	// English: Debug stores a configuration or data value for this struct.
	Debug bool `yaml:"debug" mapstructure:"debug"`
	// 中文：FlushTimeout 保存当前结构中的配置或数据值。
	// English: FlushTimeout stores a configuration or data value for this struct.
	FlushTimeout string `yaml:"flush_timeout" mapstructure:"flush_timeout"`
}

// 中文：TraceConfig 定义当前包使用的数据结构或接口。
// English: TraceConfig defines a data structure or interface used by this package.
// TraceConfig controls HTTP tracing middleware and span exporting.
type TraceConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Exporter 保存当前结构中的配置或数据值。
	// English: Exporter stores a configuration or data value for this struct.
	Exporter string `yaml:"exporter" mapstructure:"exporter"` // logger, otlp, both
	// 中文：Service 保存当前结构中的配置或数据值。
	// English: Service stores a configuration or data value for this struct.
	Service TraceServiceConfig `yaml:"service" mapstructure:"service"`
	// 中文：OTLP 保存当前结构中的配置或数据值。
	// English: OTLP stores a configuration or data value for this struct.
	OTLP TraceOTLPConfig `yaml:"otlp" mapstructure:"otlp"`
}

// 中文：TraceServiceConfig 定义当前包使用的数据结构或接口。
// English: TraceServiceConfig defines a data structure or interface used by this package.
type TraceServiceConfig struct {
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `yaml:"name" mapstructure:"name"`
}

// 中文：TraceOTLPConfig 定义当前包使用的数据结构或接口。
// English: TraceOTLPConfig defines a data structure or interface used by this package.
type TraceOTLPConfig struct {
	// 中文：Endpoint 保存当前结构中的配置或数据值。
	// English: Endpoint stores a configuration or data value for this struct.
	Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`
	// 中文：Timeout 保存当前结构中的配置或数据值。
	// English: Timeout stores a configuration or data value for this struct.
	Timeout string `yaml:"timeout" mapstructure:"timeout"`
	// 中文：Headers 保存当前结构中的配置或数据值。
	// English: Headers stores a configuration or data value for this struct.
	Headers map[string]string `yaml:"headers" mapstructure:"headers"`
}

// 中文：ElasticsearchConfig 定义当前包使用的数据结构或接口。
// English: ElasticsearchConfig defines a data structure or interface used by this package.
// ElasticsearchConfig contains Elasticsearch connection settings.
type ElasticsearchConfig struct {
	// 中文：Endpoint 保存当前结构中的配置或数据值。
	// English: Endpoint stores a configuration or data value for this struct.
	Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`
	// 中文：Username 保存当前结构中的配置或数据值。
	// English: Username stores a configuration or data value for this struct.
	Username string `yaml:"username" mapstructure:"username"`
	// 中文：Password 保存当前结构中的配置或数据值。
	// English: Password stores a configuration or data value for this struct.
	Password string `yaml:"password" mapstructure:"password"`
	// 中文：Index 保存当前结构中的配置或数据值。
	// English: Index stores a configuration or data value for this struct.
	Index string `yaml:"index" mapstructure:"index"`
	// 中文：Timeout 保存当前结构中的配置或数据值。
	// English: Timeout stores a configuration or data value for this struct.
	Timeout string `yaml:"timeout" mapstructure:"timeout"`
}

// 中文：MinIOConfig 定义当前包使用的数据结构或接口。
// English: MinIOConfig defines a data structure or interface used by this package.
// MinIOConfig contains MinIO object storage settings.
type MinIOConfig struct {
	// 中文：Endpoint 保存当前结构中的配置或数据值。
	// English: Endpoint stores a configuration or data value for this struct.
	Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`
	// 中文：AccessKey 保存当前结构中的配置或数据值。
	// English: AccessKey stores a configuration or data value for this struct.
	AccessKey string `yaml:"access_key" mapstructure:"access_key"`
	// 中文：SecretKey 保存当前结构中的配置或数据值。
	// English: SecretKey stores a configuration or data value for this struct.
	SecretKey string `yaml:"secret_key" mapstructure:"secret_key"`
	// 中文：UseSSL 保存当前结构中的配置或数据值。
	// English: UseSSL stores a configuration or data value for this struct.
	UseSSL bool `yaml:"use_ssl" mapstructure:"use_ssl"`
}

// 中文：CephConfig 定义当前包使用的数据结构或接口。
// English: CephConfig defines a data structure or interface used by this package.
// CephConfig contains S3-compatible Ceph RGW settings.
type CephConfig struct {
	// 中文：Endpoint 保存当前结构中的配置或数据值。
	// English: Endpoint stores a configuration or data value for this struct.
	Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`
	// 中文：AccessKey 保存当前结构中的配置或数据值。
	// English: AccessKey stores a configuration or data value for this struct.
	AccessKey string `yaml:"access_key" mapstructure:"access_key"`
	// 中文：SecretKey 保存当前结构中的配置或数据值。
	// English: SecretKey stores a configuration or data value for this struct.
	SecretKey string `yaml:"secret_key" mapstructure:"secret_key"`
	// 中文：UseSSL 保存当前结构中的配置或数据值。
	// English: UseSSL stores a configuration or data value for this struct.
	UseSSL bool `yaml:"use_ssl" mapstructure:"use_ssl"`
	// 中文：PublicURL 保存当前结构中的配置或数据值。
	// English: PublicURL stores a configuration or data value for this struct.
	PublicURL string `yaml:"public_url" mapstructure:"public_url"`
}

// 中文：ConfigCenterConfig 定义当前包使用的数据结构或接口。
// English: ConfigCenterConfig defines a data structure or interface used by this package.
// ConfigCenterConfig contains optional configuration-center settings.
type ConfigCenterConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Type 保存当前结构中的配置或数据值。
	// English: Type stores a configuration or data value for this struct.
	Type string `yaml:"type" mapstructure:"type"`
	// 中文：Nacos 保存当前结构中的配置或数据值。
	// English: Nacos stores a configuration or data value for this struct.
	Nacos NacosConfig `yaml:"nacos" mapstructure:"nacos"`
	// 中文：Consul 保存当前结构中的配置或数据值。
	// English: Consul stores a configuration or data value for this struct.
	Consul ConsulConfig `yaml:"consul" mapstructure:"consul"`
}

// 中文：NacosConfig 定义当前包使用的数据结构或接口。
// English: NacosConfig defines a data structure or interface used by this package.
// NacosConfig contains Nacos configuration-center settings.
type NacosConfig struct {
	// 中文：ServerAddr 保存当前结构中的配置或数据值。
	// English: ServerAddr stores a configuration or data value for this struct.
	ServerAddr string `yaml:"server_addr" mapstructure:"server_addr"`
	// 中文：Namespace 保存当前结构中的配置或数据值。
	// English: Namespace stores a configuration or data value for this struct.
	Namespace string `yaml:"namespace" mapstructure:"namespace"`
	// 中文：Group 保存当前结构中的配置或数据值。
	// English: Group stores a configuration or data value for this struct.
	Group string `yaml:"group" mapstructure:"group"`
	// 中文：DataID 保存当前结构中的配置或数据值。
	// English: DataID stores a configuration or data value for this struct.
	DataID string `yaml:"data_id" mapstructure:"data_id"`
}

// 中文：ConsulConfig 定义当前包使用的数据结构或接口。
// English: ConsulConfig defines a data structure or interface used by this package.
// ConsulConfig contains Consul configuration-center settings.
type ConsulConfig struct {
	// 中文：Address 保存当前结构中的配置或数据值。
	// English: Address stores a configuration or data value for this struct.
	Address string `yaml:"address" mapstructure:"address"`
	// 中文：Scheme 保存当前结构中的配置或数据值。
	// English: Scheme stores a configuration or data value for this struct.
	Scheme string `yaml:"scheme" mapstructure:"scheme"`
	// 中文：Token 保存当前结构中的配置或数据值。
	// English: Token stores a configuration or data value for this struct.
	Token string `yaml:"token" mapstructure:"token"`
	// 中文：Key 保存当前结构中的配置或数据值。
	// English: Key stores a configuration or data value for this struct.
	Key string `yaml:"key" mapstructure:"key"`
}

// 中文：MiddlewareConfig 定义当前包使用的数据结构或接口。
// English: MiddlewareConfig defines a data structure or interface used by this package.
// MiddlewareConfig contains global middleware switches.
type MiddlewareConfig struct {
	// 中文：Recovery 保存当前结构中的配置或数据值。
	// English: Recovery stores a configuration or data value for this struct.
	Recovery bool `yaml:"recovery" mapstructure:"recovery"`
	// 中文：CORS 保存当前结构中的配置或数据值。
	// English: CORS stores a configuration or data value for this struct.
	CORS bool `yaml:"cors" mapstructure:"cors"`
	// 中文：RequestID 保存当前结构中的配置或数据值。
	// English: RequestID stores a configuration or data value for this struct.
	RequestID bool `yaml:"request_id" mapstructure:"request_id"`
	// 中文：RateLimit 保存当前结构中的配置或数据值。
	// English: RateLimit stores a configuration or data value for this struct.
	RateLimit RateLimitConfig `yaml:"rate_limit" mapstructure:"rate_limit"`
	// 中文：CircuitBreak 保存当前结构中的配置或数据值。
	// English: CircuitBreak stores a configuration or data value for this struct.
	CircuitBreak CircuitBreakConfig `yaml:"circuit_break" mapstructure:"circuit_break"`
	// 中文：Tenant 保存当前结构中的配置或数据值。
	// English: Tenant stores a configuration or data value for this struct.
	Tenant bool `yaml:"tenant" mapstructure:"tenant"`
	// 中文：Auth 保存当前结构中的配置或数据值。
	// English: Auth stores a configuration or data value for this struct.
	Auth GlobalAuthConfig `yaml:"auth" mapstructure:"auth"`
	// 中文：Idempotent 保存当前结构中的配置或数据值。
	// English: Idempotent stores a configuration or data value for this struct.
	Idempotent IdempotentConfig `yaml:"idempotent" mapstructure:"idempotent"`
	// 中文：I18n 保存当前结构中的配置或数据值。
	// English: I18n stores a configuration or data value for this struct.
	I18n bool `yaml:"i18n" mapstructure:"i18n"`
}

// 中文：GlobalAuthConfig 定义当前包使用的数据结构或接口。
// English: GlobalAuthConfig defines a data structure or interface used by this package.
// GlobalAuthConfig controls the optional global authentication middleware.
type GlobalAuthConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：JWTSecret 保存当前结构中的配置或数据值。
	// English: JWTSecret stores a configuration or data value for this struct.
	JWTSecret string `yaml:"jwt_secret" mapstructure:"jwt_secret"`
	// 中文：PublicPaths 保存当前结构中的配置或数据值。
	// English: PublicPaths stores a configuration or data value for this struct.
	PublicPaths []string `yaml:"public_paths" mapstructure:"public_paths"`
}

// 中文：RateLimitConfig 定义当前包使用的数据结构或接口。
// English: RateLimitConfig defines a data structure or interface used by this package.
// RateLimitConfig contains rate limiter settings.
type RateLimitConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Strategy 保存当前结构中的配置或数据值。
	// English: Strategy stores a configuration or data value for this struct.
	Strategy string `yaml:"strategy" mapstructure:"strategy"`
	// 中文：Rate 保存当前结构中的配置或数据值。
	// English: Rate stores a configuration or data value for this struct.
	Rate int `yaml:"rate" mapstructure:"rate"`
	// 中文：Burst 保存当前结构中的配置或数据值。
	// English: Burst stores a configuration or data value for this struct.
	Burst int `yaml:"burst" mapstructure:"burst"`
}

// 中文：CircuitBreakConfig 定义当前包使用的数据结构或接口。
// English: CircuitBreakConfig defines a data structure or interface used by this package.
// CircuitBreakConfig contains circuit breaker settings.
type CircuitBreakConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Threshold 保存当前结构中的配置或数据值。
	// English: Threshold stores a configuration or data value for this struct.
	Threshold float64 `yaml:"threshold" mapstructure:"threshold"`
	// 中文：Timeout 保存当前结构中的配置或数据值。
	// English: Timeout stores a configuration or data value for this struct.
	Timeout string `yaml:"timeout" mapstructure:"timeout"`
}

// 中文：IdempotentConfig 定义当前包使用的数据结构或接口。
// English: IdempotentConfig defines a data structure or interface used by this package.
// IdempotentConfig contains idempotency middleware settings.
type IdempotentConfig struct {
	// 中文：Enabled 保存当前结构中的配置或数据值。
	// English: Enabled stores a configuration or data value for this struct.
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// 中文：Header 保存当前结构中的配置或数据值。
	// English: Header stores a configuration or data value for this struct.
	Header string `yaml:"header" mapstructure:"header"`
}
