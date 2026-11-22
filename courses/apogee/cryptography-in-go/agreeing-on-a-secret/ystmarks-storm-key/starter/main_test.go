package main

import (
	"bytes"
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/sha256"
	"errors"
	"net"
	"strings"
	"testing"
	"time"
)

// commandFile is this week's command file for Ystmark: three command lines for Pipit-2, each
// signed with Pipit-2's own key.
var commandFile = []string{
	"CMD pipit-2 1207 PING sig=129e7e9fc239d5f81b6f92a6e630dee35d16ceebdf81071ef41aec477be90475",
	"CMD pipit-2 1208 SET_MODE NOMINAL sig=a8170443cfcbb665afb8814991c83bbdc0b05353c29b6a3dc54928e3f7a3f6cf",
	"CMD pipit-2 1209 DUMP_FILE tlm-0907.bin sig=3d3ec78b6db6436b813a73957bb7d37276b73bdbfeb6771674377b3e58250638",
}

// seed returns the SHA-256 of label: every test key here is made from one.
func seed(label string) []byte {
	sum := sha256.Sum256([]byte(label))
	return sum[:]
}

// identity returns the test identity of the station code, made from the seed "identity <code>".
func identity(code string) Identity {
	return Identity{Code: code, Key: ed25519.NewKeyFromSeed(seed("identity " + code))}
}

// knows returns what a station holds in advance: the public identity key of each code.
func knows(codes ...string) map[string]ed25519.PublicKey {
	known := make(map[string]ed25519.PublicKey)
	for _, code := range codes {
		known[code] = identity(code).Key.Public().(ed25519.PublicKey)
	}
	return known
}

// brief returns line with each field longer than 50 characters cut to 8 and an ellipsis.
func brief(line string) string {
	fields := strings.Split(line, " ")
	for i, f := range fields {
		if len(f) > 50 {
			fields[i] = f[:8] + "…"
		}
	}
	return strings.Join(fields, " ")
}

// pipe returns the two ends of an in-memory connection, each with a deadline one second away.
func pipe() (net.Conn, net.Conn) {
	a, b := net.Pipe()
	deadline := time.Now().Add(time.Second)
	a.SetDeadline(deadline)
	b.SetDeadline(deadline)
	return a, b
}

// dialed is what dial returned, sent from its own goroutine, and the error of any DATA line
// sent after it.
type dialed struct {
	key     []byte
	err     error
	sendErr error
}

// answered is what answer returned, sent from its own goroutine.
type answered struct {
	peer string
	key  []byte
	err  error
}

// TestStormKey runs the whole exchange over one connection: both ends end with the same key,
// and Ystmark reads the command file's three lines. The pipe's deadlines end it within 1s.
func TestStormKey(t *testing.T) {
	pelEnd, ystEnd := pipe()
	pel := make(chan dialed, 1)
	go func() {
		defer pelEnd.Close()
		key, err := dial(pelEnd, identity("PEL"), "YST", knows("YST"))
		var sendErr error
		if err == nil && len(key) == 32 {
			for _, line := range commandFile {
				if sendErr = sendData(pelEnd, key, "PEL", "YST", line); sendErr != nil {
					break
				}
			}
		}
		pel <- dialed{key, err, sendErr}
	}()
	peer, key, err := answer(ystEnd, identity("YST"), knows("PEL"))
	var got []string
	var readErr error
	if err == nil && len(key) == 32 {
		for range commandFile {
			line, err := readData(ystEnd, key, "PEL", "YST")
			if err != nil {
				readErr = err
				break
			}
			got = append(got, line)
		}
	}
	ystEnd.Close()
	d := <-pel

	if d.err != nil || len(d.key) != 32 {
		t.Fatalf("dial(PEL to YST) = %d bytes, %v, want a 32-byte session key, nil", len(d.key), d.err)
	}
	if err != nil || peer != "PEL" || len(key) != 32 {
		t.Fatalf("answer(YST) = %q, %d bytes, %v, want \"PEL\", a 32-byte session key, nil", peer, len(key), err)
	}
	if !bytes.Equal(d.key, key) {
		t.Fatalf("dial's key = %x…, answer's key = %x…, want one session key at both ends", d.key[:4], key[:4])
	}
	if readErr != nil || strings.Join(got, "\n") != strings.Join(commandFile, "\n") {
		t.Errorf("readData at Ystmark = %d lines, %v, want the command file's 3 lines", len(got), readErr)
	}
}

