A send into `s` takes a slot: it goes through while the channel has room and waits once every slot is taken. A receive from `s` gives one back.

In `Acquire`, put the send in a `select` with a second case for `<-ctx.Done()`: the first case returns `nil`, the second the context's error. `Release` is a single receive from `s`.

If the test reports `most in flight = 80, want 12`, `Acquire` still returns without sending into `s`. If `Acquire` with a canceled context returns `nil` on full slots, the send has no `ctx.Done()` case beside it.
