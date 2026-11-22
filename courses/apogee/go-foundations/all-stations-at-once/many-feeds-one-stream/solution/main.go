package main

import (
	"fmt"
	"sync"
)

// merge forwards every value from every feed into one channel, closed once all feeds are done.
func merge(feeds ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, f := range feeds {
		wg.Go(func() {
			for seq := range f {
				out <- seq
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	cwr := make(chan int, 2)
	cwr <- 9001
	cwr <- 9002
	close(cwr)
	for seq := range merge(cwr) {
		fmt.Println("archived", seq)
	}
}
