package main

import (
	"encoding/xml"
	"fmt"
)

// Correction is what a dish controller adds to every pointing of one dish, in degrees, as correction.xml holds it.
type Correction struct {
	XMLName     xml.Name `xml:"correction"`
	Dish        string   `xml:"dish,attr"`
	AzOffsetDeg float64  `xml:"az_offset_deg"`
	ElOffsetDeg float64  `xml:"el_offset_deg"`
}

func main() {
	c := Correction{Dish: "LG", AzOffsetDeg: 0.18, ElOffsetDeg: -0.06}
	out, err := xml.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(xml.Header)
	fmt.Println(string(out))
}
