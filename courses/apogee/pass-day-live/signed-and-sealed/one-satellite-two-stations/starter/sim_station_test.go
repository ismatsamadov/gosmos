package main

// The simulated world's stations and satellites. A station's radio answers each call on a
// net.Pipe: its banner, then frames, as the pass line's protocol says; behind it a
// satellite checks every command by Season 5's rules and sends its telemetry and files
// down. The trouble a simCase asks for happens here.

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hkdf"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"hash/fnv"
	"io"
	"net"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Frame types and NAK codes of the pass line, as the stations speak them.
const (
	simTypeHello   byte = 0x01
	simTypeTLM     byte = 0x20
	simTypeUplink  byte = 0x30
	simTypeFile    byte = 0x40
	simTypeChunk   byte = 0x41
	simTypeEnd     byte = 0x42
	simTypeWelcome byte = 0x81
	simTypeAck     byte = 0x90
	simTypeNak     byte = 0x91
	simTypeBye     byte = 0xFF

	simNakUnknownPass byte = 1
	simNakStowed      byte = 2
	simNakBadCRC      byte = 3
	simNakHelloFirst  byte = 4
	simNakUnknownType byte = 5
	simNakRefused     byte = 6

	simChunkBytes = 256 // every chunk of every file
	simSpinLimit  = 500 // dials at one station in one instant before the world slows them down
)

// simEncode returns a frame's bytes: the length, the type, the seq, the body and the
// CRC-32 of the type, seq and body, every number big-endian.
func simEncode(typ byte, seq uint32, body []byte) []byte {
	b := make([]byte, 2+1+4+len(body)+4)
	binary.BigEndian.PutUint16(b[0:], uint16(9+len(body)))
	b[2] = typ
	binary.BigEndian.PutUint32(b[3:], seq)
	copy(b[7:], body)
	binary.BigEndian.PutUint32(b[7+len(body):], crc32.ChecksumIEEE(b[2:7+len(body)]))
	return b
}

// simFrame is one frame a station read from the rehearsal.
type simFrame struct {
	Type  byte
	Seq   uint32
	Body  []byte
	CRCOK bool
}

// simReadFrame reads one frame. A length under 9 or over 1024 is an error.
func simReadFrame(br *bufio.Reader) (simFrame, error) {
	var head [2]byte
	if _, err := io.ReadFull(br, head[:]); err != nil {
		return simFrame{}, err
	}
	n := int(binary.BigEndian.Uint16(head[:]))
	if n < 9 || n > 1024 {
		return simFrame{}, fmt.Errorf("bad length %d", n)
	}
	p := make([]byte, n)
	if _, err := io.ReadFull(br, p); err != nil {
		return simFrame{}, err
	}
	sum := binary.BigEndian.Uint32(p[n-4:])
	return simFrame{Type: p[0], Seq: binary.BigEndian.Uint32(p[1:5]), Body: p[5 : n-4], CRCOK: sum == crc32.ChecksumIEEE(p[:n-4])}, nil
}

// simHash returns a small number made from a name, the same on every run.
func simHash(name string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(name))
	return h.Sum32()
}

// simTelemetry returns the housekeeping line a satellite makes at an instant. Its frame
// number is the instant's Unix second. About one Thistle-2 line in seven reads batt_v=nan.
func simTelemetry(sat string, at time.Time) string {
	n := at.Unix()
	h := int64(simHash(sat))
	b := make([]byte, 0, 80)
	b = append(b, sat...)
	b = append(b, ' ')
	b = strconv.AppendInt(b, n, 10)
	b = append(b, ' ')
	b = at.UTC().AppendFormat(b, "2006-01-02T15:04:05Z")
	b = append(b, " batt_v="...)
	if sat == "thistle-2" && n%7 == 0 {
		b = append(b, "nan"...)
	} else {
		b = strconv.AppendFloat(b, 7.60+float64((n/60+h)%60)/100, 'f', 2, 64)
	}
	b = append(b, " temp_c="...)
	b = strconv.AppendFloat(b, -5+float64((n/60+h/7)%200)/10, 'f', 1, 64)
	b = append(b, " mode=NOMINAL"...)
	return string(b)
}

// simWorld is the simulated world of one run: every station's radio, every satellite and
// every customer's receiver, and the records tests read once the run is over.
type simWorld struct {
	c     simCase
	win   Window
	codes []string
	now   func() time.Time
	done  chan struct{} // closed when the run is over: everything in the world ends
	once  sync.Once
	wg    sync.WaitGroup
	hooks *simHooks

	planOnce sync.Once
	plan     *simPlan
	booked   map[string]Booking // the case's bookings that match a contact, by pass ID
	masters  []simMaster        // the case's uplink masters

	mu       sync.Mutex
	calls    []*simCall          // every call, in dial order
	passes   map[string]*simPass // the passes a station has welcomed, by pass ID
	sats     map[string]*simSat  // the satellites that have been commanded or welcomed
	files    map[string]*simFile // the files satellites send down, by name
	relayed  map[string]int      // TLM frames each station has relayed
	dialAt   map[string]time.Time
	dialN    map[string]int
	stopTime time.Time // when the case ended the rehearsal's context; zero: it didn't
	closed   bool      // the run is over: no call is answered any more

	// What reached the rehearsal, for tests to compare its reports with.
	newest   map[string]string // each satellite's highest-numbered good TLM line
	tlmGood  int               // TLM frames that arrived whole with a good CRC
	tlmBad   int               // TLM frames that arrived with a flipped byte
	tlmNaN   int               // of the good ones, lines that read batt_v=nan
	badSizes int               // frames with length 8 that arrived
	chunks   []simChunk        // every chunk that arrived
}

// simChunk is one CHUNK frame that reached the rehearsal.
type simChunk struct {
	File   string
	N      int
	At     time.Time
	PassID string
}

// simNewWorld builds the world for a run. now is the bubble's clock, shifted to the window.
func simNewWorld(r *simResult, now func() time.Time) *simWorld {
	w := &simWorld{c: r.c, win: r.win, codes: r.codes, now: now, done: make(chan struct{}),
		passes: map[string]*simPass{}, sats: map[string]*simSat{}, files: map[string]*simFile{},
		relayed: map[string]int{}, dialAt: map[string]time.Time{}, dialN: map[string]int{}, newest: map[string]string{}}
	w.hooks = simNewHooks(w)
	keys := r.c.keys
	if keys == "" {
		keys = simKeys
	}
	w.masters = simParseKeys(keys)
	return w
}

