package ioc

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisConfig struct {
	Host     string
	Password string
	DB       int
}

func InitRedis() *redis.Client {
	cfg := &RedisConfig{
		Host:     viper.GetString("redis.host"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic("failed to connect to redis: " + err.Error())
	}

	return redisClient
}
