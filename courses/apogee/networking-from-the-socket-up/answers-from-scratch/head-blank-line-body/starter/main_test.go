package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// stored is the station server's 201 for Furrow-A's pass, as the collector reads it.
const stored = "HTTP/1.1 201 Created\r\n" +
	"Date: Tue, 20 Jun 2028 06:40:36 GMT\r\n" +
	"Content-Type: text/plain; charset=utf-8\r\n" +
	"Content-Length: 37\r\n" +
	"Location: /recordings/furrow-a/YST/18856\r\n" +
	"Apogee-Crc32: a2912082\r\n" +
	"\r\n" +
	"stored furrow-a/YST/18856 4096 bytes\n"

func TestWriteResponse(t *testing.T) {
	fields := []Field{
		{"Content-Type", "text/plain; charset=utf-8"},
		{"Content-Length", "37"},
		{"Location", "/recordings/furrow-a/YST/18856"},
		{"Apogee-Crc32", "a2912082"},
	}
	body := []byte("stored furrow-a/YST/18856 4096 bytes\n")
	at := time.Date(2028, 6, 20, 6, 40, 36, 0, time.UTC)

	t.Run("the 201 for Furrow-A's pass", func(t *testing.T) {
		var buf bytes.Buffer
		err := writeResponse(&buf, at, 201, fields, body)
		if got := buf.String(); got != stored || err != nil {
			t.Errorf("writeResponse(201) wrote\n%q, %v\nwant\n%q, nil", got, err, stored)
		}
	})
	t.Run("read back by net/http", func(t *testing.T) {
		var buf bytes.Buffer
		writeResponse(&buf, at, 201, fields, body)
		res, err := http.ReadResponse(bufio.NewReader(&buf), nil)
		if err != nil {
			t.Fatalf("http.ReadResponse(what writeResponse wrote) = %v, want an answer", err)
		}
		got, err := io.ReadAll(res.Body)
		if res.StatusCode != 201 || string(got) != string(body) || err != nil {
			t.Errorf("http.ReadResponse(what writeResponse wrote) = %d %q, %v, want 201 %q, nil", res.StatusCode, got, err, body)
		}
	})
	t.Run("a clock two hours ahead of UTC", func(t *testing.T) {
		var buf bytes.Buffer
		writeResponse(&buf, at.In(time.FixedZone("UTC+2", 2*60*60)), 201, fields, body)
		date := "Date: Tue, 20 Jun 2028 06:40:36 GMT\r\n"
		if got := buf.String(); !strings.Contains(got, date) {
			t.Errorf("writeResponse(201, 08:40:36 at UTC+2) wrote\n%q\nwant it to hold %q", got, date)
		}
	})
	t.Run("no body", func(t *testing.T) {
		var buf bytes.Buffer
		writeResponse(&buf, at, 200, []Field{{"Content-Length", "4096"}}, nil)
		want := "HTTP/1.1 200 OK\r\nDate: Tue, 20 Jun 2028 06:40:36 GMT\r\nContent-Length: 4096\r\n\r\n"
		if got := buf.String(); got != want {
			t.Errorf("writeResponse(200, no body) wrote\n%q\nwant\n%q", got, want)
		}
	})
	t.Run("a writer that fails", func(t *testing.T) {
		gone := errors.New("collector gone")
		if err := writeResponse(failingWriter{gone}, at, 201, fields, body); !errors.Is(err, gone) {
			t.Errorf("writeResponse(a writer that fails) = %v, want %v", err, gone)
		}
	})
	t.Run("one write", func(t *testing.T) {
		var w countingWriter
		writeResponse(&w, at, 201, fields, body)
		if w.writes != 1 {
			t.Errorf("writeResponse(201) made %d writes, want 1", w.writes)
		}
	})
}

// failingWriter is a connection whose other end has gone: every write fails with err.
type failingWriter struct{ err error }

// Write fails.
func (w failingWriter) Write(p []byte) (int, error) { return 0, w.err }

// countingWriter counts the writes made on it.
type countingWriter struct{ writes int }

// Write counts the write and takes every byte.
func (w *countingWriter) Write(p []byte) (int, error) {
	w.writes++
	return len(p), nil
}
