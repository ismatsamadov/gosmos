package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"
)

// at returns the instant hh:mm:ss UTC on Rotation Day, 17 October 2028.
func at(hms string) time.Time {
	t, err := time.Parse(time.RFC3339, "2028-10-17T"+hms+"Z")
	if err != nil {
		panic(err)
	}
	return t
}

// The two versions in uplinkd's ring on Rotation Day, with the season's test secrets.
var (
	k07 = KeyVersion{
		Key:      Key{ID: "k-2028-07", Secret: []byte("uplink test master key, 32 bytes")},
		Activate: time.Date(2028, 8, 11, 0, 0, 0, 0, time.UTC),
		Retire:   time.Date(2028, 10, 24, 0, 0, 0, 0, time.UTC),
	}
	k10 = KeyVersion{
		Key:      Key{ID: "k-2028-10", Secret: []byte("uplink next master key, 32 bytes")},
		Activate: at("00:00:00"),
	}
)

// The days the stations CA issued the old certificates and Rotation Day's.
var (
	sept22 = time.Date(2028, 9, 22, 0, 0, 0, 0, time.UTC)
	oct9   = time.Date(2028, 10, 9, 0, 0, 0, 0, time.UTC)
)

// testKey returns the test key whose seed is the SHA-256 of label, as every test key this season is made.
func testKey(label string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte(label))
	return ed25519.NewKeyFromSeed(seed[:])
}

// slots returns Rotation Day's ten passes in the order they reached the log, each with its PINGs'
// seal times: one at AOS, and two for the Cirrine-1 pass across midnight.
func slots() []Slot {
	pass := func(sat, station, aos, los string) Slot {
		return Slot{Pass: Pass{SatID: sat, Station: station, AOS: at(aos), LOS: at(los)}, SealAt: []time.Time{at(aos)}}
	}
	midnight := Slot{
		Pass:   Pass{SatID: "cirrine-1", Station: "YST", AOS: time.Date(2028, 10, 16, 23, 52, 10, 0, time.UTC), LOS: at("00:03:40")},
		SealAt: []time.Time{time.Date(2028, 10, 16, 23, 59, 58, 0, time.UTC), at("00:01:00")},
	}
	return []Slot{
		pass("furrow-a", "SQT", "14:02:10", "14:11:52"),
		pass("pipit-1", "PEL", "04:12:05", "04:22:41"),
		pass("brine-3", "LAN", "09:14:20", "09:22:31"),
		midnight,
		pass("cirrine-1", "YST", "07:30:00", "07:41:12"),
		pass("pipit-3", "PEL", "00:58:30", "01:07:10"),
		pass("furrow-a", "CWR", "04:58:30", "05:08:02"),
		pass("brine-3", "LAN", "01:02:00", "01:10:40"),
		pass("pipit-2", "YST", "04:31:10", "04:41:52"),
		pass("furrow-b", "CWR", "02:51:20", "03:00:44"),
	}
}

// satellites returns the fleet's stand-ins on Rotation Day: each holds its own keys from k-2028-07
// and, since the week of 9 October, from k-2028-10, and the number of the last command it accepted.
func satellites(t *testing.T) map[string]*Satellite {
	t.Helper()
	last := map[string]uint64{"pipit-1": 830, "pipit-2": 1290, "pipit-3": 260, "furrow-a": 2150, "furrow-b": 2085, "brine-3": 3310, "cirrine-1": 4520}
	sats := map[string]*Satellite{}
	for id, n := range last {
		keys := map[string]SatKeys{}
		for _, v := range []KeyVersion{k07, k10} {
			sign, err := satKey(v.Key, "sign", id)
			if err != nil {
				t.Fatalf("satKey: %v", err)
			}
			seal, err := satKey(v.Key, "seal", id)
			if err != nil {
				t.Fatalf("satKey: %v", err)
			}
			keys[v.ID] = SatKeys{Sign: sign, Seal: seal, Retire: v.Retire}
		}
		sats[id] = &Satellite{ID: id, Keys: keys, Last: n, Now: func() time.Time { return at("00:00:00") }}
	}
	return sats
}

