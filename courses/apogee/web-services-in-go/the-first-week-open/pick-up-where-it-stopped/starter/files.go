package main

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
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

// landed opens the file that landed for booking id, the one file in its folder, or answers 404 while there is none.
func (f *Files) landed(id string) (*os.File, error) {
	dir := filepath.Join(f.Landed, id)
	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) || err == nil && len(entries) != 1 {
		return nil, &APIError{Status: http.StatusNotFound, Code: "not_found", Message: "no file for " + id + " yet"}
	}
	if err != nil {
		return nil, err
	}
	return os.Open(filepath.Join(dir, entries[0].Name()))
}

// download sends the caller the file that landed for its booking, whole or from where a broken download stopped.
func (f *Files) download(w http.ResponseWriter, r *http.Request) {
	b, err := f.booking(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	file, err := f.landed(b.ID)
	if err != nil {
		writeError(w, r, err)
		return
	}
	defer file.Close()
	sendFile(w, r, file)
}
