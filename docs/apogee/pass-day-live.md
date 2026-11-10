# Pass Day, Live: syllabus

**Authoring status:** in progress. The simulated world and the whole reference program are
built, and the investors' week runs in them exactly as the bible scripts it; every lesson's
solution is a cut of that program. Chapter 1 is in the library and is still to be read in the
app; chapters 2–6 are still to be written. Once lessons are written, the lessons in
`courses/apogee.sqlite` are the source of truth.

The capstone of the Apogee series, after Season 7; story time 2 to 20 July 2029, with the
investors' week of 23 to 29 July as its finale. The world, the characters, the code and the
writing rules are in `story.md`, whose section "The capstone: Pass Day, Live" fixes every contract
the tests check: the rehearsal's world and folder, the settings, the pass archive and the plan,
the bookings, the pass line, the commands, the files, the deliveries, the logs, numbers, pages and
report, and the investors' week as Ingrid scripts it.

This is the plan the lessons are written from. Once they are written, the lessons in
`courses/apogee.sqlite` are the source of truth. `go run ./cmd/author export` writes them out to
read.

- **Title:** Pass Day, Live
- **Slug:** `pass-day-live` · **Kind:** Guided Project
- **Short description:** Rehearse Apogee's whole ground segment for a week in one Go program,
  against twenty simulated stations and eighty satellites.
- **Description:** Apogee's investors will watch the whole network for a week and want every
  booked pass flown, so Ingrid wants that week flown first. You build the rehearsal: one program
  that plans the week for twenty stations from passplan's archive and the customers' bookings, runs
  every pass on its own clock, takes telemetry in over the wire, signs and seals every command and
  replaces a key mid-week, reassembles every file and delivers it by signed webhook, and logs,
  counts and pages through a storm at Ystmark, ending in the week's report. Everything runs in
  simulated time: a week in seconds. Standard library only.
- **Story:** a short opening, each chapter's briefing and debrief, and the field guide, told as
  `docs/storytelling.md` says.

## How the capstone differs

The learner builds one program, a lesson at a time: each lesson adds one feature, its starter is
the previous lesson's solution, and its read-only tests check the whole program, so every earlier
feature stays checked. The capstone has the fewest instructions of any course: a readme says what
must work and what the tests check, shows an output line or a frame where that settles a question,
and leaves the design to the learner. Hints carry the help. The only fixed name is `rehearse(ctx,
w World) error`; every Exercise below says what `rehearse` must do, never how.

Every test builds a `World` and runs `rehearse` inside a `testing/synctest` bubble against the
simulated world: stations' radios on `net.Pipe`, customers' receivers behind an in-memory
transport, the bubble's clock shifted to the window's start. A week passes in a moment of real
time; nothing reaches the network; no test sleeps. The tests read the console with the instant of
each line, the files `rehearse` writes, and the simulated world's records of every call, command,
chunk and delivery.

## At a glance

| # | Chapter | Slug | What gets built | Go | Lessons | tests | quiz |
|---|---|---|---|---|---|---|---|
| 1 | The Week's Plan | `the-weeks-plan` | strict settings, the pass archive, Season 7's rule for twenty stations, bookings matched to the plan | `encoding/json`, `bufio`, `slices`, `iter`, files saved whole | 4 | 3 | 1 |
| 2 | A Clock for Every Pass | `a-clock-for-every-pass` | a goroutine per contact from AOS to LOS in simulated time, calls to each station, a silent station cut off at LOS, the whole week stopped from one context | goroutines, `time`, `context.WithTimeout`, `context.AfterFunc`, `sync.WaitGroup` | 5 | 4 | 1 |
| 3 | On the Wire | `on-the-wire` | Season 4's frames on the pass line, telemetry checked strictly and kept newest-first, bad frames dropped, dead lines noticed and dialed again | `encoding/binary`, `hash/crc32`, `bufio.Reader`, `io.ReadFull`, deadlines | 5 | 4 | 1 |
| 4 | Signed and Sealed | `signed-and-sealed` | every command signed with derived keys, one at a time per satellite, Cirrine's sealed, never run twice, and a key replaced mid-week | `crypto/hkdf`, `crypto/hmac`, `crypto/aes`, `crypto/cipher`, `sync.Mutex` | 6 | 5 | 1 |
| 5 | Every File Delivered | `every-file-delivered` | downlinks reassembled and saved whole, signed webhooks, retries with backoff and jitter, missing chunks asked for again | `hash/crc32`, `os`, `net/http`, `crypto/hmac`, `math/rand/v2`, timers | 4 | 4 | 0 |
| 6 | The Storm | `the-storm` | JSON logs, the week's numbers, Season 6's pager through Ystmark's storm, and the week's report against golden files | `log/slog`, `sync/atomic`, `encoding/json`, `time.Ticker`, `testdata` | 5 | 4 | 1 |
| | **Total** | | | | **29** | **24** | **5** |

