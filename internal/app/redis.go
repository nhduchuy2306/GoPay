package app

import (
	"fmt"
	"gopay/pkg"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	cfg := pkg.RedisConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	client, err := pkg.CreateRedisConnection(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Redis:", client)
	return client
}
