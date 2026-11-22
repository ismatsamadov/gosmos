package main

import "fmt"

func main() {
	for _, line := range fieldLines(Config{}) {
		fmt.Println(line)
	}
}
