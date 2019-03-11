package applet

import (
	"github.com/BiteBit/applet"
	"github.com/BiteBit/applet/api"
)

type (
	// Applet 小程序
	Applet struct {
		name string

		AppID       string `json:"appId"`
		AppSecret   string `json:"appSecret"`
		Token       string `json:"token"`
		AesKey      string `json:"aesKey"`
		APIDomain   string `json:"apiDomain"`
		APIBasePath string `json:"apiBasePath"`

		*applet.Applet
		tokenStore api.WechatTokenStore
		before     api.Before
		after      api.After
	}
)

var (
	// Default 默认的小程序sdk
	Default = New("applet")
)

// Name 配置文件名称
func (app *Applet) Name() string {
	return app.name
}

// SetTokenStore 设置tokenstore
func (app *Applet) SetTokenStore(tokenStore api.WechatTokenStore) {
	app.tokenStore = tokenStore
}

// SetBefore set before hook
func (app *Applet) SetBefore(before api.Before) {
	app.before = before
}

// SetAfter set after hook
func (app *Applet) SetAfter(after api.After) {
	app.after = after
}

// ConfigWillLoad 配置文件将要加载
func (app *Applet) ConfigWillLoad(context.Context) {

}

// ConfigDidLoad 配置文件已经加载。做一些默认值设置
func (app *Applet) ConfigDidLoad(context.Context) {
	mp := applet.NewApplet(
		app.AppID,
		app.AppSecret,
		app.Token,
		app.AesKey,
		app.tokenStore,
	)

	mp.API.SetDomain(app.APIDomain)
	mp.API.SetBasePath(app.APIBasePath)
	mp.API.SetBefore(app.before)
	mp.API.SetAftre(app.after)

	app.Applet = mp
}

// New 新建一个加载指定配置文件的小程序sdk
func New(name string) *Applet {
	return &Applet{
		name: name,
	}
}
