package mock

import (
	"fmt"
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/protocol/proto/gate"
	"google.golang.org/protobuf/encoding/protojson"
)

var AuthMock = &authMock{
	MsgId: 0,
}

type authMock struct {
	MsgId uint32
}

func (m *authMock) GetMsgId() uint32 {
	return m.MsgId
}

func (m *authMock) PrintJson() string {
	strs := GetHelpJson(&gate.AuthenticationRequest{})
	return fmt.Sprintf("%s", strs)
}

func (m *authMock) Request(c *network.Client, connId uint32, content string) {
	if content == "help" {
		logger.Log.Infof("%s", m.PrintJson())
		return
	}

	var request = &gate.AuthenticationRequest{}
	err := protojson.Unmarshal([]byte(content), request)
	if err != nil {
		logger.Log.Errorf("%v", err)
	}

	conn, _ := c.GetConnMgr().Get(connId)
	if conn == nil {
		logger.Log.Errorf("conn not found %d", connId)
		return
	}
	msg := network.GetMessage()
	msg.Init(m.MsgId, request)
	_ = conn.SendMsg(msg)
}
