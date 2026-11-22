package main

import "testing"

func TestAppendPassID(t *testing.T) {
	tests := []struct {
		name         string
		buf          string
		sat, station string
		orbit        int
		want         string
	}{
		{"Furrow-B over Lantern Atoll on 7 February", "", "furrow-b", "LAN", 22379, "furrow-b/LAN/22379"},
		{"Pipit-1 over Port Ellery in December 2027, padded", "", "pipit-1", "PEL", 8154, "pipit-1/PEL/08154"},
		{"six digits, kept whole", "", "furrow-a", "LAN", 100000, "furrow-a/LAN/100000"},
		{"after what buf already holds", "frames/", "furrow-b", "LAN", 22379, "frames/furrow-b/LAN/22379"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(appendPassID([]byte(tt.buf), tt.sat, tt.station, tt.orbit))
			if got != tt.want {
				t.Errorf("appendPassID(%q, %q, %q, %d) = %q, want %q", tt.buf, tt.sat, tt.station, tt.orbit, got, tt.want)
			}
		})
	}
}

// night holds two passes of the night of 6–7 February, read the way uplinkd reads each frame's
// pass: from data, never from constants the compiler could fold away.
var night = []struct {
	sat, station string
	orbit        int
}{
	{"cirrine-1", "LAN", 17759},
	{"furrow-b", "LAN", 22379},
}

// TestNoAllocations builds a pass ID the way uplinkd does for every frame, in a buffer with room.
func TestNoAllocations(t *testing.T) {
	buf := make([]byte, 0, 64)
	frame := 0
	allocs := testing.AllocsPerRun(100, func() {
		p := night[frame%len(night)]
		frame++
		buf = appendPassID(buf[:0], p.sat, p.station, p.orbit)
	})
	if allocs != 0 {
		t.Errorf("testing.AllocsPerRun(100, appendPassID(buf with room, a pass of the night)) = %v, want 0", allocs)
	}
}
