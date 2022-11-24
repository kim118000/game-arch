package event

import (
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/gate/internal/constant"
	"github.com/kim118000/gate/internal/session"
)

type ListenerConnEvent struct {
	ConnMgr network.IConnManager
}

func NewListenerConnEvent(connMgr network.IConnManager) *ListenerConnEvent {
	return &ListenerConnEvent{
		ConnMgr: connMgr,
	}
}

func (lce *ListenerConnEvent) GetConnMgr() network.IConnManager {
	return lce.ConnMgr
}

func (lce *ListenerConnEvent) OnCreate(conn network.IConnection) {
	lce.ConnMgr.Add(conn)
}

func (lce *ListenerConnEvent) OnConnStart(conn network.IConnection) {
	logger.Log.Infof("%s start ....", conn)
}

func (lce *ListenerConnEvent) OnConnStop(conn network.IConnection) {
	val, err := conn.GetProperty(constant.SessionAttrKey)
	if err != nil {
		return
	}
	sess, _ := val.(*session.Session)
	sess.RemoveSession()

	logger.Log.Infof("%s stop ....", conn)
}

func (lce *ListenerConnEvent) OnClose(conn network.IConnection) {
	lce.ConnMgr.Remove(conn)
}
