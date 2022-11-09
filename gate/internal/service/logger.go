package service

import (
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/gate/internal/conf"
	"sync"
)

var (
	GateLog      logger.Logger
	GateLogError logger.Logger
	onceLogger   sync.Once
)

func initLogger(conf *conf.ServerConfig) {
	onceLogger.Do(func() {
		GateLog = logger.NewLogger(conf.LogConfig.Level)
		GateLogError = logger.NewLogger(conf.LogConfig.Level)
	})
}
