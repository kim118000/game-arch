package codec

import (
	"bufio"
	"context"
	"github.com/kim118000/core/pkg/network"
)

type LineBasedFrameDecoder struct {
	delimiter byte
}

func NewLineBasedFrameDecoder(delim byte) *LineBasedFrameDecoder {
	return &LineBasedFrameDecoder{
		delimiter: delim,
	}
}

func (lfd *LineBasedFrameDecoder) Decode(ctx context.Context, conn network.IConnection, dpc *network.DefaultPipeLineContext, msg interface{}) (*network.DefaultPipeLineContext, interface{}, error) {
	bufReader := bufio.NewReader(conn.GetTCPConnection())
	buff, readErr := bufReader.ReadBytes(lfd.delimiter)
	if readErr != nil {
		return nil, nil, readErr
	}

	message := network.GetMessage()
	message.SetData(buff)

	return dpc.GetNext(), message, nil
}
