package conf

import (
	"github.com/kim118000/core/pkg/config"
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/redis"
	"github.com/kim118000/db/conn"
)

var SC = new(ServerConfig)

type ServerConfig struct {
	config.ServerBase
	LogConfig   logger.LogConfig  `toml:"LogConfig"`
	RedisConfig redis.RedisConfig `toml:"RedisConfig"`
	DbConfig    conn.DbConfig     `toml:"DbConfig"`
}
