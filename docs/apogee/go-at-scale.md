# Go at Scale: syllabus

**Authoring status:** planned. Chapter 1 is written but not yet in the library: the season's
pass model (`predict`) is being replaced by the orbit model the earlier seasons use, which
changes the counts, pass IDs and times in this plan and in chapter 1. Once lessons are written,
the lessons in `courses/apogee.sqlite` are the source of truth.

Season 7 of Apogee, March to June 2029. Apogee grows from five stations to twenty and from
eight satellites to eighty, and the learner leads the rewrite of passplan, which must plan a
hundred days instead of one. The world, the characters, the twenty stations and eighty
satellites, the pass model, the planning rule and every shared type, signature, bound and
number are in `story.md`, section "Season 7: scale".

- **Title:** Go at Scale
- **Slug:** `go-at-scale` · **Kind:** Course
- **Short description:** Rebuild Apogee's pass planner for twenty stations and eighty
  satellites, with Go that stays correct and fast under load.
- **Description:** Apogee is growing from five dishes to twenty and from eight satellites to
  eighty, and the program that decides which dish listens to which satellite must now plan a
  hundred days instead of one. You lead its rewrite, with Nkem reviewing every change. You
  build pipelines that stop cleanly, limit how much runs at once, load things once, share
  state without locks and know why that is safe, walk half a million passes without holding
  them, write generic queues and merges, fuzz the parsers a customer's satellite can break,
  test a whole day of timers in a millisecond, measure honestly, cut the garbage, and read
  struct tags with reflection. The season ends on take-over day, when your planner takes over
  from the old one on a live day, under the heaviest schedule Apogee has ever flown. Standard
  library only, and every check runs on your Mac.
- **Story:** an opening, each chapter's briefing and debrief, and the field guide, told as
  `docs/storytelling.md` says, following the bible's spine for the season.

## At a glance

| # | Chapter | Slug | Go topics | Lessons |
|---|---|---|---|---|
| 1 | Stage by Stage | `stage-by-stage` | Build pipelines of goroutines joined by channels: stages that own their channels, composition, cancellation in every stage, a sink, fan-in, errors as a cancel cause, and parallel work that keeps its order. | 9 |
| 2 | Twelve at a Time | `twelve-at-a-time` | Bound concurrency: semaphores from buffered channels, fixed worker pools, back-pressure, sends that never wait, a group that stops at its first error, and a limit per key. | 8 |
| 3 | Once and Only Once | `once-and-only-once` | Initialize once with `sync.Once`, `OnceValues` and `OnceFunc`; count and publish with `sync/atomic`'s types and compare-and-swap; and read Go's memory model through programs whose output proves it. | 9 |
| 4 | Passes on Demand | `passes-on-demand` | Write and use iterators: `iter.Seq` and `iter.Seq2`, range over functions, stopping early, sequences made a day at a time, errors along the way and cleanup on every way out. | 8 |
| 5 | The Iterator Toolkit | `the-iterator-toolkit` | Use `slices.Collect`, `slices.Sorted` and `maps.Keys`, write generic adapters (filter, map, take), and walk two sequences in step with `iter.Pull`. | 9 |
| 6 | Next Pass First | `next-pass-first` | Build generic data structures: orderings with `cmp.Or`, `container/heap` behind a generic priority queue with methods and an iterator, and a k-way merge. | 8 |
| 7 | Fuzzing the Parsers | `fuzzing-the-parsers` | Fuzz with `testing.F`: seed corpora and corpus files, properties instead of answers, NaN and other surprises, byte-slice inputs and round trips. | 8 |
| 8 | A Day in a Millisecond | `a-day-in-a-millisecond` | Test timers, tickers, deadlines and debouncing in fake time with `testing/synctest`, `synctest.Wait`, and checks you write yourself. | 8 |
| 9 | Benchmarks That Tell the Truth | `benchmarks-that-tell-the-truth` | Write benchmarks with `b.Loop`, run them with `testing.Benchmark`, keep setup out, report metrics and allocations, and compare versions fairly. | 8 |
| 10 | Fewer Allocations | `fewer-allocations` | Measure with `testing.AllocsPerRun` and remove allocations: room made ahead, `strings.Builder`, appending into reused buffers, `sync.Pool`, struct keys. | 8 |
| 11 | Reading Struct Tags | `reading-struct-tags` | Read types, kinds, fields, values and tags with `reflect`, set fields through a pointer, and know when reflection isn't worth it. | 10 |
| 12 | A Hundred Pass Days | `a-hundred-pass-days` | Put it all together in passplan 2: lazy stations in bounded slots, a hand-over between planners, replans that change only what's new, the newest run winning, a hundred days under load, and take-over day. | 9 |

102 lessons: 8 code, 77 code_tests and 17 quizzes (83.3% hands-on). Difficulty rises from 7
to 10. Calibration starts with chapter 1; later chapters are written only after it has been
imported and read in the app, and after the shared files (`network.go`, `model.go`,
`plan.go`) have been checked against the bible's reference numbers. Concurrency checks are
correct by construction (barriers and counted jobs, never timing), leak checks count
goroutines after a cancel, time-dependent code is tested in synctest bubbles, benchmarks run a
fixed 20 iterations and are checked by iterations, allocations and their own metrics, never
nanoseconds, and allocation budgets are skipped under `-race`. No real sleep exceeds 50 ms,
nothing leaves the Mac, and no run comes near the app's 10 seconds.

## Chapter 1 · Stage by Stage

`stage-by-stage` · 9 lessons

Build pipelines of goroutines joined by channels: a stage owns and closes the channel it
returns, stages compose, every send and receive also watches the context, a sink drains the
end, a fan-in joins many feeds, a stage's error reaches the whole run as a cancel cause, and a
stage runs on several goroutines without losing its order.

#### 1.1 A Stage Owns Its Channel
`stage-owns-its-channel` · code_tests · difficulty 7
- **Objective:** Write a pipeline stage: a function that starts a goroutine, sends on the channel it returns, and closes that channel when its input ends.
- **Scenario:** Monday 5 March: Tomás has drawn twenty stations on the Fishbowl's glass, and Nkem wants tlmd's ingest rebuilt as a pipeline before the first new station's telemetry feed arrives; its first stage reads one feed.
- **Exercise:** `read(ctx, station string, r io.Reader) <-chan Heard` starts a goroutine that scans `r` and sends each line as a `Heard`, but never closes its channel (room for 16); you `defer close(out)` in the goroutine.
- **Check:** Port Ellery's three-line feed → three `Heard` values in order, each with Station `PEL`, then the channel closes (the test gives up after a second: `read did not close its channel`); an empty feed → closed at once; a 40-line feed, more than the channel's room → every line, in order. The mistake: a channel never closed, so the reader's `range` never ends.

#### 1.2 Stages Join Up
`stages-join-up` · code · difficulty 7
- **Objective:** Compose a pipeline by passing one stage's output channel as the next stage's input.
- **Scenario:** Wednesday 7 March, Nkem's drill begins with one feed: Port Ellery's morning, five lines, one of them garbled by the fog.
- **Exercise:** `main` ranges over `read`'s channel and prints raw lines; you pass that channel to the given `parse` stage and print each frame as `<satellite> <seq> <mode> <volts> V`.
- **Check:** four frame lines (the writer picks Pipit-1 and Cirrine-1 frames, one `SAFE`), the garbled line dropped by `parse`; the starter prints the five raw lines.

#### 1.3 Quick Check: Who Closes
`quick-check-who-closes` · quiz · difficulty 7
- **Objective:** Check who closes a pipeline's channels and what a channel direction promises.
- **Scenario:** Nkem's whiteboard before the drill: three stages and four arrows.
- **Exercise:** Three questions about a three-stage snippet in the readme.
- **Check:** The stage that sends on a channel closes it, never the one that reads it; a function returning `<-chan Heard` lets its caller only receive, so the caller's `close` doesn't compile; a `range` over a channel ends once the channel is closed and drained.

#### 1.4 Stop Every Stage
`stop-every-stage` · code_tests · difficulty 8
- **Objective:** Make every send and receive of a stage also wait on `ctx.Done()`, so cancelling the run stops every goroutine.
- **Scenario:** In the drill the Board read three frames and left, and the parse stage's goroutine sat blocked on its next send for the rest of the day, holding its feed open.
- **Exercise:** `parse(ctx, in <-chan Heard) <-chan TelemetryFrame` ranges over `in` and sends with a plain `out <- f`; you select on `ctx.Done()` for the receive and the send, and return when it's done.
- **Check:** a 50-line feed → every good frame, then closed; a reader that takes 3 frames of 200 and cancels → no goroutine left (leak check; the starter leaves `goroutines left after cancel = 1, want 0`); a context cancelled before the first line → closed within a second; a line `parseLine` refuses is skipped.

#### 1.5 The Newest Frame Wins
`newest-frame-wins` · code_tests · difficulty 7
- **Objective:** Write a sink that drains a channel into a map, keeping one result per key by comparing sequence numbers.
- **Scenario:** Ystmark and Tirrin Rock will both hear Cirrine-1; in the drill its frame 9013 came through Ystmark before 9012 came through Tirrin Rock's slow line, and the Board showed 9012 as the newest.
- **Exercise:** `newest(in <-chan TelemetryFrame) Health` keeps the last frame received for each satellite; you keep the one with the highest `Seq`.
- **Check:** frames in order → the last; 9013 then 9012 → 9013 stays; the same frame from two stations → one entry; satellites kept apart; an empty, closed channel → an empty, non-nil `Health`.

#### 1.6 Twenty Feeds, One Stream
`twenty-feeds-one-stream` · code_tests · difficulty 8
- **Objective:** Fan many channels into one, closing the output only after every input is done, with cancellation in each forwarding goroutine.
- **Scenario:** The drill's twenty feeds, the five real stations and fifteen made from the new stations' positions, all go into one parse stage.
- **Exercise:** `merge(ctx, feeds ...<-chan Heard) <-chan Heard` returns a channel nobody writes to; you forward each feed in a goroutine of a `sync.WaitGroup` (`wg.Go`), selecting on `ctx.Done()`, and close the output from one more goroutine after `wg.Wait()`.
- **Check:** twenty feeds of 1 to 20 lines → all 210 values, each station's lines in their own order; no feeds → closed at once; with nineteen feeds drained and one still open, the output stays open until that one closes; cancelling after 5 values → no goroutine left.

