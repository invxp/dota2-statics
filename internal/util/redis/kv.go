package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func (r *Redis) Expire(key string, lifeCycleSecond uint) (int, error) {
	return redis.Int(r.Do("Expire", key, lifeCycleSecond))
}

func (r *Redis) SetEx(key string, value interface{}, expire uint) (string, error) {
	return redis.String(r.Do("SetEX", key, expire, value))
}

func (r *Redis) Get(key string) (string, error) {
	return redis.String(r.Do("Get", key))
}

func (r *Redis) Ping() error {
	result, e := redis.String(r.Do("PING"))

	if e != nil {
		return e
	}

	if result != "PONG" {
		return fmt.Errorf("ping result was not PONG: %s", result)
	}

	return nil
}
