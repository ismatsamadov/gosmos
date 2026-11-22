package main

import (
	"fmt"
	"time"
)

// positiveDelay reads a timeout such as 30s or 2m, refusing one that is not above zero.
func positiveDelay(text string) (time.Duration, error) {
	d, err := time.ParseDuration(text)
	if err != nil {
		return 0, err
	}
	if d <= 0 {
		return 0, fmt.Errorf("timeout %q is not above zero", text)
	}
	return d, nil
}
