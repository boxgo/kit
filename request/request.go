package request

import (
	"context"
	"time"

	"github.com/BiteBit/gorequest"
	"github.com/boxgo/box/minibox"
)

type (
	// Options http request tool
	Options struct {
		Timeout   int64  `config:"timeout" desc:"Timeout millsecond, default 10s"`
		UserAgent string `config:"userAgent" desc:"Client User-Agent"`
		ShowLog   bool   `config:"showLog" desc:"Show request log"`
		Trace     bool   `config:"trace" desc:"Open prometheus trace"`

		app minibox.App
	}
)

var (
	// GlobalOptions global request options
	GlobalOptions = &Options{}

	befores []gorequest.Before
	afters  []gorequest.After
)

// Name config prefix name
func (opts *Options) Name() string {
	return "request"
}

// Exts 获取app信息
func (opts *Options) Exts() []minibox.MiniBox {
	return []minibox.MiniBox{&opts.app}
}

func (opts *Options) ConfigWillLoad(context.Context) {

}

func (opts *Options) ConfigDidLoad(context.Context) {
	if opts.UserAgent == "" {
		opts.UserAgent = opts.app.AppName
	}
}

// NewTraceRequest new a trace request
func NewTraceRequest(ctx context.Context) *gorequest.SuperAgent {
	agent := gorequest.NewWithContext(ctx)

	setup(agent)

	return agent
}

// UseBefore global use before
func UseBefore(bs ...gorequest.Before) {
	befores = append(befores, bs...)
}

// UseAfter global use after
func UseAfter(as ...gorequest.After) {
	afters = append(afters, as...)
}

func setup(agent *gorequest.SuperAgent) {
	timeout := time.Second * 10
	if GlobalOptions.Timeout != 0 {
		timeout = time.Duration(GlobalOptions.Timeout * int64(time.Millisecond))
	}

	agent.Timeout(timeout)
	agent.UseBefore(logBefore)
	agent.UseAfter(logAfter)
	agent.UseBefore(befores...)
	agent.UseAfter(afters...)
}
