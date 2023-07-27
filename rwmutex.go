package csync

import (
	"context"
	"sync/atomic"
)

// RWMutex is a reader/writer mutual exclusion lock.
// The lock can be held by an arbitrary number of readers or a single writer.
type RWMutex struct {
	Mutex                // Writer exclusion Mutex
	wait    Mutex        // Mutex for readers to wait on writer
	readers atomic.Int32 // number of pending readers. Atomic so RUnlock doesn't need to be locked
}

// NewRWMutex creates a new RWMutex.
func NewRWMutex() *RWMutex {
	return &RWMutex{
		Mutex: Mutex{
			lock: make(chan struct{}, 1),
		},
		wait: Mutex{
			lock: make(chan struct{}, 1),
		},
	}
}

// RLock locks rw for reading.
// If the lock is already locked for writing, RLock blocks until the lock is available.
// It is allowed for several readers to hold the lock simultaneously.
func (rw *RWMutex) RLock(ctx context.Context) error {
	if err := rw.wait.Lock(ctx); err != nil {
		return err
	}
	defer rw.wait.Unlock()
	if rw.readers.Load() == 0 {
		if err := rw.Mutex.Lock(ctx); err != nil {
			return err
		}
	}
	rw.readers.Add(1)
	return nil
}

// RUnlock undoes a single RLock call; it does not affect other simultaneous readers.
// It is a design error if rw is not locked for reading on entry to RUnlock.
func (rw *RWMutex) RUnlock() {
	readers := rw.readers.Add(-1)
	if readers < 0 {
		panic("RUnlock of unlocked RWMutex")
	}
	if readers == 0 {
		rw.Mutex.Unlock()
	}
}

// RLocked returns true if rw is locked for reading.
func (rw *RWMutex) RLocked() bool {
	return rw.readers.Load() > 0
}

// Readers returns the number of readers.
func (rw *RWMutex) Readers() int {
	return int(rw.readers.Load())
}
