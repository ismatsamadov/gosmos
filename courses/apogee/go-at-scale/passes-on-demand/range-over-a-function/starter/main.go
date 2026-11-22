package main

import (
	"fmt"
	"time"
)

func main() {
	july := time.Date(2029, 7, 1, 0, 0, 0, 0, time.UTC)
	for day := range Days(july, 7) {
		fmt.Println(day.Format("Mon 2 Jan"))
	}
}
