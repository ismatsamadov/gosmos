package main

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"strings"
	"testing"
	"testing/iotest"
)

// february is the start of Furrow-A's February at Salar Quieto in the pass archive.
const february = `furrow-a/SQT/22290 2029-02-01T01:55:47Z 2029-02-01T02:04:05Z 15.6
furrow-a/SQT/22291 2029-02-01T03:29:11Z 2029-02-01T03:37:50Z 18.8
furrow-a/SQT/22298 2029-02-01T13:57:19Z 2029-02-01T14:07:05Z 39.9
furrow-a/SQT/22299 2029-02-01T15:33:19Z 2029-02-01T15:38:39Z 6.5
furrow-a/SQT/22305 2029-02-02T01:37:10Z 2029-02-02T01:43:06Z 7.5
`

// damaged returns february with line n, counted from 1, replaced by bad.
func damaged(n int, bad string) string {
	lines := strings.SplitAfter(february, "\n")
	lines[n-1] = bad + "\n"
	return strings.Join(lines, "")
}

// read is one step of a loop over ReadPasses: a pass ID, or an error in its place.
type read struct {
	id  string
	err error
}

// String shows r as a message does: the pass ID, or the error in quotes.
func (r read) String() string {
	if r.err != nil {
		return fmt.Sprintf("the error %#q", r.err.Error())
	}
	return r.id
}

// readUpTo ranges over seq until it ends or has given limit pairs, and returns them with
// whatever the loop panicked with. A pair with an error must carry the zero PassView.
func readUpTo(t *testing.T, seq iter.Seq2[PassView, error], limit int) (got []read, panicked any) {
	t.Helper()
	defer func() { panicked = recover() }()
	for p, err := range seq {
		if err != nil && p != (PassView{}) {
			t.Errorf("ReadPasses yielded %s with the error %#q, want the zero PassView", p.ID, err.Error())
		}
		got = append(got, read{p.ID, err})
		if len(got) == limit {
			break
		}
	}
	return got, nil
}

// wantPasses checks that got begins with the IDs of february's first n lines, each with a nil error.
func wantPasses(t *testing.T, got []read, n int) {
	t.Helper()
	lines := strings.Split(february, "\n")
	for i := range min(n, len(got)) {
		id := strings.Fields(lines[i])[0]
		if got[i].id != id || got[i].err != nil {
			t.Fatalf("ReadPasses pair %d = %s, want %s", i, got[i], id)
		}
	}
}

func TestReadPasses(t *testing.T) {
	t.Run("Furrow-A's first three passes of February", func(t *testing.T) {
		three := strings.Join(strings.SplitAfter(february, "\n")[:3], "")
		got, panicked := readUpTo(t, ReadPasses(strings.NewReader(three)), 4)
		if panicked != nil {
			t.Fatalf("ranging over ReadPasses panicked: %v", panicked)
		}
		wantPasses(t, got, 3)
		if len(got) != 3 {
			t.Fatalf("ReadPasses(3 good lines) yielded %d pairs, want 3", len(got))
		}
	})
	cases := []struct {
		name, archive, want string
		passes              int
	}{
		{"a third line cut short in a copy", damaged(3, "furrow-a/SQT/22298 2029-02-01T13:57:19Z 2029-02-01T14:07"),
			"archive line 3: 3 fields, want 4", 2},
		{"a flipped bit in the fifth line", damaged(5, "furrow-a/SQT/22305 2029-02-02T01:37:10Z 2029-02-02T01:43:06Z 7.%"),
			`archive line 5: bad max elevation "7.%"`, 4},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, panicked := readUpTo(t, ReadPasses(strings.NewReader(c.archive)), 6)
			if panicked != nil {
				t.Fatalf("ranging over ReadPasses panicked: %v", panicked)
			}
			wantPasses(t, got, c.passes)
			if len(got) <= c.passes {
				t.Fatalf("ReadPasses yielded %d passes and no error, want %d passes, then the error %#q", len(got), c.passes, c.want)
			}
			if r := got[c.passes]; r.err == nil || r.err.Error() != c.want {
				t.Fatalf("ReadPasses pair %d = %s, want the error %#q", c.passes, r, c.want)
			}
			if len(got) > c.passes+1 {
				t.Errorf("ReadPasses yielded %d pairs, want %d: nothing after the error", len(got), c.passes+1)
			}
		})
	}
	t.Run("Hollis stops at the first pass", func(t *testing.T) {
		got, panicked := readUpTo(t, ReadPasses(strings.NewReader(february)), 1)
		if panicked != nil {
			t.Fatalf("breaking after 1 pass panicked: %v", panicked)
		}
		wantPasses(t, got, 1)
	})
	t.Run("the disk fails after two lines", func(t *testing.T) {
		errDisk := errors.New("archive disk: input/output error")
		two := strings.Join(strings.SplitAfter(february, "\n")[:2], "")
		r := io.MultiReader(strings.NewReader(two), iotest.ErrReader(errDisk))
		got, panicked := readUpTo(t, ReadPasses(r), 4)
		if panicked != nil {
			t.Fatalf("ranging over ReadPasses panicked: %v", panicked)
		}
		wantPasses(t, got, 2)
		if len(got) != 3 || !errors.Is(got[2].err, errDisk) {
			t.Fatalf("ReadPasses(2 lines, then a failed read) yielded %d pairs, want 3, the last with the error %#q", len(got), errDisk.Error())
		}
	})
}
