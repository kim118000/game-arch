package network

import (
	"context"
	"errors"
	"fmt"
	logger2 "github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/pool/byteslice"
	"github.com/kim118000/core/toolkit/queue"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

//Connection 链接
type Connection struct {
	//当前连接的socket TCP套接字
	conn net.Conn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	connID uint32

	//读数据流水线
	encode IChannelPipeLine
	decode IChannelPipeLine

	//连接事件
	connEvent IConnEvent

	//告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc

	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan INetMessage
	//暂存并发队列，丢弃最老策略
	queue queue.AsyncQueue
	//当前连接的关闭状态
	isClosed uint32
	//关闭连接
	once sync.Once

	//链接属性
	property map[string]interface{}
	////保护当前property的锁
	propertyLock sync.RWMutex

	strCache string
	//是否握手
	handshake bool
	//签名
	sign byte
}

//NewConnection 创建连接的方法
func NewConnection(conn net.Conn, connID uint32, in []EncoderHandle, out []DecoderHandle, sendMsgChanLen uint32, connevent IConnEvent) *Connection {
	//初始化Conn属性
	var c = &Connection{
		conn:        conn,
		connID:      connID,
		isClosed:    0,
		encode:      NewDefaultCodecPipeLine(),
		decode:      NewDefaultCodecPipeLine(),
		msgBuffChan: make(chan INetMessage, sendMsgChanLen),
		queue:       queue.NewLockFreeQueue(int32(sendMsgChanLen)),
		connEvent:   connevent,
		property:    nil,
	}
	if c.connEvent != nil {
		c.connEvent.OnCreate(c)
	}

	if in != nil {
		for _, proc := range in{
			c.encode.PushBack(NewDefaultEncodeContext(c.encode, proc()))
		}
	}

	if out != nil {
		for _, proc := range out{
			c.decode.PushBack(NewDefaultEncodeContext(c.decode, proc()))
		}
	}

	c.strCache = fmt.Sprintf("[connId=%d,remoteaddr=%s,localaddr=%s]", c.GetConnID(), c.conn.RemoteAddr().String(), c.conn.LocalAddr().String())
	return c
}

//StartWriter 写消息Goroutine， 用户将数据发送给客户端
func (c *Connection) StartWriter() {
	logger2.DefaultLogger.Infof("%s [writer goroutine is running]", c)
	defer logger2.DefaultLogger.Infof("%s [conn writer exit!]", c)
	defer c.Stop()

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				logger2.DefaultLogger.Infof("%s send data msgid=%d datalength=%d", c, data.GetMsgID(), data.GetDataLen())
				if _, err := c.conn.Write(data.GetData()); err != nil {
					c.reliveMessage(data)
					logger2.DefaultLogger.Errorf("%s send data error %v", c, err)
					return
				}
				c.reliveMessage(data)
			} else {
				logger2.DefaultLogger.Infof("%s MsgBuffChan is Closed", c)
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Connection) reliveMessage(message INetMessage) {
	byteslice.Put(message.GetData())
	PutMessage(message.(*Message))
}

//StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *Connection) StartReader() {
	logger2.DefaultLogger.Infof("%s [reader goroutine is running]", c)
	defer logger2.DefaultLogger.Infof("%s [conn reader exit!]", c)
	defer c.Stop()

	//读
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			//流水线读
			msg, err := c.decode.Decode(c.ctx, c, nil)
			message, ok := msg.(*Message)
			if ok {
				c.reliveMessage(message)
			}
			if err != nil{
				logger2.DefaultLogger.Debugf("%s decode process error %v", c, err)
				return
			}
		}
	}
}

//Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	if c.connEvent != nil {
		c.connEvent.OnConnStart(c)
	}
}

//Stop 停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	c.once.Do(func() {
		c.cancel()
		c.finalizer()
	})
}

func (c *Connection) IsClose() bool {
	isclose := atomic.LoadUint32(&c.isClosed)
	if isclose == 1 {
		return true
	}
	return false
}

func (c *Connection) finalizer() {
	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	if c.connEvent != nil {
		c.connEvent.OnConnStop(c)
	}

	ok := atomic.CompareAndSwapUint32(&c.isClosed, 0, 1)
	if !ok {
		return
	}

	// 关闭socket链接
	_ = c.conn.Close()

	//将链接从连接管理器中删除
	if c.connEvent != nil {
		c.connEvent.OnClose(c)
	}

	c.queue.Clear()
	//关闭该链接全部管道
	close(c.msgBuffChan)

	logger2.DefaultLogger.Infof("%s finalizer successfully", c)
}

//GetTCPConnection 从当前连接获取原始的socket
func (c *Connection) GetTCPConnection() net.Conn {
	return c.conn
}

//GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.connID
}

//RemoteAddr 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

//SendMsg
func (c *Connection) SendMsg(msg INetMessage) error {
	if c.IsClose() {
		return fmt.Errorf("connection closed when send message %s", c)
	}

	for {
		if c.queue.IsEmpty() {
			break
		}
		val := c.queue.Dequeue()
		msg, ok := val.Elem.(*Message)
		queue.PutQueueElement(val)

		if ok {
			select {
			case c.msgBuffChan <- msg:
			default:
				goto g
			}
		}
	}
g:
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()

	data, err := c.encode.Encode(c.ctx, msg)
	if err != nil {
		logger2.DefaultLogger.Errorf("%s encode msg error %v", c, err)
		return fmt.Errorf("%s encode msg error", c)
	}

	// 发送超时
	select {
	case <-ctx.Done():
		//暂存
		val := queue.GetQueueElement()
		val.Elem = data
		c.queue.Enqueue(val)

		logger2.DefaultLogger.Errorf("%s write Message timeout queue size=%d", c, c.queue.Len())
		return errors.New("write message timeout")
	case c.msgBuffChan <- data:
		return nil
	}

	return nil
}

//SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

//GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

//RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

//返回ctx，用于用户自定义的go程获取连接退出状态
func (c *Connection) Context() context.Context {
	return c.ctx
}

func (c *Connection) String() string {
	return c.strCache
}

func (c *Connection) IsHandShake() bool {
	return c.handshake
}

func (c *Connection) SetHandSign(sign byte){
	c.handshake = true
	c.sign = sign
}

func (c *Connection) GetDecodePipeLine() IChannelPipeLine {
	return c.decode
}

func (c *Connection) GetEncodePipeLine() IChannelPipeLine {
	return c.encode
}