Difficulty rises from 8 to 10: Season 7 ended at 10, and a project restarts lower, as Build gctl
and Build the Board did. Every code lesson is `code_tests`. Each lesson's tests rehearse the whole
network for the whole week at most once, one cached run of the investors' week that every check needing it shares; every other case flies one to six
stations for a day or two.

## Chapter 1 · The Week's Plan

`the-weeks-plan` · 4 lessons · Monday 2 and Tuesday 3 July

Read strict settings and passplan's pass archive, apply passplan 2's rule station by station, save
each day's plan whole, and match bookingd's export against it.

#### 1.1 The Rehearsal's Settings
`the-rehearsals-settings` · code_tests · difficulty 8
- **Objective:** Read the rehearsal's settings strictly, with defaults, and say what it rehearses.
- **Scenario:** Monday 2 July: Tomás tells the Fishbowl about the investors' week, and Ingrid wants every run of the rehearsal to say exactly which week and which stations it flies.
- **Exercise:** `rehearse` returns nil at once; it must read `rehearsal.json` from `w.Dir` by the bible's rules and print `rehearsal: 2029-07-23, 7 days, 20 stations` on `w.Out`, or return the settings error and print nothing.
- **Check:** A full file and one with only `week` (7 days, all twenty stations); `days` 1 and 7 accepted, 0 and 8 refused; every error text exactly, a missing file wrapping `fs.ErrNotExist`; the mistakes: a decoder that ignores an unknown field or a second object, and a week parsed with a time of day.

#### 1.2 Every Contact of the Week
`every-contact-of-the-week` · code_tests · difficulty 8
- **Objective:** Turn the pass archive into the week's contacts with Season 7's rule, station by station, and save a file per day.
- **Scenario:** The dish loaders need each day's contacts, and Nkem wants them to be passplan 2's own, contact for contact.
- **Exercise:** `rehearse` must read `passes.txt`, keep the candidates (servable, at least 10°, at a rehearsed station), give them dishes with the given `assign`, write `plan/<day>.txt` for every day of the window, and print `plan: 14297 contacts, 28680 dropped`.
- **Check:** All seven files of the investors' week against the reference plan, line for line; a one-day, one-station rehearsal; a contact whose LOS is after midnight in its AOS's file; a bad archive line as `archive line 3: ...`; the mistakes: Pipit-1 kept at Ystmark (no shared band), passes under 10° kept, all stations sorted into one list before assigning, lines written in station order.

