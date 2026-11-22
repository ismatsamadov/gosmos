package main

import "testing"

func TestBudget(t *testing.T) {
	tests := []struct {
		name   string
		target Target
		booked int
		want   int
	}{
		{"Lantern Atoll's quarter", Atoll, 1199, 11},
		{"Ystmark's quarter", Standard, 2208, 11},
		{"Port Ellery's quarter", Standard, 1104, 5},
		{"Cape Wren's quarter", Standard, 920, 4},
		{"a hundred passes: not one to lose", Standard, 100, 0},
		{"ten thousand passes at 99.95%", 9995, 10000, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := budget(tt.target, tt.booked); got != tt.want {
				t.Errorf("budget(%d, %d) = %d, want %d", tt.target, tt.booked, got, tt.want)
			}
		})
	}
}
