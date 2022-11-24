package logger

import "sync"

const logLevel = -1

var DefaultLogger Logger

func init() {
	DefaultLogger = NewLogger(logLevel)
}

func SetLogger(log Logger) {
	DefaultLogger = log
}

var (
	onceLogger sync.Once
	Log        Logger
	LogError   Logger
)

func InitLogger(conf *LogConfig) {
	onceLogger.Do(func() {
		Log = NewLogger(conf.Level)
		LogError = NewLogger(conf.Level)
	})
}
