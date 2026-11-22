package main

import (
	"slices"
	"testing"
)

func TestNumber(t *testing.T) {
	jog := Frame{Type: TypeJog, Body: []byte{'A', 0x00, 0x23}} // A +0.35°

	// After a WELCOME with last_seq 118, the next command is 119.
	l := &Link{next: 119}
	for _, want := range []uint32{119, 120, 121} {
		got := l.number(jog)
		if got.Seq != want {
			t.Errorf("number(JOG) = seq %d, want %d", got.Seq, want)
		}
		if got.Type != TypeJog || !slices.Equal(got.Body, jog.Body) {
			t.Errorf("number(JOG) = type 0x%02x, body % x, want the JOG unchanged but for its seq", got.Type, got.Body)
		}
	}
	if l.next != 122 {
		t.Errorf("after three commands from 119, l.next = %d, want 122", l.next)
	}
}

func TestLearn(t *testing.T) {
	tests := []struct {
		name    string
		next    uint32
		lastSeq uint32
		want    uint32
	}{
		{"a fresh link meets a new controller", 0, 0, 1},
		{"a fresh link meets YST-2 at 118", 0, 118, 119},
		{"the controller ran our last command", 122, 121, 122},
		{"the controller is ahead", 122, 200, 201},
		{"the controller is behind", 122, 57, 122},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Link{next: tt.next}
			l.learn(Welcome{Dish: 2, FwMajor: 2, FwMinor: 4, LastSeq: tt.lastSeq})
			if l.next != tt.want {
				why := "continue one past the controller's last_seq"
				if l.next < tt.next {
					why = "numbers never go back"
				}
				t.Errorf("learn(last_seq %d) with l.next at %d sets l.next to %d, want %d: %s", tt.lastSeq, tt.next, l.next, tt.want, why)
			}
		})
	}
}
