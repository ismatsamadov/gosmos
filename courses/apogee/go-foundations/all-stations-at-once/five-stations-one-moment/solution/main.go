package main

import (
	"fmt"
	"sync"
	"time"
)

// stations are the five ground stations, in the order the Board lists them.
var stations = []string{"YST", "PEL", "SQT", "LAN", "CWR"}

// waiting is how many frames each station holds right now.
var waiting = map[string]int{"YST": 3, "PEL": 2, "SQT": 4, "LAN": 2, "CWR": 5}

// poll asks one station for its waiting frames; the answer takes 20 ms to come back.
func poll(code string) int {
	time.Sleep(20 * time.Millisecond)
	return waiting[code]
}

func main() {
	counts := make([]int, len(stations))
	var wg sync.WaitGroup
	for i, code := range stations {
		wg.Go(func() {
			counts[i] = poll(code)
		})
	}
	wg.Wait()
	for i, code := range stations {
		fmt.Printf("%s: %d frames\n", code, counts[i])
	}
}
