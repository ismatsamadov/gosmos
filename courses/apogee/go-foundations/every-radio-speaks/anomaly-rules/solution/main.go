package main

import "fmt"

// TelemetryFrame is one decoded housekeeping frame.
type TelemetryFrame struct {
	SatID string
	BattV float64
	TempC float64
}

// Rule is an anomaly rule the pager runs on every frame.
type Rule interface {
	Check(f TelemetryFrame) (string, bool)
}

// BatteryLow fires when the battery is below MinV.
type BatteryLow struct {
	MinV float64
}

// Check returns a message and true when f's battery is below MinV.
func (r BatteryLow) Check(f TelemetryFrame) (string, bool) {
	if f.BattV < r.MinV {
		return fmt.Sprintf("battery %.2f V below %.2f V", f.BattV, r.MinV), true
	}
	return "", false
}

// TooHot fires when the temperature is above MaxC.
type TooHot struct {
	MaxC float64
}

// Check returns a message and true when f is hotter than MaxC.
func (r TooHot) Check(f TelemetryFrame) (string, bool) {
	if f.TempC > r.MaxC {
		return fmt.Sprintf("temperature %.1f C above %.1f C", f.TempC, r.MaxC), true
	}
	return "", false
}

// evaluate runs every rule on f and returns the messages of those that fire, in rule order.
func evaluate(rules []Rule, f TelemetryFrame) []string {
	var msgs []string
	for _, r := range rules {
		if msg, fired := r.Check(f); fired {
			msgs = append(msgs, msg)
		}
	}
	return msgs
}

func main() {
	rules := []Rule{BatteryLow{MinV: 6.6}, TooHot{MaxC: 45}}
	frame := TelemetryFrame{SatID: "thistle", BattV: 6.4, TempC: 51}
	for _, msg := range evaluate(rules, frame) {
		fmt.Println(msg)
	}
}
