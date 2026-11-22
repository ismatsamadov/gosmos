package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"strings"
)

// Frame is one pass-line frame, without its length and CRC.
type Frame struct {
	Type byte
	Seq  uint32 // a command's number; 0 on every other frame
	Body []byte
}

var (
	errBadLength = errors.New("bad length") // wrapped: "bad length 8": the stream is out of step
	errBadCRC    = errors.New("bad crc")    // the frame was damaged on the way
)

// readBanner reads the station's banner, "<CODE> READY PROTO 1" ended by CR LF or LF, and
// refuses any other line.
func readBanner(br *bufio.Reader, code string) error {
	line, err := br.ReadString('\n')
	if err != nil {
		return err
	}
	text := strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
	if text != code+" READY PROTO 1" {
		return fmt.Errorf("bad banner %q", text)
	}
	return nil
}

// encodeFrame returns f's bytes: the length, the type, the seq, the body, then the CRC-32
// of the type, seq and body, every number big-endian.
func encodeFrame(f Frame) []byte {
	b := make([]byte, 2+1+4+len(f.Body)+4)
	binary.BigEndian.PutUint16(b, uint16(MinFrameLen+len(f.Body)))
	b[2] = f.Type
	binary.BigEndian.PutUint32(b[3:], f.Seq)
	copy(b[7:], f.Body)
	binary.BigEndian.PutUint32(b[7+len(f.Body):], crc32.ChecksumIEEE(b[2:7+len(f.Body)]))
	return b
}

// writeFrame writes f in one Write.
func writeFrame(w io.Writer, f Frame) error {
	if len(f.Body) > MaxBody {
		return fmt.Errorf("frame body of %d bytes, over %d", len(f.Body), MaxBody)
	}
	_, err := w.Write(encodeFrame(f))
	return err
}

// readFrame reads one frame. A damaged frame is returned with errBadCRC, and reading can go
// on after it; a length out of range is errBadLength, and nothing after it can be trusted.
func readFrame(br *bufio.Reader) (Frame, error) {
	var head [2]byte
	if _, err := io.ReadFull(br, head[:]); err != nil {
		return Frame{}, err
	}
	n := int(binary.BigEndian.Uint16(head[:]))
	if n < MinFrameLen || n > MaxFrameLen {
		return Frame{}, fmt.Errorf("%w %d", errBadLength, n)
	}
	p := make([]byte, n)
	if _, err := io.ReadFull(br, p); err != nil {
		return Frame{}, err
	}
	f := Frame{Type: p[0], Seq: binary.BigEndian.Uint32(p[1:5]), Body: p[5 : n-4]}
	if crc32.ChecksumIEEE(p[:n-4]) != binary.BigEndian.Uint32(p[n-4:]) {
		return f, errBadCRC
	}
	return f, nil
}
