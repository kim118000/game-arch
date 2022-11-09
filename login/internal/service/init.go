package service

import "github.com/kim118000/login/internal/conf"

var (
	services []IService
)

func init() {
	services = append(services, Service)
}

func InitService(conf *conf.ServerConfig) {
	for _, v := range services {
		v.Init(conf)
	}
}

func StartService() {
	for _, v := range services {
		v.Start()
	}
}

func StopService() {
	for _, v := range services {
		v.Stop()
	}
}
