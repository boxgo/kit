package session

import (
	"errors"
	"time"

	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
)

type (
	// SentinelOptions sentinel 集群选项配置
	SentinelOptions struct {
		PoolMaxIdle         int           // 建议、默认 3
		PoolMaxActive       int           // 建议、默认 64
		PoolIdleTimeout     time.Duration // 建议、默认 240 * time.Second
		SentinelDialTimeout time.Duration // 建议、默认 500 * time.Millisecond
		Password            string        // redis 访问密码
		DB                  int           // redis 数据库
	}
)

func newSentinelPool(masterName string, addrs []string, opts SentinelOptions) *redis.Pool {
	if opts.SentinelDialTimeout == 0 {
		opts.SentinelDialTimeout = 500 * time.Millisecond
	}
	if opts.PoolMaxIdle == 0 {
		opts.PoolMaxIdle = 3
	}
	if opts.PoolMaxActive == 0 {
		opts.PoolMaxActive = 64
	}
	if opts.PoolIdleTimeout == 0 {
		opts.PoolIdleTimeout = 240 * time.Second
	}

	sntnl := &sentinel.Sentinel{
		Addrs:      addrs,
		MasterName: masterName,
		Dial: func(addr string) (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", addr, opts.SentinelDialTimeout, opts.SentinelDialTimeout, opts.SentinelDialTimeout)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	return &redis.Pool{
		Wait:        true,
		MaxIdle:     opts.PoolMaxIdle,
		MaxActive:   opts.PoolMaxActive,
		IdleTimeout: opts.PoolIdleTimeout,
		Dial: func() (redis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()

			if err != nil {
				return nil, err
			}

			c, err := redis.Dial("tcp", masterAddr)
			if err != nil {
				return nil, err
			}

			if _, err := c.Do("AUTH", opts.Password); len(opts.Password) != 0 && err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", opts.DB); err != nil {
				c.Close()
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if !sentinel.TestRole(c, "master") {
				return errors.New("Role check failed")
			}

			return nil
		},
	}
}
