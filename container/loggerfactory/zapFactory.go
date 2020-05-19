package loggerfactory

import (
	"github.com/Naist4869/awesomeProject/config"
	"github.com/Naist4869/log"
)

func ZapFactory(lc config.LogConfig) error {
	logger := log.NewLogger(lc.MaxSize, lc.MaxAge, lc.LogDir, lc.Name, lc.Console, lc.Debug, lc.MinLevel["main"])
	log.SetLogger(logger)
	return nil
}
