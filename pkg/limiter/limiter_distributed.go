package limiter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var (
	addTokenScript = redis.NewScript(`
local key = KEYS[1]
local num = tonumber(ARGV[1])
local size = tonumber(redis.call("HGET", key, "size"))
local count = tonumber(redis.call("HGET", key, "count"))
local n = 0
if count + num > size then
	n = size - count
	count = size
else
	n = num
	count = count + num
end
redis.call("HSET", key, "count", count)
return n
`)

	subTokenScript = redis.NewScript(`
local key = KEYS[1]
local num = tonumber(ARGV[1])
local size = tonumber(redis.call("HGET", key, "size"))
local count = tonumber(redis.call("HGET", key, "count"))
local n = 0
if count < num then
	n = count
	count = 0
else
	n = num
	count = count - num
end
redis.call("HSET", key, "count", count)
return n
`)
)

// distributedLimiter 分布式限流器
type distributedLimiter struct {
	name     string
	mode     int
	rate     int
	size     int
	count    int
	interval time.Duration
	client   *redis.Client
}

// newDistributedLimiter 初始化限流器
func newDistributedLimiter(o Option) (Limiter, error) {
	var (
		err    error
		client *redis.Client
	)
	if client, err = initRedis(o.Dsn); err != nil {
		return nil, err
	}

	l := &distributedLimiter{
		name:     o.Name,
		mode:     o.Mode,
		rate:     o.Rate,
		size:     o.Size,
		count:    o.Rate,
		interval: o.Interval,
		client:   client,
	}

	value := make(map[string]interface{})
	value["name"] = o.Name
	value["mode"] = o.Mode
	value["rate"] = o.Rate
	value["size"] = o.Size
	value["count"] = o.Rate
	value["interval"] = o.Interval

	if err = l.client.HSet(context.Background(), l.limiterKey(), value).Err(); err != nil {
		return nil, err
	}

	go l.run()
	return l, nil
}

// run 启动限流器
func (l *distributedLimiter) run() {
	ticker := time.NewTicker(l.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			l.addToken(l.rate)
		}
	}
}

// addToken 添加 token
func (l *distributedLimiter) addToken(num int) (int, bool) {
	var (
		n   int
		ok  bool
		err error
	)
	if num <= 0 {
		return 0, false
	}

	if n, err = l.addTokenScriptRun(l.rate); err != nil {
		log.Println("limiter error:", err)
		return n, false
	}

	if n == num {
		ok = true
	}
	return n, ok
}

// GetToken 获取 token
func (l *distributedLimiter) GetToken(num int) (int, bool) {
	var (
		n   int
		ok  bool
		err error
	)
	if num <= 0 {
		return 0, false
	}

	if n, err = l.subTokenScriptRun(num); err != nil {
		log.Println("limiter error:", err)
		return n, false
	}

	if n == num {
		ok = true
	}
	return n, ok
}

// GetMode 获取 mode
func (l *distributedLimiter) GetMode() int {
	return l.mode
}

// GetRate 获取 rate
func (l *distributedLimiter) GetRate() int {
	return l.rate
}

// initRedis 初始化 redis
func initRedis(dsn string) (*redis.Client, error) {
	options, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil
}

// limiterKey 获取键 limiter
func (l *distributedLimiter) limiterKey() string {
	return l.name + ":limiter"
}

// runAddTokenScript 执行增加 token 脚本
func (l *distributedLimiter) addTokenScriptRun(num int) (int, error) {
	return addTokenScript.Run(context.Background(), l.client, []string{l.limiterKey()}, num).Int()
}

// runAddTokenScript 执行减少 token 脚本
func (l *distributedLimiter) subTokenScriptRun(num int) (int, error) {
	return subTokenScript.Run(context.Background(), l.client, []string{l.limiterKey()}, num).Int()
}
