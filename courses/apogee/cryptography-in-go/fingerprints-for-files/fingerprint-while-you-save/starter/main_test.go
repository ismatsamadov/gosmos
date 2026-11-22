package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"runtime"
	"strings"
	"testing"
)

// errFull is the error of a delivery folder that has run out of room.
var errFull = errors.New("delivery folder full")

// countingReader counts the bytes read from r, to show the source is read once.
type countingReader struct {
	r io.Reader
	n int
}

// Read reads from r and counts what it got.
func (c *countingReader) Read(b []byte) (int, error) {
	n, err := c.r.Read(b)
	c.n += n
	return n, err
}

// fullAfter takes limit bytes, then fails every write with errFull.
type fullAfter struct {
	limit int
	saved int
}

// Write saves what still fits and fails with errFull when the rest doesn't.
func (f *fullAfter) Write(b []byte) (int, error) {
	if room := f.limit - f.saved; len(b) > room {
		f.saved += room
		return room, errFull
	}
	f.saved += len(b)
	return len(b), nil
}

// image returns n bytes standing in for a Furrow-A image.
func image(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(i % 251)
	}
	return b
}

// short shortens a digest longer than 50 characters to its hash's name and first 8 hex digits.
func short(s string) string {
	if len(s) <= 50 {
		return s
	}
	if name, digits, ok := strings.Cut(s, ":"); ok {
		return name + ":" + short(digits)
	}
	return s[:8] + "…"
}

func TestSave(t *testing.T) {
	data := image(256 << 10)
	sum := sha256.Sum256(data)
	want := "sha256:" + hex.EncodeToString(sum[:])
	src := &countingReader{r: bytes.NewReader(data)}
	var dst bytes.Buffer
	got, err := save(&dst, src)
	if err != nil || got != want {
		t.Errorf("save(256 KiB) = %q, %v, want %q, nil", short(got), err, short(want))
	}
	if !bytes.Equal(dst.Bytes(), data) {
		t.Errorf("save(256 KiB) saved %d bytes unlike the source, want its %d bytes exactly", dst.Len(), len(data))
	}
	if src.n != len(data) {
		t.Errorf("save(256 KiB) read %d bytes from the source, want %d: each byte once", src.n, len(data))
	}
}

func TestFullFolder(t *testing.T) {
	got, err := save(&fullAfter{limit: 1024}, bytes.NewReader(image(64<<10)))
	if got != "" || !errors.Is(err, errFull) {
		t.Errorf("save(to a folder full after 1 KiB) = %q, %v, want %q, %v", short(got), err, "", errFull)
	}
}

// TestNotHeldInMemory measures what save allocates for 16 MiB: the bytes pass through on their
// way to the folder and are never held whole, since a Furrow-A image is 10.9 GB.
func TestNotHeldInMemory(t *testing.T) {
	data := image(16 << 20)
	src := io.LimitReader(bytes.NewReader(data), int64(len(data)))
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	_, err := save(io.Discard, src)
	runtime.ReadMemStats(&after)
	if kib := (after.TotalAlloc - before.TotalAlloc) >> 10; err != nil || kib > 1024 {
		t.Errorf("save(16 MiB) allocated %d KiB, %v, want under 1024 KiB, nil: digest the bytes as they pass", kib, err)
	}
}
