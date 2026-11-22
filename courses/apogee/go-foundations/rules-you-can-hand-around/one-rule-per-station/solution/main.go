package main

// stationRules returns one minimum-elevation rule per station number, starting at 1.
func stationRules() (func(int) bool, func(int) bool, func(int) bool) {
	var first, second, third func(int) bool
	for minimum := 1; minimum <= 3; minimum++ {
		rule := func(elevationDeg int) bool { return elevationDeg >= minimum }
		switch minimum {
		case 1:
			first = rule
		case 2:
			second = rule
		case 3:
			third = rule
		}
	}
	return first, second, third
}
