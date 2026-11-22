package main

import "testing"

// network returns Apogee's five stations, none of them stowed.
func network() []Station {
	return []Station{
		{Code: "YST", Name: "Ystmark"},
		{Code: "PEL", Name: "Port Ellery"},
		{Code: "SQT", Name: "Salar Quieto"},
		{Code: "LAN", Name: "Lantern Atoll"},
		{Code: "CWR", Name: "Cape Wren"},
	}
}

func TestFind(t *testing.T) {
	cases := []struct{ code, want string }{
		{"YST", "Ystmark"},
		{"SQT", "Salar Quieto"},
		{"CWR", "Cape Wren"},
	}
	for _, c := range cases {
		s := find(network(), c.code)
		if s == nil {
			t.Errorf("find(stations, %q) = nil, want a pointer to %s", c.code, c.want)
			continue
		}
		if s.Name != c.want {
			t.Errorf("find(stations, %q).Name = %q, want %q", c.code, s.Name, c.want)
		}
	}
}

func TestFindPointsIntoSlice(t *testing.T) {
	stations := network()
	s := find(stations, "SQT")
	if s == nil {
		t.Fatalf(`find(stations, "SQT") = nil, want a pointer to Salar Quieto`)
	}
	s.Stowed = true
	if !stations[2].Stowed {
		t.Errorf(`after find(stations, "SQT").Stowed = true, stations[2].Stowed = false, want true: the pointer must point into the slice, not at a copy`)
	}
}

func TestFindNobodyHome(t *testing.T) {
	if s := find(network(), "XYZ"); s != nil {
		t.Errorf(`find(stations, "XYZ") = &%+v, want nil`, *s)
	}
}
