package main

import (
	"testing"
	"time"
)

// firstWithin calls first and stops the test if it hasn't returned within 1 s.
func firstWithin(t *testing.T, lan, cwr <-chan string) string {
	t.Helper()
	result := make(chan string, 1)
	go func() {
		result <- first(lan, cwr)
	}()
	select {
	case beacon := <-result:
		return beacon
	case <-time.After(time.Second):
		t.Fatal("first(lan, cwr) is still waiting after 1s, though a beacon is ready")
		return ""
	}
}

func TestFirstHeardAtCapeWren(t *testing.T) {
	t.Parallel()
	lan := make(chan string) // Lantern Atoll hears nothing
	cwr := make(chan string, 1)
	cwr <- "thistle 9003 CWR"
	if got := firstWithin(t, lan, cwr); got != "thistle 9003 CWR" {
		t.Errorf(`first(lan, cwr) = %q, want "thistle 9003 CWR"`, got)
	}
}

func TestFirstHeardAtLanternAtoll(t *testing.T) {
	t.Parallel()
	lan := make(chan string, 1)
	cwr := make(chan string) // Cape Wren hears nothing
	lan <- "thistle 9004 LAN"
	if got := firstWithin(t, lan, cwr); got != "thistle 9004 LAN" {
		t.Errorf(`first(lan, cwr) = %q, want "thistle 9004 LAN"`, got)
	}
}

func TestFirstWaitsForTheBeacon(t *testing.T) {
	t.Parallel()
	lan := make(chan string) // Thistle hasn't risen anywhere yet
	cwr := make(chan string)
	go func() {
		time.Sleep(20 * time.Millisecond) // it rises over Cape Wren a moment later
		cwr <- "thistle 9005 CWR"
	}()
	if got := firstWithin(t, lan, cwr); got != "thistle 9005 CWR" {
		t.Errorf(`first(lan, cwr) with the beacon still on its way = %q, want "thistle 9005 CWR"`, got)
	}
}
