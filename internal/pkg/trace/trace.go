package trace

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"
)

// 中文：contextKey 定义当前包使用的数据结构或接口。
// English: contextKey defines a data structure or interface used by this package.
type contextKey struct{}

// 中文：SpanContext 定义当前包使用的数据结构或接口。
// English: SpanContext defines a data structure or interface used by this package.
type SpanContext struct {
	// 中文：TraceID 保存当前结构中的配置或数据值。
	// English: TraceID stores a configuration or data value for this struct.
	TraceID string
	// 中文：SpanID 保存当前结构中的配置或数据值。
	// English: SpanID stores a configuration or data value for this struct.
	SpanID string
	// 中文：Sampled 保存当前结构中的配置或数据值。
	// English: Sampled stores a configuration or data value for this struct.
	Sampled bool
}

// 中文：Valid 执行当前包中的对应流程。
// English: Valid executes the corresponding workflow in this package.
func (sc SpanContext) Valid() bool {
	return isHex(sc.TraceID, 32) && isHex(sc.SpanID, 16)
}

// 中文：SpanSnapshot 定义当前包使用的数据结构或接口。
// English: SpanSnapshot defines a data structure or interface used by this package.
type SpanSnapshot struct {
	// 中文：TraceID 保存当前结构中的配置或数据值。
	// English: TraceID stores a configuration or data value for this struct.
	TraceID string `json:"trace_id"`
	// 中文：SpanID 保存当前结构中的配置或数据值。
	// English: SpanID stores a configuration or data value for this struct.
	SpanID string `json:"span_id"`
	// 中文：ParentSpanID 保存当前结构中的配置或数据值。
	// English: ParentSpanID stores a configuration or data value for this struct.
	ParentSpanID string `json:"parent_span_id,omitempty"`
	// 中文：Name 保存当前结构中的配置或数据值。
	// English: Name stores a configuration or data value for this struct.
	Name string `json:"name"`
	// 中文：Attributes 保存当前结构中的配置或数据值。
	// English: Attributes stores a configuration or data value for this struct.
	Attributes map[string]string `json:"attributes,omitempty"`
	// 中文：StartTime 保存当前结构中的配置或数据值。
	// English: StartTime stores a configuration or data value for this struct.
	StartTime time.Time `json:"start_time"`
	// 中文：EndTime 保存当前结构中的配置或数据值。
	// English: EndTime stores a configuration or data value for this struct.
	EndTime time.Time `json:"end_time"`
	// 中文：Duration 保存当前结构中的配置或数据值。
	// English: Duration stores a configuration or data value for this struct.
	Duration time.Duration `json:"duration"`
}

// 中文：Exporter 定义当前包使用的数据结构或接口。
// English: Exporter defines a data structure or interface used by this package.
type Exporter interface {
	// 中文：ExportSpan 声明该接口需要实现的行为。
	// English: ExportSpan declares behavior required by this interface.
	ExportSpan(ctx context.Context, span SpanSnapshot) error
}

// 中文：MultiExporter 定义当前包使用的数据结构或接口。
// English: MultiExporter defines a data structure or interface used by this package.
type MultiExporter struct {
	// 中文：exporters 保存当前结构中的配置或数据值。
	// English: exporters stores a configuration or data value for this struct.
	exporters []Exporter
}

// 中文：NewMultiExporter 创建并返回对应组件实例。
// English: NewMultiExporter creates and returns the corresponding component instance.
func NewMultiExporter(exporters ...Exporter) *MultiExporter {
	cleaned := make([]Exporter, 0, len(exporters))
	for _, exporter := range exporters {
		if exporter != nil {
			cleaned = append(cleaned, exporter)
		}
	}
	return &MultiExporter{exporters: cleaned}
}

