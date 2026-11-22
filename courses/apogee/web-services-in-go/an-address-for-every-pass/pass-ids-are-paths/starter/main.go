package main

import "net/http"

// routes registers the service's handlers on a new mux.
func (s *Service) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /stations", s.listStations)
	mux.HandleFunc("GET /stations/{code}", s.getStation)
	mux.HandleFunc("GET /passes/{id}", s.getPass)
	mux.Handle("GET /guide/", http.StripPrefix("/guide/", http.FileServerFS(s.Guide)))
	return mux
}

// getPass answers with the catalog's pass whose ID is the path's satellite, station and orbit, or 404.
func (s *Service) getPass(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	for _, p := range s.Catalog {
		if p.ID == id {
			writeJSON(w, http.StatusOK, p)
			return
		}
	}
	http.Error(w, "no pass "+id, http.StatusNotFound)
}
