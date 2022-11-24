package mock

import (
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/network"
	content2 "github.com/kim118000/protocol/proto/common/content"
	"github.com/kim118000/protocol/proto/common/id"
	"google.golang.org/protobuf/proto"
	"reflect"
	"strings"
	"sync/atomic"
)

type IMock interface {
	Request(c *network.Client, connId uint32, content string)
	PrintJson() string
}

var Mocks = make(map[uint32]IMock)
var counter int32

func init() {
	Mocks[AuthMock.GetMsgId()] = AuthMock
	Mocks[LoginMock.GetMsgId()] = LoginMock
}

func SendMsg(c *network.Client, connId uint32, msgId id.MessageId, m proto.Message) {
	content, err := proto.Marshal(m)
	if err != nil {
		logger.Log.Infof("proto marshal %v", err)
		return
	}

	requestId := atomic.AddInt32(&counter, 1)
	val := &content2.ClientInboundMessage{
		MessageId: msgId,
		RequestId: requestId,
		Content:   content,
	}
	conn, _ := c.GetConnMgr().Get(connId)
	if conn == nil {
		logger.Log.Errorf("conn not found %d", connId)
		return
	}
	msg := network.GetMessage()
	msg.Init(uint32(msgId.Number()), val)
	conn.SendMsg(msg)
}

func GetHelpJson(msg proto.Message) string {
	var strs strings.Builder
	strs.WriteString("{")
	req := reflect.TypeOf(msg).Elem()
	var numField = req.NumField()
	for i := 0; i < numField; i++ {
		name := req.Field(i).Name
		if name == "state" || name == "unknownFields" || name == "sizeCache" {
			continue
		}
		tag := req.Field(i).Tag.Get("protobuf")
		start := strings.LastIndex(tag,"=") + 1
		end := strings.LastIndex(tag, ",")
		tag = tag[start:end]
		if i == numField-1 {
			strs.WriteString("\"" + tag + "\"" + ":" + "\"\"")
			break
		}
		strs.WriteString("\"" + tag + "\"" + ":" + "\"\",")
	}

	strs.WriteString("}")
	return strs.String()
}
