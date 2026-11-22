package main

import "net/http"

// Band is a radio band.
type Band string

const (
	UHF   Band = "UHF"
	SBand Band = "S"
	XBand Band = "X"
)

// Station is one of Apogee's ground stations.
type Station struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	LatDeg     float64 `json:"lat_deg"`
	LonDeg     float64 `json:"lon_deg"`
	MaskDeg    float64 `json:"mask_deg"`
	Bands      []Band  `json:"bands"`
	Stowed     bool    `json:"stowed"`
	StowReason string  `json:"stow_reason,omitempty"`
}

// stationTable returns Apogee's five ground stations in the table's order, with the bands a customer can book.
func stationTable() []Station {
	return []Station{
		{Code: "YST", Name: "Ystmark", LatDeg: 77.8, LonDeg: 2.0, MaskDeg: 5, Bands: []Band{SBand, XBand}},
		{Code: "PEL", Name: "Port Ellery", LatDeg: 46.9, LonDeg: -38.1, MaskDeg: 5, Bands: []Band{UHF, SBand}},
		{Code: "SQT", Name: "Salar Quieto", LatDeg: -23.4, LonDeg: -104.6, MaskDeg: 3, Bands: []Band{SBand, XBand}},
		{Code: "LAN", Name: "Lantern Atoll", LatDeg: -1.8, LonDeg: 82.5, MaskDeg: 8, Bands: []Band{UHF, SBand}},
		{Code: "CWR", Name: "Cape Wren", LatDeg: -53.9, LonDeg: -95.0, MaskDeg: 5, Bands: []Band{SBand, XBand}},
	}
}

// listStations answers with every station the service knows, as JSON.
func (s *Service) listStations(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Stations)
}
