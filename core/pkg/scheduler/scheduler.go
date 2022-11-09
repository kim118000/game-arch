package scheduler

import (
	logger2 "github.com/kim118000/core/pkg/logger"
	"sync/atomic"
	"time"
)

type Task func()

var (
	chDie   = make(chan struct{})
	chExit  = make(chan struct{})

	chTasks = make(chan Task, 1<<8)
	started int32
	closed  int32
)

func try(f func()) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	f()
}

func Sched() {
	if atomic.AddInt32(&started, 1) != 1 {
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		close(chExit)
	}()

	for {
		select {
		case <-ticker.C:
			cron()

		case f := <-chTasks:
			try(f)

		case <-chDie:
			return
		}
	}
}

func Close() {
	if atomic.AddInt32(&closed, 1) != 1 {
		return
	}
	close(chDie)
	<-chExit
	logger2.DefaultLogger.Infof("scheduler stoped")
}

func PushTask(task Task) {
	chTasks <- task
}
