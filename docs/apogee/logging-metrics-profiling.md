# Logging, Metrics and Profiling: syllabus

**Authoring status:** all eleven chapters are implemented as 88 lessons. Chapter 1 was read in
the app in light and dark before the rest was written; chapters 2–11 were written from this plan
two chapters at a time and every chapter independently reviewed. Chapters 2–11 were then read in the app in light and dark, two chapters at a time. Every pass comes from the orbit model the earlier seasons use, and chapter 9
plans on the settled one Season 7 plans with. The lessons in `courses/apogee.sqlite` are the
source of truth, and the larger changes from this plan are listed under "Changes in writing and
review" at the end.

Season 6 of Apogee, November 2028 to February 2029. The learner joins the on-call rota, and on
their second night the pager wakes them at 02:46: one pass in five at Lantern Atoll is failing,
and the log says only `error`. The learner moves from a line that says nothing to programs that
explain themselves, and in February finds the slow leak behind it. The world, the characters
and every shared format are in `story.md`, section "Season 6: on call", which fixes every log
line, metric, rule, number, date, type and signature used here.

- **Title:** Logging, Metrics and Profiling
- **Slug:** `logging-metrics-profiling` · **Kind:** Course
- **Short description:** Make Apogee's programs explain themselves at three in the morning, with logs, metrics, alerts, profiles and traces.
- **Description:** You join Apogee's on-call rota, and on your second night the pager wakes you at 02:46: one pass in five at Lantern Atoll is failing, the caretaker sees nothing wrong, and the log says only "error". With Nkem reviewing, you make the programs explain themselves: structured, leveled logs that carry their pass through every program and keep every secret out; log files that never fill a disk; counters, gauges and histograms you can trust; alert rules that wake the right person once; error budgets counted in whole passes; CPU and heap profiles, benchmarks and allocation counts; and traces that follow one request from headquarters to an island. By February the night shift calls you, and you find what has been failing Lantern Atoll's passes since October. Standard library only, every check on your own Mac.
- **Story:** an opening, each chapter's briefing and debrief, and the field guide, told as `docs/storytelling.md` says, following the bible's spine for the season.

## At a glance

| # | Chapter | Slug | Go topics | Lessons |
|---|---|---|---|---|
| 1 | Logs a Machine Can Read | `logs-a-machine-can-read` | Write structured records with `log/slog`: messages and attributes, the text and JSON handlers, typed values, errors, durations in milliseconds, and an old `*log.Logger` routed through a handler. | 8 |
| 2 | Levels and Groups | `levels-and-groups` | Choose levels, filter with `HandlerOptions` and change them live with `LevelVar`, attach attributes once with `Logger.With`, group them with `slog.Group`, and skip work nobody reads. | 9 |
| 3 | Logs With Context | `logs-with-context` | Carry pass and request IDs in `context.Context`, write a `slog.Handler` that adds them to every record and survives `With`, pass them between programs in a header, and read the lines back in order. | 8 |
| 4 | Never in a Log | `never-in-a-log` | Keep secrets and personal data out of logs: redact headers, give secret types a `LogValue`, replace attributes by key with `ReplaceAttr`, mask addresses and scrub secrets hidden in text. | 8 |
| 5 | Logs That Fit the Disk | `logs-that-fit-the-disk` | Write log files that stay inside a budget: append safely, count bytes, rotate by size and by UTC day without splitting a line, prune old files, and sample repeated records. | 8 |
| 6 | Counting What Matters | `counting-what-matters` | Count with `sync/atomic`: counters, gauges, labeled counters and histograms by hand, published with `expvar`, and exact totals from readings across restarts. | 8 |
| 7 | When to Wake Someone | `when-to-wake-someone` | Turn numbers into pages: rules over a window of results, thresholds held for a time, hysteresis, one alert per problem, routing to the person on call, and silences that end. | 8 |
| 8 | Error Budgets | `error-budgets` | Turn a target into numbers in code: availability, a budget in whole passes with integer arithmetic, a rolling window, burn rates, and how long a budget lasts. | 7 |
| 9 | Where the Time Goes | `where-the-time-goes` | Write CPU profiles with `runtime/pprof`, read the hot path, measure inside a test with `testing.Benchmark` and `b.Loop`, and make a planner many times faster without changing its answer. | 8 |
| 10 | Where the Memory Goes | `where-the-memory-goes` | Write heap profiles, read memory after a collection, count allocations with `testing.AllocsPerRun`, free big arrays behind small slices, count goroutines and make every goroutine end. | 8 |
| 11 | Following One Request | `following-one-request` | Follow work through a program and between programs: `runtime/trace` tasks and regions, a flight recorder for the last seconds, spans by hand carried in a context and a header, and a timeline built from them. | 8 |

88 lessons: 4 code, 70 code_tests and 14 quizzes (84% hands-on). Calibration starts with
chapter 1; later chapters are written only after it has been imported and read in the app.
Difficulty rises from 6 to 10. Speed is measured inside tests with `testing.Benchmark` and
`testing.AllocsPerRun`, and checks compare, never time; profiles and traces are written to the
lesson folder and checked only by their first bytes, with `go tool pprof` and `go tool trace`
offered in readmes, never required. Clocks and random sources are injected, files live in
`t.TempDir()`, and no test sleeps more than 50 ms, reaches the network or touches the learner's
files.

## Chapter 1 · Logs a Machine Can Read

`logs-a-machine-can-read` · 8 lessons

Write structured records with `log/slog`: messages and attributes, the text and JSON handlers, typed values, errors, durations in milliseconds, and an old `*log.Logger` routed through a handler.

#### 1.1 A Line That Says Why
`a-line-that-says-why` · code · difficulty 6
- **Objective:** Write a record with a message and a key-value pair through a `*slog.Logger`.
- **Scenario:** Wednesday 8 November, 02:48: Lantern Atoll's uplinkd log holds three lines, and all three end in `error`.
- **Exercise:** The starter makes a text logger on standard output with the given `dropTime` (no time, so the line fits the screen) and still calls `log.Print("error")`; you replace that call with one `logger.Error` call: message `pass failed`, attribute `site` `LAN`.
- **Check:** Exactly `level=ERROR msg="pass failed" site=LAN`; the starter prints nothing on standard output, because `log` writes to standard error.

#### 1.2 Lines for Machines
`lines-for-machines` · code_tests · difficulty 6
- **Objective:** Write records as JSON, one object per line, with `slog.NewJSONHandler`.
- **Scenario:** The night shift's search reads JSON lines, and Lantern Atoll's lines must go straight into it.
- **Exercise:** `newLogger(w io.Writer) *slog.Logger` returns a text logger; you return one with a JSON handler on `w`.
- **Check:** One record is one line that decodes as an object with `level` `ERROR`, its `msg` and `"site":"LAN"`; two records make two lines; the starter's text line fails to decode, and the message quotes it.

#### 1.3 Quick Check: Log Records
`quick-check-log-records` · quiz · difficulty 6
- **Objective:** Check what a record holds and how each handler writes it.
- **Scenario:** Nkem's three questions before the new lines reach the station.
- **Exercise:** Three questions about a four-line snippet in the readme.
- **Check:** The message is fixed text and the values are attributes; an argument without a partner becomes `!BADKEY`; the text and JSON handlers write the same record in two layouts.

#### 1.4 Numbers Stay Numbers
`numbers-stay-numbers` · code_tests · difficulty 6
- **Objective:** Write typed attributes (`slog.String`, `slog.Int`, `slog.Float64`) with `Logger.LogAttrs`.
- **Scenario:** Tomás wants to search for failed passes above 30° of elevation, and a number written as text can't be compared.
- **Exercise:** `logFlown(l *slog.Logger, p PassView, lines int)` logs nothing; you log `pass flown` at INFO with `LogAttrs`: `pass_id`, `max_elev_deg` and `lines`.
- **Check:** In JSON `max_elev_deg` decodes as the number 61.4 (and 8.0 as 8, still a number), `lines` as a number and `pass_id` as a string; the mistake, `fmt.Sprint` of the elevation, writes `"61.4"` and fails; exactly one line.

