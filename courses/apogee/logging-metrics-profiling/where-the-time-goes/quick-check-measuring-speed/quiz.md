## `benchPlan` builds its inputs before `for b.Loop()`. Does the time it reports include them?

- [ ] Yes, once, spread over every run of the loop
  > The time per run counts only the loop's body: what comes before the loop runs once and isn't timed at all.
- [x] No: only the loop's body is timed
  > The inputs are built once, outside the measurement, so the time per run is the plan's alone.
- [ ] Only when `testing.Benchmark` calls `benchPlan` more than once
  > With `b.Loop`, `testing.Benchmark` calls `benchPlan` once; the loop's body repeats, and only it is timed.

## Which check of `findPasses`'s speed-up holds on any Mac?

- [x] It is at least 5 times faster than the scan, both timed in the same test
  > Both versions run on the same Mac at the same moment, so the Mac's own speed cancels out of the ratio.
- [ ] It plans twelve hours in under 2 milliseconds
  > Milliseconds depend on the Mac and on what else it is doing: the same code fails on a slower one, or a busy one.
- [ ] Its time per run is lower than in yesterday's run of the test
  > Yesterday's time came from another moment, perhaps another Mac; only times taken side by side compare fairly.

## `validID` runs once for every planned pass. What does each call cost?

- [ ] A match only: Go keeps the compiled pattern for the next call
  > `regexp.MustCompile` builds a new pattern each time it runs, and nothing is kept between calls.
- [ ] Nothing extra, since the pattern is a constant string
  > The text is constant, but compiling it into a pattern to match with is work, done again on every call.
- [x] A compile and a match: the pattern is compiled anew every time round
  > Compiled once, outside the function, the pattern leaves each call only its match.
