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

func DPanic(args ...interface{}) {
	Default.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	Default.DPanicf(template, args...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	Default.DPanicw(msg, keysAndValues...)
}

func Debug(args ...interface{}) {
	Default.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	Default.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	Default.Debugw(msg, keysAndValues...)
}

func Desugar() *zap.Logger {
	return Default.Desugar()
}

func Error(args ...interface{}) {
	Default.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	Default.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	Default.Errorw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	Default.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	Default.Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	Default.Fatalw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	Default.Info(args...)
}

func Infof(template string, args ...interface{}) {
	Default.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	Default.Infow(msg, keysAndValues...)
}

func Named(name string) *zap.SugaredLogger {
	return Default.Named(name)
}

func Panic(args ...interface{}) {
	Default.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	Default.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	Default.Panicw(msg, keysAndValues...)
}

func Sync() error {
	return Default.Sync()
}

func Warn(args ...interface{}) {
	Default.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	Default.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	Default.Warnw(msg, keysAndValues...)
}

func With(args ...interface{}) *zap.SugaredLogger {
	return Default.With(args...)
}
