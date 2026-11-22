## Passes were loaded at 04:00 with a `ttl` of 10 minutes. What does `passes()` do at exactly 04:10?

- [ ] Returns the kept passes, which go stale only after 04:10
  > At 04:10 the clock is no longer before `loadedAt` plus `ttl`, so `fresh` is already false.
- [ ] Waits a second, then loads them again
  > Nothing in the cache waits: it decides from the time its clock gives it.
- [x] Loads them again
  > 04:10 is not before 04:10, so `fresh` is false and `reload` runs.

## Why does the cache read the time through its `now` field instead of calling `time.Now`?

- [x] So a test can give it any time without waiting for it
  > A test hands the cache a clock of its own and steps it from 04:00 to 04:10 at once.
- [ ] Because `time.Now` cannot be called while a lock is held
  > `time.Now` works anywhere; the field is there so the caller chooses the clock.
- [ ] Because `time.Now` gives local time, and gctl needs UTC
  > Comparing two times works in any zone; the field is about which time, not its zone.

## Why does `passes()` return `slices.Clone(c.cached)` rather than `c.cached` itself?

- [ ] Because `Clone` also sorts the passes by AOS
  > `Clone` copies the slice in the order it has; sorting stays the caller's job.
- [x] So a caller can sort or change its passes without changing the cache
  > Each caller gets its own copy, which neither a refresh nor another caller can touch.
- [ ] Because returning a field of a locked struct does not compile
  > It compiles; the copy decides who can change the cache's slice.
