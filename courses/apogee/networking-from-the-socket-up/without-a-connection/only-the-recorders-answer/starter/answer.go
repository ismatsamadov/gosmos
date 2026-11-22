package main

import "net"

// answerTo returns the text of datagram d, sent from the address from, and
// true when d is the recorder's answer to question n: "IDLE" for "STATUS 7 IDLE".
func answerTo(n uint32, d []byte, from, recorder net.Addr) (string, bool) {
	// TODO: accept only the recorder's "STATUS <n> <text>"
	return "", false
}
