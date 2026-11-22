package main

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

// loaderLine is a contact's dish loader line, made with fmt.Sprintf
// the old way.
func loaderLine(c Contact) string {
	aos, los := c.AOS.Format("15:04:05"), c.LOS.Format("15:04:05")
	return fmt.Sprintf("%s %s %d %s\n", aos, los, c.Dish, c.ID)
}

func TestAppendContact(t *testing.T) {
	draft := takeOverDraft()
	t.Run("take-over day's first 50 contacts", func(t *testing.T) {
		for i, c := range draft[:50] {
			got, want := string(appendContact(nil, c)), loaderLine(c)
			if got != want {
				t.Fatalf("appendContact(nil, contact %d) = %q, want %q", i+1, got, want)
			}
		}
	})
	t.Run("Hollis's pass over Varsk Island", func(t *testing.T) {
		at := slices.IndexFunc(draft, func(c Contact) bool {
			return c.ID == "furrow-a/VSK/24529"
		})
		c := draft[at]
		if got, want := string(appendContact(nil, c)), loaderLine(c); got != want {
			t.Errorf("appendContact(nil, %s) = %q, want %q", c.ID, got, want)
		}
	})
	t.Run("after the day's first line", func(t *testing.T) {
		first := loaderLine(draft[0])
		dst := append(make([]byte, 0, 128), first...)
		got := string(appendContact(dst, draft[1]))
		added, ok := strings.CutPrefix(got, first)
		if !ok {
			t.Fatalf("appendContact(dst, contact 2) lost the line dst held: it starts %.20q, want %.20q", got, first)
		}
		if want := loaderLine(draft[1]); added != want {
			t.Errorf("appendContact(dst, contact 2) added %q, want %q", added, want)
		}
	})
}

func TestAppendContactAllocs(t *testing.T) {
	if raceOn {
		t.Skip("allocation counts differ under -race")
	}
	c := takeOverDraft()[0]
	buf := make([]byte, 0, 64)
	allocs := testing.AllocsPerRun(100, func() {
		buf = appendContact(buf[:0], c)
	})
	if string(buf) != loaderLine(c) {
		return // TestAppendContact says what is wrong
	}
	if allocs != 0 {
		t.Errorf("appendContact(buf[:0], contact 1) made %.0f allocations with room for the line in buf, want 0", allocs)
	}
}
