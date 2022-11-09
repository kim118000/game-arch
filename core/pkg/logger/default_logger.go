package logger

const logLevel = -1

var DefaultLogger Logger

func init() {
	DefaultLogger = NewLogger(logLevel)
}

func SetLogger(log Logger) {
	DefaultLogger = log
}
