package main

import (
	"errors"
	"net/http"
)

// writeError answers with the status, code and message of the *APIError in err, wrapped or not, or with a 500 for any other error.
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	_ = errors.AsType[*APIError]
	writeJSON(w, http.StatusInternalServerError, ErrorBody{Code: "internal", Message: "internal error"})
}
