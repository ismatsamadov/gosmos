# Apogee story bible: Season 4

### Season 4: down to the wire

After the Board project week, the learner starts the next assignment on Monday 6 March 2028.
Their new badge says **Station Software Engineer**, still No. 039, based at Ystmark for the
spring, since Monday 6 March 2028. Ystmark's dish controllers, the computers in the radomes
that turn the dishes, are worn out: Dish 2's lost its position twice in February. New ones
arrive in a crate on the same plane as the learner, with a manual that is wrong on page one and
a stick holding the maker's simulator. They speak a protocol of raw bytes over a network
connection, and Apogee drives them from Port Ellery, across the island's line to the mainland,
which drops in every storm. You write **dishd**, the program that talks to them, then a
station server for the hut that keeps every recording until headquarters collects it, built on
a bare socket so you see what every web request of Season 3 really was, and, the week before the finale, a live view that hears each pass as the recorder sends it, a datagram at a time.

Story time runs from Monday 6 March to Tuesday 20 June 2028: ten chapters, about a week each,
with installation weeks between some of them. Sunniva carries the season; Nkem reviews every
change from Port Ellery; Tomás flies up for Dish 1's calibration; Ingrid runs Ystmark's passes
from the Fishbowl; Mirela reviews the station server and the live view's answers; Yusra and Hollis bring the customers'
losses. No new people, places, satellites or customers are introduced.

#### Ystmark in spring

- **Light.** The polar night ended in mid-February; on 6 March the sun is up about eight hours.
  The midnight sun begins on 19 April: from chapter 5 on, the sun never sets. The solstice,
  Tuesday 20 June, is the finale.
- **Cold and wind.** Outside −18 °C in March; inside a radome, with its heater on, about
  −12 °C. Storms come through until May. The dishes stow when the wind reaches 72 km/h, and the
  new controllers do it themselves (the wind lock, below).
- **The hut** is the station's one building, below the ridge with the two radomes: the bench
  computer, Sunniva's terminal, the recorders' rack and the coffee machine. In readmes it is
  "the hut", a plain word, never a name.
- **The island's line.** Normally a fiber cable to the mainland: about 40 ms there and back. In
  storms the fiber fails and the station falls back to the satellite backhaul: about 0.7 s
  there and back and a few hundred kilobits a second, enough for commands and status, never for
  an 11 GB recording. Connections that were open on the fiber die silently when it drops; new
  ones must be made over the backhaul.
- **Who drives the dishes.** Ingrid runs every pass of every station from the Fishbowl, so
  Ystmark's controllers are driven from Port Ellery: dishd runs in Port Ellery's rack and
  reaches the controllers across the island's line. During a pass the controller follows the
  satellite by itself, from predictions passplan loads with the maker's own tool; dishd's
  commands are the ones between passes and in emergencies: point, stow, unstow and jog. dishd is
  the only program that sends commands to Ystmark's controllers; people watch through dishd.
- **The X-band beam.** A 3.7 m dish at Furrow-A's 8150 MHz has a beam about 0.70° wide (to
  half power, 70 λ/D). Pointed 0.35° off, the signal drops by half (3 dB); 0.70° off, to a
  sixteenth (12 dB). Tomás: "Half a beamwidth off, and Furrow-A sounds half as loud."

#### The season, chapter by chapter

You live and work at Ystmark all season. Times are UTC. Pass IDs follow the bible's rule;
orbit numbers continue Season 3's (Furrow-A 15.2 a day, Cirrine-1 14.92).

