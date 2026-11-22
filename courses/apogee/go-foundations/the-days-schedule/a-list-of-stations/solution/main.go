package main

import "fmt"

func main() {
	codes := []string{"YST", "PEL", "SQT", "LAN", "CWR"}
	fmt.Println("First:", codes[0])
	fmt.Println("Last:", codes[len(codes)-1])
	fmt.Println("Count:", len(codes))
}
