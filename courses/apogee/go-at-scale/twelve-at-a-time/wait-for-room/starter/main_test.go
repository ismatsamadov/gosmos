package main

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

// guard is how long a test waits for what a right stage does at once: one second in all.
type guard struct{ end time.Time }

// newGuard starts a guard's second.
func newGuard() guard { return guard{end: time.Now().Add(time.Second)} }

// wait is ready when the guard's second is over, and never sooner than 50 ms from now.
func (g guard) wait() <-chan time.Time {
	return time.After(max(time.Until(g.end), 50*time.Millisecond))
}

// contact returns one of Varsk Island's contacts on Monday 2 April, its first day live: its one dish, the pass ID, AOS and LOS as hh:mm:ss, and the max elevation.
func contact(id, aos, los string, elevDeg float64) Contact {
	at := func(hms string) time.Time {
		t, _ := time.Parse(time.DateTime, "2029-04-02 "+hms)
		return t
	}
	sat, _, _ := strings.Cut(id, "/")
	a, l := at(aos), at(los)
	if l.Before(a) {
		l = l.Add(24 * time.Hour) // the pass ends after midnight
	}
	return Contact{PassView: PassView{ID: id, Pass: Pass{SatID: sat, Station: "VSK", AOS: a, LOS: l, MaxElevDeg: elevDeg}}, Dish: 1}
}

// varsk is Varsk Island's 22 contacts on Monday 2 April, in AOS order.
var varsk = []Contact{
	contact("pipit-2/VSK/21125", "00:06:52", "00:12:19", 12.7),
	contact("pipit-2/VSK/21126", "01:39:05", "01:48:01", 85.9),
	contact("cirrine-1/VSK/18566", "03:13:20", "03:19:22", 13.5),
	contact("furrow-a/VSK/23205", "04:26:40", "04:33:03", 17.3),
	contact("cirrine-1/VSK/18567", "04:46:30", "04:56:00", 89.8),
	contact("furrow-b/VSK/23203", "05:13:02", "05:20:57", 35.8),
	contact("furrow-a/VSK/23206", "05:59:08", "06:07:26", 50.3),
	contact("cirrine-1/VSK/18568", "06:24:46", "06:30:01", 11.2),
	contact("furrow-b/VSK/23204", "06:47:17", "06:54:25", 21.6),
	contact("pipit-3/VSK/11078", "08:04:32", "08:10:55", 16.7),
	contact("pipit-3/VSK/11079", "09:37:38", "09:45:50", 47.4),
	contact("pipit-2/VSK/21134", "14:43:35", "14:52:29", 72.5),
	contact("pipit-2/VSK/21135", "16:19:03", "16:25:06", 14.5),
	contact("furrow-a/VSK/23213", "17:24:28", "17:28:59", 10.7),
	contact("cirrine-1/VSK/18575", "17:59:25", "18:08:55", 81.3),
	contact("pipit-3/VSK/11085", "18:47:48", "18:54:07", 17.0),
	contact("furrow-a/VSK/23214", "18:56:45", "19:05:11", 77.7),
	contact("cirrine-1/VSK/18576", "19:35:56", "19:42:21", 14.7),
	contact("furrow-b/VSK/23212", "19:44:21", "19:52:08", 32.3),
	contact("pipit-3/VSK/11086", "20:20:13", "20:28:32", 51.5),
	contact("brine-3/VSK/16616", "21:16:13", "21:23:29", 23.2),
	contact("pipit-2/VSK/21140", "23:59:56", "00:04:39", 10.9),
}

// queued takes every contact waiting in q, without waiting for more.
func queued(q chan Contact) []Contact {
	var got []Contact
	for {
		select {
		case c := <-q:
			got = append(got, c)
		default:
			return got
		}
	}
}

// checkOrder fails the test unless got holds want's contacts, in want's order.
func checkOrder(t *testing.T, call string, got, want []Contact) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s handed over %d contacts, want %d", call, len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("%s contact %d = %s, want %s, in AOS order", call, i, got[i].ID, want[i].ID)
		}
	}
}

func TestFeed(t *testing.T) {
	g := newGuard()
	t.Run("Varsk Island's first day, a loader that keeps up", func(t *testing.T) {
		q := make(chan Contact, 256)
		done := make(chan error, 1)
		go func() { done <- feed(context.Background(), varsk, q) }()
		select {
		case err := <-done:
			if err != nil {
				t.Fatalf("feed(22 contacts, room for 256) = %v, want <nil>", err)
			}
		case <-g.wait():
			t.Fatalf("feed(22 contacts, room for 256) is still waiting, want it to return once all are queued")
		}
		if n := len(q); n != 22 {
			t.Fatalf("feed returned with %d of 22 contacts in the queue, want all 22 handed over", n)
		}
		checkOrder(t, "feed(22 contacts, room for 256)", queued(q), varsk)
	})
	t.Run("Varsk Island's loader stalls, with room for 4", func(t *testing.T) {
		q := make(chan Contact, 4)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		done := make(chan error, 1)
		go func() { done <- feed(ctx, varsk, q) }()
		for len(q) < 4 {
			select {
			case err := <-done:
				t.Fatalf("feed(22 contacts, a stalled loader) = %v before the cancel, want it to wait for room", err)
			case <-g.wait():
				t.Fatalf("feed(22 contacts, room for 4) queued %d contacts, then nothing", len(q))
			case <-time.After(time.Millisecond):
			}
		}
		select {
		case err := <-done:
			t.Fatalf("feed(22 contacts, a stalled loader) = %v before the cancel, want it to wait for room", err)
		default:
		}
		cancel()
		select {
		case err := <-done:
			if !errors.Is(err, context.Canceled) {
				t.Errorf("feed(22 contacts, a stalled loader, canceled) = %v, want context.Canceled", err)
			}
		case <-g.wait():
			t.Fatalf("feed is still waiting after the run was canceled, want context.Canceled")
		}
		checkOrder(t, "feed(22 contacts, room for 4)", queued(q), varsk[:4])
	})
	t.Run("a slow loader takes one contact at a time", func(t *testing.T) {
		q := make(chan Contact)
		done := make(chan error, 1)
		go func() { done <- feed(context.Background(), varsk, q) }()
		var got []Contact
		for len(got) < len(varsk) {
			select {
			case c := <-q:
				got = append(got, c)
			case err := <-done:
				t.Fatalf("feed(22 contacts, a slow loader) = %v with %d of 22 taken, want it to wait for all 22", err, len(got))
			case <-g.wait():
				t.Fatalf("feed(22 contacts, a slow loader) handed over %d contacts, then nothing", len(got))
			}
		}
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("feed(22 contacts, a slow loader) = %v, want <nil>", err)
			}
		case <-g.wait():
			t.Fatalf("feed(22 contacts, a slow loader) is still waiting after the loader took all 22")
		}
		checkOrder(t, "feed(22 contacts, a slow loader)", got, varsk)
	})
}
