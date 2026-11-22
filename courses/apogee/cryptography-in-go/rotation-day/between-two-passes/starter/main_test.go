package main

import (
	"testing"
	"time"
)

// at returns the instant hh:mm:ss UTC on Rotation Day, 17 October 2028.
func at(hms string) time.Time {
	t, err := time.Parse(time.RFC3339, "2028-10-17T"+hms+"Z")
	if err != nil {
		panic(err)
	}
	return t
}

// rotationDay returns Rotation Day's ten passes, in the order they reached the log, not by AOS.
func rotationDay() []Pass {
	return []Pass{
		{"furrow-a", "SQT", at("14:02:10"), at("14:11:52")},
		{"pipit-1", "PEL", at("04:12:05"), at("04:22:41")},
		{"brine-3", "LAN", at("09:14:20"), at("09:22:31")},
		{"cirrine-1", "YST", at("07:30:00"), at("07:41:12")},
		{"pipit-3", "PEL", at("00:58:30"), at("01:07:10")},
		{"furrow-a", "CWR", at("04:58:30"), at("05:08:02")},
		{"cirrine-1", "YST", time.Date(2028, 10, 16, 23, 52, 10, 0, time.UTC), at("00:03:40")},
		{"brine-3", "LAN", at("01:02:00"), at("01:10:40")},
		{"pipit-2", "YST", at("04:31:10"), at("04:41:52")},
		{"furrow-b", "CWR", at("02:51:20"), at("03:00:44")},
	}
}

func TestRotationDay(t *testing.T) {
	passes := rotationDay()
	for _, tt := range []struct{ station, want string }{
		{"YST", "01:00:00"}, {"SQT", "01:00:00"}, {"CWR", "01:00:00"}, {"PEL", "01:12:10"}, {"LAN", "01:15:40"},
	} {
		t.Run(tt.station, func(t *testing.T) {
			got, err := switchTime(passes, tt.station, at("01:00:00"))
			if err != nil || !got.Equal(at(tt.want)) {
				t.Errorf("switchTime(Rotation Day's passes, %q, 01:00:00) = %s, %v, want %s, nil",
					tt.station, got.Format("15:04:05"), err, tt.want)
			}
		})
	}
}

func TestEdges(t *testing.T) {
	pass := func(aos, los string) Pass { return Pass{"furrow-b", "CWR", at(aos), at(los)} }
	tests := []struct {
		name   string
		passes []Pass
		want   string
	}{
		{"no pass at the station", nil, "01:00:00"},
		{"another station's pass at 01:00", []Pass{{"pipit-1", "PEL", at("01:00:00"), at("01:09:00")}}, "01:00:00"},
		{"a pass whose clear time starts as the switch ends", []Pass{pass("01:15:00", "01:24:00")}, "01:00:00"},
		{"the same pass, a second earlier", []Pass{pass("01:14:59", "01:23:59")}, "01:28:59"},
		{"a pass whose clear time ends at 01:00", []Pass{pass("00:45:00", "00:55:00")}, "01:00:00"},
		{"two passes with 9 clear minutes between them", []Pass{pass("01:29:00", "01:38:00"), pass("01:02:00", "01:10:00")}, "01:43:00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := switchTime(tt.passes, "CWR", at("01:00:00"))
			if err != nil || !got.Equal(at(tt.want)) {
				t.Errorf("switchTime(%s, %q, 01:00:00) = %s, %v, want %s, nil", tt.name, "CWR", got.Format("15:04:05"), err, tt.want)
			}
		})
	}
}

// TestNoGap books Cape Wren for a pass of 10 minutes every 20 minutes, from 00:40 until after 02:00 the next day.
func TestNoGap(t *testing.T) {
	var passes []Pass
	for aos := at("00:40:00"); aos.Before(at("02:01:00").Add(24 * time.Hour)); aos = aos.Add(20 * time.Minute) {
		passes = append(passes, Pass{"furrow-b", "CWR", aos, aos.Add(10 * time.Minute)})
	}
	got, err := switchTime(passes, "CWR", at("01:00:00"))
	if want := "no 10-minute gap for CWR before 2028-10-18T01:00:00Z"; err == nil || err.Error() != want {
		t.Errorf("switchTime(a pass every 20 minutes, %q, 01:00:00) = %s, %v, want an error %q", "CWR", got.Format(time.RFC3339), err, want)
	}
}
