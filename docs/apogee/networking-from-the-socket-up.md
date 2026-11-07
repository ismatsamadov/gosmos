# Networking from the Socket Up: syllabus

**Authoring status:** all nine chapters are implemented as 78 lessons. Chapter 1, which also
wrote the shared simulator (`sim.go`, hidden and byte-identical in every lesson that uses it),
was read in the app in light and dark before the rest was written; chapters 2–9 were written
from this plan two chapters at a time and each pair independently reviewed. The lessons in
`courses/apogee.sqlite` are the source of truth, and the larger changes from this plan are
listed under "Changes in writing and review" at the end.

Season 4 of Apogee, March to June 2028, at Ystmark. The station's new dish controllers speak raw
bytes across an island line that drops in every storm; the learner writes dishd, which talks to
them, then a station server built on a bare socket. The world, the characters, the controller
protocol and the HTTP subset are in `story.md`, section "Season 4: down to the wire", which fixes
every byte, type, message and number used here.

- **Title:** Networking from the Socket Up
- **Slug:** `networking-from-the-socket-up` · **Kind:** Course
- **Short description:** Talk to Ystmark's new dish controllers byte by byte, then build a web server from a bare socket.
- **Description:** Ystmark's new dish controllers arrive with a manual that is wrong on page one, and Apogee drives them from Port Ellery across a line that drops in every storm. With Sunniva at the station and Nkem reviewing, you write the code that talks to them: connections and listeners, messages framed on a stream, big-endian fields and checksums, deadlines on every read and write, many connections at once, and reconnecting without sending a command twice. Then you build the station server from nothing but a socket: requests, bodies, chunked uploads and answers that Go's own client reads. Standard library only, and every connection stays on your Mac.

## At a glance

| # | Chapter | Slug | Go topics | Lessons |
|---|---|---|---|---|
| 1 | A Socket of Your Own | `a-socket-of-your-own` | Dial, read, write and close TCP connections, listen on 127.0.0.1 with a port the system picks, accept a client, and test connection code with net.Pipe. | 8 |
| 2 | Where One Message Ends | `where-one-message-ends` | Frame messages on a stream: delimited lines with bufio.Reader, exact byte counts with io.ReadFull, length prefixes, one buffered reader per connection and one write per frame. | 9 |
| 3 | Bytes in Order | `bytes-in-order` | Read and write big-endian fields with encoding/binary, signed and scaled values, a struct shaped like a frame, and CRC-32 checks with hash/crc32. | 9 |
| 4 | Deadlines on the Wire | `deadlines-on-the-wire` | Bound every read and write with deadlines, recognize timeouts, keep a slow line alive, send heartbeats and stop a blocked read when a context ends. | 8 |
| 5 | Both Dishes at Once | `both-dishes-at-once` | Run a goroutine per connection, serve many clients from one listener, write to all of them with deadlines and shut everything down cleanly. | 8 |
| 6 | When the Line Drops | `when-the-line-drops` | Reconnect with capped, jittered backoff, tell a lost line from a refusal, number commands and settle a command in flight so none runs twice. | 9 |
| 7 | Requests from Scratch | `requests-from-scratch` | Read an HTTP/1.1 request by hand: CR LF lines with limits, the request line, header fields, Host and Content-Length, and the right status for malformed input. | 9 |
| 8 | Bodies and Chunks | `bodies-and-chunks` | Read request bodies exactly: Content-Length, chunked transfer encoding and trailers, cut streams, size limits, binary data and a CRC that comes last. | 9 |
| 9 | Answers from Scratch | `answers-from-scratch` | Write HTTP/1.1 answers by hand: status lines, headers and bodies, HEAD, persistent connections, chunked streams and a clean shutdown, read by Go's own client. | 9 |

78 lessons: 4 code, 60 code_tests and 14 quizzes (82% hands-on). Calibration starts with chapter
1; later chapters are written only after it has been imported and read in the app. Difficulty
rises from 6 to 10. Every connection is a `net.Pipe` or a listener on 127.0.0.1 with a port the
system picks; every test sets deadlines, so a mistake fails within a second or two. Chapters 4
and 6 test time on the line inside `synctest.Test` over `net.Pipe`; clocks and random sources are
injected; no real sleep exceeds 50 ms and no test touches real learner files.

## Chapter 1 · A Socket of Your Own

`a-socket-of-your-own` · 8 lessons

Dial, read, write and close TCP connections, listen on 127.0.0.1 with a port the system picks, accept a client, and test connection code with net.Pipe.

#### 1.1 A Line from the Simulator
`a-line-from-the-simulator` · code · difficulty 6
- **Objective:** Open a TCP connection with net.Dial and read what the other end sends first.
- **Scenario:** Monday 6 March: the crate is open on the hut's bench, the maker's simulator is running, and Sunniva wants to know what a dish controller says when you connect.
- **Exercise:** `main` starts the simulator on a 127.0.0.1 listener, dials it and prints `n` bytes of `buf` with `%q`, but never reads; you call `conn.Read(buf)`.
- **Check:** `21 bytes: "YST-1 READY PROTO 1\r\n"`. The starter prints `0 bytes: ""`. The `\r` comes back in chapter 2.

#### 1.2 Hang Up When Done
`hang-up-when-done` · code_tests · difficulty 6
- **Objective:** Close every connection you open, with defer conn.Close() right after the error check.
- **Scenario:** After an afternoon of banner checks, the simulator's log shows forty connections dishd opened and never closed.
- **Exercise:** `bannerOf(addr string) (string, error)` dials, reads the banner with the given `readBannerLine` and returns, leaving the connection open; you defer its close after checking the dial's error.
- **Check:** Three calls return the banner and the simulator counts three connections closed by the client within a second; a just-closed listener's address returns its dial error without a panic. The mistakes: no close, and a defer before the error check, which panics on a nil connection.

#### 1.3 Quick Check: Connections
`quick-check-connections` · quiz · difficulty 6
- **Objective:** Recognize what Dial, Read and Close promise.
- **Scenario:** Nkem's three questions from Port Ellery before dishd talks to anything real.
- **Exercise:** Answer three questions.
- **Check:** A net.Conn is both an io.Reader and an io.Writer; Read returns the bytes that have arrived, which can be fewer than the buffer holds; a deferred Close belongs after the error check.

#### 1.4 Only on This Mac
`only-on-this-mac` · code_tests · difficulty 6
- **Objective:** Listen on 127.0.0.1 with port 0, so the system picks a free port and nothing outside the Mac can connect.
- **Scenario:** dishd will have a watch port, and Sunniva wants to try one on the bench computer without opening anything to the station's network.
- **Exercise:** `listenLocal() (net.Listener, error)` listens on the fixed address `127.0.0.1:7300`; you let the system pick the port.
- **Check:** The address's host is 127.0.0.1 and its port is neither 0 nor 7300; two listeners at once get two different ports; a client connects to each. The mistake: a fixed port, which the second listener can't have; the readme's callout warns against `:0`, which listens on every network.

#### 1.5 One Watcher, One Line
`one-watcher-one-line` · code_tests · difficulty 6
- **Objective:** Accept a connection from a listener, write to it and close it.
- **Scenario:** The first thing dishd's watch port will do is greet whoever connects with the dishes it knows.
- **Exercise:** `serveOne(ln net.Listener, names []string) error` returns nil at once; you accept one connection, write the given `greeting(names)` and close the connection.
- **Check:** A test client reads exactly `dishd watching YST-1 YST-2\n` and then io.EOF; one name gives `dishd watching YST-1\n`; a closed listener returns Accept's error, matching net.ErrClosed. The mistake: closing the listener but not the connection, so the watcher never sees the end.

#### 1.6 A Connection in Memory
`a-connection-in-memory` · code_tests · difficulty 6
- **Objective:** Make two connected ends in memory with net.Pipe and close them when the test ends with t.Cleanup.
- **Scenario:** Nkem wants dishd's tests to run without a single port, on any laptop, and asks for one helper every test can use.
- **Exercise:** `pipePair(t testing.TB) (net.Conn, net.Conn)` in `pipe.go` returns nil, nil; you return both ends of net.Pipe and register their closes with t.Cleanup.
- **Check:** Bytes written on one end are read on the other (the test writes from a goroutine, since a pipe has no buffer); a greeting sent through the pair arrives whole; after the subtest that made the pair ends, a write on either end fails with io.ErrClosedPipe. The mistake: defer inside the helper, which closes the ends before the test uses them.

