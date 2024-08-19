package lock

import (
	"sync"
	"time"
)

// 超时锁，尝试 n 毫秒锁失败则返回
type TimedMutex struct {
	sync.Mutex
}

func (tm *TimedMutex) TryLock(ms int) bool {
	if ms <= 0 {
		ms = 1
	}

	for i := 0; i < ms; i++ {
		if tm.Mutex.TryLock() {
			return true
		}
		time.Sleep(time.Millisecond)
	}
	return false
}
