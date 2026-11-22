package main

import (
	"errors"
	"io"
	"slices"
	"strings"
	"testing"
)

const (
	note = "# night shift 2027-07-06, Paz"
	lg   = "2027-07-06T03:14:07Z\tLG\t183.25\t41.70"
	ch   = "2027-07-06T03:14:08Z  CH  183.31   41.66"
)

func TestCollectLines(t *testing.T) {
	want := []string{note, lg, "", ch}
	got, err := collectLines(strings.NewReader(note + "\n" + lg + "\n\n" + ch + "\n"))
	if !slices.Equal(got, want) || err != nil {
		t.Errorf("collectLines(whole log) = %q, %v, want %q, nil", got, err, want)
	}
	got, err = collectLines(strings.NewReader(""))
	if len(got) != 0 || err != nil {
		t.Errorf("collectLines(empty log) = %q, %v, want [], nil", got, err)
	}
	problem := errors.New("stick unplugged")
	damaged := io.MultiReader(strings.NewReader(lg+"\n"), readerFunc(func([]byte) (int, error) { return 0, problem }))
	got, err = collectLines(damaged)
	if !slices.Equal(got, []string{lg}) || !errors.Is(err, problem) {
		t.Errorf("collectLines(damaged copy) = %q, %v, want %q, %v", got, err, []string{lg}, problem)
	}
}

// readerFunc adapts a test function to io.Reader.
type readerFunc func([]byte) (int, error)

// Read calls the supplied test function.
func (f readerFunc) Read(p []byte) (int, error) { return f(p) }
