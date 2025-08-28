package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // contoh: "localhost:6379"
		Password: "",                      // kalau Redis ada password, isi disini
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}
	log.Println("✅ Redis connected")
}
