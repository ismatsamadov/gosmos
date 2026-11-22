The deferred `s.Close()` runs whenever `runPass` returns, after the early return's `pass failed` line as well as after `pass ok`, so every session closes exactly once.
