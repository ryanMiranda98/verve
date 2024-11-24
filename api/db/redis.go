package db

import (
	"context"

	redis "github.com/redis/go-redis/v9"
)

type RedisDB struct {
	Addr   string
	client *redis.Client
}

func NewRedisDB(addr string) *RedisDB {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisDB{
		Addr:   addr,
		client: client,
	}
}

func (r *RedisDB) Connect() error {
	_, err := r.client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisDB) Close() error {
	return r.client.Close()
}

func (r *RedisDB) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	return err
}

func (r *RedisDB) Set(ctx context.Context, key string, value interface{}) (interface{}, error) {
	return r.client.SAdd(ctx, key, value).Result()
}

func (r *RedisDB) Get(ctx context.Context, key string) (interface{}, error) {
	return r.client.SCard(ctx, key).Result()
}

func (r *RedisDB) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}
