package limiting

import (
	"container/list"
	"sync"
	"time"
)

type LeakyBucket struct {
	Capacity int
	Queue    *list.List
	Speed    int
	sync.Mutex
	c        chan struct{}
	stopChan chan struct{}
}

func NewLeakyBucket(speed, capacity int) *LeakyBucket {
	l := &LeakyBucket{
		Queue:    list.New(),
		Speed:    speed,
		Capacity: capacity,
		c:        make(chan struct{}),
		stopChan: make(chan struct{}),
	}
	go l.Start()
	return l
}

func (l *LeakyBucket) InQueue(req int) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	if l.Queue.Len() < l.Capacity {
		l.Queue.PushBack(req)
	}
}

func (l *LeakyBucket) OutQueue() (int, bool) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	if l.Queue.Len() > 0 {
		element := l.Queue.Front()
		l.Queue.Remove(element)
		return element.Value.(int), true
	}
	return 0, false
}

func (l *LeakyBucket) Start() {
	ticker := time.NewTicker(time.Duration(l.Speed) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			l.Mutex.Lock()
			if l.Queue.Len() > 0 {
				l.c <- struct{}{}
			}
			l.Mutex.Unlock()
		case <-l.stopChan:
			return
		}
	}
}

func (l *LeakyBucket) Stop() {
	close(l.stopChan)
}

func (l *LeakyBucket) WaitForNext() {
	<-l.c
}