#### 1.5 The Error Says Which Pass
`the-error-says-which-pass` · code_tests · difficulty 6
- **Objective:** Log an error value under `err`, so its whole text, wrapping included, reaches the line.
- **Scenario:** The 02:02:39 failure's error read `context deadline exceeded`: any pass at any station.
- **Exercise:** `logResult(l *slog.Logger, id string, err error)` always logs `pass flown`; you log INFO `pass flown` with `pass_id` for a nil error, and ERROR `pass failed` with `pass_id` and `err` otherwise.
- **Check:** `pass cirrine-1/LAN/16402: context deadline exceeded` appears whole under `err`; a nil error writes INFO with no `err` key (not `"err":null`); the caller's error still matches `context.DeadlineExceeded` with `errors.Is`.

#### 1.6 Milliseconds, Not Nanoseconds
`milliseconds-not-nanoseconds` · code_tests · difficulty 6
- **Objective:** Log a duration as whole milliseconds under a key that names its unit.
- **Scenario:** In JSON a `time.Duration` prints as nanoseconds, and Ingrid read the 01:55 pass's wait, 570000000000, as a phone number.
- **Exercise:** `msAttr(key string, d time.Duration) slog.Attr` returns `slog.Duration(key, d)`; you return `slog.Int64(key+"_ms", d.Milliseconds())`.
- **Check:** `("wait", 9m30s)` → `wait_ms` 570000; `("ack", 1.5s)` → `ack_ms` 1500; 999µs → 0 (whole milliseconds, rounded down); the value is a JSON number.

#### 1.7 Old Code, New Lines
`old-code-new-lines` · code_tests · difficulty 6
- **Objective:** Route an old `*log.Logger` into a slog handler with `slog.NewLogLogger`.
- **Scenario:** bookingd's access log takes a `*log.Logger` (Season 3), and Nkem wants its lines in the same JSON without touching the middleware.
- **Exercise:** `accessLog(h slog.Handler) *log.Logger` returns `log.New(os.Stderr, "", 0)`; you return `slog.NewLogLogger(h, slog.LevelInfo)`.
- **Check:** `Printf("req-000007 POST /bookings 201 236B 4ms")` arrives in the handler's JSON as `msg` at INFO; two calls make two records; nothing reaches standard error.

#### 1.8 The Line Lantern Atoll Needed
`the-line-lantern-atoll-needed` · code_tests · difficulty 7
- **Objective:** Combine a JSON handler, typed attributes, the error and milliseconds into uplinkd's pass line.
- **Scenario:** Friday 10 November, 14:00: the new uplinkd goes to Lantern Atoll, and Nkem replays the 02:02:39 failure to read it whole.
- **Exercise:** `newPassLog(w io.Writer, now func() time.Time) *slog.Logger` and `passLine(l *slog.Logger, p PassView, wait time.Duration, err error)` are stubs; you build the JSON logger through the given `stamp` and write the line with `msAttr`.
- **Check:** The replay decodes to `time` `2028-11-08T02:02:39Z`, ERROR, `pass failed`, `pass_id` `cirrine-1/LAN/16402`, `max_elev_deg` 22.6, `wait_ms` 570000 and the whole `err`; a flown pass is INFO `pass flown` with no `err`; one line per call.

## Chapter 2 · Levels and Groups

`levels-and-groups` · 9 lessons

Choose levels, filter with `HandlerOptions` and change them live with `LevelVar`, attach attributes once with `Logger.With`, group them with `slog.Group`, and skip work nobody reads.

#### 2.1 Only Warnings at Night
`only-warnings-at-night` · code · difficulty 6
- **Objective:** Set a handler's minimum level with `HandlerOptions.Level`.
- **Scenario:** Since Friday's deploy Lantern Atoll's uplinkd writes 2.1 million lines a night, and the night shift's screen should show only what needs a person.
- **Exercise:** The starter's handler options leave `Level` unset and the program logs five records, from DEBUG to ERROR; you set `Level: slog.LevelWarn`.
- **Check:** Exactly `level=WARN msg="carrier lost" site=LAN` then `level=ERROR msg="pass failed" site=LAN`.

