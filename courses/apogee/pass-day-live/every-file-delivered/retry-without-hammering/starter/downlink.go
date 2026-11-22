package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Downlink is dlinkd's reassembly: the files satellites send down on booked passes, kept
// by chunk number until each is whole.
type Downlink struct {
	mu    sync.Mutex
	files map[string]*downFile // by name
	made  sync.Once            // files/ is made at the first landing
	dir   error                // and this is why it could not be
}

// downFile is one file coming down.
type downFile struct {
	name    string  // "thistle-2-20290726-1603.dat"
	booking Booking // its pass's booking; no ID when it had none
	total   int     // its chunk count
	chunks  map[int][]byte
	landed  bool
	crcBad  bool // every chunk came and the CRC didn't match
}

// newDownlink returns an empty reassembly.
func newDownlink() *Downlink {
	return &Downlink{files: map[string]*downFile{}}
}

// frame takes a FILE, CHUNK or END frame from a pass's call; cur is the file this call's
// last FILE named, and frame returns the file the call's next frames belong to.
func (d *Downlink) frame(p *pass, cur string, f Frame) string {
	switch f.Type {
	case TypeFile:
		if len(f.Body) < 3 {
			return ""
		}
		return d.begin(p, string(f.Body[2:]), int(binary.BigEndian.Uint16(f.Body)))
	case TypeChunk:
		if cur != "" && len(f.Body) >= 2 && len(f.Body)-2 <= MaxChunk {
			d.chunk(cur, int(binary.BigEndian.Uint16(f.Body)), f.Body[2:])
		}
		return cur
	case TypeEnd:
		if cur != "" && len(f.Body) == 4 {
			d.end(p, cur, binary.BigEndian.Uint32(f.Body))
		}
		return ""
	}
	return cur
}

// begin takes a FILE frame: a file the pass is bringing down, new or sent again on a
// redial.
func (d *Downlink) begin(p *pass, name string, total int) string {
	if name == "" || filepath.Base(name) != name || name == "." || name == ".." || strings.ContainsAny(name, `/\`) || total == 0 {
		return ""
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.files[name]; !ok {
		d.files[name] = &downFile{name: name, booking: p.r.booked[p.c.ID], total: total, chunks: map[int][]byte{}}
	}
	return name
}

// chunk keeps one chunk of a file, by its number.
func (d *Downlink) chunk(name string, n int, data []byte) {
	d.mu.Lock()
	defer d.mu.Unlock()
	f := d.files[name]
	if f == nil || f.landed || n >= f.total {
		return
	}
	if _, ok := f.chunks[n]; !ok {
		f.chunks[n] = bytes.Clone(data)
	}
}

// end takes an END frame: a file with every chunk whose CRC-32 matches lands; one with
// every chunk and the wrong CRC is dropped.
func (d *Downlink) end(p *pass, name string, crc uint32) {
	d.mu.Lock()
	f := d.files[name]
	if f == nil || f.landed || len(f.chunks) < f.total {
		d.mu.Unlock()
		return
	}
	var data []byte
	for n := range f.total {
		data = append(data, f.chunks[n]...)
	}
	if crc32.ChecksumIEEE(data) != crc {
		f.crcBad, f.chunks = true, map[int][]byte{}
		d.mu.Unlock()
		return
	}
	f.landed, f.crcBad, f.chunks = true, false, nil
	d.mu.Unlock()
	d.land(p, f, data)
}

// land saves a whole file, says so, and hands its landing to delivery.
func (d *Downlink) land(p *pass, f *downFile, data []byte) {
	r := p.r
	d.made.Do(func() { d.dir = os.MkdirAll(filepath.Join(r.w.Dir, "files"), 0o755) })
	err := d.dir
	if err == nil {
		err = replaceFile(filepath.Join(r.w.Dir, "files", f.name), data)
	}
	if err != nil {
		return // not saved, so it never lands
	}
	now := r.w.Now()
	sayAt(r.w, "LANDED %s", f.name)
	if f.booking.ID != "" {
		r.deliv.landed(Landing{BookingID: f.booking.ID, File: f.name, Bytes: int64(len(data)), At: now.Truncate(time.Second)}, f.booking)
	}
}

// missing returns the numbers of the chunks a file still needs, ascending.
func (f *downFile) missing() []string {
	var out []string
	for n := range f.total {
		if _, ok := f.chunks[n]; !ok {
			out = append(out, strconv.Itoa(n))
		}
	}
	return out
}

// writeIncomplete saves reports/incomplete.txt in dir: one line per file still incomplete,
// by name: "<name> missing 21 22 23" or "<name> crc mismatch".
func (d *Downlink) writeIncomplete(dir string) error {
	d.mu.Lock()
	var lines []string
	for _, f := range d.files {
		switch {
		case f.landed:
		case f.crcBad:
			lines = append(lines, f.name+" crc mismatch")
		default:
			lines = append(lines, f.name+" missing "+strings.Join(f.missing(), " "))
		}
	}
	d.mu.Unlock()
	slices.Sort(lines)
	var buf bytes.Buffer
	for _, l := range lines {
		fmt.Fprintln(&buf, l)
	}
	return saveWhole(filepath.Join(dir, "reports", "incomplete.txt"), buf.Bytes())
}
