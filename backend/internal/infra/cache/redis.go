package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(redisURL string) *Redis {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		opts = &redis.Options{
			Addr: "localhost:6379",
			DB:  0,
		}
	}

	client := redis.NewClient(opts)
	return &Redis{Client: client}
}

func (r *Redis) Close() {
	r.Client.Close()
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}
