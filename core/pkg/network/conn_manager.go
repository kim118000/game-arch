package network

import (
	"errors"
	logger2 "github.com/kim118000/core/pkg/logger"
	"sync"
	"sync/atomic"
)

//ConnManager 连接管理模块
type ConnManager struct {
	connections map[uint32]IConnection
	connLock    sync.RWMutex
	count       int32
}

//NewConnManager 创建一个链接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]IConnection),
	}
}

func (connMgr* ConnManager) Count() int32 {
	return atomic.LoadInt32(&connMgr.count)
}

func (connMgr* ConnManager) counter(delta int32) {
	atomic.AddInt32(&connMgr.count, delta)
}

//Add 添加链接
func (connMgr *ConnManager) Add(conn IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//将conn连接添加到ConnMananger中
	connMgr.connections[conn.GetConnID()] = conn
	connMgr.counter(1)
	logger2.DefaultLogger.Infof("connection add to ConnManager connid=%d,successfully: conn num=%d", conn.GetConnID(), connMgr.Count())
}

//Remove 删除连接
func (connMgr *ConnManager) Remove(conn IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	connMgr.counter(-1)
	logger2.DefaultLogger.Infof("connection remove connid=%d,successfully: conn num=%d", conn.GetConnID(), connMgr.Count())
}

//Get 利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint32) (IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")

}

//ClearConn 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	logger2.DefaultLogger.Infof("clear all connections start....conn num=%d", connMgr.Count())

	connMgr.connLock.RLock()
	var arr []IConnection = make([]IConnection, 0, connMgr.Count())
	for _, conn := range connMgr.connections {
		arr = append(arr, conn)
	}
	connMgr.connLock.RUnlock()

	for _, conn := range arr {
		conn.Stop()
	}
	logger2.DefaultLogger.Infof("clear all connections successfully: conn num=%d", connMgr.Count())
}

//ClearOneConn  利用ConnID获取一个链接 并且删除
func (connMgr *ConnManager) ClearOneConn(connID uint32) {
	connMgr.connLock.RLock()
	connections := connMgr.connections
	conn, ok := connections[connID]
	connMgr.connLock.RUnlock()

	if !ok {
		logger2.DefaultLogger.Infof("clear connections id=%d not found", connID)
		return
	}
	conn.Stop()
	logger2.DefaultLogger.Infof("clear connections id=%d successfully", connID)
}

func (connMgr *ConnManager) Broadcast(msg INetMessage) {
	connMgr.connLock.RLock()
	var arr []IConnection = make([]IConnection, 0, connMgr.Count())
	for _, conn := range connMgr.connections {
		arr = append(arr, conn)
	}
	connMgr.connLock.RUnlock()

	for _, conn := range arr {
		conn.SendMsg(msg)
	}
}
