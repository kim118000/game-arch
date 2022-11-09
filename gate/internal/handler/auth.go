package handler

import (
	"context"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/core/pkg/pool/byteslice"
	"github.com/kim118000/gate/internal/constant"
	"github.com/kim118000/gate/internal/service"
	"github.com/kim118000/gate/internal/session"
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

	var request gate.AuthenticationRequest
	ex := proto.Unmarshal(message.GetData(), &request)
	if ex != nil {
		service.GateLog.Errorf("auth proto unmarshal error %s", ex)
		return nil, msg, ex
	}

	val, err := conn.GetProperty(constant.SessionAttrKey)
	if err != nil {
		return nil, msg, ex
	}

	sess, _ := val.(*session.Session)
	ok = sess.Auth(request.UserId, int(request.TokenTs), request.Sign)
	if !ok {
		return nil, msg, constant.ErrWrongAuthFailure
	}

	//发送认证成功消息
	var response = &gate.AuthenticationResponse{
		Ok: true,
	}

	sess.SendProtobuf(1, response)

	byteslice.Put(message.GetData())
	network.PutMessage(message)

	conn.GetDecodePipeLine().Insert(dpc, ClientMsgDispatcher)
	conn.GetDecodePipeLine().Remove(dpc)
	return nil, msg, nil
}
