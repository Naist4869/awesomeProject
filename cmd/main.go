package main

import (
	"fmt"
	"time"

	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/container/servicecontainer"
	"github.com/Naist4869/awesomeProject/model/usermodel"
	"github.com/Naist4869/awesomeProject/tool"
	"github.com/Naist4869/awesomeProject/usecase"
	"github.com/Naist4869/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	testMongo()
}
func testMongo() {
	c, err := buildContainer()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	testRegisterUser(c)
}

func buildContainer() (container.Container, error) {
	factoryMap := make(map[string]interface{})
	appConfig := config.AppConfig{}
	c := servicecontainer.ServiceContainer{FactoryMap: factoryMap, AppConfig: &appConfig}
	err := c.InitApp()
	if err != nil {
		return nil, errors.Wrap(err, "buildContainer")
	}
	return &c, nil
}

func getRegistrationUseCase(c container.Container) (usecase.IUserRegistration, error) {
	key := config.REGISTRATION
	value, err := c.BuildUseCase(key)
	if err != nil {
		// logger.Log.Errorf("%+v\n", err)
		return nil, errors.Wrap(err, "getRegistrationUseCase")
	}
	return value.(usecase.IUserRegistration), nil

}
func testRegisterUser(container container.Container) {
	ruci, err := getRegistrationUseCase(container)
	if err != nil {
		log.BaseLogger.Fatal("registration interface build failed\n", zap.Error(err))
	}
	created, err := time.Parse(tool.FormatIso8601Date, "2018-12-09")
	if err != nil {
		log.BaseLogger.Error("date format\n", zap.Error(err))
		return
	}

	user := usermodel.User{
		ID:           1,
		Phone:        "123123123123",
		NickName:     "Lan",
		PID:          0,
		Status:       0,
		AddTime:      created,
		ActivateTime: time.Time{},
		OperateTime:  time.Time{},
		AgentTrait:   usermodel.AgentTrait{},
		WeChatTrait:  usermodel.WeChatTrait{},
		Deleted:      false,
		Meta:         usermodel.DbMeta{Version: 1},
	}

	if err := ruci.RegisterUser(&user); err != nil {
		log.BaseLogger.Error("user registration failed\n", zap.Error(err))
	}
	log.BaseLogger.Info("new user registered")

}
