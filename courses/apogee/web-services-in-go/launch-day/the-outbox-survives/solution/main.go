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
	o := &Outbox{path: path, next: 1, pending: []Delivery{}}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return o, nil
	}
	if err != nil {
		return nil, err
	}
	var saved outboxFile
	if err := json.Unmarshal(data, &saved); err != nil {
		return nil, fmt.Errorf("outbox %s: %w", path, err)
	}
	o.next = saved.NextEvent
	o.pending = append(o.pending, saved.Pending...)
	return o, nil
}

// save writes the whole outbox to its file through replaceFile. The caller holds o.mu.
func (o *Outbox) save() error {
	data, err := json.Marshal(outboxFile{NextEvent: o.next, Pending: o.pending})
	if err != nil {
		return err
	}
	return replaceFile(o.path, data)
}
