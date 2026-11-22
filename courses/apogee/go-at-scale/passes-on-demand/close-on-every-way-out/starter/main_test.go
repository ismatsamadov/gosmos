package main

import (
	"errors"
	"io"
	"iter"
	"strings"
	"testing"
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

// archiveFile is a file of the pass archive as the server opens it. It counts its closes.
type archiveFile struct {
	*strings.Reader
	closes int
}

// Close counts one close.
func (f *archiveFile) Close() error {
	f.closes++
	return nil
}

// opener returns an open function that opens text as an archive file, and that file.
func opener(text string) (func() (io.ReadCloser, error), *archiveFile) {
	f := &archiveFile{Reader: strings.NewReader(text)}
	return func() (io.ReadCloser, error) { return f, nil }, f
}

// loop ranges over seq, stops at the first error or once it has stop passes, and returns
// how many passes it saw, the error it stopped at and whatever it panicked with.
func loop(seq iter.Seq2[PassView, error], stop int) (passes int, err error, panicked any) {
	defer func() { panicked = recover() }()
	for _, e := range seq {
		if e != nil {
			return passes, e, nil
		}
		if passes++; passes == stop {
			break
		}
	}
	return passes, nil, nil
}

func TestArchive(t *testing.T) {
	cases := []struct {
		name, text string
		stop       int // the loop breaks after this many passes; 0 never
		passes     int
		err        string
	}{
		{"Hollis reads the whole file", february, 0, 5, ""},
		{"Hollis stops at the first pass", february, 1, 1, ""},
		{"a third line cut short in a copy", damaged(3, "furrow-a/SQT/22298 2029-02-01T13:57:19Z 2029-02-01T14:07"), 0, 2, "archive line 3: 3 fields, want 4"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			open, f := opener(c.text)
			passes, err, panicked := loop(Archive(open), c.stop)
			if panicked != nil {
				t.Fatalf("ranging over Archive panicked: %v", panicked)
			}
			if passes != c.passes || (err == nil) != (c.err == "") || err != nil && err.Error() != c.err {
				t.Fatalf("Archive gave the loop %d passes and the error %v, want %d passes and the error %q", passes, err, c.passes, c.err)
			}
			if f.closes != 1 {
				t.Errorf("after the loop, the archive file was closed %d times, want 1", f.closes)
			}
		})
	}
	t.Run("the archive server refuses to open it", func(t *testing.T) {
		errOpen := errors.New("open pass archive: too many open files")
		open := func() (io.ReadCloser, error) { return nil, errOpen }
		passes, err, panicked := loop(Archive(open), 0)
		if panicked != nil {
			t.Fatalf("ranging over Archive panicked: %v", panicked)
		}
		if passes != 0 || !errors.Is(err, errOpen) {
			t.Errorf("Archive gave the loop %d passes and the error %v, want 0 passes and the error %q", passes, err, errOpen)
		}
	})
}
