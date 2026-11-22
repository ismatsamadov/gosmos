package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
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
	c := &Client{Timeout: 5 * time.Second, conn: conn}
	c.br = bufio.NewReader(idleReader{c})
	name, err := readBanner(c.br)
	if err != nil {
		return nil, err
	}
	c.Name = name
	return c, nil
}

// idleReader reads the client's connection, giving every read c.Timeout to bring a byte.
type idleReader struct{ c *Client }

// Read reads from the client's connection, giving the read c.Timeout from now.
func (r idleReader) Read(p []byte) (int, error) {
	r.c.conn.SetReadDeadline(time.Now().Add(r.c.Timeout))
	return r.c.conn.Read(p)
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

// send writes one frame to the controller.
func (c *Client) send(f Frame) error {
	c.conn.SetWriteDeadline(time.Now().Add(c.Timeout))
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

// hello says HELLO as dishd and returns the controller's WELCOME.
func (c *Client) hello() (Welcome, error) {
	if err := c.send(Frame{Type: TypeHello, Body: []byte("dishd")}); err != nil {
		return Welcome{}, err
	}
	ans, err := c.recv()
	switch {
	case err != nil:
		return Welcome{}, err
	case ans.Type == TypeNak:
		return Welcome{}, nakError(ans)
	case ans.Type != TypeWelcome:
		return Welcome{}, fmt.Errorf("HELLO answered with type 0x%02x", ans.Type)
	}
	return decodeWelcome(ans.Body)
}

// status asks the controller for its status frame with QUERY.
func (c *Client) status() (Status, error) {
	if err := c.send(Frame{Type: TypeQuery}); err != nil {
		return Status{}, err
	}
	ans, err := c.recv()
	switch {
	case err != nil:
		return Status{}, err
	case ans.Type == TypeNak:
		return Status{}, nakError(ans)
	case ans.Type != TypeStatus:
		return Status{}, fmt.Errorf("QUERY answered with type 0x%02x", ans.Type)
	}
	return decodeStatus(ans.Body)
}

// Close hangs up.
func (c *Client) Close() error {
	return c.conn.Close()
}

// describe turns an error from the controller called name into one line for Ingrid's log;
// timeout is how long dishd gave the line to answer.
func describe(name string, timeout time.Duration, err error) string {
	switch {
	case errors.Is(err, os.ErrDeadlineExceeded):
		return name + ": no answer in " + timeout.String()
	case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF):
		return name + ": closed"
	}
	return name + ": " + err.Error()
}

// heartbeat sends PING on every tick and requires a PONG, until ticks closes. It returns the
// first error: a dead line, or an answer that isn't a PONG.
func heartbeat(c *Client, ticks <-chan time.Time) error {
	for range ticks {
		if err := c.send(Frame{Type: TypePing}); err != nil {
			return err
		}
		f, err := c.recv()
		if err != nil {
			return err
		}
		if f.Type != TypePong {
			return fmt.Errorf("heartbeat answered with type 0x%02x, want PONG", f.Type)
		}
	}
	return nil
}
