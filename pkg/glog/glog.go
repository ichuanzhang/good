package glog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

// New 日志初始化
func New(w []io.Writer, option ...zap.Option) *zap.Logger {
	if w == nil {
		logger, _ := zap.NewProduction(option...)
		return logger
	}

	var ws []zapcore.WriteSyncer
	for _, writer := range w {
		ws = append(ws, zapcore.AddSync(writer))
	}
	core := newCore(newEncoderConfig(), ws, zapcore.InfoLevel)
	return zap.New(core, option...)
}

// newEncoderConfig 初始化解码器配置
func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// newCore 初始化 core
func newCore(encoderConfig zapcore.EncoderConfig, ws []zapcore.WriteSyncer, level zapcore.Level) zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(ws...),
		zap.NewAtomicLevelAt(level),
	)
}
