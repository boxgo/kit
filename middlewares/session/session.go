package session

import (
	"context"
	"errors"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

type (
	Mode int

	Session struct {
		Mode           Mode     `json:"mode"`       // 运行模式
		SessionKeyName string   `json:"sessionKey"` // session key name
		MasterName     string   `json:"masterName"` // 主节点名称，用于哨兵模式
		Address        []string `json:"address"`    // 地址。单节点使用主节点ip；哨兵模式使用哨兵ip列表
		Password       string   `json:"password"`   // 访问秘钥
		DB             int      `json:"db"`         // 数据库
		KeyPair        string   `json:"keyPair"`    // 加密参数
		PoolSize       int      `json:"poolSize"`   // 连接池大小
	}
)

const (
	Standalone = Mode(1)
	Sentinel   = Mode(2)
	Cluster    = Mode(3)
)

var (
	Default = &Session{}
)

func (s *Session) Name() string {
	return "middleware.session"
}

// ConfigWillLoad 配置文件将要加载
func (s *Session) ConfigWillLoad(context.Context) {

}

// ConfigDidLoad 配置文件已经加载。做一些默认值设置
func (s *Session) ConfigDidLoad(context.Context) {
	if s.Mode == 0 {
		s.Mode = Standalone
	}

	if s.SessionKeyName == "" {
		s.SessionKeyName = "session"
	}

	if s.PoolSize == 0 {
		s.PoolSize = 256
	}
}

// Session
func (s *Session) Session() gin.HandlerFunc {
	var err error
	var store redis.Store

	switch s.Mode {
	case Standalone:
		store, err = redis.NewStoreWithDB(s.PoolSize, "tcp", s.Address[0], s.Password, strconv.Itoa(s.DB), []byte(s.KeyPair))
	case Sentinel:
		store, err = redis.NewStoreWithPool(newSentinelPool(s.MasterName, s.Address, SentinelOptions{
			Password: s.Password,
			DB:       s.DB,
		}), []byte(s.KeyPair))
	default:
		err = errors.New("未支持的Session redis集群类型")
	}

	if err != nil {
		panic(err)
	}

	return sessions.Sessions(s.SessionKeyName, store)
}
