package ts

import (
	"sync"
	"sync/atomic"
	"time"
)

func NewTimer() *Timer {
	return new(Timer)
}

// Timer is a thread safe simple timer
type Timer struct {
	mu sync.Mutex

	startAt time.Time
	stopAt  time.Time

	running atomic.Bool
}

// Begin begins timer
func (t *Timer) Begin() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.running.Load() {
		return
	}
	t.running.Store(true)

	t.startAt = time.Now()
}

// Stop stops timer
func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.running.Load() {
		return
	}
	t.running.Store(false)

	t.stopAt = time.Now()
}

// Duration return the duration between startAt and stopAt.
func (t *Timer) Duration() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.running.Load() {
		return t.stopAt.Sub(t.startAt)
	} else {
		return Now().Sub(t.startAt)
	}
}

// Reset resets timer
func (t *Timer) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.running.Load() {
		return
	}

	t.startAt = Zero()
	t.stopAt = Zero()
}
