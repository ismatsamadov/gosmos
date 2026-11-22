package main

// Station is one of Apogee's ground stations.
type Station struct {
	Code string // "YST": three capital letters
	Name string // "Ystmark"
}

// stations is Apogee's station table, in the order the team lists it.
var stations = []Station{
	{Code: "YST", Name: "Ystmark"},
	{Code: "PEL", Name: "Port Ellery"},
	{Code: "SQT", Name: "Salar Quieto"},
	{Code: "LAN", Name: "Lantern Atoll"},
	{Code: "CWR", Name: "Cape Wren"},
}

// findStation returns the station with the given code, and whether there is one.
func findStation(code string) (Station, bool) {
	for _, s := range stations {
		if s.Code == code {
			return s, true
		}
	}
	return Station{}, false
}
