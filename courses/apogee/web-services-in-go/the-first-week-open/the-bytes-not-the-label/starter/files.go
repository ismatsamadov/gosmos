package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"
)

// Files serves the files of customers' bookings: the command files they send up, and the files that land for them.
type Files struct {
	Store    Store            // the bookings, to tell whose each one is
	Commands string           // the folder of command files, one per booking: bk-4921.txt
	Landed   string           // dlinkd's delivery folder: each landed file in a folder named for its booking
	LinkKey  []byte           // signs download links: 32 random bytes of its own in production
	Now      func() time.Time // the service's clock
}

// booking returns the caller's booking whose ID the path names; another customer's is not found, exactly like a missing one.
func (f *Files) booking(r *http.Request) (Booking, error) {
	c, _ := customerFrom(r.Context())
	id := r.PathValue("id")
	b, err := f.Store.Get(r.Context(), id)
	if err == nil && !owns(c, b.SatID) {
		return Booking{}, fmt.Errorf("booking %s: %w", id, ErrNotFound)
	}
	return b, err
}

// upload keeps the command file a customer sends for one of its bookings, for the station to send up during the pass.
func (f *Files) upload(w http.ResponseWriter, r *http.Request) {
	b, err := f.booking(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	data, err := readCommands(w, r)
	if err != nil {
		writeError(w, r, uploadError(err))
		return
	}
	if err := checkCommands(data); err != nil {
		writeError(w, r, err)
		return
	}
	if err := replaceFile(filepath.Join(f.Commands, b.ID+".txt"), data); err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
