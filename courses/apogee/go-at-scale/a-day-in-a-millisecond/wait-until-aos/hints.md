A timer made for `at.Sub(now())` fires exactly at `at`, so the wait needs no loop: one `time.NewTimer`, then one `select`.

Put `defer t.Stop()` right after the timer is made. The `select` has two cases: a receive from `t.C` returns `nil`, and one from `ctx.Done()` returns `ctx.Err()`, which says whether the context was canceled or ran out of time.

If the test says the wait ended at 04:13:00, the minute loop is still there. If the 30 s deadline gives `<nil>` an hour later, the `select` has no `ctx.Done()` case, or that case returns `nil`.
