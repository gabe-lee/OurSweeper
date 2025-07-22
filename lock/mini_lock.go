package lock

import (
	"context"
	"sync/atomic"
	"time"
)

type MiniLock struct {
	state uint32
}

// Lock implements sync.Locker.
func (m *MiniLock) Lock() {
	s := defaultInitialSleep
	for !atomic.CompareAndSwapUint32(&m.state, unlocked, locked) {
		time.Sleep(s)
		s >>= 1
		s |= defaultMinimumSleep
	}
}

func (m *MiniLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&m.state, unlocked, locked)
}

func (m *MiniLock) TryLockContext(context context.Context) bool {
	return m.TryLockContextAdv(context, defaultInitialSleep, defaultMinimumSleep)
}

func (m *MiniLock) TryLockContextAdv(context context.Context, initialSleep, minSleep time.Duration) bool {
	s := initialSleep
	for !atomic.CompareAndSwapUint32(&m.state, unlocked, locked) {
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

func (m *MiniLock) TryLockTimeout(timeout time.Duration) bool {
	return m.TryLockTimeoutAdv(timeout, defaultInitialSleep, defaultMinimumSleep)
}

func (m *MiniLock) TryLockTimeoutAdv(timeout, initialSleep, minSleep time.Duration) bool {
	s := initialSleep
	t := time.Duration(0)
	for !atomic.CompareAndSwapUint32(&m.state, unlocked, locked) {
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

func (m *MiniLock) Unlock() {
	atomic.StoreUint32(&m.state, unlocked)
}

var _ Locker = (*MiniLock)(nil)
