package userdataservicefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/dataservice"
	"github.com/Naist4869/awesomeProject/model"
)

var userDataServiceFactoryBuilder = map[string]Factory{
	model.MONGO: mongoDataServiceFactory,
}

type Factory func(container.Container, *config.DataConfig) (dataservice.IUserDataService, error)

// GetDataServiceFb is accessors for factoryBuilderMap
func GetUserDataServiceFb(key string) Factory {
	return userDataServiceFactoryBuilder[key]
}
