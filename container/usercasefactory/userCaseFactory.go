package usercasefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
)

var UseCaseFactoryBuilderMap = map[string]Factory{
	config.REGISTRATION: RegistrationFactory,
}

type Factory func(c container.Container, appConfig *config.AppConfig) (interface{}, error)

func GetUseCaseFb(key string) Factory {
	return UseCaseFactoryBuilderMap[key]
}
