package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ErrCodeSendTooMany error = fmt.Errorf("send too many items")
var ErrCodeVerifyTooMany error = fmt.Errorf("verify too many items")
var ErrCodeVerifyErr error = fmt.Errorf("code verify err")

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/varify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, expectedCode string) error
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{client: client}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	result, err := c.client.Eval(ctx, luaSetCode, []string{c.Key(ctx, biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch result {
	case 0:
		//	正常
		return nil
	case -1:
		return ErrCodeSendTooMany
	default:
		return fmt.Errorf("发送验证码系统错误")
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, expectedCode string) error {
	result, err := c.client.Eval(ctx, luaVerifyCode, []string{c.Key(ctx, biz, phone), expectedCode}).Int()
	if err != nil {
		return err
	}
	switch result {
	case 0:
		//	正常
		return nil
	case -1:
		return ErrCodeVerifyTooMany
	case -2:
		return ErrCodeVerifyErr
	}
	return nil
}

func (c *RedisCodeCache) Key(ctx context.Context, biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s:v1", biz, phone)
}
