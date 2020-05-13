package workwxdataservicefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/dataservice"
	"github.com/Naist4869/awesomeProject/model"
)

var Builder = map[string]Factory{
	model.MONGO: mongoDataServiceFactory,
}

type Factory func(container.Container, *config.DataConfig) (dataservice.IWorkWxDataService, error)

// GetDataServiceFb is accessors for factoryBuilderMap
func GetWorkWxDataServiceFb(key string) Factory {
	return Builder[key]
}
