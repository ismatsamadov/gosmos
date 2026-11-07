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

### The gctl project: the night shift's shell

The van comes home to Port Ellery on Saturday 23 October 2027, and by Monday every engineer
wants gctl on their own Mac. Ingrid asks for one more thing: an interactive mode, so the
night shift can type commands during a pass without starting the tool again each time. The
guided project **Build gctl** runs from Monday 25 to Friday 29 October 2027 at Port Ellery:
one chapter a day or so, ending on the night of Friday 29 October, when the night shift runs
a whole pass from your prompt. Nkem reviews; Ingrid is the customer; Cape Wren's wind supplies
the night's trouble. The badge is unchanged (Ground Tools Engineer, No. 039).

`gctl shell` is one `package main` program that grows lesson by lesson, each lesson's
starter being the previous lesson's solution:

- `shell.go`: `Shell`, its loop (`Run(in io.Reader, out io.Writer) error`) and dispatch;
  `commands.go`: the command table, a `map[string]command` from a name to its help text and
  its function; `cache.go`: the pass cache; `settings.go` and `history.go` in chapter 4;
  `data.go` (read-only): the night's passes and the station table; `main.go` wires them up.
- The prompt is `gctl> ` (no newline), written before each line is read. A line is split
  into words with `strings.Fields`: the command, then its arguments. A blank line does
  nothing. An unknown command prints `unknown command "<name>"; type help`. A command's error
  prints `error: <message>` and the session goes on; a panicking command prints
  `error: internal error in <name>` and the session goes on. `exit` or the end of the input
  prints `bye` and ends the session; `Run` returns only a failure of the input itself.
- Commands: `help` (every command and its help text, sorted by name), `exit`,
  `passes [CODE]` (the night's passes, all or one station's, one per line as
  `04:12-04:21 PEL pipit-1 max 61.4`, ordered by AOS, ending in ` (stowed)` at a stowed
  station), `stations` (`PEL Port Ellery`, and
  `stowed: <reason>` after a stowed one), `stow CODE [reason...]`, `next` (the next pass
  after now at a station that isn't stowed: `next: pipit-1 over PEL at 04:12, in 18m0s`, or
  `no more passes tonight`), `history` (every command typed this session and before,
  numbered from 1).
- The pass cache loads the night's passes through a loader function on first use, keeps
  them for a time to live (10 minutes unless the settings say otherwise), loads again when
  they are stale, and is safe for a background refresher that runs beside the commands. Its
  clock is injected (`now func() time.Time`) so tests move time without waiting.
- Settings: `gctl.json`, read strictly (Season 2's rules): `prompt`, `cache_ttl` as duration
  text, and `station`, the default for `passes`. A missing file means the defaults; a broken
  one is an error. History is appended to a file named by the program, one line per command,
  and read back at the next start.
- The night of Friday 29 October 2027 runs past midnight, so its passes are on Saturday
  30 October, UTC. Their lengths fit each orbit, mask and maximum elevation:

  | Satellite | Station | AOS–LOS | Max elevation |
  |---|---|---|---|
  | Pipit-1 | PEL | 04:12–04:21 | 61.4° (as on your first Board) |
  | Pipit-2 | YST | 04:31–04:40 | 47.2° |
  | Furrow-A | CWR | 04:58–05:06 | 33.5° |
  | Brine-3 | LAN | 05:20–05:26 | 22.8° |
  | Cirrine-1 | YST | 05:47–05:57 | 71.9° |

  At 04:40 Cape Wren's wind reaches 74 km/h and its dish stows itself (it stows from
  72 km/h); the night shift records it with `stow CWR wind 74 km/h`, and `next` must skip
  Furrow-A's pass. The relief shift takes over at 06:00, and Ingrid writes one line in the
  log: "gctl shell: GO for every shift." On Monday Yusra's message opens Season 3: "TL;DR:
  customers want to book their own passes."
- Earlier in the week, on Wednesday, a storm cuts Ystmark's fiber and Sunniva runs gctl over
  the slow satellite backhaul ("Storm. Fiber down. Nine seconds a command."), which is why the
  shell keeps its passes in a cache.
- Exact texts: `stow` confirms `CWR stowed: wind 74 km/h`; `history` prints `%3d  %s` lines;
  errors read `unknown station "XYZ"`, `passes takes at most one station code`,
  `stations takes no arguments`, `stow needs a station code`, `CWR is already stowed`,
  `next takes no arguments`, `history takes no arguments`, `internal error in <name>` and
  `history: <err>`; a broken settings file reads `gctl.json: more data after the settings`,
  `gctl.json: cache_ttl "10" is not a positive duration` or
  `gctl.json: cache_ttl 600 is not duration text`. The settings file is `gctl.json` and the
  history file `gctl_history`, both beside the program.

### Season 3: open for business

After the gctl project week, the learner starts the next assignment on Monday 1 November
2027. Their new badge says **Web Services Engineer**, still No. 039, based at Port Ellery
headquarters, since Monday 1 November 2027. Customers book passes by email: Yusra copies about
forty bookings a day into a spreadsheet by hand, and in August one typo booked Furrow-A over
the wrong station. On Tuesday 1 February 2028 Apogee opens a booking service to every customer
at once: their programs ask it for passes, book them, and hear back when their data has
landed, day and night. You build it with Go's standard library. Nkem reviews every change,
Yusra brings the customers' wishes, Mirela reviews who may do what, and Ingrid runs launch
day like a pass.

Story time runs from Monday 1 November 2027 to Tuesday 1 February 2028: eleven chapters of
about a week each, with a quiet holiday fortnight (Monday 20 December to Sunday 2 January)
between chapters 7 and 8. No new people, places, satellites or customers are introduced.
Customers speak through Yusra ("TL;DR: ...") and Furrowlight through Hollis's capitals.
Thistle re-entered in July 2027, so the service's fleet has seven satellites and Kessick
Polytechnic books nothing this season.

#### The season, chapter by chapter

Everything happens at Port Ellery headquarters. Customers try the service before launch ("the
beta"): Cirrine from chapter 3, Furrowlight from chapter 4, Brinewatch from chapter 5 (its
client goes wrong in chapter 9). Every change is tried first on a booking of a Pipit, as Apogee
tests everything.

