package di

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisConfig struct {
	Host     string
	Password string
	Port     int
}

func InitRedis() *redis.Client {
	cfg := &RedisConfig{
		Host:     viper.GetString("redis.host"),
		Password: viper.GetString("redis.password"),
		Port:     viper.GetInt("redis.port"),
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic("failed to connect to redis: " + err.Error())
	}

	return redisClient
}
