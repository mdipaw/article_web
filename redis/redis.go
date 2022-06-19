package redis

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisAddress = os.Getenv("REDIS_DSN")

type RedisClient struct {
	client *redis.Client
	exp    time.Duration
}

func NewRedis(exp time.Duration) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Password: "",
		Addr:     RedisAddress,
	})
	return &RedisClient{client: rdb, exp: exp}
}

func (r *RedisClient) Delete(key string) {
	r.client.Del(context.Background(), key)
}

func (r *RedisClient) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	sts := r.client.Set(context.Background(), key, string(data), r.exp)
	if sts.Err() != nil {
		return sts.Err()
	}
	return nil
}

func (r *RedisClient) Get(key string) (map[string]interface{}, error) {
	var result map[string]interface{}
	sts := r.client.Get(context.Background(), key)
	if sts.Err() != nil {
		return result, sts.Err()
	}
	data, err := sts.Result()
	if err != nil {
		return result, err
	}
	json.Unmarshal([]byte(data), &result)
	return result, nil
}
