package main

import "fmt"

func main() {
	for station := range 5 {
		fmt.Print(stationCode(station), " ")
		for block := range 12 {
			if visible(station, block) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