// thePlan returns the reference plan the stations fly, with the case's bookings matched to
// its contacts.
func (w *simWorld) thePlan() *simPlan {
	w.planOnce.Do(func() {
		w.plan = simPlanFor(w.win, w.codes)
		w.booked = map[string]Booking{}
		type key struct {
			sat, station string
			aos          int64
		}
		byKey := map[key]Booking{}
		for _, b := range w.c.bookings {
			byKey[key{b.SatID, b.Station, b.AOS.Unix()}] = b
		}
		for _, c := range w.plan.contacts {
			if b, ok := byKey[key{c.SatID, c.Station, c.AOS.Unix()}]; ok {
				w.booked[c.ID] = b
			}
		}
	})
	return w.plan
}

// close ends the world: every call is hung up and every goroutine returns.
func (w *simWorld) close() {
	w.once.Do(func() { close(w.done) })
	w.mu.Lock()
	w.closed = true
	calls := slices.Clone(w.calls)
	w.mu.Unlock()
	for _, cl := range calls {
		if cl.conn != nil {
			cl.conn.Close()
		}
	}
	w.wg.Wait()
}

// stopped notes the instant the case ended the rehearsal's context.
func (w *simWorld) stopped(at time.Time) {
	w.mu.Lock()
	w.stopTime = at
	w.mu.Unlock()
}

// sleepUntil waits for an instant of the world's clock; it reports false if the world
// ended first.
func (w *simWorld) sleepUntil(at time.Time) bool {
	d := at.Sub(w.now())
	if d <= 0 {
		return true
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return true
	case <-w.done:
		return false
	}
}

// eager returns the contact a glued station welcomes the moment it answers a call: the
// one contact in progress at the station then, if there is exactly one.
func (w *simWorld) eager(cl *simCall) (Contact, bool) {
	if !slices.Contains(w.c.glue, cl.station) {
		return Contact{}, false
	}
	var found []Contact
	for _, c := range w.thePlan().contacts {
		if c.Station == cl.station && !cl.dialed.Before(c.AOS) && cl.dialed.Before(c.LOS) {
			found = append(found, c)
		}
	}
	if len(found) != 1 {
		return Contact{}, false
	}
	return found[0], true
}

// quietAt reports whether a station's line is dead at an instant.
func (w *simWorld) quietAt(station string, at time.Time) bool {
	return slices.ContainsFunc(w.c.quiet, func(q simSpan) bool { return q.has(station, at) })
}

// diedSince reports whether a station's line died after a call was dialed and by now:
// such a call stays silent for good.
func (w *simWorld) diedSince(station string, dialed, now time.Time) bool {
	for _, q := range w.c.quiet {
		if q.Station == station && dialed.Before(q.From) && !now.Before(q.From) {
			return true
		}
	}
	return false
}

// slowAt returns a station's slow line at an instant, if it has one then.
func (w *simWorld) slowAt(station string, at time.Time) (simSlow, bool) {
	for _, s := range w.c.backhaul {
		if s.Station == station && !at.Before(s.From) && at.Before(s.To) {
			return s, true
		}
	}
	return simSlow{}, false
}

// stowedAt reports whether a station's dishes are stowed at an instant.
func (w *simWorld) stowedAt(station string, at time.Time) bool {
	for _, s := range w.c.stows {
		if s.Station == station && !at.Before(s.From) && (s.To.IsZero() || at.Before(s.To)) {
			return true
		}
	}
	return false
}

// nextStow returns the first instant after at when a station's dishes stow, if any.
func (w *simWorld) nextStow(station string, at time.Time) (time.Time, bool) {
	var next time.Time
	for _, s := range w.c.stows {
		if s.Station == station && s.From.After(at) && (next.IsZero() || s.From.Before(next)) {
			next = s.From
		}
	}
	return next, !next.IsZero()
}

// simCall is one call to a station's radio, as the station saw it. Tests read its fields
// once the run is over.
type simCall struct {
	w    *simWorld
	conn net.Conn      // the station's end of the pipe; nil when the dial failed
	wake chan struct{} // tells the writer there is something to write
	gone chan struct{} // closed when the call is over
	end  sync.Once

	station string
	dialed  time.Time
	failed  bool // the dial failed: no call was made
	late    bool // dialed after the rehearsal's context had ended
	pass    *simPass

	mu          sync.Mutex
	hello       string    // the body of its HELLO, "thistle-2/PEL/02136 1"; "" without one
	passID      string    // the pass its HELLO named
	dish        int       // the dish its HELLO named
	welcomed    bool      // the station welcomed it
	nak         byte      // the code the station refused its HELLO with; 0: it didn't
	received    []byte    // every byte the rehearsal wrote on it
	hungUp      time.Time // when the rehearsal closed it; zero: it never did
	closed      time.Time // when the station closed it; zero: it never did
	dropped     bool      // the station closed it without a BYE, as a broken line does
	deadlines   int       // read deadlines the rehearsal set on it
	farDeadline bool      // one of them was more than a day away
	queue       []simOut
	held        []byte // a banner waiting to go out in one write with the next message
	eager       string // the HELLO a glued station has already welcomed, before it came
}

// simOut is one message a station writes on a call.
type simOut struct {
	data      []byte
	notBefore time.Time     // not written before this instant
	piece     int           // written in pieces of at most this many bytes; 0: in one
	gap       time.Duration // between the pieces
	arrived   func()        // called once the rehearsal has all of it
	hangUp    bool          // the station hangs up after it
}

// simConn is the rehearsal's end of a call: it notes when the rehearsal hangs up and the
// read deadlines it sets.
type simConn struct {
	net.Conn
	call *simCall
	once sync.Once
}

// Close hangs up, and notes the instant the first time.
func (c *simConn) Close() error {
	c.once.Do(func() {
		cl := c.call
		cl.mu.Lock()
		cl.hungUp = cl.w.now()
		cl.mu.Unlock()
	})
	return c.Conn.Close()
}

