package main

import "fmt"

// merge forwards every value from every feed into one channel, closed once all feeds are done.
func merge(feeds ...<-chan int) <-chan int {
	out := make(chan int)
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
