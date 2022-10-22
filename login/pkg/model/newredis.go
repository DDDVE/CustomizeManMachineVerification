package model

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisConfig struct {
	Host        string
	Pass        string
	Network     string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int64
}

func NewRedis(config *RedisConfig) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			setdb := redis.DialDatabase(0)               // 设置db,默认为0
			setPasswd := redis.DialPassword(config.Pass) // 设置redis连接密码
			c, err := redis.Dial(config.Network, config.Host, setdb, setPasswd)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}
