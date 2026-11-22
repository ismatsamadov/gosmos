package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Identity is a station's lasting identity on the backhaul: its code and its Ed25519 key.
type Identity struct {
	Code string             // "PEL"
	Key  ed25519.PrivateKey // signs this station's shares
}

// ErrBadHello is the error every refused HELLO line wraps.
var ErrBadHello = errors.New("bad hello")

// maxLine is the longest line either end reads, without its newline.
const maxLine = 1024

// hello returns the HELLO line that offers share to the other end, signed with id's key.
func hello(id Identity, share *ecdh.PrivateKey) string {
	shareHex := hex.EncodeToString(share.PublicKey().Bytes())
	sig := ed25519.Sign(id.Key, []byte("backhaul share|"+id.Code+"|"+shareHex))
	return "HELLO " + id.Code + " " + shareHex + " " + hex.EncodeToString(sig)
}

// checkHello checks a HELLO line against the identity keys known in advance, and returns its code and share.
func checkHello(line string, known map[string]ed25519.PublicKey) (code string, share *ecdh.PublicKey, err error) {
	fields := strings.Fields(line)
	if len(fields) != 4 || fields[0] != "HELLO" {
		return "", nil, fmt.Errorf("%w: want HELLO <code> <share> <signature>", ErrBadHello)
	}
	code = fields[1]
	pub, ok := known[code]
	if !ok {
		return "", nil, fmt.Errorf("%w: unknown station %q", ErrBadHello, code)
	}
	share, err = parseShare(fields[2])
	if err != nil {
		return "", nil, fmt.Errorf("%w: %v", ErrBadHello, err)
	}
	sig, err := hex.DecodeString(fields[3])
	if err != nil || !ed25519.Verify(pub, []byte("backhaul share|"+code+"|"+fields[2]), sig) {
		return "", nil, fmt.Errorf("%w: signature does not verify", ErrBadHello)
	}
	return code, share, nil
}

// readHello reads the next line from r and checks it as a HELLO.
func readHello(r io.Reader, known map[string]ed25519.PublicKey) (string, *ecdh.PublicKey, error) {
	line, err := readLine(r)
	if err != nil {
		return "", nil, err
	}
	return checkHello(line, known)
}

// writeLine writes line and its newline in one Write.
func writeLine(w io.Writer, line string) error {
	_, err := io.WriteString(w, line+"\n")
	return err
}

// readLine reads one line a byte at a time, so nothing after it is taken, and returns it without its newline.
func readLine(r io.Reader) (string, error) {
	var line []byte
	b := make([]byte, 1)
	for len(line) <= maxLine {
		if _, err := io.ReadFull(r, b); err != nil {
			return "", err
		}
		if b[0] == '\n' {
			return string(line), nil
		}
		line = append(line, b[0])
	}
	return "", fmt.Errorf("line over %d bytes", maxLine)
}

// sendData seals text under the session key and writes it as one DATA line from the station from to the station to.
func sendData(w io.Writer, key []byte, from, to, text string) error {
	aead, err := newAEAD(key)
	if err != nil {
		return err
	}
	nonce := make([]byte, aead.NonceSize())
	rand.Read(nonce)
	sealed := aead.Seal(nil, nonce, []byte(text), []byte("DATA|"+from+"|"+to))
	return writeLine(w, "DATA "+hex.EncodeToString(nonce)+" "+hex.EncodeToString(sealed))
}

// readData reads one DATA line from the station from to the station to, and opens it with the session key.
func readData(r io.Reader, key []byte, from, to string) (string, error) {
	line, err := readLine(r)
	if err != nil {
		return "", err
	}
	fields := strings.Fields(line)
	if len(fields) != 3 || fields[0] != "DATA" || len(fields[1]) != 24 {
		return "", errors.New("malformed DATA line")
	}
	nonce, err1 := hex.DecodeString(fields[1])
	sealed, err2 := hex.DecodeString(fields[2])
	if err1 != nil || err2 != nil {
		return "", errors.New("malformed DATA line")
	}
	aead, err := newAEAD(key)
	if err != nil {
		return "", err
	}
	text, err := aead.Open(nil, nonce, sealed, []byte("DATA|"+from+"|"+to))
	if err != nil {
		return "", errors.New("DATA line failed authentication")
	}
	return string(text), nil
}

// newAEAD returns AES-256-GCM under a 32-byte seal key.
func newAEAD(key []byte) (cipher.AEAD, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("seal key has %d bytes, want 32", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
