package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 中文：Config 定义当前包使用的数据结构或接口。
// English: Config defines a data structure or interface used by this package.
type Config struct {
	// 中文：BaseURL 保存当前结构中的配置或数据值。
	// English: BaseURL stores a configuration or data value for this struct.
	BaseURL string
	// 中文：Timeout 保存当前结构中的配置或数据值。
	// English: Timeout stores a configuration or data value for this struct.
	Timeout time.Duration
	// 中文：Headers 保存当前结构中的配置或数据值。
	// English: Headers stores a configuration or data value for this struct.
	Headers map[string]string
	// 中文：Retries 保存当前结构中的配置或数据值。
	// English: Retries stores a configuration or data value for this struct.
	Retries int
	// 中文：RetryWait 保存当前结构中的配置或数据值。
	// English: RetryWait stores a configuration or data value for this struct.
	RetryWait time.Duration
	// 中文：Client 保存当前结构中的配置或数据值。
	// English: Client stores a configuration or data value for this struct.
	Client *http.Client
}

// 中文：Client 定义当前包使用的数据结构或接口。
// English: Client defines a data structure or interface used by this package.
type Client struct {
	// 中文：baseURL 保存当前结构中的配置或数据值。
	// English: baseURL stores a configuration or data value for this struct.
	baseURL string
	// 中文：headers 保存当前结构中的配置或数据值。
	// English: headers stores a configuration or data value for this struct.
	headers map[string]string
	// 中文：retries 保存当前结构中的配置或数据值。
	// English: retries stores a configuration or data value for this struct.
	retries int
	// 中文：retryWait 保存当前结构中的配置或数据值。
	// English: retryWait stores a configuration or data value for this struct.
	retryWait time.Duration
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *http.Client
}

// 中文：Error 定义当前包使用的数据结构或接口。
// English: Error defines a data structure or interface used by this package.
type Error struct {
	// 中文：StatusCode 保存当前结构中的配置或数据值。
	// English: StatusCode stores a configuration or data value for this struct.
	StatusCode int
	// 中文：Body 保存当前结构中的配置或数据值。
	// English: Body stores a configuration or data value for this struct.
	Body string
}

// 中文：Error 执行当前包中的对应流程。
// English: Error executes the corresponding workflow in this package.
func (e *Error) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("http status %d", e.StatusCode)
	}
	return fmt.Sprintf("http status %d: %s", e.StatusCode, e.Body)
}

// 中文：NewClient 创建并返回对应组件实例。
// English: NewClient creates and returns the corresponding component instance.
func NewClient(cfg Config) *Client {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}
	if cfg.RetryWait <= 0 {
		cfg.RetryWait = 100 * time.Millisecond
	}
	client := cfg.Client
	if client == nil {
		client = &http.Client{Timeout: cfg.Timeout}
	}
	return &Client{
		baseURL:   strings.TrimRight(cfg.BaseURL, "/"),
		headers:   cloneHeaders(cfg.Headers),
		retries:   cfg.Retries,
		retryWait: cfg.RetryWait,
		client:    client,
	}
}

// 中文：Get 执行当前包中的对应流程。
// English: Get executes the corresponding workflow in this package.
func (c *Client) Get(ctx context.Context, path string, out any) error {
	return c.Do(ctx, http.MethodGet, path, nil, out)
}

// 中文：Post 执行当前包中的对应流程。
// English: Post executes the corresponding workflow in this package.
func (c *Client) Post(ctx context.Context, path string, body any, out any) error {
	return c.Do(ctx, http.MethodPost, path, body, out)
}

// 中文：Put 执行当前包中的对应流程。
// English: Put executes the corresponding workflow in this package.
func (c *Client) Put(ctx context.Context, path string, body any, out any) error {
	return c.Do(ctx, http.MethodPut, path, body, out)
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (c *Client) Delete(ctx context.Context, path string, out any) error {
	return c.Do(ctx, http.MethodDelete, path, nil, out)
}

// 中文：Do 执行当前包中的对应流程。
// English: Do executes the corresponding workflow in this package.
func (c *Client) Do(ctx context.Context, method, path string, body any, out any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	payload, contentType, err := requestPayload(body)
	if err != nil {
		return err
	}

	attempts := c.retries + 1
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		reqBody := bytes.NewReader(payload)
		req, err := http.NewRequestWithContext(ctx, method, c.resolveURL(path), reqBody)
		if err != nil {
			return err
		}
		for key, value := range c.headers {
			req.Header.Set(key, value)
		}
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			if !shouldRetry(ctx, attempt, attempts, 0) {
				return err
			}
			c.sleep(ctx)
			continue
		}

		err = decodeResponse(resp, out)
		if err == nil {
			return nil
		}
		lastErr = err
		status := resp.StatusCode
		if !shouldRetry(ctx, attempt, attempts, status) {
			return err
		}
		c.sleep(ctx)
	}
	return lastErr
}

// 中文：URLWithQuery 执行当前包中的对应流程。
// English: URLWithQuery executes the corresponding workflow in this package.
func URLWithQuery(path string, values map[string]string) string {
	if len(values) == 0 {
		return path
	}
	u, err := url.Parse(path)
	if err != nil {
		return path
	}
	query := u.Query()
	for key, value := range values {
		query.Set(key, value)
	}
	u.RawQuery = query.Encode()
	return u.String()
}

// 中文：resolveURL 执行当前包中的对应流程。
// English: resolveURL executes the corresponding workflow in this package.
func (c *Client) resolveURL(path string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") || c.baseURL == "" {
		return path
	}
	if strings.HasPrefix(path, "/") {
		return c.baseURL + path
	}
	return c.baseURL + "/" + path
}

// 中文：sleep 执行当前包中的对应流程。
// English: sleep executes the corresponding workflow in this package.
func (c *Client) sleep(ctx context.Context) {
	timer := time.NewTimer(c.retryWait)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}

// 中文：requestPayload 执行当前包中的对应流程。
// English: requestPayload executes the corresponding workflow in this package.
func requestPayload(body any) ([]byte, string, error) {
	if body == nil {
		return nil, "", nil
	}
	switch v := body.(type) {
	case []byte:
		return v, "application/octet-stream", nil
	case string:
		return []byte(v), "text/plain; charset=utf-8", nil
	case io.Reader:
		data, err := io.ReadAll(v)
		return data, "application/octet-stream", err
	default:
		data, err := json.Marshal(v)
		return data, "application/json", err
	}
}

// 中文：decodeResponse 执行当前包中的对应流程。
// English: decodeResponse executes the corresponding workflow in this package.
func decodeResponse(resp *http.Response, out any) error {
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return &Error{StatusCode: resp.StatusCode, Body: string(data)}
	}
	if out == nil || len(data) == 0 {
		return nil
	}
	if writer, ok := out.(io.Writer); ok {
		_, err := writer.Write(data)
		return err
	}
	return json.Unmarshal(data, out)
}

// 中文：shouldRetry 执行当前包中的对应流程。
// English: shouldRetry executes the corresponding workflow in this package.
func shouldRetry(ctx context.Context, attempt, attempts int, status int) bool {
	if ctx.Err() != nil || attempt >= attempts-1 {
		return false
	}
	return status == 0 || status == http.StatusTooManyRequests || status >= http.StatusInternalServerError
}

// 中文：cloneHeaders 执行当前包中的对应流程。
// English: cloneHeaders executes the corresponding workflow in this package.
func cloneHeaders(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}
