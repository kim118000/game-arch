package network

import (
	"google.golang.org/protobuf/proto"
	"sync"
)

//Message 消息
type Message struct {
	id      uint32      //消息的ID
	data    []byte      //消息的内容
	dataLen uint32      //消息的长度
	serialData  proto.Message //序列化对象
}

var pool = sync.Pool{New: func() interface{} { return new(Message) }}

func GetMessage() *Message {
	return pool.Get().(*Message)
}

func PutMessage(msg *Message) {
	msg.serialData = nil
	msg.data = nil
	msg.id = 0
	msg.dataLen = 0
	pool.Put(msg)
}

func (msg *Message) Init(msgId uint32, serial proto.Message) {
	msg.id = msgId
	msg.serialData = serial
}

func (msg *Message) GetSerialData() proto.Message {
	return msg.serialData
}

//GetDataLen 获取消息数据段长度
func (msg *Message) GetDataLen() uint32 {
	return msg.dataLen
}

//GetMsgID 获取消息ID
func (msg *Message) GetMsgID() uint32 {
	return msg.id
}

//GetData 获取消息内容
func (msg *Message) GetData() []byte {
	return msg.data
}

//SetData 设计消息内容
func (msg *Message) SetData(data []byte) {
	msg.data = data
	msg.dataLen = uint32(len(data))
}
