package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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
func runServer(sc *servicecontainer.ServiceContainer) error {
	server := restserver.DefaultServer(&restserver.ServerConfig{
		Network:      "tcp",
		Addr:         "127.0.0.1:1237",
		Timeout:      time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	})
	wx := sc.AppConfig.UseCase.OfficialWx
	ginlogger := log.BaseLogger.With(zap.String("模块", "HTTP"))
	s := restserver.NewServer(wx.OriID, wx.AppID, wx.Token, wx.Base64AESKey, ginlogger, server.Group("app"))
	s.Http()
	return server.Start()
}
