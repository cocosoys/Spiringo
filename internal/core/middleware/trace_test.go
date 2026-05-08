package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/pkg/trace"
)

// 中文：traceCaptureExporter 定义当前包使用的数据结构或接口。
// English: traceCaptureExporter defines a data structure or interface used by this package.
type traceCaptureExporter struct {
	// 中文：spans 保存当前结构中的配置或数据值。
	// English: spans stores a configuration or data value for this struct.
	spans []trace.SpanSnapshot
}

// 中文：ExportSpan 执行当前包中的对应流程。
// English: ExportSpan executes the corresponding workflow in this package.
func (e *traceCaptureExporter) ExportSpan(ctx context.Context, span trace.SpanSnapshot) error {
	e.spans = append(e.spans, span)
	return nil
}

// 中文：TestTraceMiddlewarePropagatesTraceparent 验证相关行为符合预期。
// English: TestTraceMiddlewarePropagatesTraceparent verifies the related behavior.
func TestTraceMiddlewarePropagatesTraceparent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	exporter := &traceCaptureExporter{}
	tracer := trace.NewTracer(exporter)

	r := gin.New()
	r.Use(Trace(tracer))
	r.GET("/users/:id", func(c *gin.Context) {
		c.String(http.StatusAccepted, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	req.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d", rec.Code)
	}
	if rec.Header().Get("traceparent") == "" {
		t.Fatalf("missing response traceparent")
	}
	if len(exporter.spans) != 1 {
		t.Fatalf("exported spans = %d", len(exporter.spans))
	}
	span := exporter.spans[0]
	if span.TraceID != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Fatalf("trace id = %s", span.TraceID)
	}
	if span.ParentSpanID != "00f067aa0ba902b7" {
		t.Fatalf("parent span = %s", span.ParentSpanID)
	}
	if span.Attributes["http.route"] != "/users/:id" || span.Attributes["http.status_code"] != "202" {
		t.Fatalf("attributes = %+v", span.Attributes)
	}
}
