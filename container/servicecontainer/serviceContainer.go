package servicecontainer

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container/loggerfactory"
	"github.com/Naist4869/awesomeProject/container/usecasefactory"
	"github.com/pkg/errors"
)

type ServiceContainer struct {
	FactoryMap map[string]interface{}
	AppConfig  *config.AppConfig
}

func (s *ServiceContainer) InitApp() error {
	appConfig, err := loadConfig()
	if err != nil {
		return errors.Wrap(err, "loadConfig")
	}
	s.AppConfig = appConfig
	return loadLogger(appConfig.Log)
}

func (s *ServiceContainer) BuildUseCase(code string) (interface{}, error) {
	return usecasefactory.GetUseCaseFb(code)(s, s.AppConfig)
}

func (s *ServiceContainer) Get(code string) (interface{}, bool) {
	value, found := s.FactoryMap[code]
	return value, found
}

func (s *ServiceContainer) Put(code string, value interface{}) {
	s.FactoryMap[code] = value
}

// loads the application configurations
func loadConfig() (*config.AppConfig, error) {
	ac, err := config.LoadAppConfig()
	if err != nil {
		return nil, errors.Wrap(err, "read container")
	}
	return ac, nil
}

// loads the logger
func loadLogger(lc config.LogConfig) error {
	loggerType := lc.Code
	return loggerfactory.GetLogFactoryBuilder(loggerType)(lc)
}
