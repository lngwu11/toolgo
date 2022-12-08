package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lngwu11/toolgo"
	"time"
)

var logger = toolgo.GetLogger("gin")

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		// 执行时间
		latency := end.Sub(start)
		// 请求路径
		reqUri := c.Request.RequestURI
		// 请求IP
		clientIP := c.ClientIP()
		// 请求方法
		method := c.Request.Method
		// 状态码
		statusCode := c.Writer.Status()
		logger.Infof("| %3d | %13v | %15s | %s  %s", statusCode, latency, clientIP, method, reqUri)
	}
}
