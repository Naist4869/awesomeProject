package dataservicefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/model"
)

type Factory func(container.Container, *config.DataConfig) (interface{}, error)

var dataServiceFactoryBuilder = map[string]Factory{
	model.USER_DATA:   userDataServiceFactory,
	model.WORKWX_DATA: workWxDataServiceFactory,
}

// GetDataServiceFb is accessors for factoryBuilderMap
func GetDataServiceFb(key string) Factory {
	return dataServiceFactoryBuilder[key]
}
