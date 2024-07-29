package limiting

import (
	"fmt"
	"testing"
	"time"
)

func TestNewLeakyBucket(t *testing.T) {
	bucket := NewLeakyBucket(1, 10)
	// 启动一个goroutine从漏桶中取出请求
	go func() {
		for {
			// 等待处理信号
			bucket.WaitForNext()

			// 从队列中出队请求
			req, ok := bucket.OutQueue()
			if ok {
				fmt.Printf("Processing request %d\n", req)
			} else {
				fmt.Println("No request to process")
			}
		}
	}()

	// 向漏桶中添加请求
	for i := 0; i < 10; i++ {
		fmt.Printf("Adding request %d to the bucket\n", i)
		bucket.InQueue(i)
		time.Sleep(500 * time.Millisecond) // 每隔500毫秒添加一个请求
	}

	// 运行一段时间后停止漏桶
	time.Sleep(10 * time.Second)
	bucket.Stop()
	fmt.Println("Leaky bucket stopped")
}
