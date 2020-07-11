package usecasefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/usecase/api"
	"github.com/Naist4869/log"
	"go.uber.org/zap"
)

func ApiFactory(c container.Container, appConfig *config.AppConfig) (interface{}, error) {
	logger := log.BaseLogger.With(zap.String("模块", "API管理"))
	useCase := api.UseCase{
		Logger: logger,
	}
	return &useCase, nil
}
