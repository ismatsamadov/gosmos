package main

import "fmt"

func main() {
	for station := range 5 {
		fmt.Print(stationCode(station), " ")
		// TODO: print # or . for each of the 12 blocks
		fmt.Println()
	}
}
