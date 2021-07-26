package db

import (
	"Shopee_UMS/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type Redis struct {
	c *redis.Client
}

func NewRedis() (*Redis, error) {
	return newRedis(0)
}

func NewTestRedis() (*Redis, error) {
	return newRedis(1)
}

func newRedis(db int) (*Redis, error) {
	host := utils.MustEnv("REDIS_HOST")
	port := utils.MustEnv("REDIS_PORT")

	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := c.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &Redis{c: c}, nil
}

var ttl, _ = strconv.Atoi(utils.MustEnv("REDIS_TTL"))

func (r *Redis) StoreJson(key string, v interface{}) error {
	val, err := json.Marshal(v)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.c.Set(ctx, key, string(val), time.Duration(ttl)*time.Second).Err()
}

func (r *Redis) LoadJson(key string, v interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	val, err := r.c.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), v)
}

func (r *Redis) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.c.Del(ctx, key).Err()
}

func (r *Redis) FlushAll() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.c.FlushAll(ctx).Err()
}
