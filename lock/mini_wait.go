package lock

import (
	"context"
	"sync/atomic"
	"time"
)

type MiniWait struct {
	sema int32
}

func (m *MiniWait) Add(n int32) {
	atomic.AddInt32(&m.sema, n)
}

func (m *MiniWait) Done() {
	atomic.AddInt32(&m.sema, -1)
}

func (m *MiniWait) Init(n int32) {
	atomic.StoreInt32(&m.sema, n)
}

func (m *MiniWait) Wait(sleepInterval time.Duration) {
	for !atomic.CompareAndSwapInt32(&m.sema, 0, 0) {
		time.Sleep(sleepInterval)
	}
}

func (m *MiniWait) WaitWithTimeout(sleepInterval time.Duration, timeout time.Duration) bool {
	var t time.Duration
	for !atomic.CompareAndSwapInt32(&m.sema, 0, 0) {
		time.Sleep(sleepInterval)
		t += sleepInterval
		if t >= timeout {
			return false
		}
	}
	return true
}

func (m *MiniWait) WaitWithContext(sleepInterval time.Duration, ctx context.Context) bool {
	var t time.Duration
	for !atomic.CompareAndSwapInt32(&m.sema, 0, 0) {
		time.Sleep(sleepInterval)
		t += sleepInterval
		select {
		case <-ctx.Done():
			return false
		default:
		}
	}
	return true
}

func (m *MiniWait) RemainingWaitCount() int32 {
	return atomic.LoadInt32(&m.sema)
}
