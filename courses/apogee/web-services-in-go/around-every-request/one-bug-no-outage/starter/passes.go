package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
	// Monday's change keeps the orbit number's first five digits, to ignore anything typed
	// after them. With four digits it panics: Nkem keeps the bug until Friday's test.
	orbit := r.PathValue("orbit")[:5]
	id := r.PathValue("sat") + "/" + r.PathValue("station") + "/" + orbit
	for _, p := range s.Catalog {
		if p.ID == id {
			writeJSON(w, http.StatusOK, p)
			return
		}
	}
	writeError(w, r, &APIError{Status: http.StatusNotFound, Code: "not_found", Message: "no pass " + id})
}

// listPasses answers with the catalog's passes in the window the query asks for, of the satellite in ?sat= and the station in ?station= when they are given.
func (s *Service) listPasses(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sat, station := q.Get("sat"), q.Get("station")
	if sat != "" && !knownSat(sat) {
		writeError(w, r, &APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: fmt.Sprintf("unknown satellite %q", sat)})
		return
	}
	if station != "" && !knownStation(station) {
		writeError(w, r, &APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: fmt.Sprintf("unknown station %q", station)})
		return
	}
	from, days, err := window(q, s.Now())
	if err != nil {
		writeError(w, r, &APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: err.Error()})
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

// window reads the start and the number of days of passes a query asks for, with now and 7 as defaults.
func window(q url.Values, now time.Time) (from time.Time, days int, err error) {
	from, days = now, 7
	if v := q.Get("from"); v != "" {
		if from, err = time.Parse(time.RFC3339, v); err != nil {
			return time.Time{}, 0, errors.New("from must be an RFC 3339 time")
		}
	}
	if v := q.Get("days"); v != "" {
		if days, err = strconv.Atoi(v); err != nil || days < 1 || days > 14 {
			return time.Time{}, 0, errors.New("days must be a whole number from 1 to 14")
		}
	}
	return from, days, nil
}
