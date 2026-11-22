package main

import "testing"

func TestStowAction(t *testing.T) {
	yst := &Station{Code: "YST"}
	pel := &Station{Code: "PEL"}
	stowYst, stowPel := stowAction(yst), stowAction(pel)
	if stowYst == nil || stowPel == nil {
		t.Fatal("stowAction() = nil, want a function")
	}
	if yst.Stowed || pel.Stowed {
		t.Error("stowAction() stowed a station, want no change until called")
	}
	stowYst("wind 81 km/h")
	if !yst.Stowed || yst.StowReason != "wind 81 km/h" {
		t.Errorf("stowYst(reason) state = %v, want YST stowed for wind 81 km/h", yst)
	}
	if pel.Stowed {
		t.Error("stowYst(reason) changed PEL, want PEL unchanged")
	}
	stowPel("maintenance")
	if !pel.Stowed || pel.StowReason != "maintenance" {
		t.Errorf("stowPel(reason) state = %v, want PEL stowed for maintenance", pel)
	}
	stowYst("inspection")
	if yst.StowReason != "inspection" {
		t.Errorf("stowYst(inspection) reason = %q, want inspection", yst.StowReason)
	}
}
