package main

import (
	"context"
	"fmt"
)

// receive archives chunks until the channel closes or ctx ends, and returns how many it archived.
func receive(ctx context.Context, chunks <-chan []byte) (int, error) {
	n := 0
	for {
		select {
		case <-ctx.Done():
			return n, ctx.Err()
		case _, ok := <-chunks:
			if !ok {
				return n, nil
			}
			n++ // the chunk goes to the archive here
		}
	}
}

func main() {
	chunks := make(chan []byte, 2)
	chunks <- []byte("furrow-b chunk 0")
	chunks <- []byte("furrow-b chunk 1")
	close(chunks)
	fmt.Println(receive(context.Background(), chunks))
}
