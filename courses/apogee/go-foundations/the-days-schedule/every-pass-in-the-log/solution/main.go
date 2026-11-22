package main

import "fmt"

func main() {
	ids := []string{"pipit-1/PEL/04182", "furrow-a/SQT/11873", "thistle/LAN/02218"}
	for i, id := range ids {
		fmt.Println(i+1, id)
	}
}
