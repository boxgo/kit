package store

import (
	"time"

	"github.com/boxgo/box/logger"
	"github.com/boxgo/kit/mongo"
	"github.com/globalsign/mgo/bson"
)

type (
	// MgoTokenStore token存储器
	MgoTokenStore struct {
		name  string
		mongo *mongo.Mongo

		OriginID   string `json:"originId"`
		DB         string `json:"db"`
		Collection string `json:"collection"`
	}

	// wechatToken token db
	wechatToken struct {
		ExpireAt time.Time `bson:"expireAt"` // 参数发送时间
		ExipreIn int       `bson:"exipreIn"` // 过期周期
		Token    string    `bson:"token"`    // token值
		OriginID string    `bson:"originId"` // 原始id
	}
)

var (
	// DefaultMgoStore 默认的applet token store配置
	DefaultMgoStore = New("applet")
)

// Name token 配置文件名称
func (ts *MgoTokenStore) Name() string {
	return ts.name
}

// ConfigWillLoad 配置文件将要加载
func (ts *MgoTokenStore) ConfigWillLoad() {

}

// ConfigDidLoad 配置文件已经加载。做一些默认值设置
func (ts *MgoTokenStore) ConfigDidLoad() {
	if ts.OriginID == "" {
		panic("originId require config")
	}

	if ts.mongo == nil {
		panic("mongo require config")
	}

	if ts.DB == "" {
		ts.DB = "applet"
	}

	if ts.Collection == "" {
		ts.Collection = "applet_token"
	}
}

func (ts *MgoTokenStore) SetMongo(mongo *mongo.Mongo) {
	ts.mongo = mongo
}

// Get 获取token
func (ts *MgoTokenStore) Get() (token string, exipreIn int, expireAt time.Time) {
	doc := wechatToken{}
	err := ts.mongo.GetDB(ts.DB).C(ts.Collection).Find(bson.M{"originId": ts.OriginID}).One(&doc)
	if err != nil {
		logger.Default.Errorw("MgoTokenStore.Get.Error", "tokenStore", ts)
	}

	return doc.Token, doc.ExipreIn, doc.ExpireAt
}

// Set 保存token
func (ts *MgoTokenStore) Set(token string, exipreIn int, expireAt time.Time) {
	doc := wechatToken{
		Token:    token,
		ExipreIn: exipreIn,
		ExpireAt: expireAt,
		OriginID: ts.OriginID,
	}
	_, err := ts.mongo.GetDB(ts.DB).C(ts.Collection).Upsert(bson.M{"originId": ts.OriginID}, doc)

	if err != nil {
		logger.Default.Errorw("MgoTokenStore.Set.Error", "tokenStore", ts)
	}
}

// New 新建一个指定名称的tokenstore
func New(name string) *MgoTokenStore {
	return &MgoTokenStore{
		name: name,
	}
}
