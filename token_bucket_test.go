package limiting

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTokenBucket(t *testing.T) {
	// 创建一个令牌桶，容量为10，每秒添加1个令牌
	tb := NewTokenBucket(10, 1)

	// 模拟请求
	for i := 0; i < 20; i++ {
		if tb.AllowRequest() {
			fmt.Printf("Request %d allowed\n", i)
		} else {
			fmt.Printf("Request %d dropped\n", i)
		}
		time.Sleep(100 * time.Millisecond) // 模拟每100毫秒一个请求
	}
}
