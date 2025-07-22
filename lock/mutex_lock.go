package lock

import (
	"context"
	"sync"
	"time"
)

type MutexLock struct {
	mu sync.Mutex
}

func (m *MutexLock) Lock() {
	m.mu.Lock()
}

func (m *MutexLock) TryLock() bool {
	return m.mu.TryLock()
}

func (m *MutexLock) TryLockContext(context context.Context) bool {
	return m.TryLockContextAdv(context, defaultInitialSleep, defaultMinimumSleep)
}

func (m *MutexLock) TryLockContextAdv(context context.Context, initialSleep time.Duration, minSleep time.Duration) bool {
	s := initialSleep
	for !m.mu.TryLock() {
		time.Sleep(s)
		select {
		case <-context.Done():
			return false
		default:
		}
		s >>= 1
		s |= minSleep
	}
	return true
}

func (m *MutexLock) TryLockTimeout(timeout time.Duration) bool {
	return m.TryLockTimeoutAdv(timeout, defaultInitialSleep, defaultMinimumSleep)
}

func (m *MutexLock) TryLockTimeoutAdv(timeout time.Duration, initialSleep time.Duration, minSleep time.Duration) bool {
	s := initialSleep
	t := time.Duration(0)
	for !m.mu.TryLock() {
		time.Sleep(s)
		t += s
		if t > timeout {
			return false
		}
		s >>= 1
		s |= minSleep
	}
	return true
}

func (m *MutexLock) Unlock() {
	m.mu.Unlock()
}

var _ Locker = (*MutexLock)(nil)
