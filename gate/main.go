package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/kim118000/core/pkg/config"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/core/toolkit/file"
	"github.com/kim118000/gate/asset/arena_asset"
	"github.com/kim118000/gate/gconfig"
	"github.com/kim118000/gate/logger"
	"github.com/kim118000/gate/mgr"
	"github.com/kim118000/gate/msg/handler"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	conf := file.PathJoin("conf/conf.toml")

	var serverConfig *gconfig.ServerConfig
	if _, err := toml.DecodeFile(conf, &serverConfig); err != nil {
		fmt.Println(err)
		return
	}

	logger.InitLogger(&serverConfig.LogConfig)

	loader := config.NewFileLoader(file.PathJoin("json"))

	manager := config.NewConfigManager(loader)
	manager.RegTemplate(arena_asset.ArenaTemplate)
	manager.LoadTemplate()

	server := network.NewServer("gate", "", 8999, 1000, 1000, outProcessor(), inProcessor(), network.WithConnEvent(mgr.ConnEvent))

	server.Start()

	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	// stop server
	select {
	case signal := <-sg:
		server.Stop()
		logger.Info.Infof("got signal: %v, shutting down...", signal)
	}

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