| Ch | Week of (2028) | Where | Who needs what | What is at stake |
|---|---|---|---|---|
| 1 | Mon 6 Mar | the hut's bench, beside the crate | Sunniva: dishd must talk to the maker's simulator before Dish 1's old controller comes out on Monday 20 March | Ystmark hears a polar orbit on almost every revolution; a controller nobody can talk to is a dish nobody can use. By Friday 10 March dishd says hello: `YST-1 firmware 2.4, last command 0` |
| 2 | Mon 13 Mar | the hut, through Sunniva's test line to Port Ellery over the backhaul | Nkem: dishd must read whole frames however the line splits or glues them. The manual's page one is wrong: the banner ends in CR LF | Over the slow line an answer arrives in pieces, or two arrive in one read. In Wednesday's test dishd lost the second, a STATUS in which the simulated Dish 1 had stowed itself in the wind, and nobody saw the stow |
| 3 | Mon 27 Mar | Dish 1's radome and the hut (Dish 1's controller went in from 20 to 24 March) | Sunniva: Dish 1's first status frames, read as the manual's table says, put it at 382.15° in a 3,277 km/h wind. Tomás flies up on Tuesday 28 March to point it true | Friday 31 March, Cirrine-1's 07:02 pass (`cirrine-1/YST/13094`), Dish 1's first customer pass on the new controller |
| 4 | Mon 10 Apr | the hut, the first spring storm | Ingrid: on Tuesday 11 April the fiber dropped at 03:10 and dishd waited forty minutes on a connection that was already dead | Dish 1 sat still through Cirrine-1's 03:43 pass (`cirrine-1/YST/13256`). Yusra: "TL;DR: Cirrine lost its 03:43 pass and asks what stops it happening again." |
| 5 | Mon 24 Apr | the hut, midnight sun (Dish 2's controller went in from 17 to 21 April) | Sunniva: dishd must run both dishes at once and show them on a watch port to her terminal, Ingrid's screen and Tomás's laptop | One dish's trouble must not stop the other; Friday 28 April's deploy must leave nothing hanging |
| 6 | Mon 8 May | the hut, after the big storm | Sunniva, Tuesday 9 May: "Link down 14 min. Dish 2 waiting on your code. Coffee low." | From 07:22 to 07:36 the island had no line at all. dishd gave up and never dialed again; Ingrid restarted it at 07:51, too late for Cirrine-1's 07:28 pass (`cirrine-1/YST/13676`). Tomás's 0.35° calibration jog had left dishd at 07:22:04, two seconds before the line went; the restarted dishd sent it again under a new number, so Dish 2 moved twice and was 0.70° off for Furrow-A's 07:52 pass (`furrow-a/YST/18218`), which came down at a sixteenth of its strength. Friday 12 May at 10:00 Sunniva pulls the line on purpose, in the middle of a jog |
| 7 | Mon 22 May | the hut | Sunniva: on 9 May the fiber stayed down until 19:40; recordings had nowhere to wait but the recorder's small disk, which filled at 13:20 and lost Furrow-B's 13:21 pass. Hollis: "LOST FURROW-B 13:21. FIELD 7. WHY." Sunniva wants a station server in the hut that keeps every recording until headquarters collects it. Nkem: build it on the socket you know | The station server's first requests must be read exactly, and the crooked ones refused, never guessed at |
| 8 | Mon 5 Jun | the hut | Hollis, through Yusra: Furrowlight's June campaign, ten booked Furrow passes a day over Ystmark. The recorder streams each pass to the station server as it comes down, its CRC last. Mirela reviews the length rules on Wednesday | On Tuesday 6 June the recorder restarted in the middle of Furrow-A's 08:02 pass (`furrow-a/YST/18644`): a cut stream must never be kept as a whole recording |
| 9 | Mon 12 Jun | the hut; Ingrid's screen in the Fishbowl | Ingrid: the Fishbowl hears nothing of a Furrow pass until its recording lands, and on 6 June it learned of the recorder's restart from a cut upload. She wants every pass's frames on her screen as they come down, and the recorder asked, whenever they stop, what it is doing. The recorder gets a live port that speaks UDP; you write the live view at Port Ellery. Monday's first test is Pipit-2's 06:36 pass (`pipit-2/YST/16696`) | Over the island's line a datagram can be lost, come twice or overtake another, and nothing says so. Friday 16 June, Furrow-A's 07:52 pass (`furrow-a/YST/18796`, LOS 08:00:51) is the first shown whole: three frames lost on the line, none shown twice, and at 08:00:56 the recorder answers `IDLE` at the first try |
| 10 | Mon 19 Jun, then Tue 20 Jun | the hut; the Fishbowl on the Board | Ingrid and Yusra: headquarters' collector must hear of each recording at once and fetch it over one open line. On Tuesday 20 June the station server takes every pass from 06:00 | Furrow-A's 06:32 pass (`furrow-a/YST/18856`, LOS 06:40:35) is the first Furrowlight pass through it; its image must reach Furrowlight within a minute of LOS, as the booking service promises |

The opening (Monday 6 March) shows Ystmark in one image (the sun back for three weeks, two
radomes on a ridge, −18 °C), hands over the badge, introduces Sunniva (station chief), Nkem
(reviewer, at Port Ellery), Tomás (coming up in March), the crate and the island's line, and
sends the learner to the simulator on the hut's bench: the old controllers come out on the 20th.

The finale (chapter 10, Tuesday 20 June): Furrow-A's recording streams into the station server
during the pass; the collector at Port Ellery hears of it on the feed at 06:40:36, has every
byte and its CRC by 06:40:38, and Furrowlight's webhook from the booking service goes out at
06:40:39. Ingrid watches the Board. Sunniva: "Midnight sun. Both dishes. Every byte. Coffee for
two." The debrief points at Season 5: Mirela notes that anyone who can listen on the island's
line can read every byte on it, and Cirrine has asked that its commands be readable by Cirrine
alone.

#### New systems

- **The dish controllers.** One per dish, in the radome's rack: the computer that turns the
  dish. Named at installation by station code and dish number: `YST-1` for Dish 1, `YST-2` for
  Dish 2. Firmware 2.4. They speak the controller protocol below. In readmes: "the dish
  controller" first, "the controller" after.
- **The manual.** The maker's manual for the controllers. It is wrong in exactly two places,
  and writers never invent a third:
  1. **Page one** says each connection begins with a banner line ending in a newline (0x0A).
     It ends in a carriage return and a newline (0x0D 0x0A).
  2. **Page 14**, the status frame's table, says its fields are little-endian. Like every
     other field of the protocol they are big-endian.
- **The maker's simulator**, on a stick in the crate: a program that behaves exactly like a
  controller, run on the hut's bench before the real ones are installed. In lessons it is the
  hidden file `sim.go` (below). The maker is never named.
- **dishd**, the dish service, which this season builds: add it to the Services list ("dishd
  (talks to the dish controllers at a station, and shows them on its watch port)"). In readmes
  it is "the dish service" first, `dishd` once glossed. It runs at Port Ellery and keeps one
  connection to each of Ystmark's controllers.
- **The watch port**: dishd's listening port where anyone can watch both dishes, one line per
  dish every second (below). Sunniva keeps a terminal open on it in the hut.
- **The recorder** (a word from Season 2): the receiver's recording program in the hut's rack,
  which records each pass as it comes down. It is a Go program; it talks to the station server
  with `net/http`'s client and, from Monday 12 June, to the live view over its live port (below).
- **The station server** (chapters 7, 8 and 10): the hut's server that keeps every finished
  recording until headquarters collects it, and tells headquarters when one arrives. It has no
  program name; readmes say "the station server".
- **The collector**: dlinkd's part at Port Ellery that follows the station server's feed and
  fetches each recording. A Go program using `net/http`'s client.
- **The live view** (chapter 9): the program at Port Ellery that shows Ingrid each Ystmark pass as it comes down, every frame the recorder sends from its **live port** (below), in order and each once, and asks the recorder how it is when the frames stop. It has no program name; readmes say "the live view", and "the live port" for the recorder's UDP port.

#### The controller protocol

Written here once; every lesson uses it exactly. All multi-byte numbers are big-endian
(network byte order), the most significant byte first.

**Connection.** TCP. The controller listens on its control port; dishd dials it. On accepting
a connection the controller at once sends one line of ASCII, the **banner**, then speaks only in
frames. It serves one client at a time: when a new connection arrives, the controller sends
BYE `replaced` on the old one, closes it, and serves the new one.

```
YST-1 READY PROTO 1\r\n
```

The banner is the controller's name (1 to 16 capital letters, digits and hyphens), `READY`,
`PROTO` and the protocol version, `1`, separated by single spaces and ended by CR LF (page one
of the manual says LF alone). It is 21 bytes for `YST-1`. A reader accepts LF with or without
the CR before it, and refuses any other shape or version.

**Frames.** Everything after the banner, in both directions:

| Offset | Size | Field | Meaning |
|---|---|---|---|
| 0 | 2 | length | uint16: the number of bytes after this field, 9 + the body's length; 9 to 1024 |
| 2 | 1 | type | uint8: the frame type (table below) |
| 3 | 4 | seq | uint32: a command's sequence number; 0 on every other frame except ACK and NAK |
| 7 | n | body | 0 to 1015 bytes, by type |
| 7+n | 4 | crc | uint32: CRC-32 (IEEE, `crc32.ChecksumIEEE`) of the type, seq and body bytes |

The CRC covers exactly the bytes between the length and the CRC: not the length, not itself.
A frame is 11 bytes at the least and 1026 at the most. A length under 9 or over 1024 means the
stream is out of step: the controller sends BYE `bad length` and closes, and dishd closes too.

**Frame types.** Requests and commands go from dishd to the controller; answers come back in
the order the frames arrived. The controller never sends a frame except as an answer, and BYE
just before it closes.

| Type | Name | From | Seq | Body | Answer |
|---|---|---|---|---|---|
| 0x01 | HELLO | dishd | 0 | the client's name, 0 to 16 ASCII bytes (`dishd`) | WELCOME |
| 0x02 | PING | dishd | 0 | empty | PONG |
| 0x03 | QUERY | dishd | 0 | empty | STATUS |
| 0x10 | POINT | dishd | the command's | az uint16, el int16 (4 bytes) | ACK or NAK |
| 0x11 | STOW | dishd | the command's | empty | ACK or NAK |
| 0x12 | UNSTOW | dishd | the command's | empty | ACK or NAK |
| 0x13 | JOG | dishd | the command's | axis uint8, delta int16 (3 bytes) | ACK or NAK |
| 0x81 | WELCOME | controller | 0 | 7 bytes, below | |
| 0x82 | PONG | controller | 0 | empty | |
| 0x83 | STATUS | controller | 0 | 16 bytes, below: "the status frame" | |
| 0x90 | ACK | controller | the command's | empty: the command was executed | |
| 0x91 | NAK | controller | the refused frame's | code uint8, then the reason in ASCII | |
| 0xFF | BYE | controller | 0 | the reason in ASCII | (the controller closes) |

Requests (HELLO, PING, QUERY) and their answers carry seq 0, and a request's answer has the
request's type with the high bit set. Commands are answered by ACK or NAK, carrying the
command's seq.

**WELCOME body**, 7 bytes: 0 dish (uint8, 1 or 2); 1 firmware major (uint8, 2); 2 firmware
minor (uint8, 4); 3–6 last_seq (uint32): the sequence number of the last command this
controller executed, 0 if none.

**The status frame** (STATUS body, 16 bytes), the one page 14 of the manual gets wrong:

| Offset | Size | Type | Field | Unit and range |
|---|---|---|---|---|
| 0 | 1 | uint8 | dish | 1 or 2 |
| 1 | 1 | uint8 | state | 0 IDLE, 1 MOVING, 2 TRACKING, 3 STOWED, 4 FAULT |
| 2 | 2 | uint16 | az | azimuth in hundredths of a degree, 0 to 35999 |
| 4 | 2 | int16 | el | elevation in hundredths of a degree, 0 to 9000 (the encoders can read a little below 0) |
| 6 | 2 | uint16 | wind | the station anemometer's wind in tenths of km/h |
| 8 | 2 | int16 | temp | inside the radome, tenths of a °C |
| 10 | 2 | uint16 | flags | bit 0 heater on, bit 1 wind lock, bit 2 fault; other bits are 0 and ignored |
| 12 | 4 | uint32 | last_seq | as in WELCOME |

Monday 27 March, Dish 1 tracking Pipit-2's 09:14 test pass, the frame the chapter 3 story
starts from: `01 02 47 95 10 4a 02 80 ff 83 00 01 00 00 00 76` is dish 1, TRACKING, az 183.25°,
el 41.70°, wind 64.0 km/h, −12.5 °C, heater on, last_seq 118. Read little-endian as the manual
says: az 382.15°, el 189.60°, wind 3277.0 km/h, −3174.5 °C, last_seq 1979711488. With the right
byte order but the temperature read unsigned: 6541.1 °C. Sunniva: "Dish 1 at 382 degrees.
Wind 3277. Neither true."

**Command bodies.**

- **POINT**: az uint16 then el int16, hundredths of a degree, rounded to the nearest
  hundredth (`math.Round`, never truncated: 1.15° is 115, not 114). az 0 to 35999, el 0 to
  9000. (183.25°, 41.70°) is `47 95 10 4a`.
- **JOG** moves the dish from where it is: axis `A` (0x41, azimuth) or `E` (0x45, elevation),
  then delta int16 in hundredths, −500 to 500 and not 0. Azimuth wraps (359.90° + 0.35° is
  0.25°); an elevation that would leave 0 to 90° is refused. `A +0.35°` is `41 00 23`;
  `E −0.29°` is `45 ff e3`. JOG is the command that must never run twice: run twice, it moves
  the dish twice.
- **STOW** turns the dish to elevation 90.00° and holds it there: state STOWED. Accepted when
  already stowed. **UNSTOW** returns a stowed dish to IDLE where it is; accepted when not
  stowed (it changes nothing).

**NAK codes.** The body is the code, then its reason text, exactly:

| Code | Reason | When |
|---|---|---|
| 1 | `duplicate` | a command whose seq is at or below last_seq: nothing is executed |
| 2 | `bad argument` | a body of the wrong length, a value out of range, a JOG of 0 |
| 3 | `wind lock` | UNSTOW, POINT or JOG while the wind is at or above 72.0 km/h |
| 4 | `stowed` | POINT or JOG while stowed (UNSTOW first) |
| 5 | `unknown type` | a type not in the table, answers' types included |
| 6 | `bad crc` | a frame whose CRC doesn't match; its seq is echoed as received |
| 7 | `hello first` | anything but HELLO before the connection's first HELLO |

**BYE reasons**, exactly: `idle`, `slow frame`, `bad length`, `replaced`.

**What the controller does with a frame**, in this order, the first failure answering: the
CRC (NAK 6); HELLO first (NAK 7); a known type (NAK 5); for commands, the sequence number (NAK
1); the body (NAK 2); the wind lock (NAK 3); stowed (NAK 4); then it executes, sets last_seq to
the command's seq, and answers ACK. Only an executed command moves last_seq: a refused one
leaves it, so a refused number can be sent again. HELLO may be sent again on the same
connection; each gets a WELCOME.

**The wind lock.** At 72.0 km/h or more the controller stows its dish by itself (state
STOWED, el 90.00°), sets flag bit 1 and refuses UNSTOW, POINT and JOG with NAK 3. STOW is always
accepted. Below 72.0 the flag clears, and the dish stays stowed until UNSTOW.

**Sequence numbers, and never twice.** Seen from both ends:

1. Every command carries a sequence number. dishd numbers commands strictly upward per
   controller and sends one command at a time: the next waits for the previous one's ACK or
   NAK. Requests (HELLO, PING, QUERY) may be sent at any time, several at once.
2. The controller executes a command only if its seq is above last_seq, and keeps last_seq
   through power cuts. Anything else is NAK `duplicate`, and nothing moves.
3. After every connect and reconnect dishd sends HELLO first. From the WELCOME's last_seq:
   - new commands are numbered from `max(next, last_seq+1)`, so numbers never go back;
   - if a command was in flight when the line dropped (sent, no answer), it **ran** if
     `last_seq >= its seq`: it is done and is not sent again. Otherwise it never arrived, or was
     refused: it is sent again **with the same seq**.
4. A NAK `duplicate` for dishd's own command in flight therefore means it already ran, and
   counts as done. (This holds because dishd is the only sender of commands and sends one at a
   time.)

So a command the controller has executed is never sent again, and one it never received is
sent again under its own number, which the controller would refuse if it had run.

**Timeouts.**

| Who | Limit | Value | What happens |
|---|---|---|---|
| controller | idle: no complete frame received | 30 s | BYE `idle`, then close. A banner-only connection counts too |
| controller | slow frame: from a frame's first byte to its last | 10 s | BYE `slow frame`, then close |
| controller | a write of an answer that dishd doesn't read | 10 s | close, without BYE |
| controller | answering a complete frame | within 50 ms | (the simulator answers at once) |
| dishd | dialing a controller | 5 s | the dial fails |
| dishd | reading: each read must bring at least one byte | 5 s | `os.ErrDeadlineExceeded`: the line is lost |
| dishd | writing one frame | 5 s | `os.ErrDeadlineExceeded`: the line is lost |
| dishd | heartbeat | PING after 10 s without sending | keeps the controller's 30 s idle limit far away and finds a dead line within 15 s |
| dishd | poll | QUERY each dish every second (chapter 5) | the watch port's lines |
| dishd | reconnect | wait 1 s, doubling per failed dial in a row, at most 30 s, times a jitter in [0.5, 1.0); a handshake that fails as a lost line waits the same way, on its own count | a dial that works starts the dial count again; a WELCOME ends the handshake count |

dishd's read deadline is an idle deadline, moved forward before every read, not one deadline
for a whole frame: over the backhaul a frame may trickle in over seconds and is still alive.

**A line is lost** (dishd closes it and reconnects) on: `io.EOF`, `io.ErrUnexpectedEOF`,
`os.ErrDeadlineExceeded`, `net.ErrClosed`, `io.ErrClosedPipe`, `syscall.ECONNRESET`,
`syscall.EPIPE`, a BYE (`*ByeError`), `ErrBadLength` and `ErrBadCRC` (an answer that can't be
trusted: after the reconnect, last_seq says whether the command ran). Not lost: a NAK (an
answer) and the caller's context ending (stop, don't redial).

**Frames to copy** (hex, a space between bytes):

| Frame | Bytes |
|---|---|
| HELLO `dishd` | `00 0e 01 00 00 00 00 64 69 73 68 64 97 98 f8 ab` |
| WELCOME YST-1, 2.4, last 0 | `00 10 81 00 00 00 00 01 02 04 00 00 00 00 9e 88 e9 9f` |
| WELCOME YST-1, 2.4, last 118 | `00 10 81 00 00 00 00 01 02 04 00 00 00 76 27 ee 3d 96` |
| WELCOME YST-1, 2.4, last 10 (a 0x0a in the middle) | `00 10 81 00 00 00 00 01 02 04 00 00 00 0a 7e 5d 00 81` |
| WELCOME YST-2, 2.4, last 57 | `00 10 81 00 00 00 00 02 02 04 00 00 00 39 f0 65 7b 0a` |
| WELCOME YST-2, 2.4, last 121 | `00 10 81 00 00 00 00 02 02 04 00 00 00 79 86 b9 3a 9a` |
| WELCOME YST-2, 2.4, last 122 | `00 10 81 00 00 00 00 02 02 04 00 00 00 7a 1f b0 6b 20` |
| PING | `00 09 02 00 00 00 00 bc e2 a4 7d` |
| PONG | `00 09 82 00 00 00 00 0d 04 14 ef` |
| QUERY | `00 09 03 00 00 00 00 81 82 8d cd` |
| STATUS, Dish 1, Monday 27 March | `00 19 83 00 00 00 00 01 02 47 95 10 4a 02 80 ff 83 00 01 00 00 00 76 62 9d e0 f1` |
| STATUS, Dish 2 stowed by the wind (81.0 km/h, −3.1 °C, heater and wind lock, last 121) | `00 19 83 00 00 00 00 02 03 47 95 23 28 03 2a ff e1 00 03 00 00 00 79 f1 bf 34 97` |
| POINT seq 119, 183.25°, 41.70° | `00 0d 10 00 00 00 77 47 95 10 4a 20 eb ca 40` |
| POINT seq 119, 1.15°, 0.57° | `00 0d 10 00 00 00 77 00 73 00 39 08 89 34 01` |
| STOW seq 120 | `00 09 11 00 00 00 78 c5 7c b0 21` |
| UNSTOW seq 121 | `00 09 12 00 00 00 79 f5 db fa 67` |
| JOG seq 122, A +0.35° | `00 0c 13 00 00 00 7a 41 00 23 40 ab 3f dc` |
| JOG seq 122, E −0.29° | `00 0c 13 00 00 00 7a 45 ff e3 4f e2 a8 c2` |
| ACK seq 119 | `00 09 90 00 00 00 77 d9 45 34 92` |
| ACK seq 122 | `00 09 90 00 00 00 7a a7 f4 48 2f` |
| NAK seq 119 `stowed` | `00 10 91 00 00 00 77 04 73 74 6f 77 65 64 3f 3f 2e c5` |
| NAK seq 121 `wind lock` | `00 13 91 00 00 00 79 03 77 69 6e 64 20 6c 6f 63 6b 7e 15 64 e0` |
| NAK seq 122 `duplicate` | `00 13 91 00 00 00 7a 01 64 75 70 6c 69 63 61 74 65 63 cc 43 07` |
| BYE `idle` | `00 0d ff 00 00 00 00 69 64 6c 65 52 26 72 ec` |
| BYE `replaced` | `00 11 ff 00 00 00 00 72 65 70 6c 61 63 65 64 c9 df 77 ba` |

The mistakes these catch: a PING whose CRC also covers the length ends `c0 a3 dd d6`; one whose
CRC is written little-endian ends `7d a4 e2 bc`.

#### The simulator

`sim.go` is the maker's simulator in Go: one controller, in memory. It is written once, in the
calibration chapter (chapter 1), with every knob below, and copied unchanged into every lesson
of chapters 1 to 6 that needs a controller. It is **hidden**: the readme says in a sentence what
it does in that lesson. Every top-level name in it begins with `sim` or `Sim`, so it never
collides with the learner's names (`Frame`, `readFrame`, `Status`...).

```go
func NewSim(name string, dish uint8) *Sim // firmware 2.4, IDLE at az 180.00 el 45.00,
                                          // wind 20.0, −12.5 °C, heater on, last_seq 0
func (s *Sim) Serve(ln net.Listener) error // accepts until ln closes, then returns nil;
                                           // the newest connection wins (BYE replaced)
func (s *Sim) ServeConn(c net.Conn)        // one connection, until it ends; for net.Pipe
func (s *Sim) Executed() []string          // commands executed, in order
func (s *Sim) Received() []string          // every frame received, in order
func (s *Sim) Conns() int                  // connections accepted so far
func (s *Sim) ClosedByClient() int         // connections the client closed first
func (s *Sim) LastExecuted() uint32        // last_seq now
func (s *Sim) SetWind(kmh float64)         // at 72 or more: stow and lock at once
```

Entries of `Executed` and `Received` read `119 POINT 183.25 41.70`, `120 STOW`, `121 UNSTOW`,
`122 JOG A +0.35`; `Received` also has `HELLO dishd`, `PING`, `QUERY` and `bad crc`.

Fields a test sets before serving (each read under the simulator's mutex while it serves):

| Field | Default | Effect |
|---|---|---|
| `LastSeq uint32`, `Wind`, `TempC`, `AzDeg`, `ElDeg float64`, `State uint8`, `Flags uint16` | as `NewSim` | the controller's state |
| `Split int` | 0 | write every message (banner and frames) in pieces of at most this many bytes, one Write each |
| `Gap time.Duration` | 0 | pause before each piece after the first; at most 50 ms outside `synctest`, any length inside it |
| `Glue int` | 0 | hold the answers to the connection's first n frames until all n have arrived, then write them in one Write; later frames are answered one by one |
| `HoldBanner bool` | false | write the banner together with the first answers instead of at once (only when the test sends frames before reading the banner) |
| `Silent bool` | false | banner, then read and never answer |
| `MuteAfter int` | 0 | answer this many frames (HELLO counted), then read and never answer |
| `NoRead bool` | false | banner, then never read |
| `DropBefore uint32` | 0 | when the command with this seq arrives: close without executing or answering, once |
| `DropAfter uint32` | 0 | execute the command with this seq, then close without answering, once (the lost ACK) |
| `Before func(req string)` | nil | called before each answer is written, with the request's name (`QUERY`): a hook for barriers |
| `IdleLimit`, `FrameLimit`, `WriteLimit time.Duration` | 30 s, 10 s, 10 s | the controller's timeouts; tests set 50 ms |

The simulator moves the dish at once: POINT and JOG leave it IDLE at the new position. A test
reads its state only through the methods above. Its goroutines end when their connection or
listener closes, so it is safe inside `synctest.Test` over `net.Pipe`. Where a lesson's point is
the exact bytes on the stream (chapter 2), a test may play the controller itself, writing raw
frames from the table above.

#### dishd's code

Each lesson is its own `package main` (module `lesson`). Lessons give earlier chapters' pieces
as read-only files under these names and use only what they need. dishd uses the bible's pass
ID and station code; it adds the types below and renames nothing.

```go
// Frame is one controller-protocol frame, without its length and CRC.
type Frame struct {
	Type byte
	Seq  uint32 // a command's sequence number; 0 for requests and their answers
	Body []byte
}

// Frame types.
const (
	TypeHello   byte = 0x01
	TypePing    byte = 0x02
	TypeQuery   byte = 0x03
	TypePoint   byte = 0x10
	TypeStow    byte = 0x11
	TypeUnstow  byte = 0x12
	TypeJog     byte = 0x13
	TypeWelcome byte = 0x81
	TypePong    byte = 0x82
	TypeStatus  byte = 0x83
	TypeAck     byte = 0x90
	TypeNak     byte = 0x91
	TypeBye     byte = 0xFF
)

var (
	ErrBadLength = errors.New("bad length") // wrapped: "bad length 65535"
	ErrBadCRC    = errors.New("bad crc")
)

// NakError is a frame the controller refused.
type NakError struct {
	Seq    uint32
	Code   uint8
	Reason string // "wind lock"
}

func (e *NakError) Error() string // "seq 121 refused: wind lock"

// ByeError is the controller hanging up, and why.
type ByeError struct{ Reason string }

func (e *ByeError) Error() string // "controller closed: idle"

// DishState is what a dish is doing.
type DishState uint8

const (
	Idle DishState = iota
	Moving
	Tracking
	Stowed
	Fault
)

func (s DishState) String() string // "IDLE" ... "FAULT"; others "STATE(9)"

const (
	FlagHeater   uint16 = 1 << 0
	FlagWindLock uint16 = 1 << 1
	FlagFault    uint16 = 1 << 2
)

// Status is a decoded status frame: one dish at one moment.
type Status struct {
	Dish    uint8
	State   DishState
	AzDeg   float64 // 183.25
	ElDeg   float64 // 41.70
	WindKmh float64 // 64.0
	TempC   float64 // -12.5
	Flags   uint16
	LastSeq uint32
}

// Welcome is a decoded WELCOME frame.
type Welcome struct {
	Dish    uint8
	FwMajor uint8
	FwMinor uint8
	LastSeq uint32
}

// Client is dishd's connection to one dish controller.
type Client struct {
	Name    string        // from the banner: "YST-1"
	Timeout time.Duration // each read or write must make progress within it: 5s
	conn    net.Conn
	br      *bufio.Reader // every read goes through it, banner included, for the connection's life
}

// Result is one dish's answer to a poll.
type Result struct {
	Name   string // "YST-2"
	Status Status
	Err    error
}

// Watchers are the connections on the watch port. Safe for many goroutines.
type Watchers struct {
	Timeout time.Duration // how long one line may take to write to one watcher: 2s
	mu      sync.Mutex
	conns   map[net.Conn]bool
	closed  bool // set by CloseAll: no watcher joins after it
}

// Service is dishd: its dish clients and its watch port.
type Service struct {
	Dishes   []*Client // Dish 1 first
	Watchers *Watchers
	ln       net.Listener
	wg       sync.WaitGroup
	once     sync.Once
}

// Link keeps one controller connected through drops and numbers its commands.
type Link struct {
	Name     string
	Dial     func(ctx context.Context) (net.Conn, error)     // in tests, a fake line
	Wait     func(ctx context.Context, d time.Duration) error // sleeps, or ends with ctx
	Rand     *rand.Rand                                       // jitter (math/rand/v2)
	Timeout  time.Duration                                    // the Client's
	client   *Client
	next     uint32 // the next command's seq
	inflight *Frame // sent, not yet answered
	attempt  int    // failed dials in a row
}
```

Signatures every chapter after the one that defines them uses exactly:

```go
// chapter 1
func listenLocal() (net.Listener, error)
func greeting(names []string) string                // given: "dishd watching YST-1 YST-2\n"
func readBannerLine(conn net.Conn) (string, error) // given: the banner without CR LF, read a byte at a time
func serveOne(ln net.Listener, names []string) error
func pipePair(t testing.TB) (net.Conn, net.Conn)
func reach(name, addr string) string
func handshake(conn net.Conn) (name string, w Welcome, err error) // given in chapter 1
func firstContact(addr string) (string, error)
func bannerOf(addr string) (string, error) // 1.2: dials, reads the banner, hangs up

// chapter 2
func readBanner(br *bufio.Reader) (string, error) // the name
func readExactly(r io.Reader, n int) ([]byte, error)
func readPayload(br *bufio.Reader) ([]byte, error) // the bytes after the length: type to crc
func readFrame(br *bufio.Reader) (Frame, error)    // readPayload, then decodeFrame
func writeFrame(w io.Writer, f Frame) error        // encodeFrame, one Write
func newClient(conn net.Conn) (*Client, error)      // reads the banner; Timeout 5s
func (c *Client) send(f Frame) error
func (c *Client) recv() (Frame, error) // a BYE becomes *ByeError
func (c *Client) exchange(frames ...Frame) ([]Frame, error)
func (c *Client) Close() error
func validName(name string) bool // given in client.go: 1 to 16 capitals, digits and hyphens

// chapter 3
func encodeFrame(f Frame) []byte
func decodeFrame(payload []byte) (Frame, error)
func decodeStatus(body []byte) (Status, error)
func decodeWelcome(body []byte) (Welcome, error)
func pointBody(azDeg, elDeg float64) ([]byte, error)
func jogBody(axis byte, deltaDeg float64) ([]byte, error) // given from chapter 6
func nakError(f Frame) *NakError
func (c *Client) hello() (Welcome, error)
func (c *Client) status() (Status, error)
func (c *Client) command(f Frame) error // ACK nil, NAK *NakError

// chapter 4
type idleReader struct{ c *Client } // its Read sets the read deadline to now + c.Timeout, then reads
func describe(name string, timeout time.Duration, err error) string
func heartbeat(c *Client, ticks <-chan time.Time) error
func recvCtx(ctx context.Context, c *Client) (Frame, error)
func (c *Client) Do(ctx context.Context, f Frame) (Frame, error)

// chapter 5
func pollAll(clients []*Client) []Result
func watchLine(r Result) string
func serve(ln net.Listener, wg *sync.WaitGroup, handle func(net.Conn)) error
func (w *Watchers) Add(c net.Conn)    // closes c at once if CloseAll has run
func (w *Watchers) Remove(c net.Conn) // and closes it
func (w *Watchers) Broadcast(line string)
func (w *Watchers) CloseAll() // and marks the set closed
func (w *Watchers) Len() int
func (s *Service) Start(ln net.Listener)
func (s *Service) Shutdown()
func (s *Service) Run(ctx context.Context, ticks <-chan time.Time) error

// chapter 6
func redialDelay(attempt int, r *rand.Rand) time.Duration
func lost(err error) bool
func settle(cmd, answer Frame) error
func (l *Link) redial(ctx context.Context) (net.Conn, error)
func (l *Link) connect(ctx context.Context) error // dial, banner, HELLO, learn, resume
func (l *Link) learn(w Welcome)
func (l *Link) number(f Frame) Frame
func (l *Link) resume(w Welcome) *Frame
func (l *Link) Send(ctx context.Context, cmd Frame) error
func (l *Link) Drain(ctx context.Context, queue []Frame) ([]error, error)
```

- `exchange` writes every frame, then reads one answer per frame in order, and returns the
  answers read so far with the first error; a BYE is an error (`*ByeError`), never an answer.
- From chapter 4, `c.br` reads through an `idleReader`, so every read of a frame, banner
  included, has the client's idle deadline; `send` sets a write deadline of `Timeout`.
- `recvCtx` returns `ctx.Err()` at once for an ended context. Otherwise it registers
  `context.AfterFunc` to set the read deadline to now and then close a channel, calls `recv`,
  then `stop`; when `stop` reports false (the function has started, and `stop` doesn't wait for
  it) it waits for the channel, so a deadline set at a context's end never lands on a later
  read. A read that failed once ctx ended gives `ctx.Err()`. The deadline wakes a read that is
  waiting; if ctx ends between two reads of a frame whose bytes are still arriving, the next
  `idleReader` read sets a fresh deadline, and the frame is read to its end and returned.
- `describe` gives log lines: a deadline error is `YST-1: no answer in 5s` (the timeout as
  `time.Duration` prints it, `50ms` in tests); `io.EOF` and `io.ErrUnexpectedEOF` are
  `YST-1: closed`; anything else `YST-1: <err>`.
- `redialDelay(attempt, r)`, attempt counted from 1: 1 s doubled for each earlier failure,
  capped at 30 s, times a factor drawn uniformly from [0.5, 1.0) with `r.Float64()`. Attempt
  1 waits 0.5–1 s, 2 waits 1–2 s ... 5 waits 8–16 s; from attempt 6 on, 15–30 s. Never 0,
  and the doubling is capped before it can overflow.
- `redial` dials until a dial works, waiting `redialDelay` through `Wait` after each failure;
  its count, `l.attempt`, is of failed dials in a row, and a dial that works resets it.
  `connect` retries the whole handshake (dial, banner, HELLO) while the failure is `lost`,
  waiting `redialDelay` on its own count of failed handshakes in a row, until a HELLO
  succeeds or ctx ends.
- `Send` numbers a command (`number`), keeps it as `inflight`, and sends it with `Do`. On a
  lost line it closes the client, `connect`s again, and, if `resume` says the command already
  ran, returns nil without sending it; otherwise it sends it again with the same seq. `settle`:
  ACK of the same seq or NAK `duplicate` is nil; another NAK is its `*NakError`; anything else
  is an error. A `*NakError` is an answer, not a lost line.
- `Drain` sends a queue in order through `Send`, one outcome per command (nil or a
  `*NakError`); it stops early only when an error is neither (the context), returning the
  outcomes so far and that error.

#### The watch port

dishd's own line protocol, output only, lines ending in `\n` (LF alone):

- On connecting, a watcher gets one greeting line: `dishd watching YST-1 YST-2`.
- On every poll (each second; each tick in tests) one line per dish, Dish 1 first:
  `YST-1 TRACKING az=183.25 el=41.70 wind=64` (state as `DishState` prints it, az and el with
  two decimals, wind rounded to whole km/h), or for a dish whose poll failed
  `YST-2 DOWN no answer in 5s` (`describe`'s reason after the name).
- Watchers send nothing; anything they send is read and dropped, and a read error (they hung
  up) removes them. A watcher whose line takes longer than `Watchers.Timeout` (2 s, in tests
  too) is removed and closed: the watchers after it in the loop wait that long once, never
  again.
- `Shutdown` closes the listener, every watcher and both dish clients, then waits for every
  goroutine dishd started. A second `Shutdown` does nothing. `CloseAll` marks the set closed:
  a watcher accepted just before `Shutdown`, whose handler reaches `Add` only after `CloseAll`,
  is closed at once, so its goroutine ends instead of holding `Wait` forever.

#### The station server

A small HTTP/1.1 server written on `net` alone. It speaks only the subset below; anything else
is refused with the status given. In production a body streams to disk; in lessons it is read
whole into memory, under `MaxBody`.

**Routes.**

| Method and path | Answer |
|---|---|
| `GET /recordings` | 200 `text/plain; charset=utf-8`: one line per kept recording, `<pass ID> <bytes> <crc32>`, sorted by pass ID; an empty body when none |
| `GET /recordings/<pass ID>` | 200 `application/octet-stream`, the recording's bytes, with `Apogee-Crc32`; 404 `no recording furrow-a/YST/18856` |
| `HEAD` of either | the same status and fields as GET, Content-Length included, and no body |
| `PUT /recordings/<pass ID>` | a body by Content-Length or chunked; 201 `stored furrow-a/YST/18856 4096 bytes`, with `Location: /recordings/<pass ID>` and `Apogee-Crc32`; 409 `furrow-a/YST/18856 already stored` |
| `GET /feed` | 200 `text/plain; charset=utf-8`, chunked, `Connection: close`: one line per recording stored after the request arrived, as in the list, until the server stops |
| another known method on these paths | 405 `method not allowed` with `Allow: GET, HEAD` (`/recordings`), `Allow: GET, HEAD, PUT` (a recording), `Allow: GET` (`/feed`) |
| any other target | 404 `no such path` |

A pass ID in a path must have the bible's shape: lowercase satellite ID, `/`, three capitals,
`/`, five digits (`validPassID`, given). Targets carry no query; a `?` makes an unknown path.
`<crc32>` is the recording's CRC-32 (IEEE) in 8 lowercase hex digits.

**Reading a request**, in this order; the first failure answers:

1. **Lines.** Every line, request line and header lines alike, ends in CR LF. The server reads
   through a `bufio.Reader` of 1024 bytes (`bufio.NewReaderSize(conn, 1024)`) with
   `ReadSlice('\n')`, so a line, CR LF included, is at most 1024 bytes (a 1022-character line
   fits; 1023 doesn't). A line ending in LF without CR, or holding a CR anywhere else, is 400
   `bad line ending`.
2. **The request line**: `METHOD SP TARGET SP VERSION`, exactly three non-empty parts split by
   single spaces, the target starting with `/`, the version shaped `HTTP/<digit>.<digit>`;
   otherwise 400 `bad request line`. A request line over 1024 bytes is 414
   `request line over 1024 bytes`. A version other than `HTTP/1.1` is 505
   `version "HTTP/1.0" not supported`. A connection that ends before any byte of a request is
   not an error: the server closes quietly.
3. **Header fields**, until the blank line: `name:value`. The name is 1 or more token
   characters (letters, digits and `` !#$%&'*+-.^_`|~ ``), with nothing between it and the
   colon; the value has spaces and tabs trimmed from both ends and is otherwise kept as sent.
   A value may hold no control character but the tab: RFC 9110 §5.5 says a NUL must be refused
   or replaced, and Go's own server refuses them all (`isFieldValue`, given from 7.6). Names
   are stored in lower case; a repeated name keeps every value in order (`Header`). A line with
   no colon, an empty or crooked name, a value with a control character, or a line starting
   with a space or a tab (old line folding) is 400 `bad header line`. A header line over 1024
   bytes is 431 `header line over 1024 bytes`; a 33rd field is 431
   `more than 32 header fields`.
4. **Host** exactly once: 400 `missing host`, 400 `more than one host`.
5. **The method**: `GET`, `HEAD` or `PUT`, case-sensitive; anything else is 501
   `method "PATCH" not implemented`.
6. **The body's length**, checked in this order:
   1. Content-Length and Transfer-Encoding both present: 400
      `both content-length and transfer-encoding`. (Mirela: two programs reading one request
      two ways is how a second request gets hidden inside the first.)
   2. Transfer-Encoding: its values joined with `, ` must end in the coding `chunked`, in any
      letter case, else 400 `bad transfer-encoding "gzip"`, as RFC 9112 §6.3 requires; a
      coding before the final `chunked` is 501
      `transfer-encoding "gzip, chunked" not implemented`; on GET or HEAD, 400
      `GET has no body` (`HEAD has no body`); otherwise the body is chunked (length −1).
   3. Content-Length: once only (400 `more than one content-length`), 1 or more ASCII digits
      (400 `bad content-length "12a"`; signs and spaces inside are bad), at most `MaxBody`
      (413 `body over 65536 bytes`, also for a number too big to parse); on GET or HEAD only 0
      is allowed (400 `GET has no body`).
   4. Neither: no body, except PUT, which is 411 `length required`.

   Chapter 7's `contentLength` does steps 3 and 4; chapter 8's `bodyLength` adds 1 and 2.
7. **The body**, read through the same `bufio.Reader` that read the head (it may already hold
   the body's first bytes):
   - **Content-Length**: exactly that many bytes (`io.ReadFull`); a body that ends early is
     400 `body ended after 812 of 4096 bytes`.
   - **Chunked**: chunks of `<size in hex>[;extension]` CR LF, then exactly that many bytes,
     then CR LF; a size of 0 ends the body; then trailer fields until a blank line. The size is
     1 to 16 hex digits of either case, with spaces or tabs allowed before a `;`, and anything
     after `;` ignored; else 400 `bad chunk size "zz"`; a chunk line over 1024 bytes is 400
     `chunk line over 1024 bytes`. Data not followed by CR LF is 400 `bad chunk ending`. A
     total over `MaxBody` is 413, found as soon as a chunk's size would pass it, before
     reading that chunk. A stream that ends anywhere before the blank line after the last chunk
     is 400 `recording cut off`. Trailer lines follow the header rules (431 for the 33rd); only
     `apogee-crc32` is used.
   - **The CRC.** `Apogee-Crc32`, from the trailer if there, else from the header, is 8 hex
     digits (400 `bad apogee-crc32 "xyz"`). If present it must equal the CRC-32 of the bytes
     received: 400 `crc mismatch: sent a2912083, received a2912082`. Absent, nothing is
     checked. Nothing is stored after any failure.

The recorder sends `PUT` with a chunked body as the pass comes down, `Trailer: Apogee-Crc32` in
its head and `Apogee-Crc32: <8 hex>` after the last chunk: the CRC comes last, as a satellite's
downlink sends its file's CRC last. A chunk of HTTP is not a downlink chunk: lessons that need
both words say "HTTP chunk". Go's client sends exactly this, verified on Go 1.27.3:

```
PUT /recordings/furrow-a/YST/18856 HTTP/1.1\r\n
Host: 127.0.0.1\r\n
User-Agent: Go-http-client/1.1\r\n
Transfer-Encoding: chunked\r\n
Trailer: Apogee-Crc32\r\n
Accept-Encoding: gzip\r\n
\r\n
c\r\nfirst piece \r\n
c\r\nsecond piece\r\n
0\r\n
Apogee-Crc32: e2fc055d\r\n
\r\n
```

(The trailer's CRC is of the 24 bytes `first piece second piece`.) Go's `GET` is
`GET /recordings HTTP/1.1`, `Host`, `User-Agent: Go-http-client/1.1`, `Accept-Encoding: gzip`
and the blank line. The server ignores `Accept-Encoding` and `Expect`, and never sends
`100 Continue`. When a request's body fails, Go's client closes the connection after the
chunks it has sent, so the server sees a cut stream.

**Writing an answer.** One `bufio.Writer` per answer, flushed once (chunked answers flush
after every chunk):

```
HTTP/1.1 201 Created\r\n
Date: Tue, 20 Jun 2028 06:40:36 GMT\r\n
Content-Type: text/plain; charset=utf-8\r\n
Content-Length: 37\r\n
Location: /recordings/furrow-a/YST/18856\r\n
Apogee-Crc32: a2912082\r\n
\r\n
stored furrow-a/YST/18856 4096 bytes\n
```

- The status line is `HTTP/1.1 <code> <reason>` with reasons from the server's table (given,
  the same words `net/http` uses): 200 OK, 201 Created, 400 Bad Request, 404 Not Found,
  405 Method Not Allowed, 408 Request Timeout, 409 Conflict, 411 Length Required, 413 Request
  Entity Too Large, 414 Request URI Too Long, 431 Request Header Fields Too Large, 501 Not
  Implemented, 505 HTTP Version Not Supported.
- Then `Date` from the server's injected clock, in UTC, layout `Mon, 02 Jan 2006 15:04:05 GMT`
  (the constant `dateLayout`, given); then the answer's own fields in this order:
  `Content-Type`, `Content-Length` or `Transfer-Encoding: chunked`, then `Location`,
  `Apogee-Crc32` or `Allow`, then `Connection: close` when the server will hang up; then a blank
  line and the body.
- Every answer has a `Content-Length` except a chunked one. Every error answer is
  `text/plain; charset=utf-8`, its message and a newline, with `Connection: close`, and the
  server closes after writing it.
- A chunked answer writes each piece as `<size in lowercase hex>\r\n<bytes>\r\n`, flushes, and
  ends with `0\r\n\r\n`. Go's client shows it as `ContentLength` −1 and `TransferEncoding`
  `[chunked]`, and reads each piece as it comes.

**One connection, several requests.** After a 2xx answer (except `/feed`) to a request without
`Connection: close`, the server reads the next request on the same connection; Go's client
reuses it (three GETs, one accepted connection). It waits for the first byte of the next request
at most `Idle`, then closes quietly. Every other answer closes the connection. A request with
`Connection: close` (the token `close`, any letter case) gets `Connection: close` back.

**Deadlines.** Once a request's first byte arrives, its whole head must arrive within `Head`,
else 408 `request timeout` and close. Each read of a body must bring bytes within `Read`, else
408 and close, and nothing is stored. Each answer must be written within `Write`.

**The feed.** `GET /feed` subscribes to the store's arrivals (a channel with room for 16
lines). Each recording stored afterwards is one chunk holding its list line. A collector that
leaves is noticed at the next write; one that falls 16 lines behind is dropped (its channel
closed) and its feed ends with the last chunk, so it reconnects and reads `/recordings`. When the
server stops, every feed ends with its last chunk.

**Stopping** (`Shutdown(ctx)`, chapter 10): close the listener, so no new connection is
accepted; close connections waiting between requests; end every feed with its last chunk
(`Store.Close`); let requests in flight finish their answer, then close their connections
instead of reading another request; when all are done, return nil. If ctx ends first, close
everything and return ctx's error. `Serve` returns nil after `Shutdown`.

**Limits.** Zero means no limit, so a test sets only the limit it checks.

```go
// Limits are the station server's limits. Production values in comments.
type Limits struct {
	MaxBody int64         // 16 GiB (17179869184): a Furrow-A pass is about 11 GB. Tests: 65536
	Head    time.Duration // 5s: the whole head, from its first byte
	Read    time.Duration // 5s: each read of a body must bring bytes
	Idle    time.Duration // 60s: a kept-alive connection waits this long for the next request
	Write   time.Duration // 30s: each answer
}
```

**Types and signatures.**

```go
// Header is a request's fields: names in lower case, every value in order.
type Header map[string][]string

func (h Header) Get(name string) string // the values joined with ", "; "" when absent

// Request is one request as the station server reads it.
type Request struct {
	Method  string // "PUT"
	Target  string // "/recordings/furrow-a/YST/18856"
	Header  Header
	Length  int64  // 0 for no body, -1 for chunked
	Body    []byte // chapter 8: read by the server before routing
	Trailer Header // chapter 8: fields after a chunked body
}

// RequestError is a request the station server refuses, and its answer.
type RequestError struct {
	Status int    // 400
	Msg    string // "bad request line": the answer's body, without its newline
}

func (e *RequestError) Error() string // "400 bad request line"

// Field is one header field of an answer, written as given.
type Field struct{ Name, Value string }

// answer is what a route decides; serveConn writes it.
type answer struct {
	status int
	fields []Field // Content-Type, Content-Length and the rest, in order
	body   []byte  // nil for HEAD
	close  bool
	feed   bool // GET /feed: serveConn streams the feed instead
}

// Store keeps finished recordings by pass ID. Safe for many goroutines.
type Store struct { /* given */ }

var ErrExists = errors.New("already stored")

func NewStore() *Store
func (s *Store) Put(id string, data []byte) error // copies data; ErrExists; tells subscribers
func (s *Store) Get(id string) ([]byte, bool)
func (s *Store) List() []string                     // "<pass ID> <bytes> <crc32>", sorted
func (s *Store) Subscribe() (<-chan string, func()) // the func unsubscribes
func (s *Store) Close()                             // closes every subscription

// Server is the station server.
type Server struct {
	Store  *Store
	Limits Limits
	Now    func() time.Time
	// listener, connections (idle or busy) under a mutex, a WaitGroup, a shutting flag
}

// chapter 7
func readLine(br *bufio.Reader) (string, error) // errLineTooLong, 400 bad line ending
func parseRequestLine(line string) (method, target string, err error)
func readHeader(br *bufio.Reader) (Header, error)
func contentLength(method string, h Header, max int64) (int64, error)
func readRequest(br *bufio.Reader, max int64) (*Request, error)

// chapter 8
func readBody(br *bufio.Reader, n int64) ([]byte, error)
func chunkSize(line string) (int64, error)
func readChunked(br *bufio.Reader, max int64) ([]byte, Header, error)
func bodyLength(method string, h Header, max int64) (int64, error) // -1 for chunked
func checkCRC(body []byte, h, trailer Header) (uint32, error)
func (s *Server) put(req *Request) answer

// chapter 10
func statusLine(code int) string
func writeResponse(w io.Writer, now time.Time, status int, fields []Field, body []byte) error
func writeChunk(w *bufio.Writer, p []byte) error
func endChunks(w *bufio.Writer) error
func (s *Server) recording(req *Request) answer
func (s *Server) feed(bw *bufio.Writer, lines <-chan string) error
func (s *Server) route(req *Request) answer
func (s *Server) Serve(ln net.Listener) error
func (s *Server) serveConn(c net.Conn)
func closing(req *Request) bool // given from 10.5 (conn.go): Connection: close, any letter case
func (s *Server) Shutdown(ctx context.Context) error
```

`writeResponse` is given read-only in chapters 7 and 8 and written in chapter 10, as
`encodeFrame` is given in chapter 2 and written in chapter 3. In chapters 7 and 8 the given
`serveConn` answers one request per connection and closes, and the given route answers every
good request with the list. A chapter 8 lesson gives only what earlier lessons built, so no
editable file shows a later lesson's answer: chunked bodies reach the server from 8.4 (before,
the given `bodyLength` answers Transfer-Encoding with 501), the CRC from 8.8, and until 8.9 a
given `put` keeps any complete body and answers 201.

#### The live port (chapter 9)

From Monday 12 June the recorder has a **live port**: one UDP socket in the hut. It sends each pass's telemetry from there to the live view at Port Ellery, and answers questions there. Everything on it is ASCII, one message per datagram, with no line ending after the last line: a datagram's edges come with it.

- **Frames.** During a pass a satellite's housekeeping comes down once a second (Pipit-2's and the Furrows'), each frame numbered one past the last the satellite sent. The recorder sends each frame on as the bible's telemetry line as it records the pass. Each second's lines go out in as few datagrams as hold them (`pack`): whole lines in order, joined by LF, none over 512 bytes (`maxDatagram`). Nothing is sent twice or resent. A datagram lost on the line stays lost; the whole pass still reaches the station server in its recording.
- **Questions.** A datagram `STATUS <n>` (n from 1 to 4294967295, the asker's number) gets one datagram back, to the address it came from: `STATUS <n> RECORDING <pass ID>` while a pass is coming down, `STATUS <n> IDLE` otherwise. The recorder answers every question it receives, a repeated number included, and ignores anything else.
- **Why datagrams.** On a TCP connection one lost packet holds back every later byte until it is resent: at least a round trip of the island's line, and seconds over the backhaul. A live screen would freeze, then show old frames in a rush; a lost datagram costs only its own second.
- **Why 512 bytes.** A datagram bigger than one packet of the line (about 1,500 bytes) travels in fragments, and losing any one loses all of it; 512 bytes fit well inside one packet on any line the station uses. Tomás: "A jigsaw in three envelopes: lose one and there's no picture." macOS also refuses to send a UDP datagram over 9,216 bytes by default (`sendto: message too long`); no lesson sends one.
- **The live view** has two sockets: `Frames`, where the recorder sends, and `Questions`, which asks and receives the answers, so the two never mix. It reads each datagram into 513 bytes: one that fills them is over 512 and is dropped (`ErrTooBig`, "datagram over 512 bytes"); an empty datagram is a 0-byte datagram. It takes frames and answers only from the recorder's live port, comparing addresses by `String()`; an address can be forged, so this keeps out strays and mistakes, not liars. It shows a frame only when its number is past the last one shown (`Tracker`): repeats and older frames are dropped, and the numbers a frame jumps over count as missed. When no datagram has come for `Quiet` (5 s), the pass has gone quiet, and it asks the recorder its next question, numbered from 1 each run: each try waits `Wait` (1 s, above the backhaul's 0.7 s round trip), at most `Tries` (3), asking again under the same number; an answer to another number, or from another address, is skipped and the try's deadline stays; no answer after the last try is `silent`. It writes each frame shown, then `quiet after 84926: 8 shown, 2 missed, 3 dropped; recorder IDLE`, or `no frames; recorder IDLE`.
- **The week's numbers.** Furrow-A numbers frames as it sends them, about 8,000 a day across the network: Monday 12 June 07:38:02, 52817; Wednesday 08:32:10, 68841; Thursday 11:20:04, 77712; Friday 07:56:00, 84635, and 84917 to 84926 for 08:00:42 to 08:00:51. Furrow-B's 31877 to 31885 are the nine it owed at 08:06:10 on Tuesday. Pipit-2's frame at 06:37:15 on Monday is 30412. A Furrow frame's line is 72 bytes (Pipit-2's 70): seven fill 510 bytes, and Tuesday's nine in one datagram made 656.
- **Types and signatures** (as given in `live-in-the-fishbowl`'s starter): the `Tracker` and `Live` structs (`tracker.go`, `live.go`); `const maxDatagram = 512`; `var ErrTooBig = errors.New("datagram over 512 bytes")`; `readDatagram(pc net.PacketConn) ([]byte, net.Addr, error)`; `pack(lines []string) ([][]byte, error)`; `(t *Tracker) Accept(seq uint32) bool`; `question(n uint32) []byte` (given from 9.6); `answerTo(n uint32, d []byte, from, recorder net.Addr) (string, bool)`; `readAnswer(pc net.PacketConn, recorder net.Addr, n uint32) (string, error)` (given from 9.6); `ask(pc net.PacketConn, recorder net.Addr, n uint32, wait time.Duration, tries int) (string, error)`; `(l *Live) receive() ([]byte, error)`, `seqOf(line string) (uint32, bool)` and `lastLine(t Tracker, answer string) string` (given in 9.9); `(l *Live) Watch(out io.Writer) error`.

#### The code, by chapter

| Chapter | The learner writes (files) | Given read-only | Hidden |
|---|---|---|---|
| 1 | `listenLocal`, `serveOne`, `pipePair`, `reach`, `firstContact` (`main.go`, `pipe.go`) | `greeting`, `readBannerLine`, `handshake`, `Welcome` (`frames.go`) | `sim.go` |
| 2 | `readBanner`, `readExactly`, `readPayload`, `readFrame`, `writeFrame`, `Client`'s `recv` and `exchange` (`client.go`) | `Frame`, the types, the errors, `NakError`, `ByeError`, `encodeFrame`, `decodeFrame` (`frame.go`; this `decodeFrame` only splits a payload, and the CRC check is 3.7's); `newClient` | `sim.go` |
| 3 | `tempC` (3.2), `encodeFrame` (`encode.go`), `decodeFrame` (`decode.go`), `decodeStatus` (`status.go`), `pointBody` (`command.go`), `(*Client).command` | `client.go`, `DishState`, `rawStatus` (from 3.4), `decodeWelcome`, `nakError`, `hello`, `status` | `sim.go` |
| 4 | `describe`, `idleReader`, deadlines in `send` and `recv`, `heartbeat`, `recvCtx`, `Do` (`client.go`) | `frame.go`, `encode.go`, `decode.go`, `status.go`, `command.go` | `sim.go` |
| 5 | `pollAll`, `serve`, `Broadcast`, `Shutdown`, `Run` (`poll.go`, `watch.go`, `service.go`) | the client, `watchLine`, `greeting`, `Watchers`' other methods | `sim.go` |
| 6 | `redial`, `lost`, `learn`, `number`, `resume`, `Send`, `Drain` (`link.go`) | the client, `jogBody`, `redialDelay` (from 6.2), `settle`, `connect` (from 6.8) | `sim.go` |
| 7 | `readLine`, `parseRequestLine`, `readHeader`, `contentLength`, `readRequest` (`request.go`) | `Header`, `Request`, `RequestError`, `errLineTooLong`, `versionShaped`, `isToken` and `isFieldValue`, given in the editable `request.go`; `writeResponse`, `dateLayout`, `Store`, `validPassID`, `Server`, `serveConn` (`response.go`, `store.go`, `server.go`) | |
| 8 | `readBody`, `chunkSize`, `readChunked`, `bodyLength`, `checkCRC`, `put` (`body.go`, `routes.go`) | `request.go` and chapter 7's files | |
| 9 | `readDatagram` (`datagram.go`), `pack` (`pack.go`), `Tracker`'s `Accept` (`tracker.go`), `answerTo` (`answer.go`), `ask` (`ask.go`), `Live`'s `Watch` (`watch.go`) | `maxDatagram` and `ErrTooBig` (`datagram.go`), `Tracker` (`tracker.go`), `question` and `readAnswer` (`question.go`, from 9.6), `Live`, `receive`, `seqOf` and `lastLine` (`live.go`, 9.9) | |
| 10 | `statusLine`, `writeResponse`, `writeChunk`, `endChunks`, `recording`, `feed`, `serveConn`, `Shutdown`, `route` (`response.go`, `server.go`, `routes.go`) | `request.go`, `body.go`, `store.go`, and `closing` (`conn.go`, from 10.5) | |

#### How lessons test the wire

- Every connection is `net.Pipe` or a listener from `net.Listen("tcp", "127.0.0.1:0")`. Never a
  fixed port, never `localhost` or `:0` (all interfaces: on macOS it can ask the firewall for
  permission), never another address. Nothing reaches the network.
- Every test sets a deadline on its own end of every connection (one second) and gives the code
  under test a timeout of 50 ms where the lesson is about a timeout, so a wrong starter fails
  within a second or two and never hangs the run. HTTP clients in tests have `Timeout` one
  second, or are `http.Client{Transport: ...}` with a dial to the test listener and the same
  timeout.
- `net.Pipe` has no buffer: a Write waits until the other end reads, so a test writes from a
  goroutine or with a deadline. Closing errors differ: `io.ErrClosedPipe` on a pipe,
  `net.ErrClosed` on TCP; a read after the other end closed is `io.EOF` on both.
- **Time on the line.** Where a lesson's point is how time passes on a connection (chapters 4
  and 6), the test runs inside `synctest.Test` over `net.Pipe`, with production values: a 5 s
  idle deadline, gaps of 0.4 s between bytes, fourteen minutes of backoff, all in no real time
  (checked on Go 1.27.3). Inside a bubble: no loopback connection (it would never look idle),
  and everything the bubble starts ends before the test function returns: close the
  connections, then wait for the goroutines. Outside bubbles, deadline tests use a peer that
  never sends or never reads, so the deadline is the only way out; a "slow but alive" test
  outside a bubble uses gaps of 5 ms against a 50 ms idle limit.
- **Goroutines that must meet.** Chapter 5 runs a test in a bubble too wherever its
  goroutines must meet in a set order, or a timeout must not fire for right code, or a starter
  would fail only after a real wait: 5.1's barrier (which also releases the dishes' answers in
  reverse order), 5.3's three watchers, 5.5's watchers (2 s timeout, readers taking a byte every
  30 ms, so a deadline shared by two of them runs out in every map order), 5.6's shutdown and
  its watcher who joins as dishd stops, and 5.8's whole service. The watch port is then a test
  listener whose `Accept` hands out `net.Pipe` ends; once closed it fails with `net.ErrClosed`,
  and a second call ends the calling goroutine, so a loop that never stops fails instead of
  spinning (a spinning goroutine would stop the fake clock). A watcher that must keep reading
  has a reader goroutine of its own (a pipe Write waits for its reader), and `synctest.Wait`
  tells the test that every handler has added its watcher. 5.3 and 5.6 each keep one short
  test on a real loopback listener, outside bubbles, that a starter fails or passes at once.
- Backoff waits are injected (`Link.Wait`) and recorded by the test, or run in a bubble;
  randomness is `rand.New(rand.NewPCG(1, 2))`; clocks are `now func() time.Time`; real sleeps
  are 50 ms at most; concurrency checks are correct by construction (barriers through the
  simulator's `Before`, handlers that block until the test releases them), never by timing.
- Readers that split their bytes come from `testing/iotest` (`OneByteReader`, `HalfReader`,
  `ErrReader`, `TimeoutReader`); writers that fail or record each Write are small test types.
- Frames are compared byte for byte, messages in hex with spaces:
  `encodeFrame(PING) = 00 09 02 ..., want 00 09 02 ...`; text with `%q`; statuses as
  `PUT /recordings/furrow-a/YST/18856 = 400, want 201`.
- Recordings in tests are `testRecording(n)`, byte i being `byte(i % 256)`, so they hold 0x00,
  CR and LF: n 24 has CRC `8295a696`, 1024 `b70b4c26`, 2048 `9f5edd58`, 4096 `a2912082`,
  65536 `b11de6a1`. The finale's are Furrow-A's `furrow-a/YST/18856` (4096 bytes, `a2912082`)
  and Cirrine-1's `cirrine-1/YST/14302` (2048 bytes, `9f5edd58`).
- Program output (`code` lessons) never prints a port, a time or anything that depends on how
  reads split: requests printed in chapter 7 go to `http://127.0.0.1/...` through a Transport
  that dials the test listener, so the Host line reads `127.0.0.1`.
- Each test ends by closing every listener and connection it opened and waiting for the
  goroutines it started, so nothing leaks into the next test.
- **Datagrams** (chapter 9) stay on real sockets from `net.ListenPacket("udp", "127.0.0.1:0")`, never in a bubble: a goroutine waiting on a socket isn't durably blocked, so the bubble's clock would never move. Loss, repeats and reordering come from the test's own sends; loopback itself loses nothing and keeps their order (checked on Go 1.27.3: a thousand datagrams in a row arrive whole and in order, a datagram to a closed port vanishes without an error, and a short buffer cuts a datagram with no error). Questions are counted by a wrapper that records each WriteTo. A check that an answer arrives within a try asserts only what any timing allows. Tests use Wait 50 ms and Quiet 100 ms.

**Topics added this season** (front matter): `networking`. Lessons keep using the older ones
(`bytes`, `context`, `errors`, `goroutines`, `http`, `mutexes`, `streams`, `testing`, `time`...).

**New field guide words** (world words only; Go, network and HTTP words are taught in
readmes): backhaul, banner, beamwidth, collector, dish controller, dishd, feed, fiber, firmware,
jog, live port, live view, midnight sun, polar night, recorder, recording, simulator, station server, the hut, the
manual, watch port, wind lock. Carried over and needed again: AOS, LOS, pass, pass ID, radome,
stow, CRC, sequence number, telemetry, X-band, S-band, UTC, the Board, the Fishbowl.

**Name checks** to run before the season ships: `dishd` (a quick search on 6 November 2026 found
no program of that exact name; similar names are unrelated), and the header name `Apogee-Crc32`,
which uses the invented company's name as Season 3's headers did. `YST-1` and `YST-2` are a
station code and a dish number, the simulator and the maker are never named, and "the hut" and
"the station server" are plain words. "The live view" and "the live port" (web search, 13 November 2026): nothing in ground-station or telemetry software goes by either; the one-word "LiveView" names unrelated products, and the world writes both as lowercase plain phrases with "the", like "the hut", so they stay. No new people, places, satellites or customers.
