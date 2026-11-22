package main

// lostFrame is the reading batteryAt returns for a minute whose frame was lost.
const lostFrame = -1.0

// batteryAt returns Pipit-3's battery voltage at a minute of its first pass, or lostFrame.
func batteryAt(minute int) float64 {
	switch minute {
	case 0:
		return 7.92
	case 1:
		return 7.90
	case 3:
		return 7.86
	case 5:
		return 7.80
	default:
		return lostFrame
	}
}
