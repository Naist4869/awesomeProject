package rpcclientfactory

import (
	"context"

	"google.golang.org/grpc"

	"github.com/Naist4869/awesomeProject/adapter/rpcclient"
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
)

func grpcFactory(c container.Container, crc config.RpcConfig) (*grpc.ClientConn, error) {
	key := crc.Code
	rpcClientConfig := crc.Client
	if value, found := c.Get(key); found {
		return value.(*grpc.ClientConn), nil
	}
	client := rpcclient.DefaultClient()
	ctx := context.Background()
	if rpcClientConfig.DialTimeOut > 0 {
		ctx, _ = context.WithTimeout(ctx, rpcClientConfig.DialTimeOut)
	}
	cc, err := client.Dial(ctx, rpcClientConfig.Target)
	if err != nil {
		return nil, err
	}
	c.Put(key, cc)
	return cc, nil
}
