package main

import (
	"crypto/ed25519"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"strings"
	"time"
)

// Slot is one pass of the day and the times uplinkd seals a PING for it.
type Slot struct {
	Pass
	SealAt []time.Time // each PING reaches the satellite 5 seconds after it is sealed
}

// Report is the day's log: one line per pass, in AOS order.
type Report struct {
	Flown, Missed int
	Lines         []string
}

// Day is a day of passes through a rotation, as Ingrid runs it.
type Day struct {
	Slots      []Slot                     // the day's passes, in any order
	Ring       *KeyRing                   // uplinkd's master key versions
	Old, New   map[string]tls.Certificate // each station's certificate before and after its switch, by code
	SwitchFrom time.Time                  // no station switches before this
	List       []byte                     // revocation list 1, signed by the stations CA
	ListAt     time.Time                  // when list 1 comes into force

	roots     *x509.CertPool        // Apogee's root, which every station and uplinkd trust
	ca        *x509.Certificate     // Apogee Stations CA 2028
	uplinkd   [2]tls.Certificate    // uplinkd's client certificates: 2001, then 2101 from uplinkdAt
	uplinkdAt time.Time             // when uplinkd changes certificates: 00:30:00
	sats      map[string]*Satellite // the satellites' stand-ins, by ID
}

// Passes returns the passes of the day's slots.
func (d Day) Passes() []Pass {
	passes := make([]Pass, len(d.Slots))
	for i, s := range d.Slots {
		passes[i] = s.Pass
	}
	return passes
}

// Connect makes uplinkd's mutual TLS 1.3 handshake with the station code at the instant at: the station
// presents served, and uplinkd refuses a certificate on list (nil: no list yet). It returns the first
// refusal, uplinkd's or the station's.
func (d Day) Connect(code string, served tls.Certificate, list []byte, at time.Time) error {
	key, ok := served.PrivateKey.(ed25519.PrivateKey)
	if !ok || served.Leaf == nil {
		return fmt.Errorf("%s has no certificate to present", code)
	}
	clock := func() time.Time { return at }
	station := serverConfig(served.Leaf, d.ca, key, clock)
	harden(station)
	requireClients(station, d.roots)
	cert := d.uplinkd[0]
	if !at.Before(d.uplinkdAt) {
		cert = d.uplinkd[1]
	}
	uplinkd := clientConfig(d.roots, strings.ToLower(code)+".apogee.internal", cert, clock)
	withRevocations(uplinkd, func() []byte { return list }, d.ca, clock)
	return handshake(station, uplinkd)
}

// Ping seals a PING for satID under the key version keyID and delivers it to the satellite's stand-in
// at the instant arrive. It returns the satellite's refusal, if any.
func (d Day) Ping(satID, keyID string, arrive time.Time) error {
	sat, ok := d.sats[satID]
	if !ok {
		return fmt.Errorf("no satellite %q", satID)
	}
	line, err := sealCommand(d.Ring, keyID, Command{SatID: satID, Seq: sat.Last + 1, Op: "PING"})
	if err != nil {
		return err
	}
	sat.Now = func() time.Time { return arrive }
	_, err = sat.Accept(line)
	return err
}

// handshake runs one TLS handshake between a server and a client over an in-memory pipe.
func handshake(server, client *tls.Config) error {
	a, b := net.Pipe()
	deadline := time.Now().Add(time.Second)
	a.SetDeadline(deadline)
	b.SetDeadline(deadline)
	done := make(chan error, 1)
	go func() {
		err := tls.Server(b, server).Handshake()
		b.Close()
		done <- err
	}()
	err := tls.Client(a, client).Handshake()
	a.Close()
	if serverErr := <-done; err == nil {
		err = serverErr
	}
	return err
}
