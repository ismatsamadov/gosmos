package main

import (
	"fmt"

	"apogee/groundctl/uplink"
)

func main() {
	fmt.Println("signing with key", uplink.KeyID())
}
