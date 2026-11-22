package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"testing"
)

// pong is the controller's PONG as it goes on the wire: length 9.
const pong = "\x00\x09\x82\x00\x00\x00\x00\x0d\x04\x14\xef"

func TestReadPayload(t *testing.T) {
	longest := binary.BigEndian.AppendUint16(nil, 1024)
	for i := range 1024 {
		longest = append(longest, byte(i))
	}
	allowed := []struct {
		name  string
		frame []byte
	}{
		{"9, a PONG", []byte(pong)},
		{"1024, the longest frame", longest},
	}
	for _, a := range allowed {
		t.Run("length "+a.name, func(t *testing.T) {
			got, err := readPayload(bufio.NewReader(bytes.NewReader(a.frame)))
			if !bytes.Equal(got, a.frame[2:]) || err != nil {
				t.Errorf("readPayload(length %s) = %d bytes, %v; want its %d bytes, nil", a.name, len(got), err, len(a.frame)-2)
			}
		})
	}
	refused := []struct {
		name   string
		length uint16
	}{
		{"one short of the shortest frame", 8},
		{"one past the longest frame", 1025},
		{"zero", 0},
		{"Sunniva's cable test, ff ff", 65535},
	}
	for _, r := range refused {
		t.Run(r.name, func(t *testing.T) {
			line := &lengthThenSilence{rest: binary.BigEndian.AppendUint16(nil, r.length)}
			_, err := readPayload(bufio.NewReader(line))
			if line.readOn {
				t.Fatalf("readPayload(length %d, then nothing) read on past the length; on the line it would wait forever", r.length)
			}
			want := fmt.Sprintf("bad length %d", r.length)
			if !errors.Is(err, ErrBadLength) || err.Error() != want {
				t.Errorf("readPayload(length %d) = %v, want ErrBadLength, wrapped as %q", r.length, err, want)
			}
		})
	}
}

// lengthThenSilence is a line that has brought a frame's two length bytes and nothing
// since: any read past them is a read that would wait forever.
type lengthThenSilence struct {
	rest   []byte
	readOn bool
}

// Read hands out the length bytes, then notes that the reader went on reading.
func (l *lengthThenSilence) Read(p []byte) (int, error) {
	if len(l.rest) == 0 {
		l.readOn = true
		return 0, errors.New("nothing more has arrived")
	}
	n := copy(p, l.rest)
	l.rest = l.rest[n:]
	return n, nil
}
