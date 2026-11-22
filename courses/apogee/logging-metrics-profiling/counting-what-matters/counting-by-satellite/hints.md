Make the key first, from `a` and `b`. Look it up between `v.mu.RLock()` and `v.mu.RUnlock()`, and return the counter when it's there: reading a map that is still nil is safe, and finds nothing.

When it's missing, take `v.mu.Lock()`, with a deferred `Unlock`, and look again, returning what you find. Only if it's still missing, make the map if it is nil, then store a new counter under the key and return it.

If 40 goroutines get more than one counter for the same pair, two of them both found it missing and both added one: the second look, under `Lock`, is what stops the second.
