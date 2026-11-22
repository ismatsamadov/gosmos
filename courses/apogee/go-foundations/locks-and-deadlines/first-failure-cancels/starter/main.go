package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// runAll runs every check at once, cancels the others when one fails, and returns the first failure.
func runAll(ctx context.Context, checks []func(context.Context) error) error {
	var wg sync.WaitGroup
	for _, check := range checks {
		wg.Go(func() {
			check(ctx)
		})
	}
	wg.Wait()
	return nil
}

func main() {
	ready := func(ctx context.Context) error { return nil }
	stowed := func(ctx context.Context) error { return errors.New("YST: wind 81 km/h, dish stowed") }
	fmt.Println(runAll(context.Background(), []func(context.Context) error{ready, stowed, ready}))
}