// SetReadDeadline notes the deadline and sets it.
func (c *simConn) SetReadDeadline(t time.Time) error {
	c.call.noteDeadline(t)
	return c.Conn.SetReadDeadline(t)
}

// SetDeadline notes the deadline and sets it.
func (c *simConn) SetDeadline(t time.Time) error {
	c.call.noteDeadline(t)
	return c.Conn.SetDeadline(t)
}

// noteDeadline counts a read deadline; one more than a day away was made from the wrong
// clock.
func (cl *simCall) noteDeadline(t time.Time) {
	if t.IsZero() {
		return
	}
	cl.mu.Lock()
	cl.deadlines++
	if time.Until(t) > 24*time.Hour {
		cl.farDeadline = true
	}
	cl.mu.Unlock()
}

// dial is World.Dial: it calls a station's radio.
func (w *simWorld) dial(ctx context.Context, station string) (net.Conn, error) {
	now := w.now()
	cl := &simCall{w: w, station: station, dialed: now, wake: make(chan struct{}, 1), gone: make(chan struct{})}
	w.mu.Lock()
	w.calls = append(w.calls, cl)
	if !w.dialAt[station].Equal(now) {
		w.dialAt[station], w.dialN[station] = now, 0
	}
	w.dialN[station]++
	spinning := w.dialN[station] > simSpinLimit
	w.mu.Unlock()
	if spinning {
		// Hundreds of calls to one station in one instant: a loop that dials without
		// waiting. Let time pass, or the week would never end.
		w.sleepUntil(now.Add(time.Second))
	}
	fail := func(err error) (net.Conn, error) {
		cl.failed = true
		return nil, &net.OpError{Op: "dial", Net: "pass line", Err: err}
	}
	if err := ctx.Err(); err != nil {
		cl.late = true
		return fail(err)
	}
	if _, ok := simSite(station); !ok {
		return fail(fmt.Errorf("no station %q", station))
	}
	if slices.ContainsFunc(w.c.dialFail, func(s simSpan) bool { return s.has(station, now) }) || w.quietAt(station, now) {
		return fail(errors.New("no route to " + station))
	}
	mine, theirs := net.Pipe()
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		mine.Close()
		theirs.Close()
		return fail(errors.New("the simulated world has ended"))
	}
	cl.conn = theirs
	w.wg.Add(2)
	w.mu.Unlock()
	go cl.serve()
	go cl.write()
	return &simConn{Conn: mine, call: cl}, nil
}

// simTap notes every byte the station reads from the rehearsal.
type simTap struct{ cl *simCall }

// Read reads from the call and keeps what it read.
func (t simTap) Read(p []byte) (int, error) {
	n, err := t.cl.conn.Read(p)
	if n > 0 {
		t.cl.mu.Lock()
		t.cl.received = append(t.cl.received, p[:n]...)
		t.cl.mu.Unlock()
	}
	return n, err
}

// serve is a station answering one call: the banner, then every frame the rehearsal sends,
// until one of them hangs up.
func (cl *simCall) serve() {
	defer cl.w.wg.Done()
	defer cl.over()
	w, c := cl.w, cl.w.c
	if slices.Contains(c.silent, cl.station) {
		io.Copy(io.Discard, simTap{cl})
		return
	}
	banner := cl.station + " READY PROTO 1\r\n"
	if b, ok := c.banner[cl.station]; ok {
		banner = b
		if !strings.HasSuffix(b, "\n") {
			banner += "\r\n"
		}
	}
	if ct, ok := w.eager(cl); ok && !c.bannerOnly {
		// A glued station doesn't wait to be asked: it writes its banner and the WELCOME
		// for the contact in progress at its dish in one write.
		cl.mu.Lock()
		cl.held = []byte(banner)
		cl.mu.Unlock()
		cl.greet(fmt.Sprintf("%s %d", ct.ID, ct.Dish), w.now())
		cl.mu.Lock()
		cl.eager, cl.hello, cl.passID, cl.dish = cl.hello, "", "", 0
		cl.mu.Unlock()
	} else {
		cl.send([]byte(banner), 0, nil, false)
	}
	if c.bannerOnly {
		io.Copy(io.Discard, simTap{cl})
		return
	}
	br := bufio.NewReaderSize(simTap{cl}, 1100) // a frame is at most 1026 bytes
	for {
		f, err := simReadFrame(br)
		if err != nil {
			return
		}
		now := w.now()
		if w.diedSince(cl.station, cl.dialed, now) {
			continue // a dead line carries nothing either way
		}
		var extra time.Duration // a slow line: the frame arrives this much later
		if slow, ok := w.slowAt(cl.station, now); ok {
			extra = slow.Delay
		}
		switch {
		case !f.CRCOK:
			cl.refuse(f.Seq, simNakBadCRC, "bad crc", extra)
		case f.Type == simTypeHello && cl.eager != "":
			cl.mu.Lock()
			body := string(f.Body)
			cl.hello = body
			cl.passID, _, _ = strings.Cut(body, " ")
			matched := body == cl.eager
			cl.eager = ""
			cl.mu.Unlock()
			if !matched {
				cl.refuse(0, simNakUnknownPass, "unknown pass", extra)
			}
		case f.Type == simTypeHello:
			cl.greet(string(f.Body), now.Add(extra))
		case cl.pass == nil:
			cl.refuse(f.Seq, simNakHelloFirst, "hello first", extra)
		case f.Type == simTypeUplink:
			cl.pass.command(cl, f.Seq, string(f.Body), now.Add(extra))
		default:
			cl.refuse(f.Seq, simNakUnknownType, "unknown type", extra)
		}
	}
}

