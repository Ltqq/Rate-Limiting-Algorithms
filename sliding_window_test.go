package limiting

import (
	"fmt"

	"testing"
	"time"
)

func TestNewSlidingWindow(t *testing.T) {
	limiter := NewSlidingWindow(5, 10) // 每10秒最多5个请求

	for i := 0; i < 15; i++ {
		if limiter.Allow() {
			fmt.Println("Request", i, "accepted")
		} else {
			fmt.Println("Request", i, "rejected")
		}
		time.Sleep(2 * time.Second) // 每2秒发送一个请求
	}

}
