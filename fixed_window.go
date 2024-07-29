package limiting

import (
	"fmt"
	"sync"
	"time"
)

//How it works:
//Time is divided into fixed windows (e.g., 1-minute intervals).
//
//Each window has a counter that starts at zero.
//
//New requests increment the counter for the current window.
//
//If the counter exceeds the limit, requests are denied until the next window.

type FixedWindowAlgorithm struct {
	Current    *Window
	Size       int
	WindowSize int
}

func NewFixedWindowAlgorithm(max, size int) *FixedWindowAlgorithm {
	window := NewWindow(max)
	f := &FixedWindowAlgorithm{
		Current:    window,
		Size:       size,
		WindowSize: max,
	}
	go f.RefreshWindow()
	return f
}

func (f *FixedWindowAlgorithm) Allow() bool {
	return f.Current.IsOk()
}
func (f *FixedWindowAlgorithm) RefreshWindow() {
	ticker := time.NewTicker(time.Duration(f.Size) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		fmt.Println("refreshing window")
		f.Current = NewWindow(f.WindowSize)
	}
}

type Window struct {
	Count int
	Max   int
	sync.Mutex
}

func NewWindow(max int) *Window {
	return &Window{
		Max:   max,
		Count: 0,
	}
}

func (w *Window) IsOk() bool {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()
	if w.Count < w.Max {
		w.Count++
		return true
	}
	return false
}

// FixedWindowCounter 结构体定义
type FixedWindowCounter struct {
	limit              int           // 每个窗口的请求限制
	windowSize         time.Duration // 窗口大小
	currentWindowStart time.Time     // 当前窗口的开始时间
	counter            int           // 当前窗口的请求计数
	mutex              sync.Mutex    // 互斥锁，保护并发访问
}

// NewFixedWindowCounter 创建一个新的固定窗口计数器
func NewFixedWindowCounter(limit int, windowSize time.Duration) *FixedWindowCounter {
	return &FixedWindowCounter{
		limit:              limit,
		windowSize:         windowSize,
		currentWindowStart: time.Now(),
	}
}

// AllowRequest 检查是否允许请求并更新计数
func (fwc *FixedWindowCounter) AllowRequest() bool {
	fwc.mutex.Lock()
	defer fwc.mutex.Unlock()

	now := time.Now()

	// 检查当前时间是否已经进入下一个窗口
	if now.Sub(fwc.currentWindowStart) >= fwc.windowSize {
		fwc.currentWindowStart = now
		fwc.counter = 0
	}

	// 如果当前窗口计数超过限制，拒绝请求
	if fwc.counter >= fwc.limit {
		return false
	}

	// 否则，接受请求并增加计数
	fwc.counter++
	return true
}
