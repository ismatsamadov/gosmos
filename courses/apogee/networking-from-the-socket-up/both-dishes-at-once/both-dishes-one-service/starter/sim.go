package main

// The maker's simulator: one dish controller, in memory, speaking the controller protocol
// exactly as the manual should have described it. A test makes one with NewSim, sets any
// fields it needs, serves it on a loopback listener (Serve) or on one end of a net.Pipe
// (ServeConn), and reads what it saw through its methods. Every top-level name in this
// file begins with sim or Sim, so it never collides with dishd's own names.

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"math"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// Sim is one simulated dish controller. Set its exported fields before it serves; while
// it serves, it reads them under its own lock.
type Sim struct {
	// The controller's state.
	LastSeq uint32  // the sequence number of the last command executed
	Wind    float64 // the station's wind in km/h; at 72 or more the dish stows and locks
	TempC   float64 // inside the radome
	AzDeg   float64 // azimuth
	ElDeg   float64 // elevation
	State   uint8   // 0 IDLE, 1 MOVING, 2 TRACKING, 3 STOWED, 4 FAULT
	Flags   uint16  // bit 0 heater on, bit 1 wind lock (kept by Wind), bit 2 fault

	// How it behaves on the line.
	Split      int              // write every message in pieces of at most this many bytes, one Write each
	Gap        time.Duration    // pause before each piece after the first
	Glue       int              // hold the answers to a connection's first n frames, then write them in one Write
	HoldBanner bool             // write the banner together with the first answers instead of at once
	Silent     bool             // banner, then read and never answer
	MuteAfter  int              // answer this many frames on a connection (HELLO counted), then read and never answer
	NoRead     bool             // banner, then never read: ServeConn returns and leaves the connection open
	DropBefore uint32           // when the command with this seq arrives, close without executing or answering, once
	DropAfter  uint32           // execute the command with this seq, then close without answering, once
	Before     func(req string) // called before each answer is written, with the request's name ("QUERY")
	IdleLimit  time.Duration    // no complete frame for this long: BYE idle, then close; 0 means no limit
	FrameLimit time.Duration    // from a frame's first byte to its last: BYE slow frame, then close
	WriteLimit time.Duration    // a piece the client doesn't read within it: close, without BYE

	name string
	dish uint8

	mu            sync.Mutex
	current       *simConn // the connection being served: the newest one wins
	executed      []string
	received      []string
	conns         int
	byClient      int
	droppedBefore bool
	droppedAfter  bool
}

// NewSim returns a controller called name (YST-1) for dish number dish: firmware 2.4,
// IDLE at az 180.00 el 45.00, wind 20.0 km/h, -12.5 °C, heater on, last_seq 0.
func NewSim(name string, dish uint8) *Sim {
	return &Sim{
		Wind:       20,
		TempC:      -12.5,
		AzDeg:      180,
		ElDeg:      45,
		Flags:      simHeater,
		IdleLimit:  30 * time.Second,
		FrameLimit: 10 * time.Second,
		WriteLimit: 10 * time.Second,
		name:       name,
		dish:       dish,
	}
}

// Serve accepts connections on ln until ln is closed, then closes them, waits for them to
// end and returns nil. The newest connection wins: the one before it gets BYE replaced.
func (s *Sim) Serve(ln net.Listener) error {
	var wg sync.WaitGroup
	var conns []net.Conn
	for {
		c, err := ln.Accept()
		if err != nil {
			for _, c := range conns {
				c.Close()
			}
			wg.Wait()
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		conns = append(conns, c)
		wg.Go(func() { s.ServeConn(c) })
	}
}

// ServeConn serves one connection until it ends: the client hangs up, the controller does,
// or a newer connection replaces it. Tests use it on one end of a net.Pipe.
func (s *Sim) ServeConn(c net.Conn) {
	sc := &simConn{sim: s, conn: c, br: bufio.NewReader(c), wake: make(chan struct{}, 1)}
	banner := []byte(s.name + " READY PROTO 1\r\n")
	s.mu.Lock()
	sc.noRead, sc.glue = s.NoRead, s.Glue
	if s.HoldBanner {
		sc.banner = banner
	}
	s.conns++
	old := s.current
	s.current = sc
	s.mu.Unlock()
	if old != nil {
		old.replace()
	}
	if sc.noRead {
		if sc.write(banner) != nil {
			c.Close()
		}
		return
	}
	if sc.banner == nil {
		sc.push(simOut{data: banner})
	}
	done := make(chan struct{})
	go sc.writeLoop(done)
	sc.readLoop()
	sc.closeQueue()
	<-done
	c.Close()
	s.mu.Lock()
	if s.current == sc {
		s.current = nil
	}
	s.mu.Unlock()
}

// Executed returns the commands the controller executed, in order: "119 POINT 183.25 41.70",
// "120 STOW", "121 UNSTOW", "122 JOG A +0.35".
func (s *Sim) Executed() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.executed...)
}

