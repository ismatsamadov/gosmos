package main

import (
	"fmt"
	"hash/crc32"
	"sync"
)

// chunkOK reports whether chunk's CRC-32 matches sum.
func chunkOK(chunk []byte, sum uint32) bool {
	return crc32.ChecksumIEEE(chunk) == sum
}

// checkAll checks each chunk against its sum with workers goroutines, and reports each result by index.
func checkAll(chunks [][]byte, sums []uint32, workers int) []bool {
	ok := make([]bool, len(chunks))
	jobs := make(chan int)
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for i := range jobs {
				ok[i] = chunkOK(chunks[i], sums[i])
			}
		})
	}
	for i := range chunks {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return ok
}

func main() {
	chunks := [][]byte{[]byte("furrow-a img-0623 chunk 0000"), []byte("furrow-a img-0623 chunk 0001")}
	sums := []uint32{crc32.ChecksumIEEE(chunks[0]), 0}
	fmt.Println(checkAll(chunks, sums, 2))
}
