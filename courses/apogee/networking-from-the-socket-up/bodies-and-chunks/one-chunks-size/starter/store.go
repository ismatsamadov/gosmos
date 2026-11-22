package main

import (
	"errors"
	"fmt"
	"hash/crc32"
	"slices"
	"strings"
	"sync"
)

// ErrExists is a recording the store already keeps.
var ErrExists = errors.New("already stored")

// Store keeps finished recordings by pass ID. Safe for many goroutines.
type Store struct {
	mu     sync.Mutex
	data   map[string][]byte
	subs   map[chan string]bool // the feeds' channels
	closed bool
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{data: map[string][]byte{}, subs: map[chan string]bool{}}
}

// Put keeps a copy of data under id and sends its list line to every subscriber. An id that
// is already kept gives an error that matches ErrExists: "furrow-a/YST/18856 already stored".
func (s *Store) Put(id string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[id]; ok {
		return fmt.Errorf("%s %w", id, ErrExists)
	}
	s.data[id] = append([]byte{}, data...)
	line := fmt.Sprintf("%s %d %08x", id, len(data), crc32.ChecksumIEEE(data))
	for ch := range s.subs {
		select {
		case ch <- line:
		default: // 16 lines behind: drop it, so its feed ends
			delete(s.subs, ch)
			close(ch)
		}
	}
	return nil
}

// Get returns the recording kept under id.
func (s *Store) Get(id string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, ok := s.data[id]
	return data, ok
}

// List returns one line per recording, "<pass ID> <bytes> <crc32>", sorted by pass ID.
func (s *Store) List() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	lines := []string{}
	for id, data := range s.data {
		lines = append(lines, fmt.Sprintf("%s %d %08x", id, len(data), crc32.ChecksumIEEE(data)))
	}
	slices.Sort(lines)
	return lines
}

// Subscribe returns a channel with room for 16 lines that receives the list line of every
// recording stored from now on, and a function that unsubscribes. The channel is closed by
// Close, or once the subscriber falls 16 lines behind.
func (s *Store) Subscribe() (<-chan string, func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan string, 16)
	if s.closed {
		close(ch)
		return ch, func() {}
	}
	s.subs[ch] = true
	return ch, func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.subs[ch] {
			delete(s.subs, ch)
			close(ch)
		}
	}
}

// Close closes every subscription, so every feed ends. Recordings can still be stored.
func (s *Store) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	for ch := range s.subs {
		delete(s.subs, ch)
		close(ch)
	}
}

// validPassID reports whether id has a pass ID's shape: a lowercase satellite ID, a slash,
// three capital letters, a slash and five digits, as in furrow-a/YST/18856.
func validPassID(id string) bool {
	sat, rest, ok := strings.Cut(id, "/")
	if !ok || sat == "" || len(rest) != 9 || rest[3] != '/' {
		return false
	}
	for _, r := range sat {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-') {
			return false
		}
	}
	for _, r := range rest[:3] {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	for _, r := range rest[4:] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
