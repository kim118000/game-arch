package codec

import (
	"context"
	"encoding/binary"
	"github.com/kim118000/core/constant"
	"github.com/kim118000/core/pkg/log"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/core/pkg/pool/byteslice"
	"io"
)

type FixLengthDecoder struct {
	len int
}

func NewFixLengthDecoder(len int) *FixLengthDecoder {
	return &FixLengthDecoder{
		len: len,
	}
}

func (f *FixLengthDecoder) Decode(ctx context.Context, conn network.IConnection, dpc *network.DefaultPipeLineContext, msg interface{}) (*network.DefaultPipeLineContext, interface{}, error) {
	headData := byteslice.Get(f.len)
	defer byteslice.Put(headData)
	if _, err := io.ReadFull(conn.GetTCPConnection(), headData); err != nil {
		return nil, nil, err
	}

	msgLen := binary.LittleEndian.Uint32(headData)
	log.DefaultLogger.Debugf("conn[%s] read message length=%d", conn, msgLen)

	if msgLen < 1 || msgLen > 1024*8 {
		return nil, nil, constant.ErrWrongDatapackLength
	}

	data := byteslice.Get(int(msgLen))
	if _, err := io.ReadFull(conn.GetTCPConnection(), data); err != nil {
		byteslice.Put(data)
		log.DefaultLogger.Errorf("conn[%s] read message body error %s", conn, err.Error())
		return nil, nil, err
	}

	message := network.GetMessage()
	message.SetData(data)

	return dpc.GetNext(), message, nil
}
