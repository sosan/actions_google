package repos

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClientInterface interface {
	Set(key string, value interface{}) error
	SetEx(key string, value interface{}, expiration time.Duration) (bool, error)
	Hset(key string, field string, values interface{}) bool
	Hget(key string, field string) error
	Hexists(key string, field string) (bool, error)
	Exists(key string) (int64, error)
	Get(key string) (string, error)
	WatchToken(data string, key string, expires time.Duration) error
	AcquireLock(key, value string, expiration time.Duration) (bool, error)
	RemoveLock(key string) (int64, error)
	Hdel(key string, field string) (int64, error)
	WatchAndExecute(ctx context.Context, keys []string, txFunc func(tx *redis.Tx) error) error
	ExecuteTransaction(ctx context.Context, keys []string, txFunc func(tx *redis.Tx) error) error
	HSetNX(key string, field *string, value string) (bool, error)
}
