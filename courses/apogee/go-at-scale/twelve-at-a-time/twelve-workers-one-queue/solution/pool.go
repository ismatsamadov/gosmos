package main

import (
	"context"
	"sync"
)

// runPool calls do for every job in jobs, in workers goroutines, and returns the first error once every job is done.
func runPool(ctx context.Context, jobs <-chan Elements, workers int, do func(context.Context, Elements) error) error {
	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		first error
	)
	// fail keeps the first error a job returns.
	fail := func(err error) {
		mu.Lock()
		defer mu.Unlock()
		if first == nil {
			first = err
		}
	}
	for range workers {
		wg.Go(func() {
			for e := range jobs {
				if err := do(ctx, e); err != nil {
					fail(err)
				}
			}
		})
	}
	wg.Wait()
	return first
}