| Ch | Week of | Where | Who needs what | What is at stake |
|---|---|---|---|---|
| 1 | Mon 1 Nov 2027 | the backend team's desks beside the Fishbowl | Nkem: before building anything, learn what a customer's program will send and get back, on a stand-in service she wrote | Yusra types forty emailed bookings a day; a wrong word in a program books the wrong pass |
| 2 | Mon 8 Nov | the Fishbowl, at the go/no-go poll | Ingrid: the go/no-go poll's wind check must answer within two seconds or say why not | On Tuesday 2 November the poll waited four minutes on the weather service while Cape Wren's mast was silent, and Furrow-B's checklist finished after AOS |
| 3 | Mon 15 Nov | the backend desks | Yusra: the booking service's first answer, `GET /stations`, for Cirrine's engineers, who start their client this week; her customer guide served beside it | The first thing a customer ever asks Apogee's own service. On Friday 19 November Yusra sends Cirrine the address |
| 4 | Mon 22 Nov | the backend desks | Hollis, through Yusra: a link to every pass for his emails, and "TL;DR: Furrowlight wants to book a week at a time" | A link that points at the wrong pass books the wrong pass |
| 5 | Mon 29 Nov | the backend desks | Hollis: Furrowlight's program books its first pass, on Friday 3 December, for Furrow-A over Salar Quieto on Monday 6 December, AOS 14:02:10 | A booking the service accepts must be a pass that exists. It becomes `bk-4501`, the first booking that never touches Yusra's spreadsheet |
| 6 | Mon 6 Dec | the backend desks | Nkem: on Monday a panic in a handler answered Cirrine's program with a dropped connection and left no trace | Every request gets an answer, an ID and a line in the log, as every failure at Apogee gets a name, a reason and a line in the log |
| 7 | Mon 13 Dec | the backend desks | Yusra: on Monday a deploy restarted the beta and Furrowlight's twelve test bookings vanished | Bookings must survive a restart before the holiday fortnight. On Friday 17 December the last deploy of the year goes out |
| 8 | Mon 3 Jan 2028 | Mirela's desk, then the backend desks | Mirela's review: anyone who can reach the service can read and cancel any booking | "Brinewatch must never see Cirrine's passes." She signs off on Friday 7 January, "LGTM once it fails closed" |
| 9 | Mon 10 Jan | the backend desks | Brinewatch's beta client lost its network on Sunday night and has retried ever since with no pause: a request every millisecond | One customer's bug must not slow everyone else; Cirrine's bookings wait half a minute behind the flood |
| 10 | Mon 17 Jan | the backend desks | Yusra: customers ask `GET /bookings` every few seconds to learn whether their data has landed; Furrowlight asks every five seconds, 17,280 times a day | "Your data has landed," a second after the file lands, signed so nobody else can say it |
| 11 | Mon 24 Jan, then Tue 1 Feb | the Fishbowl | Ingrid runs launch like a pass: a rehearsal week, a full rehearsal on Monday 31 January, then launch at 14:00 UTC on Tuesday 1 February | At 14:01 Yusra finds that the customer guide's sample request says `Authorization: Token` instead of `Bearer`, and customers are copying it. The guide is built into the service, so it must restart at 14:05, while Furrow-A is over Salar Quieto (AOS 14:02:10, LOS 14:11:52) |

The opening (Monday 1 November) hands over the new badge, shows Yusra's spreadsheet and the
launch date, introduces Nkem (reviewer), Yusra (the customers' voice) and Mirela (who will
review the service in January), and sends the learner to Nkem's stand-in.

The finale (chapter 11, Tuesday 1 February): the restart at 14:05 drops no request; the
bookings and the pending webhooks survive it; Cirrine's event from its 13:30 pass, still being retried
when the old process stopped, is delivered by the new one; Furrow-A's image lands at 14:13:02
and Furrowlight's receiver hears about it at 14:13:03. By evening every customer has booked
through the service and Yusra archives her spreadsheet. The debrief points at the Board project
(the next passes on every screen in the building) and at Season 4: Sunniva's new dish
controllers at Ystmark arrive in March, and they speak only raw sockets.

#### New systems

- **bookingd**, the booking service, which this season builds: add it to the Services list
  ("bookingd (books passes for customers and tells them when their data lands)"). In readmes
  it is "the booking service" first, `bookingd` once it has been glossed.
- **The weather service**, an internal service Tomás wrote in 2024. It reads each station's
  weather mast (the anemometer and a thermometer) every minute and answers
  `GET /weather/{code}`. It is old and quirky: its errors are plain text, a silent mast answers
  503, and when a station's backhaul is slow it can stall for minutes. Chapter 2 writes its
  client; the Board project reuses that client.
- **Nkem's stand-in** (chapter 1 only): a few lines of Go that answer like the booking service
  will, about stations only, so the learner can study HTTP before the service exists.
- **The customer guide**: Yusra's pages for customers' engineers, served by the service at
  `/guide/` and built into its binary.
- **The catalog**: passplan's predicted passes for the next 14 days, which the service offers
  for booking. A customer can book only a pass in the catalog.
- **Landings**: dlinkd tells the booking service when a booking's file has been delivered
  (a `Landing`), and the service tells the customer with a webhook.

#### The code, by chapter

Each lesson is a program of its own, so the read-only file a name sits in is the lesson's
choice (`data.go` and `types.go` are fine homes for fixtures and shared types); the names,
types, signatures, routes and messages below are fixed.

Each lesson is its own `package main` (module `lesson`, as the app writes it). A lesson gives
the pieces of earlier chapters that it needs as read-only files, under these names, and uses
only the fields it needs. In chapters 4 to 6 the store (`store.go`) is a hidden file, so the
learner first meets `MemStore`'s inside in chapter 7, where they write it.

| Chapter | Names it adds (and files) |
|---|---|
| 1 | `standIn() http.Handler` (Nkem's, read-only `standin.go`); `printRequest`, `requestLine`, `printResponse` (read-only `wire.go`); `passesPath`, `advice`, `isJSON`, `fetch` |
| 2 | `Weather`, `WeatherClient`, `newWeatherClient`, `(*WeatherClient).Reading`, `decodeWeather`, `StatusError`, `windCall` (`weather.go`) |
| 3 | `Service`, `(*Service).routes`, `listStations`, `health`, `writeJSON`, the guide route; `start` (`service.go`, `stations.go`) |
| 4 | `getStation`, `listPasses`, `getPass`, `getBooking`, `cancelBooking`; `PassView`; `MemStore` (hidden `store.go`) |
| 5 | `createBooking`, `APIError`, `ErrorBody`, `writeError`, `decodeJSON`, `bodyError`, `maxBodyBytes`, `checkShape`, `checkPass`, `invalidBooking` (in `errors.go`), `knownSat`, `knownStation`, `window`, `findStation` (`errors.go`, `decode.go`, `rules.go`, `data.go`) |
| 6 | `noStore`, `requestID`, `requestIDFrom`, `statusRecorder`, `logRequests`, `recoverPanics`, `chain`, `(*Service).handler` (`middleware.go`) |
| 7 | `Store`, `ErrNotFound`, `ErrConflict`, `MemStore`, `NewMemStore`, `FileStore`, `OpenFileStore`, `firstNumber`, `(*MemStore).load`, `replaceFile`, `listBookings`, `Config`, `newService` (`store.go`, `filestore.go`) |
| 8 | `Customer`, `credentials`, `checkKey`, `signToken`, `verifyToken`, `issueToken`, `TokenResponse`, `authenticate`, `withCustomer`, `customerFrom`, `owns`, `(*Service).private`, `unauthorized`, `tokenMAC`, `findBooking` (`auth.go`, `bookings.go`) |
| 9 | `Bucket`, `Limiter`, `newLimiter`, `limit`, `retryDelay`, `Timeouts`, `newServer`, `withDeadline`, `bucketSize`, `refillPerSec`, `errRateLimited`, `(*Bucket).Take` (`limit.go`, `server.go`) |
| 10 | `Event`, `Landing`, `post`, `AnswerError`, `statusOf`, `signature`, `verifyEvent`, `backoff`, `retryable`, `Sender`, `Delivery`, `ErrUndelivered`, `Notifier`, `idempotent`, `keyStore`, `claim`, `validKey`, `replay`, `answerRecorder` (`webhook.go`, `idempotency.go`) |
| 11 | `stop`, `serve`, `run`, `Outbox`, `OpenOutbox`, `(*Outbox).NewEventID`, `Put`, `Remove`, `Pending`, `outboxFile`, `(*Notifier).Start`, `Stop`, `(*Notifier).deliver`, `sleep`, `App`, `newApp`; 11.8's `client.go` (`request`, `answer`, `errorCode`, with chapter 8's `send` and `token`) (`run.go`, `outbox.go`, `webhook.go`, `app.go`) |

