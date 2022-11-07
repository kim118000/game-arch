package codec

import (
	"context"
	"encoding/binary"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/core/pkg/pool/byteslice"
)

type FixLengthEncoder struct {
	len int
}

func NewFixLengthEncoder(len int) *FixLengthEncoder {
	return &FixLengthEncoder{
		len: len,
	}
}

func (f *FixLengthEncoder) Encode(ctx context.Context, dpc *network.DefaultPipeLineContext, msg network.INetMessage) (*network.DefaultPipeLineContext, network.INetMessage, error) {
	buff := byteslice.Get(f.len + int(msg.GetDataLen()))
	binary.LittleEndian.PutUint32(buff, msg.GetDataLen())
	body := msg.GetData()
	copy(buff[f.len:], body)
	msg.SetData(buff)
	byteslice.Put(body)
	return dpc.GetNext(), msg, nil
}
