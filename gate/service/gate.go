package service

import "github.com/kim118000/core/pkg/scheduler"

var GS *GateService

func init() {
	GS = new(GateService)
}

type GateService struct {
	
}

func (gs *GateService) Init()  {

}

func (gs *GateService) Start()  {
	go scheduler.Sched()
}

func (gs *GateService) Stop()  {
	
}