Signatures every chapter after the first one that defines them uses exactly:

```go
// chapters 1 and 2
func passesPath(sat, station, from string) string
func advice(code int) string
func isJSON(h http.Header) bool
func fetch(url string) (status, body string, err error)
func newWeatherClient(baseURL string, timeout time.Duration) *WeatherClient
func (c *WeatherClient) Reading(ctx context.Context, code string) (Weather, error)
func decodeWeather(r io.Reader) (Weather, error)
func windCall(ctx context.Context, c *WeatherClient, code string, now time.Time) string

// the booking service
func start(s *Service) (addr string, stop func() error, err error) // chapter 3 only
func writeJSON(w http.ResponseWriter, status int, v any)
func writeError(w http.ResponseWriter, r *http.Request, err error)
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error
func checkShape(b Booking) error
func checkPass(b Booking, catalog []PassView, now time.Time) error
func (s *Service) routes() *http.ServeMux
func (s *Service) handler(l *log.Logger) http.Handler // the mux inside every middleware
func (s *Service) private(h http.HandlerFunc) http.Handler // authenticate, then limit
func chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler
func noStore(next http.Handler) http.Handler
func requestID(next http.Handler) http.Handler
func requestIDFrom(ctx context.Context) string
func logRequests(l *log.Logger, now func() time.Time) func(http.Handler) http.Handler
func recoverPanics(l *log.Logger) func(http.Handler) http.Handler
func credentials(r *http.Request, scheme string) (string, error)
func checkKey(customers map[string]Customer, key string) (Customer, bool)
func signToken(customer string, expires time.Time, key []byte) string
func verifyToken(tok string, key []byte, now time.Time) (string, error)
func authenticate(customers map[string]Customer, key []byte, now func() time.Time) func(http.Handler) http.Handler
func withCustomer(ctx context.Context, c Customer) context.Context
func customerFrom(ctx context.Context) (Customer, bool)
func owns(c Customer, satID string) bool
func limit(l *Limiter) func(http.Handler) http.Handler
func withDeadline(d time.Duration) func(http.Handler) http.Handler
func newServer(h http.Handler, t Timeouts) *http.Server
func retryDelay(res *http.Response) time.Duration
func post(ctx context.Context, c *http.Client, url string, secret []byte, ts int64, ev Event) error // a non-2xx answer is an *AnswerError
func statusOf(err error) int // the status an *AnswerError in err carries, or 0 when no answer came
func signature(secret []byte, ts int64, body []byte) string
func verifyEvent(r *http.Request, secret []byte, now time.Time) (Event, error)
func backoff(attempt int, r *rand.Rand) time.Duration
func retryable(status int, err error) bool
func (s *Sender) Deliver(ctx context.Context, d *Delivery) error
func (n *Notifier) Landed(ctx context.Context, l Landing) error
func (n *Notifier) Start()      // chapter 11: resumes the outbox's deliveries, each when due
func (n *Notifier) Stop() error // chapter 11: ends deliveries in progress, keeping their attempts
func newKeyStore(now func() time.Time) *keyStore
func idempotent(keys *keyStore) func(http.Handler) http.Handler
func newService(cfg Config) (*Service, error)
func stop(srv *http.Server, within time.Duration) error
func serve(srv *http.Server, ln net.Listener) error
func OpenOutbox(path string) (*Outbox, error)
func (o *Outbox) NewEventID() (string, error) // "evt-000118", once the count is saved
func (o *Outbox) Put(d Delivery) error         // adds or replaces d's delivery, and saves
func (o *Outbox) Remove(id string) error       // drops event id's delivery, and saves
func (o *Outbox) Pending() []Delivery          // a copy of the deliveries still to make
func newApp(cfg Config) (*App, error)
func run(ctx context.Context, ln net.Listener, app *App, t Timeouts) error

// checks the learner writes, in editable non-test files
func checkStore(t *testing.T, open func() Store)                             // chapter 7, contract.go
func review(t *testing.T, base string, advance func(time.Duration))          // chapter 8, review.go
func checkService(t *testing.T, base string)                                 // chapter 11, check.go
```

`App` (chapter 11) is the assembled service: `Service *Service`, `Notifier *Notifier` and
`Handler http.Handler`, the mux inside every middleware. `Notifier` holds the `Store`, the
customers, the `Sender` and (chapter 11) the `Outbox`, and counts event IDs.

#### The service's types

The canonical types keep their fields and types. On the wire, Station and Pass carry these
tags; Booking is Season 2's, unchanged.