#### 1.7 Nobody Listening
`nobody-listening` · code_tests · difficulty 6
- **Objective:** Recognize a refused connection with errors.Is and syscall.ECONNREFUSED, and report it plainly.
- **Scenario:** Dish 2's controller is still in the crate, and dishd's log says `connect: connection refused` where Sunniva wants `YST-2: not listening`.
- **Exercise:** `reach(name, addr string) string` returns the name and the error's text for any failure; you return `YST-2: not listening` for a refused dial and `YST-2: up` after a dial that works, closing that connection.
- **Check:** A live listener gives `YST-2: up`; a just-closed listener's address gives `YST-2: not listening`; `127.0.0.1` with no port gives `YST-2: dial tcp: address 127.0.0.1: missing port in address`. The mistake: comparing the error's text instead of using errors.Is.

#### 1.8 Hello, YST-1
`hello-yst-1` · code_tests · difficulty 7
- **Objective:** Combine a dial with a timeout, a first exchange and a close in one function.
- **Scenario:** Friday 10 March: before Dish 1's old controller comes out on the 20th, Sunniva wants dishd's first real conversation with the simulator.
- **Exercise:** `firstContact(addr string) (string, error)` returns an empty string; you dial with net.DialTimeout and one second, defer the close, call the given `handshake` and return `YST-1 firmware 2.4, last command 0`.
- **Check:** The simulator as YST-1 gives `YST-1 firmware 2.4, last command 0`; as YST-2 with last_seq 57, `YST-2 firmware 2.4, last command 57`; a refused address gives an error; the simulator saw the connection closed afterwards. The mistake: naming the dish from the address instead of the banner.

## Chapter 2 · Where One Message Ends

`where-one-message-ends` · 9 lessons

Frame messages on a stream: delimited lines with bufio.Reader, exact byte counts with io.ReadFull, length prefixes, one buffered reader per connection and one write per frame.

#### 2.1 Page One Is Wrong
`page-one-is-wrong` · code_tests · difficulty 6
- **Objective:** Read a delimited line with bufio.Reader.ReadString and remove its line ending exactly.
- **Scenario:** The manual's page one says the banner ends in a newline, so dishd compares `PROTO 1` with `PROTO 1\r` and refuses a controller that is fine.
- **Exercise:** `readBanner(br *bufio.Reader) (string, error)` trims only the `\n`; you also remove a `\r` before it, check the shape `<name> READY PROTO 1` and return the name.
- **Check:** `YST-1 READY PROTO 1\r\n` gives `YST-1`; the manual's `\n` alone also gives `YST-1`; `PROTO 2` gives `unsupported protocol "2"`; `YST-1 BUSY\r\n` and a stream that ends mid-line give errors. The mistake: trimming `\n` alone.

#### 2.2 All Eleven Bytes
`all-eleven-bytes` · code_tests · difficulty 7
- **Objective:** Read exactly n bytes with io.ReadFull, however the stream splits them.
- **Scenario:** Over Sunniva's test line to Port Ellery the simulator's 11-byte PONG arrives in pieces, and dishd's single Read takes the first three bytes for the whole answer.
- **Exercise:** `readExactly(r io.Reader, n int) ([]byte, error)` makes one Read; you use io.ReadFull.
- **Check:** iotest.OneByteReader and iotest.HalfReader over a PONG give all 11 bytes; a stream ending after 5 of 11 gives io.ErrUnexpectedEOF; an empty one gives io.EOF; the simulator with `Split: 3` over net.Pipe gives the PONG's bytes. The mistake: one Read, which returns 3.

#### 2.3 The Length Comes First
`the-length-comes-first` · code_tests · difficulty 7
- **Objective:** Read a length-prefixed frame: two big-endian bytes of length, then exactly that many bytes.
- **Scenario:** Every controller frame starts with its length, because binary data can hold any byte, a newline included.
- **Exercise:** `readPayload(br *bufio.Reader) ([]byte, error)` returns nil; you read the length with io.ReadFull and binary.BigEndian.Uint16, then exactly that many bytes.
- **Check:** A PONG gives its 9-byte payload; a WELCOME whose last_seq is 10, with a 0x0a in its middle, gives its 16-byte payload whole; one byte at a time gives the same; a stream that ends inside a frame gives io.ErrUnexpectedEOF; a stream that ends between frames gives io.EOF. The mistake: ReadBytes('\n'), which stops at the 0x0a.

#### 2.4 Quick Check: Byte Streams
`quick-check-byte-streams` · quiz · difficulty 6
- **Objective:** Recognize that TCP carries a stream of bytes with no edges between messages.
- **Scenario:** Nkem asks three questions before dishd reads its first real frame.
- **Exercise:** Answer three questions.
- **Check:** One Write can arrive as several Reads and two Writes as one; io.ReadFull keeps reading until it has n bytes or an error; a delimiter suits text, while binary frames need a length prefix.

#### 2.5 One Reader per Connection
`one-reader-per-connection` · code_tests · difficulty 7
- **Objective:** Read a whole connection through one bufio.Reader, because the bytes it has buffered are no longer in the connection.
- **Scenario:** To save a round trip on the backhaul, dishd sends HELLO and QUERY together; the WELCOME and the STATUS come back glued in one read, and in Wednesday's test dishd lost the STATUS in which the simulated Dish 1 had stowed itself.
- **Exercise:** `newClient` reads the banner through `c.br`, but `(c *Client) recv() (Frame, error)` makes a new bufio.Reader on the connection for every frame; you read every frame through `c.br`.
- **Check:** The simulator with `Glue: 2` gives the banner alone, then WELCOME and STATUS glued: recv returns both in order; with `HoldBanner` as well, HELLO and QUERY sent before the banner is read, the name, WELCOME and STATUS all come through `c.br`; with `Split: 1` the same; a third recv after a PING gives PONG. The mistake: a reader per frame, which leaves the STATUS inside the first reader's buffer until the one-second guard.

#### 2.6 Lengths You Can Trust
`lengths-you-can-trust` · code_tests · difficulty 7
- **Objective:** Refuse a length prefix outside the protocol's range before reading or allocating anything.
- **Scenario:** During Sunniva's cable test a damaged length read `ff ff`, and dishd sat waiting for 65,535 bytes that would never come.
- **Exercise:** `readPayload` reads any length it is given; you return ErrBadLength, wrapped with the length, for a length under 9 or over 1024, without reading further.
- **Check:** Lengths 9 and 1024 are read; 8, 1025, 0 and 65535 give ErrBadLength at once, though nothing follows the two bytes, with the message `bad length 65535`. The mistake: checking only the upper bound.

#### 2.7 One Write per Frame
`one-write-per-frame` · code_tests · difficulty 7
- **Objective:** Build a whole frame in memory and send it with a single Write, returning its error.
- **Scenario:** Two of dishd's goroutines sent frames at the same moment, and the controller read the first half of one glued to the second half of the other.
- **Exercise:** `writeFrame(w io.Writer, f Frame) error` writes the length, then the rest, in two Writes, and ignores errors; you write the given `encodeFrame(f)` once and return the error.
- **Check:** A recording writer sees exactly one Write per frame, PING as `00 09 02 00 00 00 00 bc e2 a4 7d`; three frames make three Writes; a failing writer's error is returned. The readme says why: a connection writes each call whole, even when goroutines share it. The mistake: the length and the rest in separate Writes, between which another goroutine's frame can land.

#### 2.8 Quick Check: Frames
`quick-check-frames` · quiz · difficulty 7
- **Objective:** Reason about buffered readers and frame boundaries.
- **Scenario:** Mirela reads the frame reader from Port Ellery and asks three questions.
- **Exercise:** Answer three questions.
- **Check:** A bufio.Reader may read more from the connection than you asked for; once wrapped, every read must go through it; a length is checked before its bytes are read.

#### 2.9 A Whole Conversation
`a-whole-conversation` · code_tests · difficulty 7
- **Objective:** Send several frames, then read one answer per frame, in order, through one client.
- **Scenario:** Friday 17 March: dishd's first full session with the simulator over Sunniva's slow test line.
- **Exercise:** `(c *Client) exchange(frames ...Frame) ([]Frame, error)` returns nil; you send every frame with writeFrame, then read one answer for each, returning the answers read and the first error.
- **Check:** HELLO, QUERY and PING give WELCOME, STATUS and PONG in order, with the simulator writing one byte at a time and with all three answers glued; a BYE `idle` in place of the second answer gives one answer and a ByeError reading `controller closed: idle`; a simulator closing after two answers gives two answers and io.EOF. The mistake: returning the BYE as if it were an answer.

## Chapter 3 · Bytes in Order

`bytes-in-order` · 9 lessons

Read and write big-endian fields with encoding/binary, signed and scaled values, a struct shaped like a frame, and CRC-32 checks with hash/crc32.

