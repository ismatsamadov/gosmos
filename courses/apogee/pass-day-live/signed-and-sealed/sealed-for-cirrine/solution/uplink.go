package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hkdf"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Master is one uplink master key, a line of uplink-keys.txt.
type Master struct {
	ID       string    // "k-2029-07"
	Activate time.Time // commands are signed and sealed with it from this instant
	Retire   time.Time // nothing made with it is accepted from this instant; zero: not scheduled
	Revoked  bool
	Secret   []byte // 32 bytes; never logged or printed
}

// loadKeys reads uplink-keys.txt in dir: "#" lines and blank lines are skipped, every other
// line is "<id> <activate> <retire or -> <64 hex secret>", optionally followed by "revoked".
func loadKeys(dir string) ([]Master, error) {
	f, err := os.Open(filepath.Join(dir, "uplink-keys.txt"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var ms []Master
	sc := bufio.NewScanner(f)
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		m, err := parseMaster(line)
		if err == nil && slices.ContainsFunc(ms, func(o Master) bool { return o.ID == m.ID }) {
			err = fmt.Errorf("key %s twice", m.ID)
		}
		if err != nil {
			return nil, fmt.Errorf("uplink-keys.txt: line %d: %w", n, err)
		}
		ms = append(ms, m)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("uplink-keys.txt: %w", err)
	}
	slices.SortFunc(ms, func(a, b Master) int { return a.Activate.Compare(b.Activate) })
	return ms, nil
}

// parseMaster reads one line of uplink-keys.txt.
func parseMaster(line string) (Master, error) {
	f := strings.Fields(line)
	if len(f) < 4 || len(f) > 5 || (len(f) == 5 && f[4] != "revoked") {
		return Master{}, errors.New("want <id> <activate> <retire or -> <64 hex secret> [revoked]")
	}
	m := Master{ID: f[0], Revoked: len(f) == 5}
	var err error
	if m.Activate, err = time.Parse(time.RFC3339, f[1]); err != nil {
		return Master{}, fmt.Errorf("activate %q is not a time like 2029-07-25T00:00:00Z", f[1])
	}
	if f[2] != "-" {
		if m.Retire, err = time.Parse(time.RFC3339, f[2]); err != nil {
			return Master{}, fmt.Errorf("retire %q is not a time like 2029-08-01T00:00:00Z", f[2])
		}
	}
	if m.Secret, err = hex.DecodeString(f[3]); err != nil || len(m.Secret) != 32 {
		return Master{}, errors.New("the secret is not 64 hex digits")
	}
	return m, nil
}

// usable reports whether a master may still be loaded or used at an instant: not revoked
// and not retired.
func (m Master) usable(at time.Time) bool {
	return !m.Revoked && (m.Retire.IsZero() || at.Before(m.Retire))
}

// satKey derives one of a satellite's keys from a master: purpose is "sign" or "seal".
func satKey(m Master, purpose, satID string) []byte {
	k, err := hkdf.Key(sha256.New, m.Secret, []byte("apogee uplink keys v1"), purpose+"|"+satID, 32)
	if err != nil {
		panic(err) // HKDF-SHA256 always gives 32 bytes
	}
	return k
}

// Uplink is uplinkd: the masters, and every satellite's numbering and keys.
type Uplink struct {
	r       *run
	masters []Master // by Activate
	sats    map[string]*satUplink
}

// satUplink is one satellite as uplinkd knows it. Only the pass that holds lock may read or
// change the fields below it.
type satUplink struct {
	id     string
	sealed bool            // Cirrine's contract: every command to it is sealed
	lock   chan struct{}   // one slot: the pass that fills it has the satellite's commands
	next   uint64          // one more than the last number sent to it
	holds  map[string]bool // the masters whose keys it holds, by ID
}

// sentCommand is a command on its way up.
type sentCommand struct {
	seq uint64
}

// newUplink returns uplinkd for a window: at its start every satellite holds every master
// activated before it and neither retired nor revoked then.
func newUplink(r *run, masters []Master) *Uplink {
	u := &Uplink{r: r, masters: masters, sats: map[string]*satUplink{}}
	for _, s := range r.n.Satellites {
		su := &satUplink{id: s.ID, sealed: s.Operator == "Cirrine Atmospherics", lock: make(chan struct{}, 1), holds: map[string]bool{}}
		for _, m := range masters {
			if m.Activate.Before(r.wk.Window.From) && m.usable(r.wk.Window.From) {
				su.holds[m.ID] = true
			}
		}
		u.sats[s.ID] = su
	}
	return u
}

// uplinkRun is one pass's commands.
type uplinkRun struct {
	u        *Uplink
	sat      *satUplink
	due      bool // AOS + 10 s has come
	holding  bool // the pass holds the satellite's lock
	built    bool // its commands have been worked out
	queue    []command
	inflight *sentCommand
	last     uint64 // the satellite's last accepted number, from the latest WELCOME
	refused  bool
}

// command is one command a pass has still to send.
type command struct {
	op   string
	args []string
}

// forPass returns a pass's commands.
func (u *Uplink) forPass(p *pass) *uplinkRun {
	return &uplinkRun{u: u, sat: u.sats[p.c.SatID]}
}

// welcome takes the satellite's last accepted number from a WELCOME.
func (ur *uplinkRun) welcome(last uint64) {
	ur.last = last
}

// wants reports whether the pass needs its satellite's lock now.
func (ur *uplinkRun) wants() bool {
	return ur.due && !ur.holding && !ur.refused && (!ur.built || len(ur.queue) > 0)
}

// hold takes the satellite's lock, and works out the pass's commands the first time.
func (ur *uplinkRun) hold() {
	ur.holding = true
	if !ur.built {
		ur.built = true
		ur.queue = ur.commands()
	}
}

// commands works out what goes up on this pass: PING.
func (ur *uplinkRun) commands() []command {
	return []command{{op: "PING"}}
}

// next returns the frame to send now, if any: the next command, once the one before it has
// been answered. With nothing left it lets go of the satellite's lock.
func (ur *uplinkRun) next() (Frame, bool) {
	if !ur.holding || ur.refused || ur.inflight != nil {
		return Frame{}, false
	}
	if len(ur.queue) == 0 {
		ur.release()
		return Frame{}, false
	}
	cmd := ur.queue[0]
	ur.queue = ur.queue[1:]
	seq := max(ur.sat.next, ur.last+1)
	line, err := ur.u.line(ur.sat, cmd, seq, ur.u.r.w.Now())
	if err != nil {
		ur.refused = true
		ur.release()
		return Frame{}, false
	}
	ur.sat.next = seq + 1
	ur.inflight = &sentCommand{seq: seq}
	return Frame{Type: TypeUplink, Seq: uint32(seq), Body: []byte(line)}, true
}

// answer takes the satellite's ACK or NAK for the command in flight.
func (ur *uplinkRun) answer(f Frame) {
	c := ur.inflight
	if c == nil || f.Seq != uint32(c.seq) {
		return
	}
	ur.inflight = nil
	if f.Type == TypeNak {
		ur.refused = true
		ur.release()
	}
}

// release lets go of the satellite's lock.
func (ur *uplinkRun) release() {
	if ur.holding {
		ur.holding = false
		<-ur.sat.lock
	}
}

// finish ends the pass's commands at its LOS: it lets go of the satellite's lock.
func (ur *uplinkRun) finish() {
	ur.release()
}

// unanswered reports whether a command was still unanswered at LOS.
func (ur *uplinkRun) unanswered() bool { return ur.inflight != nil }

// line writes a command as its satellite reads it, under the newest master that has
// activated, is neither retired nor revoked and that the satellite holds: sealed for a
// Cirrine satellite, signed otherwise.
func (u *Uplink) line(sat *satUplink, cmd command, seq uint64, now time.Time) (string, error) {
	var m Master
	found := false
	for _, k := range u.masters {
		if sat.holds[k.ID] && !k.Activate.After(now) && k.usable(now) {
			m, found = k, true
		}
	}
	if !found {
		return "", fmt.Errorf("no key to sign with at %s", now.Format(time.RFC3339))
	}
	if sat.sealed {
		return seal(satKey(m, "seal", sat.id), sat.id, seq, m.ID, cmd)
	}
	return sign(satKey(m, "sign", sat.id), sat.id, seq, cmd), nil
}

// sign writes a CMD line: the HMAC-SHA256 under the satellite's sign key of the canonical
// form, "<sat>|<seq>|<OP>|<args...>".
func sign(key []byte, satID string, seq uint64, cmd command) string {
	fields := append([]string{satID, strconv.FormatUint(seq, 10), cmd.op}, cmd.args...)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(strings.Join(fields, "|")))
	return "CMD " + strings.Join(fields, " ") + " sig=" + hex.EncodeToString(mac.Sum(nil))
}

// seal writes a SEAL line: AES-256-GCM under the satellite's seal key, a fresh nonce, the
// header as associated data and "<OP>|<args...>" inside.
func seal(key []byte, satID string, seq uint64, keyID string, cmd command) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	head := fmt.Sprintf("2 %s %d %s", satID, seq, keyID)
	ad := "SEAL|" + strings.ReplaceAll(head, " ", "|")
	plain := strings.Join(append([]string{cmd.op}, cmd.args...), "|")
	sealed := gcm.Seal(nil, nonce, []byte(plain), []byte(ad))
	return "SEAL " + head + " " + hex.EncodeToString(nonce) + " " + hex.EncodeToString(sealed), nil
}
