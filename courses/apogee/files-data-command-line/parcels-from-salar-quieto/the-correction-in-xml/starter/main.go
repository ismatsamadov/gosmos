package main

import (
	"encoding/xml"
	"fmt"
)

// Correction is what a dish controller adds to every pointing of one dish, in degrees, as correction.xml holds it.
type Correction struct {
	Dish        string  `xml:"dish,attr"`
	AzOffsetDeg float64 `xml:"az_offset_deg"`
	ElOffsetDeg float64 `xml:"el_offset_deg"`
}

func main() {
	c := Correction{Dish: "LG", AzOffsetDeg: 0.18, ElOffsetDeg: -0.06}
	// TODO: name the root element, indent by two spaces and print xml.Header first.
	out, err := xml.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(out))
}
