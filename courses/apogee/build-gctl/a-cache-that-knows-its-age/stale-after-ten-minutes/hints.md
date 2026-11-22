`refresh` is the loading that `passes` already does, plus one more thing to keep: the time, from `c.now()`, in `c.loadedAt`. It returns only an error, and `passes` hands that error back.

Stale means the time to live has run out: `c.now()` is not before `c.loadedAt.Add(c.ttl)`. For passes loaded at 04:00 with 10 minutes, `Before` is false at exactly 04:10, so they are stale from that moment.

In `newShell`, set the cache's `now` to `time.Now` itself, without parentheses, and its `ttl` to `10 * time.Minute`. Add `"time"` to the imports of `shell.go`.
