package main

import "fmt"

func main() {
	reg := loadRegistry()
	delete(reg, "brine-2")
	_, listed := reg["brine-2"]
	fmt.Println("satellites:", len(reg))
	fmt.Println("brine-2 listed:", listed)
}
