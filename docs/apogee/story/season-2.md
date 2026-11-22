# Apogee story bible: Season 2

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
