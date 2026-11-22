package main

import (
	"bytes"
	"testing"
)

func TestFirstSeen(t *testing.T) {
	seen := make(map[[12]byte]bool)
	beacon := [12]byte{0x41, 0x50, 1, 0x10, 0x56}
	original := beacon
	if !firstSeen(beacon, seen) {
		t.Error("firstSeen(new) = false, want true")
	}
	if firstSeen(beacon, seen) {
		t.Error("firstSeen(repeated) = true, want false")
	}
	beacon[4]++
	if !seen[original] {
		t.Error("seen[original] = false, want true after changing caller array")
	}
	if !firstSeen(beacon, seen) {
		t.Error("firstSeen(changed) = false, want true")
	}
	if len(seen) != 2 {
		t.Errorf("len(seen) = %d, want 2", len(seen))
	}
}

func TestReadingBytes(t *testing.T) {
	beacon := [12]byte{0x41, 0x50, 1, 0x10, 0x56, 0, 0x1e, 0xe6, 0xff, 0xd3, 0, 0x82}
	got := readingBytes(&beacon)
	want := []byte{0x10, 0x56, 0, 0x1e, 0xe6}
	if !bytes.Equal(got, want) {
		t.Fatalf("readingBytes(beacon) = %v, want %v", got, want)
	}
	got[0] = 0x20
	if beacon[3] != 0x20 {
		t.Errorf("beacon[3] after editing view = %d, want 32", beacon[3])
	}
}
