package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
)

type Redis struct {
	client *redis.Client
}

func New(ctx context.Context, conf *config.Redis) (*Redis, error) {
	uri := conf.Host + ":" + conf.Port

	// connect with redis
	client := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: conf.Password,
		DB:       0, // use default DB
		Protocol: 2,
	})

	// test conn
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &Redis{
		client,
	}, nil
}

func (r *Redis) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, val, ttl).Err()
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(res), nil
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) DeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, prefix, 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			err := r.client.Del(ctx, key).Err()
			if err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}
