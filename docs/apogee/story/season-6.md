# Apogee story bible: Season 6

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
- **The first record and panics** (chapter 1, from the deploy of Friday 10 November). A
  program's first record is INFO `program started`, naming its build from
  `debug.ReadBuildInfo`; Ingrid's deploy checklist reads it at 14:01. uplinkd prepares each
  pass (its booking, its lines, its pass ID) inside `guardPass`, outside `uplink.go`: a panic
  there is recovered, so the station flies on, and ends the pass with its `pass failed` line,
  `err` `pass <id>: panic: <value>` and the whole stack under `stack`. Until that Friday the
  guard printed the stack to standard error with `debug.PrintStack`. Nkem's bench test, Friday
  at 09:30: the pass ID `cirrine-1/LAN`, its orbit number lost. A panic in the goroutine
  `runPass` starts is not recovered by the guard.

| Program | Message | Level | Its own attributes |
|---|---|---|---|
| uplinkd (station) | `pass started` | INFO | `pass` group |
| | `worker taken`, `lines sent`, `ack received` | DEBUG | `lines`, `ack_ms` |
| | `frame received` | DEBUG | `bytes`, `hex` |
| | `carrier lost` | WARN | `margin_db` |
| | `pass flown` | INFO | `pass` group, `max_elev_deg`, `wait_ms`, `ack_ms` (from chapter 3) |
| | `pass failed` | ERROR | `pass` group, `max_elev_deg`, `wait_ms`, `err`; after a panic while the pass is prepared (`guardPass`, chapter 1), `pass_id`, `err` and `stack` only |
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
| any (from 10 Nov) | `program started` | INFO | `go_version`, `module`, `module_version`: its first record (chapter 1) |
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
| `stack` | a recovered panic's stack, `string(debug.Stack())` taken in the deferred function that recovered: one value, its line breaks written `\n` |
| `go_version`, `module`, `module_version` | the build, from `debug.ReadBuildInfo`: `GoVersion` (`go1.27.3` in lessons), `Main.Path`, `Main.Version` (a tag such as `v1.8.0`; `(devel)` outside version control and in every test binary) |
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
| 1 | `newLogger`, `logFlown`, `logResult`, `msAttr`, `accessLog`, `logStart`, `guardPass`, `newPassLog`, `passLine` (`log.go`); given: `dropTime`, `stamp` (`clock.go`) |
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
- **Build info and stacks** (1.8, 1.9): the first record is compared with the test binary's
  own `debug.ReadBuildInfo` (under `go test`: `Path` `lesson.test`, `Main.Path` `lesson`,
  `Main.Version` `(devel)`, even inside a repository) and with `runtime.Version()`, never a
  fixed release. A stack is checked for its first word, `goroutine`, and for the function
  that panicked, which a test binary names `lesson.orbitOf`, not `main.orbitOf`.
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
