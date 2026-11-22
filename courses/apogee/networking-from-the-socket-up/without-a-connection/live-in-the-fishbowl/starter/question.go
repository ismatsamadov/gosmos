package main

import (
	"errors"
	"net"
	"strconv"
)

// question is question n for the recorder's live port: "STATUS 7".
func question(n uint32) []byte {
	return []byte("STATUS " + strconv.FormatUint(uint64(n), 10))
}

// readAnswer reads datagrams from pc until one is the recorder's answer to
// question n, and returns its text. It skips datagrams over maxDatagram and
// returns any other error.
func readAnswer(pc net.PacketConn, recorder net.Addr, n uint32) (string, error) {
	for {
		d, from, err := readDatagram(pc)
		if errors.Is(err, ErrTooBig) {
			continue
		}
		if err != nil {
			return "", err
		}
		if text, ok := answerTo(n, d, from, recorder); ok {
			return text, nil
		}
	}
}
