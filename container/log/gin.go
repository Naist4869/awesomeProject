package log

import (
	"time"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

// Logging colors, unused until zap implements colored logging -> https://github.com/uber-go/zap/issues/489
var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// setupLogging setups the logger to use zap
func SetupLogging(duration time.Duration, zap *Logger) {
	go func() {
		for range time.Tick(duration) {
			_ = zap.Sync()
		}
	}()
}

// Logger returns a gin handler func for all logging
func LoggerGin(duration time.Duration, logger *Logger) gin.HandlerFunc {
	SetupLogging(duration, logger)

	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		statusColor := colorForStatus(statusCode)
		path := c.Request.URL.Path

		switch {
		case statusCode >= 400 && statusCode <= 499:
			{
				logger.Warn("[GIN]",
					zap.String("statusColor", statusColor),
					zap.Int("statusCode", statusCode),
					zap.String("耗时", latency.String()),
					zap.String("客户端ip", clientIP),
					zap.String("方法", method),
					zap.String("路径", path),
					zap.String("错误", c.Errors.String()),
				)
			}
		case statusCode >= 500:
			{
				logger.Error("[GIN]",
					zap.Int("statusCode", statusCode),
					zap.String("耗时", latency.String()),
					zap.String("客户端ip", clientIP),
					zap.String("方法", method),
					zap.String("路径", path),
					zap.String("错误", c.Errors.String()),
				)
			}
		default:
			logger.Info("[GIN]",
				zap.Int("statusCode", statusCode),
				zap.String("耗时", latency.String()),
				zap.String("客户端ip", clientIP),
				zap.String("方法", method),
				zap.String("路径", path),
				zap.String("错误", c.Errors.String()),
			)
		}
	}
}

// coorForStatus returns a color based on the status code of the response
func colorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

// colorForMethod returns a color based on the HTTP method of the request
func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}
