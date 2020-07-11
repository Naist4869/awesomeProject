package rpcfactory

import (
	"github.com/Naist4869/awesomeProject/api"
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/container/rpcclientfactory"
)

func fileSystemRpcServiceFactory(c container.Container, rpcConfig config.RpcConfig) (interface{}, error) {
	key := rpcConfig.Client.Code
	conn, err := rpcclientfactory.GetRpcClientFb(key)(c, rpcConfig)
	if err != nil {
		return nil, err
	}
	fileSystemClient := api.NewFileSystemClient(conn)
	return fileSystemClient, nil
}

func tbkRpcServiceFactory(c container.Container, rpcConfig config.RpcConfig) (interface{}, error) {
	key := rpcConfig.Client.Code
	conn, err := rpcclientfactory.GetRpcClientFb(key)(c, rpcConfig)
	if err != nil {
		return nil, err
	}
	tbkClient := api.NewTBKClient(conn)
	return tbkClient, nil
}
