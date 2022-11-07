package codec

import (
	"context"
	"github.com/kim118000/core/constant"
	"github.com/kim118000/core/pkg/log"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/core/pkg/pool/byteslice"
	"github.com/kim118000/protocol/proto/gate"
	"google.golang.org/protobuf/proto"
)

type AuthenticateHandler struct {
}

func NewAuthenticateHandler() *AuthenticateHandler {
	return &AuthenticateHandler{
	}
}

func (l *AuthenticateHandler) Decode(ctx context.Context, conn network.IConnection, dpc *network.DefaultPipeLineContext, msg interface{}) (*network.DefaultPipeLineContext, interface{}, error) {
	message, ok := msg.(*network.Message)
	if !ok {
		return nil, nil, constant.ErrWrongMessageAssert
	}

	var response gate.AuthenticationResponse
	ex := proto.Unmarshal(message.GetData(), &response)
	if ex != nil {
		log.DefaultLogger.Errorf("auth proto unmarshal error %s", ex)
		return nil, msg, ex
	}

	if !response.Ok {
		return nil, msg, constant.ErrWrongAuthFailure
	}

	log.DefaultLogger.Infof("auth ok")

	byteslice.Put(message.GetData())
	network.PutMessage(message)

	conn.GetDecodePipeLine().Insert(dpc, NewClientMessageDispatcher())
	conn.GetDecodePipeLine().Remove(dpc)
	return nil, msg, nil
}
