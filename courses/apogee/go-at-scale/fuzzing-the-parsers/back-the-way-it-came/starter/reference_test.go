package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// errRefused is the reference parser's one error: it says only that a line is refused.
var errRefused = errors.New("line: refused by the reference parser")

// reference is Mirela's reference parser: it accepts exactly the lines formatLine writes for a
// frame with a lowercase satellite ID, readings in range and a known mode, and refuses the rest.
func reference(line string) (TelemetryFrame, error) {
	f := strings.Fields(line)
	if len(f) != 6 {
		return TelemetryFrame{}, errRefused
	}
	seq, errSeq := strconv.ParseUint(f[1], 10, 32)
	at, errAt := time.Parse("2006-01-02T15:04:05Z", f[2])
	batt, okBatt := strings.CutPrefix(f[3], "batt_v=")
	temp, okTemp := strings.CutPrefix(f[4], "temp_c=")
	mode, okMode := strings.CutPrefix(f[5], "mode=")
	battV, errBatt := strconv.ParseFloat(batt, 64)
	tempC, errTemp := strconv.ParseFloat(temp, 64)
	fr := TelemetryFrame{SatID: f[0], Seq: uint32(seq), At: at, BattV: battV, TempC: tempC, Mode: Mode(mode)}
	switch {
	case errSeq != nil, errAt != nil, errBatt != nil, errTemp != nil, !okBatt, !okTemp, !okMode,
		!lowercaseID(f[0]), !(battV >= 0 && battV <= 20), !(tempC >= -100 && tempC <= 150),
		mode != "NOMINAL" && mode != "SAFE" && mode != "DETUMBLE" && mode != "LEOP",
		formatLine(fr) != line:
		return TelemetryFrame{}, errRefused
	}
	return fr, nil
}

// lowercaseID reports whether id is 1 to 24 lowercase letters, digits and -, starting with a letter.
func lowercaseID(id string) bool {
	if len(id) == 0 || len(id) > 24 || id[0] < 'a' || id[0] > 'z' {
		return false
	}
	for _, r := range id {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}
	return true
}
