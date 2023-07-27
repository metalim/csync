package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/metalim/csync"
)

func main() {
	var wg csync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		}()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := wg.Wait(ctx)
	log.Println(err)
}
