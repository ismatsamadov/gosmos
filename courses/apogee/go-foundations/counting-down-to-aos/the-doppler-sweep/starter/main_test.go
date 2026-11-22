package main

import "testing"

func TestBestOffset(t *testing.T) {
	tests := []struct {
		name                 string
		fromHz, toHz, stepHz int
		want                 int
	}{
		{"full sweep in 500 Hz steps", -12000, 12000, 500, 3500},
		{"low half only", -12000, 0, 500, 0},
		{"1 kHz steps, 3000 and 4000 tie", -12000, 12000, 1000, 3000},
		{"high band, first step strongest", 5000, 9000, 1000, 5000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bestOffset(tt.fromHz, tt.toHz, tt.stepHz)
			if got != tt.want {
				t.Errorf("bestOffset(%d, %d, %d) = %d, want %d", tt.fromHz, tt.toHz, tt.stepHz, got, tt.want)
			}
		})
	}
}
