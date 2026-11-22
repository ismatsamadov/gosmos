package main

import (
	"testing"
	"time"
)

// frame returns one of the drill's frames from 7 March, as a station heard it.
func frame(sat string, seq uint32, hms string, battV, tempC float64, mode Mode) TelemetryFrame {
	at, _ := time.Parse(time.RFC3339, "2029-03-07T"+hms+"Z")
	return TelemetryFrame{SatID: sat, Seq: seq, At: at, BattV: battV, TempC: tempC, Mode: mode}
}

// The drill's frames: Cirrine-1 over Ystmark and Tirrin Rock, Pipit-1 over Port Ellery.
var (
	c9010 = frame("cirrine-1", 9010, "20:44:05", 8.02, 11.4, "NOMINAL")
	c9011 = frame("cirrine-1", 9011, "20:44:35", 8.02, 11.5, "NOMINAL")
	c9012 = frame("cirrine-1", 9012, "20:45:05", 8.02, 11.5, "NOMINAL")
	c9013 = frame("cirrine-1", 9013, "20:45:35", 8.02, 11.5, "NOMINAL")
	p6203 = frame("pipit-1", 6203, "05:25:12", 7.94, 3.5, "NOMINAL")
	p6204 = frame("pipit-1", 6204, "05:25:42", 7.93, 3.6, "SAFE")
)

func TestNewest(t *testing.T) {
	tests := []struct {
		name   string
		frames []TelemetryFrame
		want   map[string]uint32 // each satellite's newest Seq
	}{
		{"Cirrine-1's frames in order", []TelemetryFrame{c9010, c9011, c9012, c9013}, map[string]uint32{"cirrine-1": 9013}},
		{"9013 through Ystmark, then 9012 through Tirrin Rock", []TelemetryFrame{c9013, c9012}, map[string]uint32{"cirrine-1": 9013}},
		{"frame 9012 from both stations", []TelemetryFrame{c9012, c9012}, map[string]uint32{"cirrine-1": 9012}},
		{"Pipit-1 and Cirrine-1 kept apart", []TelemetryFrame{p6204, c9012, p6203}, map[string]uint32{"pipit-1": 6204, "cirrine-1": 9012}},
		{"Thistle-2's first frame, numbered 0, due on Thursday", []TelemetryFrame{{SatID: "thistle-2", Seq: 0, Mode: "DETUMBLE"}}, map[string]uint32{"thistle-2": 0}},
		{"a feed that closes with no frames", nil, map[string]uint32{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan TelemetryFrame, len(tt.frames))
			for _, f := range tt.frames {
				in <- f
			}
			close(in)
			got := newest(in)
			if got == nil {
				t.Fatalf("newest(%s) = nil, want an empty Health, ready to use", tt.name)
			}
			if len(got) != len(tt.want) {
				t.Errorf("len(newest(%s)) = %d, want %d", tt.name, len(got), len(tt.want))
			}
			for sat, seq := range tt.want {
				f, ok := got[sat]
				switch {
				case !ok:
					t.Errorf("newest(%s)[%q] is missing, want the frame with Seq %d", tt.name, sat, seq)
				case f.Seq != seq:
					t.Errorf("newest(%s)[%q].Seq = %d, want %d", tt.name, sat, f.Seq, seq)
				}
			}
		})
	}
}
