package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	want := "YST Ystmark\nPEL Port Ellery\nSQT Salar Quieto\nLAN Lantern Atoll\nCWR Cape Wren\n"
	var stdout, stderr bytes.Buffer
	status := run(nil, strings.NewReader(""), &stdout, &stderr)
	if status != 0 {
		t.Errorf("run(no arguments) = %d, want 0", status)
	}
	if stdout.String() != want {
		t.Errorf("run(no arguments) stdout = %q, want %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Errorf("run(no arguments) stderr = %q, want nothing", stderr.String())
	}
}
