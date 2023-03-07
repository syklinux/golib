package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/syklinux/golib/log"
)

var RedisConn *redis.Client

func InitRedisCon(addr, passwd string, db int) {
	redisInstance := newRedis()
	redisInstance.Init(addr, passwd, db)
}

type Redis struct{}

func newRedis() *Redis {
	return new(Redis)
}

func (r *Redis) Init(addr, passwd string, db int) {
	RedisConn = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})

	_, err := RedisConn.Ping(RedisConn.Context()).Result()
	if err != nil {
		panic(fmt.Errorf("redis client initialization failed: %w", err))
	}
}

func Close() error {
	if err := RedisConn.Close(); err != nil {
		log.Errorf("failed to shutdown connect: [REDIS]! => %v", err)
		return err
	}

	log.Infof("[REDIS] connection closed successfully")
	return nil
}
