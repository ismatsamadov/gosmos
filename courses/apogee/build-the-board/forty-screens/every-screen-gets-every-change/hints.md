Every method touches `h.subs`, which many goroutines share, so each starts with `h.mu.Lock()` and a deferred `h.mu.Unlock()`.

In `Subscribe`, make the channel with `make(chan View, 4)` before you lock, store it as both key and value, and return it: the stream gets the receiving end, and the hub keeps the end it can send on and close. `Unsubscribe` gets the receiving end back, which is exactly the key, so a comma-ok lookup finds the channel to close.

`Publish` ranges over the map's values and sends `v` on each. `Len` returns `len(h.subs)`. Delete a channel before closing it, so no `Publish` sends on it afterwards.
