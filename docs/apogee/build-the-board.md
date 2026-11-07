# Build the Board: syllabus

The guided project after Season 3 of the Apogee series, story time 7–11 February 2028. The
world, the characters, the code and the writing rules are in `story.md`, whose section "The
Board project" fixes the program's behavior, and Season 3's section fixes the types it reuses.

**Authoring status:** all four chapters are implemented as 17 lessons. One writer designed the
whole program and wrote chapter 1, which was read in the app in light and dark before chapters
2–4 were written; an independent reviewer then checked every lesson. The lessons in
`courses/apogee.sqlite` are the source of truth, and `go run ./cmd/author export` writes them
out to read; the larger changes from this plan are listed under "Changes in writing and
review" at the end.

- **Title:** Build the Board
- **Slug:** `build-the-board` · **Kind:** Guided Project
- **Short description:** Build the Board as one Go program, from a page made from a template
  to a live stream on forty screens.
- **Description:** The Board, the big screen in the Fishbowl that lists the next passes, is a
  laptop plugged into a television, and the night shift wants it on every screen in the
  building. Over one week at Port Ellery you build it as one program, a lesson at a time: a
  page made from a template that shows whatever customers write as text, each station's wind
  and who is on call, the same data as JSON that answers "nothing new" when nothing is, a
  stream that pushes every change, one source fanned out to forty screens without waiting for
  a slow one, weather fetched once a minute, and streams that say goodbye when the Board
  restarts. On Thursday night the whole night is checked against snapshots, and the Board is on
  every screen. Standard library only.
- **Story:** a short opening, each chapter's briefing and debrief, and the field guide, told
  as `docs/storytelling.md` says.

## How a guided project differs

The learner builds one program from an empty page to a finished server. Each lesson adds one
feature, its starter is the previous lesson's solution, and its read-only tests check the
program as a whole, so every earlier feature stays checked. Readmes say what the feature must
do and show less code than a course's lessons do; hints carry more of the help. Everything
else follows the bible's style guide.

Every test drives the program through HTTP with `httptest`: a page is checked by what the
server sends, never in a browser. A learner may open the Board in their own browser at
127.0.0.1, but no check depends on it.

## At a glance

| # | Chapter | Slug | What gets built | Go | Lessons | tests | quiz |
|---|---|---|---|---|---|---|---|
| 1 | One Page | `one-page` | the Board as an HTML page from a template: the next passes sorted by AOS, customers' notes escaped, a weather panel per station and who is on call | `html/template`, `net/http`, escaping, `embed` | 5 | 4 | 1 |
| 2 | Always Up to Date | `always-up-to-date` | the view as JSON with versions and `304 Not Modified`, then a stream of server-sent events that pushes every change | `encoding/json`, ETags, `http.Flusher` through `http.ResponseController`, `context`, channels | 4 | 4 | 0 |
| 3 | Forty Screens | `forty-screens` | one source of updates fanned out to every viewer, slow viewers dropped, not waited for, and forty screens at once | goroutines, `select`, buffered channels, `sync.Mutex` | 4 | 3 | 1 |
| 4 | Night Shift Ready | `night-shift-ready` | weather cached for a minute, streams that end on shutdown, the whole program run from one function, and the reception screen, checked against a night of golden files | `sync`, `time`, `http.Server.Shutdown`, `httptest`, `testdata` | 4 | 4 | 0 |
| | **Total** | | | | **17** | **15** | **2** |

Difficulty rises from 6 to 9. Every lesson but the two quizzes is `code_tests`: the test builds
the Board with `NewBoard`, fake sources and a fake clock, serves it with `httptest`, and checks
what comes back over HTTP; nothing listens on a fixed port, opens a browser or touches the
real weather service.

## Chapter 1 · One Page

`one-page` · 5 lessons · Monday 7 February

The Board as an HTML page from a template: data prepared in Go, notes escaped, and a page built
from named parts.

#### 1.1 A Page From a Template
`a-page-from-a-template` · code_tests · difficulty 6
- **Objective:** Render an HTML page from an embedded template with `html/template`, into a buffer first, and serve it as `text/html`.
- **Scenario:** The Board is a laptop plugged into the Fishbowl's television, and when it rebooted during last week's night shift the Fishbowl had no Board for twenty minutes.
- **Exercise:** `(*Board).servePage` writes `TODO`; you execute `b.page` with a `pageData` of the current view and time into a `bytes.Buffer`, then send it with `Content-Type: text/html; charset=utf-8`, answering 500 and sending none of the page when the template fails.
- **Check:** `GET /` → 200, the content type, the title and one `<tr class="pass">` row per pass in the view; a template that fails on this view (the test swaps one in) → 500 and none of the page's text. The mistake: executing straight into `w`, which sends half a page with status 200.

