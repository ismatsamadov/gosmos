package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	furrowA = Booking{SatID: "furrow-a", Station: "SQT",
		AOS: time.Date(2028, 1, 3, 14, 11, 55, 0, time.UTC), LOS: time.Date(2028, 1, 3, 14, 21, 59, 0, time.UTC)}
	furrowB = Booking{SatID: "furrow-b", Station: "SQT",
		AOS: time.Date(2028, 1, 3, 14, 59, 32, 0, time.UTC), LOS: time.Date(2028, 1, 3, 15, 9, 5, 0, time.UTC)}
)

// open returns an empty FileStore that saves to bookings.json in a new folder, and that file's path.
func open(t *testing.T) (*FileStore, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "bookings.json")
	s, err := OpenFileStore(path)
	if err != nil {
		t.Fatalf("OpenFileStore(a new folder) error = %v, want nil", err)
	}
	return s, path
}

// readSaved returns the bookings saved at path, decoded strictly, and the file's raw bytes.
func readSaved(t *testing.T, path string) ([]Booking, []byte) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading bookings.json: %v, want the saved bookings", err)
	}
	var saved []Booking
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&saved); err != nil {
		t.Fatalf("bookings.json = %s, which doesn't decode as a list of bookings: %v", data, err)
	}
	return saved, data
}

func TestSaveAfterCreate(t *testing.T) {
	s, path := open(t)
	made, err := s.Create(context.Background(), furrowA)
	if err != nil {
		t.Fatalf("Create(furrow-a over SQT) error = %v, want nil", err)
	}
	saved, data := readSaved(t, path)
	if len(saved) != 1 || saved[0] != made {
		t.Errorf("after Create, bookings.json = %s, want booking %s alone", data, made.ID)
	}
}

func TestSaveAfterDelete(t *testing.T) {
	ctx := context.Background()
	s, path := open(t)
	made, err := s.Create(ctx, furrowA)
	if err != nil {
		t.Fatalf("Create(furrow-a over SQT) error = %v, want nil", err)
	}
	if err := s.Delete(ctx, made.ID); err != nil {
		t.Fatalf("Delete(%q) error = %v, want nil", made.ID, err)
	}
	if _, data := readSaved(t, path); string(bytes.TrimSpace(data)) != "[]" {
		t.Errorf("after the last booking's Delete, bookings.json = %s, want []", data)
	}
}

func TestSavedFileIsPrivate(t *testing.T) {
	s, path := open(t)
	if _, err := s.Create(context.Background(), furrowA); err != nil {
		t.Fatalf("Create(furrow-a over SQT) error = %v, want nil", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("after Create, bookings.json: %v, want a file", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("bookings.json mode = %v, want %v", got, os.FileMode(0o600))
	}
}

func TestFailedSaveChangesNothing(t *testing.T) {
	ctx := context.Background()
	s, path := open(t)
	kept, err := s.Create(ctx, furrowA)
	if err != nil {
		t.Fatalf("Create(furrow-a over SQT) error = %v, want nil", err)
	}
	dir := filepath.Dir(path)
	if err := os.Chmod(dir, 0o500); err != nil { // no new files in the folder: every save fails
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o700) })

	lost, err := s.Create(ctx, furrowB)
	if err == nil {
		t.Errorf("Create(furrow-b) with a read-only folder = %s, nil, want an error", lost.ID)
	}
	if _, err := s.Get(ctx, "bk-4502"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get(%q) after its failed save error = %v, want ErrNotFound: the store kept a booking its file never got", "bk-4502", err)
	}
	if err := s.Delete(ctx, kept.ID); err == nil {
		t.Errorf("Delete(%q) with a read-only folder = nil, want an error", kept.ID)
	}
	if _, err := s.Get(ctx, kept.ID); err != nil {
		t.Errorf("Get(%q) after its failed Delete error = %v, want nil: the file still holds it", kept.ID, err)
	}
}

func TestFileStoreKeepsTheContract(t *testing.T) {
	dir := t.TempDir()
	n := 0
	checkStore(t, func() Store {
		n++
		s, err := OpenFileStore(filepath.Join(dir, fmt.Sprintf("bookings-%d.json", n)))
		if err != nil {
			t.Fatalf("OpenFileStore(a new file) error = %v, want nil", err)
		}
		return s
	})
}