// greet answers a HELLO that arrives at an instant: WELCOME, or NAK unknown pass when the
// station has no such contact at that dish then, or NAK stowed while its dishes are stowed.
func (cl *simCall) greet(body string, arrive time.Time) {
	w := cl.w
	id, dishText, _ := strings.Cut(body, " ")
	dish, _ := strconv.Atoi(dishText)
	cl.mu.Lock()
	cl.hello, cl.passID, cl.dish = body, id, dish
	cl.mu.Unlock()
	extra := arrive.Sub(w.now())
	c, ok := w.thePlan().byID[id]
	switch {
	case !ok || c.Station != cl.station || c.Dish != dish || arrive.Before(c.AOS) || !arrive.Before(c.LOS):
		cl.refuse(0, simNakUnknownPass, "unknown pass", extra)
		cl.mu.Lock()
		cl.nak = simNakUnknownPass
		cl.mu.Unlock()
	case w.stowedAt(cl.station, arrive):
		cl.refuse(0, simNakStowed, "stowed", extra)
		cl.mu.Lock()
		cl.nak = simNakStowed
		cl.mu.Unlock()
	default:
		cl.mu.Lock()
		cl.welcomed = true
		cl.mu.Unlock()
		cl.pass = w.pass(c)
		cl.pass.attach(cl, extra)
	}
}

// refuse sends a NAK: the refused frame's seq, a code and the reason.
func (cl *simCall) refuse(seq uint32, code byte, reason string, extra time.Duration) {
	cl.send(simEncode(simTypeNak, seq, append([]byte{code}, reason...)), extra, nil, false)
}

// dead reports whether the call's line has died: nothing gets through it any more.
func (cl *simCall) dead() bool {
	return cl.w.diedSince(cl.station, cl.dialed, cl.w.now())
}

// send queues a message for the call's writer. extra delays it; a slow line delays it
// more and breaks it into pieces.
func (cl *simCall) send(data []byte, extra time.Duration, arrived func(), hangUp bool) {
	now := cl.w.now()
	out := simOut{data: data, notBefore: now.Add(extra), arrived: arrived, hangUp: hangUp}
	if n := cl.w.c.split[cl.station]; n > 0 {
		out.piece = n
	}
	if slow, ok := cl.w.slowAt(cl.station, now); ok {
		out.notBefore = out.notBefore.Add(slow.Delay)
		out.piece, out.gap = slow.Piece, slow.Gap
	}
	cl.mu.Lock()
	if cl.held != nil {
		out.data = append(cl.held, out.data...)
		cl.held = nil
	}
	cl.queue = append(cl.queue, out)
	cl.mu.Unlock()
	select {
	case cl.wake <- struct{}{}:
	default:
	}
}

// write is the call's writer: it writes each queued message in turn, whole or in pieces.
func (cl *simCall) write() {
	defer cl.w.wg.Done()
	for {
		out, ok := cl.next()
		if !ok || !cl.w.sleepUntil(out.notBefore) {
			return
		}
		if cl.dead() {
			continue
		}
		piece := out.piece
		if piece <= 0 {
			piece = len(out.data)
		}
		for i := 0; i < len(out.data); i += piece {
			if i > 0 && out.gap > 0 && !cl.w.sleepUntil(cl.w.now().Add(out.gap)) {
				return
			}
			if _, err := cl.conn.Write(out.data[i:min(i+piece, len(out.data))]); err != nil {
				return // the rehearsal has hung up
			}
		}
		if out.arrived != nil {
			out.arrived()
		}
		if out.hangUp {
			cl.hangUp(false)
			return
		}
	}
}

// next returns the next message to write, waiting for one; false when the call is over.
func (cl *simCall) next() (simOut, bool) {
	for {
		cl.mu.Lock()
		if len(cl.queue) > 0 {
			out := cl.queue[0]
			cl.queue = cl.queue[1:]
			cl.mu.Unlock()
			return out, true
		}
		cl.mu.Unlock()
		select {
		case <-cl.wake:
		case <-cl.gone:
			return simOut{}, false
		case <-cl.w.done:
			return simOut{}, false
		}
	}
}

// hangUp closes the call from the station's end: after a BYE, or with dropped set, without
// a word, as a broken line does.
func (cl *simCall) hangUp(dropped bool) {
	cl.mu.Lock()
	if cl.closed.IsZero() {
		cl.closed, cl.dropped = cl.w.now(), dropped
	}
	cl.mu.Unlock()
	cl.over()
}

// over ends the call: its pipe is closed and its writer stops.
func (cl *simCall) over() {
	cl.end.Do(func() {
		cl.conn.Close()
		close(cl.gone)
	})
}

// simPass is one contact as its station flies it: the satellite's telemetry, its answers
// to commands and its files, sent on whichever call holds the pass now.
type simPass struct {
	w    *simWorld
	c    Contact
	sat  *simSat
	wake chan struct{}

	mu      sync.Mutex
	call    *simCall     // the call the pass speaks on; nil: none
	cmds    []simPending // commands not yet answered, in the order they arrived
	streams []*simStream // files to send, in order
	ncmd    int          // commands received so far

	// Owned by the pass's goroutine.
	frames     int       // frames sent on the pass so far, WELCOME included
	tlm        int       // TLM frames sent so far
	free       time.Time // when the satellite answered its last command
	streamFree time.Time // when the last file stream ended, plus a second
}

// simPending is a command a station has received and the satellite has not answered yet.
type simPending struct {
	call   *simCall
	seq    uint32
	line   string
	arrive time.Time
	n      int // the pass's n-th command
}

// simStream is one file going down on a pass: FILE, a slot of chunks a second, then END.
type simStream struct {
	file    *simFile
	slots   [][]int   // the chunk numbers sent in each one-second slot
	startAt time.Time // FILE goes out at this instant, or once the stream before it has ended
	crc     uint32    // what END carries
	began   time.Time // when FILE went out; zero: not yet
	slot    int       // slots sent so far
}

// pass returns the pass a contact is, making it and starting its clock the first time.
func (w *simWorld) pass(c Contact) *simPass {
	w.mu.Lock()
	defer w.mu.Unlock()
	if p, ok := w.passes[c.ID]; ok {
		return p
	}
	p := &simPass{w: w, c: c, sat: w.satLocked(c.SatID), wake: make(chan struct{}, 1)}
	if b, ok := w.booked[c.ID]; ok {
		f := w.fileLocked(c, b)
		s := &simStream{file: f, startAt: c.AOS.Add(30 * time.Second), crc: f.crc}
		missing := w.c.missing[f.name]
		for _, slot := range f.order {
			s.slots = append(s.slots, slices.DeleteFunc(slices.Clone(slot), func(n int) bool { return slices.Contains(missing, n) }))
		}
		if slices.Contains(w.c.badCRC, f.name) {
			s.crc ^= 0xFFFFFFFF
		}
		p.streams = append(p.streams, s)
	}
	w.passes[c.ID] = p
	w.wg.Add(1)
	go p.run()
	return p
}