#### 1.2 Next Passes First
`next-passes-first` · code_tests · difficulty 6
- **Objective:** Prepare a page's data in Go before the template sees it: a sorted copy, finished passes dropped, at most twelve.
- **Scenario:** The television fits twelve rows, and the night shift reads the top one first.
- **Exercise:** `nextPasses(all []Booking, now time.Time) []Booking` returns `all`; you return a new slice of the passes whose LOS is after `now`, sorted by AOS then ID, at most 12. `Refresh` already calls it.
- **Check:** Through `GET /`: rows in AOS order, ties by ID; a pass whose LOS equals now is gone and a pass in progress stays; 15 passes → 12 rows; the source's slice keeps its order, the mistake being a sort in place.

#### 1.3 A Note Is Only Text
`a-note-is-only-text` · code_tests · difficulty 7
- **Objective:** Let `html/template` escape everything a customer wrote, by never marking it as trusted HTML.
- **Scenario:** Hollis's note on tonight's booking reads `<b>RUSH</b> & PRIORITY`, and the kitchen screen showed it in bold.
- **Exercise:** `board.html` passes each note through a template function, `asHTML`, that converts it to `template.HTML`; you take it out of the template and the `FuncMap`, so the note is printed as text.
- **Check:** The note appears as `&lt;b&gt;RUSH&lt;/b&gt; &amp; PRIORITY`; a note holding `<script>` never reaches the page as a tag; quotes print as `&#34;` and `&#39;`; ordinary notes are unchanged.

#### 1.4 Quick Check: Templates
`quick-check-templates` · quiz · difficulty 6
- **Objective:** Reason about contextual escaping and rendering safely.
- **Scenario:** Nkem's questions before the page goes on a second screen.
- **Exercise:** Three questions.
- **Check:** `html/template` escapes by context (text, attribute, script); `template.HTML` turns escaping off for that one value; rendering into a buffer lets an error still become a 500.

#### 1.5 Wind on Every Panel
`wind-on-every-panel` · code_tests · difficulty 7
- **Objective:** Build a page from named parts with `define` and `template`, looping with `range` and choosing with `if` and `else`.
- **Scenario:** The night shift wants each station's wind beside the passes, and the name of whoever is on call tonight.
- **Exercise:** `board.html` shows the passes only; you define a `panel` template, use it for each station in the view, and add the on-call line, all as the bible gives them.
- **Check:** Panels in the station table's order, with class `track` at 22 km/h, `warn` at exactly 50, `stow` at exactly 72 and `none` without a reading; `YST wind 41 km/h, gusts 55`; the on-call name from the rota, last week's at Monday 07:59 and this week's at 08:00.

## Chapter 2 · Always Up to Date

`always-up-to-date` · 4 lessons · Tuesday 8 February

The same view as JSON, versions that let a screen skip what it already has, and a stream that
pushes every change.

#### 2.1 The Board as Data
`the-board-as-data` · code_tests · difficulty 7
- **Objective:** Serve the same view as JSON at a second address, from the one place that builds it.
- **Scenario:** Reloading the whole page every fifteen seconds makes the kitchen screen flicker, and Ingrid wants the page to fetch only the data.
- **Exercise:** `GET /board.json` answers 404; you register it and write the current view's JSON with `Content-Type: application/json` and `Cache-Control: no-cache`. The page's script that fetches it is given.
- **Check:** `passes`, `weather` and `on_call` match the page's rows and panels; a panel without a reading has no `weather` key and `call` `none`; no passes → `[]`, not `null`; the content type and the cache header.

#### 2.2 Nothing New, Nothing Sent
`nothing-new-nothing-sent` · code_tests · difficulty 7
- **Objective:** Name each version of the data with an `ETag`, and answer `304 Not Modified` when the caller already has it.
- **Scenario:** Forty screens asking every fifteen seconds is 160 copies a minute of a board that changes a few times an hour.
- **Exercise:** `Refresh` bumps `b.version` every time it runs; you bump it only when the view's JSON changes, and `serveJSON` sends `ETag: "v<version>"` and answers 304 with no body when `If-None-Match` matches it.
- **Check:** The first GET → 200 with `"v1"`; the same ETag → 304 with an empty body; a refresh with nothing new → still `"v1"`, the mistake being a bump on every refresh; a changed pass list → 200 with `"v2"`; an old ETag → 200 with the new data.

