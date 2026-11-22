package main

import (
	"encoding/hex"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"
)

// augustKeys is uplink-keys.txt as it stands in August, holding the test master k-2028-07.
const augustKeys = "# uplink-keys.txt: uplinkd's master keys, mode 0600\n" +
	"k-2028-07 2028-08-11T00:00:00Z - 75706c696e6b2074657374206d6173746572206b65792c203332206279746573\n"

// Cirrine-1's own keys from k-2028-07 (lesson 5.4), as the satellite holds them.
const (
	cirrineSignHex = "18668bfe29024937ff84a61df34e8a0310a4702a74834b969f56e2c6ccd00213"
	cirrineSealHex = "427a8ae78a295915acb18a87bdbb40be20d49d27c15dac44e320d4510e070171"
)

// testRing loads augustKeys.
func testRing(t *testing.T) *KeyRing {
	t.Helper()
	ring, err := loadRing(strings.NewReader(augustKeys))
	if err != nil {
		t.Fatalf("loadRing(August's uplink-keys.txt): %v", err)
	}
	return ring
}

// cirrine returns the tests' Cirrine-1 at 07:30:05 on 25 August, its last accepted number last.
func cirrine(t *testing.T, last uint64) *Satellite {
	t.Helper()
	sign, err1 := hex.DecodeString(cirrineSignHex)
	seal, err2 := hex.DecodeString(cirrineSealHex)
	if err1 != nil || err2 != nil {
		t.Fatal("bad test key")
	}
	return &Satellite{
		ID:   "cirrine-1",
		Keys: map[string]SatKeys{"k-2028-07": {Sign: sign, Seal: seal}},
		Last: last,
		Now:  func() time.Time { return time.Date(2028, 8, 25, 7, 30, 5, 0, time.UTC) },
	}
}

func TestFirstSealedCommand(t *testing.T) {
	line, err := sealCommand(testRing(t), "k-2028-07", Command{SatID: "cirrine-1", Seq: 4471, Op: "PING"})
	if err != nil {
		t.Fatalf("sealCommand(cirrine-1 4471 PING) = %v, want nil", err)
	}
	c, err := cirrine(t, 4470).Accept(line)
	if err != nil || c.SatID != "cirrine-1" || c.Seq != 4471 || c.Op != "PING" || len(c.Args) != 0 {
		t.Errorf("Cirrine-1's Accept(your PING line) = %s %d %s %q, %v, want cirrine-1 4471 PING [], nil",
			c.SatID, c.Seq, c.Op, c.Args, err)
	}
}

func TestArguments(t *testing.T) {
	want := Command{SatID: "cirrine-1", Seq: 4472, Op: "SET_MODE", Args: []string{"SAFE"}}
	line, err := sealCommand(testRing(t), "k-2028-07", want)
	if err != nil {
		t.Fatalf("sealCommand(cirrine-1 4472 SET_MODE SAFE) = %v, want nil", err)
	}
	c, err := cirrine(t, 4471).Accept(line)
	if err != nil || c.Seq != want.Seq || c.Op != want.Op || !slices.Equal(c.Args, want.Args) {
		t.Errorf("Cirrine-1's Accept(your SET_MODE line) = %d %s %q, %v, want 4472 SET_MODE [\"SAFE\"], nil",
			c.Seq, c.Op, c.Args, err)
	}
}

// TestNothingReadable reads the lines as anyone listening to the uplink can.
func TestNothingReadable(t *testing.T) {
	ring := testRing(t)
	ping, err := sealCommand(ring, "k-2028-07", Command{SatID: "cirrine-1", Seq: 4471, Op: "PING"})
	if err != nil {
		t.Fatalf("sealCommand(PING) = %v, want nil", err)
	}
	if strings.Contains(ping, "PING") {
		t.Errorf("sealCommand(PING) = %q, want a line without PING", ping)
	}
	dump := Command{SatID: "cirrine-1", Seq: 4472, Op: "DUMP_FILE", Args: []string{"occ-0825.bin"}}
	line, err := sealCommand(ring, "k-2028-07", dump)
	if err != nil {
		t.Fatalf("sealCommand(DUMP_FILE occ-0825.bin) = %v, want nil", err)
	}
	for _, part := range []string{"DUMP_FILE", "occ-0825.bin", hex.EncodeToString([]byte("DUMP_FILE|occ-0825.bin"))} {
		if strings.Contains(line, part) {
			t.Errorf("sealCommand(DUMP_FILE occ-0825.bin) = a line holding %q, want nothing of the command readable", part)
		}
	}
}

func TestFreshEachTime(t *testing.T) {
	ring := testRing(t)
	c := Command{SatID: "cirrine-1", Seq: 4471, Op: "PING"}
	a, errA := sealCommand(ring, "k-2028-07", c)
	b, errB := sealCommand(ring, "k-2028-07", c)
	if errA != nil || errB != nil {
		t.Fatalf("sealCommand(PING) twice = %v and %v, want nil", errA, errB)
	}
	if a == b {
		t.Errorf("sealCommand(PING) twice = %q both times, want two different lines", a)
	}
}

func TestRefuses(t *testing.T) {
	tests := []struct {
		name  string
		keyID string
		c     Command
	}{
		{"a bar inside an argument", "k-2028-07", Command{SatID: "cirrine-1", Seq: 4472, Op: "SET_MODE", Args: []string{"SAFE|REBOOT"}}},
		{"k-2027-03, which the ring doesn't hold", "k-2027-03", Command{SatID: "cirrine-1", Seq: 4472, Op: "PING"}},
		{"Thistle, out of the fleet since 2027", "k-2028-07", Command{SatID: "thistle", Seq: 1, Op: "PING"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if line, err := sealCommand(testRing(t), tt.keyID, tt.c); err == nil || line != "" {
				t.Errorf("sealCommand(%s) = %q, %v, want \"\" and an error", tt.name, line, err)
			}
		})
	}
	_, err := sealCommand(testRing(t), "k-2027-03", Command{SatID: "cirrine-1", Seq: 4472, Op: "PING"})
	if !errors.Is(err, ErrUnknownKey) {
		t.Errorf("sealCommand(under k-2027-03) = %v, want an error wrapping ErrUnknownKey", err)
	}
}
