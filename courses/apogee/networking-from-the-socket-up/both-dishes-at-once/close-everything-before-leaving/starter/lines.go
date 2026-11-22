package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// greeting is the first line the watch port sends a watcher: the dishes dishd knows.
func greeting(names []string) string {
	return "dishd watching " + strings.Join(names, " ") + "\n"
}

// watchLine is one dish's line on the watch port: its state, azimuth, elevation and wind, or
// DOWN and why its poll failed, as describe says it for dishd's 5s timeout.
func watchLine(r Result) string {
	if r.Err != nil {
		why := strings.TrimPrefix(describe(r.Name, 5*time.Second, r.Err), r.Name+": ")
		return r.Name + " DOWN " + why
	}
	s := r.Status
	return fmt.Sprintf("%s %v az=%.2f el=%.2f wind=%d", r.Name, s.State, s.AzDeg, s.ElDeg, int(math.Round(s.WindKmh)))
}
