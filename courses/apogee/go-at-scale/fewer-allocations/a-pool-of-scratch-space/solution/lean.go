package main

import "sync"

// maxScratch is the largest capacity of a scratch putScratch keeps.
const maxScratch = 65536

// scratch holds the scratch slices nothing is using.
var scratch = sync.Pool{
	New: func() any {
		s := make([]PassView, 0, 512)
		return &s
	},
}

// getScratch returns an empty slice to gather a station's candidates
// in, with room for any station's day.
func getScratch() *[]PassView {
	p := scratch.Get().(*[]PassView)
	*p = (*p)[:0]
	return p
}

// putScratch hands back a scratch from getScratch that nothing uses
// any more.
func putScratch(p *[]PassView) {
	if cap(*p) <= maxScratch {
		scratch.Put(p)
	}
}
