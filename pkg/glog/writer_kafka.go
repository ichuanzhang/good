package glog

import (
	"bytes"
	"context"
	"github.com/segmentio/kafka-go"
	"io"
	"log"
	"time"
)

// KafkaWriter 实现 io.Writer接口
type KafkaWriter struct {
	Writer *kafka.Writer
}

// KafkaConfig kafka 配置
type KafkaConfig struct {
	Addr         string
	Topic        string
	BatchSize    int
	BatchBytes   int
	BatchTimeout int
	WriteTimeout int
	RequiredAcks int
	Async        bool
}

// NewKafkaWriter 初始化 KafkaWriter
func NewKafkaWriter(config KafkaConfig) io.Writer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.Addr),
		Topic:        config.Topic,
		BatchSize:    config.BatchSize,
		BatchBytes:   int64(config.BatchBytes),
		BatchTimeout: time.Duration(config.BatchTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
		RequiredAcks: kafka.RequiredAcks(config.RequiredAcks),
		Async:        config.Async,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				log.Printf("kafka write message error: %v, message:%+v \n", err, messages)
			}
		},
	}
	return &KafkaWriter{writer}
}

// Write 实现 io.writer()
func (w *KafkaWriter) Write(p []byte) (n int, err error) {
	message := kafka.Message{
		Value: bytes.TrimLeft(p, "\n"),
	}
	if err = w.Writer.WriteMessages(context.Background(), message); err != nil {
		return 0, err
	}
	return len(p), nil
}
