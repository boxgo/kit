package redis

import "github.com/go-redis/redis"

type (
	// Redis config
	Redis struct {
		MasterName   string   `config:"masterName" desc:"The sentinel master name. Only failover clients."`
		Address      []string `config:"address" desc:"Either a single address or a seed list of host:port addresses of cluster/sentinel nodes."`
		Password     string   `config:"password" desc:"Redis password"`
		DB           int      `config:"db" desc:"Database to be selected after connecting to the server. Only single-node and failover clients."`
		PoolSize     int      `config:"poolSize" desc:"Connection pool size"`
		MinIdleConns int      `config:"minIdleConns" desc:"min idle connections"`

		name string
		redis.UniversalClient
	}
)

var (
	// DefaultRedis 默认的redis配置
	DefaultRedis = New("redis")
)

// Name 配置文件
func (l *Redis) Name() string {
	return l.name
}

// ConfigWillLoad 配置文件将要加载
func (l *Redis) ConfigWillLoad() {

}

// ConfigDidLoad 配置文件已经加载。做一些默认值设置
func (l *Redis) ConfigDidLoad() {
	if len(l.Address) == 0 || l.name == "" {
		panic("redis config is invalid")
	}

	l.UniversalClient = redis.NewUniversalClient(&redis.UniversalOptions{
		MasterName:   l.MasterName,
		Addrs:        l.Address,
		Password:     l.Password,
		DB:           l.DB,
		PoolSize:     l.PoolSize,
		MinIdleConns: l.MinIdleConns,
	})
}

// New a redis lock
func New(name string) *Redis {
	return &Redis{
		name: name,
	}
}
