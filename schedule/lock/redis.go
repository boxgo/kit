package lock

import (
	"time"

	"github.com/boxgo/box/logger"
	"github.com/boxgo/kit/redis"
	// "github.com/go-redis/redis"
)

type (
	// RedisLock redis lock
	RedisLock struct {
		*redis.Redis
	}
)

var (
	// DefaultRedisLock 默认的redislock配置
	DefaultRedisLock = New("redislock")
)

// Lock key锁定
func (l *RedisLock) Lock(key string, ttl time.Duration) (bool, error) {
	logger.Default.Debugf("redislock.Lock key: %s ttl: %s", key, ttl)
	return l.SetNX(key, 1, ttl).Result()
}

// IsLocked key是否被锁定
func (l *RedisLock) IsLocked(key string) (bool, error) {
	result, err := l.Get(key).Result()

	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}

		return false, err
	}

	if result == "1" {
		return true, nil
	}

	return false, nil
}

// UnLock 解锁key
func (l *RedisLock) UnLock(key string) error {
	_, err := l.Del(key).Result()
	logger.Default.Debugf("redislock.UnLock key: %s err: %s", key, err)

	return err
}

// New a redis lock
func New(name string) *RedisLock {
	return &RedisLock{
		Redis: redis.New(name),
	}
}
