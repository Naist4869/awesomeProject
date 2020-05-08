package loggerfactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/model"
)

type Factory func(lc config.LogConfig) error

// logger mapp to map logger code to logger builder
var logFactoryBuilderMap = map[string]Factory{
	model.ZAP: ZapFactory,
}

// accessors for factoryBuilderMap
func GetLogFactoryBuilder(key string) Factory {
	return logFactoryBuilderMap[key]
}