```go
type Station struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	LatDeg     float64 `json:"lat_deg"`
	LonDeg     float64 `json:"lon_deg"`
	MaskDeg    float64 `json:"mask_deg"`
	Bands      []Band  `json:"bands"` // the bands a customer can book, as in Season 2's stations.csv
	Stowed     bool    `json:"stowed"`
	StowReason string  `json:"stow_reason,omitempty"`
}

type Pass struct {
	SatID      string    `json:"sat_id"`
	Station    string    `json:"station"`
	AOS        time.Time `json:"aos"`
	LOS        time.Time `json:"los"`
	MaxElevDeg float64   `json:"max_elev_deg"`
}

// PassView is a predicted pass as the service shows it. Pass is embedded, so its
// fields sit beside "id" in the JSON.
type PassView struct {
	ID string `json:"id"` // the pass ID: "furrow-a/SQT/15864"
	Pass
}

// Booking is one reserved pass (Season 2, unchanged).
type Booking struct {
	ID      string    `json:"id"`      // "bk-4501", given by the service
	SatID   string    `json:"sat_id"`
	Station string    `json:"station"`
	AOS     time.Time `json:"aos"`
	LOS     time.Time `json:"los"`
	Note    string    `json:"note,omitempty"`
}

// ErrorBody is the JSON body of every error the service writes.
type ErrorBody struct {
	Code      string `json:"error"`                // "not_found": see the table of errors
	Message   string `json:"message"`              // for people: lowercase, no final period
	RequestID string `json:"request_id,omitempty"` // from chapter 6
}

// APIError is an error that says exactly what the caller gets.
type APIError struct {
	Status  int    // 422
	Code    string // "invalid_booking"
	Message string // "los must be after aos"
}

func (e *APIError) Error() string { return e.Message }

// Customer is one company that books passes through the service (chapter 8 on).
type Customer struct {
	ID         string   // "furrowlight": lowercase, one word
	Name       string   // "Furrowlight"
	SatIDs     []string // the satellites it operates
	KeyHash    string   // lowercase hex SHA-256 of its whole API key
	WebhookURL string   // where its events go (chapter 10)
	WebhookKey []byte   // the secret its events are signed with (chapter 10)
}

// TokenResponse is the body of POST /tokens.
type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Service is the booking service: its data, its clock and its handlers.
type Service struct {
	Stations  []Station
	Catalog   []PassView          // ordered by AOS, then ID
	Guide     fs.FS               // the customer guide's pages
	Store     Store               // a *MemStore until chapter 7 makes it this interface
	Customers map[string]Customer // by ID (chapter 8)
	TokenKey  []byte              // signs tokens (chapter 8)
	Limiter   *Limiter            // chapter 9
	Timeouts  Timeouts            // chapter 9
	Keys      *keyStore           // idempotency keys (chapter 10)
	Now       func() time.Time
	draining  atomic.Bool // chapter 11
}

// Config is everything the service is built from. Tests fill it from the fixtures.
type Config struct {
	Dir       string           // data folder: bookings.json, outbox.json (t.TempDir() in tests)
	Stations  []Station
	Catalog   []PassView
	Guide     fs.FS            // embedded in production, an fstest.MapFS in tests
	Customers []Customer       // chapter 8
	TokenKey  []byte           // chapter 8
	Now       func() time.Time
	Rand      *rand.Rand       // chapter 10: jitter, seeded in tests
	Log       io.Writer        // chapter 6: the access log
	Timeouts  Timeouts         // chapters 9 and 11
}
```

#### Fixtures every chapter shares

- **Stations**, in the bible's order, with Season 2's bookable bands: YST `S X`, PEL `UHF S`,
  SQT `S X`, LAN `UHF S`, CWR `S X`. A receive-only UHF yagi is never bookable.
  `{"code":"SQT","name":"Salar Quieto","lat_deg":-23.4,"lon_deg":-104.6,"mask_deg":3,"bands":["S","X"],"stowed":false}`
- **The fleet the service knows**: `pipit-1`, `pipit-2`, `pipit-3`, `furrow-a`, `furrow-b`,
  `brine-3`, `cirrine-1`. `thistle` is unknown to it (re-entered in July).
- **Customers and their test keys.** An API key is the customer ID, a dot, then 32 lowercase
  hex digits (16 random bytes in production). Fixtures use exactly these:

  | ID | Name | Satellites | Test API key | Test webhook secret |
  |---|---|---|---|---|
  | `apogee` | Apogee | pipit-1, pipit-2, pipit-3 | `apogee.9e4d2a7c1b0f3e5d8a6c4b2f0e1d3a5c` | `apogee webhook secret` |
  | `furrowlight` | Furrowlight | furrow-a, furrow-b | `furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b` | `furrowlight webhook secret` |
  | `brinewatch` | Brinewatch | brine-3 | `brinewatch.0f1e2d3c4b5a69788796a5b4c3d2e1f0` | `brinewatch webhook secret` |
  | `cirrine` | Cirrine Atmospherics | cirrine-1 | `cirrine.c3a1f7e95d2b4086a7e1c9f3b5d8a024` | `cirrine webhook secret` |

  The test token key is the 32 ASCII bytes `bookingd test token key 32 bytes`. Production keys
  and secrets are 32 random bytes and never appear in a lesson.
- **The catalog.** Pass IDs as the bible gives them (satellite, station, orbit number padded
  to five digits). Orbit numbers grow by about 15 a day (Furrow-A 15.21, Pipit-1 15.16,
  Brine-3 15.3, Cirrine-1 14.92 a day), continuing from Season 1's Board; a pass between
  two the catalog lists is predicted from a model of its satellite's orbit, never guessed. Furrow-A crosses Salar Quieto at
  14:02 UTC on the days this season needs. Passes every chapter may use:

  | Pass ID | AOS | LOS | Max elev |
  |---|---|---|---|
  | `pipit-1/PEL/08154` | 2027-12-06T04:12:05Z | 2027-12-06T04:21:27Z | 61.4 |
  | `cirrine-1/YST/11364` | 2027-12-06T07:30:00Z | 2027-12-06T07:40:04Z | 44.2 |
  | `brine-3/LAN/09213` | 2027-12-06T09:14:20Z | 2027-12-06T09:22:02Z | 38.5 |
  | `furrow-a/SQT/15864` | 2027-12-06T14:02:10Z | 2027-12-06T14:11:52Z | 47.9 |
  | `furrow-b/SQT/15862` | 2027-12-06T14:49:30Z | 2027-12-06T14:58:58Z | 44.0 |
  | `furrow-a/SQT/16731` | 2028-02-01T14:02:10Z | 2028-02-01T14:11:52Z | 47.9 |

  Lessons add passes of their own by the same rules (passes last 5 to 15 minutes; Brine-3
  never rises over Ystmark).
