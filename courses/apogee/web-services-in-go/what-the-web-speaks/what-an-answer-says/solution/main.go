package main

import (
	"net/http"
	"net/http/httptest"
)

func main() {
	req := httptest.NewRequest(http.MethodGet, "/stations/PEL", nil)
	rec := httptest.NewRecorder()
	standIn().ServeHTTP(rec, req)
	printResponse(rec)
}
