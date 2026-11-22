package main

import "fmt"

// elevationCall returns GO if a pass's highest point is at or above the mask, else NO-GO, followed by the numbers behind the call.
func elevationCall(maxElevDeg, maskDeg float64) string {
	var call string
	// TODO: set call to GO or NO-GO
	return fmt.Sprintf("%s (max %.1f deg, mask %.1f deg)", call, maxElevDeg, maskDeg)
}

func main() {
	fmt.Println("Pipit-1 over Port Ellery:", elevationCall(42, 5))
}
