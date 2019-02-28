package token

import (
	"time"

	"github.com/boxgo/box/logger"
	"github.com/boxgo/kit/mongo"
	"github.com/globalsign/mgo/bson"
)

type (
	// TokenStore token存储器
	TokenStore struct {
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
	// Default 默认的applet token store配置
	Default = New("applet")
)

// Name token 配置文件名称
func (ts *TokenStore) Name() string {
	return ts.name
}

// ConfigWillLoad 配置文件将要加载
func (ts *TokenStore) ConfigWillLoad() {

}

// ConfigDidLoad 配置文件已经加载。做一些默认值设置
func (ts *TokenStore) ConfigDidLoad() {
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

func (ts *TokenStore) SetMongo(mongo *mongo.Mongo) {
	ts.mongo = mongo
}

// Get 获取token
func (ts *TokenStore) Get() (token string, exipreIn int, expireAt time.Time) {
	doc := wechatToken{}
	err := ts.mongo.GetDB(ts.DB).C(ts.Collection).Find(bson.M{"originId": ts.OriginID}).One(&doc)
	if err != nil {
		logger.Default.Errorw("TokenStore.Get.Error", "tokenStore", ts)
	}

	return doc.Token, doc.ExipreIn, doc.ExpireAt
}

// Set 保存token
func (ts *TokenStore) Set(token string, exipreIn int, expireAt time.Time) {
	doc := wechatToken{
		Token:    token,
		ExipreIn: exipreIn,
		ExpireAt: expireAt,
		OriginID: ts.OriginID,
	}
	_, err := ts.mongo.GetDB(ts.DB).C(ts.Collection).Upsert(bson.M{"originId": ts.OriginID}, doc)

	if err != nil {
		logger.Default.Errorw("TokenStore.Set.Error", "tokenStore", ts)
	}
}

// New 新建一个指定名称的tokenstore
func New(name string) *TokenStore {
	return &TokenStore{
		name: name,
	}
}