// Received returns every frame the controller received, in order, as Executed writes
// commands, and "HELLO dishd", "PING", "QUERY" or "bad crc" for the rest.
func (s *Sim) Received() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.received...)
}

// Conns returns how many connections the controller has accepted so far.
func (s *Sim) Conns() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conns
}

// ClosedByClient returns how many connections the client closed before the controller did.
func (s *Sim) ClosedByClient() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.byClient
}

// LastExecuted returns last_seq now: the sequence number of the last command executed.
func (s *Sim) LastExecuted() uint32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.LastSeq
}

// SetWind sets the station's wind in km/h. At 72 or more the dish stows and locks at once.
func (s *Sim) SetWind(kmh float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Wind = kmh
	s.weather()
}

// Frame types, flags and states the simulator uses.
const (
	simHello   byte = 0x01
	simPing    byte = 0x02
	simQuery   byte = 0x03
	simPoint   byte = 0x10
	simStow    byte = 0x11
	simUnstow  byte = 0x12
	simJog     byte = 0x13
	simWelcome byte = 0x81
	simPong    byte = 0x82
	simStatus  byte = 0x83
	simAck     byte = 0x90
	simNak     byte = 0x91
	simBye     byte = 0xFF

	simHeater   uint16 = 1 << 0
	simWindLock uint16 = 1 << 1

	simIdle    uint8 = 0
	simStowed  uint8 = 3
	simLockKmh       = 72.0
)

// simEnd is why the controller ends a connection: the reason its BYE gives.
type simEnd string

// Error returns the reason.
func (e simEnd) Error() string { return string(e) }

// Reasons the controller hangs up.
const (
	simIdleEnd     simEnd = "idle"
	simSlowEnd     simEnd = "slow frame"
	simLengthEnd   simEnd = "bad length"
	simReplacedEnd simEnd = "replaced"
)

// simOut is one message for a connection's writer: its bytes, the requests they answer,
// and whether the controller hangs up after it.
type simOut struct {
	data   []byte
	reqs   []string
	hangUp bool
}

// simConn is one connection the simulator serves: a reader that answers frames as they
// arrive, and a writer that sends the answers in order, so the controller keeps reading
// while a slow client reads its answers, as it would on a real line.
type simConn struct {
	sim      *Sim
	conn     net.Conn
	br       *bufio.Reader
	replaced atomic.Bool
	noRead   bool

	// used by the reader only
	helloed  bool
	answered int
	glue     int
	banner   []byte // held until the first answers, with HoldBanner
	held     []byte
	heldReqs []string

	qmu     sync.Mutex
	queue   []simOut
	qclosed bool
	wake    chan struct{}
}

// replace tells a connection that a newer one has taken its place.
func (sc *simConn) replace() {
	if sc.noRead {
		sc.conn.Close()
		return
	}
	sc.replaced.Store(true)
	sc.conn.SetReadDeadline(time.Now())
}

// readLoop reads and answers frames until the connection ends. A frame that arrives after
// a newer connection has replaced this one is not answered.
func (sc *simConn) readLoop() {
	for {
		payload, err := sc.readFrame()
		if err == nil && sc.replaced.Load() {
			err = simReplacedEnd
		}
		if err != nil {
			sc.end(err)
			return
		}
		if !sc.handle(payload) {
			return
		}
	}
}

// readFrame reads one frame and returns the bytes after its length: type, seq, body and CRC.
func (sc *simConn) readFrame() ([]byte, error) {
	idle, frame := sc.sim.limits()
	if sc.readWithin(idle) {
		return nil, simReplacedEnd
	}
	first, err := sc.br.ReadByte()
	if err != nil {
		return nil, sc.why(err, simIdleEnd)
	}
	if sc.readWithin(frame) {
		return nil, simReplacedEnd
	}
	second, err := sc.br.ReadByte()
	if err != nil {
		return nil, sc.why(err, simSlowEnd)
	}
	n := int(first)<<8 | int(second)
	if n < 9 || n > 1024 {
		sc.sim.note(fmt.Sprintf("bad length %d", n))
		return nil, simLengthEnd
	}
	payload := make([]byte, n)
	if _, err := io.ReadFull(sc.br, payload); err != nil {
		return nil, sc.why(err, simSlowEnd)
	}
	return payload, nil
}

