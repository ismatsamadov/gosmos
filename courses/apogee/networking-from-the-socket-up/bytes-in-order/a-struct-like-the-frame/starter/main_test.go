package main

import (
	"math"
	"testing"
)

// Status bodies from the protocol's table.
const (
	monday      = "\x01\x02\x47\x95\x10\x4a\x02\x80\xff\x83\x00\x01\x00\x00\x00\x76"
	dish2Stowed = "\x02\x03\x47\x95\x23\x28\x03\x2a\xff\xe1\x00\x03\x00\x00\x00\x79"
)

func TestDecodeStatus(t *testing.T) {
	good := []struct {
		name, body string
		want       Status
	}{
		{"Dish 1 on Monday", monday, Status{Dish: 1, State: Tracking, AzDeg: 183.25, ElDeg: 41.70,
			WindKmh: 64.0, TempC: -12.5, Flags: FlagHeater, LastSeq: 118}},
		{"Dish 2 stowed by the wind", dish2Stowed, Status{Dish: 2, State: Stowed, AzDeg: 183.25, ElDeg: 90.00,
			WindKmh: 81.0, TempC: -3.1, Flags: FlagHeater | FlagWindLock, LastSeq: 121}},
	}
	for _, g := range good {
		t.Run(g.name, func(t *testing.T) {
			got, err := decodeStatus([]byte(g.body))
			if err != nil || !sameStatus(got, g.want) {
				t.Errorf("decodeStatus(% x) = %+v, %v; want %+v, nil", g.body, got, err, g.want)
			}
		})
	}
	wrong := []struct {
		name, body, want string
	}{
		{"one byte short", monday[:15], "status body is 15 bytes, want 16"},
		{"one byte too many", monday + "\x00", "status body is 17 bytes, want 16"},
	}
	for _, w := range wrong {
		t.Run(w.name, func(t *testing.T) {
			got, err := decodeStatus([]byte(w.body))
			if err == nil || err.Error() != w.want || got != (Status{}) {
				t.Errorf("decodeStatus(%d bytes) = %+v, %v; want an empty Status and %q", len(w.body), got, err, w.want)
			}
		})
	}
}

// sameStatus reports whether two statuses agree, their measurements to a thousandth.
func sameStatus(a, b Status) bool {
	near := func(x, y float64) bool { return math.Abs(x-y) < 0.001 }
	return a.Dish == b.Dish && a.State == b.State && a.Flags == b.Flags && a.LastSeq == b.LastSeq &&
		near(a.AzDeg, b.AzDeg) && near(a.ElDeg, b.ElDeg) && near(a.WindKmh, b.WindKmh) && near(a.TempC, b.TempC)
}
