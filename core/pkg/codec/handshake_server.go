package codec

import (
	"context"
	"github.com/kim118000/core/constant"
	"github.com/kim118000/core/pkg/log"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/core/pkg/pool/byteslice"
	"github.com/kim118000/core/toolkit"
	"io"
)

type ServerHandShakeDecoder struct {
}

func NewServerHandShakeDecoder() *ServerHandShakeDecoder {
	return &ServerHandShakeDecoder{
	}
}

func (shs *ServerHandShakeDecoder) Decode(ctx context.Context, conn network.IConnection, dpc *network.DefaultPipeLineContext, msg interface{}) (*network.DefaultPipeLineContext, interface{}, error) {
	buf := byteslice.Get(constant.HAND_SHAKE_ARR_LEN)
	defer byteslice.Put(buf)

	_, err := io.ReadFull(conn.GetTCPConnection(), buf)
	if err != nil {
		return nil, nil, err
	}

	suc := toolkit.CheckHandShakeByte(buf)
	if suc {
		sign := toolkit.Rand.NextBytes(constant.HAND_SHAKE_SIGN_LEN)
		conn.SetHandSign(sign[3])

		msg := network.GetMessage()
		msg.SetData(sign)
		conn.SendMsg(msg)

		conn.GetDecodePipeLine().Remove(dpc)

		log.DefaultLogger.Infof("%s hand shake success", conn)
		return nil, nil, nil
	}

	return nil, nil, constant.ErrWrongConnShake
}
