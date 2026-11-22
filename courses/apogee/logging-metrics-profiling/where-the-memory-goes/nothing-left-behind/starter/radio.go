package main

import (
	"context"
	"errors"
	"time"
)

// ErrNoAck is the radio's error when the satellite never acknowledged the pass's last line.
var ErrNoAck = errors.New("no ack before LOS")

// Report is what the radio heard during a pass.
type Report struct {
	Frames []byte    // the pass's frames, in the buffer runPass gave fly
	AckAt  time.Time // the acknowledgment of the last line; zero without one
}

// result is how a pass's flight ended, as its goroutine hands it to runPass.
type result struct {
	rep Report
	err error
}

// Uplink is a station's uplinkd: four workers, and the radio that flies each pass.
type Uplink struct {
	Station string
	workers chan struct{} // room for 4: send to take a worker, receive to give it back
	// fly is the radio: it returns at the satellite's acknowledgment of the last line, or with
	// ErrNoAck as soon as ctx ends.
	fly func(ctx context.Context, p PassView, lines []string, buf []byte) (Report, error)
}

// newUplink returns station's uplinkd, its four workers free and fly as its radio.
func newUplink(station string, fly func(context.Context, PassView, []string, []byte) (Report, error)) *Uplink {
	return &Uplink{Station: station, workers: make(chan struct{}, 4), fly: fly}
}
