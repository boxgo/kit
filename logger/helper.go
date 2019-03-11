package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func simpleConfig(level, encoding string) *zap.Config {
	return &zap.Config{
		Level:       zap.NewAtomicLevelAt(loggerLevel(level)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func loggerLevel(lvStr string) zapcore.Level {
	lv := zapcore.InfoLevel

	switch lvStr {
	case "debug":
		lv = zapcore.DebugLevel
	case "info":
		lv = zapcore.InfoLevel
	case "warn":
		lv = zapcore.WarnLevel
	case "error":
		lv = zapcore.ErrorLevel
	case "dpanic":
		lv = zapcore.DPanicLevel
	case "panic":
		lv = zapcore.PanicLevel
	case "fatal":
		lv = zapcore.FatalLevel
	}

	return lv
}
