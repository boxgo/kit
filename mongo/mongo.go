package mongo

import "github.com/globalsign/mgo"

type (
	// Mongo mongodb数据库
	Mongo struct {
		URI       string  `json:"uri"`
		DB        string  `json:"db"`
		PoolLimit uint    `json:"poolLimit"`
		Batch     uint    `json:"batch"`
		Prefetch  float64 `json:"prefetch"`
		Mode      uint    `json:"mode"`

		name    string
		session *mgo.Session
	}
)

var (
	// Default 默认的数据库连接
	Default = NewMongo("mongo")
)

// Name mongodb config name
func (m *Mongo) Name() string {
	return m.name
}

// ConfigWillLoad 配置文件将要加载
func (m *Mongo) ConfigWillLoad(context.Context) {

}

// ConfigDidLoad 配置文件已经加载。做一些默认值设置
func (m *Mongo) ConfigDidLoad(context.Context) {
	if m.PoolLimit == 0 {
		m.PoolLimit = 200
	}

	if m.Batch == 0 {
		m.Batch = 50
	}

	if m.Prefetch <= 0 {
		m.Prefetch = 0.20
	}
}

// GetSession 获取指定数据库连接
func (m *Mongo) GetSession() *mgo.Session {
	if m.session != nil {
		return m.session
	}

	sess, err := mgo.Dial(m.URI)
	if err != nil {
		panic(err)
	}

	m.session = sess
	m.session.SetMode(mgo.Mode(m.Mode), true)
	m.session.SetPoolLimit(200)
	m.session.SetBatch(50)
	m.session.SetPrefetch(0.20)

	return m.session
}

// GetDB 获取指定数据库
func (m *Mongo) GetDB(db string) *mgo.Database {
	return m.GetSession().DB(db)
}

// GetDefaultDB 获取配置文件中的db
func (m *Mongo) GetDefaultDB() *mgo.Database {
	return m.GetSession().DB(m.DB)
}

// NewMongo 新建一个mongodb链接
func NewMongo(name string) *Mongo {
	return &Mongo{
		name: name,
	}
}
