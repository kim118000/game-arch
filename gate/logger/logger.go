package logger

import (
	"github.com/kim118000/core/pkg/log"
	"sync"
)

var (
	Info log.Logger
	Error log.Logger
	once sync.Once
)

func InitLogger(conf *log.LogConfig)  {
	once.Do(func() {
		Info = log.NewLogger(conf.Level)
		Error = log.NewLogger(conf.Level)
	})
}

