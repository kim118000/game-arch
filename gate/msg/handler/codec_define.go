package handler

import (
	"github.com/kim118000/core/pkg/codec"
	"github.com/kim118000/core/pkg/network"
)

var (
	FixLengthDecoder network.Decoder = codec.NewFixLengthDecoder(4)
	AuthDecoder network.Decoder = NewAuthenticateHandler()

	ClientHandShakeDecoder network.Decoder = codec.NewClientHandShakeDecoder()
	ServerHandShakeDecoder network.Decoder = NewHandShakeHandler()

	ClientMsgDispatcher network.Decoder = NewClientMessageDispatcher()
	ClusterMsgDispatcher network.Decoder = NewClusterMessageDispatcher()
)

var (
	ProtobufEncoder network.Encoder = codec.NewProtoEncoder()
	FixLengthEncoder network.Encoder = codec.NewFixLengthEncoder(4)
)
