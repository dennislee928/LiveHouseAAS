package cache

import (
	"context"
	"time"

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

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Del(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.Client.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *Redis) Incr(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

func (r *Redis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.Client.Expire(ctx, key, expiration).Err()
}
