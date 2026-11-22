package main

import "fmt"

// clock returns a time of day as hh:mm.
func clock(h, m int) string {
	return fmt.Sprintf("%02d:%02d", h, m)
}

func main() {
	fmt.Println("AOS", clock(4, 7))
}
