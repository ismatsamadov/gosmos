Range over the map with both the key and the channel. Inside the loop, a `select` has one case that sends `v` on the channel, and a `default`.

In the `default` case, delete the channel's key from `h.subs` first and then close the channel. Deleting while ranging over a map is allowed in Go. The stream sees the closed channel, writes `bye` and returns.

When that stream returns, its deferred `Unsubscribe` gets a channel the hub no longer holds. Check that `Unsubscribe` closes only what it finds in `h.subs`: closing a channel twice panics.
