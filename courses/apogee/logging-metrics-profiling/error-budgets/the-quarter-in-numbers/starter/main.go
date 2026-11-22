package main

import "fmt"

// quarter is one station's counts for the last quarter of 2028, with its target.
type quarter struct {
	code         string
	booked, lost int
	target       Target
}

func main() {
	q4 := []quarter{
		{"YST", 2208, 2, Standard},
		{"PEL", 1104, 1, Standard},
		{"SQT", 1288, 3, Standard},
		{"LAN", 1199, 22, Atoll},
		{"CWR", 920, 4, Standard},
	}
	fmt.Printf("%-7s %7s %5s %6s %7s %5s\n", "Q4 2028", "booked", "lost", "avail", "budget", "left")
	for _, q := range q4 {
		fmt.Println(row(q.code, q.booked, q.lost, q.target))
	}
}
