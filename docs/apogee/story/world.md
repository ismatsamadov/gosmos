# Apogee story bible: the world

## The company

**Apogee** sells ground-station time by the pass. Small-satellite operators book contacts
through Apogee's API; Apogee predicts the passes, points a dish, sends the operator's signed
commands up, brings their data down and delivers it within minutes of loss of signal.

- **Mission:** "Every pass counts." A missed pass is a satellite nobody heard from for an
  orbit, and a customer whose data is ninety minutes late.
- **Size:** 38 people. Founded 2023 in Port Ellery by Tomás Ferreira-Lund (now CTO) and
  a CEO who never appears on stage.
- **Network:** five ground stations, from the Arctic to the far south.
- **Fleet:** three satellites of its own (the Pipits, used to test every change before a
  customer sees it) and five customer satellites under contract.
- **Rooms:** the ops room at headquarters is **the Fishbowl** (glass on three sides). Its
  big screen, listing the next passes, is **the Board**.
- **Culture:** everything runs in UTC; checklists before cleverness; a failure gets a name,
  a reason and a line in the log. New engineers' first program writes the message the ops
  room sends up to Pipit-1, which echoes it back on its beacon.

**The learner** is the newest software engineer, never named in lessons, always "you". Their
badge says Ground Software Engineer, No. 039 (38 people were there before them), based at
Port Ellery headquarters, since Monday 1 March 2027; the app shows it with the name the learner
gives. Story time runs from Monday 1 March 2027 (day one) to Friday 9 July 2027 (first solo
on-call week), nineteen weeks. Each course chapter is roughly a week. Pipit-3 still launches
on 5 April, in chapter 6. The following week, chapter 7 turns the go/no-go poll's rules into
values the team can hand around. Chapters 8 to 19 run from 19 April to 9 July. Season 2
starts later in July, after the first solo on-call week.

In plain words, for the opening and anywhere a newcomer reads first: Apogee is a small space
company that keeps satellites in touch with the ground. Its five ground stations are big radio
dishes; a satellite is in a dish's view for about ten minutes at a time, and that window is a
pass.

## Ground stations

| Code | Name | Lat, lon | Alt | Mask | Dishes and bands | Character |
|---|---|---|---|---|---|---|
| YST | Ystmark | 77.8° N, 2.0° E | 40 m | 5° | two 3.7 m dishes under radomes; S, X | Arctic island; sees every polar orbit on almost every revolution |
| PEL | Port Ellery | 46.9° N, 38.1° W | 25 m | 5° | one 3.0 m dish, one UHF yagi array; UHF, S | Headquarters' rooftop; the test station |
| SQT | Salar Quieto | 23.4° S, 104.6° W | 3,840 m | 3° | a 5.4 m and a 2.4 m dish, a receive-only UHF yagi; S, X, UHF to listen | High salt flat; flat horizon, dry air |
| LAN | Lantern Atoll | 1.8° S, 82.5° E | 3 m | 8° | one 3.7 m dish, UHF yagis; UHF, S | Equatorial; best for low-inclination orbits |
| CWR | Cape Wren | 53.9° S, 95.0° W | 110 m | 5° | one 3.7 m dish, a receive-only UHF yagi; S, X, UHF to listen | Remote, unstaffed; windiest site |

Coordinates are invented for exercises and deliberately land in open ocean: don't map them.
Latitude is positive north, longitude positive east.

Only Port Ellery and Lantern Atoll send on UHF. Salar Quieto and Cape Wren have receive-only
UHF yagis, enough to hear a beacon such as Thistle's or Pipit-1's but never to command a
satellite; Ystmark has no UHF at all.

- **Ystmark.** Polar night from November to February; radome heaters; arctic foxes chew
  cables, so the conduits are fox-proof. The fiber to the mainland fails in storms, and the
  station falls back to a slow satellite backhaul. Storms stow the dishes: like Cape Wren's,
  they stow when the wind reaches 72 km/h.
