package main

import (
	"os"
	"path/filepath"
	"time"
)

// Rotator writes a program's log file and rotates it.
type Rotator struct {
	Dir      string
	Program  string // "uplinkd": the current file is uplinkd.log
	MaxBytes int64
	Now      func() time.Time
	OnOpen   func() // runs each time the current file is opened: at start and after every rotation
	f        *os.File
	size     int64
	day      string // the current file's UTC day, "2028-12-04"
	n        int    // that day's rotated files so far
}

// open opens the current file and takes its size, its day (of its last change, or today's when it is empty) and that day's count, then runs OnOpen.
func (r *Rotator) open() error {
	f, err := openLog(filepath.Join(r.Dir, r.Program+".log"))
	if err != nil {
		return err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}
	files, err := rotatedFiles(r.Dir, r.Program)
	if err != nil {
		f.Close()
		return err
	}
	r.f, r.size, r.day, r.n = f, info.Size(), r.Now().UTC().Format(time.DateOnly), 0
	if r.size > 0 {
		r.day = info.ModTime().UTC().Format(time.DateOnly)
	}
	for _, rf := range files {
		if rf.Day == r.day {
			r.n = max(r.n, rf.N)
		}
	}
	if r.OnOpen != nil {
		r.OnOpen()
	}
	return nil
}

// Write writes p, one record, to the current file, rotating first when the UTC day has changed or p would take the file past MaxBytes.
func (r *Rotator) Write(p []byte) (int, error) {
	if r.f == nil {
		if err := r.open(); err != nil {
			return 0, err
		}
	}
	switch today := r.Now().UTC().Format(time.DateOnly); {
	case today > r.day && r.size == 0:
		r.day, r.n = today, 0
	case today > r.day, r.size > 0 && r.size+int64(len(p)) > r.MaxBytes:
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}
	n, err := r.f.Write(p)
	r.size += int64(n)
	return n, err
}

// rotate closes the current file, renames it under its day and the day's next number, and opens a fresh one.
func (r *Rotator) rotate() error {
	if err := r.f.Close(); err != nil {
		return err
	}
	old := filepath.Join(r.Dir, rotatedName(r.Program, r.day, r.n+1))
	if err := os.Rename(filepath.Join(r.Dir, r.Program+".log"), old); err != nil {
		return err
	}
	return r.open()
}

// Close closes the current file.
func (r *Rotator) Close() error {
	if r.f == nil {
		return nil
	}
	return r.f.Close()
}
