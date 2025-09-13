package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cprakhar/gopher-social/internal/store"
	"github.com/go-redis/redis/v8"
)

type UserStore struct {
	rdb *redis.Client
}

const UserTimeExp = time.Minute

func (u *UserStore) Get(ctx context.Context, id string) (*store.User, error) {
	cacheKey := "user:" + id
	data, err := u.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
	}
	
	return &user, nil
}

func (u *UserStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := "user:" + user.ID
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return u.rdb.SetEX(ctx, cacheKey, data, UserTimeExp).Err()
}