- **Port Ellery.** Sea fog, gulls on the yagis, and the only station you can walk to. Every
  new piece of software runs here first, against a Pipit.
- **Salar Quieto.** At 3,840 m, people walk slowly and dust gets everywhere. Its 3° mask is
  the lowest in the network; huge day-night temperature swings make the dish's pointing
  drift until it is recalibrated. Its two dishes are nicknamed La Grande and Chiquita. La
  Grande also carries a Ka-band test radio, 300 Mbps, which corrects Doppler shift itself; it
  is a test, not a band customers book. The only road up is the Quebrada Añil.
- **Lantern Atoll.** Tropical downpours fade the X-band link (rain fade), so it has none.
  Diesel power, salt corrosion, a caretaker on site. A high horizon of palm trees gives it
  the network's highest mask.
- **Cape Wren.** Nobody lives there; a contractor visits monthly. The dish stows itself when
  the wind reaches 72 km/h, and tracks with a warning from 50 km/h.
- **The test van.** A mobile station based at Port Ellery, with one 1.2 m dish, S-band only.
  It carries the dish to sites for calibration.

## The fleet

| ID | Name | Operator | Orbit | Radio | Payload and story |
|---|---|---|---|---|---|
| `pipit-1` | Pipit-1 | Apogee | 525 km, 97.5° (sun-synchronous) | UHF beacon 437.250 MHz, 9.6 kbps | 3U test cubesat; echoes uplinked messages on its beacon. Since a cold night in 2025 its beacon sits about 1.2 kHz high |
| `pipit-2` | Pipit-2 | Apogee | 550 km, 97.6° (SSO) | S-band 2245 MHz down, 2065 MHz up, 4 Mbps | 6U software-defined-radio testbed; the first satellite to get signed commands |
| `pipit-3` | Pipit-3 | Apogee | 510 km, 97.4° (SSO) | S-band, X-band test | 6U; launches on 5 April 2027 (chapter 6); its first acquisition is a Doppler hunt |
| `furrow-a` | Furrow-A | Furrowlight | 510 km, 97.4° (SSO) | X-band 8150 MHz, 150 Mbps; S-band control | 16U multispectral crop imager; big files, demanding customer |
| `furrow-b` | Furrow-B | Furrowlight | 510 km, 97.4° (SSO) | as Furrow-A | Same plane as Furrow-A, half an orbit behind |
| `brine-3` | Brine-3 | Brinewatch | 480 km, 53.0° | S-band 2230 MHz, 1 Mbps | 6U ship-tracking receiver. Never rises over Ystmark (its orbit doesn't reach far enough north). Brine-1 was lost at launch; Brine-2 never left its deployer |
| `thistle` | Thistle | Kessick Polytechnic | 415 km, 51.6°, decaying | UHF beacon 435.900 MHz, 1200 bps | 2U student cubesat with a radiation counter. Like Brine-3, never rises over Ystmark. Sends garbage frames from chapter 13; re-enters in chapter 18 |
| `cirrine-1` | Cirrine-1 | Cirrine Atmospherics | 600 km, 97.8° (SSO) | S-band 2280 MHz, 10 Mbps | 12U weather sounder (GNSS radio occultation); the customer that almost got a replayed command |

Useful numbers (with Earth radius 6371 km and μ = 398600.4418 km³/s²):

| Altitude | Period | Orbits per day | Slant range at 0° / 10° / 90° |
|---|---|---|---|
| 415 km | 92.72 min | 15.53 | — / 1479.05 km / 415 km |
| 525 km | 94.99 min | 15.16 | 2639.16 km / 1755.35 km / 525 km |
| 550 km | 95.50 min | 15.08 | — / 1815.08 km / 550 km |
| 600 km | 96.54 min | 14.92 | — |

A range rate of 7 km/s shifts a 437.25 MHz beacon by about ±10.2 kHz; a range rate of
6.5 km/s shifts a 2245 MHz downlink by about ±48.7 kHz.

Customers, all invented: **Furrowlight** (crop imaging), **Brinewatch** (ship tracking),
**Kessick Polytechnic** (a university), **Cirrine Atmospherics** (weather).

## People

Seven recurring characters. A story beat uses one of them, rarely two. First names in lessons.
Lessons already fix their pronouns, so keep them: Nkem, Ingrid, Paz, Yusra and Mirela are
"she", Tomás, Joss and Hollis "he". Sunniva and Professor Marr have none in a lesson yet:
write around one, or add it here first.

**Nkem Halloran**, staff backend engineer, your onboarding buddy and reviewer. Short, kind,
exact. Asks one question at a time and likes smaller code.
> "Works. Now make it say what it means."
> "What does the satellite see if this line runs twice?"

**Ingrid Achterberg**, flight operations lead, runs the Fishbowl. Ops voice: UTC times,
checklists, "Copy," "Negative." Calm under pressure; owns the Board and the go/no-go poll.
> "AOS 04:12. Checklist done by 04:02. Questions after LOS."

**Tomás Ferreira-Lund**, co-founder and CTO, a radio engineer. Enthusiastic, explains
physics with everyday images, keeps the orbit formulas in a battered notebook.
> "It's coming at you at seven kilometers a second. Of course its voice sounds higher!"

**Sunniva Ruud**, Ystmark station chief. Telegraphic, weather first, because the backhaul
is slow and so is everything else up there.
> "Wind 64, gusts 81. Dish 1 stowed. Coffee holding."

**Paz Ormaechea**, station engineer at Salar Quieto. Chatty voice memos, nicknames for
everything, jokes about the altitude.
> "Hola from 3,840 meters, where even the bits are short of breath."

**Yusra Belkadi**, customer success. Turns customer panic into requirements; every message
starts with a summary line.
> "TL;DR: Furrowlight wants every missing chunk listed by number before nine."

**Mirela Dunn**, security engineer, guardian of the uplink. Skeptical, precise, never
unkind. Signs off reviews with "LGTM once it fails closed."
> "Your code accepts the command. Now show me the one it refuses."

Minor characters, in a sentence at most: **Professor Elspeth Marr** and her student
**Joss Imbert** (Kessick Polytechnic, Thistle's team); **Hollis Vey** (Furrowlight's ops
manager, writes subject lines in capitals).

## The codebase

The learner's code grows into **groundctl**, the Go module behind Apogee's ground segment.
Early lessons are one `package main` file; the types below appear in the chapter that
teaches what they need, and chapter 16 splits them into packages.

- **Module:** `apogee/groundctl`. It has no dot in its first element on purpose: the app
  treats such imports as local and never runs `go mod tidy` for them.
- **Services** (named in stories, never built whole): `passplan` (predicts passes, builds
  the schedule, feeds the Board), `tlmd` (telemetry ingest and validation), `uplinkd`
  (signs and sends commands), `dlinkd` (downlink reassembly and delivery), `pager`
  (anomaly rules and on-call pages), and `gctl`, the engineers' command-line tool.
- **Packages:** `orbit`, `fleet`, `schedule`, `telemetry`, `uplink`, `downlink`, `alert`,
  `board`.

The canonical types. A lesson uses only the fields it needs, but never renames one or
changes its type.

```go
// Band is a radio band.
type Band string

const (
	UHF   Band = "UHF"
	SBand Band = "S"
	XBand Band = "X"
)

// Station is one of Apogee's ground stations.
type Station struct {
	Code       string  // "YST": three capital letters
	Name       string  // "Ystmark"
	LatDeg     float64 // north is positive
	LonDeg     float64 // east is positive
	MaskDeg    float64 // lowest usable elevation
	Bands      []Band
	Stowed     bool
	StowReason string
}

// Orbit is a circular orbit.
type Orbit struct {
	AltKm  float64
	IncDeg float64
}

// Mode is a satellite's operating mode: NOMINAL, SAFE, DETUMBLE or LEOP.
type Mode string

// Satellite is a spacecraft Apogee talks to.
type Satellite struct {
	ID       string // "pipit-1"
	Name     string // "Pipit-1"
	Operator string // "Apogee", "Furrowlight", ...
	Orbit    Orbit
	Bands    []Band
	Mode     Mode
}

// Pass is one contact between a satellite and a station.
type Pass struct {
	SatID      string
	Station    string    // station code
	AOS        time.Time // acquisition of signal
	LOS        time.Time // loss of signal
	MaxElevDeg float64
}

// TelemetryFrame is one decoded housekeeping frame.
type TelemetryFrame struct {
	SatID string
	Seq   uint32
	At    time.Time
	BattV float64
	TempC float64
	Mode  Mode
}

// Command is one telecommand on its way up.
type Command struct {
	SatID string
	Seq   uint64 // strictly increasing per satellite
	Op    string // "PING", "SET_MODE", "DUMP_FILE", "REBOOT"
	Args  []string
	Sig   string // hex HMAC-SHA256 of the canonical form
}

// Downlink is a file coming down in numbered chunks.
type Downlink struct {
	SatID  string
	File   string
	Total  int            // number of chunks
	Chunks map[int][]byte // by chunk number, from 0
	CRC32  uint32         // of the whole file
}

// Severity ranks an anomaly.
type Severity int

const (
	Info Severity = iota
	Warning
	Critical
)

// Anomaly is something the pager should hear about.
type Anomaly struct {
	SatID    string
	Rule     string
	Severity Severity
	At       time.Time
	Detail   string
}
```

Wire formats, used the same way in every lesson:

- **Telemetry line:** satellite, sequence number, UTC time, then `key=value` fields, one
  space apart:
  `thistle 4182 2027-05-11T04:13:07Z batt_v=7.91 temp_c=-4.5 mode=NOMINAL`
- **Pass ID:** satellite, station, orbit number padded to five digits: `pipit-1/PEL/04182`.
- **Canonical command** (what gets signed): fields joined by `|`, arguments last:
  `pipit-2|118|SET_MODE|SAFE`. On the wire:
  `CMD pipit-2 118 SET_MODE SAFE sig=<64 hex digits>`.
- **Pipit beacon:** 12 bytes. 0–1 the sync bytes `"AP"` (0x41 0x50); 2 the satellite number;
  3–4 sequence (big-endian uint16); 5 mode code; 6–7 battery in millivolts (big-endian
  uint16); 8–9 temperature in tenths of °C (big-endian int16); 10 flags; 11 the XOR of
  bytes 0–10.
- **Downlink chunks** are numbered from 0; the file's CRC-32 (IEEE) is sent last.

Limits ops uses everywhere: battery FULL from 8.2 V, OK from 7.2 V, LOW from 6.6 V, below
that CRITICAL (safe mode follows); Cape Wren tracks with a warning from 50 km/h of wind and
stows at 72 km/h; battery temperature limits −10 °C to 45 °C.

## Seasons

Each course of the series is a season of the learner's career at Apogee, and story time only
moves forward. Each season's opening gives the learner a new role on their badge. The roadmap
is ismatsamadov/gosmos#18; how each season is told is `docs/storytelling.md`.

| Season | Course | Story time | The season in a sentence | Issue |
|---|---|---|---|---|
| 1 | Go Foundations | Mar–Jul 2027 | From your first message to Pipit-1 to your first solo week on call | #7, #8 |
| 2 · Tools for the Van | Files, Data and the Command Line | Jul–Oct 2027 | The test van's trip up the Quebrada Añil, with no signal: gctl, from files alone | #9 |
| Project | Build gctl | Oct 2027 | The night shift's interactive tool | #10 |
| 3 · Open for Business | Web Services in Go | Nov 2027–Feb 2028 | Customers book passes through Apogee's own service, which you build | #11 |
| Project | Build the Board | Feb 2028 | The Board on every screen in the building | #12 |
| 4 · Down to the Wire | Networking from the Socket Up | Mar–Jun 2028 | Ystmark's new dish controllers, and a web server from nothing but a socket | #13 |
| 5 · Keeping Secrets | Cryptography in Go | Jul–Nov 2028 | Cirrine's commands, readable by Cirrine alone, and Rotation Day | #14 |
| 6 · On Call | Logging, Metrics and Profiling | Nov 2028–Feb 2029 | The 02:46 page from Lantern Atoll, and the slow leak behind it | #15 |
| 7 · Scale | Go at Scale | Mar–Jun 2029 | Twenty stations, eighty satellites, and passplan's rewrite | #16 |
| Capstone | Pass Day, Live | Jul 2029 | The whole ground segment, rehearsed for a week | #17 |

## Recurring situations

| Situation | What happens | Chapters |
|---|---|---|
| Pass prep | Pre-pass checklist, go/no-go poll, a countdown to AOS, reusable poll rules | 1, 5, 6, 7 |
| The Board | Making the ops room's screen readable and right | 3, 9, 12 |
| Launch week | Pipit-3's first acquisition: a Doppler hunt for a faint beacon | 6 |
| Weather | Ystmark storms and Cape Wren wind stow dishes; passes move | 5, 11, 19 |
| Thistle's decline | Garbage frames, duplicates, then a last pass before re-entry | 13, 14, 18 |
| Furrowlight escalations | Missing chunks, data volumes, a pass that moved | 2, 9, 10, 18 |
| Uplink security | Mirela's reviews, signatures, replays, mutants in the tests | 14, 15, 16 |
| On call | A page at 03:00, a stuck uplink, a cache gone wrong | 19 |

## Glossary

Gloss a term in a short clause the first time a lesson uses it, then use it freely.

| Term | Meaning |
|---|---|
| AOS | Acquisition of signal: the satellite rises into a station's view and the link starts |
| LOS | Loss of signal: the satellite sets and the link ends |
| Pass | One AOS-to-LOS window of one satellite over one station; also called a contact |
| Elevation | The satellite's angle above the horizon, 0° to 90° |
| Azimuth | Compass direction to the satellite, 0° north, clockwise |
| Mask | The lowest elevation a station can use (trees, hills, noise) |
| Max elevation | The highest point of a pass; low passes are short and weak |
| Slant range | Straight-line distance from dish to satellite |
| Footprint | The patch of Earth from which a satellite is above the mask |
| Inclination | Tilt of the orbit to the equator; it caps how far north or south the ground track reaches |
| SSO | Sun-synchronous orbit: near-polar, passes over each place at the same local time |
| Period | Time for one orbit (about 95 minutes for these satellites) |
| LEOP | Launch and early orbit phase: the first days after launch |
| Beacon | A short status message a satellite broadcasts without being asked |
| Telemetry (TLM) | Housekeeping data coming down: battery, temperatures, mode |
| Telecommand | A command going up; at Apogee, always signed |
| Uplink / downlink | Ground to space / space to ground |
| Frame | One unit of telemetry on the wire |
| Sync bytes | Fixed bytes that mark where a frame starts |
| Checksum, CRC | A short number computed from data to catch corruption |
| Doppler shift | The frequency change caused by the satellite moving toward or away from you |
| Range rate | How fast the slant range changes; negative while the satellite approaches |
| Link budget | The sum of gains and losses that says whether a signal will be strong enough |
| Path loss | How much a signal weakens over distance, in decibels |
| Margin | Received signal minus what the receiver needs; positive is good |
| dB, dBm | Decibels: a ratio, and a power relative to one milliwatt |
| Rain fade | Rain absorbing high-frequency (X-band) signals |
| Stow | Park a dish safely, usually for wind |
| Safe mode | A satellite's minimal state after a fault: beacon on, payload off |
| Detumble | Stopping a satellite's rotation, typically just after deployment |
| Sequence number | A counter in each frame or command; gaps and repeats mean trouble |
| Replay | Resending an old, validly signed command; counters stop it |
| HMAC | A signature made with a shared secret key; Apogee uses HMAC-SHA256 |
| Page | An alert that wakes the on-call engineer |
