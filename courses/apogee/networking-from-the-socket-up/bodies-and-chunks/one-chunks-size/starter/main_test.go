package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestChunkSize(t *testing.T) {
	tests := []struct {
		name, line string
		want       int64
		bad        bool
	}{
		{"a 12-byte piece", "c", 12, false},
		{"a 4,096-byte piece", "1000", 4096, false},
		{"capital hex digits", "1A", 26, false},
		{"small hex digits", "1a", 26, false},
		{"the last chunk", "0", 0, false},
		{"an extension", "400;name=value", 1024, false},
		{"a space before the extension", "400 ;x", 1024, false},
		{"a tab after the digits", "1000\t", 4096, false},
		{"not hex", "zz", 0, true},
		{"an empty line", "", 0, true},
		{"a minus sign", "-1", 0, true},
		{"a plus sign", "+5", 0, true},
		{"a 0x prefix", "0x10", 0, true},
		{"a space before the digits", " c", 0, true},
		{"17 digits", strings.Repeat("0", 13) + "1000", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chunkSize(tt.line)
			if !tt.bad {
				if got != tt.want || err != nil {
					t.Errorf("chunkSize(%q) = %d, %v, want %d, nil", tt.line, got, err, tt.want)
				}
				return
			}
			want := fmt.Sprintf("bad chunk size %q", tt.line)
			if e, ok := errors.AsType[*RequestError](err); !ok || e.Status != 400 || e.Msg != want {
				t.Errorf("chunkSize(%q) = %d, %v, want 400 %s", tt.line, got, err, want)
			}
		})
	}
}
