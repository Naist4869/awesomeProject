package usecasefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/usecase/officialwx"
)

func OfficialWxFactory(c container.Container, appConfig *config.AppConfig) (interface{}, error) {
	useCase := officialwx.UseCase{}
	return &useCase, nil
}
