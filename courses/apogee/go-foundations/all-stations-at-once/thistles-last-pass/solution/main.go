package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"
)

// merge forwards every beacon from every feed into one channel, closed once all feeds are done.
func merge(feeds ...<-chan Heard) <-chan Heard {
	out := make(chan Heard)
	var wg sync.WaitGroup
	for _, feed := range feeds {
		wg.Go(func() {
			for h := range feed {
				out <- h
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
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
		slices.Sort(heardBy[seq])
		fmt.Printf("%d heard by %s\n", seq, strings.Join(heardBy[seq], " "))
	}
	fmt.Println("last beacon:", seqs[len(seqs)-1])
}
