package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewRedisInstance(config *Config) (*redis.Client, error) {
	if config == nil {
		return nil, errors.New("*Config dependency is nil")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username: config.User,
		Password: config.Pass,
		DB:       config.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return rdb, nil
	}

	return rdb, nil
}
