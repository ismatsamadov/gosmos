package main

import (
	"slices"
	"testing"
)

func TestEvaluate(t *testing.T) {
	rules := []Rule{BatteryLow{MinV: 6.6}, TooHot{MaxC: 45}}
	strict := []Rule{BatteryLow{MinV: 7.2}, TooHot{MaxC: 40}}
	tests := []struct {
		name  string
		rules []Rule
		frame TelemetryFrame
		want  []string
	}{
		{"thistle low battery", rules, TelemetryFrame{SatID: "thistle", BattV: 6.4, TempC: 20},
			[]string{"battery 6.40 V below 6.60 V"}},
		{"thistle too hot", rules, TelemetryFrame{SatID: "thistle", BattV: 7.9, TempC: 51},
			[]string{"temperature 51.0 C above 45.0 C"}},
		{"thistle both", rules, TelemetryFrame{SatID: "thistle", BattV: 6.4, TempC: 51},
			[]string{"battery 6.40 V below 6.60 V", "temperature 51.0 C above 45.0 C"}},
		{"pipit-1 healthy", rules, TelemetryFrame{SatID: "pipit-1", BattV: 7.9, TempC: 20}, nil},
		{"pipit-2 at the limits", rules, TelemetryFrame{SatID: "pipit-2", BattV: 6.6, TempC: 45}, nil},
		{"cirrine-1 under strict rules", strict, TelemetryFrame{SatID: "cirrine-1", BattV: 7.0, TempC: 42},
			[]string{"battery 7.00 V below 7.20 V", "temperature 42.0 C above 40.0 C"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := evaluate(tt.rules, tt.frame)
			if !slices.Equal(got, tt.want) {
				t.Errorf("evaluate(%+v, %+v) = %q, want %q", tt.rules, tt.frame, got, tt.want)
			}
		})
	}
}
