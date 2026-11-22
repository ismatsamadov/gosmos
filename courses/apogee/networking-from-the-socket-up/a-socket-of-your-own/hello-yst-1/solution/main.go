package main

import (
	"fmt"
	"net"
	"time"
)

// firstContact says hello to the controller at addr and describes it in one line.
func firstContact(addr string) (string, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	name, w, err := handshake(conn)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s firmware %d.%d, last command %d", name, w.FwMajor, w.FwMinor, w.LastSeq), nil
}
