package service

import (
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/redis"
	"github.com/kim118000/core/toolkit/snowflake"
	"github.com/kim118000/db/conn"
	"github.com/kim118000/db/persistence"
	"github.com/kim118000/login/internal/conf"
)

var LoginService = new(service)
var IdSnow *snowflake.Node

type service struct{}

func (s *service) Init(conf *conf.ServerConfig) {
	var err error
	IdSnow, err = snowflake.NewNode(int64(conf.ServerId))
	if err != nil {
		panic("error creating snowflake node")
	}

	redis.InitRedis(&conf.RedisConfig)
	logger.InitLogger(&conf.LogConfig)
	persistence.InitSysDao(conn.InitDbConn(&conf.DbConfig))
}

func (s *service) Start() {

}

func (s *service) Stop() {
	redis.CloseRedis()
}
