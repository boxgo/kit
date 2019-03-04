package lock

import (
	"time"

	"github.com/boxgo/box/logger"
	"github.com/go-redis/redis"
)

type (
	// RedisLock redis lock
	RedisLock struct {
		MasterName   string   `config:"masterName" desc:"主节点名称，用于哨兵模式"`
		Address      []string `config:"address" desc:"单节点使用主节点ip；哨兵模式使用哨兵ip列表"`
		Password     string   `config:"password" desc:"访问秘钥"`
		DB           int      `config:"db" desc:"数据库索引"`
		PoolSize     int      `config:"poolSize" desc:"连接池大小"`
		MinIdleConns int      `config:"minIdleConns" desc:"最小空闲连接"`

		name        string
		redisClient redis.UniversalClient
	}
)

var (
	// DefaultRedisLock 默认的redislock配置
	DefaultRedisLock = New("redislock")
)

// Name 配置文件
func (l *RedisLock) Name() string {
	return l.name
}

// ConfigWillLoad 配置文件将要加载
func (l *RedisLock) ConfigWillLoad() {

}

// ConfigDidLoad 配置文件已经加载。做一些默认值设置
func (l *RedisLock) ConfigDidLoad() {
	l.redisClient = redis.NewUniversalClient(&redis.UniversalOptions{
		MasterName:   l.MasterName,
		Addrs:        l.Address,
		Password:     l.Password,
		DB:           l.DB,
		PoolSize:     l.PoolSize,
		MinIdleConns: l.MinIdleConns,
	})
}

// Lock key锁定
func (l *RedisLock) Lock(key string, ttl time.Duration) (bool, error) {
	logger.Default.Debugf("redislock.Lock key: %s ttl: %s", key, ttl)
	return l.redisClient.SetNX(key, 1, ttl).Result()
}

// IsLocked key是否被锁定
func (l *RedisLock) IsLocked(key string) (bool, error) {
	result, err := l.redisClient.Get(key).Result()

	if err != nil {
		if err.Error() == "redis: nil" {
			println(1)
			return false, nil
		}

		println(2)
		return false, err
	}

	if result == "1" {
		return true, nil
	}

	return false, nil
}

// UnLock 解锁key
func (l *RedisLock) UnLock(key string) error {
	_, err := l.redisClient.Del(key).Result()
	logger.Default.Debugf("redislock.UnLock key: %s err: %s", key, err)

	return err
}

// New a redis lock
func New(name string) *RedisLock {
	return &RedisLock{
		name: name,
	}
}
