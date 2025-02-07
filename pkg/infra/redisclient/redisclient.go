package redisclient

import (
	"actions_google/pkg/config"
	"context"
	"fmt"
	"log"
	"os"

	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

var ErrActionExists = fmt.Errorf("action already exists")

func NewRedisClient() *RedisClient {
	opt, err := redis.ParseURL(config.GetEnv("VAULT_URI", ""))
	if err != nil {
		log.Panicf("ERROR | Not connected to Redis. Cannot parse Redis URI %s", config.GetEnv("VAULT_URI", ""))
	}

	rdb := redis.NewClient(opt)
	// _, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if os.Getenv("TEST") == "" {
		if rdb.Ping(context.Background()).Val() != "PONG" {
			log.Panicf("ERROR | Server Redis not pong")
		}
	}

	return &RedisClient{
		Client: rdb,
		Ctx:    context.Background(),
	}
}

func (r *RedisClient) Set(key string, value interface{}) error {
	return r.Client.Set(r.Ctx, key, value, 0).Err()
}

func (r *RedisClient) SetEx(key string, value interface{}, expiration time.Duration) (bool, error) {
	result, err := r.Client.Set(r.Ctx, key, value, expiration).Result()
	return result == "OK", err
}

func (r *RedisClient) Hset(key string, field string, values interface{}) bool {
	inserted := r.Client.HSet(r.Ctx, key, field, values).Val()
	return inserted != 0
}

func (r *RedisClient) Hget(key string, field string) error {
	return r.Client.HGet(r.Ctx, key, field).Err()
}

func (r *RedisClient) Hexists(key string, field string) (bool, error) {
	return r.Client.HExists(r.Ctx, key, field).Result()
}

func (r *RedisClient) Exists(key string) (int64, error) {
	return r.Client.Exists(r.Ctx, key).Result()
}

func (r *RedisClient) Get(key string) (string, error) {
	result, err := r.Client.Get(r.Ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return result, nil
}

// REDO
func (r *RedisClient) WatchToken(data string, key string, expires time.Duration) error {
	err := r.Client.Watch(r.Ctx, func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(r.Ctx, func(pipe redis.Pipeliner) error {
			err := pipe.SetNX(r.Ctx, key, data, expires).Err()
			return err
		})
		return err
	}, key)

	return err
}

func (r *RedisClient) AcquireLock(key, value string, expiration time.Duration) (bool, error) {
	return r.Client.SetNX(r.Ctx, key, value, expiration).Result()
}

func (r *RedisClient) RemoveLock(key string) (int64, error) {
	result, err := r.Client.Del(r.Ctx, key).Result()
	return result, err
}

func (r *RedisClient) Hdel(key string, field string) (int64, error) {
	result, err := r.Client.HDel(r.Ctx, key, field).Result()
	return result, err
}

func (r *RedisClient) WatchAndExecute(ctx context.Context, keys []string, txFunc func(tx *redis.Tx) error) error {
	return r.Client.Watch(ctx, txFunc, keys...)
}

func (r *RedisClient) ExecuteTransaction(ctx context.Context, keys []string, txFunc func(tx *redis.Tx) error) error {
	return r.WatchAndExecute(ctx, keys, txFunc)
}

func (r *RedisClient) HSetNX(key string, field *string, value string) (bool, error) {
	return r.Client.HSetNX(r.Ctx, key, *field, value).Result()
}
