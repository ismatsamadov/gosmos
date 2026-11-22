// Package telemetry parses housekeeping frames from telemetry lines.
package telemetry

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// frame is one decoded housekeeping frame.
type frame struct {
	satID string
	seq   uint32
	battV float64
}

// parse reads a line such as "thistle 4182 2027-05-11T04:13:07Z batt_v=7.91 mode=NOMINAL".
func parse(line string) (frame, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return frame{}, fmt.Errorf("line %q: want satellite, seq and time first", line)
	}
	seq, err := strconv.ParseUint(fields[1], 10, 32)
	if err != nil {
		return frame{}, fmt.Errorf("seq: %w", err)
	}
	battV, err := battery(fields[3:])
	if err != nil {
		return frame{}, err
	}
	return frame{satID: fields[0], seq: uint32(seq), battV: battV}, nil
}

// battery finds the batt_v field among key=value fields and parses its value.
func battery(pairs []string) (float64, error) {
	for _, p := range pairs {
		if key, value, ok := strings.Cut(p, "="); ok && key == "batt_v" {
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return 0, fmt.Errorf("batt_v: %w", err)
			}
			return v, nil
		}
	}
	return 0, errors.New("no batt_v field")
}
