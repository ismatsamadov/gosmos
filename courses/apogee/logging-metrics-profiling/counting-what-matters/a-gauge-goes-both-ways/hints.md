Put the call that takes 1 away right after the one that adds it, as a deferred call: `defer`, then `g.Add` with -1.

With the `defer` in place, nothing is left to do after `run` but return what it returns, so the `if` and the second `g.Add(-1)` go.

If the test reads 10 after twenty passes, the failed ones still return before the gauge comes down; if it reads 1 after a panic, the call taking 1 away isn't deferred. The deferred call must come before `run` is called.
