package main

import "time"

// Manifest is a parcel's manifest.xml: the dish it is about, when the controller made it and the files it holds.
type Manifest struct {
	Version int       `xml:"version,attr"`
	Dish    string    `xml:"dish"`
	Made    time.Time `xml:"made"`
	Files   []Entry   `xml:"files>file"`
}

// Entry is one file a manifest lists: its name in the parcel and its size in bytes.
type Entry struct {
	Name string `xml:",chardata"`
	Size int64  `xml:"size,attr"`
}
