package main

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// cirrineCommands is Cirrine's command file for Thursday's 09:00 pass over Ystmark.
const cirrineCommands = "PING\nSET_MODE NOMINAL\nDUMP_FILE occ-0203-0900\n"

// form returns a multipart body: a "note" part, then a part named name holding data, and the body's Content-Type.
func form(t *testing.T, name string, data []byte) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if err := w.SetBoundary("b7"); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteField("note", "Thursday 09:00 over Ystmark"); err != nil {
		t.Fatal(err)
	}
	part, err := w.CreateFormFile(name, "commands.txt")
	if err != nil {
		t.Fatal(err)
	}
	part.Write(data)
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return &buf, w.FormDataContentType()
}

// upload calls readCommands with a POST of body as contentType, the way the booking service will.
func upload(body io.Reader, contentType string) ([]byte, error) {
	r := httptest.NewRequest(http.MethodPost, "/bookings/bk-4918/commands", body)
	r.Header.Set("Content-Type", contentType)
	return readCommands(httptest.NewRecorder(), r)
}

func TestReadCommands(t *testing.T) {
	body, ctype := form(t, "file", []byte(cirrineCommands))
	got, err := upload(body, ctype)
	if err != nil || string(got) != cirrineCommands {
		t.Errorf("readCommands(Cirrine's form) = %q, %v, want %q, nil: the bytes of the part named file and nothing else", got, err, cirrineCommands)
	}

	long := bytes.Repeat([]byte("PING\n"), 13000) // 65,000 bytes, a form just under the limit
	body, ctype = form(t, "file", long)
	if got, err := upload(body, ctype); err != nil || !bytes.Equal(got, long) {
		t.Errorf("readCommands(a 65,000-byte command file) = %d bytes, %v, want all 65000 and nil", len(got), err)
	}

	got, err = upload(strings.NewReader(cirrineCommands), "text/plain")
	if !errors.Is(err, http.ErrNotMultipart) {
		t.Errorf("readCommands(a text/plain body) = %q, %v, want the error MultipartReader gives, http.ErrNotMultipart", got, err)
	}

	var noFile bytes.Buffer
	w := multipart.NewWriter(&noFile)
	w.SetBoundary("b7")
	w.WriteField("note", "the file is coming")
	w.Close()
	got, err = upload(&noFile, w.FormDataContentType())
	if !errors.Is(err, http.ErrMissingFile) {
		t.Errorf("readCommands(a form with only a note) = %q, %v, want http.ErrMissingFile", got, err)
	}
}

// image is the 1.9 GB software image Cirrine's engineer attached, made as it is read so the test holds none of it.
// It counts the bytes read and stops with an error after 1 MiB, so no reading of it can run for long.
type image struct {
	head, tail []byte
	left       int64 // bytes of the image still to send
	read       int64
}

// errFarEnough ends a reading of the image that went far past the limit.
var errFarEnough = errors.New("test: read more than 1 MiB of the image")

// Read sends the form's head, the image's bytes, then its tail.
func (im *image) Read(p []byte) (int, error) {
	if im.read > 1<<20 {
		return 0, errFarEnough
	}
	n := 0
	switch {
	case len(im.head) > 0:
		n = copy(p, im.head)
		im.head = im.head[n:]
	case im.left > 0:
		n = min(len(p), int(min(im.left, 1<<30)))
		for i := range n {
			p[i] = byte(i)
		}
		im.left -= int64(n)
	case len(im.tail) > 0:
		n = copy(p, im.tail)
		im.tail = im.tail[n:]
	default:
		return 0, io.EOF
	}
	im.read += int64(n)
	return n, nil
}

func TestASoftwareImageIsStopped(t *testing.T) {
	var head bytes.Buffer
	w := multipart.NewWriter(&head)
	if _, err := w.CreateFormFile("file", "cirrine-1-software.img"); err != nil {
		t.Fatal(err)
	}
	headLen := head.Len()
	w.Close()
	im := &image{head: head.Bytes()[:headLen], tail: head.Bytes()[headLen:], left: 1_900_000_000}
	got, err := upload(im, w.FormDataContentType())
	if _, ok := errors.AsType[*http.MaxBytesError](err); !ok {
		t.Errorf("readCommands(a 1.9 GB form) = %d bytes, %v, want a *http.MaxBytesError", len(got), err)
	}
	if im.read > maxBodyBytes+1 {
		t.Errorf("readCommands(a 1.9 GB form) read %d bytes of it, want at most %d: the body's limit must stop the reading", im.read, maxBodyBytes+1)
	}
}
