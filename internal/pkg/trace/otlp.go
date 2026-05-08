package trace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// 中文：OTLPHTTPConfig 定义当前包使用的数据结构或接口。
// English: OTLPHTTPConfig defines a data structure or interface used by this package.
type OTLPHTTPConfig struct {
	// 中文：Endpoint 保存当前结构中的配置或数据值。
	// English: Endpoint stores a configuration or data value for this struct.
	Endpoint string
	// 中文：ServiceName 保存当前结构中的配置或数据值。
	// English: ServiceName stores a configuration or data value for this struct.
	ServiceName string
	// 中文：Timeout 保存当前结构中的配置或数据值。
	// English: Timeout stores a configuration or data value for this struct.
	Timeout time.Duration
	// 中文：Headers 保存当前结构中的配置或数据值。
	// English: Headers stores a configuration or data value for this struct.
	Headers map[string]string
}

// 中文：OTLPHTTPExporter 定义当前包使用的数据结构或接口。
// English: OTLPHTTPExporter defines a data structure or interface used by this package.
type OTLPHTTPExporter struct {
	// 中文：endpoint 保存当前结构中的配置或数据值。
	// English: endpoint stores a configuration or data value for this struct.
	endpoint string
	// 中文：serviceName 保存当前结构中的配置或数据值。
	// English: serviceName stores a configuration or data value for this struct.
	serviceName string
	// 中文：headers 保存当前结构中的配置或数据值。
	// English: headers stores a configuration or data value for this struct.
	headers map[string]string
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *http.Client
}

// 中文：NewOTLPHTTPExporter 创建并返回对应组件实例。
// English: NewOTLPHTTPExporter creates and returns the corresponding component instance.
func NewOTLPHTTPExporter(cfg OTLPHTTPConfig) (*OTLPHTTPExporter, error) {
	endpoint := strings.TrimSpace(cfg.Endpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("otlp endpoint is required")
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	serviceName := strings.TrimSpace(cfg.ServiceName)
	if serviceName == "" {
		serviceName = "spiringo"
	}
	return &OTLPHTTPExporter{
		endpoint:    endpoint,
		serviceName: serviceName,
		headers:     cloneHeaders(cfg.Headers),
		client:      &http.Client{Timeout: timeout},
	}, nil
}

// 中文：ExportSpan 执行当前包中的对应流程。
// English: ExportSpan executes the corresponding workflow in this package.
func (e *OTLPHTTPExporter) ExportSpan(ctx context.Context, span SpanSnapshot) error {
	if e == nil {
		return nil
	}
	payload := e.payload(span)
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal otlp trace payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create otlp trace request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range e.headers {
		if key != "" {
			req.Header.Set(key, value)
		}
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("send otlp trace: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("otlp trace export failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return nil
}

// 中文：payload 执行当前包中的对应流程。
// English: payload executes the corresponding workflow in this package.
func (e *OTLPHTTPExporter) payload(span SpanSnapshot) map[string]any {
	return map[string]any{
		"resourceSpans": []any{
			map[string]any{
				"resource": map[string]any{
					"attributes": []any{
						stringAttribute("service.name", e.serviceName),
					},
				},
				"scopeSpans": []any{
					map[string]any{
						"scope": map[string]any{
							"name": "github.com/spiringo/spiringo/internal/pkg/trace",
						},
						"spans": []any{
							map[string]any{
								"traceId":           span.TraceID,
								"spanId":            span.SpanID,
								"parentSpanId":      span.ParentSpanID,
								"name":              span.Name,
								"kind":              1,
								"startTimeUnixNano": unixNanoString(span.StartTime),
								"endTimeUnixNano":   unixNanoString(span.EndTime),
								"attributes":        otlpAttributes(span.Attributes),
							},
						},
					},
				},
			},
		},
	}
}

// 中文：otlpAttributes 执行当前包中的对应流程。
// English: otlpAttributes executes the corresponding workflow in this package.
func otlpAttributes(attrs map[string]string) []any {
	if len(attrs) == 0 {
		return nil
	}
	out := make([]any, 0, len(attrs))
	for key, value := range attrs {
		if key == "" {
			continue
		}
		out = append(out, stringAttribute(key, value))
	}
	return out
}

// 中文：stringAttribute 执行当前包中的对应流程。
// English: stringAttribute executes the corresponding workflow in this package.
func stringAttribute(key, value string) map[string]any {
	return map[string]any{
		"key": key,
		"value": map[string]any{
			"stringValue": value,
		},
	}
}

// 中文：unixNanoString 执行当前包中的对应流程。
// English: unixNanoString executes the corresponding workflow in this package.
func unixNanoString(t time.Time) string {
	if t.IsZero() {
		return "0"
	}
	return fmt.Sprintf("%d", t.UnixNano())
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
