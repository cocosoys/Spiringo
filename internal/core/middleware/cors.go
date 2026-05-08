package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 中文：CORSConfig 定义当前包使用的数据结构或接口。
// English: CORSConfig defines a data structure or interface used by this package.
// CORSConfig 跨域配置
type CORSConfig struct {
	// 中文：AllowOrigins 保存当前结构中的配置或数据值。
	// English: AllowOrigins stores a configuration or data value for this struct.
	AllowOrigins []string `yaml:"allow_origins"`
	// 中文：AllowMethods 保存当前结构中的配置或数据值。
	// English: AllowMethods stores a configuration or data value for this struct.
	AllowMethods []string `yaml:"allow_methods"`
	// 中文：AllowHeaders 保存当前结构中的配置或数据值。
	// English: AllowHeaders stores a configuration or data value for this struct.
	AllowHeaders []string `yaml:"allow_headers"`
	// 中文：ExposeHeaders 保存当前结构中的配置或数据值。
	// English: ExposeHeaders stores a configuration or data value for this struct.
	ExposeHeaders []string `yaml:"expose_headers"`
	// 中文：AllowCredentials 保存当前结构中的配置或数据值。
	// English: AllowCredentials stores a configuration or data value for this struct.
	AllowCredentials bool `yaml:"allow_credentials"`
	// 中文：MaxAge 保存当前结构中的配置或数据值。
	// English: MaxAge stores a configuration or data value for this struct.
	MaxAge time.Duration `yaml:"max_age"`
}

// 中文：DefaultCORSConfig 执行当前包中的对应流程。
// English: DefaultCORSConfig executes the corresponding workflow in this package.
// DefaultCORSConfig 默认跨域配置
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Accept-Language", "Authorization", "X-Language", "X-Tenant-ID", "X-Request-ID", "X-Idempotent-Key"},
		ExposeHeaders:    []string{"Content-Length", "Content-Language", "X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
}

// 中文：CORS 执行当前包中的对应流程。
// English: CORS executes the corresponding workflow in this package.
// CORS 跨域中间件
func CORS(cfg CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowed := false
		for _, o := range cfg.AllowOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", joinStrings(cfg.AllowMethods, ", "))
			c.Header("Access-Control-Allow-Headers", joinStrings(cfg.AllowHeaders, ", "))
			c.Header("Access-Control-Expose-Headers", joinStrings(cfg.ExposeHeaders, ", "))
			if cfg.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			if cfg.MaxAge > 0 {
				c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", int(cfg.MaxAge.Seconds())))
			}
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// 中文：joinStrings 执行当前包中的对应流程。
// English: joinStrings executes the corresponding workflow in this package.
func joinStrings(ss []string, sep string) string {
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
