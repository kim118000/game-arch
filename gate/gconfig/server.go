package gconfig

import (
	"github.com/kim118000/core/pkg/config"
	"github.com/kim118000/core/pkg/log"
)

type ServerConfig struct {
	config.ServerBase
	LogConfig log.LogConfig  `toml:"LogConfig"`
}
