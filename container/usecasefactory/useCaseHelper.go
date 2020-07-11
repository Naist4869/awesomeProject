package usecasefactory

import (
	"github.com/Naist4869/awesomeProject/api"
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/container/dataservicefactory"
	"github.com/Naist4869/awesomeProject/container/rpcfactory"
	"github.com/Naist4869/awesomeProject/dataservice"
	"github.com/Naist4869/awesomeProject/dataservice/officialWxDao/grpc"
)

func buildUserData(c container.Container, dc *config.DataConfig) (dataservice.IUserDataService, error) {
	dsi, err := dataservicefactory.GetDataServiceFb(dc.Code)(c, dc)
	if err != nil {
		return nil, err
	}
	impl := dsi.(dataservice.IUserDataService)
	return impl, nil
}

func buildWorkWxData(c container.Container, dc *config.DataConfig) (dataservice.IWorkWxDataService, error) {
	dsi, err := dataservicefactory.GetDataServiceFb(dc.Code)(c, dc)
	if err != nil {
		return nil, err
	}
	impl := dsi.(dataservice.IWorkWxDataService)
	return impl, nil
}

func buildOfficialWxRpc(c container.Container, officialWxConfig config.OfficialWxConfig) (impl dataservice.IOfficialWxRpcService, err error) {
	var fileSystemClient, tbkClient interface{}
	fileSystemClient, err = rpcfactory.GetRpcServiceFb(officialWxConfig.FileSystemConfig.Code)(c, officialWxConfig.FileSystemConfig)
	if err != nil {
		return
	}
	tbkClient, err = rpcfactory.GetRpcServiceFb(officialWxConfig.TBKConfig.Code)(c, officialWxConfig.TBKConfig)
	if err != nil {
		return
	}
	impl = grpc.NewOfficialWxClient(fileSystemClient.(api.FileSystemClient), tbkClient.(api.TBKClient))
	return
}
