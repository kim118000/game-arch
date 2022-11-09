package service

import (
	"github.com/kim118000/core/pkg/scheduler"
	"github.com/kim118000/gate/internal/conf"
)

var GS = new(GateService)

type GateService struct{}

func (gs *GateService) Init(conf *conf.ServerConfig) {
	initLogger(conf)
	initRedis(conf)
}

func (gs *GateService) Start() {
	go scheduler.Sched()
}

func (gs *GateService) Stop() {
	scheduler.Close()
	_ = Redis.Close()
}
