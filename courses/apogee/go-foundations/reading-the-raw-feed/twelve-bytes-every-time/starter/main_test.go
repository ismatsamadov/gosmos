package main

import "testing"

func TestSameBeacon(t *testing.T) {
	a := [12]byte{0x41, 0x50, 1, 0x10, 0x56, 0, 0x1e, 0xe6, 0xff, 0xd3, 0, 0x82}
	if !sameBeacon(a, a) {
		t.Error("sameBeacon(a, a) = false, want true")
	}
	for i := range a {
		b := a
		b[i] ^= 1
		if sameBeacon(a, b) {
			t.Errorf("sameBeacon(a, byte %d changed) = true, want false", i)
		}
	}
}

func TestWithSatellite(t *testing.T) {
	original := [12]byte{0x41, 0x50, 1, 0x10, 0x56, 0, 0x1e, 0xe6, 0xff, 0xd3, 0, 0x82}
	beacon := original
	want := original
	want[2] = 2
	if got := withSatellite(beacon, 2); got != want {
		t.Errorf("withSatellite(beacon, 2) = %v, want %v", got, want)
	}
	if beacon != original {
		t.Errorf("caller beacon = %v, want %v", beacon, original)
	}
}
