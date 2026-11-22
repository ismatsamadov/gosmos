`return <-acks, nil` has no way to stop waiting. What else could `waitAck` wait on, alongside `acks`?

`time.After(limit)` gives you a channel that delivers once `limit` has passed. Put a receive from it in a `select`, next to the receive from `acks`.
