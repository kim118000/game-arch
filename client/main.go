package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/kim118000/client/internal/conf"
	"github.com/kim118000/client/internal/handler"
	"github.com/kim118000/client/internal/service"
	"github.com/kim118000/client/mock"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/core/toolkit"
	"github.com/kim118000/core/toolkit/file"
	_ "net/http/pprof"
)

func main() {
	//auth val := `{"userId":111,"tokenTs":111,"sign":""}`
	//login {"userId":111}


	confFile := file.PathJoin("conf/conf.toml")
	if _, err := toml.DecodeFile(confFile, conf.Config); err != nil {
		fmt.Println(err)
		return
	}
	service.InitService(conf.Config)

	c := network.NewClient("client", conf.Config.ClusterAddr, 1000, outProcessor(), inProcessor())
	c.Start()

	go stdin(c)

	toolkit.RegisterSignal(func() {
		service.StopService()
	})
}

func stdin(c *network.Client) {
	for {
		var id uint32
		var content string
		_, err := fmt.Scanln(&id, &content)
		if err != nil {
			continue
		}

		mock, ok := mock.Mocks[id]
		if !ok {
			service.Log.Errorf("not found mock")
			continue
		}

		mock.Request(c, 1000, content)

	}
}

func inProcessor() []network.DecoderHandle {

	var args = []network.DecoderHandle{
		func() network.Decoder {
			return handler.ClientHandShakeDecoder
		},
		func() network.Decoder {
			return handler.FixLengthDecoder
		},
		func() network.Decoder {
			return handler.AuthDecoder
		},
		func() network.Decoder {
			return handler.ClientMsgDispatcher
		},
	}

	return args
}

func outProcessor() []network.EncoderHandle {
	var args = []network.EncoderHandle{
		func() network.Encoder {
			return handler.ProtobufEncoder
		},
		func() network.Encoder {
			return handler.FixLengthEncoder
		},
	}

	return args
}
