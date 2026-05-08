package config

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// 中文：defaultRemoteWatchInterval 声明当前包使用的常量。
// English: defaultRemoteWatchInterval declares constants used by this package.
const defaultRemoteWatchInterval = 5 * time.Second

// 中文：parseRemoteConfig 执行当前包中的对应流程。
// English: parseRemoteConfig executes the corresponding workflow in this package.
func parseRemoteConfig(data []byte, name string) (map[string]any, error) {
	configType := strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), ".")
	switch configType {
	case "yaml", "yml":
		configType = "yaml"
	case "json":
		configType = "json"
	case "toml":
		configType = "toml"
	case "":
		configType = "yaml"
	default:
		return nil, fmt.Errorf("unsupported remote config type: %s", configType)
	}

	v := viper.New()
	v.SetConfigType(configType)
	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, fmt.Errorf("parse remote config: %w", err)
	}
	return v.AllSettings(), nil
}

// 中文：doRemoteRequest 执行当前包中的对应流程。
// English: doRemoteRequest executes the corresponding workflow in this package.
func doRemoteRequest(ctx context.Context, client *http.Client, req *http.Request) ([]byte, error) {
	if client == nil {
		client = http.DefaultClient
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("remote config http %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return data, nil
}

// 中文：watchRemoteSource 执行当前包中的对应流程。
// English: watchRemoteSource executes the corresponding workflow in this package.
func watchRemoteSource(ctx context.Context, interval time.Duration, read func() (map[string]any, error), onChange func(key string, value any)) error {
	if interval <= 0 {
		interval = defaultRemoteWatchInterval
	}
	current, err := read()
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				next, err := read()
				if err != nil {
					continue
				}
				if reflect.DeepEqual(current, next) {
					continue
				}
				current = next
				for k, v := range flattenSettings(next, "") {
					onChange(k, v)
				}
			}
		}
	}()
	return nil
}

// 中文：withHTTPDefaultScheme 执行当前包中的对应流程。
// English: withHTTPDefaultScheme executes the corresponding workflow in this package.
func withHTTPDefaultScheme(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return value
	}
	return "http://" + value
}
