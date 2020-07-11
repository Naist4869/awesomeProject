package usecasefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/usecase/officialwx"
)

func OfficialWxFactory(c container.Container, appConfig *config.AppConfig) (interface{}, error) {
	uc := appConfig.UseCase.OfficialWx
	ori, err := buildOfficialWxRpc(c, uc)
	if err != nil {
		return nil, err
	}
	useCase := officialwx.NewUseCase(ori)
	return useCase, nil
}
