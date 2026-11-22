# Apogee story bible: the capstone

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
