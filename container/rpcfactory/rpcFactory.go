package rpcfactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/model"
)

type Factory func(container.Container, config.RpcConfig) (interface{}, error)

var rpcServiceFactoryBuilder = map[string]Factory{
	model.FileSystem_GRPC: fileSystemRpcServiceFactory,
	model.TBK_GRPC:        tbkRpcServiceFactory,
}

// GetRpcServiceFb is accessors for factoryBuilderMap
func GetRpcServiceFb(key string) Factory {
	return rpcServiceFactoryBuilder[key]
}
