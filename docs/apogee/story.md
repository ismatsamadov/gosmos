# Apogee: the story bible

Every Gosmos course in the Apogee series follows this bible. It fixes the world (company,
stations, satellites, people), the code the learner grows, and how lessons are written. When
a lesson needs something this file doesn't have, add it here first, then use it. How the
story is told, from the opening to the field guide, is `docs/storytelling.md`.

Everything in this world is invented: the company, its customers, the people, the places and
the satellites. Only the physics is real.

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

### Season 2: tools for the van

After the first solo on-call week ends on 9 July 2027, the learner starts the next
assignment on Monday 12 July. Their new badge says Ground Tools Engineer, still No. 039.
Paz takes the test van up the Quebrada Añil to check and correct the pointing of Salar
Quieto's dishes. Nkem reviews the tools before departure and as recordings return to
headquarters. The trip and calibration work occupy July to October; no new people,
places, satellites or customers are introduced for this season.

The command-line program **gctl** belongs to `apogee/groundctl`. It reads the same station
codes, satellite IDs, UTC timestamps and pass records as Season 1. Files contain schedules
and reports, never live radio access: the exercises work offline with recorded inputs.
Chapter 1 uses a plain-text schedule with one record per line; structured CSV and JSON
formats are introduced in their own chapters. Day labels such as `2027-07-12` come from
the program's own validated configuration, rather than arbitrary user-supplied paths.

Reports live under `reports/daily/<day>.txt`. Whole-file replacement stages a private
temporary sibling, checks write and close errors, then renames it into place; it promises
complete contents for readers, not durability after loss of power. Filesystem tests use
`t.TempDir`, and stream tests use in-memory readers and writers. No course test writes to
a real user's files. A command's behavior is tested through
`run(args []string, stdin io.Reader, stdout, stderr io.Writer) int`, never by reading the
test process's own arguments or terminating it. Exit statuses are 0 for success, 1 for an
operational failure and 2 for invalid usage; successful help returns 0.

#### The road, chapter by chapter

You ride in the van with Paz, at the laptop in the back; Nkem reviews from Port Ellery
whenever the van finds a signal, and Paz's voice memos go out at the same stops. The van
climbs the Quebrada Añil in stages, with days at a time and no signal, and story time moves
forward a week or two per chapter. Lessons give an altitude or a day, never a new place
name.

| Ch | Week of (2027) | Where the van is | What the team needs |
|---|---|---|---|
| 1 | 12 Jul | Port Ellery, packing | the week's passes on a memory stick, and safe saves |
| 2 | 19 Jul | the foot of the road, first stop | the recorder's report on the screen and in the archive |
| 3 | 26 Jul | 1,200 m | La Grande's controller log: 200,000 lines Paz brought down, the record of its drift |
| 4 | 9 Aug | 1,900 m | Yusra's booking export, so calibration never takes a booked pass |
| 5 | 23 Aug | 2,400 m | Furrowlight's delivery reports, which change shape, and gctl's settings file |
| 6 | 6 Sep | 2,800 m | Hollis's spreadsheets, in capitals |
| 7 | 20 Sep | 3,200 m | the van's laptop clock on local time, the network on UTC |
| 8 | 4 Oct | 3,500 m | gctl becomes one command with subcommands, instead of the eleven scripts on the van's laptop |
| 9 | 11 Oct | 3,700 m | gctl carries the station table and report templates inside itself |
| 10 | 18 Oct | 3,840 m, Salar Quieto | calibration day, run from gctl and files alone; La Grande points true |

Calibration compares La Grande with a reference it can trust: the van's 1.2 m dish, checked
at Port Ellery before the trip, tracks the same passes, and the difference between the two
dishes' pointing is La Grande's error. Paz calls La Grande "the sunflower": it wanders most in
the afternoon heat. The van's laptop keeps the salt flat's local time, UTC−7 all year (no
daylight saving), written `-07:00` and made with `time.FixedZone("", -7*60*60)`; everything
gctl stores is UTC.

The finale: the van climbs the last 140 m on Monday 18 October; Nkem reviews on Wednesday
evening; calibration day is Thursday 21 October, from 06:00 on the salt flat (13:00 UTC), and
La Grande points true at 23:40 UTC. The van starts home on Friday 22 October and is back at
Port Ellery on Saturday 23 October.

