package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sync"
)

// Health is each satellite's newest telemetry frame: the one with the highest number,
// whichever station relayed it and in whatever order.
type Health struct {
	mu     sync.Mutex
	newest map[string]TelemetryFrame
}

// keep takes a frame, and keeps it if it is its satellite's newest.
func (h *Health) keep(f TelemetryFrame) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if old, ok := h.newest[f.SatID]; !ok || f.Seq > old.Seq {
		h.newest[f.SatID] = f
	}
}

// write saves reports/health.txt in dir: a line per satellite of the network, in fleet
// order, its newest frame as formatLine writes it or "<id> never heard".
func (h *Health) write(dir string, n *Network) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	var buf bytes.Buffer
	for _, s := range n.Satellites {
		if f, ok := h.newest[s.ID]; ok {
			fmt.Fprintln(&buf, formatLine(f))
		} else {
			fmt.Fprintln(&buf, s.ID+" never heard")
		}
	}
	return saveWhole(filepath.Join(dir, "reports", "health.txt"), buf.Bytes())
}
