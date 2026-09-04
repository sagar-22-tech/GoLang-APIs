package main

import "github.com/redis/go-redis/v9"

var rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
