# csync
Go sync package with cancellable context

## Rationale

`sync.Mutex`, `sync.RWMutex` and `sync.WaitGroup` can block execution until another goroutine releases them. However locking can happen for a long time and we live in dynamic world, where we need to move forward even if some resources stay locked for a long time. Hence the creation of `csync` — package with similar primitives, where waiting can be cancelled by context.

## Primitives

`csync.WaitGroup`
`csync.Mutex`
`csync.RWMutex`

## Usage

```go
package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/metalim/csync"
)

const workers = 3

func main() {
	var wg csync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			// do some work in 3rd party library without context support
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		}()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := wg.Wait(ctx)
	if err != nil {
		log.Println(err)
	}
}
```
