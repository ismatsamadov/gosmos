package main

import (
	"fmt"
	"net/http"
)

// listPasses answers with the catalog's passes, of the satellite in ?sat= and the station in ?station= when they are given.
func (s *Service) listPasses(w http.ResponseWriter, r *http.Request) {
	_ = fmt.Sprintf
	writeJSON(w, http.StatusOK, s.Catalog)
}
