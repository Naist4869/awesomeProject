package restserver

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/Naist4869/log"
	"go.uber.org/zap"
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				if log.BaseLogger != nil {
					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					headers := strings.Split(string(httpRequest), "\r\n")
					for idx, header := range headers {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" {
							headers[idx] = current[0] + ": *"
						}
					}
					if brokenPipe {
						log.BaseLogger.Error("[Recovery] broken connection", zap.Error(err.(error)), zap.String("request:", string(httpRequest)))
					} else {
						log.BaseLogger.DPanic("[Recovery] panic recovered", zap.Error(err.(error)), zap.Stack("栈"))
					}
				}

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}
