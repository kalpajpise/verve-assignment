package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Endpoint           string
	Password           string
	MaxRetries         int
	MinIdleConnections int
}

type Redis struct {
	con redis.UniversalClient
}

func New(endpoint, password string, maxRetries, minIdleConnection int) (*Redis, error) {
	options := &redis.Options{
		Addr:         endpoint,
		MaxRetries:   maxRetries,
		MinIdleConns: minIdleConnection,
	}

	if password != "" {
		options.Password = password
	}

	return &Redis{
		con: redis.NewClient(options),
	}, nil
}

func (r *Redis) SetAdd(name string, value any) (int64, error) {
	return r.con.SAdd(context.Background(), name, value).Result()
}

func (r *Redis) SetLength(name string) (int64, error) {
	return r.con.SCard(context.Background(), "unique_ids").Result()
}

func (r *Redis) Delete(key string) (int64, error) {
	return r.con.Del(context.Background(), key).Result()
}
