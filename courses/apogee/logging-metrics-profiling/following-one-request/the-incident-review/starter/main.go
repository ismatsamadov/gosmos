package main

import "fmt"

func main() {
	for _, line := range timeline(benchReplay) {
		fmt.Println(line)
	}
}
