package main

import "fmt"

// Station is one of Apogee's ground stations.
type Station struct {
	Code    string
	Name    string
	MaskDeg float64
}

func main() {
	live := Station{Code: "CWR", Name: "Cape Wren", MaskDeg: 5}
	draft := live
	// TODO: try a 12-degree mask on the draft.
	fmt.Println("live mask:", live.MaskDeg)
	fmt.Println("draft mask:", draft.MaskDeg)
}
