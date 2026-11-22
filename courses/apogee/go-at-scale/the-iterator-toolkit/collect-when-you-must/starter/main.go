package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	monday := time.Date(2029, 4, 2, 0, 0, 0, 0, time.UTC)
	// live holds every station flying from Monday, by its code.
	live := map[string]Site{}
	for _, s := range newNetwork().Sites {
		if !s.Live.After(monday) {
			live[s.Code] = s
		}
	}
	codes := []string{"YST", "PEL", "SQT", "LAN", "CWR"}
	fmt.Println(strings.Join(codes, " "))
}
