package main

import (
	"net"
	"strings"
)

// answerTo returns the text of datagram d, sent from the address from, and
// true when d is the recorder's answer to question n: "IDLE" for "STATUS 7 IDLE".
func answerTo(n uint32, d []byte, from, recorder net.Addr) (string, bool) {
	if from.String() != recorder.String() {
		return "", false
	}
	text, ok := strings.CutPrefix(string(d), string(question(n))+" ")
	if !ok || text == "" {
		return "", false
	}
	return text, true
}
