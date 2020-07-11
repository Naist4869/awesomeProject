package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Naist4869/awesomeProject/adapter/restserver"

	"github.com/Naist4869/awesomeProject/usecase"

	"github.com/Naist4869/awesomeProject/model"

	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/container/servicecontainer"
	"github.com/Naist4869/log"
	"go.uber.org/zap"
)

func main() {
	flag.Parse()
	Container, err := buildContainer()
	if err != nil {
		log.BaseLogger.Fatal("容器创建失败", zap.Error(err))
	}
	if err := runServer(Container); err != nil {
		log.BaseLogger.Fatal("服务启动失败", zap.Error(err))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

//func testMongo() {
//	c, err := buildContainer()
//	if err != nil {
//		log.BaseLogger.Fatal("构建容器失败", zap.Error(err))
//		return
//	}
//	registrationUseCase, err := getRegistrationUseCase(c)
//	if err != nil {
//		log.BaseLogger.Fatal("获取用例失败", zap.Error(err))
//		return
//	}
//	controller := handler.NewUserController(registrationUseCase, views.View{})
//	r := http.NewServeMux()
//	controller.MakeUserHandler(r)
//	if err = http.ListenAndServe("127.0.0.1:1237", r); err != nil {
//		log.BaseLogger.Fatal("服务端启动失败", zap.Error(err))
//	}
//}

func buildContainer() (*servicecontainer.ServiceContainer, error) {
	factoryMap := make(map[string]interface{})
	appConfig := config.AppConfig{}
	c := servicecontainer.ServiceContainer{FactoryMap: factoryMap, AppConfig: &appConfig}
	err := c.InitApp()
	if err != nil {
		return nil, fmt.Errorf("buildContainer: %w", err)
	}
	return &c, nil
}

//func getRegistrationUseCase(c container.Container) (usecase.IUserRegistration, error) {
//	key := model.REGISTRATION
//	value, err := c.BuildUseCase(key)
//	if err != nil {
//		return nil, errors.Wrap(err, "getRegistrationUseCase")
//	}
//	return value.(usecase.IUserRegistration), nil
//}
func getOfficialWxUseCase(c container.Container) (usecase.IOfficialWx, error) {
	key := model.OfficialWx
	value, err := c.BuildUseCase(key)
	if err != nil {
		return nil, fmt.Errorf("getOfficialWxUseCase:%w", err)
	}
	return value.(usecase.IOfficialWx), nil
}
func getApiUseCase(c container.Container) (usecase.IAPI, error) {
	key := model.Api
	value, err := c.BuildUseCase(key)
	if err != nil {
		return nil, fmt.Errorf("getOfficialWxUseCase:%w", err)
	}
	return value.(usecase.IAPI), nil
}
func runServer(sc *servicecontainer.ServiceContainer) (err error) {
	var (
		useCase    usecase.IOfficialWx
		apiUseCase usecase.IAPI
	)
	server := restserver.DefaultServer(&restserver.ServerConfig{
		Network:      "tcp",
		Addr:         "127.0.0.1:1237",
		Timeout:      time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	})

	wx := sc.AppConfig.UseCase.OfficialWx
	api := sc.AppConfig.UseCase.ApiUseCase
	ginlogger := log.BaseLogger.With(zap.String("模块", "HTTP"))
	useCase, err = getOfficialWxUseCase(sc)
	if err != nil {
		return err
	}
	apiUseCase, err = getApiUseCase(sc)
	if err != nil {
		return err
	}
	s := restserver.NewServer(wx.OriID, wx.AppID, wx.Token, wx.Base64AESKey, wx.Secret, ginlogger, useCase, server.Group("app"))
	apiServer := restserver.NewApiServer(apiUseCase, api.Limit, server.Group("api"))
	for _, eachApi := range s.HTTPApis() {
		if err := apiServer.Register(eachApi); err != nil {
			log.BaseLogger.Fatal("注册API失败", zap.Error(err))
		}
	}
	return server.Start()
}
