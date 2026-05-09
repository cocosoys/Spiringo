package config

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// 中文：mapSource 定义当前包使用的数据结构或接口。
// English: mapSource defines a data structure or interface used by this package.
type mapSource struct {
	// 中文：name 保存当前结构中的配置或数据值。
	// English: name stores a configuration or data value for this struct.
	name string
	// 中文：priority 保存当前结构中的配置或数据值。
	// English: priority stores a configuration or data value for this struct.
	priority int
	// 中文：values 保存当前结构中的配置或数据值。
	// English: values stores a configuration or data value for this struct.
	values map[string]any
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (s mapSource) Name() string { return s.name }

// 中文：Priority 执行当前包中的对应流程。
// English: Priority executes the corresponding workflow in this package.
func (s mapSource) Priority() int { return s.priority }

// 中文：Read 执行当前包中的对应流程。
// English: Read executes the corresponding workflow in this package.
func (s mapSource) Read() (map[string]any, error) {
	return s.values, nil
}

// 中文：Watch 执行当前包中的对应流程。
// English: Watch executes the corresponding workflow in this package.
func (s mapSource) Watch(context.Context, func(string, any)) error { return nil }

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (s mapSource) Close() error { return nil }

// 中文：TestManagerLoadsSourcesByPriority 验证相关行为符合预期。
// English: TestManagerLoadsSourcesByPriority verifies the related behavior.
func TestManagerLoadsSourcesByPriority(t *testing.T) {
	m := NewManager()
	m.AddSource(mapSource{name: "low", priority: 10, values: map[string]any{"feature.value": "low"}})
	m.AddSource(mapSource{name: "high", priority: 20, values: map[string]any{"feature.value": "high"}})

	if err := m.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}
	if got := m.GetString("feature.value"); got != "high" {
		t.Fatalf("feature.value = %q", got)
	}
}

// 中文：TestManagerWatchAppliesChangesAndCallbacks 验证相关行为符合预期。
// English: TestManagerWatchAppliesChangesAndCallbacks verifies the related behavior.
func TestManagerWatchAppliesChangesAndCallbacks(t *testing.T) {
	changed := make(chan any, 1)
	src := &watchSource{ready: make(chan struct{})}

	m := NewManager()
	m.AddSource(src)
	m.OnChange("feature.value", func(value any) {
		changed <- value
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := m.Watch(ctx); err != nil {
		t.Fatalf("watch: %v", err)
	}

	close(src.ready)
	select {
	case got := <-changed:
		if got != "new" {
			t.Fatalf("callback value = %v", got)
		}
	case <-time.After(time.Second):
		t.Fatalf("change callback not called")
	}
	if got := m.GetString("feature.value"); got != "new" {
		t.Fatalf("manager value = %q", got)
	}
}

// 中文：TestManagerWatchRegistersKeyCallback 验证相关行为符合预期。
// English: TestManagerWatchRegistersKeyCallback verifies the related behavior.
func TestManagerWatchRegistersKeyCallback(t *testing.T) {
	changed := make(chan any, 1)
	m := NewManager()

	if err := m.Watch("feature.value", func(value any) {
		changed <- value
	}); err != nil {
		t.Fatalf("watch key: %v", err)
	}

	m.applyChange("feature.value", "new")
	select {
	case got := <-changed:
		if got != "new" {
			t.Fatalf("callback value = %v", got)
		}
	case <-time.After(time.Second):
		t.Fatalf("change callback not called")
	}
}

// 中文：TestEnvKeyVariantsIncludeSnakeCaseLeafFields 验证相关行为符合预期。
// English: TestEnvKeyVariantsIncludeSnakeCaseLeafFields verifies the related behavior.
func TestEnvKeyVariantsIncludeSnakeCaseLeafFields(t *testing.T) {
	variants := envKeyVariants("MQ_RABBITMQ_QUEUE_PREFIX")
	for _, want := range []string{"mq.rabbitmq.queue.prefix", "mq.rabbitmq.queue_prefix"} {
		if !containsString(variants, want) {
			t.Fatalf("variants missing %q: %#v", want, variants)
		}
	}
}

// 中文：TestNormalizeEnvAliases 验证历史环境名会映射到当前短名称。
// English: TestNormalizeEnvAliases verifies legacy environment names map to current short names.
func TestNormalizeEnvAliases(t *testing.T) {
	tests := map[string]string{
		"development": "dev",
		"DEV":         "dev",
		"production":  "prod",
		"Prod":        "prod",
		" local ":     "local",
		"test":        "test",
	}
	for input, want := range tests {
		if got := NormalizeEnv(input); got != want {
			t.Fatalf("NormalizeEnv(%q) = %q, want %q", input, got, want)
		}
	}
}

// 中文：TestManagerUsesSingleConfigEnvMarker 验证单一配置文件中的 app.env 会成为当前环境。
// English: TestManagerUsesSingleConfigEnvMarker verifies app.env in the single config file becomes the current environment.
func TestManagerUsesSingleConfigEnvMarker(t *testing.T) {
	dir := t.TempDir()
	configData := []byte(`app:
  name: base
  env: development
environments:
  dev:
    app:
      name: dev
    database:
      default:
        driver: sqlite
`)
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), configData, 0o600); err != nil {
		t.Fatal(err)
	}

	m := NewManager()
	m.SetConfigDir(dir)
	if err := m.Load(); err != nil {
		t.Fatal(err)
	}
	if got := m.Env(); got != "dev" {
		t.Fatalf("env = %q, want dev", got)
	}
	if got := m.GetString("app.name"); got != "dev" {
		t.Fatalf("app.name = %q, want dev", got)
	}
	if got := m.GetString("database.default.driver"); got != "sqlite" {
		t.Fatalf("database.default.driver = %q, want sqlite", got)
	}
}

