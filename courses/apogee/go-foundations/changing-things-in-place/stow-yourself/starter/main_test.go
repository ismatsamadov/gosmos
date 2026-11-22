package main

import "testing"

func TestStow(t *testing.T) {
	yst := Station{Code: "YST", Name: "Ystmark"}
	yst.Stow("wind 81 km/h")
	if !yst.Stowed {
		t.Errorf(`after yst.Stow("wind 81 km/h"), yst.Stowed = false, want true`)
	}
	if yst.StowReason != "wind 81 km/h" {
		t.Errorf(`after yst.Stow("wind 81 km/h"), yst.StowReason = %q, want "wind 81 km/h"`, yst.StowReason)
	}
}

func TestUnstow(t *testing.T) {
	yst := Station{Code: "YST", Name: "Ystmark", Stowed: true, StowReason: "wind 81 km/h"}
	yst.Unstow()
	if yst.Stowed {
		t.Errorf("after yst.Unstow(), yst.Stowed = true, want false")
	}
	if yst.StowReason != "" {
		t.Errorf(`after yst.Unstow(), yst.StowReason = %q, want ""`, yst.StowReason)
	}
}
