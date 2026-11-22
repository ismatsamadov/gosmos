package main

import "fmt"

// averageMB returns the average data volume of a pass, in megabytes.
func averageMB(totalMB, passes int) float64 {
	return float64(totalMB) / float64(passes)
}

func main() {
	fmt.Println("Furrowlight MB per pass:", averageMB(1000, 3))
}
