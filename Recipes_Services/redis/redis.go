package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	NilErr = errors.New("key not exist")
)

//go:generate mockery --name Client
type Client interface {
	// Send Ping test to the server return both status and error if present
	Ping(ctx context.Context) StatusResult
	Close() error
	Get(ctx context.Context, key string) StringResult
	Set(
		ctx context.Context,
		key string,
		value interface{},
		expiration time.Duration,
	) StatusResult
	Del(ctx context.Context, keys ...string) IntegerResult
	IncrBy(ctx context.Context, key string, value int64) IntegerResult
	Rpush(ctx context.Context, key string, values ...interface{}) IntegerResult
	LRange(ctx context.Context, key string, start, stop int64) SliceStringResult
	Llen(ctx context.Context, key string) IntegerResult
	Expire(ctx context.Context, key string, expiration time.Duration) BoolResult
}

// Result  associate with status command
type StatusResult interface {
	// full status
	String() string
	// Only the value of status
	Val() string
	// Status + error if  present
	Result() (string, error)
	Err() error
}

// Result associate with string command
type StringResult interface {
	Bool() (bool, error)
	Bytes() ([]byte, error)
	Err() error
	Float32() (float32, error)
	Float64() (float64, error)
	Int() (int, error)
	String() string
	Result() (string, error)
	Val() string
}

type IntegerResult interface {
	String() string
	Uint64() (uint64, error)
	Val() int64
	Result() (int64, error)
}

type SliceStringResult interface {
	Err() error
	Result() ([]string, error)
}

type BoolResult interface {
	Err() error
	Result() (bool, error)
}

//--------------------------------------------------------------------------------------------

// [Go-Redis] implmeneting StatusResult
type RedisResultStatus struct {
	status *redis.StatusCmd
}

func (r *RedisResultStatus) String() string {
	return r.status.String()
}
func (r *RedisResultStatus) Val() string {
	return r.status.Val()
}

func (r *RedisResultStatus) Result() (string, error) {
	return r.status.Result()
}

func (r *RedisResultStatus) Err() error {
	return r.status.Err()
}

//---------------------------------------------------------------------------------------------
// [Go-Redis] implementing StringResult

type RedisStringStatus struct {
	status *redis.StringCmd
}

func (r *RedisStringStatus) Bool() (bool, error) {
	return r.status.Bool()
}

func (r *RedisStringStatus) Bytes() ([]byte, error) {
	return r.status.Bytes()
}

func (r *RedisStringStatus) Err() error {
	err := r.status.Err()
	if err == redis.Nil {
		return NilErr
	}
	return err
}

func (r *RedisStringStatus) Float32() (float32, error) {
	return r.status.Float32()
}

func (r *RedisStringStatus) Float64() (float64, error) {
	return r.status.Float64()
}

func (r *RedisStringStatus) Int() (int, error) {
	return r.status.Int()
}

func (r *RedisStringStatus) String() string {
	return r.status.String()
}

func (r *RedisStringStatus) Result() (string, error) {
	result, err := r.status.Result()

	if err == redis.Nil {
		return "", NilErr
	}
	return result, err
}

func (r *RedisStringStatus) Val() string {

	return r.status.Val()
}

// [Go-Redis] implementing  IntegerResult
type RedisIntegerStatus struct {
	status *redis.IntCmd
}

func (r *RedisIntegerStatus) String() string {

	return r.status.String()
}

func (r *RedisIntegerStatus) Uint64() (uint64, error) {
	return r.status.Uint64()
}

func (r *RedisIntegerStatus) Val() int64 {
	return r.status.Val()
}
func (r *RedisIntegerStatus) Result() (int64, error) {
	return r.status.Result()

}

// [Go-Redis] implementing SliceStringResult
type RedisStringSliceStatus struct {
	status *redis.StringSliceCmd
}

func (r *RedisStringSliceStatus) Err() error {
	if r.status.Err() == redis.Nil {
		return NilErr
	}
	return r.status.Err()
}

func (r *RedisStringSliceStatus) Result() ([]string, error) {
	result, err := r.status.Result()
	if err == redis.Nil {
		return nil, NilErr
	}

	return result, nil
}

// [Go-Redis] implementing BoolResult
type RedisBoolStatus struct {
	status *redis.BoolCmd
}

func (r *RedisBoolStatus) Err() error {
	return r.status.Err()
}

func (r *RedisBoolStatus) Result() (bool, error) {
	return r.status.Result()
}

// [Go-Redis] implementing Client

// Concurrent Safe , can config Pool Size
func NewRedisClient(opt *redis.Options) Client {
	return &RedisClient{client: redis.NewClient(opt)}
}

type RedisClient struct {
	client *redis.Client
}

func (r *RedisClient) Ping(ctx context.Context) StatusResult {
	return &RedisResultStatus{
		status: r.client.Ping(ctx),
	}
}

func (r *RedisClient) Get(ctx context.Context, key string) StringResult {

	return &RedisStringStatus{
		status: r.client.Get(ctx, key),
	}
}

func (r *RedisClient) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) StatusResult {

	return &RedisResultStatus{
		status: r.client.Set(ctx, key, value, expiration),
	}
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) IntegerResult {
	return &RedisIntegerStatus{
		status: r.client.Del(ctx, keys...),
	}

}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) IncrBy(ctx context.Context, key string, value int64) IntegerResult {
	return &RedisIntegerStatus{
		status: r.client.IncrBy(ctx, key, value),
	}
}

func (r *RedisClient) Rpush(ctx context.Context, key string, values ...interface{}) IntegerResult {

	return &RedisIntegerStatus{
		status: r.client.RPush(ctx, key, values),
	}
}

func (r *RedisClient) LRange(ctx context.Context, key string, start, stop int64) SliceStringResult {
	return &RedisStringSliceStatus{
		status: r.client.LRange(ctx, key, start, stop),
	}
}

func (r *RedisClient) Llen(ctx context.Context, key string) IntegerResult {
	return &RedisIntegerStatus{
		status: r.client.LLen(ctx, key),
	}

}

func (r *RedisClient) Expire(ctx context.Context, key string, time time.Duration) BoolResult {
	return &RedisBoolStatus{
		status: r.client.Expire(ctx, key, time),
	}

}
