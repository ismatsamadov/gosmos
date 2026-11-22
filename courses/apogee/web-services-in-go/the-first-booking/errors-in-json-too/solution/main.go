package main

import (
	"errors"
	"net/http"
)

// writeError answers with the status, code and message of the *APIError in err, wrapped or not, or with a 500 for any other error.
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	if e, ok := errors.AsType[*APIError](err); ok {
		writeJSON(w, e.Status, ErrorBody{Code: e.Code, Message: e.Message})
		return
	}
	writeJSON(w, http.StatusInternalServerError, ErrorBody{Code: "internal", Message: "internal error"})
}
