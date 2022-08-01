package glog

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
)

// FileConfig 文件配置
type FileConfig struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
}

// NewFileWriter 初始化 FileWriter
func NewFileWriter(conf FileConfig) io.Writer {
	log.Default()
	return &lumberjack.Logger{
		Filename:   conf.Filename,
		MaxSize:    conf.MaxSize,
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,
		Compress:   conf.Compress,
	}
}
