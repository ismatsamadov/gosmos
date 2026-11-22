package main

import (
	"fmt"
	"slices"
)

func main() {
	live := []string{"04:12", "05:03", "06:47"} // AOS times on the Board
	draft := slices.Clone(live[:2])
	draft[1] = "05:30" // try Furrow-A's pass at 05:30
	fmt.Println("live:", live)
	fmt.Println("draft:", draft)
}
