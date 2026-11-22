package main

import "fmt"

func main() {
	for minute := range 10 {
		fmt.Printf("minute %d: %.1f\n", minute, elevationAt(minute))
	}
}
