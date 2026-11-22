package main

import (
	"bytes"
	"io/fs"
	"syscall"
	"testing"
)

// fullDisk is a file on a disk with room for only a few more bytes, where Lantern Atoll's was heading at 08:10.
type fullDisk struct {
	got  bytes.Buffer
	room int
}

// errNoSpace is what writing to a full disk returns.
var errNoSpace = &fs.PathError{Op: "write", Path: "uplinkd.log", Err: syscall.ENOSPC}

// Write takes what fits and fails with errNoSpace when p doesn't.
func (d *fullDisk) Write(p []byte) (int, error) {
	if len(p) <= d.room {
		d.room -= len(p)
		return d.got.Write(p)
	}
	n := d.room
	d.got.Write(p[:n])
	d.room = 0
	return n, errNoSpace
}

// line is one of uplinkd's records, 85 bytes with its newline.
const line = `{"time":"2028-12-04T08:10:00Z","level":"WARN","msg":"carrier lost","margin_db":-1.5}` + "\n"

func TestCountingWriterAddsUp(t *testing.T) {
	var buf bytes.Buffer
	c := &countingWriter{w: &buf}
	for i := 1; i <= 3; i++ {
		n, err := c.Write([]byte(line))
		if n != len(line) || err != nil {
			t.Fatalf("write %d: c.Write(line) = %d, %v, want %d, nil", i, n, err, len(line))
		}
		if want := int64(i * len(line)); c.n != want {
			t.Errorf("after write %d, %d bytes each, c.n = %d, want %d", i, len(line), c.n, want)
		}
	}
}

func TestCountingWriterStartsFromSize(t *testing.T) {
	var buf bytes.Buffer
	c := &countingWriter{w: &buf, n: 41} // a file that held 41 bytes when it was opened
	c.Write([]byte(line))
	if want := int64(41 + len(line)); c.n != want {
		t.Errorf("starting from 41 bytes, one write of %d made c.n = %d, want %d", len(line), c.n, want)
	}
}

func TestCountingWriterShortWrite(t *testing.T) {
	disk := &fullDisk{room: 10}
	c := &countingWriter{w: disk, n: 1000}
	n, err := c.Write([]byte(line))
	if n != 10 {
		t.Errorf("on a disk with room for 10 bytes, c.Write(line) wrote n = %d, want 10", n)
	}
	if err != errNoSpace {
		t.Errorf("on a full disk, c.Write(line) returned error %v, want the inner writer's own, %v", err, errNoSpace)
	}
	if c.n != 1010 {
		t.Errorf("after a write that took 10 of %d bytes, c.n = %d, want 1010: count what the file really got", len(line), c.n)
	}
}
