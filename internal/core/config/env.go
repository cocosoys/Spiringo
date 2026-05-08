package config

import (
	"context"
	"os"
	"strings"
)

// 中文：EnvSource 定义当前包使用的数据结构或接口。
// English: EnvSource defines a data structure or interface used by this package.
// EnvSource 环境变量配置源
// 环境变量前缀 SP_，分隔符用 _，如 SP_APP_NAME=spiringo → app.name=spiringo
type EnvSource struct {
	// 中文：priority 保存当前结构中的配置或数据值。
	// English: priority stores a configuration or data value for this struct.
	priority int
	// 中文：prefix 保存当前结构中的配置或数据值。
	// English: prefix stores a configuration or data value for this struct.
	prefix string
}

// 中文：NewEnvSource 创建并返回对应组件实例。
// English: NewEnvSource creates and returns the corresponding component instance.
// NewEnvSource 创建环境变量配置源
func NewEnvSource(priority int) *EnvSource {
	return &EnvSource{
		priority: priority,
		prefix:   "SP_",
	}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (s *EnvSource) Name() string { return "env" }

// 中文：Priority 执行当前包中的对应流程。
// English: Priority executes the corresponding workflow in this package.
func (s *EnvSource) Priority() int { return s.priority }

// 中文：Read 执行当前包中的对应流程。
// English: Read executes the corresponding workflow in this package.
func (s *EnvSource) Read() (map[string]any, error) {
	result := make(map[string]any)
	for _, envVar := range os.Environ() {
		if !strings.HasPrefix(envVar, s.prefix) {
			continue
		}
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimPrefix(parts[0], s.prefix)
		for _, variant := range envKeyVariants(key) {
			result[variant] = parts[1]
		}
	}
	return result, nil
}

// 中文：envKeyVariants 执行当前包中的对应流程。
// English: envKeyVariants executes the corresponding workflow in this package.
func envKeyVariants(key string) []string {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(key)), "_")
	if len(parts) == 0 {
		return nil
	}
	var result []string
	var walk func(int, string)
	walk = func(idx int, current string) {
		if idx == len(parts) {
			result = append(result, current)
			return
		}
		walk(idx+1, current+"."+parts[idx])
		walk(idx+1, current+"_"+parts[idx])
	}
	walk(1, parts[0])
	return uniqueStrings(result)
}

// 中文：uniqueStrings 执行当前包中的对应流程。
// English: uniqueStrings executes the corresponding workflow in this package.
func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

// 中文：Watch 执行当前包中的对应流程。
// English: Watch executes the corresponding workflow in this package.
func (s *EnvSource) Watch(_ context.Context, _ func(string, any)) error {
	// 环境变量在进程生命周期内通常不变，无需监听
	// 如需动态更新环境变量，建议通过配置中心（Nacos/Consul）实现
	return nil
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (s *EnvSource) Close() error { return nil }
