package main

import (
	"os"
	"strings"
	"testing"
)

func TestWriteBoard(t *testing.T) {
	tests := []struct {
		ids  []string
		want string
	}{
		{[]string{"pipit-1/PEL/04182", "furrow-a/SQT/11873", "thistle/LAN/02218"},
			"NEXT PASSES\npipit-1/PEL/04182\nfurrow-a/SQT/11873\nthistle/LAN/02218\n"},
		{[]string{"cirrine-1/YST/07455"}, "NEXT PASSES\ncirrine-1/YST/07455\n"},
		{nil, "NEXT PASSES\n"},
	}
	for _, tt := range tests {
		var b strings.Builder
		writeBoard(&b, tt.ids)
		if got := b.String(); got != tt.want {
			t.Errorf("writeBoard(&b, %q) wrote %q, want %q", tt.ids, got, tt.want)
		}
	}
}

// TestWriteBoardOnlyToW catches a board drawn on the terminal as well as into w.
func TestWriteBoardOnlyToW(t *testing.T) {
	ids := []string{"thistle/LAN/02218"}
	terminal, err := os.CreateTemp(t.TempDir(), "terminal")
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = terminal
	var b strings.Builder
	writeBoard(&b, ids)
	os.Stdout = stdout
	terminal.Close()
	printed, err := os.ReadFile(terminal.Name())
	if err != nil {
		t.Fatal(err)
	}
	if len(printed) > 0 {
		t.Errorf("writeBoard(&b, %q) printed %q to the terminal, want nothing there", ids, printed)
	}
}
