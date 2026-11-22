Count fires in an `atomic.Int32` that `fire` adds one to (import `sync/atomic`). Call `Notify` at 0, 1 and 2 seconds, with a `time.Sleep(time.Second)` between calls, so the quiet ends at 7 seconds.

Sleep until 6.999 s, call `synctest.Wait()` and check the count is 0; sleep the last millisecond, wait again and check it is 1. Then `Notify`, `Stop`, sleep well past another quiet, wait, and check it is still 1.

Each mutant fails one check: the two that fire early fail the count before 7 s, the one that fires twice the count at 7 s, and the one that ignores `Stop` the count at the end. Without `synctest.Wait()` before a check, a fire due at that very instant may not have run yet.
