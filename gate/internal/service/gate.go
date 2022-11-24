package service

import (
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/redis"
	"github.com/kim118000/core/pkg/scheduler"
	"github.com/kim118000/gate/internal/conf"
)

var GS = new(GateService)

type GateService struct{}

func (gs *GateService) Init(conf *conf.ServerConfig) {
	redis.InitRedis(&conf.RedisConfig)
	logger.InitLogger(&conf.LogConfig)
}

func (gs *GateService) Start() {
	go scheduler.Sched()
}

func (gs *GateService) Stop() {
	scheduler.Close()
	_ = redis.Client.Close()
}
