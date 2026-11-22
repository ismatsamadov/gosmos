package main

import "testing"

func TestStationCode(t *testing.T) {
	got := stationCode()
	want := "PEL"
	if got != want {
		t.Errorf("stationCode() = %q, want %q", got, want)
	}
}
