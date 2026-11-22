package main

import (
	"fmt"
	"testing"
)

func TestPointBody(t *testing.T) {
	good := []struct {
		name         string
		azDeg, elDeg float64
		want         string
	}{
		{"Tomás's check on Dish 1", 183.25, 41.70, "47 95 10 4a"},
		{"north, on the horizon", 0, 0, "00 00 00 00"},
		{"the last hundredth before north, straight up", 359.99, 90, "8c 9f 23 28"},
		{"angles that need rounding", 1.15, 0.57, "00 73 00 39"},
	}
	for _, g := range good {
		t.Run(g.name, func(t *testing.T) {
			got, err := pointBody(g.azDeg, g.elDeg)
			if fmt.Sprintf("% x", got) != g.want || err != nil {
				t.Errorf("pointBody(%v, %v) = % x, %v; want %s, nil", g.azDeg, g.elDeg, got, err, g.want)
			}
		})
	}
	refused := []struct {
		name         string
		azDeg, elDeg float64
	}{
		{"azimuth 360", 360, 41.70},
		{"azimuth 359.996, which rounds to 360.00", 359.996, 41.70},
		{"azimuth below north", -0.01, 41.70},
		{"elevation past straight up", 183.25, 90.01},
		{"elevation below the horizon", 183.25, -0.01},
	}
	for _, r := range refused {
		t.Run(r.name, func(t *testing.T) {
			if got, err := pointBody(r.azDeg, r.elDeg); err == nil || got != nil {
				t.Errorf("pointBody(%v, %v) = % x, %v; want nil and an error", r.azDeg, r.elDeg, got, err)
			}
		})
	}
}
