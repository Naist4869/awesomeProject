package rpcclientfactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/model"
	"google.golang.org/grpc"
)

type Factory func(container.Container, config.RpcConfig) (*grpc.ClientConn, error)

var rcFbMap = map[string]Factory{
	model.GRPC: grpcFactory,
}

// GetRpcClientFb is accessors for factoryBuilderMap
func GetRpcClientFb(key string) Factory {
	return rcFbMap[key]
}
