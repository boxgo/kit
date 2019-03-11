package middlewares

import (
	"context"

	"github.com/BiteBit/ginprom"
)

type (
	// Prometheus 中间件
	Prometheus struct {
		NameSpace                       string `json:"namespace"`
		SubModule                       string `json:"submodule"`
		LabelHTTPRequestTotal           string `json:"requestTotal"`
		LabelHTTPrequestDurationSeconds string `json:"requestDurSeconds"`
		LabelHTTPrequestSizeBytes       string `json:"requestSizeBytes"`
		LabelHTTPresponseSizeBytes      string `json:"responseSizeBytes"`

		*ginprom.Prom
	}
)

var (
	// DefaultPrometheus 默认的prometheus中间件
	DefaultPrometheus = &Prometheus{}
)

// Name 配置名称
func (p *Prometheus) Name() string {
	return "middleware.prometheus"
}

// ConfigWillLoad 配置文件将要加载
func (p *Prometheus) ConfigWillLoad(context.Context) {

}

// ConfigDidLoad 配置文件已经加载。做一些默认值设置
func (p *Prometheus) ConfigDidLoad(context.Context) {
	if p.LabelHTTPRequestTotal == "" {
		p.LabelHTTPRequestTotal = "http_request_total"
	}

	if p.LabelHTTPrequestDurationSeconds == "" {
		p.LabelHTTPrequestDurationSeconds = "http_request_duration_seconds"
	}

	if p.LabelHTTPrequestSizeBytes == "" {
		p.LabelHTTPrequestSizeBytes = "http_request_size_bytes"
	}

	if p.LabelHTTPresponseSizeBytes == "" {
		p.LabelHTTPresponseSizeBytes = "http_response_size_bytes"
	}

	p.Prom = ginprom.New(
		p.NameSpace,
		p.SubModule,
		p.LabelHTTPRequestTotal,
		p.LabelHTTPrequestDurationSeconds,
		p.LabelHTTPrequestSizeBytes,
		p.LabelHTTPresponseSizeBytes)
}
