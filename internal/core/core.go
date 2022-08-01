package core

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"good/internal/config"
	"good/pkg/glog"
	"good/pkg/limiter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"os"
	"time"
)

var (
	Logger  *zap.Logger
	Limiter limiter.Limiter
	Db      *gorm.DB
	Redis   *redis.Client
)

func Setup(settings *config.Settings) error {
	if err := initLogger(settings); err != nil {
		return err
	}

	if err := initLimiter(settings); err != nil {
		return err
	}

	if err := initDb(settings); err != nil {
		return err
	}

	if err := initRedis(settings); err != nil {
		return err
	}
	return nil
}

// initLogger 初始化日志
func initLogger(settings *config.Settings) error {
	var ws []io.Writer
	if settings.Log.Console.Enable {
		ws = append(ws, os.Stdout)
	}

	if settings.Log.File.Enable {
		conf := glog.FileConfig{
			Filename:   settings.Log.File.Filename,
			MaxSize:    settings.Log.File.MaxSize,
			MaxAge:     settings.Log.File.MaxAge,
			MaxBackups: settings.Log.File.MaxBackups,
			LocalTime:  settings.Log.File.LocalTime,
			Compress:   settings.Log.File.Compress,
		}
		ws = append(ws, glog.NewFileWriter(conf))
	}

	if settings.Log.Kafka.Enable {
		conf := glog.KafkaConfig{
			Addr:         settings.Log.Kafka.Addr,
			Topic:        settings.Log.Kafka.Topic,
			BatchSize:    settings.Log.Kafka.BatchSize,
			BatchBytes:   settings.Log.Kafka.BatchBytes,
			BatchTimeout: settings.Log.Kafka.BatchTimeout,
			WriteTimeout: settings.Log.Kafka.WriteTimeout,
			RequiredAcks: settings.Log.Kafka.RequiredAcks,
			Async:        settings.Log.Kafka.Async,
		}
		ws = append(ws, glog.NewKafkaWriter(conf))
	}

	Logger = glog.New(ws, zap.WithCaller(false))
	return nil
}

// initLimiter 初始化限流器
func initLimiter(settings *config.Settings) error {
	var err error
	o := limiter.Option{
		Name:     settings.Limiter.Name,
		Mode:     settings.Limiter.Mode,
		Rate:     settings.Limiter.Rate,
		Size:     settings.Limiter.Size,
		Interval: time.Duration(settings.Limiter.Interval) * time.Millisecond,
		Dsn:      settings.Limiter.Dsn,
	}
	if Limiter, err = limiter.New(o); err != nil {
		panic(err)
	}
	return nil
}

// initDb 初始化数据库链接
func initDb(settings *config.Settings) error {
	var err error
	Db, err = gorm.Open(mysql.Open(settings.Db.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	return err
}

// initRedis 初始化Redis链接
func initRedis(settings *config.Settings) error {
	options, err := redis.ParseURL(settings.Redis.Dsn)
	if err != nil {
		return err
	}
	Redis = redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = Redis.Ping(ctx).Err(); err != nil {
		return err
	}
	return nil
}
