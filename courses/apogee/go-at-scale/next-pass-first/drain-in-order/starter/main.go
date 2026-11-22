package main

import (
	"fmt"
	"time"
)

func main() {
	n := newNetwork()
	cwr := n.Sites[4] // Cape Wren, fifth in table order
	monday := time.Date(2029, 4, 9, 0, 0, 0, 0, time.UTC)
	tuesday := monday.AddDate(0, 0, 1)
	q := NewQueue(byAOS)
	for i, e := range n.Elements {
		if !canServe(cwr, n.Satellites[i]) {
			continue
		}
		for _, p := range predict(e, cwr, monday, tuesday) {
			if p.MaxElevDeg >= 10 {
				q.Push(p)
			}
		}
	}
	shown := 0
	for p := range q.Drain() {
		fmt.Println(p.AOS.Format("15:04:05"), p.ID)
		if shown++; shown == 5 {
			break
		}
	}
}
