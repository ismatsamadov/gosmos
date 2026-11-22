package main

import (
	"fmt"
	"net/http"
	"slices"
)

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
