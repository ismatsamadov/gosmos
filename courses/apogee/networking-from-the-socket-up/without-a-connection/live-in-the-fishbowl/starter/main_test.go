package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// frame is Furrow-A's telemetry line numbered seq. On Friday 16 June frame 84917 was made at
// 08:00:42, and one more each second, until LOS at 08:00:51.
func frame(seq uint32) string {
	at := time.Date(2028, 6, 16, 8, 0, 42, 0, time.UTC).Add(time.Duration(int(seq)-84917) * time.Second)
	return fmt.Sprintf("furrow-a %d %s batt_v=7.86 temp_c=12.5 mode=NOMINAL", seq, at.Format("2006-01-02T15:04:05Z"))
}

// frames returns the lines of the frames numbered from first to last.
func frames(first, last uint32) []string {
	var lines []string
	for seq := first; seq <= last; seq++ {
		lines = append(lines, frame(seq))
	}
	return lines
}

func TestWatch(t *testing.T) {
	t.Run("Friday's last ten seconds", func(t *testing.T) {
		l, rec, q := setup(t, "IDLE")
		stray := listen(t)
		send(t, rec, l, frame(84917), frame(84918))
		send(t, rec, l, frame(84917), frame(84918)) // the same datagram again
		send(t, rec, l, frame(84920))               // it overtook 84919
		send(t, rec, l, frame(84919))
		send(t, rec, l, frame(84922), frame(84923)) // 84921 was lost on the line
		send(t, stray, l, frame(99999))             // not the recorder's
		send(t, rec, l, frame(84924))
		send(t, stray, l, strings.Repeat("x", 600))
		send(t, rec, l, frame(84925), frame(84926))
		want := []string{frame(84917), frame(84918), frame(84920), frame(84922), frame(84923), frame(84924), frame(84925), frame(84926),
			"quiet after 84926: 8 shown, 2 missed, 3 dropped; recorder IDLE"}
		check(t, l, want)
		checkAsked(t, q, "STATUS 1", 1, 3)
	})
	t.Run("frames that stop mid-pass, the first question lost", func(t *testing.T) {
		l, rec, q := setup(t, "RECORDING furrow-a/YST/18796 after one lost")
		lines := frames(84635, 84639) // 07:56:00 to 07:56:04
		for _, line := range lines {
			send(t, rec, l, line)
		}
		check(t, l, append(lines, "quiet after 84639: 5 shown, 0 missed, 0 dropped; recorder RECORDING furrow-a/YST/18796"))
		checkAsked(t, q, "STATUS 1", 2, 3)
	})
	t.Run("a silent recorder", func(t *testing.T) {
		l, rec, q := setup(t, "silent")
		send(t, rec, l, frame(84917), frame(84918))
		check(t, l, []string{frame(84917), frame(84918), "quiet after 84918: 2 shown, 0 missed, 0 dropped; recorder silent"})
		checkAsked(t, q, "STATUS 1", 3, 3)
	})
	t.Run("no frames at all", func(t *testing.T) {
		l, _, q := setup(t, "IDLE")
		check(t, l, []string{"no frames; recorder IDLE"})
		checkAsked(t, q, "STATUS 1", 1, 3)
	})
	t.Run("two passes, two questions", func(t *testing.T) {
		l, rec, q := setup(t, "IDLE")
		send(t, rec, l, frame(84917))
		check(t, l, []string{frame(84917), "quiet after 84917: 1 shown, 0 missed, 0 dropped; recorder IDLE"})
		send(t, rec, l, frame(84926))
		check(t, l, []string{frame(84926), "quiet after 84926: 1 shown, 0 missed, 0 dropped; recorder IDLE"})
		q.mu.Lock()
		defer q.mu.Unlock()
		if len(q.sent) < 2 || q.sent[0] != "STATUS 1" || q.sent[len(q.sent)-1] != "STATUS 2" {
			t.Errorf("two Watches asked %q, want %q first and %q last: each question gets the next number", q.sent, "STATUS 1", "STATUS 2")
		}
	})
}

