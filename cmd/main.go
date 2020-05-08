package main

import (
	"net/http"

	"github.com/Naist4869/awesomeProject/api/views"

	"github.com/Naist4869/awesomeProject/usecase"

	"github.com/Naist4869/awesomeProject/api/handler"
	"github.com/Naist4869/awesomeProject/model"

	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/container/servicecontainer"
	"github.com/Naist4869/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	testMongo()
	select {}
}
func testMongo() {
	c, err := buildContainer()
	if err != nil {
		log.BaseLogger.Fatal("构建容器失败", zap.Error(err))
		return
	}
	registrationUseCase, err := getRegistrationUseCase(c)
	if err != nil {
		log.BaseLogger.Fatal("获取用例失败", zap.Error(err))
		return
	}
	controller := handler.NewUserController(registrationUseCase, views.View{})
	r := http.NewServeMux()
	controller.MakeUserHandler(r)
	if err = http.ListenAndServe("127.0.0.1:1237", r); err != nil {
		log.BaseLogger.Fatal("服务端启动失败", zap.Error(err))
	}
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
	key := model.REGISTRATION
	value, err := c.BuildUseCase(key)
	if err != nil {
		return nil, errors.Wrap(err, "getRegistrationUseCase")
	}
	return value.(usecase.IUserRegistration), nil
}
