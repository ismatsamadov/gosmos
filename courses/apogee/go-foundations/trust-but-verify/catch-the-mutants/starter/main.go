package main

import "fmt"

// accept reports whether a command numbered seq may follow last, the newest one accepted.
func accept(seq, last uint64) bool {
	return seq > last
}

func main() {
	fmt.Println(accept(119, 118))
}
