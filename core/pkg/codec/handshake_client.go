package codec

import (
	"context"
	"github.com/kim118000/core/internal/constant"
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/core/pkg/pool/byteslice"
	"io"
)

type ClientHandShakeDecoder struct {
}

func NewClientHandShakeDecoder() *ClientHandShakeDecoder {
	return &ClientHandShakeDecoder{
	}
}

func (chs ClientHandShakeDecoder) Decode(ctx context.Context, conn network.IConnection, dpc *network.DefaultPipeLineContext, msg interface{}) (*network.DefaultPipeLineContext, interface{}, error) {
	buf := byteslice.Get(constant.HAND_SHAKE_SIGN_LEN)
	defer byteslice.Put(buf)

	_, err := io.ReadFull(conn.GetTCPConnection(), buf)
	if err != nil {
		return nil, nil, err
	}

	conn.SetHandSign(buf[3])
	conn.GetDecodePipeLine().Remove(dpc)

	logger.DefaultLogger.Infof("%s hand shake success", conn)
	return nil, nil, nil
}