package service

import (
	"github.com/kim118000/client/internal/conf"
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/redis"
)

var ClientService = new(service)

type service struct{}

func (s *service) Init(conf *conf.ServerConfig) {
	redis.InitRedis(&conf.RedisConfig)
	logger.InitLogger(&conf.LogConfig)
}

func (s *service) Start() {

}

func (s *service) Stop() {
	redis.CloseRedis()
}
