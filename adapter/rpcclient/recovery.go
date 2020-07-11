package rpcclient

import (
	"context"
	"net"
	"os"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// recovery return a client interceptor  that recovers from any panics.
func (c *Client) recovery() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
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

				if brokenPipe {
					c.logger.Error("[Recovery] broken connection", zap.Error(err.(error)), zap.Any("request", req))
				} else {
					c.logger.DPanic("[Recovery] panic recovered", zap.Error(err.(error)), zap.Stack("栈"))
				}

			}
		}()
		err = invoker(ctx, method, req, reply, cc, opts...)
		return
	}
}
