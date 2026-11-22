package main

import "fmt"

func main() {
	reg := loadRegistry()
	// TODO: remove brine-2, which never left its deployer.
	_, listed := reg["brine-2"]
	fmt.Println("satellites:", len(reg))
	fmt.Println("brine-2 listed:", listed)
}
