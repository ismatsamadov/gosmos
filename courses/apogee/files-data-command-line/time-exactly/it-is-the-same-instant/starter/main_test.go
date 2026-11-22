package main

import "testing"

func TestNetworkTime(t *testing.T) {
	for _, tc := range []struct{ name, text, want string }{
		{"van clock", "2027-09-22T07:14:00-07:00", "2027-09-22T14:14:00Z"},
		{"network clock", "2027-09-22T14:14:00Z", "2027-09-22T14:14:00Z"},
		{"east of UTC", "2027-09-22T19:44:00+05:30", "2027-09-22T14:14:00Z"},
		{"evening in the van", "2027-09-22T20:30:00-07:00", "2027-09-23T03:30:00Z"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := networkTime(tc.text)
			if got != tc.want || err != nil {
				t.Errorf("networkTime(%q) = %q, %v, want %q, nil", tc.text, got, err, tc.want)
			}
		})
	}
	noOffset := "2027-09-22T07:14:00"
	got, err := networkTime(noOffset)
	if got != "" || err == nil {
		t.Errorf("networkTime(%q) = %q, %v, want \"\" and an error", noOffset, got, err)
	}
}
