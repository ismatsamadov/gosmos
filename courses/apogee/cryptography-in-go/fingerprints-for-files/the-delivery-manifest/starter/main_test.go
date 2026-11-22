package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

// occultation returns n bytes standing in for one of Cirrine-1's measurement files.
func occultation(n int, seed byte) []byte {
	return bytes.Repeat([]byte{seed, 'c', '1', seed + 1}, n/4)
}

// line returns a manifest line for data under name, hashed by alg.
func line(alg string, data []byte, name string) string {
	var sum []byte
	if alg == "sha3-256" {
		s := sha3.Sum256(data)
		sum = s[:]
	} else {
		s := sha256.Sum256(data)
		sum = s[:]
	}
	return fmt.Sprintf("%s:%s %d %s", alg, hex.EncodeToString(sum), len(data), name)
}

func TestCheckManifest(t *testing.T) {
	a, b := occultation(4096, 1), occultation(2048, 2)
	c, d, e := occultation(1024, 3), occultation(512, 4), occultation(256, 5)
	changed := bytes.Clone(c)
	changed[700] ^= 0x01
	fsys := fstest.MapFS{
		"occ-0721-0730.bin": {Data: a},
		"occ-0721-0907.bin": {Data: b},
		"occ-0721-1043.bin": {Data: changed},
		"occ-0721-1220.bin": {Data: d[:500]},
	}
	manifest := strings.Join([]string{
		line("sha256", a, "occ-0721-0730.bin"),
		line("sha3-256", b, "occ-0721-0907.bin"),
		line("sha256", c, "occ-0721-1043.bin"),
		line("sha256", d, "occ-0721-1220.bin"),
		line("sha3-256", e, "occ-0721-1356.bin"),
	}, "\n") + "\n"
	want := []string{
		"ok occ-0721-0730.bin",
		"ok occ-0721-0907.bin",
		"changed occ-0721-1043.bin",
		"wrong size occ-0721-1220.bin",
		"missing occ-0721-1356.bin",
	}
	got, err := checkManifest(fsys, strings.NewReader(manifest))
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("checkManifest(Friday's delivery) = %q, %v, want %q, nil", got, err, want)
	}
}

func TestBadLine(t *testing.T) {
	a := occultation(4096, 1)
	fsys := fstest.MapFS{"occ-0721-0730.bin": {Data: a}}
	good := line("sha256", a, "occ-0721-0730.bin")
	f := strings.Fields(good)
	two := f[0] + " " + f[2] // the size left out
	manifest := good + "\n" + good + "\n" + two + "\n" + good + "\n"
	const want = "manifest line 3: want <digest> <bytes> <file>"
	got, err := checkManifest(fsys, strings.NewReader(manifest))
	if err == nil || err.Error() != want {
		t.Errorf("checkManifest(line 3 with two fields) = %q, %v, want nil, %#q", got, err, want)
	}
}

func TestEmptyManifest(t *testing.T) {
	got, err := checkManifest(fstest.MapFS{}, strings.NewReader(""))
	if len(got) != 0 || err != nil {
		t.Errorf("checkManifest(empty manifest) = %q, %v, want no lines and nil", got, err)
	}
}
