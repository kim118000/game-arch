package service

import (
	"github.com/kim118000/client/internal/conf"
	"github.com/kim118000/core/pkg/logger"
	"sync"
)

var (
	onceLogger sync.Once
	Log        logger.Logger
	LogError   logger.Logger
)

func initLogger(conf *conf.ServerConfig) {
	onceLogger.Do(func() {
		Log = logger.NewLogger(conf.LogConfig.Level)
		LogError = logger.NewLogger(conf.LogConfig.Level)
	})
}
