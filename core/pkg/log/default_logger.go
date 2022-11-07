package log

import (
	"sync/atomic"
)

const logLevel = -1

var (
	DefaultLogger Logger
	setflag       uint32
)

func init() {
	DefaultLogger = NewLogger(logLevel)
}

func SetLogger(log Logger) {
	if atomic.CompareAndSwapUint32(&setflag, 0, 1) {
		DefaultLogger = log
	}
}
