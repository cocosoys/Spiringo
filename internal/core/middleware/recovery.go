package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/spiringo/spiringo/pkg/types"
)

// 中文：Recovery 执行当前包中的对应流程。
// English: Recovery executes the corresponding workflow in this package.
// Recovery 崩溃恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				c.JSON(http.StatusInternalServerError, types.Response{
					Code:    10000,
					Message: "internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