// attach makes a call the one the pass speaks on, and welcomes it: WELCOME with the dish
// and the satellite's last accepted command number, then FILE again if a file is on its
// way down.
func (p *simPass) attach(cl *simCall, extra time.Duration) {
	body := make([]byte, 9)
	body[0] = byte(p.c.Dish)
	binary.BigEndian.PutUint64(body[1:], p.sat.lastSeq())
	p.mu.Lock()
	old := p.call
	p.call = cl
	var again *simStream
	if len(p.streams) > 0 && !p.streams[0].began.IsZero() {
		again = p.streams[0]
	}
	p.mu.Unlock()
	if old != nil && old != cl {
		old.hangUp(true)
	}
	cl.send(simEncode(simTypeWelcome, 0, body), extra, nil, false)
	p.sent(cl)
	if again != nil {
		cl.send(again.file.head(), extra, nil, false)
		p.sent(cl)
	}
	p.nudge()
}

// nudge wakes the pass's goroutine to look at what is due.
func (p *simPass) nudge() {
	select {
	case p.wake <- struct{}{}:
	default:
	}
}

// sent counts a frame the pass has sent, and follows the n-th with a frame whose length is
// 8 when the case asks for one.
func (p *simPass) sent(cl *simCall) {
	p.mu.Lock()
	p.frames++
	bad := p.w.c.badLength[p.c.ID] == p.frames
	p.mu.Unlock()
	if bad {
		w := p.w
		cl.send([]byte{0, 8, 1, 2, 3, 4, 5, 6, 7, 8}, 0, func() {
			w.mu.Lock()
			w.badSizes++
			w.mu.Unlock()
		}, false)
	}
}

// live returns the call the pass can send on now, or nil when it has none or its line is
// dead: what the station relays then is gone.
func (p *simPass) live() *simCall {
	p.mu.Lock()
	cl := p.call
	p.mu.Unlock()
	if cl == nil || cl.dead() {
		return nil
	}
	select {
	case <-cl.gone:
		return nil // hung up
	default:
	}
	return cl
}

// command takes a command a call has brought: the satellite answers it two seconds after
// it arrives, one command at a time.
func (p *simPass) command(cl *simCall, seq uint32, line string, arrive time.Time) {
	p.mu.Lock()
	p.ncmd++
	p.cmds = append(p.cmds, simPending{call: cl, seq: seq, line: line, arrive: arrive, n: p.ncmd})
	p.mu.Unlock()
	p.nudge()
}

// run is the pass's clock, from its first WELCOME to its LOS: at each instant something is
// due, it does it.
func (p *simPass) run() {
	w := p.w
	defer w.wg.Done()
	late := w.c.late[p.c.Station]
	timer := time.NewTimer(time.Hour)
	defer timer.Stop()
	tick := p.c.AOS.Add(10 * time.Second) // when the satellite makes its next TLM frame
	for now := w.now(); tick.Add(late).Before(now); {
		tick = tick.Add(time.Minute) // made before the pass had a call: gone
	}
	for {
		now := w.now()
		if !now.Before(p.c.LOS) {
			p.goodbye("los")
			return
		}
		if w.stowedAt(p.c.Station, now) {
			p.goodbye("stowed")
		}
		next := p.c.LOS
		if at, ok := w.nextStow(p.c.Station, now); ok && at.Before(next) {
			next = at
		}
		if at, ok := p.answer(now); ok && at.Before(next) {
			next = at
		}
		for tick.Before(p.c.LOS) && !tick.Add(late).After(now) {
			p.telemetry(tick)
			tick = tick.Add(time.Minute)
		}
		if tick.Before(p.c.LOS) && tick.Add(late).Before(next) {
			next = tick.Add(late)
		}
		if at, ok := p.relay(now); ok && at.Before(next) {
			next = at
		}
		timer.Reset(next.Sub(now))
		select {
		case <-timer.C:
		case <-p.wake:
		case <-w.done:
			return
		}
	}
}

// goodbye sends BYE with a reason on the pass's call and hangs up.
func (p *simPass) goodbye(reason string) {
	p.mu.Lock()
	cl := p.call
	p.call = nil
	p.mu.Unlock()
	if cl != nil && !cl.dead() {
		cl.send(simEncode(simTypeBye, 0, []byte(reason)), 0, nil, true)
	}
}

// answer gives the satellite every command that is due and sends its answers. It returns
// the instant the next one is due, if one is waiting.
func (p *simPass) answer(now time.Time) (time.Time, bool) {
	for {
		p.mu.Lock()
		if len(p.cmds) == 0 {
			p.mu.Unlock()
			return time.Time{}, false
		}
		cmd := p.cmds[0]
		due := cmd.arrive
		if p.free.After(due) {
			due = p.free
		}
		due = due.Add(2 * time.Second)
		if due.After(now) {
			p.mu.Unlock()
			return due, true
		}
		p.cmds = p.cmds[1:]
		p.mu.Unlock()
		p.free = due
		p.execute(cmd, due)
	}
}

// execute hands one command to the satellite and answers it on the call that brought it:
// ACK when it ran, NAK with the satellite's reason when it refused.
func (p *simPass) execute(cmd simPending, at time.Time) {
	w, id := p.w, p.c.ID
	if w.c.dropBefore[id] == cmd.n {
		cmd.call.hangUp(true)
		return
	}
	rec := p.sat.accept(cmd.line, at, p.c.Station, id)
	if rec.Refused == "" && rec.Op == "DUMP_FILE" {
		p.redo(rec, at)
	}
	if w.c.dropAfter[id] == cmd.n {
		cmd.call.hangUp(true)
		return
	}
	cl := p.live()
	if cl == nil || cl != cmd.call {
		return // the call that brought it is gone, and its answer with it
	}
	if rec.Refused != "" {
		cl.refuse(cmd.seq, simNakRefused, rec.Refused, 0)
	} else {
		cl.send(simEncode(simTypeAck, cmd.seq, nil), 0, nil, false)
	}
	p.sent(cl)
}

