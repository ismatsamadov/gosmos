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
	"testing/iotest"
)

// errLinkDown is the error of a stream that breaks off.
var errLinkDown = errors.New("link down")

// pieces hands out what r holds 4 KiB at a time at most, as a network connection does.
type pieces struct{ r io.Reader }

// Read reads at most 4 KiB from r.
func (p pieces) Read(b []byte) (int, error) {
	if len(b) > 4096 {
		b = b[:4096]
	}
	return p.r.Read(b)
}

// image returns n bytes standing in for a Furrow-A image.
func image(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(i % 251)
	}
	return b
}

// short returns s, or its first 8 characters and an ellipsis when it is longer than 50.
func short(s string) string {
	if len(s) > 50 {
		return s[:8] + "…"
	}
	return s
}

func TestWholeStream(t *testing.T) {
	data := image(1 << 20)
	sum := sha256.Sum256(data)
	want := hex.EncodeToString(sum[:])
	got, err := fingerprintStream(pieces{bytes.NewReader(data)})
	if err != nil || got != want {
		t.Errorf("fingerprintStream(1 MiB in 4 KiB pieces) = %q, %v, want %q, nil", short(got), err, short(want))
	}
}

func TestEmptyStream(t *testing.T) {
	const want = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	got, err := fingerprintStream(strings.NewReader(""))
	if err != nil || got != want {
		t.Errorf("fingerprintStream(empty stream) = %q, %v, want %q, nil", short(got), err, short(want))
	}
}

func TestBrokenStream(t *testing.T) {
	r := pieces{io.MultiReader(bytes.NewReader(image(64<<10)), iotest.ErrReader(errLinkDown))}
	got, err := fingerprintStream(r)
	if got != "" || !errors.Is(err, errLinkDown) {
		t.Errorf("fingerprintStream(stream that breaks off after 64 KiB) = %q, %v, want %q, %v", short(got), err, "", errLinkDown)
	}
}

// TestNotHeldInMemory measures what fingerprintStream allocates for 16 MiB: a stream is hashed
// as it comes and never held whole, since a Furrow-A image is 10.9 GB.
func TestNotHeldInMemory(t *testing.T) {
	data := image(16 << 20)
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	_, err := fingerprintStream(pieces{bytes.NewReader(data)})
	runtime.ReadMemStats(&after)
	if kib := (after.TotalAlloc - before.TotalAlloc) >> 10; err != nil || kib > 1024 {
		t.Errorf("fingerprintStream(16 MiB) allocated %d KiB, %v, want under 1024 KiB, nil: hash the stream as it comes", kib, err)
	}
}