#### 2.2 Louder Without a Restart
`louder-without-a-restart` · code_tests · difficulty 7
- **Objective:** Change a running program's level with a `slog.LevelVar`.
- **Scenario:** A restart clears whatever went wrong (Ingrid's Monday restart starts this week), so the night shift turns Lantern Atoll up through uplinkd's debug port instead.
- **Exercise:** `levelHandler(v *slog.LevelVar) http.Handler` answers 204 and changes nothing; you set the level from the body of `PUT /debug/level` with `UnmarshalText`, answering 400 `unknown level "loud"` for anything else and 405 to another method.
- **Check:** Through `httptest`: `debug` enables a DEBUG record; `WARN` in capitals works; `loud` answers 400, once (nothing follows the 400), and leaves the level as it was; GET answers 405.

#### 2.3 Quick Check: Levels
`quick-check-levels` · quiz · difficulty 6
- **Objective:** Check how levels order and filter records.
- **Scenario:** Ingrid's two questions before the night shift gets the setting.
- **Exercise:** Two questions.
- **Check:** DEBUG −4, INFO 0, WARN 4, ERROR 8, and a handler at WARN drops INFO; a `LevelVar` can change while other goroutines log.

#### 2.4 Name the Program Once
`name-the-program-once` · code_tests · difficulty 6
- **Objective:** Attach attributes to every record with `Logger.With`.
- **Scenario:** Half of uplinkd's lines name their program and site, and the other half can't be found.
- **Exercise:** `programLogger(base *slog.Logger, program, site string) *slog.Logger` returns `base`; you return `base.With("program", program, "site", site)`.
- **Check:** Every record from the new logger carries both, before its own attributes; the base logger's records don't, because `With` returns a new logger; loggers for LAN and YST don't mix.

#### 2.5 One Group per Pass
`one-group-per-pass` · code_tests · difficulty 7
- **Objective:** Gather related attributes into one with `slog.Group`.
- **Scenario:** Every program will describe a pass the same way: one `pass` object with its satellite, station, AOS and LOS.
- **Exercise:** `passGroup(p PassView) slog.Attr` returns `slog.String("pass", p.ID)`; you return the bible's group of `sat_id`, `station`, `aos` and `los`.
- **Check:** JSON `"pass":{"sat_id":"cirrine-1","station":"LAN","aos":"2028-11-08T01:55:09Z","los":"2028-11-08T02:02:39Z"}`; text `pass.sat_id=cirrine-1` and the rest; exactly those four keys; times in UTC.

#### 2.6 A Logger for Each Pass
`a-logger-for-each-pass` · code_tests · difficulty 7
- **Objective:** Make a logger that carries one pass through all its records.
- **Scenario:** `runPass` writes six records per pass, and each must say which pass it belongs to; on 8 November no line did.
- **Exercise:** `passLogger(l *slog.Logger, p PassView) *slog.Logger` returns `l`; you return `l.With(slog.String("pass_id", p.ID), passGroup(p))`.
- **Check:** All six records carry `pass_id` and the group; two passes logging at once keep their own; the parent logger is unchanged.

#### 2.7 Don't Build Unread Lines
`dont-build-unread-lines` · code_tests · difficulty 7
- **Objective:** Skip expensive work for a disabled level with `Logger.Enabled`.
- **Scenario:** The DEBUG record that dumps each frame as hex costs more than the pass when nobody reads DEBUG.
- **Exercise:** `logFrame(ctx context.Context, l *slog.Logger, frame []byte, dump func([]byte) string)` always calls `dump`; you call it only when `l.Enabled(ctx, slog.LevelDebug)`.
- **Check:** At INFO `dump` is never called (a counting fake) and nothing is written; at DEBUG it runs once per frame; a `LevelVar` raised mid-run dumps the next frame.

#### 2.8 Quick Check: With and Groups
`quick-check-with-and-groups` · quiz · difficulty 7
- **Objective:** Check what `With` and groups change.
- **Scenario:** Nkem's three questions before the lines go to the station.
- **Exercise:** Three questions about a five-line snippet in the readme.
- **Check:** `With` returns a new logger and leaves the old one as it was; a group is a nested object in JSON and dotted keys in text; attributes from `With` come before the record's own.

#### 2.9 Two Lines per Pass
`two-lines-per-pass` · code_tests · difficulty 7
- **Objective:** Choose each event's level so a night at INFO reads two lines per pass: where it starts and how it ends.
- **Scenario:** Ingrid wants a night at INFO she can read at a glance; Nkem's rule: at INFO a pass is two lines, trouble is WARN or ERROR, and the rest waits at DEBUG.
- **Exercise:** `levelFor(event string) slog.Level` returns INFO for everything; you map the bible's uplinkd messages to their levels, unknown events staying INFO.
- **Check:** Replaying the night of 7–8 November's events: at INFO two lines per pass plus the warnings; at DEBUG every event; each message at the bible's level; an unknown event at INFO.

## Chapter 3 · Logs With Context

`logs-with-context` · 8 lessons

Carry pass and request IDs in `context.Context`, write a `slog.Handler` that adds them to every record and survives `With`, pass them between programs in a header, and read the lines back in order.

#### 3.1 The Pass in the Context
`the-pass-in-the-context` · code_tests · difficulty 7
- **Objective:** Carry a value through `context.Context` under an unexported key type.
- **Scenario:** Monday 20 November: Nkem and you spent the day lining up four programs' lines about one pass by hand; the pass ID must travel with the work.
- **Exercise:** `withPass(ctx context.Context, id string) context.Context` and `passFrom(ctx context.Context) (string, bool)` are stubs; you write them with `context.WithValue` and a `passKey struct{}`.
- **Check:** A round trip; a child context still finds it; a context without one, or with an empty ID, gives `"", false`; a value stored under the string `"pass_id"` is not found.

#### 3.2 A Handler That Reads Context
`a-handler-that-reads-context` · code_tests · difficulty 7
- **Objective:** Write a `slog.Handler` whose `Handle` adds attributes from the context.
- **Scenario:** Nobody should have to remember the pass ID in a log call when the context already has it.
- **Exercise:** `passHandler` wraps another handler and passes records on unchanged; in `Handle` you add `pass_id` (from `passFrom`) and `request_id` (Season 3's `requestIDFrom`, given) when present, with `Record.AddAttrs`, before calling the inner handler.
- **Check:** `InfoContext` records carry both, after their own attributes; records without them in the context don't; a DEBUG record is still dropped at INFO; `Info` without a context adds nothing.

#### 3.3 With Keeps the Wrapper
`with-keeps-the-wrapper` · code_tests · difficulty 8
- **Objective:** Return your own handler, wrapped, from `WithAttrs` and `WithGroup`.
- **Scenario:** After `l.With("program", "uplinkd")`, the pass ID vanished from every line at Lantern Atoll.
- **Exercise:** The starter's `WithAttrs` and `WithGroup` return the inner handler's result; you wrap it in a `passHandler`. The readme's note says that every attribute after `WithGroup` lands inside the group, which the next quiz checks.
- **Check:** `l.With(...).InfoContext(ctx, ...)` still carries `pass_id`, exactly once; `testing/slogtest.TestHandler` passes on your handler; the inner handler still gets the `With` attributes.

#### 3.4 Quick Check: Handlers
`quick-check-handlers` · quiz · difficulty 7
- **Objective:** Check how a handler sees records, contexts and groups.
- **Scenario:** Mirela's three questions before the handler goes into every program.
- **Exercise:** Three questions.
- **Check:** `Handle` gets the context given to `InfoContext` (a background one for `Info`); a wrapper that returns the inner handler from `WithAttrs` drops itself; an attribute added in `Handle` after `WithGroup` lands inside the group, which is why Apogee writes groups with `slog.Group`.

#### 3.5 The Pass Crosses the Wire
`the-pass-crosses-the-wire` · code_tests · difficulty 8
- **Objective:** Carry the pass ID between programs in a header: out with an `http.RoundTripper`, in with middleware.
- **Scenario:** Headquarters' uplinkd hands each pass's commands to Lantern Atoll with `POST /uplink`, and the pass ID stops at the door.
- **Exercise:** `passTransport.RoundTrip` and `acceptPass(next http.Handler) http.Handler` pass requests on unchanged; you set `Apogee-Pass` from the context on a clone of the request, and put the header's value into the handler's context.
- **Check:** Through `httptest`, the station handler's record carries `pass_id`; no pass in the context sends no header and the handler sees none; the caller's request is not modified; an empty header is ignored.

#### 3.6 Find One Pass
`find-one-pass` · code_tests · difficulty 7
- **Objective:** Read JSON lines back with `bufio.Scanner` and keep one pass's records.
- **Scenario:** Next time the night shift must find a pass's lines in a minute, not a Monday.
- **Exercise:** `follow(r io.Reader, passID string) (lines []Line, broken int, err error)` returns nothing; you decode each line into `Line`, keep that pass's lines in order, and count the lines that aren't JSON.
- **Check:** The seven lines of `cirrine-1/LAN/16402` among forty; other passes left out, and a line that only names the pass in its message; two broken lines counted, not fatal; `Raw` holds each line as written; a reader's error returned.

#### 3.7 In Time Order
`in-time-order` · code_tests · difficulty 7
- **Objective:** Merge several programs' lines into one timeline, stably.
- **Scenario:** The four programs' lines come from four computers, each in its own order.
- **Exercise:** `merge(logs ...[]Line) []Line` only joins them; you sort a new slice by time with `slices.SortStableFunc`.
- **Check:** Four logs come out in time order; equal times keep their input order (headquarters' `lines delivered` before Lantern Atoll's `lines queued`); the inputs are unchanged; no logs give an empty result.

#### 3.8 Four Programs, One Pass
`four-programs-one-pass` · code_tests · difficulty 8
- **Objective:** Combine the context, the handler and the header so one search follows a pass through every program.
- **Scenario:** Friday 24 November: Nkem replays the 8 November pass through bookingd, passplan and both uplinkds on the bench.
- **Exercise:** `programLog(w io.Writer, program, site string, now func() time.Time) *slog.Logger` returns a plain JSON logger; you wrap it with `passHandler` and `stamp` and add `program` and `site` with `With`, so the given replay's lines carry everything.
- **Check:** `follow` and `merge` return exactly the bible's seven lines, from bookingd's `booking created` on 6 November to Lantern Atoll's `pass failed` at 02:02:39, each with its `program`, its `site` and one `pass_id`; another pass's lines are left out.

## Chapter 4 · Never in a Log

`never-in-a-log` · 8 lessons

Keep secrets and personal data out of logs: redact headers, give secret types a `LogValue`, replace attributes by key with `ReplaceAttr`, mask addresses and scrub secrets hidden in text.

#### 4.1 The Key in the Headers
`the-key-in-the-headers` · code_tests · difficulty 7
- **Objective:** Redact secret headers before a request's headers are logged.
- **Scenario:** Monday 27 November: Mirela finds Furrowlight's API key in 14 of bookingd's lines, 13 of them DEBUG lines dumped with the request's headers; her test reads the log as anyone with the ops console could.
- **Exercise:** `safeHeader(h http.Header) http.Header` returns `h`; you return a copy in which every value of `Authorization`, `Cookie`, `Set-Cookie` and `Apogee-Signature` is `[redacted]`.
- **Check:** The logged line holds none of the test key's 32 hex digits; other headers are kept; the request's own header is unchanged; each of two values is redacted; a name added in lowercase with `Add` is caught.

#### 4.2 A Value That Logs Itself
`a-value-that-logs-itself` · code_tests · difficulty 7
- **Objective:** Implement `slog.LogValuer` so a secret type logs safely wherever it goes.
- **Scenario:** Mirela: "Don't count on every log call remembering. Make the key refuse to be logged."
- **Exercise:** `type APIKey string` has no `LogValue`; you write it, returning `<customer>.[redacted]`, or `[redacted]` when the key has no dot.
- **Check:** Through the JSON and text handlers, inside a group and through `With`: no secret digits and the customer ID kept (`furrowlight.[redacted]`); a key without a dot gives `[redacted]`.

#### 4.3 String Is Not Enough
`string-is-not-enough` · code_tests · difficulty 7
- **Objective:** See that the JSON handler marshals a struct's fields and ignores `String`, and fix it with `LogValue`.
- **Scenario:** Season 5's `Key` prints `k-2028-10 (secret hidden)` through fmt; Mirela logs one through the JSON handler and reads its secret in base64.
- **Exercise:** `Key` has `String`; you add `LogValue`, returning its `String()` as a string value.
- **Check:** The JSON line holds neither the base64 nor the hex of the test secret and reads `"key":"k-2028-10 (secret hidden)"`; the text handler too; fmt's `%v` is unchanged.

#### 4.4 Quick Check: Secret Values
`quick-check-secret-values` · quiz · difficulty 7
- **Objective:** Check where secrets slip into log lines.
- **Scenario:** Mirela's three questions once the key is replaced.
- **Exercise:** Three questions about a five-line snippet in the readme.
- **Check:** The JSON handler uses `LogValue` or JSON's own encoding, never `String`; `ReplaceAttr` sees every key with its groups, the built-in ones too; a secret inside an error's text passes any rule that looks only at keys.

#### 4.5 Redact by Name
`redact-by-name` · code_tests · difficulty 7
- **Objective:** Replace attributes by key with `HandlerOptions.ReplaceAttr`, at any group depth.
- **Scenario:** Fields named `token` or `password` turn up in lines nobody planned, and Mirela wants them blank by rule.
- **Exercise:** `redact(groups []string, a slog.Attr) slog.Attr` returns `a`; you replace the value of any key on the bible's list (`api_key`, `authorization`, `password`, `secret`, `token`), in any letter case, with `[redacted]`.
- **Check:** At the top level, in a group and in nested groups; `time`, `level` and `msg` untouched; `token_count`, not on the list, kept; the key itself kept.

#### 4.6 Only the Domain
`only-the-domain` · code_tests · difficulty 7
- **Objective:** Keep personal data out of logs by masking an email address.
- **Scenario:** bookingd's lines about failed webhooks carry customers' contact addresses, and a person's address is not Apogee's to copy.
- **Exercise:** `maskEmail(s string) string` returns `s`; you keep the local part's first character, then `***@` and the domain.
- **Check:** `hollis.vey@furrowlight.example` → `h***@furrowlight.example`; a one-character local part works; no `@`, an empty part or two `@` give `[redacted]`.

#### 4.7 Secrets Inside Errors
`secrets-inside-errors` · code_tests · difficulty 8
- **Objective:** Scrub secrets out of free text with a pattern compiled once.
- **Scenario:** One of the 14 lines carried the key inside an error under `err`, where no key name can catch it.
- **Exercise:** `scrub(s string) string` returns `s`; you replace API keys and tokens with `<customer>.[redacted]` using the bible's pattern, compiled once outside the function.
- **Check:** A key inside a sentence; a key and a token in one text; a 31-digit near miss and a pass ID unchanged; the test reads `redact.go` with `go/parser` and finds no `regexp.MustCompile` inside `scrub`.

#### 4.8 Mirela's Planted Lines
`mirelas-planted-lines` · code_tests · difficulty 8
- **Objective:** Combine `LogValue`, `ReplaceAttr`, masking and scrubbing in one safe logger, proved against planted secrets.
- **Scenario:** Friday 1 December, 16:00: Mirela plants twelve lines with secrets everywhere she found one, and none may reach a disk.
- **Exercise:** `safeLogger(w io.Writer) *slog.Logger` returns a plain JSON logger; you give it a `ReplaceAttr` that redacts by name, masks `email` values and scrubs every string and error value.
- **Check:** Her twelve lines (a header dump, a key in `err`, a `Key`, a password in a group, a token in the message, an address under `email`, and more) leave no secret and no full address in the output; pass IDs, times and other attributes are unchanged; twelve lines in, twelve lines out.

## Chapter 5 · Logs That Fit the Disk

`logs-that-fit-the-disk` · 8 lessons

Write log files that stay inside a budget: append safely, count bytes, rotate by size and by UTC day without splitting a line, prune old files, and sample repeated records.

#### 5.1 Append, Never Truncate
`append-never-truncate` · code_tests · difficulty 7
- **Objective:** Open a log file for appending, creating it with mode 0640.
- **Scenario:** Monday 4 December, 08:10, ten minutes into your week on call: the old `disk` rule pages you with Lantern Atoll's disk 98% full; the new log files start with how a file is opened, because a restart in the middle of a day must add to its file.
- **Exercise:** `openLog(path string) (*os.File, error)` uses `os.Create`; you use `os.OpenFile` with `O_APPEND|O_CREATE|O_WRONLY` and 0640.
- **Check:** Opening twice and writing a line each time keeps both; a new file's mode is 0640; a missing folder returns an error.

#### 5.2 Count What You Write
`count-what-you-write` · code_tests · difficulty 7
- **Objective:** Wrap an `io.Writer` to count the bytes it really wrote.
- **Scenario:** To rotate at 64 MiB, the writer must know the file's size without asking the disk after every line.
- **Exercise:** `countingWriter.Write` passes writes through without counting; you add the `n` the inner writer returns, even when it fails.
- **Check:** Several writes add up; a short write that fails counts only what was written and returns its error unchanged; counting starts from a given size.

#### 5.3 Rotate at a Size
`rotate-at-a-size` · code_tests · difficulty 8
- **Objective:** Move to a new file before a write would pass the size limit, without splitting a line.
- **Scenario:** The 41 GB `uplinkd.log` was one file that nobody could open, let alone search.
- **Exercise:** `(*Rotator).Write(p []byte)` always writes to the current file; you rotate first when the file isn't empty and `size+len(p) > MaxBytes`: close it, rename it to `uplinkd-<day>.<n>.log` and open a fresh `uplinkd.log`.
- **Check:** With `MaxBytes` 1 KiB no file passes 1 KiB, except a single longer line alone in its file; names `uplinkd-2028-12-04.1.log`, then `.2.log`; every line whole in exactly one file; the bytes on disk equal the bytes written.

#### 5.4 Quick Check: Rotation
`quick-check-rotation` · quiz · difficulty 7
- **Objective:** Check how rotation keeps lines whole and files readable.
- **Scenario:** Nkem's three questions before the rotator goes to Lantern Atoll.
- **Exercise:** Three questions.
- **Check:** A rename moves a file without copying it, and a reader holding it open keeps reading; one record is one `Write`, so rotating between writes never splits a line; the current file keeps one name that tools can follow.

#### 5.5 A New File Each Day
`a-new-file-each-day` · code_tests · difficulty 8
- **Objective:** Rotate when the UTC day changes, on an injected clock.
- **Scenario:** The night shift asks for "Tuesday's log", and Tuesday has to be one set of files, in UTC.
- **Exercise:** `Write` rotates on size only; you also rotate when the clock's UTC day is after the current file's, naming the old file by its own day.
- **Check:** 23:59:59.999 stays in the day's file and 00:00:00 starts the next; the old file is `uplinkd-2028-12-04.1.log` and the count starts again at 1; a small file still rotates at midnight; an empty file only takes the new day; a clock a second behind doesn't rotate.

#### 5.6 Keep Fourteen Days
`keep-fourteen-days` · code_tests · difficulty 8
- **Objective:** Delete old rotated files by age and by a total size budget.
- **Scenario:** Fourteen days of lines and 1 GiB per program at every station: Ingrid's budget for logs.
- **Exercise:** `prune(dir, program string, now time.Time, keepDays int, maxTotal int64) error` deletes nothing; you delete the program's rotated files more than `keepDays` old, then the oldest until the rest fit in `maxTotal`.
- **Check:** Day 14 kept and day 15 deleted; oldest first by day, then number; the current file, other programs' files and names that don't parse are never touched; a folder that can't be read returns an error.

#### 5.7 Say It Once, Then Count
`say-it-once-then-count` · code_tests · difficulty 8
- **Objective:** Sample repeated records in a handler: the first ten a minute, then one in a hundred, with a count of those dropped.
- **Scenario:** Through a downpour `carrier lost` repeats ten times a second and buries every other line of the night.
- **Exercise:** `(*Sampler).Handle` passes every record; you count per level and message within the record's UTC minute, drop as the bible's rule says and add `dropped` to the next record passed; ERROR records always pass.
- **Check:** 250 `carrier lost` warnings in one minute give 12 lines, the 110th and 210th with `dropped` 99; the next minute's first carries `dropped` 40; an ERROR is never dropped; two messages are counted apart; a logger made with `With` shares the counts.

#### 5.8 The Disk Stays Under Budget
`the-disk-stays-under-budget` · code_tests · difficulty 8
- **Objective:** Combine sampling, rotation and pruning so a station's logs stay inside their budget.
- **Scenario:** Friday 8 December: before the new logging goes to Lantern Atoll, Nkem gives it a bench night worse than any real one, 180,000 `carrier lost` warnings in fifty fades.
- **Exercise:** `newStationLog(dir string, now func() time.Time) (*slog.Logger, io.Closer, error)` returns a logger on one plain file; you wire the given `Sampler`, a JSON handler and the `Rotator` (64 KiB files and a 1 MiB budget in the test), pruning after each rotation.
- **Check:** After the replay the rotated files hold at most 1 MiB and no file passes 64 KiB; every ERROR line and each minute's first ten warnings are there; every file name follows the rule; `Close` closes the current file.

## Chapter 6 · Counting What Matters

`counting-what-matters` · 8 lessons

Count with `sync/atomic`: counters, gauges, labeled counters and histograms by hand, published with `expvar`, and exact totals from readings across restarts.

#### 6.1 A Counter That Never Slips
`a-counter-that-never-slips` · code_tests · difficulty 8
- **Objective:** Count events from many goroutines with `atomic.Uint64`.
- **Scenario:** Monday 11 December. Yusra: Brinewatch wants the exact number of its failed passes in November; Ingrid's log says 5, uplinkd's lines 3 and passplan's 6.
- **Exercise:** `Counter`'s `Inc`, `Add` and `Value` do nothing; you keep the count in an `atomic.Uint64`.
- **Check:** 50 goroutines × 200 `Inc` give exactly 10,000; `Add(5)`; a zero `Counter` reads 0; the test reads `metrics.go` with `go/parser` and finds the `atomic.Uint64` field, since a plain `uint64` can pass by luck.

#### 6.2 A Gauge Goes Both Ways
`a-gauge-goes-both-ways` · code_tests · difficulty 8
- **Objective:** Keep a reading that goes up and down, and bring it back down on every path with `defer`.
- **Scenario:** Ingrid wants `passes_running` on the Board, and after a failed pass it stayed at 1.
- **Exercise:** `track(g *Gauge, run func() error) error` adds 1 and takes it away only when `run` succeeds; you take it away with `defer`.
- **Check:** The gauge reads 1 inside `run` and 0 after a success, a failure or a panic; the error is returned; 20 runs at once end at 0.

#### 6.3 Quick Check: Counters and Gauges
`quick-check-counters-and-gauges` · quiz · difficulty 7
- **Objective:** Check which kind of number to keep, and what a drop means.
- **Scenario:** Nkem's two questions about the numbers on the Board.
- **Exercise:** Two questions.
- **Check:** A counter that drops means its program restarted; `passes_running` is a gauge because it goes down.

#### 6.4 Counting by Satellite
`counting-by-satellite` · code_tests · difficulty 8
- **Objective:** Keep a counter per pair of label values in a map behind a `sync.RWMutex`.
- **Scenario:** Brinewatch's number is one cell of a table: Brine-3's failed passes at Lantern Atoll.
- **Exercise:** `(*CounterVec).With(a, b string) *Counter` returns a new counter on every call; you look it up under `RLock`, and create it under `Lock` after looking again, keyed by `[2]string`.
- **Check:** 40 goroutines counting over seven satellites and three outcomes give exact cells; the same pair returns the same `*Counter`; `Snapshot()` is a copy; an unused pair is absent.

#### 6.5 Buckets for Waiting
`buckets-for-waiting` · code_tests · difficulty 8
- **Objective:** Count observations into fixed buckets with atomic counters.
- **Scenario:** Tomás wants to know how long satellites take to acknowledge, without keeping every number.
- **Exercise:** `(*Histogram).Observe(d time.Duration)` does nothing; you add one to the first bucket whose bound is at or above `d` in seconds (or the last, above all), to the count, and `d` in whole milliseconds to the sum.
- **Check:** 0 s and exactly 0.5 s in the first bucket; exactly 2 s in "≤ 2"; 301 s above all; `count` and `sum_ms`; 30 goroutines × 100 observations exact.

#### 6.6 Published at Debug Vars
`published-at-debug-vars` · code_tests · difficulty 8
- **Objective:** Publish metrics as JSON with `expvar.Publish` and `expvar.Func`.
- **Scenario:** The pager and the Board read every program's numbers at one address, `/debug/vars`.
- **Exercise:** `publish(m *Metrics)` publishes nothing; you publish `passes_total`, `ack_wait_seconds` and `goroutines` as `expvar.Func` values, once (the given `sync.Once`).
- **Check:** `GET /debug/vars` through `expvar.Handler()` decodes to the nested counts, the histogram's `bounds`, `counts`, `count` and `sum_ms`, and a whole number of goroutines; a later `Inc` shows, since values are read at request time; `cmdline` and `memstats` are still there; a second `publish` is harmless.

#### 6.7 Quick Check: What to Count
`quick-check-what-to-count` · quiz · difficulty 8
- **Objective:** Check how to choose metrics and labels.
- **Scenario:** Yusra's three questions before her weekly report.
- **Exercise:** Three questions.
- **Check:** Count events, not log lines (sampled lines undercount); a label's values must be few, never a pass ID; a histogram answers "how many waited over 10 s", which an average can't.

#### 6.8 A Week, Exactly
`a-week-exactly` · code_tests · difficulty 9
- **Objective:** Compute an exact total from counter readings across a restart.
- **Scenario:** Tuesday 19 December, 10:00: Brinewatch's first exact week, though Lantern Atoll's uplinkd restarted on Monday at 08:00 and its counters began again at 0.
- **Exercise:** `increase(rs []Reading) uint64` returns the last value minus the first; you add up the increases between readings, a drop counting as a restart.
- **Check:** The week's 2,017 readings give Brine-3 failed 1 and flown 23; readings 3, 5, 0, 1 give 3; a counter that never moved gives 0; fewer than two readings give 0; the starter's answer wraps around past the restart.

## Chapter 7 · When to Wake Someone

`when-to-wake-someone` · 8 lessons

Turn numbers into pages: rules over a window of results, thresholds held for a time, hysteresis, one alert per problem, routing to the person on call, and silences that end.

#### 7.1 Failures in Six Hours
`failures-in-six-hours` · code_tests · difficulty 8
- **Objective:** Compute a rule's value over a time window from pass results.
- **Scenario:** Tuesday 2 January: on the night of 27–28 December the pager woke Tomás eleven times for one problem; the new rule counts failed passes over six hours.
- **Exercise:** `failedWithin(results []Result, station string, now time.Time, d time.Duration) int` returns 0; you count the station's failed passes with LOS in `(now−d, now]`.
- **Check:** LOS exactly `now−6h` left out and exactly `now` counted; other stations and flown passes left out; skipped passes are not failures; no results give 0.

#### 7.2 Held for a Time
`held-for-a-time` · code_tests · difficulty 8
- **Objective:** Fire only after a value has stayed over its threshold for a while.
- **Scenario:** Lantern Atoll's disk touches 91% for four minutes while a recording lands, which shouldn't wake anyone; fifteen minutes should.
- **Exercise:** `(*ruleState).beforeFiring`, in `state.go`, which the given `(*Pager).Observe` calls while a rule isn't firing, fires on the first sample at `FireAt`; you keep the rule pending until the value has stayed there for `For`, and forget it when a sample falls below.
- **Check:** Over for 14 min 59 s gives no alert; exactly 15 min gives FIRING at that sample; a dip at 10 min starts the wait again; `For` 0 fires at once.

#### 7.3 Quick Check: Thresholds
`quick-check-thresholds` · quiz · difficulty 8
- **Objective:** Check what thresholds, holds and windows do to pages.
- **Scenario:** Ingrid's three questions about the old rules.
- **Exercise:** Three questions.
- **Check:** A rule checked at every LOS pages again for the same problem; holding for a while trades minutes for no false pages; two rules about one problem page twice.

#### 7.4 Fire High, Clear Low
`fire-high-clear-low` · code_tests · difficulty 8
- **Objective:** Use hysteresis: fire at one value and clear only below a lower one.
- **Scenario:** A disk bobbing between 88% and 91% would fire and resolve all night.
- **Exercise:** `(*ruleState).whileFiring`, which `Observe` calls once a rule fires, resolves when a sample falls below `FireAt`; you resolve only below `ClearBelow`.
- **Check:** The bobbing disk gives one FIRING and no RESOLVED; 0.84 gives one RESOLVED; exactly 0.85 stays firing; after resolving, a new rise fires again only after `For`.

#### 7.5 One Page per Problem
`one-page-per-problem` · code_tests · difficulty 8
- **Objective:** Send one alert per problem, per rule and station, and a reminder only after `Repeat`.
- **Scenario:** The old pager sent a page at every LOS while the problem lasted.
- **Exercise:** `whileFiring` sends FIRING for every firing sample; you send it once, again only when `Repeat` has passed, and RESOLVED once, for a Critical rule only.
- **Check:** A night of firing samples gives one FIRING; with `Repeat` 24h a reminder at exactly 24 hours; two stations firing at once get one each; `Repeat` 0 never reminds; a Warning rule sends no RESOLVED.

#### 7.6 The Right Person
`the-right-person` · code_tests · difficulty 8
- **Objective:** Route an alert by severity and the rota.
- **Scenario:** Six of the eleven pages were single failed passes that could have waited for the morning.
- **Exercise:** `route(a Alert, sev Severity, rota []Shift) Alert` leaves `To` empty; you set the on-call name for Critical (the latest shift started at or before the alert), `morning list` for Warning and `log` for Info.
- **Check:** Monday 1 January at 07:59 still Tomás Ferreira-Lund, at 08:00 `No. 039`; 28 December 00:45 Tomás; Warning and Info routed; the rota is not changed.

#### 7.7 Silences That End
`silences-that-end` · code_tests · difficulty 8
- **Objective:** Hold back a station's pages for a window that always ends.
- **Scenario:** At 07:28 Tomás muted Lantern Atoll "until morning", and nothing unmuted it until Ingrid noticed on Friday.
- **Exercise:** `silenced(s []Silence, station string, at time.Time) bool` returns false; you return true when a silence for the station covers `at` (`From` included, `Until` not); the given pager holds those alerts, with a new `held` field in `ruleState`, and routes the rest with `route`.
- **Check:** The generator's silence on 11 January holds an alert at 10:59 and not at 11:00; another station isn't held; a problem still firing pages at the first sample at or after 11:00; one that resolved by then sends nothing.

#### 7.8 The Night of Eleven Pages
`the-night-of-eleven-pages` · code_tests · difficulty 9
- **Objective:** Combine windows, holds, hysteresis, one page per problem and routing into the new rules, checked on a replayed night.
- **Scenario:** Friday 5 January: Ingrid replays the night of 27–28 December before the new rules go live on Monday.
- **Exercise:** `newRules() []Rule` returns the old rules' settings; you write the bible's `pass-failures`, `failed-pass` and `disk-full`, and the given `replay` runs the night through your pager.
- **Check:** Exactly `20:21 FIRING pass-failures LAN 2` to Tomás Ferreira-Lund, `08:45 RESOLVED pass-failures LAN 0` to him, and six lines on the morning list, and nothing else; the old rules on the same replay still give eleven pages.

## Chapter 8 · Error Budgets

`error-budgets` · 7 lessons

Turn a target into numbers in code: availability, a budget in whole passes with integer arithmetic, a rolling window, burn rates, and how long a budget lasts.

#### 8.1 Ninety-Nine and a Half
`ninety-nine-and-a-half` · code_tests · difficulty 8
- **Objective:** Compute availability as passes flown over passes booked, refusing counts that can't be.
- **Scenario:** Monday 8 January: Ingrid, Tomás and Yusra agree what "every pass counts" means in numbers.
- **Exercise:** `availability(flown, booked int) (float64, error)` returns 0; you return flown/booked, with the bible's errors for nothing booked and for more flown than booked.
- **Check:** Lantern Atoll's quarter, 1,177 of 1,199, gives 0.98165 to five places; all flown gives 1; 0 booked gives `no booked passes`; 5 flown of 4 gives `more flown than booked`.

#### 8.2 The Budget in Whole Passes
`budget-in-whole-passes` · code_tests · difficulty 8
- **Objective:** Turn a target in basis points into a budget of whole passes with integer arithmetic.
- **Scenario:** Yusra asks how many passes Lantern Atoll may lose in a quarter and still keep its promise.
- **Exercise:** `budget(t Target, booked int) int` uses floats; you compute `booked * (10000 - int(t)) / 10000`.
- **Check:** 1,199 at 9900 → 11; 1,104 at 9950 → 5; 920 at 9950 → 4; 100 at 9950 → 0; 10,000 at 9995 → 5, where the float version gives 4.

#### 8.3 Quick Check: Targets
`quick-check-targets` · quiz · difficulty 8
- **Objective:** Check what a target promises and what it leaves.
- **Scenario:** Tomás's three questions in the Fishbowl.
- **Exercise:** Three questions.
- **Check:** 100% leaves no room for a single downpour; 99.5% of 1,104 passes leaves 5; a budget left over is room for change, not a score to keep.

#### 8.4 A Rolling Thirty Days
`a-rolling-thirty-days` · code_tests · difficulty 8
- **Objective:** Count booked and lost passes in a rolling window on an injected clock.
- **Scenario:** Customers don't wait for a quarter to end; Ingrid wants the last 30 days every morning.
- **Exercise:** `window(results []Result, station string, now time.Time) (booked, lost int)` counts everything; you count the station's passes with LOS in `(now − 30 days, now]`, lost meaning failed or skipped.
- **Check:** Lantern Atoll to 8 January 00:00 gives 387 and 11; LOS exactly 30 days back left out and exactly `now` counted; skipped counts as lost; other stations and later passes left out.

#### 8.5 How Fast It Burns
`how-fast-it-burns` · code_tests · difficulty 8
- **Objective:** Compute a burn rate: the share lost over the share the target allows.
- **Scenario:** One night at Lantern Atoll spent more than a month's budget, and Ingrid wants to know how fast while it happens.
- **Exercise:** `burnRate(t Target, booked, lost int) float64` returns lost/booked; you divide by the share the target allows.
- **Check:** 2 lost of 4 at 9900 → 50; 1 of 200 at 9950 → 1; nothing booked → 0; at 10000 a loss → +Inf and no loss → 0.

#### 8.6 At This Rate
`at-this-rate` · code_tests · difficulty 8
- **Objective:** Work out how long a budget lasts at the last week's rate.
- **Scenario:** Wednesday 10 January, Yusra: "TL;DR: at this rate, when does each station run out?"
- **Exercise:** `exhaustIn(remaining, lostLast7d int) (time.Duration, bool)` returns 0, false; you return `remaining` × 7 days / `lostLast7d`.
- **Check:** 4 left at 2 a week → 14 days; 5 left at 1 → 35 days; nothing left or overspent → 0, true; nothing lost this week → false.

#### 8.7 The Quarter in Numbers
`the-quarter-in-numbers` · code · difficulty 9
- **Objective:** Combine availability, budgets and what is left into the quarter's report.
- **Scenario:** Friday 12 January: Yusra sends customers Apogee's first report in numbers; the generator's new part was fitted the day before.
- **Exercise:** The quarter's counts are given; you write `row(code string, booked, lost int, t Target) string` in the bible's layout, and `main` prints the header and five rows.
- **Check:** Exactly the bible's six-line table, from `Q4 2028  booked  lost  avail  budget  left` to `CWR         920     4  99.57       4     0`, Lantern Atoll at `98.17` with `-11` left.

## Chapter 9 · Where the Time Goes

`where-the-time-goes` · 8 lessons

Write CPU profiles with `runtime/pprof`, read the hot path, measure inside a test with `testing.Benchmark` and `b.Loop`, and make a planner many times faster without changing its answer.

#### 9.1 Four Minutes on Record
`four-minutes-on-record` · code_tests · difficulty 8
- **Objective:** Write a CPU profile of one piece of work with `pprof.StartCPUProfile` and `StopCPUProfile`.
- **Scenario:** Tuesday 16 January, 05:10: Cape Wren stowed and passplan took 3 min 52 s to plan the rest of the day. Ingrid: "Four minutes is most of a checklist."
- **Exercise:** `profileTo(path string, work func()) error` only runs `work`; you create the file, start profiling, run, stop and close, returning any error.
- **Check:** `cpu.prof` in the lesson folder starts with `1f 8b` and isn't empty; profiling was stopped (a second `StartCPUProfile` works); a path in a missing folder returns the error without running `work`; the readme offers `go tool pprof -top cpu.prof`.

#### 9.2 Quick Check: Reading a Profile
`quick-check-reading-a-profile` · quiz · difficulty 8
- **Objective:** Read a profile's top functions: flat and cumulative time.
- **Scenario:** Nkem pastes passplan's top six lines on the whiteboard.
- **Exercise:** Three questions about a six-line `-top` excerpt in the readme.
- **Check:** `elevationDeg` has the most flat time (81%); `planDay` has the most cumulative time and little of its own; a function called millions of times is a place to call less often, not only to make faster.

#### 9.3 A Benchmark Inside a Test
`a-benchmark-inside-a-test` · code_tests · difficulty 8
- **Objective:** Write a benchmark body with `b.Loop` and run it inside a test with `testing.Benchmark`.
- **Scenario:** The app runs tests, never benchmarks, so Nkem measures inside one.
- **Exercise:** `benchPlan(b *testing.B)` in `bench.go` is empty; you build the inputs once and call `planDay` inside `for b.Loop()`.
- **Check:** `testing.Benchmark(benchPlan)` runs the plan more than once and the setup exactly once (given counters); `NsPerOp` is positive; the plan's result is unchanged.

#### 9.4 Coarse Steps, Then Seconds
`coarse-steps-then-seconds` · code_tests · difficulty 9
- **Objective:** Call the hot function less: sample every 30 seconds instead of every second before bisecting each crossing.
- **Scenario:** The profile says passplan spends 81% of its time in `elevationDeg`, called every second for every satellite and station.
- **Exercise:** `findPasses(el func(time.Time) float64, from, to time.Time, maskDeg float64) []contact` samples every second and bisects each crossing, as the plan's rules say; you sample every 30 s instead, with the same bisection.
- **Check:** The same AOS and LOS, to the second, as the one-second scan for every pass of a 12-hour window over two stations; at least 20 times fewer calls to `el` (counted); at least 5 times faster by `testing.Benchmark`.

#### 9.5 An Index for Bookings
`an-index-for-bookings` · code_tests · difficulty 9
- **Objective:** Replace a search inside a loop with a map built once.
- **Scenario:** For each predicted pass, passplan read all 2,000 bookings of the next fourteen days.
- **Exercise:** `attach(passes []PassView, bookings []Booking) []Planned` scans the bookings for each pass; you build a `map[bookingKey]string` once (satellite, station and AOS in Unix seconds) and look each pass up.
- **Check:** The same booking IDs on 300 passes; a pass without a booking gets `""`; when two bookings match, the first in input order; an AOS in another time zone still matches; at least 5 times faster.

#### 9.6 Quick Check: Measuring Speed
`quick-check-measuring-speed` · quiz · difficulty 8
- **Objective:** Check how to measure speed honestly inside a test.
- **Scenario:** Mirela's three questions about the benchmark numbers.
- **Exercise:** Three questions.
- **Check:** `b.Loop` keeps setup out of what is measured; a check compares two versions on the same Mac, never milliseconds; a pattern compiled inside a loop is compiled every time round.

#### 9.7 Sort Once at the End
`sort-once-at-the-end` · code_tests · difficulty 9
- **Objective:** Sort a result once, after it is built, instead of after every insert.
- **Scenario:** passplan sorted the whole day after adding each pass: 600 passes, 600 sorts.
- **Exercise:** `buildDay(groups [][]Planned) []Planned` appends each pass and sorts every time; you append them all, then sort once by AOS, then ID.
- **Check:** The same order as before, ties broken by ID; the input slices unchanged; at least 5 times faster on 600 passes.

#### 9.8 Planned in Seconds
`planned-in-seconds` · code_tests · difficulty 9
- **Objective:** Combine coarse steps, the index and one sort into a planner that is many times faster and plans the same day.
- **Scenario:** Friday 19 January: passplan plans a day in 6.1 s, and Nkem asks for profiles on every program's debug port, for the next slow thing.
- **Exercise:** `planDay` calls the old helpers; you use `findPasses`, `attach` and one sort, and the test writes `cpu-after.prof` with your `profileTo`.
- **Check:** The same plan as the old planner for 12 hours of three satellites over two stations; at least 20 times fewer `elevationDeg` calls and at least 5 times faster; `cpu-after.prof` written and stopped.

## Chapter 10 · Where the Memory Goes

`where-the-memory-goes` · 8 lessons

Write heap profiles, read memory after a collection, count allocations with `testing.AllocsPerRun`, free big arrays behind small slices, count goroutines and make every goroutine end.

#### 10.1 A Heap Profile on Disk
`a-heap-profile-on-disk` · code_tests · difficulty 9
- **Objective:** Write a heap profile with `pprof.Lookup("heap")` after a garbage collection.
- **Scenario:** Monday 5 February: Lantern Atoll is failing again with the generator fixed, and Nkem has a new rule: nobody restarts uplinkd until its memory has been saved.
- **Exercise:** `heapProfile(path string) error` only creates the file; you run `runtime.GC`, write the heap profile into it and close it, returning any error.
- **Check:** `heap.prof` in the lesson folder starts with `1f 8b` and isn't empty; a bad path returns the error; the readme offers `go tool pprof -top -inuse_space heap.prof` and says uplinkd serves its heap profile on its debug port too.

#### 10.2 Memory in Use Now
`memory-in-use-now` · code_tests · difficulty 9
- **Objective:** Read the heap in use with `runtime.ReadMemStats` after `runtime.GC`.
- **Scenario:** Lantern Atoll's `heap_bytes` climbed 4 MiB at every failed pass and fell back only at restarts.
- **Exercise:** `inUse() uint64` returns `TotalAlloc`; you collect first and return `HeapAlloc`.
- **Check:** Holding 16 buffers of 1 MiB raises `inUse` by at least 12 MiB; letting them go brings it back within 4 MiB of where it started; the test shows why `TotalAlloc` never falls.

#### 10.3 Nothing Allocated per Frame
`nothing-allocated-per-frame` · code_tests · difficulty 9
- **Objective:** Write a function that allocates nothing, checked with `testing.AllocsPerRun`.
- **Scenario:** uplinkd files 500 frames a second at Lantern Atoll and builds a pass ID for each one.
- **Exercise:** `appendPassID(buf []byte, sat, station string, orbit int) []byte` uses `fmt.Sprintf`; you append to `buf` with `append` and `strconv.AppendInt`, padding the orbit to five digits.
- **Check:** `furrow-b/LAN/22379`; `pipit-1/PEL/08154`, padded; a six-digit orbit kept whole; with a buffer that has room, `testing.AllocsPerRun(100, …)` is 0.

#### 10.4 Quick Check: Memory
`quick-check-memory` · quiz · difficulty 9
- **Objective:** Check what heap numbers and a goroutine profile say.
- **Scenario:** Wednesday 7 February: you read the profiles you saved at 02:00.
- **Exercise:** Three questions about an eight-line excerpt of `goroutines.txt` in the readme.
- **Check:** Collect before reading the heap in use; the four goroutines on one line of `runPass` are stuck sending, not working; a goroutine stuck forever keeps everything it can reach, its 4 MiB buffer too.

#### 10.5 Small Slice, Big Array
`small-slice-big-array` · code_tests · difficulty 9
- **Objective:** Copy a small part out of a big buffer so the buffer can be freed.
- **Scenario:** Each pass report kept the last 64 bytes of its frames, and with them the whole buffer.
- **Exercise:** `tail(frames []byte, n int) []byte` returns `frames[len(frames)-n:]`; you return a copy.
- **Check:** The same bytes; keeping 32 tails of 256 KiB buffers holds under 1 MiB, where the starter holds 8 MiB; a buffer shorter than `n` comes back whole, copied.

#### 10.6 Count the Goroutines
`count-the-goroutines` · code_tests · difficulty 9
- **Objective:** Tell whether a function leaves goroutines behind with `runtime.NumGoroutine`.
- **Scenario:** The night of 6–7 February: the night shift calls you at 01:50, and Lantern Atoll's `goroutines` reads 19 where the other stations read 14 or 15.
- **Exercise:** `leftBehind(f func()) int` returns 0; you count goroutines, run `f`, then wait up to one second, checking every millisecond, for the count to come back, and return how many stayed.
- **Check:** A function whose goroutines end gives 0, even when they take 20 ms; one that leaves 3 blocked (planted by the test, released at its end) gives 3; goroutines that were there before never count.

#### 10.7 Stuck on a Send
`stuck-on-a-send` · code_tests · difficulty 10
- **Objective:** Fix a goroutine that blocks forever sending a result nobody will receive.
- **Scenario:** Four of the goroutines you saved are stuck on one line of `runPass`, `done <- result{rep, err}`, each holding a worker.
- **Exercise:** The given `runPass` is the bible's old one; you give `done` room for one and give the worker back with `defer` at the top of the goroutine.
- **Check:** 20 failed passes, the fake radio giving up only after `runPass` did, leave no goroutine behind and all four workers free; a flown pass still returns its result; the leak profile reports none from `runPass`.

#### 10.8 Nothing Left Behind
`nothing-left-behind` · code_tests · difficulty 10
- **Objective:** Combine contexts, workers and leak checks into a `runPass` that leaves nothing behind and says why a pass failed.
- **Scenario:** Friday 9 February, 10:00: the fixed uplinkd goes to every station.
- **Exercise:** `runPass` still starts a goroutine; you take the worker with a `select` on the context, give it back with `defer` and call `fly` with the context, with no goroutine at all.
- **Check:** Replaying the night of 6–7 February, every pass after a failed one gets a worker; goroutines and the heap in use come back to where they started; the failure reads `pass cirrine-1/LAN/17759: no ack before LOS` (`errors.Is(err, ErrNoAck)`); with four passes holding workers past a fifth's LOS, the fifth fails with `ErrNoWorker`.

## Chapter 11 · Following One Request

`following-one-request` · 8 lessons

Follow work through a program and between programs: `runtime/trace` tasks and regions, a flight recorder for the last seconds, spans by hand carried in a context and a header, and a timeline built from them.

#### 11.1 A Trace on Disk
`a-trace-on-disk` · code_tests · difficulty 9
- **Objective:** Record an execution trace with `trace.Start` and `trace.Stop`.
- **Scenario:** Tuesday 13 February: Nkem and you replay the night of 6–7 February on the lab bench, with the old uplinkd and a trace running.
- **Exercise:** `traceTo(path string, work func()) error` only runs `work`; you create the file, start the trace, run, stop and close.
- **Check:** `trace.out` in the lesson folder starts with `go 1.` and holds ` trace`; the trace was stopped (a second `trace.Start` works); a bad path returns the error and starts nothing; the readme offers `go tool trace trace.out`.

#### 11.2 A Task for Each Pass
`a-task-for-each-pass` · code_tests · difficulty 9
- **Objective:** Mark each pass as a task with `trace.NewTask`, ended with `defer`.
- **Scenario:** In the bench trace the night is thousands of goroutine switches; a task gathers each pass's work under one name.
- **Exercise:** `runPass` runs without a task; you start `trace.NewTask(ctx, "pass")`, log the pass ID with `trace.Log`, hand the task's context on, and `defer task.End()`.
- **Check:** The same results with and without a trace running; the test reads `uplink.go` with `go/parser` and finds `trace.NewTask`, `trace.Log`, a deferred `End` and the task's context passed to `fly`.

#### 11.3 Regions for Each Step
`regions-for-each-step` · code_tests · difficulty 9
- **Objective:** Time steps inside a task with `trace.WithRegion` and `trace.StartRegion`.
- **Scenario:** The question is where the failed passes spent their time: waiting for a worker, or flying.
- **Exercise:** `runPass` waits and flies with no regions; you wrap the wait for a worker in a `wait for worker` region and the call to `fly` in a `fly` region.
- **Check:** The same results; the parser finds both region names and every `StartRegion` ended in the function that started it; a trace written during the run starts with `go 1.`.

#### 11.4 Quick Check: Tasks and Regions
`quick-check-tasks-and-regions` · quiz · difficulty 9
- **Objective:** Check what tasks, regions and trace logs mark.
- **Scenario:** Nkem's three questions over the bench trace.
- **Exercise:** Three questions.
- **Check:** A task can span goroutines while a region stays in one; regions nest and end in the goroutine that started them; with no trace running, tasks and regions cost almost nothing (`trace.IsEnabled`).

#### 11.5 A Span by Hand
`a-span-by-hand` · code_tests · difficulty 9
- **Objective:** Start spans that find their parent in the context, with IDs and times from injected sources.
- **Scenario:** A trace file stays on the bench; Nkem wants every uplinkd to write down what each pass did, as lines the night shift can search.
- **Exercise:** `(*Tracer).Start(ctx context.Context, name string)` returns a span with no IDs; you take the trace ID and parent from the context's span (or start a trace), draw the IDs from `Rand` in the bible's order, set `Start`, copy the pass ID from the context and return a context carrying the new span.
- **Check:** A child shares its parent's trace ID and names it as parent; the first root from `rand.NewPCG(1, 2)` is trace `c4f5a58656eef510`, span `9dcec3ad`, and its child `c8d04605`; `End` sets `End` and keeps a copy; spans from two goroutines keep their own parents; every draw from `Rand` holds `t.mu` (a watching source checks it with `TryLock`).

#### 11.6 The Trace Crosses the Wire
`the-trace-crosses-the-wire` · code_tests · difficulty 10
- **Objective:** Carry a span across HTTP in the `Apogee-Trace` header and continue it on the other side.
- **Scenario:** A pass begins at headquarters, when uplinkd hands its lines to Lantern Atoll with `POST /uplink`.
- **Exercise:** `traceTransport.RoundTrip` and `continueTrace(t *Tracer, name string, next http.Handler) http.Handler` pass requests on unchanged; you write `<trace id>-<span id>` from the context's span on a clone of the request, and on arrival start a child of the header's span.
- **Check:** Through `httptest`, the station's span has headquarters' trace ID and the client's span as parent; a missing or malformed header starts a new trace and the request still succeeds; the caller's request is not modified.

#### 11.7 The Last Seconds, Kept
`the-last-seconds-kept` · code_tests · difficulty 10
- **Objective:** Keep the last seconds of a trace with `trace.FlightRecorder` and save them when something fails.
- **Scenario:** Nobody can trace every station all night, but a failed pass should leave its last seconds behind.
- **Exercise:** `watchPass(fr *trace.FlightRecorder, path string, run func() error) error` only calls `run`; you write a snapshot with `WriteTo` to `path` when `run` fails, and nothing when it succeeds.
- **Check:** A failure leaves `flight.trace` starting with `go 1.`; a success leaves no file; `run`'s error is returned, joined with the file's error when the file can't be made; the recorder is still running afterwards.

#### 11.8 The Incident Review
`the-incident-review` · code · difficulty 10
- **Objective:** Combine spans and their parents into the timeline that explains the season.
- **Scenario:** Friday 23 February, 15:00, the Fishbowl: the incident review shows the bench replay's spans, and one picture explains three months.
- **Exercise:** The bench replay's spans are given in `replay.go`; you write `timeline(spans []Span) []string`: roots by start, each followed by its children in the order they started, in the bible's layout; `main` prints the lines.
- **Check:** Exactly the bible's eight lines, from `01:16:08 pass cirrine-1/LAN/17759    11m5s` to `02:22:30   fly                       3m12s`: a `fly` span that outlives its own pass.

## Changes in writing and review

- **Two lines a pass.** At INFO a pass reads `pass started` and then `pass flown` or `pass
  failed`, where it starts and how it ends, so 2.9 is "Two Lines per Pass"
  (`two-lines-per-pass`).
- **The disk page** on Monday 4 December comes at 08:10, ten minutes into No. 039's week on call.
  At 08:25 you turn Lantern Atoll's uplinkd down to INFO, since deleting a file uplinkd holds open
  frees nothing, and Friday's deploy rotates the 41 GB file away. The file is mostly DEBUG `frame
  received` lines with their frames in hex; `carrier lost` is the line that floods a downpour,
  which 5.7's sampler is for. 5.8's flood is a night on the lab bench, worse than any real one.
- **Chapter 8's numbers** come from one rule: the model's passes, with Lantern Atoll's one S-band
  dish flying one pass at a time. That gives 387 booked in the 30 days to 8 January (budget 3,
  remaining −8) and 1,199 in the quarter (22 lost, budget 11, left −11). Elsewhere customers book
  only some of the passes the dishes could fly; Salar Quieto's quarter is 1,288, the most its
  dishes can fly.
- **Chapter 9** plans on the settled pass model: the given `orbit.go`, 1-second samples,
  crossings bisected and rounded half to even, pass IDs from `orbitNumber`. Both planners give
  the model's 2,705 passes of 16–29 January exactly. Speed checks ask for at least 5 times faster
  (measured 24 to 45 times) and call counts for at least 20 times fewer. 9.5's 2,000 test
  bookings carry no IDs.
- **Chapter 7's given pager** sends a problem held in a silence when the silence ends, and its
  reminders with it, and sends a RESOLVED only to someone who heard the FIRING; 7.4's test disk
  flaps the way its readme says.
- **Tests** accept every correct solution the reviewers tried and fail the mistakes they made:
  2.2's missing `return` after a 400, 3.6's `strings.Contains`, 11.5's missing lock. Timing
  checks hold under load: 10.6 waits for a steady goroutine count, and 10.8 and chapter 11 give
  each pass 200 ms.
- **Readmes** say only what is true: 4.7 explains the regular expression it uses, a heap profile
  samples blocks (10.1), `MaxBytes` overrides `MinAge` (11.7), and Go 1.27's leak profile reports
  four of the five stuck goroutines, since uplinkd still holds the fifth one's channel.