- **Booking IDs** are `bk-` and four digits. Yusra's spreadsheet ends at `bk-4498`; in
  chapters 4 to 6 the service holds her last imported bookings, `bk-4491` to `bk-4498`. The
  service's own count starts at `bk-4501` (Furrowlight's first, chapter 5) and is near
  `bk-4912` on launch day (Furrow-A's 14:02 pass).

#### The API

Every answer the service writes is JSON (`Content-Type: application/json`) except the guide's
pages and the router's own 404 and 405. Times are RFC 3339 in UTC. A list is `[]` when empty,
never `null`. The router's answers are plain text: `404 page not found` for a path no route
matches, and `Method Not Allowed` with an `Allow` header (`Allow: DELETE, GET, HEAD`) for a
route's path with another method. Never register a catch-all `/`: it would answer every wrong
method with its own 404 instead of 405.

| Route | Answer | Guard | Chapter |
|---|---|---|---|
| `GET /health` | 200 `{"status":"ok"}`; 503 `{"status":"draining"}` while shutting down | public | 3; 11 |
| `GET /stations` | 200 every Station in table order; `?band=UHF\|S\|X` keeps the stations that offer it | public | 3 |
| `GET /stations/{code}` | 200 the Station; 404 `no station "ZZZ"`. Codes match exactly: `sqt` is unknown | public | 4 |
| `GET /guide/` | the customer guide's pages (`http.FileServerFS` under `http.StripPrefix("/guide/", …)`); `/guide` redirects to `/guide/` with 307 | public | 3 |
| `POST /tokens` | 200 TokenResponse | API key | 8 |
| `GET /passes` | 200 PassView list, ordered by AOS then ID, with the query below | token, limit | 4; 8 |
| `GET /passes/{sat}/{station}/{orbit}` | 200 the PassView whose ID is `{sat}/{station}/{orbit}`; 404 `no pass furrow-a/SQT/15865` | token, limit | 4; 8 |
| `POST /bookings` | 201 the stored Booking, with `Location: /bookings/bk-4501` | token, limit, idempotency key | 5; 8; 10 |
| `GET /bookings` | 200 the caller's bookings, ordered by AOS then ID | token, limit | 7; 8 |
| `GET /bookings/{id}` | 200 the Booking; 404 `booking bk-5499: not found` | token, limit | 4; 8 |
| `DELETE /bookings/{id}` | 204 with no body; 404 as for GET | token, limit | 4; 8 |

`GET /passes` query, all optional, other parameters ignored:

- `sat`: a satellite of the fleet, else 400 `unknown satellite "thistle"`.
- `station`: a station code, else 400 `unknown station "ZZZ"`.
- `from`: an RFC 3339 time, default the service's `Now()`, else 400
  `from must be an RFC 3339 time`.
- `days`: a whole number from 1 to 14, default 7, else 400
  `days must be a whole number from 1 to 14`.
- It lists passes whose AOS is at or after `from` and before `from` plus `days` × 24 hours.

In chapters 3 and 4, before errors are JSON, handlers answer their errors with `http.Error`
and the same message text. From lesson 5.2 every handler error goes through `writeError` as an
`ErrorBody` (5.1 still answers with `http.Error`):

`{"error":"not_found","message":"booking bk-5499: not found","request_id":"req-000042"}`

| Status | `error` | Messages, exactly |
|---|---|---|
| 400 | `bad_request` | `request body: empty`; `request body: ` followed by the decoder's error, such as `request body: json: unknown field "sat"`; `request body: unexpected data after the JSON value`; the query messages above; `unknown band "x"`; `idempotency key must be 1 to 64 letters, digits, - or _` |
| 401 | `unauthorized` | `missing credentials`; `invalid api key`; `invalid token`; `token expired` |
| 403 | `forbidden` | `cirrine-1 is not one of your satellites` |
| 404 | `not_found` | `no station "ZZZ"`; `no pass furrow-a/SQT/15865`; `booking bk-5499: not found` |
| 409 | `conflict` | `furrow-a over SQT at 2027-12-06T14:02:10Z: already booked`; `idempotency key "wk49-mon" was used for a different request`; `a request with idempotency key "wk49-mon" is in progress` |
| 413 | `too_large` | `request body over 65536 bytes` |
| 415 | `unsupported_media_type` | `send the body as application/json` |
| 422 | `invalid_booking` | the booking rules below |
| 429 | `rate_limited` | `slow down: 10 requests at once, then 1 a second` |
| 500 | `internal` | `internal error`, whatever went wrong; the error's own text never reaches the caller |
| 503 | `unavailable` | `request took too long`; `request canceled` |

`writeError` decides the status from the error value: an `*APIError`, wrapped or not
(`errors.AsType`), gives its own status, code and message (chapter 5); `ErrNotFound` gives 404
and `ErrConflict` 409 with the error's message (chapter 7); `context.DeadlineExceeded` and
`context.Canceled` give 503 (chapter 9); anything else gives 500 `internal`. From chapter 6 it
fills `request_id` from the request's context.

`writeJSON` marshals first; if that fails it answers 500 with the plain text
`internal error`. Otherwise it sets `Content-Type: application/json`, writes the status, then
the JSON and a newline.

**Creating a booking.** The body is a Booking without `id`. The service checks, in this order,
and the first failure answers:

1. `authenticate` (chapter 8): 401.
2. `limit` (chapter 9): 429.
3. `idempotent` (chapter 10): 400 for a malformed key, 413 for a body over the limit, a
   replay, or 409.
4. `decodeJSON`: the media type of `Content-Type` must be `application/json` (any parameters,
   any letter case), else 415; at most 65,536 bytes (`http.MaxBytesReader`; a
   `*http.MaxBytesError` from either decode call is 413); one JSON object with no unknown field
   (`DisallowUnknownFields`), then nothing but white space (a second `Decode` must return
   `io.EOF`), else 400.
5. `checkShape`, 422 `invalid_booking`, in this order: `id is chosen by the service`;
   `sat_id is required`; `station is required`; `aos and los are required`;
   `los must be after aos`; `note is longer than 200 characters` (characters, not bytes).
6. Ownership (chapter 8): the caller must operate the satellite, else 403. It comes before
   the catalog check, so a customer learns nothing about another customer's passes, not even
   whether one exists (a 422 for a missing pass beside a 403 for a real one would tell).
   Nobody operates `thistle`, so from chapter 8 on booking it answers 403 `thistle is not
   one of your satellites`; checkPass's `unknown satellite` is reachable only in chapters 5
   to 7 (`GET /passes?sat=thistle` still answers 400).
7. `checkPass`, 422 `invalid_booking`, in this order: `unknown satellite "thistle"`;
   `unknown station "ZZZ"`; `no predicted pass of furrow-a over SQT at 2027-12-06T14:02:10Z`
   (the satellite, station, AOS and LOS must equal a catalog pass's, times compared with
   `Equal`); `bookings for this pass closed at 2027-12-06T13:52:10Z`. Bookings close ten
   minutes before AOS, when Ingrid's checklist must be done ("AOS 04:12. Checklist done by
   04:02."): at exactly AOS minus 10 minutes the pass is closed.
8. `Store.Create`: a booking of the same pass (satellite, station and AOS) already kept is
   409 `conflict`.
9. 201 with `Location: /bookings/<id>` and the stored Booking.

A cancellation (`DELETE`) has no time rule: it answers 204 whenever the booking exists and is
the caller's.

#### Keys and tokens (chapter 8)

- **API keys** are given to customers by Yusra and sent only to `POST /tokens`, in
  `Authorization: Key <api key>`. The service keeps only `KeyHash`, the lowercase hex SHA-256
  of the whole key, and checks a key by hashing it and comparing the two hex strings with
  `subtle.ConstantTimeCompare`. A key that doesn't parse (no dot, a secret that isn't 32
  lowercase hex digits), an unknown customer and a wrong secret all answer the same 401
  `invalid api key`.
