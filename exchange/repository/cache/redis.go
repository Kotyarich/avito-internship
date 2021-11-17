package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
	"time"
)

const rateExpireTime = time.Hour

type RedisCache struct {
	client *redis.Client
	ctx context.Context
}

func NewRedisCache() *RedisCache {
	host := os.Getenv("REDIS_HOST")
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Fatal(err.Error())
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
	})

	return &RedisCache{
		client: client,
		ctx: context.Background(),
	}
}

func (c *RedisCache) GetRubleRate(currency string) (float32, error) {
	val, err := c.client.Get(c.ctx, currency).Float32()
	if err != nil {
		return 0, err
	}

	return val, err
}

func (c *RedisCache) SetRate(currency string, rate float32) error {
	status := c.client.Set(c.ctx, currency, rate, rateExpireTime)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}