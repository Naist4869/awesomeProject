package usercasefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/usecase/userregistration"
)

func RegistrationFactory(c container.Container, appConfig *config.AppConfig) (interface{}, error) {
	uc := appConfig.UseCase.Registration
	udi, err := buildUserData(c, &uc.UserDataConfig)
	if err != nil {
		return nil, err
	}
	uruc := userregistration.UseCase{UserDataInterface: udi}
	return &uruc, nil
}
