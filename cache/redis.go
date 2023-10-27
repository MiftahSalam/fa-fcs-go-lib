package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const defaultTimeout = 10 * time.Second

type redisCache struct {
	client redis.Cmdable
	cfg    Config
}

// NewRedisClient create redis client
func NewRedisClient(cfg Config) Cache {
	var client redis.Cmdable

	poolSize := 10
	authPass := ""

	if cfg.PoolSize != 0 {
		poolSize = cfg.PoolSize
	}
	if cfg.AuthPass != "" {
		authPass = cfg.AuthPass
	}
	if len(cfg.Servers) == 0 {
		panic("no server address found")
	}

	switch cfg.Topology {
	case Standalone:
		client = redis.NewClient(&redis.Options{
			Addr:         cfg.Servers[0],
			Password:     authPass,
			DialTimeout:  cfg.Timeout,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			PoolSize:     poolSize,
			MinIdleConns: cfg.MinIdleConns,
		})
	case Cluster:
		cluster := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:          cfg.Servers,
			Password:       authPass,
			DialTimeout:    cfg.Timeout,
			ReadTimeout:    cfg.Timeout,
			WriteTimeout:   cfg.Timeout,
			PoolSize:       poolSize,
			MinIdleConns:   cfg.MinIdleConns,
			RouteByLatency: true,
			RouteRandomly:  true,
		})

		err := cluster.ForEachShard(context.Background(), func(ctx context.Context, shard *redis.Client) error {
			return shard.Ping(ctx).Err()
		})
		if err != nil {
			panic(err)
		}

		client = cluster

	case Sentinel:
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:     cfg.MasterName,
			SentinelAddrs:  cfg.Servers,
			Password:       authPass,
			DialTimeout:    cfg.Timeout,
			ReadTimeout:    cfg.Timeout,
			WriteTimeout:   cfg.Timeout,
			PoolSize:       poolSize,
			MinIdleConns:   cfg.MinIdleConns,
			RouteByLatency: true,
			RouteRandomly:  true,
		})
	}
	// ping
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	return &redisCache{
		client: client,
		cfg:    cfg,
	}
}

func (m *redisCache) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) (err error) {
	err = m.client.Set(ctx, key, val, expiration).Err()
	return
}

func (m *redisCache) Add(ctx context.Context, key string, val []byte) (err error) {
	err = m.client.Append(ctx, key, string(val)).Err()
	return
}

func (m *redisCache) Get(ctx context.Context, key string) (data string, err error) {
	data, err = m.client.Get(ctx, key).Result()
	return
}

func (m *redisCache) Delete(ctx context.Context, key string) (err error) {
	err = m.client.Del(ctx, key).Err()
	return
}

func (m *redisCache) Incr(ctx context.Context, key string) (incr int64, err error) {
	incr, err = m.client.Incr(ctx, key).Result()
	return
}

func (m *redisCache) Ping(ctx context.Context) (pong string, err error) {
	pong, err = m.client.Ping(ctx).Result()
	return
}

func (m *redisCache) GetInstance() interface{} {
	return m.client
}

func (m *redisCache) HGet(ctx context.Context, key string, field string) (data string, err error) {
	data, err = m.client.HGet(ctx, key, field).Result()
	return
}

func (m *redisCache) HGetAll(ctx context.Context, key string) (data map[string]string, err error) {
	data, err = m.client.HGetAll(ctx, key).Result()
	return
}

func (m *redisCache) HSet(ctx context.Context, key string, values map[string]interface{}, expiration time.Duration) (err error) {
	err = m.client.HSet(ctx, key, values).Err()
	if err != nil {
		m.client.Expire(ctx, key, expiration)
	}
	return
}

func (m *redisCache) SAdd(ctx context.Context, key string, members ...interface{}) (err error) {
	err = m.client.SAdd(ctx, key, members).Err()
	return
}

func (m *redisCache) HDel(ctx context.Context, key string, fields ...string) (numOfField int64, err error) {
	numOfField, err = m.client.HDel(ctx, key, fields...).Result()
	return
}