- **Tokens** are what every other customer route takes, in `Authorization: Bearer <token>`.
  A token is `<customer ID>.<expiry in Unix seconds>.<signature>`; the signature is the 64
  lowercase hex digits of HMAC-SHA256, under the service's 32-byte token key, of
  `<customer ID>.<expiry>`. A token lives one hour: `expires_at` is `Now()` plus one hour,
  truncated to the second. Furrowlight's token expiring at 2028-01-05T11:00:00Z is
  `furrowlight.1830682800.` followed by its 64 hex digits.
- `verifyToken` checks, in order: three parts separated by dots; a signature of 64 hex digits
  that `hmac.Equal` accepts; an expiry in the future. A token is valid while `now` is before
  its expiry; at the expiry second it is `token expired`. Any other failure, including a
  customer the service doesn't know, is `invalid token`.
- `credentials` reads `Authorization`: the scheme (matched with `strings.EqualFold`), one
  space, then a non-empty credential. A missing header, another scheme or an empty credential
  answers 401 `missing credentials`.
- Every 401 carries `WWW-Authenticate: Key` on `POST /tokens` and `WWW-Authenticate: Bearer`
  everywhere else. `authenticate` fails closed: it never calls the next handler after a 401.
- `authenticate` puts the Customer in the request's context (`withCustomer`, read back with
  `customerFrom`) under an unexported key type.
- **Authorization.** A customer sees and books only its own satellites (`owns`): `GET /passes`
  and `GET /bookings` list only its satellites' passes and bookings (`?sat=` of another
  customer's satellite gives `[]`); another customer's pass or booking answers 404 exactly as a
  missing one, and its cancellation changes nothing; booking another customer's satellite
  answers 403. Apogee's own key sees only the Pipits.
- `GET /health`, `GET /stations`, `GET /stations/{code}` and the guide stay public. `POST
  /tokens` is not rate limited.

#### Rate limits and timeouts (chapter 9)

- A **token bucket** per customer, applied after `authenticate` on every token route: it holds
  at most 10 tokens, refills continuously at 1 token per second, and starts full (a zero
  `Bucket` refills to full on first use). A request takes one token. Time comes from the
  injected `now`; when the clock steps backward nothing is added. Readmes say "token bucket"
  and explain once that its tokens are not chapter 8's.
- Every answer that passed the limiter carries `X-RateLimit-Limit: 10` and
  `X-RateLimit-Remaining: <whole tokens left after this request>`. A refused request answers
  429 `rate_limited` with `Retry-After: <seconds until one token is back, rounded up, at least
  1>` and `X-RateLimit-Remaining: 0`. Requests refused by `authenticate` are never counted. A
  request that reaches `limit` without a customer in its context is a wiring mistake and
  answers 500.
- Clients wait as `Retry-After` says: whole seconds, 1 second when the header is missing or
  unreadable, at most 60 seconds (`retryDelay`).

```go
// Bucket is one customer's token bucket.
type Bucket struct {
	Tokens float64   // tokens left, at most 10
	Last   time.Time // when Tokens was last brought up to date
}

// Limiter keeps a bucket per customer; it is safe for many requests at once.
type Limiter struct {
	now     func() time.Time
	mu      sync.Mutex
	buckets map[string]*Bucket
}

func newLimiter(now func() time.Time) *Limiter
func (l *Limiter) Allow(customer string) (ok bool, remaining int, retryAfter time.Duration)

// Timeouts are the service's limits, production values in comments. A zero leaves that limit unset.
type Timeouts struct {
	ReadHeader time.Duration // 5s: a client has this long to send its headers
	Read       time.Duration // 15s: ... to send the whole request
	Write      time.Duration // 30s: ... to take the whole answer
	Idle       time.Duration // 2m: a kept-alive connection may sit unused this long
	Handler    time.Duration // 10s: the deadline of every request's context (withDeadline)
	Drain      time.Duration // 5s: /health says draining this long before shutdown begins
	Shutdown   time.Duration // 20s: how long shutdown waits for requests in flight
}
```

A test sets the limits it checks, in milliseconds (50 ms; up to a second for Drain and
Shutdown), and leaves the others zero, so a request it holds open on purpose is never cut by
a limit it isn't testing.

`newServer` sets `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout` and `IdleTimeout` from
the Timeouts and `MaxHeaderBytes` to 16 KiB (16384; Go allows 4 KiB more before it answers
431, so tests send a 32 KiB header). A zero `ReadTimeout` or `WriteTimeout` means no limit; a zero `ReadHeaderTimeout` or
`IdleTimeout` falls back to `ReadTimeout`, as Go's server does. `withDeadline(0)`
adds no deadline, a zero `Drain` skips draining, and a zero `Shutdown` waits for every
request in flight.

#### Middleware, request IDs and the log (chapter 6)

- Middleware is `func(http.Handler) http.Handler`. `chain(h, a, b, c)` returns
  `a(b(c(h)))`: the first one listed runs first.
- Around the mux, outermost first: `requestID`, `logRequests`, `recoverPanics`, `noStore`,
  then `withDeadline` (chapter 9). Each private route is `chain(h, authenticate, limit)`, and
  `POST /bookings` adds `idempotent` inside them (chapter 10).
- `noStore` sets `Cache-Control: no-store` on every answer the service writes; Go's file server
  removes it from its own error answers (a missing guide page), so no test expects it there.
- `requestID` numbers requests from `req-000001` (six digits, counting up per process, from an
  `atomic.Uint64`), sets the `X-Request-Id` answer header before calling the next handler, and
  puts the ID in the request's context under an unexported key type. `requestIDFrom` returns
  `""` for a context without one.
- `statusRecorder` embeds `http.ResponseWriter`, records the first status written (200 when
  `Write` comes first, and 200 when the handler writes nothing, as the server sends) and the
  bytes written, and has `Unwrap() http.ResponseWriter` so
  `http.NewResponseController` still finds `Flush`.
- `logRequests` writes one line per request with a `*log.Logger` made with no prefix and no
  flags (safe for many goroutines), after the handler returns:
  `req-000007 POST /bookings 201 236B 4ms`. That is the request ID, the method, the URL path
  (never the query), the status, the bytes, and the time from the injected clock rounded to the
  millisecond in `time.Duration`'s own format (`0s`, `4ms`, `1.204s`).
- `recoverPanics` recovers a handler's panic, logs `req-000003 panic: <value>`, and answers 500
  `internal` through `writeError`. A panic with `http.ErrAbortHandler` is panicked again, so
  the server cuts the connection as that value asks.

#### The store and its file (chapter 7)