Smaller facts the lessons use: for one day of sweeps along the horizon, Salar Quieto's mask is
set to 0° (`{"station":"SQT","mask_deg":0}`), then back to 3°; Yusra's older export tool writes
`null` for an empty day; Furrowlight's delivery reports changed shape twice this summer; in
chapter 7 Paz raises gctl's settings timeout from `30s` to `2m` for the long sweeps ahead.
Booking IDs used this season: `bk-2041` to `bk-2045` (August, chapter 4), `bk-2071` to
`bk-2073` (24–25 August), `bk-2117` to `bk-2136` (7–9 September), `bk-3099` to `bk-3302`
(September and October, chapters 7 and 8) and `bk-3410` to `bk-3422` (calibration week).

#### Season 2's files

The formats gctl reads and writes this season. A lesson uses the fields it needs, never
renames one, and may add a lesson-only type for anything else.

- **Controller log** (chapter 3). Salar Quieto's dish controller writes one pointing
  reading per line; the log Paz brought down covers the weeks before she left, up to 6 July
  2027: the UTC time, the dish (`LG` for La Grande, `CH` for Chiquita), then azimuth and
  elevation in degrees, separated by any mix of spaces and tabs. A reading near the horizon
  can be a few tenths below zero. A line whose first character other than a space or tab is
  `#` is a note from the night shift, and blank lines separate shifts. A calibration sweep is
  kept in a file of its own: one line holding the time, the dish and every azimuth and
  elevation pair of the sweep, up to about 200 KiB (Paz's of 2 July).

  ```
  # night shift 2027-07-06, Paz
  2027-07-06T03:14:07Z	LG	183.25	41.70
  2027-07-06T03:14:08Z  CH  183.31   41.66
  ```

  The controller's own export is the same readings as `;`-separated CSV, header
  `TIME;DISH;AZ;EL;NOTE`, with `#` starting a comment line (chapter 6).
- **Booking** (chapters 4, 5, 7, 8 and 10). One reserved pass, as Yusra's booking export
  writes it: JSON, times in RFC 3339 UTC, IDs `bk-` and four digits. A day's export is a JSON
  array; the recorder's booking log is one JSON object per line.

  ```json
  {"id":"bk-2041","sat_id":"furrow-a","station":"SQT","aos":"2027-08-11T14:02:10Z","los":"2027-08-11T14:11:52Z","note":"priority"}
  ```

  ```go
  // Booking is one reserved pass, as the booking export writes it.
  type Booking struct {
  	ID      string    `json:"id"`
  	SatID   string    `json:"sat_id"`
  	Station string    `json:"station"`
  	AOS     time.Time `json:"aos"`
  	LOS     time.Time `json:"los"`
  	Note    string    `json:"note,omitempty"`
  }
  ```

- **Delivery report** (chapter 5). Furrowlight's report of a delivery, one JSON object read
  strictly: `version`, `kind`, `id` and `payload`. The ID is Furrowlight's report number, above
  2^53 (`9007199254741003`), so it is decoded as a `json.Number` into an `int64` and never
  through `float64`. The payload stays a `json.RawMessage` until the kind is known: `image`
  (`sat_id`, `station`, `file` such as `fa-0824-1402.img`, `chunks`, `missing`) or `station`
  (`station`, `day`, `passes`, `missing`). Version 1 (to July) had no kind and was always an
  image report; version 2 (August) added kinds; version 3 appeared on Thursday 26 August
  while the van was offline and is refused, `report <id>: unsupported version 3`, until
  Furrowlight sends version 2 again.

  ```json
  {"version":2,"kind":"image","id":9007199254741003,"payload":{"sat_id":"furrow-a","station":"SQT","file":"fa-0824-1402.img","chunks":2000,"missing":9}}
  ```

- **Booking sheet** (chapters 6 and 10). Hollis exports the same bookings from Furrowlight's
  spreadsheet as CSV, header first, in capitals: `ID,SATELLITE,STATION,AOS,LOS,NOTE`.
  Columns may come in any order and are found by their header. A note may hold commas,
  quotes and line breaks, and Hollis writes notes in capitals too. Times are written as in
  the JSON.
- **Settings** (chapters 5 and 10). `settings.json`, one JSON object, read strictly: an
  unknown field, a second value or anything after the object is an error, never ignored.
  `station` defaults to `SQT`, `min_elev_deg` to 3 and `timeout` to `30s`; `day` is required
  and is a date (`2006-01-02`). A timeout is duration text such as `30s` or `2m`, never a
  count of nanoseconds: `type Delay time.Duration`, written back as `Duration.String()`.

  ```json
  {"station":"SQT","day":"2027-10-21","min_elev_deg":3,"timeout":"30s"}
  ```

- **Station table** (chapters 8, 9 and 10). gctl carries `stations.csv` inside its binary,
  made from the station table above. Bands are the ones a customer can book, separated by
  spaces; a receive-only UHF yagi is not among them.

  ```
  code,name,lat_deg,lon_deg,mask_deg,bands
  YST,Ystmark,77.8,2.0,5,S X
  PEL,Port Ellery,46.9,-38.1,5,UHF S
  SQT,Salar Quieto,-23.4,-104.6,3,S X
  LAN,Lantern Atoll,-1.8,82.5,8,UHF S
  CWR,Cape Wren,-53.9,-95.0,5,S X
  ```

- **Resources** (chapter 9). gctl also carries report templates (`templates/day.txt`,
  `pass.txt`, `notices/handover.txt`, `notices/stow.txt`, with `{field}` placeholders) and a
  default `calibration.json`, the van dish's own pointing error measured at Port Ellery
  (`{"az_offset_deg":0,"el_offset_deg":0}` built in); a `calibration.json` on the stick
  overrides it, and only a missing file falls back to the default.
- **Commands** (chapters 8 and 10). `run(args []string, stdin io.Reader, stdout, stderr
  io.Writer) int` receives the arguments after the program's name. The top-level usage:

  ```
  usage: gctl <command> [flags] [arguments]

  commands:
    passes [-station CODE] [file]  list the bookings by AOS
    stations                       print the station table
    stow [-reason TEXT] CODE       record a stow
  ```

  `passes` reads the recorder's booking log (one JSON booking per line, from the named file
  or standard input) and prints `<AOS> <station> <satellite> <id>` per booking, AOS as RFC
  3339 UTC, stable-sorted by AOS. `stations` prints `<CODE> <Name>` in the table's order.
  `stow` prints `<CODE> stowed: <reason>`, the reason defaulting to `no reason given`.
  Results go to standard output and diagnostics to standard error. `-h`, `-help` and `--help`
  print the usage to standard output and return 0. A usage error prints the problem (no
  prefix: `gctl needs a command`, `unknown command "fly"`) and the usage to standard error
  and returns 2; each command's own usage starts `usage: gctl <command>`. An operational
  failure prints exactly one line, `gctl: <error>` (such as `gctl: unknown station "SQX"`),
  to standard error, nothing on standard output, and returns 1. Chapter 10's copy of gctl
  carries only `day [-format json|csv] dir`, which runs a calibration day from a folder
  holding `settings.json` and `bookings.json` or `bookings.csv`.
