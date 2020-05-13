package usecasefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/usecase/workwx"
)

func WorkWxFactory(c container.Container, appConfig *config.AppConfig) (interface{}, error) {
	ucConfig := appConfig.UseCase.WorkWx
	uc, err := buildWorkWxData(c, &ucConfig.WorkWxDataConfig)
	if err != nil {
		return nil, err
	}

	useCase := workwx.UseCase{WorkWxDataService: uc}
	return &useCase, nil

}