// 中文：TestManagerSetEnvSelectsSingleConfigProfile 验证显式环境覆盖只选择同一配置文件内的环境段。
// English: TestManagerSetEnvSelectsSingleConfigProfile verifies explicit env override selects a profile from the same config file.
func TestManagerSetEnvSelectsSingleConfigProfile(t *testing.T) {
	dir := t.TempDir()
	configData := []byte(`app:
  name: base
  env: local
environments:
  local:
    app:
      name: local
  prod:
    app:
      name: prod
    server:
      mode: release
`)
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), configData, 0o600); err != nil {
		t.Fatal(err)
	}

	m := NewManager()
	m.SetConfigDir(dir)
	m.SetEnv("production")
	if err := m.Load(); err != nil {
		t.Fatal(err)
	}
	if got := m.Env(); got != "prod" {
		t.Fatalf("env = %q, want prod", got)
	}
	if got := m.GetString("app.name"); got != "prod" {
		t.Fatalf("app.name = %q, want prod", got)
	}
	if got := m.GetString("server.mode"); got != "release" {
		t.Fatalf("server.mode = %q, want release", got)
	}
}

// 中文：TestManagerLoadsSnakeCaseEnvFields 验证相关行为符合预期。
// English: TestManagerLoadsSnakeCaseEnvFields verifies the related behavior.
func TestManagerLoadsSnakeCaseEnvFields(t *testing.T) {
	t.Setenv("SP_MQ_RABBITMQ_QUEUE_PREFIX", "events")
	t.Setenv("SP_LOCK_ZOOKEEPER_SESSION_TIMEOUT", "7s")

	m := NewManager()
	if err := m.Load(); err != nil {
		t.Fatal(err)
	}

	var mq MQConfig
	if err := m.Unmarshal("mq", &mq); err != nil {
		t.Fatal(err)
	}
	if mq.RabbitMQ.QueuePrefix != "events" {
		t.Fatalf("queue prefix = %q", mq.RabbitMQ.QueuePrefix)
	}

	var lock LockConfig
	if err := m.Unmarshal("lock", &lock); err != nil {
		t.Fatal(err)
	}
	if lock.ZooKeeper.SessionTimeout != "7s" {
		t.Fatalf("session timeout = %q", lock.ZooKeeper.SessionTimeout)
	}
}

// 中文：watchSource 定义当前包使用的数据结构或接口。
// English: watchSource defines a data structure or interface used by this package.
type watchSource struct {
	// 中文：ready 保存当前结构中的配置或数据值。
	// English: ready stores a configuration or data value for this struct.
	ready chan struct{}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (s *watchSource) Name() string { return "watch" }

// 中文：Priority 执行当前包中的对应流程。
// English: Priority executes the corresponding workflow in this package.
func (s *watchSource) Priority() int { return 1 }

// 中文：Read 执行当前包中的对应流程。
// English: Read executes the corresponding workflow in this package.
func (s *watchSource) Read() (map[string]any, error) {
	return map[string]any{}, nil
}

// 中文：Watch 执行当前包中的对应流程。
// English: Watch executes the corresponding workflow in this package.
func (s *watchSource) Watch(ctx context.Context, onChange func(string, any)) error {
	go func() {
		select {
		case <-ctx.Done():
		case <-s.ready:
			onChange("feature.value", "new")
		}
	}()
	return nil
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (s *watchSource) Close() error { return nil }

// 中文：containsString 执行当前包中的对应流程。
// English: containsString executes the corresponding workflow in this package.
func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

// 中文：TestNacosSourceReadsYAML 验证相关行为符合预期。
// English: TestNacosSourceReadsYAML verifies the related behavior.
func TestNacosSourceReadsYAML(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nacos/v1/cs/configs" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("dataId") != "config.yaml" || r.URL.Query().Get("group") != "DEFAULT_GROUP" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte("app:\n  name: remote\n"))
	}))
	defer srv.Close()

	source := NewNacosSource(NacosConfig{
		ServerAddr: strings.TrimPrefix(srv.URL, "http://"),
		DataID:     "config.yaml",
		Group:      "DEFAULT_GROUP",
	}, 50)

	values, err := source.Read()
	if err != nil {
		t.Fatalf("read nacos: %v", err)
	}
	app, ok := values["app"].(map[string]any)
	if !ok || app["name"] != "remote" {
		t.Fatalf("unexpected values: %+v", values)
	}
}

// 中文：TestConsulSourceReadsYAML 验证相关行为符合预期。
// English: TestConsulSourceReadsYAML verifies the related behavior.
func TestConsulSourceReadsYAML(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/kv/app/config.yaml" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("X-Consul-Token") != "secret" {
			t.Fatalf("missing token header")
		}
		_, _ = w.Write([]byte("app:\n  name: consul\n"))
	}))
	defer srv.Close()

	source := NewConsulSource(ConsulConfig{
		Address: strings.TrimPrefix(srv.URL, "http://"),
		Scheme:  "http",
		Token:   "secret",
		Key:     "app/config.yaml",
	}, 50)

	values, err := source.Read()
	if err != nil {
		t.Fatalf("read consul: %v", err)
	}
	app, ok := values["app"].(map[string]any)
	if !ok || app["name"] != "consul" {
		t.Fatalf("unexpected values: %+v", values)
	}
}
