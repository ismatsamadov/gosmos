// Package telemetry parses housekeeping frames from telemetry lines.
package telemetry

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Frame is one decoded housekeeping frame.
type Frame struct {
	SatID string
	Seq   uint32
	BattV float64
}

// Parse reads a line such as "thistle 4182 2027-05-11T04:13:07Z batt_v=7.91 mode=NOMINAL".
func Parse(line string) (Frame, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return Frame{}, fmt.Errorf("line %q: want satellite, seq and time first", line)
	}
	seq, err := strconv.ParseUint(fields[1], 10, 32)
	if err != nil {
		return Frame{}, fmt.Errorf("seq: %w", err)
	}
	battV, err := battery(fields[3:])
	if err != nil {
		return Frame{}, err
	}
	return Frame{SatID: fields[0], Seq: uint32(seq), BattV: battV}, nil
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
