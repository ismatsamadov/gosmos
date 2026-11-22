`s.Start(s.ln)` serves the watch port in goroutines of its own and returns at once: call it once, before the loop.

In a `for` loop, `select` on `ctx.Done()` and on `ticks`. On a tick, go through `pollAll(s.Dishes)` in order, broadcasting each result's `watchLine`: the results come back in the order of `s.Dishes`, Dish 1 first.

If the watch port accepts no watcher, it never started. If a watcher is still connected after the cancel, `Run` returned without calling `s.Shutdown()`.
