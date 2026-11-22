`Snapshot` returns `c.status` itself, so whoever holds the result can change the cache without any lock. What should it return instead?

`maps.Clone` copies a map. Take `c.mu.RLock()`, with a deferred `c.mu.RUnlock()`, around the copy, and use the same pair in `Get`.
