`fail` is a method of `*Group` whose whole body is one `g.once.Do` call: inside it, set `g.err` and call `g.cancel` with the error, so the first failure is kept and every later one changes nothing. In `Go`'s goroutine, call `fail` when `job` returns an error.

After `Go` has sent into `g.slots`, ask `context.Cause(g.ctx)`. When it isn't `nil`, the group has failed or its caller gave up: receive from `g.slots` to give the slot back, call `fail` with the cause, and return before starting a goroutine.

If `Wait` returns `nil` for a run canceled before its jobs start, the refused jobs aren't passed to `fail`; if `Go` is still waiting there, a refused job kept its slot. If `Wait` returns brine-18's error instead of brine-17's, `fail` keeps the last error rather than the first.
