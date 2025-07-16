package attempt_group

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

type AttemptGroup struct {
	Timeout      *Timeout
	currentCount atomic.Int64
	totalCount   atomic.Int64
	successCount atomic.Int64
	failureCount atomic.Int64
	Name         string
}

func NewWithTimeout(name string, duration time.Duration, initialCount int64) AttemptGroup {
	g := AttemptGroup{
		Name:    name,
		Timeout: NewTimeout(duration),
	}
	g.currentCount.Store(initialCount)
	g.totalCount.Store(initialCount)
	return g
}

func New(name string, initialCount int64) AttemptGroup {
	g := AttemptGroup{
		Name:    name,
		Timeout: NewTimeout(time.Hour * 9001),
	}
	g.currentCount.Store(initialCount)
	g.totalCount.Store(initialCount)
	return g
}

func (g *AttemptGroup) Add(count int64) {
	g.totalCount.Add(count)
	g.currentCount.Add(count)
}

func (g *AttemptGroup) Failure() {
	g.failureCount.Add(1)
	g.currentCount.Add(-1)
}

func (g *AttemptGroup) Success() {
	g.successCount.Add(1)
	g.currentCount.Add(-1)
}

func (g *AttemptGroup) Wait() error {
	keepWaiting := true
	for keepWaiting {
		select {
		case <-g.Timeout.Done():
			keepWaiting = false
		default:
			keepWaiting = !g.currentCount.CompareAndSwap(0, 1)
			if keepWaiting {
				time.Sleep(time.Nanosecond * 10)
			}
		}
	}
	g.Timeout.Cancel()
	total := g.totalCount.Load()
	succ := g.successCount.Load()
	if succ != total {
		fail := g.failureCount.Load()
		timo := total - fail - succ
		return fmt.Errorf("'%s' made %d attempts, had %d successes, %d failures, %d timeouts, ctx err: %s", g.Name, total, succ, fail, timo, g.Timeout.Err())
	}
	return nil
}

type Timeout struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTimeout(duration time.Duration) *Timeout {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	t := Timeout{
		ctx:    ctx,
		cancel: cancel,
	}
	return &t
}

func (t *Timeout) Deadline() (deadline time.Time, ok bool) {
	return t.ctx.Deadline()
}

func (t *Timeout) Err() error {
	return t.ctx.Err()
}

func (t *Timeout) Value(key any) any {
	return t.ctx.Value(key)
}

func (t *Timeout) Done() <-chan struct{} {
	return t.ctx.Done()
}

func (t *Timeout) Cancel() {
	t.cancel()
}

var _ context.Context = (*Timeout)(nil)
