package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Live is the live view at Port Ellery: it shows each Ystmark pass as the
// recorder sends it, and asks the recorder how it is when the frames stop.
type Live struct {
	Frames    net.PacketConn // the recorder sends its frames here
	Questions net.PacketConn // questions go out and answers come back here
	Recorder  net.Addr       // the recorder's live port
	Quiet     time.Duration  // no datagram for this long: the pass has gone quiet (5s)
	Wait      time.Duration  // how long each try of a question waits (1s)
	Tries     int            // tries of a question before the recorder counts as silent (3)
	asked     uint32         // the number of the last question asked
}

// receive returns the recorder's next datagram on l.Frames, waiting at most
// l.Quiet for each; it skips datagrams over maxDatagram and any from another sender.
func (l *Live) receive() ([]byte, error) {
	for {
		l.Frames.SetReadDeadline(time.Now().Add(l.Quiet))
		d, from, err := readDatagram(l.Frames)
		if errors.Is(err, ErrTooBig) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if from.String() == l.Recorder.String() {
			return d, nil
		}
	}
}

// seqOf returns a telemetry line's sequence number, its second field, and
// whether it has one.
func seqOf(line string) (uint32, bool) {
	f := strings.Fields(line)
	if len(f) < 2 {
		return 0, false
	}
	n, err := strconv.ParseUint(f[1], 10, 32)
	return uint32(n), err == nil
}

// lastLine is the live view's line for a pass gone quiet: what t counted, and
// the recorder's answer.
func lastLine(t Tracker, answer string) string {
	if t.Shown == 0 {
		return "no frames; recorder " + answer
	}
	return fmt.Sprintf("quiet after %d: %d shown, %d missed, %d dropped; recorder %s",
		t.Last, t.Shown, t.Missed, t.Dropped, answer)
}
