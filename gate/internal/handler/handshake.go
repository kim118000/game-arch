package handler

import (
	"context"
	"github.com/kim118000/core/pkg/codec"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/core/pkg/scheduler"
	"github.com/kim118000/gate/internal/constant"
	"github.com/kim118000/gate/internal/session"
)

type HandShakeHandler struct {
	codec.ServerHandShakeDecoder
}

func NewHandShakeHandler() *HandShakeHandler {
	return &HandShakeHandler{
	}
}

func (h *HandShakeHandler) Decode(ctx context.Context, conn network.IConnection, dpc *network.DefaultPipeLineContext, msg interface{}) (*network.DefaultPipeLineContext, interface{}, error) {
	next, msg, err := h.ServerHandShakeDecoder.Decode(ctx, conn, dpc, msg)
	if err != nil {
		return next, msg, err
	}

	sess := session.NewSession(conn, session.SessionMgr)
	scheduler.NewAfterTimerBySecondOnce(constant.AuthenticateTimeout, func() {
		if !sess.IsAuthSuccess() {
			sess.Kick()
		}
	})

	return next, msg, err
}
