package main

import (
	"fmt"
	"net/http"
)

// listStations answers with every station, or with those that offer the band in ?band=.
func (s *Service) listStations(w http.ResponseWriter, r *http.Request) {
	_ = fmt.Sprintf // TODO: read band from the query and answer only the stations that offer it
	writeJSON(w, http.StatusOK, s.Stations)
}
