package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parsePass reads one line of the pass archive: a pass ID, its AOS
// and LOS in RFC 3339 UTC, and its max elevation in degrees.
func parsePass(line string) (PassView, error) {
	f := strings.Fields(line)
	if len(f) != 4 {
		return PassView{}, fmt.Errorf("%d fields, want 4", len(f))
	}
	id := strings.Split(f[0], "/")
	if len(id) != 3 {
		return PassView{}, fmt.Errorf("bad pass ID %q", f[0])
	}
	aos, err := time.Parse(time.RFC3339, f[1])
	if err != nil {
		return PassView{}, fmt.Errorf("bad AOS %q", f[1])
	}
	los, err := time.Parse(time.RFC3339, f[2])
	if err != nil {
		return PassView{}, fmt.Errorf("bad LOS %q", f[2])
	}
	elev, err := strconv.ParseFloat(f[3], 64)
	if err != nil {
		return PassView{}, fmt.Errorf("bad max elevation %q", f[3])
	}
	p := Pass{
		SatID: id[0], Station: id[1],
		AOS: aos, LOS: los, MaxElevDeg: elev,
	}
	return PassView{ID: f[0], Pass: p}, nil
}
