package mgr

import (
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/gate/event"
	"github.com/kim118000/gate/service"
	"github.com/kim118000/gate/session"
)

var (
	ConnMgr    network.IConnManager
	ConnEvent  network.IConnEvent
	GateAllService []service.IService
)

func init() {
	ConnMgr = network.NewConnManager()
	ConnEvent = event.NewListenerConnEvent(ConnMgr)


	GateAllService = append(GateAllService, service.GS)
	GateAllService = append(GateAllService, session.SessionMgr)

	for _, v := range GateAllService {
		v.Start()

	}
}
