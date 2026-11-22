package main

import (
	"net/http"
	"net/http/httptest"
)

func main() {
	req := httptest.NewRequest(http.MethodGet, "/stations/PEL", nil)
	rec := httptest.NewRecorder()
	_ = req // TODO: let the stand-in answer req, writing to rec
	printResponse(rec)
}
