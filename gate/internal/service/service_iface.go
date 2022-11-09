package service

import "github.com/kim118000/gate/internal/conf"

type IService interface {
	Init(conf *conf.ServerConfig)
	Start()
	Stop()
}

