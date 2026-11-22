package main

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

// Frames from the protocol's table, as they go on the wire.
const (
	pong = "\x00\x09\x82\x00\x00\x00\x00\x0d\x04\x14\xef"
	// welcome10 is YST-1's WELCOME after command 10: a 0x0a, a newline, in its middle.
	welcome10 = "\x00\x10\x81\x00\x00\x00\x00\x01\x02\x04\x00\x00\x00\x0a\x7e\x5d\x00\x81"
)

func TestReadPayload(t *testing.T) {
	tests := []struct {
		name   string
		stream io.Reader
		want   string
	}{
		{"a PONG", strings.NewReader(pong), pong[2:]},
		{"a WELCOME with 0x0a in its middle", strings.NewReader(welcome10), welcome10[2:]},
		{"that WELCOME one byte at a time", iotest.OneByteReader(strings.NewReader(welcome10)), welcome10[2:]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readPayload(bufio.NewReader(tt.stream))
			if string(got) != tt.want || err != nil {
				t.Errorf("readPayload(%s) = % x, %v; want % x, nil", tt.name, got, err, tt.want)
			}
		})
	}
	t.Run("two frames, then the end of the stream", func(t *testing.T) {
		br := bufio.NewReader(strings.NewReader(pong + welcome10))
		for i, want := range []string{pong[2:], welcome10[2:]} {
			if got, err := readPayload(br); string(got) != want || err != nil {
				t.Fatalf("readPayload(PONG and WELCOME), frame %d = % x, %v; want % x, nil", i+1, got, err, want)
			}
		}
		if got, err := readPayload(br); !errors.Is(err, io.EOF) {
			t.Errorf("readPayload after the last frame = % x, %v; want io.EOF: the stream ended between frames", got, err)
		}
	})
	t.Run("a stream that ends inside a frame", func(t *testing.T) {
		got, err := readPayload(bufio.NewReader(strings.NewReader(pong[:6])))
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Errorf("readPayload(6 of a PONG's 11 bytes) = % x, %v; want io.ErrUnexpectedEOF", got, err)
		}
	})
}
