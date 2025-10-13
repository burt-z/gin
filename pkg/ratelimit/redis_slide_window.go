package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var slideWindowLua string

type Limiter interface {
	Limited(ctx context.Context, key string) (bool, error)
}

type RedisSlideWindowLimiter struct {
	cmd      redis.Cmdable
	interval time.Duration
	rate     int64
}

func NewRedisSlideWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int64) Limiter {
	return &RedisSlideWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (rl *RedisSlideWindowLimiter) Limited(ctx context.Context, key string) (bool, error) {
	return rl.cmd.Eval(ctx, slideWindowLua, []string{key}, rl.interval.Milliseconds(), rl.rate, time.Now().UnixMilli()).Bool()
}
