package main

import "fmt"

// Image is a picture from a payload camera.
type Image struct {
	SatID  string
	SizeMB float64
}

// Housekeeping is a frame of battery and temperature readings.
type Housekeeping struct {
	SatID string
	BattV float64
}

// Chunk is one numbered piece of a downlinked file.
type Chunk struct {
	File string
	N    int
}

// route returns where a downlinked item goes next.
func route(item any) string {
	return "quarantine"
}

// show prints an item's type and where route sends it.
func show(item any) {
	fmt.Printf("%T -> %s\n", item, route(item))
}

func main() {
	show(Image{SatID: "furrow-a", SizeMB: 412})
	show(Housekeeping{SatID: "pipit-1", BattV: 7.91})
	show(Chunk{File: "furrow-a-0412.tif", N: 17})
	show("stray line from the modem")
}
