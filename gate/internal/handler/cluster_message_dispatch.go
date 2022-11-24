package handler

import (
	"context"
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/gate/internal/session"
	"github.com/kim118000/protocol/proto/server/cluster"
	"google.golang.org/protobuf/proto"
)

type ClusterMessageDispatcher struct {
}

func NewClusterMessageDispatcher() *ClusterMessageDispatcher {
	return &ClusterMessageDispatcher{
	}
}

func (l *ClusterMessageDispatcher) Decode(ctx context.Context, conn network.IConnection, dpc *network.DefaultPipeLineContext, msg interface{}) (*network.DefaultPipeLineContext, interface{}, error) {
	message, ok := msg.(*network.Message)
	if !ok {
		logger.Log.Errorf("cluster message assert failure")
		return nil, nil, nil
	}

	var request cluster.RPCMessage
	ex := proto.Unmarshal(message.GetData(), &request)
	if ex != nil {
		logger.Log.Errorf("cluster proto unmarshal error %s", ex)
		return dpc.GetNext(), msg, nil
	}

	switch request.MessageCategory {
	case cluster.MessageCategory_MESSAGE_CATEGORY_CLIENT_OUT:
		sess, err := session.SessionMgr.Get(request.RoleId)
		if err != nil {
			logger.Log.Errorf("cluster message dispatch not found user[%d] session", request.RoleId)
		} else {
			sess.Send(nil)
		}
	}

	return dpc.GetNext(), msg, nil
}
