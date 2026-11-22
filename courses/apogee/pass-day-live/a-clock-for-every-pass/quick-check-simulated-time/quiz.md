## Where must the time on a pass's AOS line come from?

- [ ] `time.Now()`, the clock every Go program reads
  > In the tests that clock starts on 1 January 2000, so the line would name a day 29 years before the pass.
- [x] `w.Now()`, the world's clock, read as the line is written
  > The world's clock is the one the tests set to the week, and reading it as the line goes out makes the line say when it really did.
- [ ] `c.AOS`, the instant the plan gives the pass
  > That says when the line should come, not when it did: a pass that woke late would still claim to be on time.

## Why does the pass's context get `WithTimeout(ctx, untilLOS)` and not `WithDeadline(ctx, c.LOS)`?

- [ ] A context made with `WithDeadline` doesn't end when its parent does
  > It ends with its parent like any other context; that isn't the difference.
- [ ] `WithDeadline` ends a second after the instant it is given
  > A deadline ends at its instant exactly; the trouble is which clock the instant is read on.
- [x] A deadline is an instant on the `time` package's clock, which in the tests starts in 2000, 29 years before `c.LOS`
  > A duration measured from `w.Now()` means the same on any clock, so the timeout ends at exactly LOS in the tests and in the real week.

## A station never sends its banner. What ends the pass's read at LOS?

- [x] `hangUp`, which `AfterFunc` runs when `pctx` ends, closing the connection under the read
  > Closing the connection makes the waiting read return an error at once, so the pass ends at its LOS.
- [ ] `pctx` ending, which the read notices and returns
  > A read on a connection knows nothing of contexts; it waits for bytes or for the connection to close.
- [ ] Nothing: a read on a connection gives up by itself after 90 seconds
  > With no deadline set on the connection, a read waits as long as the connection stays open.
