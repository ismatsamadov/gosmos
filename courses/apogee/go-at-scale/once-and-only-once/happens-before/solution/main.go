package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

func main() {
	var (
		plan   Schedule    // what the planner fills in
		ready  atomic.Bool // raised once the plan can be read
		looked = make(chan struct{})
		done   = make(chan struct{})
	)
	go func() { // the planner
		defer close(done)
		plan.Version = 7
		plan.Contacts = make([]Contact, 2016)
		ready.Store(true)
		<-looked // so every run goes the same way, wait until the Board has looked
	}()
	for !ready.Load() { // the test Board waits for the flag
		runtime.Gosched()
	}
	fmt.Printf("plan v%d: %d contacts\n", plan.Version, len(plan.Contacts))
	close(looked)
	<-done
}
