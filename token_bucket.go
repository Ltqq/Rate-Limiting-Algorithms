package limiting

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity       int           // 令牌桶的容量
	tokens         int           // 当前的令牌数
	refillRate     int           // 令牌添加速率（每秒多少个令牌）
	refillInterval time.Duration // 令牌添加间隔
	mutex          sync.Mutex    // 互斥锁，用于并发控制
}

// NewTokenBucket 创建一个新的令牌桶
func NewTokenBucket(capacity int, refillRate int) *TokenBucket {
	tb := &TokenBucket{
		capacity:       capacity,
		tokens:         capacity,
		refillRate:     refillRate,
		refillInterval: time.Second / time.Duration(refillRate),
	}

	// 启动一个goroutine定期添加令牌
	go tb.refillTokens()

	return tb
}

// refillTokens 定期向令牌桶中添加令牌
func (tb *TokenBucket) refillTokens() {
	ticker := time.NewTicker(tb.refillInterval)
	defer ticker.Stop()

	for range ticker.C {
		tb.mutex.Lock()
		if tb.tokens < tb.capacity {
			tb.tokens++
		}
		tb.mutex.Unlock()
	}
}

// AllowRequest 尝试从令牌桶中获取一个令牌，如果成功则返回true，否则返回false
func (tb *TokenBucket) AllowRequest() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}
