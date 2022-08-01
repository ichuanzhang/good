package limiter

import (
	"errors"
	"time"
)

const (
	ModeSingle      = iota //单节点限流模式
	ModeDistributed        //分布式限流模式
)

// Limiter 限流器接口
type Limiter interface {
	// GetMode 获取限流器模式
	GetMode() (n int)
	// GetRate 获取投放令牌速率
	GetRate() (n int)
	// GetToken 获取令牌，如果获取到的令牌数量小于 num, ok 返回 false
	GetToken(num int) (n int, ok bool)
}

// Option 初始化参数
type Option struct {
	Name     string        //限流器名称，分布式限流用于生成相关键名
	Mode     int           //模式：0-代表单机限流，多节点部署也支持使用此模式；1-代表分布式限流，基于 redis 实现
	Rate     int           //速率，固定间隔时间投放令牌的数量
	Size     int           //桶大小
	Interval time.Duration //投放令牌的间隔时间
	Dsn      string        //数据源，分布式限流使用的数据源
}

// New 初始化限流器
func New(o Option) (l Limiter, err error) {
	switch o.Mode {
	case ModeSingle:
		l, err = newSingleLimiter(o)
	case ModeDistributed:
		l, err = newDistributedLimiter(o)
	default:
		err = errors.New("new limiter error: not support mode")
	}
	return l, err
}
