package main

import (
	"fmt"
	"net"
	"time"
)

// firstContact says hello to the controller at addr and describes it in one line.
func firstContact(addr string) (string, error) {
	_ = net.DialTimeout
	_ = time.Second
	_ = fmt.Sprintf
	return "", nil
}
