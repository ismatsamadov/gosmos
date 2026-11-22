package main

import "fmt"

func main() {
	for minutes := 5; minutes >= 1; minutes-- {
		fmt.Printf("T-%d\n", minutes)
	}
	fmt.Println("AOS")
}
