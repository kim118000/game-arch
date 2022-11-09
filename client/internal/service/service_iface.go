package service

import "github.com/kim118000/client/internal/conf"

type IService interface {
	Init(conf *conf.ServerConfig)
	Start()
	Stop()
}

