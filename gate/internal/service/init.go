package service

import (
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/gate/internal/conf"
	"github.com/kim118000/gate/internal/event"
	"github.com/kim118000/gate/internal/session"
)

var (
	ConnMgr    network.IConnManager
	ConnEvent  network.IConnEvent
	services []IService
)

func init() {
	ConnMgr = network.NewConnManager()
	ConnEvent = event.NewListenerConnEvent(ConnMgr)

	services = append(services, GS)
	services = append(services, session.SessionMgr)
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
