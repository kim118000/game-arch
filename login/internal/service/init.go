package service

import "github.com/kim118000/login/internal/conf"

type ServicesList []IService

var ServicesContainer = make(ServicesList, 0)

func init() {
	ServicesContainer = append(ServicesContainer, LoginService)
}

func (sl *ServicesList) InitConfig(conf *conf.ServerConfig) {
	for _, v := range *sl {
		v.Init(conf)
	}
}

func (sl *ServicesList) Start() {
	for _, v := range *sl {
		v.Start()
	}
}

func (sl *ServicesList) Stop() {
	for _, v := range *sl {
		v.Stop()
	}
}
