package codec

import (
	"context"
	"github.com/kim118000/core/constant"
	"github.com/kim118000/core/pkg/log"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/protocol/proto/common/content"
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

	var response content.ClientOutboundMessage
	ex := proto.Unmarshal(message.GetData(), &response)
	if ex != nil {
		log.DefaultLogger.Errorf("client message proto unmarshal error %s", ex)
		return nil, msg, ex
	}

	log.DefaultLogger.Infof("===== %s", response)

	return dpc.GetNext(), msg, nil
}