// redo queues the chunks a DUMP_FILE asked for: they go down from two seconds after the
// command ran, before the pass's own file if that has not begun.
func (p *simPass) redo(rec simCommand, at time.Time) {
	w := p.w
	w.mu.Lock()
	f := w.files[rec.Args[0]]
	w.mu.Unlock()
	s := &simStream{file: f, startAt: at.Add(2 * time.Second), crc: f.crc}
	for _, a := range rec.Args[1:] {
		if n, err := strconv.Atoi(a); err == nil && n >= 0 && n < len(f.chunks) {
			s.slots = append(s.slots, []int{n})
		}
	}
	if len(rec.Args) == 1 {
		for n := range f.chunks {
			s.slots = append(s.slots, []int{n})
		}
	}
	p.mu.Lock()
	at0 := len(p.streams)
	for i, old := range p.streams {
		if old.began.IsZero() && old.startAt.Equal(p.c.AOS.Add(30*time.Second)) {
			at0 = i // the pass's own file, not begun: it waits
			break
		}
	}
	p.streams = slices.Insert(p.streams, at0, s)
	p.mu.Unlock()
}

// telemetry relays the TLM frame the satellite made at an instant.
func (p *simPass) telemetry(made time.Time) {
	cl := p.live()
	if cl == nil {
		return
	}
	w := p.w
	sat, line := p.c.SatID, simTelemetry(p.c.SatID, made)
	frame := simEncode(simTypeTLM, 0, []byte(line))
	p.tlm++
	w.mu.Lock()
	w.relayed[p.c.Station]++
	every := w.c.flipEvery[p.c.Station]
	flip := w.c.corrupt[p.c.ID] == p.tlm || (every > 0 && w.relayed[p.c.Station]%every == 0)
	w.mu.Unlock()
	if flip {
		frame[9] ^= 0x20 // a byte of the line: the CRC no longer matches
	}
	cl.send(frame, 0, func() { w.heard(sat, line, made, !flip) }, false)
	p.sent(cl)
}

// heard notes a TLM frame that reached the rehearsal whole.
func (w *simWorld) heard(sat, line string, made time.Time, good bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !good {
		w.tlmBad++
		return
	}
	w.tlmGood++
	if strings.Contains(line, "batt_v=nan") {
		w.tlmNaN++
		return
	}
	if old, ok := w.newest[sat]; !ok || simFrameNumber(old) < made.Unix() {
		w.newest[sat] = line
	}
}

// simFrameNumber returns the frame number of a telemetry line.
func simFrameNumber(line string) int64 {
	f := strings.Fields(line)
	if len(f) < 2 {
		return -1
	}
	n, _ := strconv.ParseInt(f[1], 10, 64)
	return n
}

// relay sends every frame of the pass's files that is due: FILE, then a slot of chunks a
// second, then END a second after the last. It returns the instant the next is due, if any.
func (p *simPass) relay(now time.Time) (time.Time, bool) {
	for {
		p.mu.Lock()
		if len(p.streams) == 0 {
			p.mu.Unlock()
			return time.Time{}, false
		}
		s := p.streams[0]
		p.mu.Unlock()
		if s.began.IsZero() {
			at := s.startAt
			if p.streamFree.After(at) {
				at = p.streamFree
			}
			if at.After(now) {
				return at, true
			}
			p.mu.Lock()
			s.began = now
			p.mu.Unlock()
			if cl := p.live(); cl != nil {
				cl.send(s.file.head(), 0, nil, false)
				p.sent(cl)
			}
			continue
		}
		due := s.began.Add(time.Duration(s.slot+1) * time.Second)
		if due.After(now) {
			return due, true
		}
		if s.slot < len(s.slots) {
			for _, n := range s.slots[s.slot] {
				p.chunk(s.file, n)
			}
			s.slot++
			continue
		}
		if cl := p.live(); cl != nil {
			end := make([]byte, 4)
			binary.BigEndian.PutUint32(end, s.crc)
			cl.send(simEncode(simTypeEnd, 0, end), 0, nil, false)
			p.sent(cl)
		}
		p.mu.Lock()
		p.streams = p.streams[1:]
		p.mu.Unlock()
		p.streamFree = now.Add(time.Second)
	}
}

// chunk relays one chunk of a file.
func (p *simPass) chunk(f *simFile, n int) {
	cl := p.live()
	if cl == nil {
		return
	}
	w, id := p.w, p.c.ID
	body := make([]byte, 2, 2+simChunkBytes)
	binary.BigEndian.PutUint16(body, uint16(n))
	body = append(body, f.chunks[n]...)
	cl.send(simEncode(simTypeChunk, 0, body), 0, func() {
		w.mu.Lock()
		w.chunks = append(w.chunks, simChunk{File: f.name, N: n, At: w.now(), PassID: id})
		w.mu.Unlock()
	}, false)
	p.sent(cl)
}

// simFile is a file a satellite sends down on a booked pass.
type simFile struct {
	name    string   // "thistle-2-20290726-1603.dat"
	passID  string   // the booked pass it belongs to
	booking Booking  // that pass's booking
	chunks  [][]byte // 256 bytes each
	data    []byte   // the whole file
	crc     uint32   // its CRC-32
	order   [][]int  // the chunk numbers its first pass sends in each slot: now and then two swapped, or one twice
}

// head returns the file's FILE frame: its chunk count, then its name.
func (f *simFile) head() []byte {
	body := make([]byte, 2, 2+len(f.name))
	binary.BigEndian.PutUint16(body, uint16(len(f.chunks)))
	return simEncode(simTypeFile, 0, append(body, f.name...))
}

// simFileName returns the name of the file a contact's satellite sends down: the
// satellite, the AOS day and the AOS time, ".tif" for a Furrow and ".dat" for any other.
func simFileName(c Contact) string {
	ext := ".dat"
	if sat, _ := simSatellite(c.SatID); sat.Operator == "Furrowlight" {
		ext = ".tif"
	}
	return c.SatID + "-" + c.AOS.UTC().Format("20060102-1504") + ext
}

