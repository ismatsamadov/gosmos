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
forward up to a week or two per chapter (chapters 9 and 10 share the week of 11 October). Lessons give an altitude or a day, never a new place
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
| 10 | 11 Oct, Fri 15 to Sun 17 | 3,700 m, the last weekend before the climb | a stick brings two parcels from La Grande's dish controller, its settings and its log since July: read their XML, unpack them safely, write the controller's correction file |
| 11 | 18 Oct | 3,840 m, Salar Quieto | calibration day, run from gctl and files alone; La Grande points true |

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
- **Booking** (chapters 4, 5, 7, 8 and 11). One reserved pass, as Yusra's booking export
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

- **Booking sheet** (chapters 6 and 11). Hollis exports the same bookings from Furrowlight's
  spreadsheet as CSV, header first, in capitals: `ID,SATELLITE,STATION,AOS,LOS,NOTE`.
  Columns may come in any order and are found by their header. A note may hold commas,
  quotes and line breaks, and Hollis writes notes in capitals too. Times are written as in
  the JSON.
- **Settings** (chapters 5 and 11). `settings.json`, one JSON object, read strictly: an
  unknown field, a second value or anything after the object is an error, never ignored.
  `station` defaults to `SQT`, `min_elev_deg` to 3 and `timeout` to `30s`; `day` is required
  and is a date (`2006-01-02`). A timeout is duration text such as `30s` or `2m`, never a
  count of nanoseconds: `type Delay time.Duration`, written back as `Duration.String()`.

  ```json
  {"station":"SQT","day":"2027-10-21","min_elev_deg":3,"timeout":"30s"}
  ```

- **Station table** (chapters 8, 9 and 11). gctl carries `stations.csv` inside its binary,
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
- **Parcel** (chapter 10). Salar Quieto's dish controller packs files it sends into one parcel: a zip
  for a dish's settings, a gzip'd tar for its log, each with `manifest.xml`. On Friday 15 October a
  stick brought two down the road to the van at 3,700 m, both made at 2027-10-15T09:40:00Z:
  `lg-settings-2027-10-15.zip` (`manifest.xml`, `correction.xml`, `log/2027-10-14.txt`) and
  `lg-log-2027-10-15.tar.gz` (`manifest.xml` and the controller log since Paz left, one file per
  week named by its Monday, `log/2027-07-05.txt` to `log/2027-10-11.txt`). In the manifest,
  `version` is an attribute of `<parcel>`; `dish`, `made` (RFC 3339 UTC) and each `file` inside
  `files` are elements, a file's size in bytes its `size` attribute and its name its text. The
  controller's firmware, 5.3, adds a `<firmware>` element gctl skips. A manifest with no dish, time
  or file is refused (`manifest: no <dish>`, `manifest: no <made>`, `manifest: no <file>`); one name
  `filepath.IsLocal` rejects refuses the whole parcel before anything is written
  (`unsafe name "../gctl"`); files are written only through an `os.Root`, none past 4 MiB
  (`maxFileBytes`; `<name>: more than 4194304 bytes`).

  ```xml
  <parcel version="1">
    <dish>LG</dish>
    <firmware>5.3</firmware>
    <made>2027-10-15T09:40:00Z</made>
    <files>
      <file size="153">correction.xml</file>
      <file size="136">log/2027-10-14.txt</file>
    </files>
  </parcel>
  ```

- **Correction** (chapter 10). `correction.xml`: the angles La Grande's controller adds to every
  pointing, in degrees; until calibration day az 0.18, el −0.06. The controller takes a new one only
  in this shape, the XML header first and two spaces per level; calibration day replaces it.

  ```xml
  <?xml version="1.0" encoding="UTF-8"?>
  <correction dish="LG">
    <az_offset_deg>0.18</az_offset_deg>
    <el_offset_deg>-0.06</el_offset_deg>
  </correction>
  ```

- **Commands** (chapters 8 and 11). `run(args []string, stdin io.Reader, stdout, stderr
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
  to standard error, nothing on standard output, and returns 1. Chapter 11's copy of gctl
  carries only `day [-format json|csv] dir`, which runs a calibration day from a folder
  holding `settings.json` and `bookings.json` or `bookings.csv`.
- **Day report** (chapter 11). `<dir>/reports/daily/<day>.txt`, saved whole through a
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

Story time runs from Monday 1 November 2027 to Tuesday 1 February 2028: twelve chapters,
all but chapter 11 of about a week each (chapter 11 is the afternoon of Friday 21 January), with a
quiet holiday fortnight (Monday 20 December to Sunday 2 January)
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
| 11 | Fri 21 Jan, 15:00 to 18:20 UTC | the backend desks | Hollis, right after his capitals about the five-second loop: he wants to see Furrowlight's bookings himself, in a browser; Yusra wants a signed-in page beside the API before rehearsal week, and Mirela reviews it | A signed-in browser can be made to cancel a booking from any other website. At 18:20 Hollis signs in and sees only Furrowlight's bookings; Mirela's page with a hidden cancel form is refused |
| 12 | Mon 24 Jan, then Tue 1 Feb | the Fishbowl | Ingrid runs launch like a pass: a rehearsal week, a full rehearsal on Monday 31 January, then launch at 14:00 UTC on Tuesday 1 February | At 14:01 Yusra finds that the customer guide's sample request says `Authorization: Token` instead of `Bearer`, and customers are copying it. The guide is built into the service, so it must restart at 14:05, while Furrow-A is over Salar Quieto (AOS 14:02:10, LOS 14:11:52) |

The opening (Monday 1 November) hands over the new badge, shows Yusra's spreadsheet and the
launch date, introduces Nkem (reviewer), Yusra (the customers' voice) and Mirela (who will
review the service in January), and sends the learner to Nkem's stand-in.

The finale (chapter 12, Tuesday 1 February): the restart at 14:05 drops no request; the
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
- **The account page** (chapter 11): a small signed-in page that bookingd serves beside the API under
  `/account/`, where a customer's people sign in with a browser using their company's API key, see
  their own bookings and cancel one. Its sessions live only in memory, so a deploy or restart signs
  every browser out and nothing else is lost; the launch-day restart signs page users out too. In
  readmes: "the account page".
- **Furrowlight's planning page**: its engineers' own web page at `https://plan.furrowlight.example`
  (`.example` is reserved), whose script calls the API with a token that Furrowlight's server got. It
  is the one origin the API trusts for CORS (chapter 11).
- Customers reach the booking service only over HTTPS, which no lesson sets up; the cookie's `Secure`
  and the page's `Strict-Transport-Security` rely on it, and lessons use plain HTTP on 127.0.0.1.

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
| 11 | `cookieName`, `sessionLife`, `sessionCookie` (`cookie.go`); `session`, `Sessions`, `newSessions`, `(*Sessions).Start`, `Lookup`, `End` (`sessions.go`); `Account`, `(*Account).routes`, `signedIn`, `page`, `signInForm`, `signIn`, `signOut`, `cancel` (`account.go`); `pageHTML`, `signInHTML`, `render`, `renderPage`, `renderSignIn` (`page.go`); `noCrossOrigin`, `cors`, `securityHeaders` (`browser.go`); `site` (11.9's `main.go`) |
| 12 | `stop`, `serve`, `run`, `Outbox`, `OpenOutbox`, `(*Outbox).NewEventID`, `Put`, `Remove`, `Pending`, `outboxFile`, `(*Notifier).Start`, `Stop`, `(*Notifier).deliver`, `sleep`, `App`, `newApp`; 12.8's `client.go` (`request`, `answer`, `errorCode`, with chapter 8's `send` and `token`) (`run.go`, `outbox.go`, `webhook.go`, `app.go`) |

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
func (n *Notifier) Start()      // chapter 12: resumes the outbox's deliveries, each when due
func (n *Notifier) Stop() error // chapter 12: ends deliveries in progress, keeping their attempts
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

// chapter 11: the account page
func sessionCookie(id string) *http.Cookie
func newSessions(now func() time.Time) *Sessions
func (s *Sessions) Start(customer string) string     // a rand.Text() ID, ending sessionLife after now
func (s *Sessions) Lookup(id string) (string, bool) // its customer; renews it; an ended one is deleted
func (s *Sessions) End(id string)
func (a *Account) routes() *http.ServeMux
func (a *Account) signedIn(r *http.Request) (string, Customer, bool) // session ID, customer, ok
func noCrossOrigin(h http.Handler) http.Handler
func cors(origin string) func(http.Handler) http.Handler
func securityHeaders(next http.Handler) http.Handler
func site(api http.Handler, a *Account, origin string) http.Handler

// checks the learner writes, in editable non-test files
func checkStore(t *testing.T, open func() Store)                             // chapter 7, contract.go
func review(t *testing.T, base string, advance func(time.Duration))          // chapter 8, review.go
func checkService(t *testing.T, base string)                                 // chapter 12, check.go
```

`App` (chapter 12) is the assembled service: `Service *Service`, `Notifier *Notifier` and
`Handler http.Handler`, the mux inside every middleware. `Notifier` holds the `Store`, the
customers, the `Sender` and (chapter 12) the `Outbox`, and counts event IDs.

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
	draining  atomic.Bool // chapter 12
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
	Timeouts  Timeouts         // chapters 9 and 12
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
| `GET /health` | 200 `{"status":"ok"}`; 503 `{"status":"draining"}` while shutting down | public | 3; 12 |
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

- **API keys** are given to customers by Yusra and sent only to `POST /tokens`, and from chapter 11 once per sign-in in the account page's form (field
  `key`), in
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
  chapter 12 the count is kept in the outbox, so a restart never reuses an ID.
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

#### The account page (chapter 11)

- **Routes:** `GET /account/{$}` → the page, or 303 to `/account/sign-in`; `GET /account/sign-in` → the form;
  `POST /account/sign-in` → form field `key`, body at most 4,096 bytes, checked with `checkKey`: 303 to
  `/account/` with the session cookie, or 401 with the form and `invalid api key`; `POST /account/sign-out`
  → ends the session, clears the cookie, 303 to the form; `POST /account/bookings/{id}/cancel` → only the
  caller's own booking (another customer's is a plain-text 404 `booking bk-5499: not found`, exactly
  like a missing one), then 303 to `/account/`. Pages are HTML (`text/html; charset=utf-8`), rendered into
  a buffer from a template; errors are plain text.
- **Cookie:** `session=<rand.Text() ID>; Path=/account; Max-Age=1800; HttpOnly; Secure; SameSite=Lax`,
  cleared with the same name and path and `Max-Age=0`.
- **Session:** valid while now is before its end; ended at that very instant and deleted; every
  successful lookup renews it for 30 minutes, and the page re-sends the cookie.
- **Guards:** `http.CrossOriginProtection` (Go's own 403) around every account route; `cors` around the
  API only: `Vary: Origin` everywhere and, for exactly the planning page, `Access-Control-Allow-Origin`,
  its preflight answered 204 with `Allow-Methods: GET, POST, DELETE` and `Allow-Headers: Authorization,
  Content-Type, Idempotency-Key`; on every account answer `Content-Security-Policy: default-src 'self';
  frame-ancestors 'none'`, `Strict-Transport-Security: max-age=31536000`, `X-Content-Type-Options:
  nosniff`, `Referrer-Policy: same-origin`, plus `Cache-Control: no-store`.
- **Wiring:** `site` is a mux with `/account/` → `chain(a.routes(), securityHeaders, noStore,
  noCrossOrigin)` and `/` → `cors(origin)(api)`. The API's own mux still answers its 404 and 405 (the
  rule against a catch-all `/` is about the API's mux). Launch Day's lessons leave the page out of
  `newApp`.
- **Fixtures:** new bookings `bk-4776` (Furrowlight, `furrow-b/SQT/16607`, Mon 24 Jan 14:19:53–14:29:53,
  booked before Friday), `bk-4784` (Yusra's test booking, `pipit-1/PEL/08904`, Mon 24 Jan
  15:08:30–15:17:50, made Friday afternoon) and `bk-4785` (the planning page's first booking, from 11.9's
  stand-in API); `bk-4779`, `bk-4781` and `bk-4783` come from chapter 10; `bk-4786` to `bk-4800` stay free
  before Launch Day's `bk-4801`. Furrowlight's test tokens expiring at 18:00 and 19:00 UTC are
  `furrowlight.1832090400.8e8c…` and `furrowlight.1832094000.86f0…`, under the test token key. Every
  pass was checked with `cmd/passes`.

#### Starting and stopping (chapter 12)

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
- Where learners write checks themselves (chapters 7, 8 and 12), their code goes in an
  editable non-test file and the read-only test proves it against planted mutants in a
  re-executed test binary (`os.Args[0] -test.run=...` with `APOGEE_MUTANT` naming the
  mutant), as in Season 1's chapter 15.
- A flood of hundreds of requests goes through `ServeHTTP` with recorders, not the loopback.

**Topics added this season** (front matter): `http`, `json`, `security`. Lessons keep using the
older ones (`context`, `errors`, `interfaces`, `mutexes`, `testing`, `time`...).

**New field guide words** (world words only; Go and HTTP words are taught in readmes):
booking, booking service, bookingd, the beta, catalog, customer guide, weather service, weather
mast, landing, account page, planning page.

**Name checks** to run before the season ships: `bookingd` (and the header names
`Apogee-Timestamp` and `Apogee-Signature`, which use the invented company's name). Nothing else
this season is a new proper name. The planning page's host, `plan.furrowlight.example`, uses the
reserved `.example`.

### The Board project: every screen in the building

The booking service opened on Tuesday 1 February 2028. On Monday 7 February Ingrid asks for
one more thing: the Board, the big screen in the Fishbowl that lists the next passes, is still
a laptop plugged into a television, and when it rebooted during last week's night shift the
Fishbowl had no Board for twenty minutes. She wants it served by a program and shown on every
screen in the building: forty of them, in the kitchen, at reception, in the meeting rooms, in
Tomás's office and on every night-shift desk. The guided project **Build the Board** runs from
Monday 7 to Thursday 10 February 2028 at Port Ellery, one chapter a day, and ends in the night
into Friday 11 February, when the Board runs on all forty screens through Cape Wren's wind and
a deploy at 05:00. Nkem reviews; Ingrid is the customer; Tomás asks why his old weather
service is suddenly so busy; Hollis visits reception at 09:00 on Friday. The badge is
unchanged (Web Services Engineer, No. 039). No new people, places, satellites or customers.

| Ch | Day (2028) | What the team needs |
|---|---|---|
| 1 | Mon 7 Feb | the Board as a page any screen can load, with customers' notes shown as text, each station's wind and who is on call |
| 2 | Tue 8 Feb | screens that keep up without reloading: the data as JSON, then a stream that pushes every change |
| 3 | Wed 9 Feb | every change on all forty screens, and no screen able to hold the others up |
| 4 | Thu 10 Feb, and the night into Fri 11 Feb | the weather service spared, streams that end when the Board restarts, the whole program run as the night shift runs it, and a reception screen fit for a visitor |

#### The program

One `package main` program grows lesson by lesson, each lesson's starter being the previous
lesson's solution, plus the read-only test and any new read-only file the lesson gives. The
Board reuses the Season 3 types `Booking`, `Weather` and `Station` (code and name only)
exactly.

- `board.go`: `Board`, `NewBoard`, `Refresh`, `Handler` and the view; `page.go` and
  `board.html` (the template, embedded with `//go:embed`, parsed once by `NewBoard`);
  `json.go`: `/board.json`; `events.go`: `/events`; `hub.go` (chapter 3); `cache.go`
  (chapter 4); `run.go` (chapter 4); `data.go` (read-only): the station table and the rota
  for `main`; `main.go`, which wires the real sources and serves on 127.0.0.1 on a port the
  system picks, printing the address. In the building it would listen on the office network;
  no lesson runs `main`.
- Sources, given to `NewBoard` and faked by every test:

```go
// Sources is where the Board's data comes from.
type Sources struct {
	Passes  func(ctx context.Context) ([]Booking, error)            // passplan's schedule: every booked pass, in any order
	Weather func(ctx context.Context, code string) (Weather, error) // Season 3's WeatherClient.Reading
	Rota    []Shift                                                  // the on-call rota, oldest first
}

// Shift is one week on call; it starts on a Monday at 08:00 UTC.
type Shift struct {
	From time.Time
	Name string // "Nkem Halloran"
}

// View is everything the Board shows at one moment.
type View struct {
	Passes  []Booking `json:"passes"`  // LOS after now, by AOS then ID, at most 12; [] when none
	Weather []Panel   `json:"weather"` // one per station, in the station table's order
	OnCall  string    `json:"on_call"` // the latest shift that started at or before now
}

// Panel is one station's weather on the Board.
type Panel struct {
	Station string   `json:"station"`           // "CWR"
	Weather *Weather `json:"weather,omitempty"` // nil when there is no reading
	Call    string   `json:"call"`              // "track", "warn", "stow", or "none" without a reading
}

// pageData is what board.html is executed with. View is embedded, so the template
// says {{range .Passes}}; Lobby is true for the reception screen (chapter 4).
type pageData struct {
	View
	Now   time.Time
	Lobby bool
}

func NewBoard(src Sources, stations []Station, now func() time.Time) (*Board, error)
func (b *Board) Refresh(ctx context.Context) error
func (b *Board) Handler() http.Handler
func nextPasses(all []Booking, now time.Time) []Booking
func newServer(b *Board) *http.Server
func run(ctx context.Context, ln net.Listener, b *Board, ticks <-chan time.Time) error
```

- `Board` keeps its sources, stations and clock, `page` (the parsed template), the current
  view with its JSON and `version` under a mutex, `updates` (chapter 2), `hub` (chapter 3) and
  `weather`, the cache (chapter 4).
- **Refresh** builds a new view at `b.now()`: the passes through `nextPasses` (a new slice;
  the source's slice keeps its order), a panel per station (a weather error makes that panel
  `none`; it never fails the refresh), and the on-call name. A `Passes` error keeps the last
  view and is returned. The view is marshaled; only when its JSON differs from the last one
  does the version go up by one and the new view go to the streams. The version counts from 1,
  so the first refresh makes version 1; before any refresh the view is empty, version 0.
  `NewBoard` doesn't refresh; `run` does at start and on every tick (every 15 seconds in
  `main`), and tests call `Refresh` themselves.
- **The wind call** is Season 1's Cape Wren rule on the sustained wind: `track` under 50 km/h,
  `warn` from 50, `stow` from 72, and `none` without a reading.
- **On call**: the rota changes hands every Monday at 08:00 UTC: Tomás Ferreira-Lund from
  31 January (he took launch week himself), Nkem Halloran from 7 February, Mirela Dunn from
  14 February. At Monday 07:59 the Board still shows last week's name.

#### What it serves

| Route | Answer | From |
|---|---|---|
| `GET /` | the page, `text/html; charset=utf-8`, rendered into a buffer first: a template error answers 500 and sends none of the page | 1.1 |
| `GET /board.json` | the view's JSON, `Cache-Control: no-cache`, `ETag: "v<version>"`; a matching `If-None-Match` answers 304 with no body | 2.1, 2.2 |
| `GET /events` | a stream of server-sent events (below) | 2.3 |
| `GET /health` | `{"status":"ok","viewers":3}`, the streams open now | 2.3 |

`?screen=lobby` on `/`, `/board.json` and `/events` (chapter 4) gives the reception screen's
view: every booking's `id` and `note` blank, and no on-call name (`"on_call":""`). Its ETag is
`"v<version>-lobby"`; its events carry the same ids.

**The page.** Tests look for these lines, each on one line of the HTML; everything else is the
writer's choice.

- `<title>The Board</title>` and `<h1>The Board · 03:58 UTC</h1>`, the time `b.now()` at
  rendering.
- A pass: `<tr class="pass"><td>04:12</td><td>04:21</td><td>PEL</td><td>pipit-1</td><td>bk-4988</td><td>echo test</td></tr>`:
  AOS and LOS as `15:04` in UTC, the station, the satellite, the booking ID and the note,
  which is never marked as trusted HTML, so `<b>RUSH</b> & PRIORITY` appears as
  `&lt;b&gt;RUSH&lt;/b&gt; &amp; PRIORITY`. In the lobby the last two cells are empty.
- A panel: `<li class="panel track">YST wind 41 km/h, gusts 55</li>`, or without a reading
  `<li class="panel none">SQT no reading</li>`, from a `panel` template defined in
  `board.html`. Wind prints with no decimals.
- `<p class="oncall">On call: Nkem Halloran</p>`, left out in the lobby.
- From chapter 2 a small script (given, never checked) follows `/events` and falls back to
  fetching `/board.json` every 15 seconds. No check runs a browser; a learner may open the
  page in their own browser at 127.0.0.1.

**The stream.** `Content-Type: text/event-stream` and `Cache-Control: no-store`. The stream
subscribes first, then writes the current view, so no version is missed between the two; then
it writes every version it receives that is newer than the last one it wrote. (Until lesson
3.3 the given stream reads the view first and subscribes after: the bug that lesson finds.)
Each event is three lines and a blank line, the JSON on one line:

```
event: board
id: 3
data: {"passes":[...],"weather":[...],"on_call":"Nkem Halloran"}

```

It flushes after every event through `http.NewResponseController`. When its viewer leaves (the
request's context ends) it unsubscribes and returns. When its channel is closed, because the
viewer fell behind or the Board is shutting down, it writes `event: bye`, `data: {}` and a
blank line, and returns.

**The hub** (chapter 3) fans each new view out to every stream:

```go
func NewHub() *Hub
func (h *Hub) Subscribe() <-chan View    // a new channel with room for 4 views; closed at once after Close
func (h *Hub) Unsubscribe(ch <-chan View) // forget it and close it; harmless if it is already gone
func (h *Hub) Publish(v View)            // never blocks: a channel with no room is removed and closed
func (h *Hub) Close()                    // close every channel; harmless twice
func (h *Hub) Len() int                  // the channels it holds: /health's viewers
```

Only the hub closes a channel, under its mutex, and only once. In chapter 2, before the hub,
`Refresh` hands new versions to one stream through `b.updates`, a channel with room for 16, and
`/health` counts streams with a counter the stream moves.

**The weather cache** (chapter 4) keeps each station's last reading with the time it was
fetched, on the Board's clock, and gives it back for 60 seconds; at exactly 60 seconds it
fetches again (as gctl's pass cache went stale at its time to live). Errors are not kept.

**Stopping** (chapter 4): `newServer` returns an `http.Server` for `b.Handler()` with no
write timeout (streams are long) that registers `b.hub.Close` with `RegisterOnShutdown`,
because `Shutdown` waits only for connections that go idle, and a stream never does. `run`
serves it on the listener it is given and, when `ctx` ends, shuts down within 2 seconds and
returns nil. `http.ErrServerClosed` is never a failure.

#### The night into Friday 11 February 2028 (UTC)

The finale's data. On call: Nkem Halloran. Passes, as passplan sends them:

| Booking | Satellite | Station | AOS | LOS | Note |
|---|---|---|---|---|---|
| `bk-4986` | cirrine-1 | PEL | 03:20:00 | 03:30:00 | |
| `bk-4988` | pipit-1 | PEL | 04:12:05 | 04:21:27 | `echo test` |
| `bk-4990` | pipit-2 | YST | 04:31:10 | 04:40:40 | |
| `bk-4991` | furrow-a | CWR | 04:58:30 | 05:07:13 | `<b>RUSH</b> & PRIORITY` (Hollis) |
| `bk-4992` | furrow-b | SQT | 05:12:44 | 05:22:10 | `RESHOOT FIELD 12` |
| `bk-4993` | brine-3 | LAN | 05:20:15 | 05:27:00 | `harbor survey` |
| `bk-4994` | cirrine-1 | YST | 05:47:00 | 05:57:17 | |

Weather at 03:50: YST 41 km/h (gusts 55), PEL 22 (30), SQT 12 (18), LAN 18 (25), CWR 46 (58).
Cape Wren rises to 61 (gusts 77) at 04:30, `warn`, and 74 (gusts 88) at 04:40, `stow`: its
dish stows itself during the night, as it did in the gctl project. At 05:00 the Board is
deployed again so reception gets its lobby view before Hollis arrives; every stream says
`bye`, and the screens reconnect to the new process.

#### How the lessons test it

- Every test builds the whole program with `NewBoard`, fake sources and a fake clock, and
  checks it through HTTP: `httptest.NewServer(b.Handler())` and its client, or
  `httptest.NewRecorder` for a single page. Earlier features stay checked in every later
  lesson, as gctl's transcripts did.
- A stream is read by a client goroutine line by line with a guard of about a second, and
  cancelled through its request's context before the server closes (cleanups run last-in
  first-out, so register `srv.Close` first): `httptest.Server.Close` waits for handlers.
- A slow screen is a test `ResponseWriter` with `Flush` whose `Write` signals when it starts
  holding, then blocks until the test releases it; the stream handler runs with it in a
  goroutine, through `httptest.NewRequest`. In 3.2 it lets the first event through and holds
  from the second, which comes only after the stream has subscribed (the test waits for
  `/health` to count it), so the drop happens whichever order the stream uses. In 3.3 it holds
  the very first write, and the test refreshes while it does: only a stream that subscribed
  first sees that version. Both are correct by construction.
- Forty screens (3.3) are forty loopback clients; the test polls `/health` until it reports the
  count it expects, within a second, and compares each viewer's last event, since a viewer may
  join before or after any one refresh.
- The night (4.4) is `testdata/night.txt`, the moments of the night and what changes at each,
  and `testdata/*.golden`, the page, the lobby page, the JSON views and the stream's events
  written by the solution at each moment.
- No sleeps over 50 ms; every run well under the app's 10 seconds.

**Topics added**: `templates` (with Season 3's `http` and `json`).

**New names**: none. "The Board" and every person and place above are already in the bible.

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
a bare socket so you see what every web request of Season 3 really was.

Story time runs from Monday 6 March to Tuesday 20 June 2028: nine chapters, about a week each,
with installation weeks between some of them. Sunniva carries the season; Nkem reviews every
change from Port Ellery; Tomás flies up for Dish 1's calibration; Ingrid runs Ystmark's passes
from the Fishbowl; Mirela reviews the station server; Yusra and Hollis bring the customers'
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
| 9 | Mon 19 Jun, then Tue 20 Jun | the hut; the Fishbowl on the Board | Ingrid and Yusra: headquarters' collector must hear of each recording at once and fetch it over one open line. On Tuesday 20 June the station server takes every pass from 06:00 | Furrow-A's 06:32 pass (`furrow-a/YST/18856`, LOS 06:40:35) is the first Furrowlight pass through it; its image must reach Furrowlight within a minute of LOS, as the booking service promises |

The opening (Monday 6 March) shows Ystmark in one image (the sun back for three weeks, two
radomes on a ridge, −18 °C), hands over the badge, introduces Sunniva (station chief), Nkem
(reviewer, at Port Ellery), Tomás (coming up in March), the crate and the island's line, and
sends the learner to the simulator on the hut's bench: the old controllers come out on the 20th.

The finale (chapter 9, Tuesday 20 June): Furrow-A's recording streams into the station server
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
  with `net/http`'s client.
- **The station server** (chapters 7 to 9): the hut's server that keeps every finished
  recording until headquarters collects it, and tells headquarters when one arrives. It has no
  program name; readmes say "the station server".
- **The collector**: dlinkd's part at Port Ellery that follows the station server's feed and
  fetches each recording. A Go program using `net/http`'s client.

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

**Stopping** (`Shutdown(ctx)`, chapter 9): close the listener, so no new connection is
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

// chapter 9
func statusLine(code int) string
func writeResponse(w io.Writer, now time.Time, status int, fields []Field, body []byte) error
func writeChunk(w *bufio.Writer, p []byte) error
func endChunks(w *bufio.Writer) error
func (s *Server) recording(req *Request) answer
func (s *Server) feed(bw *bufio.Writer, lines <-chan string) error
func (s *Server) route(req *Request) answer
func (s *Server) Serve(ln net.Listener) error
func (s *Server) serveConn(c net.Conn)
func closing(req *Request) bool // given from 9.5 (conn.go): Connection: close, any letter case
func (s *Server) Shutdown(ctx context.Context) error
```

`writeResponse` is given read-only in chapters 7 and 8 and written in chapter 9, as
`encodeFrame` is given in chapter 2 and written in chapter 3. In chapters 7 and 8 the given
`serveConn` answers one request per connection and closes, and the given route answers every
good request with the list. A chapter 8 lesson gives only what earlier lessons built, so no
editable file shows a later lesson's answer: chunked bodies reach the server from 8.4 (before,
the given `bodyLength` answers Transfer-Encoding with 501), the CRC from 8.8, and until 8.9 a
given `put` keeps any complete body and answers 201.

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
| 9 | `statusLine`, `writeResponse`, `writeChunk`, `endChunks`, `recording`, `feed`, `serveConn`, `Shutdown`, `route` (`response.go`, `server.go`, `routes.go`) | `request.go`, `body.go`, `store.go`, and `closing` (`conn.go`, from 9.5) | |

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

**Topics added this season** (front matter): `networking`. Lessons keep using the older ones
(`bytes`, `context`, `errors`, `goroutines`, `http`, `mutexes`, `streams`, `testing`, `time`...).

**New field guide words** (world words only; Go, network and HTTP words are taught in
readmes): backhaul, banner, beamwidth, collector, dish controller, dishd, feed, fiber, firmware,
jog, midnight sun, polar night, recorder, recording, simulator, station server, the hut, the
manual, watch port, wind lock. Carried over and needed again: AOS, LOS, pass, pass ID, radome,
stow, CRC, sequence number, X-band, S-band, UTC, the Board, the Fishbowl.

**Name checks** to run before the season ships: `dishd` (a quick search on 6 November 2026 found
no program of that exact name; similar names are unrelated), and the header name `Apogee-Crc32`,
which uses the invented company's name as Season 3's headers did. `YST-1` and `YST-2` are a
station code and a dish number, the simulator and the maker are never named, and "the hut" and
"the station server" are plain words. No new people, places, satellites or customers.

### Season 5: keeping secrets

After Season 4, the learner starts the next assignment on Monday 10 July 2028. Their new badge
says **Ground Security Engineer**, still No. 039, based at Port Ellery headquarters, since
Monday 10 July 2028. On Friday 7 July Cirrine Atmospherics signed a bigger contract with three
conditions: from Friday 25 August nobody but Cirrine and Apogee may read Cirrine-1's commands,
not even someone with a dish listening to the uplink; every file delivered to Cirrine comes
with a fingerprint Cirrine's archive can check; and every key Apogee uses for Cirrine is
replaced by Tuesday 31 October, then once a year. Mirela leads the work and reviews every
change ("LGTM once it fails closed"); she plans **Rotation Day** for Tuesday 17 October: every
key in the network replaced, at every station, without missing a single pass. Yusra carries
the customers' wishes, Tomás the radio, Ingrid the passes, Sunniva and Paz their stations, and
Nkem asks the whiteboard questions.

Story time runs from Monday 10 July to Tuesday 24 October 2028: twelve chapters: seven of about
a week, four (chapters 6, 9, 10 and 12) of about two, and chapter 11 of an afternoon and a night
(Friday 6 to Saturday 7 October). No new people, places, satellites or customers
are introduced. Customers speak through Yusra ("TL;DR: ...") and Furrowlight through Hollis's
capitals. The Cape Wren contractor (the bible's "a contractor visits monthly") stays unnamed.

#### The season, chapter by chapter

| Ch | Week of | Where | Who needs what | What is at stake |
|---|---|---|---|---|
| 1 | Mon 10 Jul 2028 | Mirela's desk, Port Ellery | Mirela, before anything is sealed: "Every secret this season starts as random bytes. Show me where ours come from." The draft nonce helper for Cirrine's commands and the ops console's session tokens both come from `math/rand/v2` seeded with the clock | Whoever knows roughly when a token or nonce was made can make it again. On Friday 14 July at 16:00 Mirela makes the season's master key, `k-2028-07`, with your helper |
| 2 | Mon 17 Jul | the backend desks, with dlinkd's delivery folder | Yusra: "TL;DR: Cirrine wants proof that every file is exactly what left the satellite." Cirrine's archive checks SHA3-256, Apogee's own deliveries SHA-256 | Season 1's CRC-32 catches radio noise, not a person: Mirela changes "0.71" to "0.17" in a Furrowlight field report and the CRC still matches. Cirrine's first fingerprinted delivery is on Friday 21 July |
| 3 | Mon 24 Jul | Mirela's desk | Mirela reviews the signing path before sealing builds on it, and finds three weak spots: Pipit-3's flight software checks a homemade signature, `sha256(key + command)`; a `\|` inside an argument signs like two arguments; the replay guard counts a command before checking its signature | Pipit-3 would obey a command nobody signed, and one forged sequence number could lock a satellite out of every later command. Mirela signs off on Friday 28 July after her twelve planted lines |
| 4 | Mon 31 Jul | the Fishbowl | On Sunday 30 July the Fishbowl's spare laptop goes missing from the night-shift desk. It held a backup of the ops console's user file: 26 accounts, every password an unsalted SHA-256. Ingrid logs it; Mirela: "Assume whoever has it has a fast computer and a list of common passwords" | Every password is changed by Friday 4 August, and the console's login reopens at 16:00 that day with records that make every guess at a stolen copy cost 600,000 rounds |
| 5 | Mon 7 Aug | the lab bench at Port Ellery | Mirela's key inventory: `k-2027-03`, the one key that signs commands for every satellite, sits in uplinkd, in the lab bench's spare Pipit-2 flight computer and in two old test scripts | Whoever reads it out of any one of them can command all seven satellites. Each satellite gets keys of its own, derived from `k-2028-07`, which only uplinkd holds; they go up on the passes of Thursday 10 August, and `k-2027-03` retires at 00:00 UTC on Friday 11 August |
| 6 | Mon 14 Aug, to Fri 25 Aug | Port Ellery's roof and the test van; the Fishbowl on Friday 25 August | Tomás points the test van's 1.2 m dish at Port Ellery's uplink and reads Pipit-2's commands (S-band, 2065 MHz) off the air in plain text. Cirrine's condition starts on Friday 25 August | Cirrine-1's first sealed command, `PING`, goes up at 07:30 UTC on Friday 25 August, on pass `cirrine-1/YST/15287` (Ystmark, LOS 07:41:12), and Cirrine-1 answers |
| 7 | Mon 28 Aug | Yusra's desk | Hollis, Monday 09:02: "FURROW-A 0827 1402 NEVER ARRIVED". dlinkd says `furrow-a-20280827-1402.tif` (booking `bk-5274`) landed at 14:13:05 on Sunday; Furrowlight's receiver says nothing. Season 3's webhook signature can't settle it: Furrowlight holds the same secret, so either side could have written it | Yusra: "TL;DR: receipts nobody but Apogee can sign and anyone can check." Receipts start on Friday 1 September; that day Furrowlight checks Sunday's receipt against the file already in its own storage: it matches. The file had arrived; their receiver missed the event |
| 8 | Mon 4 Sep | Ystmark, over the backhaul | Sunday 3 September a storm cuts Ystmark's fiber, and the station falls back on the slow satellite backhaul, a link another company runs. Sunniva: "Fiber cut. Backhaul up. Everything on it readable." | This week's command file for Ystmark, and next month Rotation Day's keys, must cross a link whose operator can read and change anything. On Thursday 7 September the file crosses sealed under a key the two ends agreed over the open link |
| 9 | Mon 11 Sep, to Fri 22 Sep | Mirela's desk; Salar Quieto; the safe | Mirela's audit: anything that reaches a station's network can talk to its uplinkd and its dish controllers. On Tuesday 12 September Paz finds an old test script on a spare PC still asking Chiquita's controller for its position every minute: harmless, and nobody knew | Every station and server gets a certificate from Apogee's own certificate authority. The root is made on Monday 18 September and locked in the safe; every certificate is issued by Friday 22 September |
| 10 | Mon 25 Sep, to Fri 6 Oct | Port Ellery and every station | Tomás shows that uplinkd's traffic to Lantern Atoll crosses another company's network as plain HTTP. Mirela: TLS 1.3 everywhere, both ends proven by Apogee's certificates | Rotation Day's keys will travel this path. By Friday 6 October every station speaks TLS 1.3 with both ends proven, and a sealed Cirrine-1 command crosses it to Ystmark untouched |
| 11 | Fri 6 Oct, to Sat 7 Oct | Yusra's desk; Brinewatch's maintenance window | Brinewatch can install Rotation Day's new API key and webhook secret only in its monthly window (the first Saturday of the month, 06:00 to 08:00 UTC) and only the way its older systems allow. Its first two secrets went out XORed with a 256-byte pad its courier brought in November 2027, both from the pad's first byte | Miss the window, and on Rotation Day Brinewatch is locked out of the booking service or keeps an old secret. The pad is retired; Brinewatch's 2021 RSA key (1,024 bits) is refused; its new 3,072-bit keys arrive by 17:30; at 06:00 on Saturday both secrets go out locked to its receiving key, and Brinewatch opens them |
| 12 | Mon 9 Oct; Rotation Day Tue 17 Oct; overlap ends Tue 24 Oct | the Fishbowl | Ingrid runs Rotation Day like a pass, from the Cirrine-1 pass that crosses midnight to Furrow-A over Salar Quieto at 14:02 | Every key in the network replaced, at every station, without missing a single pass. At 10:12 the Cape Wren contractor reports that the old station computer she took away for recycling still holds its disk, and with it the key of certificate 1005 |

The opening (Monday 10 July) hands over the new badge, shows Cirrine's three conditions and
Rotation Day's date, introduces Mirela (this season's reviewer), Yusra and Tomás, and sends the
learner to Mirela's question about random numbers. Words of the first week: key, nonce,
fingerprint, the ops console.

The finale (chapter 12, Tuesday 17 October): ten passes flown, none missed; the Cirrine-1
command sealed with `k-2028-07` at 23:59:58 reaches the satellite at 00:00:03 and is accepted;
every station on its new certificate by 01:15:40; certificate 1005 refused from 10:12; on 24
October `k-2028-07` retires and the second revocation list names every old certificate. Mirela
signs the log: "Every key replaced. Every pass flown. It failed closed every time we asked it
to." The debrief points at Season 6: in November the learner carries the pager again, and
Lantern Atoll's generator has started to sound different.

#### New systems and things

- **The ops console**, Apogee's internal web tool where engineers log in to stow a dish, run a
  pass by hand or read a station's log. It is older than the learner. Its user file holds one
  password record per account (chapter 4).
- **uplinkd's station side.** uplinkd runs at headquarters, where it signs, seals and schedules
  commands, and at every station, where it queues command lines for the next pass and hands
  them to the dish. In readmes: "Ystmark's uplinkd". No new service name.
- **Apogee's certificate authority**: a root certificate kept offline in the safe at Port
  Ellery and an intermediate on a machine in the Fishbowl (chapter 9).
- **The backhaul** (already in the bible: Ystmark's satellite link when the fiber is cut) is
  run by another company, unnamed; anything on it may be read or changed.
- **Receipts**: bookingd signs one for every delivered file; customers fetch it with
  `GET /bookings/{id}/receipt` (token, owner only; 404 `no receipt for bk-5274 yet` before the
  file lands) and check it with the public keys in the customer guide (chapter 7). No lesson
  runs this route; lessons sign and check envelopes as bytes.
- **Brinewatch's older systems** (chapter 11): its security office (unnamed; it speaks through Yusra) takes secrets only on its pad or locked to an RSA key, signs with RSA, and changes its systems only in its monthly maintenance window, the first Saturday of the month, 06:00 to 08:00 UTC. The pad, 256 random bytes on a USB stick its courier brought in November 2027 (one copy at each end), carried its API key that November and its webhook secret in January 2028, both XORed from the pad's first byte; it is retired on Friday 6 October 2028, and Mirela wipes Apogee's copy on the 7th. Its 2021 RSA key has 1,024 bits and is refused; its new receiving key (OAEP) and signing key (PSS), 3,072 bits each, are made that Friday, and Yusra reads the signing key's whole fingerprint back to Brinewatch by phone.

#### Every key, its size and its encoding

Hex is always lowercase. "Random" always means `crypto/rand`; nothing secret ever comes from
`math/rand` or `math/rand/v2`. Secret key files have mode 0600 and are written whole through a
temporary sibling and a rename (Season 2's way).

| Key | Made from | Size | Written as | Lives in |
|---|---|---|---|---|
| Uplink master `k-YYYY-MM` | random | 32 bytes | 64 hex digits in `uplink-keys.txt` | uplinkd at headquarters only |
| Satellite sign key | HKDF of a master | 32 bytes | never written; LOAD_KEY carries it sealed | the satellite; uplinkd derives it when needed |
| Satellite seal key | HKDF of a master | 32 bytes (AES-256) | as above | as above |
| `k-2027-03` (Season 1) | as made in 2027 | as made in 2027 | — | uplinkd; retired 2028-08-11T00:00:00Z; never in `uplink-keys.txt` |
| Command nonce | random, fresh for every seal | 12 bytes | 24 hex digits in the SEAL line | the line itself |
| Password salt | random, per record | 16 bytes | 32 hex digits in the record | the ops console's user file |
| Password key | PBKDF2-HMAC-SHA256, 600,000 iterations | 32 bytes | 64 hex digits in the record | the user file |
| Session token | `rand.Text()` | 26 characters `A`–`Z`, `2`–`7` (at least 128 bits) | the token; the server keeps only its SHA-256 in hex | the engineer's browser |
| Receipt key `r-YYYY-MM` | random 32-byte Ed25519 seed | public 32 bytes, signature 64 | seed: 64 hex in bookingd's key file; public: 64 hex in `receipt-keys.txt` and a PEM `PUBLIC KEY` block in the customer guide; signature: 128 hex | bookingd (seed); everyone (public) |
| Station identity (backhaul) | random Ed25519 seed | public 32 bytes | 64 hex | each station; the other end knows the public key in advance |
| Backhaul share | `ecdh.X25519().GenerateKey`, fresh per session | 32 bytes | 64 hex in HELLO | never stored |
| Backhaul session key | HKDF of the X25519 secret | 32 bytes | never written | memory, one session |
| CA and certificate keys | random Ed25519 seed | public 32 bytes | PEM: `CERTIFICATE` blocks; private keys as PKCS #8 `PRIVATE KEY` blocks, mode 0600 | root: offline in the safe; intermediate: the CA machine; each leaf: its own server |
| Brinewatch's pad (to 6 Oct 2028) | random | 256 bytes | on a USB stick | Brinewatch and Apogee, one copy each; Apogee's wiped 7 Oct |
| Brinewatch's RSA keys (receiving, signing) | Brinewatch's own, 6 Oct 2028 | 3,072 bits; ciphertexts and signatures 384 bytes | PEM `PUBLIC KEY` blocks (PKIX); fingerprint `sha256:` + 64 hex of the PKIX bytes | Brinewatch (private); Apogee (public) |
| bookingd token key, webhook secrets, API keys (Season 3) | random | as Season 3 | as Season 3 | replaced on Rotation Day; Brinewatch's new API key and webhook secret reach it on 7 October locked to its RSA receiving key (chapter 11); no lesson covers the others |

Ed25519 keys are always Ed25519 (`crypto/ed25519`). The only RSA keys are Brinewatch's (chapter 11), and no lesson generates one: the tests' four were made once with `rsa.GenerateKey` (Go 1.27.3) and are embedded as PKCS #8 PEM. ML-KEM-768
sizes (chapter 8): encapsulation key 1,184 bytes, ciphertext 1,088 bytes, shared key 32 bytes,
decapsulation key seed 64 bytes.

#### Test keys and fixed vectors

Lessons never contain flight keys. Test keys are readable on purpose. Symmetric test keys are
ASCII text of the right length; every test seed or private key of Ed25519 or X25519 is the
SHA-256 of a label (`sha256.Sum256([]byte(label))`). Every value below was computed with
Go 1.27.3.

| Test key or value | Exactly |
|---|---|
| Season 1's lab bench HMAC key | ASCII `lab bench key, never flown`; `pipit-2\|118\|SET_MODE\|SAFE` signs as `512545bcb43a5bb6ccb9f9c8511f4cc8c5fedeb38c226376af5b5aaae60e3b4d` |
| Master `k-2028-07` | ASCII `uplink test master key, 32 bytes` (hex `75706c696e6b2074657374206d6173746572206b65792c203332206279746573`) |
| Master `k-2028-10` | ASCII `uplink next master key, 32 bytes` (hex `75706c696e6b206e657874206d6173746572206b65792c203332206279746573`) |
| HKDF salt (production too: it is public) | ASCII `apogee uplink keys v1` (21 bytes) |
| Derived from `k-2028-07` | `sign\|pipit-2` `82f945638cbdd9a034dfab2bc578fd7a720d73e5701d3abcf561cd825847da7c`; `sign\|furrow-a` `00a1988ac17a102e03bf1e93099ecfeb710b18ec76e69dac6e30251a674a8dd3`; `sign\|cirrine-1` `18668bfe29024937ff84a61df34e8a0310a4702a74834b969f56e2c6ccd00213`; `seal\|cirrine-1` `427a8ae78a295915acb18a87bdbb40be20d49d27c15dac44e320d4510e070171`; `seal\|pipit-2` `53fdd2ccce4f186b0bdeb36626b6e8591ac1c5f59a00a80956c348dbdbaa2c3a` |
| Derived from `k-2028-10` | `sign\|pipit-2` `514e80107e64a13250a9bf9f8c5b1370995bea12a047f057a351c5c11f7522ad`; `sign\|cirrine-1` `3c280cbe585cfec33b8ac6693e6e0802fc0c532efc2af548cb1b3a85054418ba`; `seal\|cirrine-1` `a70fbff3ac2b05cbf0a66c3fc0eb76ceb1ec55269766703f715f4ec43c78663d` |
| Fixed nonce (vectors only; production nonces are always random) | hex `a1b2c3d4e5f60718293a4b5c` |
| Sealed `PING`, `cirrine-1`, seq 4471, `k-2028-07`, fixed nonce | `SEAL 2 cirrine-1 4471 k-2028-07 a1b2c3d4e5f60718293a4b5c c1f2a9798d38081d764bbcfb5e35e2598d2584a4` |
| Password fixture | password `fails closed, always`, salt ASCII `mirela test salt` (hex `6d6972656c6120746573742073616c74`); 600,000 iterations → key `6d91ceab779ff574b23820b133a740fa76d8b6a0deb5d6f0da2f94a6508dd938`; 1,000 iterations → `5b893367890cac1547d6c07d906a4b46b8a89733bf1f7fb3d7ba24f8e16174b3` |
| Old console record | `kiosk`: `sha256$1c4bb463c52d458ae6e9225430c45c753c5ae417bb4e7e4530773aaa0cda201b` (the password `fishbowl`) |
| Receipt keys | seed = SHA-256 of `receipts r-2028-08` → public `ba6c6226d2c3c056684e9d9a5e489db44a148e27e3fcb8f2ba16a6bddb0a86b2`; of `receipts r-2028-10` → public `3c3ae60d0f3d9f2f20ae8d68ee3b8996b233ed694c73b20f9091e3d14a494f12` |
| Station identities | seeds = SHA-256 of `identity PEL`, `identity YST`, `identity MID` (the one in the middle) → public `1497e6fc165f9bcd9fa8540b6a74089da8e02b669244eb2c9386fdd80588bcae`, `cacb8869a23066fe505e739fc60a7246f839d6a25e9994c1cd9a073de7cdd355`, `90aa0565e70d224291d92e937c4227ed71c6806cbda864e28bc7a6f9e8713b95` |
| Backhaul shares (vectors) | X25519 private keys = SHA-256 of `ephemeral PEL`, `ephemeral YST`, `ephemeral MID` → shares `a10c608ed294ad1d437b121ce63f146ee081f200b8b71126bc8f5642e69aa93e`, `b6f8a5a974e292828807a52f32334ac3925ea872ca82de318f53e971fafc5a66`, `1c85634a8ffb99861b81aa64cdc22c9f386010f5795f8212d89551372b360021`; PEL with YST agree on `c50d64c87ddd744d6c24842274ac98d74c33be2bdadd89463eae753cae91cc7a`; session key `24e57dc80a28df831da79d5b65636ba94ed13a4274908279576a15428a06958e` |
| CA and certificate keys | seeds = SHA-256 of `ca root 2028`, `ca stations 2028`, `cert <name> <serial>` (such as `cert yst.apogee.internal 1001`, `cert uplinkd 2001`); impostors: `ca someone else`, `cert impostor` |
| Digests | SHA-256 of `Apogee ground calling Pipit-1` `fe5e0f73525ffc9f862167a80cb444fbc55de2a916041901a86964737b2fbd66`, of `...Pipit-2` `8cbc29efabfcd2fb143763fd987194d89474355d64403f55e0def74fab81e72f`; SHA3-256 of `...Pipit-1` `7f7e6c3e2e3a3f1120a59350ba37cf6b382ad4644e840a4e15a35a179d9584ea`; SHA-256 of nothing `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| Brinewatch's RSA test keys (PEM in each RSA lesson's `keys_test.go`, byte for byte the same) | receiving, 3,072 bits, `sha256:d2499efc1976897aa14ea5d3948c4524bb826ad42b6ac1b43c5772f6b2ab8c1e`; signing, 3,072 bits, `sha256:e47dfcb6d8f2b0e68d9cc9ebfc169c6dbe786e4930e5f12e74238ff24d65f508`; the 2021 key, 1,024 bits, `sha256:6998de4f989149257401c3037508a903a2a573483f60b3eeb980c9abe5d2b2fe`; the impostor's, 2,048 bits, `sha256:66eb315a02a52f08680f0a365829b8e5108739e09f78673b27c6bfee0ea5e87e` (SHA-256 of each PKIX public key) |
| Brinewatch's test pad | 256 bytes: the SHA-256 of `brinewatch pad 0` to `brinewatch pad 7`, one after another |
| Brinewatch's new secrets (test values) | API key `brinewatch.20281017a1b2c3d4e5f60718293a4b5c`; webhook secret ASCII `brinewatch webhook secret 2028-10` (the old ones are Season 3's) |
| Short codes (lesson 11.2) | texts `brinewatch <n>`: 2 bytes, 48 and 307 share `bc8c` (308 tries); 3 bytes, 20 and 4774 share `f9f338` (4,775 tries); matching text 0's code: 2 bytes at text 96636, 3 bytes none in 100,000 |

Since Go 1.26, `ecdh` key generation ignores a reader passed to it, so a fixed "random"
reader never makes fixed keys: fixed keys come from `NewKeyFromSeed` (Ed25519) and
`ecdh.X25519().NewPrivateKey` (X25519). No test uses `testing/cryptotest.SetGlobalRandom` to
expect fixed output.

#### Random values (chapter 1)

- The helpers in `random.go`, used by every later chapter: `newKey(id string) Key` (32 random
  bytes), `newNonce() []byte` (12), `newSalt() []byte` (16), `newSessionToken() string`
  (`rand.Text()`). `crypto/rand.Read` never returns an error and always fills its slice; if the
  system can't supply randomness, the program stops. Code calls it without checking a result.
- A secret has at least 128 bits: length × log2(alphabet size), when every symbol is picked
  at random. Six digits are 19.9 bits;
  `rand.Text()` is 26 base32 characters, 130 bits.
- Random 12-byte nonces stay safe for 2^32 seals under one key; Apogee's busiest key seals a
  few thousand a year. A nonce counter kept in memory starts again after a restart, so Apogee
  never uses counters for nonces.

#### Digests (chapter 2)

- A **digest string** names its hash: `sha256:<64 hex>` or `sha3-256:<64 hex>`. Nothing else
  is accepted: no capitals, no other hash, no bare hex. `parseDigest` errors:
  `unknown hash "md5"`, `digest must be 64 lowercase hex digits`, `digest must be <hash>:<hex>`.
- Files are hashed as streams (`io.Copy` into the `hash.Hash`), never read whole: a Furrow-A
  image is 10.9 GB. Digests of public files are compared with `==` or `bytes.Equal`;
  constant-time comparison is for secrets.
- The **delivery manifest** (`manifest.txt` beside the files): one line per file,
  `<digest> <bytes> <file>`, single spaces; file names never contain spaces. The check reports
  one line per manifest line, in order: `ok <file>`, `wrong size <file>` (checked before
  hashing), `changed <file>`, `missing <file>`. A line without exactly three fields, a bad
  digest or a bad size is an error, `manifest line 3: want <digest> <bytes> <file>`.
- **The forged field report** (lessons 2.2 and 2.4): the original is the 20 bytes
  `NDVI field 12: 0.71\n` (hex `4e445649206669656c642031323a20302e37310a`, SHA-256
  `72198855a7d55dd3cbd5fff4d06168339fa8d32086a6cd18899400210b610aae`). Mirela's forgery is
  `NDVI field 12: 0.17\n` followed by the four bytes `1a 15 bc d1` (hex
  `4e445649206669656c642031323a20302e31370a1a15bcd1`, SHA-256
  `4b402892e69f9ddac861d7107795c4a44f68499ae7f4bc0d1ee1acea71093b8f`). Both have CRC-32
  `f6ced6c0`.

#### Signed commands, hardened (chapter 3)

Season 1's formats are unchanged: the canonical command `pipit-2|118|SET_MODE|SAFE`, the wire
line `CMD pipit-2 118 SET_MODE SAFE sig=<64 hex digits>`, HMAC-SHA256, sequence numbers
strictly increasing per satellite. Season 5 adds rules, never a new layout:

- **Canonical fields** (`canonical(c Command) (string, error)`): the satellite ID is 1 or more
  lowercase letters, digits and hyphens; `Seq` is at least 1; the operation is 1 or more
  capital letters and `_`;
  each argument is 1 or more printable ASCII characters from `!` to `~` other than `|`. Any
  other field is refused (`argument 1 contains "|"`), so no two commands share a canonical form.
- **Order of checks** for a CMD line: parse the line (exactly one `sig=` field, last, 64 hex
  digits); build the canonical form; verify the HMAC with `hmac.Equal`; only then compare
  `Seq` with the last accepted; only then record it. A refusal changes nothing.
- **The homemade signature** Mirela finds in Pipit-3's flight software is `sha256(key ||
  canonical)`. Lesson 3.1's read-only test forges it without the key by SHA-256 length
  extension: it rebuilds a hash state from a tag (the 4 bytes `sha\x03`, the 32 tag bytes, 64
  zero bytes, then the big-endian count of bytes hashed so far including padding) with
  `encoding.BinaryUnmarshaler` on `sha256.New()`, writes `|REBOOT`, and sends the command with
  SHA-256's padding and `|REBOOT` appended. HMAC refuses that forgery.
- **One HMAC per message.** A `hash.Hash` keeps everything written to it; signers start a new
  `hmac.New` for every message, which also makes them safe for many goroutines at once.
- **Comparing secrets.** MACs: `hmac.Equal` on bytes. Secrets of different lengths (a typed
  phrase): compare their SHA-256 digests with `subtle.ConstantTimeCompare`, because it returns
  at once when lengths differ.
- Errors, checked with `errors.Is`: `ErrForged` (`command failed authentication`, never more
  detail) and `ErrReplay` (`replayed sequence number: 118, last 118`), shared with sealed
  commands below.

#### The ops console's passwords (chapter 4)

- **Record**: `pbkdf2-sha256$<iterations>$<32 hex salt>$<64 hex key>`, from
  `pbkdf2.Key(sha256.New, password, salt, iterations, 32)` with a fresh 16-byte salt. New
  records use 600,000 iterations. A record with a count from 1,000 to 10,000,000 is checked
  with its own count; anything else, any other algorithm, or a salt or key of the wrong length
  doesn't parse. Mirela's: `pbkdf2-sha256$600000$6d6972656c6120746573742073616c74$6d91ceab779ff574b23820b133a740fa76d8b6a0deb5d6f0da2f94a6508dd938`.
- **Old records** (before 4 August): `sha256$<64 hex>`, the unsalted SHA-256 of the password.
  Lesson 4.1's intermediate step, `salted$<32 hex salt>$<64 hex of sha256(salt + password)>`,
  is lesson-only and never stored by the console.
- **The user file** `console-users.txt`: one line per account, `<user> <record>`; user names
  are lowercase, such as `mirela`, `kiosk`, `bench`.
- **Login rules**: derive again with the record's own count and salt and compare with
  `subtle.ConstantTimeCompare`; an unknown user and a wrong password give the same error,
  `invalid user or password` (`errInvalidLogin`), after the same work: one derivation against a
  fixed dummy record for an unknown user, and one for an old `sha256$` record too, which would
  otherwise answer at once and show which accounts still have one (`matches`, from 4.8). A
  record under 600,000 iterations still answers faster until it is upgraded. A successful login
  against an old record or one under 600,000 iterations saves a new record (upgrade at login);
  a session token is `rand.Text()`, and the console keeps only its SHA-256 in hex.
- **Cost in tests**: one derivation at 600,000 iterations takes about 60 ms on an Apple-silicon
  Mac. A lesson's tests run at most six of them; other cases use 1,000-iteration records.

```go
// Record is a parsed password record.
type Record struct {
	Iterations int    // 600000
	Salt       []byte // 16 bytes
	Key        []byte // 32 bytes
}

// Console is the ops console's login (lessons 4.7 and 4.8).
type Console struct {
	records  map[string]string                                    // user → record text
	kdf      func(pw string, salt []byte, iter int) ([]byte, error) // PBKDF2-HMAC-SHA256; tests count calls
	save     func(user, record string) error                      // stores an upgraded record
	sessions map[string]string                                    // hex SHA-256 of a token → user
}
```

The dummy record an unknown user is checked against is a fixed, valid 600,000-iteration record
whose password nobody knows; the check always fails after one `kdf` call, like a wrong password.

#### Keys from keys (chapter 5)

- A satellite key is `hkdf.Key(sha256.New, master.Secret, []byte("apogee uplink keys v1"),
  purpose+"|"+satID, 32)`, with purpose `sign` (its HMAC key for CMD lines) or `seal` (its
  AES-256-GCM key for SEAL lines). One master, many unrelated keys: a satellite holds only its
  own two keys for each master version it knows. HKDF is for keys that are already random; a
  password is never an HKDF input (that is PBKDF2's job).
- **`uplink-keys.txt`** (uplinkd only, mode 0600): `#` lines and blank lines are ignored; each
  other line is `<id> <activate> <retire or -> <64 hex secret>`, optionally followed by the word
  `revoked`; times are RFC 3339 UTC. A duplicate ID, a bad time or a secret that isn't 64 hex
  digits is an error naming the line. In August the file holds one line:
  `k-2028-07 2028-08-11T00:00:00Z - 75706c69...6573`.
- **The fleet** for derivations: `pipit-1`, `pipit-2`, `pipit-3`, `furrow-a`, `furrow-b`,
  `brine-3`, `cirrine-1` (Thistle re-entered in 2027).
- **Printing keys.** `Key` has `String()` returning `k-2028-07 (secret hidden)` and
  `GoString()` returning `Key{ID: "k-2028-07", Secret: hidden}`. fmt can't call methods on a
  value in an unexported field and prints its bytes, so `KeyRing` has its own `String()`:
  `key ring: k-2028-07, k-2028-10`, IDs in activation order, and a `GoString()` with the same
  text, since `%#v` calls GoString. The number verbs (`%d`, `%o`, `%b`, `%c`, `%U`) skip both
  methods and still print a key's bytes, so keys are logged with `%v` (lesson 5.5).

#### Sealed commands (chapter 6)

Uplink format 2 (format 1 is Season 1's signed CMD line). Every satellite's flight software
opens SEAL lines since August; Cirrine-1's commands are always sealed from 25 August; other
operators' routine commands may stay CMD lines; LOAD_KEY is always sealed.

```go
// Sealed is a command sealed for one satellite: only its header is readable.
type Sealed struct {
	Version    int    // 2
	SatID      string // "cirrine-1"
	Seq        uint64 // as Command.Seq: strictly increasing per satellite
	KeyID      string // "k-2028-07": the master version its key comes from
	Nonce      []byte // 12 bytes from crypto/rand, new for every seal
	Ciphertext []byte // AES-256-GCM of "OP|ARG...", its 16-byte tag included
}
```

- **The line**: `SEAL 2 <sat> <seq> <key id> <nonce: 24 hex> <ciphertext: hex>`, seven fields
  separated by single spaces. The ciphertext is 17 to 272 bytes (34 to 544 hex digits): a
  plaintext of 1 to 256 bytes plus the tag. Readmes show it shortened: `SEAL 2 cirrine-1 4471
  k-2028-07 a1b2…4b5c c1f2…84a4`.
- **The key**: the satellite's `seal` key from the named master (chapter 5).
- **Associated data**: the header as text, `SEAL|2|cirrine-1|4471|k-2028-07`, authenticated
  but not hidden, so the station can schedule the command without reading it, and nobody can
  move a sealed command to another satellite, number or key.
- **Plaintext**: the operation and arguments joined by `|`, checked by chapter 3's canonical
  rules: `PING`, `SET_MODE|SAFE`. After opening, the full canonical command is
  `cirrine-1|4471|PING`.
- **Order of checks** (the satellite and its stand-in in tests): parse (`ErrMalformed`, also
  for a line that doesn't start with `SEAL`); version 2 (`ErrVersion`, only for a SEAL line);
  satellite is its own (`ErrWrongSat`); key ID held (`ErrUnknownKey`) and not revoked
  (`ErrRevoked`) or retired (`ErrRetired`); `Seq` above the last accepted (`ErrReplay`); `Open`
  with the associated data (`ErrForged`); only then record `Seq`. A
  refusal changes nothing, so a forged line numbered 18,446,744,073,709,551,615 can't lock the
  satellite out. (Chapter 3 verifies before it compares `Seq`; here a number already seen is
  refused before any work is spent opening it, and 6.7's readme says so.)
- **Errors**, all `errors.New` sentinels checked with `errors.Is`, wrapped with context:

  ```go
  var (
  	ErrMalformed = errors.New("malformed command line") // "malformed command line: want 7 fields"
  	ErrVersion   = errors.New("unsupported version")    // "unsupported version: 3"
  	ErrWrongSat  = errors.New("not for this satellite") // `not for this satellite: "pipit-2"`
  	// ErrUnknownKey means the ring holds no version with the key ID asked for.
  	ErrUnknownKey = errors.New("unknown key")
  	ErrRevoked    = errors.New("key revoked") // `key revoked: "k-2028-10"`
  	ErrRetired    = errors.New("key retired") // `key retired: "k-2028-07"`
  	// ErrReplay refuses a command whose sequence number isn't above the last one accepted.
  	ErrReplay = errors.New("replayed sequence number")
  	// ErrForged refuses a command whose signature doesn't check out, and never says more.
  	ErrForged = errors.New("command failed authentication")
  )
  ```

  `ErrUnknownKey`, `ErrReplay` and `ErrForged` keep the text of their first declarations
  (chapters 5 and 3). Wrapped, they read `unknown key: "k-2027-03"` and
  `replayed sequence number: 4471, last 4471`.

- **LOAD_KEY** (chapter 12): the plaintext `LOAD_KEY|<new key id>|<64 hex sign key>|<64 hex
  seal key>` (148 bytes for `k-2028-10`), sealed under the satellite's current seal key. After
  opening it, the satellite accepts the new version at once; it signs nothing, so "activate"
  doesn't apply to it.
- The station never holds a satellite's seal key: Ystmark's uplinkd queues Cirrine-1's SEAL
  lines without being able to read them.

#### Receipts (chapter 7)

```go
// Receipt is Apogee's signed word that a booking's file was delivered.
type Receipt struct {
	Version   int       `json:"v"`          // 1
	KeyID     string    `json:"key_id"`     // "r-2028-08": the key that signed it
	BookingID string    `json:"booking_id"` // "bk-5274"
	SatID     string    `json:"sat_id"`
	Station   string    `json:"station"`
	File      string    `json:"file"`       // "furrow-a-20280827-1402.tif"
	Bytes     int64     `json:"bytes"`
	Digest    string    `json:"digest"`     // "sha256:<64 hex>", or "sha3-256:" for Cirrine
	LandedAt  time.Time `json:"landed_at"`  // from the Landing (Season 3)
	SignedAt  time.Time `json:"signed_at"`  // bookingd's clock, truncated to the second
}

// Envelope carries a receipt and its signature.
type Envelope struct {
	KeyID   string          `json:"key_id"`  // the same as the receipt's
	Receipt json.RawMessage `json:"receipt"` // the exact bytes that were signed
	Sig     string          `json:"sig"`     // 128 hex digits: Ed25519 over Receipt's bytes
}

// ReceiptKey is one published receipt key and its signing window.
type ReceiptKey struct {
	ID       string            // "r-2028-08"
	Public   ed25519.PublicKey // 32 bytes
	Activate time.Time         // receipts signed from here...
	Retire   time.Time         // ...until here (zero: still signing)
	Revoked  bool
}

// Notary is bookingd's receipt signer.
type Notary struct {
	KeyID string
	priv  ed25519.PrivateKey
}
```

- **Signing**: `json.Marshal(receipt)` (no newline, no indentation) is signed with
  `ed25519.Sign`; the envelope is `json.Marshal(Envelope{...})` with those exact bytes as
  `receipt`. Example receipt bytes:
  `{"v":1,"key_id":"r-2028-08","booking_id":"bk-5274","sat_id":"furrow-a","station":"SQT","file":"furrow-a-20280827-1402.tif","bytes":10874220000,"digest":"sha256:<64 hex>","landed_at":"2028-08-27T14:13:05Z","signed_at":"2028-09-01T09:00:00Z"}`.
  The archived file is `<file>.receipt.json`, the envelope as sent.
- **Checking**, in this order: decode the envelope strictly (unknown fields refused; 64 KiB at
  most) else `ErrBadEnvelope`; find the key by the envelope's `key_id` in your own published list,
  never a key the envelope brings, else `ErrUnknownSigner`; refuse a revoked key
  (`ErrRevoked`); `sig` must be 128 hex digits and verify over the raw `receipt` bytes, else
  `ErrBadSignature`; only then decode the receipt strictly; its `key_id` must equal the
  envelope's (`ErrBadEnvelope`); `signed_at` must fall in `[Activate, Retire)` of that key,
  else `ErrOutsideWindow`. Checking a file: the same size, and the digest of the file, hashed
  with the receipt's own algorithm, equal to the receipt's, else `ErrWrongFile`.
- `ed25519.Verify` panics on a public key that isn't 32 bytes: check the length first.
- **`receipt-keys.txt`** (published in the customer guide, beside each key's PEM block):
  `<id> <activate> <retire or -> <64 hex public key>`, `revoked` appended when revoked:

  ```
  r-2028-08 2028-09-01T00:00:00Z 2028-10-17T00:00:00Z ba6c6226...86b2
  r-2028-10 2028-10-17T00:00:00Z - 3c3ae60d...4f12
  ```

- Errors: `ErrBadEnvelope` (`malformed receipt`), `ErrUnknownSigner` (`unknown receipt key`),
  `ErrBadSignature` (`receipt signature does not verify`), `ErrOutsideWindow` (`receipt signed
  outside its key's window`), `ErrWrongFile` (`file does not match receipt`), and `ErrRevoked`
  from the sealed commands.
- `newNotary(keyID string, seed []byte) (*Notary, error)` refuses a seed that isn't 32 bytes
  (`ed25519.NewKeyFromSeed` would panic).

#### The backhaul handshake (chapter 8)

Port Ellery headquarters (`PEL`) dials; the station (`YST`) answers. Each side holds its own
Ed25519 identity (`Identity{Code string; Key ed25519.PrivateKey}`) and knows the other's
public identity key in advance (`map[string]ed25519.PublicKey`, by station code).

1. Each side makes a fresh X25519 key for this session: `ecdh.X25519().GenerateKey(rand.Reader)`.
2. Each sends one line, PEL first: `HELLO <code> <64 hex share> <128 hex signature>`. The
   signature is the side's identity key over the bytes `backhaul share|<code>|<64 hex share>`.
   A line is checked before its share is used: exactly four fields, a known code that isn't
   your own, 64 lowercase hex digits, a signature that verifies, else `ErrBadHello` (`bad hello`,
   wrapped: `bad hello: unknown station "MID"`, `bad hello: signature does not verify`).
3. Both compute the X25519 secret (an all-zero result is an error, never a key) and derive
   `hkdf.Key(sha256.New, secret, nil, "apogee backhaul v1|<PEL share hex>|<YST share hex>", 32)`:
   the dialer's share first.
4. Data lines: `DATA <24 hex nonce> <hex ciphertext>`, AES-256-GCM under the session key with
   associated data `DATA|<from code>|<to code>`, a fresh random nonce each. A line that doesn't
   open ends the session with an error. DATA lines carry no sequence number, so a machine in
   the middle can drop, repeat or reorder them within a session, though never read or change
   one; the CMD lines inside carry their own numbers. The HKDF info binds the two shares, not
   the station codes, which is sound because each share's signature covers its code.
5. Nothing of the session is kept afterwards: a stolen station later can't open a recording.

Tests run both sides over `net.Pipe()` with a deadline of 1 s; the side in the middle is the
test, holding the `identity MID` key.

#### Apogee's certificate authority (chapter 9)

All keys Ed25519; all times UTC; serial numbers fixed in tests (in production, random 128-bit
from `crypto/rand`, with `rand.Int`).

| Certificate | Subject (CN, O=Apogee) | Issuer | Serial | Valid | Shape |
|---|---|---|---|---|---|
| Root | `Apogee Root CA 2028` | itself | 1 | 2028-09-18T00:00:00Z to 2038-09-18T00:00:00Z (10 years) | CA, `MaxPathLen` 1, key usage cert sign and CRL sign |
| Intermediate | `Apogee Stations CA 2028` | the root | 2 | 2028-09-18T00:00:00Z to 2030-09-18T00:00:00Z (2 years) | CA, `MaxPathLen` 0 with `MaxPathLenZero`, cert sign and CRL sign |
| Stations | `yst.apogee.internal`, `pel.`, `sqt.`, `lan.`, `cwr.` | the intermediate | 1001 YST, 1002 PEL, 1003 SQT, 1004 LAN, 1005 CWR | 2028-09-22T00:00:00Z for 90 days, to 2028-12-21T00:00:00Z | DNS name = CN; server auth; digital signature |
| uplinkd (headquarters) | `uplinkd`, DNS `uplinkd.apogee.internal` | the intermediate | 2001 | as stations | client auth; digital signature |
| passplan | `passplan`, DNS `passplan.apogee.internal` | the intermediate | 2002 | as stations | client auth (it may read status, never queue) |
| dishd (Port Ellery) | `dishd`, DNS `dishd.apogee.internal` | the intermediate | 2003 | as stations | client auth (it dials Ystmark's dish controllers) |
| Rotation Day's | the same names | the intermediate | 1101–1105, 2101–2103 | 2028-10-09T00:00:00Z to 2029-01-07T00:00:00Z | as above |

- Host names end in `.apogee.internal`; `.internal` is reserved for private networks, so no
  name can belong to anyone outside.
- A certificate is valid from `NotBefore` to `NotAfter`, **both instants included** (crypto/x509
  refuses only a time before `NotBefore` or after `NotAfter`).
- Files: `<name>.pem` holds the leaf then the intermediate (`CERTIFICATE` blocks);
  `<name>.key` the PKCS #8 private key (`PRIVATE KEY`), mode 0600; `roots.pem` the root.
  Names: `yst`, `pel`, `sqt`, `lan`, `cwr`, `uplinkd`, `passplan`, `dishd`. `tls.X509KeyPair`
  loads a pair.
- `Verify` always gets `Roots` (a pool with the root only), `Intermediates` (the chain the peer
  sent), `DNSName` where a name is checked, `KeyUsages` (server or client auth) and
  `CurrentTime` from the injected clock. Never leave `Roots` nil: on a Mac that hands the
  check to the system verifier, which doesn't know Apogee's root and refuses Ed25519
  certificates.
- The impostor in tests is a root also named `Apogee Root CA 2028` (key `ca someone else`)
  issuing `yst.apogee.internal`: same names, wrong key.
- `Issuer{Cert *x509.Certificate; Key ed25519.PrivateKey}` is a CA's certificate and key.

#### Every connection encrypted (chapter 10)

- **Both ends**: TLS 1.3 only (`MinVersion: tls.VersionTLS13`); the server presents its leaf
  and the intermediate; the client sets `RootCAs` (Apogee's root) and `ServerName`
  (`yst.apogee.internal`), since tests dial 127.0.0.1; the server sets `ClientAuth:
  tls.RequireAndVerifyClientCert` and `ClientCAs` (Apogee's root); both set `Time` to the
  injected clock. Go's default key exchange is the hybrid `X25519MLKEM768`; nobody sets
  `CurvePreferences`. `InsecureSkipVerify` never appears outside lesson 10.1's starter.
- **Who may do what**: the caller is the common name of `r.TLS.VerifiedChains[0][0]` (never
  `PeerCertificates`, which is only what the client sent). Any verified Apogee client may read
  status; only `uplinkd` may queue commands. A request without `r.TLS` is refused, never a panic.
- **Ystmark's uplinkd over HTTPS** (the station side; all answers `text/plain`):

  | Route | Answer |
  |---|---|
  | `GET /status` | 200 `ok YST` and a newline |
  | `POST /uplink` | body: one CMD or SEAL line, at most 1024 bytes, syntax checked only; 202 `queued cirrine-1 4471`; 403 `only uplinkd may queue commands`; 400 `bad command line` |

- **The dish controllers** speak their own raw socket protocol (Season 4), and dishd at Port
  Ellery is still the only program that talks to Ystmark's. This season wraps its connection to
  a controller in TLS (`tls.Dialer`, or `tls.Client` with `HandshakeContext`): dishd proves
  itself with its client certificate, 2003, and a controller presents its station's
  certificate, 1001 for Ystmark. In tests the controller is a stand-in on
  `net.Listen("tcp", "127.0.0.1:0")` behind `tls.Server` that echoes each line it reads,
  `PING yst` in, `PING yst` out; it never claims to be the real protocol.

#### Brinewatch's older keys (chapter 11)

- **The pad:** `Pad{key []byte; used int}`; `(*Pad).Seal` XORs with the next unused bytes and refuses with `pad has 7 bytes left, need 25`, spending nothing. A refusal changes nothing.
- **Short codes:** `shortCode(text, n)` is the hex of the first n bytes of SHA-256; Brinewatch's tool shows 3 bytes (six hex digits, 16.7 million codes).
- **`readRSA(pemText []byte) (*rsa.PublicKey, error)`:** one `PUBLIC KEY` block with only white space after it (`want one PUBLIC KEY block`), then `x509.ParsePKIXPublicKey` and a type switch. Errors: `want at least 2048 bits, got 1024`, `want an RSA key, got Ed25519`, `want an RSA key, got X25519`, `want an RSA key, got %T`.
- **OAEP:** `encryptTo` and `decryptWith` use `sha256.New()` and labels `apogee api key` and `apogee webhook secret`; a 3,072-bit key gives 384-byte ciphertexts and locks at most 318 bytes.
- **PSS:** `signRequest` and `checkRequest` work on the SHA-256 digest with `crypto.SHA256` and `nil` options; Brinewatch's tool uses `PSSSaltLengthEqualsHash`, which `nil` verification accepts.
- **The request:** `Request{Text, Sig, Key []byte}`, with Text `brinewatch handover sha256:<64 hex>` (`keyFingerprint`: `sha256:` + SHA-256 of `x509.MarshalPKIXPublicKey`).
- **The handover:** check the signature, read the key, compare the line exactly, then two lines `api-key <hex>` and `webhook-secret <hex>` (768 hex digits each). Any failure gives an error and no lines.

#### Rotation Day (chapter 12)

```go
// KeyVersion is one version of a key and the times it may be used.
type KeyVersion struct {
	Key                // Season 1's Key: ID and Secret
	Activate time.Time // signers and sealers use it from this instant
	Retire   time.Time // nothing it made is accepted from this instant; zero: not scheduled
	Revoked  bool      // refused at once, whatever the times; never undone
}
```

**The rules**, the same for every kind of key:

1. Verifiers accept a version from the moment it is installed, before its `Activate`, so
   nothing sealed after the switch arrives at a verifier that doesn't know the key yet.
2. A signer uses the newest version (by `Activate`) that has activated, isn't revoked and
   hasn't retired. With none, it refuses to sign: `no key to sign with at <RFC 3339 time>`.
3. Windows are half open: a version signs from `Activate` included and is refused from
   `Retire` included.
4. **Overlap**: an uplink master's `Retire` is its successor's `Activate` plus 7 days, so a
   command sealed with the old key just before the switch is still accepted after it.
   Receipt keys don't overlap in signing (the old one's `Retire` is the new one's `Activate`),
   and stay published for as long as their receipts are kept: a receipt is judged by its
   `signed_at`, a command by when it arrives.
5. **Revocation** wins over every window, at once, and is never undone; a revoked receipt key
   makes all its receipts unacceptable. Certificates are revoked through the intermediate's
   revocation list.
6. A CMD line names no key: during an overlap a satellite tries its keys newest first. SEAL
   lines name theirs.

**Revocation lists** (`x509.CreateRevocationList`, signed by `Apogee Stations CA 2028`):
`Number` counts from 1; `NextUpdate` is `ThisUpdate` plus 7 days; a list is current from
`ThisUpdate` to `NextUpdate`, both included. A checker refuses a list whose signature doesn't
verify against the intermediate (`CheckSignatureFrom`) or that isn't current, and then refuses
every certificate rather than trusting it (fail closed). The TLS check runs in
`VerifyConnection` (it also runs on resumed sessions) and answers
`certificate 1005 is revoked`.

**The timeline** (UTC):

| When | What |
|---|---|
| Fri 2028-07-14 16:00 | `k-2028-07` made (Activate 2028-08-11T00:00:00Z) |
| Thu 2028-08-10 | every satellite's own `k-2028-07` keys go up inside its operator's August flight-software update, which also teaches it to open SEAL lines |
| Fri 2028-08-11 00:00 | `k-2027-03` retires; satellites answer only to their own keys |
| Mon 2028-08-28 | `r-2028-08` made; published; Activate 2028-09-01T00:00:00Z |
| Fri 2028-10-06 | Brinewatch's pad retired; its new RSA keys made, and its new API key and webhook secret |
| Sat 2028-10-07 06:00 | Brinewatch's maintenance window: both new secrets handed over locked to its receiving key; they take effect on Rotation Day |
| Mon 2028-10-09 | `k-2028-10` made (Activate 2028-10-17T00:00:00Z); `k-2028-07` gets Retire 2028-10-24T00:00:00Z; `r-2028-10` published (Activate 2028-10-17T00:00:00Z) and `r-2028-08` gets Retire 2028-10-17T00:00:00Z; certificates 1101–1105 and 2101–2103 issued |
| Mon 9 to Fri 13 Oct | each satellite gets LOAD_KEY for `k-2028-10` on its first pass, sealed under its `k-2028-07` seal key |
| Tue 2028-10-17 00:00:00 | Rotation Day: uplinkd seals and signs with `k-2028-10`; bookingd signs receipts with `r-2028-10` |
| 00:30:00 | uplinkd's client certificate 2101 replaces 2001, and dishd's 2103 replaces 2003 (stations and controllers accept both until 24 October) |
| 01:00:00 to 01:15:40 | each station switches to its new certificate in its first quiet gap (below) |
| 10:12:00 | revocation list 1: certificate 1005 (Cape Wren's old station computer left with its disk) |
| Tue 2028-10-24 00:00:00 | `k-2028-07` retires; revocation list 2: 1001 to 1005 and 2001 to 2003 |

**Rotation Day's passes** (Season 1's `Pass`; lessons 12.7 and 12.8). Each pass gets one sealed
`PING`; the Cirrine-1 pass across midnight gets two, sealed at 23:59:58 and 00:01:00. A command
reaches its satellite 5 seconds after it is sealed.

| Satellite | Station | AOS | LOS |
|---|---|---|---|
| cirrine-1 | YST | 2028-10-16T23:52:10Z | 2028-10-17T00:03:40Z |
| pipit-3 | PEL | 00:58:30 | 01:07:10 |
| brine-3 | LAN | 01:02:00 | 01:10:40 |
| furrow-b | CWR | 02:51:20 | 03:00:44 |
| pipit-1 | PEL | 04:12:05 | 04:22:41 |
| pipit-2 | YST | 04:31:10 | 04:41:52 |
| furrow-a | CWR | 04:58:30 | 05:08:02 |
| cirrine-1 | YST | 07:30:00 | 07:41:12 |
| brine-3 | LAN | 09:14:20 | 09:22:31 |
| furrow-a | SQT | 14:02:10 | 14:11:52 |

**The switch rule** (`switchTime`): switches start from 01:00:00; a switch takes 10 minutes and
must keep 5 minutes clear of every pass at that station on each side, so it can't overlap
`[AOS − 5 min, LOS + 5 min)`. The first such start is `from` itself or 5 minutes after some
pass's LOS. Results: YST 01:00:00, SQT 01:00:00, CWR 01:00:00, PEL 01:12:10, LAN 01:15:40. No
fit within 24 hours of `from` is an error: `no 10-minute gap for CWR before
2028-10-18T01:00:00Z`.

**The finale's report** (`Report{Flown, Missed int; Lines []string}`), one line per pass in AOS
order: its AOS as `15:04:05`, the satellite, the station, `flown`, the key IDs its commands
were sealed with (joined by `+` when there are two) and the serial of the station certificate
in force at AOS: `04:12:05 pipit-1 PEL flown k-2028-10 1102`,
`23:52:10 cirrine-1 YST flown k-2028-07+k-2028-10 1001`. A miss names its first refusal:
`09:14:20 brine-3 LAN missed: certificate 1004 is revoked`. The real day: 10 flown, 0 missed.

#### The code, by chapter

Each lesson is its own `package main` (module `lesson`). A lesson gives the earlier pieces it
needs as read-only files under these names, and uses only the fields it needs. `Key`,
`Command`, `Pass` and Season 3's `Booking` and `Landing` are unchanged.

| Chapter | Names it adds (and files) |
|---|---|
| 1 | `newKey`, `newNonce`, `newSalt`, `newSessionToken` (`random.go`); `bits` and `Noncer` are lesson-only |
| 2 | `parseDigest`, `newHash`, `fingerprintAs` (`digest.go`), `checkManifest` (`manifest.go`); `fingerprint`, `fingerprintStream` and `save` are lesson-only |
| 3 | `Command` (`command.go`), `tag`, `verify`, `canonical` and its helpers `checkField`, `satIDRune`, `opRune`, `argRune` (`mac.go`), `Guard`, `(*Guard).Accept`, `ErrForged`, `ErrReplay` (`guard.go`), `checkLine` (`line.go`); `samePhrase`, `Signer` and `NewSigner` stay in their lessons' `main.go` |
| 4 | `Record`, `derive`, `formatRecord`, `parseRecord`, `checkPassword`, `newRecord`, `login`, `Console`, `(*Console).Check`, `(*Console).Login`, `(*Console).matches`, `(*Console).newRecord`, `needsUpgrade`, `dummyRecord`, `errInvalidLogin` (`password.go`, `console.go`); `checkStored` stays in 4.6's `main.go` |
| 5 | `signKey`, `satKey`, `fleetKeys`, `KeyVersion`, `KeyRing`, `loadRing`, `(*KeyRing).SatKey`, `(*KeyRing).Versions`, `(*KeyRing).String`, `(*KeyRing).GoString`, `Key.String`, `Key.GoString`, `parseVersion`, `(*KeyRing).add`, `ErrUnknownKey` (`derive.go`, `ring.go`; `%#v` calls GoString, so both types need one) |
| 6 | `seal`, `open`, `sealCmd`, `Sealed`, `(Sealed).AD`, `formatSealed`, `parseSealed`, `Satellite`, `SatKeys`, `(*Satellite).Accept`, `sealCommand`, the `Err…` sentinels (`seal.go`, `satellite.go`) |
| 7 | `Notary`, `newNotary`, `(*Notary).Sign`, `(*Notary).Published`, `checkReceipt`, `publishKey`, `readKey`, `Receipt`, `Envelope`, `ReceiptKey`, `wrapReceipt`, `openEnvelope`, `verifyEnvelope`, `issue`, `checkDelivery` (`receipt.go`) |
| 8 | `agree`, `parseShare`, `agreeHex`, `Identity`, `hello`, `checkHello`, `ErrBadHello`, `sessionKey`, `Backhaul`, `(*Backhaul).NewSession`, `wrapFor`, `unwrap`, `dial`, `answer` (`agree.go`, `backhaul.go`). Never a type named `Station`: that name is the bible's ground station |
| 9 | `trusted`, `Issuer`, `newRoot`, `newIntermediate`, `issueServer`, `issueClient`, `verifyStation`, `checkAt`, `issueAll` (`ca.go`) |
| 10 | `newClient`, `serverConfig`, `harden`, `requireClients`, `clientConfig`, `caller`, `Queue`, `stationHandler`, `dialController`, `uplinkClient`, `stationServer`, `send` (`tlsconf.go`, `station.go`) |
| 11 | `Pad`, `(*Pad).Seal`, `shortCode`, `firstMatch` (lesson-only, `main.go`); `readRSA`, `encryptTo`, `decryptWith`, `signRequest`, `checkRequest`, `keyFingerprint`, `Request`, `apiKeyLabel`, `webhookLabel` (`brinewatch.go`); `handOver` (`main.go`) |
| 12 | `accepts`, `signingKey`, `loadKeyCommand`, `revoke`, `isRevoked`, `withRevocations`, `switchTime` (`rotate.go`), `(*Satellite).Receive` and `load` (`satellite.go`), `Slot`, `Day`, `(Day).Passes`, `(Day).Connect`, `(Day).Ping`, `handshake`, `Report`, `runDay` (`day.go`, `main.go`) |

Signatures every chapter after the one that defines them uses exactly:

```go
func newKey(id string) Key
func newNonce() []byte
func newSalt() []byte
func newSessionToken() string
func parseDigest(s string) (alg string, sum []byte, err error)
func newHash(alg string) (hash.Hash, error)
func fingerprintAs(r io.Reader, alg string) (string, error) // "sha256:<hex>"
func canonical(c Command) (string, error)
func verify(canonical, sigHex string, key []byte) bool // Season 1's, unchanged
func (g *Guard) Accept(c Command, key []byte) error
func checkLine(line string, key []byte, g *Guard) (Command, error)
func derive(pw string, salt []byte) ([]byte, error)
func parseRecord(s string) (Record, error)
func formatRecord(r Record) string
func checkPassword(r Record, pw string) (bool, error)
func signKey(master Key, satID string) ([]byte, error)          // 5.1: the "sign" key only
func satKey(master Key, purpose, satID string) ([]byte, error) // from 5.4: purpose "sign" or "seal"
func loadRing(r io.Reader) (*KeyRing, error)
func (r *KeyRing) SatKey(keyID, purpose, satID string) ([]byte, error)
func (r *KeyRing) Versions() []KeyVersion // by Activate
func (s Sealed) AD() []byte
func formatSealed(s Sealed) string
func parseSealed(line string) (Sealed, error)
func (s *Satellite) Accept(line string) (Command, error)
func sealCommand(ring *KeyRing, keyID string, c Command) (string, error)
func newNotary(keyID string, seed []byte) (*Notary, error)
func (n *Notary) Sign(msg []byte) []byte
func (n *Notary) Published() []byte // what customers are given
func wrapReceipt(n *Notary, r Receipt) ([]byte, error)
func verifyEnvelope(keys map[string]ed25519.PublicKey, env []byte) (Receipt, error)
func checkDelivery(keys []ReceiptKey, env []byte, file io.Reader) (Receipt, error)
func hello(id Identity, share *ecdh.PrivateKey) string
func checkHello(line string, known map[string]ed25519.PublicKey) (code string, share *ecdh.PublicKey, err error)
func sessionKey(shared, dialerShare, answererShare []byte) ([]byte, error)
func dial(conn io.ReadWriter, me Identity, peer string, known map[string]ed25519.PublicKey) ([]byte, error)
func answer(conn io.ReadWriter, me Identity, known map[string]ed25519.PublicKey) (peer string, key []byte, err error)
func newRoot(key ed25519.PrivateKey, now time.Time) (*x509.Certificate, error)
func newIntermediate(root Issuer, key ed25519.PrivateKey, now time.Time) (*x509.Certificate, error)
func issueServer(ca Issuer, host string, serial int64, pub ed25519.PublicKey, now time.Time) (*x509.Certificate, error)
func issueClient(ca Issuer, name string, serial int64, pub ed25519.PublicKey, now time.Time) (*x509.Certificate, error)
func verifyStation(leaf *x509.Certificate, chain []*x509.Certificate, roots *x509.CertPool, host string, now time.Time) error
func serverConfig(leaf, inter *x509.Certificate, key ed25519.PrivateKey, now func() time.Time) *tls.Config
func clientConfig(roots *x509.CertPool, host string, cert tls.Certificate, now func() time.Time) *tls.Config
func caller(r *http.Request) (string, bool)
func stationHandler(code string, q *Queue) http.Handler
func dialController(ctx context.Context, addr string, cfg *tls.Config) (net.Conn, error)
func send(ctx context.Context, c *http.Client, url, line string) error
func readRSA(pemText []byte) (*rsa.PublicKey, error)
func encryptTo(pub *rsa.PublicKey, secret []byte, label string) ([]byte, error)
func checkRequest(pub *rsa.PublicKey, req, sig []byte) error
func keyFingerprint(pub *rsa.PublicKey) (string, error) // "sha256:<hex>"
func accepts(ring []KeyVersion, id string, now time.Time) bool
func signingKey(ring []KeyVersion, now time.Time) (KeyVersion, error)
func loadKeyCommand(ring *KeyRing, satID string, seq uint64, next KeyVersion) (string, error)
func revoke(ca Issuer, number int64, entries []x509.RevocationListEntry, now time.Time) ([]byte, error)
func isRevoked(crlDER []byte, ca, cert *x509.Certificate, now time.Time) (bool, error)
func withRevocations(cfg *tls.Config, current func() []byte, ca *x509.Certificate, now func() time.Time)
func switchTime(passes []Pass, station string, from time.Time) (time.Time, error)
```

`Satellite` (the flight software's stand-in) holds `ID`, `Keys map[string]SatKeys` (by key ID;
`SatKeys{Sign, Seal []byte; Retire time.Time; Revoked bool}`), `Last uint64` and
`Now func() time.Time`; from chapter 12 its `Receive` checks a line as `Accept` does, then
installs the keys a LOAD_KEY carries (`load`). `Guard` holds `last map[string]uint64` (a zero
`Guard` works). `Queue` holds a mutex and the queued lines.

#### How lessons test it

- **Never on random output.** Tests check properties (sizes, alphabets, two draws differing,
  round trips, a changed byte refused) or fixed vectors from the table above, computed in the
  test with the standard library where the algorithm is deterministic (SHA-2, SHA-3, HMAC,
  PBKDF2, HKDF, AES-GCM with a fixed nonce, Ed25519, X25519 from fixed private keys, ML-KEM
  round trips). A property that random output could break by chance (two 12-byte nonces equal)
  has odds below 2^-90. RSA-OAEP and RSA-PSS are randomized, so tests open or verify their output and never compare its bytes.
- **The attacker is the test.** One lesson in every chapter shows the mistake the chapter
  guards against, in Mirela's voice, before the right way: lessons 1.1, 2.2, 3.1, 4.1, 5.1,
  6.1, 7.1, 8.4 (after two lessons of basics), 9.1, 10.1, 11.1 and 12.1. The starter holds the
  mistake, and the read-only test plays the person who exploits it (predicts the nonce, forges
  the CRC, extends the hash, cracks the table, signs for another satellite, flips the
  counter-mode bytes, forges the shared-secret receipt, sits in the middle, presents a
  same-named certificate, impersonates a server, reads one pad-sealed secret from another, or arrives a second after midnight). The
  solution makes the attack fail.
- **Reading the code.** Where right and wrong code give the same answers (constant-time
  comparison, no `math/rand`, extract once), the test parses the learner's file with
  `go/parser`, as Season 1's lesson 14.10 did, and says what it looked for.
- **PBKDF2 budget**: at most six 600,000-iteration derivations per lesson run.
- **Servers and connections**: `httptest.NewUnstartedServer` with `TLS` set, then
  `StartTLS()`, or `tls.Server` on `net.Listen("tcp", "127.0.0.1:0")`; never a fixed port. A
  server's `Config.ErrorLog` is `log.New(io.Discard, "", 0)`, so refused handshakes print
  nothing. Clients have `RootCAs`, `ServerName` and a timeout of at most 2 s; clocks are
  injected through `tls.Config.Time` and `VerifyOptions.CurrentTime`; every wait has a guard
  of about a second; sleeps are 50 ms at most.
- **Files** live under `t.TempDir()`; filesystems for manifests are `fstest.MapFS`.
- Messages read `call = got, want want`, with hex shortened to its first 8 digits and `…` when
  the full value is longer than 50 characters.
- Readme examples and outputs stay within 50 characters a line: hex is shown as its first and
  last four digits around `…`, and PEM blocks (60-character lines) are never shown whole.

**Topics added this season** (front matter): `crypto`. Lessons keep using `bytes`, `errors`,
`files`, `http`, `json`, `security`, `streams`, `strings`, `testing`, `time` and the rest.

**New field guide words** (world words, in plain language): Ground Security Engineer, the ops
console, Cirrine's contract, Rotation Day, random bytes, nonce, fingerprint (digest), salt,
password record, iteration count, master key, key ID, derived key, sealed command, associated
data, tag, receipt, public key, private key, key agreement, share, the middle (a machine in the
middle), forward secrecy, certificate, certificate authority, root, intermediate, chain,
revocation list, overlap, retire, revoke, TLS, mutual TLS, post-quantum, pad, birthday bound, RSA, OAEP, PSS, API key, webhook secret, Brinewatch.

**Name checks** to run before the season ships: `Apogee Root CA 2028`, `Apogee Stations CA
2028`, `apogee.internal` (a reserved top-level domain, so only the company name matters), and
"Rotation Day" as a phrase. No new people, places, satellites or customers.

### Season 6: on call

After Season 5's Rotation Day, the learner starts the next assignment on Monday 6 November
2028. Their new badge says **Ground Reliability Engineer**, still No. 039, based at Port Ellery
headquarters, since Monday 6 November 2028. That morning at 08:00 Mirela hands them the pager:
they join the on-call rota, a week at a time, with Nkem, Mirela and Tomás. On their second
night, at 02:46 on Wednesday 8 November, the pager wakes them: one pass in five at Lantern Atoll
is failing, the caretaker sees nothing wrong, and uplinkd's log says only `error`. Over the
season the learner makes Apogee's programs explain themselves: logs a machine can search,
numbers that show trouble before a customer does, alerts that wake the right person once, and
profiles that show where the time and the memory go. In February the night shift calls them,
and they find Lantern Atoll's culprit: a slow leak, one stuck goroutine per failed pass, that
grew for weeks between restarts.

Story time runs from Monday 6 November 2028 to Friday 23 February 2029: eleven chapters, ten of
about a week and the last of two, with a quiet holiday fortnight (Tuesday 19 December to Monday
1 January) between chapters 6 and 7. Nkem reviews every change and carries the season; Ingrid
runs the Fishbowl, the night shift and the pager's rules; Tomás knows the radio, suspects Lantern
Atoll's generator and takes the pager over Christmas; Mirela reads the logs for secrets; Yusra
turns customers' questions into numbers; Hollis writes in capitals; Sunniva and Paz have a line
each in February. No new people, places, satellites or customers are introduced. Lantern Atoll's
caretaker (the bible's "a caretaker on site") and the night shift in the Fishbowl stay unnamed:
write around their pronouns. Every time is UTC.

#### The season, chapter by chapter

You work at the backend desks beside the Fishbowl; on-call nights you are at home with the pager
and a laptop, reading logs through the ops console (Season 5). Every change is tried first on
the lab bench at Port Ellery, against a Pipit's passes.

| Ch | Week of | Where | Who needs what | What is at stake |
|---|---|---|---|---|
| 1 | Mon 6 Nov 2028 | your desk; at home on the night of 7–8 Nov | Nkem, Wednesday morning: "`error` is not a log line." uplinkd's lines must say what failed, where and for which pass, in a form a program can search | At 02:46 three of Lantern Atoll's passes had failed overnight, Cirrine-1's sealed command had never gone up, and the log said `error` three times. A restart at 03:14 fixed it and nobody knows why. The new lines reach Lantern Atoll on Friday 10 Nov at 14:00 |
| 2 | Mon 13 Nov | the Fishbowl, with the night shift | Ingrid: a night at INFO she can read at a glance (two lines a pass: where it starts, how it ends), everything when someone asks, and no restart to change it. From Monday 08:00 her checklist restarts Lantern Atoll's uplinkd every Monday "until we know why" | Since Friday's deploy Lantern Atoll's uplinkd runs at DEBUG, to catch the next failure, and writes 2.1 million lines a night; the ops console's log view times out, and the night shift can't find a pass |
| 3 | Mon 20 Nov | the backend desks, then the bench | Nkem: follow the 8 November pass through bookingd, passplan, uplinkd at headquarters and uplinkd at Lantern Atoll in one search | It took Nkem and you all of Monday to line up four programs' lines by hand; next time the night shift must do it in a minute |
| 4 | Mon 27 Nov | Mirela's desk | Mirela finds Furrowlight's API key in 14 of bookingd's lines from Wednesday 22 Nov: 13 DEBUG header dumps and one `request failed` error | Anyone who can read the logs could have booked as Furrowlight. Yusra and Hollis replace the key by 16:00 (Hollis: "NEW KEY. WHY."). On Friday 1 Dec at 16:00 Mirela plants twelve lines with secrets: none may reach a disk |
| 5 | Mon 4 Dec | at home at 08:10, ten minutes into your week on call, then the backend desks | The old `disk` rule pages you at 08:10, and every ten minutes after: Lantern Atoll's disk is 98% full, `uplinkd.log` alone 41 GB. Since 13 Nov its settings start it at DEBUG, "until we know why", and nobody set them back | A full disk leaves the recorder nowhere to keep the next pass, as Ystmark's did in May. The new log files go to every station on Friday 8 Dec |
| 6 | Mon 11 Dec | Yusra's desk, the backend desks | Yusra: Brinewatch wants the exact number of its failed passes in November, and every week from now on. Ingrid's log says 5, uplinkd's lines 3, passplan's 6; bookingd's records say 4 | Logs are not counts. The counters reach every station on Tuesday 12 Dec at 10:00, and Brinewatch's first exact week goes out on Tuesday 19 Dec at 10:00 |
| 7 | Tue 2 Jan 2029 | the Fishbowl | The night of 27–28 Dec: Lantern Atoll's passes failed from evening to dawn, and the pager woke Tomás eleven times for one problem | The new rules go live on Monday 8 Jan at 08:00: one page per problem, to the person on call, and silences that end |
| 8 | Mon 8 Jan | the Fishbowl, Yusra's desk | Ingrid, Tomás and Yusra agree what "every pass counts" means in numbers; Yusra sends customers the quarter's report on Friday 12 Jan | Lantern Atoll lost 22 of 1,199 passes in the last quarter, eleven over its budget. The generator's new part arrives by boat on Thursday 11 Jan |
| 9 | Mon 15 Jan | the backend desks | Tuesday 16 Jan, 05:10: Cape Wren stowed in 74 km/h of wind and passplan took 3 min 52 s to plan the rest of the day. Ingrid: "Four minutes is most of a checklist." | Until 05:14 the Board showed Brine-3's 05:29 pass at Cape Wren as booked, and the night shift started its checklist. By Friday 19 Jan a day plans in 6.1 s, and every program serves profiles on its debug port |
| 10 | Mon 5 Feb | at home on the night of 6–7 Feb; the backend desks | The failures are back at Lantern Atoll, with the generator fixed. Ingrid's note on the Board: "LAN failing: call No. 039 before anyone restarts anything" | The night shift calls you at 01:50; you save the profiles at 02:00 and restart at 02:10; on Wednesday you find the culprit; the fix reaches every station on Friday 9 Feb at 10:00 |
| 11 | Mon 12 Feb, to Fri 23 Feb | the lab bench at Port Ellery; the Fishbowl | Nkem: the incident review on Friday 23 Feb at 15:00 must explain everything since 17 October | On Tuesday 13 Feb the bench replay of 6–7 Feb, traced, shows each failed pass waiting for a worker; spans reach every uplinkd on Monday 19 Feb; the review shows one picture |

**The opening** (Monday 6 November): Apogee's five stations fly passes day and night from the
Fishbowl, and when something breaks at three in the morning the programs' own notes are all
anyone has; the badge; the cast (Nkem, your reviewer; Ingrid, the Fishbowl and the night shift;
Mirela, who hands you the pager; Lantern Atoll, the equatorial island station with its diesel
generator and its caretaker); the rota (you are on call from 08:00, for a week); and, as its last
page, Wednesday 02:46: the pager goes off, one pass in five at Lantern Atoll is failing, the
caretaker sees nothing wrong, the log says `error`. Its button opens lesson 1.1. Words of the
first week: on call, the rota, page, the night shift, log line.

**How each chapter ends** (debrief facts): 1, the new lines at Lantern Atoll from Friday 14:00,
and Ingrid's checklist gets "Monday 08:00: restart LAN uplinkd"; 2, the night shift's view shows
two lines a pass, where it starts and how it ends, and can turn a station up or down without a
restart, while Lantern Atoll's file
keeps everything at DEBUG; 3, the night shift's search
finds the 8 November pass's seven lines in seconds: it died inside Lantern Atoll's uplinkd,
before its lines reached the radio; 4, Mirela signs off on Friday at 16:00: "LGTM once it fails
closed"; 5, Lantern Atoll's logs hold at 1 GiB, the disk about two thirds full; 6, Brinewatch's week, exactly:
62 booked, 61 flown, 1 failed (Wednesday 13 Dec, a downpour); 7, the replayed night pages once;
8, the quarter's report goes out, and the generator's part is fitted on Thursday 11 Jan from 09:00
to 11:00; 9, a day plans in 6.1 s; Ingrid strikes the Monday restart from the checklist on
Monday 15 Jan ("Generator fixed. No more Monday restarts."); 10, the fix at every station;
Sunniva: "Ystmark: one stuck. Wind 40. Coffee holding." Paz: "Salar Quieto, zero stuck. Even the
goroutines are short of breath up here."; 11, the finale.

**The finale** (chapter 11, Friday 23 February, 15:00, the Fishbowl): the incident review. Its
picture is the bench replay's timeline (below): a `fly` span that outlives its own pass by six
minutes, and a pass that waits eight and a half minutes for a worker that never comes. Since Friday 9 February
Lantern Atoll's uplinkd has run fourteen days without a restart, its goroutines at 14 every hour;
two passes were lost to downpours, and nothing stayed stuck. Ingrid strikes the runbook's last
line, "restart uplinkd". The debrief closes the season and points at Season 7, in the Seasons
table's words: twenty stations, eighty satellites, and passplan's rewrite.

#### Lantern Atoll this season

- **The wet season.** November to February brings the atoll's downpours, most in the evening and
  at night. They fade even S-band for minutes at a time: the radio logs `carrier lost` ten times
  a second while the signal is under its threshold.
- **Its passes.** About nineteen a day in view, two or three from each satellite, Brine-3's the
  most (its 53° orbit suits an equatorial station); about thirteen of them are booked. Its UHF yagis and its dish can
  fly two passes at once. About one pass in a hundred has always ended without the satellite's
  acknowledgment (a downpour, a low pass behind the palms over the 8° mask): Lantern Atoll's
  natural losses, more than any other station's (about one in a thousand elsewhere).
- **The generator.** In late October it started to sound different (Season 5's last debrief): a
  worn fuel injector. The power stayed steady throughout and was never the cause. Tomás suspects
  it from 8 November ("A dip in the power, the radio hangs, the pass times out."), and after
  Christmas, when the DEBUG lines show passes waiting for a worker, he has a reason for that too
  ("The radio hung, so the worker never came back."). The part is ordered by boat on 9 November,
  arrives on Thursday 11 January and is fitted from 09:00 to 11:00 with the station off the air.
- **The station computer**, in the station building: 4 cores, 8 GiB of memory, a 120 GB disk
  shared by the logs and the recordings waiting for headquarters to collect them.
- **The caretaker** checks the dish, the radio, the yagis and the generator when the night shift
  asks, in the atoll's morning. Unnamed; no pronoun; at most one line per chapter.

#### The leak

**In plain words.** uplinkd at each station runs every pass on one of four workers. When a pass
ended without its acknowledgment, the code that ran it gave up at LOS as it should, but the
goroutine it had started stayed behind, stuck forever, still holding its worker. Each pass that
failed left one more stuck goroutine; after four natural losses all four workers were held, and
from then on every pass at the station waited for a worker until LOS and failed, leaving one more
goroutine behind. A restart cleared everything, so it came back weeks later.

**The code** (`uplink.go`, written in 2025, the same at every station until 9 February 2029):

```go
// runPass runs one pass on a free worker and waits for it until LOS.
func (u *Uplink) runPass(ctx context.Context, p PassView, lines []string) error {
	buf := make([]byte, 0, 4<<20) // the pass's frames, up to 4 MiB
	done := make(chan result)     // unbuffered: the bug
	go func() {
		u.workers <- struct{}{}          // take one of four workers; wait while all are busy
		rep, err := u.fly(p, lines, buf) // no context: the radio gives up 15 min after its last line
		done <- result{rep, err}         // after LOS nobody receives: blocks forever
		<-u.workers                      // so the worker is never given back
	}()
	select {
	case res := <-done:
		return res.err
	case <-ctx.Done(): // the pass's context ends at LOS (Season 1)
		return fmt.Errorf("pass %s: %w", p.ID, ctx.Err())
	}
}
```

- `runPass` starts 2 minutes before AOS (pass prep) with a context whose deadline is LOS.
- `fly` waits for AOS, hands the lines to the radio 10 seconds after AOS and returns at the
  satellite's acknowledgment of the last line (usually within two minutes), or with `ErrNoAck`
  when the radio gives up, 15 minutes after its last line: always after LOS.
- A flown pass holds its worker from pass prep to its acknowledgment, a few minutes: at a
  station flying two passes at once, three held workers still leave room, so three stuck
  goroutines change nothing anyone can see.
- A pass with no acknowledgment: `runPass` returns at LOS; 15 minutes after the last line `fly`
  returns, and the goroutine blocks on `done <- result{rep, err}` forever, holding its worker and,
  through `rep.Frames`, its 4 MiB buffer.
- A pass that finds all four workers held: its goroutine blocks on `u.workers <- struct{}{}`
  forever (no worker is ever given back), holding its 4 MiB buffer; `runPass` returns at LOS.
- Every failure therefore reads `pass <id>: context deadline exceeded`, whether the satellite was
  silent or the pass never got a worker; logs and counters could not tell them apart, which is why
  the failures looked like the radio, and the generator took the blame.
- Each failed pass: one stuck goroutine and 4 MiB more in use. uplinkd's normal count is 14
  goroutines and about 9 MiB in use.
- **Why Lantern Atoll:** it has the most natural losses, so its workers ran out first; Ystmark's
  uplinkd held one stuck goroutine on 9 February, the others none.
- **Why restarts helped:** a restart ends every goroutine. Deploys restarted it too, and the
  Monday restarts (13 November to 8 January) kept it under four.

**Restarts of Lantern Atoll's uplinkd**, the arithmetic of the whole season: Tue 17 Oct 2028
01:15:40 (onto its new certificate, Rotation Day); Wed 8 Nov 03:14 (you, at the page); the
deploys of Fri 10 Nov, Fri 17 Nov, Fri 24 Nov, Fri 1 Dec and Fri 8 Dec at 14:00 and Tue 12 Dec at
10:00 (chapters 1 to 6); the Monday restarts at 08:00 on 13, 20 and 27 Nov, 4, 11 and 18 Dec, 1
and 8 Jan (25 Dec skipped: Christmas Day); Thu 28 Dec 03:10 (Tomás); Fri 19 Jan 14:00 (the deploy
that puts profiles on every program's debug port, chapter 9); Wed 7 Feb 02:10 (you, after the
profiles); Fri 9 Feb 10:00 (the fix). No deploy reaches the station during the holiday fortnight
or between 19 January and 7 February.

**Passes Lantern Atoll lost.** Natural (no acknowledgment): in October 3, 11, 24 and 30 Oct; in
November 5 Nov (`furrow-b/LAN/20957`), 8 Nov (`brine-3/LAN/14383`), and 15, 22 and 29 Nov
(Brine-3); in December 6 and 13 Dec (Brine-3), 20, 23 and 26 Dec, 27 Dec (`brine-3/LAN/15145`);
in January 3, 19 (20:48, after that day's deploy) and 30 Jan; in February 4 Feb and 7 Feb
(`cirrine-1/LAN/17759`). Lost to the leak (no worker): 8 Nov `cirrine-1/LAN/16402` and
`furrow-a/LAN/20998`, the five of 27–28 Dec and 7 Feb `furrow-b/LAN/22379`: **eight passes**.
Brinewatch's losses in November (bookingd's records): 4, on 8, 15, 22 and 29 Nov.

**The profiles you saved on 7 February at 02:00** (from the debug port, chapter 9's endpoints):

- `goroutines.txt` (`/debug/pprof/goroutine?debug=1`): `goroutine profile: total 19`; 4 in
  `main.(*Uplink).runPass.func1` on the line `done <- result{rep, err}`; 1 on the line
  `u.workers <- struct{}{}`; 14 of the program's own.
- `heap.prof` (`/debug/pprof/heap`): 29 MiB in use, 20 MiB of it five buffers made in `runPass`.
- `/debug/vars`: `goroutines` 19 (Ystmark 15, the other stations 14), `heap_bytes` 30408704.
- The leak profile (`pprof.Lookup("goroutineleak")`, Go 1.27) reports four of them, the ones on
  `done <- result{rep, err}`: nothing else can reach `done`. The fifth waits on `u.workers`, which
  uplinkd itself still holds, so the leak profile can't call it stuck; the goroutine count and
  `goroutines.txt` show it.

**The fix** (Friday 9 February 10:00, every station; chapter 10): no goroutine at all. Take the
worker with a `select` on the context, give it back with `defer`, and pass the context to `fly`,
which stops at LOS. Failures now say why:

```go
var (
	ErrNoAck    = errors.New("no ack before LOS")
	ErrNoWorker = errors.New("no free worker before LOS")
)

func (u *Uplink) runPass(ctx context.Context, p PassView, lines []string) error {
	select {
	case u.workers <- struct{}{}:
	case <-ctx.Done():
		return fmt.Errorf("pass %s: %w", p.ID, ErrNoWorker)
	}
	defer func() { <-u.workers }()
	_, err := u.fly(ctx, p, lines, make([]byte, 0, 4<<20)) // ErrNoAck when ctx ends first
	if err != nil {
		return fmt.Errorf("pass %s: %w", p.ID, err)
	}
	return nil
}
```

Lesson 10.7's intermediate fix keeps the goroutine but gives `done` room for one
(`make(chan result, 1)`) and gives the worker back with `defer`: nothing is stuck forever, but a
failed pass still holds its worker until the radio gives up.

#### The 02:46 page

The night of Tuesday 7 to Wednesday 8 November 2028, Lantern Atoll. Its uplinkd had run since
Rotation Day, and three natural losses (24 Oct, 30 Oct, Sun 5 Nov 14:48:54) held three of its
four workers.

| Pass | Booking | AOS | LOS | Max elev | What happened |
|---|---|---|---|---|---|
| `brine-3/LAN/14383` | `bk-5608` | Wed 8 Nov 01:48:20 | 01:54:41 | 19.1° | prep 01:46:20 takes the last worker; lines at 01:48:30; a downpour, no acknowledgment; failed at LOS; the radio gives up at 02:03:30 and the goroutine blocks: four of four workers held |
| `cirrine-1/LAN/16402` | `bk-5611` | 01:55:09 | 02:02:39 | 22.6° | prep 01:53:09, waits for a worker; failed at LOS; Cirrine-1's sealed `PING` (seq 4688, sealed with `k-2028-10` at headquarters at 01:44:40, queued at Lantern Atoll at 01:48:30) never went up |
| `furrow-a/LAN/20998` | `bk-5613` | 02:39:45 | 02:46:24 | 21.7° | prep 02:37:45, waits; failed at LOS |
| `brine-3/LAN/14384` | `bk-5609` | 03:26:17 | 03:32:10 | 16.6° | flown, after the restart (acknowledged at 03:27:38) |

The night's earlier passes flew on the one free worker. Furrow-B's 01:53:06 pass overlaps
Brine-3's and Cirrine-1's on the station's one S-band dish and wasn't booked; no other pass was
due between 02:03 and 03:26.

| Time | What |
|---|---|
| Wed 8 Nov 01:54:41 | the first failure; uplinkd writes `2028/11/08 01:54:41 error` |
| 02:02:39 | the second: `2028/11/08 02:02:39 error` |
| 02:46:24 | the third: `2028/11/08 02:46:24 error` |
| 02:46:34 | the old `pass-failures` rule pages you: `LAN: 3 of the last 15 passes failed` |
| 02:48 | you open Lantern Atoll's uplinkd log in the ops console: three lines, each ending `error` |
| 02:54 | the night shift calls the caretaker (morning on the atoll) |
| 03:05 | the caretaker: dish, radio and yagis fine; the generator sounds different, the power steady |
| 03:14 | you follow the runbook's last line and restart uplinkd; six stuck goroutines and 24 MiB go with it, unseen |
| 03:32:10 | `brine-3/LAN/14384` flies |
| 08:00 | the morning meeting. Ingrid: "One pass in five at Lantern Atoll. Cause unknown. A restart fixed it." Nkem: "Then the log is the first thing we fix." Yusra: "TL;DR: Cirrine asks why its 01:55 command never left the ground." |
| Thu 9 Nov | Ingrid adds the rule `lan-failed-pass`, "until we know why"; Tomás orders the generator part |

The old lines come from Go's `log` package with its standard flags, on a station clock set to
UTC: `2028/11/08 02:46:24 error`.

#### The night of eleven pages

Wednesday 27 to Thursday 28 December 2028, Lantern Atoll; Tomás on call. The Monday restart of
25 December was skipped; downpours on 20, 23 and 26 December (the last at 08:41, before the
night's last fifteen passes) held three workers.

| Pass | AOS | LOS | Result |
|---|---|---|---|
| `brine-3/LAN/15145` | 27 Dec 19:33:58 | 19:41:28 | failed (natural); its goroutine keeps the last worker from 19:49:08 |
| `pipit-1/LAN/14031` | 20:13:04 | 20:20:49 | failed (no worker) |
| `cirrine-1/LAN/17147` | 28 Dec 00:37:36 | 00:45:16 | failed (no worker) |
| `furrow-a/LAN/21758` | 01:49:55 | 01:57:49 | failed (no worker) |
| `cirrine-1/LAN/17148` | 02:14:08 | 02:19:14 | failed (no worker) |
| `furrow-b/LAN/21756` | 02:37:22 | 02:44:57 | failed (no worker) |
| `brine-3/LAN/15152` | 07:18:19 | 07:26:04 | flown (Tomás restarted uplinkd at 03:10: nine stuck goroutines, 36 MiB) |
| `pipit-1/LAN/14038` | 08:04:26 | 08:12:34 | flown |

The old pager sent, ten seconds after each LOS, `lan-failed-pass` for each failure and
`pass-failures` whenever 3 or more of the last 15 passes had failed: 19:41:38 `lan-failed-pass`;
20:20:59 `lan-failed-pass`; 00:45:26 both; 01:57:59 both; 02:19:24 both; 02:45:07 both; and
07:26:14 `pass-failures` for a pass that had flown. **Eleven pages.** At 07:28 Tomás muted
Lantern Atoll's pages with no end; Ingrid found the mute on Friday 29 December at 09:00. Tomás:
"Eleven times. Like a buoy bell in a swell, and the same problem every time." The `goroutines`
gauge (chapter 6) read 23 at 03:00; no rule looked at it.

The new rules (below), replayed on the same night, send exactly: `20:21 FIRING pass-failures
LAN 2` to Tomás Ferreira-Lund; six `failed-pass` lines to the morning list (19:42, 20:21,
00:46, 01:58, 02:20, 02:45); and `08:45 RESOLVED pass-failures LAN 0` to Tomás, which doesn't
wake him. Nothing else.

#### The night of 6–7 February

Lantern Atoll's uplinkd has run since the deploy of Friday 19 January at 14:00 (the Monday
restarts ended on 8 January); natural losses on 19 January (20:48), 30 January and 4 February
hold three workers (goroutines 15, 16, 17). Nkem is on call.

| Time | What |
|---|---|
| Wed 7 Feb 01:16:08 | `cirrine-1/LAN/17759` (AOS 01:18:08, LOS 01:27:13, booking `bk-5851`) takes the last worker; lines at 01:18:18; no acknowledgment |
| 01:27:13 | failed at LOS |
| 01:33:18 | the radio gives up; the goroutine blocks: four of four workers held |
| 01:36:52 | `furrow-b/LAN/22379` (AOS 01:38:52, LOS 01:45:23, `bk-5853`) waits for a worker |
| 01:45:23 | failed at LOS |
| 01:46:00 | `01:46 FIRING pass-failures LAN 2` to Nkem Halloran |
| 01:50 | the night shift calls you, as Ingrid's note asks |
| 02:00 | you save `heap.prof` and `goroutines.txt` from Lantern Atoll's uplinkd |
| 02:10 | you restart it |
| 02:22:30 | `furrow-a/LAN/22382` (AOS 02:24:30, LOS 02:32:07, `bk-5855`) takes a worker; acknowledged at 02:25:42: flown |
| Wed 7 Feb, morning | you read the profiles and find the line; the generator was innocent |
| Thu 8 Feb | Nkem reviews the fix: "What does the worker see if this pass never ends?" |
| Fri 9 Feb 10:00 | the fix at every station, with the new `goroutines` rule |

#### On call

- **The rota.** Four engineers, a week each, handing over on Mondays at 08:00 UTC (the Board's
  rule). In lessons the learner's weeks read `No. 039` (lessons never use the learner's name).

  | From (Monday 08:00) | On call |
  |---|---|
  | 2028-10-30 | Mirela Dunn |
  | 2028-11-06, 2028-12-04, 2029-01-01, 2029-01-29 | No. 039 |
  | 2028-11-13, 2028-12-11, 2029-01-08, 2029-02-05 | Nkem Halloran |
  | 2028-11-20, 2028-12-18, 2029-01-15, 2029-02-12 | Mirela Dunn |
  | 2028-11-27, 2028-12-25, 2029-01-22, 2029-02-19 | Tomás Ferreira-Lund (he takes the holiday week, as he took launch week) |

  As `[]Shift` (the Board's type), oldest first. At Monday 07:59 the previous week's name still
  holds.
- **The pager** (the bible's `pager` service) receives every pass's result from the stations'
  uplinkd ten seconds after LOS, and reads every program's `/debug/vars`. A page wakes the
  person on call; **the morning list** is Ingrid's list of things that can wait for daylight,
  read at the 08:00 meeting.
- **The runbook** for Lantern Atoll ended, until 23 February, in "if passes fail and the station
  looks fine, restart uplinkd". From 5 February it starts with "save `/debug/pprof/heap` and
  `/debug/pprof/goroutine?debug=1`, then call the owner".
- **The Monday restart**: Ingrid's checklist line "Monday 08:00: restart LAN uplinkd", from
  13 November 2028 to 8 January 2029.

#### Logs

Every program logs with `log/slog`, one record per event, to its log file.

- **Handlers.** Production: `slog.NewJSONHandler` with `HandlerOptions{Level: <the program's
  *slog.LevelVar>, ReplaceAttr: <chapter 4's redaction>}`, writing to chapter 5's `Rotator`,
  wrapped (outermost first) by chapter 5's `Sampler` and chapter 3's `passHandler`, and the
  logger made with `.With("program", <name>, "site", <site>)`. People at a terminal read the
  same records through `slog.NewTextHandler`.
- **Time.** Records carry UTC. In lessons the given `stamp` writes the story's clock on every
  record. The JSON handler writes a time in RFC 3339 with the fraction trimmed
  (`2028-11-08T02:02:39Z`, `2028-11-08T02:02:39.4Z`); the text handler always writes
  milliseconds (`2028-11-08T02:02:39.000Z`). Time attributes such as `aos` print the same way in
  each. `code` lessons leave the time out with the given `dropTime`, so a line fits the screen:
  `level=ERROR msg="pass failed" site=LAN`.
- **Levels.** `DEBUG`, `INFO`, `WARN`, `ERROR` only. Programs start at INFO. The night shift
  changes a running program's level with `PUT /debug/level` on its debug port (below).
- **Messages** are short, lowercase, fixed text with no values inside; the values are
  attributes.
- **Groups** are written with `slog.Group`, never `Logger.WithGroup`: an attribute that a
  handler adds in `Handle` lands inside any group opened with `WithGroup`.

| Program | Message | Level | Its own attributes |
|---|---|---|---|
| uplinkd (station) | `pass started` | INFO | `pass` group |
| | `worker taken`, `lines sent`, `ack received` | DEBUG | `lines`, `ack_ms` |
| | `frame received` | DEBUG | `bytes`, `hex` |
| | `carrier lost` | WARN | `margin_db` |
| | `pass flown` | INFO | `pass` group, `max_elev_deg`, `wait_ms`, `ack_ms` (from chapter 3) |
| | `pass failed` | ERROR | `pass` group, `max_elev_deg`, `wait_ms`, `err` |
| | `pass skipped` | WARN | `pass` group, `reason` |
| | `lines queued` | INFO | `lines` (on `POST /uplink`) |
| uplinkd (headquarters) | `command sealed` | INFO | `seq` |
| | `lines delivered` | INFO | `status` |
| bookingd | `booking created` | INFO | `booking_id` |
| | `request failed` | ERROR | `err` |
| | `request headers` | DEBUG | `header`: the line Mirela found; removed on 27 Nov |
| | Season 3's access line, through `slog.NewLogLogger` | INFO | the text as `msg`: `req-000007 POST /bookings 201 236B 4ms` |
| passplan | `pass planned` | INFO | `booking_id` |
| | `day planned` | INFO | `passes`, `ms` |
| pager | `alert sent` | INFO | `rule`, `station`, `kind`, `to` |
| any (chapter 11) | `span` | INFO | `trace_id`, `span_id`, `parent_id`, `name`, `start`, `ms` |

**Attribute keys**, the same in every program. Units are in the name; durations are whole
milliseconds in an `_ms` key (`slog.Int64(key+"_ms", d.Milliseconds())`), because the JSON
handler writes a `time.Duration` as nanoseconds (`570000000000`).

| Key | Value |
|---|---|
| `program` | `uplinkd`, `bookingd`, `passplan`, `pager` |
| `site` | where the program runs: a station code (`LAN`), or `HQ` for headquarters' rack |
| `pass_id` | the pass ID, `cirrine-1/LAN/16402`; from chapter 3 added by `passHandler` from the context, after the record's own attributes, never with `With` |
| `request_id` | Season 3's `req-000042`, from the context the same way |
| `pass` | group: `sat_id`, `station`, `aos`, `los` (times) |
| `booking_id` | `bk-5611` |
| `max_elev_deg` | number, `22.6` |
| `wait_ms`, `ack_ms`, `ms` | whole milliseconds: from pass prep to the end, from the lines to the acknowledgment, a step's length |
| `lines`, `bytes`, `seq`, `status`, `passes`, `dropped` | integers |
| `hex` | a frame's bytes as lowercase hex, on DEBUG `frame received` only (chapter 2) |
| `margin_db` | number, the radio's signal margin |
| `err` | the error's text (`slog.Any("err", err)`) |
| `reason` | `stowed` |
| `rule`, `station`, `kind`, `to` | the pager's alert (below) |
| `trace_id`, `span_id`, `parent_id`, `name`, `start` | a span (chapter 11) |

The canonical failure line from chapter 3 on (one line in the file):

```
{"time":"2028-11-08T02:02:39Z","level":"ERROR","msg":"pass failed","program":"uplinkd","site":"LAN","pass":{"sat_id":"cirrine-1","station":"LAN","aos":"2028-11-08T01:55:09Z","los":"2028-11-08T02:02:39Z"},"max_elev_deg":22.6,"wait_ms":570000,"err":"pass cirrine-1/LAN/16402: context deadline exceeded","pass_id":"cirrine-1/LAN/16402"}
```

Chapter 1's version, before `With`, groups and the context: `{"time":"2028-11-08T02:02:39Z",
"level":"ERROR","msg":"pass failed","pass_id":"cirrine-1/LAN/16402","max_elev_deg":22.6,
"wait_ms":570000,"err":"pass cirrine-1/LAN/16402: context deadline exceeded"}`.

**The 8 November pass in four programs** (chapter 3's replay; seven lines, other passes'
lines among them):

| Time | Program | Site | Message | Also |
|---|---|---|---|---|
| Mon 6 Nov 10:14:03 | bookingd | HQ | `booking created` | `booking_id` `bk-5611`, `request_id` `req-004211` |
| Tue 7 Nov 00:00:41 | passplan | HQ | `pass planned` | `booking_id` `bk-5611` |
| Wed 8 Nov 01:44:40 | uplinkd | HQ | `command sealed` | `seq` 4688 |
| 01:48:30 | uplinkd | HQ | `lines delivered` | `status` 202 |
| 01:48:30 | uplinkd | LAN | `lines queued` | `lines` 1 |
| 01:53:09 | uplinkd | LAN | `pass started` | `pass` group |
| 02:02:39 | uplinkd | LAN | `pass failed` | `wait_ms` 570000, `err` |

Lines with equal times keep their input order (headquarters' log is given first).

The replay's other passes (chapter 3's fixtures): bookings `bk-5608` (`brine-3/LAN/14383`, Mon 6
Nov 09:41:27, `req-004187`), `bk-5609` (`brine-3/LAN/14384`, 09:41:28, `req-004188`) and `bk-5613`
(`furrow-a/LAN/20998`, 15:20:55, `req-004302`); passplan plans all four at 00:00:41, by AOS;
Brine-3's lines are delivered and queued at 01:31:05 and 03:18:40 (`lines` 2 each); `carrier
lost` ten times at 01:51:02 and five at 01:52:47, `margin_db` -1.5; `wait_ms` 501000, 570000,
519000 and 201000; the 03:27:38 `pass flown` has `ack_ms` 71000.

**Headers between programs.** `Apogee-Pass: cirrine-1/LAN/16402` carries the pass ID on
`POST /uplink` from headquarters to a station (chapter 3); `Apogee-Trace` carries a span
(chapter 11). An empty or missing header means none.

**The debug port.** Every program serves, on a private port (the station network only; in
lessons an `httptest` server on 127.0.0.1):

| Route | Answer | From |
|---|---|---|
| `PUT /debug/level` | body `debug`, `info`, `warn` or `error`, any case (`slog.LevelVar.UnmarshalText`, which also takes offsets such as `info+2`): 204; a name it doesn't know 400 `unknown level "loud"`; another method 405 | ch 2 |
| `GET /debug/vars` | `expvar.Handler()` | ch 6 |
| `GET /debug/pprof/heap`, `/debug/pprof/goroutine` (`?debug=1` for text), `/debug/pprof/profile` | `net/http/pprof`'s `Handler(name)` and `Profile`, mounted on the program's own mux | 19 Jan (ch 9) |

**Never in a log** (chapter 4), from Friday 1 December 2028:

- Secret headers, redacted before a header is logged: `Authorization`, `Cookie`, `Set-Cookie`,
  `Apogee-Signature`; each value becomes `[redacted]`, on a copy. `safeLogger`'s `ReplaceAttr`
  hands any `http.Header` value to `safeHeader`, because the JSON handler writes a header's values
  without calling `ReplaceAttr` on them.
- Attribute keys always redacted by `ReplaceAttr`, in any letter case and any group: `api_key`,
  `authorization`, `password`, `secret`, `token`. The value becomes `[redacted]`; the key stays.
- API keys and tokens inside any text (messages, string values, error texts) are scrubbed with
  one pattern compiled once, `\b([a-z]+)\.(?:[0-9a-f]{32}|[0-9]+\.[0-9a-f]{64})\b`, replaced by
  `$1.[redacted]`: `furrowlight.5b0e…7c6b` becomes `furrowlight.[redacted]`.
- Types that hold a secret implement `slog.LogValuer`: `APIKey` logs `furrowlight.[redacted]`
  (`[redacted]` when it has no dot); Season 5's `Key` logs its `String()`,
  `k-2028-10 (secret hidden)`. The JSON handler ignores `String` and marshals exported fields,
  so without `LogValue` a `Key` prints its secret in base64.
- Personal data: customers' contact addresses, kept by bookingd for webhook failure notices
  (`hollis.vey@furrowlight.example`, `ops@brinewatch.example`, `uplink@cirrine.example`,
  `fishbowl@apogee.example`: `.example` is reserved, so none is real), are logged under `email`
  and masked to the local part's first character, `***@` and the domain:
  `h***@furrowlight.example`. Anything that isn't one address with one `@` and both parts
  non-empty becomes `[redacted]`.
- Test secrets are Season 3's and Season 5's fixtures: Furrowlight's test key
  `furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b`, its token shape
  `furrowlight.1830682800.<64 hex>`, `k-2028-10`'s test secret `uplink next master key, 32 bytes`
  (base64 `dXBsaW5rIG5leHQgbWFzdGVyIGtleSwgMzIgYnl0ZXM=`), Season 5's password
  `fails closed, always`, and, for other customers, Season 3's
  `brinewatch.0f1e2d3c4b5a69788796a5b4c3d2e1f0`, `cirrine.c3a1f7e95d2b4086a7e1c9f3b5d8a024` and the
  webhook secrets `furrowlight webhook secret` and `cirrine webhook secret`. Production keys never
  appear.

**Log files** (chapter 5), from Friday 8 December 2028:

- Folder `/var/log/apogee/` on each computer (`t.TempDir()` in lessons). The current file is
  `<program>.log`; rotated files are `<program>-<UTC day>.<n>.log`, n counting from 1 within
  that day: `uplinkd-2028-12-04.1.log`. Files are opened with
  `O_APPEND|O_CREATE|O_WRONLY`, mode 0640.
- A record is one `Write`, so a line is never split. Before a write, the `Rotator` rotates when
  the current file is not empty and `size+len(p) > MaxBytes` (64 MiB; 1 KiB or 64 KiB in tests),
  or when the clock's UTC day is after the current file's day; the old file is renamed under its
  own day (an empty current file is not renamed, it only takes the new day). A single record
  longer than `MaxBytes` goes alone into a fresh file. A clock that steps back never rotates. An existing current file opened at start takes its day from its
  modification time, and n continues after the highest of that day's rotated files.
- Pruning, each time the current file is opened (at start and after each rotation, through
  `OnOpen`): delete the program's rotated files whose day is more than
  `KeepDays` (14) days before today's, then the oldest (by day, then n) until the rotated files
  fit `MaxTotal` (1 GiB per program; 1 MiB in tests). The current file and other programs' files
  are never touched; a name that doesn't parse is left alone.
- **Sampling.** Records are counted per level and message within the record's UTC minute. The
  first 10 of a minute pass; after that one in a hundred (the 110th, 210th, ...); ERROR records
  always pass. A passed record after drops carries `dropped`, the number dropped since that
  level and message last passed, carried into the next minute: 250 `carrier lost` warnings in
  one minute write 12 lines (the 110th and 210th with `dropped` 99), and the next minute's first
  carries `dropped` 40. The sampler shares its counts with the loggers `With` makes, and keeps a
  total for `log_lines_dropped_total`.
- On 4 December at 08:10 the disk held 41 GB of `uplinkd.log` (most of it DEBUG `frame received`
  lines, each with its frame in hex; `carrier lost` repeats ten times a second through a downpour), 64 GB of recordings waiting and 13 GB of system: 98%. At 08:25 you turn
  Lantern Atoll's uplinkd down to INFO on its debug port and the file stops growing (deleting it
  would free nothing while uplinkd holds it open); Friday's deploy rotates it on its first line and
  prunes it, and after that the disk stays about two thirds full. A program's starting level comes
  from its settings file; the
  night shift had set Lantern Atoll's uplinkd to start at DEBUG on 13 November, so restarts and
  deploys kept it there.

#### Metrics

Built by hand with `sync/atomic` (chapter 6), published with `expvar` at `GET /debug/vars` on
each program's debug port, read by the pager every five minutes. Names are lowercase with a unit
or `_total`; label values are few (satellites, stations, outcomes, reasons), never a pass ID.

```go
// Counter counts events; it only goes up. Safe for many goroutines.
type Counter struct{ n atomic.Uint64 }

func (c *Counter) Inc()
func (c *Counter) Add(n uint64)
func (c *Counter) Value() uint64

// Gauge is a reading that goes up and down.
type Gauge struct{ v atomic.Int64 }

func (g *Gauge) Set(v int64)
func (g *Gauge) Add(d int64)
func (g *Gauge) Value() int64

// CounterVec keeps a counter per pair of label values, such as satellite and outcome.
type CounterVec struct {
	mu sync.RWMutex
	m  map[[2]string]*Counter // made at first use
}

func (v *CounterVec) With(a, b string) *Counter                // the same *Counter for the same pair
func (v *CounterVec) Snapshot() map[string]map[string]uint64 // {"brine-3": {"failed": 1, "flown": 23}}

// Histogram counts observations into fixed buckets. Safe for many goroutines.
type Histogram struct {
	Bounds []float64       // upper bounds in seconds, ascending
	counts []atomic.Uint64 // one per bound, then one above every bound
	count  atomic.Uint64
	sumMs  atomic.Int64
}

func NewHistogram(bounds ...float64) *Histogram
func (h *Histogram) Observe(d time.Duration) // the first bucket whose bound >= d in seconds
func (h *Histogram) Snapshot() HistogramSnapshot

type HistogramSnapshot struct {
	Bounds []float64 `json:"bounds"`
	Counts []uint64  `json:"counts"` // len(Bounds)+1
	Count  uint64    `json:"count"`
	SumMs  int64     `json:"sum_ms"`
}

// Metrics are a station uplinkd's numbers.
type Metrics struct {
	Passes   CounterVec // passes_total: by satellite and outcome
	Failures CounterVec // failures_total: by satellite and reason
	AckWait  *Histogram // ack_wait_seconds
	Running  Gauge      // passes_running
}

// Reading is one value of a counter, read by the pager every five minutes.
type Reading struct {
	At    time.Time
	Value uint64
}
```

| Name | Kind | Unit and labels | From |
|---|---|---|---|
| `passes_total` | counter | passes, by satellite and outcome (`flown`, `failed`, `skipped`) | each station's uplinkd |
| `failures_total` | counter | failed passes, by satellite and reason: `deadline` until 9 Feb; then `no_ack`, `no_worker`, `radio` | each station's uplinkd |
| `ack_wait_seconds` | histogram | seconds from lines sent to acknowledgment; bounds 0.5, 1, 2, 5, 10, 30, 60, 120, 300 | each station's uplinkd |
| `passes_running` | gauge | passes now between prep and their end | each station's uplinkd |
| `goroutines` | gauge | `runtime.NumGoroutine()`, read when asked (14 for a healthy uplinkd) | every program |
| `heap_bytes` | gauge | `HeapAlloc` from `runtime.ReadMemStats`, read when asked | every program |
| `log_lines_dropped_total` | counter | records the sampler dropped | every program |
| `disk_used_ratio` | gauge | 0 to 1, the log and recordings disk | each station computer |

Exactly 2 s goes in the `≤ 2` bucket; 0 s in the first; 301 s above all. `sum_ms` adds whole
milliseconds. A published counter or vec is an `expvar.Func` returning its snapshot, so the JSON
is read at request time; `cmdline` and `memstats` stay in it:

```
"passes_total": {"brine-3": {"failed": 1, "flown": 23}, "pipit-1": {"flown": 12}},
"ack_wait_seconds": {"bounds": [0.5,1,2,5,10,30,60,120,300], "counts": [0,3,40,15,2,1,0,0,0,0], "count": 61, "sum_ms": 143210},
"goroutines": 14
```

- **Totals from readings.** A counter starts again at 0 when its program restarts. A total over a
  period is the sum of the increases between readings, a drop counting as a restart (the new
  value is the increase since): readings 3, 5, 0, 1 make 3. Last minus first is wrong, and with
  `uint64` it wraps around.
- **Brinewatch's first exact week**, Tuesday 12 December 10:00 to Tuesday 19 December 10:00:
  2,017 readings of Lantern Atoll's counters, five minutes apart, the Monday restart at 08:00:30
  on 18 December between the 08:00 and 08:05 readings. Brine-3 at Lantern Atoll: flown 23,
  failed 1 (Wednesday 13 December, a downpour). All stations: 62 booked, 61 flown, 1 failed.
- **November's argument**: Ingrid's log says 5 (the 7 November failure written twice), uplinkd's
  lines 3 (its lines before 10 November say only `error`), passplan's 6 (it counts replanned
  passes as failed), bookingd's records 4.

#### Alerts

The pager (chapter 7). Rules are evaluated every minute on the minute; a pass's result counts
from its LOS.

```go
// Rule decides when a station's metric is worth telling someone.
type Rule struct {
	Name       string        // "pass-failures"
	Severity   Severity      // Critical pages; Warning goes on the morning list; Info is only logged
	FireAt     float64       // the problem starts when the value is at or above this
	ClearBelow float64       // and is over only when the value falls below this
	For        time.Duration // the value must stay at or above FireAt this long first
	Repeat     time.Duration // a problem still firing is sent again after this long; 0: never
}

// Sample is one value of a rule's metric for one station.
type Sample struct {
	At      time.Time
	Station string
	Value   float64
}

// Alert is one message the pager sends.
type Alert struct {
	At      time.Time
	Rule    string
	Station string
	Kind    string  // "FIRING" or "RESOLVED"
	Value   float64
	To      string  // the on-call name, "morning list" or "log"
}

func (a Alert) String() string // "20:21 FIRING pass-failures LAN 2": At as 15:04, Value as %g

// Silence holds back a station's pages for a while. It always ends.
type Silence struct {
	Station string
	From    time.Time // included
	Until   time.Time // excluded
}

// Pager keeps each rule's state per station.
type Pager struct {
	Rules    []Rule
	Rota     []Shift
	Silences []Silence
	state    map[[2]string]*ruleState // by rule name and station
}

func (p *Pager) Observe(rule string, s Sample) (Alert, bool)
```

**Each rule and station is in one of three states**:

1. **Inactive.** A sample at or above `FireAt` makes it pending since that sample's time; with
   `For` 0 it fires at once.
2. **Pending.** A sample below `FireAt` makes it inactive, sending nothing. A sample at or above
   `FireAt` whose time is at least `For` after the pending start makes it firing and sends
   FIRING.
3. **Firing.** A sample below `ClearBelow` makes it inactive and sends RESOLVED (Critical rules
   only; Warning and Info send no RESOLVED). Otherwise, when `Repeat` isn't 0 and at least
   `Repeat` has passed since the last FIRING sent, it sends FIRING again.

**Routing.** Critical goes to the on-call name at the alert's time (the latest `Shift` whose
`From` is at or before it); a RESOLVED goes to whoever is on call at its own time and doesn't wake them; Warning
goes to `morning list`; Info to `log`. **Silences.** While a silence covers the station, alerts
are held back but the state still changes; a problem that started during the silence and is still
firing when it ends is sent at the first sample at or after `Until`, and one that resolved by then
sends nothing, RESOLVED included; a reminder that fell due during the silence is sent then too.

**The rules from Monday 8 January 2029, 08:00**:

| Rule | Metric, per station | FireAt | ClearBelow | For | Repeat | Severity |
|---|---|---|---|---|---|---|
| `pass-failures` | failed passes with LOS in the last 6 hours, `(now−6h, now]` | 2 | 1 | 0 | 24h | Critical |
| `failed-pass` | failed passes with LOS in the last minute | 1 | 1 | 0 | 0 | Warning |
| `disk-full` | `disk_used_ratio` | 0.90 | 0.85 | 15m | 12h | Critical |
| `goroutines` (from Fri 9 Feb) | a program's goroutines minus its baseline (14 for uplinkd) | 3 | 1 | 1h | 24h | Warning |

Skipped passes are not failures. Replayed, the `goroutines` rule would have put Lantern Atoll on
the morning list on Monday 5 February, two days before the night.

**The old rules, until 8 January** (no state, no holding, no de-duplication; mutes with no end):
`pass-failures`, ten seconds after each LOS, pages when at least 3 of the station's last 15
passes failed (`LAN: 3 of the last 15 passes failed`); `lan-failed-pass` (from Thursday
9 November) pages for every failed pass at Lantern Atoll (`LAN: pass brine-3/LAN/15145
failed`); `disk`, every 10 minutes, pages while a station's disk is at least 98% full
(`LAN: disk 98% full`).

**The generator's silence**: Lantern Atoll, Thursday 11 January 2029 09:00 to 11:00.

#### Error budgets

Agreed on Monday 8 January 2029 by Ingrid, Tomás and Yusra (chapter 8).

- **Targets** are basis points of booked passes flown: 9950 (99.5%) at every station, 9900
  (99.0%) at Lantern Atoll, whose palms and downpours alone cost about one pass in a hundred.
  `type Target int`; `Standard Target = 9950`, `Atoll Target = 9900`.
- **Lost** is booked minus flown: failed and skipped passes both count ("every pass counts").
- **Availability** = flown / booked; nothing booked is an error, `no booked passes`; more flown
  than booked is `more flown than booked`.
- **Budget** = whole passes, in integers: `booked * (10000 - int(t)) / 10000`. Floats get it
  wrong: `int((1-0.9995)*10000)` is 4 at run time where 10,000 passes at 9995 allow 5.
- **Remaining** = budget − lost; negative means overspent.
- **Window**: a station's passes with LOS in `(now − 30×24h, now]`.
- **Burn rate** = (lost / booked) / ((10000 − target) / 10000); 0 when nothing was booked; with
  a target of 10000, infinite (`math.Inf(1)`) once anything is lost. `pass-failures` is a burn
  rate: two of Lantern Atoll's four passes in six hours burn at 50, a month's budget in fourteen
  hours.
- **At this rate**: `exhaustIn(remaining, lostLast7d)` is `remaining` days' worth of the last
  week's losses, `remaining × 7 days / lostLast7d`; 0 when nothing remains; no answer (false)
  when nothing was lost in the week.
- **Lantern Atoll, the 30 days to Monday 8 January 00:00**: booked 387, lost 11 (natural 13, 20,
  23, 26 and 27 Dec and 3 Jan; five on 27–28 Dec), budget 3, remaining −8. **Booked at Lantern Atoll**: every pass of two minutes or
  more, in the model, of Pipit-1, Furrow-A, Furrow-B, Brine-3 and Cirrine-1: Pipit-1's on the UHF
  yagis; the one S-band dish flies one pass at a time, Brine-3's first, then the pass that ends
  first. Pipit-2 and Pipit-3 never take the dish from a customer, and Brinewatch did not book
  Brine-3's 3-minute pass at 21:13 on 27 December (`brine-3/LAN/15146`). Elsewhere customers book
  only some of the passes the dishes could fly (Brinewatch 38 of Brine-3's 110 at Port Ellery,
  Salar Quieto and Cape Wren in its first exact week).

**The quarter, 1 October to 31 December 2028** (Yusra's report, Friday 12 January):

```
Q4 2028  booked  lost  avail  budget  left
YST        2208     2  99.91      11     9
PEL        1104     1  99.91       5     4
SQT        1288     3  99.77       6     3
LAN        1199    22  98.17      11   -11
CWR         920     4  99.57       4     0
```

Header `%-7s %7s %5s %6s %7s %5s`; rows `%-7s %7d %5d %6.2f %7d %5d`, availability in percent.
Cape Wren's four were all skipped in wind. Lantern Atoll's 1,199 are its booked passes with LOS in
the quarter; Salar Quieto's 1,288 fit under the 1,436 its dishes could fly. Lantern Atoll's 22: 15 natural, 7 to the leak (nobody
knows that yet).

#### passplan

passplan plans the day ahead each morning and after every stow, from Tomás's notebook model
(chapter 9). On Tuesday 16 January 2029 planning a day took 3 min 52 s on the planning machine
in the Fishbowl's rack; on Friday 19 January, 6.1 s. The fixes keep the plan identical.

```go
// Elements place a satellite on its orbit for the pass model.
type Elements struct {
	Epoch    time.Time // the moment the numbers below describe; no pass before it
	Rev0     int       // the orbit number at Epoch
	NodeDeg  float64   // east longitude of the northbound equator crossing, at Epoch
	PhaseDeg float64   // how far past that crossing the satellite is at Epoch, 0 to 360
	// RevsPerDay is the satellite's orbits a day, where tracking measured them; zero means
	// the period its altitude gives.
	RevsPerDay float64
}

// elevationDeg is the satellite's elevation over the station at t, in degrees (given, orbit.go).
func elevationDeg(sat Satellite, el Elements, st Station, t time.Time) float64

// orbitNumber is the number of the orbit on which the satellite, near t, crosses the station's
// latitude: Rev0 at Epoch, one more at each northbound equator crossing (given, orbit.go).
func orbitNumber(sat Satellite, el Elements, st Station, t time.Time) int

// Planned is one predicted pass in the day's plan.
type Planned struct {
	PassView        // ID and Pass (Season 3)
	Booking  string // "bk-5803", or "" when nobody booked it
}
```

**The model**, exactly (the one Season 3's reviewers settled; Season 7's `model.go` is the same
code): a = 6371 + AltKm km; the satellite moves along its orbit at n = √(398600.4418 / a³)
rad/s, or at RevsPerDay·2π/86400 rad/s where `RevsPerDay` isn't 0; the orbit's plane turns at
the mean Sun's rate, 2π/(365.2422·86400) rad/s, when IncDeg is over 90 (sun-synchronous), and
otherwise at −1.5·n·J2·(6371/a)²·cos i rad/s, with J2 = 1.08263e−3; Earth turns at
2π/86164.0905 rad/s. With dt the seconds since Epoch, u = PhaseDeg (as radians) + n·dt,
Ω = NodeDeg (as radians) + (the plane's rate − Earth's rate)·dt and i = IncDeg, the satellite is
at P = a(cos u cos Ω − sin u cos i sin Ω, cos u sin Ω + sin u cos i cos Ω, sin u sin i) km from
Earth's center, in a frame that turns with Earth; the station is on a sphere of 6371 km,
S = 6371(cos φ cos λ, cos φ sin λ, sin φ); elevation = asin((P − S)·S / (|P − S|·6371)). No
refraction. Brine-3 never rises over Ystmark in it. The orbit number near t is that of the
orbit on which the satellite crosses the station's latitude nearest t (the top or bottom of
its track, for a station its track never reaches). `orbit.go` writes the arithmetic as Season
7's `model.go` does, so the two give the same elevation to the last bit.

Lessons give this `orbit.go`, checked against Season 7's `model.go` on 200,000 random
elevations, none different:

```go
package main

import (
	"math"
	"time"
)

const (
	earthRadiusKm = 6371.0
	muKm3PerS2    = 398600.4418
	j2            = 1.08263e-3                       // how strongly Earth's bulge turns an orbit
	earthRate     = 2 * math.Pi / 86164.0905         // Earth's turn, radians a second
	sunRate       = 2 * math.Pi / (365.2422 * 86400) // the mean Sun's, radians a second
)

func rad(d float64) float64 { return d * math.Pi / 180 }
func deg(r float64) float64 { return r * 180 / math.Pi }

// skyTrack is one satellite over one station. Its times are seconds after the Epoch.
type skyTrack struct {
	a, cosI, sinI   float64    // the orbit's radius in km, and its tilt
	n, u0           float64    // radians a second along the orbit, and how far along at Epoch
	node0, nodeRate float64    // the node's longitude at Epoch, and how fast it moves over Earth
	st, up          [3]float64 // the station, km from Earth's center, and straight up there
}

func newSkyTrack(o Orbit, e Elements, s Station) *skyTrack {
	k := &skyTrack{a: earthRadiusKm + o.AltKm}
	i := rad(o.IncDeg)
	k.cosI, k.sinI = math.Cos(i), math.Sin(i)
	k.n = math.Sqrt(muKm3PerS2 / (k.a * k.a * k.a))
	if e.RevsPerDay > 0 {
		k.n = e.RevsPerDay * 2 * math.Pi / 86400
	}
	drift := -1.5 * k.n * j2 * (earthRadiusKm / k.a) * (earthRadiusKm / k.a) * k.cosI
	if o.IncDeg > 90 {
		drift = sunRate // sun-synchronous: the plane keeps its angle to the Sun
	}
	k.node0, k.nodeRate, k.u0 = rad(e.NodeDeg), drift-earthRate, rad(e.PhaseDeg)
	la, lo := rad(s.LatDeg), rad(s.LonDeg)
	k.up = [3]float64{math.Cos(la) * math.Cos(lo), math.Cos(la) * math.Sin(lo), math.Sin(la)}
	for j := range 3 {
		k.st[j] = earthRadiusKm * k.up[j]
	}
	return k
}

// pos is the satellite's place over the turning Earth, km from its center.
func (k *skyTrack) pos(t float64) (x, y, z float64) {
	su, cu := math.Sincos(k.u0 + k.n*t)
	sn, cn := math.Sincos(k.node0 + k.nodeRate*t)
	return k.a * (cn*cu - sn*su*k.cosI), k.a * (sn*cu + cn*su*k.cosI), k.a * su * k.sinI
}

// elev is the satellite's elevation over the station at t, in degrees.
func (k *skyTrack) elev(t float64) float64 {
	x, y, z := k.pos(t)
	dx, dy, dz := x-k.st[0], y-k.st[1], z-k.st[2]
	d := math.Sqrt(dx*dx + dy*dy + dz*dz)
	return deg(math.Asin((dx*k.up[0] + dy*k.up[1] + dz*k.up[2]) / d))
}

// orbit is the number of the orbit on which the satellite, near t, crosses the station's
// latitude: Rev0 at Epoch, one more at each northbound equator crossing. A satellite whose
// track never reaches the station's latitude counts from the top or bottom of its track.
func (k *skyTrack) orbit(rev0 int, t float64) int {
	u := k.u0 + k.n*t
	var marks []float64 // where along an orbit the satellite is at the station's latitude
	if s := k.up[2] / k.sinI; s >= 1 {
		marks = []float64{math.Pi / 2}
	} else if s <= -1 {
		marks = []float64{3 * math.Pi / 2}
	} else {
		marks = []float64{math.Asin(s), math.Pi - math.Asin(s)}
	}
	best := math.Inf(1)
	var at float64
	for _, m := range marks {
		c := m + 2*math.Pi*math.Round((u-m)/(2*math.Pi)) // the nearest such point to u
		if d := math.Abs(c - u); d < best {
			best, at = d, c
		}
	}
	return rev0 + int(math.Floor(at/(2*math.Pi)))
}

// elevationDeg is the satellite's elevation over the station at t, in degrees.
func elevationDeg(sat Satellite, el Elements, st Station, t time.Time) float64 {
	return newSkyTrack(sat.Orbit, el, st).elev(t.Sub(el.Epoch).Seconds())
}

// orbitNumber is the number of the orbit on which the satellite, near t, crosses the
// station's latitude.
func orbitNumber(sat Satellite, el Elements, st Station, t time.Time) int {
	return newSkyTrack(sat.Orbit, el, st).orbit(el.Rev0, t.Sub(el.Epoch).Seconds())
}
```

Keep the arithmetic exactly as written: the same formula written another way gives other
last bits in about a third of all elevations on arm64, because Go fuses multiplies and adds
differently when an expression is written differently. It changed no pass, but there is no
reason to risk it.

**Elements at Epoch 2029-01-16T00:00:00Z** (NodeDeg, PhaseDeg, Rev0; `RevsPerDay` 0 unless
given): pipit-1 201.4926001313, 322.4045275738, 14321; pipit-2 234.8439457617, 41.6844523733,
19979; pipit-3 149.9999975496, 93.2834560607, 9917; furrow-a 294.5055026568, 109.4557673805,
22046, `RevsPerDay` 867/57; furrow-b 297.5851207259, 287.9312984850, 22043, `RevsPerDay`
867/57; brine-3 94.5696529379, 61.6707270863, 15439; cirrine-1 282.5236519320, 176.5142936516,
17430. They are the settled fits carried to the epoch (Season 3's for Pipit-1, the Furrows,
Brine-3 and Cirrine-1; Season 4's for Pipit-2, fitted to the Build gctl night's Ystmark pass;
Pipit-3 from Apogee's own tracking, as in Season 7), so the plan's passes are exactly the
catalog's. Lantern Atoll gets about nineteen passes a day in it (19.0 from 1 October to 31
December 2028; 18.0 of 10° or more), two or three from each satellite. One `elevationDeg` call
costs under 100 ns on a laptop, so a test plans a few satellites over a few hours in
milliseconds.

**The plan's rules**: the elevation is sampled every second from `from` (included) to `to`
(excluded); where it crosses the station's mask between two samples, the crossing is bisected
to a millionth of a second or finer and rounded to the nearest second (half to even): AOS
where it rises, LOS where it sets; only passes that start and end inside the window; a pass
shorter than 2 minutes is too short to book and is left out; `MaxElevDeg` comes from the given
`peak`: the highest elevation at the unrounded AOS and at each whole second after it until
LOS, to a tenth as `strconv.FormatFloat(v, 'f', 1, 64)` rounds it (the same in every version);
the pass ID uses `orbitNumber` at the peak's instant; a booking matches by satellite, station
and AOS to the second (the first in input order when two match); the plan is sorted by AOS,
then ID.

**Where the four minutes went**, and the fixes (lessons 9.4 to 9.8): `elevationDeg` called every
second for every satellite and station (81% of the profile) → step 30 s instead of 1 s before
bisecting each crossing (the same AOS and LOS; about 28 times fewer calls); a linear scan of the next
fourteen days' bookings for every pass (9.5's test generates 2,000, which carry no IDs; the
real ones on Tuesday 16 January are `bk-5788` to `bk-5808`) → a map built once, keyed by
`bookingKey{SatID, Station string; AOSUnix int64}` (never a `time.Time` key: `==` compares its
location too); the whole day sorted after every insert → once at the end. The plan's pass-ID check
compiled its pattern on every call; Mirela asks about it in the quiz. Each lesson checks at least
5 times faster (measured 24×, 45×, 39× and 25× for 9.4, 9.5, 9.7 and 9.8), and 9.4 and 9.8 also
count at least 20 times fewer calls (measured 28× and 27×).

#### Traces and spans

**runtime/trace** (chapter 11): files in the lesson folder, `trace.out` from `trace.Start`, and
`flight.trace` from a `trace.FlightRecorder` snapshot (production: `MinAge` 1 minute, `MaxBytes`
16 MiB; tests `MinAge` 1 s). In uplinkd: a task `pass` per pass (`trace.NewTask(ctx, "pass")`,
ended with `defer`), `trace.Log(ctx, "pass", <pass ID>)`, and regions `wait for worker` and `fly`.
A trace file begins with `go 1.` and holds ` trace` in its header (Go 1.27.3 writes
`go 1.26 trace`); checks never depend on the version.

**Spans by hand**, in every uplinkd from Monday 19 February 2029:

```go
// Span is one timed step of work, recorded by hand.
type Span struct {
	TraceID  string    // 16 lowercase hex digits, the same for every span of one pass
	SpanID   string    // 8 lowercase hex digits
	ParentID string    // the enclosing span's SpanID; "" for a root
	Name     string    // "deliver", "queue", "pass", "worker", "fly"
	PassID   string    // from the context (chapter 3's passFrom), "" without one
	Start    time.Time
	End      time.Time
}

// Tracer starts spans and keeps the finished ones. Safe for many goroutines.
type Tracer struct {
	Now  func() time.Time
	Rand *rand.Rand // math/rand/v2: trace IDs aren't secrets; used only under mu
	mu   sync.Mutex
	done []Span
}

func (t *Tracer) Start(ctx context.Context, name string) (context.Context, *Span)
func (t *Tracer) End(s *Span)   // sets End from Now and keeps a copy
func (t *Tracer) Spans() []Span // the finished spans, in the order they ended
func spanFrom(ctx context.Context) (*Span, bool)
```

- `Start` takes the trace ID and parent from the span in `ctx`; without one it starts a trace,
  drawing the trace ID first (`fmt.Sprintf("%016x", Rand.Uint64())`); then every span draws its
  ID (`fmt.Sprintf("%08x", Rand.Uint32())`). With `rand.New(rand.NewPCG(1, 2))` the first root
  has trace `c4f5a58656eef510` and span `9dcec3ad`; the next two spans are `c8d04605` and
  `cbedc0dc`.
- **Across programs**: `Apogee-Trace: <trace id>-<span id>`, such as
  `Apogee-Trace: c4f5a58656eef510-9dcec3ad`, set by a `RoundTripper` from the context's span on a
  clone of the request. The receiving side puts a remote parent (a `*Span` with only `TraceID` and
  `SpanID`) in the request's context and starts its span as its child. A missing or malformed
  header (anything but 16 and 8 lowercase hex digits around one `-`) starts a new trace; it never
  refuses the request. Headquarters' `deliver` span is the parent of the station's `queue` span
  (`POST /uplink`); the station keeps the IDs with the queued lines and starts the pass's `pass`
  span as the `queue` span's child.
- **A finished span is one log line**: msg `span` at INFO with `trace_id`, `span_id`,
  `parent_id`, `name`, `start` and `ms`, and `pass_id` from the context.

**The bench replay** (Tuesday 13 February): the night of 6–7 February with the old `runPass`,
three workers held at the start, spans added, the restart at 02:10 included. Its spans, and the
incident review's timeline (lesson 11.8): each root by start time, its children below in the
order they started; root lines `%s %-24s %8s` (start as `15:04:05`, `pass <pass ID>`, duration),
child lines `%s   %-22s %8s`:

```
01:16:08 pass cirrine-1/LAN/17759    11m5s
01:16:08   worker                       0s
01:16:08   fly                      17m10s
01:36:52 pass furrow-b/LAN/22379     8m31s
01:36:52   worker                    8m31s
02:22:30 pass furrow-a/LAN/22382     3m12s
02:22:30   worker                       0s
02:22:30   fly                       3m12s
```

The `fly` span that outlives its pass by 6 min 5 s is the leak; `furrow-b/LAN/22379` has no `fly`
at all.

#### The code, by chapter

Each lesson is its own `package main` (module `lesson`). A lesson gives the earlier pieces it
needs as read-only files and uses only the fields it needs. `Station`, `Satellite`, `Orbit`,
`Pass` and `PassView`, `Booking`, `Key` (Season 5: `ID`, `Secret []byte`, `String()`),
`Severity` and `Shift` are unchanged.

| Ch | Names it adds (files) |
|---|---|
| 1 | `newLogger`, `logFlown`, `logResult`, `msAttr`, `accessLog`, `newPassLog`, `passLine` (`log.go`); given: `dropTime`, `stamp` (`clock.go`) |
| 2 | `levelHandler`, `programLogger`, `passGroup`, `passLogger`, `logFrame`, `levelFor` (`levels.go`) |
| 3 | `passKey`, `withPass`, `passFrom`, `passHandler`, `passTransport`, `acceptPass`, `Line`, `follow`, `merge`, `programLog` (`context.go`, `handler.go`, `follow.go`, `line.go`); given: Season 3's `requestIDKey` and `requestIDFrom` (`request.go`), `replayNovember8` (`replay.go`) |
| 4 | `safeHeader`, `APIKey`, `(APIKey).LogValue`, `(Key).LogValue`, `redact`, `maskEmail`, `secretPattern`, `secrets`, `scrub`, `safeLogger` (`redact.go`) |
| 5 | `openLog`, `countingWriter`, `Rotator`, `prune`, `Sampler`, `NewSampler`, `(*Sampler).Dropped`, `newStationLog` (`files.go`, `rotate.go`, `sample.go`); given: `rotatedName`, `parseRotated`, `rotated`, `rotatedFiles`, `fileBytes`, `totalBytes`, `keepDays` (`files.go`), `(*Rotator).open` and `OnOpen`, and the sampler's `sampling`, `tallyKey`, `tally`, `tallyOf` |
| 6 | `Counter`, `Gauge`, `track`, `CounterVec`, `Histogram`, `NewHistogram`, `HistogramSnapshot`, `Metrics`, `publish`, `Reading`, `increase` (`metrics.go`, `track.go`, `vec.go`, `histogram.go`, `publish.go`, `readings.go`) |
| 7 | `Result`, `failedWithin`, `Rule`, `Sample`, `Alert`, `Silence`, `Pager`, `ruleState` (`since`, `firing`, `sent`; `held` from 7.7), `beforeFiring`, `whileFiring`, `route`, `silenced`, `newRules` (`results.go`, `rules.go`, `state.go`, `pager.go`, `route.go`, `silence.go`, `rota.go`); given: `replay` (`replay.go`) |
| 8 | `Target`, `availability`, `budget`, `window`, `burnRate`, `exhaustIn`, `row` (`budget.go`) |
| 9 | `profileTo`, `benchPlan`, `contact`, `findPasses`, `bookingKey`, `attach`, `Planned`, `buildDay`, `planDay` (`plan.go`, `bench.go`); given: `Elements` (`fleet.go`), `elevationDeg`, `orbitNumber` (`orbit.go`), `elevation`, `crossing`, `toSecond`, `scanEverySecond`, `peak`, `passViews`, `scanAllBookings`, `earlier`, `sortAfterEach` (`passplan.go`) |
| 10 | `heapProfile`, `inUse`, `appendPassID`, `tail`, `leftBehind`, `Report`, `Uplink`, `runPass`, `ErrNoAck`, `ErrNoWorker` (`memory.go`, `uplink.go`) |
| 11 | `traceTo`, `Span`, `Tracer`, `spanFrom`, `traceTransport`, `continueTrace`, `watchPass`, `timeline` (`trace.go`, `span.go`, `review.go`) |

```go
// Result is how one pass ended, as a station's uplinkd reports it to the pager.
type Result struct {
	PassID  string    // "cirrine-1/LAN/16402"
	Station string    // "LAN"
	LOS     time.Time // a result counts from its LOS
	Outcome string    // "flown", "failed" or "skipped"
	Reason  string    // "" when flown; "deadline", "no_ack", "no_worker", "radio"; "stowed" when skipped
}

// Line is one JSON log line read back (chapter 3).
type Line struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Msg     string    `json:"msg"`
	Program string    `json:"program"`
	Site    string    `json:"site"`
	PassID  string    `json:"pass_id"`
	Raw     string    `json:"-"` // the line as written
}

// Report is what the radio heard during a pass (chapter 10).
type Report struct {
	Frames []byte    // the pass's frames, in the buffer runPass gave fly
	AckAt  time.Time // the acknowledgment of the last line; zero without one
}

// Uplink is a station's uplinkd as chapter 10 sees it.
type Uplink struct {
	Station string
	workers chan struct{} // room for 4: send to take a worker, receive to give it back
	fly     func(ctx context.Context, p PassView, lines []string, buf []byte) (Report, error) // before 10.8: no ctx
}

// Rotator writes a program's log file and rotates it (chapter 5).
type Rotator struct {
	Dir      string
	Program  string // "uplinkd": the current file is uplinkd.log
	MaxBytes int64
	Now      func() time.Time
	OnOpen   func() // runs each time the current file is opened: at start and after every rotation
	f        *os.File
	size     int64
	day      string // the current file's UTC day, "2028-12-04"
	n        int    // that day's rotated files so far
}
```

Signatures every chapter after the one that defines them uses exactly:

```go
func dropTime(groups []string, a slog.Attr) slog.Attr
func stamp(h slog.Handler, now func() time.Time) slog.Handler // keeps itself through WithAttrs and WithGroup
func newLogger(w io.Writer) *slog.Logger
func logFlown(l *slog.Logger, p PassView, lines int)
func logResult(l *slog.Logger, id string, err error)
func msAttr(key string, d time.Duration) slog.Attr // ("wait", d) → wait_ms
func accessLog(h slog.Handler) *log.Logger
func newPassLog(w io.Writer, now func() time.Time) *slog.Logger
func passLine(l *slog.Logger, p PassView, wait time.Duration, err error)
func levelHandler(v *slog.LevelVar) http.Handler
func programLogger(base *slog.Logger, program, site string) *slog.Logger
func passGroup(p PassView) slog.Attr
func passLogger(l *slog.Logger, p PassView) *slog.Logger // pass_id and the group in ch 2; the group only from ch 3
func logFrame(ctx context.Context, l *slog.Logger, frame []byte, dump func([]byte) string)
func levelFor(event string) slog.Level
func withPass(ctx context.Context, id string) context.Context
func passFrom(ctx context.Context) (string, bool)
func acceptPass(next http.Handler) http.Handler
func follow(r io.Reader, passID string) (lines []Line, broken int, err error)
func merge(logs ...[]Line) []Line
func programLog(w io.Writer, program, site string, now func() time.Time) *slog.Logger
func safeHeader(h http.Header) http.Header
func redact(groups []string, a slog.Attr) slog.Attr
func maskEmail(s string) string
func scrub(s string) string
func safeLogger(w io.Writer) *slog.Logger
func openLog(path string) (*os.File, error)
func (r *Rotator) Write(p []byte) (int, error)
func (r *Rotator) Close() error
func prune(dir, program string, now time.Time, keepDays int, maxTotal int64) error
func NewSampler(inner slog.Handler) *Sampler
func newStationLog(dir string, now func() time.Time) (*slog.Logger, io.Closer, error)
func track(g *Gauge, run func() error) error
func publish(m *Metrics) // once per test binary: expvar names can't be published twice
func increase(rs []Reading) uint64
func failedWithin(results []Result, station string, now time.Time, d time.Duration) int
func route(a Alert, sev Severity, rota []Shift) Alert
func silenced(s []Silence, station string, at time.Time) bool
func newRules() []Rule
func availability(flown, booked int) (float64, error)
func budget(t Target, booked int) int
func window(results []Result, station string, now time.Time) (booked, lost int)
func burnRate(t Target, booked, lost int) float64
func exhaustIn(remaining, lostLast7d int) (time.Duration, bool)
func row(code string, booked, lost int, t Target) string
func profileTo(path string, work func()) error
func benchPlan(b *testing.B)
func findPasses(el func(time.Time) float64, from, to time.Time, maskDeg float64) []contact // contact{AOS, LOS}: the crossings, before rounding
func peak(el func(time.Time) float64, aos, los time.Time) (float64, time.Time) // MaxElevDeg and the moment of it
func attach(passes []PassView, bookings []Booking) []Planned
func buildDay(groups [][]Planned) []Planned
func planDay(sats []Satellite, els map[string]Elements, stations []Station, bookings []Booking, from, to time.Time) []Planned
func heapProfile(path string) error
func inUse() uint64
func appendPassID(buf []byte, sat, station string, orbit int) []byte
func tail(frames []byte, n int) []byte
func leftBehind(f func()) int
func newUplink(station string, fly func(context.Context, PassView, []string, []byte) (Report, error)) *Uplink // 10.6 and 10.7: fly without ctx
func (u *Uplink) runPass(ctx context.Context, p PassView, lines []string) error
func traceTo(path string, work func()) error
func continueTrace(t *Tracer, name string, next http.Handler) http.Handler
func watchPass(fr *trace.FlightRecorder, path string, run func() error) error
func timeline(spans []Span) []string
```

#### How lessons test it

- **Clocks** are `now func() time.Time` from a small fake clock with `Now` and `Add` (Season 3's,
  under a mutex when goroutines read it); log records get the story's time through `stamp`.
  Randomness is `rand.New(rand.NewPCG(1, 2))`. No test reads the wall clock for a value it
  checks.
- **Log output** goes to a `bytes.Buffer`. Tests decode each JSON line (`json.Unmarshal` into a
  `map[string]any` or `Line`) and compare fields, never whole lines, except where a key's order or
  count is the point (`pass_id` exactly once, counted in the raw line). Handlers are tested
  directly with `slog.NewRecord` and `Handle` where the record's time matters;
  `testing/slogtest.TestHandler` checks a handler of your own. `code` lessons print time-free
  text lines of at most 50 characters.
- **Secrets** are searched for in the whole output (the 32 hex digits, the base64, the password,
  the full address), so a secret anywhere fails the test.
- **Files** live under `t.TempDir()`, with limits of 1 KiB or 64 KiB so rotation happens in a few
  hundred lines. **Profiles and traces** are written to the lesson folder (`cpu.prof`,
  `cpu-after.prof`, `heap.prof`, `trace.out`, `flight.trace`); a check looks only at the file
  being there, not empty, and starting with its format's bytes (`1f 8b` for a profile, `go 1.`
  for a trace), and that the profile or trace was stopped (a second `StartCPUProfile` or
  `trace.Start` on `io.Discard` succeeds). Readmes may offer `go tool pprof -top cpu.prof`,
  `go tool pprof -sample_index=inuse_space heap.prof` or `go tool trace trace.out` in a
  terminal; no check depends on them.
- **Speed** is measured inside a test: the test first sets `flag.Set("test.benchtime",
  "100ms")`, then compares two `testing.Benchmark` results on the same Mac. Checks are relations
  with wide margins, never times: require at most a fifth of the speedup measured on the
  author's Mac (5× where 25× was measured). Where calls can be counted, count them (exact).
  **Allocations** are `testing.AllocsPerRun(100, f)`, which also holds under `-race`; never
  `AllocsPerOp` from a benchmark, which `-race` changes.
- **Memory** is `runtime.GC()` then `runtime.ReadMemStats` (`HeapAlloc`), compared with margins
  (at least 12 MiB more against under 4 MiB), using at most 32 MiB.
- **Goroutines**: count with `runtime.NumGoroutine()` before; release the fake radio (a channel
  the test closes) only after `runPass` has returned; then wait, polling every millisecond for at
  most one second, for the count to come back. A message reads
  `goroutines after 20 failed passes = 34, want 14`. A starter's stuck goroutines stay stuck
  until the test binary ends, so every test compares against its own starting count. The leak
  profile is read only for goroutines a test planted itself, polled until it reports them.
  Never `synctest` for leaks: a bubble with a stuck goroutine panics, and starters must fail,
  not panic.
- **Concurrency** tests are correct by construction (a `WaitGroup`, then exact totals; values
  checked from inside the work, as in a gauge read during `run`); sleeps are 50 ms at most.
- **Servers** are `httptest` on 127.0.0.1, clients `srv.Client()` with a timeout of at most two
  seconds; `expvar.Handler()` is called directly or through `httptest`; `expvar` names are
  published once per test binary through the given `sync.Once`.
- **Reading the code**: where right and wrong code behave alike (a pattern compiled once, a task
  ended with `defer`, an `atomic.Uint64` field), the test parses the learner's file with
  `go/parser` and says what it looked for, as Seasons 1 and 5 did.
- Messages read `call = got, want want`: `budget(9995, 10000) = 4, want 5`. Every run ends well
  under the app's 10 seconds; benchmark lessons take about a second.

**Topics added this season** (front matter): `logging`, `metrics`, `profiling`, `tracing`.
Lessons keep using the older ones (`context`, `errors`, `files`, `goroutines`, `http`, `json`,
`mutexes`, `security`, `testing`, `time`...).

**New field guide words** (world words, in plain language): Ground Reliability Engineer, the
rota, the night shift, the morning list, runbook, the Monday restart, the debug port, the
caretaker, the generator, downpour, log line, log level, structured log, redaction, personal data,
log rotation, sampling, metric, counter, gauge, histogram, alert, threshold, flapping,
hysteresis, silence, availability target, error budget, burn rate, profile, hot path, benchmark,
heap, leak, stuck goroutine, worker, trace, span, flight recorder. Carried over and needed again:
on-call, page, pager, incident review, pass, pass ID, AOS, LOS, UTC, stow, rain fade, the
Fishbowl, the Board, the ops console, uplinkd, bookingd, passplan, sealed command, API key.

**Name checks** to run before the season ships: the job title "Ground Reliability Engineer"; the
header names `Apogee-Pass` and `Apogee-Trace`, in Season 3's `Apogee-` style; "the morning list"
and "the Monday restart" as phrases; the contact addresses use the reserved `.example` domain.
No new people, places, satellites or customers. Booking IDs this season run from `bk-5600` to
`bk-5899`, above Season 5's `bk-5274` and Season 3's example `bk-5499`.
Name checks (web search, 10 November 2026): nothing real was found for the badge "Ground
Reliability Engineer" (only generic postings such as "Reliability Engineer, Ground Systems") or
for the headers `Apogee-Pass` and `Apogee-Trace`.

### Season 7: scale

On Monday 5 March 2029 Tomás tells the building that investors have put money into Apogee,
and draws the network it pays for on the Fishbowl's glass: fifteen new stations by June, twenty
in all, and the satellites of four customers' launches, eighty in all by the end of June. The
learner's new badge says **Senior Planning Engineer**, still No. 039, based at Port Ellery
headquarters, since Monday 5 March 2029. Their job is passplan's rewrite. The passplan Apogee
has flown since 2023, which Nkem wrote in Apogee's first summer, plans one day of passes at a
time, in one goroutine: 6.1 seconds for Season 6's seven satellites at five stations, 7 seconds
for eight. For eighty satellites at twenty stations it needs 4 minutes 39 seconds a day, nearly
eight hours for a hundred days, and the new one must plan a hundred days, correctly and fast. You lead the rewrite,
**passplan 2**; Nkem reviews every change and asks, one question at a time, what happens when
two planners touch the same pass. The season ends on take-over day, Thursday 28 June, when
passplan 2 takes over from the old planner on a live day, under the heaviest schedule Apogee
has ever flown.

Story time runs from Monday 5 March to Thursday 28 June 2029: twelve chapters, nine of about
a week and three (chapters 8, 9 and 11) of about two; chapter 12 runs from Monday 11 June to
take-over day. Nkem carries the season; Tomás draws the network and makes one claim too many
about speed; Ingrid owns take-over day; Yusra brings the customers' wishes and Hollis his
capitals; Mirela reviews the telemetry parser; Joss Imbert and Professor Marr of Kessick
Polytechnic return with Thistle-2. No new people are introduced. The new places are fifteen
stations, each an Apogee dish on a host's land; the new satellites are Thistle-2 and three
customers' constellations, named by rule from names already in the bible. The hosts stay
unnamed, as the Cape Wren contractor does.

In plain words, for the opening: Apogee is growing from five dishes to twenty and from a
handful of satellites to eighty, and the program that decides which dish listens to which
satellite, and when, has to be rebuilt for that size without missing a pass.

#### The season, chapter by chapter

Everything happens at Port Ellery headquarters; the new stations speak through their hosts'
forms, and Sunniva and Paz only in passing. Times are UTC.

| Ch | Week of (2029) | Where | Who needs what | What is at stake |
|---|---|---|---|---|
| 1 | Mon 5 Mar | the Fishbowl's glass, then the backend desks | Nkem: before passplan 2, a smaller pipeline built the same way. tlmd's ingest reads the station telemetry feeds in one loop; with twenty feeds it falls behind, and when two stations hear the same satellite its frames arrive out of order, so "the newest frame" goes backward. Wednesday 7 March, from 21:00, is Nkem's drill: twenty feeds, the five real stations and fifteen made from the new stations' positions, replaying the day's lines. Thistle-2 leaves its rocket at 10:20 on Thursday 8 March | The Board and the pager read every satellite's newest frame from this feed. Thistle-2's first beacon reaches Port Ellery at 17:08 on Thursday (`thistle-2/PEL/00005`), and Joss and Professor Marr are in the Fishbowl to hear it |
| 2 | Mon 12 Mar | the backend desks; passplan's new machine in Port Ellery's rack, twelve cores | Nkem: passplan 2's first draft predicts each satellite in its own goroutine, and each day in another: eight thousand at once. On Friday night, 9 March, the test machine ran out of memory. "Twelve planner slots. Never more." The station dish loaders and the Board take what passplan sends at their own pace | A planner that starves its machine stops the Board that shares it; a queue nobody bounds grows until the machine falls over |
| 3 | Mon 19 Mar | Ingrid's desk in the Fishbowl, then the backend desks | Ingrid: on Monday the Board froze for six minutes while a hundred callers each loaded the station table. On Wednesday 21 March two of passplan 2's planning runs finished out of order, and on the test Board beside Ingrid's the older plan replaced the newer one | The test Board showed a plan twenty minutes old; on take-over day that would be the real one. Nkem's whiteboard question: "Why does a program without a race behave?" |
| 4 | Mon 26 Mar | Yusra's desk | Yusra: "TL;DR: customers want to book a whole season ahead." Hollis: "WE PLAN A SEASON NOT A WEEK." bookingd's catalog offers fourteen days today; the third quarter (1 July to 30 September) holds 672,292 passes at stations that can serve them | The old passplan hands every service a slice of every pass: about eighty-four megabytes for each one that asks |
| 5 | Mon 2 Apr | the Fishbowl | The first wave of new stations goes live on Monday 2 April: Tirrin Rock, Varsk Island, Calidra Atoll, Ilha Brumosa, Ostrel Island. From the same day passplan 2 runs in shadow, beside the old planner, and the two schedules are compared contact by contact every hour. The Board wants one station's contacts, bookingd one customer's, Ingrid the next twelve | Two planners that disagree about one pass would fly it twice or not at all; the first wave's schedules must be right from their first pass. Sunniva: "Tirrin up. Two dishes. Ice on both." |
| 6 | Mon 9 Apr | the backend desks | Nkem: the shadow run sorts every candidate pass of a hundred days into one slice per station. A priority queue merges the eighty satellites' passes one at a time | Cirrine's twelve new satellites leave their rocket on Tuesday 17 April; the planner must be ready for them |
| 7 | Mon 16 Apr | Mirela's desk; Kessick Polytechnic by message | On Saturday 14 April Thistle-2's student flight software sent `batt_v=nan` when a sensor glitched; tlmd's parser accepted it, and the pager's battery rule never fired, since nothing compares below a value that isn't a number. At 05:10 a line with a missing field crashed the parser and stopped every station's telemetry for forty minutes. Joss apologizes; Mirela: "Your parser accepts the line. Now show me the lines it refuses" | Cirrine's twelve send their first telemetry through this parser on Tuesday 17 April |
| 8 | Mon 23 Apr, to Fri 4 May | the backend desks | Nkem: passplan 2's timers (a replan every ten minutes, a deadline on every planning run, a wait until each contact's AOS, a pause for news to settle) are tested with real sleeps: the suite takes eleven minutes and fails one run in five | The second wave goes live on Monday 7 May; a test that lies about time would let a broken replan through |
| 9 | Mon 7 May, to Fri 18 May | the Fishbowl's glass, the backend desks | The second wave goes live on Monday 7 May: Ilvane Spit, Cape Solvane, Saltreach Island, Corith Point, Nimrie Point. Tomás writes on the glass: "passplan 2: forty times faster!" Nkem: "Faster at what?" Brinewatch's thirty-six leave their rocket on Wednesday 16 May | A planner that only looks faster would be trusted on take-over day |
| 10 | Mon 21 May | the backend desks | With fifty-six satellites, every replan of a hundred days leaves the machine's memory full of garbage, and the test Board's updates stutter while the collector cleans up. Nkem: "Halve the planner's garbage before the third wave" | The third wave adds five stations and the last launch adds twenty-four satellites |
| 11 | Mon 28 May, to Fri 8 Jun | the backend desks | The third wave's hosts send their forms: Vorrel Island's says 0 dishes, Ilha Calmaria's lists a band "K". On Friday 1 June a test deploy of passplan 2 read `slots: 0` and planned nothing. Nkem: one validator for every settings struct, its rules in tags | The third wave goes live on Monday 4 June: Hasku Point, Ilha Calmaria, Quindle Island, Punta Ventisca, Vorrel Island. A wrong form takes a station's passes away |
| 12 | Mon 11 Jun, to Thu 28 Jun | the Fishbowl | Ingrid runs take-over day like a pass: a full rehearsal in shadow on Monday 25 June, Furrowlight's launch on Wednesday 27 June, and the take-over at 00:00 UTC on Thursday 28 June | The heaviest schedule Apogee has ever flown: 2,011 contacts in one day at twenty stations, with twenty-four new Furrows in their first day of life |

The opening (Monday 5 March) shows Tomás's drawing on the glass (five dots becoming twenty,
eight satellites becoming eighty), hands over the badge, introduces Nkem (reviewer, author of
the old passplan), Tomás, Ingrid (take-over day is hers) and Yusra, and sends the learner to
Nkem's first question about tlmd's ingest. Words of the first week: passplan, planner slot,
telemetry feed, the drill.

The finale (chapter 12, Thursday 28 June): passplan 2 takes over at midnight without touching
one of the old planner's twelve contacts in progress; flies 2,011 contacts through Cape Wren's
wind; and every one of Furrowlight's twenty-four new satellites, first heard the evening
before, is heard all day. Ingrid writes one line in the log at 00:00 on Friday: "passplan 2:
GO for every day." Hollis, through Yusra: "TWENTY-FOUR FURROWS. EVERY ONE HEARD." The debrief
points at July: the whole ground segment, rehearsed live for a week.

#### The network grows

Apogee's network at the end of June 2029: the five stations of the bible's table and fifteen
new ones, each an Apogee dish on a **host's** land (a lighthouse authority, a weather
observatory, a fisheries harbor, a research camp): the host gives power, a line to the internet
and a pair of hands; nobody from Apogee lives there. Each new station arrives as a dish and a
rack in a shipping container and goes live in one of three **waves**, at 00:00 UTC on a
Monday. Coordinates are invented and land in open ocean, as the bible's do. A new station has
one 3.7 m dish for S and X unless the table says otherwise; X-band is left out where tropical
rain would fade it, as at Lantern Atoll. Only Port Ellery and Lantern Atoll send on UHF.

| Code | Name | Lat, lon | Mask | Dishes | Bands | Live | Character |
|---|---|---|---|---|---|---|---|
| YST | Ystmark | 77.8° N, 2.0° E | 5° | 2 | S X | flying | the bible's |
| PEL | Port Ellery | 46.9° N, 38.1° W | 5° | 1 | UHF S | flying | the bible's |
| SQT | Salar Quieto | 23.4° S, 104.6° W | 3° | 2 | S X | flying | the bible's |
| LAN | Lantern Atoll | 1.8° S, 82.5° E | 8° | 1 | UHF S | flying | the bible's |
| CWR | Cape Wren | 53.9° S, 95.0° W | 5° | 1 | S X | flying | the bible's |
| TIR | Tirrin Rock | 69.5° N, 12.0° W | 5° | 2 | S X | Mon 2 Apr | A bare rock in the Norwegian Sea with a lighthouse and two radomes; spray ices them in winter. Ystmark's partner for polar passes |
| VSK | Varsk Island | 51.0° N, 160.0° W | 7° | 1 | S X | Mon 2 Apr | A North Pacific fishing harbor; fog most mornings, and a hill to the north |
| CAL | Calidra Atoll | 0.5° N, 150.0° W | 8° | 1 | S | Mon 2 Apr | An equatorial atoll; palms on the horizon and rain fade, as at Lantern Atoll |
| BRM | Ilha Brumosa | 30.0° S, 15.0° W | 5° | 1 | S X | Mon 2 Apr | A misty South Atlantic island ("brumosa" means misty); its host is a weather observatory |
| OST | Ostrel Island | 42.5° S, 70.0° E | 5° | 1 | S X | Mon 2 Apr | In the westerly gales of the southern Indian Ocean; the dish stows at 72 km/h, like Cape Wren's |
| ILV | Ilvane Spit | 72.0° N, 165.0° W | 3° | 2 | S X | Mon 7 May | A gravel spit in the Chukchi Sea, flat to every horizon; sea ice until June |
| CSL | Cape Solvane | 35.0° N, 45.0° W | 5° | 1 | S X | Mon 7 May | A volcanic cape in the mid-Atlantic; its host keeps vines on the slope below the dish |
| SRH | Saltreach Island | 17.0° N, 135.0° E | 6° | 1 | S | Mon 7 May | Typhoon country in the western Pacific |
| CRP | Corith Point | 10.0° N, 65.0° E | 6° | 1 | S | Mon 7 May | A point of an Arabian Sea island; the monsoon arrives in June |
| NIM | Nimrie Point | 45.0° S, 150.0° W | 5° | 1 | S X | Mon 7 May | A lonely point in the South Pacific; the supply ship calls twice a year |
| HSK | Hasku Point | 32.0° N, 170.0° E | 4° | 1 | S X | Mon 4 Jun | An old airstrip on a North Pacific island |
| ICM | Ilha Calmaria | 8.0° N, 30.0° W | 7° | 1 | S | Mon 4 Jun | In the calm belt of the tropical Atlantic ("calmaria" means calm), with sudden squalls |
| QND | Quindle Island | 25.0° S, 95.0° E | 5° | 1 | S X | Mon 4 Jun | An island of the eastern Indian Ocean, in the trade winds |
| PVT | Punta Ventisca | 60.5° S, 40.0° E | 5° | 1 | S X | Mon 4 Jun | A cape in the Southern Ocean ("ventisca" means blizzard); the radome heater never stops |
| VRL | Vorrel Island | 64.0° S, 140.0° W | 4° | 2 | S X | Mon 4 Jun | An icebound island in the far south: the south's Ystmark |

The table order (the five, then the waves in this order) is the **table order** every list
of stations uses. In readmes a new station is its name first and its code after: "Tirrin Rock
(TIR)".

**The fleet** grows from eight to eighty: the seven the bible's fleet still has (Thistle
re-entered in July 2027), Thistle-2 on Thursday 8 March, then three launches. Satellites of a
launch leave the rocket ("deploy") one after another and fly in a line around the same orbit;
a new satellite's first days are its LEOP. IDs and names follow the bible's rule: `cirrine-2`
is Cirrine-2, `brine-17` is Brine-17, `furrow-q` is Furrow-Q.

| Satellites | Operator | Orbit | Bands | Deployed (Epoch below) | Story |
|---|---|---|---|---|---|
| `pipit-1`, `pipit-2`, `pipit-3`, `furrow-a`, `furrow-b`, `brine-3`, `cirrine-1` | as in the bible | as in the bible | Pipit-1 UHF; Pipit-2 S; Pipit-3 S X; the Furrows X S; Brine-3 S; Cirrine-1 S | flying | the bible's |
| `thistle-2` | Kessick Polytechnic | 505 km, 97.4° | UHF | Thu 8 Mar, 10:20:00 | Kessick's second student cubesat; Joss now leads its flight software. First beacon at Port Ellery, 17:08:39 (`thistle-2/PEL/00005`) |
| `cirrine-2` to `cirrine-13` (12) | Cirrine Atmospherics | 600 km, 97.8° | S | Tue 17 Apr, 09:40:00 | the bigger contract of Season 5; first contact `cirrine-3/VSK/00001` at 09:41:01 |
| `brine-4` to `brine-39` (36) | Brinewatch | 480 km, 53.0° | S | Wed 16 May, 13:30:00 | ship tracking; like Brine-3, never over Ystmark |
| `furrow-c` to `furrow-z` (24) | Furrowlight | 510 km, 97.4° | X S | Wed 27 Jun, 20:40:00 | lift-off 18:05, deployed one about every four minutes from 19:10; first contact `furrow-g/YST/00001` at 20:40:16, the last of the 24 first heard at 22:28:50 (`furrow-c/TIR/00002`) |

The fleet order (the seven, Thistle-2, then the three launches in this order, each by
number or letter) is the **fleet order**.

#### The pass model

passplan predicts passes from each satellite's **elements**, the few numbers that place it on
its orbit (`Elements`, below: Season 6's type, to which this season adds `SatID` and `Orbit`). The model is written here once and copied unchanged into every
lesson that predicts a pass, as `network.go` (the twenty stations and eighty satellites) and
`model.go` (`predict`); `go run ./cmd/passes` runs the same code. It is the model Season 3's
reviewers settled, ported to Go: a circular orbit around a round Earth, the orbit's plane turning with the mean
Sun for a sun-synchronous satellite (one tilted past 90°) and by Earth's bulge for any other,
Earth turning beneath it, and a station's mask. It watches the satellite's elevation every 20
seconds (less often while the satellite is far below the horizon), finds each instant it
crosses the mask to a hundred-millionth of a second and rounds it to the second, and takes the
highest of the elevations a second apart from AOS as the maximum. A pass shorter than two
minutes doesn't count. Near the top of its track a satellite crosses a high station's latitude
twice on one orbit; the station sees it rise once and set once: one pass. Pass IDs follow the
bible's rule: the orbit is the one on which the satellite crosses the station's latitude during
the pass (the top or bottom of its track where it never reaches that latitude), counted from
`Rev0` at `Epoch`; a new satellite's orbits count from 1 at its deployment. (Season 3's first
version numbered a pass by the orbit at its AOS, which gave two passes one ID at stations just
north of the equator; the rest agree.) Elements of the eighty:

| Satellites | Epoch | Rev0 | NodeDeg | PhaseDeg |
|---|---|---|---|---|
| `pipit-1` | 2029-01-01T00:00:00Z | 14094 | 201.4926025817 | 176.9565679403 |
| `pipit-2` | 2029-01-01T00:00:00Z | 19752 | 234.8439482122 | 339.4065995135 |
| `pipit-3` | 2029-01-01T00:00:00Z | 9689 | 150 | 40 |
| `furrow-a` | 2029-01-01T00:00:00Z | 21818 | 294.5055051072 | 52.6136621172 |
| `furrow-b` | 2029-01-01T00:00:00Z | 21815 | 297.5851231763 | 231.0891932221 |
| `brine-3` | 2029-01-01T00:00:00Z | 15209 | 179.2270109042 | 188.3141532885 |
| `cirrine-1` | 2029-01-01T00:00:00Z | 17206 | 282.5236543824 | 268.6744811484 |
| `thistle-2` | 2029-03-08T10:20:00Z | 1 | −120 | 0 |
| `cirrine-n`, n = 2 to 13 | 2029-04-17T09:40:00Z | 1 | −150 | 30 × (n − 2) |
| `brine-n`, n = 4 to 39 | 2029-05-16T13:30:00Z | 1 | 120 | 10 × (n − 4) |
| `furrow-c` … `furrow-z`, k = 0 to 23 | 2029-06-27T20:40:00Z | 1 | 10 | 15 × k |

Pipit-1, Furrow-A, Furrow-B, Brine-3 and Cirrine-1 are Season 3's settled fits, and Pipit-2 is
Season 4's, fitted to the Build gctl night's pass over Ystmark (`pipit-2/YST/13287`, 04:31:00 to
04:40:33 on 30 October 2027, 47.2°) with its orbits counted from 19752 at 1 January 2029; all six
are carried to 1 January 2029 (Furrow-A and Furrow-B with `RevsPerDay` 867/57: 867 orbits in 57
days, a little faster than their altitude gives); every other satellite has `RevsPerDay` 0. So the
orbit numbers continue the earlier seasons', and Furrow-A still rises over Salar Quieto at
14:02:10 on 28 June (`furrow-a/SQT/24534`, to 14:12:04, 47.9°), as on launch day in 2028, though
that day both of the station's dishes are taken by then. Pipit-1 crosses Port Ellery at 03:39:12
and 05:13:13 on 28 June (`pipit-1/PEL/16795`, `/16796`), each time a few minutes after Thistle-2,
which gets the station's one dish.

```go
// Site is a station as passplan sees it.
type Site struct {
	Station
	Dishes int       // the dishes passplan can give passes to
	Live   time.Time // the first moment it flies a pass; zero means always
}

// Elements place a satellite on its orbit for the pass model.
type Elements struct {
	SatID    string
	Orbit    Orbit     // altitude and inclination
	Epoch    time.Time // the moment the numbers below describe; no pass before it
	Rev0     int       // the orbit number at Epoch
	NodeDeg  float64   // east longitude of the northbound equator crossing, at Epoch
	PhaseDeg float64   // how far past that crossing the satellite is at Epoch, 0 to 360
	// RevsPerDay is the satellite's orbits a day, where tracking measured them; zero means
	// the period its altitude gives.
	RevsPerDay float64
}

// Network is everything passplan plans for.
type Network struct {
	Sites      []Site      // the twenty, in table order
	Satellites []Satellite // the eighty, in fleet order
	Elements   []Elements  // the same eighty, in the same order
}
```

`network.go`, exactly (it needs `fmt` and `time`):

```go
// newNetwork returns the network as it stands at the end of June 2029.
func newNetwork() *Network {
	sx, s, us := []Band{SBand, XBand}, []Band{SBand}, []Band{UHF, SBand}
	w1 := time.Date(2029, 4, 2, 0, 0, 0, 0, time.UTC)
	w2 := time.Date(2029, 5, 7, 0, 0, 0, 0, time.UTC)
	w3 := time.Date(2029, 6, 4, 0, 0, 0, 0, time.UTC)
	site := func(code, name string, lat, lon, mask float64, bands []Band, dishes int, live time.Time) Site {
		st := Station{Code: code, Name: name, LatDeg: lat, LonDeg: lon, MaskDeg: mask, Bands: bands}
		return Site{Station: st, Dishes: dishes, Live: live}
	}
	var always time.Time
	n := &Network{Sites: []Site{
		site("YST", "Ystmark", 77.8, 2.0, 5, sx, 2, always),
		site("PEL", "Port Ellery", 46.9, -38.1, 5, us, 1, always),
		site("SQT", "Salar Quieto", -23.4, -104.6, 3, sx, 2, always),
		site("LAN", "Lantern Atoll", -1.8, 82.5, 8, us, 1, always),
		site("CWR", "Cape Wren", -53.9, -95.0, 5, sx, 1, always),
		site("TIR", "Tirrin Rock", 69.5, -12.0, 5, sx, 2, w1),
		site("VSK", "Varsk Island", 51.0, -160.0, 7, sx, 1, w1),
		site("CAL", "Calidra Atoll", 0.5, -150.0, 8, s, 1, w1),
		site("BRM", "Ilha Brumosa", -30.0, -15.0, 5, sx, 1, w1),
		site("OST", "Ostrel Island", -42.5, 70.0, 5, sx, 1, w1),
		site("ILV", "Ilvane Spit", 72.0, -165.0, 3, sx, 2, w2),
		site("CSL", "Cape Solvane", 35.0, -45.0, 5, sx, 1, w2),
		site("SRH", "Saltreach Island", 17.0, 135.0, 6, s, 1, w2),
		site("CRP", "Corith Point", 10.0, 65.0, 6, s, 1, w2),
		site("NIM", "Nimrie Point", -45.0, -150.0, 5, sx, 1, w2),
		site("HSK", "Hasku Point", 32.0, 170.0, 4, sx, 1, w3),
		site("ICM", "Ilha Calmaria", 8.0, -30.0, 7, s, 1, w3),
		site("QND", "Quindle Island", -25.0, 95.0, 5, sx, 1, w3),
		site("PVT", "Punta Ventisca", -60.5, 40.0, 5, sx, 1, w3),
		site("VRL", "Vorrel Island", -64.0, -140.0, 4, sx, 2, w3),
	}}
	add := func(id, name, op string, o Orbit, bands []Band, epoch time.Time, rev0 int, node, phase, revs float64) {
		n.Satellites = append(n.Satellites, Satellite{ID: id, Name: name, Operator: op, Orbit: o, Bands: bands, Mode: "NOMINAL"})
		n.Elements = append(n.Elements, Elements{SatID: id, Orbit: o, Epoch: epoch, Rev0: rev0, NodeDeg: node, PhaseDeg: phase, RevsPerDay: revs})
	}
	jan := time.Date(2029, 1, 1, 0, 0, 0, 0, time.UTC)
	xs := []Band{XBand, SBand}
	// Pipit-1, Pipit-2, Furrow-A, Furrow-B, Brine-3 and Cirrine-1 fitted to passes the
	// stations recorded (Furrow-A and Furrow-B make 867 orbits in 57 days); Pipit-3 from
	// Apogee's own tracking.
	add("pipit-1", "Pipit-1", "Apogee", Orbit{525, 97.5}, []Band{UHF}, jan, 14094, 201.4926025817, 176.9565679403, 0)
	add("pipit-2", "Pipit-2", "Apogee", Orbit{550, 97.6}, s, jan, 19752, 234.8439482122, 339.4065995135, 0)
	add("pipit-3", "Pipit-3", "Apogee", Orbit{510, 97.4}, sx, jan, 9689, 150, 40, 0)
	add("furrow-a", "Furrow-A", "Furrowlight", Orbit{510, 97.4}, xs, jan, 21818, 294.5055051072, 52.6136621172, 867.0/57)
	add("furrow-b", "Furrow-B", "Furrowlight", Orbit{510, 97.4}, xs, jan, 21815, 297.5851231763, 231.0891932221, 867.0/57)
	add("brine-3", "Brine-3", "Brinewatch", Orbit{480, 53.0}, s, jan, 15209, 179.2270109042, 188.3141532885, 0)
	add("cirrine-1", "Cirrine-1", "Cirrine Atmospherics", Orbit{600, 97.8}, s, jan, 17206, 282.5236543824, 268.6744811484, 0)
	add("thistle-2", "Thistle-2", "Kessick Polytechnic", Orbit{505, 97.4}, []Band{UHF},
		time.Date(2029, 3, 8, 10, 20, 0, 0, time.UTC), 1, -120, 0, 0)
	for k := 2; k <= 13; k++ {
		add(fmt.Sprintf("cirrine-%d", k), fmt.Sprintf("Cirrine-%d", k), "Cirrine Atmospherics", Orbit{600, 97.8}, s,
			time.Date(2029, 4, 17, 9, 40, 0, 0, time.UTC), 1, -150, float64(30*(k-2)), 0)
	}
	for k := 4; k <= 39; k++ {
		add(fmt.Sprintf("brine-%d", k), fmt.Sprintf("Brine-%d", k), "Brinewatch", Orbit{480, 53.0}, s,
			time.Date(2029, 5, 16, 13, 30, 0, 0, time.UTC), 1, 120, float64(10*(k-4)), 0)
	}
	for k := 0; k < 24; k++ {
		add(fmt.Sprintf("furrow-%c", 'c'+k), fmt.Sprintf("Furrow-%c", 'C'+k), "Furrowlight", Orbit{510, 97.4}, xs,
			time.Date(2029, 6, 27, 20, 40, 0, 0, time.UTC), 1, 10, float64(15*k), 0)
	}
	return n
}
```

`model.go`, exactly (it needs `fmt`, `math`, `strconv` and `time`):

```go
// The pass model Season 3's reviewers settled: a circular orbit around a round Earth, the
// orbit's plane turning with the mean Sun for a sun-synchronous satellite (one tilted past
// 90°) and by Earth's bulge for any other, and Earth turning beneath it. A pass is the time the
// satellite's elevation over a station stays at or above the station's mask: its AOS and LOS
// are the instants the elevation crosses the mask, rounded to the second, and its maximum
// elevation is the highest of its elevations a second apart from AOS, to a tenth of a degree.

const (
	earthRadiusKm = 6371.0
	muKm3PerS2    = 398600.4418
	j2            = 1.08263e-3                       // how strongly Earth's bulge turns an orbit
	earthRate     = 2 * math.Pi / 86164.0905         // Earth's turn, radians a second
	sunRate       = 2 * math.Pi / (365.2422 * 86400) // the mean Sun's, radians a second
)

func rad(d float64) float64 { return d * math.Pi / 180 }
func deg(r float64) float64 { return r * 180 / math.Pi }

// periodS returns an orbit's period in seconds, as its altitude gives it.
func periodS(o Orbit) float64 {
	r := earthRadiusKm + o.AltKm
	return 2 * math.Pi * math.Sqrt(r*r*r/muKm3PerS2)
}

// skyTrack is one satellite over one station: everything predict needs, worked out once.
// Its times are seconds after the satellite's Epoch.
type skyTrack struct {
	a, cosI, sinI   float64    // the orbit's radius in km, and its tilt
	n, u0           float64    // radians a second along the orbit, and how far along at Epoch
	node0, nodeRate float64    // the node's longitude at Epoch, and how fast it moves over Earth
	st, up          [3]float64 // the station, km from Earth's center, and straight up there
	mask, reach     float64    // the mask, and the widest angle at Earth's center it sees out to
	speed           float64    // no satellite and station draw apart or together faster, radians a second
}

func newSkyTrack(e Elements, s Site) *skyTrack {
	k := &skyTrack{a: earthRadiusKm + e.Orbit.AltKm, mask: s.MaskDeg}
	i := rad(e.Orbit.IncDeg)
	k.cosI, k.sinI = math.Cos(i), math.Sin(i)
	k.n = math.Sqrt(muKm3PerS2 / (k.a * k.a * k.a))
	if e.RevsPerDay > 0 {
		k.n = e.RevsPerDay * 2 * math.Pi / 86400
	}
	drift := -1.5 * k.n * j2 * (earthRadiusKm / k.a) * (earthRadiusKm / k.a) * k.cosI
	if e.Orbit.IncDeg > 90 {
		drift = sunRate // sun-synchronous: the plane keeps its angle to the Sun
	}
	k.node0, k.nodeRate, k.u0 = rad(e.NodeDeg), drift-earthRate, rad(e.PhaseDeg)
	la, lo := rad(s.LatDeg), rad(s.LonDeg)
	k.up = [3]float64{math.Cos(la) * math.Cos(lo), math.Cos(la) * math.Sin(lo), math.Sin(la)}
	for j := range 3 {
		k.st[j] = earthRadiusKm * k.up[j]
	}
	m := rad(s.MaskDeg)
	k.reach = math.Acos(earthRadiusKm*math.Cos(m)/k.a) - m
	k.speed = k.n + earthRate + math.Abs(drift)
	return k
}

// pos is the satellite's place over the turning Earth, km from its center.
func (k *skyTrack) pos(t float64) (x, y, z float64) {
	su, cu := math.Sincos(k.u0 + k.n*t)
	sn, cn := math.Sincos(k.node0 + k.nodeRate*t)
	return k.a * (cn*cu - sn*su*k.cosI), k.a * (sn*cu + cn*su*k.cosI), k.a * su * k.sinI
}

// elev is the satellite's elevation over the station at t, in degrees.
func (k *skyTrack) elev(t float64) float64 {
	x, y, z := k.pos(t)
	dx, dy, dz := x-k.st[0], y-k.st[1], z-k.st[2]
	d := math.Sqrt(dx*dx + dy*dy + dz*dz)
	return deg(math.Asin((dx*k.up[0] + dy*k.up[1] + dz*k.up[2]) / d))
}

// apart is the angle at Earth's center between the satellite and the station.
func (k *skyTrack) apart(t float64) float64 {
	x, y, z := k.pos(t)
	return math.Acos(max(-1, min(1, (x*k.up[0]+y*k.up[1]+z*k.up[2])/k.a)))
}

// above is the elevation over the mask: negative while the satellite is below it.
func (k *skyTrack) above(t float64) float64 { return k.elev(t) - k.mask }

// orbit is the number of the orbit on which the satellite, near t, crosses the station's
// latitude: Rev0 at Epoch, one more at each northbound equator crossing. A satellite whose
// track never reaches the station's latitude counts from the top or bottom of its track.
func (k *skyTrack) orbit(rev0 int, t float64) int {
	u := k.u0 + k.n*t
	var marks []float64 // where along an orbit the satellite is at the station's latitude
	if s := k.up[2] / k.sinI; s >= 1 {
		marks = []float64{math.Pi / 2}
	} else if s <= -1 {
		marks = []float64{3 * math.Pi / 2}
	} else {
		marks = []float64{math.Asin(s), math.Pi - math.Asin(s)}
	}
	best := math.Inf(1)
	var at float64
	for _, m := range marks {
		c := m + 2*math.Pi*math.Round((u-m)/(2*math.Pi)) // the nearest such point to u
		if d := math.Abs(c - u); d < best {
			best, at = d, c
		}
	}
	return rev0 + int(math.Floor(at/(2*math.Pi)))
}

// predict returns e's passes over s whose AOS is in [from, to), by AOS.
// No pass is before Epoch or before the station goes live, and a pass shorter than two
// minutes doesn't count. Near the top of its track a satellite crosses a high station's
// latitude twice on one orbit; the station sees it rise once and set once: one pass.
func predict(e Elements, s Site, from, to time.Time) []PassView {
	from = later(from, later(e.Epoch, s.Live))
	if !from.Before(to) {
		return nil
	}
	k := newSkyTrack(e, s)
	start, end := from.Sub(e.Epoch).Seconds(), to.Sub(e.Epoch).Seconds()
	at := func(t float64) time.Time { // to the nearest second
		return e.Epoch.Add(time.Duration(math.RoundToEven(t)) * time.Second)
	}
	tenths := func(v float64) float64 { // to a tenth, as the settled model prints it
		f, _ := strconv.ParseFloat(strconv.FormatFloat(v, 'f', 1, 64), 64)
		return f
	}
	var out []PassView
	const step = 20.0
	t := start - 1800 // half an hour early: a pass in progress at from is seen rising
	prev := k.above(t)
	for t < end+60 {
		if prev < 0 {
			// Far below the mask the satellite can't come up for a while: skip what it can't.
			if gap := k.apart(t) - k.reach - 0.01; gap > 0 {
				if jump := math.Floor(gap/k.speed/step) - 1; jump >= 1 {
					t += jump * step
				}
			}
		}
		t += step
		cur := k.above(t)
		if prev < 0 && cur >= 0 {
			aos := k.crossing(t-step, t)
			peak, top := k.highest(aos, int(2*k.reach/(k.n-earthRate))+60)
			end := aos + float64(2*peak) + 120
			for k.above(end) >= 0 {
				end += 60
			}
			los := k.crossing(aos+float64(peak), end)
			p := Pass{SatID: e.SatID, Station: s.Code, AOS: at(aos), LOS: at(los), MaxElevDeg: tenths(top)}
			id := fmt.Sprintf("%s/%s/%05d", e.SatID, s.Code, k.orbit(e.Rev0, aos+float64(peak)))
			out = append(out, PassView{ID: id, Pass: p})
			t = los
			cur = k.above(t)
		}
		prev = cur
	}
	kept := out[:0]
	for _, p := range out {
		if !p.AOS.Before(from) && p.AOS.Before(to) && p.LOS.Sub(p.AOS) >= 2*time.Minute {
			kept = append(kept, p)
		}
	}
	return kept
}

// crossing returns the instant between lo and hi where the elevation crosses the mask, below
// it at one end and not at the other: the Illinois method, to a hundred-millionth of a second.
func (k *skyTrack) crossing(lo, hi float64) float64 {
	f := k.above
	flo, fhi := f(lo), f(hi)
	side := 0
	for range 100 {
		m := (lo*fhi - hi*flo) / (fhi - flo)
		if !(m > lo && m < hi) {
			m = (lo + hi) / 2
		}
		fm := f(m)
		if (fm >= 0) == (fhi >= 0) {
			hi, fhi = m, fm
			if side == 1 {
				flo /= 2
			}
			side = 1
		} else {
			lo, flo = m, fm
			if side == -1 {
				fhi /= 2
			}
			side = -1
		}
		if hi-lo < 1e-8 || fm == 0 {
			return m
		}
	}
	return (lo + hi) / 2
}

// highest returns the whole second i in [0, n] after aos when the elevation is largest, and
// that elevation. The elevation rises to one peak, then falls, so a binary search on which
// way it is heading finds the peak.
func (k *skyTrack) highest(aos float64, n int) (int, float64) {
	f := func(i int) float64 { return k.elev(aos + float64(i)) }
	lo, hi := 0, n
	for lo < hi {
		m := (lo + hi) / 2
		if f(m+1) > f(m) {
			lo = m + 1
		} else {
			hi = m
		}
	}
	return lo, f(lo)
}

// later returns the later of two times.
func later(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// earlier returns the earlier of two times.
func earlier(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
```

The model is hidden where a lesson only needs passes and read-only where a lesson calls it by
name. Its type is `Predictor` when code takes it as a value:
`type Predictor func(e Elements, s Site, from, to time.Time) []PassView`.

**Reference numbers.** A copy of the shared files is right only if it gives these, all
computed with Go 1.27.3 (the calibration writer checks them before chapter 4 is written):

| What | Number |
|---|---|
| Passes on Thursday 28 June (every satellite over every station) | 7,485 (Furrow-A: 105) |
| ... at stations that can serve the satellite (`canServe`, below) | 7,289 |
| ... of those, candidates: max elevation at least 10° | 6,102 |
| Servable passes, 1 to 7 July | 51,174 (the first: `brine-27/VSK/00697` 00:00:04, `furrow-s/SRH/00049` 00:00:15, `brine-16/CRP/00697` 00:00:18) |
| Third quarter, 1 July to 30 September: all passes / servable | 690,261 / 672,292 |
| Furrowlight's 26 satellites, third quarter, servable | 249,045 (1 to 14 July: 37,917, the first `furrow-s/SRH/00049` 00:00:15) |
| Contacts planned for Monday 5 March (5 stations, 7 satellites) | 128, from 142 candidates |
| Contacts for Monday 2 April (10 stations, 8 satellites) | 260, from 286 candidates |
| Tirrin Rock alone on Tuesday 17 April | 120 contacts, from 148 candidates |
| Contacts for Wednesday 27 June (the old planner's last day) | 1,653 |
| Contacts for 28 June planned from scratch (no hand-over, no stow) | 2,020, from 6,102 candidates, 4,082 dropped; Cape Wren 76 of its 360 |
| From the hand-over, 28 June to 5 July (lesson 12.8's week) | 14,282 contacts, from 42,853 candidates, 28,571 dropped |

#### Planning: the rule

A **candidate** is a predicted pass at a station that can serve the satellite (they share a
band), whose max elevation is at least `Config.MinElevDeg` (10°), and that no stow overlaps.
Each station is planned on its own, in **AOS order, then pass ID**: a pass takes the
lowest-numbered dish that is free at its AOS, a dish being free once its last contact's LOS
plus the **turnaround** (2 minutes, to turn and settle) has passed; a pass no dish can take
is **dropped**. A pass given a dish is a **contact**. Stations never share a dish, so
stations can be planned at once, in any order, and the plan is the same. `plan.go`, exactly
(it needs `cmp`, `iter`, `slices`, `strings` and `time`):

```go
// byAOS orders passes by AOS, then by pass ID.
func byAOS(a, b PassView) int {
	return cmp.Or(a.AOS.Compare(b.AOS), strings.Compare(a.ID, b.ID))
}

// canServe reports whether the station offers one of the satellite's bands.
func canServe(s Site, sat Satellite) bool {
	return slices.ContainsFunc(sat.Bands, func(b Band) bool { return slices.Contains(s.Bands, b) })
}

// stowedDuring reports whether a stow at the pass's station overlaps the pass.
func stowedDuring(p Pass, stows []Stow) bool {
	for _, st := range stows {
		if st.Station == p.Station && (st.To.IsZero() || p.AOS.Before(st.To)) && st.From.Before(p.LOS) {
			return true
		}
	}
	return false
}

// assign gives one station's candidate passes dishes. The passes come in AOS
// order, then ID; free[d] is when dish d+1 is next free (zero: free now).
// It returns the contacts in AOS order and how many passes found no dish.
func assign(s Site, passes iter.Seq[PassView], free []time.Time, turnaround time.Duration) ([]Contact, int) {
	next := make([]time.Time, s.Dishes)
	copy(next, free)
	var out []Contact
	dropped := 0
	for p := range passes {
		d := slices.IndexFunc(next, func(t time.Time) bool { return !p.AOS.Before(t) })
		if d < 0 {
			dropped++
			continue
		}
		next[d] = p.LOS.Add(turnaround)
		out = append(out, Contact{PassView: p, Dish: d + 1})
	}
	return out, dropped
}
```

`byAOS` is the comparator everywhere: schedules, merges, catalogs. Chapter 6 teaches it;
`plan.go` is given from chapter 6 on. Before that, a lesson that needs it gives `byAOS` read-only or hidden as it stands until 6.1, comparing the AOS alone (`// byAOS orders passes by AOS.` above `return a.AOS.Compare(b.AOS)`), and no check depends on the order of passes that rise in the same second; 6.1 adds the pass ID after Port Ellery's tie of Sunday 8 April, the first two candidates at one station to rise in the same second since 2 April.

**The window.** passplan 2 plans a **window** of a hundred days: from 00:00 UTC of the current
day to 00:00 a hundred days later (28 June to 6 October on take-over day). A schedule holds
every contact whose AOS is in its window.

**The hand-over** (two planners, one dish). At a **cut-over** instant, contacts of the current
schedule whose AOS is before it stay as they are, flown by whoever planned them; planning
starts at the cut-over with each dish free at the latest LOS plus turnaround of its contacts
with AOS before the cut-over. The same rule makes every replan of passplan 2 a **continuation**:
replanning at any instant with nothing new gives exactly the contacts the schedule already
had from that instant on. At 00:00 on 28 June the old planner's schedule has twelve contacts in
progress and two that ended less than two minutes before midnight; planned without the
hand-over, ten of passplan 2's first contacts would take a dish still busy with one of the
old planner's (the first: `furrow-j/SRH/00003` at 00:00:40, on the dish `furrow-k/SRH/00003`
holds until 00:05:00). Run 1 is itself a replan, of the old planner's schedule from the cut-over, so passplan 2's schedules on take-over day also carry the old planner's contacts, which rose before their window, and every later replan's hand-over still waits for their dishes. Planned from the hand-over alone, version 1 would forget them, and the 00:10 tick would give Ilvane Spit's dish 2 to `cirrine-4/ILV/01069` (AOS 00:10:23) while `furrow-f/ILV/00003`'s turnaround lasts until 00:10:26: a fourth version, and one dish touched by both planners.

**Publishing.** Every planning run takes the next **run number** when it starts (counting from
1 when passplan 2 starts) and its schedule carries it as `Version`. A run whose contacts equal
the current schedule's publishes nothing. The publisher keeps the current schedule in an
`atomic.Pointer` and accepts a schedule only from a newer run (compare-and-swap); a schedule is
never changed after it is published. passplan 2 never runs two plans at once on purpose: a new
trigger cancels the run in progress and starts a fresh one, and the publisher still refuses
anything older than what it holds.

```go
// Window is the span a plan covers: From included, To excluded.
type Window struct {
	From time.Time
	To   time.Time
}

// Stow is a span when a station's dishes are parked. A zero To lasts until further notice.
type Stow struct {
	Station string    // "CWR"
	From    time.Time // 2029-06-28T09:47:00Z
	To      time.Time // 2029-06-28T13:23:00Z
	Reason  string    // "wind 74 km/h"
}

// Contact is a pass given a dish: what a station will fly.
type Contact struct {
	PassView     // the pass and its ID
	Dish     int `json:"dish"` // 1 or 2
}

// Schedule is one published plan. Nobody changes it after it is published.
type Schedule struct {
	Version  uint64    // the run that made it
	Window             // what it plans
	Made     time.Time // when its run started
	Contacts []Contact // every contact with AOS in the window, by AOS then ID
	Dropped  int       // candidates no dish could take
}
```

#### passplan 2's settings, bounds and stages

```go
// Config is passplan 2's settings.
type Config struct {
	Slots      int           `check:"min=1,max=64" default:"12"`
	Days       int           `check:"min=1,max=120" default:"100"`
	MinElevDeg float64       `check:"min=0,max=45" default:"10"`
	Turnaround time.Duration `check:"min=0s,max=10m" default:"2m"`
	Replan     time.Duration `check:"min=1m,max=1h" default:"10m"`
	RunLimit   time.Duration `check:"min=1s,max=5m" default:"30s"`
	Queue      int           `check:"min=1,max=4096" default:"256"`
	Quiet      time.Duration `check:"min=0s,max=1m" default:"5s"`
}
```

The bounds, everywhere:

| Bound | Value | Where |
|---|---|---|
| Planner slots | 12: jobs (a satellite to predict, a station to plan) running at once, one per core of passplan's machine | chapters 2, 12 |
| Window | 100 days | chapters 4, 12 |
| Turnaround | 2 minutes | the rule |
| Minimum max elevation for a contact | 10° | the rule |
| Replan | every 10 minutes, and at once on news | chapters 8, 12 |
| Run limit | 30 seconds per planning run; past it the run is canceled | chapter 8 |
| Quiet | 5 seconds: news arriving closer together than this makes one replan | chapter 8 |
| Dish loader queue | 256 contacts per station; when full, passplan waits | chapter 2 |
| Dish loader uploads | 2 at once per station (the hosts' lines are thin) | chapter 2 |
| The Board's queue | 4 schedules; when full, the oldest waiting one is dropped, never the newest | chapter 2 |
| tlmd's stage channels | room for 16 values each | chapter 1 |
| tlmd's parse workers | 4, at most 4 lines in flight | chapter 1 |

The **stages** of a planning run, in order:

1. **Load** the network once for every caller (`sync.OnceValues`, chapter 3).
2. **Predict**: for each station, every satellite it can serve, lazily, a day at a time
   (`Passes`, chapter 4), at most twelve stations at once (chapter 2's group, chapter 12).
3. **Merge**: each station's eighty sorted streams into one stream in AOS order (`MergeBy`,
   a priority queue, chapter 6), filtered to candidates (chapter 5's adapters).
4. **Assign**: the rule above, per station, from the dishes' free times (chapter 12).
5. **Collect**: the twenty stations' contacts merged into one schedule by AOS (`MergeBy`).
6. **Publish**: compare-and-swap into the publisher, newer runs only (chapter 3).
7. **Deliver**: each station's contacts to its dish loader through a queue of 256 (wait when
   full), the schedule to the Board through a queue of 4 (drop the oldest), chapter 2.

A run is cancellable at every stage: every goroutine selects on its context, the predictor
is asked again only between days, and a canceled run publishes nothing and leaves no
goroutine behind.

#### tlmd's ingest (chapter 1)

The ingest is a pipeline: one `read` stage per telemetry feed, `merge`, `parse`, then the
`newest` sink. A **telemetry feed** is the stream of telemetry lines one station hears; lines
are the bible's telemetry lines. A satellite heard by two stations appears in both feeds, so its
frames can arrive twice and out of order; the newest frame is the one with the highest
sequence number, whatever order it arrived in.

```go
// Heard is one telemetry line as a station heard it.
type Heard struct {
	Station string // "TIR"
	Line    string // "cirrine-1 9012 2029-03-07T20:45:05Z batt_v=8.02 temp_c=11.5 mode=NOMINAL"
}

// Health is each satellite's newest frame (its highest Seq), by satellite ID.
type Health map[string]TelemetryFrame
```

- Each stage owns the channel it returns: it alone sends on it and it closes it when it is done.
  Channels have room for 16 values. Every send and every receive also selects on `ctx.Done()`.
- `parseLine` is given read-only in chapter 1: Season 1's reading of a line, by `strings.Fields`,
  `strconv.ParseUint`, `time.Parse(time.RFC3339, ...)` and `strconv.ParseFloat`. The `parse` stage
  drops a line `parseLine` refuses and goes on. Chapter 7 makes `parseLine` strict.
- `readFeed` reports a broken feed by cancelling the run with `context.WithCancelCause`'s
  cancel, the cause wrapped as `feed YST: <err>`; `ingest` returns `context.Cause(ctx)` after
  every stage has stopped, or nil and the Health when every feed ended.
- `parseInOrder` parses with four workers and sends frames in the order their lines arrived:
  each line gets a channel of its own for its result, queued in order (room for `workers`), and
  the sender reads those channels in turn. At most `workers` lines are parsed at once.

#### Telemetry lines, strictly (chapter 7)

A line tlmd accepts is exactly what `formatLine` writes: six fields, one space apart, nothing
before or after.

```go
// formatLine writes f as a telemetry line, the one form tlmd accepts.
func formatLine(f TelemetryFrame) string {
	return fmt.Sprintf("%s %d %s batt_v=%.2f temp_c=%.1f mode=%s",
		f.SatID, f.Seq, f.At.UTC().Format("2006-01-02T15:04:05Z"), f.BattV, f.TempC, f.Mode)
}
```

- The satellite ID is 1 to 24 characters of lowercase letters, digits and `-`, starting with a
  letter. The sequence number is digits only, no sign, no leading zero, at most 4294967295. The
  time is `2006-01-02T15:04:05Z` exactly. `batt_v` has exactly two decimals and is from 0.00 to
  20.00; `temp_c` exactly one and is from −100.0 to 150.0, and never `NaN` or an infinity (a
  NaN compares false with every number, which is how `batt_v=nan` slipped past the range check).
  The mode is `NOMINAL`, `SAFE`, `DETUMBLE` or `LEOP`.
- The strict `parseLine` checks fields in order and returns the first problem, wrapped as
  `line: <problem>`: `want 6 fields, got 5`; `satellite "Pipit-1": want lowercase letters,
  digits and -`; `sequence "-1": want a whole number from 0 to 4294967295`; `time
  "2029-04-14T03:12:00+00:00": want 2006-01-02T15:04:05Z`; `field 4: want batt_v=`; `batt_v
  "nan": want volts like 7.91`; `batt_v 31.20: want 0.00 to 20.00`; `temp_c "1e1": want
  degrees like -4.5`; `mode "nominal": want NOMINAL, SAFE, DETUMBLE or LEOP`. Its last check
  is the round trip: a line that parses but isn't `formatLine` of its own frame (such as
  `batt_v=+7.91` or `batt_v=07.91`) is `line: not canonical`.
- Thistle-2's lines of Saturday 14 April, as Joss forwarded them, include
  `thistle-2 2214 2029-04-14T03:36:12Z batt_v=nan temp_c=1.0 mode=SAFE` (heard at Port Ellery,
  pass `thistle-2/PEL/00560`) and the line that crashed the old parser on the next pass, `thistle-2 2290 2029-04-14T05:10:41Z batt_v=7.40` (four fields). Thistle-2 numbers a frame about every 75 seconds. Mirela keeps four of Saturday's lines as corpus files: those two (`saturday-0336`, `saturday-0510`), `thistle-2 2216 2029-04-14T03:38:42Z batt_v=7.42 temp_c=-Inf mode=SAFE` (`saturday-0338`) and the 05:13:11 line that lost its satellite ID and starts with a space (`saturday-0513`). Joss's fourteen odd lines (7.8) are frame 2214 with the battery at 7.91 V written fourteen wrong ways: `batt_v=+7.91`, `batt_v=07.91`, a time ending `+00:00`, the sequence `02214`, `batt_v=7.9`, `temp_c=1`, `temp_c=1e0`, `batt_v=NaN`, `mode=safe`, `Thistle-2`, the sequence `-1`, a time ending `.000Z`, two spaces after the satellite ID, and a seventh field `crc=6e`.
- The **Pipit beacon** (chapter 7.6) is the bible's 12 bytes, decoded into
  `Beacon{Sat uint8; Seq uint16; Mode uint8; BattMV uint16; TempDeciC int16; Flags uint8}`:
  `decodeBeacon(b []byte) (Beacon, error)` refuses a length other than 12 (`beacon: want 12
  bytes, got 7`), wrong sync bytes (`beacon: no sync`) and a bad XOR (`beacon: bad check
  byte`); `encodeBeacon(bc Beacon) []byte` writes it back. Every Pipit sends this beacon, and Pipit-1 is satellite number 1 and Pipit-2 number 2, as in Season 1.

How fuzz lessons run (the app's `go test` runs a fuzz test's seed corpus only, never the
fuzzer): the read-only `main_test.go` holds `FuzzParseLine` (and `FuzzBeacon`). It adds every
seed in the learner's `seeds.go` (`var lineSeeds []string`) with `f.Add`, and `go test` also
runs every file in `testdata/fuzz/FuzzParseLine/` (format: the line `go test fuzz v1`, then
`string("...")`). Its target calls the learner's property, `checkLine(t, parse, line)` in
`checks.go`, through `safely`, which turns a panic into the failure `parseLine("...") panicked:
<value>`, so a crash reads as a failed check rather than a stack dump. Mirela's reference parser, which accepts exactly the lines `formatLine` writes with a valid satellite ID, readings in range and a known mode, is hidden in `reference_test.go`; mutant parsers built from it live in the test file; tests that prove the seeds or the property catch them re-run the test binary
(`os.Args[0] -test.run=^FuzzParseLine$`, `APOGEE_MUTANT=<name>`) and expect it to fail, as
Season 1's chapter 15 did. Readmes say once how to fuzz for real in a terminal:
`go test -fuzz=FuzzParseLine -fuzztime=30s`, and that a failure it finds is saved under
`testdata/fuzz/FuzzParseLine/` and runs as a seed from then on.

#### Iterators and generic types (chapters 4 to 6)

```go
type Predictor func(e Elements, s Site, from, to time.Time) []PassView

func Days(from time.Time, n int) iter.Seq[time.Time]                               // 4.1: n midnights
func Passes(predict Predictor, e Elements, s Site, w Window) iter.Seq[PassView]     // 4.2: a day at a time
func (sch *Schedule) At(station string) iter.Seq2[int, Contact]                     // 4.4: (dish, contact) by AOS
func Catalog(n *Network, w Window, predict Predictor) iter.Seq[PassView]            // 4.5: servable, in byAOS order
func ReadPasses(r io.Reader) iter.Seq2[PassView, error]                             // 4.6: the pass archive
func parsePass(line string) (PassView, error)                                       // 4.6: given, one archive line
func Archive(open func() (io.ReadCloser, error)) iter.Seq2[PassView, error]         // 4.7: closes on every way out
func Bookable(n *Network, w Window, predict Predictor, operator string) iter.Seq[PassView] // 4.8

func Filter[V any](seq iter.Seq[V], keep func(V) bool) iter.Seq[V]   // 5.2
func Map[V, W any](seq iter.Seq[V], f func(V) W) iter.Seq[W]         // 5.3
func Take[V any](seq iter.Seq[V], n int) iter.Seq[V]                 // 5.3
func Zip[V any](a, b iter.Seq[V]) iter.Seq2[V, V]                    // 5.6: pairs, until either ends
func Merge[V any](a, b iter.Seq[V], cmp func(V, V) int) iter.Seq[V] // 5.7: two sorted into one
func Compare(old, new iter.Seq[Contact]) iter.Seq[Change]            // 5.9: the shadow run
func changes(olds, news []Contact) []Change                          // 5.9: given, one second's contacts matched by pass ID

// Change is one difference between two schedules.
type Change struct {
	Kind    string  // "old only", "new only" or "dish"
	Contact Contact // the old planner's for "old only", the new planner's otherwise
}

type passHeap []PassView // 6.2: heap.Interface by byAOS

// Queue is a priority queue: Pop returns the smallest value by cmp.
type Queue[T any] struct{ h *heapOf[T] } // heapOf implements heap.Interface for a slice and cmp

func NewQueue[T any](cmp func(a, b T) int) *Queue[T] // 6.4
func (q *Queue[T]) Push(v T)
func (q *Queue[T]) Pop() (T, bool)  // the zero value and false when empty
func (q *Queue[T]) Len() int
func (q *Queue[T]) Drain() iter.Seq[T] // 6.5: pops until empty or the loop stops
func (q *Queue[T]) Peek() (T, bool)    // 6.6

func MergeBy[T any](cmp func(a, b T) int, seqs ...iter.Seq[T]) iter.Seq[T] // 6.7: k sorted into one
func stationPasses(n *Network, s Site, w Window, predict Predictor, minElevDeg float64) iter.Seq[PassView] // 6.8
```

- An iterator stops the moment `yield` returns false and never calls it again (Go panics with `runtime error: range function continued iteration after function for loop body returned false`).
  An iterator that opened something closes it on every way out, with `defer` inside the
  iterator function.
- Generated sequences are made a day at a time: `Passes` asks the predictor for one day, yields
  it, and only then asks for the next; `Catalog` predicts every pair for one day, sorts that
  day by `byAOS`, yields it, and moves on. Breaking out of either after the first pass costs
  one day of prediction, never the quarter.
- `Merge` and `MergeBy` keep ties stable: on equal keys the earlier sequence's value comes
  first. They use `iter.Pull` and call every `stop` they got, whether the loop ends or breaks.
  `MergeBy` holds one value per input sequence at a time, in a `Queue`.
- `Compare` walks both schedules a second at a time: it gathers each schedule's contacts with the earliest AOS either has left and yields what `changes` finds among them (the old planner's `old only` and `dish` changes in its order, then `new only` in passplan 2's), so contacts at different stations that rise in the same second (5 April, 05:00:11: `brine-3/SQT/16651` and `furrow-a/VSK/23251`) match in whatever order each planner lists them. On Thursday 5 April Tirrin Rock's host parked both dishes from 07:30 to 08:24 to clear their ice, a stow only the old planner heard of: its 273 contacts and passplan 2's 276 differ in four, all at Tirrin Rock (`new only` `furrow-b/TIR/23250` 07:38:06, `pipit-2/TIR/21175` 07:52:22 and `pipit-3/TIR/11124` 08:23:30, then `dish` `furrow-a/TIR/23253` 08:24:16, on dish 2 for passplan 2 and dish 1 for the old planner). passplan 2 hears of stows from Friday 6 April, and planned again with the stow, Thursday has 273 contacts in both.
- The **pass archive** is one pass per line, `<pass ID> <AOS> <LOS> <max elev>`, times in RFC
  3339 UTC, elevation with one decimal:
  `furrow-a/SQT/24534 2029-06-28T14:02:10Z 2029-06-28T14:12:04Z 47.9`. A bad line yields its
  error, `archive line 3: <problem>`, and ends the sequence. The archive is Apogee's file of every past pass; lessons quote it from 1 January 2029 on, where the model's passes begin. `parsePass` reads a line; its problems read `3 fields, want 4`, `bad pass ID "…"`, `bad AOS "…"`, `bad LOS "…"` and `bad max elevation "…"`. A failed read yields the reader's error, unwrapped, after the passes before it, and `Archive` yields `open`'s error unwrapped.

#### Time in a bubble (chapter 8)

Lessons test timers with `testing/synctest`. A bubble's clock starts at 2000-01-01 00:00 UTC,
so code takes `now func() time.Time` for instants and uses the `time` package (timers,
tickers, `Sleep`, context deadlines) for waiting; a test shifts the bubble's clock to 2029:

```go
synctest.Test(t, func(t *testing.T) {
	shift := time.Date(2029, 6, 28, 0, 0, 0, 0, time.UTC).Sub(time.Now())
	now := func() time.Time { return time.Now().Add(shift).UTC() }
	// ... a day of passes, in a millisecond of real time
})
```

```go
func waitUntil(ctx context.Context, now func() time.Time, at time.Time) error             // 8.1
func every(ctx context.Context, d time.Duration, f func(at time.Time))                     // 8.4
func planWithin(ctx context.Context, limit time.Duration, plan func(context.Context) (*Schedule, error)) (*Schedule, error) // 8.5
func newDebouncer(quiet time.Duration, fire func()) *Debouncer                             // 8.6
func (d *Debouncer) Notify()
func (d *Debouncer) Stop()
func flyDay(ctx context.Context, now func() time.Time, contacts []Contact, fly func(Contact)) error // 8.8
```

- `waitUntil` returns nil at `at` exactly (at once when `at` has passed) and `ctx.Err()` if the
  context ends first; it stops its timer either way.
- `every` calls `f` at each tick until the context ends, then stops its ticker and returns: a day
  has 144 ticks of 10 minutes.
- `planWithin` gives the run a context that ends after `limit` and cancels it when the run returns; a run past its limit returns `plan: took longer than 30s: context deadline exceeded` (the message names `limit`), which wraps `context.DeadlineExceeded`; any other error, a caller's cancel included, comes back unchanged.
- A `Debouncer` calls `fire` once, `quiet` after the last of a burst of `Notify` calls; `Stop` cancels a pending fire, and a later `Notify` arms it again. It holds one `time.AfterFunc` timer, made by the first `Notify` and `Reset` by later ones under a mutex; passplan 2's quiet is 5 seconds.
- `flyDay` waits until each contact's AOS in order and calls `fly`; it returns `ctx.Err()` when
  the context ends, having flown exactly the contacts whose AOS came first.
- When a timer and `ctx.Done()` are ready together, a `select` takes either, so no test ends a context at the instant a timer is due: 8.4 cancels at 24:00:01 and 09:47, between ticks, and 8.8 at noon, between contacts at 11:59:51 and 12:00:01.
- Learners write synctest checks in `checks.go` (8.2 and 8.7), proved against mutants as above, re-run with `APOGEE_MUTANT` as in chapter 7. The test runs the learner's check twenty times, stopping at the first failure, so a check that passes only by luck (one that looks without `synctest.Wait`) fails; and it fails a check that calls the code on the real clock, which it tells by a real-clock reading taken just before (a call within a minute of it is outside any bubble; a bubble's clock starts in 2000 and may be moved to any year). A count that a timer's function changes is an atomic, or under a mutex: the race detector sees no order between a test's read and a timer's later write, even in fake time. Inside a bubble a real network read keeps the clock still (in-memory fakes such as `net.Pipe` don't), `T.Run` isn't allowed, and every goroutine the bubble starts must end before the test function returns.

#### Benchmarks and allocations (chapters 9 and 10)

- A learner's benchmark is a function in `bench.go`, such as `benchPlanDay(b *testing.B)`, with
  `for b.Loop()`. Read-only tests run it with `testing.Benchmark` after
  `flag.Set("test.benchtime", "20x")`, so it runs exactly 20 iterations in milliseconds.
  Tests never check nanoseconds: they check iterations, allocations, bytes and the learner's
  own metrics (`b.ReportMetric`). Readmes give the terminal command for the timings: `go test -bench=. -benchmem -count=10`. Tests that set `-test.benchtime` put it back when they end, so a terminal's benchmarks keep their time. A sub-benchmark doesn't inherit `b.ReportAllocs` from its parent; lesson 9.8's test reads the six lines of a re-run of its test binary (`-test.run=^$ -test.bench=^BenchmarkWindows$ -test.benchtime=2x -test.cpu=1`).
- Allocation budgets use `testing.AllocsPerRun` and `BenchmarkResult.AllocsPerOp`, computed on
  Go 1.27.3 by the writer, never guessed. Counts drift under `-race`, so every budget test skips
  there: the test files include `race_on_test.go` (`//go:build race`, `const raceOn = true`)
  and `race_off_test.go` (`//go:build !race`, `const raceOn = false`), and a budget test
  begins `if raceOn { t.Skip("allocation counts differ under -race") }`. Without `-race`, code that formats with `fmt` or pulls with `iter.Pull` still allocates a little more when the collector runs (`planNew` up to about 0.4%; exact with `GOGC=off`), so a budget on such code allows a margin (9.8: 1.5%, its benchmark re-run on one CPU); `planDay` and `assign` count exactly. A range-over-func loop counts the same in a test only when the compiler inlines its function there: out of line, the loop's closure, its state and the variables it writes go on the heap (3 allocations for 10.1's `collect` written with `slices.Grow`), so 10.1 subtracts what `collect` makes for no passes.
- The **old planner** (`planOld` in `planners.go`, given read-only in chapters 9 and 10; a copy in chapter 10 is 9's file byte for byte) is the 2023 passplan's planning brought to the rule and run on passplan 2's pass model: it predicts every pair into one slice, sorts that slice by station, then AOS, then ID, and assigns station by station, growing its slices as it goes and keeping each contact's dish in a map keyed by pass ID. The **new planner** (`planNew`) is per station: chapter 6's `stationPasses` (candidates merged by `MergeBy`), assigned, then the stations' contacts merged by `MergeBy` into a slice sized ahead. Both return the same contacts for the same input, and the new one sorts nothing in place. On one core it is no faster, and heavier: predicting passes is nearly all of both (91 to 96% of `planOld`'s time, 96 to 98% of its allocations), and `MergeBy` starts an `iter.Pull` for every satellite at every station (about six allocations each) and boxes every value it queues (two allocations, through `container/heap`'s `any`). For 6, 12 and 24 hours from midnight on Friday 18 May (303, 619 and 1,237 contacts), `planNew` takes up to about 15% longer and makes 2.1 to 3.3 times the allocations: about 23,600, 30,500 and 43,100 against 7,150, 12,040 and 20,970. No planner that predicts its own passes can make a fraction of `planOld`'s allocations: predicting the 1,564 pairs that share a band is 93 to 94% of them. That is chapter 9's truth on the glass ("Not faster. Heavier. That we can fix."): no later lesson or story line says passplan 2 plans faster than the old planner on one core, or allocates less. Where a story line sets passplan 2 against the old passplan's 4 minutes 39 seconds a day, the difference is the old passplan's own, older pass prediction and passplan 2's twelve planner slots, never passplan 2's planning. Chapter 10 cuts garbage around the merge: `planStationFat` is its own function (`planNew`, readable in 9.6 and 9.8, has no string keys), and chapter 6's `queue.go` and `merge.go` stay as they are, since chapter 12 gives them read-only.
- The **dish loader line** (chapter 10), one contact per line, a day per file: AOS and LOS as `15:04:05`, the dish, the pass ID: `05:29:26 05:37:52 1 furrow-a/VSK/24529`.
- Chapter 10's pieces: a host's report (10.2) is `CWR contacts: 76`, then one line a contact in AOS order, `cirrine-1/CWR/19862 dish 1`; the export (10.5) writes each station's dish loader lines in one write from an 8 KB buffer (Ystmark's 199 lines on take-over day, 7,811 bytes, are the longest day); a planner slot's scratch holds 512 passes (Ystmark's 439 candidates are the most), and `putScratch` keeps none over 65,536 (6.8 MB). The per-station step `planStationFat` (10.8) copies the candidates into a slice grown from nil, skips repeats with a map keyed by `fmt.Sprintf("%s %s %d", SatID, Station, AOS.Unix())`, sorts the copy and assigns it: on take-over day 1,471 allocations and 213 KB at Cape Wren, 1,789 and 251 KB at Ystmark; `planStationLean` 11 and 60 KB, 13 and 93 KB. Planning take-over day at twenty stations from a network loaded once (each station's candidates predicted, the step, the host report and the export of its dish loader lines) makes 81,657 allocations and 12.0 MB with Monday 21 May's code and 42,705 and 6.0 MB with Friday 25 May's (the debrief's "6 MB of garbage instead of 12"); of what is left, predicting passes is 99% of the allocations and 77% of the bytes. No line compares these with `planOld`.
- A `sync.Pool` may lose its values to any collection, keeps a part for each P (a CPU the program runs on), and after a collection allocates that part again on its next use (2 allocations, `New` not called). No check needs a value back from a pool through a collection: 10.6 counts its rounds with the collector off (`debug.SetGCPercent(-1)`) and checks that a 65,536-pass scratch is kept on one CPU (`runtime.GOMAXPROCS(1)`) with the collector off, skipped under `-race`, where `Put` drops one value in four; 10.8 measures on one CPU. With a forced collection before every pool call, every check of 10.6 and 10.8 still passes but 10.6's count, which then reads the pool's own 4 allocations a round.

#### Struct tags (chapter 11)

One validator reads the `check` tag of every exported field; `applyDefaults` reads `default`.
A `check` tag is rules separated by commas: `required` (not the zero value; a slice not empty),
`min=<n>` and `max=<n>` (ints, floats, and durations written as duration text), `len=<n>` (a
string's length in characters, or each string of a `[]string`), `oneof=<a> <b> ...` (a string,
or each string of a `[]string`). A field whose type has no use for a rule is a mistake in the tag, and says so: `Config.Slots: rule len=3 does not fit int`; so does a `min` or `max` bound that doesn't parse for the field's type (`Config.Slots: rule max=6A does not fit int`). An unknown rule is `HostConfig.Name: unknown rule "sometimes"`. A `time.Duration` is told by its type (its kind is `int64`); `oneof` takes words separated by spaces and `len` counts characters, both on a `string` or each string of a `[]string`. A field with no `check` key, or an empty one, has no rules.

```go
// HostConfig is the form a station's host fills in when its dish is installed.
type HostConfig struct {
	Code    string   `check:"required,len=3"`
	Name    string   `check:"required"`
	LatDeg  float64  `check:"min=-90,max=90"`
	LonDeg  float64  `check:"min=-180,max=180"`
	MaskDeg float64  `check:"min=0,max=15"`
	Dishes  int      `check:"min=1,max=2"`
	Bands   []string `check:"required,oneof=UHF S X"`
}

// PipelineConfig is tlmd's ingest settings.
type PipelineConfig struct {
	Workers int      `check:"min=1,max=32" default:"4"`
	Buffer  int      `check:"min=0,max=1024" default:"16"`
	Feeds   []string `check:"required,len=3"` // station codes
}
```

- `Validate(v any) error` takes a struct or a pointer to one and checks every exported field in order and walks into each struct field, embedded or not, right after that field's own rules (pointer fields and slices are not followed), returning every problem joined with `errors.Join`, or nil; a nested struct's problems name its own type (`HostConfig.Dishes` inside a struct that embeds `HostConfig`). Each
  problem is `<Type>.<Field>: <problem>`: `HostConfig.Dishes: 0 is below min 1`,
  `Config.Slots: 70 is above max 64`, `Config.Turnaround: 15m0s is above max 10m0s`,
  `HostConfig.Bands: "K" is not one of UHF S X`, `HostConfig.Code: "VRLX" is not 3 characters`,
  `HostConfig.Name: required`. Ints print with `%d`, floats with `%g`, durations with
  `String()`, strings with `%q`. A value that isn't a struct is `validate: want a struct, got %T`: `got int`, `got []main.HostConfig`.
- `applyDefaults(ptr any) error` sets each zero exported field that has a `default` tag, through
  a pointer (`defaults: want a pointer to a struct, got main.Config`); it leaves set fields, unexported fields and fields without a `default` key alone. Zero means "use the default", so a `Buffer: 0` becomes 16; a default the field can't read returns `setText`'s error, wrapped as `defaults: <Field>: …`.
- The forms of the third wave: Vorrel Island's says `Dishes: 0` (the host meant two); Ilha Calmaria's first form (Thursday 31 May) left the station's name blank and listed bands `S` and `K`, and its corrected form (Friday 1 June) still listed `S` and `K`; the other three are right, and all five were right by go-live. The test deploy of Friday 1 June
  had a `Config` with `Slots: 0`.

#### Take-over day

Wednesday 27 June and Thursday 28 June 2029, UTC:

| When | What |
|---|---|
| Wed 18:05 | Furrowlight's rideshare lifts off with twenty-four Furrows |
| 19:10–20:40 | They leave the rocket, one about every four minutes; elements valid from 20:40:00 |
| 19:56 | The launcher's elements arrive; the old passplan starts replanning the day for eighty satellites |
| 20:00:39 | Its plan is ready, after 4 minutes 39 seconds |
| 20:40:16 | First contact of the twenty-four, sixteen seconds after their elements begin: `furrow-g/YST/00001`; by 22:28:50 (`furrow-c/TIR/00002`) all have been heard |
| 23:30 | Ingrid's go/no-go poll for the take-over: GO |
| Thu 00:00:00 | Cut-over. passplan 2's run 1 replans the old planner's schedule from the cut-over, over 28 June to 6 October, and publishes version 1: the old planner's contacts that rose before 00:00 stay in it, flown by the old planner, and the rest is planned from the hand-over: 2,016 contacts on the day, 201,774 in the window (613,496 candidates, 411,722 dropped). The old planner keeps its twelve contacts in progress |
| 00:01:04 | passplan 2's first contact: `cirrine-4/TIR/01069` |
| 00:08:26 | The old planner's last contact ends (`furrow-f/ILV/00003`); at 00:10 it is switched off |
| 05:29:26 | `furrow-a/VSK/24529`, the day's 417th contact: Furrow-A over Varsk Island at 74.9°, Hollis's pass |
| 07:48:33 | `pipit-1/LAN/16797`, the 612th: Pipit-1's first contact of the day. Its two passes over Port Ellery went to Thistle-2, which rose a few minutes earlier each time |
| 09:47:00 | Cape Wren's wind reaches 74 km/h and its dish stows itself; the stow is recorded with no end. Run 60 publishes version 60: Cape Wren dark until further notice, 1,982 contacts on the day (193,239 in the window) |
| 13:23:00 | The wind has stayed under 72 km/h since 13:15, and Cape Wren unstows; the stow now ends at 13:23. Run 83 publishes version 83: 2,011 contacts on the day (201,769 in the window). Six Cape Wren passes were lost to the wind (`pipit-2/CWR/22443` 10:44:34, `cirrine-1/CWR/19868` 10:54:48, `furrow-b/CWR/24530` 11:51:46, `cirrine-1/CWR/19869` 12:28:44, `furrow-g/CWR/00011` 13:11:36, and `furrow-d/CWR/00011` 13:22:45, still above the horizon at 13:23); the afternoon's 21 Cape Wren contacts became 22 others, from `furrow-b/CWR/24531` at 13:23:30. Cape Wren ends the day with 71 contacts instead of 76 |
| 23:59:49 | The day's last AOS: `cirrine-2/BRM/01083` |
| Fri 00:00 | Ingrid writes in the log: "passplan 2: GO for every day." |

The day in numbers: 146 runs (run 1 at the cut-over, 143 ten-minute ticks from 00:10 to
23:50, and the two at 09:47 and 13:23); three published versions (1, 60, 83), the other runs
finding nothing new; 2,011 contacts flown at twenty stations, 819 of them for the twenty-four
new Furrows; no pass flown twice, none touched by both planners. Lessons that run the whole
day test it with their own `Config` (a replan every hour, a one-day window), so it ends in under a second. 12.9's test also puts a `dayModel` in `passModel`'s place: each pair's passes of the day predicted once and every replan's question answered from them, exactly what `predict` gives (checked for every pair from every ten-minute instant and from 09:47 and 13:23), so the day costs one day of prediction: about 0.35 s and 0.65 s of CPU, against 0.7 s and 2.3 s.

**Bookings.** Where a lesson needs a booking (Season 3's `Booking`), its ID is from `bk-7001`
to `bk-7999`: Season 5 reached `bk-5274`, and the numbers below `bk-7000` are left to Season 6.
Furrowlight's booking of `furrow-a/VSK/24529` on take-over day is `bk-7412`; Kessick
Polytechnic's of Thistle-2's first pass over Port Ellery (`thistle-2/PEL/00005`) is `bk-7001`.

#### The code, by chapter

Each lesson is its own `package main` (module `lesson`). A lesson gives the pieces it needs
from earlier chapters as read-only files under these names; `network.go` and `model.go` are
hidden unless the readme sends the learner to read them. `types.go` holds the types (Season 1's and this season's: `Site`, `Elements`, `Network`, `Window`, `Contact` and the rest), cut to what the lesson needs; `network.go` is the `newNetwork` block alone.

| Chapter | The learner writes (files) | Given |
|---|---|---|
| 1 | `read`, `parse` (cancellable), `newest`, `merge`, `ingest`, `parseInOrder` (`ingest.go`) | `Heard`, `Health`, `parseLine`, `readFeed` from 1.7 (`telemetry.go`) |
| 2 | `Slots`, `newSlots`, `Acquire`, `Release`, `predictAll`, `runPool`, `feed`, `offer`, `Group`, `newGroup`, `Go`, `fail`, `Wait`, `Uploads`, `newUploads`, `(*Uploads).Do` (`slots.go`, `predict.go`, `pool.go`, `queue.go`, `group.go`, `uploads.go`) | `Elements`, `PassView`, `Contact`, `Schedule`; 2.1's `predictAll` (`predict.go`); `network.go` and `model.go` in 2.8; fake predictors in tests |
| 3 | `Tables`, `(*Tables).Network`, `networkOnce`, `Run`, `newRun`, `Stop`, `Done`, `Stats`, `(*Stats).Add`, `Snapshot`, `Publisher`, `Current`, `Update`, `Offer` (`tables.go`, `run.go`, `stats.go`, `publish.go`) | `loadNetwork() (*Network, error)` (the slow loader: in tests it counts its calls), `network.go` |
| 4 | `Days`, `Passes`, `(*Schedule).At`, `Catalog`, `ReadPasses`, `Archive`, `Bookable` (`seq.go`, `archive.go`) | `network.go` and `model.go` (hidden); `Predictor` above `Passes` in `seq.go`; `byAOS` as it stands until 6.1 (the AOS alone) and `canServe` (`plan.go`) in 4.5 and 4.8; `parsePass` (`passline.go`) in 4.6 and 4.7; in 4.4 the schedule of 2 April, planned by a hidden test file from hidden `network.go`, `model.go` and `plan.go` |
| 5 | `Filter`, `Map`, `Take`, `Zip`, `Merge`, `Compare` (`adapt.go`, `shadow.go`); `main` in 5.1 and 5.5 | `types.go`; the earlier adapters back in the editable `adapt.go`, read-only in 5.5; `Change` above `Compare` in `shadow.go`, and `changes` (`match.go`, read-only) in 5.9; the schedule of 2 April, planned from scratch by a hidden test file (a hidden `wave.go` in 5.5) from hidden `network.go`, `model.go` and `plan.go`, with `byAOS` as it stands until 6.1; in 5.9 both planners' schedules of 5 April, planned the same way by a hidden test file |
| 6 | `byAOS`, `passHeap`, `Queue`, `heapOf`, `NewQueue` and methods, `MergeBy`, `stationPasses` (`plan.go`, `queue.go`, `merge.go`) | `network.go`, `model.go`, `plan.go`, `Predictor` and `Passes` (`seq.go`), `Filter` (`adapt.go`) |
| 7 | `parseLine` (strict, `parse.go`), `lineSeeds` (`seeds.go`), `checkLine` (`checks.go`), `decodeBeacon` (`decode.go`) | `formatLine` (`telemetry.go`), `encodeBeacon` and `Beacon` (`beacon.go`), the fuzz targets, Mirela's reference parser (hidden), corpus files |
| 8 | `waitUntil`, `checkWaitUntil` (`checks.go`), `every`, `planWithin`, `Debouncer`, `checkDebouncer` (`checks.go`), `flyDay` (`timers.go`) | `Contact`, `Schedule` (`types.go`); the earlier `timers.go`, read-only, in 8.2 and 8.7; `notifier` (`checks.go`); in 8.8 the day's contacts as a fixture, planned by a hidden test file from hidden `network.go`, `model.go` and `plan.go`; Nkem's mutants in the tests |
| 9 | `benchPlanDay`, `benchFresh`, `benchStream`, `sameContacts`, `benchWindows` (`bench.go`) | `june28`, `stationDay`, `dayCounts`, `planned` and `planDay` (`day.go`, read-only in 9.1, 9.2, 9.4 and 9.5); `planOld` and `planNew` (`planners.go`, read-only in 9.6 and 9.8) with chapter 6's `seq.go`, `adapt.go`, `queue.go` and `merge.go`, hidden; `plan.go` (read-only in 9.1 to 9.5, hidden in 9.6 and 9.8); `network.go` and `model.go`, hidden; `Config` in `types.go` in 9.6 and 9.8 |
| 10 | `collect`, `stationReport`, `appendContact`, `exportStations`, `getScratch`, `putScratch`, `passKey`, `dedupe`, `planStationLean` (`lean.go`) | `planStationFat` (`fat.go`, 10.8); earlier solutions read-only: `appendContact` (`loader.go`, 10.5), the pool with `getScratch`, `putScratch` and `maxScratch` (`scratch.go`, 10.8), `passKey` (`keys.go`, 10.8); `plan.go` (read-only in 10.8, hidden elsewhere); `types.go`; hidden `network.go`, `model.go`, chapter 9's `day.go` (10.1, 10.6–10.8) and chapter 8's `day_test.go` (10.2, 10.4, 10.5) |
| 11 | `fieldLines`, `tagLines`, `describe`, `checkRule`, `measure`, `applyDefaults`, `Validate` (`validate.go`) | `Config`, `HostConfig`, `PipelineConfig` (`types.go`); `main`, with 11.2's `Note` (`main.go`, in 11.1 and 11.2); from 11.6, `checkRule` cutting each rule at `=` with the `min` and `max` messages and `misfit`, in the editable file; `setText` (`defaults.go`) in 11.8; in 11.10, read-only, `checkRule`, `misfit` and `measure` (`rules.go`), `applyDefaults` and `setText` (`defaults.go`) and the hosts' forms as first sent, `hostForms` and `form` (`forms.go`) |
| 12 | `planStation`, `Plan`, `handover`, `replanFrom`, `(*Replanner).Run`, `heard`, `takeOver` (`planner.go`, `replan.go`) | everything above, read-only (`network.go` and `model.go` hidden): `plan.go` (with `stowedDuring`), 4.8's `seq.go`, 5.7's `adapt.go`, chapter 6's `queue.go` and `merge.go`, `group.go`, `publish.go`, `timers.go`; `passmodel.go` (`passModel`, the model `planStation` asks, which tests replace); `Stow`, `Schedule`, `Config` and `Event` in `types.go`; the rest of the `Replanner` (`newReplanner`, `publish`) in `replan.go` from 12.6; the earlier `planner.go`, read-only, in 12.6, 12.7 and 12.9; answer keys in a hidden `ref_test.go` (all but 12.6) |

Signatures every lesson after the one that defines them uses exactly (with the ones listed in
the sections above):

```go
// chapter 1
func read(ctx context.Context, station string, r io.Reader) <-chan Heard
func readFeed(ctx context.Context, cancel context.CancelCauseFunc, station string, r io.Reader) <-chan Heard
func parse(ctx context.Context, in <-chan Heard) <-chan TelemetryFrame
func newest(in <-chan TelemetryFrame) Health
func merge(ctx context.Context, feeds ...<-chan Heard) <-chan Heard
func ingest(ctx context.Context, feeds map[string]io.Reader) (Health, error)
func parseInOrder(ctx context.Context, in <-chan Heard, workers int, parse func(string) (TelemetryFrame, error)) <-chan TelemetryFrame

// chapter 2
type Slots chan struct{}
func newSlots(n int) Slots
func (s Slots) Acquire(ctx context.Context) error // nil once a slot is held; ctx.Err() if ctx ends first
func (s Slots) Release()
func predictAll(ctx context.Context, els []Elements, slots int, predict func(context.Context, Elements) ([]PassView, error)) (map[string][]PassView, error)
func runPool(ctx context.Context, jobs <-chan Elements, workers int, do func(context.Context, Elements) error) error
func feed(ctx context.Context, contacts []Contact, q chan<- Contact) error
func offer(q chan *Schedule, s *Schedule) (dropped int)
func newGroup(ctx context.Context, limit int) (*Group, context.Context)
func (g *Group) Go(job func(ctx context.Context) error) // waits for a slot; runs nothing once the group's context has ended, and fails the group with that cause
func (g *Group) Wait() error                            // the first error, after every job has returned
func newUploads(perStation int) *Uploads
func (u *Uploads) Do(ctx context.Context, station string, upload func() error) error

// chapter 3
func loadNetwork() (*Network, error)
func (t *Tables) Network() (*Network, error)
func networkOnce(load func() (*Network, error)) func() (*Network, error)
func newRun() *Run
func (r *Run) Stop()
func (r *Run) Done() <-chan struct{}
func (s *Stats) Add(planned, dropped int) // one run's counts
func (s *Stats) Snapshot() (planned, dropped, runs int64)
func (p *Publisher) Current() *Schedule                      // nil before the first schedule
func (p *Publisher) Update(contacts []Contact) *Schedule     // a new schedule, Version one more
func (p *Publisher) Offer(s *Schedule) bool                  // only a newer Version replaces the current

// chapter 7
func parseLine(line string) (TelemetryFrame, error)
func checkLine(t *testing.T, parse func(string) (TelemetryFrame, error), line string)
func decodeBeacon(b []byte) (Beacon, error)
func encodeBeacon(bc Beacon) []byte

// chapter 8
func checkWaitUntil(t *testing.T, wait func(context.Context, func() time.Time, time.Time) error)
type notifier interface{ Notify(); Stop() }
func checkDebouncer(t *testing.T, newD func(quiet time.Duration, fire func()) notifier)

// chapter 9
func planOld(n *Network, w Window, cfg Config) []Contact
func planNew(n *Network, w Window, cfg Config) []Contact
func stationDay(code string, day time.Time) (Site, []PassView)                   // the station and its candidates that day, in fleet order
func planDay(cands []PassView, sites []Site, turnaround time.Duration) []Contact // sorts cands in place, then assigns
func benchPlanDay(b *testing.B)
func benchFresh(b *testing.B)
func benchStream(b *testing.B)
func sameContacts(old, new []Contact) error // "contact 213: old cirrine-5/ILV/00445 dish 1, new cirrine-5/ILV/00445 dish 2"
func benchWindows(b *testing.B)             // sub-benchmarks old/hours=6 ... new/hours=24

// chapter 10
func collect(seq iter.Seq[PassView], n int) []PassView
func stationReport(code string, contacts []Contact) string // "CWR contacts: 76", then "cirrine-1/CWR/19862 dish 1" a contact
func appendContact(dst []byte, c Contact) []byte
func exportStations(w io.Writer, byStation map[string][]Contact, codes []string) error
func getScratch() *[]PassView
func putScratch(p *[]PassView) // keeps none over maxScratch (65536) passes
type passKey struct {          // a pass's identity without formatting its ID
	SatID, Station string
	AOS            int64 // Unix seconds
}
func dedupe(passes []PassView) []PassView
func planStationFat(s Site, cands []PassView, turnaround time.Duration) []Contact
func planStationLean(s Site, cands []PassView, turnaround time.Duration) []Contact

// chapter 11
func fieldLines(v any) []string // "Slots int int"
func tagLines(v any) []string   // "Slots: min=1,max=64"
func describe(v any) []string   // "Slots=12"
func checkRule(name string, v reflect.Value, rule string) error
func misfit(name, rule string, t reflect.Type) error             // "Config.Slots: rule len=3 does not fit int"
func measure(v reflect.Value, text string) (int, any, any, bool) // cmp.Compare's result, v and the bound for %v; false when the bound doesn't fit v
func applyDefaults(ptr any) error
func setText(v reflect.Value, text string) error // a duration, an int, a float or a string, from text
func Validate(v any) error

// chapter 12
var passModel Predictor = predict // passmodel.go: the pass model planStation asks; tests replace it
func planStation(ctx context.Context, n *Network, s Site, w Window, free []time.Time, stows []Stow, cfg Config) ([]Contact, int, error)
func Plan(ctx context.Context, n *Network, w Window, free map[string][]time.Time, stows []Stow, cfg Config) (*Schedule, error)
func handover(cur *Schedule, at time.Time, turnaround time.Duration) map[string][]time.Time
func replanFrom(ctx context.Context, n *Network, cur *Schedule, at time.Time, stows []Stow, cfg Config) (*Schedule, error)
// Event is news for the replanner: a tick, or a stow that begins or ends.
type Event struct {
	At   time.Time
	Stow *Stow // nil for a tick; a stow with a zero To begins, one with a To ends it
}
func newReplanner(n *Network, pub *Publisher, cfg Config, now func() time.Time) *Replanner
func (r *Replanner) Run(ctx context.Context, events <-chan Event) error
func heard(stows []Stow, st Stow) []Stow // 12.7: st begun, or its station's open stow ended, in a new slice
func takeOver(ctx context.Context, n *Network, old *Schedule, at time.Time, events <-chan Event, cfg Config, now func() time.Time, fly func(Contact)) error
```

- `Plan` runs `planStation` for every station in a group of `cfg.Slots`, merges the twenty
  contact lists with `MergeBy`, and returns a schedule whose `Version` the caller sets; its
  first error cancels the other stations and is returned.
- `replanFrom` keeps `cur`'s contacts with AOS before `at`, plans from `at` with `handover(cur,
  at, cfg.Turnaround)`, and returns both together in a new slice, by AOS then ID, over `cur`'s window; its `Dropped` counts the passes from `at` on. `handover` gives a station an entry for each dish up to the highest its contacts before `at` used; a dish with no entry is free now, as `assign` reads `free`.
- A `Replanner` takes the next run number for every event (a tick or a stow), cancels the run in
  progress, plans from the event's instant with every stow it has heard of, and offers the
  result; a run whose contacts equal the current schedule's offers nothing. Its runs are numbered after the publisher's current version; `Run` also replans every `cfg.Replan` on a ticker of its own, plans each run through `planWithin(ctx, cfg.RunLimit, …)`, waits for its runs before returning (when `events` closes, the last one still finishes), and leaves word in its `news` channel (room for one) after each schedule it publishes, which the station runner reads.
- `takeOver` replans a copy of `old` given passplan 2's window (`cfg.Days` days from `at`) from `at` and publishes it as version 1; it runs a `Replanner` on `events` and flies each contact of the newest schedule at its AOS, the first after the last one flown in `byAOS` order (starting from a pass at `at` with no ID), waking on the `Replanner`'s `news` to read the schedule again; once `ctx` ends it waits for `Run` and returns `ctx.Err()`.

#### How lessons test it

- Everything stays in the process: no network, no files outside `t.TempDir()`, no real clock in
  a check. Clocks are injected or come from a synctest bubble; randomness is
  `rand.New(rand.NewPCG(1, 2))`; real sleeps are 50 ms at most.
- Concurrency checks are correct by construction: limits are measured with an `atomic.Int64`
  of jobs in flight and its highest value (a compare-and-swap loop), and parallelism with a
  barrier (fake jobs that wait until the expected number are in flight, with a one-second
  guard), never by timing.
- A **leak check** takes `runtime.NumGoroutine()` before the code runs and, after cancelling,
  polls every millisecond for up to a second until the count is back; a leak fails as
  `goroutines left after cancel = 3, want 0`. (A leaked goroutine inside a synctest bubble
  ends the whole test binary with a deadlock panic, so leak checks never use a bubble.)
- Every wait in a test has a guard of about a second, so a wrong starter fails within a second
  or two and never hangs the run. The guard is for what right code does at once: a wait on real work, such as 2.8's day of predictions (up to 0.4 s under `-race` on a busy machine), gets five seconds of its own, and the guard starts after it. A wait on a loop that can run on without yielding also caps memory: 5.9's walk fails once it has allocated 64 MB (a right walk of a day allocates 64 KB), since a loop that appends without pulling takes gigabytes within the guard's second, and the table stops at the first walk still running. 12.9's station runner test ends the day after 3,000 flights, more than a day holds, so a runner that flies one contact over and over stops. A check that can fail a wrong starter only by waiting out its guard (12.2's twelve stations at once) runs first, with a guard of its own, and the checks after it run only once it passes, so the starter fails in about a second.
- Tests plan a day or a week, lesson 12.8 included. The model watches each satellite's elevation,
  about 4 µs a pass on one core: a day of every pass (7,485) takes 35 ms, a week of servable
  passes 0.2 s, a hundred days 3 s, so a hundred days and a quarter stay story numbers.
- Where learners write checks (chapters 7 and 8) or benchmarks (chapter 9), their code goes in an
  editable non-test file and the read-only test calls it, as Season 1's chapter 15 did.
- Messages read `call = got, want want`; schedules are compared by pass ID and dish:
  `Plan(28 June)[289] = pipit-2/YST/22439 dish 2, want pipit-2/YST/22439 dish 1`.

**Topics added this season** (front matter): `iterators`, `performance`, `reflection`. Fuzzing
and synctest lessons use `testing`; atomics use `goroutines` and `mutexes`.

**New field guide words** (world words only; Go words are taught in readmes): passplan,
passplan 2, planner slot, candidate (pass), contact, dropped pass, turnaround, window, shadow
run, cut-over, hand-over, take-over day, host, wave, rideshare, deployment, elements, first
contact, telemetry feed, the drill, dish loader, pass archive, Thistle-2, planning run, test Board, station runner. Carried over and
needed again: AOS, LOS, pass, pass ID, LEOP, stow, mask, max elevation, catalog, the Board, the
Fishbowl, booking, bookingd, tlmd, UTC, benchmark (Season 6's), band (Season 1's), deploy (Season 4's).

**Name checks** to run before the season ships: the fifteen station names. A web search on
6 November 2026 found nothing real under Tirrin Rock, Varsk Island, Calidra Atoll, Ilha
Brumosa, Ostrel Island, Ilvane Spit, Cape Solvane, Saltreach Island, Corith Point, Nimrie
Point, Hasku Point, Ilha Calmaria or Punta Ventisca (near misses recorded: Cape Solander in
Australia, Kaskulu Point in Alaska, Torrin on Skye, a restaurant called Ventisca in Punta
Arenas, a Calmaria guesthouse in Portugal); Quindle Island and Vorrel Island are not yet
checked. The satellites are numbered or lettered members of names already checked (Thistle,
Cirrine, Brine, Furrow); "passplan 2" is the bible's own service with a number. No new people
or customers.

### The capstone: Pass Day, Live

After take-over day, the learner starts the capstone on Monday 2 July 2029. Their new badge
says **Ground Segment Lead**, still No. 039, based at Port Ellery headquarters, since Monday
2 July 2029. That morning Tomás tells the Fishbowl that the investors who put money into Apogee
in March will watch the whole network for one week, from Monday 23 July 00:00 to Monday 30 July
00:00 UTC, **the investors' week**, and want to see it run without a single missed pass: every
booked pass flown. Ingrid: "Then we fly that week first. Every pass of it, before they see one."
The learner builds **the rehearsal**: one Go program that runs the whole ground segment
(passplan's planning, uplinkd's commands, tlmd's telemetry, dlinkd's files and deliveries, the
pager) against **the simulated world**, twenty stations, eighty satellites and the customers'
receivers, and runs the investors' week in simulated time, in seconds. Nkem reviews every change,
as she did on the first day; Ingrid owns the week's go/no-go; Sunniva brings the storm; Mirela
the keys; Yusra and Hollis the bookings and the files; Joss and Professor Marr Thistle-2's first
science; Tomás the investors.

Story time runs from Monday 2 July to Friday 20 July 2029: six chapters of two to five days.
The investors' week itself, 23 to 29 July, is the finale's debrief. No new people, places,
satellites or customers are introduced. Kessick Polytechnic, which has booked Thistle-2's passes
since March (`bk-7001`), gets a webhook receiver for the first time. The hosts stay unnamed.

In plain words, for the opening: before anyone outside Apogee watches the network for a week,
the team flies that same week inside a computer, every pass and every file, with storms and
broken lines put in on purpose, and checks that nothing is missed.

**Where the passes come from.** Every pass in this section and in the lessons (pass ID, AOS, LOS,
maximum elevation) comes from Season 7's pass model, `model.go`: Season 3's reviewers' settled
model ported to Go, on Season 7's `network.go`, whose elements are Season 3's fits for
Pipit-1, Cirrine-1, Brine-3, Furrow-A and Furrow-B, Season 4's for Pipit-2 and Season 7's own for
the rest. The hidden test file `sim_passes_test.go` holds a copy of `model.go`. It gives exactly
what Season 3's model gives (all 933 passes of the five fitted in Season 3 over the five original
stations in the investors' week, and all 213 of Pipit-2's, to the second and the tenth of a
degree; `cmd/passes`'s tests hold it to those counts), and, as in Season 7, a pass shorter than two minutes doesn't count. No pass is invented,
and no pass is longer than its orbit allows. The rehearsal itself never predicts a pass: it reads them from
passplan's pass archive (below).

#### The capstone, chapter by chapter

| Ch | Days (2029) | Where | Who needs what | What is at stake |
|---|---|---|---|---|
| 1 | Mon 2 – Tue 3 Jul | the backend desks | Tomás and Ingrid: the investors' week planned for twenty stations from passplan's pass archive, and checked against Yusra's export of the week's bookings | A booked pass with no dish is lost before the week begins. Two of the export's 1,495 bookings name passes the plan can't fly: Hollis's Furrow-A over Cape Wren at 04:58:30, a time copied from an old email, and a booking for Thistle, which re-entered in July 2027. Yusra fixes both by Tuesday noon |
| 2 | Wed 4 – Thu 5 Jul | the Fishbowl | Ingrid: every pass on the Fishbowl's screen as it starts and ends, the whole week in seconds so she can run it again after every change, and a stop button | Calidra Atoll's host router answers a call and never says a word. A pass that outlives its LOS holds its dish, as the stuck goroutines held Lantern Atoll's workers in February |
| 3 | Fri 6 – Mon 9 Jul | the backend desks; Sunniva by link from Ystmark | Sunniva and Joss: every satellite's telemetry from every station's radio, frame by frame, and the newest frame of each | Thistle-2 still sends `batt_v=nan` when its sensor glitches; Varsk Island's host line flips a bit now and then; Ystmark's fiber dies without a sound in storms, and the backhaul takes ten minutes to come up |
| 4 | Tue 10 – Wed 11 Jul | Mirela's desk | Mirela: every command signed with its satellite's own key, Cirrine's sealed, numbered one at a time per satellite and never run twice. The investors' security reviewer asked to see a key replaced during the live week, so `k-2029-07` replaces `k-2028-10` at 00:00 on Wednesday 25 July | A satellite that misses its new keys can't be commanded once `k-2028-10` retires on 1 August |
| 5 | Thu 12 – Fri 13 Jul | Yusra's desk | Yusra and Hollis: every file reassembled, checked and delivered by signed webhook within a minute; Joss: Kessick's first receiver, for Thistle-2's first science | The investors will ask customers how fast their data arrives. On Tuesday 24 July Furrow-K's image loses three chunks at Vorrel Island, and Cirrine's receiver restarts at 09:32 for its own deploy |
| 6 | Mon 16 – Fri 20 Jul | the Fishbowl | Sunniva's forecast on Monday: a storm at Ystmark on Thursday 26 July. Ingrid: logs, numbers and pages through it, and the week's report at the full rehearsal on Friday 20 July, 15:00 | "Every pass counted, or we don't show them." |

**The opening** (Monday 2 July): Tomás's announcement in the Fishbowl and the investors' week; the
badge; the cast: Nkem (reviewer, as on the first day), Ingrid (the week's go/no-go), Tomás (the
investors), Sunniva (Ystmark's weather); the simulated world, which Nkem grew from her twenty-feed
drill of March; and the first assignment: say what week the rehearsal flies, and plan it. Words of
the first week: rehearsal, the investors' week, the simulated world, simulated time.

**Chapter 2 ends** on Thursday 5 July, wind 12 at Ystmark. Sunniva: "Wind 12. Every call on time. Radios still waiting for a hello."

**Chapter 3 ends** on Monday 9 July, wind 9 at Ystmark. Sunniva: "Wind 9. Fiber cut at 03:31, as scripted. One Pipit-2 pass lost. The rest flew the backhaul."

**Chapter 4 ends** on Wednesday 11 July. Every command goes up signed or sealed, one at a time per satellite, and none runs twice; in the rehearsal all eighty satellites hold `k-2029-07` by Monday 09:27:31. Mirela: "LGTM once it fails closed. It does."

**Chapter 5 ends** on Friday 13 July. In the rehearsal every file of the week lands whole and reaches its customer: Furrow-K's image at 08:12:30 on Tuesday, Thistle-2's first science at Kessick at 16:04:25 on Thursday. Yusra: "TL;DR: 1,484 files, 1,484 events, none lost."

**The finale** (chapter 6, Friday 20 July, 15:00): the full rehearsal of the investors' week with
Ingrid's script (below). The fiber at Ystmark dies at 03:31 on Thursday and one pass fails; the
dishes stow in the wind from 10:50 to 12:00, during the new Furrows' highest passes of the week;
the pager pages No. 039 once, at 10:56, and resolves at 16:56; Thistle-2's first science file
reaches Kessick at 16:04; the week's report says NO-GO: nine booked passes lost, all Furrowlight's,
all at Ystmark. Hollis: "NINE FURROWS. ONE STORM." Yusra moves each booking to its Furrow's next free
contact elsewhere; the second run says GO. The debrief closes the series: on Thursday 26 July the
storm comes as Sunniva said, Ystmark stows at 10:50, the pager rings once, and by 00:03 on Monday 30 July every booked pass of the week has flown (the last, `bk-9493` over Lantern Atoll, sets at 00:02:42). Ingrid writes one line in the log: "Every pass
counts. Every pass counted."

**Who is on call**: Season 6's rota, a week each from Monday 08:00 UTC, except that Nkem swapped
weeks with No. 039 so the rehearsal's author carries the pager for the investors' week. As `[]Shift`
(the Board's type), in the read-only `rota.go` from lesson 6.3:

| From (Monday 08:00) | On call |
|---|---|
| 2029-07-09 | Tomás Ferreira-Lund |
| 2029-07-16 | Nkem Halloran |
| 2029-07-23 | No. 039 |
| 2029-07-30 | Mirela Dunn |

#### New things

- **The rehearsal**: the program the learner builds. It is the ground segment's own code, run
  against the simulated world instead of the real one; in production `main` hands it the real
  clock, TCP to each station's radio and HTTPS to customers. In readmes "the rehearsal", a plain
  word, never a program name.
- **The simulated world**: Nkem's program that plays every station's radio, every satellite and
  every customer's receiver, with whatever trouble Ingrid writes into it. In lessons it is the
  hidden test files (`sim_*_test.go`); readmes say in a sentence what it does in that lesson.
- **The investors' week**: Monday 23 July 00:00 to Monday 30 July 00:00 UTC.
- **The pass line**: the connection the rehearsal holds with a station's radio for one pass, and
  the frames it carries (below).
- **Simulated time**: the week's clock, which the tests run in a `testing/synctest` bubble, so a
  week passes in a moment of real time.
- **Kessick's receiver**: the webhook receiver Joss wrote for Thistle-2's files.

#### The rehearsal's world

Each lesson is one `package main` program (module `lesson`). Its only fixed entry point is
`rehearse`; everything else is the learner's design, and readmes never require a name inside it.
The read-only `world.go`, given from lesson 1.1, holds `World`; `rehearse`, below, is the learner's, in `rehearse.go`:

```go
// World is everything outside the rehearsal: main builds the real one, and the tests
// build a simulated one.
type World struct {
	Dir  string                                                      // the rehearsal's folder: its inputs, and the plan and reports it writes
	Now  func() time.Time                                            // the clock: every instant comes from it
	Dial func(ctx context.Context, station string) (net.Conn, error) // a station's radio, for the pass line (chapter 2)
	HTTP *http.Client                                                // customers' webhook receivers (chapter 5)
	Out  io.Writer                                                   // the console: one line per event, safe for many goroutines
}

// rehearse runs the week rehearsal.json describes against w and returns when every pass,
// file, delivery and report is done; when ctx ends first it stops and returns
// context.Cause(ctx).
func rehearse(ctx context.Context, w World) error
```

- **The clock rule.** Every instant comes from `w.Now()`: AOS, LOS, a landing, an alert's
  minute. Every wait is a duration measured from it, `c.AOS.Sub(w.Now())`, waited with the `time`
  package: timers, `time.Sleep`, `context.WithTimeout(ctx, c.LOS.Sub(w.Now()))`. Never
  `time.Until`, `context.WithDeadline` or `time.Now()` with a story instant: inside the tests'
  bubble the `time` package's clock starts at 2000-01-01, so a 2029 instant is 29 years away.
  Connection deadlines are the exception: `conn.SetReadDeadline(time.Now().Add(90 *
  time.Second))`, an instant on the `time` package's own clock, never `w.Now()`.
- **The console.** Every line is one `Write` of a whole line ending in `\n`, so lines from many
  goroutines never mix. Times on it are `Mon 15:04:05` from `w.Now()`.
- **Names.** Every top-level name in the test files begins with `sim` (hidden files) or `test` /
  `Test` (`main_test.go`), so the learner's own names never collide with them; lesson 1.1's readme
  says so in one sentence.
- **Given read-only**: `world.go` (1.1); `main.go` (1.1), which dials `radio.<code>.apogee.internal:7301`; `types.go`, the bible's canonical types the rehearsal
  uses, unchanged (Band, Station, Orbit, Mode, Satellite, Pass, TelemetryFrame), Season 3's
  `PassView`, `Booking`, `Landing` and `Event`, and Season 7's `Contact`, `Window` and `Stow` (1.1);
  Season 7's `network.go` (`Site`, `Elements`, `Network`, `newNetwork`) and `plan.go` (`byAOS`,
  `canServe`, `stowedDuring`, `assign`), unchanged (1.1 and 1.2); `telemetry.go`, Season 7's strict
  `parseLine` and `formatLine` (3.2); `passline.go`, the frame types, NAK codes and limits below as
  constants with doc comments (3.1); `rota.go` (6.3). Season 7's `model.go` is not given: the
  rehearsal reads its passes from the archive.

#### The rehearsal's folder

| File in `w.Dir` | What | From |
|---|---|---|
| `rehearsal.json` | the settings, read strictly | 1.1 |
| `passes.txt` | passplan's pass archive for the window | 1.2 |
| `bookings.json` | bookingd's export of the window's bookings | 1.3 |
| `uplink-keys.txt` | the uplink masters, Season 5's format (mode 0600) | 4.1 |
| `customers.json` | each customer's webhook URL and secret (mode 0600) | 5.2 |
| `plan/<YYYY-MM-DD>.txt` | written: the day's contacts | 1.2 |
| `reports/unplanned.txt` | written: booked passes the plan can't fly | 1.3 |
| `reports/health.txt` | written: each satellite's newest frame | 3.2 |
| `files/<file>` | written: each downlink, whole | 5.1 |
| `reports/incomplete.txt`, `reports/undelivered.txt` | written: files still missing chunks, events never delivered | 5.1, 5.3 |
| `logs/uplinkd.log`, `logs/dlinkd.log`, `logs/pager.log` | written: JSON log records | 6.1 (`pager.log` 6.3) |
| `reports/metrics.json`, `reports/alerts.txt`, `reports/week.txt` | written: the numbers, the pages, the week's report | 6.2, 6.3, 6.5 |

Every file the rehearsal writes is saved whole through a temporary sibling and a rename (Season
2's way), folders created as needed, and written even when empty. When `ctx` ends during the run,
the reports written after the run (health, incomplete, undelivered, metrics, alerts, week) are not
written. A missing `bookings.json` means no bookings; every other missing input is an error that
wraps `fs.ErrNotExist`. Input errors name the file: `bookings.json: json: unknown field "sat"`. A second value is `bookings.json: more data after the bookings` (`customers.json: ... the customers`); `uplink-keys.txt`'s are `uplink-keys.txt: line 3: <problem>`; the archive's name their line, `archive line 3: want 4 fields, got 3`. The logs are written into temporary siblings as the week runs and renamed into place when it ends, stopped or not; a stopped pass has no end record and counts in no metric.

#### Settings

`rehearsal.json`, one JSON object read strictly (Season 2's rules: an unknown field, a second value
or anything after the object is an error):

```json
{"week":"2029-07-23","days":7,"stations":["YST","TIR"]}
```

- `week` is required, a date `2006-01-02`: the rehearsal starts at 00:00 UTC that day (any day,
  not only a Monday). `days` is 1 to 7, default 7. The **window** is `week` 00:00 to `days` × 24
  hours later, From included, To excluded (Season 7's `Window`). `stations` lists the station codes
  to rehearse, in any order; empty or absent means all twenty. The rehearsal works only with the
  rehearsed stations, always in Season 7's table order.
- Errors, exactly: `rehearsal.json: json: unknown field "dayz"`; `rehearsal.json: more data after
  the settings`; `rehearsal.json: week is required`; `rehearsal.json: week "2029-7-23" is not a
  date like 2006-01-02`; `rehearsal.json: days 8 is not from 1 to 7`; `rehearsal.json: unknown
  station "YSX"`. On any input error `rehearse` returns it and prints nothing.
- The first console line: `rehearsal: 2029-07-23, 7 days, 20 stations`; one day or one station
  takes the singular (`rehearsal: 2029-07-26, 1 day, 2 stations`).

#### The pass archive and the plan

`passes.txt` is Season 7's pass archive: one pass per line, `<pass ID> <AOS> <LOS> <max elev>`,
times RFC 3339 UTC, elevation with one decimal, by AOS then pass ID; every predicted pass of every
satellite over each rehearsed station (any band, any elevation from the mask up, two minutes or
longer) with AOS in the window. A bad line is an error, `archive line 3: <problem>`, as Season 7's `ReadPasses` said.

```
thistle-2/PEL/02136 2029-07-26T16:03:46Z
    2029-07-26T16:09:59Z 11.7
```

(In the file that is one line; readmes split it.)

**The rule**, Season 7's, unchanged: a **candidate** is an archived pass at a rehearsed station
that can serve its satellite (`canServe`) with a maximum elevation of at least 10°. Each station is
planned on its own, candidates in AOS order then pass ID (`byAOS`): a pass takes the
lowest-numbered dish free at its AOS, a dish being free once its last contact's LOS plus 2 minutes
has passed; a pass no dish can take is **dropped** (`assign`, every dish free at the window's
start). The plan knows no stows: weather is found during the week, never planned around.

- `plan/<day>.txt`, one file per day of the window: the contacts whose AOS falls on that UTC day,
  by AOS then pass ID, as Season 7's dish loader lines: `16:03:46 16:09:59 1 thistle-2/PEL/02136`
  (AOS, LOS, dish, pass ID). A contact whose LOS is after midnight stays in its AOS's day.
- Console: `plan: 14297 contacts, 28680 dropped`.
- **The investors' week** (all twenty stations, seven days): 52,596 archived passes; 42,977
  candidates; **14,297 contacts** (2,028 to 2,053 a day); 28,680 dropped.

#### Bookings

`bookings.json` is bookingd's export (Season 3's file, the store's `bookings.json`): one JSON array
of Season 3's `Booking`, read strictly. A booking **matches** the contact with the same satellite,
station and AOS to the second (Season 6's rule); bookings at stations the rehearsal doesn't fly, or
with AOS outside the window, are left out of everything. An unmatched booking is **unplanned**,
with the first of these reasons that fits:

| Reason | When |
|---|---|
| `unknown satellite` | its `sat_id` isn't in the network (`thistle`: re-entered in July 2027) |
| `unknown station` | its `station` isn't in the network |
| `no such pass` | no archived pass of that satellite over that station has that AOS |
| `not servable` | the station and satellite share no band (Pipit-1, UHF only, over Ystmark) |
| `too low` | the pass's maximum elevation is under 10° |
| `no dish free` | the rule dropped it |

- `reports/unplanned.txt`: one line per unplanned booking, by AOS then booking ID:
  `bk-9902 Thu 04:58:30 furrow-a/CWR no such pass` (booking ID, AOS as `Mon 15:04:05`,
  satellite/station, reason). Empty when every booking matched.
- Console: `bookings: 1493, unplanned 0`.
- **The week's bookings**, `bk-8001` to `bk-9493` (after Season 7's `bk-7999`), numbered in the
  plan's order (AOS, then pass ID): every contact of a Furrowlight, Brinewatch or Cirrine
  Atmospherics satellite with a maximum elevation of at least 70°, and every Thistle-2 contact from
  its first science pass on. 1,493 bookings: Furrowlight 713, Brinewatch 421, Cirrine 345, Kessick
  Polytechnic 14. Apogee books none of its own passes.
- Chapter 1's export (Tuesday 3 July) held two more, which Yusra removes: `bk-9901` for `thistle`
  (`unknown satellite`) and `bk-9902`, Furrow-A over Cape Wren at Thursday 04:58:30 (`no such
  pass`). Tests plant one booking per reason, from the archive: Pipit-1 over Ystmark
  (`pipit-1/YST/17218`, Thursday 01:09:32, `not servable`), Cirrine-1 over Punta Ventisca
  (`cirrine-1/PVT/20279`, Thursday 00:17:43, 8.2°, `too low`), Furrow-A over Quindle Island
  (`furrow-a/QND/24951`, Thursday 00:02:58, 12.4°, `no dish free`), and station `SQX` (`unknown
  station`).

#### Running the week

Every contact flies on its own clock: the reference program starts a goroutine per contact that
waits until its AOS, and a pass's context ends at its LOS. Contacts at different stations, and at
a station's two dishes, overlap freely.

- **Console**, every contact exactly two lines: at its AOS `Thu 16:03:46 AOS thistle-2/PEL/02136
  dish 1`, and at its LOS `Thu 16:09:59 LOS thistle-2/PEL/02136 flown`. The LOS line always
  comes at LOS, whatever happened earlier in the pass, with the outcome: `flown`, `failed
  <reason>` or `skipped stowed`; or at the moment the rehearsal stops, `stopped`.
- **Outcomes** (each chapter adds what it checks): a pass is **flown** when its station welcomed it
  (from 3.1; in chapter 2, when its banner arrived), every command it sent was acknowledged (from
  4.1) and its line wasn't down at LOS. Otherwise, the first that fits: **skipped** `stowed` (the
  station refused the pass at HELLO: its dishes were stowed); **failed** `refused` (a NAK at HELLO other than stowed); `no_answer` (no banner, a wrong banner, every dial failed or neither a WELCOME nor a NAK before LOS); `refused` (a command refused); `stowed` (the station cut the pass with BYE `stowed`, or refused a redial as stowed);
  `line_lost` (lost and not back before LOS); `no_ack` (a command still unanswered at LOS).
- **The end**: after the last contact's LOS and the last delivery, `rehearse` writes its reports and
  prints `week done: 14282 flown, 3 failed, 12 skipped`, then returns nil.
- **Stopping**: when `ctx` ends, no pass dials again, every call open hangs up at once and prints
  its LOS line `stopped`, deliveries stop, and once every goroutine has ended the console's last
  line is `rehearsal stopped at Tue 12:00:00` and `rehearse` returns `context.Cause(ctx)`.
- **Chapter 2's runs** (the settled model): the investors' week's first contact is `furrow-r/ICM/00384` at Monday 00:00:06, and its last LOS is `furrow-m/ILV/00490`'s at Monday 30 July 00:07:45, when `rehearse` returns. On Thursday at Ystmark, in the week's plan and in a Thursday-only one, dish 2 takes Furrow-M (`furrow-m/YST/00436`, 10:46:20–10:55:32) while dish 1 still flies Furrow-N (`furrow-n/YST/00436`, 10:42:25–10:51:36). Calidra Atoll and Port Ellery fly 500 and 614 contacts in the week; with Calidra Atoll silent, `week done: 614 flown, 500 failed, 0 skipped`. Ystmark and Tirrin Rock on Monday and Tuesday, stopped at Tuesday 12:00:00, have three passes open (`furrow-g/TIR/00406`, `furrow-f/YST/00406`, `furrow-f/TIR/00406`) and none starting or ending then; the next AOS is 12:00:58.
- **Chapter 3's runs**: on Thursday at Ystmark `pipit-2/YST/22861` (03:28:42–03:37:44) is the only contact in progress when the fiber dies at 03:31; its last frame before then is its TLM frame of 03:30:52, so the dead line is noticed at 03:32:22, and the next contact there, `cirrine-1/YST/20282` at 04:25:02, is the first to fly the backhaul. At Varsk Island on Thursday `cirrine-5/VSK/01487` flies 00:42:43–00:47:38 (its third TLM frame at 00:44:53) and `pipit-2/VSK/22860` 01:39:45–01:48:41; the station relays 747 TLM frames that day, and in the investors' week its flips spoil 13 of its 5,280. At Tirrin Rock `pipit-2/TIR/22864` (08:15:11–08:23:23) has dish 1 and `furrow-b/TIR/24954` (08:17:25–08:25:39) dish 2, no other candidate within four minutes. Thistle-2 flies only at Port Ellery and Lantern Atoll (10 and 15 contacts in the week); its first frame at Port Ellery on Thursday is `thistle-2 1879776236 2029-07-26T16:03:56Z batt_v=7.66 temp_c=8.4 mode=NOMINAL`, and the one made at 16:06:56 reads `batt_v=nan`. A Thursday at Ystmark, Tirrin Rock and Port Ellery never hears brine-6, -10, -15, -18, -23, -31, -33 and -38; the week hears all eighty. Ilha Brumosa's shortest Thursday pass is `furrow-j/BRM/00442` (21:28:38–21:34:16) and Corith Point's first `brine-20/CRP/01080` (00:54:03–00:59:02).
- **Chapter 4's runs**: Tuesday at Port Ellery, before `k-2029-07` activates, has 93 contacts, 17 of them Cirrine's; Pipit-2's first is `pipit-2/PEL/22832` at 05:28:22, and Thistle-2 flies `thistle-2/PEL/02106` (16:45:15–16:53:57) and `thistle-2/PEL/02107` (18:19:14–18:26:53). Wednesday at Port Ellery has 88. On Thursday Furrow-O rises over Ystmark at 18:26:27 (dish 2) and over Tirrin Rock at 18:26:29, and Furrow-L over Tirrin Rock at 18:38:01 and over Ystmark at 18:38:10; no other satellite's two passes there begin within 12 s. In the week five pairs of one satellite's contacts begin less than 2 s apart: brine-24 at CRP and LAN (Sat 13:11:23), brine-3 at CSL and PEL (Tue 19:43:44), brine-32 at LAN and CRP (Wed 14:17:54/55), cirrine-1 at YST and TIR (Sat 15:50:06/07), furrow-w at TIR and YST (Sun 18:29:54/55). Cirrine-7's only dish-1 pass at Ystmark on Thursday is `cirrine-7/YST/01488` (01:51:00–01:57:56). Over the backhaul `cirrine-1/YST/20282`'s PING, sent at 04:25:12, reaches Ystmark at 04:25:12.35 and runs at 04:25:14.35; a bad length loses the line at 04:25:12.6, and the copy is refused as replayed at 04:25:16.35. The week's last LOAD_KEY runs at Mon 09:27:31 (`thistle-2/LAN/02086`, 09:27:17–09:33:19), and no command runs within 5 s of Wed 00:00:00.
- **Chapter 5's runs**: the investors' week lands 1,483 files until lesson 5.4 (`reports/incomplete.txt`: `furrow-k-20290724-0740.tif missing 21 22 23`) and all 1,484 from it, with 1,484 events; each file named below is the only landing in its second. The first is `cirrine-6-20290723-0003.dat` (`bk-8001`, ILV, 6,144 bytes) at Monday 00:04:17: `evt-000001`, timestamp `1879459457`, signature `bfeca88736851964481525450c7153269ce7e2d67a93666937780cf880046127`. Until 5.4 Cirrine-5's event is `evt-000288` and Kessick's `evt-000782` (signature `d2f673e98b5e195ab686c78182693b38487c53632de5cdf753c9d6b145310aa5`); from 5.4 Furrow-K's is `evt-000284` (station VRL; its DUMP_FILE runs on `furrow-k/VSK/00404` at 08:12:24), Cirrine-5's `evt-000289` and Kessick's `evt-000783`. On their first pass 362 of the week's files have two chunks swapped and 319 one chunk twice (Thistle-2's first science file: chunk 6). A Thursday-only plan at Port Ellery drops `brine-5/PEL/01080` (`bk-8674`); among the booked contacts it keeps are `cirrine-8/PEL/01488` (`bk-8676`, 01:32:40–01:42:59, its file landing at 01:33:35), `brine-14/PEL/01083` (`bk-8697`, 05:38:10–05:47:29), `furrow-a/PEL/24958` (`bk-8739`, 10:44:10–10:53:23), `cirrine-8/PEL/01496` (`bk-8782`, 14:47:42–14:58:00) and `brine-3/PEL/18374` (`bk-8798`, 17:06:06–17:15:25). After its first science pass Thistle-2's next contacts are `thistle-2/LAN/02138` (Thursday 20:12:03–20:18:39, `bk-8840`; a DUMP_FILE sent there runs at 20:12:17) and `thistle-2/LAN/02146` (Friday 08:02:49–08:08:36, `bk-8930`). Four contacts work out their commands while another station still brings their satellite's file down: `furrow-u/TIR/00441` (Thursday, AOS 18:03:28, during `furrow-u/YST/00441`'s file over the backhaul), `furrow-k/TIR/00471` and `furrow-j/TIR/00471` (Saturday), `furrow-e/TIR/00486` (Sunday). The week's shortest contact is `furrow-j/LAN/00461` (3 min 14 s), the shortest booked `thistle-2/LAN/02153` (3 min 18 s); on `furrow-t/LAN/00431` (Thursday 02:50:16–02:53:35) a file lands at 02:51:35 and its event's 8th attempt always comes after LOS.
- **Chapter 6's runs**: the investors' week writes 28,595 records to `logs/uplinkd.log` (14,297 `pass started`, 14,282 `pass flown`, 3 `pass failed`, 12 `pass skipped`, one `line lost`, `pipit-2/YST/22861` at Thursday 03:32:22), 2,968 to `logs/dlinkd.log` (1,484 `file landed`, 1,484 `event delivered`) and 5 to `logs/pager.log`; `frames_total` 124,056, `bad_frames_total` 13, `bad_lines_total` 27, `files_landed_total` 1,484, `ack_wait_seconds` 14,239 in ≤ 2 and 127 in ≤ 5 (count 14,366, `sum_ms` 28,820,900); delivered Brinewatch 421, Cirrine 345, Furrowlight 704, Kessick 14. Booked, in table order: YST 210, PEL 74, SQT 82, LAN 59, CWR 87, TIR 110, VSK 123, CAL 44, BRM 60, OST 64, ILV 66, CSL 56, SRH 43, CRP 43, NIM 70, HSK 50, ICM 41, QND 48, PVT 51, VRL 112; the second run YST 19, TIR 16, CSL 13, ICM 7, QND 8. Thursday at YST, VSK and PEL with the storm, `flipEvery` VSK 400 and Brinewatch's receiver down: 381 contacts, 242 answers in ≤ 2 and 127 in ≤ 5, 14 Brinewatch events undelivered, 1 bad frame, 1 bad line, 46 files. A one-day plan of Monday 23 July at YST, PEL, NIM and VRL flies `pipit-2/YST/22817` 05:25:47–05:35:30, `furrow-b/NIM/24906` 05:27:08–05:35:18, `furrow-z/VRL/00387` 05:28:35–05:35:02, `furrow-n/VRL/00388` 07:48:38–07:58:16, `furrow-o/NIM/00388` 07:50:05–07:59:18 and `cirrine-10/PEL/01451` 14:07:54–14:17:51. A Thursday-only plan at Port Ellery holds 11 of the week's bookings. Lantern Atoll (one dish) flies 532 contacts in the week, 59 booked (Calidra Atoll 500 and 44), so neither reaches the 100 bookings at which 99.0% first allows a lost pass. The last booked pass, `bk-9493` (`brine-36/LAN/01141`, Sunday 23:54:36 to Monday 00:02:42), sets after the window closes.

#### The pass line

The rehearsal calls a pass's station at its AOS (`w.Dial(ctx, code)`); in production a TCP
connection to the station's radio, in the tests one end of a `net.Pipe` inside the bubble. The
station answers each call on its own, so a station takes several at once (two dishes, redials).
The protocol is Season 4's controller protocol in shape: a text banner, then framed binary messages.

**Banner.** On answering, the station at once sends one ASCII line, `<CODE> READY PROTO 1\r\n`
(`PEL READY PROTO 1\r\n`, 19 bytes). A reader accepts LF with or without the CR before it and
refuses any other shape or version (`bad banner "BRM READY PROTO 2"`). The banner and the frames
that follow are read through one `bufio.Reader` for the life of the call: the station may write the
banner and its first frame in one write.

**Frames**, both directions, exactly Season 4's layout: length uint16 (the bytes after it, 9 + the
body's length, 9 to 1024), type uint8, seq uint32, body 0 to 1015 bytes, CRC-32 (IEEE) of the type,
seq and body; every number big-endian. A frame whose CRC doesn't match is dropped and counted, and
reading goes on with the next frame; a length under 9 or over 1024 means the stream is out of step:
the line is lost.

| Type | Name | From | Seq | Body |
|---|---|---|---|---|
| 0x01 | HELLO | rehearsal | 0 | `<pass ID> <dish>` in ASCII: `thistle-2/PEL/02136 1` |
| 0x81 | WELCOME | station | 0 | 9 bytes: the dish (uint8), then `last` (uint64): the satellite's last accepted command number |
| 0x20 | TLM | station | 0 | one telemetry line, ASCII |
| 0x30 | UPLINK | rehearsal | the command's number | one command line, ASCII: a CMD or SEAL line |
| 0x90 | ACK | station | the command's | empty: the satellite accepted and ran it |
| 0x91 | NAK | station | the refused frame's | a code (uint8), then the reason in ASCII |
| 0x40 | FILE | station | 0 | the file's chunk count (uint16), then its name in ASCII |
| 0x41 | CHUNK | station | 0 | the chunk's number (uint16, from 0), then its data, at most 1013 bytes |
| 0x42 | END | station | 0 | the whole file's CRC-32 (uint32) |
| 0xFF | BYE | station | 0 | the reason in ASCII: `los` or `stowed`; then the station closes |

CHUNK and END belong to the last FILE frame on that call. **NAK codes**: 1 `unknown pass` (no such
contact at this station and dish between its AOS and LOS); 2 `stowed`; 3 `bad crc`; 4 `hello
first`; 5 `unknown type`; 6 a command refused, the satellite's own error text as the reason
(Season 5's: `command failed authentication`, `replayed sequence number: 413, last 413`, `unknown
key: "k-2029-07"`, and `command must be sealed` from a Cirrine satellite given a CMD line).

**What a station does during a contact**, in the simulated world (lessons may rely on it):

1. Answers HELLO at once: WELCOME, or NAK 2 while its dishes are stowed, or NAK 1.
2. Relays a TLM frame at AOS + 10 s and every 60 s after, while before LOS.
3. Answers each command 2 s after it arrives (ACK or NAK), one at a time, in the order received.
4. On a booked pass, relays the satellite's file: FILE at AOS + 30 s, the k-th chunk k seconds after it, END one second after the last chunk. Chunks asked for again (DUMP_FILE) come first: FILE again, then the
   asked chunks, from 2 s after the command's ACK, and the pass's own file waits until they are done.
5. At LOS sends BYE `los` and closes; when its dishes stow during a call, BYE `stowed` and closes.
6. On a new call during a pass (a redial), sends FILE again before continuing a file in progress;
   frames it relayed while the line was dead are gone.

**Limits.** Every read must bring at least one byte within 90 s (an idle deadline, moved before
each read, never one deadline for a whole frame: over Ystmark's backhaul a frame trickles in over a
second). A line is **lost** on that timeout, on the call ending before LOS without BYE, or on a bad
length; the rehearsal closes it and dials again at once, then waits Season 4's `redialDelay` between
failed attempts (1 s doubled per failure, at most 30 s, times a factor drawn from [0.5, 1.0)), until
a call is welcomed or LOS comes. A redial says HELLO for the same pass and dish.

**Frames to copy** (hex):

| Frame | Bytes |
|---|---|
| HELLO `thistle-2/PEL/02136 1` | `00 1e 01 00 00 00 00 74 68 69 73 74 6c 65 2d 32 2f 50 45 4c 2f 30 32 31 33 36 20 31 ba 2c bf 54` |
| WELCOME dish 1, last 412 | `00 12 81 00 00 00 00 01 00 00 00 00 00 00 01 9c 38 1e 8d 3e` |
| WELCOME dish 1, last 0 | `00 12 81 00 00 00 00 01 00 00 00 00 00 00 00 00 d8 bc 63 10` |
| ACK seq 413 | `00 09 90 00 00 01 9d 80 81 0e b5` |
| NAK `stowed` | `00 10 91 00 00 00 00 02 73 74 6f 77 65 64 fa 7a 83 76` |
| NAK `unknown pass` | `00 16 91 00 00 00 00 01 75 6e 6b 6e 6f 77 6e 20 70 61 73 73 ea f5 e2 7f` |
| FILE 8 `thistle-2-20290726-1603.dat` | `00 26 40 00 00 00 00 00 08 74 68 69 73 74 6c 65 2d 32 2d 32 30 32 39 30 37 32 36 2d 31 36 30 33 2e 64 61 74 ac d1 42 c3` |
| END crc `1a2b3c4d` | `00 0d 42 00 00 00 00 1a 2b 3c 4d 66 4f 65 07` |
| BYE `los` | `00 0c ff 00 00 00 00 6c 6f 73 29 cd bc 25` |
| BYE `stowed` | `00 0f ff 00 00 00 00 73 74 6f 77 65 64 90 f6 b8 05` |

The mistakes these catch are Season 4's: a CRC that also covers the length, a length written
little-endian.

#### Telemetry

A TLM line is a satellite's housekeeping frame as Season 7's strict `parseLine` reads it (given in
`telemetry.go`); a line it refuses is dropped and counted. In the simulated world a satellite's
frame number is the Unix second it was made, so the newest frame is the one with the highest
number, whatever station relayed it and in whatever order; a satellite over two stations at once
arrives twice. About one Thistle-2 frame in seven reads `batt_v=nan`.

`reports/health.txt`: one line per satellite of the network, in fleet order: `formatLine` of its
newest frame that parsed, or `<id> never heard`. Thistle-2's first TLM frame of its science pass,
relayed at AOS + 10 s, as the simulated world makes it:

```
thistle-2 1879776236 2029-07-26T16:03:56Z
    batt_v=7.66 temp_c=8.4 mode=NOMINAL
```

(One line on the wire and in the file; readmes split it.)

#### Commands

`uplink-keys.txt` is Season 5's file: `#` lines and blank lines ignored, each other line `<id>
<activate> <retire or -> <64 hex secret>`, optionally `revoked`. For the investors' week:

```
k-2028-07 2028-08-11T00:00:00Z 2028-10-24T00:00:00Z 7570…6573
k-2028-10 2028-10-17T00:00:00Z 2029-08-01T00:00:00Z 7570…6573
k-2029-07 2029-07-25T00:00:00Z - 7570…6573
```

The test masters are ASCII, 32 bytes: `k-2028-07` `uplink test master key, 32 bytes` and `k-2028-10`
`uplink next master key, 32 bytes` (Season 5's), and the new **`k-2029-07`** `uplink 2029 master
key, 32 bytes` (hex `75706c696e6b2032303239206d6173746572206b65792c203332206279746573`).
Mirela made `k-2029-07` on Monday 9 July with Activate Wednesday 25 July 00:00, and set
`k-2028-10`'s Retire to its Activate plus seven days (Season 5's overlap rule).

- **Keys from keys**, Season 5's rule: a satellite's key is `hkdf.Key(sha256.New, master,
  []byte("apogee uplink keys v1"), purpose+"|"+satID, 32)`, purpose `sign` or `seal`.
- **Signed**: Season 1's CMD line, `CMD <sat> <seq> <OP> [args...] sig=<64 hex>`, the HMAC-SHA256
  under the satellite's `sign` key of the canonical form `<sat>|<seq>|<OP>|<args...>` (Season 5's
  canonical rules). **Sealed**: Season 5's SEAL line, `SEAL 2 <sat> <seq> <key id> <24 hex nonce>
  <hex ciphertext>`, AES-256-GCM under the satellite's `seal` key, a fresh 12-byte nonce from
  `crypto/rand` for every seal, associated data `SEAL|2|<sat>|<seq>|<key id>`, plaintext
  `<OP>|<args...>`. A resent line (below) is the same bytes, not a new seal.
- **Who gets what**: every command to a Cirrine Atmospherics satellite (all thirteen) is sealed,
  as Cirrine's contract says; LOAD_KEY is always sealed; every other command is a CMD line.
- **A contact's commands**, sent from AOS + 10 s once the station has welcomed the pass, each after
  the previous one's answer: `PING`; then `LOAD_KEY|<id>|<64 hex sign key>|<64 hex seal key>` for
  each master the satellite doesn't hold yet (148 bytes for `k-2029-07`); then `DUMP_FILE <file>
  [chunk...]` for each of the satellite's incomplete files (chapter 5).
- **Numbers**: strictly increasing per satellite. The next command takes `max(next, last+1)`, where
  `last` comes from the WELCOME and `next` is one more than the rehearsal's own last number for that
  satellite (Season 4's rule). **One command at a time per satellite across the whole network**: a
  satellite over two stations at once (common at the polar stations) never has two commands in
  flight, so its numbers reach it in order. The lock is a channel with room for one, taken in the same `select` that reads the call: inside the tests' bubble a goroutine waiting on a `sync.Mutex` isn't durably blocked, so a mutex held across a round trip stops simulated time for good; the harness then fails the test and says where the mutex was asked for.
- **Never twice** (Season 4's rule): when the line is lost with a command in flight, the new
  WELCOME's `last` says what happened: `last` at or above its number means it ran, and counts as
  acknowledged; otherwise the same line is sent again, same number. A NAK `replayed sequence number` to that line means its first copy reached the satellite after the line dropped and ran: it counts as acknowledged (Season 4's NAK `duplicate` rule). A command still unanswered at
  LOS is settled the same way on the satellite's next contact: done if `last` reached it; a PING is
  then dropped, and a LOAD_KEY or DUMP_FILE still needed is sent again, numbered afresh.
- **Which master**: at the start of the window each satellite holds every master activated before
  the window and not yet retired (`k-2028-10` in the investors' week). A master it doesn't hold, not
  retired and not revoked, goes up as LOAD_KEY on its first contact, sealed under the newest master
  it holds; it holds it once that LOAD_KEY is acknowledged. Signing and sealing use the newest master
  that has activated, hasn't retired, isn't revoked and that the satellite holds (Season 5's
  `signingKey`, plus "holds"). So in the investors' week every satellite gets `k-2029-07` on its
  first contact (all eighty hold it by Monday 09:27:31; the last is Thistle-2, whose first contact, at Lantern Atoll, begins at 09:27:17), commands use `k-2028-10` until Tuesday
  23:59:59 and `k-2029-07` from Wednesday 00:00:00, and `k-2028-07` is never used.
- **The satellites** in the simulated world check every line by Season 5's order (parse, version,
  satellite, key held and valid, number above `last`, signature or seal, only then record), try
  their keys newest first for a CMD line, refuse a CMD line to a Cirrine (`command must be sealed`),
  and keep `last` across contacts. Their numbers at the start of the window are the simulated
  world's (`pipit-2` 6314, `cirrine-1` 9712, `thistle-2` 412 in the examples); the others: pipit-1 5518, pipit-3 2907, furrow-a 8841, furrow-b 8790, brine-3 7426, and for each 2029 satellite 40 plus its ID's FNV-1a hash mod 300.
- **Vectors** (Go 1.27.3): `CMD pipit-2 6315 PING sig=0cb110af1256712defcb30409a53e47e402b944c6a13ab67a236315160034810`
  (`k-2028-10`, sign key `514e8010…22ad`, Season 5's); with Season 5's fixed nonce (vectors only),
  `SEAL 2 cirrine-1 9713 k-2029-07 a1b2c3d4e5f60718293a4b5c 840b65d7152c90505bbbec4f7dad949a2e2a4fee`
  and `SEAL 2 cirrine-1 9713 k-2028-10 a1b2c3d4e5f60718293a4b5c fc5a75921a95c7faaf06adb91ebc66ea2d2d379c`.
  Derived from `k-2029-07`: `sign|pipit-2` `d504962086a135791928b0886641af259de287658a5c43fe7f12c9fd5e24e49e`,
  `seal|cirrine-1` `6e309f1cf6854eb92c5f0e6dba01e4af5b118407769d43591479121913a5a930`,
  `sign|thistle-2` `f1502f436e73bef882e08a0dca866619827b0aa33b7916d563da9981b14c84e0`,
  `seal|thistle-2` `2b8e335f4ef5a1161f14d014bd8c56df6d20b1dcce34035bf6806c98bc3ac459`; from
  `k-2028-10`: `seal|thistle-2` `c158b5a257a7752abb87a211edf6dc49e0c6cd045bd4c786cbea77267f2249f9`,
  `seal|cirrine-7` `b632863e687b1ca1dc88e6126d31fae2542000ce5d203c634cdb676448f7555a`.

#### Files

On a booked pass the satellite sends its file down: FILE (the chunk count and the name), CHUNK
frames, END with the whole file's CRC-32. The name is Season 3's, the satellite, the AOS day and the
AOS time: `<sat>-<YYYYMMDD>-<HHMM>.<ext>`, `.tif` for Furrowlight's satellites and `.dat` for every
other: `furrow-k-20290724-0740.tif`, `thistle-2-20290726-1603.dat`. Every chunk holds 256 bytes;
a Furrowlight file has 48 (12,288 bytes), Cirrine's 24 (6,144), Brinewatch's 16 (4,096), Kessick's
8 (2,048). In the simulated world chunks sometimes come out of order or twice.

- **Reassembled**: chunks are kept by number (Season 1's `Downlink`). At END, a file with every
  chunk whose bytes' CRC-32 equals END's **lands**: it is written whole as `files/<name>` (never a
  partial file), the console says `Tue 09:32:18 LANDED cirrine-5-20290724-0931.dat`, and its
  `Landing` (Season 3's: booking ID, file, bytes, the landing instant truncated to the second) goes
  to delivery. A file missing chunks waits for them; one whose chunks are all there but whose CRC
  doesn't match is dropped and asked for whole.
- **Asked again** (5.4): on the satellite's next contact (any station), after its PING, `DUMP_FILE
  <name> <n> <n>...`, the missing numbers ascending, or the name alone for the whole file; it lands when the asked chunks have come down. A file is asked for only while no pass is bringing it down: a pass holds a file from its FILE, or from its own DUMP_FILE's run, until the file's END or the pass's LOS; a DUMP_FILE an earlier contact left unanswered holds nothing, so the next contact asks again.
- `reports/incomplete.txt`, at the end: one line per file still incomplete, by name: `<name>
  missing 21 22 23` or `<name> crc mismatch`. Empty in the investors' week.

#### Deliveries

`customers.json`, read strictly: one JSON array of `{"id","name","webhook_url","webhook_secret"}`.
A booking's customer is the one whose `name` is its satellite's operator. The five, with their test
secrets (Season 3's, and Kessick's new one):

| `id` | `name` | `webhook_url` | `webhook_secret` |
|---|---|---|---|
| `apogee` | Apogee | `https://hooks.apogee.example/apogee` | `apogee webhook secret` |
| `furrowlight` | Furrowlight | `https://hooks.furrowlight.example/apogee` | `furrowlight webhook secret` |
| `brinewatch` | Brinewatch | `https://hooks.brinewatch.example/apogee` | `brinewatch webhook secret` |
| `cirrine` | Cirrine Atmospherics | `https://hooks.cirrine.example/apogee` | `cirrine webhook secret` |
| `kessick` | Kessick Polytechnic | `https://hooks.kessick.example/apogee` | `kessick webhook secret` |

`.example` is reserved, so no address is real. The secret is the text's bytes.

- **The event and its request**, exactly Season 3's webhook: `POST` to the customer's URL, the
  `Event` as `json.Marshal` writes it (type `data.landed`), `Content-Type: application/json`,
  `Apogee-Timestamp: <Unix seconds of this attempt, from w.Now()>`, `Apogee-Signature: <64 hex>`, the
  HMAC-SHA256 under the customer's secret of the timestamp, a dot and the exact body. Event IDs are
  `evt-` and six digits, counting from `evt-000001` in the order files land (files landing in the
  same second may take either order). Kessick's first:
  `{"id":"evt-000783","type":"data.landed","booking_id":"bk-8792","sat_id":"thistle-2","station":"PEL","file":"thistle-2-20290726-1603.dat","bytes":2048,"landed_at":"2029-07-26T16:04:25Z"}`,
  at timestamp `1879776265`, signs as `340bfc1649ec9f7fae412dddb919258dad72fc043e8bf0f45b2c8b06fd03c62e`: the 783rd file to land in the investors' week, the only one at 16:04:25 (the 782nd, `evt-000782`, until lesson 5.4 lands Furrow-K's Tuesday image). The event's `station` is the booking's, even when the file lands on a later contact elsewhere.
- **Season 3's schedule**: the first attempt as soon as the file lands; retried on a transport
  error, 408, 429 and 5xx; at most 8 attempts, after failed attempt n waiting 2 s × 2^(n−1), capped
  at 5 minutes, times a factor drawn uniformly from [0.5, 1.0); every attempt signed again with its
  own timestamp, the event ID unchanged; any other answer, and the 8th failure, make the event
  **undelivered**. The receivers in the simulated world refuse a timestamp more than 300 s from their
  clock and a signature `hmac.Equal` rejects, and drop an event ID they already have.
- Console: `Tue 09:34:02 DELIVERED cirrine-5-20290724-0931.dat` when a receiver answers 2xx;
  `UNDELIVERED <file>` when it never will. `reports/undelivered.txt`, at the end: `<file>
  <customer> <last answer>`, the last status code or `error`. The week ends only when every
  delivery has finished.

#### Logs, numbers and pages (chapter 6)

**Logs**, Season 6's conventions: `slog.NewJSONHandler` at INFO, one record per event, `time` from
`w.Now()` (never the wall clock), durations as whole milliseconds in `_ms` keys, the `pass` group
written with `slog.Group`, and nothing secret anywhere (no master's text or hex, no derived key, no
webhook secret; tests search every log for each).

| File (`program`, `site`) | Message | Level | Attributes besides `program` and `site` |
|---|---|---|---|
| `logs/uplinkd.log` (`uplinkd`, the station) | `pass started` | INFO | `pass` group (`sat_id`, `station`, `aos`, `los`), `pass_id`, `dish` |
| | `pass flown` | INFO | `pass`, `pass_id`, `max_elev_deg`, `frames`, `commands` |
| | `pass failed` | ERROR | `pass`, `pass_id`, `max_elev_deg`, `reason`, `err` |
| | `pass skipped` | WARN | `pass`, `pass_id`, `reason` |
| | `line lost` | WARN | `pass_id`, `err` |
| `logs/dlinkd.log` (`dlinkd`, `HQ`) | `file landed` | INFO | `file`, `bytes`, `booking_id`, `pass_id` |
| | `file not saved` | ERROR | `file`, `err` |
| | `event delivered` | INFO | `event_id`, `customer`, `attempts` |
| | `event undelivered` | ERROR | `event_id`, `customer`, `attempts`, `err` |
| `logs/pager.log` (`pager`, `HQ`) | `alert sent` | INFO | `rule`, `station`, `kind`, `value`, `to` |

`frames` counts the pass's TLM frames with a good CRC and `commands` the commands it sent that its satellite ran; a file's `pass_id` is the pass it landed on, `customer` the customer's ID and `attempts` the requests sent for the event.

**Numbers**: `reports/metrics.json`, one JSON object written at the end, in Season 6's shapes (a
counter by two labels is `{"<a>":{"<b>":n}}`; the histogram is Season 6's `HistogramSnapshot`):
`passes_total` by satellite and outcome (`flown`, `failed`, `skipped`); `failures_total` by
satellite and reason; `frames_total` (TLM frames with a good CRC); `bad_frames_total` (frames of any
type with a bad CRC); `bad_lines_total` (TLM lines `parseLine` refused); `commands_total` by
satellite and `accepted` or `refused`; `files_landed_total`; `deliveries_total` by customer and
`delivered` or `undelivered`; `ack_wait_seconds` from a command sent to its answer, bounds 0.5, 1, 2,
5, 10, 30, 60, 120, 300 (an answer 2 s after its command goes in the ≤ 2 bucket; over Ystmark's
backhaul answers take 2.7 s and land in ≤ 5).

**Pages**, Season 6's pager from 8 January 2029, with its two pass rules, per rehearsed station,
evaluated every minute on the minute of `w.Now()`; from the window's first minute until the week is done, each minute counting every pass whose LOS is at or before it; a pass's result counts from its LOS; skipped and
stopped passes are not failures:

| Rule | Value per station | FireAt | ClearBelow | For | Repeat | Severity |
|---|---|---|---|---|---|---|
| `pass-failures` | failed passes with LOS in `(now−6h, now]` | 2 | 1 | 0 | 24h | Critical |
| `failed-pass` | failed passes with LOS in `(now−1m, now]` | 1 | 1 | 0 | 0 | Warning |

Critical goes to whoever is on call at the alert's minute (the rota above; a RESOLVED to the same person, unlike Season 6's pager, so the one woken hears it is over), Warning to `morning list`. `reports/alerts.txt`, one line per alert in time order, the same
minute by station (table order) then rule name: `Thu 10:56 FIRING pass-failures YST 2 to No. 039`
(minute as `Mon 15:04`, kind, rule, station, value as `%g`, recipient).

**The week's report**, `reports/week.txt`, Season 6's report in a week's numbers:

```
2029-W30  booked  lost  avail  budget  left
YST          210     9  95.71       1    -8
```

- The title is the window's first day's ISO week (`2029-W30`). Header `%-8s %7s %5s %6s %7s %5s`;
  one row per rehearsed station in table order, `%-8s %7d %5d %6.2f %7d %5d`; then `all`, the
  columns summed and its availability from the sums. **Booked** counts the station's bookings in
  the window (matched or not); **lost** is booked minus flown, so unplanned, failed and skipped
  passes all count (Season 6: "every pass counts"); availability is flown over booked in percent, `-`
  (`%6s`) with nothing booked; **budget** is Season 6's integer formula, `booked × (10000 − target)
  / 10000`, target 9900 at Lantern Atoll and Calidra Atoll (palms and rain fade) and 9950 elsewhere;
  **left** is budget minus lost.
- Then one line per lost booking, by AOS then booking ID: `bk-8743 Thu 11:05:54 furrow-h/YST
  skipped stowed` (booking ID, AOS as `Mon 15:04:05`, satellite/station, outcome; unplanned ones as
  `unplanned <reason>`).
- Last, the verdict: `GO: every booked pass flown`, or `NO-GO: 9 booked passes lost` (`NO-GO: 1 booked pass lost` for one).

#### The investors' week, as Ingrid scripts it

The finale's world, and every fact chapters 2 to 6 may use. All times UTC, 2029.

| When | Where | What the simulated world does |
|---|---|---|
| Mon 23 Jul 00:00 | everywhere | the window opens; every satellite holds `k-2028-10` |
| all week | Varsk Island | about one TLM frame in 400 arrives with a flipped byte (bad CRC) |
| all week | Thistle-2 | about one frame in seven reads `batt_v=nan` |
| Tue 24 Jul 07:40:40 | Vorrel Island | `furrow-k/VRL/00403` (`bk-8279`, 84.8°): chunks 21, 22 and 23 of `furrow-k-20290724-0740.tif` never arrive. Its next contact, `furrow-k/VSK/00404` at 08:12:10, asks for them; the file lands at 08:12:30 |
| Tue 24 Jul 09:32:00 – 09:33:00 | Cirrine's receiver | restarts for its own deploy: every request fails. `cirrine-5-20290724-0931.dat` (`bk-8289`, `cirrine-5/CAL/01463`) lands at 09:32:18 and is delivered on its 6th or 7th attempt |
| Wed 25 Jul 00:00:00 | everywhere | `k-2029-07` activates |
| Thu 26 Jul 03:31:00 | Ystmark | the fiber dies: open calls go silent, new calls fail until the backhaul is up at 03:41:00; then, until the fiber is back at 19:40:00, every message takes 0.35 s longer each way (Season 4's 0.7 s round trip) and the station's arrive in pieces of at most 64 bytes, 0.25 s apart. `pipit-2/YST/22861` (03:28:42–03:37:44) can't get back before LOS: failed `line_lost` |
| Thu 26 Jul 10:50:00 – 12:00:00 | Ystmark | wind lock: both dishes stow. `furrow-n/YST/00436` (10:42:25–10:51:36) and `furrow-m/YST/00436` (10:46:20–10:55:32) are cut: failed `stowed`. The twelve contacts with AOS from 10:54:10 to 11:56:46 are skipped, nine of them booked |
| Thu 26 Jul 16:03:46 | Port Ellery | `thistle-2/PEL/02136` (`bk-8792`, 11.7°, LOS 16:09:59): Thistle-2's first science file, `thistle-2-20290726-1603.dat`, lands at 16:04:25 and goes to Kessick's receiver. (Its higher pass at 17:36:21 goes to Pipit-2, which takes the station's one dish three minutes earlier.) |

**What the rehearsal reports**, with the settled model and these bookings:

- `week done: 14282 flown, 3 failed, 12 skipped`.
- `reports/alerts.txt`: `Thu 03:38 FIRING failed-pass YST 1 to morning list`; `Thu 10:52 FIRING
  failed-pass YST 1 to morning list`; `Thu 10:56 FIRING failed-pass YST 1 to morning list`; `Thu
  10:56 FIRING pass-failures YST 2 to No. 039`; `Thu 16:56 RESOLVED pass-failures YST 0 to No. 039`.
- `reports/week.txt`: YST 210 booked, 9 lost, 95.71, budget 1, left −8; every other station 0 lost;
  `all` 1493, 9, 99.40, 1, −8; the nine lost lines; `NO-GO: 9 booked passes lost`. The nine:
  `bk-8743` furrow-h 11:05:54, `bk-8746` furrow-g 11:09:49, `bk-8750` furrow-e 11:17:39, `bk-8753`
  furrow-z 11:29:23, `bk-8755` furrow-y 11:33:18, `bk-8757` furrow-w 11:41:07, `bk-8759` furrow-v
  11:45:02, `bk-8762` furrow-t 11:52:51, `bk-8763` furrow-s 11:56:46, all `skipped stowed` at
  Ystmark. (Pipit-2's `pipit-2/YST/22866`, 11:21:07, is skipped too; Apogee books none of its own
  passes.)
- `reports/incomplete.txt` and `reports/undelivered.txt` empty.
- **Yusra's moves**, Friday 20 July: each of the nine to its Furrow's next unbooked contact that day at another station of the second run: furrow-h to `furrow-h/CSL/00437` 12:52:04, furrow-g to `furrow-g/ICM/00436` 11:28:56,
  furrow-e to `furrow-e/CSL/00436` 11:31:23, furrow-z to `furrow-z/CSL/00437` 11:42:34, furrow-y to
  `furrow-y/ICM/00437` 11:52:47, furrow-w to `furrow-w/CSL/00437` 11:53:55, furrow-v to `furrow-v/QND/00438` 14:22:57, furrow-t to `furrow-t/TIR/00437` 11:55:14, furrow-s to
  `furrow-s/QND/00438` 14:34:16. The second run, Thursday only at YST, CSL, ICM, TIR and QND (63 bookings, 605 contacts), with the same storm: `GO: every booked pass flown`. (Not Saltreach Island: a plan that starts at 00:00 Thursday gives its one dish to a different train of Furrows in the first hour and drops `bk-8669`.)

#### How the lessons test it

- **One entry point.** Every test builds a `World` and calls `rehearse`, inside
  `synctest.Test`: the simulated world's clock is the bubble's, shifted to the window's start
  (`now := func() time.Time { return time.Now().Add(shift).UTC() }`, Season 7's way), `Dial` returns
  one end of a `net.Pipe` whose other end a simulated station serves, `HTTP` is an `*http.Client`
  (Timeout 5 s, never following redirects) whose `Transport` hands each request to the customer's
  simulated receiver in memory, and `Out` records each line with the instant it was written. Nothing
  reaches the network; there is no loopback inside a bubble (Seasons 4 and 7).
- **The pass archive** is written by the test from Season 7's model (`sim_passes_test.go`), only
  for the stations and days a case rehearses, and cached per test binary: the full investors' week (52,596 passes) takes about 0.075 s on ten cores, 0.22 s on one and 1.1 s under `-race` (0.28 s for the four stations of a `-race` run); a day at
  two stations, 4 ms on one core.
- **Guards.** The harness runs `rehearse` in a goroutine and waits for it until the window's end
  plus one hour of simulated time; past that it fails the test (`rehearse still running at Mon 30
  Jul 01:00`), cancels `ctx`, closes every call and waits for `rehearse` to return, so a starter that
  waits forever on a silent station fails cleanly instead of deadlocking the bubble. Once `rehearse` returns, the world sleeps until that same instant (at least ten simulated minutes) before the bubble ends, so passes a rehearsal left asleep wake and are checked instead of panicking the bubble; a learner's goroutine stuck on a channel of its own can still panic it, and starters never do. A watch outside the bubble, in real time, ends the one hang the guard can't see: when the console has been quiet for half a second and, at two looks in a row, no goroutine in the bubble can run and one waits for a `sync.Mutex`, it fails the run's test (`simulated time stopped after Tue 19:43:51 while a goroutine waited for a sync.Mutex at line.go:152: a mutex held across a wait stops the clock`) and ends the test binary.
- **The simulated world's knobs**, set per case (each defaults to off): the bookings, keys and
  customers files; `silent` (stations that answer and never send a byte); `banner` (a station's
  banner text); `dialFail` (a station's calls fail during a window, all week without one); `split` and `glue` (messages written in pieces of at most n bytes, or, at a station with one contact in progress, the banner and that contact's WELCOME written in one write as the call is answered, before HELLO); `late` (a station relays each TLM frame late, so a satellite's frames arrive out of number order); `flipEvery` (a station's every n-th TLM frame arrives with a flipped byte); `corrupt` (a pass's
  n-th TLM frame with a flipped byte); `badLength` (a frame with length 8 after a pass's n-th
  frame); `quiet` and `backhaul` (a station's dead-line window and its slow window); `dropBefore` and
  `dropAfter` (a pass's call dropped just before the satellite gets a command, or just after it runs
  it); `stows` (a station's wind-lock window); `missing` (chunks a file never sends on its first
  pass); `badCRC` (a file whose END carries a wrong CRC); `receivers` (a customer's receiver down,
  or answering a status, during a window); `lastSeq` (a satellite's `last` at the start).
- **What tests read**: the console with its instants; the files in `w.Dir`; and the simulated
  world's records: each call (station, pass, dish, dialed and hung-up instants, who hung up), the raw
  bytes it received, every command each satellite ran or refused (number, operation, key ID,
  instant, station), every chunk sent, and every delivery attempt (customer, event, instant, signature checked). Lesson 4.3 also fixes `crypto/rand`'s stream with `testing/cryptotest.SetGlobalRandom` for one small run and checks that every nonce is a piece of it: a nonce from `math/rand` seals as well as one from `crypto/rand`, and `cipher.NewGCMWithRandomNonce` draws from `crypto/rand` without importing it.
- **Size.** Each lesson's tests fly the whole network for the whole week at most once: one cached run of the investors' week as Ingrid scripts it, with its 1,493 bookings (`simWeekRun`), shared by every check that needs it from lesson 1.1 to 6.5; every other case is one to six stations for one or two days, except 2.3's silent Calidra Atoll, which flies the week beside Port Ellery (1,114 contacts). On an Apple-silicon Mac (M5) that week takes 1.24 s with every chapter's features (0.89 s without bookings), the bubble pinned to one CPU, which is faster there than ten. Under `-race` the four-station week takes 2.4 s, plus 0.28 s for its archive, and
  `scripts/pkgcheck.sh` runs every test three times that way, so the test files include Season
  7's `race_on_test.go` and `race_off_test.go`, and with `raceOn` the full-week cases rehearse four
  stations (YST, TIR, PEL, VRL) and check invariants instead of golden files.
- **Golden files** (6.5): `testdata/week.golden`, `alerts.golden`, `health.golden`, written by the
  solution and compared byte for byte; `metrics.json` is decoded and compared with the simulated
  world's own counts. Console lines written in the same instant are compared as a set.
- Every wait is in simulated time; no real sleep anywhere. Concurrency checks are correct by
  construction: per-satellite order and once-only execution are read from the satellites' records,
  never inferred from timing. Messages read `call = got, want want`:
  `LOS line for thistle-2/PEL/02136 at Thu 16:09:58, want Thu 16:09:59`.

#### For writers: the reference program

The readmes leave the design to the learner, but each starter is the previous lesson's official
solution, so the official solutions share one design. The calibration writer builds the hidden
files (`sim_world_test.go`, the harness; `sim_station_test.go`, stations and satellites;
`sim_passes_test.go`, the settled model, the archive and the reference plan; `sim_hooks_test.go`,
the receivers) and the whole reference program first, then writes chapter 1; each writer cuts their
lessons from the reference program, lesson n's solution being it with lesson n's features and no
later ones. Its shape: `rehearse.go` (the order of the work), `settings.go`, `archive.go`,
`week.go` (the plan and its files), `bookings.go`, `run.go` (a goroutine per contact, the console),
`wire.go` (frames and the banner), `line.go` (one pass on the pass line, redials), `health.go`,
`uplink.go` (keys, signing, sealing, numbering under one lock per satellite), `downlink.go`,
`deliver.go`, `logs.go`, `metrics.go`, `pager.go`, `report.go`, and `main.go`, which builds the
real world and is never run by a lesson.

**Topics** (front matter): the lessons use the older ones (`context`, `crypto`, `files`,
`goroutines`, `http`, `json`, `logging`, `metrics`, `mutexes`, `networking`, `security`, `streams`,
`testing`, `time`...); none is added.

**New field guide words** (world words, in plain language): Ground Segment Lead, the rehearsal,
the simulated world, simulated time, the investors' week, pass line, radio (a station's), go/no-go
for the week, Kessick's receiver, science file. Carried over and needed again: pass, contact,
candidate, dropped pass, pass archive, booking, AOS, LOS, banner, frame, CRC, backhaul, fiber, stow,
wind lock, master key, derived key, sealed command, LOAD_KEY, overlap, retire, webhook, landing, on
call, page, the morning list, error budget, night shift, uplinkd, dlinkd, the Fishbowl, the Board.

**Name checks** (web search, 10 November 2026): "Ground Segment Lead" is a common job title in
real space-industry postings, a generic title like "ground station engineer", not a company,
product or person, so it stays; nothing real was found for "the investors' week" (only generic
"Investor Week" events), "pass line" as a satellite term (elsewhere a generic phrase) or the
customer ID `kessick`; the receiver hosts use the reserved `.example` domain, and the customer
names were checked before. The master `k-2029-07` follows Season 5's pattern. No new people,
places, satellites or customers.

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
