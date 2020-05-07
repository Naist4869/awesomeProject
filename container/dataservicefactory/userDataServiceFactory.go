package dataservicefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/container/dataservicefactory/userdataservicefactory"
)

func userDataServiceFactory(c container.Container, dataConfig *config.DataConfig) (interface{}, error) {
	key := dataConfig.DataStoreConfig.Code
	udsi, err := userdataservicefactory.GetUserDataServiceFb(key)(c, dataConfig)
	if err != nil {
		return nil, err
	}
	return udsi, nil
}
