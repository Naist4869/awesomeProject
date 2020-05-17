package restserver

import (
	"net/http"
	"time"

	"github.com/Naist4869/awesomeProject/tool/metadata"

	"github.com/Naist4869/log"

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
func SetupLogging(zap *log.Logger) {
	go func() {
		for range time.Tick(time.Second) {
			_ = zap.Sync()
		}
	}()
}

// Logger returns a gin handler func for all logging
func Logger(logger *log.Logger) HandlerFunc {
	SetupLogging(logger)
	const noUser = "no_user"
	return func(c *Context) {
		t := time.Now()
		//req := c.Request
		path := c.Request.URL.Path
		var quota float64
		if deadline, ok := c.Context.Deadline(); ok {
			quota = time.Until(deadline).Seconds()
		}

		c.Next()

		latency := time.Since(t)
		clientIP := c.RemoteIP()
		method := c.Request.Method
		statusCode := c.status
		statusColor := colorForStatus(statusCode)
		err := c.Errors.String()

		caller := metadata.String(c, metadata.Caller)
		if caller == "" {
			caller = noUser
		}

		if len(c.RoutePath) > 0 {
			//_metricServerReqCodeTotal.Inc(c.RoutePath[1:], caller, req.Method, c.Errors.ByType(ErrorTypeWxServer).String())
			//_metricServerReqDur.Observe(int64(latency/time.Millisecond), c.RoutePath[1:], caller, req.Method)
		}

		switch {
		case statusCode >= 400 && statusCode <= 499:
			{
				logger.Warn("[GIN]",
					zap.String("statusColor", statusColor),
					zap.Int("statusCode", statusCode),
					zap.String("耗时", latency.String()),
					zap.Float64("超时", quota),
					zap.String("客户端ip", clientIP),
					zap.String("方法", method),
					zap.String("路径", path),
					zap.String("用户", caller),
					zap.String("错误", err),
				)
			}
		case statusCode >= 500 || err != "":
			{
				logger.Error("[GIN]",
					zap.String("statusColor", statusColor),
					zap.Int("statusCode", statusCode),
					zap.String("耗时", latency.String()),
					zap.Float64("超时", quota),
					zap.String("客户端ip", clientIP),
					zap.String("方法", method),
					zap.String("路径", path),
					zap.String("用户", caller),
					zap.String("错误", err),
				)
			}
		default:
			logger.Info("[GIN]",
				zap.String("statusColor", statusColor),
				zap.Int("statusCode", statusCode),
				zap.String("耗时", latency.String()),
				zap.Float64("超时", quota),
				zap.String("客户端ip", clientIP),
				zap.String("方法", method),
				zap.String("路径", path),
				zap.String("用户", caller),
			)
		}
	}
}

// coorForStatus returns a color based on the status code of the response
func colorForStatus(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
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
