package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// 中文：WebhookNotifier 定义当前包使用的数据结构或接口。
// English: WebhookNotifier defines a data structure or interface used by this package.
type WebhookNotifier struct {
	// 中文：url 保存当前结构中的配置或数据值。
	// English: url stores a configuration or data value for this struct.
	url string
	// 中文：headers 保存当前结构中的配置或数据值。
	// English: headers stores a configuration or data value for this struct.
	headers map[string]string
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *http.Client
}

// 中文：NewWebhookNotifier 创建并返回对应组件实例。
// English: NewWebhookNotifier creates and returns the corresponding component instance.
func NewWebhookNotifier(url string, timeout time.Duration, headers map[string]string) *WebhookNotifier {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &WebhookNotifier{
		url:     url,
		headers: headers,
		client:  &http.Client{Timeout: timeout},
	}
}

// 中文：Notify 执行当前包中的对应流程。
// English: Notify executes the corresponding workflow in this package.
func (n *WebhookNotifier) Notify(ctx context.Context, a Alert) error {
	if n.url == "" {
		return fmt.Errorf("alert webhook url is required")
	}
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	body, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("marshal alert: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range n.headers {
		req.Header.Set(key, value)
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("alert webhook http %d: %s", resp.StatusCode, string(data))
	}
	return nil
}
