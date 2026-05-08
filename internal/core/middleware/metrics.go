package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/internal/pkg/metrics"
)

// 中文：Metrics 执行当前包中的对应流程。
// English: Metrics executes the corresponding workflow in this package.
func Metrics(registry *metrics.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		if registry == nil {
			return
		}
		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		labels := metrics.Labels{
			"method": c.Request.Method,
			"path":   route,
			"status": strconv.Itoa(c.Writer.Status()),
		}
		registry.IncCounter("http_requests_total", labels)
		registry.ObserveSummary("http_request_duration_seconds", labels, time.Since(start).Seconds())
	}
}
