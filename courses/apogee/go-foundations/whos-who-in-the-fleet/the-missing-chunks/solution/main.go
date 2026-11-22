package main

import "fmt"

// missing returns, in increasing order, the chunk numbers from 0 to total-1 that chunks lacks.
func missing(chunks map[int]string, total int) []int {
	var gaps []int
	for n := range total {
		if _, ok := chunks[n]; !ok {
			gaps = append(gaps, n)
		}
	}
	return gaps
}

func main() {
	chunks := map[int]string{0: "img-0", 1: "img-1", 3: "img-3"}
	fmt.Println("missing:", missing(chunks, 5))
}
