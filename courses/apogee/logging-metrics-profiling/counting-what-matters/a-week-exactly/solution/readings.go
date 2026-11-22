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
	var total uint64
	for i := 1; i < len(rs); i++ {
		if rs[i].Value < rs[i-1].Value {
			total += rs[i].Value
		} else {
			total += rs[i].Value - rs[i-1].Value
		}
	}
	return total
}
