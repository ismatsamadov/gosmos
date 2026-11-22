package main

import "fmt"

func main() {
	live := []string{"04:12", "05:03", "06:47"} // AOS times on the Board
	// TODO: give the draft its own copy of the first two times.
	draft := live[:2]
	draft[1] = "05:30" // try Furrow-A's pass at 05:30
	fmt.Println("live:", live)
	fmt.Println("draft:", draft)
}