// simFileChunks returns how many chunks an operator's files have: Furrowlight's 48,
// Cirrine's 24, Brinewatch's 16 and Kessick's, or anyone else's, 8.
func simFileChunks(operator string) int {
	switch operator {
	case "Furrowlight":
		return 48
	case "Cirrine Atmospherics":
		return 24
	case "Brinewatch":
		return 16
	}
	return 8
}

// fileLocked returns the file a booked contact sends down, making it the first time. The
// world's mutex is held.
func (w *simWorld) fileLocked(c Contact, b Booking) *simFile {
	name := simFileName(c)
	if f, ok := w.files[name]; ok {
		return f
	}
	sat, _ := simSatellite(c.SatID)
	n := simFileChunks(sat.Operator)
	f := &simFile{name: name, passID: c.ID, booking: b, data: make([]byte, n*simChunkBytes)}
	x := uint64(simHash(name))<<32 | uint64(simHash(c.ID)) | 1
	for i := range f.data {
		x ^= x << 13
		x ^= x >> 7
		x ^= x << 17
		f.data[i] = byte(x >> 24)
	}
	for i := range n {
		f.chunks = append(f.chunks, f.data[i*simChunkBytes:(i+1)*simChunkBytes])
		f.order = append(f.order, []int{i})
	}
	f.crc = crc32.ChecksumIEEE(f.data)
	h := int(simHash(name) >> 4)
	if h%4 == 1 {
		i := h / 4 % (n - 1) // two chunks come the other way round
		f.order[i], f.order[i+1] = f.order[i+1], f.order[i]
	}
	if h%5 == 2 {
		i := h / 5 % n // one chunk comes twice
		f.order[i] = append(f.order[i], f.order[i][0])
	}
	w.files[name] = f
	return f
}

// simMaster is one line of uplink-keys.txt.
type simMaster struct {
	ID       string
	Activate time.Time
	Retire   time.Time // zero: not scheduled
	Revoked  bool
	Secret   []byte
}

// simParseKeys reads uplink-keys.txt as the case wrote it, skipping what it can't read.
func simParseKeys(text string) []simMaster {
	var out []simMaster
	for _, line := range strings.Split(text, "\n") {
		f := strings.Fields(line)
		if len(f) < 4 || strings.HasPrefix(f[0], "#") {
			continue
		}
		m := simMaster{ID: f[0], Revoked: len(f) > 4 && f[4] == "revoked"}
		var err error
		if m.Activate, err = time.Parse(time.RFC3339, f[1]); err != nil {
			continue
		}
		if f[2] != "-" {
			if m.Retire, err = time.Parse(time.RFC3339, f[2]); err != nil {
				continue
			}
		}
		if m.Secret, err = hex.DecodeString(f[3]); err != nil {
			continue
		}
		out = append(out, m)
	}
	slices.SortFunc(out, func(a, b simMaster) int { return a.Activate.Compare(b.Activate) })
	return out
}

// simSatKey derives a satellite's key from a master, by Season 5's rule: purpose is "sign"
// or "seal".
func simSatKey(master []byte, purpose, sat string) []byte {
	k, err := hkdf.Key(sha256.New, master, []byte("apogee uplink keys v1"), purpose+"|"+sat, 32)
	if err != nil {
		panic(err) // 32 bytes from SHA-256 is always possible
	}
	return k
}

// simSat is a satellite: the keys it holds, its last accepted command number, and every
// command it has run or refused.
type simSat struct {
	w      *simWorld
	id     string
	sealed bool // a Cirrine Atmospherics satellite: it refuses a CMD line

	mu       sync.Mutex
	last     uint64
	keys     map[string]simSatKeys // by master ID
	commands []simCommand
	nonces   map[string]int // how often each SEAL nonce has come
}

// simSatKeys are a satellite's two keys from one master.
type simSatKeys struct {
	Sign, Seal []byte
}

// simCommand is one command line a satellite was given.
type simCommand struct {
	Sat     string
	Seq     uint64
	Op      string   // "PING"; "" when the line could not be read that far
	Args    []string // after the operation
	KeyID   string   // the master whose key checked it; "" when none did
	Sealed  bool     // it came as a SEAL line
	At      time.Time
	Station string
	PassID  string
	Line    string // as it arrived
	Refused string // why the satellite refused it; "" when it ran
}

// simFirstLast returns a satellite's last accepted command number at the start of a run.
func simFirstLast(id string) uint64 {
	switch id {
	case "pipit-1":
		return 5518
	case "pipit-2":
		return 6314
	case "pipit-3":
		return 2907
	case "furrow-a":
		return 8841
	case "furrow-b":
		return 8790
	case "brine-3":
		return 7426
	case "cirrine-1":
		return 9712
	case "thistle-2":
		return 412
	}
	return 40 + uint64(simHash(id)%300) // the satellites launched this spring
}

// satLocked returns a satellite, making it the first time: it holds every master activated
// before the window and not retired or revoked then. The world's mutex is held.
func (w *simWorld) satLocked(id string) *simSat {
	if s, ok := w.sats[id]; ok {
		return s
	}
	sat, _ := simSatellite(id)
	s := &simSat{w: w, id: id, sealed: sat.Operator == "Cirrine Atmospherics", last: simFirstLast(id),
		keys: map[string]simSatKeys{}, nonces: map[string]int{}}
	if n, ok := w.c.lastSeq[id]; ok {
		s.last = n
	}
	for _, m := range w.masters {
		if m.Activate.Before(w.win.From) && !m.Revoked && (m.Retire.IsZero() || m.Retire.After(w.win.From)) {
			s.keys[m.ID] = simSatKeys{simSatKey(m.Secret, "sign", id), simSatKey(m.Secret, "seal", id)}
		}
	}
	w.sats[id] = s
	return s
}

// sat returns a satellite of the run.
func (w *simWorld) sat(id string) *simSat {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.satLocked(id)
}

// lastSeq returns the satellite's last accepted command number.
func (s *simSat) lastSeq() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.last
}