// 中文：ExportSpan 执行当前包中的对应流程。
// English: ExportSpan executes the corresponding workflow in this package.
func (e *MultiExporter) ExportSpan(ctx context.Context, span SpanSnapshot) error {
	var firstErr error
	for _, exporter := range e.exporters {
		if err := exporter.ExportSpan(ctx, span); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// 中文：Tracer 定义当前包使用的数据结构或接口。
// English: Tracer defines a data structure or interface used by this package.
type Tracer struct {
	// 中文：exporter 保存当前结构中的配置或数据值。
	// English: exporter stores a configuration or data value for this struct.
	exporter Exporter
}

// 中文：NewTracer 创建并返回对应组件实例。
// English: NewTracer creates and returns the corresponding component instance.
func NewTracer(exporter Exporter) *Tracer {
	return &Tracer{exporter: exporter}
}

// 中文：Start 执行当前包中的对应流程。
// English: Start executes the corresponding workflow in this package.
func (t *Tracer) Start(ctx context.Context, name string, attrs map[string]string) (context.Context, *Span) {
	if ctx == nil {
		ctx = context.Background()
	}
	parent := SpanContextFromContext(ctx)
	traceID := parent.TraceID
	if traceID == "" {
		traceID = randomHex(16)
	}
	sampled := parent.Sampled
	if parent.TraceID == "" {
		sampled = true
	}

	sc := SpanContext{
		TraceID: traceID,
		SpanID:  randomHex(8),
		Sampled: sampled,
	}
	if name == "" {
		name = "span"
	}
	span := &Span{
		tracer:       t,
		context:      sc,
		parentSpanID: parent.SpanID,
		name:         name,
		attributes:   cloneAttrs(attrs),
		startTime:    time.Now(),
	}
	return ContextWithSpanContext(ctx, sc), span
}

// 中文：Span 定义当前包使用的数据结构或接口。
// English: Span defines a data structure or interface used by this package.
type Span struct {
	// 中文：tracer 保存当前结构中的配置或数据值。
	// English: tracer stores a configuration or data value for this struct.
	tracer *Tracer
	// 中文：context 保存当前结构中的配置或数据值。
	// English: context stores a configuration or data value for this struct.
	context SpanContext
	// 中文：parentSpanID 保存当前结构中的配置或数据值。
	// English: parentSpanID stores a configuration or data value for this struct.
	parentSpanID string
	// 中文：name 保存当前结构中的配置或数据值。
	// English: name stores a configuration or data value for this struct.
	name string
	// 中文：startTime 保存当前结构中的配置或数据值。
	// English: startTime stores a configuration or data value for this struct.
	startTime time.Time

	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.Mutex
	// 中文：attributes 保存当前结构中的配置或数据值。
	// English: attributes stores a configuration or data value for this struct.
	attributes map[string]string
	// 中文：endOnce 保存当前结构中的配置或数据值。
	// English: endOnce stores a configuration or data value for this struct.
	endOnce sync.Once
}

// 中文：Context 执行当前包中的对应流程。
// English: Context executes the corresponding workflow in this package.
func (s *Span) Context() SpanContext {
	if s == nil {
		return SpanContext{}
	}
	return s.context
}

// 中文：SetAttribute 执行当前包中的对应流程。
// English: SetAttribute executes the corresponding workflow in this package.
func (s *Span) SetAttribute(key, value string) {
	if s == nil || key == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.attributes == nil {
		s.attributes = map[string]string{}
	}
	s.attributes[key] = value
}

// 中文：End 执行当前包中的对应流程。
// English: End executes the corresponding workflow in this package.
func (s *Span) End(ctx context.Context) {
	if s == nil {
		return
	}
	s.endOnce.Do(func() {
		if ctx == nil {
			ctx = context.Background()
		}
		end := time.Now()
		s.mu.Lock()
		attrs := cloneAttrs(s.attributes)
		s.mu.Unlock()

		snapshot := SpanSnapshot{
			TraceID:      s.context.TraceID,
			SpanID:       s.context.SpanID,
			ParentSpanID: s.parentSpanID,
			Name:         s.name,
			Attributes:   attrs,
			StartTime:    s.startTime,
			EndTime:      end,
			Duration:     end.Sub(s.startTime),
		}
		if s.tracer != nil && s.tracer.exporter != nil {
			_ = s.tracer.exporter.ExportSpan(ctx, snapshot)
		}
	})
}

// 中文：ContextWithSpanContext 执行当前包中的对应流程。
// English: ContextWithSpanContext executes the corresponding workflow in this package.
func ContextWithSpanContext(ctx context.Context, sc SpanContext) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, contextKey{}, sc)
}