- **Day report** (chapter 10). `<dir>/reports/daily/<day>.txt`, saved whole through a
  temporary sibling (as chapter 1 taught) before it is printed: a heading, a count, then the
  day's bookings at the settings' station, by AOS then ID, notes on one line.

  ```
  Calibration day 2027-10-21 at SQT (Salar Quieto)
  Booked passes (UTC): 4
  14:02:10-14:13:52 furrow-a bk-3410 PRIORITY, FIELD 7
  ```

## Seasons

Each course of the series is a season of the learner's career at Apogee, and story time only
moves forward. Each season's opening gives the learner a new role on their badge. The roadmap
is jobloom/gosmos#18; how each season is told is `docs/storytelling.md`.

| Season | Course | Story time | The season in a sentence | Issue |
|---|---|---|---|---|
| 1 | Go Foundations | Mar–Jul 2027 | From your first message to Pipit-1 to your first solo week on call | written; #7, #8 |
| 2 · Tools of the Trade | Files, Data and the Command Line | Jul–Oct 2027 | The test van's trip up the Quebrada Añil, with no signal: gctl, from files alone | #9 |
| Project | Build gctl | Oct 2027 | The night shift's interactive tool | #10 |
| 3 · Open for Business | Web Services in Go | Nov 2027–Feb 2028 | Customers book passes through Apogee's own service, which you build | #11 |
| Project | Build the Board | Feb 2028 | The Board on every screen in the building | #12 |
| 4 · Down to the Wire | Networking from the Socket Up | Mar–Jun 2028 | Ystmark's new dish controllers, and a web server from nothing but a socket | #13 |
| 5 · Keeping Secrets | Cryptography in Go | Jul–Oct 2028 | Cirrine's commands, readable by Cirrine alone, and Rotation Day | #14 |
| 6 · On Call | Logging, Metrics and Profiling | Nov 2028–Feb 2029 | The 03:12 page from Lantern Atoll, and the slow leak behind it | #15 |
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

