package main

import (
	"errors"
	"net"
)

// maxDatagram is the most bytes one datagram on the live port may hold.
const maxDatagram = 512

// ErrTooBig is a datagram over maxDatagram bytes, which the read cut short.
var ErrTooBig = errors.New("datagram over 512 bytes")

// readDatagram reads one datagram from pc and returns its bytes and its sender;
// for one over maxDatagram bytes, it returns nil, the sender and ErrTooBig.
func readDatagram(pc net.PacketConn) ([]byte, net.Addr, error) {
	buf := make([]byte, maxDatagram+1)
	n, from, err := pc.ReadFrom(buf)
	if err != nil {
		return nil, nil, err
	}
	if n > maxDatagram {
		return nil, from, ErrTooBig
	}
	return buf[:n], from, nil
}
