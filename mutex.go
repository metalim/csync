package csync

import (
	"context"
)

// Same as sync.Mutex but with a context that can be cancelled.
type Mutex struct {
	lock chan struct{}
}

// NewMutex creates a new Mutex.
func NewMutex() *Mutex {
	return &Mutex{
		lock: make(chan struct{}, 1),
	}
}

// Lock locks m.
func (m *Mutex) Lock(ctx context.Context) error {
	// implicit init
	if m.lock == nil {
		m.lock = make(chan struct{}, 1)
	}

	select {
	case m.lock <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Unlock unlocks m.
// It is a design error if m is not locked on entry to Unlock.
func (m *Mutex) Unlock() {
	<-m.lock
}

// Locked returns true if m is locked.
func (m *Mutex) Locked() bool {
	return len(m.lock) > 0 // locked or not
}