#### 2.3 A Stream That Stays Open
`a-stream-that-stays-open` · code_tests · difficulty 8
- **Objective:** Hold a response open and push a server-sent event, flushing it through the writer's `http.Flusher` with `http.ResponseController`.
- **Scenario:** Even fifteen seconds is too late when Cape Wren stows during a pass.
- **Exercise:** `GET /events` answers 501; you set `Content-Type: text/event-stream`, write the current view as one `board` event with its `id`, flush it, and wait until the viewer leaves.
- **Check:** A loopback client reads the first event within a second: the headers, `event: board`, `id: 1` and the view's JSON on one `data:` line; when the client cancels, `/health` shows 0 viewers again within a second. The mistake: no flush, so the event waits in a buffer and the client's guard fails.

#### 2.4 Every Change as It Happens
`every-change-as-it-happens` · code_tests · difficulty 8
- **Objective:** Feed a stream from a channel with `select`, ending on the request's context.
- **Scenario:** The night shift wants a stowed dish on the screen within a second of the refresh that saw it.
- **Exercise:** `Refresh` sends each new version on `b.updates`; `serveEvents` writes only the first event; you loop with `select` over `b.updates` and `r.Context().Done()`, writing each version as an event.
- **Check:** One viewer reads `id: 1`, then `id: 2` and `id: 3` after two refreshes that changed something; a refresh with nothing new sends nothing; the handler ends when the client leaves. (A second viewer would take turns with the first on the one channel: the next chapter's problem.)

## Chapter 3 · Forty Screens

`forty-screens` · 4 lessons · Wednesday 9 February

One source of updates fanned out to every viewer, without ever waiting for the slowest.

#### 3.1 Every Screen Gets Every Change
`every-screen-gets-every-change` · code_tests · difficulty 8
- **Objective:** Fan one stream of values out to many receivers, each with its own buffered channel, behind a mutex.
- **Scenario:** The meeting room got the 03:58 change and the kitchen didn't, because two streams were taking turns on one channel.
- **Exercise:** `Hub`'s methods are stubs, and the stream, `Refresh` and `/health` already use the hub instead of `b.updates`; you write `Subscribe` (a new channel with room for 4, kept in a map under the hub's mutex), `Unsubscribe`, `Publish` (a send on every channel) and `Len`.
- **Check:** Three loopback viewers each read versions 2 and 3; a viewer that leaves is unsubscribed and `/health`'s count drops; `Publish` with nobody subscribed returns at once.

#### 3.2 Never Wait for Slow Screens
`never-wait-for-slow-screens` · code_tests · difficulty 8
- **Objective:** Send without ever blocking, with `select` and `default`, and drop a receiver that has fallen behind.
- **Scenario:** The reception laptop sleeps with its stream open, and at 04:12 `Refresh` hung behind it while the Fishbowl waited for Pipit-1.
- **Exercise:** `Publish` sends with a plain send; you send with `select` and `default`, and a channel with no room is removed and closed, which ends its stream with `bye`.
- **Check:** A stream whose writer blocks (a test `ResponseWriter` that lets the first event through, then holds every `Write` until released) is dropped once its 4 slots fill, and every `Refresh` returns at once; a loopback viewer beside it reads every version; the blocked stream's handler returns once released; `Unsubscribe` after a drop doesn't panic, the mistake being a channel closed twice.

#### 3.3 Forty at Once
`forty-at-once` · code_tests · difficulty 9
- **Objective:** Subscribe before reading the shared state, so nothing published in between is missed when many goroutines join at once.
- **Scenario:** Every morning the cleaners switch on all forty screens in the same minute as the 08:00 refresh, and three of them showed the 07:59 board until the next change.
- **Exercise:** `serveEvents` writes the current view and then subscribes; you subscribe first, and skip any version no newer than the one the stream has already written.
- **Check:** A stream whose very first write the test holds while a refresh publishes version 2 still writes version 2 once released, where the starter never does; forty loopback viewers connecting between refreshes all end on the latest version and none gets a version twice; `/health` says 40, and 1 after 39 leave.

