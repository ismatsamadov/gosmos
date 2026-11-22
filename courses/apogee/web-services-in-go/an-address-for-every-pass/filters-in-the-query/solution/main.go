package main

import (
	"fmt"
	"net/http"
)

// listPasses answers with the catalog's passes, of the satellite in ?sat= and the station in ?station= when they are given.
func (s *Service) listPasses(w http.ResponseWriter, r *http.Request) {
	sat, station := r.URL.Query().Get("sat"), r.URL.Query().Get("station")
	if sat != "" && !knownSat(sat) {
		http.Error(w, fmt.Sprintf("unknown satellite %q", sat), http.StatusBadRequest)
		return
	}
	if station != "" && !knownStation(station) {
		http.Error(w, fmt.Sprintf("unknown station %q", station), http.StatusBadRequest)
		return
	}
	passes := []PassView{}
	for _, p := range s.Catalog {
		if (sat == "" || p.SatID == sat) && (station == "" || p.Station == station) {
			passes = append(passes, p)
		}
	}
	writeJSON(w, http.StatusOK, passes)
}
