# Apogee story bible: the Board project

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
