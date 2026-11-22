package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// Client is dishd's connection to one dish controller.
type Client struct {
	Name    string        // from the banner: "YST-1"
	Timeout time.Duration // each read or write must make progress within it: 5s
	conn    net.Conn
	br      *bufio.Reader // every read goes through it, banner included, for the connection's life
}

// newClient reads the controller's banner and returns a client for the connection.
func newClient(conn net.Conn) (*Client, error) {
	br := bufio.NewReader(conn)
	name, err := readBanner(br)
	if err != nil {
		return nil, err
	}
	return &Client{Name: name, Timeout: 5 * time.Second, conn: conn, br: br}, nil
}

// exchange sends every frame, then reads one answer per frame, in order. It returns the
// answers read and the first error.
func (c *Client) exchange(frames ...Frame) ([]Frame, error) {
	return nil, nil
}

// send writes one frame to the controller.
func (c *Client) send(f Frame) error {
	return writeFrame(c.conn, f)
}

// recv reads the controller's next frame; a BYE becomes a *ByeError.
func (c *Client) recv() (Frame, error) {
	f, err := readFrame(c.br)
	if err != nil {
		return Frame{}, err
	}
	if f.Type == TypeBye {
		return Frame{}, &ByeError{Reason: string(f.Body)}
	}
	return f, nil
}

// Close hangs up.
func (c *Client) Close() error {
	return c.conn.Close()
}

// readBanner reads the controller's banner line and returns the controller's name.
func readBanner(br *bufio.Reader) (string, error) {
	line, err := br.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	parts := strings.Split(line, " ")
	if len(parts) != 4 || parts[1] != "READY" || parts[2] != "PROTO" || !validName(parts[0]) {
		return "", fmt.Errorf("not a banner: %q", line)
	}
	if parts[3] != "1" {
		return "", fmt.Errorf("unsupported protocol %q", parts[3])
	}
	return parts[0], nil
}

// validName reports whether name is 1 to 16 capital letters, digits and hyphens.
func validName(name string) bool {
	return len(name) >= 1 && len(name) <= 16 && strings.Trim(name, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-") == ""
}

// readPayload reads one frame and returns the bytes after its length: type, seq, body
// and CRC.
func readPayload(br *bufio.Reader) ([]byte, error) {
	var length [2]byte
	if _, err := io.ReadFull(br, length[:]); err != nil {
		return nil, err
	}
	n := int(binary.BigEndian.Uint16(length[:]))
	if n < 9 || n > 1024 {
		return nil, fmt.Errorf("%w %d", ErrBadLength, n)
	}
	payload := make([]byte, n)
	if _, err := io.ReadFull(br, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// readFrame reads one frame from br.
func readFrame(br *bufio.Reader) (Frame, error) {
	payload, err := readPayload(br)
	if err != nil {
		return Frame{}, err
	}
	return decodeFrame(payload)
}

// writeFrame sends f to w in a single Write.
func writeFrame(w io.Writer, f Frame) error {
	_, err := w.Write(encodeFrame(f))
	return err
}
