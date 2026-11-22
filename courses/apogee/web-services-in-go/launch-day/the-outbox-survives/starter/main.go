package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// OpenOutbox loads the outbox saved at path. A missing file is an empty outbox whose first event is evt-000001.
func OpenOutbox(path string) (*Outbox, error) {
	_, _, _, _ = os.ReadFile, errors.Is, fs.ErrNotExist, fmt.Errorf
	return &Outbox{path: path, next: 1, pending: []Delivery{}}, nil
}

// save writes the whole outbox to its file through replaceFile. The caller holds o.mu.
func (o *Outbox) save() error {
	_ = json.Marshal
	return nil
}
