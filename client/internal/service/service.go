package service

import "github.com/kim118000/client/internal/conf"

var Service = new(service)

type service struct{}

func (s *service) Init(conf *conf.ServerConfig) {
	initLogger(conf)
	initRedis(conf)
}

func (s *service) Start() {

}

func (s *service) Stop() {
	_ = Redis.Close()
}
