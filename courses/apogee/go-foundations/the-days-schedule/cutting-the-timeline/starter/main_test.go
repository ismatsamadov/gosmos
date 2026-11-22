package main

import (
	"slices"
	"testing"
)

func TestNextN(t *testing.T) {
	ids := []string{"pipit-1/PEL/04182", "furrow-a/SQT/11873", "thistle/LAN/02218", "cirrine-1/YST/07455", "pipit-3/YST/00214"}
	tests := []struct {
		name string
		n    int
		want []string
	}{
		{"next three", 3, []string{"pipit-1/PEL/04182", "furrow-a/SQT/11873", "thistle/LAN/02218"}},
		{"exactly five", 5, []string{"pipit-1/PEL/04182", "furrow-a/SQT/11873", "thistle/LAN/02218", "cirrine-1/YST/07455", "pipit-3/YST/00214"}},
		{"more than there are", 9, []string{"pipit-1/PEL/04182", "furrow-a/SQT/11873", "thistle/LAN/02218", "cirrine-1/YST/07455", "pipit-3/YST/00214"}},
		{"none", 0, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextN(ids, tt.n); !slices.Equal(got, tt.want) {
				t.Errorf("nextN(ids, %d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestLastN(t *testing.T) {
	readings := []float64{7.9, 7.8, 7.7}
	tests := []struct {
		name string
		n    int
		want []float64
	}{
		{"last two", 2, []float64{7.8, 7.7}},
		{"last one", 1, []float64{7.7}},
		{"none", 0, []float64{}},
		{"more than there are", 5, []float64{7.9, 7.8, 7.7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lastN(readings, tt.n); !slices.Equal(got, tt.want) {
				t.Errorf("lastN(%v, %d) = %v, want %v", readings, tt.n, got, tt.want)
			}
		})
	}
}
