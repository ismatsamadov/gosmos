package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"
	"time"
)

// bench returns testdata/bench.log: the bench replay's lines from bookingd, passplan and uplinkd at
// headquarters and at Lantern Atoll, one log after another, with two lines that aren't JSON, the
// bench's own heading and a line cut short where the copy stopped.
func bench(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/bench.log")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// want is one line follow must return.
type want struct {
	time, program, site, msg string
}

// describe writes a line as the messages show it: time, program, site and message.
func describe(l Line) string {
	return l.Time.UTC().Format(time.RFC3339) + " " + l.Program + " " + l.Site + " " + l.Msg
}

func TestFollow(t *testing.T) {
	tests := []struct {
		pass string
		want []want
	}{
		{"cirrine-1/LAN/16402", []want{
			{"2028-11-06T10:14:03Z", "bookingd", "HQ", "booking created"},
			{"2028-11-07T00:00:41Z", "passplan", "HQ", "pass planned"},
			{"2028-11-08T01:44:40Z", "uplinkd", "HQ", "command sealed"},
			{"2028-11-08T01:48:30Z", "uplinkd", "HQ", "lines delivered"},
			{"2028-11-08T01:48:30Z", "uplinkd", "LAN", "lines queued"},
			{"2028-11-08T01:53:09Z", "uplinkd", "LAN", "pass started"},
			{"2028-11-08T02:02:39Z", "uplinkd", "LAN", "pass failed"},
		}},
		{"furrow-a/LAN/20998", []want{
			{"2028-11-06T15:20:55Z", "bookingd", "HQ", "booking created"},
			{"2028-11-07T00:00:41Z", "passplan", "HQ", "pass planned"},
			{"2028-11-08T02:37:45Z", "uplinkd", "LAN", "pass started"},
			{"2028-11-08T02:46:24Z", "uplinkd", "LAN", "pass failed"},
		}},
	}
	data := bench(t)
	written := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	for _, tt := range tests {
		t.Run(tt.pass, func(t *testing.T) {
			call := `follow(bench.log, "` + tt.pass + `")`
			lines, broken, err := follow(bytes.NewReader(data), tt.pass)
			if err != nil {
				t.Fatalf("%s returned the error %v, want nil", call, err)
			}
			if broken != 2 {
				t.Errorf("%s counted %d broken lines, want 2: the heading and the line cut short", call, broken)
			}
			if len(lines) != len(tt.want) {
				t.Fatalf("%s returned %d lines of %d in the file, want %d", call, len(lines), len(written), len(tt.want))
			}
			for i, w := range tt.want {
				l := lines[i]
				if got := describe(l); got != w.time+" "+w.program+" "+w.site+" "+w.msg || l.PassID != tt.pass {
					t.Errorf("%s line %d = %q of %q, want %q of %q", call, i+1, got, l.PassID, w.time+" "+w.program+" "+w.site+" "+w.msg, tt.pass)
				}
				found := false
				for _, raw := range written {
					found = found || raw == l.Raw
				}
				if !found {
					t.Errorf("%s line %d has Raw %q, want the line exactly as the file has it", call, i+1, l.Raw)
				}
			}
		})
	}
}

// TestFollowByKey reads two of bookingd's lines: its access line for the request that looked
// Cirrine-1's pass up, which names the pass in its message but carries no pass_id, and the
// booking that followed.
func TestFollowByKey(t *testing.T) {
	const bookingd = `{"time":"2028-11-06T10:13:58Z","level":"INFO","msg":"req-004209 GET /passes/cirrine-1/LAN/16402 200 144B 1ms","program":"bookingd","site":"HQ"}
{"time":"2028-11-06T10:14:03Z","level":"INFO","msg":"booking created","program":"bookingd","site":"HQ","booking_id":"bk-5611","pass_id":"cirrine-1/LAN/16402","request_id":"req-004211"}
`
	lines, _, err := follow(strings.NewReader(bookingd), "cirrine-1/LAN/16402")
	if err != nil {
		t.Fatalf("follow(bookingd's access line and booking) returned the error %v, want nil", err)
	}
	if len(lines) != 1 || lines[0].Msg != "booking created" {
		t.Errorf(`follow(bookingd's access line and booking) returned %d lines, want 1, "booking created": the access line names the pass in its message, but it has no pass_id`, len(lines))
	}
}

// errLinkDown is the error a station's log stream ends with when its link drops.
var errLinkDown = errors.New("link to LAN lost")

func TestFollowReadError(t *testing.T) {
	data := bench(t)
	first := data[:bytes.IndexByte(data, '\n')+1]
	r := io.MultiReader(bytes.NewReader(first), iotest.ErrReader(errLinkDown))
	if _, _, err := follow(r, "cirrine-1/LAN/16402"); !errors.Is(err, errLinkDown) {
		t.Errorf("follow(a stream that ends in %q) returned the error %v, want that error", errLinkDown, err)
	}
}
