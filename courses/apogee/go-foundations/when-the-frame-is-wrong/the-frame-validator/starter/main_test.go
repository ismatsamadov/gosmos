package main

import (
	"errors"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		frame     TelemetryFrame
		duplicate bool // errors.Is(err, ErrDuplicate) holds
		limit     bool // err holds a *LimitError for temp_c
		lines     int  // lines in the message; 0 means err is nil
	}{
		{"healthy", TelemetryFrame{SatID: "thistle", Seq: 4183, TempC: 21.5}, false, false, 0},
		{"duplicate", TelemetryFrame{SatID: "thistle", Seq: 4182, TempC: 21.5}, true, false, 1},
		{"900 C", TelemetryFrame{SatID: "thistle", Seq: 4183, TempC: 900}, false, true, 1},
		{"duplicate and 900 C", TelemetryFrame{SatID: "thistle", Seq: 4182, TempC: 900}, true, true, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.frame, 4182)
			if tt.lines == 0 {
				if err != nil {
					t.Errorf("validate(%+v, 4182) = %q, want nil", tt.frame, err)
				}
				return
			}
			if err == nil {
				t.Fatalf("validate(%+v, 4182) = nil, want an error", tt.frame)
			}
			if got := errors.Is(err, ErrDuplicate); got != tt.duplicate {
				t.Errorf("errors.Is(err, ErrDuplicate) = %v, want %v; err is %q", got, tt.duplicate, err)
			}
			le, ok := errors.AsType[*LimitError](err)
			if ok != tt.limit {
				t.Errorf("errors.AsType[*LimitError](err) found one = %v, want %v; err is %q", ok, tt.limit, err)
			} else if ok && (le.Field != "temp_c" || le.Value != tt.frame.TempC) {
				t.Errorf("LimitError = %+v, want Field %q and Value %v", *le, "temp_c", tt.frame.TempC)
			}
			if got := strings.Count(err.Error(), "\n") + 1; got != tt.lines {
				t.Errorf("message %q has %d lines, want %d", err, got, tt.lines)
			}
		})
	}
}