// TestMiddle puts a machine in the middle of the backhaul. It keeps each end's HELLO and passes
// on one of its own instead: its own share, under the other end's name, signed with the only
// identity key it has, "identity MID". Both ends must refuse it, before any DATA line.
func TestMiddle(t *testing.T) {
	pelEnd, midPEL := pipe() // Port Ellery to the middle
	midYST, ystEnd := pipe() // the middle to Ystmark
	pel := make(chan dialed, 1)
	go func() {
		defer pelEnd.Close()
		key, err := dial(pelEnd, identity("PEL"), "YST", knows("YST"))
		var sendErr error
		if err == nil && len(key) == 32 {
			sendErr = sendData(pelEnd, key, "PEL", "YST", commandFile[0])
		}
		pel <- dialed{key, err, sendErr}
	}()
	yst := make(chan answered, 1)
	go func() {
		defer ystEnd.Close()
		peer, key, err := answer(ystEnd, identity("YST"), knows("PEL"))
		yst <- answered{peer, key, err}
	}()

	midShare, err := ecdh.X25519().NewPrivateKey(seed("ephemeral MID"))
	if err != nil {
		t.Fatalf("ecdh.X25519().NewPrivateKey(SHA-256 of %q) = %v", "ephemeral MID", err)
	}
	midKey := identity("MID").Key
	var leaked string
	if _, err := readLine(midPEL); err == nil { // PEL's HELLO, kept back
		writeLine(midYST, hello(Identity{Code: "PEL", Key: midKey}, midShare))
		readLine(midYST) // YST's HELLO, if Ystmark answers at all
		if writeLine(midPEL, hello(Identity{Code: "YST", Key: midKey}, midShare)) == nil {
			leaked, _ = readLine(midPEL) // anything PEL sends next
		}
	}
	midPEL.Close()
	midYST.Close()
	d, a := <-pel, <-yst

	if !errors.Is(a.err, ErrBadHello) {
		t.Errorf("answer(YST, sent the middle's HELLO as PEL) = %q, %d bytes, %v, want ErrBadHello", a.peer, len(a.key), a.err)
	}
	if !errors.Is(d.err, ErrBadHello) {
		t.Errorf("dial(PEL to YST, answered by the middle as YST) = %d bytes, %v, want ErrBadHello", len(d.key), d.err)
	}
	if leaked != "" {
		t.Errorf("after the middle's HELLO, PEL sent %q, want no DATA line", brief(leaked))
	}
}

// TestWrongStation dials SQT, but Ystmark answers: a real, signed HELLO from the wrong station.
func TestWrongStation(t *testing.T) {
	pelEnd, ystEnd := pipe()
	yst := make(chan answered, 1)
	go func() {
		defer ystEnd.Close()
		peer, key, err := answer(ystEnd, identity("YST"), knows("PEL"))
		yst <- answered{peer, key, err}
	}()
	key, err := dial(pelEnd, identity("PEL"), "SQT", knows("YST"))
	pelEnd.Close()
	<-yst
	if !errors.Is(err, ErrBadHello) {
		t.Errorf("dial(PEL to SQT, answered by YST) = %d bytes, %v, want ErrBadHello", len(key), err)
	}
}

// handshake runs dial and answer over a fresh connection and returns both ends' keys, stopping
// the test unless they are one 32-byte key.
func handshake(t *testing.T) (pelKey, ystKey []byte) {
	t.Helper()
	pelEnd, ystEnd := pipe()
	pel := make(chan dialed, 1)
	go func() {
		defer pelEnd.Close()
		key, err := dial(pelEnd, identity("PEL"), "YST", knows("YST"))
		pel <- dialed{key, err, nil}
	}()
	_, ystKey, err := answer(ystEnd, identity("YST"), knows("PEL"))
	ystEnd.Close()
	d := <-pel
	if d.err != nil || err != nil || len(d.key) != 32 || !bytes.Equal(d.key, ystKey) {
		t.Fatalf("dial(PEL to YST) = %d bytes, %v; answer(YST) = %d bytes, %v; want one 32-byte key at both ends", len(d.key), d.err, len(ystKey), err)
	}
	return d.key, ystKey
}

// TestTampered changes one hex digit of a DATA line on its way to Ystmark: readData must refuse it.
func TestTampered(t *testing.T) {
	pelKey, ystKey := handshake(t)
	var buf bytes.Buffer
	if err := sendData(&buf, pelKey, "PEL", "YST", commandFile[0]); err != nil {
		t.Fatalf("sendData(PEL's key, line 1) = %v", err)
	}
	line := buf.Bytes()
	if got, err := readData(bytes.NewReader(line), ystKey, "PEL", "YST"); err != nil || got != commandFile[0] {
		t.Fatalf("readData(line 1 as sent) = %q, %v, want %q", brief(got), err, brief(commandFile[0]))
	}
	i := len(line) - 2 // the ciphertext's last digit, before the newline
	if line[i] == '0' {
		line[i] = '1'
	} else {
		line[i] = '0'
	}
	if got, err := readData(bytes.NewReader(line), ystKey, "PEL", "YST"); err == nil {
		t.Errorf("readData(line 1 with one digit changed) = %q, nil, want an error", brief(got))
	}
}
