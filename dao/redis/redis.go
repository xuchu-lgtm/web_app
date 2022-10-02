package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"web_app/settings"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	_, err = client.Ping().Result()
	if err != nil {
		zap.L().Error("redis ping err", zap.Error(err))
	}
	return
}

func Close() {
	_ = client.Close()
}
