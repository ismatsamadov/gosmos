Put the attempt inside a `for` loop: count it in `d.Attempts`, post with a fresh `s.Now().Unix()`, and return nil as soon as `post` does.

After a failure, check in this order: `ctx.Err()` first, returned as it is; then `retryable(statusOf(err), err)` and the attempt count, which end with the `ErrUndelivered` line; only then wait.

`s.schedule(d)` works out the wait and sets `d.NextAt`, and `s.Wait` returns an error when `ctx` ends during it: return that error at once. If the receiver that is down all day gets 9 attempts, compare `d.Attempts` with 8 after counting the attempt.
