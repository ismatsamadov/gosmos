package main

import "fmt"

func main() {
	for _, code := range []int{200, 201, 404, 413} {
		fmt.Printf("%q\n", statusLine(code))
	}
}
