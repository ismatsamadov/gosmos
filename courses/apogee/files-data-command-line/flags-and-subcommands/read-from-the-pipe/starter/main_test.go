package main

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// fromRecorder is what the recorder writes on the pipe.
const fromRecorder = `{"id":"bk-3301","sat_id":"pipit-2","station":"SQT","aos":"2027-10-06T08:20:00Z"}
{"id":"bk-3302","sat_id":"furrow-b","station":"SQT","aos":"2027-10-06T13:47:30Z"}
`

// fromStick is a booking log saved on the memory stick.
const fromStick = `{"id":"bk-3290","sat_id":"cirrine-1","station":"PEL","aos":"2027-10-06T05:02:00Z"}
`

// ids lists the bookings' IDs in order.
func ids(bookings []Booking) []string {
	var out []string
	for _, b := range bookings {
		out = append(out, b.ID)
	}
	return out
}

// readerFunc adapts a test function to io.Reader.
type readerFunc func([]byte) (int, error)

// Read calls the supplied test function.
func (f readerFunc) Read(p []byte) (int, error) { return f(p) }

func TestLoadBookings(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "week.json"), []byte(fromStick), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, file, pipe string
		want             []string
		err              error
	}{
		{"the recorder's pipe", "", fromRecorder, []string{"bk-3301", "bk-3302"}, nil},
		{"an empty pipe", "", "", nil, nil},
		{"a file from the stick", "week.json", "", []string{"bk-3290"}, nil},
		{"a missing file", "missing.json", "", nil, fs.ErrNotExist},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdin io.Reader = strings.NewReader(tc.pipe)
			name := tc.file
			if name != "" {
				name = filepath.Join(dir, name)
				stdin = readerFunc(func([]byte) (int, error) {
					t.Errorf("loadBookings(%q, stdin) read stdin, want only the file read", tc.file)
					return 0, io.EOF
				})
			}
			got, err := loadBookings(name, stdin)
			if !slices.Equal(ids(got), tc.want) || !errors.Is(err, tc.err) {
				t.Errorf("loadBookings(%q, stdin) = %v, %v, want %v, %v", tc.file, ids(got), err, tc.want, tc.err)
			}
		})
	}
}