// 中文：SpanContextFromContext 执行当前包中的对应流程。
// English: SpanContextFromContext executes the corresponding workflow in this package.
func SpanContextFromContext(ctx context.Context) SpanContext {
	if ctx == nil {
		return SpanContext{}
	}
	sc, _ := ctx.Value(contextKey{}).(SpanContext)
	return sc
}

// 中文：ParseTraceparent 执行当前包中的对应流程。
// English: ParseTraceparent executes the corresponding workflow in this package.
func ParseTraceparent(value string) (SpanContext, error) {
	parts := strings.Split(strings.TrimSpace(value), "-")
	if len(parts) != 4 {
		return SpanContext{}, fmt.Errorf("invalid traceparent")
	}
	if parts[0] != "00" {
		return SpanContext{}, fmt.Errorf("unsupported traceparent version: %s", parts[0])
	}
	sc := SpanContext{
		TraceID: parts[1],
		SpanID:  parts[2],
		Sampled: strings.HasSuffix(parts[3], "1"),
	}
	if !sc.Valid() || sc.TraceID == strings.Repeat("0", 32) || sc.SpanID == strings.Repeat("0", 16) {
		return SpanContext{}, fmt.Errorf("invalid traceparent identifiers")
	}
	return sc, nil
}

// 中文：FormatTraceparent 执行当前包中的对应流程。
// English: FormatTraceparent executes the corresponding workflow in this package.
func FormatTraceparent(sc SpanContext) string {
	flags := "00"
	if sc.Sampled {
		flags = "01"
	}
	return fmt.Sprintf("00-%s-%s-%s", sc.TraceID, sc.SpanID, flags)
}

// 中文：LoggerExporter 定义当前包使用的数据结构或接口。
// English: LoggerExporter defines a data structure or interface used by this package.
type LoggerExporter struct {
	// 中文：logger 保存当前结构中的配置或数据值。
	// English: logger stores a configuration or data value for this struct.
	logger *slog.Logger
}

// 中文：NewLoggerExporter 创建并返回对应组件实例。
// English: NewLoggerExporter creates and returns the corresponding component instance.
func NewLoggerExporter(logger *slog.Logger) *LoggerExporter {
	if logger == nil {
		logger = slog.Default()
	}
	return &LoggerExporter{logger: logger}
}

// 中文：ExportSpan 执行当前包中的对应流程。
// English: ExportSpan executes the corresponding workflow in this package.
func (e *LoggerExporter) ExportSpan(ctx context.Context, span SpanSnapshot) error {
	e.logger.DebugContext(ctx, "trace span",
		"trace_id", span.TraceID,
		"span_id", span.SpanID,
		"parent_span_id", span.ParentSpanID,
		"name", span.Name,
		"duration", span.Duration.String(),
		"attributes", span.Attributes,
	)
	return nil
}

// 中文：randomHex 执行当前包中的对应流程。
// English: randomHex executes the corresponding workflow in this package.
func randomHex(bytesLen int) string {
	buf := make([]byte, bytesLen)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%0*x", bytesLen*2, time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

// 中文：isHex 执行当前包中的对应流程。
// English: isHex executes the corresponding workflow in this package.
func isHex(value string, length int) bool {
	if len(value) != length {
		return false
	}
	for _, r := range value {
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') {
			continue
		}
		return false
	}
	return true
}

// 中文：cloneAttrs 执行当前包中的对应流程。
// English: cloneAttrs executes the corresponding workflow in this package.
func cloneAttrs(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
