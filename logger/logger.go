package logger

import (
	"context"

	"go.uber.org/zap"
)

type (
	// Logger logger option
	Logger struct {
		Level          string `json:"level" desc:"Levels: debug,info,warn,error,dpanic,panic,fatal"`
		Encoding       string `json:"encoding" desc:"PS: console or json"`
		TraceUID       string `json:"traceUid" desc:"Name as trace uid in context"`
		TraceRequestID string `json:"traceRequestId" desc:"Name as trace requestId in context"`
		*zap.SugaredLogger
	}
)

var (
	// Default the default logger
	Default = &Logger{}
)

func init() {
	if Default.SugaredLogger == nil {
		l, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}

		Default.SugaredLogger = l.Sugar()
	}
}

// Name logger config name
func (logger *Logger) Name() string {
	return "logger"
}

// ConfigWillLoad logger config
func (logger *Logger) ConfigWillLoad(context.Context) {
	if logger.Level == "" {
		logger.Level = "debug"
	}
	if logger.Encoding == "" {
		logger.Encoding = "console"
	}
	if logger.TraceUID == "" {
		logger.TraceUID = "uid"
	}
	if logger.TraceRequestID == "" {
		logger.TraceRequestID = "requestId"
	}

	log, err := simpleConfig(logger.Level, logger.Encoding).Build()
	if err != nil {
		panic(err)
	}

	defer log.Sync()

	logger.SugaredLogger = log.Sugar()
}

// ConfigDidLoad did load
func (logger *Logger) ConfigDidLoad(ctx context.Context) {
	logger.ConfigWillLoad(ctx)
}

// Trace logger with requestId and uid
func (logger *Logger) Trace(ctx context.Context) *zap.SugaredLogger {
	uid := ctx.Value(logger.TraceUID)
	requestID := ctx.Value(logger.TraceRequestID)

	return logger.SugaredLogger.With("requestId", requestID).With("uid", uid)
}
