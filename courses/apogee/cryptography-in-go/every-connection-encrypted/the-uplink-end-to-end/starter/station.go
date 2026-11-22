package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Queue holds a station's command lines until the next pass. It is safe for many requests at once.
type Queue struct {
	mu    sync.Mutex
	lines []string
}

// Add queues one line.
func (q *Queue) Add(line string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.lines = append(q.lines, line)
}

// Lines returns the queued lines, oldest first.
func (q *Queue) Lines() []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	return append([]string(nil), q.lines...)
}

// caller returns the name on the client's verified certificate, and whether there is one.
func caller(r *http.Request) (string, bool) {
	if r.TLS == nil || len(r.TLS.VerifiedChains) == 0 || len(r.TLS.VerifiedChains[0]) == 0 {
		return "", false
	}
	return r.TLS.VerifiedChains[0][0].Subject.CommonName, true
}

// stationHandler answers uplinkd at the station with this code: its status, and its queue.
func stationHandler(code string, q *Queue) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "ok %s\n", code)
	})
	mux.HandleFunc("POST /uplink", func(w http.ResponseWriter, r *http.Request) {
		if name, ok := caller(r); !ok || name != "uplinkd" {
			http.Error(w, "only uplinkd may queue commands", http.StatusForbidden)
			return
		}
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1024))
		line := strings.TrimSuffix(string(body), "\n")
		sat, seq, lineErr := lineHeader(line)
		if err != nil || lineErr != nil {
			http.Error(w, "bad command line", http.StatusBadRequest)
			return
		}
		q.Add(line)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "queued %s %d\n", sat, seq)
	})
	return mux
}

// lineHeader checks the syntax of a CMD or SEAL line and returns its satellite and sequence number.
func lineHeader(line string) (string, uint64, error) {
	if strings.HasPrefix(line, "SEAL ") {
		s, err := parseSealed(line)
		return s.SatID, s.Seq, err
	}
	f := strings.Split(line, " ")
	if len(f) < 5 || f[0] != "CMD" || !strings.HasPrefix(f[len(f)-1], "sig=") || len(f[len(f)-1]) != 68 {
		return "", 0, errors.New("bad command line")
	}
	seq, err := strconv.ParseUint(f[2], 10, 64)
	if err != nil || seq == 0 {
		return "", 0, errors.New("bad command line")
	}
	return f[1], seq, nil
}

// send posts one command line to the station's uplinkd at url, and fails unless the station queued it.
func send(ctx context.Context, c *http.Client, url, line string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/uplink", strings.NewReader(line))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("%s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}
