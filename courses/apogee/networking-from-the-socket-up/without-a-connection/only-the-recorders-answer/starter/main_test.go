package main

import (
	"net"
	"testing"
	"time"
)

func TestAnswerTo(t *testing.T) {
	ip := net.IPv4(127, 0, 0, 1)
	recorder := &net.UDPAddr{IP: ip, Port: 7401}
	sameAddress := &net.UDPAddr{IP: ip, Port: 7401} // the recorder's address, another value
	otherPort := &net.UDPAddr{IP: ip, Port: 7402}
	for _, c := range []struct {
		name string
		n    uint32
		d    string
		from net.Addr
		text string
		ok   bool
	}{
		{"the recorder recording", 7, "STATUS 7 RECORDING furrow-a/YST/18781", sameAddress, "RECORDING furrow-a/YST/18781", true},
		{"the recorder idle", 7, "STATUS 7 IDLE", sameAddress, "IDLE", true},
		{"question 4294967295", 4294967295, "STATUS 4294967295 IDLE", sameAddress, "IDLE", true},
		{"from another port", 7, "STATUS 7 IDLE", otherPort, "", false},
		{"a late answer to question 6", 7, "STATUS 6 IDLE", sameAddress, "", false},
		{"the answer to question 70", 7, "STATUS 70 IDLE", sameAddress, "", false},
		{"no text after the number", 7, "STATUS 7", sameAddress, "", false},
		{"a frame", 7, "furrow-a 77712 2028-06-15T11:20:04Z batt_v=7.86 temp_c=12.5 mode=NOMINAL", sameAddress, "", false},
	} {
		t.Run(c.name, func(t *testing.T) {
			text, ok := answerTo(c.n, []byte(c.d), c.from, recorder)
			if ok != c.ok || ok && text != c.text {
				t.Errorf("answerTo(%d, %q from %v) = %q, %t; want %q, %t", c.n, c.d, c.from, text, ok, c.text, c.ok)
			}
		})
	}
	t.Run("through readAnswer", func(t *testing.T) {
		view, rec, stray := listen(t), listen(t), listen(t)
		send(t, stray, view, "STATUS 7 RECORDING furrow-a/YST/00000") // from anyone
		send(t, rec, view, "STATUS 6 IDLE")                           // a question that gave up
		send(t, rec, view, "STATUS 7 RECORDING furrow-a/YST/18783")
		view.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
		text, err := readAnswer(view, rec.LocalAddr(), 7)
		if text != "RECORDING furrow-a/YST/18783" || err != nil {
			t.Errorf("readAnswer(question 7) after a stray's answer and question 6's = %q, %v; want %q, nil", text, err, "RECORDING furrow-a/YST/18783")
		}
	})
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

// send sends s as one datagram from one socket to another.
func send(t *testing.T, from, to net.PacketConn, s string) {
	t.Helper()
	if _, err := from.WriteTo([]byte(s), to.LocalAddr()); err != nil {
		t.Fatalf("sending %q: %v", s, err)
	}
}
