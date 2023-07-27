package csync

import (
	"context"
	"sync/atomic"
)

// same as sync.WaitGroup but with a context that can be cancelled
type WaitGroup struct {
	state atomic.Int32
	done  atomic.Pointer[chan struct{}]
}

func New() *WaitGroup {
	return &WaitGroup{}
}

func (wg *WaitGroup) Add(delta int) {
	state := wg.state.Add(int32(delta))
	if state < 0 {
		panic("sync: negative WaitGroup counter")
	}
	switch state {
	case 0:
		close(*wg.done.Load())
	case int32(delta):
		ch := make(chan struct{}) // NB! not new(chan struct{})
		wg.done.Store(&ch)
	}
}

func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

func (wg *WaitGroup) Wait(ctx context.Context) error {
	ptr := wg.done.Load()
	if ptr == nil {
		return nil
	}
	select {
	case <-*ptr:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
