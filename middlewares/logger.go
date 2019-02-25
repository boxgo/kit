package middlewares

import (
	"time"

	"github.com/BiteBit/ginlog"
	"github.com/boxgo/box/logger"
	"github.com/gin-gonic/gin"
)

type (
	// Logger gin中间件
	Logger struct {
		TimeFormat        string `json:"timeFormat"`
		RequestBodyLimit  uint   `json:"requestBodyLimit"`
		RequestQueryLimit uint   `json:"requestQueryLimit"`
		ResponseBodyLimit uint   `json:"responseBodyLimit"`
	}
)

var (
	// DefaultLogger 默认的日志中间件
	DefaultLogger = &Logger{}
)

// Name 日志中间件配置名称
func (l *Logger) Name() string {
	return "middlewareLogger"
}

// ConfigWillLoad 配置文件将要加载
func (l *Logger) ConfigWillLoad() {

}

// ConfigDidLoad 配置文件已经加载。做一些默认值设置
func (l *Logger) ConfigDidLoad() {
	if l.TimeFormat == "" {
		l.TimeFormat = time.RFC3339
	}

	if l.RequestBodyLimit == 0 {
		l.RequestBodyLimit = 2000
	}

	if l.RequestQueryLimit == 0 {
		l.RequestQueryLimit = 2000
	}

	if l.ResponseBodyLimit == 0 {
		l.ResponseBodyLimit = 2000
	}
}

// Logger zap
func (l *Logger) Logger() gin.HandlerFunc {
	logger := logger.Default.Desugar()

	return ginlog.Logger(logger, ginlog.Options{
		TimeFormat:        time.RFC3339,
		RequestBodyLimit:  int(l.RequestBodyLimit),
		RequestQueryLimit: int(l.RequestQueryLimit),
		ResponseBodyLimit: int(l.ResponseBodyLimit),
	})
}
