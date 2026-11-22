package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Station is one of Apogee's ground stations, with the two fields the stand-in shows.
type Station struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// standIn answers like the booking service will, about stations only.
func standIn() http.Handler {
	stations := []Station{
		{Code: "YST", Name: "Ystmark"},
		{Code: "PEL", Name: "Port Ellery"},
		{Code: "SQT", Name: "Salar Quieto"},
		{Code: "LAN", Name: "Lantern Atoll"},
		{Code: "CWR", Name: "Cape Wren"},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /stations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stations)
	})
	mux.HandleFunc("GET /stations/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		for _, s := range stations {
			if s.Code == code {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(s)
				return
			}
		}
		http.Error(w, fmt.Sprintf("no station %q", code), http.StatusNotFound)
	})
	return mux
}
