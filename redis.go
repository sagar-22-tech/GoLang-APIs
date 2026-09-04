package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
}
