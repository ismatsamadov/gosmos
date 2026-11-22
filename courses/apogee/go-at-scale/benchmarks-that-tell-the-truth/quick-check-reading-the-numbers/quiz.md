## What does `9 allocs/op` on PlanDay's line count?

- [ ] The values `planDay` left on the heap when the run ended
  > allocs/op counts every allocation, kept or not; what stays in use is a question for a heap profile.
- [x] The allocations one turn of the loop made, on average
  > The run's allocations divided by its 110263 turns: each plan allocates 9 times, 42000 bytes in all, the B/op beside it.
- [ ] The allocations of the whole run of 110263 turns
  > The whole run made about a million; the line divides them by the turns, as ns/op divides the time.

## Why does `BenchmarkHigh` report 0.3589 ns/op?

- [ ] Comparing one number with 10 takes about a third of a nanosecond
  > Kept inside `for b.Loop()`, the same call takes several times as long; under 0.4 ns is the loop counting, with nothing in it.
- [ ] It ran 12 calls at once, one on each core
  > The -12 is how many cores the run may use; a benchmark's loop runs on one goroutine, one call at a time.
- [x] Nothing used `high`'s result, so the compiler dropped the call
  > The loop timed only its own counting. Inside `for b.Loop()` every call and its result are kept, so the call would be timed.

## PlanDay says 11078 ns/op this morning and 10406 last night. What goes on the glass?

- [x] Nothing yet: run it with `-count=10` and look at the spread
  > Ten runs of the same code show how far its time moves on its own; only a change bigger than that spread means something.
- [ ] planDay got 6% slower overnight
  > Two runs can't tell a change from the noise: on a busy machine the same code moves that much from one run to the next.
- [ ] 10742 ns/op, the average of the two
  > Two runs can't show how far the number wanders, so their average hides as much as it tells.
