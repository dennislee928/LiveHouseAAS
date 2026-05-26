package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dennis-lee/LiveHouseAAS/backend/internal/infra/cache"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	redis    *cache.Redis
}

func NewRateLimiter(redis *cache.Redis) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		redis:    redis,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, times := range rl.requests {
			var valid []time.Time
			for _, t := range times {
				if now.Sub(t) < time.Minute {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = valid
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) allowInMemory(key string, limit int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	times := rl.requests[key]

	var valid []time.Time
	for _, t := range times {
		if now.Sub(t) < window {
			valid = append(valid, t)
		}
	}

	if len(valid) >= limit {
		rl.requests[key] = valid
		return false
	}

	valid = append(valid, now)
	rl.requests[key] = valid
	return true
}

func (rl *RateLimiter) allowRedis(key string, limit int, window time.Duration) bool {
	windowSec := int(window.Seconds())
	if windowSec < 1 {
		windowSec = 1
	}
	ctx := context.Background()
	redisKey := "rl:" + key

	count, err := rl.redis.Incr(ctx, redisKey)
	if err != nil {
		return rl.allowInMemory(key, limit, window)
	}

	if count == 1 {
		rl.redis.Expire(ctx, redisKey, time.Duration(windowSec)*time.Second)
	}

	return count <= int64(limit)
}

func (rl *RateLimiter) allow(key string, limit int, window time.Duration) bool {
	if rl.redis != nil && rl.redis.Client != nil {
		if err := rl.redis.Ping(context.Background()); err == nil {
			return rl.allowRedis(key, limit, window)
		}
	}
	return rl.allowInMemory(key, limit, window)
}

func (rl *RateLimiter) Middleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			key = userID.(string)
		}

		if !rl.allow(key, limit, window) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded, try again later",
			})
			return
		}
		c.Next()
	}
}

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	rl := NewRateLimiter(nil)
	return rl.Middleware(limit, window)
}
