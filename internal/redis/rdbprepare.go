package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	ctx       context.Context
	rdsClient *redis.Client
}

const (
	hostRds     = "localhost"
	portRds     = "6379"
	passwordRds = ""
	dbRds       = 0
)

func (r *RedisClient) AddData(key string, val []byte) {
	r.rdsClient.Set(r.ctx, key, val, 0)
}

func (r *RedisClient) GetData(key string) ([]byte, error) {
	res, err := r.rdsClient.Get(r.ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func PrepareRedis(ctx context.Context) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", hostRds, portRds),
		Password: passwordRds,
		DB:       dbRds,
	})
	m := &RedisClient{rdsClient: client, ctx: ctx}
	return m
}
