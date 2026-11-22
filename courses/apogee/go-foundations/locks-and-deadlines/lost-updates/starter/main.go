package main

import (
	"fmt"
	"sync"
)

// Counter counts frames per satellite, for many feeds at once.
type Counter struct {
	mu     sync.Mutex
	frames map[string]int
}

// NewCounter returns an empty Counter.
func NewCounter() *Counter {
	return &Counter{frames: map[string]int{}}
}

// Add counts one frame from sat.
func (c *Counter) Add(sat string) {
}

// Get returns how many frames sat has sent.
func (c *Counter) Get(sat string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.frames[sat]
}

func main() {
	c := NewCounter()
	c.Add("pipit-2")
	c.Add("pipit-2")
	fmt.Println("pipit-2:", c.Get("pipit-2"), "frames")
}
