package gin_middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gostring"
)

// SetTraceId 设置跟踪编号 https://www.jianshu.com/p/2a1a74ad3c3a
func SetTraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.Request.Header.Get("X-Request-Id")
		if requestId == "" {
			requestId = gostring.GetUuId()
		}
		c.Set("trace_id", requestId)
		c.Writer.Header().Set("X-Request-Id", requestId)
		c.Next()
	}
}

// GetTraceId 通过gin中间件获取跟踪编号
func GetTraceId(c *gin.Context) string {
	return fmt.Sprintf("%s", c.MustGet("trace_id"))
}
