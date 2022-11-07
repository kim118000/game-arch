package codec

import (
	"context"
	"github.com/kim118000/core/pkg/log"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/core/pkg/pool/byteslice"
	"google.golang.org/protobuf/proto"
)

type ProtoBufEncoder struct {
}

func NewProtoEncoder() *ProtoBufEncoder {
	return &ProtoBufEncoder{
	}
}

func (pb *ProtoBufEncoder) Encode(ctx context.Context, dpc *network.DefaultPipeLineContext, msg network.INetMessage) (*network.DefaultPipeLineContext, network.INetMessage, error) {
	if msg.GetData() != nil && len(msg.GetData()) > 0 {
		return nil, msg, nil
	}

	protoMsg := msg.GetSerialData()
	size := proto.Size(protoMsg)
	body := byteslice.Get(size)
	var prototools proto.MarshalOptions
	_, err := prototools.MarshalAppend(body[:0], protoMsg)

	if err != nil {
		byteslice.Put(body)
		return nil, msg, err
	}

	msg.SetData(body)
	log.DefaultLogger.Debugf("encode protobuf message %s", protoMsg)

	return dpc.GetNext(), msg, nil
}
