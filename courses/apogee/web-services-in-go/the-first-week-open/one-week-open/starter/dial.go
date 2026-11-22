package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"syscall"
	"time"
)

// errInside is the refusal of a webhook address that is not out on the internet.
var errInside = errors.New("webhook address is inside a network")

// checkDial refuses an address inside a network: loopback, private, link-local or unspecified, and one that doesn't parse. As a net.Dialer's Control hook it sees the very address about to be dialed.
func checkDial(network, address string, _ syscall.RawConn) error {
	ap, err := netip.ParseAddrPort(address)
	ip := ap.Addr().Unmap()
	if err != nil || ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified() {
		return fmt.Errorf("%w: %s", errInside, address)
	}
	return nil
}

// webhookClient returns the client webhooks go out with: 5 seconds at most, no redirect followed, and no connection inside a network, whatever name led there.
func webhookClient() *http.Client {
	dialer := &net.Dialer{Timeout: 5 * time.Second, Control: checkDial}
	return &http.Client{
		Timeout:       5 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		Transport:     &http.Transport{DialContext: dialer.DialContext},
	}
}
