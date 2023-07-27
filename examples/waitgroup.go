package main

import (
	"context"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/metalim/csync"
)

// imagine this is a long running task from a 3rd party library
// and we can't pass context into it
func work() {
	time.Sleep(time.Duration(rand.Intn(7)) * time.Second)
}

const workers = 3

func main() {
	var wg csync.WaitGroup
	wg.Add(workers)
	var done atomic.Int32
	for i := 0; i < workers; i++ {
		go func(i int) {
			defer wg.Done()
			log.Printf("Worker %d started", i)
			work()
			log.Printf("Worker %d finished", i)
			done.Add(1)
		}(i)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := wg.Wait(ctx)
	if err != nil {
		log.Println(err)
	}

	log.Printf("%d / %d workers finished", done.Load(), workers)
	log.Println("Moving on...")
	// ...
}