// rotationDay builds Rotation Day as it happened. With rehearsal, it builds the day Ingrid rehearsed
// instead: Lantern Atoll never switches, and list 1 comes at 05:00, naming 1004 as well as 1005.
func rotationDay(t *testing.T, rehearsal bool) Day {
	t.Helper()
	caDay := time.Date(2028, 9, 18, 0, 0, 0, 0, time.UTC)
	rootKey, interKey := testKey("ca root 2028"), testKey("ca stations 2028")
	root, err := newRoot(rootKey, caDay)
	if err != nil {
		t.Fatalf("newRoot: %v", err)
	}
	inter, err := newIntermediate(Issuer{Cert: root, Key: rootKey}, interKey, caDay)
	if err != nil {
		t.Fatalf("newIntermediate: %v", err)
	}
	ca := Issuer{Cert: inter, Key: interKey}
	issue := func(mint func(Issuer, string, int64, ed25519.PublicKey, time.Time) (*x509.Certificate, error), name string, serial int64, day time.Time) tls.Certificate {
		key := testKey(fmt.Sprintf("cert %s %d", name, serial))
		leaf, err := mint(ca, name, serial, key.Public().(ed25519.PublicKey), day)
		if err != nil {
			t.Fatalf("issuing %s %d: %v", name, serial, err)
		}
		return tls.Certificate{Certificate: [][]byte{leaf.Raw, inter.Raw}, PrivateKey: key, Leaf: leaf}
	}
	d := Day{
		Slots:      slots(),
		Ring:       &KeyRing{versions: []KeyVersion{k07, k10}},
		Old:        map[string]tls.Certificate{},
		New:        map[string]tls.Certificate{},
		SwitchFrom: at("01:00:00"),
		ListAt:     at("10:12:00"),
		roots:      x509.NewCertPool(),
		ca:         inter,
		uplinkd:    [2]tls.Certificate{issue(issueClient, "uplinkd", 2001, sept22), issue(issueClient, "uplinkd", 2101, oct9)},
		uplinkdAt:  at("00:30:00"),
		sats:       satellites(t),
	}
	d.roots.AddCert(root)
	for i, code := range []string{"YST", "PEL", "SQT", "LAN", "CWR"} {
		host := strings.ToLower(code) + ".apogee.internal"
		d.Old[code] = issue(issueServer, host, int64(1001+i), sept22)
		d.New[code] = issue(issueServer, host, int64(1101+i), oct9)
	}
	revoked := []int64{1005}
	if rehearsal {
		d.New["LAN"] = d.Old["LAN"]
		d.ListAt = at("05:00:00")
		revoked = []int64{1004, 1005}
	}
	var entries []x509.RevocationListEntry
	for _, serial := range revoked {
		entries = append(entries, x509.RevocationListEntry{SerialNumber: big.NewInt(serial), RevocationTime: d.ListAt})
	}
	if d.List, err = revoke(ca, 1, entries, d.ListAt); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	return d
}

// check compares a day's report with the one Ingrid kept.
func check(t *testing.T, day string, r Report, flown, missed int, lines []string) {
	t.Helper()
	if r.Flown != flown || r.Missed != missed {
		t.Errorf("runDay(%s) = %d flown, %d missed, want %d flown, %d missed", day, r.Flown, r.Missed, flown, missed)
	}
	if len(r.Lines) != len(lines) {
		t.Errorf("runDay(%s) = %d lines, want %d", day, len(r.Lines), len(lines))
	}
	for i := 0; i < len(lines) && i < len(r.Lines); i++ {
		if r.Lines[i] != lines[i] {
			t.Errorf("runDay(%s).Lines[%d] = %q, want %q", day, i, r.Lines[i], lines[i])
		}
	}
}

func TestRotationDay(t *testing.T) {
	d := rotationDay(t, false)
	start := time.Now()
	r := runDay(d)
	if took := time.Since(start); took > time.Second {
		t.Errorf("runDay(Rotation Day) took %v, want well under a second", took.Round(time.Millisecond))
	}
	check(t, "Rotation Day", r, 10, 0, []string{
		"23:52:10 cirrine-1 YST flown k-2028-07+k-2028-10 1001",
		"00:58:30 pipit-3 PEL flown k-2028-10 1002",
		"01:02:00 brine-3 LAN flown k-2028-10 1004",
		"02:51:20 furrow-b CWR flown k-2028-10 1105",
		"04:12:05 pipit-1 PEL flown k-2028-10 1102",
		"04:31:10 pipit-2 YST flown k-2028-10 1101",
		"04:58:30 furrow-a CWR flown k-2028-10 1105",
		"07:30:00 cirrine-1 YST flown k-2028-10 1101",
		"09:14:20 brine-3 LAN flown k-2028-10 1104",
		"14:02:10 furrow-a SQT flown k-2028-10 1103",
	})
	if got := d.sats["cirrine-1"].Now(); !got.Equal(at("07:30:05")) {
		t.Errorf("after runDay(Rotation Day), Cirrine-1's last PING arrived at %s, want 07:30:05: each PING reaches the satellite 5 seconds after it is sealed", got.Format("15:04:05"))
	}
}

// TestRehearsal runs the day as Ingrid rehearsed it: Lantern Atoll never switches, and list 1, at
// 05:00, names its old certificate too.
func TestRehearsal(t *testing.T) {
	d := rotationDay(t, true)
	start := time.Now()
	r := runDay(d)
	if took := time.Since(start); took > time.Second {
		t.Errorf("runDay(Ingrid's rehearsal) took %v, want well under a second", took.Round(time.Millisecond))
	}
	check(t, "Ingrid's rehearsal", r, 9, 1, []string{
		"23:52:10 cirrine-1 YST flown k-2028-07+k-2028-10 1001",
		"00:58:30 pipit-3 PEL flown k-2028-10 1002",
		"01:02:00 brine-3 LAN flown k-2028-10 1004",
		"02:51:20 furrow-b CWR flown k-2028-10 1105",
		"04:12:05 pipit-1 PEL flown k-2028-10 1102",
		"04:31:10 pipit-2 YST flown k-2028-10 1101",
		"04:58:30 furrow-a CWR flown k-2028-10 1105",
		"07:30:00 cirrine-1 YST flown k-2028-10 1101",
		"09:14:20 brine-3 LAN missed: certificate 1004 is revoked",
		"14:02:10 furrow-a SQT flown k-2028-10 1103",
	})
}
