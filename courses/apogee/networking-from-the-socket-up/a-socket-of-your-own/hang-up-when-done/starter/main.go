package main

import "net"

// bannerOf dials the controller at addr and returns its banner.
func bannerOf(addr string) (string, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	return readBannerLine(conn)
}
