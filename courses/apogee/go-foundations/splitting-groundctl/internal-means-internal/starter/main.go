package main

import (
	"fmt"

	"apogee/groundctl/uplink/internal/keys"
)

func main() {
	fmt.Println("signing with key", keys.Current().ID)
}
