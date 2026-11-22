package main

import (
	"net/http"
	"time"
)

// Pass is one contact between a satellite and a station.
type Pass struct {
	SatID      string    `json:"sat_id"`
	Station    string    `json:"station"`
	AOS        time.Time `json:"aos"`
	LOS        time.Time `json:"los"`
	MaxElevDeg float64   `json:"max_elev_deg"`
}

// PassView is a predicted pass as the service shows it. Pass is embedded, so its fields sit beside "id" in the JSON.
type PassView struct {
	ID string `json:"id"` // the pass ID: "furrow-a/SQT/15864"
	Pass
}

// getPass answers with the catalog's pass whose ID is the path's satellite, station and orbit, or 404.
func (s *Service) getPass(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("sat") + "/" + r.PathValue("station") + "/" + r.PathValue("orbit")
	for _, p := range s.Catalog {
		if p.ID == id {
			writeJSON(w, http.StatusOK, p)
			return
		}
	}
	http.Error(w, "no pass "+id, http.StatusNotFound)
}