// check runs l.Watch, at most two seconds, and compares what it wrote with want, line by line.
func check(t *testing.T, l *Live, want []string) {
	t.Helper()
	var out strings.Builder
	done := make(chan error, 1)
	go func() { done <- l.Watch(&out) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Watch() = %v, want nil", err)
		}
	case <-time.After(2 * time.Second):
		l.Frames.Close()
		l.Questions.Close()
		<-done
		t.Fatalf("Watch() still running after 2s, with Quiet at 100ms and Wait at 50ms: stop when receive fails with a deadline")
	}
	got := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
	if out.Len() == 0 {
		got = nil
	}
	for i := range max(len(got), len(want)) {
		g, w := "(nothing)", "(nothing)"
		if i < len(got) {
			g = got[i]
		}
		if i < len(want) {
			w = want[i]
		}
		if g != w {
			t.Errorf("Watch() line %d = %q, want %q (it wrote %d lines, want %d)", i+1, g, w, len(got), len(want))
			return
		}
	}
}

// sentLog is the live view's questions socket, keeping a copy of every datagram sent from it.
type sentLog struct {
	net.PacketConn
	mu   sync.Mutex
	sent []string
}

// WriteTo sends p to addr and keeps a copy of it if it was sent.
func (s *sentLog) WriteTo(p []byte, addr net.Addr) (int, error) {
	n, err := s.PacketConn.WriteTo(p, addr)
	if err == nil {
		s.mu.Lock()
		s.sent = append(s.sent, string(p))
		s.mu.Unlock()
	}
	return n, err
}

// checkAsked checks that between least and most questions were sent, each one want.
func checkAsked(t *testing.T, q *sentLog, want string, least, most int) {
	t.Helper()
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.sent) < least || len(q.sent) > most {
		t.Errorf("Watch asked %q, %d in all; want %d to %d questions", q.sent, len(q.sent), least, most)
	}
	for _, s := range q.sent {
		if s != want {
			t.Errorf("Watch asked %q, want %q", s, want)
			return
		}
	}
}

// setup builds a live view on loopback, with Quiet at 100ms, Wait at 50ms and 3 tries, and a
// stand-in for the recorder's live port that answers its questions by plan: "IDLE" to every
// one, "RECORDING furrow-a/YST/18796 after one lost" to all but the first, or "silent". It
// returns the view, the recorder's socket and the view's questions socket.
func setup(t *testing.T, plan string) (*Live, net.PacketConn, *sentLog) {
	t.Helper()
	rec := listen(t)
	q := &sentLog{PacketConn: listen(t)}
	l := &Live{Frames: listen(t), Questions: q, Recorder: rec.LocalAddr(),
		Quiet: 100 * time.Millisecond, Wait: 50 * time.Millisecond, Tries: 3}
	done := make(chan struct{})
	go func() {
		defer close(done)
		buf := make([]byte, 1024)
		for heard := 0; ; heard++ {
			k, from, err := rec.ReadFrom(buf)
			if err != nil {
				return
			}
			n, ok := strings.CutPrefix(string(buf[:k]), "STATUS ")
			if _, err := strconv.ParseUint(n, 10, 32); !ok || err != nil {
				continue
			}
			switch {
			case plan == "IDLE":
				rec.WriteTo([]byte("STATUS "+n+" IDLE"), from)
			case strings.HasSuffix(plan, " after one lost") && heard > 0:
				rec.WriteTo([]byte("STATUS "+n+" "+strings.TrimSuffix(plan, " after one lost")), from)
			}
		}
	}()
	t.Cleanup(func() {
		rec.Close()
		<-done
	})
	return l, rec, q
}

// send sends the lines, joined by newlines, as one datagram from pc to the live view's frames.
func send(t *testing.T, pc net.PacketConn, l *Live, lines ...string) {
	t.Helper()
	d := strings.Join(lines, "\n")
	if _, err := pc.WriteTo([]byte(d), l.Frames.LocalAddr()); err != nil {
		t.Fatalf("sending a datagram of %d bytes: %v", len(d), err)
	}
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