```go
var (
	ErrNotFound = errors.New("not found")      // wrapped: "booking bk-5499: not found"
	ErrConflict = errors.New("already booked") // wrapped: "furrow-a over SQT at 2027-12-06T14:02:10Z: already booked"
)

// Store keeps bookings. Its methods are safe to call from many goroutines at once.
type Store interface {
	// Create gives b the next ID and keeps it. It returns ErrConflict when a booking
	// of the same pass (satellite, station and AOS) is already kept.
	Create(ctx context.Context, b Booking) (Booking, error)
	// Get returns the booking with this ID, or ErrNotFound.
	Get(ctx context.Context, id string) (Booking, error)
	// List returns every booking, ordered by AOS and then ID, in a new slice.
	List(ctx context.Context) ([]Booking, error)
	// Delete removes the booking with this ID, or returns ErrNotFound.
	Delete(ctx context.Context, id string) error
}
```

- Every method returns `ctx.Err()` first when the context is already done.
- `MemStore` (`NewMemStore()`) is a map behind a `sync.Mutex`, with the next ID number; the
  conflict check and the insert happen under one lock. IDs count up from `bk-4501`. An
  unexported `load([]Booking) error` restores saved bookings with their IDs (duplicate IDs are
  an error) and sets the next number to the larger of the current one and one more than the
  highest held: never backward, never below 4501.
- `FileStore` (`OpenFileStore(path)`) keeps a `*MemStore` and the path of `bookings.json` in
  the service's data folder. After every successful Create or Delete it saves `List`'s
  bookings as one JSON array through `replaceFile(path, data)`: Season 2's write to a private
  temporary sibling, close, then rename, mode 0600. When saving fails it undoes the change and
  returns the error, so the store and the file never disagree. Opening: a missing file is an
  empty store; a file that doesn't decode, strictly, is an error, never an empty store.

#### Webhooks and idempotency keys (chapter 10)

```go
// Landing is dlinkd's word that a booking's file has been delivered.
type Landing struct {
	BookingID string    // "bk-4912"
	File      string    // "furrow-a-20280201-1402.tif": satellite, AOS day and time
	Bytes     int64
	At        time.Time // when the file was complete, UTC
}

// Event is what the service tells a customer when something happens to a booking.
type Event struct {
	ID        string    `json:"id"`         // "evt-000118": the same on every attempt
	Type      string    `json:"type"`       // "data.landed", the only type so far
	BookingID string    `json:"booking_id"`
	SatID     string    `json:"sat_id"`
	Station   string    `json:"station"`
	File      string    `json:"file"`
	Bytes     int64     `json:"bytes"`
	LandedAt  time.Time `json:"landed_at"`
}

// Delivery is one event on its way to one customer.
type Delivery struct {
	Customer string    `json:"customer"` // "furrowlight"
	Event    Event     `json:"event"`
	Attempts int       `json:"attempts"` // made so far
	NextAt   time.Time `json:"next_at"`  // when the next attempt is due
}

var ErrUndelivered = errors.New("undelivered")

// AnswerError is a receiver's answer outside 200 to 299: the event did not arrive.
type AnswerError struct {
	StatusCode int    // 503
	Status     string // "503 Service Unavailable"
}

// Sender delivers events to customers' webhooks.
type Sender struct {
	Client    *http.Client                                     // Timeout 5s; never follows a redirect
	Customers map[string]Customer                              // by ID: each one's WebhookURL and WebhookKey
	Now       func() time.Time                                 // each attempt's timestamp
	Rand      *rand.Rand                                       // jitter; used only under mu
	Wait      func(ctx context.Context, d time.Duration) error // sleeps in production; tests record d
	mu        sync.Mutex
}
```

- **The request.** `POST <the customer's WebhookURL>` with the event as `json.Marshal` writes
  it (no newline) and these headers: `Content-Type: application/json`;
  `Apogee-Timestamp: <Unix seconds of this attempt>`, such as `1833027183`;
  `Apogee-Signature: <64 lowercase hex digits>`, the HMAC-SHA256 under the customer's webhook
  secret of the timestamp, a dot, then the exact body bytes.
  `{"id":"evt-000118","type":"data.landed","booking_id":"bk-4912","sat_id":"furrow-a","station":"SQT","file":"furrow-a-20280201-1402.tif","bytes":10912500000,"landed_at":"2028-02-01T14:13:02Z"}`
- **Delivered** means any 2xx answer. **Retried**: transport errors (refused connections,
  timeouts), 408, 429 and 500 to 599. **Not retried**: any other answer, 3xx included (the
  client never follows redirects); the event is undelivered at once. A canceled context (the
  service stopping) is not a failure: the delivery keeps its attempt count for later.
- **The schedule.** At most 8 attempts. The first goes out as soon as the landing arrives (on launch day
  a second after the file lands, 71 s after LOS). After failed attempt n (n = 1 to 7) the sender waits `backoff(n)`:
  2 s doubled for each earlier failure (`2s × 2^(n−1)`), capped at 5 minutes, times a jitter
  factor drawn uniformly from [0.5, 1.0) from the sender's seeded `math/rand/v2` source. So
  the waits fall in 1–2 s, 2–4 s, 4–8 s, 8–16 s, 16–32 s, 32–64 s and 64–128 s. Each attempt is
  signed again with its own timestamp. After the 8th failure the event is undelivered
  (`ErrUndelivered`) and goes on Yusra's morning list.
- **Receivers** (Yusra's sample for customers, `verifyEvent`): read at most 64 KiB of body;
  refuse a timestamp more than 300 seconds from their own clock either way (exactly 300 is
  accepted); recompute the signature and compare with `hmac.Equal` before trusting anything in
  the body; drop an event whose `id` they have already handled; answer 2xx within 5 seconds.
- **Event IDs** are `evt-` and six digits, counting up per service from `evt-000001`. From
  chapter 11 the count is kept in the outbox, so a restart never reuses an ID.
- **Idempotency keys.** `POST /bookings` may carry `Idempotency-Key`: 1 to 64 letters, digits,
  `-` or `_` (else 400). Keys belong to the customer (Brinewatch's `k1` is not Cirrine's
  `k1`) and are honored for 24 hours on the injected clock; after that the same key starts afresh. The service keeps the
  request body's SHA-256 with the answer's status, `Content-Type` and `Location` headers and
  body. The same key with the same body gets that answer again with
  `Idempotent-Replayed: true` and books nothing; the same key with another body answers 409;
  a second request while the first is still running answers 409 (in progress). Answers of 500
  and above are not kept, so a retry of those runs for real. `idempotent` reads the body
  through `http.MaxBytesReader` (413 when too long) and hands the handler a fresh reader of
  the same bytes.

#### Starting and stopping (chapter 11)

- `run` serves `app.Handler` on the given listener with `newServer`, until `ctx` is done. Then,
  in this order: the service starts draining (`GET /health` answers 503 `draining`, every
  other route still works) for `Timeouts.Drain`; `Shutdown` runs with a deadline of
  `Timeouts.Shutdown`; the notifier stops and saves the outbox; `run` returns nil. `Serve`'s
  `http.ErrServerClosed` is never returned as a failure.
- Handlers' contexts come from a base context (`BaseContext`) that `run` cancels only if
  `Shutdown`'s deadline passes; then it calls `Close` and returns the error wrapped as
  `shutdown: context deadline exceeded`. Handlers pass `r.Context()` to the store, never
  `context.Background()`; deliveries run under the notifier's own context, which `Stop` cancels.
- `main` turns SIGINT and SIGTERM into a canceled context with `signal.NotifyContext`, and is
  the only code that ever listens on a real address. No lesson runs it.
- **The outbox**: `outbox.json` in the data folder, saved whole through `replaceFile` on every
  change: `{"next_event":118,"pending":[<Delivery>...]}`. At start, pending deliveries whose
  `next_at` has passed are tried at once and the rest at their time; event IDs continue from
  `next_event` even when nothing is pending. A missing file is an empty outbox with
  `next_event` 1; a broken one stops the start.
- On launch day the outbox holds Cirrine's `evt-000117` (the file of Cirrine-1's 13:30 pass over Ystmark,
  `cirrine-1/YST/12218`, which landed at 14:04:57 over Ystmark's slow backhaul; its receiver
  was down for its own deploy, and two attempts were made) when the service restarts; Furrowlight's
  launch-day event is `evt-000118`.

