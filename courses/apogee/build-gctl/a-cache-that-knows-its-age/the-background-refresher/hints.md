Lock at the top of `passes` and `refresh`, and defer the unlock on the next line. A `sync.Mutex` needs no setup: its zero value is unlocked, so `newShell` doesn't change.

If `passes` held the lock and called `refresh`, which locks too, it would wait for itself forever. That is why both call `reload`, which never locks: only a method already holding `c.mu` calls it.

Return `slices.Clone(c.cached)` from `passes`, still under the lock, so that a caller who sorts or changes the passes changes only a copy. Add `"slices"` and `"sync"` to the imports of `cache.go`.
