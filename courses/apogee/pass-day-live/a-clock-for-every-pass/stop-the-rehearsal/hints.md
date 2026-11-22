A pass still asleep until its AOS has to wake at the stop: make the wait a `time.NewTimer` and a `select` on its channel and on `ctx.Done()`, and when the context ends first, return without the AOS line or the call. `Stop` the timer on the way out.

An open call already hangs up when its pass's context ends, and the pass's context ends with `ctx`, its parent. After the pass, `ctx.Err()` is not nil if the rehearsal was stopped: the LOS line then says `stopped` instead of the outcome.

Check `ctx.Err()` in `rehearse` only after every pass has returned, so the stop's line comes last. Print it with `w.Now()`, and return `context.Cause(ctx)`: the cause the run was stopped with, or `context.Canceled` when it was stopped without one.
