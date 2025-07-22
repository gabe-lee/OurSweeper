package lock

import (
	"context"
	"time"
)

const (
	defaultInitialSleep = time.Nanosecond << 16
	defaultMinimumSleep = time.Nanosecond * 16

	unlocked uint32 = 0
	locked   uint32 = 1
)

type Locker interface {
	Lock()
	TryLock() bool
	TryLockTimeout(timeout time.Duration) bool
	TryLockTimeoutAdv(timeout, initialSleep, minSleep time.Duration) bool
	TryLockContext(context context.Context) bool
	TryLockContextAdv(context context.Context, initialSleep, minSleep time.Duration) bool
	Unlock()
}
