package main

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestFixtureFS(t *testing.T) {
	if err := fstest.TestFS(fixtureFS(), "stations.csv", "templates/notices/stow.txt"); err != nil {
		t.Fatalf("fstest.TestFS(fixtureFS(), ...) = %v, want nil", err)
	}

	got, err := stowNotice(fixtureFS(), "CWR", "wind 81 km/h")
	if want := "CWR stowed: wind 81 km/h\n"; got != want || err != nil {
		t.Errorf("stowNotice(fixture, %q) = %q, %v, want %q, nil", "CWR", got, err, want)
	}

	if _, err := stowNotice(fixtureFS(), "SQT", "wind 40 km/h"); err == nil {
		t.Errorf("stowNotice(fixture, %q) error = nil, want an error: the fixture's table holds Cape Wren alone", "SQT")
	}

	noNotice := fixtureFS()
	delete(noNotice, "templates/notices/stow.txt")
	if _, err := stowNotice(noNotice, "CWR", "wind 81 km/h"); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("stowNotice(fixture without the notice, %q) error = %v, want a not-exist error", "CWR", err)
	}
}
