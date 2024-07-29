package limiting

import (
	"sync"
	"time"
)

//How it works:
//Keep track of request count for the current and previous window.
//
//Calculate the weighted sum of requests based on the overlap with the sliding window.
//
//If the weighted sum is less than the limit, allow the request.

type SlidingWindow struct {
	CurrentWindowCount int
	LastWindowCount    int
	WindowSize         time.Duration
	WindowStartTime    time.Time
	MaxRequests        int
	mutex              sync.Mutex
}

// NewSlidingWindow 创建一个新的滑动窗口实例
func NewSlidingWindow(maxRequests int, windowSizeInSeconds int) *SlidingWindow {
	return &SlidingWindow{
		WindowSize:      time.Duration(windowSizeInSeconds) * time.Second,
		WindowStartTime: time.Now(),
		MaxRequests:     maxRequests,
	}
}

// Allow 检查是否允许请求
func (sw *SlidingWindow) Allow() bool {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(sw.WindowStartTime)
	if elapsed >= sw.WindowSize {
		// 移动窗口
		sw.WindowStartTime = now
		sw.LastWindowCount = sw.CurrentWindowCount
		sw.CurrentWindowCount = 0
		elapsed = 0
	}

	// 计算加权请求数量
	weight := float64(sw.LastWindowCount)*(1-float64(elapsed)/float64(sw.WindowSize)) + float64(sw.CurrentWindowCount)
	if weight+1 > float64(sw.MaxRequests) {
		return false
	}

	sw.CurrentWindowCount++
	return true
}
