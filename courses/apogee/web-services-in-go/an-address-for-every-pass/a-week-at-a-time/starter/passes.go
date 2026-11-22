package main

import (
	"fmt"
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

// listPasses answers with the catalog's passes in the window the query asks for, of the satellite in ?sat= and the station in ?station= when they are given.
func (s *Service) listPasses(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sat, station := q.Get("sat"), q.Get("station")
	if sat != "" && !knownSat(sat) {
		http.Error(w, fmt.Sprintf("unknown satellite %q", sat), http.StatusBadRequest)
		return
	}
	if station != "" && !knownStation(station) {
		http.Error(w, fmt.Sprintf("unknown station %q", station), http.StatusBadRequest)
		return
	}
	from, days, err := window(q, s.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	to := from.Add(time.Duration(days) * 24 * time.Hour)
	passes := []PassView{}
	for _, p := range s.Catalog {
		inWindow := !p.AOS.Before(from) && p.AOS.Before(to)
		if inWindow && (sat == "" || p.SatID == sat) && (station == "" || p.Station == station) {
			passes = append(passes, p)
		}
	}
	writeJSON(w, http.StatusOK, passes)
}
