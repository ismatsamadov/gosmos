package main

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"
)

func TestReadBanner(t *testing.T) {
	tests := []struct {
		name, stream, want string
	}{
		{"the controller's own line", "YST-1 READY PROTO 1\r\n", "YST-1"},
		{"page one's line, LF alone", "YST-1 READY PROTO 1\n", "YST-1"},
		{"Dish 2's controller", "YST-2 READY PROTO 1\r\n", "YST-2"},
		{"busy, not ready", "YST-1 BUSY\r\n", ""},
		{"a space before the line ending", "YST-1 READY PROTO 1 \r\n", ""},
		{"two carriage returns", "YST-1 READY PROTO 1\r\r\n", ""},
		{"a name in small letters", "yst-1 READY PROTO 1\r\n", ""},
		{"a stream that ends mid-line", "YST-1 READY PRO", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readBanner(bufio.NewReader(strings.NewReader(tt.stream)))
			if tt.want == "" && (got != "" || err == nil) {
				t.Errorf("readBanner(%q) = %q, %v, want \"\" and an error", tt.stream, got, err)
			}
			if tt.want != "" && (got != tt.want || err != nil) {
				t.Errorf("readBanner(%q) = %q, %v, want %q, nil", tt.stream, got, err, tt.want)
			}
		})
	}
	t.Run("protocol 2", func(t *testing.T) {
		stream := "YST-2 READY PROTO 2\r\n"
		got, err := readBanner(bufio.NewReader(strings.NewReader(stream)))
		if got != "" || err == nil || err.Error() != `unsupported protocol "2"` {
			t.Errorf("readBanner(%q) = %q, %v, want \"\" and unsupported protocol \"2\"", stream, got, err)
		}
	})
	t.Run("the frame after the banner stays in the reader", func(t *testing.T) {
		br := bufio.NewReader(strings.NewReader("YST-1 READY PROTO 1\r\n\x00\x09"))
		if got, err := readBanner(br); got != "YST-1" || err != nil {
			t.Fatalf("readBanner(banner and a frame) = %q, %v, want \"YST-1\", nil", got, err)
		}
		if next, err := br.ReadByte(); next != 0x00 || err != nil {
			t.Errorf("after readBanner, the next byte is %#02x, %v; want 0x00, the frame's first byte", next, err)
		}
	})
	t.Run("the simulator's banner", func(t *testing.T) {
		conn, end := net.Pipe()
		served := make(chan struct{})
		go func() {
			NewSim("YST-1", 1).ServeConn(end)
			close(served)
		}()
		defer func() {
			conn.Close()
			<-served
		}()
		conn.SetDeadline(time.Now().Add(time.Second))
		if got, err := readBanner(bufio.NewReader(conn)); got != "YST-1" || err != nil {
			t.Errorf("readBanner(simulator) = %q, %v, want \"YST-1\", nil", got, err)
		}
	})
}
