package service

import (
	"github.com/kim118000/client/internal/conf"
	"github.com/kim118000/core/pkg/redis"
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

