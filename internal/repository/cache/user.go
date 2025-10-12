package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"jike_gin/internal/domain"
	"time"
)

var ErrKeyNotFound = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
}

type RedisUserCache struct {
	client redis.Cmdable
	expire time.Duration
}

func NewUserCache(redisClient redis.Cmdable) UserCache {
	return &RedisUserCache{
		client: redisClient,
		expire: time.Minute * 30,
	}
}
func (c *RedisUserCache) key(ctx context.Context, id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func (u *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := u.key(ctx, id)
	data, err := u.client.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("Get===>", err.Error())
		return domain.User{}, err
	}
	user := domain.User{}
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		fmt.Println("Get===>", err.Error())
		return domain.User{}, err
	}
	return user, nil
}

func (u *RedisUserCache) Set(ctx context.Context, user domain.User) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		fmt.Println("==>", err)
		return err
	}
	key := u.key(ctx, user.Id)
	u.client.Set(ctx, key, string(userBytes), u.expire)
	return nil
}
