package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

// owed is the nine frames of Furrow-B the recorder owed on Tuesday at 08:06, 72 bytes each.
var owed = func() []string {
	var lines []string
	at := time.Date(2028, 6, 13, 8, 6, 10, 0, time.UTC)
	for i := range 9 {
		lines = append(lines, fmt.Sprintf("furrow-b %d %s batt_v=7.91 temp_c=11.0 mode=NOMINAL",
			31877+i, at.Add(time.Duration(i)*time.Second).Format("2006-01-02T15:04:05Z")))
	}
	return lines
}()

func TestPack(t *testing.T) {
	a255, b256, c256 := strings.Repeat("a", 255), strings.Repeat("b", 256), strings.Repeat("c", 256)
	full := strings.Repeat("d", 512)
	for _, c := range []struct {
		name  string
		lines []string
		want  [][]string // each datagram's lines
	}{
		{"the nine frames Furrow-B owed", owed, [][]string{owed[:7], owed[7:]}},
		{"one frame", owed[:1], [][]string{owed[:1]}},
		{"512 bytes: 255, a newline and 256", []string{a255, b256}, [][]string{{a255, b256}}},
		{"513 bytes: 256, a newline and 256", []string{b256, c256}, [][]string{{b256}, {c256}}},
		{"one line of 512 bytes", []string{full}, [][]string{{full}}},
		{"no lines", nil, nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := pack(c.lines)
			if err != nil {
				t.Fatalf("pack(%s) = %v, want %s and nil", c.name, err, describe(c.want))
			}
			want := make([]string, len(c.want))
			for i, lines := range c.want {
				want[i] = strings.Join(lines, "\n")
			}
			if len(got) != len(want) {
				t.Fatalf("pack(%s) = %s, want %s", c.name, describeBytes(got), describe(c.want))
			}
			for i := range want {
				if string(got[i]) != want[i] {
					t.Errorf("pack(%s): datagram %d = %d bytes starting %q, want %d bytes starting %q", c.name, i+1, len(got[i]), start(string(got[i])), len(want[i]), start(want[i]))
				}
			}
		})
	}
	t.Run("a line of 513 bytes", func(t *testing.T) {
		got, err := pack([]string{owed[0], strings.Repeat("e", 513)})
		if !errors.Is(err, ErrTooBig) {
			t.Errorf("pack(a frame and a line of 513 bytes) = %s, %v; want an error wrapping ErrTooBig", describeBytes(got), err)
		}
	})
	t.Run("through the live view", func(t *testing.T) {
		got, err := pack(owed)
		if err != nil {
			t.Fatalf("pack(the nine frames) = %v, want nil", err)
		}
		view, rec := listen(t), listen(t)
		view.SetReadDeadline(time.Now().Add(time.Second))
		for _, d := range got {
			if _, err := rec.WriteTo(d, view.LocalAddr()); err != nil {
				t.Fatalf("sending a datagram of %d bytes: %v", len(d), err)
			}
		}
		var lines []string
		for range got {
			d, _, err := readDatagram(view)
			if err != nil {
				t.Fatalf("readDatagram() of a packed datagram = %v, want nil", err)
			}
			lines = append(lines, strings.Split(string(d), "\n")...)
		}
		if strings.Join(lines, "\n") != strings.Join(owed, "\n") {
			t.Errorf("the live view read %d lines, want the nine frames in order", len(lines))
		}
	})
}

// describe sums up datagrams given as lines: "2 datagrams: 7 lines (510 bytes), 2 lines (145 bytes)".
func describe(datagrams [][]string) string {
	var parts []string
	for _, lines := range datagrams {
		parts = append(parts, fmt.Sprintf("%s (%d bytes)", count(len(lines), "line"), len(strings.Join(lines, "\n"))))
	}
	return fmt.Sprintf("%s: %s", count(len(datagrams), "datagram"), strings.Join(parts, ", "))
}

// describeBytes sums up datagrams as pack returns them.
func describeBytes(datagrams [][]byte) string {
	var parts []string
	for _, d := range datagrams {
		parts = append(parts, fmt.Sprintf("%s (%d bytes)", count(strings.Count(string(d), "\n")+1, "line"), len(d)))
	}
	return fmt.Sprintf("%s: %s", count(len(datagrams), "datagram"), strings.Join(parts, ", "))
}

// count writes n things: "1 line", "2 lines".
func count(n int, thing string) string {
	if n == 1 {
		return "1 " + thing
	}
	return fmt.Sprintf("%d %ss", n, thing)
}

// start is the first 20 bytes of s.
func start(s string) string {
	return s[:min(len(s), 20)]
}

// listen opens a UDP socket on 127.0.0.1, on a port the system picks, until the test ends.
func listen(t *testing.T) net.PacketConn {
	t.Helper()
	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { pc.Close() })
	return pc
}
