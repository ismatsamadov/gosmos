package main

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// contact returns one of take-over day's contacts as passplan 2 plans it: the pass ID, AOS and LOS as hh:mm:ss on Thursday 28 June, the max elevation and the dish.
func contact(id, aos, los string, elevDeg float64, dish int) Contact {
	at := func(hms string) time.Time {
		t, _ := time.Parse(time.DateTime, "2029-06-28 "+hms)
		return t
	}
	parts := strings.Split(id, "/")
	p := Pass{SatID: parts[0], Station: parts[1], AOS: at(aos), LOS: at(los), MaxElevDeg: elevDeg}
	return Contact{PassView: PassView{ID: id, Pass: p}, Dish: dish}
}

// takeOver is the first five contacts of take-over day, as the test Board shows them.
var takeOver = []Contact{
	contact("cirrine-4/TIR/01069", "00:01:04", "00:11:04", 41.0, 2),
	contact("cirrine-4/YST/01069", "00:03:43", "00:12:51", 23.6, 2),
	contact("cirrine-11/PVT/01069", "00:05:16", "00:12:28", 13.3, 1),
	contact("furrow-x/CSL/00004", "00:05:47", "00:14:57", 61.5, 1),
	contact("cirrine-6/HSK/01069", "00:06:06", "00:16:25", 37.8, 1),
}

// ids returns the pass IDs of contacts, in order.
func ids(contacts []Contact) []string {
	out := make([]string, len(contacts))
	for i, c := range contacts {
		out[i] = c.ID
	}
	return out
}

func TestUpdate(t *testing.T) {
	kept := t.Run("the test Board holds version 1 while version 2 is published", func(t *testing.T) {
		var p Publisher
		p.Update(takeOver[:2])
		held := p.Current()
		p.Update(takeOver)
		if held.Version != 1 || !slices.Equal(ids(held.Contacts), ids(takeOver[:2])) {
			t.Errorf("the schedule the Board held became version %d with %d contacts, want version 1 with its 2, unchanged", held.Version, len(held.Contacts))
		}
	})
	t.Run("the Board's next look gets version 2", func(t *testing.T) {
		var p Publisher
		p.Update(takeOver[:2])
		second := p.Update(takeOver)
		cur := p.Current()
		if cur != second {
			t.Fatalf("Current() after Update is not the schedule Update returned")
		}
		if cur.Version != 2 || !slices.Equal(ids(cur.Contacts), ids(takeOver)) {
			t.Errorf("Current() after two Updates = version %d with %d contacts, want version 2 with 5", cur.Version, len(cur.Contacts))
		}
	})
	t.Run("a plan's contacts are its own", func(t *testing.T) {
		var p Publisher
		contacts := slices.Clone(takeOver)
		s := p.Update(contacts)
		contacts[0].Dish = 1 // the caller reuses its slice
		if s.Contacts[0].Dish != 2 {
			t.Errorf("Update(contacts)'s first contact changed to dish %d with the caller's slice, want dish 2: a slice of its own", s.Contacts[0].Dish)
		}
	})
	t.Run("an empty publisher", func(t *testing.T) {
		var p Publisher
		if s := p.Current(); s != nil {
			t.Errorf("Current() before any Update = version %d, want nil", s.Version)
		}
	})
	if !kept {
		return // readers of a schedule that changes under them would only fail the same way
	}
	t.Run("eight screens read while a hundred plans are published", func(t *testing.T) {
		var p Publisher
		var mu sync.Mutex
		var seen string // the first schedule a reader saw half-made
		stop := make(chan struct{})
		var wg sync.WaitGroup
		for range 8 {
			wg.Go(func() {
				for {
					select {
					case <-stop:
						return
					default:
					}
					if s := p.Current(); s != nil && len(s.Contacts) != int(s.Version) {
						mu.Lock()
						if seen == "" {
							seen = fmt.Sprintf("version %d with %d contacts", s.Version, len(s.Contacts))
						}
						mu.Unlock()
					}
				}
			})
		}
		for n := 1; n <= 100; n++ {
			p.Update(make([]Contact, n)) // version n has n contacts
		}
		close(stop)
		wg.Wait()
		if seen != "" {
			t.Errorf("a screen read %s, want version n with n contacts, whole", seen)
		}
	})
}
