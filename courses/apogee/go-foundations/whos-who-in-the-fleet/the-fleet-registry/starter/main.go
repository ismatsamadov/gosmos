package main

import "fmt"

func main() {
	operators := map[string]string{
		"pipit-1":   "Apogee",
		"pipit-2":   "Apogee",
		"pipit-3":   "Apogee",
		"furrow-a":  "Furrowlight",
		"furrow-b":  "Furrowlight",
		"brine-3":   "Brinewatch",
		"cirrine-1": "Cirrine Atmospherics",
	}
	fmt.Println("pipit-1:", operators["pipit-1"])
	fmt.Println("furrow-b:", "?")
	fmt.Println("thistle:", "?")
}