#### 3.4 Quick Check: Fan-Out
`quick-check-fan-out` · quiz · difficulty 8
- **Objective:** Reason about channels shared by many goroutines.
- **Scenario:** Nkem's questions after the cleaners' test.
- **Exercise:** Three questions.
- **Check:** A value sent on a channel reaches one receiver only; `select` with `default` never blocks; only the hub, which owns a channel, closes it, because closing a channel twice panics.

## Chapter 4 · Night Shift Ready

`night-shift-ready` · 4 lessons · Thursday 10 February and the night into Friday 11 February

Weather fetched once a minute, streams that end on shutdown, the whole program run from one
function, and the reception screen, checked against a whole night.

#### 4.1 Weather Once a Minute
`weather-once-a-minute` · code_tests · difficulty 8
- **Objective:** Cache each key's result for a time to live on an injected clock, behind a mutex.
- **Scenario:** Tomás asks why his old weather service now gets twenty calls a minute from the Board alone.
- **Exercise:** `(*weatherCache).reading(ctx, code)` calls the weather service every time; you keep each station's reading with the time it came, give it back for 60 seconds, and fetch again after, keeping no errors.
- **Check:** Four refreshes 15 seconds apart → one call per station; at exactly 60 seconds → fetched again; an error → `none` on that panel and a new try at the next refresh; each station is cached on its own.

#### 4.2 Streams Say Goodbye
`streams-say-goodbye` · code_tests · difficulty 9
- **Objective:** End long-lived streams when the server shuts down, with `RegisterOnShutdown`, because `Shutdown` waits only for connections that go idle.
- **Scenario:** Tonight's deploy at 05:00 must not wait for forty open streams.
- **Exercise:** `Hub.Close` is empty, and `newServer(b *Board) *http.Server` sets only the handler; you close every channel once and refuse new subscribers, and in `newServer` register `b.hub.Close` with `RegisterOnShutdown`.
- **Check:** With three loopback viewers of a `newServer` serving a 127.0.0.1 listener, `Shutdown` with a 2 s limit returns nil in well under a second; each viewer reads `event: bye` and then the end of the stream; `Publish`, `Subscribe` and `Close` after `Close` are harmless. The mistake: no `RegisterOnShutdown`, so `Shutdown` runs out its 2 seconds and fails.

#### 4.3 The Whole Board, Running
`the-whole-board-running` · code_tests · difficulty 9
- **Objective:** Run the whole program from one `run` function: refresh on a ticker, serve, and stop cleanly when the context ends.
- **Scenario:** Thursday evening: Nkem wants the Board started exactly as the night shift will start it.
- **Exercise:** `run(ctx context.Context, ln net.Listener, b *Board, ticks <-chan time.Time) error` only serves a `newServer`; you refresh once at start and on every tick, and when `ctx` ends you shut down within 2 seconds and return nil.
- **Check:** Through HTTP: page, JSON and stream answer; a tick after a change in the sources reaches a viewer; a refresh that fails keeps the last view and the run going; cancelling `ctx` ends `run` within a second, and the viewer reads `bye`.

#### 4.4 Every Screen in the Building
`every-screen-in-the-building` · code_tests · difficulty 9
- **Objective:** Finish the Board and check it against a whole night recorded in golden files.
- **Scenario:** The night into Friday 11 February: forty screens, Cape Wren's wind rising, a deploy at 05:00, and Hollis at reception at 09:00.
- **Exercise:** Add the reception screen: `?screen=lobby` on `/`, `/board.json` and `/events` leaves out notes, booking IDs and the on-call name, the last thing the night's snapshots need; then every snapshot in `testdata` matches.
- **Check:** Each `testdata/*.golden` (the page, the lobby page, both JSON views and the stream's events) at its moment of the night, with the clock and the sources moved as `testdata/night.txt` says; the lobby's ETag is `"v<n>-lobby"`; no note or booking ID appears anywhere in a lobby answer.

## Changes in writing and review

- **The night's passes** end at times their orbits allow. Ten of the sixteen ran longer than
  any pass their satellite can make over that station; every AOS is kept and the LOS moved
  (Pipit-1's 04:12 pass over Port Ellery now ends at 04:21:27, so the page shows `04:12` to
  `04:21`). Each test moment still drops exactly one pass. The passes are each possible, but
  a few satellites' passes don't all come from one orbit; no lesson or check depends on that.
- **3.4** A channel delivers each value sent on it to exactly one receiver; the quiz no longer
  says a channel never copies a value.
- **Field guide** Dishes stow when the wind reaches 72 km/h, as in Build gctl.
