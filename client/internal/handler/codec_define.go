package handler

import (
	"github.com/kim118000/core/pkg/codec"
	"github.com/kim118000/core/pkg/network"
)

var (
	FixLengthDecoder network.Decoder = codec.NewFixLengthDecoder(4)
	ClientHandShakeDecoder network.Decoder = codec.NewClientHandShakeDecoder()
	ClientMsgDispatcher network.Decoder = NewClientMessageDispatcher()
	AuthDecoder network.Decoder = NewAuthenticateHandler()
)

var (
	ProtobufEncoder network.Encoder = codec.NewProtoEncoder()
	FixLengthEncoder network.Encoder = codec.NewFixLengthEncoder(4)
)
