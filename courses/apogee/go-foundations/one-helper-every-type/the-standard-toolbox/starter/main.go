package main

import "fmt"

func main() {
	codes := []string{"YST", "PEL", "SQT", "LAN", "CWR"}
	masksDeg := []float64{5, 5, 3, 8, 5}

	// TODO: replace each loop with one call from the slices package.
	hasSQT := false
	for _, c := range codes {
		hasSQT = c == "SQT"
	}
	fmt.Println("has SQT:", hasSQT)

	indexLAN := 0
	for _, c := range codes {
		indexLAN++
		if c == "LAN" {
			break
		}
	}
	fmt.Println("index of LAN:", indexLAN)

	highest := masksDeg[0]
	for _, m := range masksDeg {
		if m < highest {
			highest = m
		}
	}
	fmt.Println("highest mask:", highest)
}
