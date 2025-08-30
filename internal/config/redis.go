package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	redisURL := os.Getenv("REDIS_URL")

	if redisURL == "" {
		log.Println("⚠️ REDIS_URL not found, Redis features will be disabled for development")
		return
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("⚠️ Failed to parse Redis URL: %v. Redis features disabled.", err)
		return
	}

	RedisClient = redis.NewClient(opt)

	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Printf("⚠️ Failed to connect Redis: %v. Redis features disabled.", err)
		RedisClient = nil
		return
	}
	log.Println("✅ Redis connected")
}
