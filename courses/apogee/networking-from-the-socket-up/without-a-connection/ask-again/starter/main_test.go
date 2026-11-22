package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const recording = "RECORDING furrow-a/YST/18783"

func TestAsk(t *testing.T) {
	for _, c := range []struct {
		name string
		plan string // how the stand-in recorder answers
	}{
		{"an answer at once", "every"},
		{"the first question lost", "all but the first"},
		{"the first answer late", "the first, late"},
	} {
		t.Run(c.name, func(t *testing.T) {
			q, rec := questions(t), standIn(t, c.plan)
			text, err := askWithin(q, rec, 7, 50*time.Millisecond, 3)
			if text != recording || err != nil {
				t.Errorf("ask(question 7, wait 50ms, 3 tries) with %s = %q, %v; want %q, nil", c.name, text, err, recording)
			}
			checkSent(t, q, 1, 3)
		})
	}
	t.Run("a silent recorder, 3 tries", func(t *testing.T) {
		q, rec := questions(t), standIn(t, "none")
		start := time.Now()
		text, err := askWithin(q, rec, 7, 50*time.Millisecond, 3)
		if !errors.Is(err, os.ErrDeadlineExceeded) {
			t.Errorf("ask(question 7, wait 50ms, 3 tries) with no answer = %q, %v; want a deadline error", text, err)
		}
		if took := time.Since(start); took < 150*time.Millisecond {
			t.Errorf("ask with no answer gave up after %v, want 3 waits of 50ms", took.Round(time.Millisecond))
		}
		checkSent(t, q, 3, 3)
	})
	t.Run("a silent recorder, 1 try", func(t *testing.T) {
		q, rec := questions(t), standIn(t, "none")
		_, err := askWithin(q, rec, 7, 50*time.Millisecond, 1)
		if !errors.Is(err, os.ErrDeadlineExceeded) {
			t.Errorf("ask(question 7, wait 50ms, 1 try) with no answer = %v, want a deadline error", err)
		}
		checkSent(t, q, 1, 1)
	})
	t.Run("a closed socket", func(t *testing.T) {
		q, rec := questions(t), standIn(t, "every")
		q.Close()
		_, err := askWithin(q, rec, 7, 50*time.Millisecond, 3)
		if !errors.Is(err, net.ErrClosed) {
			t.Errorf("ask from a closed socket = %v, want an error matching net.ErrClosed", err)
		}
	})
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

// questions opens the live view's questions socket on loopback, until the test ends.
func questions(t *testing.T) *sentLog {
	t.Helper()
	return &sentLog{PacketConn: listen(t)}
}

// checkSent checks that between least and most questions were sent, every one of them STATUS 7.
func checkSent(t *testing.T, q *sentLog, least, most int) {
	t.Helper()
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.sent) < least || len(q.sent) > most {
		want := fmt.Sprintf("%d to %d", least, most)
		if least == most {
			want = fmt.Sprint(least)
		}
		t.Errorf("ask sent %q, %d in all; want %s", q.sent, len(q.sent), want)
	}
	for _, s := range q.sent {
		if s != "STATUS 7" {
			t.Errorf("ask sent %q, want %q: a question asked again keeps its number", s, "STATUS 7")
			return
		}
	}
}

// standIn plays the recorder's live port until the test ends, answering each question by
// its plan: "every" question, "all but the first", "the first, late" (once the second question
// arrives, and no other), or "none". It returns the port's address.
func standIn(t *testing.T, plan string) net.Addr {
	t.Helper()
	pc := listen(t)
	done := make(chan struct{})
	go func() {
		defer close(done)
		buf := make([]byte, 1024)
		var first string // the first question, held back by "the first, late"
		for heard := 0; ; heard++ {
			k, from, err := pc.ReadFrom(buf)
			if err != nil {
				return
			}
			q := string(buf[:k])
			n, ok := strings.CutPrefix(q, "STATUS ")
			if _, err := strconv.ParseUint(n, 10, 32); !ok || err != nil {
				continue
			}
			answer := ""
			switch {
			case plan == "every", plan == "all but the first" && heard > 0:
				answer = q
			case plan == "the first, late" && heard == 0:
				first = q
			case plan == "the first, late" && heard == 1:
				answer = first
			}
			if answer != "" {
				pc.WriteTo([]byte(answer+" "+recording), from)
			}
		}
	}()
	t.Cleanup(func() {
		pc.Close()
		<-done
	})
	return pc.LocalAddr()
}

// askWithin calls ask and waits for it at most a second; if it is still waiting then, it
// closes pc so the call ends, and fails the test.
func askWithin(pc *sentLog, recorder net.Addr, n uint32, wait time.Duration, tries int) (string, error) {
	type result struct {
		text string
		err  error
	}
	got := make(chan result, 1)
	go func() {
		text, err := ask(pc, recorder, n, wait, tries)
		got <- result{text, err}
	}()
	select {
	case r := <-got:
		return r.text, r.err
	case <-time.After(time.Second):
		pc.Close()
		<-got
		return "", errors.New("ask still waiting after a second")
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
