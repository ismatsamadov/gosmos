package main

// maxScratch is the largest capacity of a scratch putScratch keeps.
const maxScratch = 65536

// getScratch returns an empty slice to gather a station's candidates
// in, with room for any station's day.
func getScratch() *[]PassView {
	s := make([]PassView, 0, 512)
	return &s
}

// putScratch hands back a scratch from getScratch that nothing uses
// any more.
func putScratch(p *[]PassView) {}
