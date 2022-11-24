package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
	"unsafe"
)

type Options redis.Options

type RedisClient struct {
	*redis.Client
	*Options
}

type Option func(op *Options)

func WithPoolSize(size int) Option {
	return func(op *Options) {
		op.PoolSize = size
	}
}

func WithPassWord(pwd string) Option {
	return func(op *Options) {
		op.Password = pwd
	}
}

func WithDB(db int) Option {
	return func(op *Options) {
		op.DB = db
	}
}

func NewRedisClient(addr string, ops ...Option) *RedisClient {
	var option *Options = new(Options)

	for _, fn := range ops {
		fn(option)
	}

	option.Addr = addr

	c := redis.NewClient((*redis.Options)(unsafe.Pointer(option)))

	redis := &RedisClient{
		c,
		option,
	}

	return redis
}

func (rc *RedisClient) Set(key string, value interface{}, expiration time.Duration) {
	val, _ := json.Marshal(value)
	rc.Client.Set(context.Background(), key, val, expiration)
}

func (rc *RedisClient) Get(key string) string {
	return rc.Client.Get(context.Background(), key).Val()
}