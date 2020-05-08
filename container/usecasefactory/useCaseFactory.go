package usecasefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/model"
)

var BuilderMap = map[string]Factory{
	model.REGISTRATION: RegistrationFactory,
}

type Factory func(c container.Container, appConfig *config.AppConfig) (interface{}, error)

func GetUseCaseFb(key string) Factory {
	return BuilderMap[key]
}