#### 1.7 When a Feed Breaks
`when-a-feed-breaks` · code_tests · difficulty 8
- **Objective:** Carry one stage's error to the whole run with `context.WithCancelCause`, and return `context.Cause` once every stage has stopped.
- **Scenario:** In the drill Ystmark's line dropped in the middle of its feed, and the old ingest carried on as if Ystmark's satellites had gone quiet.
- **Exercise:** `ingest(ctx, feeds map[string]io.Reader) (Health, error)` wires the given `readFeed`, `merge`, `parse` and `newest` but never notices a failure; you make the run's context with `context.WithCancelCause`, hand its cancel to every `readFeed`, and after the sink returns give back `context.Cause(ctx)` when it isn't nil.
- **Check:** three good feeds → every satellite's newest frame, nil; Ystmark's feed failing after two lines with `connection reset` → the error `feed YST: connection reset` (`errors.Is` the reader's error) and no goroutine left; a parent context already cancelled → `context.Canceled`. The mistake: returning `ctx.Err()`, which says "context canceled" and loses which feed broke.

#### 1.8 Quick Check: Pipelines
`quick-check-pipelines` · quiz · difficulty 7
- **Objective:** Check cancellation, causes and fan-in.
- **Scenario:** Nkem's three questions after the drill.
- **Exercise:** Three questions.
- **Check:** A stage blocked on a send nobody receives stays forever unless the send also selects on `ctx.Done()`; `context.Cause` returns the error given to the cancel function, while `ctx.Err()` only says `context.Canceled`; a fan-in closes its output after every forwarding goroutine has returned, never when the first input ends.

#### 1.9 Four Parsers, One Order
`four-parsers-one-order` · code_tests · difficulty 9
- **Objective:** Run one stage on several goroutines and still send its results in input order, with a queue of per-item result channels.
- **Scenario:** Friday 9 March: parsing is the drill's slowest stage, and Nkem wants it on four cores, but each satellite's frames must still come out in the order a station heard them.
- **Exercise:** `parseInOrder(ctx, in <-chan Heard, workers int, parse func(string) (TelemetryFrame, error)) <-chan TelemetryFrame` parses one line at a time; you give each line a result channel of its own, queue those channels in order (room for `workers`), parse in up to `workers` goroutines, and send results by reading the queue in turn.
- **Check:** with a `parse` that holds the first four lines until all four are in flight and releases them last-first → frames in input order, four in flight at once and never five; 2,000 drill lines → every good frame, in order; a refused line skipped without breaking the order; cancelling after 10 frames → no goroutine left. The mistake: a plain worker pool, whose results come out in finishing order.

## Chapter 2 · Twelve at a Time

`twelve-at-a-time` · 8 lessons

Bound concurrency: a semaphore from a buffered channel, a fixed pool of workers, back-pressure
through a bounded queue, a queue that drops its oldest value instead of waiting, a group that
stops at its first error, and a separate limit for each key.

#### 2.1 A Slot for Each Job
`a-slot-for-each-job` · code_tests · difficulty 8
- **Objective:** Limit how many goroutines work at once with a semaphore made of a buffered channel.
- **Scenario:** Monday 12 March: passplan 2's first draft predicted eighty satellites in eighty goroutines at once, and on Friday night its test machine ran out of memory; Nkem: "Twelve planner slots. Never more."
- **Exercise:** `Slots` is a `chan struct{}` made by `newSlots(n)`, and the given `predictAll` acquires one before each job; `Acquire(ctx) error` and `Release()` do nothing; you send into the channel in `Acquire`, selecting on `ctx.Done()`, and receive in `Release`.
- **Check:** 80 fake predictions that hold until released → 12 in flight at the start (barrier) and never more (the starter: `most in flight = 80, want 12`), all 80 results; `Acquire` on a full `Slots` with a cancelled context → `context.Canceled` without waiting; `Release` frees a slot another `Acquire` was waiting for.

#### 2.2 Twelve Workers, One Queue
`twelve-workers-one-queue` · code_tests · difficulty 8
- **Objective:** Bound work with a fixed pool of goroutines that take jobs from one channel.
- **Scenario:** Each planner slot keeps a day of predictions in memory; Nkem wants twelve workers that live for the whole run instead of eighty short goroutines, so the memory is twelve days at most.
- **Exercise:** `runPool(ctx, jobs <-chan Elements, workers int, do func(context.Context, Elements) error) error` starts a goroutine per job; you start exactly `workers` goroutines that range over `jobs`, wait for them, and return the first error.
- **Check:** 80 jobs and 12 workers → 80 calls, at most 12 at once, and no more goroutines than workers while the jobs hold (counted with `runtime.NumGoroutine`); a failing job → its error, returned after the other jobs finish; an empty, closed `jobs` → nil at once.

#### 2.3 Quick Check: Slots and Pools
`quick-check-slots-and-pools` · quiz · difficulty 7
- **Objective:** Check how semaphores and pools bound work.
- **Scenario:** Nkem's whiteboard: slots or workers?
- **Exercise:** Three questions.
- **Check:** With a channel of capacity 12 as a semaphore, the 13th `Acquire` waits until a slot is released; a slot not released on an error path is lost for good, so release it with `defer`; a pool of 12 workers holds 12 goroutines however many jobs there are, while a semaphore starts a goroutine per job and lets 12 work at once.

#### 2.4 Wait for Room
`wait-for-room` · code_tests · difficulty 8
- **Objective:** Apply back-pressure: send into a bounded channel, let the sender wait while it's full, and give up only when the context ends.
- **Scenario:** Varsk Island's dish loader takes contacts over a thin line; the first draft kept every contact for it in a growing slice, and on Wednesday 14 March one slow station's queue held 61,000 contacts.
- **Exercise:** `feed(ctx, contacts []Contact, q chan<- Contact) error` appends everything to a slice that a goroutine of its own drains, and returns at once; you send each contact into `q` in order, selecting on `ctx.Done()`.
- **Check:** a loader that takes everything → every contact in order, nil; a stalled loader and a queue with room for 4 → exactly 4 sent when the test cancels, then `context.Canceled`; a slow loader (a release per contact) → nothing lost or reordered. The mistake: returning before the loader has its contacts.

#### 2.5 Never Wait for the Board
`never-wait-for-the-board` · code_tests · difficulty 8
- **Objective:** Send without waiting with `select` and `default`, dropping the oldest waiting value when the queue is full.
- **Scenario:** The Board's queue holds four schedules; when a screen is slow, a new schedule must push out the oldest waiting one, because the Board only ever needs the newest, and passplan must never wait for it.
- **Exercise:** `offer(q chan *Schedule, s *Schedule) (dropped int)` sends with a plain `q <- s`; you try the send with `default` and, when the queue is full, take one waiting schedule out without waiting and try again, counting what you dropped.
- **Check:** an empty queue → sent, 0 dropped; a full queue holding versions 1 to 4 → offering 5 returns at once with 1 dropped and the queue holds 2, 3, 4, 5; ten offers into a queue nobody reads → versions 7 to 10 remain, in order. The mistake: dropping the new value (`default` alone), which leaves the Board on an old plan.

#### 2.6 First Failure Stops All
`first-failure-stops-all` · code_tests · difficulty 9
- **Objective:** Build a small bounded group: jobs run in at most `limit` goroutines, and the first error cancels the others and is the one `Wait` returns.
- **Scenario:** When one satellite's elements are missing, the other seventy-nine predictions are wasted; Nkem wants the run to stop at its first failure, as your first pass day did.
- **Exercise:** `Group` with `newGroup(ctx, limit)`, `Go(job)` and `Wait()` runs jobs within the limit but ignores their errors; you record the first error once (`sync.Once`), cancel the group's context with it as the cause, and make `Go` run nothing once the group has failed.
- **Check:** twenty jobs that succeed → nil, all run, never more than `limit` at once; the third job failing while the others wait on their context → `Wait` returns that error, every running job saw its context end, and a job given to `Go` afterward never runs; two failures → only the first returned.

#### 2.7 Two Uploads per Station
`two-uploads-per-station` · code_tests · difficulty 9
- **Objective:** Keep a separate limit per key: a map of semaphores made on first use, under a mutex.
- **Scenario:** A host's line carries two uploads at a time, and on Thursday 15 March three uploads to Calidra Atoll's dish loader at once timed out together.
- **Exercise:** `newUploads(perStation int)` and `(*Uploads).Do(ctx, station, upload func() error) error` run every upload at once; you keep one `Slots` per station in a map under a mutex, acquire before the upload and release after.
- **Check:** six uploads each to CAL and OST at once → never more than 2 in flight per station, and 4 in flight together at the barrier; a context cancelled while waiting → its error, and that upload never runs; an upload's error is returned. The mistake: one limit shared by every station, which makes OST wait for CAL.

#### 2.8 Eighty Satellites, Twelve Slots
`eighty-satellites-twelve-slots` · code_tests · difficulty 9
- **Objective:** Combine a bounded group with real work: predict the whole fleet in twelve slots and stop at the first failure.
- **Scenario:** Friday 16 March: passplan 2's predict stage runs on its new machine for the first time, eighty satellites over twenty stations for one day.
- **Exercise:** `predictAll(ctx, els, slots, predict)` predicts one satellite after another; you run each in the given `Group` with limit `slots` and collect the results in a map under a mutex.
- **Check:** 28 June with the model → 80 entries, `furrow-a` with 105 passes, 7,483 in all; never more than 12 at once (a counting predictor); a predictor that fails for `brine-17` with `no elements for brine-17` → that error, every running prediction saw its context end, and none started after it.

## Chapter 3 · Once and Only Once

`once-and-only-once` · 9 lessons

Initialize once with `sync.Once`, `sync.OnceValues` and `sync.OnceFunc`; count with
`atomic.Int64`, publish with `atomic.Pointer` and replace only with something newer by
compare-and-swap; and see, through Go's memory model, why a program without a race behaves.

#### 3.1 Load It Once
`load-it-once` · code_tests · difficulty 8
- **Objective:** Run an expensive initialization exactly once, however many goroutines ask, with `sync.Once`.
- **Scenario:** Monday 19 March: the Board froze for six minutes while a hundred callers each loaded the station table from the network share.
- **Exercise:** `(*Tables).Network() (*Network, error)` calls the given `loadNetwork` on every call; you call it inside `t.once.Do` and keep its result and error in the struct.
- **Check:** 100 goroutines at once → one load (the test's loader counts) and the same `*Network` for every caller; a loader that fails → the same error for every caller, still one load; a zero `Tables` is ready to use.

#### 3.2 One Value for Every Caller
`one-value-for-every-caller` · code_tests · difficulty 8
- **Objective:** Replace a hand-built once with `sync.OnceValues`, and know it remembers an error as firmly as a value.
- **Scenario:** Nkem: "Three fields and a `Once` for one table. Go has a function for that."
- **Exercise:** `networkOnce(load func() (*Network, error)) func() (*Network, error)` returns `load` itself; you return `sync.OnceValues(load)`.
- **Check:** 100 goroutines → 1 load, one pointer; a failing load → the same error for every caller and no second try; two results of `networkOnce` load independently, once each.

#### 3.3 Stop Twice, Close Once
`stop-twice-close-once` · code_tests · difficulty 8
- **Objective:** Make a stop function safe to call any number of times with `sync.OnceFunc`.
- **Scenario:** The Board's shutdown and a cancelled run both call a planning run's `Stop`; on Tuesday 20 March the second call closed an already closed channel and the process died.
- **Exercise:** `newRun()` sets `r.stop` to a function that closes `r.done`; you wrap it in `sync.OnceFunc`.
- **Check:** `Stop` twice → no panic (the test recovers one and reports `second Stop panicked: close of closed channel`) and `Done()` closed; ten goroutines calling `Stop` at once → no panic (each recovers its own); before any `Stop`, `Done()` is still open.

#### 3.4 Quick Check: Once
`quick-check-once` · quiz · difficulty 8
- **Objective:** Check what `sync.Once` and its helpers promise.
- **Scenario:** Nkem's questions before the Board restarts with the new table.
- **Exercise:** Three questions.
- **Check:** When `once.Do` returns, the writes its function made are visible to the caller, even a caller whose call didn't run it; `OnceValues` returns its first error forever, so a load that may fail and later succeed needs something else; a call to `once.Do` from inside its own function waits forever.

#### 3.5 Happens Before
`happens-before` · code · difficulty 8
- **Objective:** See that an atomic store publishes the writes made before it, and only those, to a goroutine whose load sees it.
- **Scenario:** Wednesday 21 March: the test Board read a plan the moment its ready flag went up and got version 0 with no contacts, because the planner raised the flag first and filled the plan after.
- **Exercise:** The planner goroutine stores `ready` true, waits until the Board has looked, then fills `plan`; the Board waits for `ready`, prints and says it has looked; you move the two lines that fill the plan above the `Store`.
- **Check:** `plan v7: 2010 contacts`. The starter prints `plan v0: 0 contacts` on every run: the Board looks after the flag and before the fill. Both versions are free of races (the readme says `go run -race` agrees): the order alone was wrong.

#### 3.6 Counting Without Locks
`counting-without-locks` · code_tests · difficulty 8
- **Objective:** Count from many goroutines with `atomic.Int64`'s `Add` and `Load`.
- **Scenario:** Twelve slots report how many passes they planned and dropped, and the counters, plain ints, came out different on every run.
- **Exercise:** `(*Stats).Add(planned, dropped int)` and `Snapshot()` do nothing; the fields are `atomic.Int64`; you add to `Planned` and `Dropped`, count the run in `Runs`, and load all three.
- **Check:** 12 goroutines × 1,000 calls of `Add(3, 1)` → exactly 36,000 planned, 12,000 dropped and 12,000 runs, every time; a zero `Stats` → 0, 0, 0. (The readme says what `go test -race` reports for plain ints, in a terminal.)

#### 3.7 Publish, Don't Mutate
`publish-dont-mutate` · code_tests · difficulty 9
- **Objective:** Share a changing value by publishing new copies through `atomic.Pointer`, never changing one a reader may hold.
- **Scenario:** After Wednesday's mix-up Nkem reads the publisher and finds it appending each plan's contacts to the very schedule the test Board is reading.
- **Exercise:** `(*Publisher).Update(contacts)` changes the current schedule in place; you build a new `Schedule` (`Version` one more, a slice of its own) and `Store` it.
- **Check:** a schedule taken before `Update` keeps its version and contacts after it; `Current()` after `Update` is the new one; 8 readers loading while 1 writer updates 100 times → every schedule seen has as many contacts as its version (the test's version n has n contacts); `Current()` before any update → nil.

#### 3.8 Quick Check: The Memory Model
`quick-check-memory-model` · quiz · difficulty 8
- **Objective:** Check what Go's memory model promises, and what a data race costs.
- **Scenario:** Nkem's whiteboard question: why does a program without a race behave?
- **Exercise:** Three questions about two short snippets in the readme.
- **Check:** Two goroutines touching one variable, at least one writing, with nothing ordering them, is a data race; a program without races behaves as if its goroutines' steps were taken one at a time in some order, because each unlock, channel operation and atomic operation orders the memory around it; a program with a race promises nothing, and can even see half of a string another goroutine is writing.

#### 3.9 Only a Newer Plan
`only-a-newer-plan` · code_tests · difficulty 9
- **Objective:** Replace a shared value only with a newer one, using a compare-and-swap loop that is right in any interleaving.
- **Scenario:** Wednesday's real trouble: two planning runs finished out of order, and the older plan replaced the newer on the test Board. "What happens when two planners touch the same plan?"
- **Exercise:** `(*Publisher).Offer(s *Schedule) bool` stores `s` whatever its version; you load the current schedule, refuse `s` unless its `Version` is greater, and `CompareAndSwap`, trying again when another offer got in first.
- **Check:** the first offer into an empty publisher → true; offering version 2, then 1 → 2 stays and the second returns false; 100 goroutines offering versions 1 to 100 in a shuffled order (seeded) → version 100 every time; an offer equal to the current version → false.

## Chapter 4 · Passes on Demand

`passes-on-demand` · 8 lessons

Write and use iterators: functions you range over, `iter.Seq` and `iter.Seq2`, stopping the
moment the loop stops, long sequences made a day at a time, errors yielded along the way, and
files closed on every way out of a loop.

#### 4.1 Range Over a Function
`range-over-a-function` · code · difficulty 8
- **Objective:** Write an `iter.Seq` and range over it.
- **Scenario:** Monday 26 March: Yusra's TL;DR says customers want to book a whole season ahead, and Hollis wants the days of his summer listed before he books them.
- **Exercise:** `Days(from time.Time, n int) iter.Seq[time.Time]` returns a function that yields nothing; you yield `n` midnights from `from`, a day apart; `main` prints the first week of the third quarter as `Mon 2 Jan`.
- **Check:** seven lines, `Sun 1 Jul` to `Sat 7 Jul`; the starter prints nothing.

#### 4.2 When the Loop Stops
`when-the-loop-stops` · code_tests · difficulty 8
- **Objective:** Return from an iterator as soon as `yield` returns false.
- **Scenario:** Hollis wanted only Furrow-A's next three passes over Salar Quieto, and the catalog's draft iterator went on predicting the rest of his summer after his loop had stopped, then panicked.
- **Exercise:** `Passes(predict, e, s, w) iter.Seq[PassView]` asks `predict` for one day at a time and yields every pass, ignoring what `yield` returns; you return when it returns false.
- **Check:** breaking after 3 of Furrow-A's July passes over SQT → 3 passes, no panic (the test recovers and reports `breaking after 3 passes panicked: range function continued iteration after function for loop body returned false`) and the predictor asked about 1 July only; all of 1 to 7 July → every pass the model gives, in AOS order; an empty window → nothing, and no call to the predictor.

#### 4.3 Quick Check: Iterators
`quick-check-iterators` · quiz · difficulty 7
- **Objective:** Check how a range loop and an iterator function work together.
- **Scenario:** Nkem's questions before the catalog becomes lazy.
- **Exercise:** Three questions about an eight-line iterator in the readme.
- **Check:** The loop's body runs inside the call to `yield`, so `yield` returns when the body ends; `break` makes `yield` return false, and the iterator must stop; an iterator built from a function starts again from the beginning each time it's ranged over, unless it reads something that runs out.

#### 4.4 Pairs with Seq2
`pairs-with-seq2` · code_tests · difficulty 8
- **Objective:** Yield two values at a time with `iter.Seq2`, as a method of a type.
- **Scenario:** The Board's station panel shows each contact with its dish; Tirrin Rock has two dishes, and the panel's code walked the whole schedule twice to find them.
- **Exercise:** `(*Schedule).At(station string) iter.Seq2[int, Contact]` yields nothing; you yield `(dish, contact)` for each contact at that station, in the schedule's order.
- **Check:** TIR in 2 April's schedule → its 37 contacts in AOS order, 5 of them on dish 2; an unknown station → nothing; breaking after 2 → exactly 2; ranging twice → the same pairs twice.

#### 4.5 A Day at a Time
`a-day-at-a-time` · code_tests · difficulty 9
- **Objective:** Make a long sequence lazily: produce one day, yield it in order, and only then produce the next.
- **Scenario:** The third quarter's catalog holds 672,283 passes, and the old planner hands each service a slice of all of them, about eighty-four megabytes a copy.
- **Exercise:** `Catalog(n, w, predict) iter.Seq[PassView]` predicts the whole window into one slice, sorts it with the given `byAOS` and yields it; you predict every servable pair for one day, sort that day, yield it, then move to the next.
- **Check:** 1 to 7 July → 51,176 passes beginning `brine-27/VSK/00697`, `furrow-s/SRH/00049`, `brine-16/CRP/00697`, every one servable and in `byAOS` order; breaking after 10 passes → every call to the predictor was for 1 July (a counting predictor); the starter asks about every day before yielding one.

#### 4.6 Errors Along the Way
`errors-along-the-way` · code_tests · difficulty 8
- **Objective:** Report failures from inside a sequence with `iter.Seq2[V, error]`, and stop after one.
- **Scenario:** Yusra's service reads the pass archive, a file of last quarter's passes, and a line damaged in a copy made the old reader skip it without a word.
- **Exercise:** `ReadPasses(r io.Reader) iter.Seq2[PassView, error]` yields only the lines it can parse; you yield `(PassView{}, err)` for a bad line, wrapped as `archive line 3: <problem>`, then stop, and yield a scanner's error the same way.
- **Check:** three good lines → three passes with nil errors; a bad third line → two passes, then the error, then nothing; breaking after the first → stops; a reader that fails (`iotest.ErrReader`) → its error.

#### 4.7 Close on Every Way Out
`close-on-every-way-out` · code_tests · difficulty 8
- **Objective:** Release what an iterator opened on every way out (the end, a break, an error) with `defer` inside the iterator function.
- **Scenario:** After a morning of Hollis's questions the archive server reported 412 open files, one for every loop that broke early.
- **Exercise:** `Archive(open func() (io.ReadCloser, error)) iter.Seq2[PassView, error]` opens the file, yields from `ReadPasses` and closes it after the last line; you `defer` the close right after a successful open.
- **Check:** a full read → closed once; a loop that breaks after one pass → closed once; a bad line → closed; an open error → yielded, nothing to close (the test's file counts its closes).

#### 4.8 Every Pass of the Quarter
`every-pass-of-the-quarter` · code_tests · difficulty 9
- **Objective:** Combine lazy generation, a filter and early stopping into a sequence a service can range over.
- **Scenario:** Friday 30 March: Furrowlight books its summer campaign, and Hollis's program wants every pass of Furrowlight's twenty-six satellites in the third quarter.
- **Exercise:** `Bookable(n, w, predict, operator string) iter.Seq[PassView]` is a stub; you build it on the day-at-a-time pattern, predicting only the operator's satellites at stations that can serve them.
- **Check:** Furrowlight, the quarter's first two weeks, 1 to 14 July (a whole quarter takes about a second to predict) → 37,917 passes in `byAOS` order, the first `furrow-s/SRH/00049`; an unknown operator → nothing, without one call to the predictor; breaking after Hollis's first 20 → the predictor asked about 1 July only.

## Chapter 5 · The Iterator Toolkit

`the-iterator-toolkit` · 9 lessons

Use the standard library's iterator helpers (`slices.Collect`, `slices.Sorted`, `maps.Keys`),
write generic adapters of your own (filter, map, take) that stop their sources when they stop,
and walk two sequences in step with `iter.Pull`.

#### 5.1 Collect When You Must
`collect-when-you-must` · code · difficulty 8
- **Objective:** Turn a sequence into a slice with `slices.Collect` and `slices.Sorted`, and get a map's keys as a sequence with `maps.Keys`.
- **Scenario:** Monday 2 April: the first wave goes live, and Ingrid wants the Board's station list in code order, from the map of live sites.
- **Exercise:** `main` prints the codes from a fixed list of the five original stations; you print `slices.Sorted(maps.Keys(live))`, joined by spaces.
- **Check:** `BRM CAL CWR LAN OST PEL SQT TIR VSK YST`; the starter prints the five.

#### 5.2 A Filter of Your Own
`a-filter-of-your-own` · code_tests · difficulty 8
- **Objective:** Write a generic iterator adapter that passes on only the values a function keeps, and stops its source as soon as its own loop stops.
- **Scenario:** The Board wants one station's contacts, bookingd one customer's, and every service has written its own loop with an `if` inside.
- **Exercise:** `Filter[V any](seq iter.Seq[V], keep func(V) bool) iter.Seq[V]` returns `seq` unchanged; you range over `seq`, yield the values `keep` keeps, and return when `yield` says stop.
- **Check:** TIR's contacts of 2 April's schedule → only TIR's 37, in order; nothing kept → nothing; breaking after 2 kept values → the source yields nothing more (a counting source); works for ints too.

#### 5.3 Map and Take
`map-and-take` · code_tests · difficulty 8
- **Objective:** Write the two other everyday adapters: one that changes each value and one that stops after n.
- **Scenario:** Ingrid's corner of the Board shows the next twelve contacts as lines of text.
- **Exercise:** `Map[V, W any](seq, f)` and `Take[V any](seq, n)` return empty sequences; you write both.
- **Check:** `Take` of 2 April's contacts, 12 → the first twelve, and the source asked for exactly 12 values, not 13 (the mistake: counting after asking for one more); `Take` with 0 → nothing, and the source never asked; `Map` of contacts to their pass IDs → the IDs in order; `Take(Map(...), 12)` → twelve strings.

#### 5.4 Quick Check: Adapters
`quick-check-adapters` · quiz · difficulty 8
- **Objective:** Check how adapters chain and when their work happens.
- **Scenario:** Nkem's question: what has run once `Take(Map(Filter(...)))` is built but not yet ranged over?
- **Exercise:** Three questions about a chain in the readme.
- **Check:** Nothing runs until the outer sequence is ranged over; values travel through the chain one at a time, not stage by stage; when `Take` stops, every adapter below it stops and the source's `yield` returns false.

#### 5.5 A Chain for Ingrid
`a-chain-for-ingrid` · code · difficulty 8
- **Objective:** Compose adapters and a collector into one readable chain.
- **Scenario:** Tuesday 3 April: Ingrid wants the first five Furrowlight contacts at the new stations, as Board lines, before the morning poll.
- **Exercise:** `main` loops by hand over 2 April's schedule and prints every contact; you print `slices.Collect(Take(Map(Filter(contacts, furrowAtNew), line), 5))` with the given adapters, one per line.
- **Check:** five lines `<AOS 15:04:05> <pass ID> dish <n>` (the writer takes them from the model's 2 April schedule); the starter prints all 237.

#### 5.6 Pulling Two Schedules
`pulling-two-schedules` · code_tests · difficulty 9
- **Objective:** Turn a push iterator into a pull iterator with `iter.Pull`, and always call its `stop`.
- **Scenario:** From Monday the shadow run compares the old and new planners' schedules side by side, and a range loop can only walk one sequence.
- **Exercise:** `Zip[V any](a, b iter.Seq[V]) iter.Seq2[V, V]` is a stub; you pull from `b` while ranging over `a`, yield pairs until either ends, and `defer stop()`.
- **Check:** schedules of 5 and 7 contacts → 5 pairs in order; one empty → no pairs; breaking after 2 pairs → both sources stopped (the test's sources record that their loops ended); the mistake, a missing `stop`, leaves `b`'s loop unfinished.

#### 5.7 Merge Two Schedules
`merge-two-schedules` · code_tests · difficulty 9
- **Objective:** Merge two sorted sequences into one sorted sequence with two pull iterators.
- **Scenario:** The shadow run walks both planners' contacts as one stream in AOS order.
- **Exercise:** `Merge[V any](a, b iter.Seq[V], cmp func(V, V) int) iter.Seq[V]` yields all of `a`, then all of `b`; you pull from both and always yield the smaller head, `a`'s on a tie.
- **Check:** two interleaved lists of ints → one sorted list; two schedules by `byAOS` → sorted; equal keys → `a`'s first; one empty → the other unchanged; breaking early → both sources stopped.

#### 5.8 Quick Check: Pull Iterators
`quick-check-pull-iterators` · quiz · difficulty 8
- **Objective:** Check the rules of `iter.Pull`.
- **Scenario:** Nkem's questions about the shadow run's merge.
- **Exercise:** Three questions.
- **Check:** Once the sequence has ended, `next` returns the zero value and false, as often as it's called; `stop` must be called when you stop early, and calling it twice is fine; `next` and `stop` must never be called from two goroutines at once.

#### 5.9 The Shadow Run
`the-shadow-run` · code_tests · difficulty 9
- **Objective:** Walk two sorted schedules in step and report every difference, lazily.
- **Scenario:** Thursday 5 April, the shadow run's first full day: the two planners must agree contact by contact before passplan 2 is trusted, and at Tirrin Rock they don't.
- **Exercise:** `Compare(old, new iter.Seq[Contact]) iter.Seq[Change]` is a stub; you pull from both in `byAOS` order and yield `old only`, `new only` and `dish` changes, as the bible defines them.
- **Check:** two identical schedules → no change; 5 April's fixtures, where the new planner was given a turnaround one second short → exactly the planted differences at TIR, in AOS order (the writer computes them from the model); a contact missing from the new schedule → `old only`; breaking after the first change → both sources stopped.

## Chapter 6 · Next Pass First

`next-pass-first` · 8 lessons

Build generic data structures: a comparison with several keys from `cmp.Or`, `container/heap`
behind a generic priority queue with methods and an iterator, and a merge of any number of
sorted sequences.

#### 6.1 Order by Several Keys
`order-by-several-keys` · code_tests · difficulty 8
- **Objective:** Write a three-way comparison with several keys from `cmp.Or`, `time.Time.Compare` and `strings.Compare`.
- **Scenario:** Monday 9 April: two passes with the same AOS came out in different orders on two shadow runs, and the comparison reported a difference that wasn't one.
- **Exercise:** `byAOS(a, b PassView) int` compares only the AOS; you add the pass ID as the second key with `cmp.Or`.
- **Check:** an earlier AOS first; equal AOS → by pass ID (`brine-34/LAN/00403` before `cirrine-10/LAN/00827`, both at 19:09:47 on 11 June); identical passes → 0; `slices.SortFunc` with it gives one order from any starting order (three shuffles with a seeded source).

#### 6.2 A Heap of Passes
`a-heap-of-passes` · code_tests · difficulty 8
- **Objective:** Implement `heap.Interface` and pop passes in AOS order with `container/heap`.
- **Scenario:** Nkem: "A queue that always knows its next pass, without sorting a hundred days."
- **Exercise:** `passHeap` has `Len`, `Less` (by `byAOS`) and `Swap`; `Push` and `Pop` are stubs; you append in `Push`, asserting the `any` to a `PassView`, and remove the last element in `Pop`.
- **Check:** 2 April's 253 candidates pushed in a shuffled order → popped in `byAOS` order; `heap.Init` on an unsorted slice, then pops → sorted; a heap of one → that pass, then empty. The mistake: removing the first element in `Pop` (`container/heap` has already moved the smallest to the end).

#### 6.3 Quick Check: Heaps
`quick-check-heaps` · quiz · difficulty 8
- **Objective:** Check how `container/heap` uses your methods.
- **Scenario:** Nkem's whiteboard: whose `Pop` is whose?
- **Exercise:** Three questions.
- **Check:** You call `heap.Push(h, x)` and `heap.Pop(h)`, and the package calls your methods; your `Pop` removes the last element because the package has moved the smallest there; a push or a pop takes time in proportion to log n, so a queue of 20,000 passes needs about fifteen comparisons per pop.

#### 6.4 A Queue for Any Type
`a-queue-for-any-type` · code_tests · difficulty 9
- **Objective:** Hide `container/heap` behind a generic type with a comparison function and methods.
- **Scenario:** The planner needs queues of passes, of contacts and of a merge's heads; Nkem wants one queue for all of them.
- **Exercise:** `Queue[T]` holds an unexported `heapOf[T]` (a slice and a `cmp`) that implements `heap.Interface`; `NewQueue`, `Push`, `Pop` and `Len` are stubs; you write them on `heap.Push` and `heap.Pop`.
- **Check:** passes by `byAOS` → AOS order; ints by `cmp.Compare` → ascending; `Pop` on an empty queue → the zero value and false; 1,000 ints from a seeded source → sorted.

#### 6.5 Drain in Order
`drain-in-order` · code · difficulty 8
- **Objective:** Give a generic type an iterator method.
- **Scenario:** Ingrid wants Cape Wren's next passes on the Board, straight from the planner's queue.
- **Exercise:** `(*Queue[T]).Drain() iter.Seq[T]` yields nothing; you pop until the queue is empty or the loop stops; `main` prints the first five of Cape Wren's candidates for 9 April (max elevation 10° or more), from a queue.
- **Check:** five lines `<AOS 15:04:05> <pass ID>`, the first `00:26:02 cirrine-1/CWR/18668` (the writer takes the rest from the model); the starter prints nothing.

#### 6.6 Peek Before You Pop
`peek-before-you-pop` · code_tests · difficulty 8
- **Objective:** Read a queue's smallest value without removing it, returning the zero value of a type parameter when there is none.
- **Scenario:** A merge chooses which stream to read next by looking at the queue's head, and must not take it until it's sure.
- **Exercise:** `(*Queue[T]).Peek() (T, bool)` is a stub; you return the first element of the heap's slice, or `var zero T` and false.
- **Check:** after pushing 5, 2 and 8 → 2, true, and `Len` still 3; empty → 0 and false, and "" for a queue of strings; after a `Pop` → the next smallest.

#### 6.7 Merge Eighty Streams
`merge-eighty-streams` · code_tests · difficulty 9
- **Objective:** Merge any number of sorted sequences with a priority queue of their heads and pull iterators.
- **Scenario:** Wednesday 11 April: the shadow run's sorted slices per station give way to eighty lazy streams, one per satellite, merged a pass at a time.
- **Exercise:** `MergeBy[T any](cmp func(a, b T) int, seqs ...iter.Seq[T]) iter.Seq[T]` concatenates its sequences; you pull each, push each head with its sequence's index into a `Queue`, pop the smallest, yield it, and push that sequence's next value.
- **Check:** 80 sorted streams of ints → one sorted stream of every value; the 80 satellites' `Passes` over CWR on 9 April → `byAOS` order, the same as collecting and sorting; equal keys → the earlier sequence's first; breaking after 10 → every source stopped; no sequences → nothing.

#### 6.8 One Station's Day, in Order
`one-stations-day-in-order` · code_tests · difficulty 9
- **Objective:** Combine lazy sequences, a filter and a k-way merge into one station's candidates, then hand them to the planning rule.
- **Scenario:** On Tuesday 17 April at 09:40 Cirrine's twelve leave their rocket, and Tirrin Rock's first day with them must be planned from eighty streams without one big sort.
- **Exercise:** `stationPasses(n, s, w, predict, minElevDeg)` is a stub; you merge the `Passes` of every satellite the station can serve, filtered to at least `minElevDeg`, with `MergeBy` and `byAOS`.
- **Check:** TIR on 17 April → 113 candidates in `byAOS` order, and the given `assign` makes them 102 contacts; VRL, not live until June → nothing; breaking early → no prediction for a later day.

## Chapter 7 · Fuzzing the Parsers

`fuzzing-the-parsers` · 8 lessons

Fuzz with `testing.F`: what `go test` runs (seeds and corpus files) and what `go test -fuzz`
adds, seeds that catch bugs, properties instead of expected answers, NaN, byte-slice inputs,
and a round trip that only a strict parser passes.

#### 7.1 The Fuzzer's Find
`the-fuzzers-find` · code_tests · difficulty 8
- **Objective:** Read a corpus file a fuzz run saved under `testdata/fuzz`, and fix the code until the fuzz test's seeds pass.
- **Scenario:** Monday 16 April: on Saturday a four-field line from Thistle-2 crashed tlmd's parser and stopped every station's telemetry for forty minutes, and Mirela fuzzed the parser in a terminal overnight.
- **Exercise:** `testdata/fuzz/FuzzParseLine/` holds `thistle-2 2290 2029-04-14T05:10:41Z batt_v=7.40`; `parseLine` reads fields 4 and 5 without counting the fields; you return `line: want 6 fields, got 4` when there aren't six.
- **Check:** `FuzzParseLine` over its corpus passes (the starter fails with `parseLine("thistle-2 2290 2029-04-14T05:10:41Z batt_v=7.40") panicked: runtime error: index out of range [4] with length 4`); "" → `line: want 6 fields, got 0`; seven fields → refused; a good line still parses.

#### 7.2 Seeds That Catch Bugs
`seeds-that-catch-bugs` · code_tests · difficulty 8
- **Objective:** Choose seed inputs that make a fuzz test fail on broken code, knowing `go test` runs only the seeds.
- **Scenario:** Mirela plants three broken copies of her reference parser and asks for seeds that would have caught Saturday before Saturday.
- **Exercise:** `lineSeeds` in `seeds.go` holds one good line; the read-only `FuzzParseLine` adds each seed with `f.Add` and checks that the parser never panics and that what it accepts has a battery from 0 to 20 and a known mode; you add seeds.
- **Check:** the reference parser passes; each mutant (one panics on five fields, one accepts `batt_v=nan`, one accepts `mode=nominal`) fails a re-run of `FuzzParseLine` (`APOGEE_MUTANT=<name>`); at least three seeds, all different.

#### 7.3 Quick Check: Fuzz Tests
`quick-check-fuzz-tests` · quiz · difficulty 8
- **Objective:** Check what runs where when a test is a fuzz test.
- **Scenario:** Mirela's three questions before anyone runs the fuzzer overnight again.
- **Exercise:** Three questions.
- **Check:** `go test` runs a fuzz test's seeds, its `f.Add` values and its corpus files, and nothing else; `go test -fuzz=FuzzParseLine` keeps inventing inputs until one fails or its time runs out; the failing input is saved as a file under `testdata/fuzz/FuzzParseLine/`, so plain `go test` runs it from then on.

#### 7.4 A Property, Not an Answer
`a-property-not-an-answer` · code_tests · difficulty 9
- **Objective:** Write a fuzz target's check as a property that holds for every input, instead of an expected value.
- **Scenario:** Mirela: the fuzzer invents lines nobody wrote an answer for, so the check must say what is true of every answer.
- **Exercise:** `checkLine(t, parse, line)` in `checks.go` only calls `parse`; you add: when it accepts the line, the satellite ID isn't empty, the battery is a number from 0 to 20 (`math.IsNaN`), the temperature a number from −100 to 150, and the mode one of the four.
- **Check:** the reference parser passes on the given seeds; mutants accepting `batt_v=NaN`, `temp_c=+Inf`, `mode=nominal` and an empty satellite ID each fail a re-run; a refused line passes the check whatever its error says.

#### 7.5 Not a Number
`not-a-number` · code_tests · difficulty 8
- **Objective:** Refuse NaN and infinities from `strconv.ParseFloat`, with range checks a NaN can't slip through.
- **Scenario:** Saturday's other line, `batt_v=nan` from Thistle-2, sailed through the parser, and the pager's `BattV < 6.6` rule never fired.
- **Exercise:** `parseLine` (from 7.1) refuses `v < 0 || v > 20`, which a NaN passes; you refuse any value for which `v >= 0 && v <= 20` is false, and the same for the temperature, as `batt_v NaN: want 0.00 to 20.00`.
- **Check:** Thistle-2's `batt_v=nan` line and a `temp_c=-Inf` line → refused; 0.00 and 20.00 → accepted; 20.01 → refused; `FuzzParseLine` over Joss's corpus files with 7.4's property → passes.

#### 7.6 Twelve Bytes, Any Bytes
`twelve-bytes-any-bytes` · code_tests · difficulty 9
- **Objective:** Fuzz a binary decoder with `[]byte` inputs and a round-trip property.
- **Scenario:** Pipit beacons now come from new stations' receivers too, and some hand tlmd a few bytes short when a pass ends in the middle of a frame.
- **Exercise:** `decodeBeacon(b []byte) (Beacon, error)` reads `b[11]` before checking the length; you refuse any length but 12 first, as `beacon: want 12 bytes, got 7`.
- **Check:** `FuzzBeacon` over its seeds and corpus (a 7-byte tail) → no panic; every accepted beacon encodes back to the same 12 bytes with `encodeBeacon`; wrong sync bytes and a bad check byte → refused; a valid Pipit-1 beacon → its fields.

#### 7.7 Back the Way It Came
`back-the-way-it-came` · code_tests · difficulty 9
- **Objective:** Use a round trip as a fuzz property: an accepted line must be exactly what formatting its frame writes.
- **Scenario:** Joss's next flight software writes `batt_v=+7.91` and `batt_v=07.91`, which tlmd accepts and writes back as `7.91`: two lines for one frame, and the archive can't tell them apart.
- **Exercise:** `checkLine` (from 7.4) checks values only; you add the round trip: when `parse` accepts `line`, `formatLine(frame)` must equal `line`.
- **Check:** the strict reference parser passes; mutants accepting `+7.91`, `07.91`, a time ending `+00:00` and a sequence with a leading zero each fail a re-run (seeds given).

#### 7.8 Lines Thistle-2 Should Never Send
`lines-thistle-should-never-send` · code_tests · difficulty 9
- **Objective:** Make a parser strict enough that every line it accepts round-trips, and prove it with the whole fuzz harness.
- **Scenario:** Thursday 19 April: Cirrine's twelve send their first telemetry, and Mirela signs off on the parser only when every one of Joss's odd lines is refused and every good line comes back unchanged.
- **Exercise:** `parseLine` (from 7.5) accepts whatever `strconv` accepts; you check each field's form as the bible gives it and finish with the round trip, `line: not canonical`.
- **Check:** `FuzzParseLine` with Joss's fourteen lines, the fuzzer's finds and 7.7's property → passes; the bible's errors for `-1`, `2029-04-14T03:12:00+00:00`, `07.91` and `nominal`; 50 good lines from the drill → accepted and unchanged by the round trip.

## Chapter 8 · A Day in a Millisecond

`a-day-in-a-millisecond` · 8 lessons

Test time-dependent code in fake time with `testing/synctest`: timers, tickers, context
deadlines and a debouncer, bubbles that jump the clock when everyone waits, `synctest.Wait`,
and checks you write yourself.

#### 8.1 Wait Until AOS
`wait-until-aos` · code_tests · difficulty 8
- **Objective:** Wait for an instant with a `time.Timer`, and give up when the context ends.
- **Scenario:** Monday 23 April: the station runner waits for each contact's AOS in a loop that sleeps a minute at a time, and its tests sleep for real.
- **Exercise:** `waitUntil(ctx, now, at) error` sleeps in one-minute steps until `now()` passes `at`; you start one timer for `at.Sub(now())`, select on it and `ctx.Done()`, and stop the timer either way.
- **Check:** (in a synctest bubble shifted to 28 June) waiting for 04:12:25, `furrow-s/CRP/00006`'s AOS, returns at 04:12:25 exactly (the starter returns at 04:13:00); an instant already past → at once; a 30 s deadline on a wait of an hour → `context.DeadlineExceeded` at 30 s; all of it in under a second of real time.

#### 8.2 Your First Bubble
`your-first-bubble` · code_tests · difficulty 9
- **Objective:** Write a test with `synctest.Test`, whose fake clock jumps whenever every goroutine in it is waiting.
- **Scenario:** Nkem: "You'll test timers for the rest of the season. Write one bubble yourself."
- **Exercise:** `checkWaitUntil(t, wait)` in `checks.go` calls `wait` against the real clock; you run it inside `synctest.Test` with the clock shifted to 28 June, and check that a wait for 04:12:25 ends exactly then, that a past instant returns at once, and that a cancelled context gives `context.Canceled`.
- **Check:** the real `waitUntil` passes; mutants (one returns a second early, one ignores the context, one returns at once) each fail a re-run; the check takes under a second.

#### 8.3 Quick Check: Bubbles
`quick-check-bubbles` · quiz · difficulty 8
- **Objective:** Check what time does inside a synctest bubble.
- **Scenario:** Nkem's whiteboard: what does `time.Now()` say inside a bubble?
- **Exercise:** Three questions.
- **Check:** A bubble's clock starts at midnight UTC on 1 January 2000 and moves only when every goroutine in it is blocked on something only the bubble can wake (a channel, a timer, a `Sleep`); a goroutine waiting on a mutex or on real I/O doesn't count, so the clock won't jump for it; a goroutine still blocked when the test function returns makes the bubble report a deadlock.

#### 8.4 Every Ten Minutes
`every-ten-minutes` · code_tests · difficulty 8
- **Objective:** Run work on a `time.Ticker` and stop the ticker when the context ends.
- **Scenario:** passplan 2 replans every ten minutes, and the draft's loop of `time.After` made a new timer every round.
- **Exercise:** `every(ctx, d, f)` calls `f` once and returns; you make a ticker, call `f` with each tick's time until the context ends, and `defer` its `Stop`.
- **Check:** (bubble) a day at 10 minutes → 144 calls, at 00:10, 00:20 and so on to 24:00; cancelling at 09:47 → 58 calls, and `every` has returned (`synctest.Wait`); 6 hours → 4 calls.

#### 8.5 Thirty Seconds to Plan
`thirty-seconds-to-plan` · code_tests · difficulty 8
- **Objective:** Bound a call with `context.WithTimeout` and report a deadline as a clear error.
- **Scenario:** A planning run that hangs holds up every replan behind it, and Nkem wants every run over within thirty seconds.
- **Exercise:** `planWithin(ctx, limit, plan)` calls `plan(ctx)` as it is; you give it a context with a timeout of `limit`, cancel it when done, and wrap a deadline as `plan: took longer than 30s`.
- **Check:** (bubble) a plan that takes 29 s → its schedule; one that takes 31 s and watches its context → the error at exactly 30 s, `errors.Is` `context.DeadlineExceeded`; a parent cancelled at 10 s → `context.Canceled`, not the timeout message.

#### 8.6 Wait for Quiet
`wait-for-quiet` · code_tests · difficulty 9
- **Objective:** Gather a burst of events into one action with a timer that each event pushes back.
- **Scenario:** On Friday 27 April a gale at Ostrel Island sent eleven stow and unstow messages in four seconds, and passplan 2 replanned eleven times.
- **Exercise:** `newDebouncer(quiet, fire)`, `Notify` and `Stop`; `Notify` calls `fire` at once; you start a `time.AfterFunc` on the first `Notify`, `Reset` it on later ones under a mutex, and stop it in `Stop`.
- **Check:** (bubble, `synctest.Wait`) eleven notifies 400 ms apart → one fire, 5 s after the last; two bursts 10 s apart → two fires; `Stop` before the quiet ends → no fire.

#### 8.7 Everyone Is Waiting
`everyone-is-waiting` · code_tests · difficulty 9
- **Objective:** Use `synctest.Wait` to check a moment when every other goroutine is blocked, without sleeping.
- **Scenario:** Nkem wants the debouncer's test written by you, so the next change to it can't break it unnoticed.
- **Exercise:** `checkDebouncer(t, newD)` in `checks.go` is empty; you write a bubble that notifies three times a second apart, checks with `synctest.Wait` that nothing has fired a moment before the quiet ends, waits past it and checks exactly one fire, then checks that `Stop` cancels a pending fire.
- **Check:** the real debouncer passes; mutants (one fires on the first notify, one never resets, one fires twice, one ignores `Stop`) each fail a re-run.

#### 8.8 Fly the Day
`fly-the-day` · code_tests · difficulty 9
- **Objective:** Combine timers, a context and a fake clock to run, and test, a whole day of contacts.
- **Scenario:** Friday 4 May: the station runner flies a day of contacts from the plan, and Ingrid wants it tested on take-over day's plan before any host sees it.
- **Exercise:** `flyDay(ctx, now, contacts, fly)` calls `fly` for every contact at once; you wait for each contact's AOS with `waitUntil`, in order, and stop when the context ends.
- **Check:** (bubble) 28 June's plan as it stands on 4 May (2,013 contacts) → every `fly` at its AOS to the second, in order; cancelled at 12:00 → `context.Canceled`, and exactly the contacts with AOS before 12:00 flown; the whole day in under a second of real time.

## Chapter 9 · Benchmarks That Tell the Truth

`benchmarks-that-tell-the-truth` · 8 lessons

Write benchmarks with `for b.Loop()`, run them with `testing.Benchmark`, keep setup and input
copies out of the measurement with the timer, report your own metrics and allocations, check
that two versions agree before comparing them, and compare them with sub-benchmarks.

#### 9.1 A Benchmark of Your Own
`a-benchmark-of-your-own` · code_tests · difficulty 8
- **Objective:** Write a benchmark function with `for b.Loop()`.
- **Scenario:** Monday 7 May: the second wave goes live, Tomás writes on the glass that passplan 2 is forty times faster, and Nkem asks for a benchmark anyone can run.
- **Exercise:** `benchPlanDay(b *testing.B)` in `bench.go` calls `planDay` once; you build Cape Wren's 28 June candidates with the given helper before the loop and call `planDay` inside `for b.Loop()`.
- **Check:** `testing.Benchmark(benchPlanDay)` at 20 iterations → `N` is 20 and `planDay` ran 20 times (it counts), returning Cape Wren's 75 contacts each time.

#### 9.2 Setup Before the Loop
`setup-before-the-loop` · code_tests · difficulty 8
- **Objective:** Keep setup out of what a benchmark measures.
- **Scenario:** Tomás's forty came from a benchmark that predicted the day inside the loop for the old planner and outside it for the new one.
- **Exercise:** `benchPlanDay` predicts Cape Wren's candidates inside the loop; you move the prediction before `for b.Loop()`.
- **Check:** `AllocsPerOp` at most the planner's own (measured in the test with `testing.AllocsPerRun` on `planDay` alone) plus 2; `N` is 20; every iteration still plans 246 candidates; skipped under `-race`.

#### 9.3 Quick Check: Reading the Numbers
`quick-check-reading-the-numbers` · quiz · difficulty 8
- **Objective:** Read a benchmark's output, and know what `b.Loop` guards against.
- **Scenario:** The output of `go test -bench=. -benchmem` is taped to the glass beside Tomás's forty.
- **Exercise:** Three questions about three output lines in the readme.
- **Check:** `ns/op` is time per iteration, `B/op` the bytes allocated per iteration and `allocs/op` the allocations; inside `for b.Loop()` the compiler keeps every call and its result, while an old `for i := 0; i < b.N; i++` loop could have an unused result optimized away and report almost nothing; one run proves little, so compare several (`-count=10`).

#### 9.4 Fresh Input, Untimed
`fresh-input-untimed` · code_tests · difficulty 9
- **Objective:** Give each iteration fresh input without timing its preparation, with `b.StopTimer` and `b.StartTimer`.
- **Scenario:** `planDay` sorts its candidates in place, so from the second iteration on the benchmark planned an already sorted day: the easy case, every time.
- **Exercise:** `benchFresh(b)` hands `planDay` the same slice every iteration; you copy the unsorted candidates into a working slice with the timer stopped, then start it and plan.
- **Check:** every iteration's input arrived unsorted (`planDay` records it); `AllocsPerOp` within the planner's own, so the copy wasn't counted; `N` is 20.

#### 9.5 Count What Matters
`count-what-matters` · code_tests · difficulty 9
- **Objective:** Report a benchmark's own unit of work with `b.ReportMetric`, and catch an input that runs out after one iteration.
- **Scenario:** The new planner's benchmark read its candidates from a stream that was empty once used: from the second iteration on it planned nothing, very fast.
- **Exercise:** `benchStream(b)` makes the candidate stream once, before the loop, from a pull iterator; you range a reusable sequence instead (`slices.Values` of the day) and report `contacts/op`.
- **Check:** the reported `contacts/op` is 75 for Cape Wren's 28 June; `N` is 20; the starter reports 3.75 (75 over 20 iterations).

#### 9.6 Same Input, Same Answer
`same-input-same-answer` · code_tests · difficulty 8
- **Objective:** Check that two implementations agree before comparing their speed.
- **Scenario:** Wednesday 16 May, 13:30: Brinewatch's thirty-six leave their rocket, and before anyone compares the planners on the week ahead Nkem asks whether they agree.
- **Exercise:** `sameContacts(old, new []Contact) error` returns nil; you compare the lengths and each contact's pass ID and dish, and return the first difference as the bible writes it.
- **Check:** equal lists → nil; a dish swapped → `contact 212: old pipit-2/TIR/22438 dish 1, new pipit-2/TIR/22438 dish 2`-style message for the planted index; different lengths → `old has 2010 contacts, new has 2009`; the given `planOld` and `planNew` on 17 May → nil.

#### 9.7 Quick Check: Fair Comparisons
`quick-check-fair-comparisons` · quiz · difficulty 8
- **Objective:** Check what makes a comparison of two versions fair.
- **Scenario:** Nkem's last questions before the glass is wiped.
- **Exercise:** Three questions.
- **Check:** Both versions run on the same input, on the same machine, in the same session, several times over; a version that gives a different answer isn't faster, it's wrong; a benchmark whose code sorts its input in place measures the sorted case from the second iteration on.

#### 9.8 Is It Really Faster?
`is-it-really-faster` · code_tests · difficulty 9
- **Objective:** Compare two versions fairly with sub-benchmarks (`b.Run`) over the same inputs and sizes.
- **Scenario:** Friday 18 May: the truth replaces the forty on the glass, for six hours, half a day and a day.
- **Exercise:** `benchWindows(b)` in `bench.go` benchmarks the new planner on one day; you run `b.Run` for `old` and `new`, each with sub-benchmarks `hours=6`, `hours=12` and `hours=24` on the same network and windows (longer windows would spend the run predicting passes), and report allocations.
- **Check:** a re-run of the test binary with `-test.bench=^BenchmarkWindows$ -test.benchtime=2x -test.benchmem` lists all six sub-benchmarks (`BenchmarkWindows/old/hours=6` and so on) at 2 iterations; at every size the new planner's `allocs/op` is at most a quarter of the old one's (skipped under `-race`); the two give the same contacts (`sameContacts`).

## Chapter 10 · Fewer Allocations

`fewer-allocations` · 8 lessons

Measure allocations with `testing.AllocsPerRun` and remove them without changing results: room
made ahead with `make`, `strings.Builder`, appending numbers and times into a byte slice, one
buffer reused, `sync.Pool` for scratch space, and struct keys instead of formatted ones.

#### 10.1 Make Room First
`make-room-first` · code_tests · difficulty 8
- **Objective:** Allocate a slice once with `make` and a capacity when you know its size, instead of letting `append` grow it.
- **Scenario:** Monday 21 May: with fifty-six satellites, every replan of a hundred days leaves the planner's memory full of garbage, and the Board's updates stutter while the collector cleans up.
- **Exercise:** `collect(seq, n)` appends to a nil slice; you start from `make([]PassView, 0, n)`.
- **Check:** 1,000 passes → the same passes and 1 allocation (`testing.AllocsPerRun`; the writer records the starter's count, about a dozen); an `n` smaller than the count → still every pass; skipped under `-race` for the counts.

#### 10.2 Build the Report Once
`build-the-report-once` · code_tests · difficulty 8
- **Objective:** Build a long string with `strings.Builder` and `Grow` instead of `+=`.
- **Scenario:** The planner writes each host a text report on every replan, and Cape Wren's 75 lines were built with `+=`, each copying everything before it.
- **Exercise:** `stationReport(code, contacts)` adds each line with `+=`; you write into a `strings.Builder`, grown once for the whole report.
- **Check:** Cape Wren's 28 June report → the same text as before; at most 2 allocations (the starter makes one per line); an empty list → the heading alone.

#### 10.3 Quick Check: Where Garbage Comes From
`quick-check-where-garbage-comes` · quiz · difficulty 8
- **Objective:** Check which lines of code allocate.
- **Scenario:** Nkem's whiteboard: four lines, and which of them allocate?
- **Exercise:** Three questions about four short lines in the readme.
- **Check:** A value that outlives its function (a returned pointer, a slice that escapes) is put on the heap; passing a non-pointer value as `any`, to `fmt.Println` say, usually allocates a copy; converting a `[]byte` to a `string` copies the bytes (and `go build -gcflags=-m` in a terminal shows what escapes).

#### 10.4 Append, Don't Format
`append-dont-format` · code_tests · difficulty 9
- **Objective:** Write numbers and times into a byte slice with `strconv.AppendInt`, `strconv.AppendFloat` and `time.Time.AppendFormat`, without allocating.
- **Scenario:** Each dish loader gets its day as lines like `05:29:26 05:37:52 1 furrow-a/VSK/24529`, each made with `fmt.Sprintf`, 2,006 times a replan.
- **Exercise:** `appendContact(dst, c)` returns `append(dst, fmt.Sprintf(...)...)`; you append each part directly.
- **Check:** the same bytes as the `fmt` version for 28 June's first 50 contacts; 0 allocations when `dst` has room; the line ends in a newline.

#### 10.5 One Buffer for Every Station
`one-buffer-every-station` · code_tests · difficulty 9
- **Objective:** Reuse one buffer across many writes by truncating it with `buf[:0]` instead of making a new one.
- **Scenario:** The export to twenty dish loaders made a new buffer for every station on every replan.
- **Exercise:** `exportStations(w, byStation, codes)` makes a buffer for each station; you keep one, reset it with `buf = buf[:0]` before each station, and write it.
- **Check:** the bytes written equal the old version's for 28 June's twenty stations; the second export allocates at most 2 times in all; a writer's error is returned.

#### 10.6 A Pool of Scratch Space
`a-pool-of-scratch-space` · code_tests · difficulty 9
- **Objective:** Share scratch buffers between goroutines with `sync.Pool`, storing pointers to slices and keeping no oversized ones.
- **Scenario:** Each of the twelve slots makes a fresh slice for a station's candidates every run; Nkem: "Twelve slots, twelve buffers, not a thousand."
- **Exercise:** `getScratch()` makes a new slice every time and `putScratch` throws it away; you keep a `sync.Pool` of `*[]PassView`, hand them out with length 0, and put back only those of capacity 65,536 or less.
- **Check:** after a warm-up, 100 rounds of get, fill and put → under one allocation a round on average; a scratch from the pool has length 0; no get ever returns a capacity above 65,536; 12 goroutines using it at once each see only their own data.

#### 10.7 Keys Without Formatting
`keys-without-formatting` · code_tests · difficulty 8
- **Objective:** Use a comparable struct as a map key instead of a formatted string.
- **Scenario:** The shadow run removes duplicate passes with a map keyed by `fmt.Sprintf` of each pass's parts: three allocations a lookup.
- **Exercise:** `dedupe(passes)` keys its map by a formatted string; you key it by `passKey{SatID, Station, AOS}` (AOS in Unix seconds).
- **Check:** 2 April's 253 candidates, each twice → the 253 once, in first-seen order; a lookup allocates nothing; two satellites with the same AOS at one station → both kept.

#### 10.8 Half the Garbage
`half-the-garbage` · code_tests · difficulty 9
- **Objective:** Combine the chapter's techniques to halve a hot function's allocations without changing its output.
- **Scenario:** Friday 25 May: before the third wave, Nkem's target is the planner's per-station step at half its garbage.
- **Exercise:** `planStationLean(s, cands, turnaround)` is a copy of the given `planStationFat`; you size its slices ahead, use a pooled scratch slice, and drop its per-pass string keys, keeping its output identical.
- **Check:** CWR and YST on 28 June → the same contacts as `planStationFat`; `AllocsPerRun` at most half the fat version's (skipped under `-race`); no candidates → no contacts.

## Chapter 11 · Reading Struct Tags

`reading-struct-tags` · 10 lessons

Read types, kinds, fields, values and struct tags with `reflect`, including the `Fields`
iterators; check rules from tags on ints, floats, durations, strings and slices; set fields
through a pointer; and know when a plain method beats reflection.

#### 11.1 What Kind of Value
`what-kind-of-value` · code · difficulty 8
- **Objective:** Look at a struct's fields with `reflect.TypeOf` and the `Fields` iterator: names, types and kinds.
- **Scenario:** Monday 28 May: the third wave's hosts send their forms, and Nkem wants one validator for every settings struct; first, see what `reflect` sees.
- **Exercise:** `fieldLines(v any) []string` returns nil; you range over `reflect.TypeOf(v).Fields()` and write `<name> <type> <kind>`; `main` prints them for `Config`.
- **Check:** eight lines, from `Slots int int` to `Quiet time.Duration int64` (`Turnaround time.Duration int64` among them); the starter prints nothing.

#### 11.2 Reading a Tag
`reading-a-tag` · code · difficulty 8
- **Objective:** Read a field's tag with `StructTag.Get`, and tell a missing key from an empty one with `Lookup`.
- **Scenario:** Each field's rules are already written in its `check` tag; the validator's first job is to read them.
- **Exercise:** `tagLines(v any) []string` writes each field's whole tag; you write `<name>: <check>`, or `<name>: no check` when the key is missing; `main` prints `PipelineConfig` and the lesson's `Note` struct (one `Text` field with only a `json` tag).
- **Check:** `Workers: min=1,max=32` / `Buffer: min=0,max=1024` / `Feeds: required,len=3` / `Text: no check`.

#### 11.3 Quick Check: Types and Kinds
`quick-check-types-and-kinds` · quiz · difficulty 8
- **Objective:** Check the difference between a type and a kind, and between `Get` and `Lookup`.
- **Scenario:** Nkem's whiteboard: `time.Duration` or `int64`?
- **Exercise:** Three questions.
- **Check:** A type is what was declared (`time.Duration`) and a kind what it's made of (`int64`), so a rule chooses by kind and then by type; `Tag.Get("check")` returns "" both for a missing key and an empty one, and `Lookup` tells them apart; `reflect.TypeOf` of a pointer is a pointer type, and `Elem()` reaches the struct.

#### 11.4 Values Behind the Fields
`values-behind-the-fields` · code_tests · difficulty 9
- **Objective:** Read field values with `reflect.ValueOf` and `Value.Fields`, skipping unexported fields and following a pointer.
- **Scenario:** passplan 2 logs its settings when it starts, and the log line, written by hand, had forgotten two new fields.
- **Exercise:** `describe(v any) []string` returns nil; you range over the value's fields (through a pointer when given one) and write `<name>=<value>` for each exported field, formatted with `%v`.
- **Check:** the default `Config` → `Slots=12`, `Days=100`, `MinElevDeg=10`, `Turnaround=2m0s`, `Replan=10m0s`, `RunLimit=30s`, `Queue=256`, `Quiet=5s`; a `*Config` → the same; a struct with an unexported field → skipped, no panic; a non-struct → nil.

#### 11.5 Required Means Not Zero
`required-means-not-zero` · code_tests · difficulty 8
- **Objective:** Check a value against its type's zero value with `Value.IsZero`, and a slice for elements with `Len`.
- **Scenario:** Ilha Calmaria's first form came back with no station name at all.
- **Exercise:** `checkRule(name, v, rule)` knows no rules yet; you handle `required`: `<name>: required` for a zero value or an empty slice.
- **Check:** `HostConfig.Name` "" → `HostConfig.Name: required`; "Ilha Calmaria" → nil; `Bands` nil and empty → required; `[]string{"S"}` → nil; an unknown rule → `HostConfig.Name: unknown rule "sometimes"`.

#### 11.6 Numbers in Range
`numbers-in-range` · code_tests · difficulty 9
- **Objective:** Compare numbers of any kind with bounds from a tag, durations read as duration text.
- **Scenario:** Vorrel Island's form says 0 dishes, and on Friday 1 June a test deploy's `Slots: 0` planned nothing at all.
- **Exercise:** `checkRule` handles `required`; you add `min=` and `max=` for ints (`Int`), floats (`Float`) and `time.Duration` fields (bounds through `time.ParseDuration`), with the bible's messages.
- **Check:** `Dishes` 0 → `HostConfig.Dishes: 0 is below min 1`; 2 → nil; `Slots` 70 → `Config.Slots: 70 is above max 64`; `MaskDeg` 15.5 → above max 15; `Turnaround` 15m → `Config.Turnaround: 15m0s is above max 10m0s`; `len=3` on an int → `Config.Slots: rule len=3 does not fit int`.

#### 11.7 Only These Values
`only-these-values` · code_tests · difficulty 9
- **Objective:** Check strings and slices of strings against `oneof` and `len` rules.
- **Scenario:** Ilha Calmaria's corrected form lists bands S and K, and Apogee books UHF, S and X.
- **Exercise:** `checkRule` handles numbers; you add `oneof=` (a string, or each string of a `[]string`) and `len=` (characters, not bytes).
- **Check:** `Bands` S, K → `HostConfig.Bands: "K" is not one of UHF S X`; S, X → nil; `Code` "VRLX" → `HostConfig.Code: "VRLX" is not 3 characters`; `Feeds` TIR, YS → `PipelineConfig.Feeds: "YS" is not 3 characters`; a three-letter code with an accented letter → nil (characters, not bytes).

#### 11.8 Defaults from Tags
`defaults-from-tags` · code_tests · difficulty 9
- **Objective:** Set fields through reflection: only through a pointer, only exported fields, only when they can be set.
- **Scenario:** The test deploy's `Slots: 0` was a setting someone left out, not a wish; Nkem: zero means "use the default", and the defaults are already in the tags.
- **Exercise:** `applyDefaults(ptr any) error` changes nothing; you take the pointer's `Elem()` and, for each exported zero field with a `default` tag, set it with `SetInt` (through `time.ParseDuration` for durations), `SetFloat` or `SetString`.
- **Check:** an empty `Config` → `Slots` 12, `Days` 100, `MinElevDeg` 10, `Turnaround` 2m, `Quiet` 5s; `Slots` 6 → kept; a `Config` value, not a pointer → `defaults: want a pointer to a struct, got main.Config`; an unexported field with a `default` tag → left alone, no panic.

#### 11.9 Quick Check: The Price of Reflection
`quick-check-price-of-reflection` · quiz · difficulty 8
- **Objective:** Check when reflection is worth it and when a plain method is better.
- **Scenario:** Nkem's question before the validator ships: why not a `Validate` method on each struct?
- **Exercise:** Three questions.
- **Check:** A `Validate() error` method on one type is checked by the compiler and easier to read, so reflection pays only when many types share the same rules; a mistake in a tag shows up when the code runs, never when it builds; reflection is slower than the code it replaces, so it belongs in start-up checks, not in the planner's inner loop.

#### 11.10 One Validator for Every Config
`one-validator-every-config` · code_tests · difficulty 9
- **Objective:** Combine the field iterator, tag rules and `errors.Join` into one validator for any struct.
- **Scenario:** Friday 8 June: the third wave went live on Monday with every form checked, and from now on passplan 2's own settings go through the same validator when it starts.
- **Exercise:** `Validate(v any) error` is a stub; you check every exported field's rules in order with the given `checkRule`, go into nested structs, and join every problem.
- **Check:** the fifteen hosts' forms → only Vorrel Island's and Ilha Calmaria's problems, in order; the default `Config` → nil; an empty `Config` after `applyDefaults` → nil; a struct with two bad fields → both, joined with `errors.Join` in field order; an int → `validate: want a struct, got int`.

## Chapter 12 · A Hundred Pass Days

`a-hundred-pass-days` · 9 lessons

Everything above in passplan 2: one station planned from lazy merged streams, twenty stations
in twelve slots, the hand-over between two planners, replans that change only what is new,
the newest run winning, stows that come and go, a hundred days under load, and take-over day.

#### 12.1 One Station, Lazily
`one-station-lazily` · code_tests · difficulty 9
- **Objective:** Plan one station from merged lazy streams, watching the context as you go.
- **Scenario:** Monday 11 June: Ingrid sets take-over day for Thursday 28 June, and Nkem asks for the planner's last pieces, starting with one station.
- **Exercise:** `planStation(ctx, n, s, w, free, stows, cfg)` is a stub; you feed `stationPasses`, less the stowed passes, into `assign` with the dishes' free times, checking the context every 1,000 passes, and return the contacts and drops, or `ctx.Err()`.
- **Check:** CWR on 28 June from scratch → 75 contacts, 285 dropped; with the stow from 09:47 to 13:23 → 71; SRH with its dish free at 00:07:00 → first contact `furrow-h/SRH/00003` at 00:08:21; a context cancelled at the start → `context.Canceled` and no contacts.

#### 12.2 Twenty Stations, Twelve Slots
`twenty-stations-twelve-slots` · code_tests · difficulty 9
- **Objective:** Run a job per station in a bounded group and merge their sorted results into one.
- **Scenario:** Tuesday 12 June: the twenty stations, planned twelve at a time on passplan's twelve cores.
- **Exercise:** `Plan(ctx, n, w, free, stows, cfg)` plans the stations one after another and sorts everything; you run `planStation` per station in a `Group` of `cfg.Slots`, keep each station's contacts by its index, and merge them with `MergeBy` and `byAOS`.
- **Check:** 28 June from scratch → 2,013 contacts and 4,087 dropped, the same list as one station at a time; never more than 12 stations in flight; a predictor failing for HSK → its error, the other stations cancelled; `cfg.Slots` 1 → the same schedule.

#### 12.3 Quick Check: The Plan
`quick-check-the-plan` · quiz · difficulty 8
- **Objective:** Check why the planner can work the way it does.
- **Scenario:** Ingrid's questions before the rehearsal, in plain words.
- **Exercise:** Three questions.
- **Check:** Stations can be planned at the same time because no two stations share a dish; the final merge holds one contact per station at a time, twenty heads in its queue; a cancelled run publishes nothing, so the Board keeps the last good schedule.

#### 12.4 Two Planners, One Dish
`two-planners-one-dish` · code_tests · difficulty 9
- **Objective:** Hand over shared state between two owners at an instant, so neither touches what the other holds.
- **Scenario:** Monday 25 June, the rehearsal: planned from scratch at midnight, nine of passplan 2's first contacts took a dish the old planner was still using. "What happens when two planners touch the same pass?"
- **Exercise:** `handover(cur, at, turnaround) map[string][]time.Time` returns nil; you give each station one entry per dish, set from the latest LOS plus turnaround of `cur`'s contacts with AOS before `at`.
- **Check:** the old planner's 27 June schedule at 00:00 on the 28th → 12 dishes busy past midnight, SRH's dish 1 free at 00:07:00; `Plan` with it → 2,010 contacts and none within two minutes of an old contact on the same dish; without it → 10 such conflicts (the test counts them); contacts of `cur` with AOS after `at` → ignored.

#### 12.5 Nothing New, Nothing Changed
`nothing-new-nothing-changed` · code_tests · difficulty 9
- **Objective:** Replan from any instant so that, with no news, the result equals what was already published.
- **Scenario:** passplan 2 replans every ten minutes, and in Monday's rehearsal each replan moved contacts nobody had asked to move, so the Board blinked.
- **Exercise:** `replanFrom(ctx, n, cur, at, stows, cfg)` plans the whole window again from its start; you keep `cur`'s contacts with AOS before `at`, plan from `at` with `handover(cur, at, cfg.Turnaround)`, and join the two.
- **Check:** version 1 of 28 June replanned at 09:47 with no stow → exactly version 1's contacts, 1,229 of them from 09:47 on with the same IDs and dishes; with Cape Wren stowed from 09:47 with no end → 1,978 on the day; the contacts before `at` never change.

#### 12.6 The Newest Run Wins
`the-newest-run-wins` · code_tests · difficulty 10
- **Objective:** Combine cancellation and compare-and-swap so that the newest news always wins, whatever order runs finish in.
- **Scenario:** Tuesday 26 June, rehearsing: a stow arrived halfway through a ten-minute run, the tick's run finished second, and it replaced the stow's plan for four seconds.
- **Exercise:** `(*Replanner).Run(ctx, events)` plans each event in turn and publishes with `Store`; you number the runs, cancel the run in progress when an event arrives, plan in a goroutine, and publish through the publisher's `Offer`.
- **Check:** (bubble, a fake planner that takes 20 s) events at 0 s and 5 s → the first run cancelled, only run 2 published; a run whose plan equals the current one publishes nothing; 30 events arriving in a shuffled order (seeded) → the current schedule comes from the last event; no goroutine left after the context ends.

#### 12.7 Wind at Cape Wren
`wind-at-cape-wren` · code_tests · difficulty 9
- **Objective:** Turn events into replans with the state they carry: a stow with no end, then its end.
- **Scenario:** Wednesday 27 June: Cape Wren's forecast says 74 km/h on Thursday morning, and Ingrid has passplan 2 rehearse it on Thursday's plan: a stow from 09:47 with no end, ended at 13:23.
- **Exercise:** The given `Replanner` plans every event without stows; you keep the stows it has heard of (a stow event with a zero `To` adds one, a later one for the same station fills in its `To`) and plan each event with them.
- **Check:** the start → 2,010 contacts on the day; the stow → a newer version with 1,978, Cape Wren's 32 contacts after 09:47 gone; its end → a newer version with 2,006, Cape Wren at 71, and the five lost passes the bible lists absent.

#### 12.8 A Hundred Days Under Load
`a-hundred-days-under-load` · code_tests · difficulty 10
- **Objective:** Run the planner at its real size, lazily, and cancellable within a day's work at any moment.
- **Scenario:** Wednesday 27 June, 21:00: while Furrowlight's twenty-four are being heard for the first time, the rehearsal's last run plans a hundred days from 28 June with all eighty satellites, and Nkem's test of it plans the first week.
- **Exercise:** `Plan` (from 12.2) checks its context only between stations; you check it before each day's prediction inside the merged streams, so a cancel stops every station within one day of work.
- **Check:** 28 June to 5 July from the hand-over (a hundred days take three seconds to predict, a week a fifth of one) → 14,239 contacts, 42,859 candidates, 28,620 dropped; a predictor that cancels the run when first asked about 2 July → `context.Canceled`, no station asks about a day after 3 July, and no goroutine left; at most 12 stations in flight.

#### 12.9 Take-Over Day
`take-over-day` · code_tests · difficulty 10
- **Objective:** Run passplan 2's whole day: the hand-over, replans on ticks and events, and the newest plan flown, in fake time.
- **Scenario:** Thursday 28 June 2029, 00:00:00 UTC, the cut-over: twenty-four Furrows heard for the first time last night, eighty satellites, twenty stations, and Cape Wren's wind due at 09:47.
- **Exercise:** `takeOver(ctx, n, old, at, events, cfg, now, fly)` is a stub; you hand over from `old` at `at`, run the replanner on the events and a ticker of `cfg.Replan`, and fly each contact of the current schedule at its AOS, reading the schedule again after every replan.
- **Check:** (bubble, a one-day window, a replan every hour, the stow at 09:47 and its end at 13:23) 2,006 contacts flown, each at its AOS, none twice and none on a dish the old planner still held; none of Cape Wren's five lost passes flown; every one of the twenty-four new Furrows flown, 813 contacts; `pipit-1/LAN/16797` flown at 07:48:33; the whole day in under a second of real time.