## Style guide

### Voice

Second person, present tense: "You open the pass log." Concise, warm, precise. Explain like
a good colleague next to the learner: no hype, no jokes at the learner's expense, no "simply"
or "just". Prefer concrete numbers to adjectives. American spelling.

### A lesson is short and hands-on

Lessons should feel like doing, not reading. Each one teaches **one idea** and has the
learner use it right away.

- **Readme: 60 to 200 words, most about 100**, counting code. A concept that needs more than
  about 200 words becomes two lessons.
- **Shape**, always in this order:
  1. **Story beat:** one or two sentences that say why this matters today. Never a page.
  2. **The idea:** one to three sentences.
  3. **One tiny example:** at most eight lines, about something other than the task, and
     never the solution.
  4. **`## Your task`:** one to three sentences, or up to three numbered steps. Name the
     function or line to change.
  5. Optionally, one callout of one or two sentences (`> [!TIP]`, `> [!NOTE]` or
     `> [!WARNING]`) for a trap or a term.
- No H1 heading: the app shows the title. `## Your task` is usually the only heading.
- Leave out what this lesson doesn't need, even if it's true. When a learner might wonder,
  point ahead in a sentence ("Why the `&`? Next lesson.") instead of explaining now.
- No lists of rules, no history, no tangents, no "in this lesson you will learn".

Template:

````markdown
<Story beat: one or two sentences.>

<The idea: one to three sentences.>

```go
<one tiny example>
```

## Your task

<What to change or write, in one to three sentences.>
````

### Quizzes are quick checks

- A quiz follows the lesson or two whose idea it checks. It is never where an idea is taught.
- **1 to 3 questions**, each with 2 to 5 answers, exactly one right. Most have 3 answers.
- A question is one line (the format allows nothing else). Put any code the questions refer
  to in the quiz's readme, eight lines at most, and keep that readme to 60 to 120 words.
- Wrong answers are plausible mistakes, not jokes. No "all of the above", no trick wording,
  no answer that is right under some reading. Vary the position of the right answer.
- Title them "Quick Check: <topic>".

### Mix and rhythm

- About 83% of lessons are code (`code`, checked by output) or tests (`code_tests`), about
  17% quizzes.
- Use `code` when the output shows the idea at a glance; use `code_tests` when edge cases
  matter or the learner writes a function.
- Most chapters end with a slightly larger lesson that combines the chapter's ideas, still
  within the word limit.

### Code style

- gofmt, tabs, standard library only, Go 1.27. Idiomatic names (Effective Go): MixedCaps,
  initialisms in one case (`ID`, `AOS`, `LOS`, `CRC`, `HMAC`, `UTC`).
- Units in names: `AltKm`, `MaskDeg`, `FreqHz`, `BattV`, `BattMV`, `TempC`, `RateMbps`,
  `MarginDB`. Times are `time.Time` in UTC, durations `time.Duration`.
- Small functions: a solution's function is rarely longer than 15 lines.
- Every function in starter code except `main` has a one-line doc comment saying what it does.
- Errors: lowercase, no final punctuation, wrapped with `%w` when they carry context.
- No global variables except constants and the lesson's fixed data.
- Never print the wall clock, random values or map iteration order.

### Starter code and solutions

- The starter builds, runs and fails the check for the lesson's reason: wrong output or a
  failing test, never a panic. The exception is a lesson about something the compiler
  checks (an unused import, a missing method, an unexported name, an import cycle): then
  the starter fails to build with exactly the message the lesson explains, and the
  syllabus entry names it.
- The place to change is obvious: a function returning a zero value, a placeholder string,
  or at most one `// TODO:` line naming the action. The solution removes the TODO.
- The solution changes only what the task asks, usually 1 to 10 lines. It never adds a
  file: every file a solution touches is in the starter, and the solution includes every
  editable starter file, changed or not.
- Lessons in the same chapter may reuse the previous lesson's solution as the next starter.

### Tests

