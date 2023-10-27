package cache

import (
	"context"
	"time"
)

// Cache key value interface
type Cache interface {
	Set(ctx context.Context, key string, val interface{}, expiration time.Duration) (err error)
	Add(ctx context.Context, key string, val []byte) (err error)
	Get(ctx context.Context, key string) (data string, err error)
	Delete(ctx context.Context, key string) (err error)
	Ping(ctx context.Context) (pong string, err error)
	GetInstance() interface{}
	HGet(ctx context.Context, key string, field string) (data string, err error)
	HGetAll(ctx context.Context, key string) (data map[string]string, err error)
	HSet(ctx context.Context, key string, values map[string]interface{}, expiration time.Duration) (err error)
	SAdd(ctx context.Context, key string, members ...interface{}) (err error)
	HDel(ctx context.Context, key string, fields ...string) (numOfField int64, err error)
}
