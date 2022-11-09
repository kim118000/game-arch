package mock

import (
	"fmt"
	"github.com/kim118000/client/internal/service"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/protocol/proto/common/id"
	"github.com/kim118000/protocol/proto/game"
	"google.golang.org/protobuf/encoding/protojson"
)

var LoginMock = &login{
	MsgId: id.MessageId_TO_SERVER_LOGIN,
}

type login struct {
	MsgId id.MessageId
}

func (m *login) GetMsgId() uint32 {
	return uint32(m.MsgId.Number())
}

func (m *login) PrintJson() string {
	strs := GetHelpJson(&game.LoginRequest{})
	return fmt.Sprintf("%s", strs)
}

func (m *login) Request(c *network.Client, connId uint32, content string) {
	if content == "help" {
		service.Log.Infof("%s", m.PrintJson())
		return
	}

	var request = &game.LoginRequest{}
	err := protojson.Unmarshal([]byte(content), request)
	if err != nil {
		service.Log.Errorf("%v", err)
	}

	SendMsg(c, connId, m.MsgId, request)
}
