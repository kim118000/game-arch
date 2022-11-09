package service

import (
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/login/internal/conf"
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
