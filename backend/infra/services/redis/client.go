package redisclient

import (
	"context"
	"log"
	"strconv"
	"time"

	"ai_hub.com/app/infra/config"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func ConnectRedis() *redis.Client {
	if client != nil {
		return client
	}

	timeoutMS, err := strconv.Atoi(config.Env.Timeout)
	if err != nil || timeoutMS <= 0 {
		timeoutMS = 5000
	}

	opts, err := redis.ParseURL(config.Env.RedisURL)
	if err != nil {
		log.Fatalf("[Redis] Invalid REDIS_URL: %v", err)
	}

	opts.DialTimeout = time.Duration(timeoutMS) * time.Millisecond

	opts.OnConnect = func(ctx context.Context, cn *redis.Conn) error {
		log.Println("[Redis] OnConnect fired")
		return nil
	}

	c := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.Ping(ctx).Err(); err != nil {
		log.Fatalf("[Redis] Connection error: %v", err)
	} else {
		log.Println("[Redis] Connected successfully")
	}

	client = c
	return client
}

func EnsureRedis() *redis.Client {
	if client == nil {
		return ConnectRedis()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("[Redis] Connection lost, reconnecting... (%v)", err)
		client = nil
		return ConnectRedis()
	}
	return client
}

func Get(ctx context.Context, key string) (string, error) {
	c := EnsureRedis()
	val, err := c.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func Set(ctx context.Context, key, value string, expiration time.Duration) error {
	c := EnsureRedis()
	return c.Set(ctx, key, value, expiration).Err()
}
