package toolkit

import (
	"github.com/kim118000/core/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func RegisterSignal(fn func()) {
	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	// stop server
	select {
	case signal := <-sg:
		fn()
		logger.DefaultLogger.Infof("got signal: %v, shutting down...", signal)
	}
}