#### 1.3 Bookings Against the Plan
`bookings-against-the-plan` · code_tests · difficulty 9
- **Objective:** Match each booking to a contact by satellite, station and AOS to the second, and name every booked pass the plan can't fly, with its reason.
- **Scenario:** Tuesday 3 July: Yusra exports the week's 1,495 bookings from bookingd, and a booked pass without a dish is lost before the week begins.
- **Exercise:** `rehearse` must read `bookings.json` strictly, write `reports/unplanned.txt` with a line per unplanned booking and the first reason that fits, and print `bookings: 1493, unplanned 0`.
- **Check:** The clean export (no line); one planted booking per reason (`thistle`, `SQX`, Hollis's Cape Wren 04:58:30, Pipit-1 over Ystmark, Cirrine-1 at 8.2°, Furrow-A over Quindle Island dropped) in AOS order; bookings at other stations or outside the window left out; a booking one second off its pass is `no such pass`; a broken export named in the error; the mistake: matching by the minute or by station alone.

#### 1.4 Quick Check: The Week's Plan
`quick-check-the-weeks-plan` · quiz · difficulty 8
- **Objective:** Reason about candidates, dishes and matching.
- **Scenario:** Nkem's three questions before the plan meets the pass line.
- **Exercise:** Three questions.
- **Check:** A candidate with a free band and 40° can still be dropped (its station's dishes are busy); stations can be planned in any order because they never share a dish; a booking carries no pass ID, only its satellite, station and AOS, so the AOS must match to the second.

## Chapter 2 · A Clock for Every Pass

`a-clock-for-every-pass` · 5 lessons · Wednesday 4 and Thursday 5 July

Run every contact in its own goroutine on the simulated clock, call its station at AOS, hang up at
LOS, cut off a silent station, and stop the whole week from one context.

#### 2.1 From AOS to LOS
`from-aos-to-los` · code_tests · difficulty 8
- **Objective:** Run each contact concurrently from its AOS to its LOS, with waits measured from the injected clock.
- **Scenario:** Ingrid wants each pass on the Fishbowl's screen as it starts and ends, and the whole week in seconds, so she can run it again after every change.
- **Exercise:** After planning, `rehearse` must run the window: an `AOS` line at each contact's AOS and a `LOS ... flown` line at its LOS, then `week done: 14297 flown, 0 failed, 0 skipped` after the last LOS, and only then return.
- **Check:** Inside the bubble, every line's instant equal to its AOS or LOS for the whole investors' week; two dishes' contacts overlapping at Ystmark; `rehearse` returning at the last LOS, not before; the mistakes: `time.Until(c.AOS)` or `context.WithDeadline(ctx, c.LOS)`, 29 years away on the bubble's clock (the harness's guard reports `rehearse still running`), and passes flown one after another.

#### 2.2 Call the Station
`call-the-station` · code_tests · difficulty 8
- **Objective:** Dial a pass's station at AOS, read its banner, and hang up at LOS.
- **Scenario:** Sunniva wants the rehearsal's calls to reach each station's radio on time and leave on time, as the real ones will.
- **Exercise:** `rehearse` must call `w.Dial` for each contact at its AOS, read the banner (`PEL READY PROTO 1`, CR LF or LF), hold the call and close it at LOS; a wrong banner or a failed dial makes the pass `failed no_answer`.
- **Check:** The simulated world's records: one call per contact, dialed at AOS to the right station, hung up by the rehearsal at LOS; a banner sent one byte at a time; Ilha Brumosa's `BRM READY PROTO 2` and Corith Point's failing dials as `failed no_answer`, still printed at LOS; the mistakes: refusing a banner that ends in LF alone, accepting `PROTO 2`, and calling a station other than the pass's.

#### 2.3 Cut Off at LOS
`cut-off-at-los` · code_tests · difficulty 9
- **Objective:** End a read blocked on a silent connection by closing it when the pass's context ends.
- **Scenario:** On Wednesday's run the rehearsal hung for a simulated week behind Calidra Atoll, whose host's router answers a call and never says a word.
- **Exercise:** A pass waiting on a station that never sends its banner must end at its LOS as `failed no_answer`, and the week go on without it.
- **Check:** Calidra Atoll silent all week: each of its passes hung up at exactly its LOS, every other station's passes flown, `rehearse` back at the window's end; a station silent after its banner held to LOS and flown; the mistake: the previous lesson's read with no way out, which the harness's guard reports and then releases.

#### 2.4 Quick Check: Simulated Time
`quick-check-simulated-time` · quiz · difficulty 8
- **Objective:** Reason about instants, durations and cancellation in a bubble.
- **Scenario:** Nkem asks how a week can pass in a second without lying about time.
- **Exercise:** Three questions.
- **Check:** Instants come from `w.Now()` and waits are durations from it; `context.WithTimeout(ctx, c.LOS.Sub(w.Now()))` is the pass's deadline, not `WithDeadline(ctx, c.LOS)`; closing a connection is what ends a blocked read.

#### 2.5 Stop the Rehearsal
`stop-the-rehearsal` · code_tests · difficulty 9
- **Objective:** Stop the whole week from one context: no new call, every call ended at once, nothing left running, and the cause returned.
- **Scenario:** Ingrid stops a run halfway when she spots a wrong setting, and the next run must not meet a pass the last one left behind.
- **Exercise:** When `ctx` ends, `rehearse` must dial nothing more, hang up every open call with its `LOS ... stopped` line, print `rehearsal stopped at Tue 12:00:00` once everything has ended, and return `context.Cause(ctx)`.
- **Check:** A cancel at Tuesday 12:00:00 with a cause: every open call hung up at 12:00:00, no dial and no `AOS` line after it, the cause returned (`errors.Is` for a plain cancel); a cancel before the first AOS dials nothing; the mistakes: returning nil, and passes asleep until their AOS that wake and dial after the stop.

## Chapter 3 · On the Wire

`on-the-wire` · 5 lessons · Friday 6 to Monday 9 July

Speak Season 4's frames on the pass line, read telemetry strictly and keep each satellite's newest
frame, drop what fails its CRC, and notice a dead line and dial it again.

#### 3.1 Hello in Frames
`hello-in-frames` · code_tests · difficulty 9
- **Objective:** Write and read length-prefixed binary frames with a CRC, after the banner, through one buffered reader.
- **Scenario:** Sunniva's radios now want each call to name its pass and dish before they say anything more.
- **Exercise:** After the banner, `rehearse` must send HELLO (`thistle-2/PEL/02136 1`) and read the WELCOME; a NAK at HELLO makes the pass `failed refused`, or `skipped stowed` when its code is 2, BYE `stowed` makes it `failed stowed`, and no WELCOME by LOS `failed no_answer`; frames the station sends later are read and set aside until the next lesson.
- **Check:** HELLO's bytes against the bible's example; the banner and WELCOME arriving in one write, and WELCOME one byte at a time; a pass HELLO'd for the wrong dish refused with NAK 1; the mistakes: a CRC that covers the length, a length written little-endian, a new reader for the frames after the banner.

#### 3.2 Telemetry In
`telemetry-in` · code_tests · difficulty 9
- **Objective:** Read TLM frames until LOS, parse each line with the strict parser, and keep each satellite's newest frame.
- **Scenario:** The Board and the pager read every satellite's newest frame, and Joss says Thistle-2 still sends `batt_v=nan` when its sensor glitches.
- **Exercise:** `rehearse` must parse every TLM line with the given `parseLine`, drop what it refuses, keep the frame with the highest number per satellite, and write `reports/health.txt` at the end.
- **Check:** `health.txt` against the simulated world's newest frame per satellite for a day at Ystmark and Tirrin Rock and for the whole week; a satellite over both stations at once (newest by number, not by arrival); never-heard satellites; Thistle-2's `nan` lines never kept; the mistake: keeping the last frame received.

#### 3.3 Frames That Lie
`frames-that-lie` · code_tests · difficulty 9
- **Objective:** Drop a frame whose CRC fails and read on, and treat a length out of range as a lost line.
- **Scenario:** Varsk Island's host line flips a bit now and then, and on Friday one flipped length sent the rehearsal reading a thousand bytes of nothing.
- **Exercise:** A frame whose CRC doesn't match must be dropped and the pass go on; a length under 9 or over 1024 must end the line as lost, the pass `failed line_lost` at LOS.
- **Check:** A flipped byte in a Varsk Island pass's third TLM frame: the frames after it read, the pass flown; a length of 8 after a pass's second frame: `failed line_lost`; health never holding a corrupted line; the mistake: ending the pass at the first bad CRC.

#### 3.4 When the Line Goes Quiet
`when-the-line-goes-quiet` · code_tests · difficulty 10
- **Objective:** Notice a dead connection with an idle read deadline, and dial again with growing waits until the pass is welcomed again or LOS comes.
- **Scenario:** Ystmark's fiber dies without a sound in storms: calls stay open and nothing arrives, and the backhaul takes ten minutes to come up.
- **Exercise:** Every read must bring a byte within 90 s; a lost line must be closed and dialed again (at once, then Season 4's `redialDelay` between failures), say HELLO for the same pass, and read on until LOS.
- **Check:** Ystmark quiet from Thursday 03:31 with dials failing until 03:41: the loss noticed within 90 s, the waits between failed dials in their ranges, the passes after 03:41 flown over the backhaul, `pipit-2/YST/22861` `failed line_lost`; frames trickling 0.25 s apart never tripping the deadline; the mistakes: a deadline from `w.Now()` (never fires in the bubble), one deadline for a whole frame, redialing with no wait.

#### 3.5 Quick Check: On the Wire
`quick-check-on-the-wire` · quiz · difficulty 9
- **Objective:** Reason about framing, deadlines and lost lines.
- **Scenario:** Sunniva's questions, telegraphic, before the backhaul meets the rehearsal.
- **Exercise:** Three questions.
- **Check:** One `bufio.Reader` for the banner and the frames because it may already hold the first frame; an idle deadline moved before each read lets a slow frame live; a bad length means nothing after it can be trusted, a bad CRC only that frame.

## Chapter 4 · Signed and Sealed

`signed-and-sealed` · 6 lessons · Tuesday 10 and Wednesday 11 July

Sign every command with its satellite's derived key, number commands one at a time per satellite,
seal Cirrine's, resume after a lost line without running anything twice, and replace a key
mid-week.

#### 4.1 Every Command Signed
`every-command-signed` · code_tests · difficulty 9
- **Objective:** Derive a satellite's sign key from the uplink master and send a signed, numbered command, waiting for its answer while the station keeps talking.
- **Scenario:** Mirela: the investors' security reviewer will read every line the rehearsal sends up.
- **Exercise:** `rehearse` must read `uplink-keys.txt`, and on every contact, from AOS + 10 s, send `CMD <sat> <last+1> PING sig=<hex>` with the satellite's `sign` key and wait for its ACK; a NAK makes the pass `failed refused`, and no answer by LOS `failed no_ack`.
- **Check:** A day at Port Ellery: every PING run by its satellite; pipit-2's first line equal to the bible's vector; a satellite starting at `last` 0; a planted NAK; the mistakes: signing with the master itself, numbering from 1, and waiting for the ACK without reading frames, so the station's telemetry blocks it and the pass ends `no_ack`.

#### 4.2 One Satellite, Two Stations
`one-satellite-two-stations` · code_tests · difficulty 9
- **Objective:** Keep one satellite's commands in order across concurrent passes with one lock per satellite, a channel with room for one.
- **Scenario:** Nkem: "What does the satellite see if two stations send it a command at once?" Over the poles, Ystmark and Tirrin Rock often hold the same Furrow.
- **Exercise:** A satellite must never have two commands in flight, whichever stations hold it, and every command must take the next number.
- **Check:** A day at Ystmark and Tirrin Rock: every satellite's commands run in strictly increasing numbers, none refused as replayed, none run twice; the mistake: numbering per pass, so two passes send `last+1` and the satellite refuses the second.

#### 4.3 Sealed for Cirrine
`sealed-for-cirrine` · code_tests · difficulty 9
- **Objective:** Seal a command with AES-256-GCM under a derived key, a fresh nonce and the header as associated data.
- **Scenario:** Cirrine's contract: nobody but Cirrine and Apogee reads its commands, now all thirteen satellites' of them.
- **Exercise:** Every command to a Cirrine Atmospherics satellite must go up as Season 5's SEAL line; every other satellite's stays a CMD line.
- **Check:** Every Cirrine opening every line by Season 5's checks; no nonce used twice in the week; the bible's `k-2028-10` vector with the fixed nonce; the mistakes: a CMD line to a Cirrine (`command must be sealed`), a nonce from `math/rand`, associated data in another order (`command failed authentication`).

#### 4.4 Never Twice
`never-twice` · code_tests · difficulty 10
- **Objective:** After a lost line, settle a command in flight from the satellite's `last`: done if it ran, sent again with the same number if not.
- **Scenario:** On Monday's run a PING to Cirrine-7 went up twice under two numbers when Ystmark's line dropped between the command and its answer.
- **Exercise:** With a command in flight when the line is lost, `rehearse` must read the new WELCOME's `last`, count the command as answered if it ran, and otherwise send the same line again.
- **Check:** A call dropped just before the satellite gets a PING: sent again, same number, run once; dropped just after it ran: not sent again, the pass flown; the same across two contacts for an unanswered LOAD_KEY; the mistakes: a new number for the lost command (it runs twice) and no resend (`no_ack`).

#### 4.5 A Key Replaced Mid-Week
`a-key-replaced-mid-week` · code_tests · difficulty 10
- **Objective:** Load a new master's keys before its switch and change signing and sealing at the instant it activates, for satellites that hold it.
- **Scenario:** Mirela made `k-2029-07` on Monday 9 July for 00:00 on Wednesday 25 July, the investors' reviewer watching.
- **Exercise:** On each satellite's first contact, after its PING, `rehearse` must send LOAD_KEY for every master it doesn't hold, sealed under the newest one it does, and from each master's Activate sign and seal with it once the satellite holds it.
- **Check:** The investors' week: all eighty loaded on their first contacts, by Monday 09:27:17; every command until Tuesday 23:59:59 under `k-2028-10` and from Wednesday 00:00:00 under `k-2029-07`; the retired `k-2028-07` never used; a rehearsal starting on Wednesday (LOAD_KEY first under `k-2028-10`, then the new key); the mistakes: switching before the satellite holds the key (`unknown key "k-2029-07"`) and LOAD_KEY sent as a CMD line.

#### 4.6 Quick Check: Keys and Numbers
`quick-check-keys-and-numbers` · quiz · difficulty 9
- **Objective:** Reason about derived keys, numbering and rotation.
- **Scenario:** Mirela's questions, ending in "LGTM once it fails closed".
- **Exercise:** Three questions.
- **Check:** A stolen satellite key commands only that satellite; the same number again after a lost line is what keeps a command from running twice; the old master keeps working for seven days after the switch so a satellite that missed its LOAD_KEY can still get it.

## Chapter 5 · Every File Delivered

`every-file-delivered` · 4 lessons · Thursday 12 and Friday 13 July

Reassemble downlinks from numbered chunks and save them whole, tell each customer by Season 3's
signed webhook, retry on Season 3's schedule in simulated time, and ask the satellite again for
what never came.

#### 5.1 Chunks Into Files
`chunks-into-files` · code_tests · difficulty 9
- **Objective:** Keep chunks by number, check the whole file's CRC-32 when its end arrives, and save it whole.
- **Scenario:** Hollis: "EVERY CHUNK OR NOTHING." Booked passes bring their files down in numbered chunks, some out of order, some twice.
- **Exercise:** `rehearse` must reassemble each FILE's chunks, write `files/<name>` only when every chunk is there and the CRC matches, print `LANDED <name>`, and list what is still incomplete in `reports/incomplete.txt` at the end.
- **Check:** Every booked file of a day at three stations byte for byte, then the whole week's count; chunks out of order and twice; a missing chunk (`missing 21 22 23`) and a wrong END CRC (`crc mismatch`), neither written; a redial in the middle of a file; the mistakes: writing chunks as they come (a partial file) and a CRC over the chunks in arrival order.

#### 5.2 Your Data Has Landed
`your-data-has-landed` · code_tests · difficulty 9
- **Objective:** Send Season 3's signed webhook for every landed file, to the booking's customer.
- **Scenario:** Yusra: the investors will ask Furrowlight how soon its images arrive; and Joss has written Kessick's first receiver, for Thistle-2's first science on Thursday 26 July.
- **Exercise:** `rehearse` must read `customers.json` and, as each booked file lands, POST its `data.landed` event with `Apogee-Timestamp` and `Apogee-Signature`, printing `DELIVERED <file>` on a 2xx.
- **Check:** Every receiver verifying every signature and timestamp; one event per landed file, within 60 s of landing; event IDs from `evt-000001`, unique, with no gaps; Kessick's event for `bk-8792`; the mistakes: signing a re-encoded body or one with a newline, a timestamp in milliseconds.

#### 5.3 Retry Without Hammering
`retry-without-hammering` · code_tests · difficulty 10
- **Objective:** Retry failed deliveries on Season 3's backoff with jitter, in simulated time, signing every attempt again.
- **Scenario:** Cirrine's receiver restarts for its own deploy on Tuesday 24 July at 09:32, eighteen seconds before Cirrine-5's file lands.
- **Exercise:** `rehearse` must retry transport errors, 408, 429 and 5xx up to 8 attempts, waiting as Season 3 says, the same event ID with a fresh timestamp and signature each time; any other answer, or the 8th failure, makes it `UNDELIVERED`, listed in `reports/undelivered.txt`; the week ends only when every delivery has.
- **Check:** Cirrine-5's event delivered once its receiver is back, every wait inside its range; a receiver down all week: 8 attempts, then undelivered; a 400 never retried, a 503 retried; the mistakes: reusing the first attempt's timestamp (refused after 300 s), retrying a 4xx, waits with no jitter (each at the top of its range, outside it).

#### 5.4 Missing Chunks, Asked Again
`missing-chunks-asked-again` · code_tests · difficulty 10
- **Objective:** Turn a file's missing chunks into a command on the satellite's next contact, and land the file when they arrive.
- **Scenario:** On Tuesday 24 July Furrow-K's 07:40 image over Vorrel Island comes down three chunks short, and Hollis wants it whole.
- **Exercise:** On a satellite's next contact after a file stays incomplete, `rehearse` must send `DUMP_FILE <name> <n>...` (signed or sealed like its other commands), finish the file from the chunks that come, and deliver it.
- **Check:** `furrow-k-20290724-0740.tif` asked for on `furrow-k/VSK/00404`, landing at 08:12:30 and delivered once; a CRC mismatch asked for whole; a satellite with no later contact left in `incomplete.txt`; the mistakes: asking again before the first pass's END, and asking on every later contact after the file is complete.

## Chapter 6 · The Storm

`the-storm` · 5 lessons · Monday 16 to Friday 20 July

Write structured logs on the simulated clock, count the week, page through Season 6's rules as
Ystmark's storm arrives, and write the week's report, checked against golden files.

#### 6.1 A Line for Every Pass
`a-line-for-every-pass` · code_tests · difficulty 9
- **Objective:** Write JSON log records with `log/slog`, stamped with the simulated clock, with nothing secret in them.
- **Scenario:** Monday 16 July: Sunniva forecasts a storm at Ystmark for Thursday 26 July, and Ingrid wants to read the rehearsal's version of that day the way the night shift would.
- **Exercise:** `rehearse` must write the bible's records to `logs/uplinkd.log`, `logs/dlinkd.log` and `logs/pager.log`, every record's `time` from `w.Now()`.
- **Check:** Exactly one `pass started` and one end record per contact, with the bible's attributes; a failed and a skipped pass's levels, `reason` and `err`; one `file landed` per file and one `event delivered` per event; every master's text and hex, every derived key and every webhook secret absent from every log; the mistakes: the wall clock in `time`, a duration as nanoseconds.

#### 6.2 Counted, Not Guessed
`counted-not-guessed` · code_tests · difficulty 9
- **Objective:** Count the week with counters safe for many goroutines and publish them in Season 6's shapes.
- **Scenario:** Yusra: the investors will ask how many passes flew, and "about" is not a number.
- **Exercise:** `rehearse` must keep the bible's counters and histogram as the week runs and write `reports/metrics.json` at the end.
- **Check:** Every counter equal to the simulated world's own count for a day at three stations and for the whole week; a 2 s answer in the ≤ 2 bucket and Ystmark's 2.7 s ones in ≤ 5; skipped passes under `skipped`, never `failed`; the mistake: a redialed pass counted twice.

#### 6.3 One Page per Problem
`one-page-per-problem` · code_tests · difficulty 10
- **Objective:** Evaluate stateful alert rules every minute on the simulated clock and route each alert by the rota.
- **Scenario:** Ingrid: the storm must wake whoever is on call once, not once a pass, and you carry the pager for the investors' week.
- **Exercise:** `rehearse` must evaluate `pass-failures` and `failed-pass` for every rehearsed station on every minute, route them with the given rota, and write `reports/alerts.txt` and the `alert sent` records.
- **Check:** Ystmark's Thursday: exactly the bible's five alerts, the page at 10:56 to `No. 039` and its RESOLVED at 16:56; skipped passes never counted; a failure at Monday 07:59 paging Nkem and one at 08:00 paging No. 039; same-minute order; the mistakes: a page per failure, a RESOLVED for a Warning, a failure counted from its last frame instead of its LOS.

#### 6.4 Quick Check: Pages and Counts
`quick-check-pages-and-counts` · quiz · difficulty 9
- **Objective:** Reason about what the pager and the report count.
- **Scenario:** Tomás asks what the investors will see on Friday.
- **Exercise:** Three questions.
- **Check:** A stowed pass is skipped, not failed, so it never pages, yet it is lost to its booking; Ystmark's week allows one lost pass, from 210 booked at 99.5% in integers; counts come from counters, not from reading the logs back.

#### 6.5 Every Pass Counts
`every-pass-counts` · code_tests · difficulty 10
- **Objective:** Finish the rehearsal and check the whole investors' week against golden files.
- **Scenario:** Friday 20 July, 15:00, the Fishbowl: the full rehearsal with Ingrid's script, the fiber cut and the wind lock at Ystmark, Cirrine's restart, Furrow-K's lost chunks and Thistle-2's first science.
- **Exercise:** `rehearse` must write `reports/week.txt` (rows, `all`, the lost bookings, the verdict) as the bible gives it, the last thing the week's golden files need.
- **Check:** The investors' week against `testdata/week.golden` (NO-GO, the nine Furrows), `alerts.golden` and `health.golden`, `metrics.json` against the simulated world's counts, both lists empty; Yusra's moves on Thursday at five stations: GO; a station with nothing booked showing `-`; the mistakes: lost counting failures only, a budget computed in floats, unplanned bookings left out.
