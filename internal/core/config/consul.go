package config

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 中文：ConsulSource 定义当前包使用的数据结构或接口。
// English: ConsulSource defines a data structure or interface used by this package.
type ConsulSource struct {
	// 中文：cfg 保存当前结构中的配置或数据值。
	// English: cfg stores a configuration or data value for this struct.
	cfg ConsulConfig
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

// 中文：NewConsulSource 创建并返回对应组件实例。
// English: NewConsulSource creates and returns the corresponding component instance.
func NewConsulSource(cfg ConsulConfig, priority int) *ConsulSource {
	return &ConsulSource{
		cfg:      cfg,
		priority: priority,
		client:   &http.Client{Timeout: 5 * time.Second},
		interval: defaultRemoteWatchInterval,
	}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (s *ConsulSource) Name() string {
	return "consul:" + s.cfg.Key
}

// 中文：Priority 执行当前包中的对应流程。
// English: Priority executes the corresponding workflow in this package.
func (s *ConsulSource) Priority() int {
	return s.priority
}

// 中文：Read 执行当前包中的对应流程。
// English: Read executes the corresponding workflow in this package.
func (s *ConsulSource) Read() (map[string]any, error) {
	if strings.TrimSpace(s.cfg.Address) == "" {
		return nil, fmt.Errorf("consul address is required")
	}
	if strings.TrimSpace(s.cfg.Key) == "" {
		return nil, fmt.Errorf("consul key is required")
	}

	scheme := strings.TrimSpace(s.cfg.Scheme)
	if scheme == "" {
		scheme = "http"
	}
	base, err := url.Parse(scheme + "://" + strings.TrimSpace(s.cfg.Address))
	if err != nil {
		return nil, fmt.Errorf("parse consul address: %w", err)
	}
	base.Path = "/v1/kv/" + strings.TrimLeft(s.cfg.Key, "/")
	q := base.Query()
	q.Set("raw", "")
	base.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, base.String(), nil)
	if err != nil {
		return nil, err
	}
	if s.cfg.Token != "" {
		req.Header.Set("X-Consul-Token", s.cfg.Token)
	}

	data, err := doRemoteRequest(context.Background(), s.client, req)
	if err != nil {
		return nil, fmt.Errorf("read consul config: %w", err)
	}
	return parseRemoteConfig(data, s.cfg.Key)
}

// 中文：Watch 执行当前包中的对应流程。
// English: Watch executes the corresponding workflow in this package.
func (s *ConsulSource) Watch(ctx context.Context, onChange func(key string, value any)) error {
	return watchRemoteSource(ctx, s.interval, s.Read, onChange)
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (s *ConsulSource) Close() error {
	return nil
}