#### The weather service (chapter 2, and the Board)

- `GET /weather/{code}` answers 200 with one reading:
  `{"station":"CWR","at":"2027-11-08T13:55:00Z","wind_kmh":64,"gust_kmh":81,"temp_c":4.5}`.
  It sometimes adds fields (`humidity_pct` since last spring), so its client ignores unknown
  fields: the opposite of the booking service's own strict input, on purpose.
- Its errors are plain text: 404 `no station "CRW"`; 503
  `mast silent since 2027-11-08T12:51:00Z` with `Retry-After: 60`.
- A reading is fresh for 15 minutes after its `at`: exactly 15 minutes old is still fresh.
- Its client never waits more than 2 seconds (`http.Client{Timeout: 2 * time.Second}`), and
  each call also takes the poll's context.
- The go/no-go wind call uses the sustained wind (`wind_kmh`), as Season 1 did: under 50 km/h
  `GO: wind 22 km/h`; from 50 `GO: wind 61 km/h, warning`; from 72
  `NO-GO: wind 74 km/h, stow`; a stale reading `NO-GO: reading 47m0s old`; a 503
  `NO-GO: mast silent`; any other failure `NO-GO: no weather`. Wind prints with no decimals.
  The rules apply in this order, the first that fits answering: a 503, any other failure, a
  stale reading, then the wind from 72 km/h down.
- `decodeWeather` decodes leniently and refuses a reading without its station or its time:
  `reading has no "station"`, `reading has no "at"`.

```go
// Weather is one station's latest reading from the weather service.
type Weather struct {
	Station string    `json:"station"`
	At      time.Time `json:"at"`
	WindKmh float64   `json:"wind_kmh"`
	GustKmh float64   `json:"gust_kmh"`
	TempC   float64   `json:"temp_c"`
}

// WeatherClient asks the weather service for station readings.
type WeatherClient struct {
	BaseURL string       // the service's address, in tests an httptest server's URL
	HTTP    *http.Client // with a Timeout
}

func (c *WeatherClient) Reading(ctx context.Context, code string) (Weather, error)

// StatusError is an answer from the weather service other than 200 OK.
type StatusError struct {
	StatusCode int    // 503
	Status     string // "503 Service Unavailable"
	Body       string // at most the first 512 bytes it sent, trimmed
}

func (e *StatusError) Error() string { return e.Status + ": " + e.Body }
```

`Reading` wraps its errors with the station: `weather for CWR: 503 Service Unavailable: mast
silent since 2027-11-08T12:51:00Z`.

#### Nkem's stand-in (chapter 1)

`standIn()` returns an `http.Handler`, a mux with two routes: `GET /stations` answers 200 with
`[{"code":"YST","name":"Ystmark"},{"code":"PEL","name":"Port Ellery"},...]` (code and name
only, table order) and `GET /stations/{code}` answers `{"code":"PEL","name":"Port Ellery"}`,
or 404 `no station "ZZZ"` in plain text. Both JSON answers end in a newline. Requests built
only to be printed use URLs on `http://127.0.0.1`, so a printed `Host` line reads
`127.0.0.1`; no lesson prints a port.

#### How lessons run servers and test them

- Every server is `httptest.NewServer` or `httptest.NewUnstartedServer` (which listen on
  127.0.0.1 on a port the system picks), or an `http.Server` serving
  `net.Listen("tcp", "127.0.0.1:0")`. Never a fixed port, never `localhost` or `:0` (all
  interfaces), never an address outside 127.0.0.1. Nothing reaches the network, and nothing
  asks macOS for permission.
- One handler call goes through `httptest.NewRecorder` and `httptest.NewRequest`. A check of
  a 204 or a HEAD answer goes through `httptest.NewServer`: a recorder keeps a body that a real
  server drops.
- Clients in tests are `srv.Client()` or have a timeout of at most 2 seconds. Every wait in a
  test has a guard of about one second, so a starter fails within a second or two and never
  hangs the run.
- `httptest.Server.Close` waits for requests in flight. A handler blocked on purpose waits on
  a channel the test closes or on `r.Context().Done()`, and streams are canceled first: cleanups
  run last-in first-out, so register `srv.Close` before the cleanups that release handlers.
- Clocks are `now func() time.Time` from a small fake clock with `Now` and `Add` (under a mutex
  when goroutines read it); randomness is `rand.New(rand.NewPCG(1, 2))`; sleeps are 50 ms at
  most; a timeout under test is 50 ms, and a request that must succeed uses the test server's
  `srv.Client()` or a client with a 2-second timeout, so a busy Mac never fails it.
- Tests compare JSON after decoding it, except where a lesson's point is the exact bytes.
  Messages read `call = got, want want`: `GET /stations/ZZZ = 200, want 404`.
- Data files live under `t.TempDir()`. The guide in tests is an `fstest.MapFS`.
- Where learners write checks themselves (chapters 7, 8 and 11), their code goes in an
  editable non-test file and the read-only test proves it against planted mutants in a
  re-executed test binary (`os.Args[0] -test.run=...` with `APOGEE_MUTANT` naming the
  mutant), as in Season 1's chapter 15.
- A flood of hundreds of requests goes through `ServeHTTP` with recorders, not the loopback.

**Topics added this season** (front matter): `http`, `json`, `security`. Lessons keep using the
older ones (`context`, `errors`, `interfaces`, `mutexes`, `testing`, `time`...).

**New field guide words** (world words only; Go and HTTP words are taught in readmes):
booking, booking service, bookingd, the beta, catalog, customer guide, weather service, weather
mast, landing.

**Name checks** to run before the season ships: `bookingd` (and the header names
`Apogee-Timestamp` and `Apogee-Signature`, which use the invented company's name). Nothing else
this season is a new proper name.

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