// readWithin gives the next read d to finish, or no limit when d is 0, and reports whether
// a newer connection has replaced this one.
func (sc *simConn) readWithin(d time.Duration) bool {
	var t time.Time
	if d > 0 {
		t = time.Now().Add(d)
	}
	sc.conn.SetReadDeadline(t)
	return sc.replaced.Load()
}

// why turns a failed read into the reason the connection ends: replaced, the limit that
// ran out, or the error itself.
func (sc *simConn) why(err error, limit simEnd) error {
	if sc.replaced.Load() {
		return simReplacedEnd
	}
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return limit
	}
	return err
}

// end finishes a connection whose reading stopped with err.
func (sc *simConn) end(err error) {
	if why, ok := err.(simEnd); ok {
		if why == simReplacedEnd && sc.peerGone() {
			sc.sim.closedByClient()
			return
		}
		sc.flushHeld()
		sc.push(simOut{data: simFrame(simBye, 0, []byte(why)), hangUp: true})
		return
	}
	if simPeerClosed(err) {
		sc.sim.closedByClient()
	}
	// Anything else means the controller closed the connection itself.
}

// peerGone reports whether the client has already hung up. It looks once, for a moment.
func (sc *simConn) peerGone() bool {
	sc.conn.SetReadDeadline(time.Now().Add(5 * time.Millisecond))
	_, err := sc.br.ReadByte()
	return simPeerClosed(err)
}

// handle answers one frame. It reports false when the frame ends the connection.
func (sc *simConn) handle(p []byte) bool {
	s := sc.sim
	typ, seq, body := p[0], binary.BigEndian.Uint32(p[1:5]), p[5:len(p)-4]
	crcOK := crc32.ChecksumIEEE(p[:len(p)-4]) == binary.BigEndian.Uint32(p[len(p)-4:])
	req := simName(typ)
	if !crcOK {
		req = "bad crc"
	}
	s.mu.Lock()
	if crcOK {
		s.received = append(s.received, simDescribe(typ, seq, body))
	} else {
		s.received = append(s.received, req)
	}
	if s.Silent || (s.MuteAfter > 0 && sc.answered >= s.MuteAfter) {
		s.mu.Unlock()
		return true
	}
	answer, ok := sc.answer(typ, seq, body, crcOK)
	s.mu.Unlock()
	if !ok {
		sc.flushHeld()
		sc.push(simOut{hangUp: true})
		return false
	}
	sc.answered++
	sc.held = append(sc.held, answer...)
	sc.heldReqs = append(sc.heldReqs, req)
	if sc.glue > 1 {
		sc.glue--
		return true
	}
	sc.glue = 0
	sc.flushHeld()
	return true
}

// answer decides the controller's answer to one frame, checking in the protocol's order:
// the CRC, HELLO first, a known type, then a command's own checks. It reports false when
// the frame drops the line instead (DropBefore, DropAfter). It runs under s.mu.
func (sc *simConn) answer(typ byte, seq uint32, body []byte, crcOK bool) ([]byte, bool) {
	s := sc.sim
	s.weather()
	switch {
	case !crcOK:
		return simNakFrame(seq, 6, "bad crc"), true
	case !sc.helloed && typ != simHello:
		return simNakFrame(seq, 7, "hello first"), true
	}
	switch typ {
	case simHello:
		if len(body) > 16 || !simASCII(body) {
			return simNakFrame(seq, 2, "bad argument"), true
		}
		sc.helloed = true
		w := binary.BigEndian.AppendUint32([]byte{s.dish, 2, 4}, s.LastSeq)
		return simFrame(simWelcome, 0, w), true
	case simPing, simQuery:
		if len(body) != 0 {
			return simNakFrame(seq, 2, "bad argument"), true
		}
		if typ == simPing {
			return simFrame(simPong, 0, nil), true
		}
		return simFrame(simStatus, 0, s.status()), true
	case simPoint, simStow, simUnstow, simJog:
		return sc.command(typ, seq, body)
	}
	return simNakFrame(seq, 5, "unknown type"), true
}

