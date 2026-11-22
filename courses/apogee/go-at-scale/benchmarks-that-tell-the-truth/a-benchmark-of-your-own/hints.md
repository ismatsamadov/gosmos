Wrap `for b.Loop() { ... }` around the one line that calls `planDay`. The line with `stationDay` stays above it, where it runs once.

Don't count the turns yourself: `b.Loop` returns false on its own once the run has had the turns it asked for, 20 in the tests.

If `planDay` got 2 calls, the call is still outside the loop. If it got 21, a call is left outside the loop, or the loop counts to `b.N`, and then `testing.Benchmark` calls `benchPlanDay` twice, the first time with `b.N` at 1.
