package handler

import (
	"context"
	"github.com/kim118000/core/pkg/log"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/gate/constant"
	"github.com/kim118000/gate/logger"
	"github.com/kim118000/gate/session"
	"github.com/kim118000/protocol/proto/common/content"
	"github.com/kim118000/protocol/proto/server/cluster"
	"google.golang.org/protobuf/proto"
)

type ClientMessageDispatcher struct {
}

func NewClientMessageDispatcher() *ClientMessageDispatcher {
	return &ClientMessageDispatcher{
	}
}

func (l *ClientMessageDispatcher) Decode(ctx context.Context, conn network.IConnection, dpc *network.DefaultPipeLineContext, msg interface{}) (*network.DefaultPipeLineContext, interface{}, error) {
	message, ok := msg.(*network.Message)
	if !ok {
		return nil, nil, constant.ErrWrongMessageAssert
	}

	var request content.ClientInboundMessage
	ex := proto.Unmarshal(message.GetData(), &request)
	if ex != nil {
		log.DefaultLogger.Errorf("client message dispatch proto unmarshal error %s", ex)
		return nil, msg, ex
	}

	var endPointType cluster.EndPointType
	if request.MessageId > 1 && request.MessageId < 30000 {
		endPointType = cluster.EndPointType_ENDPOINT_GAME
	}

	val, err := conn.GetProperty(constant.SessionAttrKey)
	if err != nil {
		log.DefaultLogger.Errorf("client message dispatch session not found [%s]", conn)
		return nil, msg, ex
	}

	sess, _ := val.(*session.Session)

	switch endPointType {
	case cluster.EndPointType_ENDPOINT_GAME:
		gameNode := sess.GetGameNode()
		if gameNode != nil {

		} else {
			var rpcMsg = new(cluster.RPCMessage)
			rpcMsg.MessageCategory = cluster.MessageCategory_MESSAGE_CATEGORY_CLIENT_IN
			rpcMsg.MessageId = request.MessageId
			rpcMsg.ClientRequestId = request.RequestId
			rpcMsg.UserId = sess.GetUserId()
			rpcMsg.Content = request.Content

			logger.Info.Debugf("=========%s", rpcMsg)
			//data := network.GetMessage()
			//data.Init(uint32(request.MessageId.Number()), rpcMsg)
			//gameNode.SendMsg(data)
		}
	}

	return dpc.GetNext(), msg, nil
}
