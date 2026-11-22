package main

import "time"

// Manifest is a parcel's manifest.xml: the dish it is about, when the controller made it and the files it holds.
type Manifest struct {
	// TODO: tag every field here and in Entry with where it sits in manifest.xml.
	Version int
	Dish    string
	Made    time.Time
	Files   []Entry
}

// Entry is one file a manifest lists: its name in the parcel and its size in bytes.
type Entry struct {
	Name string
	Size int64
}
