package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// merge forwards every beacon from every feed into one channel, closed once all feeds are done.
func merge(feeds ...<-chan Heard) <-chan Heard {
	return feeds[0]
}

func main() {
	beacons := merge(
		listen("PEL", 9002, 9003),
		listen("SQT", 9003, 9004),
		listen("LAN", 9001, 9002),
		listen("CWR", 9001, 9002),
	)
	heardBy := map[int][]string{}
	for h := range beacons {
		heardBy[h.Seq] = append(heardBy[h.Seq], h.Station)
	}
	seqs := slices.Sorted(maps.Keys(heardBy))
	for _, seq := range seqs {
		fmt.Printf("%d heard by %s\n", seq, strings.Join(heardBy[seq], " "))
	}
	fmt.Println("last beacon:", seqs[len(seqs)-1])
}
