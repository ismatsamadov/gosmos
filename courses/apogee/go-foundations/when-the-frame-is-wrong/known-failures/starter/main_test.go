package main

import (
	"errors"
	"testing"
)

func TestAccept(t *testing.T) {
	if err := accept(4182, 4182); !errors.Is(err, ErrDuplicate) {
		t.Errorf("accept(4182, 4182) = %v, want ErrDuplicate", err)
	}
	if err := accept(4183, 4182); err != nil {
		t.Errorf("accept(4183, 4182) = %v, want nil", err)
	}
}

func TestShouldPage(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"ErrDuplicate", ErrDuplicate, false},
		{`errors.New("bad crc")`, errors.New("bad crc"), true},
		// Same text as ErrDuplicate, but a different error: only the sentinel itself is known.
		{`errors.New("duplicate frame")`, errors.New("duplicate frame"), true},
		{"nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldPage(tt.err); got != tt.want {
				t.Errorf("shouldPage(%s) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
