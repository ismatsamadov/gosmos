package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"strings"
)

// Welcome is a decoded WELCOME frame.
type Welcome struct {
	Dish    uint8
	FwMajor uint8
	FwMinor uint8
	LastSeq uint32
}

// helloDishd is HELLO from dishd as it goes on the wire: its length, its type, its
// sequence number, the name "dishd" and a CRC.
const helloDishd = "\x00\x0e\x01\x00\x00\x00\x00dishd\x97\x98\xf8\xab"

// readBannerLine reads the controller's first line a byte at a time, so it takes nothing
// after the line from the connection, and returns the line without its line ending.
func readBannerLine(conn net.Conn) (string, error) {
	var line []byte
	b := make([]byte, 1)
	for len(line) < 64 {
		if _, err := io.ReadFull(conn, b); err != nil {
			return "", err
		}
		if b[0] == '\n' {
			return strings.TrimSuffix(string(line), "\r"), nil
		}
		line = append(line, b[0])
	}
	return "", fmt.Errorf("no line ending in the first 64 bytes: %q", line)
}

// handshake reads the controller's banner, says HELLO and reads the WELCOME that answers
// it. It returns the controller's name from the banner and the WELCOME.
func handshake(conn net.Conn) (name string, w Welcome, err error) {
	banner, err := readBannerLine(conn)
	if err != nil {
		return "", Welcome{}, err
	}
	name, ok := strings.CutSuffix(banner, " READY PROTO 1")
	if !ok || name == "" {
		return "", Welcome{}, fmt.Errorf("not a controller's banner: %q", banner)
	}
	if _, err := io.WriteString(conn, helloDishd); err != nil {
		return "", Welcome{}, err
	}
	// A WELCOME is 18 bytes: its length (16), type 0x81, a sequence number of 0,
	// the dish, the firmware's major and minor numbers, the last command, then a CRC.
	frame := make([]byte, 18)
	if _, err := io.ReadFull(conn, frame); err != nil {
		return "", Welcome{}, err
	}
	crc := crc32.ChecksumIEEE(frame[2:14])
	if string(frame[:3]) != "\x00\x10\x81" || binary.BigEndian.Uint32(frame[14:]) != crc {
		return "", Welcome{}, fmt.Errorf("not a WELCOME: % x", frame)
	}
	w = Welcome{
		Dish:    frame[7],
		FwMajor: frame[8],
		FwMinor: frame[9],
		LastSeq: binary.BigEndian.Uint32(frame[10:14]),
	}
	return name, w, nil
}
