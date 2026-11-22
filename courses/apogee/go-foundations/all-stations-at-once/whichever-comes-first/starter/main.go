package main

import "fmt"

// first returns the first beacon heard at either Lantern Atoll or Cape Wren.
func first(lan, cwr <-chan string) string {
	return <-lan
}

func main() {
	lan, cwr := make(chan string, 1), make(chan string, 1)
	lan <- "thistle 9001 LAN"
	fmt.Println(first(lan, cwr))
}
