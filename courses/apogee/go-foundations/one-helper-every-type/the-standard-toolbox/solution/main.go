package main

import (
	"fmt"
	"slices"
)

func main() {
	codes := []string{"YST", "PEL", "SQT", "LAN", "CWR"}
	masksDeg := []float64{5, 5, 3, 8, 5}

	hasSQT := slices.Contains(codes, "SQT")
	fmt.Println("has SQT:", hasSQT)

	indexLAN := slices.Index(codes, "LAN")
	fmt.Println("index of LAN:", indexLAN)

	highest := slices.Max(masksDeg)
	fmt.Println("highest mask:", highest)
}
