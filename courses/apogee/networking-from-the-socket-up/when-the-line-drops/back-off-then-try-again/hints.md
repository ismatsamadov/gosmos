The count lives in `l.attempt`, so it survives from one failed dial to the next. Add one to it each time `l.Dial` returns an error.

Wait `redialDelay(l.attempt, l.Rand)` in place of `time.Second`, after the count has gone up, and call it once per wait: every call draws a jitter from `l.Rand`.

If the second test's first wait comes out between 15 and 30 seconds, the count from the first drop was still there: set `l.attempt` to 0 right before returning a connection that works.