// command checks one command and executes it: the sequence number, the body, the wind
// lock, stowed. It runs under s.mu.
func (sc *simConn) command(typ byte, seq uint32, body []byte) ([]byte, bool) {
	s := sc.sim
	if seq != 0 && seq == s.DropBefore && !s.droppedBefore {
		s.droppedBefore = true
		return nil, false
	}
	if seq <= s.LastSeq {
		return simNakFrame(seq, 1, "duplicate"), true
	}
	az := int(math.Round(s.AzDeg * 100))
	el := int(math.Round(s.ElDeg * 100))
	bad := simNakFrame(seq, 2, "bad argument")
	switch typ {
	case simPoint:
		if len(body) != 4 {
			return bad, true
		}
		az = int(binary.BigEndian.Uint16(body))
		el = int(int16(binary.BigEndian.Uint16(body[2:])))
		if az > 35999 || el < 0 || el > 9000 {
			return bad, true
		}
	case simStow, simUnstow:
		if len(body) != 0 {
			return bad, true
		}
	case simJog:
		if len(body) != 3 {
			return bad, true
		}
		delta := int(int16(binary.BigEndian.Uint16(body[1:])))
		if delta == 0 || delta < -500 || delta > 500 {
			return bad, true
		}
		switch body[0] {
		case 'A':
			az = ((az+delta)%36000 + 36000) % 36000
		case 'E':
			el += delta
			if el < 0 || el > 9000 {
				return bad, true
			}
		default:
			return bad, true
		}
	}
	if typ != simStow && s.Wind >= simLockKmh {
		return simNakFrame(seq, 3, "wind lock"), true
	}
	if (typ == simPoint || typ == simJog) && s.State == simStowed {
		return simNakFrame(seq, 4, "stowed"), true
	}
	switch typ {
	case simPoint, simJog:
		s.AzDeg, s.ElDeg, s.State = float64(az)/100, float64(el)/100, simIdle
	case simStow:
		s.ElDeg, s.State = 90, simStowed
	case simUnstow:
		if s.State == simStowed {
			s.State = simIdle
		}
	}
	s.LastSeq = seq
	s.executed = append(s.executed, simDescribe(typ, seq, body))
	if seq != 0 && seq == s.DropAfter && !s.droppedAfter {
		s.droppedAfter = true
		return nil, false
	}
	return simFrame(simAck, seq, nil), true
}

// weather applies the wind lock: at 72 km/h or more the dish stows and the lock flag is
// set; below it the flag clears and the dish stays stowed until UNSTOW. It runs under s.mu.
func (s *Sim) weather() {
	if s.Wind >= simLockKmh {
		s.State, s.ElDeg = simStowed, 90
		s.Flags |= simWindLock
	} else {
		s.Flags &^= simWindLock
	}
}

// status is the 16-byte STATUS body for the controller's state now. It runs under s.mu.
func (s *Sim) status() []byte {
	b := []byte{s.dish, s.State}
	b = binary.BigEndian.AppendUint16(b, uint16(math.Round(s.AzDeg*100)))
	b = binary.BigEndian.AppendUint16(b, uint16(int16(math.Round(s.ElDeg*100))))
	b = binary.BigEndian.AppendUint16(b, uint16(math.Round(s.Wind*10)))
	b = binary.BigEndian.AppendUint16(b, uint16(int16(math.Round(s.TempC*10))))
	b = binary.BigEndian.AppendUint16(b, s.Flags)
	return binary.BigEndian.AppendUint32(b, s.LastSeq)
}

// limits returns the controller's idle and frame limits.
func (s *Sim) limits() (idle, frame time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.IdleLimit, s.FrameLimit
}

// line returns how the controller writes: Split, Gap and WriteLimit.
func (s *Sim) line() (split int, gap, limit time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Split, s.Gap, s.WriteLimit
}

// before returns the Before hook, if any.
func (s *Sim) before() func(string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Before
}

// note adds an entry to the list of frames received.
func (s *Sim) note(entry string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.received = append(s.received, entry)
}

// closedByClient counts a connection the client closed first.
func (s *Sim) closedByClient() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byClient++
}

// flushHeld queues the held banner and answers as one message.
func (sc *simConn) flushHeld() {
	if len(sc.banner) == 0 && len(sc.held) == 0 {
		return
	}
	data := append(sc.banner, sc.held...)
	sc.push(simOut{data: data, reqs: sc.heldReqs})
	sc.banner, sc.held, sc.heldReqs = nil, nil, nil
}

// push queues a message for the writer.
func (sc *simConn) push(out simOut) {
	sc.qmu.Lock()
	sc.queue = append(sc.queue, out)
	sc.qmu.Unlock()
	sc.poke()
}

// closeQueue tells the writer that no more messages will come.
func (sc *simConn) closeQueue() {
	sc.qmu.Lock()
	sc.qclosed = true
	sc.qmu.Unlock()
	sc.poke()
}

// poke wakes the writer if it is waiting.
func (sc *simConn) poke() {
	select {
	case sc.wake <- struct{}{}:
	default:
	}
}

