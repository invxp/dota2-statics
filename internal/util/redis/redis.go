package redis

import (
	"fmt"
	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"time"
)

/*
工具包
Redis库，支持Sentinel
*/

type Redis struct {
	p   *redis.Pool
	log *log.Logger
}

func (r *Redis) Do(command string, args ...interface{}) (interface{}, error) {
	conn := r.p.Get()

	if conn == nil {
		return nil, fmt.Errorf("redis connection was nil")
	}

	if conn.Err() != nil {
		err := conn.Close()
		return nil, fmt.Errorf("redis connection status error: %v, %v", conn.Err(), err)
	}

	result, e := conn.Do(command, args...)

	_ = conn.Close()

	r.log.Printf("redis do: %s@%v, result: %v@%v", command, args, result, e)

	return result, e
}

func connectToRedis(host, pwd string, sentinels []string, sentinelName string, db, idle, active, timeout int) *redis.Pool {
	if len(sentinels) > 0 {
		st := &sentinel.Sentinel{
			Addrs:      sentinels,
			MasterName: sentinelName,
			Dial: func(addr string) (redis.Conn, error) {
				return redis.Dial("tcp", addr)
			},
		}
		return &redis.Pool{
			MaxIdle:     idle,
			MaxActive:   active,
			Wait:        true,
			IdleTimeout: time.Duration(timeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				masterAddr, err := st.MasterAddr()
				if err != nil {
					return nil, err
				}
				c, err := redis.Dial("tcp", masterAddr, redis.DialPassword(pwd), redis.DialDatabase(db), redis.DialConnectTimeout(5*time.Second))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if _, err := c.Do("PING"); err != nil {
					return err
				}
				return nil
			},
		}
	} else {
		return &redis.Pool{
			MaxIdle:     idle,
			MaxActive:   active,
			IdleTimeout: time.Duration(timeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", host, redis.DialPassword(pwd), redis.DialDatabase(db), redis.DialConnectTimeout(5*time.Second))
			},
			Wait: true,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if _, err := c.Do("PING"); err != nil {
					return err
				}
				return nil
			},
		}
	}
}

func New(host, password, sentinelName string, sentinels []string, database, maxIdle, maxActive, idleTimeout int, logPath string) *Redis {
	filePath := fmt.Sprintf("%s", logPath)
	if err := os.MkdirAll(filePath, 0755); err != nil {
		log.Panic(err)
	}
	fileName := fmt.Sprintf("%s/%s", filePath, "redis.log")
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Panic(err)
	}

	redisInstance := &Redis{}
	redisInstance.p = connectToRedis(host, password, sentinels, sentinelName, database, maxIdle, maxActive, idleTimeout)
	redisInstance.log = log.New(file, "", log.LstdFlags|log.Lshortfile)

	if len(sentinels) > 0 {
		redisInstance.log.Printf("connect to redis sentinel: %s@%v|%s/%d\n", password, sentinels, sentinelName, database)
	} else {
		redisInstance.log.Printf("connect to redis host: %s@%s/%d\n", password, host, database)
	}

	if e := redisInstance.Ping(); e != nil {
		redisInstance.log.Fatalf("redis ping failed: %v\n", e)
	}

	redisInstance.log.Printf("redis ping success\n")

	go func() {
		totalFault := 0
		for {
			if e := redisInstance.Ping(); e != nil {
				redisInstance.log.Printf("redis ping failed: %v\n", e)
				totalFault++
			} else {
				if totalFault > 0 {
					redisInstance.log.Printf("redis ping recover\n")
				}
				totalFault = 0
			}
			if totalFault > 10 {
				redisInstance.log.Fatalf("redis ping fatal exit")
			}
			time.Sleep(time.Second * 30)
		}
	}()

	return redisInstance
}

func (r *Redis) Stop() {
	_ = r.p.Close()
}
