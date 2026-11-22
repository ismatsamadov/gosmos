package main

import (
	"fmt"
	"net/http"
	"slices"
)

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

// listStations answers with every station, or with those that offer the band in ?band=.
func (s *Service) listStations(w http.ResponseWriter, r *http.Request) {
	band := Band(r.URL.Query().Get("band"))
	if band != "" && band != UHF && band != SBand && band != XBand {
		http.Error(w, fmt.Sprintf("unknown band %q", band), http.StatusBadRequest)
		return
	}
	list := []Station{}
	for _, st := range s.Stations {
		if band == "" || slices.Contains(st.Bands, band) {
			list = append(list, st)
		}
	}
	writeJSON(w, http.StatusOK, list)
}

// findStation returns the station whose code is exactly code, and whether there is one.
func findStation(stations []Station, code string) (Station, bool) {
	for _, st := range stations {
		if st.Code == code {
			return st, true
		}
	}
	return Station{}, false
}
