package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"
)

// FileStore keeps bookings in a MemStore and saves all of them to a file after every change.
type FileStore struct {
	mu   sync.Mutex // held through each change and its save, so saves happen in order
	mem  *MemStore
	path string
}

// OpenFileStore opens the bookings saved at path: no file is an empty store, a broken one an error.
func OpenFileStore(path string) (*FileStore, error) {
	f := &FileStore{mem: NewMemStore(), path: path}
	file, err := os.Open(path)
	if errors.Is(err, fs.ErrNotExist) {
		return f, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var saved []Booking
	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&saved); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if dec.Decode(&struct{}{}) != io.EOF {
		return nil, fmt.Errorf("%s: unexpected data after the bookings", path)
	}
	if err := f.mem.load(saved); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return f, nil
}

// Create keeps b and saves; when saving fails, it forgets b again and returns the error.
func (f *FileStore) Create(ctx context.Context, b Booking) (Booking, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	b, err := f.mem.Create(ctx, b)
	if err != nil {
		return Booking{}, err
	}
	if err := f.save(); err != nil {
		f.mem.Delete(context.WithoutCancel(ctx), b.ID)
		return Booking{}, err
	}
	return b, nil
}

// Get returns the booking with this ID, or ErrNotFound.
func (f *FileStore) Get(ctx context.Context, id string) (Booking, error) {
	return f.mem.Get(ctx, id)
}

// List returns every booking, ordered by AOS and then ID, in a new slice.
func (f *FileStore) List(ctx context.Context) ([]Booking, error) {
	return f.mem.List(ctx)
}

// Delete removes the booking and saves; when saving fails, it puts the booking back and returns the error.
func (f *FileStore) Delete(ctx context.Context, id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	b, err := f.mem.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := f.mem.Delete(ctx, id); err != nil {
		return err
	}
	if err := f.save(); err != nil {
		f.mem.load([]Booking{b})
		return err
	}
	return nil
}

// save writes every booking to the file as one JSON array, replacing the file whole.
func (f *FileStore) save() error {
	list, err := f.mem.List(context.Background())
	if err != nil {
		return err
	}
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return replaceFile(f.path, append(data, '\n'))
}
