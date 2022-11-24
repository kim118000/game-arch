package redis

import "sync"

var (
	onceRedis sync.Once
	Client     *RedisClient
)

func InitRedis(conf *RedisConfig) {
	onceRedis.Do(func() {
		op := []Option{
			WithPassWord(conf.Pwd),
			WithDB(conf.DB),
		}
		Client = NewRedisClient(conf.Addr, op...)
	})
}

func CloseRedis()  {
	Client.Close()
}