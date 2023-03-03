package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lngwu11/toolgo/loggo"
	"time"
)

var logger = loggo.GetLogger("gin")

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()             // 开始时间
		c.Next()                        // 处理请求
		end := time.Now()               // 结束时间
		statusCode := c.Writer.Status() // 状态码
		latency := end.Sub(start)       // 执行时间
		clientIP := c.ClientIP()        // 请求IP
		method := c.Request.Method      // 请求方法
		reqUri := c.Request.RequestURI  // 请求路径
		logger.Infof("| %3d | %13v | %15s | %s  %s", statusCode, latency, clientIP, method, reqUri)
	}
}
