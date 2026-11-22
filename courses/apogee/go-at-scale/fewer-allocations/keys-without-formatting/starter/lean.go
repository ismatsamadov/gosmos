package main

import "fmt"

// passKey is a pass's identity without formatting its ID.
type passKey struct {
	SatID, Station string
	AOS            int64 // Unix seconds
}

// dedupe returns passes without repeats, each where it first came. A
// pass is its satellite, its station and its AOS.
func dedupe(passes []PassView) []PassView {
	seen := map[string]bool{}
	var out []PassView
	for _, p := range passes {
		key := fmt.Sprintf("%s %s %d", p.SatID, p.Station, p.AOS.Unix())
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, p)
	}
	return out
}
