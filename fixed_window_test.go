package limiting

import (
	"fmt"
	"testing"
	"time"
)

func TestNewFixedWindowAlgorithm(t *testing.T) {
	algorithm := NewFixedWindowAlgorithm(10, 2)

	for i := 0; i < 30; i++ {
		if algorithm.Allow() {
			fmt.Printf("Request %d allowed\n", i)
		} else {
			fmt.Printf("Request %d dropped\n", i)
		}
		time.Sleep(100 * time.Millisecond) // 模拟每
	}
}

func TestNewFixedWindowAlgorithm2(t *testing.T) {
	limit := 3
	windowSize := time.Second * time.Duration(3)
	counter := NewFixedWindowCounter(limit, windowSize)

	for i := 0; i < 10; i++ {
		if counter.AllowRequest() {
			fmt.Println("Request", i, "accepted")
		} else {
			fmt.Println("Request", i, "rejected")
		}
		time.Sleep(time.Millisecond * 500)
	}
}
