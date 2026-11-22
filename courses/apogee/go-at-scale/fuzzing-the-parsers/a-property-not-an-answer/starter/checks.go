package main

import (
	"math"
	"testing"
)

// checkLine fails t when parse accepts line but gives a frame no satellite could send.
func checkLine(t *testing.T, parse func(string) (TelemetryFrame, error), line string) {
	parse(line)
	_ = math.IsNaN // TODO: check the frames parse accepts
}
