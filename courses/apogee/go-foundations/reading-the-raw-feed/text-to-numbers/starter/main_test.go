package main

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"testing"
)

func TestBatteryReadsVoltage(t *testing.T) {
	line := "thistle 4182 2027-05-11T04:13:07Z batt_v=7.91 temp_c=-4.5 mode=NOMINAL"
	got, err := battery(line)
	if err != nil {
		t.Fatalf("battery(%q) error = %v, want nil", line, err)
	}
	if math.Abs(got-7.91) > 0.01 {
		t.Errorf("battery(%q) = %.2f, want 7.91", line, got)
	}
}

func TestBatteryRejectsBadNumber(t *testing.T) {
	line := "thistle 4183 2027-05-11T04:13:37Z batt_v=7.9l temp_c=-4.5 mode=NOMINAL"
	_, err := battery(line)
	if err == nil {
		t.Fatalf("battery(%q) error = nil, want an error that names batt_v", line)
	}
	if !strings.Contains(err.Error(), "batt_v") {
		t.Errorf("battery(%q) error = %q, want it to name batt_v", line, err)
	}
	if !errors.Is(err, strconv.ErrSyntax) {
		t.Errorf("errors.Is(err, strconv.ErrSyntax) = false, want true: wrap the parse error with %%w")
	}
}

func TestBatteryRejectsMissingField(t *testing.T) {
	line := "thistle 4374 2027-05-11T05:49:07Z temp_c=3.0 mode=NOMINAL"
	if _, err := battery(line); err == nil {
		t.Errorf("battery(%q) error = nil, want an error", line)
	}
}