#### 3.1 Dish 1 at 382 Degrees
`dish-1-at-382-degrees` · code · difficulty 7
- **Objective:** Read multi-byte fields in the right order: big-endian, the most significant byte first.
- **Scenario:** Monday 27 March: Dish 1's first status frame from its new controller, read as page 14 of the manual says, puts the dish at 382.15° in a 3,277 km/h wind.
- **Exercise:** `main` decodes Monday's 16-byte status body, a constant, with binary.LittleEndian and prints azimuth, elevation and wind; you read it with binary.BigEndian.
- **Check:** `az 183.25` / `el 41.70` / `wind 64.0`. The starter prints `az 382.15` / `el 189.60` / `wind 3277.0`.

#### 3.2 Below Zero in the Radome
`below-zero-in-the-radome` · code_tests · difficulty 7
- **Objective:** Read a signed 16-bit field by converting the uint16 to int16, then scale it.
- **Scenario:** dishd's screen says Dish 1's radome is at 6,541.1 °C; inside, Sunniva can see her breath at −12.5 °C.
- **Exercise:** `tempC(body []byte) float64` divides the unsigned bytes 8–9 by 10; you convert them through int16 first.
- **Check:** `ff 83` gives −12.5; `00 7d` 12.5; `00 00` 0; `ff e1` −3.1; `80 00` −3276.8; `7f ff` 3276.7. The mistakes: no conversion (6541.1) and integer division before converting to float (−12).

#### 3.3 Quick Check: Byte Order
`quick-check-byte-order` · quiz · difficulty 7
- **Objective:** Check big-endian order, single bytes and signed conversions.
- **Scenario:** Tomás, off the plane on Tuesday 28 March, asks three questions over coffee in the hut.
- **Exercise:** Answer three questions.
- **Check:** Big-endian puts the most significant byte first, so `47 95` is 18325; a one-byte field has no byte order; int16 of 0xff83 is −125.

#### 3.4 A Struct Like the Frame
`a-struct-like-the-frame` · code_tests · difficulty 7
- **Objective:** Decode a fixed binary layout into a struct with binary.Decode, check the length first, then convert units.
- **Scenario:** Sunniva wants every field of the status frame on her screen, not only the azimuth and the temperature.
- **Exercise:** `rawStatus` mirrors the layout with fixed-size fields; `decodeStatus(body []byte) (Status, error)` returns a zero Status; you require 16 bytes, decode with binary.Decode and big-endian order, and fill Status in degrees, km/h and °C.
- **Check:** Monday's body gives dish 1, TRACKING, 183.25, 41.70, 64.0, −12.5, the heater flag and last_seq 118; Dish 2 stowed in the wind gives dish 2, STOWED, 183.25, 90.00, 81.0, −3.1, heater and wind lock, 121; 15 bytes give `status body is 15 bytes, want 16`, and 17 bytes `status body is 17 bytes, want 16`. The mistake: leaving the length to binary.Decode, which accepts 17 bytes.

#### 3.5 Point Dish 1
`point-dish-1` · code_tests · difficulty 7
- **Objective:** Encode big-endian values with binary.BigEndian.AppendUint16, rounding to the protocol's hundredths and refusing values out of range.
- **Scenario:** Tomás asks for Dish 1 at azimuth 183.25°, elevation 41.70°, to check that a POINT and the next STATUS agree.
- **Exercise:** `pointBody(azDeg, elDeg float64) ([]byte, error)` returns nil; you round to hundredths with math.Round, require 0 to 359.99 and 0 to 90.00 after rounding, and append the azimuth as uint16 and the elevation as int16.
- **Check:** (183.25, 41.70) gives `47 95 10 4a`; (0, 0) `00 00 00 00`; (359.99, 90) `8c 9f 23 28`; (1.15, 0.57) `00 73 00 39`; 360, 359.996 (which rounds to 360.00), −0.01 and 90.01 give errors. The mistake: converting without rounding, which turns 1.15 into 114 (the test passes values as variables, not constants).

#### 3.6 The CRC at the End
`the-crc-at-the-end` · code_tests · difficulty 8
- **Objective:** Compute a CRC-32 with crc32.ChecksumIEEE over exactly the bytes the protocol names, and append it big-endian.
- **Scenario:** The simulator answers dishd's first real POINT with NAK `bad crc`.
- **Exercise:** `encodeFrame(f Frame) []byte` in `encode.go` builds the length, type, seq and body, then a CRC of zero; you compute the CRC of the type, seq and body and append it.
- **Check:** PING gives `00 09 02 00 00 00 00 bc e2 a4 7d`; POINT seq 119 gives `00 0d 10 00 00 00 77 47 95 10 4a 20 eb ca 40`; HELLO `dishd` gives its 16 bytes; a nil and an empty body give the same frame; the simulator ACKs the POINT. The mistakes: a CRC that also covers the length (`c0 a3 dd d6`) and one written little-endian (`7d a4 e2 bc`).

#### 3.7 Refuse a Damaged Frame
`refuse-a-damaged-frame` · code_tests · difficulty 8
- **Objective:** Check a frame's CRC on arrival and refuse the frame before trusting any field.
- **Scenario:** Sunniva's cable test flipped one bit of a STATUS, and dishd showed Dish 1 at 203.73° while it sat at 183.25°.
- **Exercise:** `decodeFrame(payload []byte) (Frame, error)` in `decode.go` splits the payload without checking; you refuse a payload under 9 bytes with ErrBadLength and compare its last four bytes with the CRC of the rest, returning ErrBadCRC.
- **Check:** Monday's STATUS payload gives type 0x83, seq 0 and its 16-byte body; that payload with any single bit flipped, each of its 200 bits in turn, gives ErrBadCRC; 8 bytes give ErrBadLength; a NAK gives its seq and body. The mistake: decoding the fields first and checking the CRC after.

