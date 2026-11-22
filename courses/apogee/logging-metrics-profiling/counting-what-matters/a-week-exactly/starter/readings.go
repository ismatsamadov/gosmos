package main

import "time"

// Reading is one value of a counter, read by the pager every five minutes.
type Reading struct {
	At    time.Time
	Value uint64
}

// increase returns how much a counter went up across rs, oldest first. A drop means its
// program restarted, and the new value is what it counted since.
func increase(rs []Reading) uint64 {
	if len(rs) < 2 {
		return 0
	}
	return rs[len(rs)-1].Value - rs[0].Value
}
