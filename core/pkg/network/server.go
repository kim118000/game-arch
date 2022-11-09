package network

import (
	"errors"
	"fmt"
	logger2 "github.com/kim118000/core/pkg/logger"
	"net"
	"sync/atomic"
	"time"
)

var counter uint32

//Server 接口实现，定义一个Server服务类
type Server struct {
	//服务器的名称
	name string
	//服务绑定的IP地址
	ip string
	//服务绑定的端口
	port uint16

	//最大连接
	maxConn uint32
	//连接发送的最大chan大小
	connSendMaxBuffLen uint32

	//读数据流水线
	encode []EncoderHandle
	decode []DecoderHandle

	//连接事件
	connEvent IConnEvent
	exitChan  chan struct{}
}

//NewServer 创建一个服务器句柄
func NewServer(name string, ip string, port uint16, maxconn uint32, sendbufflen uint32, encode []EncoderHandle, decode []DecoderHandle, opts ...ServerOption) IServer {

	s := &Server{
		name:               name,
		ip:                 ip,
		port:               port,
		maxConn:            maxconn,
		connSendMaxBuffLen: sendbufflen,
		encode:             encode,
		decode:             decode,
		exitChan:           nil,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) init(conn *net.TCPConn) {
	conn.SetNoDelay(true)
}

//Start 开启网络服务
func (s *Server) Start() {
	logger2.DefaultLogger.Infof("start server name=%s listenner at [ip:%s,port:%d] is starting", s.name, s.ip, s.port)
	s.exitChan = make(chan struct{})

	//开启一个go去做服务端Linster业务
	go func() {
		// 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
		if err != nil {
			logger2.DefaultLogger.Fatalf("resolve tcp addr err: %v", err)
			return
		}

		// 监听服务器地址
		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			logger2.DefaultLogger.Fatalf("resolve tcp addr err: %v", err)
			return
		}

		go func() {
			var tempDelay time.Duration // how long to sleep on accept failure
			// 启动server网络连接业务
			for {
				conn, err := listener.AcceptTCP()
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						logger2.DefaultLogger.Infof("server listener closed %s", s.name)
						return
					}
					if ne, ok := err.(net.Error); ok && ne.Temporary() {
						if tempDelay == 0 {
							tempDelay = 5 * time.Millisecond
						} else {
							tempDelay *= 2
						}
						if max := 1 * time.Second; tempDelay > max {
							tempDelay = max
						}
						logger2.DefaultLogger.Infof("http: Accept error: %v; retrying in %v", err, tempDelay)
						time.Sleep(tempDelay)
						continue
					}
					return
				}

				connId := atomic.AddUint32(&counter, 1)
				//设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
				if connId >= s.maxConn {
					conn.Close()
					continue
				}

				//初始化参数
				s.init(conn)
				tempDelay = 0
				dealConn := NewConnection(conn, connId, s.encode, s.decode, s.connSendMaxBuffLen, s.connEvent)
				go dealConn.Start()
			}
		}()

		select {
		case <-s.exitChan:
			err := listener.Close()
			if err != nil {
				logger2.DefaultLogger.Errorf("server listener close error %v", err)
			}
		}
	}()
}

//Stop 停止服务
func (s *Server) Stop() {
	logger2.DefaultLogger.Infof("stop server name %s", s.name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.exitChan <- struct{}{}
	close(s.exitChan)

	//清理连接
	s.connEvent.GetConnMgr().ClearConn()
}

//Serve 运行服务
func (s *Server) Serve() {
	s.Start()
}