#### 3.8 Quick Check: Checksums
`quick-check-checksums` · quiz · difficulty 7
- **Objective:** Know what a CRC protects against and what it doesn't.
- **Scenario:** Mirela reads the frame code and asks three questions.
- **Exercise:** Answer three questions.
- **Check:** A CRC catches accidental damage but not a deliberate change, since anyone can compute it again (signatures, as on Pipit-2's commands, do that); the check comes before any field is used; crc32.ChecksumIEEE is the same CRC-32 that downlinked files carry.

#### 3.9 Dish 1 Points True
`dish-1-points-true` · code_tests · difficulty 8
- **Objective:** Send a command and handle every kind of answer: ACK, NAK and BYE.
- **Scenario:** Friday 31 March, 06:50: Cirrine-1's 07:02 pass will be Dish 1's first customer pass on the new controller, and the go/no-go poll needs the dish pointed.
- **Exercise:** `(c *Client) command(f Frame) error` sends the frame and returns nil; you read the answer: an ACK of the same seq is nil, a NAK is the given `nakError(ans)`, a BYE its ByeError, anything else an error naming the type and seq.
- **Check:** POINT 119 gives nil, and the next status says 183.25° and 41.70°; a stowed simulator gives a NakError (errors.AsType) with code 4 and `seq 119 refused: stowed`; wind at 81 km/h gives code 3, `wind lock`; an ACK for seq 118 gives an error; a damaged answer gives ErrBadCRC. The mistake: taking any ACK as this command's.

## Chapter 4 · Deadlines on the Wire

`deadlines-on-the-wire` · 8 lessons

Bound every read and write with deadlines, recognize timeouts, keep a slow line alive, send heartbeats and stop a blocked read when a context ends.

#### 4.1 Never Wait Forever
`never-wait-forever` · code_tests · difficulty 7
- **Objective:** Put a deadline on every read with SetReadDeadline, so a silent peer can't hold you.
- **Scenario:** Tuesday 11 April, 03:10: the storm cut the fiber, and dishd waited forty minutes on an answer from Dish 1 that was never coming.
- **Exercise:** `(c *Client) recv() (Frame, error)` reads with no deadline; you set the read deadline to now plus `c.Timeout` before reading.
- **Check:** A simulator silent after its banner, with Timeout 50 ms, gives an error within a second; an answering simulator gives the frame; five exchanges 50 ms apart with a 200 ms timeout all succeed. The mistake: one deadline set in newClient, an instant that has passed by the third exchange.

#### 4.2 Is It a Timeout?
`is-it-a-timeout` · code_tests · difficulty 7
- **Objective:** Recognize a deadline with errors.Is and os.ErrDeadlineExceeded, and keep other failures as they are.
- **Scenario:** Ingrid's log must say `YST-1: no answer in 5s` when the line goes quiet, and something exact for everything else.
- **Exercise:** `describe(name string, timeout time.Duration, err error) string` returns the name and the error's text; you return `no answer in <timeout>` for a deadline and `closed` for io.EOF and io.ErrUnexpectedEOF.
- **Check:** A real deadline error from a net.Pipe read with 50 ms gives `YST-1: no answer in 50ms`, wrapped or not; io.EOF gives `YST-1: closed`; ErrBadCRC gives `YST-1: bad crc`; context.DeadlineExceeded gives `YST-1: context deadline exceeded`. The mistake: searching the error's text for "timeout".

#### 4.3 Quick Check: Deadlines
`quick-check-deadlines` · quiz · difficulty 7
- **Objective:** Check what a deadline is and how it behaves.
- **Scenario:** Sunniva's three questions, telegraph style, over the backhaul.
- **Exercise:** Answer three questions.
- **Check:** A deadline is an instant, not a duration, and stays until changed; SetDeadline with the zero time removes it; after a timeout the connection works again once a later deadline is set.

#### 4.4 Slow Is Not Dead
`slow-is-not-dead` · code_tests · difficulty 8
- **Objective:** Use an idle deadline that moves forward before each read, so a slow but living line is never cut.
- **Scenario:** Over the backhaul a 27-byte STATUS trickles in over ten seconds, and a deadline for the whole frame cuts every one.
- **Exercise:** `idleReader` wraps the client's connection and its Read only reads; you set the read deadline to now plus the client's Timeout before every Read, so `c.br` reads through it.
- **Check:** Inside synctest.Test over net.Pipe with a 5 s timeout: a STATUS at one byte every 0.4 s (10.4 s in all) arrives whole; a byte every 6 s gives os.ErrDeadlineExceeded at 5 s; a silent simulator gives the deadline at exactly 5 s of fake time. The mistake: one deadline for the whole frame, which cuts the slow STATUS at 5 s.

#### 4.5 Writes Wait Too
`writes-wait-too` · code_tests · difficulty 8
- **Objective:** Put a deadline on writes with SetWriteDeadline, because a peer that stops reading blocks your Write.
- **Scenario:** During the storm a POINT sat in dishd's Write for six minutes, its bytes going nowhere.
- **Exercise:** `(c *Client) send(f Frame) error` writes with no deadline; you set a write deadline of `c.Timeout` first.
- **Check:** Over net.Pipe a simulator with `NoRead` gives os.ErrDeadlineExceeded within a second at 50 ms; a reading simulator gets the exact bytes; five sends 50 ms apart with a 200 ms timeout succeed. The mistake: a read deadline where the write needs one.

#### 4.6 Keep the Line Warm
`keep-the-line-warm` · code_tests · difficulty 8
- **Objective:** Send a heartbeat on every tick and treat a missing answer as a dead line.
- **Scenario:** Over the storm link a dead connection looks exactly like a quiet one until the next command, maybe an hour away, and the controller hangs up on a client silent for 30 s.
- **Exercise:** `heartbeat(c *Client, ticks <-chan time.Time) error` returns nil at once; on each tick you send PING and require a PONG, return the first error, and return nil when the ticks channel closes.
- **Check:** Three ticks put three PINGs in the simulator's Received list and return nil; with `MuteAfter: 2` the second tick gives a deadline error; with no ticks and the controller's IdleLimit at 50 ms, the next recv gives `controller closed: idle`. The mistake: sending PING without reading the PONG, which leaves it for the next command to read.

#### 4.7 Stop When the Pass Ends
`stop-when-the-pass-ends` · code_tests · difficulty 8
- **Objective:** Make a blocked read return when a context ends, with context.AfterFunc setting the read deadline to now.
- **Scenario:** At LOS a pass is over, and dishd must stop waiting on Dish 1 at once, not five seconds later.
- **Exercise:** `recvCtx(ctx context.Context, c *Client) (Frame, error)` reads with only the client's timeout; you register context.AfterFunc to set the deadline to now, stop it when the read returns, and return ctx.Err() when the context ended the read.
- **Check:** In a bubble, a silent simulator, a 5 s timeout and a context ending after 20 ms give context.DeadlineExceeded at 20 ms; an already cancelled context gives context.Canceled at once; an answering simulator gives the frame, and the next read still works after that context is cancelled. The mistake: never calling the stop function, so a late cancel cuts a later read.

#### 4.8 Through the Storm
`through-the-storm` · code_tests · difficulty 8
- **Objective:** Combine write and idle deadlines with a context in one call that fails fast and says why.
- **Scenario:** Friday 14 April: Sunniva replays Tuesday's storm on the simulator: a slow line, a silent one, and one that stops reading.
- **Exercise:** `(c *Client) Do(ctx context.Context, f Frame) (Frame, error)` returns an error; you send with the write deadline and receive with recvCtx.
- **Check:** In a bubble with a 5 s timeout: the slow line gives the STATUS; the silent line gives an error that describe turns into `YST-1: no answer in 5s`, at 5 s; `NoRead` gives a write deadline error; a context ending at 2 s gives context.DeadlineExceeded at 2 s; a calm line answers at once. The mistake: reading the answer when the send has failed.

## Chapter 5 · Both Dishes at Once

`both-dishes-at-once` · 8 lessons

Run a goroutine per connection, serve many clients from one listener, write to all of them with deadlines and shut everything down cleanly.

#### 5.1 One Goroutine per Dish
`one-goroutine-per-dish` · code_tests · difficulty 8
- **Objective:** Run each connection's exchange in its own goroutine and wait for all of them with sync.WaitGroup.
- **Scenario:** Dish 2's controller went in last week, and dishd asks Dish 1, then Dish 2, so over the backhaul every screen is a second and a half behind.
- **Exercise:** `pollAll(clients []*Client) []Result` polls one client after the other; you poll each in its own goroutine with wg.Go and store each result at its index.
- **Check:** Two simulators whose Before hook answers a QUERY only once both have been asked give both statuses, in client order, within a second; one client gives its result; none gives an empty slice. The mistake: polling in turn, which the barrier holds until the guard fires.

#### 5.2 One Dish Down
`one-dish-down` · code_tests · difficulty 8
- **Objective:** Keep a result or an error per connection, so one failure never hides the other's answer.
- **Scenario:** Dish 2's controller rebooted for a firmware check, and dishd's poll reported nothing for either dish.
- **Exercise:** `pollAll` stops at the first error and returns nil; you record each dish's Status or Err in its own Result.
- **Check:** Dish 1 answering and Dish 2's simulator closed give Result 0 with its status and Result 1 with an error; both failing give two errors; the given watchLine turns them into `YST-1 TRACKING az=183.25 el=41.70 wind=64` and `YST-2 DOWN closed`. The mistake: stopping at the first error, which drops Dish 1's good status.

#### 5.3 Many Watchers at Once
`many-watchers-at-once` · code_tests · difficulty 8
- **Objective:** Serve many connections at once with an accept loop that starts a goroutine per connection and ends cleanly when its listener closes.
- **Scenario:** Sunniva's terminal in the hut, Ingrid's screen in the Fishbowl and Tomás's laptop all want the watch port, and the second waits until the first hangs up.
- **Exercise:** `serve(ln net.Listener, wg *sync.WaitGroup, handle func(net.Conn)) error` handles each connection before accepting the next; you start each handler with wg.Go and keep accepting, returning nil once the listener is closed and any other Accept error as it is.
- **Check:** Three clients connected at once each get the greeting while the test's handlers stay blocked until released; closing the listener returns nil, and wg.Wait returns once the handlers are released. The mistake: returning the closed listener's error as a failure.

#### 5.4 Every Line to Every Watcher
`every-line-to-every-watcher` · code_tests · difficulty 8
- **Objective:** Write the same line to every connection in a set kept under a mutex, and drop the ones that fail.
- **Scenario:** Every second each dish's line must reach everyone on the watch port.
- **Exercise:** `(w *Watchers) Broadcast(line string)` writes to nobody; you copy the set under the lock, write the line and a newline to each watcher outside it, and Remove each watcher whose write fails.
- **Check:** Three watchers each read `YST-1 TRACKING az=183.25 el=41.70 wind=64\n`; a watcher that has hung up is removed while the others get ten broadcasts in a row; with no watchers nothing happens; Len reports the set. The mistake: returning at the first failed write, so the watchers after it miss the line.

#### 5.5 A Slow Watcher Waits Alone
`a-slow-watcher-waits-alone` · code_tests · difficulty 9
- **Objective:** Bound every write to a watcher with a write deadline, and drop a watcher that can't keep up.
- **Scenario:** Tomás's laptop went to sleep with the watch port open, and every other screen froze with it.
- **Exercise:** Broadcast writes with no deadline; you set a write deadline of `w.Timeout` before each write and Remove the watcher whose write fails.
- **Check:** Over net.Pipe, with one watcher that never reads and two that do and Timeout 50 ms, Broadcast returns within a second, both readers have the line, and the sleeper is removed and closed; the next broadcast reaches the two at once. The mistake: one deadline set when the watcher joined, which soon expires for everyone.

#### 5.6 Close Everything Before Leaving
`close-everything-before-leaving` · code_tests · difficulty 9
- **Objective:** Shut down a server: close the listener so Accept returns, close every connection so blocked reads return, then wait for every goroutine.
- **Scenario:** Friday 28 April's deploy restarts dishd, and the old process hangs on, holding both controllers and three watchers.
- **Exercise:** `(s *Service) Shutdown()` closes only the listener; you also close every watcher and both dish clients, then wait on `s.wg`, all inside `s.once`.
- **Check:** With two simulators and three watchers connected, Shutdown returns within a second; each watcher reads io.EOF; each simulator counts its connection closed by the client; the accept loop has returned; a second Shutdown returns at once. The mistake: waiting before closing, so the watchers' read loops never end.

#### 5.7 Quick Check: Shutting Down
`quick-check-shutting-down` · quiz · difficulty 8
- **Objective:** Reason about the order of a clean stop.
- **Scenario:** Nkem's three questions before Friday's deploy.
- **Exercise:** Answer three questions.
- **Check:** Closing a listener makes a blocked Accept return an error matching net.ErrClosed; closing a connection makes a blocked Read on it return; Wait comes after every close, or it waits forever.

#### 5.8 Both Dishes, One Service
`both-dishes-one-service` · code_tests · difficulty 9
- **Objective:** Run both controllers and the watch port in one service: a goroutine per dish poll, a goroutine per watcher, one clean shutdown.
- **Scenario:** Friday 28 April at midnight, the sun still up: dishd runs Dish 1 and Dish 2 together for the first time, and Sunniva watches both from the hut.
- **Exercise:** `(s *Service) Run(ctx context.Context, ticks <-chan time.Time) error` returns nil at once; you start the watch port, then on each tick poll both dishes and broadcast each watchLine, and when the context ends shut down and return nil.
- **Check:** Two simulators, two watchers and three ticks give each watcher its greeting and six lines, YST-1 before YST-2 at each tick; Dish 2's simulator closed after the first tick gives `YST-2 DOWN closed` while Dish 1's lines go on; a cancel makes Run return nil within a second with every connection closed. The mistake: returning on cancel without shutting down.

## Chapter 6 · When the Line Drops

`when-the-line-drops` · 9 lessons

Reconnect with capped, jittered backoff, tell a lost line from a refusal, number commands and settle a command in flight so none runs twice.

#### 6.1 Dial Again
`dial-again` · code_tests · difficulty 8
- **Objective:** Reconnect in a loop until a dial succeeds or the context ends.
- **Scenario:** Tuesday 9 May, 07:36: the line was back and dishd wasn't; it had given up at 07:22 and never dialed again.
- **Exercise:** `(l *Link) redial(ctx context.Context) (net.Conn, error)` dials once; you dial again after `l.Wait(ctx, time.Second)` until a dial works, returning the wait's error when the context ends.
- **Check:** A fake line that refuses three dials and then connects gives the connection after four dials and three 1 s waits; a context cancelled during a wait gives context.Canceled and no further dial; a dial that works at once waits for nothing. The mistake: returning the last dial's error after the context ended, so the caller can't tell stop from down.

#### 6.2 Back Off, Then Try Again
`back-off-then-try-again` · code_tests · difficulty 8
- **Objective:** Space reconnection attempts with a capped, jittered backoff, and start again from the first delay after a good connection.
- **Scenario:** Dialing every second through fourteen minutes is 840 attempts over a line Sunniva needs for everything else.
- **Exercise:** The given `redialDelay(attempt, r)` is the bible's backoff; redial waits 1 s every time; you count failed dials in `l.attempt`, wait `redialDelay(l.attempt, l.Rand)`, and set `l.attempt` back to 0 when a dial succeeds.
- **Check:** With rand.New(rand.NewPCG(1, 2)), six refused dials then success give waits in 0.5–1 s, 1–2, 2–4, 4–8, 8–16 and 15–30 s, equal to the test's own computation from the same seed; after that success and a new drop the first wait is 0.5–1 s again; sixty refused dials never wait over 30 s. The mistake: never resetting, so the next storm starts at 30 s.

#### 6.3 Quick Check: Reconnecting
`quick-check-reconnecting` · quiz · difficulty 8
- **Objective:** Reason about backoff, its cap and when to stop trying.
- **Scenario:** Ingrid's questions before dishd is trusted with a storm alone.
- **Exercise:** Answer three questions.
- **Check:** Jitter keeps the two dishes' links from redialing in step; the cap keeps a fourteen-minute outage from turning into a long wait after the line returns; a cancelled context means stop, not try again.

#### 6.4 Lost or Just Refused?
`lost-or-just-refused` · code_tests · difficulty 8
- **Objective:** Sort errors into a lost line, which means reconnect, and everything else, which doesn't.
- **Scenario:** In the storm, dishd took a NAK `wind lock` for a dead line and redialed forty times, while the controller had answered perfectly well.
- **Exercise:** `lost(err error) bool` returns `err != nil`; you return true for the bible's list of lost-line errors, found with errors.Is and errors.AsType, and false for the rest.
- **Check:** io.EOF, io.ErrUnexpectedEOF, a deadline, net.ErrClosed, io.ErrClosedPipe, ECONNRESET, EPIPE, a ByeError, ErrBadLength and ErrBadCRC, plain and wrapped with %w, give true; a NakError `wind lock`, context.Canceled, context.DeadlineExceeded and nil give false. The mistake: `err != nil`, which redials on a refusal and on a cancel.

#### 6.5 Number Every Command
`number-every-command` · code_tests · difficulty 8
- **Objective:** Give each command a strictly increasing sequence number, continuing from the controller's own count after every HELLO.
- **Scenario:** Tomás asks the question Nkem asked about Pipit-1 a year ago: what does the dish see if this line runs twice?
- **Exercise:** `(l *Link) number(f Frame) Frame` returns the frame unchanged and `(l *Link) learn(w Welcome)` does nothing; you set the frame's Seq from `l.next` and count up, and make learn set `l.next` to the larger of itself and `w.LastSeq + 1`.
- **Check:** After a WELCOME with last_seq 118, three commands get 119, 120 and 121; a later WELCOME with 200 makes the next 201; one with 57 when the link is at 122 leaves it at 122; a fresh link after last_seq 0 starts at 1. The mistake: `l.next = w.LastSeq + 1`, which goes back to 58.

#### 6.6 The Command in Flight
`the-command-in-flight` · code_tests · difficulty 9
- **Objective:** After reconnecting, decide from the controller's count whether the command in flight ran, and send it again only if it didn't, under its own number.
- **Scenario:** At 07:22:04 Tomás's jog left dishd, the line dropped two seconds later, and nobody knew whether Dish 2 had moved.
- **Exercise:** `(l *Link) resume(w Welcome) *Frame` always returns `l.inflight`; you return nil and clear it when `w.LastSeq` is at or past its seq.
- **Check:** JOG 122 in flight with last_seq 122 gives nil (it ran); with 121, the JOG, still seq 122; with 150, nil; nothing in flight, nil; against the simulator with `DropAfter: 122`, then with `DropBefore: 122`, the jog is executed exactly once. The mistake: always sending it again, or under a new number, which the simulator executes twice.

#### 6.7 Quick Check: Never Twice
`quick-check-never-twice` · quiz · difficulty 8
- **Objective:** Reason about lost answers, sequence numbers and commands that must not repeat.
- **Scenario:** Mirela, who stopped a replayed command to Cirrine-1 last year, asks three questions.
- **Exercise:** Answer three questions.
- **Check:** A lost ACK doesn't mean a lost command; a command sent again keeps its number, so the controller can refuse a repeat; a POINT run twice ends in the same place, a JOG run twice doesn't.

#### 6.8 Send Through the Drop
`send-through-the-drop` · code_tests · difficulty 9
- **Objective:** Send a command so that a dropped line costs a reconnect, never a second execution.
- **Scenario:** The pass plan's STOW for Dish 2 left dishd at 07:21:58, and the old dishd lost it with the line.
- **Exercise:** `(l *Link) Send(ctx context.Context, cmd Frame) error` numbers the command and sends it once with Do; you keep it in flight and, on a lost line, close the client, call the given `connect`, return nil if resume says it ran, or send it again, letting `settle` judge the answer.
- **Check:** The simulator dropping after executing the STOW gives nil, one execution and two connections; dropping before executing gives nil and one execution; a NAK `stowed` gives the NakError with no reconnect; a context cancelled while the line is down gives the context's error. The mistake: numbering the command again on every attempt.

#### 6.9 Fourteen Minutes Without the Mainland
`fourteen-minutes-without-the-mainland` · code_tests · difficulty 9
- **Objective:** Run a queue of commands through a long outage: back off, reconnect, settle the command in flight and finish the queue in order.
- **Scenario:** Friday 12 May, 10:00: Sunniva pulls the island's line in the middle of Tomás's jog and keeps it out for fourteen minutes.
- **Exercise:** `(l *Link) Drain(ctx context.Context, queue []Frame) ([]error, error)` sends only the first command; you send each in order through Send, record each outcome (nil or a NakError), and stop early only on any other error.
- **Check:** In a bubble: STOW, UNSTOW, POINT and JOG, the jog executed and its answer lost as the line drops, then every dial refused for fourteen minutes of fake time, give four nil outcomes, each command once in order in the simulator's Executed, every wait at most 30 s and the waits adding up to at least 14 minutes; with the wind at 81 km/h, UNSTOW and POINT come back `wind lock` and the queue still finishes; a context ending at minute 5 gives its error and the outcomes so far. The mistake: stopping the queue at the first refusal.

## Chapter 7 · Requests from Scratch

`requests-from-scratch` · 9 lessons

Read an HTTP/1.1 request by hand: CR LF lines with limits, the request line, header fields, Host and Content-Length, and the right status for malformed input.

#### 7.1 What Go Sends
`what-go-sends` · code · difficulty 8
- **Objective:** See an HTTP request as bytes on a socket: a request line, header lines and a blank line, each ending in CR LF.
- **Scenario:** Monday 22 May: the recorder and headquarters' collector are Go programs using net/http, and before writing the station server you watch what one of them sends.
- **Exercise:** `main` listens on 127.0.0.1, sends `GET /recordings` from an http.Client in a goroutine (its Transport dials the listener, so Host reads 127.0.0.1), and prints the first line it reads with `%q`; you print every line through the blank one, reading with a bufio.Reader. `main` then answers an empty 200 so the client finishes.
- **Check:** `"GET /recordings HTTP/1.1\r\n"` / `"Host: 127.0.0.1\r\n"` / `"User-Agent: Go-http-client/1.1\r\n"` / `"Accept-Encoding: gzip\r\n"` / `"\r\n"`. The starter prints only the first line.

#### 7.2 Three Parts, Two Spaces
`three-parts-two-spaces` · code_tests · difficulty 8
- **Objective:** Split a request line into method, target and version, and refuse every other shape.
- **Scenario:** The recorder's request lines are always tidy, but the station server will hear from anything on the hut's network.
- **Exercise:** `parseRequestLine(line string) (method, target string, err error)` returns the whole line as the method; you split on single spaces into three non-empty parts, require a target starting with `/` and a version shaped `HTTP/<digit>.<digit>` (400 `bad request line`), and answer any version but HTTP/1.1 with 505.
- **Check:** `GET /recordings HTTP/1.1` gives GET and /recordings; a PUT of `/recordings/furrow-a/YST/18856` gives its target; two spaces, a tab, a missing version, `GET recordings HTTP/1.1` and `HTTP/1` give 400; `HTTP/1.0` and `HTTP/2.0` give 505 `version "HTTP/1.0" not supported`; `get` is returned as the method, checked later. The mistake: strings.Fields, which accepts two spaces and a tab.

#### 7.3 Quick Check: Request Lines
`quick-check-request-lines` · quiz · difficulty 8
- **Objective:** Name the parts of a request line and what each refusal means.
- **Scenario:** Nkem's three questions, with Season 3's booking service in mind.
- **Exercise:** Answer three questions.
- **Check:** The target is the path as the client wrote it; 505 says the version is wrong, 400 that the line itself is; method names are case-sensitive, so `get` is not GET.

#### 7.4 Two Bytes End a Line
`two-bytes-end-a-line` · code_tests · difficulty 8
- **Objective:** Read one protocol line with bufio.Reader.ReadSlice, requiring CR LF and refusing a line longer than the reader's buffer.
- **Scenario:** Mirela's first rule for the station server: a line that doesn't end in CR LF, or never ends at all, is refused, never guessed at.
- **Exercise:** `readLine(br *bufio.Reader) (string, error)` uses ReadString and accepts any ending; on a reader made with bufio.NewReaderSize and 1024 bytes you use ReadSlice('\n'), return errLineTooLong for bufio.ErrBufferFull, 400 `bad line ending` when no `\r` comes before the `\n`, and the line without its CR LF, copied into a string.
- **Check:** `GET / HTTP/1.1\r\n` gives the line; 1022 characters and CR LF give the line; 1023 and CR LF give errLineTooLong; `Host: x\n` gives 400; a stream ending mid-line gives io.ErrUnexpectedEOF; an empty one gives io.EOF. The mistake: ReadString, which keeps reading a 64 KiB line with no end.

#### 7.5 Headers Until the Blank Line
`headers-until-the-blank-line` · code_tests · difficulty 8
- **Objective:** Read header fields up to the blank line: a name, a colon and a value trimmed of spaces, names in lower case.
- **Scenario:** Headquarters' collector says who it is and what it wants in its header fields, and the station server must read every one.
- **Exercise:** `readHeader(br *bufio.Reader) (Header, error)` returns an empty Header; you read lines until the empty one, split each at its first colon, lower-case the name, trim spaces and tabs around the value and append it to the name's values.
- **Check:** Go's GET head gives host 127.0.0.1, user-agent Go-http-client/1.1 and accept-encoding gzip; `Content-Type:   text/plain  ` gives text/plain; `X-A: 1` then `x-a: 2` gives `1, 2` from Get; `Host: 127.0.0.1:8080` keeps everything after the first colon; `X-Empty:` gives an empty value. The mistake: splitting at the last colon.

#### 7.6 Refuse a Crooked Header
`refuse-a-crooked-header` · code_tests · difficulty 9
- **Objective:** Refuse malformed header lines and too many fields with the right status.
- **Scenario:** Mirela's second rule: when two programs could read a header two ways, refuse it.
- **Exercise:** readHeader accepts any line with a colon; you refuse a line without a colon, an empty name or one with a space, a tab or another character outside the token set (400 `bad header line`), a value holding a control character other than the tab, such as NUL (400, with the given `isFieldValue`), a line starting with a space or a tab (400), a header line over 1024 bytes (431 `header line over 1024 bytes`) and a 33rd field (431 `more than 32 header fields`).
- **Check:** Each case gives its status and message as a RequestError; `Host :` and `Ho st:` give 400; a NUL, a bell or a DEL in a value gives 400, while a tab and bytes from 0x80 up are kept; 32 fields are read and 33 give 431. The mistake: trimming the name, which turns `Host :` into host.

#### 7.7 Quick Check: Header Fields
`quick-check-header-fields` · quiz · difficulty 8
- **Objective:** Reason about header names, values and why strictness matters.
- **Scenario:** Mirela's questions after her review.
- **Exercise:** Answer three questions.
- **Check:** Names match in any letter case and values keep theirs; a repeated field's values are kept in order; refusing `Host :` is safer than guessing what the next program will read.

#### 7.8 How Long Is the Body?
`how-long-is-the-body` · code_tests · difficulty 9
- **Objective:** Decide a body's length from Content-Length, refusing every value two programs could read differently.
- **Scenario:** The recorder's first test uploads say how long they are, and one of Tuesday's said `Content-Length: 4096, 4096`.
- **Exercise:** `contentLength(method string, h Header, max int64) (int64, error)` returns 0; you follow the bible's steps 3 and 4: one value of ASCII digits, at most max, only 0 on GET and HEAD, and 411 for a PUT with none.
- **Check:** PUT with 4096 or 0004096 gives 4096; `12a`, `-1`, `+5`, `4 096`, an empty value and `4096, 4096` give 400 `bad content-length "..."`; 65537 with a max of 65536 and a 20-digit number give 413; two Content-Length lines give 400; GET with 5 gives 400 `GET has no body` and GET with 0 gives 0; PUT with none gives 411. The mistake: strconv.Atoi, which takes `+5`.

#### 7.9 Read the Whole Head
`read-the-whole-head` · code_tests · difficulty 9
- **Objective:** Assemble a request reader from the line, header and length pieces, and answer every refusal with its status.
- **Scenario:** Friday 26 May: the station server answers its first requests in the hut; for now every good request gets the list of recordings, and every bad one its reason.
- **Exercise:** `readRequest(br *bufio.Reader, max int64) (*Request, error)` returns an empty request; you read the request line (414 for one over 1024 bytes) and the header, require exactly one Host (400 `missing host`, `more than one host`), accept only GET, HEAD and PUT (501), and set the length with contentLength.
- **Check:** Through the given server on a loopback port with Head at 50 ms: Go's own GET gives 200 and the list; `PATCH /recordings HTTP/1.1` gives 501; no Host and two Host lines give 400; a 1100-byte request line gives 414; 33 fields give 431; half a head and then silence gives 408 within a second; every refusal carries `Connection: close` and the connection ends. The mistake: checking Host with Get alone, which lets two Host lines through.

## Chapter 8 · Bodies and Chunks

`bodies-and-chunks` · 9 lessons

Read request bodies exactly: Content-Length, chunked transfer encoding and trailers, cut streams, size limits, binary data and a CRC that comes last.

#### 8.1 Exactly That Many Bytes
`exactly-that-many-bytes` · code_tests · difficulty 8
- **Objective:** Read a body of known length through the reader that read the head, with io.ReadFull, and refuse one that ends early.
- **Scenario:** Monday 5 June: the recorder uploads a 4,096-byte Pipit-2 test recording, and the station server must keep exactly those bytes, zeros, CRs and LFs included.
- **Exercise:** `readBody(br *bufio.Reader, n int64) ([]byte, error)` reads lines until it has n bytes, as chapter 7 read the head; you read exactly n bytes with io.ReadFull and answer a short body with 400 `body ended after 812 of 4096 bytes`.
- **Check:** testRecording(4096) sent in the same write as the head gives the same 4096 bytes; 0 gives an empty body with no read; 812 bytes and the end of the stream give the 400; Go's client's PUT of 1024 bytes through the given server is stored whole, CRC b70b4c26. The mistake: reading lines, which waits for a newline the recording doesn't end with.

#### 8.2 One Chunk's Size
`one-chunks-size` · code_tests · difficulty 9
- **Objective:** Parse a chunk-size line: hex digits, then any extension after a semicolon.
- **Scenario:** Go's client streams the recorder's pass in chunks, each one's size first, in hex: `c` for 12 bytes, `1000` for 4,096.
- **Exercise:** `chunkSize(line string) (int64, error)` uses strconv.Atoi; you cut at `;`, trim spaces and tabs after the digits, require 1 to 16 hex digits and parse them in base 16.
- **Check:** `c` gives 12; `1000` gives 4096; `1A` and `1a` give 26; `0` gives 0; `400;name=value` and `400 ;x` give 1024; `zz`, an empty line, `-1`, `+5`, `0x10` and 17 digits give 400 `bad chunk size "zz"`. The mistake: base 10, which reads `1000` as 1000 and refuses `c`.

#### 8.3 Quick Check: Bodies and Chunks
`quick-check-bodies-and-chunks` · quiz · difficulty 8
- **Objective:** Recognize where a body ends under each framing.
- **Scenario:** Nkem's three questions, with the booking service's 413s in mind.
- **Exercise:** Answer three questions.
- **Check:** A Content-Length body ends after exactly that many bytes; a chunked body ends at a zero-size chunk and the blank line after its trailer; chunked lets the recorder start before it knows how long the pass will be.

#### 8.4 Chunk After Chunk
`chunk-after-chunk` · code_tests · difficulty 9
- **Objective:** Read a chunked body: a size line, the data, CR LF, until the zero chunk, then the trailer.
- **Scenario:** Furrow-B's morning pass streams in as it comes down, in chunks of whatever size the recorder has ready.
- **Exercise:** `readChunked(br *bufio.Reader, max int64) ([]byte, Header, error)` returns nothing; you loop over readLine, chunkSize and io.ReadFull, require an empty line after each chunk's data (400 `bad chunk ending`), stop at size 0 and read the trailer with readHeader.
- **Check:** Go's own two-chunk upload gives `first piece second piece` and a trailer with apogee-crc32 e2fc055d; testRecording(4096) in chunks of 1, 7 and 4088 gives the bytes; `xx\r\n` after a chunk's data gives 400; a lone zero chunk gives an empty body; a chunk made of CR and LF bytes is kept. The mistake: reading a chunk's data with readLine, which splits at a CR LF inside it.

#### 8.5 Cut Off Mid-Pass
`cut-off-mid-pass` · code_tests · difficulty 9
- **Objective:** Tell a chunked body that ended from one that was cut: only the zero chunk and its blank line end it.
- **Scenario:** Tuesday 6 June: the recorder restarted in the middle of Furrow-A's 08:02 pass, and its stream simply stopped.
- **Exercise:** readChunked returns what it has, with no error, when the stream ends between chunks; you answer 400 `recording cut off` for an end anywhere before the final blank line, and the given handler then keeps nothing.
- **Check:** A stream ending after a whole chunk, inside a chunk, after the zero chunk and inside the trailer each give `recording cut off`; Go's client whose body fails mid-stream leaves nothing in the store or in `/recordings`; a complete stream is stored. The mistake: taking the end of the stream for the end of the body.

#### 8.6 One Length, Never Two
`one-length-never-two` · code_tests · difficulty 9
- **Objective:** Combine Content-Length and Transfer-Encoding into one length, refusing a request with both or with an encoding the server doesn't know.
- **Scenario:** Wednesday: Mirela shows how a request with both lengths, read one way by the station server and another by the next program, hides a second request inside the first.
- **Exercise:** `bodyLength(method string, h Header, max int64) (int64, error)` returns contentLength's answer; you add the bible's steps 1 and 2: both present give 400, a Transfer-Encoding whose last coding isn't chunked gives 400, a coding before a final chunked gives 501, chunked on GET or HEAD gives 400, and chunked gives −1.
- **Check:** `chunked` and `Chunked` give −1; `gzip` gives 400 `bad transfer-encoding "gzip"`, and so does `chunked, gzip`; `gzip, chunked` gives 501 `transfer-encoding "gzip, chunked" not implemented`; two Transfer-Encoding lines of chunked give 501; both headers give 400 `both content-length and transfer-encoding`, whichever comes first; GET with chunked gives 400; Content-Length alone behaves as before. The mistake: letting Content-Length win when both are present, which is the trick itself.

#### 8.7 Too Big to Keep
`too-big-to-keep` · code_tests · difficulty 9
- **Objective:** Refuse a body over the limit before reading it, and hang up, since its unread bytes are still on the line.
- **Scenario:** A misconfigured test recorder sent a whole day of passes as one body, and the bench computer ran out of memory reading it.
- **Exercise:** readChunked accepts chunks of any size; you answer 413 `body over 65536 bytes` as soon as a chunk's size would take the total past max, before reading that chunk.
- **Check:** Chunked bodies totaling 65536 bytes are stored and 65537 give 413 with the last chunk unread (the test sends only its size line and waits for the answer); a Content-Length of 65537 gives 413 with nothing sent after the head; both answers carry `Connection: close` and the connection ends. The mistake: adding a chunk's size to the total after reading the chunk.

#### 8.8 The CRC Comes Last
`the-crc-comes-last` · code_tests · difficulty 9
- **Objective:** Check a body against a CRC that arrives in a trailer, after the data it covers.
- **Scenario:** Thursday: Nkem points out that the recorder knows a recording's CRC only at LOS, so it sends it last, the way a satellite sends a file's CRC after its chunks.
- **Exercise:** `checkCRC(body []byte, h, trailer Header) (uint32, error)` returns the body's CRC without checking; you take apogee-crc32 from the trailer, else from the header, require 8 hex digits (400 `bad apogee-crc32 "xyz"`), and refuse a mismatch with 400 `crc mismatch: sent ..., received ...`.
- **Check:** Go's client with req.Trailer set is accepted, e2fc055d; a trailer one digit off gives 400 and keeps nothing; a CRC in a Content-Length upload's header is checked; no CRC at all is accepted; `A2912082` in capitals is accepted; `xyz` gives 400. The mistake: reading only the header, so a chunked upload's CRC is never checked.

#### 8.9 The Recorder Streams a Pass
`the-recorder-streams-a-pass` · code_tests · difficulty 9
- **Objective:** Accept an upload end to end: the length rules, either kind of body, the CRC, the store and the answer.
- **Scenario:** Friday 9 June: ten Furrow passes stream in, and Hollis wants every one: "EVERY JUNE PASS. EVERY FIELD."
- **Exercise:** `(s *Server) put(req *Request) answer` answers 201 to everything; you check the pass ID (404 `no such path`) and the CRC, keep the body with s.Store.Put, and answer 201 with Location, Apogee-Crc32 and `stored <pass ID> <n> bytes`, or 409 `<pass ID> already stored`.
- **Check:** Through Go's client: a chunked upload of furrow-a/YST/18856 (4096 bytes, trailer a2912082) gives 201 with those fields, and `/recordings` lists it; the same upload again gives 409; a Content-Length upload of cirrine-1/YST/13676 gives 201; a wrong CRC gives 400 and keeps nothing; a bad pass ID gives 404; a cut stream keeps nothing. The mistake: putting the body in the store before checking its CRC.

## Chapter 9 · Answers from Scratch

`answers-from-scratch` · 9 lessons

Write HTTP/1.1 answers by hand: status lines, headers and bodies, HEAD, persistent connections, chunked streams and a clean shutdown, read by Go's own client.

#### 9.1 A Status Line
`a-status-line` · code · difficulty 8
- **Objective:** Write a status line: the version, the code and its reason, ended by CR LF.
- **Scenario:** Monday 19 June: the station server's answers came from a helper until now, and this week you write them yourself.
- **Exercise:** `statusLine(code int) string` returns `"HTTP/1.1 200 OK\r\n"` for every code; you look the reason up in the given `reasons` table, and `main` prints four with `%q`.
- **Check:** `"HTTP/1.1 200 OK\r\n"` / `"HTTP/1.1 201 Created\r\n"` / `"HTTP/1.1 404 Not Found\r\n"` / `"HTTP/1.1 413 Request Entity Too Large\r\n"`. The starter prints the first line four times.

#### 9.2 Head, Blank Line, Body
`head-blank-line-body` · code_tests · difficulty 9
- **Objective:** Write an answer's head, a blank line and its body through one buffered writer, flushed once.
- **Scenario:** The collector at Port Ellery reads these bytes over the island's line, so the station server must send them exactly and all at once.
- **Exercise:** `writeResponse(w io.Writer, now time.Time, status int, fields []Field, body []byte) error` writes only the status line; you add Date from now in dateLayout, the fields in order, the blank line and the body through one bufio.Writer, returning Flush's error.
- **Check:** The bible's 201 byte for byte at 06:40:36 UTC; http.ReadResponse reads the same bytes back as 201 with its body; a now in another time zone gives the same Date in GMT; a failing writer's error is returned; a counting writer sees one Write. The mistakes: no blank line, and Date in local time.

#### 9.3 Quick Check: Status Lines
`quick-check-status-lines` · quiz · difficulty 9
- **Objective:** Check the shape of an answer on the wire.
- **Scenario:** Nkem's questions before Go's own client meets the station server.
- **Exercise:** Answer three questions.
- **Check:** The reason phrase is for people, and a client reads the code; the blank line ends the head; without Content-Length or chunks a client can only read until the connection closes, and can't tell a cut answer from a whole one.

#### 9.4 The Head Without the Body
`the-head-without-the-body` · code_tests · difficulty 9
- **Objective:** Answer HEAD with GET's status and fields, Content-Length included, and no body.
- **Scenario:** Before fetching an 11 GB recording over the line, the collector asks how big it is.
- **Exercise:** `(s *Server) recording(req *Request) answer` returns the body for HEAD too; you keep every field, Content-Length the full size, and leave the body nil for HEAD.
- **Check:** GET of furrow-a/YST/18856 through Go's client gives 200, 4096 bytes, application/octet-stream and Apogee-Crc32 a2912082; HEAD gives ContentLength 4096 through Go's client and nothing after the blank line on a raw connection; an unknown recording gives 404 to both. The mistake: Content-Length 0 on HEAD.

#### 9.5 Keep the Line Open
`keep-the-line-open` · code_tests · difficulty 9
- **Objective:** Serve several requests on one connection: read the next after a successful answer, and close after an error, after `Connection: close`, or when the line sits idle.
- **Scenario:** Over the backhaul each new connection costs a 0.7 s round trip, and the collector makes three requests for every recording.
- **Exercise:** `(s *Server) serveConn(c net.Conn)` answers one request and closes; you loop: wait for the next request's first byte for at most Idle, read and answer it, and go on only after a 2xx answer to a request without `Connection: close`.
- **Check:** Go's client's GET, HEAD and list GET use one accepted connection; a request with `Connection: close` gets it back and the stream ends; a 404 closes the connection; a raw connection silent after its first answer is closed within a second at Idle 50 ms. The mistake: going on after an error answer, whose request may still have a body on the line.

#### 9.6 Answer in Pieces
`answer-in-pieces` · code_tests · difficulty 10
- **Objective:** Write a chunked answer: each piece as its hex size, its bytes and CR LF, flushed at once, then the zero chunk.
- **Scenario:** The collector will keep one request open on `/feed` and hear of each recording the moment it lands, instead of asking every few seconds.
- **Exercise:** `writeChunk(w *bufio.Writer, p []byte) error` and `endChunks(w *bufio.Writer) error` are empty, and `(s *Server) feed(bw *bufio.Writer, lines <-chan string) error` writes nothing after the head; you write each line as one flushed chunk and the zero chunk when lines closes.
- **Check:** Go's client sees TransferEncoding chunked and ContentLength −1 and reads the first line before the test sends the second; after the channel closes the body ends without error; pieces of 4096 and 1 bytes get sizes `1000` and `1`; a collector that leaves makes feed return the write's error. The mistake: no flush, so each line waits in the buffer for the next.

#### 9.7 Quick Check: Open Lines
`quick-check-open-lines` · quiz · difficulty 9
- **Objective:** Reason about persistent connections and streamed answers.
- **Scenario:** Ingrid's questions before the station server takes every pass.
- **Exercise:** Answer three questions.
- **Check:** A client can reuse a connection only when it knows where each answer ends; a request's body must be read whole before the next request; the zero chunk tells a client a streamed answer ended, and its absence that it was cut.

#### 9.8 Stop Without Cutting
`stop-without-cutting` · code_tests · difficulty 10
- **Objective:** Stop a server without cutting work: refuse new connections, close idle ones, end the feeds and let requests in flight finish, within a limit.
- **Scenario:** The station server's first update is due during Tuesday's passes, and it must never cut an upload.
- **Exercise:** `(s *Server) Shutdown(ctx context.Context) error` closes only the listener; you also close connections waiting between requests, call s.Store.Close, wait for requests in flight, and when the context ends first, close everything and return its error.
- **Check:** An upload in flight when Shutdown starts gets 201 and is kept, then Shutdown returns nil; a kept-alive idle connection closes at once; a feed ends with its zero chunk; a new dial is refused; an upload that never finishes, with a 100 ms context, gives context.DeadlineExceeded and its connection closed; Serve returns nil. The mistake: closing every connection at once, which cuts the upload.

#### 9.9 The Collector Comes Calling
`the-collector-comes-calling` · code_tests · difficulty 10
- **Objective:** Put every route in one place and run the whole station server against Go's own clients.
- **Scenario:** Tuesday 20 June, midsummer: Furrow-A's 06:32 pass is the first Furrowlight pass through the station server, and its image must reach Furrowlight within a minute of LOS.
- **Exercise:** `(s *Server) route(req *Request) answer` sends every request to the list; you route `/recordings`, a recording and `/feed` by method, answering 404 `no such path` and 405 `method not allowed` with each path's Allow.
- **Check:** The collector follows `/feed` while the recorder streams furrow-a/YST/18856 with its trailer; the feed's line arrives within a second and the collector's GET returns 4096 bytes whose CRC matches, on one connection; Cirrine-1's upload cut mid-stream never reaches the feed; `PUT /recordings` gives 405 with `Allow: GET, HEAD`; HEAD /feed gives 405 with `Allow: GET`; DELETE gives 501; Shutdown ends the feed cleanly and closes everything within a second. The mistake: 404 for a known path with the wrong method.

## Changes in writing and review

- **Given code** keeps the text the learner last saw, lesson after lesson, checked over the
  whole course. Chapter 2's read-only `decodeFrame` only splits a payload: the CRC check is
  3.7's task, and 3.7's story has dishd showing a damaged STATUS.
- **Chapter 4** `recvCtx` waits for its `AfterFunc` to finish before it returns, as the
  context package's own example does, so a deadline set as the context ends can't cut the next
  read. Deadline lessons run in `testing/synctest` bubbles with the production timeout.
- **Chapter 5** `Watchers` gains a `closed` flag set by `CloseAll`: a watcher accepted just
  before Shutdown is closed at once instead of hanging `Wait`. Its tests run in bubbles, so
  they never depend on map order or real waits.
- **Chapter 6** `l.attempt` counts failed dials in a row and resets on a dial that works; the
  given `connect` keeps its own count of failed handshakes, so their backoff grows too.
- **Chapter 7** `readLine` refuses a lone CR as well as a bare LF, and `readHeader` a value
  holding a control character other than the tab, as RFC 9110 §5.5 asks of a NUL and as Go's
  own server does. The readmes say where the station server is stricter than HTTP: HTTP/1.1
  only, targets starting with `/`, an LF alone, and a `Content-Length` list such as
  `4096, 4096`, which Mirela's crooked client sends.
- **Chapter 8** A request whose last transfer coding isn't `chunked` gets 400, as RFC 9112
  requires; a coding before a final `chunked` stays 501. Chunk sizes of 1 to 16 hex digits are
  the station server's limit, not HTTP's. `Request.Close` is left out: chapter 9's `closing`
  helper, given from 9.5, decides when to hang up.
- **Passes** all come from one orbit model; Pipit-2's orbit was fitted for this season, and
  passes that model contradicted moved to real ones. Chapters 3, 4 and 6 keep their passes at
  the model's times: Cirrine-1's 07:02, 03:43 and 07:28 and Furrow-A's 07:52, so 3.9 opens at
  06:50.