// next waits for the next message; it reports false once the queue is closed and empty.
func (sc *simConn) next() (simOut, bool) {
	for {
		sc.qmu.Lock()
		if len(sc.queue) > 0 {
			out := sc.queue[0]
			sc.queue = sc.queue[1:]
			sc.qmu.Unlock()
			return out, true
		}
		closed := sc.qclosed
		sc.qmu.Unlock()
		if closed {
			return simOut{}, false
		}
		<-sc.wake
	}
}

// writeLoop writes the queued messages in order until the queue is closed.
func (sc *simConn) writeLoop(done chan<- struct{}) {
	defer close(done)
	stopped := false
	for {
		out, ok := sc.next()
		if !ok {
			return
		}
		if stopped {
			continue
		}
		if hook := sc.sim.before(); hook != nil {
			for _, req := range out.reqs {
				hook(req)
			}
		}
		if err := sc.write(out.data); err != nil {
			stopped = true
			if errors.Is(err, os.ErrDeadlineExceeded) {
				sc.conn.Close() // the client stopped reading: hang up without BYE
			}
			continue
		}
		if out.hangUp {
			stopped = true
			sc.conn.Close()
		}
	}
}

// write sends one message in pieces of at most Split bytes, pausing Gap before each piece
// after the first; the client must read each piece within WriteLimit.
func (sc *simConn) write(data []byte) error {
	split, gap, limit := sc.sim.line()
	for first := true; len(data) > 0; first = false {
		n := len(data)
		if split > 0 && n > split {
			n = split
		}
		if !first && gap > 0 {
			time.Sleep(gap)
		}
		var t time.Time
		if limit > 0 {
			t = time.Now().Add(limit)
		}
		sc.conn.SetWriteDeadline(t)
		if _, err := sc.conn.Write(data[:n]); err != nil {
			return err
		}
		data = data[n:]
	}
	return nil
}

// simFrame encodes one frame: its length, type, seq, body, then the CRC-32 of the type,
// seq and body.
func simFrame(typ byte, seq uint32, body []byte) []byte {
	b := binary.BigEndian.AppendUint16(nil, uint16(9+len(body)))
	b = append(b, typ)
	b = binary.BigEndian.AppendUint32(b, seq)
	b = append(b, body...)
	return binary.BigEndian.AppendUint32(b, crc32.ChecksumIEEE(b[2:]))
}

// simNakFrame encodes a NAK of the frame with this seq: its code, then its reason.
func simNakFrame(seq uint32, code byte, reason string) []byte {
	return simFrame(simNak, seq, append([]byte{code}, reason...))
}

// simName is a frame type's name, as Before receives it.
func simName(typ byte) string {
	switch typ {
	case simHello:
		return "HELLO"
	case simPing:
		return "PING"
	case simQuery:
		return "QUERY"
	case simPoint:
		return "POINT"
	case simStow:
		return "STOW"
	case simUnstow:
		return "UNSTOW"
	case simJog:
		return "JOG"
	}
	return fmt.Sprintf("unknown type 0x%02x", typ)
}

// simDescribe is one frame as Received and Executed list it: "HELLO dishd", "PING",
// "119 POINT 183.25 41.70", "120 STOW", "122 JOG A +0.35".
func simDescribe(typ byte, seq uint32, body []byte) string {
	name := simName(typ)
	switch typ {
	case simHello:
		if len(body) == 0 {
			return name
		}
		return name + " " + string(body)
	case simPing, simQuery:
		if len(body) == 0 {
			return name
		}
		return fmt.Sprintf("%s % x", name, body)
	case simPoint:
		if len(body) == 4 {
			az := float64(binary.BigEndian.Uint16(body)) / 100
			el := float64(int16(binary.BigEndian.Uint16(body[2:]))) / 100
			return fmt.Sprintf("%d POINT %.2f %.2f", seq, az, el)
		}
	case simStow, simUnstow:
		if len(body) == 0 {
			return fmt.Sprintf("%d %s", seq, name)
		}
	case simJog:
		if len(body) == 3 {
			delta := float64(int16(binary.BigEndian.Uint16(body[1:]))) / 100
			return fmt.Sprintf("%d JOG %c %+.2f", seq, body[0], delta)
		}
	default:
		return name
	}
	return fmt.Sprintf("%d %s % x", seq, name, body)
}

// simASCII reports whether every byte of b is ASCII.
func simASCII(b []byte) bool {
	for _, c := range b {
		if c >= 0x80 {
			return false
		}
	}
	return true
}

// simPeerClosed reports whether err means the client closed its end of the connection.
func simPeerClosed(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, syscall.ECONNRESET)
}
