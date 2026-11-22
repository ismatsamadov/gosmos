package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"maps"
	"strings"
	"testing"
)

const (
	// week1 and week2 are two weeks of the controller log, shortened.
	week1 = "# night shift 2027-07-07\n" +
		"2027-07-07T08:02:11Z\tLG\t96.40\t22.15\n" +
		"2027-07-07T08:02:12Z\tLG\t96.71\t22.48\n"
	week2 = "# night shift 2027-07-13\n" +
		"2027-07-13T07:41:52Z\tLG\t143.08\t35.62\n" +
		"2027-07-13T07:41:53Z\tLG\t143.30\t35.91\n"
	// manifest lists the two weeks, as the controller lists its log.
	manifest = `<?xml version="1.0" encoding="UTF-8"?>
<parcel version="1">
  <dish>LG</dish>
  <firmware>5.3</firmware>
  <made>2027-10-15T09:40:00Z</made>
  <files>
    <file size="97">log/2027-07-05.txt</file>
    <file size="99">log/2027-07-12.txt</file>
  </files>
</parcel>
`
)

// file is one file to pack into a test parcel.
type file struct {
	name string
	data []byte
}

// tarGz packs files, in order, into a gzip'd tar in memory.
func tarGz(t *testing.T, files ...file) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		t.Fatal(err)
	}
	tw := tar.NewWriter(zw)
	for _, f := range files {
		h := &tar.Header{Name: f.name, Mode: 0o600, Size: int64(len(f.data)), Typeflag: tar.TypeReg}
		if err := tw.WriteHeader(h); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(f.data); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// countingReader counts the bytes read through it.
type countingReader struct {
	r io.Reader
	n int
}

// Read reads from r and counts the bytes it got.
func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += n
	return n, err
}

func TestReadLogs(t *testing.T) {
	parcel := tarGz(t,
		file{"manifest.xml", []byte(manifest)},
		file{"log/2027-07-05.txt", []byte(week1)},
		file{"log/2027-07-12.txt", []byte(week2)})
	want := map[string]string{"manifest.xml": manifest, "log/2027-07-05.txt": week1, "log/2027-07-12.txt": week2}
	got, err := readLogs(bytes.NewReader(parcel), 64<<10)
	if err != nil || !maps.Equal(got, want) {
		t.Errorf("readLogs(the log parcel, 64 KiB) = %q, %v, want manifest.xml and both weeks, nil", got, err)
	}
	got, err = readLogs(bytes.NewReader(tarGz(t)), 64<<10)
	if err != nil || len(got) != 0 {
		t.Errorf("readLogs(a parcel with no files, 64 KiB) = %q, %v, want no files, nil", got, err)
	}
}

func TestReadLogsLimit(t *testing.T) {
	limit := int64(len(week1))
	got, err := readLogs(bytes.NewReader(tarGz(t, file{"log/2027-07-05.txt", []byte(week1)})), limit)
	if err != nil || got["log/2027-07-05.txt"] != week1 {
		t.Errorf("readLogs(a week of exactly %d bytes, %d) = %q, %v, want the week, nil", limit, limit, got, err)
	}
	longer := week1 + "\n"
	got, err = readLogs(bytes.NewReader(tarGz(t, file{"log/2027-07-05.txt", []byte(longer)})), limit)
	if err == nil || !strings.Contains(err.Error(), "log/2027-07-05.txt") {
		t.Errorf("readLogs(a week of %d bytes, %d) = %q, %v, want an error naming log/2027-07-05.txt", len(longer), limit, got, err)
	}
}

func TestReadLogsStopsAtTheLimit(t *testing.T) {
	parcel := tarGz(t, file{"log/2027-10-11.txt", make([]byte, 32<<20)}) // 32 MiB of zeros
	r := &countingReader{r: bytes.NewReader(parcel)}
	got, err := readLogs(r, 64<<10)
	if err == nil {
		t.Fatalf("readLogs(32 MiB of zeros, 64 KiB) = %d files, nil, want an error", len(got))
	}
	if r.n > len(parcel)/2 {
		t.Errorf("readLogs(32 MiB of zeros, 64 KiB) read %d of the parcel's %d bytes, want it to stop soon after 64 KiB", r.n, len(parcel))
	}
}

func TestReadLogsBrokenParcels(t *testing.T) {
	parcel := tarGz(t, file{"log/2027-07-05.txt", []byte(week1)}, file{"log/2027-07-12.txt", []byte(week2)})
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"a log that isn't gzip'd", []byte(week1)},
		{"a parcel cut off halfway", parcel[:len(parcel)/2]},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readLogs(bytes.NewReader(tc.data), 64<<10)
			if err == nil {
				t.Errorf("readLogs(%s, 64 KiB) = %q, nil, want an error", tc.name, got)
			}
		})
	}
}
