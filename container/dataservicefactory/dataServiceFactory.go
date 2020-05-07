package dataservicefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
)

type Factory func(container.Container, *config.DataConfig) (interface{}, error)

var dataServiceFactoryBuilder = map[string]Factory{
	config.USER_DATA: userDataServiceFactory,
}

// GetDataServiceFb is accessors for factoryBuilderMap
func GetDataServiceFb(key string) Factory {
	return dataServiceFactoryBuilder[key]
}
