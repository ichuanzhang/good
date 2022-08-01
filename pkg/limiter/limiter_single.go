package limiter

import (
	"sync"
	"time"
)

// singleLimiter 单节点限流器
type singleLimiter struct {
	mode     int
	rate     int
	size     int
	count    int
	interval time.Duration
	mu       sync.Mutex
}

// newSingleLimiter 初始化限流器
func newSingleLimiter(o Option) (Limiter, error) {
	l := &singleLimiter{
		mode:     o.Mode,
		rate:     o.Rate,
		interval: o.Interval,
		size:     o.Size,
		count:    o.Rate,
	}
	go l.run()
	return l, nil
}

// run 启动限流器
func (l *singleLimiter) run() {
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
func (l *singleLimiter) addToken(num int) (n int, ok bool) {
	if num <= 0 {
		return 0, false
	}
	l.mu.Lock()
	if l.count+num > l.size {
		n = l.size - l.count
		l.count = l.size
	} else {
		n = num
		l.count += n
	}
	l.mu.Unlock()
	return n, true
}

// GetToken 获取 token
func (l *singleLimiter) GetToken(num int) (n int, ok bool) {
	if num <= 0 {
		return 0, false
	}
	l.mu.Lock()
	if l.count < num {
		n = l.count
		l.count = 0
		ok = false
	} else {
		n = num
		l.count -= n
		ok = true
	}
	l.mu.Unlock()
	return n, ok
}

// GetMode 获取 mode
func (l *singleLimiter) GetMode() int {
	return l.mode
}

// GetRate 获取 rate
func (l *singleLimiter) GetRate() int {
	return l.rate
}
