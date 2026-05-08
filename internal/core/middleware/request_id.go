package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：RequestID 执行当前包中的对应流程。
// English: RequestID executes the corresponding workflow in this package.
// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 注入到context和响应头
		ctx := types.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
