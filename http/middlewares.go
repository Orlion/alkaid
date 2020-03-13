package http

import (
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Orlion/alkaid/client"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func logrusLogger(log *client.Log) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 日志
		log.Trace(logrus.Fields{
			"statusCode":  statusCode,
			"latencyTime": latencyTime,
			"clientIP":    clientIP,
			"reqMethod":   reqMethod,
			"reqUri":      reqUri,
		}, "[App] Http access")
	}
}

func recovery(log *client.Log) gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				log.Error(logrus.Fields{
					"err": err,
				}, "[App] Http recovery")

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					context.Error(err.(error)) // nolint: errcheck
					context.Abort()
				} else {
					context.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		context.Next()
	}
}
