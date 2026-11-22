package main

import (
	"math"
	"testing"
)

func TestClampVolts(t *testing.T) {
	if got := Clamp(Volts(9.1), 6.6, 8.4); math.Abs(float64(got)-8.4) > 0.01 {
		t.Errorf("Clamp(Volts(9.1), 6.6, 8.4) = %.2f, want 8.40", got)
	}
}

func TestClampKilometers(t *testing.T) {
	if got := Clamp(Kilometers(300), 400, 600); math.Abs(float64(got)-400) > 0.01 {
		t.Errorf("Clamp(Kilometers(300), 400, 600) = %.2f, want 400.00", got)
	}
}

func TestClampFloat64(t *testing.T) {
	if got := Clamp(7.5, 6.6, 8.4); math.Abs(got-7.5) > 0.01 {
		t.Errorf("Clamp(7.5, 6.6, 8.4) = %.2f, want 7.50", got)
	}
}

func TestClampKmPerHour(t *testing.T) {
	if got := Clamp(KmPerHour(131), 0, 120); math.Abs(float64(got)-120) > 0.01 {
		t.Errorf("Clamp(KmPerHour(131), 0, 120) = %.2f, want 120.00", got)
	}
}

// Celsius is declared here, not in main.go: Real has to admit it without naming it.
type Celsius float32

func TestClampCelsius(t *testing.T) {
	if got := Clamp(Celsius(-12.5), -10, 45); math.Abs(float64(got)+10) > 0.01 {
		t.Errorf("Clamp(Celsius(-12.5), -10, 45) = %.2f, want -10.00", got)
	}
}
