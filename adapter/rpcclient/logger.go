package rpcclient

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/peer"

	"github.com/Naist4869/log"
	"google.golang.org/grpc"
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

func Logger(logger *log.Logger) grpc.UnaryClientInterceptor {
	SetupLogging(logger)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		t := time.Now()
		var (
			peerInfo peer.Peer
			ip       string
		)
		opts = append(opts, grpc.Peer(&peerInfo))
		var quota float64
		if deadline, ok := ctx.Deadline(); ok {
			quota = time.Until(deadline).Seconds()
		}

		// invoker requests
		err = invoker(ctx, method, req, reply, cc, opts...)

		latency := time.Since(t)
		if peerInfo.Addr != nil {
			ip = peerInfo.Addr.String()
		}

		// todo metrics
		if err == nil {
			logger.Info("[GRPC]",
				zap.String("方法", method),
				zap.String("服务端ip", ip),
				zap.String("耗时", latency.String()),
				zap.Float64("超时", quota),
			)
		} else {
			logger.Info("[GRPC]",
				zap.String("方法", method),
				zap.String("服务端ip", ip),
				zap.String("耗时", latency.String()),
				zap.Float64("超时", quota),
				zap.String("请求", req.(fmt.Stringer).String()),
				zap.Error(err),
				zap.Stack("栈"),
			)
		}

		return
	}
}

//type logOption struct {
//	grpc.EmptyDialOption
//	grpc.EmptyCallOption
//	flag int8
//}
//
//func extractLogDialOption(opts []grpc.DialOption) (flag int8) {
//	for _, opt := range opts {
//		if logOpt, ok := opt.(logOption); ok {
//			return logOpt.flag
//		}
//	}
//	return
//}
