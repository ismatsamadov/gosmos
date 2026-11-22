package main

// Rule decides whether a pass at the given elevation is allowed.
type Rule func(float64) bool

// poll calls every rule in order, stopping at the first refusal.
func poll(elevationDeg float64, rules []Rule) bool {
	for _, rule := range rules {
		if !rule(elevationDeg) {
			return false
		}
	}
	return true
}