- Read-only `main_test.go` in `package main`, standard library only. Every `*_test.go` file
  is read-only, so a learner never edits a test file.
- Each `TestXxx` checks one function. More than three cases go in a table with names from
  the story, run with `t.Run`.
- Failure messages read `call = got, want want`: `stationCode() = "TODO", want "PEL"`. Use
  `%q` for strings, `%v` otherwise, and a tolerance for floats:
  `math.Abs(got-want) > 0.01` with `%.2f` in the message.
- `t.Errorf` when the other checks still make sense; `t.Fatalf` when they don't.
- Cover the normal case, a boundary, and the mistake learners actually make.
- Concurrency tests are correct by construction for a right solution: they never rely on
  timing luck or the race detector (the app doesn't use `-race`). Sleeps are 50 ms at most.
- In the testing chapter, where learners write checks themselves, their code goes in an
  editable non-test file (`checks.go`, `cases.go`) and the read-only test calls it. The
  test proves the checks work by running them against planted mutants in a re-executed test
  binary (`os.Args[0] -test.run=...` with an environment variable naming the mutant).

### Expected output

- Only stdout is compared, and the program must exit with status 0.
- `expected_output.txt` is exactly what the solution prints, its final newline included:
  `library verify` compares it byte for byte. The app, grading a learner, ignores line endings,
  trailing spaces and final newlines, and verify checks that the starter fails even then.
- ASCII unless the lesson is about text encoding; at most about 15 lines.

### What the app does that writers must know

- Code runs in a folder whose `go.mod` reads `module lesson` and the learner's Go version
  (`go 1.27`), unless the starter brings its own `go.mod`. A `code` lesson is built and
  run, and Submit compares its stdout; a `code_tests` lesson runs
  `go test -count=1 -timeout 10s` on Run and on Submit.
- `go test` runs part of `go vet`: starter code must pass it. Never write
  `fmt.Printf(msg)` with a variable format string.
- The app reads `go test` output for cases between lines of ten or more dashes and for lines
  like "3 passed, 1 failed". Tests must not print either.
- Every run takes well under 2 seconds on a laptop; the app stops runs at 10.

## Naming rules

- **Everything invented.** No real companies, people, products or places. Before a new name
  ships, search the web for its exact spelling. If it belongs to a real company, product,
  place, flown or announced space mission, or a person with a public profile, pick another.
  Generic words (birds, plants) may turn up in tiny or cancelled projects; record such
  overlaps below and keep the name only if nobody would mistake one for the other.
- **Name checks** (web search, 2026-11-01): nothing real was found for Furrowlight,
  Brinewatch, Cirrine Atmospherics, Kessick Polytechnic, Ystmark, Port Ellery, Salar
  Quieto, Lantern Atoll, Cape Wren, Thistle as a satellite, or the cast's full names.
  Known overlaps: a small satellite project called "Pipit" was cancelled before launch and
  never had a numbered series, so Pipit-1 to Pipit-3 stay. (A rename should keep five
  letters, so lesson 1.3's "29 bytes" stays true.) The Salar Quieto engineer was renamed
  from an earlier draft because that name belongs to a public figure.
- **Satellites:** IDs lowercase with hyphens (`pipit-1`, `furrow-a`); display names as in
  the fleet table (`Pipit-1`, `Furrow-A`).
- **Stations:** three capital letters (`YST`) and the name (`Ystmark`). Never abbreviate a
  name otherwise.
- **Characters:** first name in story beats; spellings exactly as above, accents included
  (Tomás).
- **Course, chapter and lesson slugs:** lowercase letters, digits and single hyphens, two
  to five words. Lesson slugs are unique across the whole course, not just the chapter.
- **Titles:** Title Case, two to five words, specific: "The Pre-Pass Checklist", not
  "Statements". Quizzes are "Quick Check: <topic>".
- **Files:** `main.go`, `main_test.go`; other files named for what they hold
  (`station.go`, `checks.go`); package folders are one lowercase word.
- **Topics** (front matter), from this list, alphabetical in each lesson: arrays, bytes, channels,
  comments, compiler-errors, conditionals, constants, context, conversions, defer, errors,
  files, formatting, functions, generics, goroutines, interfaces, loops, maps, methods, modules,
  multiple-returns, mutexes, packages, pointers, printing, programs, runes, select, slices,
  sorting, statements, streams, strings, structs, switch, testing, time, types, variables.
