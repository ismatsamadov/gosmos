package main

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCheckCommands(t *testing.T) {
	cases := []struct {
		name string
		data string
		kind string // what the sniffer sees, "" when the file must be kept
	}{
		{"Furrowlight's command list", "PING\nSET_MODE NOMINAL\nDUMP_FILE field-12\n", ""},
		{"a list that starts with a byte-order mark", "\xef\xbb\xbfPING\n", ""},
		{"the list saved as a PDF", "%PDF-1.7\n%\xe2\xe3\xcf\xd3\n1 0 obj\n", "application/pdf"},
		{"the list as a web page", "<!DOCTYPE html><html><body>PING</body></html>", "text/html; charset=utf-8"},
		{"the list as XML", "<?xml version=\"1.0\"?><commands/>", "text/xml; charset=utf-8"},
		{"the list as a spreadsheet", "PK\x03\x04\x14\x00\x06\x00", "application/zip"},
		{"a photo of the list", "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR", "image/png"},
		{"the list in UTF-16", "\xff\xfeP\x00I\x00N\x00G\x00", "text/plain; charset=utf-16le"},
		{"a file of raw bytes", "PING\x00\x01\x02SET_MODE", "application/octet-stream"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := checkCommands([]byte(c.data))
			if c.kind == "" {
				if err != nil {
					t.Errorf("checkCommands(%q) = %v, want nil", c.data, err)
				}
				return
			}
			want := "command file must be plain text, not " + c.kind
			apiErr, ok := errors.AsType[*APIError](err)
			if !ok || apiErr.Status != http.StatusUnsupportedMediaType || apiErr.Code != "unsupported_media_type" || apiErr.Message != want {
				t.Errorf("checkCommands(%q) = %#v, want a 415 unsupported_media_type %q", c.data, err, want)
			}
		})
	}
}

// bookings is a store holding two bookings, enough for the upload handler.
type bookings map[string]Booking

// Create is not used here.
func (s bookings) Create(ctx context.Context, b Booking) (Booking, error) {
	return Booking{}, errors.New("not in this test")
}

// Get returns the booking with this ID, or ErrNotFound.
func (s bookings) Get(ctx context.Context, id string) (Booking, error) {
	b, ok := s[id]
	if !ok {
		return Booking{}, ErrNotFound
	}
	return b, nil
}

// List is not used here.
func (s bookings) List(ctx context.Context) ([]Booking, error) {
	return nil, errors.New("not in this test")
}

// Delete is not used here.
func (s bookings) Delete(ctx context.Context, id string) error { return errors.New("not in this test") }

// post sends a command form for booking id, its file part labeled partType, as Furrowlight.
func post(t *testing.T, f *Files, id, partType string, data []byte) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	h := make(map[string][]string)
	h["Content-Disposition"] = []string{`form-data; name="file"; filename="commands.txt"`}
	h["Content-Type"] = []string{partType}
	part, err := w.CreatePart(h)
	if err != nil {
		t.Fatal(err)
	}
	part.Write(data)
	w.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /bookings/{id}/commands", f.upload)
	r := httptest.NewRequest(http.MethodPost, "/bookings/"+id+"/commands", &body)
	r.Header.Set("Content-Type", w.FormDataContentType())
	furrowlight := Customer{ID: "furrowlight", Name: "Furrowlight", SatIDs: []string{"furrow-a", "furrow-b"}}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, r.WithContext(withCustomer(r.Context(), furrowlight)))
	return rec
}

func TestUploadLooksAtTheBytes(t *testing.T) {
	f := &Files{
		Store: bookings{
			"bk-4921": {ID: "bk-4921", SatID: "furrow-a", Station: "SQT", AOS: time.Date(2028, 2, 4, 14, 36, 38, 0, time.UTC), LOS: time.Date(2028, 2, 4, 14, 46, 32, 0, time.UTC)},
		},
		Commands: t.TempDir(),
	}
	kept := filepath.Join(f.Commands, "bk-4921.txt")

	rec := post(t, f, "bk-4921", "text/plain", []byte("%PDF-1.7\n%\xe2\xe3\xcf\xd3\n"))
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Errorf("POST /bookings/bk-4921/commands with a PDF labeled text/plain = %d %s, want 415", rec.Code, rec.Body)
	}
	if _, err := os.Stat(kept); err == nil {
		t.Errorf("bk-4921.txt was kept after the PDF, want nothing kept")
	}

	commands := []byte("PING\nSET_MODE NOMINAL\nDUMP_FILE field-12\n")
	rec = post(t, f, "bk-4921", "application/octet-stream", commands)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("POST /bookings/bk-4921/commands with text labeled application/octet-stream = %d %s, want 204: the bytes decide, not the label", rec.Code, rec.Body)
	}
	if got, err := os.ReadFile(kept); err != nil || !bytes.Equal(got, commands) {
		t.Errorf("bk-4921.txt = %q, %v, want %q", got, err, commands)
	}
}
