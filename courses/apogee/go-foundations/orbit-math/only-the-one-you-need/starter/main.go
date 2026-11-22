package main

import "fmt"

// hoursMinutes splits minutes after midnight into hours and minutes.
func hoursMinutes(total int) (int, int) {
	return total / 60, total % 60
}

func main() {
	h, m := hoursMinutes(252)
	fmt.Printf("Minutes past the hour: %d\n", m)
}