// usable returns the masters whose keys the satellite holds and may still use at an
// instant, newest first.
func (s *simSat) usable(at time.Time) []simMaster {
	var out []simMaster
	for _, m := range s.w.masters {
		if _, ok := s.keys[m.ID]; ok && !m.Revoked && (m.Retire.IsZero() || at.Before(m.Retire)) {
			out = append(out, m)
		}
	}
	slices.Reverse(out)
	return out
}

// accept checks one command line by Season 5's order (the line, the version, the
// satellite, the key, the number, the signature or the seal) and runs it only if all of it
// holds. It records the command either way.
func (s *simSat) accept(line string, at time.Time, station, passID string) simCommand {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec := simCommand{Sat: s.id, At: at, Station: station, PassID: passID, Line: line}
	rec.Refused = s.check(&rec, at)
	if rec.Refused == "" {
		rec.Refused = s.perform(&rec)
	}
	if rec.Refused == "" {
		s.last = rec.Seq
	}
	s.commands = append(s.commands, rec)
	return rec
}

// check reads and verifies a line, filling rec; it returns why the satellite refuses it,
// or "".
func (s *simSat) check(rec *simCommand, at time.Time) string {
	f := strings.Split(rec.Line, " ")
	switch f[0] {
	case "CMD":
		if len(f) < 5 || !strings.HasPrefix(f[len(f)-1], "sig=") {
			return "malformed command line: want CMD <sat> <seq> <OP> [args] sig=<hex>"
		}
		sig, err := hex.DecodeString(strings.TrimPrefix(f[len(f)-1], "sig="))
		seq, err2 := strconv.ParseUint(f[2], 10, 64)
		if err != nil || len(sig) != 32 || err2 != nil || seq == 0 {
			return "malformed command line: want CMD <sat> <seq> <OP> [args] sig=<hex>"
		}
		rec.Seq, rec.Op, rec.Args = seq, f[3], f[4:len(f)-1]
		if f[1] != s.id {
			return fmt.Sprintf("not for this satellite: %q", f[1])
		}
		if s.sealed || rec.Op == "LOAD_KEY" {
			return "command must be sealed"
		}
		usable := s.usable(at)
		if len(usable) == 0 {
			return "command failed authentication"
		}
		if seq <= s.last {
			return fmt.Sprintf("replayed sequence number: %d, last %d", seq, s.last)
		}
		canonical := strings.Join(append([]string{f[1], f[2], rec.Op}, rec.Args...), "|")
		for _, m := range usable {
			mac := hmac.New(sha256.New, s.keys[m.ID].Sign)
			mac.Write([]byte(canonical))
			if hmac.Equal(mac.Sum(nil), sig) {
				rec.KeyID = m.ID
				return ""
			}
		}
		return "command failed authentication"
	case "SEAL":
		rec.Sealed = true
		if len(f) != 7 {
			return "malformed command line: want 7 fields"
		}
		if f[1] != "2" {
			return "unsupported version: " + f[1]
		}
		seq, err := strconv.ParseUint(f[3], 10, 64)
		nonce, err2 := hex.DecodeString(f[5])
		sealed, err3 := hex.DecodeString(f[6])
		if err != nil || seq == 0 || err2 != nil || len(nonce) != 12 || err3 != nil {
			return "malformed command line: want SEAL 2 <sat> <seq> <key id> <nonce> <ciphertext>"
		}
		rec.Seq = seq
		if f[2] != s.id {
			return fmt.Sprintf("not for this satellite: %q", f[2])
		}
		keys, held := s.keys[f[4]]
		if !held {
			return fmt.Sprintf("unknown key: %q", f[4])
		}
		for _, m := range s.w.masters {
			if m.ID == f[4] && m.Revoked {
				return fmt.Sprintf("key revoked: %q", f[4])
			}
			if m.ID == f[4] && !m.Retire.IsZero() && !at.Before(m.Retire) {
				return fmt.Sprintf("key retired: %q", f[4])
			}
		}
		if seq <= s.last {
			return fmt.Sprintf("replayed sequence number: %d, last %d", seq, s.last)
		}
		block, err := aes.NewCipher(keys.Seal)
		if err != nil {
			return "command failed authentication"
		}
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return "command failed authentication"
		}
		plain, err := gcm.Open(nil, nonce, sealed, []byte(strings.Join([]string{"SEAL", "2", f[2], f[3], f[4]}, "|")))
		if err != nil {
			return "command failed authentication"
		}
		s.nonces[f[5]]++
		parts := strings.Split(string(plain), "|")
		rec.KeyID, rec.Op, rec.Args = f[4], parts[0], parts[1:]
		return ""
	}
	return "malformed command line: want CMD or SEAL"
}

// perform runs a command that checked out; it returns why the satellite can't, or "".
func (s *simSat) perform(rec *simCommand) string {
	switch rec.Op {
	case "PING":
		if len(rec.Args) != 0 {
			return "malformed command line: PING takes no arguments"
		}
	case "LOAD_KEY":
		if len(rec.Args) != 3 {
			return "malformed command line: LOAD_KEY wants a key ID and two keys"
		}
		sign, err := hex.DecodeString(rec.Args[1])
		seal, err2 := hex.DecodeString(rec.Args[2])
		if err != nil || err2 != nil || len(sign) != 32 || len(seal) != 32 {
			return "malformed command line: LOAD_KEY wants a key ID and two keys"
		}
		s.keys[rec.Args[0]] = simSatKeys{sign, seal}
	case "DUMP_FILE":
		if len(rec.Args) == 0 {
			return "malformed command line: DUMP_FILE wants a file"
		}
		s.w.mu.Lock()
		f, ok := s.w.files[rec.Args[0]]
		s.w.mu.Unlock()
		if !ok || !strings.HasPrefix(f.name, s.id+"-") {
			return fmt.Sprintf("no such file: %q", rec.Args[0])
		}
	default:
		return fmt.Sprintf("unknown operation: %q", rec.Op)
	}
	return ""
}

// client returns World.HTTP: a client whose requests go to the customers' simulated
// receivers and nowhere else.
func (w *simWorld) client() *http.Client {
	return &http.Client{
		Transport: w.hooks,
		Timeout:   5 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
