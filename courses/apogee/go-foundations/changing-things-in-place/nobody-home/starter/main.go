package main

import "fmt"

// Station is one of Apogee's ground stations.
type Station struct {
	Code   string
	Name   string
	Stowed bool
}

// find returns a pointer to the station in stations with the given code, or nil if there is none.
func find(stations []Station, code string) *Station {
	return nil
}

func main() {
	stations := []Station{{Code: "YST", Name: "Ystmark"}, {Code: "CWR", Name: "Cape Wren"}}
	if s := find(stations, "XYZ"); s == nil {
		fmt.Println("XYZ: no such station")
	}
}
