package main

import (
	"iter"
	"slices"
	"testing"
)

// counted yields vals one at a time, as a schedule does, and counts in *n the values it has
// yielded.
func counted[V any](vals []V, n *int) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range vals {
			*n++
			if !yield(v) {
				return
			}
		}
	}
}

// upTo ranges over seq until it ends or has given limit values, and returns them with
// whatever the loop panicked with.
func upTo[V any](seq iter.Seq[V], limit int) (got []V, panicked any) {
	defer func() { panicked = recover() }()
	for v := range seq {
		got = append(got, v)
		if len(got) == limit {
			break
		}
	}
	return got, nil
}

// ids returns the pass IDs of cs.
func ids(cs []Contact) []string {
	var out []string
	for _, c := range cs {
		out = append(out, c.ID)
	}
	return out
}

// passID is a contact's pass ID, for Map.
func passID(c Contact) string { return c.ID }

func TestTake(t *testing.T) {
	monday := firstWaveDay()
	first12 := ids(monday.Contacts[:12])

	t.Run("Ingrid's next twelve", func(t *testing.T) {
		asked := 0
		got, panicked := upTo(Take(counted(monday.Contacts, &asked), 12), 13)
		if panicked != nil {
			t.Fatalf("ranging over Take(Monday's contacts, 12) panicked: %v", panicked)
		}
		if len(got) > 12 {
			t.Fatalf("Take(Monday's contacts, 12) yielded more than 12 contacts, want 12")
		}
		if len(got) != 12 {
			t.Fatalf("Take(Monday's contacts, 12) yielded %d contacts, want 12", len(got))
		}
		for i := range got {
			if got[i].ID != first12[i] {
				t.Fatalf("Take(Monday's contacts, 12) contact %d = %s, want %s", i, got[i].ID, first12[i])
			}
		}
		if asked != 12 {
			t.Errorf("Take(Monday's contacts, 12) asked the schedule for %d contacts, want 12", asked)
		}
	})
	t.Run("a corner with no room", func(t *testing.T) {
		asked := 0
		got, panicked := upTo(Take(counted(monday.Contacts, &asked), 0), 1)
		if panicked != nil {
			t.Fatalf("ranging over Take(Monday's contacts, 0) panicked: %v", panicked)
		}
		if len(got) != 0 || asked != 0 {
			t.Errorf("Take(Monday's contacts, 0) yielded %d and asked the schedule for %d, want 0 and 0", len(got), asked)
		}
	})
	t.Run("a loop that stops at five", func(t *testing.T) {
		asked := 0
		got, panicked := upTo(Take(counted(monday.Contacts, &asked), 12), 5)
		if panicked != nil {
			t.Fatalf("breaking after 5 of Take(Monday's contacts, 12) panicked: %v", panicked)
		}
		if len(got) != 5 || asked != 5 {
			t.Errorf("Take(Monday's contacts, 12) gave %d contacts to a loop that stops at 5, after asking the schedule for %d, want 5 and 5", len(got), asked)
		}
	})
	t.Run("one Take, ranged twice", func(t *testing.T) {
		seq := Take(slices.Values(monday.Contacts), 12)
		first, _ := upTo(seq, 13)
		second, _ := upTo(seq, 13)
		if !slices.Equal(ids(first), first12) || !slices.Equal(ids(second), first12) {
			t.Errorf("ranging twice over one Take(Monday's contacts, 12) gave %d contacts, then %d, want the first 12 both times", len(first), len(second))
		}
	})
	t.Run("more than the day holds", func(t *testing.T) {
		got, panicked := upTo(Take(slices.Values(monday.Contacts), 300), 301)
		if panicked != nil {
			t.Fatalf("ranging over Take(Monday's contacts, 300) panicked: %v", panicked)
		}
		if len(got) != len(monday.Contacts) {
			t.Errorf("Take(Monday's %d contacts, 300) yielded %d contacts, want all %d", len(monday.Contacts), len(got), len(monday.Contacts))
		}
	})
}

func TestMap(t *testing.T) {
	monday := firstWaveDay()
	all := ids(monday.Contacts)

	t.Run("Monday's pass IDs", func(t *testing.T) {
		got, panicked := upTo(Map(slices.Values(monday.Contacts), passID), len(all)+1)
		if panicked != nil {
			t.Fatalf("ranging over Map(Monday's contacts, passID) panicked: %v", panicked)
		}
		if len(got) != len(all) {
			t.Fatalf("Map(Monday's %d contacts, passID) yielded %d IDs, want %d", len(all), len(got), len(all))
		}
		for i := range got {
			if got[i] != all[i] {
				t.Fatalf("Map(Monday's contacts, passID) ID %d = %q, want %q", i, got[i], all[i])
			}
		}
	})
	t.Run("a loop that stops at three", func(t *testing.T) {
		asked := 0
		got, panicked := upTo(Map(counted(monday.Contacts, &asked), passID), 3)
		if panicked != nil {
			t.Fatalf("breaking after 3 of Map(Monday's contacts, passID) panicked: %v", panicked)
		}
		if !slices.Equal(got, all[:3]) || asked != 3 {
			t.Errorf("Map(Monday's contacts, passID) gave %d IDs to a loop that stops at 3, after asking the schedule for %d contacts, want the first 3 after 3", len(got), asked)
		}
	})
	t.Run("the corner's twelve IDs", func(t *testing.T) {
		asked := 0
		got, panicked := upTo(Take(Map(counted(monday.Contacts, &asked), passID), 12), 13)
		if panicked != nil {
			t.Fatalf("ranging over Take(Map(Monday's contacts, passID), 12) panicked: %v", panicked)
		}
		if !slices.Equal(got, all[:12]) || asked != 12 {
			t.Errorf("Take(Map(Monday's contacts, passID), 12) yielded %d IDs after asking the schedule for %d contacts, want the first 12 after 12", len(got), asked)
		}
	})
}
