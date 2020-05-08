package datastorefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
	"github.com/Naist4869/awesomeProject/model"
)

type Factory func(container.Container, config.DataStoreConfig) (interface{}, error)

var dsFbMap = map[string]Factory{
	model.MONGO: mongoFactory,
}

// GetDataStoreFb is accessors for factoryBuilderMap
func GetDataStoreFb(key string) Factory {
	return dsFbMap[key]
}
