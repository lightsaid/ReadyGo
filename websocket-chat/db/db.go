package db

import (
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}
