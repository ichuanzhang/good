package glog

import (
	"go.uber.org/zap"
	"io"
	"os"
	"testing"
	"time"
)

var (
	loggerConsole *zap.Logger
	loggerFile    *zap.Logger
	loggerKafka   *zap.Logger
)

func TestMain(m *testing.M) {
	//loggerConsole = New(nil)
	loggerConsole = New([]io.Writer{os.Stdout})

	fileConfig := FileConfig{
		Filename:   "test.log",
		MaxSize:    30,
		MaxAge:     7,
		MaxBackups: 30,
		LocalTime:  true,
		Compress:   false,
	}
	loggerFile = New([]io.Writer{NewFileWriter(fileConfig)})

	kafkaConf := KafkaConfig{
		Addr:         "10.0.0.11:9092",
		Topic:        "good-log",
		RequiredAcks: 1,
		BatchSize:    100,
		BatchBytes:   1048576,
		BatchTimeout: 1,
	}
	loggerKafka = New([]io.Writer{NewKafkaWriter(kafkaConf)})

	os.Exit(m.Run())
}

func BenchmarkConsole(b *testing.B) {
	start := time.Now()
	for i := 0; i < b.N; i++ {
		loggerConsole.Info("hello world",
			zap.String("method", "GET"),
			zap.Int("status", 200),
			zap.String("url", "/ping?name=zhangsan"),
			zap.Duration("duration", time.Since(start)),
		)
	}
}

func BenchmarkFile(b *testing.B) {
	start := time.Now()
	for i := 0; i < b.N; i++ {
		loggerFile.Info("hello world",
			zap.String("method", "GET"),
			zap.Int("status", 200),
			zap.String("url", "/ping?name=zhangsan"),
			zap.Duration("duration", time.Since(start)),
		)
	}
}

func BenchmarkKafka(b *testing.B) {
	start := time.Now()
	for i := 0; i < b.N; i++ {
		loggerKafka.Info("hello world",
			zap.String("method", "GET"),
			zap.Int("status", 200),
			zap.String("url", "/ping?name=zhangsan"),
			zap.Duration("duration", time.Since(start)),
		)
	}
}
