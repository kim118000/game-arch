package main

import (
	"github.com/BurntSushi/toml"
	"github.com/kim118000/core/pkg/config"
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/core/toolkit"
	"github.com/kim118000/core/toolkit/file"
	"github.com/kim118000/gate/internal/asset/arena_asset"
	"github.com/kim118000/gate/internal/conf"
	"github.com/kim118000/gate/internal/handler"
	"github.com/kim118000/gate/internal/service"
)

func main() {
	confFile := file.PathJoin("conf/conf.toml")

	if _, err := toml.DecodeFile(confFile, conf.Config); err != nil {
		logger.Log.Errorf("load server config err %s", err)
		return
	}

	service.InitService(conf.Config)
	service.StartService()

	loader := config.NewFileLoader(file.PathJoin("json"))
	manager := config.NewConfigManager(loader)
	manager.RegTemplate(arena_asset.ArenaTemplate)
	manager.LoadTemplate()

	server := network.NewServer(conf.Config.ServerName, "", uint16(conf.Config.ServerPort), conf.Config.MaxConn, conf.Config.SendBuffLength, outProcessor(), inProcessor(), network.WithConnEvent(service.ConnEvent))
	server.Start()

	toolkit.RegisterSignal(func() {
		service.StopService()
	})
}

func inProcessor() []network.DecoderHandle {
	var args = []network.DecoderHandle{
		func() network.Decoder {
			return handler.ServerHandShakeDecoder
		},
		func() network.Decoder {
			return handler.FixLengthDecoder
		},
		func() network.Decoder {
			return handler.AuthDecoder
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
