package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/pkg/trace"
)

// 中文：Trace 执行当前包中的对应流程。
// English: Trace executes the corresponding workflow in this package.
func Trace(tracer *trace.Tracer) gin.HandlerFunc {
	if tracer == nil {
		tracer = trace.NewTracer(nil)
	}
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if remote, err := trace.ParseTraceparent(c.GetHeader("traceparent")); err == nil {
			ctx = trace.ContextWithSpanContext(ctx, remote)
		}

		ctx, span := tracer.Start(ctx, "http.request", map[string]string{
			"http.method": c.Request.Method,
			"http.path":   c.Request.URL.Path,
		})
		c.Request = c.Request.WithContext(ctx)
		c.Header("traceparent", trace.FormatTraceparent(span.Context()))

		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		span.SetAttribute("http.route", route)
		span.SetAttribute("http.status_code", strconv.Itoa(c.Writer.Status()))
		span.End(c.Request.Context())
	}
}
