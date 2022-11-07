package main

import (
	"fmt"
	codec2 "github.com/kim118000/client/codec"
	"github.com/kim118000/core/pkg/codec"
	"github.com/kim118000/core/pkg/log"
	"github.com/kim118000/core/pkg/network"
	content2 "github.com/kim118000/protocol/proto/common/content"
	"github.com/kim118000/protocol/proto/gate"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)


var (
	FixLengthDecoder network.Decoder = codec.NewFixLengthDecoder(4)
	ClientHandShakeDecoder network.Decoder = codec.NewClientHandShakeDecoder()
	ClientMsgDispatcher network.Decoder = codec2.NewClientMessageDispatcher()
	AuthDecoder network.Decoder = codec2.NewAuthenticateHandler()
)

var (
	ProtobufEncoder network.Encoder = codec.NewProtoEncoder()
	FixLengthEncoder network.Encoder = codec.NewFixLengthEncoder(4)
)


func main() {
	addr := []network.ClusterAddr{
		{
			Name:      "kim",
			Addr:      "127.0.0.1:8999",
			Id:        1000,
			IsHand:    true,
			Reconnect: false,
		},
	}

	c := network.NewClient("client", addr, 1000, outProcessor(), inProcessor())
	c.Start()

	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	go stdin(c)

	time.Sleep(2 * time.Second)

	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	select {
	case signal := <-sg:
		c.Stop()
		log.DefaultLogger.Infof("got signal: %v, shutting down...", signal)
	}
}

func stdin(c *network.Client) {
	for {
		var id int
		var content string
		n, err := fmt.Scanln(&id, &content)
		if err != nil {
			log.DefaultLogger.Infof("%v", err)
		}

		if id == 1 {
			val := &gate.AuthenticationRequest{
				Sign:    "a",
				TokenTs: 1111,
				UserId:  111111111,
			}
			conn, _ := c.GetConnMgr().Get(1000)
			msg := network.GetMessage()
			msg.Init(1, val)
			conn.SendMsg(msg)
		}

		if id == 2 {
			val := &content2.ClientInboundMessage{
				MessageId: 1000,
				RequestId: 1000,
				Content: []byte("afdafafda"),
			}
			conn, _ := c.GetConnMgr().Get(1000)
			msg := network.GetMessage()
			msg.Init(1000, val)
			conn.SendMsg(msg)
		}

		log.DefaultLogger.Infof("%d id = %d c = %s", n, id, content)
	}
}

func inProcessor() []network.DecoderHandle {

	var args = []network.DecoderHandle{
		func() network.Decoder {
			return ClientHandShakeDecoder
		},
		func() network.Decoder {
			return FixLengthDecoder
		},
		func() network.Decoder {
			return AuthDecoder
		},
		func() network.Decoder {
			return ClientMsgDispatcher
		},
	}

	return args
}

func outProcessor() []network.EncoderHandle {
	var args = []network.EncoderHandle{
		func() network.Encoder {
			return ProtobufEncoder
		},
		func() network.Encoder {
			return FixLengthEncoder
		},
	}

	return args
}
