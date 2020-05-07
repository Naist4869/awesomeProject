package datastorefactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/awesomeProject/container"
)

type Factory func(container.Container, config.DataStoreConfig) (interface{}, error)

var dsFbMap = map[string]Factory{
	config.MONGO: mongoFactory,
}

// GetDataStoreFb is accessors for factoryBuilderMap
func GetDataStoreFb(key string) Factory {
	return dsFbMap[key]
}
