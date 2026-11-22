`synctest.Test(t, func(t *testing.T) { ... })` holds all three checks. Make `day`, `shift` and `now` inside it as the readme does: `time.Now()` reads the bubble's clock only there.

After each wait, `now()` is the instant it returned, so compare it with `Equal`: 04:12:25 after the first wait, and still 04:12:25 after the wait for 03:12:25. For the third, cancel a context made with `context.WithCancel(t.Context())` before the call, and test the error with `errors.Is` (add `errors` to the imports).

Each mutant fails one check: the one that wakes a second early needs the exact time of the first wait, the one that waits a minute needs the time checked after the second, and the one that ignores its context needs the canceled wait. Wait for an instant still ahead there: with a canceled context and an instant already past, either answer is right.
