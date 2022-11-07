package network

import (
	"github.com/kim118000/core/pkg/log"
	"github.com/kim118000/core/toolkit"
	"net"
	"sync/atomic"
	"time"
)

const CONN_CLUSTER_ADDR = "CONN_CLUSTER_ADDR"

type Client struct {
	Name      string
	addr      []ClusterAddr
	timeout   time.Duration
	closeFlag int32

	//连接发送的最大chan大小
	connSendMaxBuffLen uint32

	//读数据流水线
	encode []EncoderHandle
	decode []DecoderHandle

	connMgr IConnManager
	//连接事件
	connEvent IConnEvent
}

func NewClient(name string, addr []ClusterAddr, sendBuffLen uint32, encode []EncoderHandle, decode []DecoderHandle, opts ...ClientOption) *Client {
	c := &Client{
		Name:               name,
		addr:               addr,
		connSendMaxBuffLen: sendBuffLen,
		encode:             encode,
		decode:             decode,
		connMgr:            NewConnManager(),
	}

	for _, opt := range opts {
		opt(c)
	}

	c.connEvent = c
	return c
}

func (c *Client) connect(addr ClusterAddr) {
begin:
	if c.IsClose() {
		return
	}

	conn, err := net.DialTimeout("tcp", addr.Addr, c.timeout)
	if err != nil {
		log.DefaultLogger.Warnf("net connect failure %s error %s", addr.String(), err)
		if addr.Reconnect {
			time.Sleep(2 * time.Second)
			goto begin
		}
		return
	}

	dealConn := NewConnection(conn, addr.Id, c.encode, c.decode, c.connSendMaxBuffLen, c.connEvent)
	dealConn.SetProperty(CONN_CLUSTER_ADDR, addr)
	go dealConn.Start()
}

func (c *Client) IsClose() bool {
	close := atomic.LoadInt32(&c.closeFlag)
	if close == 1 {
		return true
	}
	return false
}

func (c *Client) Close() {
	atomic.StoreInt32(&c.closeFlag, 1)
}

func (c *Client) Start() {
	log.DefaultLogger.Infof("start client work name: %s", c.Name)

	for _, addr := range c.addr {
		go c.connect(addr)
	}
}

func (c *Client) Stop() {
	c.Close()
	//清理连接
	c.connMgr.ClearConn()
	log.DefaultLogger.Infof("stop client work name: %s", c.Name)
}

func (c *Client) GetConnMgr() IConnManager {
	return c.connMgr
}

func (c *Client) OnCreate(conn IConnection) {
	c.connMgr.Add(conn)
}

func (c *Client) OnConnStart(conn IConnection) {
	val, _ := conn.GetProperty(CONN_CLUSTER_ADDR)
	addr, ok := val.(ClusterAddr)
	if ok && addr.IsHand {
		msg := GetMessage()
		msg.SetData(toolkit.HandshakeByte())
		conn.SendMsg(msg)
	}
}

func (c *Client) OnConnStop(conn IConnection) {
}

func (c *Client) OnClose(conn IConnection) {
	c.connMgr.Remove(conn)
	val, _ := conn.GetProperty(CONN_CLUSTER_ADDR)
	addr, ok := val.(ClusterAddr)
	if ok && addr.Reconnect {
		go c.connect(addr)
	}
}
