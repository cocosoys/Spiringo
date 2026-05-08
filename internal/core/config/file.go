package config

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 中文：FileSource 定义当前包使用的数据结构或接口。
// English: FileSource defines a data structure or interface used by this package.
// FileSource 文件配置源
type FileSource struct {
	// 中文：name 保存当前结构中的配置或数据值。
	// English: name stores a configuration or data value for this struct.
	name string
	// 中文：path 保存当前结构中的配置或数据值。
	// English: path stores a configuration or data value for this struct.
	path string
	// 中文：priority 保存当前结构中的配置或数据值。
	// English: priority stores a configuration or data value for this struct.
	priority int
	// 中文：watcher 保存当前结构中的配置或数据值。
	// English: watcher stores a configuration or data value for this struct.
	watcher *fsnotify.Watcher
}

// 中文：NewFileSource 创建并返回对应组件实例。
// English: NewFileSource creates and returns the corresponding component instance.
// NewFileSource 创建文件配置源
func NewFileSource(path string, priority int) *FileSource {
	return &FileSource{
		name:     fmt.Sprintf("file:%s", path),
		path:     path,
		priority: priority,
	}
}

// 中文：Name 执行当前包中的对应流程。
// English: Name executes the corresponding workflow in this package.
func (s *FileSource) Name() string { return s.name }

// 中文：Priority 执行当前包中的对应流程。
// English: Priority executes the corresponding workflow in this package.
func (s *FileSource) Priority() int { return s.priority }

// 中文：Read 执行当前包中的对应流程。
// English: Read executes the corresponding workflow in this package.
func (s *FileSource) Read() (map[string]any, error) {
	v := viper.New()
	v.SetConfigFile(s.path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read file %s: %w", s.path, err)
	}
	return v.AllSettings(), nil
}

// 中文：Watch 执行当前包中的对应流程。
// English: Watch executes the corresponding workflow in this package.
func (s *FileSource) Watch(ctx context.Context, onChange func(key string, value any)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("file watch: create watcher: %w", err)
	}
	s.watcher = watcher

	// 监听配置文件所在目录（因为vim等编辑器会替换文件而非修改）
	dir := filepath.Dir(s.path)
	if err := watcher.Add(dir); err != nil {
		return fmt.Errorf("file watch: add dir %s: %w", dir, err)
	}

	// 防抖：避免编辑器保存时触发多次
	var debounceTimer *time.Timer
	debounceDelay := 500 * time.Millisecond

	go func() {
		defer watcher.Close()
		for {
			select {
			case <-ctx.Done():
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// 只关心目标文件的写入和创建事件
				if filepath.Base(event.Name) != filepath.Base(s.path) {
					continue
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Rename) {
					if debounceTimer != nil {
						debounceTimer.Stop()
					}
					debounceTimer = time.AfterFunc(debounceDelay, func() {
						s.reloadAndNotify(onChange)
					})
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.Error("file config watch error", "path", s.path, "error", err)
			}
		}
	}()

	slog.Info("file config watcher started", "path", s.path)
	return nil
}

// 中文：reloadAndNotify 执行当前包中的对应流程。
// English: reloadAndNotify executes the corresponding workflow in this package.
// reloadAndNotify 重新读取配置并通知变更
func (s *FileSource) reloadAndNotify(onChange func(key string, value any)) {
	settings, err := s.Read()
	if err != nil {
		slog.Error("file config reload failed", "path", s.path, "error", err)
		return
	}

	// 逐键通知变更
	for k, v := range flattenSettings(settings, "") {
		onChange(k, v)
	}
	slog.Info("file config reloaded", "path", s.path)
}

// 中文：flattenSettings 执行当前包中的对应流程。
// English: flattenSettings executes the corresponding workflow in this package.
// flattenSettings 将嵌套map展平为点号分隔的key
func flattenSettings(settings map[string]any, prefix string) map[string]any {
	result := make(map[string]any)
	for k, v := range settings {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		if nested, ok := v.(map[string]any); ok {
			for nk, nv := range flattenSettings(nested, fullKey) {
				result[nk] = nv
			}
		} else {
			result[fullKey] = v
		}
	}
	return result
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (s *FileSource) Close() error {
	if s.watcher != nil {
		return s.watcher.Close()
	}
	return nil
}
