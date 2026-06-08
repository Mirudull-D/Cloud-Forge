package queue

import (
	"CloudHub/config"
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Envs.RedisConn,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Redis connected Successfully ...! ")
	return rdb
}
