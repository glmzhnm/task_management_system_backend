package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		statusColor := "\033[32m"
		if statusCode >= 400 && statusCode < 500 {
			statusColor = "\033[33m"
		} else if statusCode >= 500 {
			statusColor = "\033[31m"
		}
		reset := "\033[0m"

		fmt.Printf("%s[%d]%s %s | %s | %v | %s\n",
			statusColor, statusCode, reset,
			method, path, latency, clientIP,
		)
	}
}
