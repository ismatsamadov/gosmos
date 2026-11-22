package main

import "fmt"

// stationCode returns the code the scheduler uses for Port Ellery.
func stationCode() string {
	return "PEL"
}

func main() {
	fmt.Println(stationCode())
}
