package config

import (
	"context"
	"fmt"
	"khanhanhtr/sample/redis"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func NewRedisClient(options *goredis.Options) (redis.Client, error) {
	redis := redis.NewRedisClient(options)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	status := redis.Ping(ctx)
	fmt.Println(status.String())

	return redis, status.Err()
}
