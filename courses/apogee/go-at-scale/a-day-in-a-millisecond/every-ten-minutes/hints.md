Make the ticker with `time.NewTicker(d)`, and `defer` its `Stop` on the next line: that turns it off however `every` returns.

Loop forever around a `select` with two cases: a receive from `tick.C`, whose value is the tick's time and goes to `f`, and `<-ctx.Done()`, which returns.

If the test counts 1 call, `f(time.Now())` is still there. If replans of 20 seconds give 139 calls, the loop makes a new timer each round instead of reading one ticker. If `every` calls `f` after its context ended, the `ctx.Done()` case doesn't return.
