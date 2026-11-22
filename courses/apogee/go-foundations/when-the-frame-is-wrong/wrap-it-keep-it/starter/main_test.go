package main

import (
	"errors"
	"testing"
)

func TestIngest(t *testing.T) {
	tests := []struct {
		sat       string
		seq, last uint32
		want      string // "" means nil
	}{
		{"thistle", 4182, 4182, "thistle seq 4182: duplicate frame"},
		{"pipit-1", 900, 1200, "pipit-1 seq 900: duplicate frame"},
		{"thistle", 4183, 4182, ""},
	}
	for _, tt := range tests {
		err := ingest(tt.sat, tt.seq, tt.last)
		switch {
		case tt.want == "" && err != nil:
			t.Errorf("ingest(%q, %d, %d) = %q, want nil", tt.sat, tt.seq, tt.last, err)
		case tt.want != "" && err == nil:
			t.Errorf("ingest(%q, %d, %d) = nil, want %q", tt.sat, tt.seq, tt.last, tt.want)
		case tt.want != "" && err.Error() != tt.want:
			t.Errorf("ingest(%q, %d, %d) = %q, want %q", tt.sat, tt.seq, tt.last, err, tt.want)
		case tt.want != "" && !errors.Is(err, ErrDuplicate):
			t.Errorf("errors.Is(ingest(%q, %d, %d), ErrDuplicate) = false, want true", tt.sat, tt.seq, tt.last)
		}
	}
}
