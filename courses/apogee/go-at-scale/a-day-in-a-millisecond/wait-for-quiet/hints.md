Keep the timer in `d.timer`. On the first `Notify` it is `nil`: make it with `time.AfterFunc(d.quiet, d.fire)`. Every later `Notify` calls `Reset` on it instead.

Lock `d.mu` at the top of `Notify` and `Stop`, and `defer` the unlock. In `Stop`, check that `d.timer` isn't `nil` before stopping it: a `Stop` can come before any message.

If `fire` runs eleven times, `Notify` still calls `d.fire` itself, or makes a timer on every call. If it runs at 5s instead of 9s, later calls don't reset the timer. If it runs after `Stop`, `Stop` doesn't stop `d.timer`.
