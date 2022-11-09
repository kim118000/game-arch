package service

import "github.com/kim118000/login/internal/conf"

type IService interface {
	Init(conf *conf.ServerConfig)
	Start()
	Stop()
}

