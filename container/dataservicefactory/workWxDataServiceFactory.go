package dataservicefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/container/dataservicefactory/workwxdataservicefactory"
)

func workWxDataServiceFactory(c container.Container, dataConfig *config.DataConfig) (interface{}, error) {
	key := dataConfig.DataStoreConfig.Code
	impl, err := workwxdataservicefactory.GetWorkWxDataServiceFb(key)(c, dataConfig)
	if err != nil {
		return nil, err
	}
	return impl, nil
}
