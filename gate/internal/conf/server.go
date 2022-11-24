package conf

import (
	"github.com/kim118000/core/pkg/config"
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/redis"
)

var Config = new(ServerConfig)

type ServerConfig struct {
	config.ServerBase
	MaxConn uint32
	SendBuffLength uint32
	LogConfig   logger.LogConfig  `toml:"LogConfig"`
	RedisConfig redis.RedisConfig `toml:"RedisConfig"`
}
