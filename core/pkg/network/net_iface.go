package network

import (
	"context"
	"google.golang.org/protobuf/proto"
	"net"
)

type IConnEvent interface {
	OnCreate(conn IConnection)
	OnConnStart(conn IConnection)
	OnConnStop(conn IConnection)
	OnClose(conn IConnection)
	GetConnMgr() IConnManager //得到链接管理
}

type IConnManager interface {
	Add(conn IConnection)                   //添加链接
	Remove(conn IConnection)                //删除连接
	Get(connID uint32) (IConnection, error) //利用ConnID获取链接
	Count() int32                           //获取当前连接
	ClearConn()                             //删除并停止所有链接
	ClearOneConn(connID uint32)             //删除指定id并停止连接
	Broadcast(msg INetMessage)              //广播
}

type INetMessage interface {
	Init(msgId uint32, serial proto.Message)
	GetSerialData() proto.Message
	GetDataLen() uint32
	GetMsgID() uint32
	GetData() []byte
	SetData(data []byte)
}

//定义连接接口
type IConnection interface {
	Start()                   //启动连接，让当前连接开始工作
	Stop()                    //停止连接，结束当前连接状态M
	IsClose() bool            //是否已关闭
	Context() context.Context //返回ctx，用于用户自定义的go协程获取连接退出状态

	GetTCPConnection() net.Conn //从当前连接获取原始的socket
	GetConnID() uint32          //获取当前连接ID
	RemoteAddr() net.Addr       //获取远程客户端地址信息
	GetDecodePipeLine() IChannelPipeLine
	GetEncodePipeLine() IChannelPipeLine
	SendMsg(msg INetMessage) error //发送消息到缓冲区

	IsHandShake() bool
	SetHandSign(sign byte)
	SetProperty(key string, value interface{})   //设置链接属性
	GetProperty(key string) (interface{}, error) //获取链接属性
	RemoveProperty(key string)                   //移除链接属性
}

//定义服务接口
type IServer interface {
	Start() //启动服务器方法
	Stop()  //停止服务器方法
	Serve() //开启业务服务方法
}

type IClient interface {
	Start()
	Stop()
}
