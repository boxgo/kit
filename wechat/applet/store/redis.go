package store

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/boxgo/kit/logger"
	"github.com/boxgo/kit/redis"
)

/**
token在redis中的存储格式

key: prefix + originID
value: token + 有效时长（秒）+ 过期时间（时间戳）
*/

type (
	// RedisTokenStore token存储器
	RedisTokenStore struct {
		name     string
		redis    *redis.Redis
		OriginID string `json:"originId"`
	}
)

const (
	prefix = "applet.token."
)

var (
	DefaultRedisStore = New("redis")
)

// Name token 配置文件名称
func (ts *RedisTokenStore) Name() string {
	return ts.name
}

func (ts *RedisTokenStore) SetRedis(redis *redis.Redis) {
	ts.redis = redis
}

// ConfigWillLoad 配置文件将要加载
func (ts *RedisTokenStore) ConfigWillLoad() {

}

// ConfigDidLoad 配置文件已经加载。做一些默认值设置
func (ts *RedisTokenStore) ConfigDidLoad() {
	if ts.OriginID == "" {
		panic("originId require config")
	}
}

// Get 获取token
func (ts *RedisTokenStore) Get() (token string, exipreIn int, expireAt time.Time) {
	tokenValue, err := ts.redis.Get(ts.key()).Result()

	if err != nil {
		logger.Default.Errorw("RedisTokenStore.Get.Error", "originID", ts.OriginID, "err", err)
		return
	}

	return parseToken(tokenValue)
}

// Set 保存token
func (ts *RedisTokenStore) Set(token string, exipreIn int, expireAt time.Time) {
	str := formatToken(token, exipreIn, expireAt)
	_, err := ts.redis.Set(ts.key(), str, 0).Result()

	if err != nil {
		logger.Default.Errorw("RedisTokenStore.Set.Error", "originID", ts.OriginID, "err", err)
		return
	}
}

func (ts *RedisTokenStore) key() string {
	return prefix + ts.OriginID
}

func parseToken(raw string) (token string, exipreIn int, expireAt time.Time) {
	tokeninfo := strings.Split(raw, ":")
	if len(tokeninfo) != 3 {
		logger.Default.Errorw("parseToken.error", "token", raw)
		return
	}

	token = tokeninfo[0]
	exipreIn, _ = strconv.Atoi(tokeninfo[1])
	expireAtInt, _ := strconv.Atoi(tokeninfo[2])

	expireAt = time.Unix(int64(expireAtInt), 0)

	return
}

func formatToken(token string, exipreIn int, expireAt time.Time) string {
	return fmt.Sprintf("%s:%d:%d", token, exipreIn, expireAt.Unix())
}
