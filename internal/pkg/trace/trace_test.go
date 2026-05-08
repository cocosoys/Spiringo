package trace

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// 中文：captureExporter 定义当前包使用的数据结构或接口。
// English: captureExporter defines a data structure or interface used by this package.
type captureExporter struct {
	// 中文：spans 保存当前结构中的配置或数据值。
	// English: spans stores a configuration or data value for this struct.
	spans []SpanSnapshot
}

// 中文：ExportSpan 执行当前包中的对应流程。
// English: ExportSpan executes the corresponding workflow in this package.
func (e *captureExporter) ExportSpan(ctx context.Context, span SpanSnapshot) error {
	e.spans = append(e.spans, span)
	return nil
}

// 中文：TestTraceparentRoundTrip 验证相关行为符合预期。
// English: TestTraceparentRoundTrip verifies the related behavior.
func TestTraceparentRoundTrip(t *testing.T) {
	sc := SpanContext{
		TraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanID:  "00f067aa0ba902b7",
		Sampled: true,
	}
	got, err := ParseTraceparent(FormatTraceparent(sc))
	if err != nil {
		t.Fatalf("parse traceparent: %v", err)
	}
	if got != sc {
		t.Fatalf("trace context = %+v", got)
	}
}

// 中文：TestTracerExportsEndedSpan 验证相关行为符合预期。
// English: TestTracerExportsEndedSpan verifies the related behavior.
func TestTracerExportsEndedSpan(t *testing.T) {
	exporter := &captureExporter{}
	tracer := NewTracer(exporter)

	ctx, span := tracer.Start(context.Background(), "job", map[string]string{"kind": "test"})
	childCtx, child := tracer.Start(ctx, "child", nil)
	child.End(childCtx)
	span.End(ctx)

	if len(exporter.spans) != 2 {
		t.Fatalf("exported spans = %d", len(exporter.spans))
	}
	if exporter.spans[0].ParentSpanID != span.Context().SpanID {
		t.Fatalf("child parent = %q", exporter.spans[0].ParentSpanID)
	}
	if exporter.spans[1].Attributes["kind"] != "test" {
		t.Fatalf("parent attributes = %+v", exporter.spans[1].Attributes)
	}
}

// 中文：TestOTLPHTTPExporterPostsTracePayload 验证相关行为符合预期。
// English: TestOTLPHTTPExporterPostsTracePayload verifies the related behavior.
func TestOTLPHTTPExporterPostsTracePayload(t *testing.T) {
	var got map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("content-type = %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("X-Test") != "ok" {
			t.Fatalf("missing custom header")
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	exporter, err := NewOTLPHTTPExporter(OTLPHTTPConfig{
		Endpoint:    server.URL,
		ServiceName: "test-service",
		Headers:     map[string]string{"X-Test": "ok"},
	})
	if err != nil {
		t.Fatalf("NewOTLPHTTPExporter returned error: %v", err)
	}
	err = exporter.ExportSpan(context.Background(), SpanSnapshot{
		TraceID:   "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanID:    "00f067aa0ba902b7",
		Name:      "GET /health",
		StartTime: time.Unix(10, 0),
		EndTime:   time.Unix(10, int64(time.Millisecond)),
		Attributes: map[string]string{
			"http.method": "GET",
		},
	})
	if err != nil {
		t.Fatalf("ExportSpan returned error: %v", err)
	}

	resourceSpans := got["resourceSpans"].([]any)
	scopeSpans := resourceSpans[0].(map[string]any)["scopeSpans"].([]any)
	spans := scopeSpans[0].(map[string]any)["spans"].([]any)
	span := spans[0].(map[string]any)
	if span["traceId"] != "4bf92f3577b34da6a3ce929d0e0e4736" || span["name"] != "GET /health" {
		t.Fatalf("span payload = %#v", span)
	}
}

// 中文：TestOTLPHTTPExporterReportsStatusError 验证相关行为符合预期。
// English: TestOTLPHTTPExporterReportsStatusError verifies the related behavior.
func TestOTLPHTTPExporterReportsStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad", http.StatusBadGateway)
	}))
	defer server.Close()

	exporter, err := NewOTLPHTTPExporter(OTLPHTTPConfig{Endpoint: server.URL})
	if err != nil {
		t.Fatalf("NewOTLPHTTPExporter returned error: %v", err)
	}
	if err := exporter.ExportSpan(context.Background(), SpanSnapshot{Name: "job"}); err == nil {
		t.Fatal("ExportSpan returned nil error, want status error")
	}
}
