package service

import (
	"github.com/kim118000/core/pkg/redis"
	"github.com/kim118000/gate/internal/conf"
	"sync"
)

var (
	onceRedis     sync.Once
	Redis    *redis.RedisClient
)

func initRedis(conf *conf.ServerConfig) {
	onceRedis.Do(func() {
		op := []redis.Option{
			redis.WithPassWord(conf.RedisConfig.Pwd),
			redis.WithDB(conf.RedisConfig.DB),
		}
		Redis = redis.NewRedisClient(conf.RedisConfig.Addr, op...)
	})
}

