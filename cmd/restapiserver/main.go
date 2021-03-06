package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Naist4869/awesomeProject/model"
	"github.com/Naist4869/awesomeProject/usecase"

	"github.com/Naist4869/log"
	"go.uber.org/zap"

	"github.com/Naist4869/awesomeProject/adapter/restserver"
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/container/servicecontainer"
)

type Service struct {
	container container.Container
}

func main() {
	flag.Parse()
	container, err := buildContainer()
	if err != nil {
		log.BaseLogger.Fatal("容器创建失败", zap.Error(err))
	}
	if err := runServer(container); err != nil {
		log.BaseLogger.Fatal("服务启动失败", zap.Error(err))
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-c

}
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
