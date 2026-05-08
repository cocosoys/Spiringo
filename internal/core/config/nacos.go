package config

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 中文：NacosSource 定义当前包使用的数据结构或接口。
// English: NacosSource defines a data structure or interface used by this package.
type NacosSource struct {
	// 中文：cfg 保存当前结构中的配置或数据值。
	// English: cfg stores a configuration or data value for this struct.
	cfg NacosConfig
	// 中文：priority 保存当前结构中的配置或数据值。
	// English: priority stores a configuration or data value for this struct.
	priority int
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *http.Client
	// 中文：interval 保存当前结构中的配置或数据值。
	// English: interval stores a configuration or data value for this struct.
	interval time.Duration
}

// 中文：NewNacosSource 创建并返回对应组件实例。
// English: NewNacosSource creates and returns the corresponding component instance.
func NewNacosSource(cfg NacosConfig, priority int) *NacosSource {
	return &NacosSource{
		cfg:      cfg,
		priority: priority,
		client:   &http.Client{Timeout: 5 * time.Second},
		interval: defaultRemoteWatchInterval,
	}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (s *NacosSource) Name() string {
	return "nacos:" + s.cfg.DataID
}

// 中文：Priority 执行当前包中的对应流程。
// English: Priority executes the corresponding workflow in this package.
func (s *NacosSource) Priority() int {
	return s.priority
}

// 中文：Read 执行当前包中的对应流程。
// English: Read executes the corresponding workflow in this package.
func (s *NacosSource) Read() (map[string]any, error) {
	if strings.TrimSpace(s.cfg.ServerAddr) == "" {
		return nil, fmt.Errorf("nacos server_addr is required")
	}
	if strings.TrimSpace(s.cfg.DataID) == "" {
		return nil, fmt.Errorf("nacos data_id is required")
	}

	base, err := url.Parse(withHTTPDefaultScheme(s.cfg.ServerAddr))
	if err != nil {
		return nil, fmt.Errorf("parse nacos server_addr: %w", err)
	}
	base.Path = "/nacos/v1/cs/configs"
	q := base.Query()
	q.Set("dataId", s.cfg.DataID)
	if s.cfg.Group != "" {
		q.Set("group", s.cfg.Group)
	} else {
		q.Set("group", "DEFAULT_GROUP")
	}
	if s.cfg.Namespace != "" {
		q.Set("tenant", s.cfg.Namespace)
	}
	base.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, base.String(), nil)
	if err != nil {
		return nil, err
	}
	data, err := doRemoteRequest(context.Background(), s.client, req)
	if err != nil {
		return nil, fmt.Errorf("read nacos config: %w", err)
	}
	return parseRemoteConfig(data, s.cfg.DataID)
}

// 中文：Watch 执行当前包中的对应流程。
// English: Watch executes the corresponding workflow in this package.
func (s *NacosSource) Watch(ctx context.Context, onChange func(key string, value any)) error {
	return watchRemoteSource(ctx, s.interval, s.Read, onChange)
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (s *NacosSource) Close() error {
	return nil
}
